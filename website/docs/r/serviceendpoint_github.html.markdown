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
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_github" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example GitHub Personal Access Token"

  auth_personal {
    # Also can be set with AZDO_GITHUB_SERVICE_CONNECTION_PAT environment variable
    personal_access_token = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
}
```

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_github" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example GitHub"
  auth_oauth {
    oauth_configuration_id = "00000000-0000-0000-0000-000000000000"
  }
}
```

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_github" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example GitHub Apps: Azure Pipelines"
  # Note Github Apps do not support a description and will always be empty string. Must be explicitly set to override the default value.  
  description = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_name` - (Required) The Service Endpoint name.

---

* `auth_oauth` - (Optional) An `auth_oauth` block as documented below. Allows connecting using an Oauth token.

* `auth_personal` - (Optional) An `auth_personal` block as documented below. Allows connecting using a personal access token.

* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

    ~>**NOTE:** GitHub Apps can not be created or updated via terraform. You must install and configure the app on GitHub and then import it. You must also set the `description` to "" explicitly."

---

`auth_personal` block supports the following:

* `personal_access_token` - (Required) The Personal Access Token for GitHub.

---

An `auth_oauth` block supports the following:

* `oauth_configuration_id` - (Required) The OAuth Configuration ID.

  ~>**NOTE:** GitHub OAuth flow can not be performed via terraform. You must create this on Azure DevOps and then import it.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Service Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the GitHub Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the GitHub Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the GitHub Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the GitHub Service Endpoint.

## Import

Azure DevOps GitHub Service Endpoint can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_github.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
