//go:build (all || core || resource_project) && !exclude_resource_project
// +build all core resource_project
// +build !exclude_resource_project

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// Verifies that the following sequence of events occurs without error:
//
//		(1) TF apply creates project
//		(2) TF state values are set
//		(3) project can be queried by ID and has expected name
//	 (4) TF apply update project with changing name
//	 (5) project can be queried by ID and has expected name
//		(6) TF destroy deletes project
//		(7) project can no longer be queried by ID
func TestAccProject_CreateAndUpdate(t *testing.T) {
	projectNameFirst := testutils.GenerateResourceName()
	projectNameSecond := testutils.GenerateResourceName()
	tfNode := "azuredevops_project.project"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclProjectResource(projectNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_template_id"),
					resource.TestCheckResourceAttr(tfNode, "name", projectNameFirst),
					resource.TestCheckResourceAttr(tfNode, "version_control", "Git"),
					resource.TestCheckResourceAttr(tfNode, "visibility", "private"),
					resource.TestCheckResourceAttr(tfNode, "work_item_template", "Agile"),
					testutils.CheckProjectExists(projectNameFirst),
				),
			},
			{
				Config: testutils.HclProjectResource(projectNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_template_id"),
					resource.TestCheckResourceAttr(tfNode, "name", projectNameSecond),
					resource.TestCheckResourceAttr(tfNode, "version_control", "Git"),
					resource.TestCheckResourceAttr(tfNode, "visibility", "private"),
					resource.TestCheckResourceAttr(tfNode, "work_item_template", "Agile"),
					testutils.CheckProjectExists(projectNameSecond),
				),
			},
			{
				// Resource Acceptance Testing https://www.terraform.io/docs/extend/resources/import.html#resource-acceptance-testing-implementation
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProject_CreateAndUpdateWithFeatures(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_project.project"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclProjectResourceWithFeature(projectName, "disabled", "disabled"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_template_id"),
					resource.TestCheckResourceAttr(tfNode, "name", projectName),
					resource.TestCheckResourceAttr(tfNode, "version_control", "Git"),
					resource.TestCheckResourceAttr(tfNode, "visibility", "private"),
					resource.TestCheckResourceAttr(tfNode, "work_item_template", "Agile"),
					resource.TestCheckResourceAttr(tfNode, "features.testplans", "disabled"),
					resource.TestCheckResourceAttr(tfNode, "features.artifacts", "disabled"),
					testutils.CheckProjectExists(projectName),
				),
			},
			{
				Config: testutils.HclProjectResourceWithFeature(projectName, "enabled", "disabled"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_template_id"),
					resource.TestCheckResourceAttr(tfNode, "name", projectName),
					resource.TestCheckResourceAttr(tfNode, "version_control", "Git"),
					resource.TestCheckResourceAttr(tfNode, "visibility", "private"),
					resource.TestCheckResourceAttr(tfNode, "work_item_template", "Agile"),
					resource.TestCheckResourceAttr(tfNode, "features.testplans", "enabled"),
					resource.TestCheckResourceAttr(tfNode, "features.artifacts", "disabled"),
					testutils.CheckProjectExists(projectName),
				),
			},
		},
	})
}
