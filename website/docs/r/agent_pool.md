# azuredevops_project
Manages an agent pool within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_agent_pool" "pool" {
    name = "staging_pool"
    auto_provision = false
    is_hosted = false
}
```

## Arugument Reference

The following arguments are supported:

* `name` - (Required) The name of the agent pool.
* `auto_provision` - (Optional) Specifies whether to auto provision the agent pool in new projects. - default is false.
* `is_hosted` - (Optional) Specifies whether the agent pool is hosted or private. - default is false
* `pool_type` - (Optional) Specifies whether the agent pool type is Automation or Deployment.  default is "automation"

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the agent pool.

## Relevant Links
* [Azure DevOps Service REST API 5.1 - Agent Pools](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/pools?view=azure-devops-rest-5.1)

## Import
Azure DevOps Agent Pools can be imported using the agent pool id, e.g.

```
 terraform import azuredevops_agent_pool.pool 42
```