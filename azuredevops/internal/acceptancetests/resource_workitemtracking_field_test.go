package acceptancetests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

func TestAccWorkItemTrackingField_Basic(t *testing.T) {
	fieldName := generateFieldName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldBasic(fieldName),
				Check: resource.ComposeTestCheckFunc(
					// Computed attributes
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttrSet(tfNode, "supported_operations.#"),
					// Default values
					resource.TestCheckResourceAttr(tfNode, "usage", "workItem"),
					resource.TestCheckResourceAttr(tfNode, "read_only", "false"),
					resource.TestCheckResourceAttr(tfNode, "can_sort_by", "true"),
					resource.TestCheckResourceAttr(tfNode, "is_queryable", "true"),
					resource.TestCheckResourceAttr(tfNode, "is_identity", "false"),
					resource.TestCheckResourceAttr(tfNode, "is_picklist", "false"),
					resource.TestCheckResourceAttr(tfNode, "is_picklist_suggested", "false"),
					resource.TestCheckResourceAttr(tfNode, "is_locked", "false"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkItemTrackingField_Complete(t *testing.T) {
	fieldName := generateFieldName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldComplete(fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkItemTrackingField_Boolean(t *testing.T) {
	fieldName := generateFieldName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldBoolean(fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkItemTrackingField_Html(t *testing.T) {
	fieldName := generateFieldName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldHtml(fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkItemTrackingField_Integer(t *testing.T) {
	fieldName := generateFieldName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldWithType(fieldName, "integer"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkItemTrackingField_DateTime(t *testing.T) {
	fieldName := generateFieldName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldWithType(fieldName, "dateTime"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkItemTrackingField_PlainText(t *testing.T) {
	fieldName := generateFieldName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldWithType(fieldName, "plainText"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkItemTrackingField_Double(t *testing.T) {
	fieldName := generateFieldName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldWithType(fieldName, "double"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkItemTrackingField_Identity(t *testing.T) {
	fieldName := generateFieldName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldWithType(fieldName, "identity"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkItemTrackingField_TreePath(t *testing.T) {
	fieldName := generateFieldName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldWithType(fieldName, "treePath"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkItemTrackingField_History(t *testing.T) {
	fieldName := generateFieldName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldWithType(fieldName, "history"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkItemTrackingField_Guid(t *testing.T) {
	fieldName := generateFieldName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldWithType(fieldName, "guid"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkItemTrackingField_Lock(t *testing.T) {
	fieldName := generateFieldName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldBasic(fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: lockField(fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkItemTrackingField_Restore(t *testing.T) {
	fieldName := generateFieldName()
	tfNode := "azuredevops_workitemtracking_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fieldBasic(fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: "# empty config to delete the field",
			},
			{
				Config: restoreField(fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:            tfNode,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"restore"},
			},
		},
	})
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

func fieldWithType(name string, fieldType string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtracking_field" "test" {
  name           = "%s"
  reference_name = "Custom.%s"
  type           = "%s"
}
`, name, name, fieldType)
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

func fieldHtml(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtracking_field" "test" {
  name           = "%s"
  reference_name = "Custom.%s"
  type           = "html"
  description    = "A html field for testing"
}
`, name, name)
}

func lockField(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtracking_field" "test" {
  name           = "%s"
  reference_name = "Custom.%s"
  type           = "string"
  is_locked      = true
}
`, name, name)
}

func restoreField(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtracking_field" "test" {
  name           = "%s"
  reference_name = "Custom.%s"
  type           = "string"
  restore        = true
}
`, name, name)
}

// generateFieldName generates a valid field name without hyphens or other invalid characters
func generateFieldName() string {
	return strings.ReplaceAll(testutils.GenerateResourceName(), "-", "")
}

// checkFieldDestroyed verifies that all fields referenced in the state are destroyed. This will be invoked
// *after* terraform destroys the resource but *before* the state is wiped clean.
func checkFieldDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

	for _, res := range s.RootModule().Resources {
		if res.Type != "azuredevops_workitemtracking_field" {
			continue
		}

		referenceName := res.Primary.ID

		_, err := clients.WorkItemTrackingClient.GetWorkItemField(clients.Ctx, workitemtracking.GetWorkItemFieldArgs{
			FieldNameOrRefName: &referenceName,
		})
		if utils.ResponseWasNotFound(err) {
			continue
		}

		if err != nil {
			return err
		}
	}

	return nil
}
