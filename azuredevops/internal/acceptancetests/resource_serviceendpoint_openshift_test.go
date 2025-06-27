//go:build (all || resource_serviceendpoint_openshift) && !exclude_resource_serviceendpoint_openshift

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointOpenshift_authorizationBasic_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_openshift"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceConnectionOpenshiftAuthorizationBasic(projectName, serviceEndpointName, "https://ado.test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_basic.0.username", "username"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_basic.0.password", "password"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://ado.test"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_basic.#", "auth_basic.0.%", "auth_basic.0.username", "auth_basic.0.password"},
			},
		},
	})
}

func TestAccServiceEndpointOpenshift_authorizationBasic_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_openshift"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceConnectionOpenshiftAuthorizationBasicComplete(projectName, serviceEndpointName, "https://ado.test", "/opt/tmp/file", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_basic.0.username", "username2"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_basic.0.password", "password2"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://ado.test"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_basic.#", "auth_basic.0.%", "auth_basic.0.username", "auth_basic.0.password"},
			},
		},
	})
}

func TestAccServiceEndpointOpenshift_authorizationBasic_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_openshift"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceConnectionOpenshiftAuthorizationBasic(projectName, serviceEndpointName, "https://ado.test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://ado.test"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_basic.#", "auth_basic.0.%", "auth_basic.0.username", "auth_basic.0.password"},
			},
			{
				Config: hclServiceConnectionOpenshiftAuthorizationBasicComplete(projectName, serviceEndpointName, "https://ado.test", "/opt/tmp/file", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_basic.0.username", "username2"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_basic.0.password", "password2"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://ado.test"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "certificate_authority_file", "/opt/tmp/file"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "accept_untrusted_certs", "true"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_basic.#", "auth_basic.0.%", "auth_basic.0.username", "auth_basic.0.password"},
			},
			{
				Config: hclServiceConnectionOpenshiftAuthorizationBasic(projectName, serviceEndpointName, "https://ado.test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_basic.0.username", "username"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_basic.0.password", "password"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://ado.test"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_basic.#", "auth_basic.0.%", "auth_basic.0.username", "auth_basic.0.password"},
			},
		},
	})
}

func TestAccServiceEndpointOpenshift_authorizationToken_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_openshift"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceConnectionOpenshiftAuthorizationTokenBasic(projectName, serviceEndpointName, "https://ado.test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_token.0.token", "token"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://ado.test"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_token.#", "auth_token.0.%", "auth_token.0.token"},
			},
		},
	})
}

func TestAccServiceEndpointOpenshift_authorizationToken_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_openshift"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceConnectionOpenshiftAuthorizationTokenComplete(projectName, serviceEndpointName, "https://ado.test", "/opt/tmp/file", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_token.0.token", "token2"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://ado.test"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "certificate_authority_file", "/opt/tmp/file"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "accept_untrusted_certs", "true"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_token.#", "auth_token.0.%", "auth_token.0.token"},
			},
		},
	})
}

func TestAccServiceEndpointOpenshift_authorizationToken_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_openshift"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceConnectionOpenshiftAuthorizationTokenBasic(projectName, serviceEndpointName, "https://ado.test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_token.0.token", "token"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://ado.test"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_token.#", "auth_token.0.%", "auth_token.0.token"},
			},
			{
				Config: hclServiceConnectionOpenshiftAuthorizationTokenComplete(projectName, serviceEndpointName, "https://ado.test", "/opt/tmp/file", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_token.0.token", "token2"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://ado.test"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "certificate_authority_file", "/opt/tmp/file"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "accept_untrusted_certs", "true"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_token.#", "auth_token.0.%", "auth_token.0.token"},
			},
			{
				Config: hclServiceConnectionOpenshiftAuthorizationTokenBasic(projectName, serviceEndpointName, "https://ado.test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_token.0.token", "token"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://ado.test"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_token.#", "auth_token.0.%", "auth_token.0.token"},
			},
		},
	})
}

func TestAccServiceEndpointOpenshift_authorizationNone_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_openshift"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceConnectionOpenshiftAuthorizationNoneBasic(projectName, serviceEndpointName, "https://ado.test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_none.0.%", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_none.0.kube_config", "config"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://ado.test"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_none.#", "auth_none.0.%", "auth_none.0.kube_config"},
			},
		},
	})
}

func TestAccServiceEndpointOpenshift_authorizationNone_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_openshift"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceConnectionOpenshiftAuthorizationNoneBasic(projectName, serviceEndpointName, "https://ado.test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_none.0.%", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_none.0.kube_config", "config"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://ado.test"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_none.#", "auth_none.0.%", "auth_none.0.kube_config"},
			},
			{
				Config: hclServiceConnectionOpenshiftAuthorizationNoneUpdate(projectName, serviceEndpointName, "https://ado.test2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_none.0.%", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_none.0.kube_config", "config2"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://ado.test2"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_none.#", "auth_none.0.%", "auth_none.0.kube_config"},
			},
			{
				Config: hclServiceConnectionOpenshiftAuthorizationNoneBasic(projectName, serviceEndpointName, "https://ado.test2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_none.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "server_url", "https://ado.test2"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				ResourceName:            tfSvcEpNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_none.#", "auth_none.0.%", "auth_none.0.kube_config"},
			},
		},
	})
}

func TestAccServiceEndpointOpenshift_requireImportError(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_openshift"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceConnectionOpenshiftAuthorizationBasic(projectName, serviceEndpointName, "https://ado.test"),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				Config:      hclServiceConnectionOpenshiftRequireImport(projectName, serviceEndpointName, "https://ado.test"),
				ExpectError: testutils.RequiresImportError(serviceEndpointName),
			},
		},
	})
}

func hclServiceConnectionOpenshiftAuthorizationBasic(projectName, seName, url string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_openshift" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  server_url            = "%s"
  auth_basic {
    username = "username"
    password = "password"
  }
}`, projectName, seName, url)
}

func hclServiceConnectionOpenshiftAuthorizationBasicComplete(projectName, seName, url, authFile string, untrustedCert bool) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_openshift" "test" {
  project_id                 = azuredevops_project.test.id
  service_endpoint_name      = "%s"
  server_url                 = "%s"
  certificate_authority_file = "%s"
  accept_untrusted_certs     = "%t"
  auth_basic {
    username = "username2"
    password = "password2"
  }
}`, projectName, seName, url, authFile, untrustedCert)
}

func hclServiceConnectionOpenshiftAuthorizationTokenBasic(projectName, seName, url string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_openshift" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  server_url            = "%s"
  auth_token {
    token = "token"
  }
}`, projectName, seName, url)
}

func hclServiceConnectionOpenshiftAuthorizationTokenComplete(projectName, seName, url, authFile string, untrustedCert bool) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_openshift" "test" {
  project_id                 = azuredevops_project.test.id
  service_endpoint_name      = "%s"
  server_url                 = "%s"
  certificate_authority_file = "%s"
  accept_untrusted_certs     = "%t"
  auth_token {
    token = "token2"
  }
}`, projectName, seName, url, authFile, untrustedCert)
}

func hclServiceConnectionOpenshiftAuthorizationNoneBasic(projectName, seName, url string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_openshift" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  server_url            = "%s"
  auth_none {
    kube_config = "config"
  }
}`, projectName, seName, url)
}

func hclServiceConnectionOpenshiftAuthorizationNoneUpdate(projectName, seName, url string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_openshift" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  server_url            = "%s"
  auth_none {
    kube_config = "config2"
  }
}`, projectName, seName, url)
}

func hclServiceConnectionOpenshiftRequireImport(projectName, seName, url string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_openshift" "import" {
  project_id            = azuredevops_serviceendpoint_openshift.test.project_id
  service_endpoint_name = azuredevops_serviceendpoint_openshift.test.service_endpoint_name
  server_url            = azuredevops_serviceendpoint_openshift.test.server_url
  auth_basic {
    username = tolist(azuredevops_serviceendpoint_openshift.test.auth_basic)[0].username
    password = tolist(azuredevops_serviceendpoint_openshift.test.auth_basic)[0].password
  }
}`, hclServiceConnectionOpenshiftAuthorizationBasic(projectName, seName, url))
}
