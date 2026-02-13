---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_jfrog_xray_v2"
description: |-
  Gets information about an existing JFrog XRay V2 Service Endpoint. 
---

# Data Source : azuredevops_serviceendpoint_jfrog_xray_v2

Use this data source to access information about an existing JFrog XRay V2 Service Endpoint.

## Example Usage

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

data "azuredevops_serviceendpoint_jfrog_xray_v2" "example" {
  project_id            = data.azuredevops_project.example.id
  service_endpoint_name = "Example JFrog XRay V2"
}

output "service_endpoint_id" {
  value = data.azuredevops_serviceendpoint_jfrog_xray_v2.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_id` - (Optional) the ID of the Service Endpoint.

* `service_endpoint_name` - (Optional) the Name of the Service Endpoint.

~> **NOTE:** One of either `service_endpoint_id` or `service_endpoint_name` must be specified.

## Attributes Reference

In addition to the Arguments list above - the following Attributes are exported:

* `url` - The URL of the Artifactory server to connect with.
* `description` - Specifies the description of the Service Endpoint.
