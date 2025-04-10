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
    username = "username"
    password = "password"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project. Changing this forces a new Service Connection Maven to be created.

* `service_endpoint_name` - (Required) The name of the service endpoint. Changing this forces a new Service Connection Maven to be created.

* `url` - (Required) The URL of the Maven Repository.

* `repository_id` - (Required) The ID of the server that matches the id element of the `repository/mirror` that Maven tries to connect to.

---

* `description` - (Optional) The Service Endpoint description. Defaults to Managed by Terraform.

* `authentication_token` - (Optional) A `authentication_token` block as documented below.

* `authentication_basic` - (Optional) A `authentication_basic` block as documented below.

---

A `authentication_token` block supports the following:

* `token` - (Required) Authentication Token generated through maven repository.

---

A `authentication_basic` block supports the following:

* `username` - (Required) The Username of the Maven Repository.

* `password` - (Required) The password Maven Repository.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Maven Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the Maven Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the Maven Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the Maven Service Endpoint.

## Import

Azure DevOps Maven Service Connection can be imported using the `projectId/id` or `projectName/id`, e.g.

```shell
terraform import azuredevops_serviceendpoint_maven.example projectName/00000000-0000-0000-0000-000000000000
```
