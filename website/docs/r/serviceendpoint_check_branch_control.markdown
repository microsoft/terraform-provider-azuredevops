---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_check_branch_control"
description: |-
  Manages a Branch Control check for service endpoints.
---

# azuredevops_serviceendpoint_check_branch_control
Manages a Branch Control check for service endpoints.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_jfrog_artifactory_v2" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example JFrog Artifactory V2"
  description           = "Managed by Terraform"
  url                   = "https://artifactory.my.com"
  authentication_token {
    token = "0000000000000000000000000000000000000000"
  }
}

resource "azuredevops_serviceendpoint_check_branch_control" "example" {
  project_id                       = azuredevops_project.example.id
  endpoint_id                      = azuredevops_serviceendpoint_jfrog_artifactory_v2.example.id
  display_name                     = "Protected branches only"
  allowed_branches                 = ["refs/heads/releases/*", "refs/heads/main/*", "refs/heads/master/*"]
  verify_branch_protection         = true
  ignore_unknown_protection_status = false
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.
* `endpoint_id` - (Required) The Service Endpoint id.
* `allowed_branches` -
* `verify_branch_protection` -
* `ignore_unknown_protection_status` -


## Attributes Reference

The following attributes are exported:
TODO

## Relevant Links
TODO: Add branch protection doc URL
* [Azure DevOps Service Connections](https://docs.microsoft.com/en-us/azure/devops/pipelines/library/service-endpoints?view=azure-devops&tabs=yaml)
* [Artifactory User Token](https://docs.artifactory.org/latest/user-guide/user-token/)

## Import
TODO: Determine if this can be supported, are the check IDs available anywhere outside the API(?)
