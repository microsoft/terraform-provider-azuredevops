---
layout: "azuredevops"
page_title: "Azure DevOps Provider: Authenticating to a Service Principal with an OIDC Token"
description: |-
  This guide will cover how to use an oidc token to authenticate to a service principal for use with Azure DevOps.
---

# Azure DevOps Provider: Authenticating to a Service Principal with an OIDC Token

The Azure DevOps provider supports service principals through a variety of authentication methods, including workload identity federation from any OIDC compliant token issuer. The OIDC token will be used to exchange for an AAD access token, which will be used for the authentication.

The OIDC token can be provided via three different ways:

- By plain OIDC token: This is the most straight forward, whilst restricted method, in that the OIDC token is designed to be short lived. If the AAD access token expires, the provider will use this OIDC token to exchange for another AAD access token, which will most likely fail due to the OIDC token has expired.

- By OIDC token file: This is mainly for supporting AAD authentication when [running in Azure Kubernetes Service clusters](https://learn.microsoft.com/en-us/azure/aks/workload-identity-overview?tabs=dotnet).

- By OIDC request token: This is the recommended way to authenticate when running in Github Action, or Azure DevOps Pipeline.

## Service Principal Configuration

1. Create a service principal in [Azure portal](https://learn.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal) or
using [Azure PowerShell](https://learn.microsoft.com/en-us/azure/active-directory/develop/howto-authenticate-service-principal-powershell). Ignore steps about application roles and certificates.

2. [Configure your app registration to trust your identity provider.](https://learn.microsoft.com/en-us/azure/active-directory/workload-identities/workload-identity-federation-create-trust?pivots=identity-wif-apps-methods-azp#other-identity-providers)

3. [Add the service principal to your Azure DevOps Organization.](https://learn.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/service-principal-managed-identity?view=azure-devops#2-add-and-manage-service-principal-in-an-azure-devops-organization)

## Provider Configuration

The provider will need the Directory (tenant) ID and the Application (client) ID from the Azure AD app registration, which are provided via `tenant_id` and `client_id`. Meanwhile, the `use_oidc` must be set to `true` to use OIDC token flows.

For the different OIDC token flows, different configurations are needed:

- Plain OIDC token: `oidc_token`
- OIDC token file: `oidc_token_file_path`
- OIDC request token:
    - Azure DevOps Pipeline:
        - `oidc_request_token`
        - `oidc_request_url`: Not necessary as it can be sourced from `SYSTEM_OIDCREQUESTURI`, which is populated by ADO pipelines.
        - `oidc_azure_service_connection_id`: Not necessary as it can be sourced from `AZURESUBSCRIPTION_SERVICE_CONNECTION_ID`, which is populated by tasks like `AzureCLI@2`.
    - Github Action: 
        - `oidc_request_token`: Not necessary as it can be sourced from `ACTIONS_ID_TOKEN_REQUEST_TOKEN`, which is populated by Github Action. 
        - `oidc_request_url`: Not necessary as it can be sourced from `ACTIONS_ID_TOKEN_REQUEST_URL`, which is populated by Github Action.

## Examples

### Plain OIDC token

```hcl
terraform {
  required_providers {
    azuredevops = {
      source  = "microsoft/azuredevops"
      version = ">=1"
    }
  }
}

provider "azuredevops" {
  org_service_url = "https://dev.azure.com/my-org"
  client_id  = "00000000-0000-0000-0000-000000000001"
  tenant_id  = "00000000-0000-0000-0000-000000000001"
  use_oidc   = true
  oidc_token = "top-secret-base64-encoded-oidc-token-string"
}

resource "azuredevops_project" "project" {
  name        = "Test Project"
  description = "Test Project Description"
}
```

### OIDC token file

```hcl
terraform {
  required_providers {
    azuredevops = {
      source  = "microsoft/azuredevops"
      version = ">=1"
    }
  }
}

provider "azuredevops" {
  org_service_url      = "https://dev.azure.com/my-org"
  client_id            = "00000000-0000-0000-0000-000000000001"
  tenant_id            = "00000000-0000-0000-0000-000000000001"
  use_oidc             = true
  oidc_token_file_path = "C:\\my_oidc_token.txt"
}

resource "azuredevops_project" "project" {
  name        = "Test Project"
  description = "Test Project Description"
}
```

### GitHub Actions

For GitHub Actions workflows, you'll need to ensure the workflow has `write` permissions for the `id-token`.

```yaml
permissions:
  id-token: write
```

Meanwhile, follow the [Azure document](https://learn.microsoft.com/en-us/entra/workload-id/workload-identity-federation-create-trust?pivots=identity-wif-apps-methods-azp) to configure your app to trust an external identity provider.

For more information about OIDC in GitHub Actions, see [official documentation](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-cloud-providers).


```hcl
terraform {
  required_providers {
    azuredevops = {
      source  = "microsoft/azuredevops"
      version = ">=1"
    }
  }
}

provider "azuredevops" {
  org_service_url      = "https://dev.azure.com/my-org"
  client_id            = "00000000-0000-0000-0000-000000000001"
  tenant_id            = "00000000-0000-0000-0000-000000000001"
  use_oidc             = true
}

resource "azuredevops_project" "project" {
  name        = "Test Project"
  description = "Test Project Description"
}
```

#### Azure Pipelines

Follow the [ADO document](https://learn.microsoft.com/en-gb/azure/devops/pipelines/release/configure-workload-identity?view=azure-devops&tabs=app-registration) to set a workload identity service connection.

It is recommend to use the `AzureCLI@2` task as below (note the azureSubscription input parameter):

```yaml
- task: AzureCLI@2
  inputs:
    azureSubscription: $(SERVICE_CONNECTION_ID)
    scriptType: bash
    scriptLocation: "inlineScript"
    inlineScript: |
      # Terraform commands
  env:
    ARM_USE_OIDC: true
    SYSTEM_ACCESSTOKEN: $(System.AccessToken)
    SYSTEM_OIDCREQUESTURI: $(System.OidcRequestUri)
    ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID: $(SERVICE_CONNECTION_ID)
```

As a result, the only configuration needed is as follows:

```hcl
terraform {
  required_providers {
    azuredevops = {
      source  = "microsoft/azuredevops"
      version = ">=1"
    }
  }
}

provider "azuredevops" {
  org_service_url      = "https://dev.azure.com/my-org"
  client_id            = "00000000-0000-0000-0000-000000000001"
  tenant_id            = "00000000-0000-0000-0000-000000000001"
}
```
