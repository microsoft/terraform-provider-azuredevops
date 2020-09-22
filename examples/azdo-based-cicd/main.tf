# Make sure to set the following environment variables:
#   AZDO_PERSONAL_ACCESS_TOKEN
#   AZDO_ORG_SERVICE_URL
provider "azuredevops" {
  version = ">= 0.0.1"
}


// This section creates a project
resource "azuredevops_project" "project" {
  name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}


// This section assigns users from AAD into a pre-existing group in AzDO
data "azuredevops_group" "group" {
  project_id = azuredevops_project.project.id
  name       = "Build Administrators"
}

resource "azuredevops_user_entitlement" "users" {
  for_each             = toset(var.aad_users)
  principal_name       = each.value
  account_license_type = "stakeholder"
}

resource "azuredevops_group_membership" "membership" {
  group   = data.azuredevops_group.group.descriptor
  members = values(azuredevops_user_entitlement.users)[*].descriptor
}



// This section configures variable groups and a build definition
resource "azuredevops_build_definition" "build" {
  project_id = azuredevops_project.project.id
  name       = "Sample Build Definition"
  path       = "\\ExampleFolder"

  repository {
    repo_type   = "TfsGit"
    repo_id     = azuredevops_git_repository.repository.id
    branch_name = azuredevops_git_repository.repository.default_branch
    yml_path    = "azure-pipelines.yml"
  }

  variable_groups = [azuredevops_variable_group.vg.id]
}

// This section configures an Azure DevOps Variable Group
resource "azuredevops_variable_group" "vg" {
  project_id   = azuredevops_project.project.id
  name         = "Sample VG 1"
  description  = "A sample variable group."
  allow_access = true

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
}

// This section configures an Azure DevOps Git Repository with branch policies
resource "azuredevops_git_repository" "repository" {
  project_id = azuredevops_project.project.id
  name       = "Sample Repo"
  initialization {
    init_type = "Clean"
  }
}

// Configuration of AzureRm service end point
resource "azuredevops_serviceendpoint_azurerm" "endpoint1" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "TestServiceAzureRM"
  credentials {
    serviceprincipalid  = "00000000-0000-0000-0000-000000000000"
    serviceprincipalkey = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
  azurerm_spn_tenantid      = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_id   = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_name = "Microsoft Azure DEMO"
}

resource "azuredevops_serviceendpoint_bitbucket" "bitbucket_account" {
  project_id            = "vanilla-sky"
  username              = "xxxx"
  password              = "xxxx"
  service_endpoint_name = "test-bitbucket"
  description           = "test"
}

resource "azuredevops_resource_authorization" "bitbucket_account_authorization" {
  project_id  = azuredevops_project.project.id
  resource_id = azuredevops_serviceendpoint_bitbucket.bitbucket_account.id
  authorized  = true
}

resource "azuredevops_serviceendpoint_kubernetes" "kubeendpoint1" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "Sample Kubernetes"
  apiserver_url         = "https://sample-kubernetes-cluster.hcp.westeurope.azmk8s.io"
  authorization_type    = "AzureSubscription"

  azure_subscription {
    subscription_id   = "1c020621-d7a3-457d-b0cc-5d8e6e12d4e6" # a fake GUID
    subscription_name = "Microsoft Azure DEMO"
    tenant_id         = "e46643be-eb78-472f-9780-e01d8190ba10" # a fake GUID
    resourcegroup_id  = "sample-rg"
    namespace         = "default"
    cluster_name      = "sample-aks"
  }
}

resource "azuredevops_serviceendpoint_kubernetes" "kubeendpoint2" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "Sample Kubernetes"
  apiserver_url         = "https://sample-aks.hcp.westeurope.azmk8s.io"
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
}

resource "azuredevops_serviceendpoint_kubernetes" "serviceendpoint" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "Sample Kubernetes"
  apiserver_url         = "https://sample-kubernetes-cluster.hcp.westeurope.azmk8s.io"
  authorization_type    = "ServiceAccount"

  service_account {
    token   = "bXktYXBw[...]K8bPxc2uQ=="
    ca_cert = "Mzk1MjgkdmRnN0pi[...]mHHRUH14gw4Q=="
  }
}
