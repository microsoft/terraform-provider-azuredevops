//go:build (all || resource_serviceendpoint_jfrog_distribution_v2) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_jfrog_distribution_v2
// +build !exclude_serviceendpoints

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointJFrogDistributionV2_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_jfrog_distribution_v2"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointJFrogDistributionV2ResourceBasic(projectName, serviceEndpointName, t.Name()),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "url"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
				),
			},
		},
	})
}

func TestAccServiceEndpointJFrogDistributionV2_basic_usernamepassword(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_jfrog_distribution_v2"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointJFrogDistributionV2ResourceBasicUsernamePassword(projectName, serviceEndpointName, t.Name()),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authentication_basic.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
				),
			},
		},
	})
}

func TestAccServiceEndpointJFrogDistributionV2_complete_token(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	description := t.Name()

	resourceType := "azuredevops_serviceendpoint_jfrog_distribution_v2"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointJFrogDistributionV2ResourceComplete(projectName, serviceEndpointName, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authentication_token.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://url.com/1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
				),
			},
		},
	})
}

func TestAccServiceEndpointJFrogDistributionV2_complete_usernamepassword(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	description := t.Name()

	resourceType := "azuredevops_serviceendpoint_jfrog_distribution_v2"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointJFrogDistributionV2ResourceCompleteUsernamePassword(projectName, serviceEndpointName, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authentication_basic.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://url.com/1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
				),
			},
		},
	})
}

func TestAccServiceEndpointJFrogDistributionV2_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()

	description := t.Name()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_jfrog_distribution_v2"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointJFrogDistributionV2ResourceBasic(projectName, serviceEndpointNameFirst, t.Name()),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst), resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
				),
			},
			{
				Config: hclSvcEndpointJFrogDistributionV2ResourceUpdate(projectName, serviceEndpointNameSecond, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authentication_token.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://url.com/2"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
				),
			},
		},
	})
}

func TestAccServiceEndpointJFrogDistributionV2_update_usernamepassword(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()

	description := t.Name()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_jfrog_distribution_v2"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointJFrogDistributionV2ResourceBasicUsernamePassword(projectName, serviceEndpointNameFirst, t.Name()),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst), resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
				),
			},
			{
				Config: hclSvcEndpointJFrogDistributionV2ResourceUpdateUsernamePassword(projectName, serviceEndpointNameSecond, description),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authentication_basic.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "url", "https://url.com/2"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
				),
			},
		},
	})
}

func TestAccServiceEndpointJFrogDistributionV2_RequiresImportErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	resourceType := "azuredevops_serviceendpoint_jfrog_distribution_v2"
	tfSvcEpNode := resourceType + ".test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointJFrogDistributionV2ResourceBasic(projectName, serviceEndpointName, t.Name()),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				Config:      hclSvcEndpointJFrogDistributionV2ResourceRequiresImport(projectName, serviceEndpointName, t.Name()),
				ExpectError: testutils.RequiresImportError(serviceEndpointName),
			},
		},
	})
}

func TestAccServiceEndpointJFrogDistributionV2_RequiresImportErrorStepUsernamePassword(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	resourceType := "azuredevops_serviceendpoint_jfrog_distribution_v2"
	tfSvcEpNode := resourceType + ".test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcEndpointJFrogDistributionV2ResourceBasicUsernamePassword(projectName, serviceEndpointName, t.Name()),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
			{
				Config:      hclSvcEndpointJFrogDistributionV2ResourceRequiresImport(projectName, serviceEndpointName, t.Name()),
				ExpectError: testutils.RequiresImportError(serviceEndpointName),
			},
		},
	})
}

func hclSvcEndpointJFrogDistributionV2ResourceBasic(projectName string, serviceEndpointName string, description string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_jfrog_distribution_v2" "test" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	authentication_token {
		token			   	   = "redacted"
	}
	url			   		   = "http://url.com/1"
	description 		   = "%s"
}`, serviceEndpointName, description)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func hclSvcEndpointJFrogDistributionV2ResourceBasicUsernamePassword(projectName string, serviceEndpointName string, description string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_jfrog_distribution_v2" "test" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	authentication_basic {
		username			   = "u"
		password			   = "redacted"
	}
	url			   		   = "http://url.com/1"
	description 		   = "%s"
}`, serviceEndpointName, description)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func hclSvcEndpointJFrogDistributionV2ResourceCompleteUsernamePassword(projectName string, serviceEndpointName string, description string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_jfrog_distribution_v2" "test" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	description            = "%s"
	authentication_basic {
		username			   = "u"
		password			   = "redacted"
	}
	url			   		   = "https://url.com/1"
}`, serviceEndpointName, description)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func hclSvcEndpointJFrogDistributionV2ResourceComplete(projectName string, serviceEndpointName string, description string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_jfrog_distribution_v2" "test" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	description            = "%s"
	authentication_token {
		token          = "redacted"
	}
	  url			   		   = "https://url.com/1"
}`, serviceEndpointName, description)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func hclSvcEndpointJFrogDistributionV2ResourceUpdate(projectName string, serviceEndpointName string, description string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_jfrog_distribution_v2" "test" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	description            = "%s"
	authentication_token {
		token          = "redacted2"
	}
	  url			   		   = "https://url.com/2"
}`, serviceEndpointName, description)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func hclSvcEndpointJFrogDistributionV2ResourceUpdateUsernamePassword(projectName string, serviceEndpointName string, description string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_jfrog_distribution_v2" "test" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	description            = "%s"
	authentication_basic {
		username			   = "u2"
		password			   = "redacted2"
	}
	url			   		   = "https://url.com/2"
}`, serviceEndpointName, description)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func hclSvcEndpointJFrogDistributionV2ResourceRequiresImport(projectName string, serviceEndpointName string, description string) string {
	template := hclSvcEndpointJFrogDistributionV2ResourceBasic(projectName, serviceEndpointName, description)
	return fmt.Sprintf(`
%s
resource "azuredevops_serviceendpoint_jfrog_distribution_v2" "import" {
  project_id                = azuredevops_serviceendpoint_jfrog_distribution_v2.test.project_id
  service_endpoint_name = azuredevops_serviceendpoint_jfrog_distribution_v2.test.service_endpoint_name
  description            = azuredevops_serviceendpoint_jfrog_distribution_v2.test.description
  url          = azuredevops_serviceendpoint_jfrog_distribution_v2.test.url
  authentication_token {
	  token          = "redacted"
  }
}
`, template)
}
func hclSvcEndpointJFrogDistributionV2ResourceRequiresImportUsernamePassword(projectName string, serviceEndpointName string, description string) string {
	template := hclSvcEndpointJFrogDistributionV2ResourceBasicUsernamePassword(projectName, serviceEndpointName, description)
	return fmt.Sprintf(`
%s
resource "azuredevops_serviceendpoint_jfrog_distribution_v2" "import" {
  project_id                = azuredevops_serviceendpoint_jfrog_distribution_v2.test.project_id
  service_endpoint_name = azuredevops_serviceendpoint_jfrog_distribution_v2.test.service_endpoint_name
  description            = azuredevops_serviceendpoint_jfrog_distribution_v2.test.description
  url          	= azuredevops_serviceendpoint_jfrog_distribution_v2.test.url
  authentication_basic {
	username			   = "u"
	password			   = "redacted"
  }
}
`, template)
}
