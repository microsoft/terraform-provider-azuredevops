package core

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
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

	var top int
	data, ok := d.GetOk("top")
	if ok {
		top = data.(int)
	}

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
