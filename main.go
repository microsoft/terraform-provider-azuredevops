package main

import (
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return azuredevops.Provider()
		},
	})
}
