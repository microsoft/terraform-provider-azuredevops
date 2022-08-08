---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_maven"
description: |-
  Manages a Service Connection for Maven.
---

# azuredevops_serviceendpoint_maven

Manages a Maven service endpoint within Azure DevOps, which can be used as a resource in YAML pipelines to connect to a Maven instance.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_maven" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "maven-example"
  description           = "Service Endpoint for 'Maven' (Managed by Terraform)"
  url                   = "https://example.com"
  repository_id         = "example"

  authentication_token {
    token = "0000000000000000000000000000000000000000"
  }
}
```

Alternatively a username and password may be used.

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_maven" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "maven-example"
  description           = "Service Endpoint for 'Maven' (Managed by Terraform)"
  url                   = "https://example.com"
  repository_id         = "example"

  authentication_basic {
    username              = "username"
    password              = "password"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project. Changing this forces a new Service Connection Maven to be created.
* `service_endpoint_name` - (Required) The name of the service endpoint. Changing this forces a new Service Connection Maven to be created.
* `description` - (Optional) The Service Endpoint description. Defaults to Managed by Terraform.
* `url` - (Required) The Service Endpoint url.
* `repository_id` - (Required) The Repository id.
* either `authentication_token` or `authentication_basic` (one is required)
  * `authentication_token`
    * `token` - Authentication Token generated through maven repository.
  * `authentication_basic`
    * `username` - Maven Repository Username.
    * `password` - Maven Repository Password.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.

## Import

Service Connection Maven can be imported using the `projectId/id` or or `projectName/id`, e.g.

```shell
terraform import azuredevops_serviceendpoint_maven.example 00000000-0000-0000-0000-000000000000
```
