package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider - The top level Azure DevOps Provider definition.
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{

			// Why is the key/value named the way they are?
			"azuredevops_foo": resourceFoo(),

		},
	}
}
