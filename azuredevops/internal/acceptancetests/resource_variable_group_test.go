package acceptancetests

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func TestAccVariableGroup_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	vgName := testutils.GenerateResourceName()
	tfVarGroupNode := "azuredevops_variable_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkVariableGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclVariableGroupBasic(projectName, vgName),
				Check: resource.ComposeTestCheckFunc(
					checkVariableGroupExists(vgName, false),
				),
			},
			{
				ResourceName:      tfVarGroupNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfVarGroupNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVariableGroup_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	vgName := testutils.GenerateResourceName()
	vgName2 := testutils.GenerateResourceName()
	tfVarGroupNode := "azuredevops_variable_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkVariableGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclVariableGroupBasic(projectName, vgName),
				Check: resource.ComposeTestCheckFunc(
					checkVariableGroupExists(vgName, false),
				),
			},
			{
				ResourceName:      tfVarGroupNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfVarGroupNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: hclVariableGroupUpdate(projectName, vgName2),
				Check: resource.ComposeTestCheckFunc(
					checkVariableGroupExists(vgName2, true),
				),
			},
			{
				ResourceName:            tfVarGroupNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfVarGroupNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_variable.0.value", "secret_variable.1.value", "secret_variable.2.value"},
			},
			{
				Config: hclVariableGroupBasic(projectName, vgName),
				Check: resource.ComposeTestCheckFunc(
					checkVariableGroupExists(vgName, false),
				),
			},
			{
				ResourceName:      tfVarGroupNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfVarGroupNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVariableGroup_keyVault_basic(t *testing.T) {
	if os.Getenv("TEST_SERVICE_PRINCIPAL_ID") == "" || os.Getenv("TEST_SERVICE_PRINCIPAL_KEY") == "" ||
		os.Getenv("TEST_ARM_TENANT_ID") == "" || os.Getenv("TEST_ARM_SUBSCRIPTION_ID") == "" ||
		os.Getenv("TEST_ARM_SUBSCRIPTION_NAME") == "" || os.Getenv("TEST_ARM_KV_NAME") == "" {
		t.Skip("Skip test as `TEST_SERVICE_PRINCIPAL_ID` or `TEST_SERVICE_PRINCIPAL_KEY` or `TEST_ARM_TENANT_ID` or `TEST_ARM_SUBSCRIPTION_ID` or `TEST_ARM_SUBSCRIPTION_NAME` or `TEST_ARM_KV_NAME` is not set")
	}
	projectName := testutils.GenerateResourceName()

	vgKeyVault := testutils.GenerateResourceName()
	tfVarGroupNode := "azuredevops_variable_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkVariableGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclVariableGroupAzureKeyVault(projectName, vgKeyVault),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfVarGroupNode, "project_id"),
					resource.TestCheckResourceAttr(tfVarGroupNode, "name", vgKeyVault),
					checkVariableGroupExists(vgKeyVault, false),
				),
			}, {
				ResourceName:            tfVarGroupNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfVarGroupNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_vault.0.search_depth"},
			},
		},
	})
}

func checkVariableGroupExists(expectedName string, expectedAllowAccess bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		varGroup, ok := s.RootModule().Resources["azuredevops_variable_group.test"]
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
				return fmt.Errorf("reference should be not empty for allow access true")
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
func checkVariableGroupDestroyed(s *terraform.State) error {
	for _, res := range s.RootModule().Resources {
		if res.Type != "azuredevops_variable_group" {
			continue
		}

		// Indicates the variable group still exists -- this should fail the test
		if _, err := getVariableGroupFromResource(res); err == nil {
			return fmt.Errorf("Unexpectedly found a variable group that should be deleted")
		}

		// Indicates the definition reference still exists -- this should fail the test
		if _, err := getDefinitionResourceFromVariableGroupResource(res); err == nil {
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
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
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
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

	return clients.BuildClient.GetProjectResources(
		clients.Ctx,
		build.GetProjectResourcesArgs{
			Project: &projectID,
			Type:    converter.String("variablegroup"),
			Id:      &resource.Primary.ID,
		},
	)
}

func hclVariableGroupBasic(projectName, variableGroupName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}
resource "azuredevops_variable_group" "test" {
  project_id   = azuredevops_project.test.id
  name         = "%s"
  description  = "test description"
  allow_access = false
  variable {
    name  = "key1"
    value = "value1"
  }
}`, projectName, variableGroupName)
}

func hclVariableGroupUpdate(projectName, variableGroupName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_variable_group" "test" {
  project_id   = azuredevops_project.test.id
  name         = "%s"
  description  = "update description"
  allow_access = true
  variable {
    name  = "key1"
    value = "value1"
  }
  variable {
    name  = "key2"
    value = "value2"
  }
  variable {
    name = "key3"
  }

  secret_variable {
    name  = "skey1"
    value = "value1"
  }
  secret_variable {
    name  = "skey2"
    value = "value2"
  }
  secret_variable {
    name = "skey3"
  }
}`, projectName, variableGroupName)
}

func hclVariableGroupAzureKeyVault(projectName, variableGroupName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_azurerm" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%sAzureRM"
  credentials {
    serviceprincipalid  = "%s"
    serviceprincipalkey = "%s"
  }
  azurerm_spn_tenantid                   = "%s"
  azurerm_subscription_id                = "%s"
  azurerm_subscription_name              = "%s"
  service_endpoint_authentication_scheme = "ServicePrincipal"
}

resource "azuredevops_variable_group" "test" {
  project_id   = azuredevops_project.test.id
  name         = "%s"
  description  = "A sample variable group."
  allow_access = false
  key_vault {
    name                = "%s"
    service_endpoint_id = azuredevops_serviceendpoint_azurerm.test.id
  }
  variable {
    name = "key1"
  }
}
`, projectName, projectName, os.Getenv("TEST_SERVICE_PRINCIPAL_ID"), os.Getenv("TEST_SERVICE_PRINCIPAL_KEY"),
		os.Getenv("TEST_ARM_TENANT_ID"), os.Getenv("TEST_ARM_SUBSCRIPTION_ID"), os.Getenv("TEST_ARM_SUBSCRIPTION_NAME"),
		variableGroupName, os.Getenv("TEST_ARM_KV_NAME"))
}
