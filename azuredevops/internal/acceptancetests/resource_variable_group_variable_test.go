package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	taskagentsvc "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/taskagent"
)

func TestAccVariableGroupVariable_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	vgName := testutils.GenerateResourceName()
	node := "azuredevops_variable_group_variable.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkVariableGroupVariableDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclVariableGroupVariableBasic(projectName, vgName, "foo"),
				Check: resource.ComposeTestCheckFunc(
					checkVariableGroupVariableExists(node),
				),
			},
			{
				ResourceName:      node,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: hclVariableGroupVariableBasic(projectName, vgName, "bar"),
				Check: resource.ComposeTestCheckFunc(
					checkVariableGroupVariableExists(node),
				),
			},
			{
				ResourceName:      node,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVariableGroupVariable_secret(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	vgName := testutils.GenerateResourceName()
	node := "azuredevops_variable_group_variable.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkVariableGroupVariableDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclVariableGroupVariableSecret(projectName, vgName, "foo"),
				Check: resource.ComposeTestCheckFunc(
					checkVariableGroupVariableExists(node),
				),
			},
			{
				ResourceName:            node,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_value"},
			},
			{
				Config: hclVariableGroupVariableSecret(projectName, vgName, "bar"),
				Check: resource.ComposeTestCheckFunc(
					checkVariableGroupVariableExists(node),
				),
			},
			{
				ResourceName:            node,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_value"},
			},
		},
	})
}

func checkVariableGroupVariableDestroyed(s *terraform.State) error {
	for _, res := range s.RootModule().Resources {
		if res.Type != "azuredevops_variable_group_variable" {
			continue
		}

		ok, err := checkVariableGroupVariableFromState(res)
		if err == nil && ok {
			return fmt.Errorf("variable still exists")
		}
	}

	return nil
}

func checkVariableGroupVariableExists(node string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		state, ok := s.RootModule().Resources[node]
		if !ok {
			return fmt.Errorf("Did not find a variable group in the TF state")
		}

		ok, err := checkVariableGroupVariableFromState(state)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("%q doesn't exist", node)
		}
		return nil
	}
}

func checkVariableGroupVariableFromState(resource *terraform.ResourceState) (bool, error) {
	projectId, groupId, name, err := taskagentsvc.ResourceVariableGroupVariableParseId(resource.Primary.ID)
	if err != nil {
		return false, err
	}

	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
	resp, err := clients.TaskAgentClient.GetVariableGroup(
		clients.Ctx,
		taskagent.GetVariableGroupArgs{
			GroupId: &groupId,
			Project: &projectId,
		},
	)
	if err != nil {
		return false, err
	}

	if resp.Variables == nil {
		return false, fmt.Errorf("unexpected null variables in group response")
	}

	vars := *resp.Variables
	_, ok := vars[name]
	return ok, nil
}

func TestAccVariableGroupVariable_ForEach_ConcurrentCreate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	vgName := testutils.GenerateResourceName()

	node1 := "azuredevops_variable_group_variable.example1"
	node2 := "azuredevops_variable_group_variable.example2"
	node3 := "azuredevops_variable_group_variable.example3"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkVariableGroupVariableDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclVariableGroupVariableForEach(projectName, vgName),
				Check: resource.ComposeTestCheckFunc(
					checkVariableGroupVariableExists(node1),
					checkVariableGroupVariableExists(node2),
					checkVariableGroupVariableExists(node3),
				),
			},
			{
				ResourceName:      node1,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      node2,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      node3,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclVariableGroupVariableBasic(projectName, variableGroupName, val string) string {
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
  variable {
    name         = "skey1"
    secret_value = "svalue1"
    is_secret    = true
  }
  lifecycle {
    ignore_changes = [variable]
  }
}
resource "azuredevops_variable_group_variable" "test" {
  project_id        = azuredevops_project.test.id
  variable_group_id = azuredevops_variable_group.test.id
  name              = "test-key"
  value             = "%s"
}
`, projectName, variableGroupName, val)
}

func hclVariableGroupVariableSecret(projectName, variableGroupName, val string) string {
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
  variable {
    name         = "skey1"
    secret_value = "svalue1"
    is_secret    = true
  }
  lifecycle {
    ignore_changes = [variable]
  }
}
resource "azuredevops_variable_group_variable" "test" {
  project_id        = azuredevops_project.test.id
  variable_group_id = azuredevops_variable_group.test.id
  name              = "test-key"
  secret_value      = "%s"
}
`, projectName, variableGroupName, val)
}

func hclVariableGroupVariableForEach(projectName, variableGroupName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_variable_group" "test" {
  project_id   = azuredevops_project.test.id
  name         = "%s"
  description  = "test description"
  allow_access = false

  # Seed variables; changes ignored to allow adding separate resources.
  variable {
    name  = "seed"
    value = "seed"
  }
  lifecycle {
    ignore_changes = [variable]
  }
}

resource "azuredevops_variable_group_variable" "example1" {
  project_id        = azuredevops_project.test.id
  variable_group_id = azuredevops_variable_group.test.id
  name              = "key1"
  value             = "val1"
}

resource "azuredevops_variable_group_variable" "example2" {
  project_id        = azuredevops_project.test.id
  variable_group_id = azuredevops_variable_group.test.id
  name              = "key2"
  value             = "val2"
}


resource "azuredevops_variable_group_variable" "example3" {
  project_id        = azuredevops_project.test.id
  variable_group_id = azuredevops_variable_group.test.id
  name              = "key3"
  value             = "val3"
}
`, projectName, variableGroupName)
}
