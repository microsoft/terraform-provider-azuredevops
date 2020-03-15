package testhelper

import (
	"fmt"
	"strings"
)

// TestAccAzureGitRepoResource HCL describing an AzDO GIT repository resource
func TestAccAzureGitRepoResource(projectName string, gitRepoName string, initType string) string {
	azureGitRepoResource := fmt.Sprintf(`
resource "azuredevops_git_repository" "gitrepo" {
	project_id      = azuredevops_project.project.id
	name            = "%s"
	initialization {
		init_type = "%s"
	}
}`, gitRepoName, initType)

	projectResource := TestAccProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, azureGitRepoResource)
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
	if projectName == "" {
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

// TestAccServiceEndpointAzureRMResource HCL describing an AzDO service endpoint
func TestAccServiceEndpointAzureRMResource(projectName string, serviceEndpointName string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_azurerm" "serviceendpointrm" {
	project_id             = azuredevops_project.project.id
	service_endpoint_name  = "%s"
	azurerm_spn_clientid 	="e318e66b-ec4b-4dff-9124-41129b9d7150"
	azurerm_spn_tenantid      = "9c59cbe5-2ca1-4516-b303-8968a070edd2"
    azurerm_subscription_id   = "3b0fee91-c36d-4d70-b1e9-fc4b9d608c3d"
    azurerm_subscription_name = "Microsoft Azure DEMO"
    azurerm_scope             = "/subscriptions/3b0fee91-c36d-4d70-b1e9-fc4b9d608c3d"
	azurerm_spn_clientsecret ="d9d210dd-f9f0-4176-afb8-a4df60e1ae72"

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

// TestAccAgentPoolResource HCL describing an AzDO Agent Pool
func TestAccAgentPoolResource(poolName string) string {
	return fmt.Sprintf(`
resource "azuredevops_agent_pool" "pool" {
	name           = "%s"
	auto_provision = false
	pool_type      = "automation"
	}`, poolName)
}

// TestAccBuildDefinitionResource HCL describing an AzDO build definition
func TestAccBuildDefinitionResource(projectName string, buildDefinitionName string, buildPath string) string {
	buildDefinitionResource := fmt.Sprintf(`
resource "azuredevops_build_definition" "build" {
	project_id      = azuredevops_project.project.id
	name            = "%s"
	agent_pool_name = "Hosted Ubuntu 1604"
	path			= "%s"

	repository {
	  repo_type             = "GitHub"
	  repo_name             = "repoOrg/repoName"
	  branch_name           = "branch"
	  yml_path              = "path/to/yaml"
	}
}`, buildDefinitionName, strings.ReplaceAll(buildPath, `\`, `\\`))

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
