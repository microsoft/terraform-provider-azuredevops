package core

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/identity"
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

	projectID := d.Get("project_id").(string)
	teamID := d.Get("team_id").(string)

	team, err := clients.CoreClient.GetTeam(clients.Ctx, core.GetTeamArgs{
		ProjectId:      converter.String(projectID),
		TeamId:         converter.String(teamID),
		ExpandIdentity: converter.Bool(true),
	})

	if err != nil {
		return err
	}

	if strings.EqualFold(d.Get("mode").(string), "overwrite") {
		members := tfhelper.ExpandStringSet(d.Get("members").(*schema.Set))
		updateTeamMembers(clients, team, &members)
	} else {
		membersToAdd := d.Get("members").(*schema.Set)
		err = addTeamMembers(clients, team, linq.From(membersToAdd.List()))
		if err != nil {
			return err
		}
	}

	// The ID for this resource is meaningless so we can just assign a random ID
	d.SetId(fmt.Sprintf("%d", rand.Int()))

	return resourceTeamMembersRead(d, m)
}

func resourceTeamMembersRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	teamID := d.Get("team_id").(string)

	team, err := clients.CoreClient.GetTeam(clients.Ctx, core.GetTeamArgs{
		ProjectId:      converter.String(projectID),
		TeamId:         converter.String(teamID),
		ExpandIdentity: converter.Bool(true),
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	membershipList, err := readTeamMembers(clients, team)
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

	projectID := d.Get("project_id").(string)
	teamID := d.Get("team_id").(string)

	team, err := clients.CoreClient.GetTeam(clients.Ctx, core.GetTeamArgs{
		ProjectId:      converter.String(projectID),
		TeamId:         converter.String(teamID),
		ExpandIdentity: converter.Bool(true),
	})

	if err != nil {
		return err
	}

	if strings.EqualFold(d.Get("mode").(string), "overwrite") {
		members := tfhelper.ExpandStringSet(d.Get("members").(*schema.Set))
		updateTeamMembers(clients, team, &members)
	} else {
		oldData, newData := d.GetChange("members")

		// members that need to be added will be missing from the old data, but present in the new data
		membersToAdd := newData.(*schema.Set).Difference(oldData.(*schema.Set))
		err = addTeamMembers(clients, team, linq.From(membersToAdd.List()))
		if err != nil {
			return err
		}

		// members that need to be removed will be missing from the new data, but present in the old data
		membersToRemove := oldData.(*schema.Set).Difference(newData.(*schema.Set))
		err = removeTeamMembers(clients, team, linq.From(membersToRemove.List()))
		if err != nil {
			return err
		}
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

	if strings.EqualFold("overwrite", d.Get("mode").(string)) {
		log.Printf("[TRACE] Removing all members from team %s", *team.Name)

		members, err := clients.IdentityClient.ReadMembers(clients.Ctx, identity.ReadMembersArgs{
			ContainerId: converter.String(team.Id.String()),
		})
		if err != nil {
			return err
		}
		for _, id := range *members {
			_, err := clients.IdentityClient.RemoveMember(clients.Ctx, identity.RemoveMemberArgs{
				ContainerId: converter.String(team.Id.String()),
				MemberId:    converter.String(id),
			})
			if err != nil {
				return fmt.Errorf("Error removing member %s from team %s: %+v", id, *team.Name, err)
			}
		}
	} else {
		members := tfhelper.ExpandStringSet(d.Get("members").(*schema.Set))
		err := removeTeamMembers(clients, team, linq.From(members))
		if err != nil {
			return err
		}
	}

	d.SetId("")
	return nil
}
