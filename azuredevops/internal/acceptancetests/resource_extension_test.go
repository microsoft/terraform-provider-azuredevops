package acceptancetests

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/extensionmanagement"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

func TestAccExtension_basic(t *testing.T) {
	publisherId := "ms-securitydevops"
	extensionId := "microsoft-security-devops-azdevops"
	tfNode := "azuredevops_extension.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkExtensionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclExtensionBasic(publisherId, extensionId),
				Check: resource.ComposeTestCheckFunc(
					checkExtensionExist(extensionId),
					resource.TestCheckResourceAttrSet(tfNode, "extension_id"),
					resource.TestCheckResourceAttrSet(tfNode, "publisher_id"),
					resource.TestCheckResourceAttrSet(tfNode, "publisher_name"),
					resource.TestCheckResourceAttrSet(tfNode, "extension_name"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s/%s", publisherId, extensionId),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccExtension_complete(t *testing.T) {
	publisherId := "ms-securitydevops"
	extensionId := "microsoft-security-devops-azdevops"
	tfNode := "azuredevops_extension.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkExtensionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclExtensionComplete(publisherId, extensionId),
				Check: resource.ComposeTestCheckFunc(
					checkExtensionExist(extensionId),
					resource.TestCheckResourceAttrSet(tfNode, "extension_id"),
					resource.TestCheckResourceAttrSet(tfNode, "publisher_id"),
					resource.TestCheckResourceAttrSet(tfNode, "publisher_name"),
					resource.TestCheckResourceAttrSet(tfNode, "extension_name"),
					resource.TestCheckResourceAttrSet(tfNode, "scope.#"),
					resource.TestCheckResourceAttrSet(tfNode, "version"),
					resource.TestCheckResourceAttrSet(tfNode, "disabled"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s/%s", publisherId, extensionId),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccExtension_update(t *testing.T) {
	publisherId := "ms-securitydevops"
	extensionId := "microsoft-security-devops-azdevops"
	tfNode := "azuredevops_extension.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkExtensionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclExtensionBasic(publisherId, extensionId),
				Check: resource.ComposeTestCheckFunc(
					checkExtensionExist(extensionId),
					resource.TestCheckResourceAttrSet(tfNode, "extension_id"),
					resource.TestCheckResourceAttrSet(tfNode, "publisher_id"),
					resource.TestCheckResourceAttrSet(tfNode, "publisher_name"),
					resource.TestCheckResourceAttrSet(tfNode, "extension_name"),
					resource.TestCheckResourceAttrSet(tfNode, "scope.#"),
					resource.TestCheckResourceAttrSet(tfNode, "version"),
					resource.TestCheckResourceAttrSet(tfNode, "disabled"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s/%s", publisherId, extensionId),
				ImportStateVerify: true,
			},
			{
				Config: hclExtensionUpdate(publisherId, extensionId, true),
				Check: resource.ComposeTestCheckFunc(
					checkExtensionExist(extensionId),
					resource.TestCheckResourceAttrSet(tfNode, "extension_id"),
					resource.TestCheckResourceAttrSet(tfNode, "publisher_id"),
					resource.TestCheckResourceAttrSet(tfNode, "publisher_name"),
					resource.TestCheckResourceAttrSet(tfNode, "extension_name"),
					resource.TestCheckResourceAttr(tfNode, "disabled", "true"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s/%s", publisherId, extensionId),
				ImportStateVerify: true,
			},
			{
				Config: hclExtensionUpdate(publisherId, extensionId, false),
				Check: resource.ComposeTestCheckFunc(
					checkExtensionExist(extensionId),
					resource.TestCheckResourceAttrSet(tfNode, "extension_id"),
					resource.TestCheckResourceAttrSet(tfNode, "publisher_id"),
					resource.TestCheckResourceAttrSet(tfNode, "publisher_name"),
					resource.TestCheckResourceAttrSet(tfNode, "extension_name"),
					resource.TestCheckResourceAttrSet(tfNode, "scope.#"),
					resource.TestCheckResourceAttr(tfNode, "disabled", "false"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s/%s", publisherId, extensionId),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccExtension_requireImportError(t *testing.T) {
	publisherId := "ms-securitydevops"
	extensionId := "microsoft-security-devops-azdevops"
	tfNode := "azuredevops_extension.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkExtensionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclExtensionBasic(publisherId, extensionId),
				Check: resource.ComposeTestCheckFunc(
					checkExtensionExist(extensionId),
					resource.TestCheckResourceAttrSet(tfNode, "extension_id"),
					resource.TestCheckResourceAttrSet(tfNode, "publisher_id"),
					resource.TestCheckResourceAttrSet(tfNode, "publisher_name"),
					resource.TestCheckResourceAttrSet(tfNode, "extension_name"),
					resource.TestCheckResourceAttrSet(tfNode, "scope.#"),
					resource.TestCheckResourceAttrSet(tfNode, "version"),
					resource.TestCheckResourceAttrSet(tfNode, "disabled"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s/%s", publisherId, extensionId),
				ImportStateVerify: true,
			},
			{
				Config:      hclExtensionImportError(publisherId, extensionId),
				ExpectError: requiresExtensionImportError(publisherId, extensionId),
			},
		},
	})
}

func checkExtensionDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
	for _, res := range s.RootModule().Resources {
		if res.Type != "azuredevops_extension" {
			continue
		}
		ids := strings.Split(res.Primary.ID, "/")

		_, err := clients.ExtensionManagementClient.GetInstalledExtensionByName(clients.Ctx, extensionmanagement.GetInstalledExtensionByNameArgs{
			PublisherName: &ids[0],
			ExtensionName: &ids[1],
		})

		if err == nil {
			return fmt.Errorf("Extension with Publisher ID=%s , Extension ID: %s should not exist", ids[0], ids[1])
		}
	}
	return nil
}

func checkExtensionExist(expectedExtensionId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources["azuredevops_extension.test"]
		if !ok {
			return fmt.Errorf("Did not find `azuredevops_extension` in the TF state")
		}

		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
		ids := strings.Split(res.Primary.ID, "/")

		extension, err := clients.ExtensionManagementClient.GetInstalledExtensionByName(clients.Ctx, extensionmanagement.GetInstalledExtensionByNameArgs{
			PublisherName: &ids[0],
			ExtensionName: &ids[1],
		})
		if err != nil {
			return fmt.Errorf("Extension with Publisher ID=%s , Extension ID: %s cannot be found!. Error=%v", ids[0], ids[1], err)
		}

		if *extension.ExtensionId != expectedExtensionId {
			return fmt.Errorf("Extension with Publisher ID=%s has Extension ID=%s, but expected Extension ID=%s", *extension.PublisherId, *extension.ExtensionId, expectedExtensionId)
		}
		return nil
	}
}

func requiresExtensionImportError(publisherId, extensionId string) *regexp.Regexp {
	message := "Installing extension for Publisher: %s, Name: %s. Error: TF1590010: Extension %s.%s is already installed in this organization"
	return regexp.MustCompile(fmt.Sprintf(message, publisherId, extensionId, publisherId, extensionId))
}

func hclExtensionBasic(publisherId, extensionId string) string {
	return fmt.Sprintf(`
resource "azuredevops_extension" "test" {
  publisher_id = "%s"
  extension_id = "%s"
}`, publisherId, extensionId)
}

func hclExtensionComplete(publisherId, extensionId string) string {
	return fmt.Sprintf(`
resource "azuredevops_extension" "test" {
  publisher_id = "%s"
  extension_id = "%s"
  disabled     = false
}`, publisherId, extensionId)
}

func hclExtensionUpdate(publisherId, extensionId string, disabled bool) string {
	return fmt.Sprintf(`
resource "azuredevops_extension" "test" {
  publisher_id = "%s"
  extension_id = "%s"
  disabled     = %t
}`, publisherId, extensionId, disabled)
}

func hclExtensionImportError(publisherId, extensionId string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_extension" "import" {
  publisher_id = azuredevops_extension.test.publisher_id
  extension_id = azuredevops_extension.test.extension_id
}`, hclExtensionBasic(publisherId, extensionId))
}
