## 0.1.4 (Unreleased)
FEATURES:
* **New Resource** `azuredevops_serviceendpoint_ssh` [#270](https://github.com/microsoft/terraform-provider-azuredevops/issues/270)
* **New Resource** `azuredevops_serviceendpoint_npm` [#334](https://github.com/microsoft/terraform-provider-azuredevops/issues/334)
* **New Resource** `azuredevops_serviceendpoint_azuredevops` [#339](https://github.com/microsoft/terraform-provider-azuredevops/issues/339)
* **New Resource** `azuredevops_serviceendpoint_github_enterprise` [#210](https://github.com/microsoft/terraform-provider-azuredevops/issues/210)

IMPROVEMENTS:
* `azuredevops_group` - Support for changing group display names [#356](https://github.com/microsoft/terraform-provider-azuredevops/issues/356)
  
BUG FIX:
  `azuredevops_group` - `scope` will be suppressed during `plan` and `apply`  [#345](https://github.com/microsoft/terraform-provider-azuredevops/issues/345)

## 0.1.3
FEATURES:
* **New Resource** `azuredevops_branch_policy_merge_types` [#300](https://github.com/microsoft/terraform-provider-azuredevops/issues/300)

IMPROVEMENTS:
* Support darwin/arm64 (Apple Silicon) [#332](https://github.com/microsoft/terraform-provider-azuredevops/issues/332)
* All service endpoints - Description accept any string between 0~1024 in length [#295](https://github.com/microsoft/terraform-provider-azuredevops/issues/295)
* `azuredevops_git_repository` - Support import Azure Git repository resource [#43](https://github.com/microsoft/terraform-provider-azuredevops/issues/43)
* `azuredevops_serviceendpoint_azurecr` - Support expose service principal ID [#317](https://github.com/microsoft/terraform-provider-azuredevops/issues/317)
* `azuredevops_serviceendpoint_github` - Compatible with GitHub App service connection [#326](https://github.com/microsoft/terraform-provider-azuredevops/issues/326)

BUG FIX:
* `azuredevops_serviceendpoint_azurecr` - Fix unable to update the description  [#312](https://github.com/microsoft/terraform-provider-azuredevops/issues/312)
* `azuredevops_branch_policy_build_validation` - Handle deleted policy [#330](https://github.com/microsoft/terraform-provider-azuredevops/issues/330)

## 0.1.2

FEATURES:
* **New Resource** `azuredevops_serviceendpoint_artifactory` [#256](https://github.com/microsoft/terraform-provider-azuredevops/issues/256)
* **New Resource** `azuredevops_serviceendpoint_sonarqube` [#257](https://github.com/microsoft/terraform-provider-azuredevops/issues/257)

IMPROVEMENTS:
* `azuredevops_serviceendpoint_azurecr` - Change docker registry login server to lowercase [#277](https://github.com/microsoft/terraform-provider-azuredevops/issues/277)
* `azuredevops_serviceendpoint_github` - Enhance `auth_...` configuration block check [#275](https://github.com/microsoft/terraform-provider-azuredevops/issues/275)
* `azuredevops_branch_policy_min_reviewers` - Support new configurations [#255](https://github.com/microsoft/terraform-provider-azuredevops/issues/255)
  - `last_pusher_cannot_approve` - Prohibit the most recent pusher from approving their own changes. Defaults to false.
  - `allow_completion_with_rejects_or_waits` - Allow completion even if some reviewers vote to wait or reject. Defaults to false.
  - `on_push_reset_approved_votes` - When new changes are pushed reset all approval votes (does not reset votes to reject or wait). Defaults to false.
  - `on_push_reset_all_votes` - When new changes are pushed reset all code reviewer votes. Defaults to false.
  - `on_last_iteration_require_vote` - On last iteration require vote. Defaults to false.

BUG FIX:
* All service endpoint resources - Add resource status check during creation and deletion [#261](https://github.com/microsoft/terraform-provider-azuredevops/issues/261)
* `azuredevops_variable_group` - Key vault variables will be verified with Azure key vault secrets [#252](https://github.com/microsoft/terraform-provider-azuredevops/issues/252)

## 0.1.1 

FEATURES:
* **New Resource** `azuredevops_build_definition_permissions` [#254](https://github.com/microsoft/terraform-provider-azuredevops/issues/254)
* **New Resource** `azuredevops_serviceendpoint_runpipeline` [#182](https://github.com/microsoft/terraform-provider-azuredevops/issues/182)

IMPROVEMENTS:   
`azuredevops_serviceendpoint_kubernetes` - Support `cluster_admin` in Kubernetes service connections [#218](https://github.com/microsoft/terraform-provider-azuredevops/issues/218)
`azuredevops_git_repository` - Remove `source_type` default value [#265](https://github.com/microsoft/terraform-provider-azuredevops/issues/265)

## 0.1.0

FEATURES:
* **New Resource** `azuredevops_git_permissions` [#18](https://github.com/microsoft/terraform-provider-azuredevops/issues/18)
* **New Resource** `azuredevops_project_permissions` [#18](https://github.com/microsoft/terraform-provider-azuredevops/issues/18)
* **New Resource** `azuredevops_serviceendpoint_aws` [#58](https://github.com/microsoft/terraform-provider-azuredevops/issues/58)
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
