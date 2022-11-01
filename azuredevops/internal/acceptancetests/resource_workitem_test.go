//go:build (all || core || workitem) && !exclude_resource_workitem
// +build all core workitem
// +build !exclude_resource_workitem

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkitem_CreateAndUpdate(t *testing.T) {
	workitemTitle := testutils.GenerateResourceName()
	projectNameFirst := testutils.GenerateResourceName()
	projectNameSecond := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitem.workitem"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclWorkitemResource(projectNameFirst, workitemTitle),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectNameFirst),
					resource.TestCheckResourceAttr(tfNode, "title", workitemTitle),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttr(tfNode, "type", "Issue"),
					resource.TestCheckResourceAttr(tfNode, "state", "Active"),
				),
			},
			{
				Config: testutils.HclWorkitemResource(projectNameSecond, workitemTitle),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectNameSecond),
					resource.TestCheckResourceAttr(tfNode, "title", workitemTitle),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttr(tfNode, "type", "Issue"),
					resource.TestCheckResourceAttr(tfNode, "state", "Active"),
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
