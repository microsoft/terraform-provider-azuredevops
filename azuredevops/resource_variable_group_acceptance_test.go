// +build all resource_variable_group

package azuredevops

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/testhelper"
)

func TestAccVariableGroup_CreateAndUpdate(t *testing.T) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	vargroupNameFirst := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	vargroupNameSecond := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	vargroupNameNoSecret := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	allowAccessTrue := true
	allowAccessFalse := false

	tfVarGroupNode := "azuredevops_variable_group.vg"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testhelper.TestAccPreCheck(t, nil) },
		Providers:    testAccProviders,
		CheckDestroy: testAccVariableGroupCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testhelper.TestAccVariableGroupResource(projectName, vargroupNameFirst, allowAccessTrue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfVarGroupNode, "project_id"),
					resource.TestCheckResourceAttr(tfVarGroupNode, "name", vargroupNameFirst),
					testAccCheckVariableGroupResourceExists(vargroupNameFirst, allowAccessTrue),
				),
			}, {
				Config: testhelper.TestAccVariableGroupResource(projectName, vargroupNameSecond, allowAccessFalse),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfVarGroupNode, "project_id"),
					resource.TestCheckResourceAttr(tfVarGroupNode, "name", vargroupNameSecond),
					testAccCheckVariableGroupResourceExists(vargroupNameSecond, allowAccessFalse),
				),
			}, {
				Config: testhelper.TestAccVariableGroupResourceNoSecrets(projectName, vargroupNameNoSecret, allowAccessFalse),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfVarGroupNode, "project_id"),
					resource.TestCheckResourceAttr(tfVarGroupNode, "name", vargroupNameNoSecret),
					testAccCheckVariableGroupResourceExists(vargroupNameNoSecret, allowAccessFalse),
				),
			}, {
				// Resource Acceptance Testing https://www.terraform.io/docs/extend/resources/import.html#resource-acceptance-testing-implementation
				ResourceName:      tfVarGroupNode,
				ImportStateIdFunc: testAccImportStateIDFunc(tfVarGroupNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVariableGroupKeyVault_CreateAndUpdate(t *testing.T) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	vargroupKeyvault := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	keyVaultName := "key-vault-name"
	allowAccessFalse := false
	tfVarGroupNode := "azuredevops_variable_group.vg"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testhelper.TestAccPreCheck(t, nil) },
		Providers:    testAccProviders,
		CheckDestroy: testAccVariableGroupCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testhelper.TestAccVariableGroupResourceKeyVaultWithProject(projectName, vargroupKeyvault, allowAccessFalse, keyVaultName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfVarGroupNode, "project_id"),
					resource.TestCheckResourceAttr(tfVarGroupNode, "name", vargroupKeyvault),
					testAccCheckVariableGroupResourceExists(vargroupKeyvault, allowAccessFalse),
				),
			}, {
				// Resource Acceptance Testing https://www.terraform.io/docs/extend/resources/import.html#resource-acceptance-testing-implementation
				ResourceName:      tfVarGroupNode,
				ImportStateIdFunc: testAccImportStateIDFunc(tfVarGroupNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Given an AzDO variable group name, this will return a function that will check whether
// or not the definition (1) exists in the state, (2) exists in AzDO, and (3) has the correct
// or expected name
func testAccCheckVariableGroupResourceExists(expectedName string, expectedAllowAccess bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		varGroup, ok := s.RootModule().Resources["azuredevops_variable_group.vg"]
		if !ok {
			return fmt.Errorf("Did not find a variable group in the TF state")
		}

		variableGroup, err := getVariableGroupFromResource(varGroup)
		if err != nil {
			return err
		}

		if *variableGroup.Name != expectedName {
			return fmt.Errorf("Variable Group has Name=%s, but expected %s", *variableGroup.Name, expectedName)
		}

		// testing Allow access with definition reference AzDo object
		definitionReference, err := getDefinitionResourceFromVariableGroupResource(varGroup)
		if err != nil {
			return err
		}

		if expectedAllowAccess {
			if len(*definitionReference) == 0 {
				return fmt.Errorf("Definition reference should be not empty for allow access true")
			}
			if len(*definitionReference) > 0 && *(*definitionReference)[0].Authorized != expectedAllowAccess {
				return fmt.Errorf("Variable Group has Allow_access=%t, but expected %t", *(*definitionReference)[0].Authorized, expectedAllowAccess)
			}
		} else {
			if len(*definitionReference) > 0 {
				return fmt.Errorf("Definition reference should be empty for allow access false")
			}
		}
		return nil
	}
}

// Verifies that all variable groups referenced in the state are destroyed. This will be
// invoked *after* Terraform destroys the resource but *before* the state is wiped clean.
func testAccVariableGroupCheckDestroy(s *terraform.State) error {
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_variable_group" {
			continue
		}

		// Indicates the variable group still exists -- this should fail the test
		if _, err := getVariableGroupFromResource(resource); err == nil {
			return fmt.Errorf("Unexpectedly found a variable group that should be deleted")
		}

		// Indicates the definition reference still exists -- this should fail the test
		if _, err := getDefinitionResourceFromVariableGroupResource(resource); err == nil {
			return fmt.Errorf("Unexpectedly found a definition reference for allow access that should be deleted")
		}
	}

	return nil
}

// Given a resource from the state, return a variable group (and error)
func getVariableGroupFromResource(resource *terraform.ResourceState) (*taskagent.VariableGroup, error) {
	variableGroupID, err := strconv.Atoi(resource.Primary.ID)
	if err != nil {
		return nil, err
	}

	projectID := resource.Primary.Attributes["project_id"]
	clients := testAccProvider.Meta().(*config.AggregatedClient)
	return clients.TaskAgentClient.GetVariableGroup(
		clients.Ctx,
		taskagent.GetVariableGroupArgs{
			GroupId: &variableGroupID,
			Project: &projectID,
		},
	)
}

// Given a resource from the state, return a definition Reference (and error)
func getDefinitionResourceFromVariableGroupResource(resource *terraform.ResourceState) (*[]build.DefinitionResourceReference, error) {
	projectID := resource.Primary.Attributes["project_id"]
	clients := testAccProvider.Meta().(*config.AggregatedClient)

	return clients.BuildClient.GetProjectResources(
		clients.Ctx,
		build.GetProjectResourcesArgs{
			Project: &projectID,
			Type:    converter.String("variablegroup"),
			Id:      &resource.Primary.ID,
		},
	)
}

// Set the Imported ID
func testAccImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["id"]), nil
	}
}

func init() {
	InitProvider()
}
