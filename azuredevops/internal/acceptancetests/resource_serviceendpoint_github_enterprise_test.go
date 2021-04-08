// +build all resource_serviceendpoint_github_enterprise
// +build !exclude_serviceendpoints

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointGitHubEnterprise_PersonalTokenBasic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_github_enterprise"
	tfSvcEpNode := resourceType + ".serviceendpoint"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, &[]string{"AZDO_GITHUB_ENTERPRISE_SERVICE_CONNECTION_PAT"}) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: personTokenConfigBasicGithubEnterprise(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_personal.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://github.contoso.com"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
		},
	})
}

func TestAccServiceEndpointGitHubEnterprise_PersonalTokenUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()
	description := "Manage by Terraform Update"

	resourceType := "azuredevops_serviceendpoint_github_enterprise"
	tfSvcEpNode := resourceType + ".serviceendpoint"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, &[]string{"AZDO_GITHUB_ENTERPRISE_SERVICE_CONNECTION_PAT"}) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: personTokenConfigBasicGithubEnterprise(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_personal.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://github.contoso.com"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
				),
			}, {
				Config: personTokenConfigUpdateGithubEnterprise(projectName, serviceEndpointNameSecond, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_personal.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://github.contoso.com"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
				),
			},
		},
	})
}

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccServiceEndpointGitHubEnterprise_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_github_enterprise"
	tfSvcEpNode := resourceType + ".serviceendpoint"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, &[]string{"AZDO_GITHUB_ENTERPRISE_SERVICE_CONNECTION_PAT"}) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointGitHubEnterpriseResource(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_personal.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://github.contoso.com"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
				),
			}, {
				Config: testutils.HclServiceEndpointGitHubEnterpriseResource(projectName, serviceEndpointNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_personal.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://github.contoso.com"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
				),
			}, {
				// Resource Acceptance Testing https://www.terraform.io/docs/extend/resources/import.html#resource-acceptance-testing-implementation
				ResourceName:            tfSvcEpNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_personal"},
			},
		},
	})
}

func personTokenConfigBasicGithubEnterprise(projectName string, serviceEndpointName string) string {
	projectResource := testutils.HclProjectResource(projectName)

	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_github_enterprise" "serviceendpoint" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%[1]s"
	url                    = "https://github.contoso.com"
	auth_personal {
		personal_access_token= "test_token_basic"
	}
}`, serviceEndpointName)

	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func personTokenConfigUpdateGithubEnterprise(projectName string, serviceEndpointName string, description string) string {
	projectResource := testutils.HclProjectResource(projectName)

	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_github_enterprise" "serviceendpoint" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%[1]s"
	url                    = "https://github.contoso.com"
	auth_personal {
		personal_access_token= "test_token_update"
	}
	description = "%[2]s"
}`, serviceEndpointName, description)

	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}
