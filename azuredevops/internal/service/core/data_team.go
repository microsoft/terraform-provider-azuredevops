package core

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataTeam() *schema.Resource {
	baseSchema := ResourceTeam()
	for k, v := range baseSchema.Schema {
		if k != "name" && k != "project_id" {
			baseSchema.Schema[k] = &schema.Schema{
				Type:     v.Type,
				Computed: true,
			}
		}
	}
	return &schema.Resource{
		Read:   dataTeamRead,
		Schema: baseSchema.Schema,
	}
}

func dataTeamRead(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("Not implemented")
}
