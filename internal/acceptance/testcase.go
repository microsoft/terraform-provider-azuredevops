package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance/testclient"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance/types"
	"github.com/microsoft/terraform-provider-azuredevops/internal/provider"
)

func (d TestData) DataSourceTest(t *testing.T, steps []resource.TestStep) {
	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		Steps:    steps,
	}
	d.runAcceptanceTest(t, testCase)
}

func (d TestData) ResourceTest(t *testing.T, testResource types.TestResource, steps []resource.TestStep) {
	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		CheckDestroy: func(s *terraform.State) error {
			client := testclient.New(t)
			for label, resourceState := range s.RootModule().Resources {
				if resourceState.Type != d.ResourceType {
					continue
				}
				if label != d.ResourceLabel {
					continue
				}

				ok, err := testResource.Exists(t.Context(), client, resourceState.Primary)
				if err != nil {
					return fmt.Errorf("checking existence: %v", err)
				}
				if ok {
					return fmt.Errorf("%q still exists", d.ResourceAddr())
				}
			}

			return nil
		},
		Steps: steps,
	}
	d.runAcceptanceTest(t, testCase)
}

func (d TestData) runAcceptanceTest(t *testing.T, testCase resource.TestCase) {
	testCase.ExternalProviders = d.externalProviders()
	testCase.ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"azuredevops": providerserver.NewProtocol6WithError(provider.New()),
	}

	resource.ParallelTest(t, testCase)
}

func (d TestData) externalProviders() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"azuread": {
			VersionConstraint: "=3.7.0",
			Source:            "registry.terraform.io/hashicorp/azuread",
		},
		"time": {
			VersionConstraint: "=0.13.1",
			Source:            "registry.terraform.io/hashicorp/time",
		},
		"azuredevops-v1": {
			VersionConstraint: "=1.13.0",
			Source:            "registry.terraform.io/microsoft/azuredevops",
		},
	}
}
