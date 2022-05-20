//go:build (all || resource_build_folder) && (!exclude_permissions || !exclude_resource_build_folder)
// +build all resource_build_folder
// +build !exclude_permissions !exclude_resource_build_folder

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestBuildFolder(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := testutils.HclBuildFolder(projectName, "\\test-folder", "Acceptance Test Folder")

	tfNode := "azuredevops_build_folder.test_folder"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "path"),
					resource.TestCheckResourceAttrSet(tfNode, "description"),
				),
			},
		},
	})
}
