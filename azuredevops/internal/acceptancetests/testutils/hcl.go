package testutils

import (
	"fmt"
	"strings"
)

func getGitRepoResource(gitRepoName string, initType string) string {
	return fmt.Sprintf(`
resource "azuredevops_git_repository" "repository" {
	project_id      = azuredevops_project.project.id
	name            = "%s"
	initialization {
		init_type = "%s"
	}
}`, gitRepoName, initType)
}

// HclGitRepoResource HCL describing an AzDO GIT repository resource
func HclGitRepoResource(projectName string, gitRepoName string, initType string) string {
	azureGitRepoResource := getGitRepoResource(gitRepoName, initType)

	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, azureGitRepoResource)
}

// HclForkedGitRepoResource HCL describing an AzDO GIT repository resource
func HclForkedGitRepoResource(projectName string, gitRepoName string, gitForkedRepoName string, initType string, forkedInitType string) string {
	azureGitRepoResource := fmt.Sprintf(`
	resource "azuredevops_git_repository" "gitforkedrepo" {
		project_id      		= azuredevops_project.project.id
		parent_repository_id    = azuredevops_git_repository.repository.id
		name            		= "%s"
		initialization {
			init_type = "%s"
		}
	}`, gitForkedRepoName, forkedInitType)
	gitRepoResource := HclGitRepoResource(projectName, gitRepoName, initType)
	return fmt.Sprintf("%s\n%s", gitRepoResource, azureGitRepoResource)
}

// HclGitRepoFileResource HCl describing a file in an AzDO GIT repository
func HclGitRepoFileResource(projectName, gitRepoName, initType, branch, file, content string) string {
	gitRepoFileResource := fmt.Sprintf(`
	resource "azuredevops_git_repository_file" "file" {
		repository_id = azuredevops_git_repository.repository.id
		file          = "%s"
		content       = "%s"
		branch        = "%s"
	}`, file, content, branch)
	gitRepoResource := HclGitRepoResource(projectName, gitRepoName, initType)
	return fmt.Sprintf("%s\n%s", gitRepoFileResource, gitRepoResource)
}

// HclProjectResource HCL describing an AzDO project
func HclProjectResource(projectName string) string {
	if strings.EqualFold(projectName, "") {
		return ""
	}
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
	name       = "%[1]s"
	description        = "%[1]s-description"
	visibility         = "private"
	version_control    = "Git"
	work_item_template = "Agile"
}`, projectName)
}

// HclServicePrincipleEntitlementResource HCL describing an AzDO service principal entitlement
func HclServicePrincipleEntitlementResource(servicePrincipalObjectId string) string {
	if servicePrincipalObjectId == "" {
		panic("Parameter: servicePrincipalObjectId cannot be empty")
	}
	return fmt.Sprintf(`
resource "azuredevops_service_principal_entitlement" "test" {
	origin_id = "%[1]s"
	origin	  = "aad"
}`, servicePrincipalObjectId)
}

// HclSecurityroleDefinitionsDataSource HCL describing a data source for securityrole definitions
func HclSecurityroleDefinitionsDataSource() string {
	return `
data "azuredevops_securityrole_definitions" "definitions-list" {
	scope = "distributedtask.environmentreferencerole"
}`
}

// HclServiceEndpointGitHubResource HCL describing an AzDO service endpoint
func HclServiceEndpointGitHubResource(projectName string, serviceEndpointName string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_github" "serviceendpoint" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	auth_personal {
	}
}`, serviceEndpointName)

	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// HclServiceEndpointGitHubDataSourceWithServiceEndpointID HCL describing a data source for an AzDO service endpoint
func HclServiceEndpointGitHubDataSourceWithServiceEndpointID() string {
	return `
data "azuredevops_serviceendpoint_github" "serviceendpoint" {
  project_id = azuredevops_project.project.id
  service_endpoint_id         = azuredevops_serviceendpoint_github.serviceendpoint.id
}`
}

// HclServiceEndpointGitHubDataSourceWithServiceEndpointName HCL describing a data source for an AzDO service endpoint
func HclServiceEndpointGitHubDataSourceWithServiceEndpointName(serviceEndpointName string) string {
	return fmt.Sprintf(`
data "azuredevops_serviceendpoint_github" "serviceendpoint" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  depends_on            = [azuredevops_serviceendpoint_github.serviceendpoint]
}
`, serviceEndpointName)
}

func HclServiceEndpointGitHubEnterpriseResource(projectName string, serviceEndpointName string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_github_enterprise" "serviceendpoint" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	url                    = "https://github.contoso.com"
	auth_personal {
		personal_access_token = "hcl_test_token_basic"
	}
}`, serviceEndpointName)

	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// HclServiceEndpointRunPipelineResource HCL describing an AzDO service endpoint
func HclServiceEndpointRunPipelineResourceSimple(serviceEndpointName string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_runpipeline" "serviceendpoint" {
  project_id             = azuredevops_project.project.id
  organization_name      = "example"
  service_endpoint_name  = "%[1]s"
	auth_personal {
	}
}`, serviceEndpointName)

	return serviceEndpointResource
}

func HclServiceEndpointRunPipelineResource(serviceEndpointName string, accessToken string, description string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_runpipeline" "serviceendpoint" {
  project_id             = azuredevops_project.project.id
  organization_name      = "example"
  service_endpoint_name  = "%[1]s"
  auth_personal {
    personal_access_token= "%[2]s"
  }
	description = "%[3]s"
}`, serviceEndpointName, accessToken, description)

	return serviceEndpointResource
}

// HclServiceEndpointDockerRegistryResource HCL describing an AzDO service endpoint
func HclServiceEndpointDockerRegistryResource(projectName string, serviceEndpointName string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_dockerregistry" "serviceendpoint" {
	docker_email           = "test@email.com"
	docker_username        = "testuser"
	docker_password        = "secret"
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"

}`, serviceEndpointName)

	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// HclServiceEndpointAzureRMDataSourceWithServiceEndpointID HCL describing a data source for an AzDO service endpoint
func HclServiceEndpointAzureRMDataSourceWithServiceEndpointID() string {
	return `
data "azuredevops_serviceendpoint_azurerm" "serviceendpointrm" {
  project_id = azuredevops_project.project.id
  service_endpoint_id         = azuredevops_serviceendpoint_azurerm.serviceendpointrm.id
}`
}

// HclServiceEndpointAzureRMDataSourceWithServiceEndpointName HCL describing a data source for an AzDO service endpoint
func HclServiceEndpointAzureRMDataSourceWithServiceEndpointName(serviceEndpointName string) string {
	return fmt.Sprintf(`
data "azuredevops_serviceendpoint_azurerm" "serviceendpointrm" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  depends_on            = [azuredevops_serviceendpoint_azurerm.serviceendpointrm]
}
`, serviceEndpointName)
}

// HclServiceEndpointAzureRMResource HCL describing an AzDO service endpoint
func HclServiceEndpointAzureRMResource(projectName string, serviceEndpointName string, serviceprincipalid string, serviceprincipalkey string, serviceEndpointAuthenticationScheme string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_azurerm" "serviceendpointrm" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  credentials {
    serviceprincipalid  = "%s"
    serviceprincipalkey = "%s"
  }
  azurerm_spn_tenantid                   = "9c59cbe5-2ca1-4516-b303-8968a070edd2"
  azurerm_subscription_id                = "3b0fee91-c36d-4d70-b1e9-fc4b9d608c3d"
  azurerm_subscription_name              = "Microsoft Azure DEMO"
  service_endpoint_authentication_scheme = "%s"
}
`, serviceEndpointName, serviceprincipalid, serviceprincipalkey, serviceEndpointAuthenticationScheme)

	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func HclServiceEndpointAzureRMResourceWithValidate(projectName string, serviceEndpointName string, serviceprincipalid string, serviceprincipalkey string, serviceEndpointAuthenticationScheme string, validate bool) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_azurerm" "serviceendpointrm" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  credentials {
    serviceprincipalid  = "%s"
    serviceprincipalkey = "%s"
  }
  azurerm_spn_tenantid                   = "9c59cbe5-2ca1-4516-b303-8968a070edd2"
  azurerm_subscription_id                = "3b0fee91-c36d-4d70-b1e9-fc4b9d608c3d"
  azurerm_subscription_name              = "Microsoft Azure DEMO"
  service_endpoint_authentication_scheme = "%s"
  features {
	validate = %v
  }
}
`, serviceEndpointName, serviceprincipalid, serviceprincipalkey, serviceEndpointAuthenticationScheme, validate)

	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// HclServiceEndpointAzureRMResource HCL describing an AzDO service endpoint
func HclServiceEndpointAzureRMNoKeyResource(projectName string, serviceEndpointName string, serviceprincipalid string, serviceEndpointAuthenticationScheme string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_azurerm" "serviceendpointrm" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  credentials {
    serviceprincipalid  = "%s"
  }
  azurerm_spn_tenantid                   = "9c59cbe5-2ca1-4516-b303-8968a070edd2"
  azurerm_subscription_id                = "3b0fee91-c36d-4d70-b1e9-fc4b9d608c3d"
  azurerm_subscription_name              = "Microsoft Azure DEMO"
  service_endpoint_authentication_scheme = "%s"
}
`, serviceEndpointName, serviceprincipalid, serviceEndpointAuthenticationScheme)

	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// HclServiceEndpointAzureRMResourceMG HCL describing an AzDO service endpoint
func HclServiceEndpointAzureRMResourceWithMG(projectName string, serviceEndpointName string, serviceprincipalid string, serviceprincipalkey string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_azurerm" "serviceendpointrm" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  credentials {
    serviceprincipalid  = "%s"
    serviceprincipalkey = "%s"
  }
  azurerm_spn_tenantid                   = "9c59cbe5-2ca1-4516-b303-8968a070edd2"
  azurerm_management_group_id            = "Microsoft_Azure_Demo_MG"
  azurerm_management_group_name          = "Microsoft Azure Demo MG"
  service_endpoint_authentication_scheme = "ServicePrincipal"
}
`, serviceEndpointName, serviceprincipalid, serviceprincipalkey)

	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// HclServiceEndpointAzureRMAutomaticResourceWithProject HCL describing an AzDO service endpoint
func HclServiceEndpointAzureRMAutomaticResourceWithProject(projectName string, serviceEndpointName string, serviceEndpointAuthenticationScheme string, subscriptionId string, subscriptionName string, tenantId string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_azurerm" "serviceendpointrm" {
  project_id                             = azuredevops_project.project.id
  service_endpoint_name                  = "%s"
  azurerm_spn_tenantid                   = "%s"
  azurerm_subscription_id                = "%s"
  azurerm_subscription_name              = "%s"
  service_endpoint_authentication_scheme = "%s"
}
`, serviceEndpointName, tenantId, subscriptionId, subscriptionName, serviceEndpointAuthenticationScheme)

	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// HclServiceEndpointServiceFabricResource HCL describing an AzDO service endpoint
func HclServiceEndpointServiceFabricResource(projectName string, serviceEndpointName string, authorizationType string) string {
	var serviceEndpointResource string
	switch authorizationType {
	case "Certificate":
		serviceEndpointResource = fmt.Sprintf(`
resource "azuredevops_serviceendpoint_servicefabric" "serviceendpoint" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  cluster_endpoint      = "tcp://test"
  certificate {
    server_certificate_lookup     = "Thumbprint"
    server_certificate_thumbprint = "test"
    client_certificate            = "test"
    client_certificate_password   = "test"
  }
}`, serviceEndpointName)
	case "UsernamePassword":
		serviceEndpointResource = fmt.Sprintf(`
resource "azuredevops_serviceendpoint_servicefabric" "serviceendpoint" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  cluster_endpoint      = "tcp://test"
  azure_active_directory {
    server_certificate_lookup     = "Thumbprint"
    server_certificate_thumbprint = "test"
    username                      = "test"
    password                      = "test"
  }
}`, serviceEndpointName)
	case "None":
		serviceEndpointResource = fmt.Sprintf(`
resource "azuredevops_serviceendpoint_servicefabric" "serviceendpoint" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "%s"
  cluster_endpoint      = "tcp://test"
  none {
    unsecured   = false
    cluster_spn = "test"
  }
}`, serviceEndpointName)
	}
	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// HclServiceEndpointGenericResource HCL describing an AzDO service endpoint
func HclServiceEndpointGenericResource(projectName string, serviceEndpointName string, serverUrl string, username string, password string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_generic" "test" {
	project_id            = azuredevops_project.project.id
	service_endpoint_name = "%s"
	description           = "test"
	server_url            = "%s"
	username              = "%s"
	password              = "%s"
}`, serviceEndpointName, serverUrl, username, password)

	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// HclVariableGroupResource HCL describing an AzDO group
func HclVariableGroupResource(variableGroupName string, allowAccess bool) string {
	return fmt.Sprintf(`
resource "azuredevops_variable_group" "vg" {
	project_id  = azuredevops_project.project.id
	name        = "%s"
	description = "A sample variable group."
	allow_access = %t
	variable {
		name   = "key1"
		value  = "value1"
	}
	variable {
		name  = "key2"
		value = "value2"
	}
	variable {
		name = "key3"
	}

	secret_variable {
		name   = "skey1"
		value  = "value1"
	}
	secret_variable {
		name  = "skey2"
		value = "value2"
	}
	secret_variable {
		name = "skey3"
	}
}`, variableGroupName, allowAccess)
}

// HclVariableGroupResourceKeyVault HCL describing an AzDO variable group with key vault
func HclVariableGroupResourceKeyVault(variableGroupName string, allowAccess bool, keyVaultName string) string {
	return fmt.Sprintf(`
resource "azuredevops_variable_group" "vg" {
	project_id  = azuredevops_project.project.id
	name        = "%s"
	description = "A sample variable group."
	allow_access = %t
	key_vault {
        name = "%s"
        service_endpoint_id  = azuredevops_serviceendpoint_azurerm.serviceendpointrm.id
    }
	variable {
		name = "key1"
	}
}`, variableGroupName, allowAccess, keyVaultName)
}

// HclVariableGroupDataSource HCL describing a data source for an AzDO Variable Group
func HclVariableGroupDataSource() string {
	return `
data "azuredevops_variable_group" "vg" {
	project_id  = azuredevops_project.project.id
	name        = azuredevops_variable_group.vg.name
}`
}

// HclAgentPoolResource HCL describing an AzDO Agent Pool
func HclAgentPoolResource(poolName string) string {
	return fmt.Sprintf(`
resource "azuredevops_agent_pool" "pool" {
	name           = "%s"
	auto_provision = false
	auto_update    = false
	pool_type      = "automation"
	}`, poolName)
}

// HclAgentPoolResourceAppendPoolNameToResourceName HCL describing an AzDO Agent Pool with agent pool name appended to resource name
func HclAgentPoolResourceAppendPoolNameToResourceName(poolName string) string {
	return fmt.Sprintf(`
resource "azuredevops_agent_pool" "pool_%[1]s" {
	name           = "%[1]s"
	auto_provision = false
	auto_update    = false
	pool_type      = "automation"
	}`, poolName)
}

// HclAgentPoolDataSource HCL describing a data source for an AzDO Agent Pool
func HclAgentPoolDataSource() string {
	return `
data "azuredevops_agent_pool" "pool" {
	name = azuredevops_agent_pool.pool.name
}`
}

// HclAgentPoolsDataSource HCL describing a data source for an AzDO Agent Pools
func HclAgentPoolsDataSource() string {
	return `
data "azuredevops_agent_pools" "pools" {
}`
}

// HclAgentQueueDataSource HCL describing a data source for an AzDO Agent Queue
func HclAgentQueueDataSource(projectName, queueName string) string {
	return fmt.Sprintf(`
%s

data "azuredevops_agent_queue" "queue" {
	project_id = azuredevops_project.project.id
	name = "%s"
}`, HclProjectResource(projectName), queueName)
}

// HclAgentQueueResource HCL describing an AzDO Agent Pool and Agent Queue
func HclAgentQueueResource(projectName, poolName string) string {
	poolHCL := HclAgentPoolResource(poolName)
	queueHCL := fmt.Sprintf(`
resource "azuredevops_project" "p" {
	name = "%s"
}

resource "azuredevops_agent_queue" "q" {
	project_id    = azuredevops_project.p.id
	agent_pool_id = azuredevops_agent_pool.pool.id
}`, projectName)

	return fmt.Sprintf("%s\n%s", poolHCL, queueHCL)
}

// HclBuildDefinitionResourceGitHub HCL describing an AzDO build definition sourced from GitHub
func HclBuildDefinitionResourceGitHub(projectName string, buildDefinitionName string, buildPath string) string {
	return HclBuildDefinitionResourceWithProject(
		projectName,
		buildDefinitionName,
		buildPath,
		"GitHub",
		"repoOrg/repoName",
		"refs/heads/master",
		"path/to/yaml",
		"")
}

// HclBuildDefinitionResourceBitbucket HCL describing an AzDO build definition sourced from Bitbucket
func HclBuildDefinitionResourceBitbucket(projectName string, buildDefinitionName string, buildPath string, serviceConnectionID string) string {
	return HclBuildDefinitionResourceWithProject(
		projectName,
		buildDefinitionName,
		buildPath,
		"Bitbucket",
		"repoOrg/repoName",
		"master",
		"path/to/yaml",
		serviceConnectionID)
}

// HclBuildDefinitionResourceTfsGit HCL describing an AzDO build definition sourced from AzDo Git Repo
func HclBuildDefinitionResourceTfsGit(projectName string, gitRepoName string, buildDefinitionName string, buildPath string) string {
	buildDefinitionResource := HclBuildDefinitionResourceWithProject(
		projectName,
		buildDefinitionName,
		buildPath,
		"TfsGit",
		"${azuredevops_git_repository.repository.id}",
		"refs/heads/master",
		"path/to/yaml",
		"")

	azureGitRepoResource := getGitRepoResource(gitRepoName, "Clean")

	return fmt.Sprintf("%s\n%s", azureGitRepoResource, buildDefinitionResource)
}

// HclBuildDefinitionResource HCL describing an AzDO build definition
func HclBuildDefinitionResource(
	buildDefinitionName string,
	buildPath string,
	repoType string,
	repoID string,
	branchName string,
	yamlPath string,
	serviceConnectionID string,
) string {
	return fmt.Sprintf(`
	resource "azuredevops_build_definition" "build" {
		project_id      = azuredevops_project.project.id
		name            = "%s"
		agent_pool_name = "Azure Pipelines"
		path			= "%s"

		repository {
			repo_type             = "%s"
			repo_id               = "%s"
			branch_name           = "%s"
			yml_path              = "%s"
			service_connection_id = "%s"
		}
	}`, buildDefinitionName, buildPath, repoType, repoID, branchName, yamlPath, serviceConnectionID)
}

// HclBuildDefinitionDataSource HCL describing a data source for an AzDO Variable Group
func HclBuildDefinitionDataSource(path string) string {
	return fmt.Sprintf(`
data "azuredevops_build_definition" "build" {
	project_id  = azuredevops_project.project.id
	name        = azuredevops_build_definition.build.name
	path        = "%s"
}`, path)
}

// HclBuildDefinitionResourceWithProject HCL describing an AzDO build definition and a project
func HclBuildDefinitionResourceWithProject(
	projectName string,
	buildDefinitionName string,
	buildPath string,
	repoType string,
	repoID string,
	branchName string,
	yamlPath string,
	serviceConnectionID string,
) string {
	escapedBuildPath := strings.ReplaceAll(buildPath, `\`, `\\`)
	buildDefinitionResource := HclBuildDefinitionResource(buildDefinitionName, escapedBuildPath, repoType, repoID, branchName, yamlPath, serviceConnectionID)
	projectResource := HclProjectResource(projectName)

	return fmt.Sprintf("%s\n%s", projectResource, buildDefinitionResource)
}

// HclBuildDefinitionWithVariables A build definition with variables
func HclBuildDefinitionWithVariables(varValue, secretVarValue, name string) string {
	buildDefinitionResource := fmt.Sprintf(`
	resource "azuredevops_build_definition" "build" {
		project_id = azuredevops_project.project.id
		name       = "%s"
		repository {
			repo_type   = "TfsGit"
			repo_id     = azuredevops_git_repository.repository.id
			branch_name = azuredevops_git_repository.repository.default_branch
			yml_path    = "azure-pipelines.yml"
		}

		variable {
			name  = "FOO_VAR"
			value = "%s"
		}

		variable {
			name      = "BAR_VAR"
			secret_value     = "%s"
			is_secret = true
		}
	}`, name, varValue, secretVarValue)
	repoAndProjectResource := HclGitRepoResource(name, name+"-repo", "Clean")

	return fmt.Sprintf("%s\n%s", repoAndProjectResource, buildDefinitionResource)
}

// HclResourceAuthorization HCL describing a resource authorization
func HclResourceAuthorization(resourceID string, authorized bool) string {
	return fmt.Sprintf(`
resource "azuredevops_resource_authorization" "auth" {
	project_id  = azuredevops_project.project.id
	resource_id = %s
	authorized  = %t
	type = "endpoint"
}`, resourceID, authorized)
}

// HclDefinitionResourceAuthorization HCL describing a resource authorization
func HclDefinitionResourceAuthorization(resourceID, definitionID, resourceType string, authorized bool) string {
	return fmt.Sprintf(`
resource "azuredevops_resource_authorization" "auth" {
	project_id  = azuredevops_project.project.id
	resource_id = %s
	definition_id = %s
	type = "%s"
	authorized  = %t
}`, resourceID, definitionID, resourceType, authorized)
}

// HclProjectPermissions creates HCL for testing to set permissions for a AzDO project
func HclProjectPermissions(projectName string) string {
	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

data "azuredevops_group" "tf-project-readers" {
	project_id = azuredevops_project.project.id
	name       = "Readers"
}

resource "azuredevops_project_permissions" "project-permissions" {
	project_id  = azuredevops_project.project.id
	principal   = data.azuredevops_group.tf-project-readers.id
	permissions = {
	  DELETE              = "Deny"
	  EDIT_BUILD_STATUS   = "NotSet"
	  WORK_ITEM_MOVE      = "Allow"
	  DELETE_TEST_RESULTS = "Deny"
	}
}
`, projectResource)
}

// HclBuildFolder creates HCL for testing Build Folders
func HclBuildFolder(projectName string, path string, description string) string {
	projectResource := HclProjectResource(projectName)

	escapedBuildPath := strings.ReplaceAll(path, `\`, `\\`)
	return fmt.Sprintf(`
%s

resource "azuredevops_build_folder" "test_folder" {
	project_id  = azuredevops_project.project.id
	path        = "%s"
	description = "%s"
}
`, projectResource, escapedBuildPath, description)
}

func HclTeamConfiguration(projectName string, teamName string, teamDescription string, teamAdministrators *[]string, teamMembers *[]string) string {
	var teamResource string
	projectResource := HclProjectResource(projectName)
	if teamDescription != "" {
		teamResource = fmt.Sprintf(`
%s

resource "azuredevops_team" "team" {
	project_id = azuredevops_project.project.id
	name = "%s"
	description = "%s"
`, projectResource, teamName, teamDescription)
	} else {
		teamResource = fmt.Sprintf(`
%s

resource "azuredevops_team" "team" {
	project_id = azuredevops_project.project.id
	name = "%s"
`, projectResource, teamName)
	}

	if teamAdministrators != nil {
		teamResource = fmt.Sprintf(`
%s
	administrators = [
		%s
	]
`, teamResource, strings.Join(*teamAdministrators, ","))
	}

	if teamMembers != nil {
		teamResource = fmt.Sprintf(`
%s
	members = [
		%s
	]
`, teamResource, strings.Join(*teamMembers, ","))
	}

	return fmt.Sprintf(`
%s
}
`, teamResource)
}

func getEnvironmentResource(environmentName string) string {
	return fmt.Sprintf(`
resource "azuredevops_environment" "environment" {
	project_id = azuredevops_project.project.id
	name       = "%s"
}`, environmentName)
}

// HclEnvironmentResource HCL describing an AzDO environment resource
func HclEnvironmentResource(projectName string, environmentName string) string {
	azureEnvironmentResource := getEnvironmentResource(environmentName)

	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, azureEnvironmentResource)
}

// HclServicehookStorageQeueuePipelinesResource HCL describing an AzDO subscription resource
func HclServicehookStorageQeueuePipelinesResourceWithStageEvent(projectName, accountKey, queueName, stateFilter, resultFilter string) string {
	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_servicehook_storage_queue_pipelines" "test" {
  project_id   = azuredevops_project.project.id
  account_name = "teststorageacc"
  account_key  = "%s"
  queue_name   = "%s"
  stage_state_changed_event {
	stage_state_filter = "%s"
	stage_result_filter = "%s"
  }
}
`, projectResource, accountKey, queueName, stateFilter, resultFilter)
}

func HclServicehookStorageQeueuePipelinesResourceWithoutEventConfig(projectName, accountKey, queueName, eventType string) string {
	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_servicehook_storage_queue_pipelines" "test" {
  project_id   = azuredevops_project.project.id
  account_name = "teststorageacc"
  account_key  = "%s"
  queue_name   = "%s"
  %s {}
}
`, projectResource, accountKey, queueName, eventType)
}
