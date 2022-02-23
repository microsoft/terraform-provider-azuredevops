---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_aws"
description: |-
  Manages a AWS service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_aws
Manages a AWS service endpoint within Azure DevOps. Using this service endpoint requires you to first install [AWS Toolkit for Azure DevOps](https://marketplace.visualstudio.com/items?itemName=AmazonWebServices.aws-vsts-tools).

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_aws" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example AWS"
  access_key_id         = "00000000-0000-0000-0000-000000000000"
  secret_access_key     = "accesskey"
  description           = "Managed by AzureDevOps"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The project ID or project name.
* `service_endpoint_name` - (Required) The Service Endpoint name.
* `access_key_id` - (Required) The AWS access key ID for signing programmatic requests.
* `secret_access_key` - (Required) The AWS secret access key for signing programmatic requests.
* `session_token` - (Optional) The AWS session token for signing programmatic requests.
* `role_to_assume` - (Optional) The Amazon Resource Name (ARN) of the role to assume.
* `role_session_name` - (Optional) Optional identifier for the assumed role session.
* `external_id` - (Optional) A unique identifier that is used by third parties when assuming roles in their customers' accounts, aka cross-account role access.
* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The project ID or project name.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links
* [aws-toolkit-azure-devops](https://github.com/aws/aws-toolkit-azure-devops)
* [Azure DevOps Service REST API 6.0 - Agent Pools](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-6.0)

## Import
Azure DevOps Service Endpoint AWS can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```
 terraform import azuredevops_serviceendpoint_aws.serviceendpoint 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
