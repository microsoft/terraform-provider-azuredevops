package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccCheckRestAPI_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceConnectionName := testutils.GenerateResourceName()
	displayName := testutils.GenerateResourceName()

	tfCheckNode := "azuredevops_check_rest_api.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckPipelineCheckDestroyed("azuredevops_check_rest_api"),
		Steps: []resource.TestStep{
			{
				Config: hclCheckRestAPIResourceBasic(projectName, serviceConnectionName, displayName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckPipelineCheckExistsWithName(tfCheckNode, displayName),
					resource.TestCheckResourceAttr(tfCheckNode, "connected_service_name_selector", "connectedServiceName"),
					resource.TestCheckResourceAttr(tfCheckNode, "connected_service_name", "se_"+serviceConnectionName),
					resource.TestCheckResourceAttr(tfCheckNode, "method", "GET"),
				),
			},
		},
	})
}

func TestAccCheckRestAPI_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	variableGroupName := testutils.GenerateResourceName()
	displayName := testutils.GenerateResourceName()
	serviceConnectionName := testutils.GenerateResourceName()

	tfCheckNode := "azuredevops_check_rest_api.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckPipelineCheckDestroyed("azuredevops_check_rest_api"),
		Steps: []resource.TestStep{
			{
				Config: hclCheckRestAPIResourceComplete(projectName, serviceConnectionName, displayName, variableGroupName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckPipelineCheckExistsWithName(tfCheckNode, displayName),
					resource.TestCheckResourceAttr(tfCheckNode, "connected_service_name_selector", "connectedServiceName"),
					resource.TestCheckResourceAttr(tfCheckNode, "connected_service_name", "se_"+serviceConnectionName),
					resource.TestCheckResourceAttr(tfCheckNode, "method", "POST"),
					resource.TestCheckResourceAttr(tfCheckNode, "headers", "{\"contentType\":\"application/json\"}"),
					resource.TestCheckResourceAttr(tfCheckNode, "body", "{\"params\":\"value\"}"),
					resource.TestCheckResourceAttr(tfCheckNode, "completion_event", "ApiResponse"),
					resource.TestCheckResourceAttr(tfCheckNode, "success_criteria", "eq(root['status'], '200')"),
					resource.TestCheckResourceAttr(tfCheckNode, "url_suffix", "user/1"),
					resource.TestCheckResourceAttr(tfCheckNode, "retry_interval", "4000"),
					resource.TestCheckResourceAttr(tfCheckNode, "variable_group_name", variableGroupName),
					resource.TestCheckResourceAttr(tfCheckNode, "timeout", "40000"),
				),
			},
		},
	})
}

func TestAccCheckRestAPI_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	variableGroupName := testutils.GenerateResourceName()
	displayName := testutils.GenerateResourceName()
	serviceConnectionName := testutils.GenerateResourceName()

	tfCheckNode := "azuredevops_check_rest_api.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckPipelineCheckDestroyed("azuredevops_check_rest_api"),
		Steps: []resource.TestStep{
			{
				Config: hclCheckRestAPIResourceBasic(projectName, serviceConnectionName, displayName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckPipelineCheckExistsWithName(tfCheckNode, displayName),
					resource.TestCheckResourceAttr(tfCheckNode, "connected_service_name_selector", "connectedServiceName"),
					resource.TestCheckResourceAttr(tfCheckNode, "connected_service_name", "se_"+serviceConnectionName),
					resource.TestCheckResourceAttr(tfCheckNode, "method", "GET"),
				),
			},
			{
				Config: hclCheckRestAPIResourceComplete(projectName, serviceConnectionName, displayName, variableGroupName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckPipelineCheckExistsWithName(tfCheckNode, displayName),
					resource.TestCheckResourceAttr(tfCheckNode, "connected_service_name_selector", "connectedServiceName"),
					resource.TestCheckResourceAttr(tfCheckNode, "connected_service_name", "se_"+serviceConnectionName),
					resource.TestCheckResourceAttr(tfCheckNode, "method", "POST"),
					resource.TestCheckResourceAttr(tfCheckNode, "headers", "{\"contentType\":\"application/json\"}"),
					resource.TestCheckResourceAttr(tfCheckNode, "body", "{\"params\":\"value\"}"),
					resource.TestCheckResourceAttr(tfCheckNode, "completion_event", "ApiResponse"),
					resource.TestCheckResourceAttr(tfCheckNode, "success_criteria", "eq(root['status'], '200')"),
					resource.TestCheckResourceAttr(tfCheckNode, "url_suffix", "user/1"),
					resource.TestCheckResourceAttr(tfCheckNode, "retry_interval", "4000"),
					resource.TestCheckResourceAttr(tfCheckNode, "variable_group_name", variableGroupName),
					resource.TestCheckResourceAttr(tfCheckNode, "timeout", "40000"),
				),
			},
		},
	})
}

func hclCheckRestAPIResourceTemplate(projectName string, serviceConnectionName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_generic" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  server_url            = "https://dev.azure.com/"
  username              = "username"
  password              = "dummy"
}`, projectName, serviceConnectionName)
}

func hclCheckRestAPIResourceBasic(projectName, serviceConnectionName, displayName string) string {
	template := hclCheckRestAPIResourceTemplate(projectName, serviceConnectionName)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic" "test2" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "se_%s"
  server_url            = "https://dev.azure.com/"
  username              = "username"
  password              = "dummy"
}

resource "azuredevops_check_rest_api" "test" {
  project_id                      = azuredevops_project.test.id
  target_resource_id              = azuredevops_serviceendpoint_generic.test.id
  target_resource_type            = "endpoint"
  display_name                    = "%s"
  connected_service_name_selector = "connectedServiceName"
  connected_service_name          = azuredevops_serviceendpoint_generic.test2.service_endpoint_name
  method                          = "GET"
}`, template, serviceConnectionName, displayName)
}

func hclCheckRestAPIResourceComplete(projectName, serviceConnectionName, displayName, variableGroupName string) string {
	template := hclCheckRestAPIResourceTemplate(projectName, serviceConnectionName)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic" "test2" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "se_%s"
  server_url            = "https://dev.azure.com/"
  username              = "username"
  password              = "dummy"
}

resource "azuredevops_variable_group" "test" {
  project_id   = azuredevops_project.test.id
  name         = "%s"
  allow_access = true
  variable {
    name  = "FOO"
    value = "BAR"
  }
}

resource "azuredevops_check_rest_api" "test" {
  project_id           = azuredevops_project.test.id
  target_resource_id   = azuredevops_serviceendpoint_generic.test.id
  target_resource_type = "endpoint"

  display_name                    = "%s"
  connected_service_name_selector = "connectedServiceName"
  connected_service_name          = azuredevops_serviceendpoint_generic.test2.service_endpoint_name
  method                          = "POST"
  headers                         = "{\"contentType\":\"application/json\"}"
  body                            = "{\"params\":\"value\"}"
  completion_event                = "ApiResponse"
  success_criteria                = "eq(root['status'], '200')"
  url_suffix                      = "user/1"
  retry_interval                  = 4000
  variable_group_name             = azuredevops_variable_group.test.name
  timeout                         = "40000"
}`, template, serviceConnectionName, variableGroupName, displayName)
}
