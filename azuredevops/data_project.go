package azuredevops

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataProject() *schema.Resource {
	baseSchema := resourceProject()
	for k, v := range baseSchema.Schema {
		if k != "project_name" {
			baseSchema.Schema[k] = &schema.Schema{
				Type:     v.Type,
				Computed: true,
			}
		}
	}
	return &schema.Resource{
		Read:   baseSchema.Read,
		Schema: baseSchema.Schema,
	}
}
