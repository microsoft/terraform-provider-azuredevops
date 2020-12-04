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
data "azuredevops_agent_pools" "pools" {
}

output "agent_pool_name" {
  value = data.azuredevops_agent_pools.pools.agent_pools.*.name
}

output "auto_provision" {
  value = data.azuredevops_agent_pools.pools.agent_pools.*.auto_provision
}

output "pool_type" {
  value = data.azuredevops_agent_pools.pools.agent_pools.*.pool_type
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

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Agent Pools - Get](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/pools/get?view=azure-devops-rest-5.1)
