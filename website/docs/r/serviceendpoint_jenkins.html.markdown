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
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "jenkins-example"
  description           = "Service Endpoint for 'Jenkins' (Managed by Terraform)"
  url                   = "https://example.com"
  accept_untrusted_certs  = false

  authentication_basic {
    username              = "username"
    password              = "password"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project. Changing this forces a new Service Connection Jenkins to be created.
* `service_endpoint_name` - (Required) The name of the service endpoint. Changing this forces a new Service Connection Jenkins to be created.
* `description` - (Optional) The Service Endpoint description. Defaults to Managed by Terraform.
* `url` - (Required) The Service Endpoint url.
* `accept_untrusted_certs` - (Optional) Allows the Jenkins clients to accept self-signed SSL server certificates.
* `username` - (Required) The Service Endpoint username to authenticate at the Jenkins Instance. 
* `password` - (Required) The Service Endpoint password to authenticate at the Jenkins Instance.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.

## Import

Service Connection Jenkins can be imported using the `resource id`, e.g.

```shell
terraform import azuredevops_serviceendpoint_jenkins.example 00000000-0000-0000-0000-000000000000
```
