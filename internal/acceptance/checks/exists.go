package checks

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance/testclient"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance/types"
)

func DoesNotExistInAzure(t *testing.T, testResource types.TestResource, resourceAddr string) resource.TestCheckFunc {
	return existsFunc(t, testResource, resourceAddr, false)
}

func ExistsInAzure(t *testing.T, testResource types.TestResource, resourceAddr string) resource.TestCheckFunc {
	return existsFunc(t, testResource, resourceAddr, true)
}

func existsFunc(t *testing.T, testResource types.TestResource, resourceAddr string, shouldExist bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testclient.New(t)
		rs, ok := s.RootModule().Resources[resourceAddr]
		if !ok {
			return fmt.Errorf("%q was not found in the state", resourceAddr)
		}

		ok, err := testResource.Exists(t.Context(), client, rs.Primary)
		if err != nil {
			return fmt.Errorf("checking existence for %q: %+v", resourceAddr, err)
		}

		if ok != shouldExist {
			if !shouldExist {
				return fmt.Errorf("%q still exists", resourceAddr)
			}

			return fmt.Errorf("%q did not exist", resourceAddr)
		}

		return nil
	}
}
