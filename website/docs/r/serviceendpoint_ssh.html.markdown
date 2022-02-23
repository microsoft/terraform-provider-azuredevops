---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_ssh"
description: |- 
  Manages a SSH service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_ssh

Manages a SSH service endpoint within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_ssh" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example SSH"
  host                  = "1.2.3.4"
  username              = "username"
  description           = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The project ID or project name.
- `service_endpoint_name` - (Required) The Service Endpoint name.
- `host` - (Required) The Host name or IP address of the remote machine.
- `username` - (Required) Username for connecting to the endpoint.
- `port` - (Optional) Port number on the remote machine to use for connecting. Defaults to `22`.
- `password` - (Optional) Password for connecting to the endpoint.
- `private_key` - (Optional) Private Key for connecting to the endpoint.
- `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The project ID or project name.
- `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Service Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-6.0)

## Import

Azure DevOps Service Endpoint SSH can be imported using **projectID/serviceEndpointID** or **
projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_ssh.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
