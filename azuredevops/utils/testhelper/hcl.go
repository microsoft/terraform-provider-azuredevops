package testhelper

import (
	"fmt"
	"strings"
)

func getAzureGitRepoResource(gitRepoName string, initType string) string {
	return fmt.Sprintf(`
resource "azuredevops_git_repository" "gitrepo" {
	project_id      = azuredevops_project.project.id
	name            = "%s"
	initialization {
		init_type = "%s"
	}
}`, gitRepoName, initType)
}

// TestAccAzureGitRepoResource HCL describing an AzDO GIT repository resource
func TestAccAzureGitRepoResource(projectName string, gitRepoName string, initType string) string {
	azureGitRepoResource := getAzureGitRepoResource(gitRepoName, initType)

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, azureGitRepoResource)
}

// TestAccAzureForkedGitRepoResource HCL describing an AzDO GIT repository resource
func TestAccAzureForkedGitRepoResource(projectName string, gitRepoName string, gitForkedRepoName string, initType string, forkedInitType string) string {
	azureGitRepoResource := fmt.Sprintf(`
	resource "azuredevops_git_repository" "gitforkedrepo" {
		project_id      		= azuredevops_project.project.id
		parent_repository_id    = azuredevops_git_repository.gitrepo.id
		name            		= "%s"
		initialization {
			init_type = "%s"
		}
	}`, gitForkedRepoName, forkedInitType)
	gitRepoResource := TestAccAzureGitRepoResource(projectName, gitRepoName, initType)
	return fmt.Sprintf("%s\n%s", gitRepoResource, azureGitRepoResource)
}

// TestAccGroupDataSource HCL describing an AzDO Group Data Source
func TestAccGroupDataSource(projectName string, groupName string) string {
	dataSource := fmt.Sprintf(`
data "azuredevops_group" "group" {
	project_id = azuredevops_project.project.id
	name       = "%s"
}`, groupName)

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, dataSource)
}

// TestAccProjectResource HCL describing an AzDO project
func TestAccProjectResource(projectName string) string {
	if strings.EqualFold(projectName, "") {
		return ""
	}
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
	project_name       = "%s"
	description        = "%s-description"
	visibility         = "private"
	version_control    = "Git"
	work_item_template = "Agile"
}`, projectName, projectName)
}

// TestAccProjectDataSource HCL describing a data source for an AzDO project
func TestAccProjectDataSource(projectName string) string {
	return fmt.Sprintf(`
data "azuredevops_project" "project" {
	project_name = "%s"
}`, projectName)
}

// TestAccProjectGitRepositories HCL describing a data source for an AzDO git repo
func TestAccProjectGitRepositories(projectName string, gitRepoName string) string {
	return fmt.Sprintf(`
data "azuredevops_project" "project" {
	project_name = "%s"
}

data "azuredevops_git_repositories" "repositories" {
	project_id = data.azuredevops_project.project.id
	name = "%s"
}`, projectName, gitRepoName)
}

// TestAccUserEntitlementResource HCL describing an AzDO UserEntitlement
func TestAccUserEntitlementResource(principalName string) string {
	return fmt.Sprintf(`
resource "azuredevops_user_entitlement" "user" {
	principal_name     = "%s"
	account_license_type = "express"
}`, principalName)
}

// TestAccServiceEndpointGitHubResource HCL describing an AzDO service endpoint
func TestAccServiceEndpointGitHubResource(projectName string, serviceEndpointName string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_github" "serviceendpoint" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	auth_personal {
	}
}`, serviceEndpointName)

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// TestAccServiceEndpointDockerHubResource HCL describing an AzDO service endpoint
func TestAccServiceEndpointDockerHubResource(projectName string, serviceEndpointName string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_dockerhub" "serviceendpoint" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
}`, serviceEndpointName)

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// TestAccServiceEndpointKubernetesResource HCL describing an AzDO kubernetes service endpoint
func TestAccServiceEndpointKubernetesResource(projectName string, serviceEndpointName string, authorizationType string) string {
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
	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// TestAccServiceEndpointAzureRMResource HCL describing an AzDO service endpoint
func TestAccServiceEndpointAzureRMResource(projectName string, serviceEndpointName string) string {
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

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// TestAccServiceEndpointAzureRMAutomaticResource HCL describing an AzDO service endpoint
func TestAccServiceEndpointAzureRMAutomaticResource(projectName string, serviceEndpointName string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_azurerm" "serviceendpointrm" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	azurerm_spn_tenantid      = "9c59cbe5-2ca1-4516-b303-8968a070edd2"
    azurerm_subscription_id   = "3b0fee91-c36d-4d70-b1e9-fc4b9d608c3d"
    azurerm_subscription_name = "Microsoft Azure DEMO"

}`, serviceEndpointName)

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

// TestAccVariableGroupResource HCL describing an AzDO variable group
func TestAccVariableGroupResource(projectName string, variableGroupName string, allowAccess bool) string {
	variableGroupResource := fmt.Sprintf(`
resource "azuredevops_variable_group" "vg" {
	project_id  = azuredevops_project.project.id
	name        = "%s"
	description = "A sample variable group."
	allow_access = %t
	variable {
		name      = "key1"
		value     = "value1"
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

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, variableGroupResource)
}

// TestAccVariableGroupResourceNoSecrets Similar to TestAccVariableGroupResource, but without a secret variable
func TestAccVariableGroupResourceNoSecrets(projectName string, variableGroupName string, allowAccess bool) string {
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

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, variableGroupResource)
}

// TestAccAgentPoolResource HCL describing an AzDO Agent Pool
func TestAccAgentPoolResource(poolName string) string {
	return fmt.Sprintf(`
resource "azuredevops_agent_pool" "pool" {
	name           = "%s"
	auto_provision = false
	pool_type      = "automation"
	}`, poolName)
}

// TestAccBuildDefinitionResourceGitHub HCL describing an AzDO build definition sourced from GitHub
func TestAccBuildDefinitionResourceGitHub(projectName string, buildDefinitionName string, buildPath string) string {
	return TestAccBuildDefinitionResource(
		projectName,
		buildDefinitionName,
		buildPath,
		"GitHub",
		"repoOrg/repoName",
		"master",
		"path/to/yaml",
		"")
}

// TestAccBuildDefinitionResourceBitbucket HCL describing an AzDO build definition sourced from Bitbucket
func TestAccBuildDefinitionResourceBitbucket(projectName string, buildDefinitionName string, buildPath string, serviceConnectionID string) string {
	return TestAccBuildDefinitionResource(
		projectName,
		buildDefinitionName,
		buildPath,
		"Bitbucket",
		"repoOrg/repoName",
		"master",
		"path/to/yaml",
		serviceConnectionID)
}

// TestAccBuildDefinitionResourceTfsGit HCL describing an AzDO build definition sourced from AzDo Git Repo
func TestAccBuildDefinitionResourceTfsGit(projectName string, gitRepoName string, buildDefinitionName string, buildPath string) string {
	buildDefinitionResource := TestAccBuildDefinitionResource(
		projectName,
		buildDefinitionName,
		buildPath,
		"TfsGit",
		"${azuredevops_git_repository.gitrepo.id}",
		"master",
		"path/to/yaml",
		"")

	azureGitRepoResource := getAzureGitRepoResource(gitRepoName, "Clean")

	return fmt.Sprintf("%s\n%s", azureGitRepoResource, buildDefinitionResource)
}

// TestAccBuildDefinitionResource HCL describing an AzDO build definition
func TestAccBuildDefinitionResource(
	projectName string,
	buildDefinitionName string,
	buildPath string,
	repoType string,
	repoID string,
	branchName string,
	yamlPath string,
	serviceConnectionID string,
) string {
	repositoryBlock := fmt.Sprintf(`
repository {
	repo_type             = "%s"
	repo_id               = "%s"
	branch_name           = "%s"
	yml_path              = "%s"
	service_connection_id = "%s"
}`, repoType, repoID, branchName, yamlPath, serviceConnectionID)

	buildDefinitionResource := fmt.Sprintf(`
resource "azuredevops_build_definition" "build" {
	project_id      = azuredevops_project.project.id
	name            = "%s"
	agent_pool_name = "Hosted Ubuntu 1604"
	path			= "%s"

	%s
}`, buildDefinitionName, strings.ReplaceAll(buildPath, `\`, `\\`), repositoryBlock)

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, buildDefinitionResource)
}

// TestAccGroupMembershipResource full terraform stanza to standup a group membership
func TestAccGroupMembershipResource(projectName, groupName, userPrincipalName string) string {
	membershipDependenciesStanza := TestAccGroupMembershipDependencies(projectName, groupName, userPrincipalName)
	membershipStanza := `
resource "azuredevops_group_membership" "membership" {
	group = data.azuredevops_group.group.descriptor
	members = [azuredevops_user_entitlement.user.descriptor]
}`

	return membershipDependenciesStanza + "\n" + membershipStanza
}

// TestAccGroupMembershipDependencies all the dependencies needed to configure a group membership
func TestAccGroupMembershipDependencies(projectName, groupName, userPrincipalName string) string {
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

// TestAccGroupResource HCL describing an AzDO group, if the projectName is empty, only a azuredevops_group instance is returned
func TestAccGroupResource(groupResourceName, projectName, groupName string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_group" "%s" {
	scope        = azuredevops_project.project.id
	display_name = "%s"
}

output "group_id_%s" {
	value = azuredevops_group.%s.id
}
`, TestAccProjectResource(projectName), groupResourceName, groupName, groupResourceName, groupResourceName)
}
