//go:build (all || resource_workitemtrackingprocess || resource_workitemtrackingprocess_inherited_page) && !exclude_resource_workitemtrackingprocess
// +build all resource_workitemtrackingprocess resource_workitemtrackingprocess_inherited_page
// +build !exclude_resource_workitemtrackingprocess

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

func TestAccWorkitemtrackingprocessInheritedPage_Basic(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_inherited_page.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicInheritedPage(workItemTypeName, processName, "Custom label"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: inheritedPageImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkitemtrackingprocessInheritedPage_Update(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_inherited_page.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicInheritedPage(workItemTypeName, processName, "Custom label"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: inheritedPageImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: updatedInheritedPage(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: inheritedPageImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkitemtrackingprocessInheritedPage_Revert(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	inheritedPageNode := "azuredevops_workitemtrackingprocess_inherited_page.test"
	customLabel := "Custom label"

	var pageId string
	var processId string
	var witRefName string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicInheritedPage(workItemTypeName, processName, customLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(inheritedPageNode, "id", func(value string) error {
						pageId = value
						return nil
					}),
					resource.TestCheckResourceAttrWith(inheritedPageNode, "process_id", func(value string) error {
						processId = value
						return nil
					}),
					resource.TestCheckResourceAttrWith(inheritedPageNode, "work_item_type_reference_name", func(value string) error {
						witRefName = value
						return nil
					}),
				),
			},
			{
				Config: removedWorkItemType(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					checkPageLabelReverted(&processId, &witRefName, &pageId, customLabel),
				),
			},
		},
	})
}

func basicInheritedPage(workItemTypeName string, processName string, label string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_inherited_page" "test" {
  process_id                    = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  page_id                       = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].id
  label                         = "%s"
}
`, workItemType, label)
}

func updatedInheritedPage(workItemTypeName string, processName string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_inherited_page" "test" {
  process_id                    = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  page_id                       = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].id
  label                         = "Updated label"
}
`, workItemType)
}

func removedWorkItemType(workItemTypeName string, processName string) string {
	return basicWorkItemType(workItemTypeName, processName)
}

func checkPageLabelReverted(processIdStr *string, witRefName *string, pageId *string, customLabel string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

		processId, err := uuid.Parse(*processIdStr)
		if err != nil {
			return fmt.Errorf("invalid process_id: %w", err)
		}

		workItemType, err := clients.WorkItemTrackingProcessClient.GetProcessWorkItemType(clients.Ctx, workitemtrackingprocess.GetProcessWorkItemTypeArgs{
			ProcessId:  &processId,
			WitRefName: witRefName,
			Expand:     &workitemtrackingprocess.GetWorkItemTypeExpandValues.Layout,
		})
		if err != nil {
			return fmt.Errorf("getting work item type: %w", err)
		}

		if workItemType == nil {
			return fmt.Errorf("work item type is nil")
		}
		if workItemType.Layout == nil {
			return fmt.Errorf("work item type layout is nil")
		}
		if workItemType.Layout.Pages == nil {
			return fmt.Errorf("work item type layout pages is nil")
		}

		var page *workitemtrackingprocess.Page
		for i := range *workItemType.Layout.Pages {
			p := &(*workItemType.Layout.Pages)[i]
			if p.Id != nil && *p.Id == *pageId {
				page = p
				break
			}
		}

		if page == nil {
			return fmt.Errorf("page %s not found", *pageId)
		}

		if page.Label == nil {
			return fmt.Errorf("page label is nil")
		}

		if *page.Label == customLabel {
			return fmt.Errorf("page label should have reverted but is still %q", customLabel)
		}

		return nil
	}
}

func inheritedPageImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		processId := rs.Primary.Attributes["process_id"]
		witRefName := rs.Primary.Attributes["work_item_type_reference_name"]
		pageId := rs.Primary.ID
		return fmt.Sprintf("%s/%s/%s", processId, witRefName, pageId), nil
	}
}
