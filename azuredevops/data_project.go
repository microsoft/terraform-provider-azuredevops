package azuredevops

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataProject() *schema.Resource {
	baseSchema := resourceProject()
	return &schema.Resource{
		Read:   baseSchema.Read,
		Schema: baseSchema.Schema,
	}
}
