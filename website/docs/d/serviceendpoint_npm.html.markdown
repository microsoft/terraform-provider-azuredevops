---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_npm"
description: |-
  Gets information about an existing NPM Service Endpoint. 
---

# Data Source : azuredevops_serviceendpoint_npm

Use this data source to access information about an existing NPM Service Endpoint.

## Example Usage

```hcl
data "azuredevops_serviceendpoint_npm" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example npm"
}

output "service_endpoint_id" {
  value = data.azuredevops_serviceendpoint_npm.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

---

* `service_endpoint_id` - (Optional) the ID of the Service Endpoint.

* `service_endpoint_name` - (Optional) the Name of the Service Endpoint.

~> **NOTE:** One of either `service_endpoint_id` or `service_endpoint_name` must be specified.

## Attributes Reference

In addition to the Arguments list above - the following Attributes are exported:

* `authorization` - Specifies the Authorization Scheme Map.

* `url` - Specifies the URL of the npm registry to connect with.

* `description` - Specifies the description of the Service Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minute) Used when retrieving the NPM Service Endpoint.
