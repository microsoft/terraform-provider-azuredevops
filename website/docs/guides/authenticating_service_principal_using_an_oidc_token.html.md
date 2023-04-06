---
layout: "azuredevops"
page_title: "Azure DevOps Provider: Authenticating to a Service Principal with an OIDC Token"
description: |-
  This guide will cover how to use an oidc token to authenticate to a service principal for use with Azure DevOps.
---

# Azure DevOps Provider: Authenticating to a Service Principal with an OIDC Token

The Azure DevOps provider supports service principals through a variety of authentication methods, including workload identity federation from any OIDC compliant token issuer.

## Service Principal Configuration

1. Create a service principal in [Azure portal](https://learn.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal) or
using [Azure PowerShell](https://learn.microsoft.com/en-us/azure/active-directory/develop/howto-authenticate-service-principal-powershell). Ignore steps about application roles and certificates.

2. [Configure your app registration to trust your identity provider.](https://learn.microsoft.com/en-us/azure/active-directory/workload-identities/workload-identity-federation-create-trust?pivots=identity-wif-apps-methods-azp#other-identity-providers)

3. [Add the service principal to your Azure DevOps Organization.](https://learn.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/service-principal-managed-identity?view=azure-devops#2-add-and-manage-service-principal-in-an-azure-devops-organization)

## Provider Configuration

The provider will need the Directory (tenant) ID and the Application (client) ID from the Azure AD app registration. They may be provided via the `AZDO_SP_TENANT_ID` and `AZDO_SP_CLIENT_ID` environment variables, or in the provider configuration block with the `sp_tenant_id` and `sp_client_id` attributes.

The token may be provided as a base64 encoded string, or by a file on the filesystem with the `AZDO_SP_OIDC_TOKEN` or `AZDO_SP_OIDC_TOKEN_PATH` environment variables, or in the provider configuration block with the `sp_oidc_token` or `sp_client_oidc_token_path` attributes.

### Providing the token through the file system

```hcl
terraform {
  required_providers {
    azuredevops = {
      source = "microsoft/azuredevops"
      version = ">=0.1.0"
    }
  }
}

provider "azuredevops" {
  org_service_url       = "https://dev.azure.com/my-org"

  sp_client_id              = "00000000-0000-0000-0000-000000000001"
  sp_tenant_id              = "00000000-0000-0000-0000-000000000001"
  sp_client_oidc_token_path = "C:\\my_oidc_token.txt"
}

resource "azuredevops_project" "project" {
  name        = "Test Project"
  description = "Test Project Description"
}
```

### Providing the token directly as a string

```hcl
terraform {
  required_providers {
    azuredevops = {
      source = "microsoft/azuredevops"
      version = ">=0.1.0"
    }
  }
}

provider "azuredevops" {
  org_service_url                = "https://dev.azure.com/my-org"

  sp_client_id  = "00000000-0000-0000-0000-000000000001"
  sp_tenant_id  = "00000000-0000-0000-0000-000000000001"
  sp_oidc_token = "top-secret-base64-encoded-oidc-token-string"
}

resource "azuredevops_project" "project" {
  name               = "Test Project"
  description        = "Test Project Description"
}
```
