## 0.0.2 (Unreleased)

FEATURES:
* **New Resource**  `azuredevops_git_permissions` [#18](https://github.com/terraform-providers/terraform-provider-azuredevops/pull/18)
* **New Resource**  `azuredevops_project_permissions` [#18](https://github.com/terraform-providers/terraform-provider-azuredevops/issues/18)
* **New Resource**  `azuredevops_serviceendpoint_aws` [#58](https://github.com/terraform-providers/terraform-provider-azuredevops/issues/58)
* **New Resource** `azuredevops_branch_policy_auto_reviewers` [#71](https://github.com/terraform-providers/terraform-provider-azuredevops/pull/71)
* **New Data Resource**  `azuredevops_git_repository` [#18](https://github.com/terraform-providers/terraform-provider-azuredevops/issues/18)

IMPROVEMENTS:
* **All resources: Remove from `.state` if project has been deleted** [#25](https://github.com/terraform-providers/terraform-provider-azuredevops/issues/25)
* **Add `path_filter` support for repository build policy** [#62](https://github.com/terraform-providers/terraform-provider-azuredevops/issues/62)

BUG FIXS:
* **GitHub service connection API breaking change** [#72](https://github.com/terraform-providers/terraform-provider-azuredevops/issues/72)
* **Terraform crash when the service connection description is set to an empty string** [#60](https://github.com/terraform-providers/terraform-provider-azuredevops/pull/60)

## 0.0.1 (June 18, 2020)

NOTES:
* The Azure DevOps provider can be used to configure Azure DevOps project in [Microsoft Azure](https://azure.microsoft.com/en-us/) using [Azure DevOps Service REST API](https://docs.microsoft.com/en-us/rest/api/azure/devops/?view=azure-devops-rest-5.1)
    
FEATURES:
* **New Resource** `azuredevops_build_definition`                               
* **New Resource** `azuredevops_project`                                                 
* **New Resource** `azuredevops_variable_group` [#21](https://github.com/microsoft/terraform-provider-azuredevops/issues/21)
* **New Resource** `azuredevops_serviceendpoint_github` [#3](https://github.com/microsoft/terraform-provider-azuredevops/issues/3)
* **New Resource** `azuredevops_serviceendpoint_dockerregistry` [#297](https://github.com/microsoft/terraform-provider-azuredevops/issues/3)
* **New Resource** `azuredevops_serviceendpoint_azurerm` [#3](https://github.com/microsoft/terraform-provider-azuredevops/issues/3)
* **New Resource** `azuredevops_git_repository` [#94](https://github.com/microsoft/terraform-provider-azuredevops/issues/94) [#95](https://github.com/microsoft/terraform-provider-azuredevops/issues/95) [#96](https://github.com/microsoft/terraform-provider-azuredevops/issues/96) [#97](https://github.com/microsoft/terraform-provider-azuredevops/issues/97)
* **New Resource** `azuredevops_user_entitlement` [#125](https://github.com/microsoft/terraform-provider-azuredevops/issues/125)
* **New Resource** `azuredevops_group_membership` [#74](github.com/microsoft/terraform-provider-azuredevops/issues/74)
* **New Resource** `azuredevops_agent_pool` [#22](https://github.com/microsoft/terraform-provider-azuredevops/issues/22)
* **New Resource** `azuredevops_group` [#103](https://github.com/microsoft/terraform-provider-azuredevops/issues/103)
* **New Data Source** `azuredevops_group` [#126](https://github.com/microsoft/terraform-provider-azuredevops/issues/126)
* **New Data Source** `azuredevops_projects` [#17](https://github.com/microsoft/terraform-provider-azuredevops/issues/17)
