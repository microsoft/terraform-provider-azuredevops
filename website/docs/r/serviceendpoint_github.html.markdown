---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_github"
description: |-
  Manages a GitHub service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_github

Manages a GitHub service endpoint within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_github" "serviceendpoint_gh_1" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "Sample GithHub Personal Access Token"

  auth_personal {
    # Also can be set with AZDO_GITHUB_SERVICE_CONNECTION_PAT environment variable
    personal_access_token = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
}
```

```hcl
resource "azuredevops_serviceendpoint_github" "serviceendpoint_gh_2" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "Sample GithHub Grant"

  auth_oauth {
    oauth_configuration_id = "00000000-0000-0000-0000-000000000000"
  }
}
```

```hcl
resource "azuredevops_serviceendpoint_github" "serviceendpoint_gh_3" {
  project_id = azuredevops_project.project.id
  service_endpoint_name = "Sample GithHub Apps: Azure Pipelines"
  # Note Github Apps do not support a description and will always be empty string. Must be explicitly set to override the default value.
  description = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The project ID or project name.
- `service_endpoint_name` - (Required) The Service Endpoint name.
- `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.
- `auth_personal` - (Optional) An `auth_personal` block as documented below. Allows connecting using a personal access token.
- `auth_oauth` - (Optional) An `auth_oauth` block as documented below. Allows connecting using an Oauth token.

**NOTE: Github Apps can not be created or updated via terraform. You must install and configure the app on Github and then import it. You must also set the `description` to "" explicitly."**

`auth_personal` block supports the following:

- `personal_access_token` - (Required) The Personal Access Token for Github.

`auth_oauth` block supports the following:

- `oauth_configuration_id` - (Required) **NOTE: Github OAuth flow can not be performed via terraform. You must create this on Azure DevOps and then import it.** The OAuth Configuration ID.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The project ID or project name.
- `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Service Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-5.1)

## Import

Azure DevOps Service Endpoint GitHub can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
$ terraform import azuredevops_serviceendpoint_github.serviceendpoint 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
