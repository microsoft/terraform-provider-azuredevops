package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccSecureFileDataSource_Basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	secureFileName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_secure_file.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclSecureFileDataSourceBasic(projectName, secureFileName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", secureFileName),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "file_hash_sha1"),
					resource.TestCheckResourceAttrSet(tfNode, "file_hash_sha256"),
				),
			},
		},
	})
}

func TestAccSecureFileDataSource_WithProperties(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	secureFileName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_secure_file.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclSecureFileDataSourceWithProperties(projectName, secureFileName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", secureFileName),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "file_hash_sha1"),
					resource.TestCheckResourceAttrSet(tfNode, "file_hash_sha256"),
					resource.TestCheckResourceAttr(tfNode, "properties.custom_key", "custom_value"),
					resource.TestCheckResourceAttr(tfNode, "properties.environment", "production"),
				),
			},
		},
	})
}

func hclSecureFileDataSourceBasic(projectName, secureFileName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
	name               = "%[1]s"
	description        = "%[1]s-description"
	visibility         = "private"
	version_control    = "Git"
	work_item_template = "Agile"
}

resource "azuredevops_secure_file" "test" {
	project_id = azuredevops_project.project.id
	name       = "%[2]s"
	content    = "test file content"
}

data "azuredevops_secure_file" "test" {
	project_id = azuredevops_project.project.id
	name       = azuredevops_secure_file.test.name
}
`, projectName, secureFileName)
}

func hclSecureFileDataSourceWithProperties(projectName, secureFileName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
	name               = "%[1]s"
	description        = "%[1]s-description"
	visibility         = "private"
	version_control    = "Git"
	work_item_template = "Agile"
}

resource "azuredevops_secure_file" "test" {
	project_id = azuredevops_project.project.id
	name       = "%[2]s"
	content    = "test file content with properties"
	properties = {
		custom_key  = "custom_value"
		environment = "production"
	}
}

data "azuredevops_secure_file" "test" {
	project_id = azuredevops_project.project.id
	name       = azuredevops_secure_file.test.name
}
`, projectName, secureFileName)
}
