---
layout: "azuredevops"
page_title: "Provider: Azure DevOps"
description: |-
  The Azure DevOps provider is used to interact with Azure DevOps organization resources.
---

# Azure DevOps provider

The Azure DevOps provider can be used to configure Azure DevOps project in [Microsoft Azure](https://azure.microsoft.com/en-us/) using [Azure DevOps Service REST API](https://docs.microsoft.com/en-us/rest/api/azure/devops/?view=azure-devops-rest-7.0)

Use the navigation to the left to read about the available resources.

Interested in the provider's latest features, or want to make sure you're up to date? Check out the [changelog](https://github.com/microsoft/terraform-provider-azuredevops/blob/master/CHANGELOG.md) for version information and release notes.

## Example Usage

```hcl
terraform {
  required_providers {
    azuredevops = {
      source  = "microsoft/azuredevops"
      version = ">= 0.1.0"
    }
  }
}

resource "azuredevops_project" "project" {
  name        = "Project Name"
  description = "Project Description"
}
```

## Authentication

Authentication may be accomplished using an [Azure AD service principal](https://learn.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/service-principal-managed-identity) if your organization is connected to Entra ID, or by a [personal access token](https://learn.microsoft.com/en-us/azure/devops/organizations/accounts/use-personal-access-tokens-to-authenticate).

The provider will use the first available authentication method that is available. They are discovered in the following order:

* Personal Access Token
* With `use_oidc = true`
  * OIDC Token
  * OIDC Token File Path
  * OIDC Token Request URL
  * TFC Cloud Workload Identity Token
* Client Certificate Path
* Client Certificate
* Client Secret Path
* Client Secret
* With `use_msi = true`
  * Managed Service Identity

The OIDC service principal authentication methods allow for secure passwordless authentication from [Terraform Cloud](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/dynamic-provider-credentials) & [GitHub Actions](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect).

* [Authenticating to a Service Principal with a Managed Identity](guides/authenticating_managed_identity.html)
* [Authenticating to a Service Principal with a Client Certificate](guides/authenticating_service_principal_using_a_client_certificate.html)
* [Authenticating to a Service Principal with a Client Secret](guides/authenticating_service_principal_using_a_client_secret.html)
* [Authenticating to a Service Principal with an OIDC Token](guides/authenticating_service_principal_using_an_oidc_token.html)
* [Authenticating using a Personal Access Token](guides/authenticating_using_the_personal_access_token.html)

## Argument Reference

The following arguments are supported in the `provider` block:

- `org_service_url` - (Required) This is the Azure DevOps organization url. It can also be
  sourced from the `AZDO_ORG_SERVICE_URL` environment variable.

- `personal_access_token` - This is the Azure DevOps organization personal access
  token. The account corresponding to the token will need "owner" privileges for this
  organization. It can also be sourced from the `AZDO_PERSONAL_ACCESS_TOKEN` environment variable.

- `client_id` - The client id used when authenticating to a service principal or the principal id when
authenticating with a user specified managed service identity. It can also be sourced from
the `ARM_CLIENT_ID` or `AZURE_CLIENT_ID` environment variable.

- `client_id_file_path` - The path to a file containing a client id to authenticate. It can also be sourced from the `ARM_CLIENT_ID_FILE_PATH` environment variable.

- `tenant_id` - The tenant id used when authenticating to a service principal.
It can also be sourced from the `ARM_TENANT_ID` environment variable.

- `auxiliary_tenant_ids` - List of auxiliary Tenant IDs required for multi-tenancy and cross-tenant scenarios. This can also be sourced from the `ARM_AUXILIARY_TENANT_IDS` environment variable.

- `client_secret` - The client secret used to authenticate to a service principal.
It can also be sourced from the `ARM_CLIENT_SECRET` environment variable.

- `client_secret_path` - The path to a file containing a client secret to authenticate to a service principal.
It can also be sourced from the `ARM_CLIENT_SECRET_PATH` or `ARM_CLIENT_SECRET_FILE_PATH` environment variable.

- `client_certificate_path` - The path to a file containing a certificate to authenticate to a service
principal, typically a .pfx file.
It can also be sourced from the `ARM_CLIENT_CERTIFICATE_PATH` environment variable.

- `client_certificate` - A base64 encoded certificate to authentiate to a service principal.
It can also be sourced from the `ARM_CLIENT_CERTIFICATE` environment variable.

- `client_certificate_password` - This is the password associated with a certificate provided
by `client_certificate_path` or `client_certificate`. It can also be sourced
from the `ARM_CLIENT_CERTIFICATE_PASSWORD` environment variable.

- `oidc_token` - An OIDC token to authenticate to a service principal.
It can also be sourced from the `ARM_OIDC_TOKEN` environment variable.

- `oidc_token_file_path` - The path to a file containing nn OIDC token to authenticate to a service principal.
It can also be sourced from the `AZDO_TOKEN_PATH`, or `AZURE_FEDERATED_TOKEN_FILE` environment variable.

- `oidc_request_token` - The bearer token for the request to the OIDC provider. For use when authenticating as a Service Principal using OpenID Connect.
It can also be sourced from the `ARM_OIDC_REQUEST_TOKEN`, `ACTIONS_ID_TOKEN_REQUEST_TOKEN`, `SYSTEM_ACCESSTOKEN` environment variables.

- `oidc_request_url` - The URL for the OIDC provider from which to request an ID token. For use when authenticating as a Service Principal using OpenID Connect.
It can also be sourced from the `ARM_OIDC_REQUEST_URL`, `ACTIONS_ID_TOKEN_REQUEST_URL` or `SYSTEM_OIDCREQUESTURI` environment variables.

- `oidc_azure_service_connection_id` - The Azure Pipelines Service Connection ID to use for authentication. This can also be sourced from the `ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID`, `ARM_OIDC_AZURE_SERVICE_CONNECTION_ID`, or `AZURESUBSCRIPTION_SERVICE_CONNECTION_ID` environment variables.

- `use_oidc` - Boolean, enables OIDC auth methods. It can also be sourced from the `ARM_USE_OIDC` environment variable.

- `use_msi` - Boolean, enables authentication with a Managed Service Identity in Azure. It can also be sourced from the `ARM_USE_MSI` environment variable.

- `use_cli` - Should Azure CLI be used for authentication? This can also be sourced from the `ARM_USE_CLI` environment variable. Defaults to `true`.
