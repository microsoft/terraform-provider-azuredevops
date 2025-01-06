package core

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func DataTeam() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataTeamRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
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
				Deprecated:   "This property is deprecated and will be removed in the feature", // TODO remove
			},
		},
	}
}

func dataTeamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	teamName := d.Get("name").(string)

	team, err := clients.CoreClient.GetTeam(clients.Ctx, core.GetTeamArgs{
		ProjectId: converter.String(projectID),
		TeamId:    converter.String(teamName),
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf(" Get Team (Team Name: %s). Error: %+v", teamName, err))
	}

	members, err := getTeamMembers(clients, team)
	if err != nil {
		return diag.FromErr(fmt.Errorf(" Get Team members (Team Name: %s). Error: %+v", teamName, err))
	}

	administrators, err := getTeamAdministrators(d, clients, team)
	if err != nil {
		return diag.FromErr(fmt.Errorf(" Get Team administrators (Team Name: %s). Error: %+v", teamName, err))
	}

	descriptor, err := clients.GraphClient.GetDescriptor(clients.Ctx, graph.GetDescriptorArgs{
		StorageKey: team.Id,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf(" Get Team descriptor (Team Name: %s). Error: %+v", teamName, err))
	}

	d.SetId(team.Id.String())
	d.Set("name", team.Name)
	d.Set("description", team.Description)
	d.Set("administrators", administrators)
	d.Set("members", members)
	d.Set("descriptor", descriptor.Value)
	return nil
}
