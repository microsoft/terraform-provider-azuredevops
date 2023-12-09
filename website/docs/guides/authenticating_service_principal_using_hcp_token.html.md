---
layout: "azuredevops"
page_title: "Azure DevOps Provider: Authenticating to a Service Principal with a Terraform Cloud Workload Identity Token"
description: |-
  This guide will cover how to use a terraform cloud workload identity to authenticate to a service principal for use with Azure DevOps.
---

# Azure DevOps Provider: Authenticating to a Service Principal with a Terraform Cloud workload idenity token.

The Azure DevOps provider supports service principals through a variety of authentication methods, including the [OIDC identity token issued by Terraform Cloud](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/dynamic-provider-credentials).

## Service Principal Configuration

1. Create a service principal in [Azure portal](https://learn.microsoft.com/en-us/azure/active-directory/develop/howto-create-service-principal-portal) or
using [Azure PowerShell](https://learn.microsoft.com/en-us/azure/active-directory/develop/howto-authenticate-service-principal-powershell). Ignore steps about application roles and certificates.

2. Configure your app registration to trust your workspaces. On the Azure AD application page go to **Certificates & secrets**. Then click to the **Federated credentials** tab. Click **+ Add Credential**. Select **Other issuer** from the drop-down. The `issuer` is `https://app.terraform.io`, **make sure it starts with `https://` and does not have a trailing slash**. The `Subject Identifier` will be in the form `organization:my-org-name:project:my-project-name:workspace:my-workspace-name:run_phase:plan`. Note that the project will be `Default Project` if none is configured for your workspace. Both the plan and apply phase will need to be configured separately and may be on different service principals. Give your credential a name, this will not be changeable later. The audience will be `api://AzureADTokenExchange` by default, it must match the value configured in the Terraform workspace by setting the `TFC_WORKLOAD_IDENTITY_AUDIENCE` environment variable.

3. [Add the service principal to your Azure DevOps Organization.](https://learn.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/service-principal-managed-identity?view=azure-devops#2-add-and-manage-service-principal-in-an-azure-devops-organization)

4. Set the `TFC_WORKLOAD_IDENTITY_AUDIENCE` environment variable to `api://AzureADTokenExchange` in the Terraform cloud workspace, or a custom audience which you configured in #2. **This is required even if you intend to use the standard value.**

## Provider Configuration

The provider will need the Directory (tenant) ID and the Application (client) ID from the Azure AD app registration. They may be provided via the `ARM_TENANT_ID` and `ARM_CLIENT_ID` environment variables, or in the provider configuration block with the `tenant_id` and `client_id` attributes. Then the provider is configured to use the Terraform Cloud identity by either setting the `ARM_OIDC_HCP` environment variable to `true`, or the `oidc_hcp` provider attribute.

Separate service principals may used for the plan & apply phases by using the `ARM_TENANT_ID_PLAN`, `ARM_CLIENT_ID_PLAN`, `ARM_TENANT_ID_APPLY`, and `ARM_CLIENT_ID_APPLY` environment variables, or their respective provider attributes: `tenant_id_plan`, `client_id_plan`, `tenant_id_apply`, and `client_id_apply`.

### Configure the provider to authenticate with the Terraform Cloud workload idenity token

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

  client_id = "00000000-0000-0000-0000-000000000001"
  tenant_id = "00000000-0000-0000-0000-000000000001"
  oidc_hcp  = true
}

resource "azuredevops_project" "project" {
  name        = "Test Project"
  description = "Test Project Description"
}
```

### Configure the provider to authenticate with the Terraform Cloud workload idenity token with different plan & apply service principals

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

  client_id_plan  = "00000000-0000-0000-0000-000000000001"
  client_id_apply = "00000000-0000-0000-0000-000000000001"
  tenant_id_plan  = "00000000-0000-0000-0000-000000000001"
  tenant_id_apply = "00000000-0000-0000-0000-000000000001"
  oidc_hcp  = true
}

resource "azuredevops_project" "project" {
  name        = "Test Project"
  description = "Test Project Description"
}
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

  client_id                    = "00000000-0000-0000-0000-000000000001"
  tenant_id                    = "00000000-0000-0000-0000-000000000001"
  oidc_github_actions          = true
  oidc_github_actions_audience = "my-special-audience"
}

resource "azuredevops_project" "project" {
  name        = "Test Project"
  description = "Test Project Description"
}
```
