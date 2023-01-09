//go:build (all || resource_variable_group) && !exclude_resource_variable_group
// +build all resource_variable_group
// +build !exclude_resource_variable_group

package acceptancetests

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

var config = `data "azuredevops_project" "example" {
  name = "ADDTest"
}


resource "azuredevops_variable_group" "example" {
  project_id   = data.azuredevops_project.example.id
  name         = "Example Variable Group"
  description  = "Example Variable Group Description"
  allow_access = true

  variable {
    name  = "key1"
  }

  lifecycle {
    ignore_changes = [variable]
  }
}

resource "azuredevops_serviceendpoint_azurerm" "example" {
  project_id            = data.azuredevops_project.example.id
  service_endpoint_name = "Example AzureRM"
  description           = "Managed by Terraform"
  credentials {
    serviceprincipalid  = "52b594f4-8ebf-4371-8e93-a1fdeb9822fb"
    serviceprincipalkey = "Hl48Q~MEW8tBfM3YjJgDLXdvPk8cY3rTsRfLycP3"
  }
  azurerm_spn_tenantid      = "8c08a505-9ab1-4f32-910c-8a8639e53c4f"
  azurerm_subscription_id   = "e393adb3-b5be-4789-bdc9-848367f0d152"
  azurerm_subscription_name = "Visual Studio Enterprise"
}

resource "azuredevops_variable_group_values" "variables" {
  project_id        = data.azuredevops_project.example.id
  variable_group_id = azuredevops_variable_group.example.id

  key_vault {
    name                = "xz-ado-test-kv"
    service_endpoint_id = azuredevops_serviceendpoint_azurerm.example.id
  }

  variable {
    name = "var04"
  }

  variable {
    name = "var05"
  }
}`

func TestAccVariableGroup_DDDD(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkVariableGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccVariableGroup_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()

	vargroupNameFirst := testutils.GenerateResourceName()
	vargroupNameSecond := testutils.GenerateResourceName()
	vargroupNameNoSecret := testutils.GenerateResourceName()

	tfVarGroupNode := "azuredevops_variable_group.vg"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkVariableGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclVariableGroupResourceWithProject(projectName, vargroupNameFirst, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfVarGroupNode, "project_id"),
					resource.TestCheckResourceAttr(tfVarGroupNode, "name", vargroupNameFirst),
					checkVariableGroupExists(vargroupNameFirst, true),
				),
			}, {
				Config: testutils.HclVariableGroupResourceWithProject(projectName, vargroupNameSecond, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfVarGroupNode, "project_id"),
					resource.TestCheckResourceAttr(tfVarGroupNode, "name", vargroupNameSecond),
					checkVariableGroupExists(vargroupNameSecond, false),
				),
			}, {
				Config: testutils.HclVariableGroupResourceNoSecretsWithProject(projectName, vargroupNameNoSecret, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfVarGroupNode, "project_id"),
					resource.TestCheckResourceAttr(tfVarGroupNode, "name", vargroupNameNoSecret),
					checkVariableGroupExists(vargroupNameNoSecret, false),
				),
			}, {
				// Resource Acceptance Testing https://www.terraform.io/docs/extend/resources/import.html#resource-acceptance-testing-implementation
				ResourceName:      tfVarGroupNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfVarGroupNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVariableGroupKeyVault_CreateAndUpdate(t *testing.T) {
	t.Skip("Skipping test TestAccVariableGroupKeyVault_CreateAndUpdate: azure key vault not provisioned on test infrastructure")
	projectName := testutils.GenerateResourceName()

	vargroupKeyvault := testutils.GenerateResourceName()
	keyVaultName := "key-vault-name"
	allowAccessFalse := false
	tfVarGroupNode := "azuredevops_variable_group.vg"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkVariableGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclVariableGroupResourceKeyVaultWithProject(projectName, vargroupKeyvault, allowAccessFalse, keyVaultName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfVarGroupNode, "project_id"),
					resource.TestCheckResourceAttr(tfVarGroupNode, "name", vargroupKeyvault),
					checkVariableGroupExists(vargroupKeyvault, allowAccessFalse),
				),
			}, {
				// Resource Acceptance Testing https://www.terraform.io/docs/extend/resources/import.html#resource-acceptance-testing-implementation
				ResourceName:            tfVarGroupNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfVarGroupNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_vault.0.search_depth"},
			},
		},
	})
}

// Given an AzDO variable group name, this will return a function that will check whether
// or not the definition (1) exists in the state, (2) exists in AzDO, and (3) has the correct
// or expected name
func checkVariableGroupExists(expectedName string, expectedAllowAccess bool) resource.TestCheckFunc {
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
func checkVariableGroupDestroyed(s *terraform.State) error {
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
