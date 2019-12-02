# azuredevops_user_entitlement
Manages a user entitlement within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_user_entitlement" "user" {
    principal_name     = "foo@contoso.com"
}
```

## Argument Reference

* `principal_name` - (Optional) The principal name is the PrincipalName of a graph member from the source provider. Usually, e-mail address.
* `origin_id` - (Optional) The unique identifier from the system of origin. Typically a sid, object id or Guid. e.g. Used for member of other tenant on Azure Active Directory.
* `origin` - (Optional) The type of source provider for the origin identifier. Valid values: `aad` (Azure Active Directory) or `ghb` (GitHub). Defaults to `aad`.
* `account_license_type` - (Optional) Type of Account License. Valid values: `advanced`, `earlyAdopter`, `express`, `none`, `professional`, or `stakeholder`. Defaults to `express`.

**NOTE:** Set `principal_name` or `origin_id`. Set both values are not allowed.
**NOTE:** Currently `Update` is not supported. If you change these arguments, it will delete and create a new resource.

## Attributes Reference

The following attributes are exported:

* `id` - The userId of the User.
* `descriptor` - The descriptor is the primary way to reference the graph subject while the system is running. This field will uniqely identify the user graph subject.

## Relevant Links
* [Azure DevOps Service REST API 5.1 - User Entitlements - Add](https://docs.microsoft.com/en-us/rest/api/azure/devops/memberentitlementmanagement/user%20entitlements/add?view=azure-devops-rest-5.1)

## Import

Not supported.

## PAT Permissions Required

- **Member Entitlement Management**: Read & Write