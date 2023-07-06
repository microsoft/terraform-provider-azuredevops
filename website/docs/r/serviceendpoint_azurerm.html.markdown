---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_azurerm"
description: |-
  Manages a AzureRM service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_azurerm

Manages Manual or Automatic AzureRM service endpoint within Azure DevOps.

## Requirements (Manual AzureRM Service Endpoint)

Before to create a service end point in Azure DevOps, you need to create a Service Principal in your Azure subscription.

For detailed steps to create a service principal with Azure cli see the [documentation](https://docs.microsoft.com/en-us/cli/azure/create-an-azure-service-principal-azure-cli?view=azure-cli-latest)

## Example Usage

### Service Principal Manual AzureRM Service Endpoint (Subscription Scoped)

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_azurerm" "example" {
  project_id                    = azuredevops_project.example.id
  service_endpoint_name         = "Example AzureRM"
  description                   = "Managed by Terraform"
  azurerm_service_endpoint_type = "ServicePrincipal"
  credentials {
    serviceprincipalid  = "00000000-0000-0000-0000-000000000000"
    serviceprincipalkey = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
  azurerm_spn_tenantid      = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_id   = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_name = "Example Subscription Name"
}
```

### Service Principal Manual AzureRM Service Endpoint (ManagementGroup Scoped)

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_azurerm" "example" {
  project_id                    = azuredevops_project.example.id
  service_endpoint_name         = "Example AzureRM"
  description                   = "Managed by Terraform"
  azurerm_service_endpoint_type = "ServicePrincipal"
  credentials {
    serviceprincipalid  = "00000000-0000-0000-0000-000000000000"
    serviceprincipalkey = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
  azurerm_spn_tenantid          = "00000000-0000-0000-0000-000000000000"
  azurerm_management_group_id   = "managementGroup"
  azurerm_management_group_name = "managementGroup"
}
```

### Service Principal Automatic AzureRM Service Endpoint

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_azurerm" "example" {
  project_id                    = azuredevops_project.example.id
  service_endpoint_name         = "Example AzureRM"
  azurerm_service_endpoint_type = "ServicePrincipal"
  azurerm_spn_tenantid          = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_id       = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_name     = "Example Subscription Name"
}
```

### Workload Identity Federation Manual AzureRM Service Endpoint (Subscription Scoped)

```hcl
locals {
  service_connection_name = "example-federated-sc"
}

resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azurerm_resource_group" "identity" {
  name     = "identity"
  location = "UK South"
}

resource "azurerm_user_assigned_identity" "example" {
  location            = var.location
  name                = "example-identity"
  resource_group_name = "azurerm_resource_group.identity.name"
}

resource "azurerm_federated_identity_credential" "example" {
  name                = "example-federated-credential"
  resource_group_name = azurerm_resource_group.identity.name
  audience            = ["api://AzureADTokenExchange"]
  issuer              = "https://app.vstoken.visualstudio.com"
  parent_id           = azurerm_user_assigned_identity.example.id
  subject             = "sc://${var.azure_devops_organisation}/${azuredevops_project.example.name}/${local.service_connection_name}"
}

resource "azuredevops_serviceendpoint_azurerm" "example" {
  project_id                    = azuredevops_project.example.id
  service_endpoint_name         = local.service_connection_name
  description                   = "Managed by Terraform"
  azurerm_service_endpoint_type = "WorkloadIdentityFederation"
  credentials {
    serviceprincipalid  = azurerm_user_assigned_identity.example.client_id
  }
  azurerm_spn_tenantid      = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_id   = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_name = "Example Subscription Name"
}

#NOTE: The federated credential subject is formed from the Azure DevOps Organisation, Project and the Service Connection name.
```

### Workload Identity Federation Automatic AzureRM Service Endpoint

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_azurerm" "example" {
  project_id                    = azuredevops_project.example.id
  service_endpoint_name         = "Example AzureRM"
  azurerm_service_endpoint_type = "WorkloadIdentityFederation"
  azurerm_spn_tenantid          = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_id       = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_name     = "Example Subscription Name"
}
```

### Managed Identity AzureRM Service Endpoint

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_azurerm" "example" {
  project_id                    = azuredevops_project.example.id
  service_endpoint_name         = "Example AzureRM"
  azurerm_service_endpoint_type = "ManagedServiceIdentity"
  azurerm_spn_tenantid          = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_id       = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_name     = "Example Subscription Name"
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The ID of the project.
- `service_endpoint_name` - (Required) The Service Endpoint Name.
- `azurerm_spn_tenantid` - (Required) The Tenant ID if the service principal.
- `azurerm_service_endpoint_type` - (Optional) Specifies the type of azurerm endpoint, either `WorkloadIdentityFederation`, `ManagedServiceIdentity` or `ServicePrincipal`. Defaults to `ServicePrincipal` for backwards compatibility.
- `azurerm_management_group_id` - (Optional) The Management group ID of the Azure targets.
- `azurerm_management_group_name` - (Optional) The Management group Name of the targets.
- `azurerm_subscription_id` - (Optional) The Subscription ID of the Azure targets.
- `azurerm_subscription_name` - (Optional) The Subscription Name of the targets.
- `environment` - (Optional) The Cloud Environment to use. Defaults to `AzureCloud`. Possible values are `AzureCloud`, `AzureChinaCloud`. Changing this forces a new resource to be created.

~> **NOTE:** One of either `Subscription` scoped i.e. `azurerm_subscription_id`, `azurerm_subscription_name` or `ManagementGroup` scoped i.e. `azurerm_management_group_id`, `azurerm_management_group_name` values must be specified.

- `description` - (Optional) Service connection description.
- `credentials` - (Optional) A `credentials` block.
- `resource_group` - (Optional) The resource group used for scope of automatic service endpoint.

---

A `credentials` block supports the following:

- `serviceprincipalid` - (Required) The service principal application Id
- `serviceprincipalkey` - (Optional) The service principal secret. This not required if `azurerm_service_endpoint_type` is set to `WorkloadIdentityFederation`.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The ID of the project.
- `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Service End points](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-6.0)

## Import

Azure DevOps Service Endpoint Azure Resource Manage can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_azurerm.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
