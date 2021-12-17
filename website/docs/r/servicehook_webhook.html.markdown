---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_servicehook_webhook"
description: |-
  Manages a Webhook service hook Azure DevOps organization.
---

# azuredevops_servicehook_webhook (Resource)

Manages a Webhook service hook Azure DevOps organization.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  name               = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "repo" {
  project_id = azuredevops_project.project.id
  name       = "Sample Empty Git Repository"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_servicehook_webhook" "webhook" {
  project_id = azuredevops_project.project.id
  event_type = "git.push"
  url        = "https://my-webhooks.org"

  # optional
  basic_auth {
    username = "my_username"
    password = "my_password"
  }

  # optional
  filters = {
    repository = azuredevops_git_repository.repo.id
  }

  # optional
  http_headers = {
    Authorization = "Bearer bearing"
    X-My-Header   = "header value"
  }
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The project ID or project name.
- `event_type` - (Required) Event type.
- `url` - (Required) The url of the hook to invoke.
- `basic_auth` - (Optional) Basic authentication.
  - `username`
  - `password`
- `filters` - (Optional) Filters that depend on event type.
- `http_headers` - (Optional) HTTP headers included in request.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service hook webhook.
- `project_id` - The project ID or project name.
- `event_type` - Event type.
- `url` - The url of the hook to invoke.

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Service hooks](https://docs.microsoft.com/en-us/rest/api/azure/devops/hooks/?view=azure-devops-rest-5.1)
