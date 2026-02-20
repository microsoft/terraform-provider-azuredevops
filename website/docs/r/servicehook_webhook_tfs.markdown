---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_servicehook_webhook_tfs"
description: |-
  Manages a Webhook TFS Service Hook.
---

# azuredevops_servicehook_webhook_tfs

Manages a Webhook TFS Service Hook that sends HTTP POST requests to a specified URL when Azure DevOps events occur.

## Example Usage

### Git Push Event

```hcl
resource "azuredevops_project" "example" {
  name = "example-project"
}

resource "azuredevops_servicehook_webhook_tfs" "example" {
  project_id = azuredevops_project.example.id
  url        = "https://example.com/webhook"
  
  git_push {
    branch        = "refs/heads/main"
    repository_id = azuredevops_git_repository.example.id
  }
}
```

### Build Completed Event with Authentication

```hcl
resource "azuredevops_servicehook_webhook_tfs" "example" {
  project_id           = azuredevops_project.example.id
  url                  = "https://example.com/webhook"
  basic_auth_username  = "webhook_user"
  basic_auth_password  = var.webhook_password
  accept_untrusted_certs = false
  
  build_completed {
    definition_name = "CI Build"
    build_status    = "Succeeded"
  }
}
```

### Pull Request Created Event with HTTP Headers

```hcl
resource "azuredevops_servicehook_webhook_tfs" "example" {
  project_id = azuredevops_project.example.id
  url        = "https://example.com/webhook"
  
  http_headers = {
    "X-Custom-Header" = "my-value"
    "Authorization"   = "Bearer ${var.api_token}"
  }
  
  git_pull_request_created {
    repository_id = azuredevops_git_repository.example.id
    branch        = "refs/heads/develop"
  }
}
```

### Work Item Updated Event

```hcl
resource "azuredevops_servicehook_webhook_tfs" "example" {
  project_id                = azuredevops_project.example.id
  url                       = "https://example.com/webhook"
  resource_details_to_send  = "all"
  messages_to_send          = "text"
  detailed_messages_to_send = "markdown"
  
  work_item_updated {
    work_item_type = "Bug"
    area_path      = "MyProject\\Area"
    changed_fields = "System.State"
  }
}
```

An empty configuration block will trigger on all events of that type:

```hcl
resource "azuredevops_servicehook_webhook_tfs" "example" {
  project_id = azuredevops_project.example.id
  url        = "https://example.com/webhook"
  
  git_push {}
}
```


## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project. Changing this forces a new Service Hook Webhook TFS to be created.

* `url` - (Required) The URL to send HTTP POST to.

---

* `accept_untrusted_certs` - (Optional) Accept untrusted SSL certificates. Defaults to `false`.

* `basic_auth_username` - (Optional) Basic authentication username.

* `basic_auth_password` - (Optional) Basic authentication password.

* `http_headers` - (Optional) HTTP headers as key-value pairs to include in the webhook request.

* `resource_details_to_send` - (Optional) Resource details to send - `all`, `minimal`, or `none`. Defaults to `all`.

* `messages_to_send` - (Optional) Messages to send - `all`, `text`, `html`, `markdown`, or `none`. Defaults to `all`.

* `detailed_messages_to_send` - (Optional) Detailed messages to send - `all`, `text`, `html`, `markdown`, or `none`. Defaults to `all`.

* `resource_version` - (Optional) The resource version for the webhook subscription. Defaults to `latest`.

---

### Event Types

Exactly one of the following event type blocks must be specified:

* `build_completed` - (Optional) A `build_completed` block as defined below.

* `git_pull_request_commented` - (Optional) A `git_pull_request_commented` block as defined below.

* `git_pull_request_created` - (Optional) A `git_pull_request_created` block as defined below.

* `git_pull_request_merge_attempted` - (Optional) A `git_pull_request_merge_attempted` block as defined below.

* `git_pull_request_updated` - (Optional) A `git_pull_request_updated` block as defined below.

* `git_push` - (Optional) A `git_push` block as defined below.

* `repository_created` - (Optional) A `repository_created` block as defined below.

* `repository_deleted` - (Optional) A `repository_deleted` block as defined below.

* `repository_forked` - (Optional) A `repository_forked` block as defined below.

* `repository_renamed` - (Optional) A `repository_renamed` block as defined below.

* `repository_status_changed` - (Optional) A `repository_status_changed` block as defined below.

* `service_connection_created` - (Optional) A `service_connection_created` block as defined below.

* `service_connection_updated` - (Optional) A `service_connection_updated` block as defined below.

* `tfvc_checkin` - (Optional) A `tfvc_checkin` block as defined below.

* `work_item_commented` - (Optional) A `work_item_commented` block as defined below.

* `work_item_created` - (Optional) A `work_item_created` block as defined below.

* `work_item_deleted` - (Optional) A `work_item_deleted` block as defined below.

* `work_item_restored` - (Optional) A `work_item_restored` block as defined below.

* `work_item_updated` - (Optional) A `work_item_updated` block as defined below.

---

A `build_completed` block supports the following:

* `definition_name` - (Optional) Include only events for completed builds for a specific pipeline.

* `build_status` - (Optional) Include only events for completed builds that have a specific completion status. Valid values: `Succeeded`, `PartiallySucceeded`, `Failed`, `Stopped`.

---

A `git_pull_request_commented` block supports the following:

* `repository_id` - (Optional) Include only events for pull requests in a specific repository (repository ID). If not specified, all repositories in the project will trigger the event.

* `branch` - (Optional) Include only events for pull requests in a specific branch.

---

A `git_pull_request_created` block supports the following:

* `repository_id` - (Optional) Include only events for pull requests in a specific repository (repository ID). If not specified, all repositories in the project will trigger the event.

* `branch` - (Optional) Include only events for pull requests in a specific branch.

* `pull_request_created_by` - (Optional) Include only events for pull requests created by users in a specific group.

* `pull_request_reviewers_contains` - (Optional) Include only events for pull requests with reviewers in a specific group.

---

A `git_pull_request_merge_attempted` block supports the following:

* `repository_id` - (Optional) Include only events for pull requests in a specific repository (repository ID). If not specified, all repositories in the project will trigger the event.

* `branch` - (Optional) Include only events for pull requests in a specific branch.

* `pull_request_created_by` - (Optional) Include only events for pull requests created by users in a specific group.

* `pull_request_reviewers_contains` - (Optional) Include only events for pull requests with reviewers in a specific group.

* `merge_result` - (Optional) Include only events for pull requests with a specific merge result. Valid values: `Succeeded`, `Unsuccessful`, `Conflicts`, `Failure`, `RejectedByPolicy`.

---

A `git_pull_request_updated` block supports the following:

* `repository_id` - (Optional) Include only events for pull requests in a specific repository (repository ID). If not specified, all repositories in the project will trigger the event.

* `branch` - (Optional) Include only events for pull requests in a specific branch.

* `pull_request_created_by` - (Optional) Include only events for pull requests created by users in a specific group.

* `pull_request_reviewers_contains` - (Optional) Include only events for pull requests with reviewers in a specific group.

* `notification_type` - (Optional) Include only events for pull requests with a specific change. Valid values: `PushNotification`, `ReviewersUpdateNotification`, `StatusUpdateNotification`, `ReviewerVoteNotification`.

---

A `git_push` block supports the following:

* `repository_id` - (Optional) Include only events for code pushes to a specific repository (repository ID). If not specified, all repositories in the project will trigger the event.

* `branch` - (Optional) Include only events for code pushes to a specific branch.

* `pushed_by` - (Optional) Include only events for code pushes by users in a specific group.

---

A `repository_created` block supports the following:

* `project_id` - (Optional) Include only events for repositories created in a specific project.

---

A `repository_deleted` block supports the following:

* `repository_id` - (Optional) Include only events for repositories with a specific repository ID.

---

A `repository_forked` block supports the following:

* `repository_id` - (Optional) Include only events for repositories with a specific repository ID.

---

A `repository_renamed` block supports the following:

* `repository_id` - (Optional) Include only events for repositories with a specific repository ID.

---

A `repository_status_changed` block supports the following:

* `repository_id` - (Optional) Include only events for repositories with a specific repository ID.

---

A `service_connection_created` block supports the following:

* `project_id` - (Optional) Include only events for service connections created in a specific project.

---

A `service_connection_updated` block supports the following:

* `project_id` - (Optional) Include only events for service connections updated in a specific project.

---

A `tfvc_checkin` block supports the following:

* `path` - (Required) Include only events for check-ins that change files under a specific path.

---

A `work_item_commented` block supports the following:

* `work_item_type` - (Optional) Include only events for work items of a specific type.

* `area_path` - (Optional) Include only events for work items under a specific area path.

* `tag` - (Optional) Include only events for work items that contain a specific tag.

* `comment_pattern` - (Optional) Include only events for work items with a comment that contains a specific string.

---

A `work_item_created` block supports the following:

* `work_item_type` - (Optional) Include only events for work items of a specific type.

* `area_path` - (Optional) Include only events for work items under a specific area path.

* `tag` - (Optional) Include only events for work items that contain a specific tag.

---

A `work_item_deleted` block supports the following:

* `work_item_type` - (Optional) Include only events for work items of a specific type.

* `area_path` - (Optional) Include only events for work items under a specific area path.

* `tag` - (Optional) Include only events for work items that contain a specific tag.

---

A `work_item_restored` block supports the following:

* `work_item_type` - (Optional) Include only events for work items of a specific type.

* `area_path` - (Optional) Include only events for work items under a specific area path.

* `tag` - (Optional) Include only events for work items that contain a specific tag.

---

A `work_item_updated` block supports the following:

* `work_item_type` - (Optional) Include only events for work items of a specific type.

* `area_path` - (Optional) Include only events for work items under a specific area path.

* `tag` - (Optional) Include only events for work items that contain a specific tag.

* `changed_fields` - (Optional) Include only events for work items with a change in a specific field.

* `links_changed` - (Optional) Include only events for work items with one or more links added or removed.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Hook Webhook TFS.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Webhook TFS Service Hook.
* `read` - (Defaults to 5 minutes) Used when retrieving the Webhook TFS Service Hook.
* `update` - (Defaults to 10 minutes) Used when updating the Webhook TFS Service Hook.
* `delete` - (Defaults to 10 minutes) Used when deleting the Webhook TFS Service Hook.

## Import

Webhook TFS Service Hook can be imported using the `resource id`, e.g.

```shell
terraform import azuredevops_servicehook_webhook_tfs.example 00000000-0000-0000-0000-000000000000
```
