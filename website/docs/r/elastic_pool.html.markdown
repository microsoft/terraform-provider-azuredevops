---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_elastic_pool"
description: |-
  Manages Elastic pool within Azure DevOps organization.
---

# azuredevops_agent_pool

Manages Elastic pool within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_azurerm" "example" {
  project_id                             = azuredevops_project.example.id
  service_endpoint_name                  = "Example Azure Connection"
  description                            = "Managed by Terraform"
  service_endpoint_authentication_scheme = "ServicePrincipal"
  credentials {
    serviceprincipalid  = "00000000-0000-0000-0000-000000000000"
    serviceprincipalkey = "00000000-0000-0000-0000-000000000000"
  }
  azurerm_spn_tenantid      = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_id   = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_name = "Subscription Name"
}

resource "azuredevops_elastic_pool" "example" {
  name                   = "Example Elastic Pool"
  service_endpoint_id    = azuredevops_serviceendpoint_azurerm.example.id
  service_endpoint_scope = azuredevops_project.example.id
  desired_idle           = 2
  max_capacity           = 3
  azure_resource_id      = "/subscriptions/<Subscription Id>/resourceGroups/<Resource Name>/providers/Microsoft.Compute/virtualMachineScaleSets/<VMSS Name>"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the Elastic pool.

- `azure_resource_id` - (Required) The ID of the Azure resource.

- `service_endpoint_id` - (Required) The ID of Service Endpoint used to connect to Azure.

- `service_endpoint_scope` - (Required) The Project ID of Service Endpoint belongs to.

- `desired_idle` - (Required) Number of agents to keep on standby.

- `max_capacity` - (Required) Maximum number of virtual machines in the scale set.

---
- `recycle_after_each_use` - (Optional) Tear down virtual machines after every use. Defaults to `false`.

- `time_to_live_minutes` - (Optional) Delay in minutes before deleting excess idle agents. Defaults to `30`.

- `agent_interactive_ui` - (Optional) Set whether agents should be configured to run with interactive UI. Defaults to `false`.

- `auto_provision` - (Optional) Specifies whether a queue should be automatically provisioned for each project collection. Defaults to `false`.

- `auto_update` - (Optional) Specifies whether or not agents within the pool should be automatically updated. Defaults to `true`.

- `project_id` - (Optional) The ID of the project where a new Elastic Pool will be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Elastic pool.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Elastic Pools](https://learn.microsoft.com/en-us/rest/api/azure/devops/distributedtask/elasticpools/create?view=azure-devops-rest-7.0)

## Import

Azure DevOps Agent Pools can be imported using the Elastic pool ID, e.g.

```sh
terraform import azuredevops_elastic_pool.example 0
```
