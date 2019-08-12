package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider - The top level Azure DevOps Provider definition.
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"azuredevops_foo": resourceFoo(),
		},
	}
}
