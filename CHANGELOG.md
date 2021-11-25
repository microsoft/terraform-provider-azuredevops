## 0.1.8 
FEATURES:
* **New Resource** `azuredevops_git_repository_file ` [#225](https://github.com/microsoft/terraform-provider-azuredevops/issues/225)
* **New Resource** `azuredevops_serviceendpoint_permissions ` [#249](https://github.com/microsoft/terraform-provider-azuredevops/issues/249)
* **New Data Resource** `azuredevops_groups ` [#483](https://github.com/microsoft/terraform-provider-azuredevops/issues/483)

IMPROVEMENTS:
* `azuredevops_build_definition`
  - Support scheduled triggers. [#445](https://github.com/microsoft/terraform-provider-azuredevops/issues/445)
  - Default agent pool has been updated from `Hosted Ubuntu 1604` to `Azure Pipelines`. [#466](https://github.com/microsoft/terraform-provider-azuredevops/issues/466)
* `azuredevops_serviceendpoint_azuredevops` - Extension [Configurable Pipeline Runner](https://marketplace.visualstudio.com/items?itemName=CSE-DevOps.RunPipelines) should be installed as documented. [#454](https://github.com/microsoft/terraform-provider-azuredevops/issues/454)
* `azuredevops_git_repository` - `initialization` should be ignored when importing as documented. [#467](https://github.com/microsoft/terraform-provider-azuredevops/issues/467)
* `azuredevops_branch_policy_status_check` - Support new property `genre`. [#472](https://github.com/microsoft/terraform-provider-azuredevops/issues/472)
* **Data Resource** `azuredevops_users` - Support export user IDs. [#400](https://github.com/microsoft/terraform-provider-azuredevops/issues/400)
* **Data Resource** `azuredevops_group` - Allow generic groups to be returned when searching the organization. [#485](https://github.com/microsoft/terraform-provider-azuredevops/issues/485)

BUG FIX:
* `azuredevops_user_entitlement` -
  - `principal_name` Suppress case sensitive. [#446](https://github.com/microsoft/terraform-provider-azuredevops/issues/446)
  - If user status is `Delete` or `None`, this resource will be removed from `.tfstate`. [#447](https://github.com/microsoft/terraform-provider-azuredevops/issues/447)
* All service endpoints:
  - Enhance service endpoint status handler. [#474](https://github.com/microsoft/terraform-provider-azuredevops/issues/474)
  - Compatible with when `Authorizaiton` is not returned by service. [#460](https://github.com/microsoft/terraform-provider-azuredevops/issues/460)

## 0.1.7
FEATURES:
* **New Resource** `azuredevops_team ` [#121](https://github.com/microsoft/terraform-provider-azuredevops/issues/121)
* **New Resource** `azuredevops_team_members` [#121](https://github.com/microsoft/terraform-provider-azuredevops/issues/121)
* **New Resource** `azuredevops_team_administrators` [#121](https://github.com/microsoft/terraform-provider-azuredevops/issues/121)
* **New Resource** `azuredevops_repository_policy_case_enforcement` [#113](https://github.com/microsoft/terraform-provider-azuredevops/issues/113)
* **New Resource** `azuredevops_repository_policy_reserved_names` [#113](https://github.com/microsoft/terraform-provider-azuredevops/issues/113)
* **New Resource** `azuredevops_repository_policy_max_path_length` [#113](https://github.com/microsoft/terraform-provider-azuredevops/issues/113)
* **New Resource** `azuredevops_repository_policy_max_file_size` [#113](https://github.com/microsoft/terraform-provider-azuredevops/issues/113)
* **New Resource** `azuredevops_repository_policy_check_credentials` [#113](https://github.com/microsoft/terraform-provider-azuredevops/issues/113)
* **New Data Resource** `azuredevops_team` [#121](https://github.com/microsoft/terraform-provider-azuredevops/issues/121)
* **New Data Resource** `azuredevops_teams` [#121](https://github.com/microsoft/terraform-provider-azuredevops/issues/121)

BREAKING CHANGES:
* **Resource** `azuredevops_repository_policy_author_email_pattern` - Remove `settings` and `scope`, policy `scope` can be set by  [repository ID](https://github.com/microsoft/terraform-provider-azuredevops/blob/master/website/docs/r/repository_policy_author_email_pattern.html.markdown) [#436](https://github.com/microsoft/terraform-provider-azuredevops/issues/436)
* **Resource** `azuredevops_repository_policy_file_path_pattern` - Remove `settings` and `scope`, policy `scope` can be set by  [repository ID](https://github.com/microsoft/terraform-provider-azuredevops/blob/master/website/docs/r/repository_policy_file_path_pattern.html.markdown) [#436](https://github.com/microsoft/terraform-provider-azuredevops/issues/436)

## 0.1.6
FEATURES:
* **New Resource** `serviceendpoint_generic` [#402](https://github.com/microsoft/terraform-provider-azuredevops/issues/402)
* **New Resource** `serviceendpoint_generic_git` [#402](https://github.com/microsoft/terraform-provider-azuredevops/issues/402)

IMPROVEMENTS:
* `resource_git_repository` - Support import private repository. [#236](https://github.com/microsoft/terraform-provider-azuredevops/issues/236)
* `azuredevops_git_permissions` - Can create permissions on non-existent branches. [#411](https://github.com/microsoft/terraform-provider-azuredevops/issues/411)
* `azuredevops_repository_policy_author_email_pattern` - Support project level repository policy setting [#113](https://github.com/microsoft/terraform-provider-azuredevops/issues/113)
* `azuredevops_repository_policy_file_path_pattern` - Support project level repository policy setting  [#113](https://github.com/microsoft/terraform-provider-azuredevops/issues/113)

BUG FIX:
* `azuredevops_git_repository` - Changing the `init_type` will recreate the repository. [#406](https://github.com/microsoft/terraform-provider-azuredevops/issues/406)
* `azuredevops_serviceendpoint_kubernetes` - Import crash.  [#414](https://github.com/microsoft/terraform-provider-azuredevops/issues/414)

## 0.1.5
FEATURES:
* **New Resource** `azuredevops_serviceendpoint_servicefabric` [#38](https://github.com/microsoft/terraform-provider-azuredevops/issues/38)
* **New Resource** `azuredevops_repository_policy_author_email_pattern` [#113](https://github.com/microsoft/terraform-provider-azuredevops/issues/113)
* **New Resource** `azuredevops_repository_policy_file_path_pattern` [#113](https://github.com/microsoft/terraform-provider-azuredevops/issues/113)
* **New Resource** `azuredevops_branch_policy_status_check` [#352](https://github.com/microsoft/terraform-provider-azuredevops/issues/352)

IMPROVEMENTS:
* `azuredevops_serviceendpoint_azurerm` - Credentials can be updated without recreate.  [#387](https://github.com/microsoft/terraform-provider-azuredevops/issues/387)

BUG FIX:
* `azuredevops_group` - Fix group scope not set  [#366](https://github.com/microsoft/terraform-provider-azuredevops/issues/366)
* `azuredevops_serviceendpoint_azurecr` - Fix container registry name cannot be updated.  [#391](https://github.com/microsoft/terraform-provider-azuredevops/issues/391)

## 0.1.4
FEATURES:
* **New Resource** `azuredevops_serviceendpoint_ssh` [#270](https://github.com/microsoft/terraform-provider-azuredevops/issues/270)
* **New Resource** `azuredevops_serviceendpoint_npm` [#334](https://github.com/microsoft/terraform-provider-azuredevops/issues/334)
* **New Resource** `azuredevops_serviceendpoint_azuredevops` [#339](https://github.com/microsoft/terraform-provider-azuredevops/issues/339)
* **New Resource** `azuredevops_serviceendpoint_github_enterprise` [#210](https://github.com/microsoft/terraform-provider-azuredevops/issues/210)

IMPROVEMENTS:
* `azuredevops_group` - Support for changing group display names [#356](https://github.com/microsoft/terraform-provider-azuredevops/issues/356)
  
BUG FIX:
  `azuredevops_group` - `scope` will be suppressed during `plan` and `apply`  [#345](https://github.com/microsoft/terraform-provider-azuredevops/issues/345)
  `azuredevops_variable_group` - handle non-existent variable groups [#359](https://github.com/microsoft/terraform-provider-azuredevops/issues/359)

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
