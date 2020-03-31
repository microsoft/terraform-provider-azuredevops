---
layout: "azuredevops"
page_title: "Provider: Azure DevOps"
description: |-
  The Azure DevOps provider is used to interact with Azure DevOps organization resources.
---

# Azure DevOps provider

The Azure DevOps provider can be used to configure Azure DevOps project in [Microsoft Azure](https://azure.microsoft.com/en-us/) using [Azure DevOps Service REST API](https://docs.microsoft.com/en-us/rest/api/azure/devops/?view=azure-devops-rest-5.1)

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
provider "azuredevops" {
  version = ">= 0.0.1"
}

resource "azuredevops_project" "project" {
  project_name       = "Project Name"
  description        = "Project Description"
}
```

# Argument Reference

The following arguments are supported in the `provider` block:

* `org_service_url` - (Required) This is the Azure DevOps organization url. It can also be
  sourced from the `AZDO_ORG_SERVICE_URL` environment variable. 

* `personal_access_token` - (Required) This is the Azure DevOps organization personal access 
  token. The account corresponding to the token will need "owner" privileges for this 
  organization. It can also be sourced from the `AZDO_PERSONAL_ACCESS_TOKEN` environment variable.
