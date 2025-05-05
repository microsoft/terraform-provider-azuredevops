---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_jenkins"
description: |-
  Manages a Service Connection for Jenkins.
---

# azuredevops_serviceendpoint_jenkins

Manages a Jenkins service endpoint within Azure DevOps, which can be used as a resource in YAML pipelines to connect to a Jenkins instance.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_jenkins" "example" {
  project_id             = azuredevops_project.example.id
  service_endpoint_name  = "jenkins-example"
  description            = "Service Endpoint for 'Jenkins' (Managed by Terraform)"
  url                    = "https://example.com"
  accept_untrusted_certs = false
  username               = "username"
  password               = "password"

}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project. Changing this forces a new Service Connection Jenkins to be created.

* `service_endpoint_name` - (Required) The name of the service endpoint. Changing this forces a new Service Connection Jenkins to be created.

* `url` - (Required) The Service Endpoint url.

* `username` - (Required) The Service Endpoint username to authenticate at the Jenkins Instance.

* `password` - (Required) The Service Endpoint password to authenticate at the Jenkins Instance.

---

* `description` - (Optional) The Service Endpoint description. Defaults to Managed by Terraform.

* `accept_untrusted_certs` - (Optional) Allows the Jenkins clients to accept self-signed SSL server certificates. Defaults to `false.`


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Jenkins Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the Jenkins Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the Jenkins Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the Jenkins Service Endpoint.

## Import

Azure DevOps Jenkins Service Endpoint can be imported using the `projectId/id` or `projectName/id`, e.g.

```shell
terraform import azuredevops_serviceendpoint_jenkins.example projectName/00000000-0000-0000-0000-000000000000
```
