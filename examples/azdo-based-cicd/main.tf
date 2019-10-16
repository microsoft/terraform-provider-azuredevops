# Note: This template is *aspirational* and is included in order to exemplify the value that
# this provider will eventually provide. Based on this, keep the following in mind:
#  - Not every stanza is backed by an implementation in the provider
#  - The schema for each resource/data source will almost certainly change once implemented

# As the provider is built, the idea is that this file will slowly become uncommented and updated
# once the features are implemented.

provider "azuredevops" {
  version = ">= 0.0.1"
}




// The following section defines an AzDO project that hosts a repository and
// a service connection

// Defines the project in AzDO. This project will host Git repositories
resource "azuredevops_project" "project" {
  project_name = "Super Awesome Project"
  description  = "A project to track super awesome things."
  visibility   = "private"
  # enable_tfvc        = false
  work_item_template = "Agile"
}

// Defines a Git repository hosted in the project
resource "azuredevops_azure_git_repository" "repository" {
  project        = azuredevops_project.project.id
  name           = "main-repo"
  default_branch = "master"
}

// Defines an ARM service connection
# resource "azuredevops_serviceendpoint" "arm" {
#   project_id            = azuredevops_project.project.id
#   service_endpoint_name = "ARM Service Connection"
#   service_endpoint_type = "arm"

#   configuration = {
#     service_principal_username = var.service_principal_username
#     service_principal_password = var.service_principal_password
#     subscription_id            = var.subscription_id
#     tenant_id                  = var.tenant_id
#   }
# }




// The following section defines a build pipeline that will use
// some variable groups

// Defines variable groups that can be used by builds/releases
# resource "azuredevops_variable_group" "group1" {
#   project_id = azuredevops_project.project.id
#   name       = "My First Variable Group"
#   values = [{
#     name      = "SERVICE_CONNECTION_ID"
#     value     = azuredevops_serviceendpoint.arm.id
#     is_secret = false
#     }, {
#     name      = "MY_SECRET_1"
#     value     = "foo"
#     is_secret = true
#   }]
# }

// Defines variable groups that can be used by builds/releases
# resource "azuredevops_variable_group" "group2" {
#   project_id = azuredevops_project.project.id
#   name       = "My Second Variable Group"
#   values = [{
#     name      = "MY_SECRET_2"
#     value     = "bar"
#     is_secret = true
#   }]
# }

// A build that kicks off an infrastructure provisioning pipeline, using variable groups
resource "azuredevops_build_definition" "cicd" {
  project_id            = azuredevops_project.project.id
  agent_pool_name       = "Hosted Ubuntu 1604"
  name                  = "CICD Pipeilne"

  # variables_groups = [
  #   azuredevops_variable_group.group1.id,
  #   azuredevops_variable_group.group2.id
  # ]

  repository {
    # repo_type   = azuredevops_git_repo.repository.type
    # repo_name   = azuredevops_git_repo.repository.name
    # branch_name = azuredevops_git_repo.repository.default_branch
    yml_path    = "cicd/azure-pipelines-infra.yml"
  }
}




// The following section will configure branching policies for the created repositories
// note: see the types of branch policies here:
//   https://dev.azure.com/$ORG/$PROJECT/_apis/policy/types?api-version=5.0

// Look up the ID of the branch policy that ensures each PR is linked to a work item, then
// configure the policy
# data "azuredevops_branch_policy_type" "work_item_policy" {
#   name = "Work item linking"
# }
# resource "azuredevops_branch_policy" "work_item_linked_policy" {
#   project_id    = azuredevops_project.project.id
#   type          = azuredevops_branch_policy_type.work_item_policy.id
#   repository_id = azuredevops_git_repo.repository.id
#   branch        = azuredevops_git_repo.repository.default_branch
# }

// Look up the ID of the branch policy that ensures each PR reviewed by 2 people, then
// configure the policy
# data "azuredevops_branch_policy_type" "min_reviewer_count" {
#   name = "Minimum number of reviewers"
# }
# resource "azuredevops_branch_policy" "min_reviewer_count_policy" {
#   project_id    = azuredevops_project.project.id
#   type          = azuredevops_branch_policy_type.min_reviewer_count.id
#   repository_id = azuredevops_git_repo.repository.id
#   branch        = azuredevops_git_repo.repository.default_branch
#   settings = {
#     "minimumApproverCount" = 2
#   }
# }

// Look up the ID of the branch policy that ensures each PR does not break a build, then
// configure the policy
# data "azuredevops_branch_policy_type" "build" {
#   name = "Minimum number of reviewers"
# }
# resource "azuredevops_branch_policy" "build_policy" {
#   project_id    = azuredevops_project.project.id
#   type          = azuredevops_branch_policy_type.build.id
#   repository_id = azuredevops_git_repo.repository.id
#   branch        = azuredevops_git_repo.repository.default_branch
#   settings = {
#     "buildDefinitionId" = azuredevops_build_definition.cicd.id
#   }
# }




// The following section will import an AAD group into AzDO, and assign it to
// a project specific group

# data "azuredevops_group" "azdo_group" {
#   principal_name = format("[%s]\\Build Administrators", azuredevops_project.project.name)
# }

// add existing AAD group to AzDO
//  https://docs.microsoft.com/en-us/rest/api/azure/devops/graph/groups/create?view=azure-devops-rest-5.0#add-an-aad-group-by-oid
# resource "azuredevops_group" "aad_group" {
#   originId = variables.aad_group_id
# }

# resource "azuredevops_group_membership" "membership" {
#   group_descriptor = azuredevops_group.azdo_group.descriptor
#   members = [
#     azuredevops_group.aad_group.descriptor
#   ]
# }
