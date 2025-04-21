## 1.9.0

FEATURES:

* **New Resource** `azuredevops_extension` [#1328](https://github.com/microsoft/terraform-provider-azuredevops/pull/1328)
* **New Resource** `azuredevops_serviceendpoint_openshift` [#1336](https://github.com/microsoft/terraform-provider-azuredevops/pull/1336)
* **New Data Resource** `azuredevops_git_repository_file` [#1335](https://github.com/microsoft/terraform-provider-azuredevops/pull/1335)

IMPROVEMENTS:

* `azuredevops_git_repository_file` - Add support for git author and committer. [#1340](https://github.com/microsoft/terraform-provider-azuredevops/pull/1340)
* `azuredevops_serviceendpoint_github` - Add support for oauth2. [#1353](https://github.com/microsoft/terraform-provider-azuredevops/pull/1353)
* `azuredevops_git_repository` - Set the branch wait timeout to creation timeout, customizable via `timeouts` in HCL. [#1356](https://github.com/microsoft/terraform-provider-azuredevops/pull/1356)
* `azuredevops_feed_permission`
  - Add support for import. [#1339](https://github.com/microsoft/terraform-provider-azuredevops/pull/1339)
  - Optimize error message. [#1350](https://github.com/microsoft/terraform-provider-azuredevops/pull/1350)
* `azuredevops_git_repository` - Add support for ephemeral password. [#1343](https://github.com/microsoft/terraform-provider-azuredevops/pull/1343)
  > **WARNING**: [To use write-only arguments, you must use Terraform v.1.11 or later and use a resource that supports write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral/write-only#requirements). 
  
* Document - Fixed typo (manage to managed) and updated Azure AD to Entra ID. [#1341](https://github.com/microsoft/terraform-provider-azuredevops/pull/1341)
* Bump `azcore` to `v1.17.1` [#1330](https://github.com/microsoft/terraform-provider-azuredevops/pull/1330)
* Bump `github.com/golang-jwt/jwt/v5` from `v5.2.1` to `v5.2.2` [#1332](https://github.com/microsoft/terraform-provider-azuredevops/pull/1332)
* Replace `golang/mock` with `uber-go/mock` [#1333](https://github.com/microsoft/terraform-provider-azuredevops/pull/1333)
* Bump golang to `v1.24.1` and update CI images [#1334](https://github.com/microsoft/terraform-provider-azuredevops/pull/1334)
* Add `terrafmt` to CI. [#1348](https://github.com/microsoft/terraform-provider-azuredevops/pull/1348)
* Bump `golang.org/x/net` from `v0.37.0` to `v0.38.0`. [#1354](https://github.com/microsoft/terraform-provider-azuredevops/pull/1354)

## 1.8.1

BUG FIX:

* `azuredevops_variable_group` - Fix validation conflict with `ignore_changes`. [#1325](https://github.com/microsoft/terraform-provider-azuredevops/pull/1325)

## 1.8.0

FEATURES:

* **New Data Resource** `azuredevops_group_membership` [#1307](https://github.com/microsoft/terraform-provider-azuredevops/pull/1307)

BUG FIX:

* `azuredevops_project` - Fix the bug where `id` was set to the project name. [#1316](https://github.com/microsoft/terraform-provider-azuredevops/pull/1316)
* All service connection resources - Check if service connection has been deleted. [#1318](https://github.com/microsoft/terraform-provider-azuredevops/pull/1318)

IMPROVEMENTS:

* `azuredevops_client_config` - Add support for export organization ID. [#1301](https://github.com/microsoft/terraform-provider-azuredevops/pull/1301)
* `azuredevops_service_principal_entitlement` - Suppress case difference for `origin`. [#1303](https://github.com/microsoft/terraform-provider-azuredevops/pull/1303)
* Bump `terraform-plugin-sdk/v2` sdk to `v2.36.1` and `azidentity` to `v1.8.2`. [#1310](https://github.com/microsoft/terraform-provider-azuredevops/pull/1310)
* `azuredevops_build_definition` 
  - Add support for other Git(`Git`) to `repository.repo_type`. [#1312](https://github.com/microsoft/terraform-provider-azuredevops/pull/1312)
  - Add support for create classic agent jobs. [#1312](https://github.com/microsoft/terraform-provider-azuredevops/pull/1312)
* **Data source** `azuredevops_build_definition`
  - Add support for exporting other Git(`Git`). [#1312](https://github.com/microsoft/terraform-provider-azuredevops/pull/1312)
  - Add support for exporting classic agent jobs. [#1312](https://github.com/microsoft/terraform-provider-azuredevops/pull/1312)
* `azuredevops_variable_group` - Check secret variables during plan. [#1321](https://github.com/microsoft/terraform-provider-azuredevops/pull/1321)
* `azuredevops_users` - Update documentation. [#1302](https://github.com/microsoft/terraform-provider-azuredevops/pull/1302)
* `azuredevops_team` - Fix incorrect name in Terraform registry. [#1304](https://github.com/microsoft/terraform-provider-azuredevops/pull/1304)
* `azuredevops_teams` - Fix incorrect name in Terraform registry. [#1304](https://github.com/microsoft/terraform-provider-azuredevops/pull/1304)
* `serviceendpoint_azurecr` - Documentation update. [#1308](https://github.com/microsoft/terraform-provider-azuredevops/pull/1308)

## 1.7.0

FEATURES:

* **New Resource** `azuredevops_dashboard` [#1284](https://github.com/microsoft/terraform-provider-azuredevops/pull/1284)
* **New Data Resource** `azuredevops_descriptor` [#1294](https://github.com/microsoft/terraform-provider-azuredevops/pull/1294)
* **New Data Resource** `azuredevops_storage_key` [#1294](https://github.com/microsoft/terraform-provider-azuredevops/pull/1294)
* **New Data Resource** `azuredevops_user` [#1296](https://github.com/microsoft/terraform-provider-azuredevops/pull/1296)

BUG FIX:

* `azuredevops_project` - Fix name unchanged but updated.  [#1285](https://github.com/microsoft/terraform-provider-azuredevops/pull/1285)
* Permission resources 
  - Fix `descriptor` filter bug, cannot set permission for AAD groups. [#1297](https://github.com/microsoft/terraform-provider-azuredevops/pull/1297)
  - Fix collection level groups/users cannot set permission bug. [#1299](https://github.com/microsoft/terraform-provider-azuredevops/pull/1299)

IMPROVEMENTS:

* `azuredevops_identity_groups`  
  - Add support for `descriptor`. [#1279](https://github.com/microsoft/terraform-provider-azuredevops/pull/1279)
  - Add support for `subject_descriptor`. [#1292](https://github.com/microsoft/terraform-provider-azuredevops/pull/1292)
* `azuredevops_identity_group` - Add support for `subject_descriptor`. [#1292](https://github.com/microsoft/terraform-provider-azuredevops/pull/1292)
* `azuredevops_identity_user` - Add support for `subject_descriptor`. [#1293](https://github.com/microsoft/terraform-provider-azuredevops/pull/1293)

BREAKING CHANGE:

* All service endpoint resources - Change `authorization` to compute only, not configurable. [#1298](https://github.com/microsoft/terraform-provider-azuredevops/pull/1298)

## 1.6.0

FEATURES:

* **New Resource** `azuredevops_service_principal_entitlement` [#1253](https://github.com/microsoft/terraform-provider-azuredevops/pull/1253)
* **New Resource** `azuredevops_feed_retention_policy` [#1257](https://github.com/microsoft/terraform-provider-azuredevops/pull/1257)
* **New Resource** `azuredevops_project_tags` [#1259](https://github.com/microsoft/terraform-provider-azuredevops/pull/1259)
* **New Resource** `azuredevops_serviceendpoint_checkmarx_sca` [#1267](https://github.com/microsoft/terraform-provider-azuredevops/pull/1267)
* **New Resource** `azuredevops_serviceendpoint_checkmarx_sast` [#1268](https://github.com/microsoft/terraform-provider-azuredevops/pull/1268)
* **New Resource** `azuredevops_serviceendpoint_checkmarx_one` [#1269](https://github.com/microsoft/terraform-provider-azuredevops/pull/1269)
* **New Resource** `azuredevops_check_rest_api` [#1274](https://github.com/microsoft/terraform-provider-azuredevops/pull/1274)
* **New Data Resource** `azuredevops_service_principal` [#1253](https://github.com/microsoft/terraform-provider-azuredevops/pull/1253)

BUG FIX:

* `azuredevops_securityrole_assignment` - Fix inconsistent result after apply.  [#1255](https://github.com/microsoft/terraform-provider-azuredevops/pull/1255)
* `azuredevops_wiki` - Fix documentation typos.  [#1264](https://github.com/microsoft/terraform-provider-azuredevops/pull/1264)
* `azuredevops_git_repository` - Fix branch not found bug.  [#1270](https://github.com/microsoft/terraform-provider-azuredevops/pull/1270)
* Permission resources - Add support for identity filtering.  [#1256](https://github.com/microsoft/terraform-provider-azuredevops/pull/1256)

IMPROVEMENTS:

* `azuredevops_project` - Update documentation. [#1258](https://github.com/microsoft/terraform-provider-azuredevops/pull/1258)
* `azuredevops_feed_retention_policy` - Add support for organization level feed retention policy. [#1261](https://github.com/microsoft/terraform-provider-azuredevops/pull/1261)
* **Data Resource** `azuredevops_team` - Optimize the read operation, use `GetTeam` instead of `GetTeams` [#1262](https://github.com/microsoft/terraform-provider-azuredevops/pull/1262)
* All resource documentation - Add timeout documentation. [#1273](https://github.com/microsoft/terraform-provider-azuredevops/pull/1273)
* Update dependencies and bump go to `v1.23` [#1277](https://github.com/microsoft/terraform-provider-azuredevops/pull/1277)
* Documentation
  - Update document format [#1278](https://github.com/microsoft/terraform-provider-azuredevops/pull/1278)
  - Fix documentation errors and add missing properties. [#1280](https://github.com/microsoft/terraform-provider-azuredevops/pull/1280)

## 1.5.0

FEATURES:

* **New Resource** `azuredevops_serviceendpoint_snyk` [#1224](https://github.com/microsoft/terraform-provider-azuredevops/pull/1224)
* **New Resource** `azuredevops_serviceendpoint_dynamics_lifecycle_services` [#1240](https://github.com/microsoft/terraform-provider-azuredevops/pull/1240)
* **New Resource** `azuredevops_serviceendpoint_azure_service_bus` [#1242](https://github.com/microsoft/terraform-provider-azuredevops/pull/1242)
* **New Resource** `azuredevops_serviceendpoint_gitlab` [#1243](https://github.com/microsoft/terraform-provider-azuredevops/pull/1243)
* **New Resource** `azuredevops_serviceendpoint_visualstudiomarketplace` [#1246](https://github.com/microsoft/terraform-provider-azuredevops/pull/1246)
* **New Data Resource** `azuredevops_serviceendpoint_bitbucket` [#1200](https://github.com/microsoft/terraform-provider-azuredevops/pull/1200)

BUG FIX:

* `azuredevops_serviceendpoint_github_enterprise` - Add `nil` check.  [#1209](https://github.com/microsoft/terraform-provider-azuredevops/pull/1209)
* `azuredevops_serviceendpoint_generic` - Relax `server_url` restrictions.  [#1210](https://github.com/microsoft/terraform-provider-azuredevops/pull/1210)
* All service connection resources - Fix import share service connection not point to the right project.  [#1211](https://github.com/microsoft/terraform-provider-azuredevops/pull/1211)
* `azuredevops_group_entitlement` 
  - Detect group deleted.  [#1212](https://github.com/microsoft/terraform-provider-azuredevops/pull/1212)  
  - Fix group import crash bug.  [#1220](https://github.com/microsoft/terraform-provider-azuredevops/pull/1220)
* `azuredevops_check_branch_control` - Remove the required check for `ignore_unknown_protection_status`. [#1222](https://github.com/microsoft/terraform-provider-azuredevops/pull/1222)
* `azuredevops_serviceendpoint_kubernetes` - Fix crash bug. [#1228](https://github.com/microsoft/terraform-provider-azuredevops/pull/1228)

IMPROVEMENTS:

* SDK update - Update `resource.StateChangeConf to` `retry.StateChangeConf` [#1204](https://github.com/microsoft/terraform-provider-azuredevops/pull/1204)
* `azuredevops_securityrole_assignment` - Change `resource_id` to `forceNew=true`  [#1205](https://github.com/microsoft/terraform-provider-azuredevops/pull/1205)
* Add client initialization error handle  [#1207](https://github.com/microsoft/terraform-provider-azuredevops/pull/1207)
* `azuredevops_user_entitlement` - Update documentation  [#1208](https://github.com/microsoft/terraform-provider-azuredevops/pull/1208)
* `azuredevops_serviceendpoint_azurerm`
  - Add support `server_url` and cloud environment`AzureStack` [#1213](https://github.com/microsoft/terraform-provider-azuredevops/pull/1213) 
  - Add support for `credentials.serviceprincipalcertificate`[#1225](https://github.com/microsoft/terraform-provider-azuredevops/pull/1225)
  - Add support for `credentials.serviceprincipalcertificate`[#1225](https://github.com/microsoft/terraform-provider-azuredevops/pull/1225)
* `azuredevops_git_repository` 
  - Add support for initialize of uninitialized repository [#1218](https://github.com/microsoft/terraform-provider-azuredevops/pull/1218)
  - Update document [#1221](https://github.com/microsoft/terraform-provider-azuredevops/pull/1221)
  - Support importing repository via username/password [#1223](https://github.com/microsoft/terraform-provider-azuredevops/pull/1223)
* `azuredevops_build_definition` - Add support for `build_completion_trigger` [#1226](https://github.com/microsoft/terraform-provider-azuredevops/pull/1226)
* `azuredevops_serviceendpoint_kubernetes` - Add support for `service_account.accept_untrusted_certs` [#1229](https://github.com/microsoft/terraform-provider-azuredevops/pull/1229)
* All service connections - Remove `forceNew` for `service_endpoint_name` [#1238](https://github.com/microsoft/terraform-provider-azuredevops/pull/1238)
* `azuredevops_serviceendpoint_aws` - Add `nil` check in resource read [#1239](https://github.com/microsoft/terraform-provider-azuredevops/pull/1239)
* `azuredevops_serviceendpoint_azurecr` - Change `serviceprincipalid` to `forceNew=true` [#1247](https://github.com/microsoft/terraform-provider-azuredevops/pull/1247)
* go.mod - Bump `golang.org/x/crypto` from `v0.24.0` to `v0.31.0` [#1252](https://github.com/microsoft/terraform-provider-azuredevops/pull/1252)

BREAKING CHANGE:

* `azuredevops_build_definition` - Change `name` from optional to required. [#1185](https://github.com/microsoft/terraform-provider-azuredevops/pull/1185)


## 1.4.0

FEATURES:

* **New Data Resource** `azuredevops_serviceendpoint_bitbucket` [#1200](https://github.com/microsoft/terraform-provider-azuredevops/pull/1200)

BUG FIX:

* `azuredevops_agent_queue` - Fix `name` not set bug.  [#1157](https://github.com/microsoft/terraform-provider-azuredevops/pull/1157)
* `azuredevops_serviceendpoint_sonarqube` - Adding nil check to project ID. [#1159](https://github.com/microsoft/terraform-provider-azuredevops/pull/1159)
* `azuredevops_group` - Detect that group has been deleted. [#1196](https://github.com/microsoft/terraform-provider-azuredevops/pull/1196)
* All service connection - Detect that service connection is not fully returned and this appears to be a permission issue. [#1193](https://github.com/microsoft/terraform-provider-azuredevops/pull/1193)

IMPROVEMENTS:

* `azuredevops_wiki` - Add support for delete project type wiki [#1166](https://github.com/microsoft/terraform-provider-azuredevops/pull/1166)
* `azuredevops_agent_queue` - Add `name` validation [#1184](https://github.com/microsoft/terraform-provider-azuredevops/pull/1184)
* **Data Source** `azuredevops_agent_queue` - Add `name` validation [#1184](https://github.com/microsoft/terraform-provider-azuredevops/pull/1184)
* `azuredevops_git_repository` 
  - Add support for enable/disable repository [#1181](https://github.com/microsoft/terraform-provider-azuredevops/pull/1181)
  - Update test case [#1188](https://github.com/microsoft/terraform-provider-azuredevops/pull/1188)
  - Optimize resource import [#1194](https://github.com/microsoft/terraform-provider-azuredevops/pull/1194)
* **Data Source** `azuredevops_git_repository` - Optimize resource acquisition. [#1197](https://github.com/microsoft/terraform-provider-azuredevops/pull/1197)
* `azuredevops_repository_policy_max_file_size` - Add support for max file size `50M` [#1168](https://github.com/microsoft/terraform-provider-azuredevops/pull/1168)
* `azuredevops_feed_permission` - Sync permissions after create/update [#1169](https://github.com/microsoft/terraform-provider-azuredevops/pull/1169)
* `azuredevops_branch_policy_build_validation` - Update document [#1172](https://github.com/microsoft/terraform-provider-azuredevops/pull/1172)
* `serviceendpoint_azurecr` - Fix document error [#1163](https://github.com/microsoft/terraform-provider-azuredevops/pull/1163)
* `azuredevops_build_definition_permissions` - Update document [#1195](https://github.com/microsoft/terraform-provider-azuredevops/pull/1195)

BREAKING CHANGE:

* `azuredevops_build_definition` - Change `name` from optional to required. [#1185](https://github.com/microsoft/terraform-provider-azuredevops/pull/1185)

## 1.3.0

BUG FIX:

* `azuredevops_serviceendpoint_azurecr`
  - Fix `tenant_id` not set as expected.  [#1115](https://github.com/microsoft/terraform-provider-azuredevops/pull/1115)
  - Fix `tenant_id` not set bug.  [#1142](https://github.com/microsoft/terraform-provider-azuredevops/pull/1142)
* **Data Source** `azuredevops_users` - Return empty list if user not found.  [#1116](https://github.com/microsoft/terraform-provider-azuredevops/pull/1116)
* `azuredevops_securityrole_assignment` - Detecting role assignment revoke.  [#1120](https://github.com/microsoft/terraform-provider-azuredevops/pull/1120)
* `azuredevops_serviceendpoint_kubernetes` - Enhance `nil` check.  [#1127](https://github.com/microsoft/terraform-provider-azuredevops/pull/1127)
* `azuredevops_team` - Fix idempotency add members issue.  [#1130](https://github.com/microsoft/terraform-provider-azuredevops/pull/1130)
* `azuredevops_serviceendpoint_azurecr` - Expect `serviceprincipalkey` only if ServicePrincipal authentication is used. [#1134](https://github.com/microsoft/terraform-provider-azuredevops/pull/1134)
* `azuredevops_build_folder` - Fix import bug. [#1143](https://github.com/microsoft/terraform-provider-azuredevops/pull/1143)
* `azuredevops_serviceendpoint_dockerregistry` - Enhance `nil` check. [#1146](https://github.com/microsoft/terraform-provider-azuredevops/pull/1146)
* `azuredevops_group` - Add support for `group_id`. [#1147](https://github.com/microsoft/terraform-provider-azuredevops/pull/1147)
* **Data Source** `azuredevops_group` - Add support for `group_id`. [#1149](https://github.com/microsoft/terraform-provider-azuredevops/pull/1149)

IMPROVEMENTS:

* `azuredevops_feed` Support import [#1119](https://github.com/microsoft/terraform-provider-azuredevops/pull/1119)
* Add default timeout [#1114](https://github.com/microsoft/terraform-provider-azuredevops/pull/1114)
* Update Task Agent resources  [#1128](https://github.com/microsoft/terraform-provider-azuredevops/pull/1128)

## 1.2.0

FEATURES:

* **New Resource** `azuredevops_wiki` [#1032](https://github.com/microsoft/terraform-provider-azuredevops/pull/1032)

BUG FIX:

* `azuredevops_check_exclusive_lock` - Add example.  [#1054](https://github.com/microsoft/terraform-provider-azuredevops/pull/1054)
* `azuredevops_users` - Fix user not found bug.  [#1110](https://github.com/microsoft/terraform-provider-azuredevops/pull/1110)
* `azuredevops_git_repository`
  - Fix repository not found bug.  [#1065](https://github.com/microsoft/terraform-provider-azuredevops/pull/1065)
  - Detect repository deleted outside of Terraform  [#1087](https://github.com/microsoft/terraform-provider-azuredevops/pull/1087)
* `azuredevops_pipeline_authorization` - Check Pipeline Project for Resource Permissions.  [#1059](https://github.com/microsoft/terraform-provider-azuredevops/pull/1059)
* `azuredevops_serviceendpoint_kubernetes` - Enhance parameter `nil` checking.  [#1091](https://github.com/microsoft/terraform-provider-azuredevops/pull/1091)
* `azuredevops_git_repository_file` - Check branch status.  [#1100](https://github.com/microsoft/terraform-provider-azuredevops/pull/1100)

IMPROVEMENTS:

* `azuredevops_serviceendpoint_azurerm` - Add support for `AzureUSGovernment` and `AzureGermanCloud` clouds. [#1061](https://github.com/microsoft/terraform-provider-azuredevops/pull/1061)
* `azuredevops_variable_group`
  - Add validation that variable can have either only `value` attribute or both `is_secret` and `secret_value` attributes. [#1075](https://github.com/microsoft/terraform-provider-azuredevops/pull/1075)
  - Update document.  [#1044](https://github.com/microsoft/terraform-provider-azuredevops/pull/1044)
* `azuredevops_serviceendpoint_azurecr` - Add support for `WorkloadIdentityFederation`.  [#1105](https://github.com/microsoft/terraform-provider-azuredevops/pull/1105)
* `azuredevops_git_repository` - Fix typo error.  [#1111](https://github.com/microsoft/terraform-provider-azuredevops/pull/1111)
* Dependency upgrade -  [#1083](https://github.com/microsoft/terraform-provider-azuredevops/pull/1083)
* `azuredevops_check_approval` - Update tests  [#1092](https://github.com/microsoft/terraform-provider-azuredevops/pull/1092)
* `azuredevops_check_exclusive_lock` - Add default timeout and update tests.  [#1092](https://github.com/microsoft/terraform-provider-azuredevops/pull/1092)
* `azuredevops_check_branch_control` - Add default timeout and update tests.  [#1092](https://github.com/microsoft/terraform-provider-azuredevops/pull/1092)
* `azuredevops_check_business_hours` - Add default timeout and update tests.  [#1092](https://github.com/microsoft/terraform-provider-azuredevops/pull/1092)
* `azuredevops_check_required_template` - Add default timeout and update tests.  [#1092](https://github.com/microsoft/terraform-provider-azuredevops/pull/1092)
* `azuredevops_build_definition` - Update unit tests.  [#1094](https://github.com/microsoft/terraform-provider-azuredevops/pull/1094)
* **Data Source** `azuredevops_build_definition` - Update tests.  [#1094](https://github.com/microsoft/terraform-provider-azuredevops/pull/1094)
* `azuredevops_build_folder` - Add default timeout and update tests.  [#1094](https://github.com/microsoft/terraform-provider-azuredevops/pull/1094)
* `azuredevops_pipeline_authorization` - Add update tests.  [#1094](https://github.com/microsoft/terraform-provider-azuredevops/pull/1094)
* `azuredevops_resource_authorization ` - Add update tests.  [#1094](https://github.com/microsoft/terraform-provider-azuredevops/pull/1094)
* **Data Source** `azuredevops_project` - Add default timeout and update tests.  [#1095](https://github.com/microsoft/terraform-provider-azuredevops/pull/1095)
* **Data Source** `azuredevops_projects` - Update unit tests.  [#1095](https://github.com/microsoft/terraform-provider-azuredevops/pull/1095)
* **Data Source** `azuredevops_team` - Update unit tests.  [#1095](https://github.com/microsoft/terraform-provider-azuredevops/pull/1095)
* **Data Source** `azuredevops_teams` - Update unit tests.  [#1095](https://github.com/microsoft/terraform-provider-azuredevops/pull/1095)
* `azuredevops_project ` - Add update unit tests.  [#1095](https://github.com/microsoft/terraform-provider-azuredevops/pull/1095)
* `azuredevops_project_features ` - Add update unit tests.  [#1095](https://github.com/microsoft/terraform-provider-azuredevops/pull/1095)
* `azuredevops_project_pipeline_settings ` - Add update unit tests.  [#1095](https://github.com/microsoft/terraform-provider-azuredevops/pull/1095)
* `azuredevops_team ` - Add update unit tests.  [#1095](https://github.com/microsoft/terraform-provider-azuredevops/pull/1095)
* `azuredevops_team_members ` - Add update unit tests.  [#1095](https://github.com/microsoft/terraform-provider-azuredevops/pull/1095)
* `azuredevops_team_administrators ` - Add update unit tests.  [#1095](https://github.com/microsoft/terraform-provider-azuredevops/pull/1095)
* `azuredevops_feed` - Optimize code and update tests.  [#1098](https://github.com/microsoft/terraform-provider-azuredevops/pull/1098)
* **Data Source** `azuredevops_feed` - Optimize code and update tests.  [#1098](https://github.com/microsoft/terraform-provider-azuredevops/pull/1098)
* `azuredevops_feed_permission` - Optimize code and update tests.  [#1098](https://github.com/microsoft/terraform-provider-azuredevops/pull/1098)
* **Data Source** `azuredevops_git_repositories` - Add timeout and update tests.  [#1103](https://github.com/microsoft/terraform-provider-azuredevops/pull/1103)
* **Data Source** `azuredevops_git_repository` - Add timeout and update tests.  [#1103](https://github.com/microsoft/terraform-provider-azuredevops/pull/1103)
* `azuredevops_git_repository` - Add timeout and update tests.  [#1103](https://github.com/microsoft/terraform-provider-azuredevops/pull/1103)
* `azuredevops_git_repository_branch` - Add timeout and update tests.  [#1103](https://github.com/microsoft/terraform-provider-azuredevops/pull/1103)
* `azuredevops_git_repository_file` - Add timeout and update tests.  [#1103](https://github.com/microsoft/terraform-provider-azuredevops/pull/1103)
* `azuredevops_group_membership` - Add timeout and update tests.  [#1107](https://github.com/microsoft/terraform-provider-azuredevops/pull/1107)
* `azuredevops_group` - Add timeout and update tests.  [#1107](https://github.com/microsoft/terraform-provider-azuredevops/pull/1107)
* **Data Source** `azuredevops_users` - Add timeout and update tests.  [#1107](https://github.com/microsoft/terraform-provider-azuredevops/pull/1107)
* **Data Source** `azuredevops_group` - Add timeout and update tests.  [#1107](https://github.com/microsoft/terraform-provider-azuredevops/pull/1107)
* **Data Source** `azuredevops_groups` - Add timeout and update tests.  [#1107](https://github.com/microsoft/terraform-provider-azuredevops/pull/1107)
* **Data Source** `azuredevops_identity_user` - Add timeout and update tests.  [#1108](https://github.com/microsoft/terraform-provider-azuredevops/pull/1108)
* **Data Source** `azuredevops_identity_group` - Add timeout and update tests.  [#1108](https://github.com/microsoft/terraform-provider-azuredevops/pull/1108)
* **Data Source** `azuredevops_identity_groups` - Add timeout and update tests.  [#1108](https://github.com/microsoft/terraform-provider-azuredevops/pull/1108)
* `azuredevops_user_entitlement` - Add timeout and update tests.  [#1109](https://github.com/microsoft/terraform-provider-azuredevops/pull/1109)
* `azuredevops_group_entitlement` - Add timeout and update tests.  [#1109](https://github.com/microsoft/terraform-provider-azuredevops/pull/1109)


## 1.1.1

BUG FIX:

* `azuredevops_team_members` - Optimize `descriptor` read operation.  [#1048](https://github.com/microsoft/terraform-provider-azuredevops/pull/1048)
* `azuredevops_team` - Optimize `descriptor` read operation.  [#1048](https://github.com/microsoft/terraform-provider-azuredevops/pull/1048)
* `azuredevops_teams` - Optimize `descriptor` read operation.   [#1048](https://github.com/microsoft/terraform-provider-azuredevops/pull/1048)

## 1.1.0

FEATURES:

* **New Data Resource** `azuredevops_identity_user` [#956](https://github.com/microsoft/terraform-provider-azuredevops/pull/956)
* **New Data Resource** `azuredevops_identity_group` [#956](https://github.com/microsoft/terraform-provider-azuredevops/pull/956)
* **New Data Resource** `azuredevops_identity_groups` [#956](https://github.com/microsoft/terraform-provider-azuredevops/pull/956)
* **New Resource** `azuredevops_securityrole_assignment` [#982](https://github.com/microsoft/terraform-provider-azuredevops/pull/982)
* **New Data Resource** `azuredevops_securityrole_definitions` [#982](https://github.com/microsoft/terraform-provider-azuredevops/pull/982)
* **New Resource** `azuredevops_feed` [#1011](https://github.com/microsoft/terraform-provider-azuredevops/pull/1011)
* **New Resource** `azuredevops_feed_permission` [#1011](https://github.com/microsoft/terraform-provider-azuredevops/pull/1011)
* **New Data Resource** `azuredevops_feed` [#1011](https://github.com/microsoft/terraform-provider-azuredevops/pull/1011)

IMPROVEMENTS:

* `azuredevops_pipeline_authorization` - Allow pipeline authorization across projects. [#973](https://github.com/microsoft/terraform-provider-azuredevops/pull/973) 
* `azuredevops_git_repository` - Support export repository status. [#1024](https://github.com/microsoft/terraform-provider-azuredevops/pull/1024)
* **Data Resource** `azuredevops_git_repositories` - Support export repository status. [#1024](https://github.com/microsoft/terraform-provider-azuredevops/pull/1024)
* **Data Resource** `azuredevops_git_repository` - Support export repository status. [#1024](https://github.com/microsoft/terraform-provider-azuredevops/pull/1024)
* **Document** `azuredevops_elastic_pool` - Fix document title. [#1037](https://github.com/microsoft/terraform-provider-azuredevops/pull/1037)
* **Document** Adding information about use in Azure Pipelines. [#1019](https://github.com/microsoft/terraform-provider-azuredevops/pull/1019)

BUG FIX:

* `azuredevops_serviceendpoint_azurerm` - Fix `azurerm_subscription_id` conflicts with `azurerm_management_group_id`.  [#1004](https://github.com/microsoft/terraform-provider-azuredevops/pull/1004) 
* `azuredevops_team_members` - Optimize `descriptor` read operation.  [#1014](https://github.com/microsoft/terraform-provider-azuredevops/pull/1014)
* `azuredevops_team` - Optimize `descriptor` read operation.  [#1014](https://github.com/microsoft/terraform-provider-azuredevops/pull/1014)
* `azuredevops_teams` - Optimize `descriptor` read operation.   [#1014](https://github.com/microsoft/terraform-provider-azuredevops/pull/1014)
* `azuredevops_group_membership` - Fix group entitlement not found error.   [#1015](https://github.com/microsoft/terraform-provider-azuredevops/pull/1015)
* `azuredevops_git_repository` - Fix cannot set `default_branch` on update.   [#1020](https://github.com/microsoft/terraform-provider-azuredevops/pull/1020)

## 1.0.1

FEATURES:

* Fix AzureAD authorization and OIDC validationOIDC validation [#993](https://github.com/microsoft/terraform-provider-azuredevops/pull/993)

## 1.0.0 


FEATURES:

* **New Resource** `azuredevops_environment_resource_kubernetes` [#935](https://github.com/microsoft/terraform-provider-azuredevops/pull/935)
* **New Resource** `azuredevops_library_permissions` [#740](https://github.com/microsoft/terraform-provider-azuredevops/pull/740)
* **New Resource** `azuredevops_variable_group_permissions` [#740](https://github.com/microsoft/terraform-provider-azuredevops/pull/740)
* Add support for Service Principal, Identity, OIDC etc. authorization [#747](https://github.com/microsoft/terraform-provider-azuredevops/pull/747)

IMPROVEMENTS:

* `azuredevops_check_required_template` - Add support for `githubenterprise` repository type. [#962](https://github.com/microsoft/terraform-provider-azuredevops/pull/962)
* `azuredevops_elastic_pool` - Add support for `project_id`. [#966](https://github.com/microsoft/terraform-provider-azuredevops/pull/966)
* `azuredevops_pipeline_authorization` - Update document. [#960](https://github.com/microsoft/terraform-provider-azuredevops/pull/960)
* **Data Resource**`azuredevops_groups ` - Add support for group `id`. [#980](https://github.com/microsoft/terraform-provider-azuredevops/pull/980)

BUG FIX:

* `azuredevops_serviceendpoint_azurecr` - Fix potential nil exception.  [#972](https://github.com/microsoft/terraform-provider-azuredevops/pull/972)
* `azuredevops_serviceendpoint_azurerm` - Fix import error.  [#967](https://github.com/microsoft/terraform-provider-azuredevops/pull/967)
* `azuredevops_variable_group` - Exclude Key Vault disabled secrets.  [#947](https://github.com/microsoft/terraform-provider-azuredevops/pull/947)
* `azuredevops_git_repository` - Fix default branch not set when `init_type=Clean` or `init_type=Fork`.  [#946](https://github.com/microsoft/terraform-provider-azuredevops/pull/946)
* `azuredevops_check_approval` - Add missing `version` property.  [#977](https://github.com/microsoft/terraform-provider-azuredevops/pull/977)
* `azuredevops_check_branch_control` - Add missing `version` property.  [#977](https://github.com/microsoft/terraform-provider-azuredevops/pull/977)
* `azuredevops_check_business_hours` - Add missing `version` property.  [#977](https://github.com/microsoft/terraform-provider-azuredevops/pull/977)
* `azuredevops_check_exclusive_lock` - Add missing `version` property.  [#977](https://github.com/microsoft/terraform-provider-azuredevops/pull/977)
* `azuredevops_check_required_template` - Add missing `version` property.  [#977](https://github.com/microsoft/terraform-provider-azuredevops/pull/977)
* `azuredevops_pipeline_authorization` - Fix pipeline authorization not set.  [#986](https://github.com/microsoft/terraform-provider-azuredevops/pull/986)


## 0.11.0

FEATURES:

* **New Resource** `azuredevops_servicehook_storage_queue_pipelines` [#914](https://github.com/microsoft/terraform-provider-azuredevops/pull/914)

IMPROVEMENTS:

* `azuredevops_serviceendpoint_azurerm` - Add support for `featuure` to verify the service connection. [#865](https://github.com/microsoft/terraform-provider-azuredevops/pull/865)
* `azuredevops_build_definition` - Add support for `queue_status`. [#916](https://github.com/microsoft/terraform-provider-azuredevops/pull/916)
* `azuredevops_pipeline_authorization` - Enhance authorization status check. [#929](https://github.com/microsoft/terraform-provider-azuredevops/pull/929)
* `azuredevops_agent_queue` - Add support for `name`. [#906](https://github.com/microsoft/terraform-provider-azuredevops/pull/906)
* `azuredevops_users` - Improve read operation performance. [#939](https://github.com/microsoft/terraform-provider-azuredevops/pull/939)
* **Data Resource** `azuredevops_environment` - Add support for fetch environment by name. [#917](https://github.com/microsoft/terraform-provider-azuredevops/pull/917)

BUG FIX:

* `azuredevops_serviceendpoint_azurerm` - Fix resource deleted but state not removed.  [#921](https://github.com/microsoft/terraform-provider-azuredevops/pull/921)
* `azuredevops_git_repository_file` - Fix apply for non-project resources fails.  [#925](https://github.com/microsoft/terraform-provider-azuredevops/pull/925)
* `azuredevops_build_definition` - Fix `skip_first_run` to work for all repo types. [#928](https://github.com/microsoft/terraform-provider-azuredevops/pull/928)
* `azuredevops_git_repository` - Fix `default_branch` not set when `init_type=Clean` or `init_type=Fork`. [#946](https://github.com/microsoft/terraform-provider-azuredevops/pull/946)
* `azuredevops_variable_group` - Exclude disabled secrets. [#947](https://github.com/microsoft/terraform-provider-azuredevops/pull/947)


## 0.10.0

IMPROVEMENTS:

* `azuredevops_pipeline_authorization` - Add support for `repository` authorization  [#883](https://github.com/microsoft/terraform-provider-azuredevops/pull/883) 
* `azuredevops_elastic_pool` - Support set `time_to_live_minutes` to `0` [#885](https://github.com/microsoft/terraform-provider-azuredevops/pull/885)
* `azuredevops_serviceendpoint_azurerm` - Support export `service_principal_id` [#902](https://github.com/microsoft/terraform-provider-azuredevops/pull/902)
* `azuredevops_area_permissions` - Update document [#909](https://github.com/microsoft/terraform-provider-azuredevops/pull/909)

BUG FIX:

* `azuredevops_build_folder_permissions` - Fix root folder permissions for builds not set [#893](https://github.com/microsoft/terraform-provider-azuredevops/pull/893)
* `azuredevops_project_pipeline_settings` - Fix `enforce_referenced_repo_scoped_token` not set [#898](https://github.com/microsoft/terraform-provider-azuredevops/pull/898)
  


## 0.9.1

FEATURES:

* **New Resource** `azuredevops_group_entitlement` [#870](https://github.com/microsoft/terraform-provider-azuredevops/pull/870)

## 0.9.0

FEATURES:

* **New Resource** `azuredevops_serviceendpoint_nuget` [#866](https://github.com/microsoft/terraform-provider-azuredevops/pull/866)
* **New Data Resource** `azuredevops_serviceendpoint_azurecr` [#867](https://github.com/microsoft/terraform-provider-azuredevops/pull/867)

IMPROVEMENTS:

* `azuredevops_serviceendpoint_azurerm` - Add support for `workload_identity_federation_issuer` and `workload_identity_federation_subject` [#861](https://github.com/microsoft/terraform-provider-azuredevops/pull/861)
* `azuredevops_build_definition` - Add support for `skip_first_run` [#871](https://github.com/microsoft/terraform-provider-azuredevops/pull/871)
* All service connections - Decouple create/read/update/delete from generic functions [#863](https://github.com/microsoft/terraform-provider-azuredevops/pull/863)
* Update API link [#869](https://github.com/microsoft/terraform-provider-azuredevops/pull/869)

## 0.8.0

FEATURES:

* **New Resource** `azuredevops_serviceendpoint_maven` [#617](https://github.com/microsoft/terraform-provider-azuredevops/pull/617)
* **New Resource** `azuredevops_serviceendpoint_jenkins` [#617](https://github.com/microsoft/terraform-provider-azuredevops/pull/617)
* **New Resource** `azuredevops_serviceendpoint_nexus` [#617](https://github.com/microsoft/terraform-provider-azuredevops/pull/617)
* **New Data Resource** `azuredevops_environment` [#838](https://github.com/microsoft/terraform-provider-azuredevops/pull/838)

IMPROVEMENTS:

* `azuredevops_check_branch_control` - Add support for `timeout` [#834](https://github.com/microsoft/terraform-provider-azuredevops/pull/834)
* `azuredevops_check_business_hours` - Add support for `timeout` [#834](https://github.com/microsoft/terraform-provider-azuredevops/pull/834)
* `azuredevops_group ` - Upgrade the API from v5 to v7  [#854](https://github.com/microsoft/terraform-provider-azuredevops/pull/854)
* **Data Resource** `azuredevops_team` - Add support for `top`, custom the number of teams returned [#778](https://github.com/microsoft/terraform-provider-azuredevops/pull/778)
* **Data Resource** `azuredevops_teams` - Add support for `top`, custom the number of teams returned [#778](https://github.com/microsoft/terraform-provider-azuredevops/pull/778)

BUG FIX:

* `azuredevops_git_permissions` - Fix branch name tokenization [#842](https://github.com/microsoft/terraform-provider-azuredevops/pull/842)

BREAKING CHANGE:

Deprecate hash properties, all the hash properties have been removed.
  * `azuredevops_serviceendpoint_aws` [#856](https://github.com/microsoft/terraform-provider-azuredevops/pull/856)
  * `azuredevops_serviceendpoint_azuredevops` [#856](https://github.com/microsoft/terraform-provider-azuredevops/pull/856)
  * `azuredevops_serviceendpoint_azurerm` [#856](https://github.com/microsoft/terraform-provider-azuredevops/pull/856)
  * `azuredevops_serviceendpoint_bitbucket` [#856](https://github.com/microsoft/terraform-provider-azuredevops/pull/856)
  * `azuredevops_serviceendpoint_dockerregistry` [#856](https://github.com/microsoft/terraform-provider-azuredevops/pull/856)
  * `azuredevops_serviceendpoint_generic` [#856](https://github.com/microsoft/terraform-provider-azuredevops/pull/856)
  * `azuredevops_serviceendpoint_generic_git` [#856](https://github.com/microsoft/terraform-provider-azuredevops/pull/856)
  * `azuredevops_serviceendpoint_github` [#856](https://github.com/microsoft/terraform-provider-azuredevops/pull/856)
  * `azuredevops_serviceendpoint_github_enterprise` [#856](https://github.com/microsoft/terraform-provider-azuredevops/pull/856)
  * `azuredevops_serviceendpoint_incomingwebhook` [#856](https://github.com/microsoft/terraform-provider-azuredevops/pull/856)
  * `azuredevops_serviceendpoint_kubernetes` [#856](https://github.com/microsoft/terraform-provider-azuredevops/pull/856)
  * `azuredevops_serviceendpoint_runpipeline` [#856](https://github.com/microsoft/terraform-provider-azuredevops/pull/856)
  * `azuredevops_serviceendpoint_sonarqube` [#856](https://github.com/microsoft/terraform-provider-azuredevops/pull/856)
  * `azuredevops_serviceendpoint_ssh` [#856](https://github.com/microsoft/terraform-provider-azuredevops/pull/856)

## 0.7.0

FEATURES:
* **New Resource** `azuredevops_elastic_pool  ` [#835](https://github.com/microsoft/terraform-provider-azuredevops/pull/835)
* **New Resource** `azuredevops_check_exclusive_lock` [#832](https://github.com/microsoft/terraform-provider-azuredevops/pull/832)
* **New Resource** `azuredevops_check_required_template` [#806](https://github.com/microsoft/terraform-provider-azuredevops/pull/806)

IMPROVEMENTS:
* `azuredevops_project` - Get description from service connection not project [#824](https://github.com/microsoft/terraform-provider-azuredevops/pull/824)
* `azuredevops_check_business_hours` - Resolved description for business hours check. [#831](https://github.com/microsoft/terraform-provider-azuredevops/pull/831)
* `azuredevops_serviceendpoint_azurerm` - Support workload identity. [#818](https://github.com/microsoft/terraform-provider-azuredevops/pull/818)
* **Data Resource** `azuredevops_serviceendpoint_azurerm` - Add support for managed identity and federated workload. [#818](https://github.com/microsoft/terraform-provider-azuredevops/pull/818)

BUG FIX:
* `azuredevops_pipeline_authorization` - Fix resource recreate with `pipeline_id` not configured [#809](https://github.com/microsoft/terraform-provider-azuredevops/pull/809)
* `azuredevops_serviceendpoint_azurerm` - Fix imported resource force recreate bug [#827](https://github.com/microsoft/terraform-provider-azuredevops/pull/827)
* `azuredevops_branch_policy_status_check` - Fixed `filename_patterns` order [#828](https://github.com/microsoft/terraform-provider-azuredevops/pull/828)
* `azuredevops_git_repository` - Set `default_branch` for imported repository [#829](https://github.com/microsoft/terraform-provider-azuredevops/pull/829)

## 0.6.0

FEATURES:
* **New Resource** `azuredevops_check_approval` [#728](https://github.com/microsoft/terraform-provider-azuredevops/pull/728)
* **New Resource** `azuredevops_serviceendpoint_gcp_terraform` [#742](https://github.com/microsoft/terraform-provider-azuredevops/pull/742)
* **New Resource** `azuredevops_pipeline_authorization` - Alternative to `azuredevops_resource_authorization` [#787](https://github.com/microsoft/terraform-provider-azuredevops/pull/787)
* **New Data Resource** `azuredevops_serviceendpoint_npm` [#795](https://github.com/microsoft/terraform-provider-azuredevops/pull/795)
* **New Data Resource** `azuredevops_serviceendpoint_sonarcloud` [#796](https://github.com/microsoft/terraform-provider-azuredevops/pull/796)

IMPROVEMENTS:
* `azuredevops_workitem` - Add support for `area_path` and `iteration_path` [#750](https://github.com/microsoft/terraform-provider-azuredevops/pull/750)
* `azuredevops_check_approval` - Set `timeout` default value [#760](https://github.com/microsoft/terraform-provider-azuredevops/pull/760)
* `azuredevops_git_repository` - Uppercase the name of `readme.md` file [#761](https://github.com/microsoft/terraform-provider-azuredevops/pull/761)
* `azuredevops_project_pipeline_settings` - Add support for `enforce_job_scope_for_release`[#777](https://github.com/microsoft/terraform-provider-azuredevops/pull/777)
* Upgrade API from v6 to v7. [#774](https://github.com/microsoft/terraform-provider-azuredevops/pull/774)
* Upgrade legacy API from v5 to v7. [#785](https://github.com/microsoft/terraform-provider-azuredevops/pull/785)

BUG FIX:
* `azuredevops_branch_policy_min_reviewers` - Fix `on_push_reset_approved_votes` cannot set to `true` [#792](https://github.com/microsoft/terraform-provider-azuredevops/pull/792)
* `azuredevops_project` - Fix state inconsistent after apply [#793](https://github.com/microsoft/terraform-provider-azuredevops/pull/793)

## 0.5.0

FEATURES:
* **New Resource** `azuredevops_serviceendpoint_jfrog_distribution_v2` [#705](https://github.com/microsoft/terraform-provider-azuredevops/pull/705)
* **New Resource** `azuredevops_serviceendpoint_jfrog_artifactory_v2` [#705](https://github.com/microsoft/terraform-provider-azuredevops/pull/705)
* **New Resource** `azuredevops_serviceendpoint_jfrog_platform_v2` [#705](https://github.com/microsoft/terraform-provider-azuredevops/pull/705)
* **New Resource** `azuredevops_serviceendpoint_jfrog_xray_v2` [#705](https://github.com/microsoft/terraform-provider-azuredevops/pull/705)

IMPROVEMENTS: 
* `azuredevops_serviceendpoint_azurerm` - Add support for resource state migration created prior to v0.4.0. [#754](https://github.com/microsoft/terraform-provider-azuredevops/pull/754)
* `azuredevops_variable_group` - Enhance create state handler. [#756](https://github.com/microsoft/terraform-provider-azuredevops/pull/756)
* **Data Resource** `azuredevops_team` - Support export `descriptor`. [#753](https://github.com/microsoft/terraform-provider-azuredevops/pull/753)

## 0.4.0

FEATURES:
* **New Resource** `azuredevops_workitem` [#659](https://github.com/microsoft/terraform-provider-azuredevops/pull/659)
* **New Resource** `azuredevops_serviceendpoint_externaltfs` [#676](https://github.com/microsoft/terraform-provider-azuredevops/pull/676)
* **New Resource** `azuredevops_check_branch_control` [#706](https://github.com/microsoft/terraform-provider-azuredevops/pull/706)
* **New Resource** `azuredevops_check_business_hours` [#706](https://github.com/microsoft/terraform-provider-azuredevops/pull/706)
* **New Resource** `azuredevops_git_repository_branch` [#713](https://github.com/microsoft/terraform-provider-azuredevops/pull/713)

BUG FIX:
* `azuredevops_git_repository_file` - Create new file if deleted. [#680](https://github.com/microsoft/terraform-provider-azuredevops/pull/680)
* `azuredevops_serviceendpoint_npm` - Fix `access_token` not updated after change. [#708](https://github.com/microsoft/terraform-provider-azuredevops/pull/708)
* `azuredevops_serviceendpoint_artifactory` - Fix unit test. [#725](https://github.com/microsoft/terraform-provider-azuredevops/pull/725)
* `azuredevops_build_folder` - Fix `path` cannot be updated. [#730](https://github.com/microsoft/terraform-provider-azuredevops/pull/730)

IMPROVEMENTS:
* `azuredevops_build_folder_permissions` - Check if the folder exists. [#714](https://github.com/microsoft/terraform-provider-azuredevops/pull/714)
* `azuredevops_branch_policy_auto_reviewers` - Support config minimum number of reviewers. [#672](https://github.com/microsoft/terraform-provider-azuredevops/pull/672)
* `azuredevops_agent_pool` - Enhance create/update handler. [#716](https://github.com/microsoft/terraform-provider-azuredevops/pull/716)
* `azuredevops_serviceendpoint_azurerm` - Support for `environment` property. [#699](https://github.com/microsoft/terraform-provider-azuredevops/pull/699)
* `azuredevops_agent_pool` - Support for `auto_update` property. [#690](https://github.com/microsoft/terraform-provider-azuredevops/pull/690)
* **Date Resource** `azuredevops_agent_pool` - Support for `auto_update` property. [#690](https://github.com/microsoft/terraform-provider-azuredevops/pull/690)

## 0.3.0

FEATURES:
* **New Resource** `azuredevops_serviceendpoint_sonarcloud` [#658](https://github.com/microsoft/terraform-provider-azuredevops/pull/658)
* **New Data Resource** `azuredevops_serviceendpoint_azurerm` [#623](https://github.com/microsoft/terraform-provider-azuredevops/pull/623)
* **New Data Resource** `azuredevops_serviceendpoint_github` [#627](https://github.com/microsoft/terraform-provider-azuredevops/pull/627)

BUG FIX:
* `azuredevops_project` - Fall back to organization default template if template ID not found. [#626](https://github.com/microsoft/terraform-provider-azuredevops/pull/626)
* `azuredevops_serviceendpoint_kubernetes` - Fix plugin crash when the `cluster_context` attribute was not specified. [#638](https://github.com/microsoft/terraform-provider-azuredevops/pull/638)
* `azuredevops_build_definition_permissions` - Recreate the resource if relate build definition not found. [#644](https://github.com/microsoft/terraform-provider-azuredevops/pull/644)
* `azuredevops_serviceendpoint_artifactory` - Fix token lost when update other properties. [#656](https://github.com/microsoft/terraform-provider-azuredevops/pull/656)

IMPROVEMENTS:
* `azuredevops_variable_group` - Support custom Key Vault secrets search depth. [#654](https://github.com/microsoft/terraform-provider-azuredevops/pull/654)
* `azuredevops_team` - Support export team `descriptor`. [#648](https://github.com/microsoft/terraform-provider-azuredevops/pull/648)
* Upgrade Terraform Plugin SDK to `v2.23.0` - [#587](https://github.com/microsoft/terraform-provider-azuredevops/issues/587)

## 0.2.2

FEATURES:
* **New Resource** `azuredevops_serviceendpoint_octopusdeploy` [#529](https://github.com/microsoft/terraform-provider-azuredevops/issues/529)
* **New Resource** `azuredevops_serviceendpoint_incomingwebhook ` [#531](https://github.com/microsoft/terraform-provider-azuredevops/issues/531)
* **New Data Resource** `azuredevops_build_definitions ` [#562](https://github.com/microsoft/terraform-provider-azuredevops/issues/562)


BUG FIX:
* `azuredevops_serviceendpoint_kubernetes` - Does not update `service_account` values when changed. [#576](https://github.com/microsoft/terraform-provider-azuredevops/issues/576)
* `azuredevops_project_features` - Fix concurrent modification error. [#593](https://github.com/microsoft/terraform-provider-azuredevops/issues/593)
* `azuredevops_project` - Fix concurrent modification error.  [#593](https://github.com/microsoft/terraform-provider-azuredevops/issues/593)
* `azuredevops_project` - Handling 404 error code.  [#614](https://github.com/microsoft/terraform-provider-azuredevops/issues/614)

IMPROVEMENTS:
* `azuredevops_serviceendpoint_azurerm` - Support for management group scope. [#527](https://github.com/microsoft/terraform-provider-azuredevops/issues/527)
* `azuredevops_branch_policy_build_validation"` - Adding `DefaultBranch` into `match_type` setting. [#305](https://github.com/microsoft/terraform-provider-azuredevops/issues/305)
* `azuredevops_branch_policy_min_reviewers"` - Adding `DefaultBranch` into `match_type` setting. [#305](https://github.com/microsoft/terraform-provider-azuredevops/issues/305)
* `azuredevops_branch_policy_auto_reviewers"` - Adding `DefaultBranch` into `match_type` setting. [#305](https://github.com/microsoft/terraform-provider-azuredevops/issues/305)
* `azuredevops_branch_policy_work_item_linking"` - Adding `DefaultBranch` into `match_type` setting. [#305](https://github.com/microsoft/terraform-provider-azuredevops/issues/305)
* `azuredevops_branch_policy_comment_resolution"` - Adding `DefaultBranch` into `match_type` setting. [#305](https://github.com/microsoft/terraform-provider-azuredevops/issues/305)
* `azuredevops_branch_policy_merge_types"` - Adding `DefaultBranch` into `match_type` setting. [#305](https://github.com/microsoft/terraform-provider-azuredevops/issues/305)
* `azuredevops_branch_policy_status_check"` - Adding `DefaultBranch` into `match_type` setting. [#305](https://github.com/microsoft/terraform-provider-azuredevops/issues/305)
* `azuredevops_project_pipeline_settings` - Replace deprecated APIs with latest SDK APIs. [#609](https://github.com/microsoft/terraform-provider-azuredevops/issues/609)
* Upgrade Terraform Plugin SDK to `v2.11.0` - [#587](https://github.com/microsoft/terraform-provider-azuredevops/issues/587)

BREAKING CHANGE:
* `azuredevops_serviceendpoint_servicefabric` - Remove sensitive data hashes. [#613](https://github.com/microsoft/terraform-provider-azuredevops/issues/613)

## 0.2.1

FEATURES:
* **New Resource** `azuredevops_project_pipeline_settings` [#556](https://github.com/microsoft/terraform-provider-azuredevops/issues/556)

BUG FIX:
* `azuredevops_group` - Fix scope not set [#542](https://github.com/microsoft/terraform-provider-azuredevops/issues/542)
* `azuredevops_branch_policy_build_validation` - Fix `filename_patterns` disordered.  [#539](https://github.com/microsoft/terraform-provider-azuredevops/issues/539)
* `azuredevops_variable_group` - Fix create 401 authorization error.  [#541](https://github.com/microsoft/terraform-provider-azuredevops/issues/541)
* `azuredevops_group` - Can not create group at project level.  [#558](https://github.com/microsoft/terraform-provider-azuredevops/issues/558)
* `azuredevops_project` - Unable disable/enable project feature artifacts.  [#568](https://github.com/microsoft/terraform-provider-azuredevops/issues/568)

IMPROVEMENTS:
* Update document - [#543](https://github.com/microsoft/terraform-provider-azuredevops/issues/543)
* Deprecate `azuredevops_serviceendpoint_azuredevops`, use `azuredevops_serviceendpoint_runpipeline` instead - [#565](https://github.com/microsoft/terraform-provider-azuredevops/issues/565)

## 0.2.0 

FEATURES:
* **New Resource** `azuredevops_servicehook_permissions` [#504](https://github.com/microsoft/terraform-provider-azuredevops/issues/504)
* **New Resource** `azuredevops_tagging_permissions ` [#510](https://github.com/microsoft/terraform-provider-azuredevops/issues/510)
* **New Resource** `azuredevops_serviceendpoint_argocd ` [#501](https://github.com/microsoft/terraform-provider-azuredevops/issues/501)
* **New Resource** `azuredevops_environment` [#143](https://github.com/microsoft/terraform-provider-azuredevops/issues/143)
* **New Data Resource** `azuredevops_variable_group` [#311](https://github.com/microsoft/terraform-provider-azuredevops/issues/311)

BUG FIX:
* `azuredevops_serviceconnection_azurerm` - Service principal secret will not be updated when update other settings. [#495](https://github.com/microsoft/terraform-provider-azuredevops/issues/495)
* `azuredevops_build_definition`
  - Enhance repository check. [#493](https://github.com/microsoft/terraform-provider-azuredevops/issues/493)
  - `path` cannot end with backslash. [#513](https://github.com/microsoft/terraform-provider-azuredevops/issues/513)
* `azuredevops_git_repository` - `default_branch` cannot set with initialize type `Uninitialized`. [#498](https://github.com/microsoft/terraform-provider-azuredevops/issues/498)
* `azuredevops_variable_group` - Support search top 500 Key Vault secrets. [#388](https://github.com/microsoft/terraform-provider-azuredevops/issues/388)
* `azuredevops_group` - Import group not set scope. [#345](https://github.com/microsoft/terraform-provider-azuredevops/issues/345)

IMPROVEMENTS:
* `data_project` - Optimize read operation [#524](https://github.com/microsoft/terraform-provider-azuredevops/issues/524)
* Document scaffold - Generate document from source code [#503](https://github.com/microsoft/terraform-provider-azuredevops/issues/503)
* Upgrade Azure DevOps API to V6  [#494](https://github.com/microsoft/terraform-provider-azuredevops/issues/494)
* **All permission resources**
  - Refactor the implementation of `SecurityNamespace` and the according helper functions. [#149](https://github.com/microsoft/terraform-provider-azuredevops/pull/149)
  - All permission resources will now clear the `Id` on a `Read` operation when the connected ACLs not found. [#149](https://github.com/microsoft/terraform-provider-azuredevops/pull/149)
  
BREAKING CHANGE:
* All service endpoint - Service endpoint `project_id` only support project ID, project name is no longer supported since v0.2.0. [#494](https://github.com/microsoft/terraform-provider-azuredevops/issues/494)

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
