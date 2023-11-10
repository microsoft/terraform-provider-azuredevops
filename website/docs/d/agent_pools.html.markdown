---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_agent_pools"
description: |-
  Use this data source to access information about existing Agent Pools within Azure DevOps.
---

# Data Source: azuredevops_agent_pools

Use this data source to access information about existing Agent Pools within Azure DevOps.

## Example Usage

```hcl
data "azuredevops_agent_pools" "example" {
}

output "agent_pool_name" {
  value = data.azuredevops_agent_pools.example.agent_pools.*.name
}

output "auto_provision" {
  value = data.azuredevops_agent_pools.example.agent_pools.*.auto_provision
}

output "auto_update" {
  value = data.azuredevops_agent_pools.example.agent_pools.*.auto_update
}

output "pool_type" {
  value = data.azuredevops_agent_pools.example.agent_pools.*.pool_type
}
```

## Argument Reference

This data source has no arguments

## Attributes Reference

The following attributes are exported:

- `agent_pools` - A list of existing agent pools in your Azure DevOps Organization with the following details about every agent pool:
  - `name` - The name of the agent pool
  - `pool_type` - Specifies whether the agent pool type is Automation or Deployment.
  - `auto_provision` - Specifies whether or not a queue should be automatically provisioned for each project collection.
  - `auto_update` - Specifies whether or not agents within the pool should be automatically updated.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Agent Pools - Get](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/pools/get?view=azure-devops-rest-7.0)
