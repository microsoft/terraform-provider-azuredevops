package core

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/core"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceTeamAdministrators() *schema.Resource {
	return &schema.Resource{
		Create: resourceTeamAdministratorsCreate,
		Read:   resourceTeamAdministratorsRead,
		Update: resourceTeamAdministratorsUpdate,
		Delete: resourceTeamAdministratorsDelete,
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
			"administrators": {
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

func resourceTeamAdministratorsCreate(d *schema.ResourceData, m interface{}) error {
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

	if strings.EqualFold(d.Get("mode").(string), "overwrite") {
		administrators := tfhelper.ExpandStringSet(d.Get("administrators").(*schema.Set))
		err := updateTeamAdministrators(d, clients, team, &administrators)
		if err != nil {
			return err
		}
	} else {
		administratorsToAdd := d.Get("administrators").(*schema.Set)
		err := setTeamAdministratorsPermissions(d, clients, team, linq.From(administratorsToAdd.List()), securityhelper.PermissionTypeValues.Allow)
		if err != nil {
			return err
		}
	}

	// The ID for this resource is meaningless so we can just assign a random ID
	d.SetId(fmt.Sprintf("%d", rand.Int()))

	return resourceTeamAdministratorsRead(d, m)
}

func resourceTeamAdministratorsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	teamID := d.Get("team_id").(string)

	team, err := clients.CoreClient.GetTeam(clients.Ctx, core.GetTeamArgs{
		ProjectId:      converter.String(projectID),
		TeamId:         converter.String(teamID),
		ExpandIdentity: converter.Bool(false),
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	administratorList, err := readTeamAdministrators(d, clients, team)
	if err != nil {
		return err
	}

	mode := d.Get("mode").(string)
	stateMembers := d.Get("administrators").(*schema.Set)
	administrators := make([]string, 0)
	for _, administratorship := range administratorList.List() {
		if strings.EqualFold("overwrite", mode) || stateMembers.Contains(administratorship) {
			administrators = append(administrators, administratorship.(string))
		}
	}

	d.Set("project_id", team.ProjectId.String())
	d.Set("team_id", team.Id.String())
	d.Set("administrators", administrators)

	return nil
}

func resourceTeamAdministratorsUpdate(d *schema.ResourceData, m interface{}) error {
	if !d.HasChange("administrators") && !d.HasChange("mode") {
		return nil
	}

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

	if strings.EqualFold(d.Get("mode").(string), "overwrite") {
		administrators := tfhelper.ExpandStringSet(d.Get("administrators").(*schema.Set))
		err = updateTeamAdministrators(d, clients, team, &administrators)
		if err != nil {
			return err
		}
	} else {
		oldData, newData := d.GetChange("administrators")

		// administrators that need to be added will be missing from the old data, but present in the new data
		administratorsToAdd := newData.(*schema.Set).Difference(oldData.(*schema.Set))
		err = setTeamAdministratorsPermissions(d, clients, team, linq.From(administratorsToAdd.List()), securityhelper.PermissionTypeValues.Allow)
		if err != nil {
			return err
		}

		// administrators that need to be removed will be missing from the new data, but present in the old data
		administratorsToRemove := oldData.(*schema.Set).Difference(newData.(*schema.Set))
		err = setTeamAdministratorsPermissions(d, clients, team, linq.From(administratorsToRemove.List()), securityhelper.PermissionTypeValues.NotSet)
		if err != nil {
			return err
		}
	}
	return resourceTeamAdministratorsRead(d, m)
}

func resourceTeamAdministratorsDelete(d *schema.ResourceData, m interface{}) error {
	var err error
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

	var administratorList *schema.Set
	if strings.EqualFold("overwrite", d.Get("mode").(string)) {
		log.Printf("[TRACE] Removing all administrators from team %s", *team.Name)

		administratorList, err = readTeamAdministrators(d, clients, team)
		if err != nil {
			return err
		}
	} else {
		administratorList = d.Get("administrators").(*schema.Set)
	}

	err = setTeamAdministratorsPermissions(d, clients, team, linq.From(administratorList.List()), securityhelper.PermissionTypeValues.NotSet)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
