---
layout: "azuredevops"
page_title: "Azure DevOps Provider: Authenticating to a Service Principal with a GitHub Actions OIDC Token"
description: |-
  This guide will cover how to use a github actions oidc token to authenticate to a service principal for use with Azure DevOps.
---

# Azure DevOps Provider: Authenticating to a Service Principal with a GitHub Actions OIDC Token

The Azure DevOps provider supports service principals through a variety of authentication methods, including the [OIDC identity token issued by GitHub Actions](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect).

## Service Principal Configuration

1. Create a service principal in [Azure portal](https://learn.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal) or
using [Azure PowerShell](https://learn.microsoft.com/en-us/azure/active-directory/develop/howto-authenticate-service-principal-powershell). Ignore steps about application roles and certificates.

2. [Configure your app registration to trust your workflow.](https://learn.microsoft.com/en-us/azure/active-directory/workload-identities/workload-identity-federation-create-trust?pivots=identity-wif-apps-methods-azp#github-actions)

3. [Add the service principal to your Azure DevOps Organization.](https://learn.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/service-principal-managed-identity?view=azure-devops#2-add-and-manage-service-principal-in-an-azure-devops-organization)

## Provider Configuration

The provider will need the Directory (tenant) ID and the Application (client) ID from the Azure AD app registration. They may be provided via the `AZDO_SP_TENANT_ID` and `AZDO_SP_CLIENT_ID` environment variables, or in the provider configuration block with the `sp_tenant_id` and `sp_client_id` attributes. Then the provider is configured to use the workflows identity by either setting the `AZDO_SP_OIDC_GITHUB_ACTIONS` environment variable to `true`, or the `sp_oidc_github_actions` provider attribute. The audience of the token may be customized by setting the `AZDO_SP_OIDC_GITHUB_ACTIONS_AUDIENCE` environment variable, or the `sp_oidc_github_actions_audience` provider attribute. The configured audience must match on the app registration and the Terraform provider configuration, the default values are normally acceptable.

The workflow, or specific terraform step in the workflow, must have the `id-token` permission which can be granted with:
```yaml
permissions:
  id-token: write
```

### Configure the provider to authenticate with the GitHub Action's identity token

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
  org_service_url                 = "https://dev.azure.com/my-org"

  sp_client_id                    = "00000000-0000-0000-0000-000000000001"
  sp_tenant_id                    = "00000000-0000-0000-0000-000000000001"
  sp_oidc_github_actions          = true
  sp_oidc_github_actions_audience = "my-special-audience"
}

resource "azuredevops_project" "project" {
  name        = "Test Project"
  description = "Test Project Description"
}
```
