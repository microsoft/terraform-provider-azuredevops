---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_generic_v2"
description: |-
  Manages a Generic Service Endpoint (v2) within Azure DevOps.
---

# Resource: azuredevops_serviceendpoint_generic_v2

Manages a Generic Service Endpoint (v2) within Azure DevOps, which can be used to connect to various external services with custom authentication mechanisms.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

# Basic username/password authentication
resource "azuredevops_serviceendpoint_generic_v2" "example" {
  project_id            = azuredevops_project.example.id
  name                  = "Example Generic Service Endpoint"
  description           = "Managed by Terraform"
  service_endpoint_type = "generic"
  server_url            = "https://example.com"
  authorization_scheme  = "UsernamePassword"
  authorization_parameters = {
    username = "username"
    password = "password"
  }
}

# Token-based authentication
resource "azuredevops_serviceendpoint_generic_v2" "token_example" {
  project_id            = azuredevops_project.example.id
  name                  = "Token-based Service Endpoint"
  description           = "Managed by Terraform"
  service_endpoint_type = "generic"
  server_url            = "https://api.example.com"
  authorization_scheme  = "Token"
  authorization_parameters = {
    apitoken = "your-api-token"
  }

  parameters = {
    releaseUrl = "https://releases.example.com"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project to which the service endpoint belongs.
* `name` - (Required) The name of the service endpoint.
* `type` - (Required) The type of the service endpoint. This can be any valid service endpoint type, such as "generic", "artifactory", etc.
* `shared_project_ids` - (Optional) A list of project IDs where the service endpoint should be shared.
* `description` - (Optional) The description of the service endpoint. Defaults to "Managed by Terraform".
* `server_url` - (Required) The URL of the server associated with the service endpoint.
* `authorization_scheme` - (Required) The authorization scheme to use. Common values include "UsernamePassword", "Token", "OAuth", etc.
* `authorization_parameters` - (Optional) Map of key/value pairs for the specific authorization scheme. These often include sensitive data like tokens, usernames, and passwords.
* `parameters` - (Optional) Additional data associated with the service endpoint. This is a map of key/value pairs.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - The ID of the service endpoint.

## Import

Service endpoints can be imported using the project ID and service endpoint ID:

```
terraform import azuredevops_serviceendpoint_generic_v2.example <project_id>/<id>
```
