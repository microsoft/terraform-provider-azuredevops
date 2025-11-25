package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointProjectPermissions_Update(t *testing.T) {
	projectName1 := testutils.GenerateResourceName()
	projectName2 := testutils.GenerateResourceName()
	projectName3 := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	tfNode := "azuredevops_serviceendpoint_project_permissions.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			// Stap 1: Initieel delen met Project 2
			{
				Config: hclServiceEndpointPermissionsBuilder(projectName1, projectName2, projectName3, serviceEndpointName, `
				project_reference {
					project_id            = azuredevops_project.p2.id
					service_endpoint_name = "shared-connection"
					description           = "Initial share"
				}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "project_reference.#", "1"),
					resource.TestCheckResourceAttr(tfNode, "project_reference.0.service_endpoint_name", "shared-connection"),
				),
			},
			// Stap 2: Update - Project 3 toevoegen (Upsert logica)
			{
				Config: hclServiceEndpointPermissionsBuilder(projectName1, projectName2, projectName3, serviceEndpointName, `
				project_reference {
					project_id            = azuredevops_project.p2.id
					service_endpoint_name = "shared-connection"
					description           = "Initial share"
				}
				project_reference {
					project_id            = azuredevops_project.p3.id
					service_endpoint_name = "shared-connection-p3"
					description           = "Added via update"
				}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "project_reference.#", "2"),
				),
			},
			// Stap 3: Update - Project 2 verwijderen EN alias van Project 3 wijzigen
			{
				Config: hclServiceEndpointPermissionsBuilder(projectName1, projectName2, projectName3, serviceEndpointName, `
				project_reference {
					project_id            = azuredevops_project.p3.id
					service_endpoint_name = "renamed-connection-p3"
					description           = "Updated description"
				}
				`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "project_reference.#", "1"),
					resource.TestCheckResourceAttr(tfNode, "project_reference.0.service_endpoint_name", "renamed-connection-p3"),
					resource.TestCheckResourceAttr(tfNode, "project_reference.0.description", "Updated description"),
				),
			},
		},
	})
}

func hclServiceEndpointPermissionsBuilder(p1, p2, p3, seName, permissionsBlock string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "p1" {
  name = "%[1]s"
}

resource "azuredevops_project" "p2" {
  name = "%[2]s"
}

resource "azuredevops_project" "p3" {
  name = "%[3]s"
}

resource "azuredevops_serviceendpoint_azuredevops" "example" {
  project_id            = azuredevops_project.p1.id
  service_endpoint_name = "%[4]s"
  org_url               = "https://dev.azure.com/testorganization"
  release_api_url       = "https://vsrm.dev.azure.com/testorganization"
  personal_access_token = "0000000000000000000000000000000000000000000000000000"
}

resource "azuredevops_serviceendpoint_project_permissions" "test" {
  project_id         = azuredevops_project.p1.id
  service_endpoint_id = azuredevops_serviceendpoint_azuredevops.example.id

  %[5]s
}
`, p1, p2, p3, seName, permissionsBlock)
}
