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
data "azuredevops_agent_pool" "example" {
  name = "Example Agent Pool"
}

output "name" {
  value = data.azuredevops_agent_pool.example.name
}

output "pool_type" {
  value = data.azuredevops_agent_pool.example.pool_type
}

output "auto_provision" {
  value = data.azuredevops_agent_pool.example.auto_provision
}

output "auto_update" {
  value = data.azuredevops_agent_pool.example.auto_update
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Agent Pool.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the agent pool

* `pool_type` - Specifies whether the agent pool type is Automation or Deployment.

* `auto_provision` - Specifies whether a queue should be automatically provisioned for each project collection.

* `auto_update` - Specifies whether or not agents within the pool should be automatically updated.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Agent Pools - Get](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/pools/get?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minute) Used when retrieving the Agent Pool.