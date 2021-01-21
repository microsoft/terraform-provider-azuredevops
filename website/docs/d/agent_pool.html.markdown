---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_agent_pool"
description: |-
  Use this data source to access information about an existing Agent Pool within Azure DevOps.
---

# Data Source: azuredevops_agent_pool

Use this data source to access information about an existing Agent Pool within Azure DevOps.

## Example Usage

```hcl
data "azuredevops_agent_pool" "pool" {
  name = "Sample Agent Pool"
}

output "name" {
  value = data.azuredevops_agent_pool.pool.name
}

output "pool_type" {
  value = data.azuredevops_agent_pool.pool.pool_type
}

output "auto_provision" {
  value = data.azuredevops_agent_pool.pool.auto_provision
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Name of the Agent Pool.

## Attributes Reference

The following attributes are exported:

`name` - The name of the agent pool
`pool_type` - Specifies whether the agent pool type is Automation or Deployment.
`auto_provision` - Specifies whether or not a queue should be automatically provisioned for each project collection.

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Agent Pools - Get](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/pools/get?view=azure-devops-rest-5.1)
