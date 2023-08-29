---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_agent_pool"
description: |-
  Manages an agent pool within Azure DevOps organization.
---

# azuredevops_agent_pool

Manages an agent pool within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_agent_pool" "example" {
  name           = "Example-pool"
  auto_provision = false
  auto_update    = false
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the agent pool.
- `auto_provision` - (Optional) Specifies whether a queue should be automatically provisioned for each project collection. Defaults to `false`.
- `pool_type` - (Optional) Specifies whether the agent pool type is Automation or Deployment. Defaults to `automation`.
- `auto_update` - (Optional) Specifies whether or not agents within the pool should be automatically updated. Defaults to `true`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the agent pool.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Agent Pools](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/pools?view=azure-devops-rest-7.0)

## Import

Azure DevOps Agent Pools can be imported using the agent pool ID, e.g.

```sh
terraform import azuredevops_agent_pool.example 0
```
