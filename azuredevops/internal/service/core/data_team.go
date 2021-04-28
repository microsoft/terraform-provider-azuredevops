package core

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

func DataTeam() *schema.Resource {
	baseSchema := ResourceTeam()
	for k := range baseSchema.Schema {
		if k != "name" && k != "project_id" {
			baseSchema.Schema[k].Computed = true
			baseSchema.Schema[k].Required = false
		}
	}
	return &schema.Resource{
		Read:   dataTeamRead,
		Schema: baseSchema.Schema,
	}
}

func dataTeamRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	teamName := d.Get("name").(string)

	team, members, administrators, err := readTeamByName(d, clients, projectID, teamName)
	if err != nil {
		return err
	}

	flattenTeam(d, team, members, administrators)
	return nil
}
