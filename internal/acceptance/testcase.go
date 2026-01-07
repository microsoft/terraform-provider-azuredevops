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

func (d TestData) ResourceTest(t *testing.T, testResource types.TestResource, steps []resource.TestStep) {
	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"azuredevops": providerserver.NewProtocol6WithError(provider.New()),
		},
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
	resource.ParallelTest(t, testCase)
}
