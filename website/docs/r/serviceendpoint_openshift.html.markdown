---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_openshift"
description: |-
  Manages an Openshift service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_openshift

Manages an Openshift service endpoint within Azure DevOps organization. Using this service endpoint requires you to first install the [OpenShift Extension](https://marketplace.visualstudio.com/items?itemName=redhat.openshift-vsts).

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_openshift" "example" {
  project_id                 = azuredevops_project.example.id
  service_endpoint_name      = "Example Openshift"
  server_url                 = "https://example.server"
  certificate_authority_file = "/opt/file"
  accept_untrusted_certs     = true
  auth_basic {
    username = "username"
    password = "password"
  }
}
```

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_openshift" "example" {
  project_id                 = azuredevops_project.example.id
  service_endpoint_name      = "Example Openshift"
  server_url                 = "https://example.server"
  certificate_authority_file = "/opt/file"
  accept_untrusted_certs     = true
  auth_token {
    token = "username"
  }
}
```
```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_openshift" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Openshift"
  server_url            = "https://example.server"
  auth_none {
    kube_config = "config"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_name` - (Required) The Service Endpoint name.

---

* `server_url` - (Optional) The URL for the OpenShift cluster to connect to.

* `accept_untrusted_certs` - (Optional) Set this option to allow clients to accept a self-signed certificate. Available when using `auth_basic` or `auth_token` authorization.

* `certificate_authority_file` - (Optional) The path to a certificate authority file to correctly and securely authenticates with an OpenShift server that uses HTTPS. Available when using `auth_basic` or `auth_token` authorization.

* `auth_basic` - (Optional) An `auth_basic` block as documented below.

* `auth_token` - (Optional) An `auth_token` block as documented below.

* `auth_none` - (Optional) An `auth_none` block as documented below.

* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

---

`auth_basic` block supports the following:

* `username` - (Required) The name of the user.

* `password` - (Required) The password of the user.

---

`auth_token` block supports the following:

* `token` - (Required) The API token.

---

`auth_none` block supports the following:

* `kube_config` - (Optional) The kubectl config


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The ID of the project.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Service Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Openshift Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the Openshift Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the Openshift Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the Openshift Service Endpoint.

## Import

Azure DevOps Openshift Service Endpoint can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_openshift.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
