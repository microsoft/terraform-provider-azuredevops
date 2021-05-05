// Configuração do novo projeto
terraform {
  required_providers {
    azuredevops = {
      source = "microsoft/azuredevops"
      version = ">=0.1.0"
    }
  }
}

/*
provider "azuredevops" {
  org_service_url       = var.org_url
  personal_access_token = var.org_token
}
*/

resource "azuredevops_project" "terraform_add_project" {
  name               = var.project_name
  description        = var.description
  visibility         = var.visibility
  version_control    = var.version_control
  work_item_template = var.work_item_template
  features = {
      "boards"       = var.boards
      "repositories" = var.repositories
      "pipelines"    = var.pipelines
      "testplans"    = var.testplans
      "artifacts"    = var.artifacts
  }
}

// Adicionando grupos padrão da azure
data "azuredevops_group" "tf_padrao_projectadm" {
  project_id = azuredevops_project.terraform_add_project.id
  name       = "Project Administrators"
}

// Adicionando grupos padrão da Stone
resource "azuredevops_group" "developers" {
  scope        = azuredevops_project.terraform_add_project.id
  display_name = "Developers"
  description  = "Members of this group can add modify and delete code builds and workitems and view releases within the team project."

  members = [
    data.azuredevops_group.tf_padrao_projectadm.descriptor
  ]
}

resource "azuredevops_group" "infrastructureengineers" {
  scope        = azuredevops_project.terraform_add_project.id
  display_name = "InfrastructureEngineers"
  description  = "Members of this group can add modify and delete code builds workitems and releases within the team project."

  members = [
    data.azuredevops_group.tf_padrao_projectadm.descriptor
  ]
}

resource "azuredevops_group" "techleads" {
  scope        = azuredevops_project.terraform_add_project.id
  display_name = "TechLeads"
  description  = "Members of this group can add modify and delete code builds and workitems and view trigger releases within the team project."

  members = [
    data.azuredevops_group.tf_padrao_projectadm.descriptor
  ]
}

resource "azuredevops_group" "databaseadmins" {
  scope        = azuredevops_project.terraform_add_project.id
  display_name = "DatabaseAdmins"
  description  = "Members of this group can add modify and delete code builds workitems and releases within the team project."

  members = [
    data.azuredevops_group.tf_padrao_projectadm.descriptor
  ]
}

resource "azuredevops_group" "externaldevelopers" {
  scope        = azuredevops_project.terraform_add_project.id
  display_name = "ExternalDevelopers"
  description  = "Members of this group can add modify and delete code builds and workitems and view releases within the team project."

  members = [
    data.azuredevops_group.tf_padrao_projectadm.descriptor
  ]
}

resource "azuredevops_group" "productmanagers" {
  scope        = azuredevops_project.terraform_add_project.id
  display_name = "ProductManagers"
  description  = "Members of this group can add modify and delete workitems and view code builds and releases within the team project."

  members = [
    data.azuredevops_group.tf_padrao_projectadm.descriptor
  ]
}

// Gerenciando as definições de Build
resource "azuredevops_build_definition" "build" {
  project_id = azuredevops_project.terraform_add_project.id
  name       = "Default Build Definition"
}

// Configurando permissões dos grupos padrões Stone
resource "azuredevops_project_permissions" "perm-developers" {
  project_id = azuredevops_project.terraform_add_project.id
  principal  = azuredevops_group.developers.id
  permissions = {
    WORK_ITEM_PERMANENTLY_DELETE = "Deny"
  }
}

resource "azuredevops_project_permissions" "perm-techleads" {
  project_id = azuredevops_project.terraform_add_project.id
  principal  = data.azuredevops_group.techleads.id
  permissions = {
    WORK_ITEM_PERMANENTLY_DELETE = "Deny"
  }
}
/*
resource "azuredevops_project_permissions" "perm-infrastructureengs" {
  project_id = azuredevops_project.terraform_add_project.id
  principal  = data.azuredevops_group.infrastructureengineers.id
  permissions = {

  }
}
*/
resource "azuredevops_project_permissions" "perm-databaseadmins" {
  project_id = azuredevops_project.terraform_add_project.id
  principal  = data.azuredevops_group.databaseadmins.id
  permissions = {
    WORK_ITEM_PERMANENTLY_DELETE = "Deny"
  }
}

resource "azuredevops_project_permissions" "perm-extdevelopers" {
  project_id = azuredevops_project.terraform_add_project.id
  principal  = data.azuredevops_group.externaldevelopers.id
  permissions = {
    WORK_ITEM_PERMANENTLY_DELETE = "Deny"
  }
}

resource "azuredevops_project_permissions" "perm-productmanagers" {
  project_id = azuredevops_project.terraform_add_project.id
  principal  = data.azuredevops_group.productmanagers.id
  permissions = {
    WORK_ITEM_PERMANENTLY_DELETE = "Deny"
  }
}


// Configuração dos endpoints
resource "azuredevops_serviceendpoint_github" "service_endpoint_github" {
  project_id            = azuredevops_project.terraform_add_project.id
  service_endpoint_name = "GithHub Access"
  description = "Setting connection to service endpoint github"

  auth_personal {
    personal_access_token = var.github_token
  }
}

// Configuração dos grupo de variaveis
resource "azuredevops_variable_group" "var_group" {
  project_id   = azuredevops_project.terraform_add_project.id
  name         = "Default Variable Group"
  allow_access = true

  # var group default
  variable {
    name   = "StoneCo.NssmPath"
    value  = var.NssmPath
  }

  variable {
    name   = "StoneCo.DefaultInstallationDir"
    value  = var.DefaultInstallationDir
  }

  variable {
    name   = "StoneCo.ChocolateySource"
    value  = var.ChocolateySource
  }

  # Atlanta production network
  variable {
    name   = "ZabbixEndpoint_atlanta"
    value  = var.ZabbixEndpoint_atlanta
  }

  variable {
    name  = "F5LTMName_atlanta"
    value = var.F5LTMName_atlanta
  }

  variable {
    name          = "F5UserName_atlanta"
    secret_value  = var.F5UserName_atlanta
    is_secret     = true
  }

  variable {
    name          = "F5Password_atlanta"
    secret_value  = var.F5Password_atlanta
    is_secret     = true
  }

  variable {
    name          = "TaskGroupAutoCreated_atlanta"
    secret_value  = var.TaskGroupAutoCreated_atlanta
    is_secret     = true
  }

  # Chicago production network
  variable {
    name   = "ZabbixEndpoint_chicago"
    value  = var.ZabbixEndpoint_chicago
  }

  variable {
    name   = "F5LTMName_chicago"
    value  = var.F5LTMName_chicago
  }

  variable {
    name          = "F5UserName_chicago"
    secret_value  = var.F5UserName_chicago
    is_secret     = true
  }

  variable {
    name          = "F5Password_chicago"
    secret_value  = var.F5Password_chicago
    is_secret     = true
  }

  variable {
    name          = "TaskGroupAutoCreated_chicago"
    secret_value  = var.TaskGroupAutoCreated_chicago
    is_secret     = true
  }


  # Central production network
  variable {
    name   = "ZabbixEndpoint_central"
    value  = var.ZabbixEndpoint_central
  }

  variable {
    name   = "F5LTMName_central"
    value  = var.F5LTMName_central
  }

  variable {
    name          = "F5UserName_central"
    secret_value  = var.F5UserName_central
    is_secret     = true
  }

  variable {
    name          = "F5Password_central"
    secret_value  = var.F5Password_central
    is_secret     = true
  }


  # Tier point non production network
  variable {
    name   = "ZabbixEndpoint_tierpoint_nonproduction"
    value  = var.ZabbixEndpoint_tierpoint_nonproduction
  }

  variable {
    name   = "F5LTMName_tierpoint_nonproduction"
    value  = var.F5LTMName_tierpoint_nonproduction
  }

  variable {
    name          = "F5UserName_tierpoint_nonproduction"
    secret_value  = var.F5UserName_tierpoint_nonproduction
    is_secret     = true
  }

  variable {
    name          = "F5Password_tierpoint_nonproduction"
    secret_value  = var.F5Password_tierpoint_nonproduction
    is_secret     = true
  }


  # Tier point production network
  variable {
    name   = "ZabbixEndpoint_tierpoint_production"
    value  = var.ZabbixEndpoint_tierpoint_production
  }

  variable {
    name   = "F5LTMName_tierpoint_production"
    value  = var.F5LTMName_tierpoint_production
  }

  variable {
    name          = "F5UserName_tierpoint_production"
    secret_value  = var.F5UserName_tierpoint_production
    is_secret     = true
  }

  variable {
    name          = "F5Password_tierpoint_production"
    secret_value  = var.F5Password_tierpoint_production
    is_secret     = true
  }
}
