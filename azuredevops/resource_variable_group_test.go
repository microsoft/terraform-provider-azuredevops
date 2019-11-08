package azuredevops

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/stretchr/testify/require"
)

var testVarGroupProjectID = uuid.New().String()

// This definition matches the overall structure of what a configured git repository would
// look like. Note that the ID and Name attributes match -- this is the service-side behavior
// when configuring a GitHub repo.
var testVariableGroup = taskagent.VariableGroup{
	Id:          converter.Int(100),
	Name:        converter.String("Name"),
	Description: converter.String("This is a test variable group."),
	Variables: &map[string]taskagent.VariableValue{
		"var1": {
			Value:    converter.String("value1"),
			IsSecret: converter.Bool(false),
		},
	},
}

/**
 * Begin unit tests
 */
// verifies that the flatten/expand round trip yields the same build definition
func TestAzureDevOpsVariableGroup_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceVariableGroup().Schema, nil)
	flattenVariableGroup(resourceData, &testVariableGroup, &testVarGroupProjectID)

	variableGroupParams, projectID := expandVariableGroupParameters(resourceData)

	require.Equal(t, *testVariableGroup.Name, *variableGroupParams.Name)
	require.Equal(t, *testVariableGroup.Description, *variableGroupParams.Description)
	require.Equal(t, *testVariableGroup.Variables, *variableGroupParams.Variables)
	require.Equal(t, testVarGroupProjectID, *projectID)
}

/**
 * Begin acceptance tests
 */

func TestAccAzureDevOpsVariableGroup_CreateAndUpdate(t *testing.T) {
	projectName := testAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	vargroupNameFirst := testAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	vargroupNameSecond := testAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	tfVarGroupNode := "azuredevops_variable_group.vg"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccVariableGroupCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVariableGroupResource(projectName, vargroupNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfVarGroupNode, "project_id"),
					resource.TestCheckResourceAttr(tfVarGroupNode, "name", vargroupNameFirst),
					testAccCheckVariableGroupResourceExists(vargroupNameFirst),
				),
			}, {
				Config: testAccVariableGroupResource(projectName, vargroupNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfVarGroupNode, "project_id"),
					resource.TestCheckResourceAttr(tfVarGroupNode, "name", vargroupNameSecond),
					testAccCheckVariableGroupResourceExists(vargroupNameSecond),
				),
			},
		},
	})
}

// HCL describing an AzDO variable group
func testAccVariableGroupResource(projectName string, variableGroupName string) string {
	variableGroupResource := fmt.Sprintf(`
resource "azuredevops_variable_group" "vg" {
	project_id  = azuredevops_project.project.id
	name        = "%s"
	description = "A sample variable group."

	variable {
		name      = "key1"
		value     = "value1"
		is_secret = true
	}
	
	variable {
		name  = "key2"
		value = "value2"
	}

	variable {
		name = "key3"
	}
}`, variableGroupName)

	projectResource := testAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, variableGroupResource)
}

// Given an AzDO variable group name, this will return a function that will check whether
// or not the definition (1) exists in the state, (2) exists in AzDO, and (3) has the correct
// or expected name
func testAccCheckVariableGroupResourceExists(expectedName string) resource.TestCheckFunc {
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
