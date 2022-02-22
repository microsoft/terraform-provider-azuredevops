---
layout: "azuredevops"
page_title: "Provider: Azure DevOps"
description: |-
  The Azure DevOps provider is used to interact with Azure DevOps organization resources.
---

# Azure DevOps provider

The Azure DevOps provider can be used to configure Azure DevOps project in [Microsoft Azure](https://azure.microsoft.com/en-us/) using [Azure DevOps Service REST API](https://docs.microsoft.com/en-us/rest/api/azure/devops/?view=azure-devops-rest-6.0)

Use the navigation to the left to read about the available resources.

Interested in the provider's latest features, or want to make sure you're up to date? Check out the [changelog](https://github.com/microsoft/terraform-provider-azuredevops/blob/master/CHANGELOG.md) for version information and release notes.

## Example Usage

```hcl
terraform {
  required_providers {
    azuredevops = {
      source = "microsoft/azuredevops"
      version = ">=0.1.0"
    }
  }
}

resource "azuredevops_project" "project" {
  name       = "Project Name"
  description        = "Project Description"
}
```

## Argument Reference

The following arguments are supported in the `provider` block:

- `org_service_url` - (Required) This is the Azure DevOps organization url. It can also be
  sourced from the `AZDO_ORG_SERVICE_URL` environment variable.

- `personal_access_token` - (Required) This is the Azure DevOps organization personal access
  token. The account corresponding to the token will need "owner" privileges for this
  organization. It can also be sourced from the `AZDO_PERSONAL_ACCESS_TOKEN` environment variable.
