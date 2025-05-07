//go:build (all || core || workitem) && !exclude_resource_workitem
// +build all core workitem
// +build !exclude_resource_workitem

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkItem_basic(t *testing.T) {
	workItemTitle := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitem.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: workItemBasic(projectName, workItemTitle),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttr(tfNode, "title", workItemTitle),
					resource.TestCheckResourceAttr(tfNode, "type", "Issue"),
					resource.TestCheckResourceAttr(tfNode, "state", "Active"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
			},
		},
	})
}

func TestAccWorkItem_titleUpdate(t *testing.T) {
	workItemTitle := testutils.GenerateResourceName()
	workItemTitleUpdated := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitem.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: workItemBasic(projectName, workItemTitle),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "title", workItemTitle),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
			},
			{
				Config: workItemBasic(projectName, workItemTitleUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "title", workItemTitleUpdated),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
			},
		},
	})
}

func TestAccWorkItem_tagUpdate(t *testing.T) {
	workItemTitle := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitem.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: workItemBasic(projectName, workItemTitle),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "title", workItemTitle),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
			},
			{
				Config: workItemTagUpdate(projectName, workItemTitle),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "title", workItemTitle),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
			},
			{
				Config: workItemBasic(projectName, workItemTitle),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "title", workItemTitle),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
			},
		},
	})
}

func TestAccWorkItem_parent(t *testing.T) {
	workItemTitle := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitem.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: workItemParent(projectName, workItemTitle),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "title", workItemTitle),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttr(tfNode, "type", "Issue"),
					resource.TestCheckResourceAttr(tfNode, "state", "Active"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
			},
		},
	})
}

func TestAccWorkItem_parentUpdate(t *testing.T) {
	workItemTitle := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitem.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: workItemParent(projectName, workItemTitle),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "title", workItemTitle),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttr(tfNode, "type", "Issue"),
					resource.TestCheckResourceAttr(tfNode, "state", "Active"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
			},
			{
				Config: workItemParentUpdate(projectName, workItemTitle),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "title", workItemTitle),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttr(tfNode, "type", "Issue"),
					resource.TestCheckResourceAttr(tfNode, "state", "Active"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
			},
		},
	})
}

func TestAccWorkItem_parentDelete(t *testing.T) {
	workItemTitle := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitem.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: workItemParent(projectName, workItemTitle),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "title", workItemTitle),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttr(tfNode, "type", "Issue"),
					resource.TestCheckResourceAttr(tfNode, "state", "Active"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
			},
			{
				Config: workItemParentDelete(projectName, workItemTitle),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "title", workItemTitle),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttr(tfNode, "type", "Issue"),
					resource.TestCheckResourceAttr(tfNode, "state", "Active"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
			},
		},
	})
}

func workItemTemplate(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
  name               = "%[1]s"
  description        = "%[1]s-description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}
`, name)
}

func workItemBasic(projectNane string, title string) string {
	template := workItemTemplate(projectNane)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitem" "test" {
  title      = "%s"
  project_id = azuredevops_project.project.id
  type       = "Issue"
}
`, template, title)
}

func workItemTagUpdate(projectNane string, title string) string {
	template := workItemTemplate(projectNane)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitem" "test" {
  title      = "%s"
  project_id = azuredevops_project.project.id
  type       = "Issue"
  state      = "Active"
  tags       = ["tag1", "tag2"]
}
`, template, title)
}

func workItemParent(projectNane string, title string) string {
	template := workItemTemplate(projectNane)
	return fmt.Sprintf(`
%[1]s

resource "azuredevops_workitem" "parent" {
  title      = "%[2]s Parent"
  project_id = azuredevops_project.project.id
  type       = "Issue"
}

resource "azuredevops_workitem" "test" {
  title      = "%[2]s"
  project_id = azuredevops_project.project.id
  type       = "Issue"
  parent_id  = azuredevops_workitem.parent.id
}
`, template, title)
}

func workItemParentDelete(projectNane string, title string) string {
	template := workItemTemplate(projectNane)
	return fmt.Sprintf(`
%[1]s

resource "azuredevops_workitem" "parent" {
  title      = "%[2]s Parent"
  project_id = azuredevops_project.project.id
  type       = "Issue"
}

resource "azuredevops_workitem" "test" {
  title      = "%[2]s"
  project_id = azuredevops_project.project.id
  type       = "Issue"
}
`, template, title)
}

func workItemParentUpdate(projectNane string, title string) string {
	template := workItemTemplate(projectNane)
	return fmt.Sprintf(`
%[1]s

resource "azuredevops_workitem" "parent" {
  title      = "%[2]s Parent"
  project_id = azuredevops_project.project.id
  type       = "Issue"
}

resource "azuredevops_workitem" "parent2" {
  title      = "%[2]s Parent2"
  project_id = azuredevops_project.project.id
  type       = "Issue"
}

resource "azuredevops_workitem" "test" {
  project_id = azuredevops_project.project.id
  title      = "%[2]s"
  type       = "Issue"
  parent_id  = azuredevops_workitem.parent2.id
}
`, template, title)
}
