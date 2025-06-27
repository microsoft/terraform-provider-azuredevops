//go:build (all || core || resource_project) && !exclude_resource_project

package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func TestAccProject_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_project.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclProjectBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					checkProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "process_template_id"),
					resource.TestCheckResourceAttr(tfNode, "name", projectName),
					resource.TestCheckResourceAttr(tfNode, "version_control", "Git"),
					resource.TestCheckResourceAttr(tfNode, "visibility", "private"),
					resource.TestCheckResourceAttr(tfNode, "work_item_template", "Agile"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  checkImportProject(),
			},
		},
	})
}

func TestAccProject_importByName(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_project.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclProjectBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					checkProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "process_template_id"),
					resource.TestCheckResourceAttr(tfNode, "name", projectName),
					resource.TestCheckResourceAttr(tfNode, "version_control", "Git"),
					resource.TestCheckResourceAttr(tfNode, "visibility", "private"),
					resource.TestCheckResourceAttr(tfNode, "work_item_template", "Agile"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return projectName, nil
				},
				ImportStateCheck: checkImportProject(),
			},
		},
	})
}

func TestAccProject_update(t *testing.T) {
	projectNameFirst := testutils.GenerateResourceName()
	projectNameSecond := testutils.GenerateResourceName()
	tfNode := "azuredevops_project.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclProjectBasic(projectNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_template_id"),
					resource.TestCheckResourceAttr(tfNode, "name", projectNameFirst),
					resource.TestCheckResourceAttr(tfNode, "version_control", "Git"),
					resource.TestCheckResourceAttr(tfNode, "visibility", "private"),
					resource.TestCheckResourceAttr(tfNode, "work_item_template", "Agile"),
					checkProjectExists(projectNameFirst),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  checkImportProject(),
			},
			{
				Config: hclProjectUpdate(projectNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_template_id"),
					resource.TestCheckResourceAttr(tfNode, "name", projectNameSecond),
					resource.TestCheckResourceAttr(tfNode, "version_control", "Git"),
					resource.TestCheckResourceAttr(tfNode, "visibility", "public"),
					resource.TestCheckResourceAttr(tfNode, "work_item_template", "Agile"),
					checkProjectExists(projectNameSecond),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  checkImportProject(),
			},
		},
	})
}

func TestAccProject_features(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_project.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclProjectFeature(projectName, "disabled", "disabled"),
				Check: resource.ComposeTestCheckFunc(
					checkProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "process_template_id"),
					resource.TestCheckResourceAttr(tfNode, "name", projectName),
					resource.TestCheckResourceAttr(tfNode, "version_control", "Git"),
					resource.TestCheckResourceAttr(tfNode, "visibility", "private"),
					resource.TestCheckResourceAttr(tfNode, "work_item_template", "Agile"),
					resource.TestCheckResourceAttr(tfNode, "features.testplans", "disabled"),
					resource.TestCheckResourceAttr(tfNode, "features.artifacts", "disabled"),
				),
			},
			{
				ResourceName:            tfNode,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"features.testplans", "features.artifacts", "features.%"},
				ImportStateCheck:        checkImportProject(),
			},
			{
				Config: hclProjectFeature(projectName, "enabled", "disabled"),
				Check: resource.ComposeTestCheckFunc(
					checkProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "process_template_id"),
					resource.TestCheckResourceAttr(tfNode, "name", projectName),
					resource.TestCheckResourceAttr(tfNode, "version_control", "Git"),
					resource.TestCheckResourceAttr(tfNode, "visibility", "private"),
					resource.TestCheckResourceAttr(tfNode, "work_item_template", "Agile"),
					resource.TestCheckResourceAttr(tfNode, "features.testplans", "enabled"),
					resource.TestCheckResourceAttr(tfNode, "features.artifacts", "disabled"),
				),
			},
			{
				ResourceName:            tfNode,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"features.testplans", "features.artifacts", "features.%"},
				ImportStateCheck:        checkImportProject(),
			},
		},
	})
}

func TestAccProject_requireImportError(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_project.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclProjectBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					checkProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "process_template_id"),
					resource.TestCheckResourceAttr(tfNode, "name", projectName),
					resource.TestCheckResourceAttr(tfNode, "version_control", "Git"),
					resource.TestCheckResourceAttr(tfNode, "visibility", "private"),
					resource.TestCheckResourceAttr(tfNode, "work_item_template", "Agile"),
				),
			},
			{
				Config:      hclProjectImport(projectName),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Error:  creating project: TF200019: The following project already exists on the Azure DevOps Server: %s. You cannot create a new project with the same name as an existing project. Provide a different name`, projectName)),
			},
		},
	})
}

func checkProjectExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		state, ok := s.RootModule().Resources["azuredevops_project.test"]
		if !ok {
			return fmt.Errorf("Did not find a project in the TF state")
		}

		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
		id := state.Primary.ID
		project, err := clients.CoreClient.GetProject(clients.Ctx, core.GetProjectArgs{
			ProjectId:           &id,
			IncludeCapabilities: converter.Bool(true),
			IncludeHistory:      converter.Bool(false),
		})

		if err != nil {
			return fmt.Errorf("Project with ID=%s cannot be found!. Error=%v", id, err)
		}

		if *project.Name != expectedName {
			return fmt.Errorf("Project with ID=%s has Name=%s, but expected Name=%s", id, *project.Name, expectedName)
		}

		return nil
	}
}

func checkImportProject() resource.ImportStateCheckFunc {
	return func(states []*terraform.InstanceState) error {
		if len(states) != 1 {
			return fmt.Errorf("Expected project imported but not found in the state.")
		}
		state := states[0]
		projectId := state.ID
		if err := uuid.Validate(projectId); err != nil {
			return fmt.Errorf("Project ID should be a valid UUID. Got: %s", projectId)
		}
		return nil
	}
}

func hclProjectBasic(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "%[1]s-description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}`, name)
}

func hclProjectUpdate(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "%[1]s-description-update"
  visibility         = "public"
  version_control    = "Git"
  work_item_template = "Agile"
}`, name)
}

func hclProjectFeature(projectName, testPlans, stateArtifacts string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%s"
  description        = "%s-description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"

  features = {
    "testplans" = "%s"
    "artifacts" = "%s"
  }
}`, projectName, projectName, testPlans, stateArtifacts)
}

func hclProjectImport(name string) string {
	template := hclProjectBasic(name)
	return fmt.Sprintf(`
%s

resource "azuredevops_project" "import" {
  name               = azuredevops_project.test.name
  description        = azuredevops_project.test.description
  visibility         = azuredevops_project.test.visibility
  version_control    = azuredevops_project.test.version_control
  work_item_template = azuredevops_project.test.work_item_template
}`, template)
}
