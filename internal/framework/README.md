# Intro

The `framework` package is a wrapper around the native `Resource`, designed to accelerate the resource authoring and maintain the consistency, while still keeps all the features provided by the terraform-plugin-framework.

## Operation Flow

A wrapped resource has an opinionated flow for each operation.

### Create

- Underlying Create
- Create Poll, one or multiple times
- Set Identity
- Set State
- (Optional)
  - Post Create
  - Post Create Poll, one or multiple times
  - Set State

The underlying Create doesn't need to handle the protocol response if the req.Plan contains all the information for a follow-up Read(). Otherwise, the implement shall set those Read() related attributes to the state.

The `Create` and `Create Poll` are meant for creating the *main* resource, while the `Post Create` and `Post Create Poll` are meant for creating associated resources under the *main* resource, or updating the main resource for some properties that can only be configured once the resource is created. This split makes this *main* resource be tracked by terraform right after the `Create Poll`, since it is set to the state. If anything goes wrong after creating the *main* resource, then it will be marked as *tainted* in the state and the next `terraform apply` will recreate the resource properly.

### Read

- Underlying Read
- If doesn't exist remotely, remove the resource from state and quit
- Set Identity

The underlying Read must handle the protocol response (e.g. set the state).

### Update

- Underlying Update
- Update Poll, one or multiple times
- Set State
- (Optional)
  - Post Update
  - Post Update Poll, one or multiple times
  - Set State

The underlying Update doesn't need to handle the protocol response.

The `Update` and `Update Poll` are meant for updating the *main* resource, while the `Post Update` and `Post Update Poll` are meant for updating the associated resources under the *main* resource.

### Delete

- Underlying Delete
- (Optional) DeletePoll

The underlying Delete doesn't need to handle the protocol response.

## Configuration

The configuration of a resource can be implemented by simply embedding a struct `framework.ImplSetMeta` to the implementor of the `Resource` interface.

## Metadata

The resource metadata method is handled by the resource wrapper, the `Resource` implementor only needs to do:

- Implement the `ResourceType()` method
- Embed the dummy `framework.ImplMetadata` struct (which does nothing, but to meet the interface requirement)

## Timeout

The timeout feature is handled by the `framework`. By default, it will define a timeout as below for the resource:

- Create: 5min
- Read: 5min
- Update: 5min
- Delete: 5min

A resource only needs to define a `Timeouts` field in the data model, e.g.:

        type myModel struct {
            ...
            Timeouts timeouts.Value `tfsdk:"timeouts"`
        }

The `ctx` passed to the Create/Read/Update/Delete has the deadline already set.

Otherwise, if the default timeout doesn't work for a resource, the author can simply implement the interface `ResourceWithTimeout` to define the proper timeouts.

## Log

The user can simply embed a struct `framework.ImplLog[*myResource]` in the implementor of the `Resource` interface, i.e. `myResource`. This struct will then implement the log related methods:

- Info
- Warn
- Error

These are properly setup with the terraform log subsystem.

Additionally, the resource wrapper will emit lifecycle related logs, e.g. `Start to create the resource`.

## Identity

Each wrapped resource will support "Resource Identity" at the first place, with the following efforts:

- Implement the `IdentitySchema()` method to define the identity schema
- Define an identity model that implements the `ResourceIdentity` interface
- Implement the `Identy()` method to return the identity model

With above, the framework will be able to properly set the identity during Create and Read.

## Import

This requires the `Identity` setup above, which will then implement the `ImportState()` method automatically, with the support for both import by id and import by identity.
