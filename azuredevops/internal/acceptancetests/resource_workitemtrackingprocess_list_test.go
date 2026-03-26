package acceptancetests

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

func TestAccWorkitemtrackingprocessList_Basic(t *testing.T) {
	listName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_list.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkListDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicList(listName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkitemtrackingprocessList_Update(t *testing.T) {
	listName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_list.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkListDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicList(listName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: updatedList(listName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: basicList(listName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkitemtrackingprocessList_Integer(t *testing.T) {
	listName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_list.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkListDestroyed,
		Steps: []resource.TestStep{
			{
				Config: integerList(listName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func basicList(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_list" "test" {
  name  = "%s"
  items = ["Red", "Green", "Blue"]
}
`, name)
}

func updatedList(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_list" "test" {
  name         = "%s"
  items        = ["Red", "Green", "Blue", "Yellow"]
  is_suggested = true
}
`, name)
}

func integerList(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_list" "test" {
  name  = "%s"
  type  = "integer"
  items = ["1", "2", "3"]
}
`, name)
}

func checkListDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
	timeout := 10 * time.Second

	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_workitemtrackingprocess_list" {
			continue
		}

		id, err := uuid.Parse(resource.Primary.ID)
		if err != nil {
			return err
		}

		err = retry.RetryContext(clients.Ctx, timeout, func() *retry.RetryError {
			_, err := clients.WorkItemTrackingProcessClient.GetList(clients.Ctx, workitemtrackingprocess.GetListArgs{
				ListId: &id,
			})
			if err == nil {
				return retry.RetryableError(fmt.Errorf("list with ID %s should not exist", id.String()))
			}
			if utils.ResponseWasNotFound(err) {
				return nil
			}

			return retry.NonRetryableError(err)
		})
		if err != nil {
			return err
		}
	}

	return nil
}
