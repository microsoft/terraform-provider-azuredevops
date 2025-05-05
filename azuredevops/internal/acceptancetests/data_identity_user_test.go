package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccIdentityUsersDataSource_general(t *testing.T) {
	userName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_identity_user.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclIdentityUsersDataSourceGeneral(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
					resource.TestCheckResourceAttrSet(tfNode, "subject_descriptor"),
				),
			},
		},
	})
}

func TestAccIdentityUsersDataSource_mailAddress(t *testing.T) {
	userName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_identity_user.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclIdentityUsersDataSourceMailAddress(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
					resource.TestCheckResourceAttrSet(tfNode, "subject_descriptor"),
				),
			},
		},
	})
}

func TestAccIdentityUsersDataSource_displayName(t *testing.T) {
	userName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_identity_user.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclIdentityUsersDataSourceDisplayName(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
					resource.TestCheckResourceAttrSet(tfNode, "subject_descriptor"),
				),
			},
		},
	})
}

func TestAccIdentityUsersDataSource_accountName(t *testing.T) {
	userName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_identity_user.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclIdentityUsersDataSourceAccountName(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
					resource.TestCheckResourceAttrSet(tfNode, "subject_descriptor"),
				),
			},
		},
	})
}

func hclIdentityUsersDataSourceGeneral(userName string) string {
	return fmt.Sprintf(`
resource "azuredevops_user_entitlement" "test" {
  principal_name       = "%s"
  account_license_type = "express"
}

data "azuredevops_identity_user" "test" {
  name = azuredevops_user_entitlement.test.principal_name
}`, fmt.Sprintf(`%s@foo.com`, userName))
}

func hclIdentityUsersDataSourceMailAddress(userName string) string {
	return fmt.Sprintf(`
resource "azuredevops_user_entitlement" "test" {
  principal_name       = "%[1]s"
  account_license_type = "express"
}

data "azuredevops_identity_user" "test" {
  name          = "%[1]s"
  search_filter = "MailAddress"
  depends_on    = [azuredevops_user_entitlement.test]
}`, fmt.Sprintf(`%s@foo.com`, userName))
}

func hclIdentityUsersDataSourceDisplayName(userName string) string {
	return fmt.Sprintf(`
resource "azuredevops_user_entitlement" "test" {
  principal_name       = "%[1]s"
  account_license_type = "express"
}

data "azuredevops_identity_user" "test" {
  name          = "%[1]s"
  search_filter = "DisplayName"
  depends_on    = [azuredevops_user_entitlement.test]
}`, fmt.Sprintf(`%s@foo.com`, userName))
}

func hclIdentityUsersDataSourceAccountName(userName string) string {
	return fmt.Sprintf(`
resource "azuredevops_user_entitlement" "test" {
  principal_name       = "%[1]s"
  account_license_type = "basic"
}

data "azuredevops_identity_user" "test" {
  name          = "%[1]s"
  search_filter = "AccountName"
  depends_on    = [azuredevops_user_entitlement.test]
}`, fmt.Sprintf(`%s@foo.com`, userName))
}
