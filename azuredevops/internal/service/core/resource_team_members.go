package core

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/ahmetb/go-linq"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceTeamMembers() *schema.Resource {
	return &schema.Resource{
		Create: resourceTeamMembersCreate,
		Read:   resourceTeamMembersRead,
		Update: resourceTeamMembersUpdate,
		Delete: resourceTeamMembersDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"team_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "add",
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					"add", "overwrite",
				}, true),
			},
			"members": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
				MinItems:   1,
				Required:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Set:        schema.HashString,
			},
		},
	}
}

func resourceTeamMembersCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	team, err := clients.CoreClient.GetTeam(clients.Ctx, core.GetTeamArgs{
		ProjectId:      converter.String(d.Get("project_id").(string)),
		TeamId:         converter.String(d.Get("team_id").(string)),
		ExpandIdentity: converter.Bool(true),
	})

	if err != nil {
		return err
	}

	var membersToAdd *schema.Set = nil
	mode := d.Get("mode").(string)
	if strings.EqualFold(mode, "overwrite") {
		membersToAdd = d.Get("members").(*schema.Set)
		members := tfhelper.ExpandStringSet(membersToAdd)
		err = setTeamMembers(clients, team, &members)
		if err != nil {
			return err
		}
	} else {
		membersToAdd = d.Get("members").(*schema.Set)
		err = addTeamMembers(clients, team, linq.From(membersToAdd.List()), true)
		if err != nil {
			return err
		}
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"Waiting"},
		Target:  []string{"Synched"},
		Refresh: func() (interface{}, string, error) {
			clients = m.(*client.AggregatedClient)
			state := "Waiting"
			actualMemberships, err := getTeamMembers(clients, team)
			if err != nil {
				return nil, "", fmt.Errorf("reading team memberships: %+v", err)
			}
			if membersToAdd == nil || actualMemberships.Intersection(membersToAdd).Len() == membersToAdd.Len() {
				state = "Synched"
			}
			if strings.EqualFold(mode, "overwrite") && membersToAdd != nil && actualMemberships.Len() != membersToAdd.Len() {
				state = "Waiting"
			}
			return state, state, nil
		},
		Timeout:                   60 * time.Minute,
		MinTimeout:                5 * time.Second,
		Delay:                     5 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	if _, err := stateConf.WaitForState(); err != nil { //nolint:staticcheck
		return fmt.Errorf("waiting for distribution of adding members. %v ", err)
	}

	// The ID for this resource is meaningless so we can just assign a random ID
	d.SetId(fmt.Sprintf("%d", rand.Int()))
	return resourceTeamMembersRead(d, m)
}

func resourceTeamMembersRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	team, err := clients.CoreClient.GetTeam(clients.Ctx, core.GetTeamArgs{
		ProjectId:      converter.String(d.Get("project_id").(string)),
		TeamId:         converter.String(d.Get("team_id").(string)),
		ExpandIdentity: converter.Bool(true),
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	membershipList, err := getTeamMembers(clients, team)
	if err != nil {
		return err
	}

	mode := d.Get("mode").(string)
	stateMembers := d.Get("members").(*schema.Set)
	members := make([]string, 0)
	for _, membership := range membershipList.List() {
		if strings.EqualFold("overwrite", mode) || stateMembers.Contains(membership) {
			members = append(members, membership.(string))
		}
	}

	d.Set("project_id", team.ProjectId.String())
	d.Set("team_id", team.Id.String())
	d.Set("members", members)
	return nil
}

func resourceTeamMembersUpdate(d *schema.ResourceData, m interface{}) error {
	if !d.HasChange("members") && !d.HasChange("mode") {
		return nil
	}

	clients := m.(*client.AggregatedClient)

	team, err := clients.CoreClient.GetTeam(clients.Ctx, core.GetTeamArgs{
		ProjectId:      converter.String(d.Get("project_id").(string)),
		TeamId:         converter.String(d.Get("team_id").(string)),
		ExpandIdentity: converter.Bool(true),
	})

	if err != nil {
		return err
	}

	var membersToAdd *schema.Set = nil
	var membersToRemove *schema.Set = nil

	mode := d.Get("mode").(string)
	if strings.EqualFold(mode, "overwrite") {
		membersToAdd := d.Get("members").(*schema.Set)
		members := tfhelper.ExpandStringSet(membersToAdd)
		err = setTeamMembers(clients, team, &members)
		if err != nil {
			return err
		}
	} else {
		oldData, newData := d.GetChange("members")

		// members that need to be added will be missing from the old data, but present in the new data
		membersToAdd = newData.(*schema.Set).Difference(oldData.(*schema.Set))
		err = addTeamMembers(clients, team, linq.From(membersToAdd.List()), true)
		if err != nil {
			return err
		}

		// members that need to be removed will be missing from the new data, but present in the old data
		membersToRemove = oldData.(*schema.Set).Difference(newData.(*schema.Set))
		err = removeTeamMembers(clients, team, linq.From(membersToRemove.List()))
		if err != nil {
			return err
		}
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"Waiting"},
		Target:  []string{"Synched"},
		Refresh: func() (interface{}, string, error) {
			clients = m.(*client.AggregatedClient)
			state := "Waiting"
			actualMemberships, err := getTeamMembers(clients, team)
			if err != nil {
				return nil, "", fmt.Errorf("Error reading team memberships: %+v", err)
			}
			if (membersToAdd == nil || actualMemberships.Intersection(membersToAdd).Len() == membersToAdd.Len()) &&
				(membersToRemove == nil || actualMemberships.Intersection(membersToRemove).Len() <= 0) {
				state = "Synched"
			}
			if strings.EqualFold(mode, "overwrite") && membersToAdd != nil && actualMemberships.Len() != membersToAdd.Len() {
				state = "Waiting"
			}

			return state, state, nil
		},
		Timeout:                   60 * time.Minute,
		MinTimeout:                5 * time.Second,
		Delay:                     5 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	if _, err = stateConf.WaitForState(); err != nil { //nolint:staticcheck
		return fmt.Errorf("waiting for distribution of member list update. %v ", err)
	}

	return resourceTeamMembersRead(d, m)
}

func resourceTeamMembersDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	teamID := d.Get("team_id").(string)

	team, err := clients.CoreClient.GetTeam(clients.Ctx, core.GetTeamArgs{
		ProjectId:      converter.String(projectID),
		TeamId:         converter.String(teamID),
		ExpandIdentity: converter.Bool(false),
	})

	if err != nil {
		return err
	}

	var membersToRemove *schema.Set = nil

	if strings.EqualFold("overwrite", d.Get("mode").(string)) {
		log.Printf("[TRACE] Removing all members from team %s", *team.Name)

		err := setTeamMembers(clients, team, nil)
		if err != nil {
			return err
		}
	} else {
		membersToRemove = d.Get("members").(*schema.Set)
		members := tfhelper.ExpandStringSet(membersToRemove)
		err := removeTeamMembers(clients, team, linq.From(members))
		if err != nil {
			return err
		}
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"Waiting"},
		Target:  []string{"Synched"},
		Refresh: func() (interface{}, string, error) {
			clients = m.(*client.AggregatedClient)
			state := "Waiting"
			actualMemberships, err := getTeamMembers(clients, team)
			if err != nil {
				return nil, "", fmt.Errorf("Error reading team memberships: %+v", err)
			}
			if (membersToRemove == nil && actualMemberships.Len() <= 0) ||
				(membersToRemove != nil && actualMemberships.Intersection(membersToRemove).Len() <= 0) {
				state = "Synched"
			}

			return state, state, nil
		},
		Timeout:                   60 * time.Minute,
		MinTimeout:                5 * time.Second,
		Delay:                     5 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	if _, err = stateConf.WaitForState(); err != nil { //nolint:staticcheck
		return fmt.Errorf("waiting for distribution of member list update. %v ", err)
	}

	d.SetId("")
	return nil
}
