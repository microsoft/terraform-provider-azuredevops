---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_workload_identity"
description: |-
  Manages a native Azure DevOps service connection using a Microsoft Entra workload identity.
---

# azuredevops_serviceendpoint_workload_identity

Manages an Azure DevOps service connection of type `workloadidentityuser`, using workload identity federation with an existing app registration or user-assigned managed identity. Pipelines can use this connection to access Azure DevOps repositories, artifacts and REST APIs without a PAT.

This is the native **Azure DevOps** connection type. The deprecated `azuredevops_serviceendpoint_azuredevops` resource manages the `AZDOAPI` connection supplied by the Configurable Pipeline Runner extension. `azuredevops_serviceendpoint_azurerm` manages a different connection type for Azure Resource Manager.

~> **Note** This connection type is available in Azure DevOps Services. The identity must already exist in Microsoft Entra ID and be added to Azure DevOps with the permissions needed to access the target resources. Cross-organization connections require the target organization to use the same Entra tenant and the identity to have access there.

## Example Usage

### App registration

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_service_principal_entitlement" "example" {
  origin_id = "00000000-0000-0000-0000-000000000001"
}

resource "azuredevops_serviceendpoint_workload_identity" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Azure DevOps workload identity"
  org_url               = "https://dev.azure.com/example"
  tenant_id             = "00000000-0000-0000-0000-000000000002"
  service_principal_id  = "00000000-0000-0000-0000-000000000003"
  object_id             = "00000000-0000-0000-0000-000000000004"
  subject_descriptor    = azuredevops_service_principal_entitlement.example.descriptor
}
```

`service_principal_id` is the application (client) ID. For `Application`, `object_id` is the **app registration object ID**, which Azure DevOps uses to create its federated credential. Organization membership uses the separate **enterprise application/service principal object ID** (`origin_id` above), and `subject_descriptor` refers to that membership. For `ManagedIdentity`, `object_id` is the managed identity's **principal object ID**. Grant the identity access to the relevant projects and resources separately.

### User-assigned managed identity

```hcl
resource "azuredevops_serviceendpoint_workload_identity" "example" {
  project_id                   = "00000000-0000-0000-0000-000000000001"
  service_endpoint_name        = "Azure DevOps managed identity"
  org_url                      = "https://dev.azure.com/example"
  tenant_id                    = "00000000-0000-0000-0000-000000000002"
  service_principal_id         = "00000000-0000-0000-0000-000000000003"
  object_id                    = "00000000-0000-0000-0000-000000000004"
  subject_descriptor           = "aadsp.example"
  service_principal_type       = "ManagedIdentity"
  managed_identity_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000005/resourceGroups/example/providers/Microsoft.ManagedIdentity/userAssignedIdentities/example"
}
```

### Federated credential setup and drafts

Azure DevOps attempts to configure the federated identity credential on the existing identity. The caller needs the corresponding Entra or Azure permissions. The service can save a draft if automatic credential creation or verification fails, including when called with some service principal authentication configurations. Provider support does not bypass these upstream permission requirements.

For a manual federation setup, use two applies:

1. Set `draft = true` and apply. The provider requests a draft using `data.isDraft` and `isReady`, preserves its ID, and exports the server-generated issuer and subject. A draft is not a usable pipeline connection; check `is_ready`.
2. Configure a federated identity credential on the app registration or managed identity using `workload_identity_federation_issuer`, `workload_identity_federation_subject` and the audience `api://AzureADTokenExchange`. After the credential exists and has propagated, set `draft = false` and apply again to finalize the same connection.

`draft` controls whether an unfinished connection is acceptable during create/update; `is_ready` reports the actual server state. Updates use the API's draft-validation path even when `draft = false`; Azure DevOps finalizes the connection after successful validation. Draft mode does not disable Azure DevOps' automatic federation behavior, and unrelated API failures still produce an error. If the service has already created the credential automatically, manage or import that credential instead of creating a duplicate.

Azure DevOps can return an HTTP error after saving a draft, including while a new federated credential propagates. For recognized federation/draft errors, the provider retrieves the matching draft and records its ID. With `draft = true`, it reports a warning so that federation can be completed without recreating the connection. When readiness is required, known Entra federated-credential propagation errors are retried within the configured create/update timeout. Other API errors still fail. With `draft = false`, a saved draft produces an error and Terraform can mark the resource as tainted. Inspect the connection and state before retrying: an ordinary retry can replace a tainted resource. Start with `draft = true` when a manual federation setup is expected.

Renaming an already-ready app-registration connection can cause Azure DevOps to remove or recreate its federated credential. This operation can fail under app-only provider authentication even when manual federation and finalization succeeded. It may require delegated authentication with the corresponding Entra permissions. If an update leaves a draft while `draft = false`, the next plan schedules validation again.

The resource does not create an app registration, managed identity, Azure subscription or organization membership. Deletion uses `deep=false` to avoid asking Azure DevOps to delete the backing service principal. Manage separately owned federated credentials through the appropriate Entra/Azure provider.

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) ID of the project that owns the service connection. Changing this forces a new resource.
* `service_endpoint_name` - (Required) Name of the service connection.
* `org_url` - (Required) URL of the target Azure DevOps organization. Changing this forces a new resource.
* `tenant_id` - (Required) Microsoft Entra tenant ID. Changing this forces a new resource.
* `service_principal_id` - (Required) Application (client) ID of the app registration or managed identity. Changing this forces a new resource.
* `object_id` - (Required) App registration object ID for `Application`, or principal object ID for `ManagedIdentity`. Changing this forces a new resource.
* `subject_descriptor` - (Optional) Azure DevOps graph descriptor of the identity in the organization. Available from `azuredevops_service_principal_entitlement.descriptor`, or resolved by Azure DevOps when omitted. Changing this forces a new resource.
* `service_principal_type` - (Optional) `Application` or `ManagedIdentity`. Defaults to `Application`. Changing this forces a new resource.
* `managed_identity_resource_id` - (Optional) Azure resource ID of the user-assigned managed identity. Only permitted for `ManagedIdentity`. When omitted, Azure DevOps attempts to resolve it through Azure Resource Graph; supplying it explicitly avoids depending on that discovery. Changing this forces a new resource.
* `description` - (Optional) Description of the connection. Defaults to `Managed by Terraform`.
* `draft` - (Optional) Allow a draft connection while federation is configured. Defaults to `false`. Set to `false` after configuring the credential to finalize the connection.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the service connection.
* `authorization` - Authorization metadata, including `scheme = WorkloadIdentityFederation`.
* `is_ready` - Whether Azure DevOps reports the connection as ready.
* `workload_identity_federation_issuer` - Issuer returned by Azure DevOps for the federated identity credential.
* `workload_identity_federation_subject` - Subject returned by Azure DevOps for the federated identity credential.

## Timeouts

* `create` - (Defaults to 2 minutes) Used when creating the connection.
* `read` - (Defaults to 1 minute) Used when reading the connection.
* `update` - (Defaults to 2 minutes) Used when updating or finalizing the connection.
* `delete` - (Defaults to 2 minutes) Used when deleting the connection.

## Import

Connections can be imported using `projectID/serviceEndpointID` or `projectName/serviceEndpointID`:

```sh
terraform import azuredevops_serviceendpoint_workload_identity.example 00000000-0000-0000-0000-000000000001/00000000-0000-0000-0000-000000000002
```

Only native `workloadidentityuser` connections can be imported. For an unfinished connection, configure `draft = true` while setting up federation, then apply with `draft = false` to finalize it.

## Relevant Links

* [Access Azure DevOps with Microsoft Entra workload identity](https://learn.microsoft.com/en-us/azure/devops/pipelines/library/add-devops-entra-service-connection?view=azure-devops)
* [Service endpoint types REST API](https://learn.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/types/get?view=azure-devops-rest-7.1)
* [Upstream feature request and authentication reports](https://github.com/microsoft/terraform-provider-azuredevops/issues/1577)
