---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_incomingwebhook"
description: |-
  Manages a Service Connection Incoming WebHook.
---

# azuredevops_serviceendpoint_incomingwebhook

Manages a Service Connection Incoming WebHook.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_incomingwebhook" "example" {
  project_id            = azuredevops_project.example.id
  webhook_name          = "example_webhook"
  secret                = "secret"
  http_header           = "X-Hub-Signature"
  service_endpoint_name = "Example IncomingWebhook"
  description           = "Managed by Terraform"
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project. Changing this forces a new Service Connection Incoming WebHook to be created.
* `webhook_name` - (Required) The name of the webhook.
* `secret` - (Optional) Secret for the webhook. WebHook service will use this secret to calculate the payload checksum.
* `http_header` - (Optional) Http header name on which checksum will be sent.
* `service_endpoint_name` - (Required) The name of the service endpoint. Changing this forces a new Service Connection Incoming WebHook to be created.
* `description` - (Optional) The Service Endpoint description. Defaults to Managed by Terraform.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Service Connection Incoming WebHook.
* `read` - (Defaults to 1 minute) Used when retrieving the Service Connection Incoming WebHook.
* `update` - (Defaults to 2 minutes) Used when updating the Service Connection Incoming WebHook.
* `delete` - (Defaults to 2 minutes) Used when deleting the Service Connection Incoming WebHook.

## Import

Service Connection Incoming WebHooks can be imported using the `resource id`, e.g.

```shell
terraform import azuredevops_serviceendpoint_incomingwebhook.example 00000000-0000-0000-0000-000000000000
```
