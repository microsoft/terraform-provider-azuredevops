---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_powerplatform"
description: |-
  Manages a PowerPlatform Service Endpoint.
---

# azuredevops_serviceendpoint_powerplatform

Manages a PowerPlatform Service Endpoint.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_powerplatform" "example" {
  project_id                       = data.azuredevops_project.project.id
  service_endpoint_name            = "PowerPlaform-connection"
  description                      = "Managed by Terraform"
  url                              = "https://dev-environment.crm11.dynamics.com/"
 credentials {
   serviceprincipalid    = "00000000-0000-0000-0000-000000000000"
   serviceprincipalkey   = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
   tenantId              = "00000000-0000-0000-0000-000000000000"
 }
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project. Changing this forces a new PowerPlatform Service Endpoint to be created.
* `service_endpoint_name` - (Required) The Service Endpoint Name.
* `url` - (Required) The Service Endpoint url.

---

* `credentials` - (Optional) A `credentials` block as defined below.
* `description` - (Optional) Service connection description.

---

A `credentials` block supports the following:

* `serviceprincipalid` - (Required) The service principal application ID.
* `serviceprincipalkey` - (Required) The service principal application key.
* `tenantId` - (Required) The service principal tenant id.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the PowerPlatform Service Endpoint.

* `authorization` - A `authorization` block as defined below.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the PowerPlatform Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the PowerPlatform Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the PowerPlatform Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the PowerPlatform Service Endpoint.

## Import

PowerPlatform Service Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azuredevops_serviceendpoint_powerplatform.example 00000000-0000-0000-0000-000000000000
```
