package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkItemTrackingField_Basic(t *testing.T) {
	fieldName := testutils.GenerateFieldName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldBasic(fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", fieldName),
					resource.TestCheckResourceAttrSet(tfNode, "reference_name"),
					resource.TestCheckResourceAttr(tfNode, "type", "string"),
					resource.TestCheckNoResourceAttr(tfNode, "description"),
					resource.TestCheckResourceAttr(tfNode, "usage", "workItem"),
					resource.TestCheckResourceAttr(tfNode, "read_only", "false"),
					resource.TestCheckResourceAttr(tfNode, "can_sort_by", "true"),
					resource.TestCheckResourceAttr(tfNode, "is_queryable", "true"),
					resource.TestCheckResourceAttr(tfNode, "is_identity", "false"),
					resource.TestCheckResourceAttr(tfNode, "is_picklist", "false"),
					resource.TestCheckResourceAttr(tfNode, "is_picklist_suggested", "false"),
					resource.TestCheckNoResourceAttr(tfNode, "picklist_id"),
					resource.TestCheckResourceAttr(tfNode, "is_locked", "false"),
					resource.TestCheckResourceAttr(tfNode, "is_deleted", "false"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttrSet(tfNode, "supported_operations.#"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: computeFieldImportID(tfNode),
			},
		},
	})
}

func TestAccWorkItemTrackingField_Complete(t *testing.T) {
	fieldName := testutils.GenerateFieldName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldComplete(fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", fieldName),
					resource.TestCheckResourceAttrSet(tfNode, "reference_name"),
					resource.TestCheckResourceAttr(tfNode, "type", "string"),
					resource.TestCheckResourceAttr(tfNode, "description", "Test field description"),
					resource.TestCheckResourceAttr(tfNode, "usage", "workItem"),
					resource.TestCheckResourceAttr(tfNode, "read_only", "false"),
					resource.TestCheckResourceAttr(tfNode, "can_sort_by", "true"),
					resource.TestCheckResourceAttr(tfNode, "is_queryable", "true"),
					resource.TestCheckResourceAttr(tfNode, "is_identity", "false"),
					resource.TestCheckResourceAttr(tfNode, "is_picklist", "false"),
					resource.TestCheckResourceAttr(tfNode, "is_picklist_suggested", "false"),
					resource.TestCheckResourceAttr(tfNode, "is_locked", "false"),
					resource.TestCheckResourceAttr(tfNode, "is_deleted", "false"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttrSet(tfNode, "supported_operations.#"),
					resource.TestCheckNoResourceAttr(tfNode, "picklist_id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: computeFieldImportID(tfNode),
			},
		},
	})
}

func TestAccWorkItemTrackingField_Boolean(t *testing.T) {
	fieldName := testutils.GenerateFieldName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldBoolean(fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", fieldName),
					resource.TestCheckResourceAttr(tfNode, "type", "boolean"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: computeFieldImportID(tfNode),
			},
		},
	})
}

func TestAccWorkItemTrackingField_Update(t *testing.T) {
	fieldName := testutils.GenerateFieldName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldBasic(fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", fieldName),
					resource.TestCheckResourceAttr(tfNode, "is_locked", "false"),
					resource.TestCheckResourceAttr(tfNode, "is_deleted", "false"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: computeFieldImportID(tfNode),
			},
			{
				Config: fieldUpdated(fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", fieldName),
					resource.TestCheckResourceAttr(tfNode, "is_locked", "true"),
					resource.TestCheckResourceAttr(tfNode, "is_deleted", "false"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: computeFieldImportID(tfNode),
			},
		},
	})
}

func TestAccWorkItemTrackingField_ProjectScoped(t *testing.T) {
	fieldName := testutils.GenerateFieldName()
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldProjectScoped(projectName, fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", fieldName),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: computeFieldImportID(tfNode),
			},
		},
	})
}

func computeFieldImportID(resourceNode string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceNode]
		if !ok {
			return "", fmt.Errorf("Resource node not found: %s", resourceNode)
		}
		projectID := rs.Primary.Attributes["project_id"]
		if projectID != "" {
			return fmt.Sprintf("%s/%s", projectID, rs.Primary.Attributes["id"]), nil
		}
		return rs.Primary.Attributes["id"], nil
	}
}

func fieldBasic(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtracking_field" "test" {
  name           = "%s"
  reference_name = "Custom.%s"
  type           = "string"
}
`, name, name)
}

func fieldComplete(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtracking_field" "test" {
  name                  = "%s"
  reference_name        = "Custom.%s"
  type                  = "string"
  description           = "Test field description"
  usage                 = "workItem"
  read_only             = false
  can_sort_by           = true
  is_queryable          = true
  is_identity           = false
  is_picklist           = false
  is_picklist_suggested = false
  is_locked             = false
  is_deleted            = false
}
`, name, name)
}

func fieldBoolean(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtracking_field" "test" {
  name           = "%s"
  reference_name = "Custom.%s"
  type           = "boolean"
  description    = "A boolean field for testing"
}
`, name, name)
}

func fieldUpdated(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtracking_field" "test" {
  name           = "%s"
  reference_name = "Custom.%s"
  type           = "string"
  is_locked      = true
}
`, name, name)
}

func fieldProjectScoped(projectName, fieldName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%s"
  description        = "Test project for field"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_workitemtracking_field" "test" {
  project_id     = azuredevops_project.test.id
  name           = "%s"
  reference_name = "Custom.%s"
  type           = "string"
}
`, projectName, fieldName, fieldName)
}
