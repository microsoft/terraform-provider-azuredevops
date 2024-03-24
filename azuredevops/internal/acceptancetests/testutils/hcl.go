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

// HclGroupDataSource HCL describing an AzDO Group Data Source
func HclGroupDataSource(projectName string, groupName string) string {
	if projectName == "" {
		return fmt.Sprintf(`
data "azuredevops_group" "group" {
	name       = "%s"
}`, groupName)
	}
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
	name       = "%[1]s"
	description        = "%[1]s-description"
	visibility         = "private"
	version_control    = "Git"
	work_item_template = "Agile"
}`, projectName)
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
	name       = "%s"
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

// HclProjectFeatures HCL describing an AzDO project including feature setup using azuredevops_git_repositories
func HclProjectPipelineSettings(projectName string, enforceJobAuthScope, enforceReferencedRepoScopedToken, enforceSettableVar, publishPipelineMetadata, statusBadgesArePrivate, enforceJobAuthScopeForReleases bool) string {
	projectPipelineSettings := fmt.Sprintf(`
resource "azuredevops_project_pipeline_settings" "this" {
	project_id = azuredevops_project.project.id

	enforce_job_scope = %t
	enforce_referenced_repo_scoped_token = %t
	enforce_settable_var = %t
	publish_pipeline_metadata = %t
	status_badges_are_private = %t
	enforce_job_scope_for_release = %t
}`, enforceJobAuthScope, enforceReferencedRepoScopedToken, enforceSettableVar, publishPipelineMetadata, statusBadgesArePrivate, enforceJobAuthScopeForReleases)

	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, projectPipelineSettings)
}

// HclProjectsDataSource HCL describing a data source for multiple AzDO projects
func HclProjectsDataSource(projectName string) string {
	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

data "azuredevops_projects" "project-list" {
	name = azuredevops_project.project.name
}
`, projectResource)
}

// HclProjectsDataSourceWithStateAndInvalidName creates HCL for a multi value data source for AzDo projects
func HclProjectsDataSourceWithStateAndInvalidName() string {
	return `data "azuredevops_projects" "project-list" {
		name = "invalid_name"
		state = "wellFormed"
	}`
}

// HclProjectGitRepository HCL describing a single-value data source for an AzDO git repository
func HclProjectGitRepository(projectName string, gitRepoName string) string {
	return fmt.Sprintf(`
data "azuredevops_project" "project" {
	name = "%s"
}

data "azuredevops_git_repository" "repository" {
	project_id = data.azuredevops_project.project.id
	name = "%s"
}`, projectName, gitRepoName)
}

// HclProjectGitRepositories HCL describing a multi value data source for AzDO git repositories
func HclProjectGitRepositories(projectName string, gitRepoName string) string {
	return fmt.Sprintf(`
data "azuredevops_project" "project" {
	name = azuredevops_project.project.name
}

data "azuredevops_git_repositories" "repositories" {
	project_id = data.azuredevops_project.project.id
	name = "%s"
}`, gitRepoName)
}

// HclProjectGitRepositoryImport HCL describing a AzDO git repositories
func HclProjectGitRepositoryImport(gitRepoName string, projectName string) string {
	azureGitRepoResource := fmt.Sprintf(`
	resource "azuredevops_git_repository" "repository" {
		project_id      = azuredevops_project.project.id
		name            = "%s"
		initialization {
		   init_type = "Import"
		   source_type = "Git"
		   source_url = "https://github.com/microsoft/terraform-provider-azuredevops.git"
		 }
	}`, gitRepoName)
	projectResource := HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, azureGitRepoResource)
}

func HclProjectGitRepoImportPrivate(projectName, gitRepoName, gitImportRepoName, serviceEndpointName string) string {
	gitRepoResource := HclGitRepoResource(projectName, gitRepoName, "Clean")
	serviceEndpointResource := fmt.Sprintf(`
	resource "azuredevops_serviceendpoint_generic_git" "serviceendpoint" {
		project_id            = azuredevops_project.project.id
		service_endpoint_name = "%s"
		repository_url        = azuredevops_git_repository.repository.remote_url
	}
	`, serviceEndpointName)
	importGitRepoResource := fmt.Sprintf(`
	resource "azuredevops_git_repository" "import" {
		project_id      = azuredevops_project.project.id
		name            = "%s"
		initialization {
		   init_type             = "Import"
		   source_type           = "Git"
		   source_url            = azuredevops_git_repository.repository.remote_url
		   service_connection_id = azuredevops_serviceendpoint_generic_git.serviceendpoint.id
		 }
	}`, gitImportRepoName)
	return fmt.Sprintf("%s\n%s\n%s", gitRepoResource, serviceEndpointResource, importGitRepoResource)
}

// HclSecurityroleDefinitionsDataSource HCL describing a data source for securityrole definitions
func HclSecurityroleDefinitionsDataSource() string {
	return `
data "azuredevops_securityrole_definitions" "definitions-list" {
	scope = "distributedtask.environmentreferencerole"
}
`
}

// HclUserEntitlementResource HCL describing an AzDO UserEntitlement
func HclUserEntitlementResource(principalName string) string {
	return fmt.Sprintf(`
resource "azuredevops_user_entitlement" "user" {
	principal_name     = "%s"
	account_license_type = "express"
}`, principalName)
}

// HclGroupEntitlementResource HCL describing an AzDO GroupEntitlement
func HclGroupEntitlementResource(displayName string) string {
	return fmt.Sprintf(`
resource "azuredevops_group_entitlement" "group" {
	display_name = "%s"
	account_license_type = "express"
}`, displayName)
}

// HclGroupEntitlementResource HCL describing an AzDO GroupEntitlement linked
// with Azure AD
func HclGroupEntitlementResourceAAD(originId string) string {
	return fmt.Sprintf(`
resource "azuredevops_group_entitlement" "group_aad" {
	origin_id = "%s"
	origin = "aad"
	account_license_type = "express"
}`, originId)
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
}
`
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

// HclServiceEndpointAzureCRResource HCL describing an AzDO service endpoint
func HclServiceEndpointAzureCRResource(projectName string, serviceEndpointName string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_azurecr" "serviceendpoint" {
	project_id                = azuredevops_project.project.id
	service_endpoint_name     = "%s"
	azurecr_spn_tenantid      = "9c59cbe5-2ca1-4516-b303-8968a070edd2"
	azurecr_subscription_id   = "3b0fee91-c36d-4d70-b1e9-fc4b9d608c3d"
	azurecr_subscription_name = "Microsoft Azure DEMO"
	resource_group            = "testrg"
	azurecr_name              = "testacr"
}`, serviceEndpointName)

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

// HclServiceEndpointAzureRMDataSourceWithServiceEndpointID HCL describing a data source for an AzDO service endpoint
func HclServiceEndpointAzureRMDataSourceWithServiceEndpointID() string {
	return `
data "azuredevops_serviceendpoint_azurerm" "serviceendpointrm" {
  project_id = azuredevops_project.project.id
  service_endpoint_id         = azuredevops_serviceendpoint_azurerm.serviceendpointrm.id
}
`
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
	projectAndServiceEndpoint := HclServiceEndpointAzureRMResource(projectName, "test-service-connection", "e318e66b-ec4b-4dff-9124-41129b9d7150", "d9d210dd-f9f0-4176-afb8-a4df60e1ae72", "ServicePrincipal")

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
  name = "%s"
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

func getEnvironmentResourceKubernetes(resourceName string) string {
	return fmt.Sprintf(`
resource "azuredevops_environment_resource_kubernetes" "kubernetes" {
	project_id          = azuredevops_project.project.id
	environment_id      = azuredevops_environment.environment.id
	service_endpoint_id = azuredevops_serviceendpoint_kubernetes.serviceendpoint.id
	
	name         = "%s"
	namespace    = "default"
	cluster_name = "example-aks"
	tags         = ["tag1", "tag2"]
}`, resourceName)
}

// HclEnvironmentResourceKubernetesResource HCL describing an AzDO environment kubernetes resource
func HclEnvironmentResourceKubernetes(projectName string, environmentName string, serviceEndpointName string, resourceName string) string {
	serviceEndpointResource := HclServiceEndpointKubernetesResource(projectName, serviceEndpointName, "ServiceAccount")
	azureEnvironmentResource := getEnvironmentResource(environmentName)
	environmentKubernetesResource := getEnvironmentResourceKubernetes(resourceName)
	return fmt.Sprintf("%s\n%s\n%s", serviceEndpointResource, azureEnvironmentResource, environmentKubernetesResource)
}
