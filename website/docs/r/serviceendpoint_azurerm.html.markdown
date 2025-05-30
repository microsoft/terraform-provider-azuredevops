---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_azurerm"
description: |-
  Manages a Azure Resource Manager service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_azurerm

Manages Manual or Automatic Azure Resource Manager service endpoint within Azure DevOps.

~>**NOTE:**
If you receive an error message like:```Failed to obtain the Json Web Token(JWT) using service principal client ID. Exception message: A configuration issue is preventing authentication - check the error message from the server for details.```
You should check the secret of this Application or if you recently rotate the secret, wait a few minutes for Azure to propagate the secret.

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
  project_id                             = azuredevops_project.example.id
  service_endpoint_name                  = "Example AzureRM"
  description                            = "Managed by Terraform"
  service_endpoint_authentication_scheme = "ServicePrincipal"
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
  project_id                             = azuredevops_project.example.id
  service_endpoint_name                  = "Example AzureRM"
  description                            = "Managed by Terraform"
  service_endpoint_authentication_scheme = "ServicePrincipal"
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
  project_id                             = azuredevops_project.example.id
  service_endpoint_name                  = "Example AzureRM"
  service_endpoint_authentication_scheme = "ServicePrincipal"
  azurerm_spn_tenantid                   = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_id                = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_name              = "Example Subscription Name"
}
```

### Workload Identity Federation Manual AzureRM Service Endpoint (Subscription Scoped)

```hcl
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=3.0.0"
    }
  }
}

provider "azurerm" {
  features {}
}

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
  location            = azurerm_resource_group.identity.location
  name                = "example-identity"
  resource_group_name = "azurerm_resource_group.identity.name"
}

resource "azuredevops_serviceendpoint_azurerm" "example" {
  project_id                             = azuredevops_project.example.id
  service_endpoint_name                  = local.service_connection_name
  description                            = "Managed by Terraform"
  service_endpoint_authentication_scheme = "WorkloadIdentityFederation"
  credentials {
    serviceprincipalid = azurerm_user_assigned_identity.example.client_id
  }
  azurerm_spn_tenantid      = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_id   = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_name = "Example Subscription Name"
}

resource "azurerm_federated_identity_credential" "example" {
  name                = "example-federated-credential"
  resource_group_name = azurerm_resource_group.identity.name
  parent_id           = azurerm_user_assigned_identity.example.id
  audience            = ["api://AzureADTokenExchange"]
  issuer              = azuredevops_serviceendpoint_azurerm.example.workload_identity_federation_issuer
  subject             = azuredevops_serviceendpoint_azurerm.example.workload_identity_federation_subject
}
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
  project_id                             = azuredevops_project.example.id
  service_endpoint_name                  = "Example AzureRM"
  service_endpoint_authentication_scheme = "WorkloadIdentityFederation"
  azurerm_spn_tenantid                   = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_id                = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_name              = "Example Subscription Name"
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
  project_id                             = azuredevops_project.example.id
  service_endpoint_name                  = "Example AzureRM"
  service_endpoint_authentication_scheme = "ManagedServiceIdentity"
  azurerm_spn_tenantid                   = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_id                = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_name              = "Example Subscription Name"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_name` - (Required) The Service Endpoint Name.

* `azurerm_spn_tenantid` - (Required) The Tenant ID of the service principal.

---

* `service_endpoint_authentication_scheme` - (Optional) Specifies the type of Azure Resource Manager Service Endpoint. Possible values are `WorkloadIdentityFederation`, `ManagedServiceIdentity` or `ServicePrincipal`. Defaults to `ServicePrincipal` for backwards compatibility.

    ~> **NOTE:** The `WorkloadIdentityFederation` authentication scheme is currently in private preview. Your organisation must be part of the preview and the feature toggle must be turned on to use it. More details can be found [here](https://aka.ms/azdo-rm-workload-identity).

* `azurerm_management_group_id` - (Optional) The Management group ID of the Azure targets.

* `azurerm_management_group_name` - (Optional) The Management group Name of the targets.

* `azurerm_subscription_id` - (Optional) The Subscription ID of the Azure targets.

* `azurerm_subscription_name` - (Optional) The Subscription Name of the targets.

* `environment` - (Optional) The Cloud Environment to use. Defaults to `AzureCloud`. Possible values are `AzureCloud`, `AzureChinaCloud`, `AzureUSGovernment`, `AzureGermanCloud` and `AzureStack`. Changing this forces a new resource to be created.

* `server_url` - (Optional) The server URL of the service endpoint. Changing this forces a new resource to be created.

* `shared_project_ids` - (Optional) list of project IDs you want the service connection shared with. 

~> **NOTE:** One of either `Subscription` scoped i.e. `azurerm_subscription_id`, `azurerm_subscription_name` or `ManagementGroup` scoped i.e. `azurerm_management_group_id`, `azurerm_management_group_name` values must be specified.

* `credentials` - (Optional) A `credentials` block as defined below.

* `description` - (Optional) Service connection description.

* `resource_group` - (Optional) The resource group used for scope of automatic service endpoint.

* `features` - (Optional) A `features` block as defined below.

---

A `credentials` block supports the following:

* `serviceprincipalid` - (Required) The service principal application ID

* `serviceprincipalkey` - (Optional) The service principal secret. This not required if `service_endpoint_authentication_scheme` is set to `WorkloadIdentityFederation`.

* `serviceprincipalcertificate` - (Optional) The service principal certificate. This not required if `service_endpoint_authentication_scheme` is set to `WorkloadIdentityFederation`.

---

A `features` block supports the following:

* `validate` - (Optional) Whether or not to validate connection with Azure after create or update operations. Defaults to `false`

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.
* `service_endpoint_name` - The Service Endpoint name.
* `service_principal_id` - The Application(Client) ID of the Service Principal.
* `workload_identity_federation_issuer` - The issuer if `service_endpoint_authentication_scheme` is set to `WorkloadIdentityFederation`. This looks like `https://vstoken.dev.azure.com/00000000-0000-0000-0000-000000000000`, where the GUID is the Organization ID of your Azure DevOps Organisation.
* `workload_identity_federation_subject` - The subject if `service_endpoint_authentication_scheme` is set to `WorkloadIdentityFederation`. This looks like `sc://<organisation>/<project>/<service-connection-name>`.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Service End points](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Azure Resource Manager Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the Azure Resource Manager Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the Azure Resource Manager Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the Azure Resource Manager Service Endpoint.

## Import

Azure DevOps Azure Resource Manager Service Endpoint can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_azurerm.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
