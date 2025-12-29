package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttrSet(tfNode, "supported_operations.#"),
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
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttrSet(tfNode, "supported_operations.#"),
					resource.TestCheckNoResourceAttr(tfNode, "picklist_id"),
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
	fieldName := "testaccb2i4ttpqo0"
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
			},
		},
	})
}

func TestAccWorkItemTrackingField_Lock(t *testing.T) {
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
					resource.TestCheckResourceAttr(tfNode, "name", fieldName),
					resource.TestCheckResourceAttr(tfNode, "is_locked", "true"),
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
					resource.TestCheckResourceAttr(tfNode, "name", fieldName),
					resource.TestCheckResourceAttr(tfNode, "is_locked", "false"),
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
