---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_servicefabric"
description: |-
  Manages a Service Fabric service endpoint Azure DevOps organization.
---

# azuredevops_serviceendpoint_servicefabric

Manages a Service Fabric service endpoint within Azure DevOps.

## Example Usage

### Client Certificate Authentication

```hcl
data "azuredevops_project" "p" {
  name = "contoso"
}

resource "azuredevops_serviceendpoint_servicefabric" "test" {
  project_id            = data.azuredevops_project.p.id
  service_endpoint_name = "test"
  description           = "test"
  cluster_endpoint      = "tcp://test"

  certificate {
    server_certificate_lookup     = "Thumbprint"
    server_certificate_thumbprint = "test"
    client_certificate            = filebase64("${path.module}/certificate.pfx")
    client_certificate_password   = "password"
  }
}
```

### Azure Active Directory Authentication

```hcl
data "azuredevops_project" "p" {
  name = "contoso"
}

resource "azuredevops_serviceendpoint_servicefabric" "test" {
  project_id            = data.azuredevops_project.p.id
  service_endpoint_name = "test"
  description           = "test"
  cluster_endpoint      = "tcp://test"

  azure_active_directory {
    server_certificate_lookup     = "Thumbprint"
    server_certificate_thumbprint = "test"
    username                      = "username"
    password                      = "password"
  }
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The project ID or project name.
- `service_endpoint_name` - (Required) The Service Endpoint name.
- `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.
- One of either `certificate` or `azure_active_directory` blocks
- `certificate`
  - `server_certificate_lookup` - (Required) Verification mode for the cluster. Possible values include `Thumbprint` or `CommonName`.
  - `server_certificate_thumbprint` - (Optional) The thumbprint(s) of the cluster's certificate(s). This is used to verify the identity of the cluster. This value overrides the publish profile. Separate multiple thumbprints with a comma (',')
  - `server_certificate_common_name` - (Optional) The common name(s) of the cluster's certificate(s). This is used to verify the identity of the cluster. This value overrides the publish profile. Separate multiple common names with a comma (',')
  - `client_certificate` - (Required) Base64 encoding of the cluster's client certificate file.
  - `client_certificate_password` - (Required) Password for the certificate.
- `azure_active_directory`
  - `server_certificate_lookup` - (Required) Verification mode for the cluster. Possible values include `Thumbprint` or `CommonName`.
  - `server_certificate_thumbprint` - (Optional) The thumbprint(s) of the cluster's certificate(s). This is used to verify the identity of the cluster. This value overrides the publish profile. Separate multiple thumbprints with a comma (',')
  - `server_certificate_common_name` - (Optional) The common name(s) of the cluster's certificate(s). This is used to verify the identity of the cluster. This value overrides the publish profile. Separate multiple common names with a comma (',')
  - `username` - (Required) - Specify an Azure Active Directory account.
  - `password` - (Required) - Password for the Azure Active Directory account.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The project ID or project name.
- `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-5.1)

## Import

Azure DevOps Service Endpoint Service Fabric can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
 terraform import azuredevops_serviceendpoint_servicefabric.serviceendpoint 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
