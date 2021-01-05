## 0.1.1 (Unreleased)

FEATURES:
* **New Resource** `azuredevops_build_definition_permissions` [#254](https://github.com/microsoft/terraform-provider-azuredevops/issues/254)

IMPROVEMENTS:   
`azuredevops_serviceendpoint_kubernetes` - Support `cluster_admin` in Kubernetes service connections [#218](https://github.com/microsoft/terraform-provider-azuredevops/issues/218)

## 0.1.0

FEATURES:
* **New Resource** `azuredevops_git_permissions` [#18](https://github.com/microsoft/terraform-provider-azuredevops/issues/18)
* **New Resource** `azuredevops_project_permissions` [#18](https://github.com/microsoft/terraform-provider-azuredevops/issues/18)
* **New Resource** `azuredevops_serviceendpoint_aws` [#58](https://github.com/microsoft/terraform-provider-azuredevops/issues/58)
* **New Resource** `azuredevops_serviceendpoint_runpipeline` [#182](https://github.com/microsoft/terraform-provider-azuredevops/issues/182)
* **New Resource** `azuredevops_branch_policy_auto_reviewers` [#71](https://github.com/microsoft/terraform-provider-azuredevops/issues/71)
* **New Resource** `azuredevops_workitemquery_permissions` [#79](https://github.com/microsoft/terraform-provider-azuredevops/issues/79)
* **New Resource** `azuredevops_serviceendpoint_azurecr` [#119](https://github.com/microsoft/terraform-provider-azuredevops/issues/119/)
* **New Resource** `azuredevops_area_permissions` [#85](https://github.com/microsoft/terraform-provider-azuredevops/issues/85)
* **New Resource** `azuredevops_iteration_permissions` [#85](https://github.com/microsoft/terraform-provider-azuredevops/issues/85)
* **New Resource** `azuredevops_branch_policy_work_item_linking` [#144](https://github.com/microsoft/terraform-provider-azuredevops/issues/144)
* **New Resource** `azuredevops_branch_policy_comment_resolution` [#144](https://github.com/microsoft/terraform-provider-azuredevops/issues/144)
* **New Data Resource** `azuredevops_git_repository` [#18](https://github.com/microsoft/terraform-provider-azuredevops/issues/18)
* **New Data Resource** `azuredevops_area` [#85](https://github.com/microsoft/terraform-provider-azuredevops/issues/85)
* **New Data Resource** `azuredevops_iteration` [#85](https://github.com/microsoft/terraform-provider-azuredevops/issues/85)
* **New Data Resource** `azuredevops_agent_queue` [#175](https://github.com/microsoft/terraform-provider-azuredevops/issues/175)

IMPROVEMENTS:

* All resources - remove from `.tfstate` if project has been deleted [#25](https://github.com/microsoft/terraform-provider-azuredevops/issues/25)
* Data source `azuredevops_build_definition` - support export `origin` and `origin_id` [#177](https://github.com/microsoft/terraform-provider-azuredevops/issues/177)
* Data source `azuredevops_project` - add `project_id` for data source configuration [#163](https://github.com/microsoft/terraform-provider-azuredevops/issues/163)
* `azuredevops_branch_policy_build_validation`  - add `filename_patterns` support for repository build policy [#62](https://github.com/microsoft/terraform-provider-azuredevops/issues/62)
* `azuredevops_git_repository`
    - Use `default_branch` as the name of an initialized branch [#89](https://github.com/microsoft/terraform-provider-azuredevops/issues/89)
    - Add support for import Git repository [#45](https://github.com/microsoft/terraform-provider-azuredevops/issues/45)
* `azuredevops_build_definition`
    - Add Support for GitHub enterprise as a build definition repository type [#97](https://github.com/microsoft/terraform-provider-azuredevops/issues/97)
    - Add Support for report build status configuration [#63](https://github.com/microsoft/terraform-provider-azuredevops/issues/63)
* Data Resource `azuredevops_group` support search for project collection groups [#200](https://github.com/microsoft/terraform-provider-azuredevops/issues/200)

BUG FIX:
* All service connection resources - Terraform crashes when the service connection description is set to an empty string [#60](https://github.com/microsoft/terraform-provider-azuredevops/issues/60)
* Resource import - set the project ID to `project_id` [#172](https://github.com/microsoft/terraform-provider-azuredevops/issues/172)
* `azuredevops_build_definition` - build Definition creation failed when repository type is GitHub [#65](https://github.com/microsoft/terraform-provider-azuredevops/issues/65)
* `azuredevops_serviceendpoint_github` - GitHub service connection API breaking change [#72](https://github.com/microsoft/terraform-provider-azuredevops/issues/72)

BREAKING CHANGES:
* `azuredevops_git_repository` - `initialization` is a required configuration [#54](https://github.com/microsoft/terraform-provider-azuredevops/issues/54)
* `azuredevops_project` - rename `project_name` to `name` [#179](https://github.com/microsoft/terraform-provider-azuredevops/issues/179)

## 0.0.1 (June 18, 2020)

NOTES:
* The Azure DevOps provider can be used to configure Azure DevOps project in [Microsoft Azure](https://azure.microsoft.com/en-us/) using [Azure DevOps Service REST API](https://docs.microsoft.com/en-us/rest/api/azure/devops/?view=azure-devops-rest-5.1)

FEATURES:
* **New Resource** `azuredevops_build_definition`                               
* **New Resource** `azuredevops_project`                                                 
* **New Resource** `azuredevops_variable_group`
* **New Resource** `azuredevops_serviceendpoint_github`
* **New Resource** `azuredevops_serviceendpoint_dockerregistry`
* **New Resource** `azuredevops_serviceendpoint_azurerm`
* **New Resource** `azuredevops_git_repository`
* **New Resource** `azuredevops_user_entitlement`
* **New Resource** `azuredevops_group_membership`
* **New Resource** `azuredevops_agent_pool`
* **New Resource** `azuredevops_group`
* **New Data Source** `azuredevops_group`
* **New Data Source** `azuredevops_projects`
