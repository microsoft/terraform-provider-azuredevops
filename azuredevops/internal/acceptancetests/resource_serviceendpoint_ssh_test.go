// +build all resource_serviceendpoint_ssh
// +build !resource_serviceendpoint_ssh

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointSSH_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_ssh"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointSSHResourceBasic(projectName, serviceEndpointName, "1.2.3.4", "username"),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
				),
			},
		},
	})
}

func TestAccServiceEndpointSSH_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	description := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_ssh"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointSSHResourceComplete(projectName, serviceEndpointName, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "private_key"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "password"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "host", "1.2.3.4"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "port", "22"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "username", "username"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
				),
			},
		},
	})
}

func TestAccServiceEndpointSSH_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()

	description := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_ssh"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointSSHResourceBasic(projectName, serviceEndpointNameFirst, "1.2.3.4", "username"),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst), resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
				),
			},
			{
				Config: hclSvcEndpointSSHResourceUpdate(projectName, serviceEndpointNameSecond, "2.2.3.4", 23, "newname", description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "host", "2.2.3.4"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "port", "23"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "username", "newname"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
				),
			},
		},
	})
}

func TestAccServiceEndpointSSH_RequiresImportErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	resourceType := "azuredevops_serviceendpoint_ssh"
	tfSvcEpNode := resourceType + ".test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointSSHResourceBasic(projectName, serviceEndpointName, "1.2.3.4", "username"),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				Config:      hclSvcEndpointSSHResourceRequiresImport(projectName, serviceEndpointName, "1.2.3.4", "username"),
				ExpectError: testutils.RequiresImportError(serviceEndpointName),
			},
		},
	})
}

func hclSvcEndpointSSHResourceBasic(projectName string, serviceEndpointName string, host string, username string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
  name               = "%[1]s"
  description        = "%[1]s-description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_ssh" "test" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%[2]s"
  host                  = "%[3]s"
  username              = "%[4]s"
}
`, projectName, serviceEndpointName, host, username)
}

func hclSvcEndpointSSHResourceComplete(projectName string, serviceEndpointName string, description string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
  name               = "%[1]s"
  description        = "%[1]s-description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_ssh" "test" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%[2]s"
  host                  = "1.2.3.4"
  port                  = 22
  private_key           = "privateKey"
  username              = "username"
  password              = "password"
  description           = "%[3]s"
}
`, projectName, serviceEndpointName, description)
}

func hclSvcEndpointSSHResourceUpdate(projectName string, serviceEndpointName string, host string, port int, username string, description string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
  name               = "%[1]s"
  description        = "%[1]s-description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_ssh" "test" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%[2]s"
  host                  = "%[3]s"
  port                  = "%[4]d"
  private_key           = "privateKey"
  username              = "%[5]s"
  password              = "password"
  description           = "%[6]s"
}
`, projectName, serviceEndpointName, host, port, username, description)
}

func hclSvcEndpointSSHResourceRequiresImport(projectName string, serviceEndpointName string, host string, username string) string {
	template := hclSvcEndpointSSHResourceBasic(projectName, serviceEndpointName, host, username)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_ssh" "import" {
  project_id            = azuredevops_serviceendpoint_ssh.test.project_id
  service_endpoint_name = azuredevops_serviceendpoint_ssh.test.service_endpoint_name
  description           = azuredevops_serviceendpoint_ssh.test.description
  host                  = azuredevops_serviceendpoint_ssh.test.host
  username              = azuredevops_serviceendpoint_ssh.test.username
}
`, template)
}
