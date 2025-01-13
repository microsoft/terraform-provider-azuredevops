---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_check_rest_api"
description: |-
  Manages a Rest API check.
---

# azuredevops_check_rest_api

Manages a Rest API check on a resource within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_serviceendpoint_generic" "example" {
  project_id            = azuredevops_project.example.id
  server_url            = "https://some-server.example.com"
  service_endpoint_name = "Example Generic"
  username              = "username"
  password              = "password"
  description           = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_generic" "example_azure" {
  project_id            = azuredevops_project.example.id
  server_url            = "https://dev.azure.com/"
  service_endpoint_name = "Example Generic Azure"
  username              = "username"
  password              = "dummy"
}

resource "azuredevops_variable_group" "example" {
  project_id   = azuredevops_project.example.id
  name         = "Example Variable Group"
  allow_access = true

  variable {
    name  = "FOO"
    value = "BAR"
  }
}

resource "azuredevops_check_rest_api" "example" {
  project_id                      = azuredevops_project.example.id
  target_resource_id              = azuredevops_serviceendpoint_generic.example.id
  target_resource_type            = "endpoint"
  display_name                    = "Example REST API Check"
  connected_service_name_selector = "connectedServiceName"
  connected_service_name          = azuredevops_serviceendpoint_generic.example_azure.service_endpoint_name
  method                          = "POST"
  headers                         = "{\"contentType\":\"application/json\"}"
  body                            = "{\"params\":\"value\"}"
  completion_event                = "ApiResponse"
  success_criteria                = "eq(root['status'], '200')"
  url_suffix                      = "user/1"
  retry_interval                  = 4000
  variable_group_name             = azuredevops_variable_group.example.name
  timeout                         = "40000"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project. Changing this forces a new resource to be created.
 
* `target_resource_id` - (Required) The ID of the resource being protected by the check. Changing this forces a new resource to be created

* `target_resource_type` - (Required) The type of resource being protected by the check. Possible values: `endpoint`, `environment`, `queue`, `repository`, `securefile`, `variablegroup`. Changing this forces a new resource to be created.

* `connected_service_name_selector` - (Required) The type of the Service Connection used to invoke the REST API. Possible values: `connectedServiceName`(**Generic** type service connection) and `connectedServiceNameARM`(**Azure Resource Manager** type service connection).
  
* `connected_service_name` - (Required) The name of the Service Connection.

* `display_name` - (Required) The Name of the Rest API check.

* `method` - (Required) The HTTP method of the request. Possible values: `OPTIONS`, `GET`, `HEAD`, `POST`, `PUT`, `DELETE`, `TRACE`, `PATCH`

---
* `body` - (Optional) The Rest API request body.

* `headers` - (Optional) The headers of the request in JSON format.
  
* `retry_interval` - (Optional) The time between evaluations (minutes). 
  
    ~>**NOTE** 1) The retry times should less them 10 based on the timeout. For example: `timeout` is `4000` then `retry_interval` should be `0` or no less then `400`.
    <br>2) `retry_interval` is not required when `completion_event=Callback`.

* `success_criteria` - (Optional) The Criteria which defines when to pass the task. No criteria means response content does not influence the result.

  ~>**NOTE** `success_criteria` is used when `completion_event=ApiResponse`

* `url_suffix` - (Optional) The URL suffix and parameters.

* `variable_group_name` - (Optional) The name of the Variable Group.

* `completion_event` - (Optional) The completion event of the Rest API call. Possible values: `Callback`, `ApiResponse`. Defaults to `Callback`.

* `timeout` - (Optional) The timeout in minutes for the Rest API check. Defaults to `1440`.

## Attributes Reference

In addition to all arguments above the following attributes are exported:

* `id` - The ID of the check.
* `version` - The version of the Rest API check.

## Relevant Links

- [Define approvals and checks](https://learn.microsoft.com/en-us/azure/devops/pipelines/process/approvals?view=azure-devops&tabs=check-pass)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Rest API Check.
* `read` - (Defaults to 1 minute) Used when retrieving the Rest API Check.
* `update` - (Defaults to 2 minutes) Used when updating the Rest API Check.
* `delete` - (Defaults to 2 minutes) Used when deleting the Rest API Check.

## Import

Importing this resource is not supported.
