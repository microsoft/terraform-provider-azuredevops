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
resource "azuredevops_agent_pool" "pool" {
  name           = "sample-pool"
  auto_provision = false
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the agent pool.
- `auto_provision` - (Optional) Specifies whether or not a queue should be automatically provisioned for each project collection. Defaults to `false`.
- `pool_type` - (Optional) Specifies whether the agent pool type is Automation or Deployment. Defaults to `automation`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the agent pool.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Agent Pools](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/pools?view=azure-devops-rest-6.0)

## Import

Azure DevOps Agent Pools can be imported using the agent pool ID, e.g.

```sh
terraform import azuredevops_agent_pool.pool 42
```
