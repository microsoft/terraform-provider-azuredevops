---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_incomingwebhook"
description: |-
  Manages a Service Connection Incoming WebHook.
---

# azuredevops_serviceendpoint_incomingwebhook

Manages an Incoming WebHook service endpoint within Azure DevOps, which can be used as a resource in YAML pipelines to subscribe to a webhook event.

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
* `webhook_name` - (Required) The name of the WebHook.
* `secret` - (Optional) Secret for the WebHook. WebHook service will use this secret to calculate the payload checksum.
* `http_header` - (Optional) Http header name on which checksum will be sent.
* `service_endpoint_name` - (Required) The name of the service endpoint. Changing this forces a new Service Connection Incoming WebHook to be created.
* `description` - (Optional) The Service Endpoint description. Defaults to Managed by Terraform.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeout) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Incoming WebHook Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the Incoming WebHook Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the Incoming WebHook Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the Incoming WebHook Service Endpoint.

## Import

Azure DevOps Incoming WebHook Service Endpoint can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```shell
terraform import azuredevops_serviceendpoint_incomingwebhook.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
