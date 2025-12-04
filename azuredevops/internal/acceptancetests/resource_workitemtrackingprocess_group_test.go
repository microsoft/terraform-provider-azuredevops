//go:build (all || resource_workitemtrackingprocess || resource_workitemtrackingprocess_group) && !exclude_resource_workitemtrackingprocess
// +build all resource_workitemtrackingprocess resource_workitemtrackingprocess_group
// +build !exclude_resource_workitemtrackingprocess

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkitemtrackingprocessGroup_Basic(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicGroup(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_id"),
					resource.TestCheckResourceAttrSet(tfNode, "work_item_type_reference_name"),
					resource.TestCheckResourceAttrSet(tfNode, "page_id"),
					resource.TestCheckResourceAttrSet(tfNode, "section_id"),
					resource.TestCheckResourceAttr(tfNode, "label", "Test Group"),
					resource.TestCheckResourceAttr(tfNode, "visible", "true"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: groupImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkitemtrackingprocessGroup_Update(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicGroup(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_id"),
					resource.TestCheckResourceAttrSet(tfNode, "work_item_type_reference_name"),
					resource.TestCheckResourceAttrSet(tfNode, "page_id"),
					resource.TestCheckResourceAttrSet(tfNode, "section_id"),
					resource.TestCheckResourceAttr(tfNode, "label", "Test Group"),
					resource.TestCheckResourceAttr(tfNode, "visible", "true"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: groupImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: updatedGroup(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_id"),
					resource.TestCheckResourceAttrSet(tfNode, "work_item_type_reference_name"),
					resource.TestCheckResourceAttrSet(tfNode, "page_id"),
					resource.TestCheckResourceAttrSet(tfNode, "section_id"),
					resource.TestCheckResourceAttr(tfNode, "label", "Updated Group"),
					resource.TestCheckResourceAttr(tfNode, "visible", "false"),
					resource.TestCheckResourceAttr(tfNode, "order", "2"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: groupImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkitemtrackingprocessGroup_WithCustomId(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: groupWithCustomId(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_id"),
					resource.TestCheckResourceAttrSet(tfNode, "work_item_type_reference_name"),
					resource.TestCheckResourceAttrSet(tfNode, "page_id"),
					resource.TestCheckResourceAttrSet(tfNode, "section_id"),
					resource.TestCheckResourceAttr(tfNode, "label", "Custom ID Group"),
					resource.TestCheckResourceAttr(tfNode, "id", "custom-group-id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: groupImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkitemtrackingprocessGroup_Move(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_group.test"

	var originalSectionId string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicGroup(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_id"),
					resource.TestCheckResourceAttrSet(tfNode, "work_item_type_reference_name"),
					resource.TestCheckResourceAttrSet(tfNode, "page_id"),
					resource.TestCheckResourceAttrSet(tfNode, "section_id"),
					resource.TestCheckResourceAttr(tfNode, "label", "Test Group"),
					resource.TestCheckResourceAttr(tfNode, "visible", "true"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrWith(tfNode, "section_id", func(value string) error {
						originalSectionId = value
						return nil
					}),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: groupImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: movedGroup(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_id"),
					resource.TestCheckResourceAttrSet(tfNode, "work_item_type_reference_name"),
					resource.TestCheckResourceAttrSet(tfNode, "page_id"),
					resource.TestCheckResourceAttrSet(tfNode, "section_id"),
					resource.TestCheckResourceAttr(tfNode, "label", "Test Group"),
					resource.TestCheckResourceAttr(tfNode, "visible", "true"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrWith(tfNode, "section_id", func(value string) error {
						if value == originalSectionId {
							return fmt.Errorf("section_id should have changed, but is still %s", originalSectionId)
						}
						return nil
					}),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: groupImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func basicGroup(workItemTypeName string, processName string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_group" "test" {
  process_id                     = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  page_id                        = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].id
  section_id                     = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].sections[0].id
  label                          = "Test Group"
}
`, workItemType)
}

func updatedGroup(workItemTypeName string, processName string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_group" "test" {
  process_id                     = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  page_id                        = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].id
  section_id                     = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].sections[0].id
  label                          = "Updated Group"
  visible                        = false
  order                          = 2
}
`, workItemType)
}

func groupWithCustomId(workItemTypeName string, processName string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_group" "test" {
  process_id                     = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  page_id                        = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].id
  section_id                     = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].sections[0].id
  label                          = "Custom ID Group"
  id                             = "custom-group-id"
}
`, workItemType)
}

func movedGroup(workItemTypeName string, processName string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_group" "test" {
  process_id                     = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  page_id                        = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].id
  section_id                     = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].sections[1].id
  label                          = "Test Group"
}
`, workItemType)
}

func groupImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		processId := rs.Primary.Attributes["process_id"]
		witRefName := rs.Primary.Attributes["work_item_type_reference_name"]
		pageId := rs.Primary.Attributes["page_id"]
		sectionId := rs.Primary.Attributes["section_id"]
		groupId := rs.Primary.ID
		return fmt.Sprintf("%s/%s/%s/%s/%s", processId, witRefName, pageId, sectionId, groupId), nil
	}
}
