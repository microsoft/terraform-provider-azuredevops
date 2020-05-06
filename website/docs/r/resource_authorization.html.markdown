# azuredevops_resource_authorization
Manages authorization of resources, e.g. for access in build pipelines.

Currently supported resources: service endpoint (aka service connection, endpoint).

## Example Usage

```hcl
resource "azuredevops_serviceendpoint_bitbucket" "bitbucket_account" {
  project_id = "vanilla-sky"
  username               = "xxxx"
  password               = "xxxx"
  service_endpoint_name  = "test-bitbucket"
  description            = "test"
}

resource "azuredevops_resource_authorization" "auth" {
  project_id = "phoenix-project"
  resource_id = azuredevops_serviceendpoint_bitbucket.bitbucket_account.id
  authorized = true
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The project ID or project name. Type: string.
* `service_endpoint_name` - (Required) The Service Endpoint name. Type: string.
* `authorized` - (Required) Set to true to allow public access in the project. Type: boolean.
* `type` - (Optional) The type of the resource to authorize. Type: string. Valid values: `endpoint`. Default value: `endpoint`.

## Attributes Reference

The following attributes are exported: 

n/a

## Relevant Links
* [Azure DevOps Service REST API 5.1 - Authorize Definition Resource](https://docs.microsoft.com/en-us/rest/api/azure/devops/build/resources/authorize%20definition%20resources?view=azure-devops-rest-5.1)
