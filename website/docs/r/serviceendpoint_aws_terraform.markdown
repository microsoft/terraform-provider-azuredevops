---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_aws_terraform"
description: |-
  Manages an Azure "AWS for Terraform" service endpoint within an Azure DevOps organization.
---

# azuredevops_serviceendpoint_aws_terraform

Manages a "AWS service endpoint" within Azure DevOps. Using this service endpoint requires you to first install the [Azure Terraform Extension for Azure DevOps from Microsoft DevLabs](https://marketplace.visualstudio.com/items?itemName=ms-devlabs.custom-terraform-tasks\).

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_aws_terraform" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example AWS"
  access_key_id         = "00000000-0000-0000-0000-000000000000"
  secret_access_key     = "accesskey"
  region                = "us-east-1"
  description           = "Managed by AzureDevOps"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.
* `service_endpoint_name` - (Required) The Service Endpoint name.
* `access_key_id` - (Required) The AWS access key ID for signing programmatic requests.
* `secret_access_key` - (Required) The AWS secret access key for signing programmatic requests.
* `region` - (Required) The AWS region to send requests to.
* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links
* [azure-pipelines-terraform](https://github.com/microsoft/azure-pipelines-terraform)
* [Azure DevOps Service REST API 6.0 - Agent Pools](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-6.0)

## Import
Azure DevOps Service Endpoint AWS for Terraform can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
 terraform import azuredevops_serviceendpoint_aws_terraform.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```

