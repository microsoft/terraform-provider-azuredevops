package core

import (
	"fmt"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func DataTeam() *schema.Resource {
	return &schema.Resource{
		Read: dataTeamRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"administrators": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Set:        schema.HashString,
			},
			"members": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Set:        schema.HashString,
			},
			"descriptor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"top": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      100,
				ValidateFunc: validation.IntAtLeast(1),
			},
		},
	}
}

func dataTeamRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	teamName := d.Get("name").(string)

	top := d.Get("top").(int)

	team, members, administrators, err := readTeamByName(d, clients, projectID, teamName, top)
	if err != nil {
		return err
	}

	descriptor, err := clients.GraphClient.GetDescriptor(clients.Ctx, graph.GetDescriptorArgs{
		StorageKey: team.Id,
	})
	if err != nil {
		return fmt.Errorf(" get team descriptor. Error: %+v", err)
	}

	d.SetId(team.Id.String())
	d.Set("name", team.Name)
	d.Set("description", team.Description)
	d.Set("administrators", administrators)
	d.Set("members", members)
	d.Set("descriptor", descriptor.Value)

	return nil
}

func readTeamByName(d *schema.ResourceData, clients *client.AggregatedClient, projectID string, teamName string, top int) (*core.WebApiTeam, *schema.Set, *schema.Set, error) {
	teamList, err := clients.CoreClient.GetTeams(clients.Ctx, core.GetTeamsArgs{
		ProjectId:      converter.String(projectID),
		Mine:           converter.Bool(false),
		Top:            converter.Int(top),
		ExpandIdentity: converter.Bool(false),
	})

	if err != nil {
		return nil, nil, nil, err
	}

	if teamList == nil || len(*teamList) <= 0 {
		return nil, nil, nil, fmt.Errorf("Project [%s] does not contain any teams", projectID)
	}

	iTeam := linq.From(*teamList).
		FirstWith(func(v interface{}) bool {
			item := v.(core.WebApiTeam)
			return strings.EqualFold(*item.Name, teamName)
		})
	if iTeam == nil {
		return nil, nil, nil, fmt.Errorf("Unable to find Team with name [%s] in project with ID [%s]", teamName, projectID)
	}

	team := iTeam.(core.WebApiTeam)
	members, err := readTeamMembers(clients, &team)
	if err != nil {
		return nil, nil, nil, err
	}
	administrators, err := readTeamAdministrators(d, clients, &team)
	if err != nil {
		return nil, nil, nil, err
	}

	return &team, members, administrators, nil
}
