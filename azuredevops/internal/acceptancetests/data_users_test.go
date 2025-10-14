package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccUsers_DataSource(t *testing.T) {
	userName := "foo@email.com"
	tfNode := "data.azuredevops_users.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclDataUsersBasic(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "users.#", "1"),
					resource.TestCheckResourceAttr(tfNode, "principal_name", "foo@email.com"),
				),
			},
		},
	})
}

func TestAccUsers_DataSource_AllSvc(t *testing.T) {
	tfNode := "data.azuredevops_users.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclDataUserAllSvc(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "users.0.id"),
					resource.TestCheckResourceAttrSet(tfNode, "users.0.origin_id"),
				),
			},
		},
	})
}

func TestAccUsers_DataSource_All_WithFeatures(t *testing.T) {
	tfNode := "data.azuredevops_users.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclDataUserAllWithFeatures(3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "users.0.id"),
					resource.TestCheckResourceAttrSet(tfNode, "users.0.origin_id"),
				),
			},
		},
	})
}

func TestAccUsers_DataSource_userNotFound(t *testing.T) {
	tfNode := "data.azuredevops_users.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclDataUserUserNotFound(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "users.#", "0"),
				),
			},
		},
	})
}

func hclDataUserAllWithFeatures(numWorkers int) string {
	return fmt.Sprintf(`
data "azuredevops_users" "test" {
  features {
    concurrent_workers = %d
  }
}`, numWorkers)
}

func hclDataUserAllSvc() string {
	return `
data "azuredevops_users" "test" {
  subject_types = ["aad"]
}`
}

func hclDataUsersBasic(uname string) string {
	return fmt.Sprintf(`
resource "azuredevops_user_entitlement" "test" {
  principal_name       = "%[1]s"
  account_license_type = "basic"
}

data "azuredevops_users" "test" {
  principal_name = "%[1]s"
  depends_on     = [azuredevops_user_entitlement.test]
}`, uname)
}

func hclDataUserUserNotFound() string {
	return `
data "azuredevops_users" "test" {
  principal_name = "dummy"
}`
}
