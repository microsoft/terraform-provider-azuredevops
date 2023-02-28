//go:build (all || core || workitem) && !exclude_resource_workitem
// +build all core workitem
// +build !exclude_resource_workitem

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkitem_Create(t *testing.T) {
	workitemTitle := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitem.workitem"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclWorkitemResource(projectName, workitemTitle),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "title", workitemTitle),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttr(tfNode, "type", "Issue"),
					resource.TestCheckResourceAttr(tfNode, "state", "Active"),
					resource.TestCheckResourceAttr(tfNode, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(tfNode, "tags.1", "tag2=value"),
				),
			},
		},
	})
}

func TestAccWorkitem_Update(t *testing.T) {
	workitemTitle := testutils.GenerateResourceName()
	workitemTitleUpdated := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitem.workitem"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclWorkitemResource(projectName, workitemTitle),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "title", workitemTitle),
				),
			},
			{
				Config: testutils.HclWorkitemResource(projectName, workitemTitleUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "title", workitemTitleUpdated),
				),
			},
		},
	})
}

func TestAccWorkitem_Import(t *testing.T) {

	tfNode := "azuredevops_workitem.workitem"
	workitemTitle := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclWorkitemResource(projectName, workitemTitle),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
				),
			},
			{
				// Resource Acceptance Testing https://www.terraform.io/docs/extend/resources/import.html#resource-acceptance-testing-implementation
				ResourceName:      tfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
