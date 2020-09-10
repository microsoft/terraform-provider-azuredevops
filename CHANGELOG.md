## 0.0.2 (Unreleased)

FEATURES:
* **New Resource**  `azuredevops_git_permissions` [#18](https://github.com/microsoft/terraform-provider-azuredevops/issues/18)
* **New Resource**  `azuredevops_project_permissions` [#18](https://github.com/microsoft/terraform-provider-azuredevops/issues/18)
* **New Resource**  `azuredevops_serviceendpoint_aws` [#58](https://github.com/microsoft/terraform-provider-azuredevops/issues/58)
* **New Resource**  `azuredevops_serviceendpoint_devops` [#58](https://github.com/microsoft/terraform-provider-azuredevops/issues/182)
* **New Resource** `azuredevops_branch_policy_auto_reviewers` [#71](https://github.com/microsoft/terraform-provider-azuredevops/issues/71)
* **New Resource** `azuredevops_workitemquery_permissions` [#79](https://github.com/microsoft/terraform-provider-azuredevops/issues/79)
* **New Resource** `azuredevops_serviceendpoint_azurecr` [#119](https://github.com/microsoft/terraform-provider-azuredevops/issues/119/)
* **New Resource** `azuredevops_area_permissions` [#85](https://github.com/microsoft/terraform-provider-azuredevops/issues/85)
* **New Resource** `azuredevops_iteration_permissions` [#85](https://github.com/microsoft/terraform-provider-azuredevops/issues/85)
* **New Resource** `azuredevops_branch_policy_work_item_linking` [#144](https://github.com/microsoft/terraform-provider-azuredevops/issues/144)
* **New Resource** `azuredevops_branch_policy_comment_resolution` [#144](https://github.com/microsoft/terraform-provider-azuredevops/issues/144)
* **New Data Resource**  `azuredevops_git_repository` [#18](https://github.com/microsoft/terraform-provider-azuredevops/issues/18)
* **New Data Resource**  `azuredevops_area` [#85](https://github.com/microsoft/terraform-provider-azuredevops/issues/85)
* **New Data Resource**  `azuredevops_iteration` [#85](https://github.com/microsoft/terraform-provider-azuredevops/issues/85)

IMPROVEMENTS:
* **All resources: Remove from `.state` if project has been deleted** [#25](https://github.com/microsoft/terraform-provider-azuredevops/issues/25)
* **`azuredevops_branch_policy_build_validation`: Add `filename_patterns` support for repository build policy** [#62](https://github.com/microsoft/terraform-provider-azuredevops/issues/62)
* **`azuredevops_git_repository`:
    - Use `default_branch` as the name of an initialized branch [#89](https://github.com/microsoft/terraform-provider-azuredevops/pull/89)
    - Add support for import Git repository [#45](https://github.com/microsoft/terraform-provider-azuredevops/issues/45)
* **`azuredevops_build_definition`:**
    - Add Support for GitHub enterprise as a build definition repository type [#97](https://github.com/microsoft/terraform-provider-azuredevops/pull/97)
    - Add Support for report build status configuration [#63](https://github.com/microsoft/terraform-provider-azuredevops/issues/63)

BUG FIX:
* **`azuredevops_serviceendpoint_github`: GitHub service connection API breaking change** [#72](https://github.com/microsoft/terraform-provider-azuredevops/issues/72)
* **All service connection resources: Terraform crashes when the service connection description is set to an empty string** [#60](https://github.com/microsoft/terraform-provider-azuredevops/pull/60)
* **`azuredevops_build_definition`: Build Definition creation failed when repository type is GitHub** [#65](https://github.com/microsoft/terraform-provider-azuredevops/issues/65)

BREAKING CHANGES:
* `azuredevops_git_repository` - `initialization` now is a required configuration.  [#54](https://github.com/microsoft/terraform-provider-azuredevops/issues/54)

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
* **New Resource** `azuredevops_group_membership` [#74](https://github.com/microsoft/terraform-provider-azuredevops/issues/74)
* **New Resource** `azuredevops_agent_pool` [#22](https://github.com/microsoft/terraform-provider-azuredevops/issues/22)
* **New Resource** `azuredevops_group` [#103](https://github.com/microsoft/terraform-provider-azuredevops/issues/103)
* **New Data Source** `azuredevops_group` [#126](https://github.com/microsoft/terraform-provider-azuredevops/issues/126)
* **New Data Source** `azuredevops_projects` [#17](https://github.com/microsoft/terraform-provider-azuredevops/issues/17)
