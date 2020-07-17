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

// HclGroupDataSource HCL describing an AzDO Group Data Source
func HclGroupDataSource(projectName string, groupName string) string {
	dataSource := fmt.Sprintf(`
data "azuredevops_group" "group" {
	project_id = azuredevops_project.project.id
	name       = "%s"
}`, groupName)

	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, dataSource)
}

// HclProjectResource HCL describing an AzDO project
func HclProjectResource(projectName string) string {
	if strings.EqualFold(projectName, "") {
		return ""
	}
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
	project_name       = "%[1]s"
	description        = "%[1]s-description"
	visibility         = "private"
	version_control    = "Git"
	work_item_template = "Agile"
}`, projectName)
}

// HclProjectDataSource HCL describing a data source for an AzDO project
func HclProjectDataSource(projectName string) string {
	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

data "azuredevops_project" "project" {
	project_name = azuredevops_project.project.project_name
}`, projectResource)
}

// HclProjectResourceWithFeature HCL describing an AzDO project including internal feature setup
func HclProjectResourceWithFeature(projectName string, featureStateTestplans string, featureStateArtifacts string) string {
	if projectName == "" {
		panic("Parameter: projectName cannot be empty")
	}
	if featureStateTestplans == "" {
		panic("Parameter: featureStateTestplans cannot be empty")
	}
	if featureStateArtifacts == "" {
		panic("Parameter: featureStateArtifacts cannot be empty")
	}
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
	project_name       = "%s"
	description        = "%s-description"
	visibility         = "private"
	version_control    = "Git"
	work_item_template = "Agile"

	features = {
		"testplans" = "%s"
		"artifacts" = "%s"
	}
}`, projectName, projectName, featureStateTestplans, featureStateArtifacts)
}

// HclProjectFeatures HCL describing an AzDO project including feature setup using azuredevops_git_repositories
func HclProjectFeatures(projectName string, featureStateTestplans string, featureStateArtifacts string) string {
	projectFeatures := fmt.Sprintf(`
resource "azuredevops_project_features" "project-features" {
	project_id = azuredevops_project.project.id
	features = {
		"testplans" = "%s"
		"artifacts" = "%s"
	}
}`, featureStateTestplans, featureStateArtifacts)

	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, projectFeatures)
}

// HclProjectsDataSource HCL describing a data source for multiple AzDO projects
func HclProjectsDataSource(projectName string) string {
	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

data "azuredevops_projects" "project-list" {
	project_name = azuredevops_project.project.project_name
}
`, projectResource)
}

// HclProjectsDataSourceWithStateAndInvalidName creates HCL for a multi value data source for AzDo projects
func HclProjectsDataSourceWithStateAndInvalidName() string {
	return `data "azuredevops_projects" "project-list" {
		project_name = "_invalid_project_name"
		state = "wellFormed"
	}`
}

// HclProjectGitRepository HCL describing a single-value data source for an AzDO git repository
func HclProjectGitRepository(projectName string, gitRepoName string) string {
	return fmt.Sprintf(`
data "azuredevops_project" "project" {
	project_name = "%s"
}

data "azuredevops_git_repository" "repository" {
	project_id = data.azuredevops_project.project.id
	name = "%s"
}`, projectName, gitRepoName)
}

// HclProjectGitRepositories HCL describing a multivalue data source for AzDO git repositories
func HclProjectGitRepositories(projectName string, gitRepoName string) string {
	return fmt.Sprintf(`
data "azuredevops_project" "project" {
	project_name = azuredevops_project.project.project_name
}

data "azuredevops_git_repositories" "repositories" {
	project_id = data.azuredevops_project.project.id
	name = "%s"
}`, gitRepoName)
}

// HclUserEntitlementResource HCL describing an AzDO UserEntitlement
func HclUserEntitlementResource(principalName string) string {
	return fmt.Sprintf(`
resource "azuredevops_user_entitlement" "user" {
	principal_name     = "%s"
	account_license_type = "express"
}`, principalName)
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

// HclServiceEndpointDockerRegistryResource HCL describing an AzDO service endpoint
func HclServiceEndpointDockerRegistryResource(projectName string, serviceEndpointName string /*, username string, password string*/) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_dockerregistry" "serviceendpoint" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"

}`, serviceEndpointName /*, username, password*/)

	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// HclServiceEndpointKubernetesResource HCL describing an AzDO kubernetes service endpoint
func HclServiceEndpointKubernetesResource(projectName string, serviceEndpointName string, authorizationType string) string {
	var serviceEndpointResource string
	switch authorizationType {
	case "AzureSubscription":
		serviceEndpointResource = fmt.Sprintf(`
resource "azuredevops_serviceendpoint_kubernetes" "serviceendpoint" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	apiserver_url = "https://sample-kubernetes-cluster.hcp.westeurope.azmk8s.io"
	authorization_type = "AzureSubscription"
	azure_subscription {
		subscription_id = "8a7aace5-66b1-66b1-66b1-8968a070edd2"
		subscription_name = "Microsoft Azure DEMO"
		tenant_id = "2e3a33f9-66b1-66b1-66b1-8968a070edd2"
		resourcegroup_id = "sample-rg"
		namespace = "default"
		cluster_name = "sample-aks"
	}
}`, serviceEndpointName)
	case "ServiceAccount":
		serviceEndpointResource = fmt.Sprintf(`
resource "azuredevops_serviceendpoint_kubernetes" "serviceendpoint" {
	project_id            = azuredevops_project.project.id
	service_endpoint_name = "%s"
	apiserver_url         = "https://sample-kubernetes-cluster.hcp.westeurope.azmk8s.io"
	authorization_type    = "ServiceAccount"
	service_account {
	  token   = "kubernetes_TEST_api_token"
	  ca_cert = "kubernetes_TEST_ca_cert"
	}
}`, serviceEndpointName)
	case "Kubeconfig":
		serviceEndpointResource = fmt.Sprintf(`
resource "azuredevops_serviceendpoint_kubernetes" "serviceendpoint" {
	project_id            = azuredevops_project.project.id
	service_endpoint_name = "%s"
	apiserver_url         = "https://sample-kubernetes-cluster.hcp.westeurope.azmk8s.io"
	authorization_type    = "Kubeconfig"
	kubeconfig {
		kube_config            = <<EOT
								apiVersion: v1
								clusters:
								- cluster:
									certificate-authority: fake-ca-file
									server: https://1.2.3.4
								name: development
								contexts:
								- context:
									cluster: development
									namespace: frontend
									user: developer
								name: dev-frontend
								current-context: dev-frontend
								kind: Config
								preferences: {}
								users:
								- name: developer
								user:
									client-certificate: fake-cert-file
									client-key: fake-key-file
								EOT
		accept_untrusted_certs = true
		cluster_context        = "dev-frontend"
	}
}`, serviceEndpointName)
	}
	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// HclServiceEndpointAzureRMResource HCL describing an AzDO service endpoint
func HclServiceEndpointAzureRMResource(projectName string, serviceEndpointName string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_azurerm" "serviceendpointrm" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	credentials {
		serviceprincipalid 	="e318e66b-ec4b-4dff-9124-41129b9d7150"
		serviceprincipalkey ="d9d210dd-f9f0-4176-afb8-a4df60e1ae72"
	}
	azurerm_spn_tenantid      = "9c59cbe5-2ca1-4516-b303-8968a070edd2"
    azurerm_subscription_id   = "3b0fee91-c36d-4d70-b1e9-fc4b9d608c3d"
    azurerm_subscription_name = "Microsoft Azure DEMO"

}`, serviceEndpointName)

	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// HclServiceEndpointAzureRMAutomaticResourceWithProject HCL describing an AzDO service endpoint
func HclServiceEndpointAzureRMAutomaticResourceWithProject(projectName string, serviceEndpointName string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_azurerm" "serviceendpointrm" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	azurerm_spn_tenantid      = "9c59cbe5-2ca1-4516-b303-8968a070edd2"
    azurerm_subscription_id   = "3b0fee91-c36d-4d70-b1e9-fc4b9d608c3d"
    azurerm_subscription_name = "Microsoft Azure DEMO"

}`, serviceEndpointName)

	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func HclVariableGroupResource(variableGroupName string, allowAccess bool) string {
	return fmt.Sprintf(`
resource "azuredevops_variable_group" "vg" {
	project_id  = azuredevops_project.project.id
	name        = "%s"
	description = "A sample variable group."
	allow_access = %t
	variable {
		name      = "key1"
		secret_value  = "value1"
		is_secret = true
	}

	variable {
		name  = "key2"
		value = "value2"
	}

	variable {
		name = "key3"
	}
}`, variableGroupName, allowAccess)
}

// HclVariableGroupResourceWithProject HCL describing an AzDO variable group
func HclVariableGroupResourceWithProject(projectName string, variableGroupName string, allowAccess bool) string {
	variableGroupResource := HclVariableGroupResource(variableGroupName, allowAccess)
	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, variableGroupResource)
}

// HclVariableGroupResourceNoSecretsWithProject Similar to HclVariableGroupResource, but without a secret variable
func HclVariableGroupResourceNoSecretsWithProject(projectName string, variableGroupName string, allowAccess bool) string {
	variableGroupResource := fmt.Sprintf(`
resource "azuredevops_variable_group" "vg" {
	project_id  = azuredevops_project.project.id
	name        = "%s"
	description = "A sample variable group."
	allow_access = %t
	variable {
		name      = "key1"
		value     = "value1"
	}
}`, variableGroupName, allowAccess)

	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, variableGroupResource)
}

// HclVariableGroupResourceKeyVaultWithProject HCL describing an AzDO project and variable group with key vault
func HclVariableGroupResourceKeyVaultWithProject(projectName string, variableGroupName string, allowAccess bool, keyVaultName string) string {
	projectAndServiceEndpoint := HclServiceEndpointAzureRMResource(projectName, "test-service-connection")

	return fmt.Sprintf("%s\n%s", projectAndServiceEndpoint, HclVariableGroupResourceKeyVault(variableGroupName, allowAccess, keyVaultName))
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

// HclAgentPoolResource HCL describing an AzDO Agent Pool
func HclAgentPoolResource(poolName string) string {
	return fmt.Sprintf(`
resource "azuredevops_agent_pool" "pool" {
	name           = "%s"
	auto_provision = false
	pool_type      = "automation"
	}`, poolName)
}

// HclAgentPoolResourceAppendPoolNameToResourceName HCL describing an AzDO Agent Pool with agent pool name appended to resource name
func HclAgentPoolResourceAppendPoolNameToResourceName(poolName string) string {
	return fmt.Sprintf(`
resource "azuredevops_agent_pool" "pool_%[1]s" {
	name           = "%[1]s"
	auto_provision = false
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

// HclAgentQueueResource HCL describing an AzDO Agent Pool and Agent Queue
func HclAgentQueueResource(projectName, poolName string) string {
	poolHCL := HclAgentPoolResource(poolName)
	queueHCL := fmt.Sprintf(`
resource "azuredevops_project" "p" {
	project_name = "%s"
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
		"master",
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
		"master",
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
	escapedBuildPath := strings.ReplaceAll(buildPath, `\`, `\\`)

	return fmt.Sprintf(`
	resource "azuredevops_build_definition" "build" {
		project_id      = azuredevops_project.project.id
		name            = "%s"
		agent_pool_name = "Hosted Ubuntu 1604"
		path			= "%s"

		repository {
			repo_type             = "%s"
			repo_id               = "%s"
			branch_name           = "%s"
			yml_path              = "%s"
			service_connection_id = "%s"
		}
	}`, buildDefinitionName, escapedBuildPath, repoType, repoID, branchName, yamlPath, serviceConnectionID)
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

// HclGroupMembershipResource full terraform stanza to standup a group membership
func HclGroupMembershipResource(projectName, groupName, userPrincipalName string) string {
	membershipDependenciesStanza := HclGroupMembershipDependencies(projectName, groupName, userPrincipalName)
	membershipStanza := `
resource "azuredevops_group_membership" "membership" {
	group = data.azuredevops_group.group.descriptor
	members = [azuredevops_user_entitlement.user.descriptor]
}`

	return membershipDependenciesStanza + "\n" + membershipStanza
}

// HclGroupMembershipDependencies all the dependencies needed to configure a group membership
func HclGroupMembershipDependencies(projectName, groupName, userPrincipalName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
	project_name = "%s"
}
data "azuredevops_group" "group" {
	project_id = azuredevops_project.project.id
	name       = "%s"
}
resource "azuredevops_user_entitlement" "user" {
	principal_name       = "%s"
	account_license_type = "express"
}

output "group_descriptor" {
	value = data.azuredevops_group.group.descriptor
}
output "user_descriptor" {
	value = azuredevops_user_entitlement.user.descriptor
}
`, projectName, groupName, userPrincipalName)
}

// HclGroupResource HCL describing an AzDO group, if the projectName is empty, only a azuredevops_group instance is returned
func HclGroupResource(groupResourceName, projectName, groupName string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_group" "%s" {
	scope        = azuredevops_project.project.id
	display_name = "%s"
}

output "group_id_%s" {
	value = azuredevops_group.%s.id
}
`, HclProjectResource(projectName), groupResourceName, groupName, groupResourceName, groupResourceName)
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

// HclGitPermissions creates HCl for testing to set permissions for a the all Git repositories of AzDO project
func HclGitPermissions(projectName string) string {
	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

data "azuredevops_group" "project-readers" {
	project_id = azuredevops_project.project.id
	name       = "Readers"
}

resource "azuredevops_git_permissions" "git-permissions" {
	project_id  = azuredevops_project.project.id
	principal   = data.azuredevops_group.project-readers.id
	permissions = {
		CreateRepository = "Deny"
		DeleteRepository = "Deny"
		RenameRepository = "NotSet"
	}
}
`, projectResource)
}

// HclGitPermissionsForRepository creates HCl for testing to set permissions for a the all Git repositories of AzDO project
func HclGitPermissionsForRepository(projectName string, gitRepoName string) string {
	projectResource := HclProjectResource(projectName)
	gitRepository := getGitRepoResource(gitRepoName, "clean")

	return fmt.Sprintf(`
%s

%s

data "azuredevops_group" "project-readers" {
	project_id = azuredevops_project.project.project_id
	name       = "Readers"
}

resource "azuredevops_git_permissions" "git-permissions" {
	project_id    = azuredevops_project.project.project_id
	repository_id = azuredevops_git_repository.gitrepo.id
	principal     = data.azuredevops_group.project-readers.id
	permissions   = {
		CreateRepository = "Deny"
		DeleteRepository = "Deny"
		RenameRepository = "NotSet"
	}
}
`, projectResource, gitRepository)
}
