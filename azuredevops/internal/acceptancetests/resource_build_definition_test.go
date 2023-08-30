//go:build (all || resource_build_definition) && !exclude_resource_build_definition
// +build all resource_build_definition
// +build !exclude_resource_build_definition

package acceptancetests

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/build"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

func TestAccBuildDefinition_Basic(t *testing.T) {
	name := testutils.GenerateResourceName()

	tfBuildDefNode := "azuredevops_build_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclBuildDefinitionPath(name, `\\`),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", name),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", `\`),
				),
			}, {
				ResourceName:            tfBuildDefNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfBuildDefNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_first_run"},
			},
		},
	})
}

func TestAccBuildDefinition_PathUpdate(t *testing.T) {
	name := testutils.GenerateResourceName()

	pathFirst := `\\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	pathSecond := `\\` + name + `\\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	tfBuildDefNode := "azuredevops_build_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclBuildDefinitionPath(name, pathFirst),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", name),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", strings.ReplaceAll(pathFirst, `\\`, `\`)),
				),
			},
			{
				Config: hclBuildDefinitionPath(name, pathSecond),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", name),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", strings.ReplaceAll(pathSecond, `\\`, `\`)),
				),
			}, {
				ResourceName:            tfBuildDefNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfBuildDefNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_first_run"},
			},
		},
	})
}

// Verifies a build for with variables can create and update, including secret variables
func TestAccBuildDefinition_WithVariables(t *testing.T) {
	name := testutils.GenerateResourceName()
	tfNode := "azuredevops_build_definition.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclBuildDefinitionVariable("foo1", "bar1", name),
				Check:  checkForVariableValues(tfNode, "foo1", "bar1"),
			}, {
				Config: hclBuildDefinitionVariable("foo2", "bar2", name),
				Check:  checkForVariableValues(tfNode, "foo2", "bar2"),
			},
		},
	})
}

func TestAccBuildDefinition_Schedules(t *testing.T) {
	name := testutils.GenerateResourceName()
	tfNode := "azuredevops_build_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclBuildDefinitionSchedules(name),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "revision"),
					resource.TestCheckResourceAttrSet(tfNode, "repository.0.repo_id"),
					resource.TestCheckResourceAttr(tfNode, "schedules.#", "1"),
					resource.TestCheckResourceAttr(tfNode, "schedules.0.days_to_build.#", "1"),
					resource.TestCheckResourceAttr(tfNode, "name", name),
				),
			},
		},
	})
}

// Checks that the expected variable values exist in the state
func checkForVariableValues(tfNode string, expectedVals ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rootModule := s.RootModule()
		resource, ok := rootModule.Resources[tfNode]
		if !ok {
			return fmt.Errorf("Did not find resource in TF state")
		}

		is := resource.Primary
		if is == nil {
			return fmt.Errorf("No primary instance: %s in %s", tfNode, rootModule.Path)
		}

		for _, expectedVal := range expectedVals {
			found := false
			for _, value := range is.Attributes {
				if value == expectedVal {
					found = true
				}
			}

			if !found {
				return fmt.Errorf("Did not find variable with value %s", expectedVal)
			}

		}

		return nil
	}
}

func checkBuildDefinitionExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		buildDef, ok := s.RootModule().Resources["azuredevops_build_definition.test"]
		if !ok {
			return fmt.Errorf("Did not find a build definition in the TF state")
		}

		buildDefinition, err := getBuildDefinitionFromResource(buildDef)
		if err != nil {
			return err
		}

		if *buildDefinition.Name != expectedName {
			return fmt.Errorf("Build Definition has Name=%s, but expected Name=%s", *buildDefinition.Name, expectedName)
		}

		return nil
	}
}

// verifies that all build definitions referenced in the state are destroyed. This will be invoked
// *after* terraform destroys the resource but *before* the state is wiped clean.
func checkBuildDefinitionDestroyed(s *terraform.State) error {
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_build_definition" {
			continue
		}

		// indicates the build definition still exists - this should fail the test
		if _, err := getBuildDefinitionFromResource(resource); err == nil {
			return fmt.Errorf("Unexpectedly found a build definition that should be deleted")
		}
	}

	return nil
}

// given a resource from the state, return a build definition (and error)
func getBuildDefinitionFromResource(resource *terraform.ResourceState) (*build.BuildDefinition, error) {
	buildDefID, err := strconv.Atoi(resource.Primary.ID)
	if err != nil {
		return nil, err
	}

	projectID := resource.Primary.Attributes["project_id"]
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
	return clients.BuildClient.GetDefinition(clients.Ctx, build.GetDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &buildDefID,
	})
}

func hclBuildDefinitionTemplate(name string) string {
	return fmt.Sprintf(`

resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "%[1]s-description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "acc-%[1]s"
  initialization {
    init_type = "Clean"
  }
}`, name)
}

func hclBuildDefinitionPath(name, path string) string {
	template := hclBuildDefinitionTemplate(name)
	return fmt.Sprintf(`
%s

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  path       = "%[3]s"

  ci_trigger {
    override {
      batch = true
      branch_filter {
        include = ["master"]
      }
      path_filter {
        include = ["*/**.ts"]
      }
      max_concurrent_builds_per_branch = 2
      polling_interval                 = 0
    }
  }

  repository {
    repo_type   = "TfsGit"
    repo_id     = azuredevops_git_repository.test.id
    branch_name = azuredevops_git_repository.test.default_branch
    yml_path    = "azure-pipelines.yml"
  }
}
`, template, name, path)
}

func hclBuildDefinitionVariable(name, varVal, secretVarVal string) string {
	template := hclBuildDefinitionTemplate(name)
	return fmt.Sprintf(`
%s

resource "azuredevops_build_definition" "test" {
		project_id = azuredevops_project.test.id
		name       = "%[2]s"
		repository {
			repo_type   = "TfsGit"
			repo_id     = azuredevops_git_repository.test.id
			branch_name = azuredevops_git_repository.test.default_branch
			yml_path    = "azure-pipelines.yml"
		}

		variable {
			name  = "FOO_VAR"
			value = "%[3]s"
		}

		variable {
			name      = "BAR_VAR"
			secret_value     = "%[4]s"
			is_secret = true
		}
	}`, template, name, varVal, secretVarVal)
}

func hclBuildDefinitionSchedules(name string) string {
	template := hclBuildDefinitionTemplate(name)
	return fmt.Sprintf(`
%s

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  path       = "\\ExampleFolder"

  ci_trigger {
    override {
      batch = true
      branch_filter {
        include = ["master"]
      }
      path_filter {
        include = ["*/**.ts"]
      }
      max_concurrent_builds_per_branch = 2
      polling_interval                 = 0
    }
  }

  schedules {
    branch_filter {
      include = ["master"]
    }

    days_to_build              = ["Mon"]
    schedule_only_with_changes = true
    start_hours                = 0
    start_minutes              = 0
    time_zone                  = "(UTC) Coordinated Universal Time"
  }

  repository {
    repo_type   = "TfsGit"
    repo_id     = azuredevops_git_repository.test.id
    branch_name = azuredevops_git_repository.test.default_branch
    yml_path    = "azure-pipelines.yml"
  }
}
`, template, name)
}
