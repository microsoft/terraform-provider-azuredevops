package framework

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-azuredevops/internal/meta"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/ctxutil"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/errorutil"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/retry"
)

var _ resource.Resource = resourceWrapper{}
var _ resource.ResourceWithConfigure = resourceWrapper{}
var _ resource.ResourceWithImportState = resourceWrapper{}
var _ resource.ResourceWithConfigValidators = resourceWrapper{}
var _ resource.ResourceWithModifyPlan = resourceWrapper{}
var _ resource.ResourceWithMoveState = resourceWrapper{}
var _ resource.ResourceWithUpgradeState = resourceWrapper{}
var _ resource.ResourceWithValidateConfig = resourceWrapper{}
var _ resource.ResourceWithUpgradeIdentity = resourceWrapper{}
var _ ResourceWithTimeout = resourceWrapper{}

// The followings are unsafe interfaces. This requires additional wrappers around this resourceWrapper and opt-in.

type resourceWrapper struct {
	Resource
}

func WrapResource(in Resource) func() resource.Resource {
	return func() resource.Resource {
		return resourceWrapper{Resource: in}
	}
}

func (r resourceWrapper) Timeout() ResourceTimeout {
	if r, ok := r.Resource.(ResourceWithTimeout); ok {
		return r.Timeout()
	}
	return ResourceTimeout{
		Create: 5 * time.Minute,
		Read:   5 * time.Minute,
		Update: 5 * time.Minute,
		Delete: 5 * time.Minute,
	}
}

func (r resourceWrapper) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.Resource.Metadata(ctx, req, resp)
	resp.TypeName = r.Resource.ResourceType()
}

func (r resourceWrapper) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.Resource.Schema(ctx, req, resp)
	timeout := r.Timeout()
	resp.Schema.Attributes["timeouts"] = timeouts.Attributes(ctx, timeouts.Opts{
		Create:            true,
		Read:              true,
		Update:            true,
		Delete:            true,
		CreateDescription: fmt.Sprintf("(Defaults to %s) Used when creating this resource.", timeout.Create),
		ReadDescription:   fmt.Sprintf("(Defaults to %s) Used when reading this resource.", timeout.Read),
		UpdateDescription: fmt.Sprintf("(Defaults to %s) Used when updating this resource.", timeout.Update),
		DeleteDescription: fmt.Sprintf("(Defaults to %s) Used when deleting this resource.", timeout.Delete),
	})
}

func (r resourceWrapper) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	identity := r.Identity()
	if req.ID != "" {
		// Import via ID
		if err := identity.FromId(req.ID); err != nil {
			resp.Diagnostics.AddError("Converting identity from id string", err.Error())
			return
		}
	} else {
		// Import via Identity
		resp.Diagnostics.Append(req.Identity.Get(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	for _, field := range identity.Fields() {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, field.PathState, field.Value)...)
	}
}

func (r resourceWrapper) WritePoll(ctx context.Context, operation WriteOperation, plan tfsdk.Plan, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), fmt.Sprintf("Start to poll the resource (after %s)", operation))
	defer tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), fmt.Sprintf("Finish to poll the resource (after %s)", operation))
	var (
		pollOption = func(ctx context.Context) retry.RetryOption {
			return retry.RetryOption{Timeout: ctxutil.UntilDeadline(ctx)}
		}
		pollCheck = func(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) error {
			return nil
		}
		pollRetryableDiags = func(diag.Diagnostics) bool {
			return false
		}
	)

	switch operation {
	case WriteOperationCreate:
		if r, ok := r.Resource.(ResourceWithCreatePoll); ok {
			pollOption = r.CreatePollOption
			pollCheck = r.CreatePollCheck
			pollRetryableDiags = r.CreatePollRetryableDiags
		}
	case WriteOperationPostCreate:
		if r, ok := r.Resource.(ResourceWithPostCreatePoll); ok {
			pollOption = r.PostCreatePollOption
			pollCheck = r.PostCreatePollCheck
			pollRetryableDiags = r.PostCreatePollRetryableDiags
		}
	case WriteOperationUpdate:
		if r, ok := r.Resource.(ResourceWithUpdatePoll); ok {
			pollOption = r.UpdatePollOption
			pollCheck = r.UpdatePollCheck
			pollRetryableDiags = r.UpdatePollRetryableDiags
		}
	case WriteOperationPostUpdate:
		if r, ok := r.Resource.(ResourceWithPostUpdatePoll); ok {
			pollOption = r.PostUpdatePollOption
			pollCheck = r.PostUpdatePollCheck
			pollRetryableDiags = r.PostUpdatePollRetryableDiags
		}
	default:
		panic(fmt.Sprintf("unknown operation for polling: %s", operation))
	}

	oldDiags := resp.Diagnostics
	if err := retry.RetryContext(ctx, pollOption(ctx), func() *retry.RetryError {
		tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), fmt.Sprintf("Start to read the resource (after %s)", operation))

		// Reset rresp's diags before every retry
		resp.Diagnostics = slices.Clone(oldDiags)

		r.Resource.Read(ctx, req, resp)
		if resp.Diagnostics.HasError() {
			if pollRetryableDiags(resp.Diagnostics) {
				return retry.RetryableError(errorutil.DiagsToError(resp.Diagnostics))
			} else {
				return retry.NonRetryableError(errorutil.DiagsToError(resp.Diagnostics))
			}
		}

		tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), fmt.Sprintf("Finish to read the resource (after %s)", operation))
		if err := pollCheck(ctx, plan, resp.State); err != nil {
			return retry.RetryableError(err)
		}
		return nil
	}); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Polling failed (after %s)", operation), err.Error())
		return
	}
}

func (r resourceWrapper) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var timeout timeouts.Value
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("timeouts"), &timeout)...)
	if resp.Diagnostics.HasError() {
		return
	}

	duration, diags := timeout.Create(ctx, r.Timeout().Create)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	ctx = tflog.NewSubsystem(ctx, r.Resource.ResourceType())
	ctx = tflog.SubsystemSetField(ctx, r.Resource.ResourceType(), "operation", "Create")

	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Start to create the resource")
	r.Resource.Create(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Finish to create the resource")

	// If the inner Create() doesn't set the state, temporarily set the plan to state, so that we can use the state to construct the read request below.
	if resp.State.Raw.IsNull() {
		resp.Diagnostics.Append(resp.State.Set(ctx, req.Plan.Raw)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	rreq := resource.ReadRequest{
		State:              resp.State,
		Private:            resp.Private,
		Identity:           resp.Identity,
		ProviderMeta:       req.ProviderMeta,
		ClientCapabilities: resource.ReadClientCapabilities{},
	}
	rresp := resource.ReadResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
		Identity:    resp.Identity,
		Private:     resp.Private,
		Deferred:    nil,
	}

	// Create Poll
	r.WritePoll(ctx, WriteOperationCreate, req.Plan, rreq, &rresp)
	*resp = resource.CreateResponse{
		State:       rresp.State,
		Identity:    rresp.Identity,
		Private:     rresp.Private,
		Diagnostics: rresp.Diagnostics,
	}

	// Set the identity
	rresp.Diagnostics = append(rresp.Diagnostics, r.setIdentity(ctx, rresp.State, rresp.Identity)...)

	// Post Create
	if rr, ok := r.Resource.(ResourceWithPostCreate); ok && rr.ShouldPostCreate(ctx, req) {
		tflog.SubsystemInfo(ctx, r.ResourceType(), "Start to post create the resource")
		rr.PostCreate(ctx, req, resp)
		if resp.Diagnostics.HasError() {
			return
		}
		tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Finish to post create the resource")

		r.WritePoll(ctx, WriteOperationPostCreate, req.Plan, rreq, &rresp)
		*resp = resource.CreateResponse{
			State:       rresp.State,
			Identity:    rresp.Identity,
			Private:     rresp.Private,
			Diagnostics: rresp.Diagnostics,
		}
	}
}

func (r resourceWrapper) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var timeout timeouts.Value
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("timeouts"), &timeout)...)
	if resp.Diagnostics.HasError() {
		return
	}

	duration, diags := timeout.Read(ctx, r.Timeout().Read)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	ctx = tflog.NewSubsystem(ctx, r.Resource.ResourceType())
	ctx = tflog.SubsystemSetField(ctx, r.Resource.ResourceType(), "operation", "Read")

	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Start to read the resource")
	r.Resource.Read(ctx, req, resp)
	// If the resource doesn't exist, remove it from the state and return.
	if slices.ContainsFunc(resp.Diagnostics, IsDiagResourceNotFound) {
		tflog.SubsystemWarn(ctx, r.Resource.ResourceType(), "Resource not found, removing it from the state and return")
		resp.Diagnostics = slices.DeleteFunc(resp.Diagnostics, IsDiagResourceNotFound)
		resp.State.RemoveResource(ctx)

		// Set the identity to avoid error message about lacking of identity after successfully returning from read.
		// This happens when a resource that has no identity before, and it has disappeared in remote.
		// We assume the "req" has the adequate state to form the identity.
		resp.Diagnostics = append(resp.Diagnostics, r.setIdentity(ctx, req.State, resp.Identity)...)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Finish to read the resource")

	// Set the identity
	resp.Diagnostics = append(resp.Diagnostics, r.setIdentity(ctx, resp.State, resp.Identity)...)
}

func (r resourceWrapper) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var timeout timeouts.Value
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("timeouts"), &timeout)...)
	if resp.Diagnostics.HasError() {
		return
	}

	duration, diags := timeout.Update(ctx, r.Timeout().Update)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	ctx = tflog.NewSubsystem(ctx, r.Resource.ResourceType())
	ctx = tflog.SubsystemSetField(ctx, r.Resource.ResourceType(), "operation", "Update")

	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Start to update the resource")
	r.Resource.Update(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Finish to update the resource")

	// Temporarily set the plan to state, so that we can use the state to construct the read request below.
	resp.Diagnostics.Append(resp.State.Set(ctx, req.Plan.Raw)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rreq := resource.ReadRequest{
		State:              resp.State,
		Private:            resp.Private,
		Identity:           resp.Identity,
		ProviderMeta:       req.ProviderMeta,
		ClientCapabilities: resource.ReadClientCapabilities{},
	}

	rresp := resource.ReadResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
		Identity:    resp.Identity,
		Private:     resp.Private,
		Deferred:    nil,
	}

	// Update Poll
	r.WritePoll(ctx, WriteOperationUpdate, req.Plan, rreq, &rresp)
	*resp = resource.UpdateResponse{
		State:       rresp.State,
		Identity:    rresp.Identity,
		Private:     rresp.Private,
		Diagnostics: rresp.Diagnostics,
	}

	// Post Update
	if rr, ok := r.Resource.(ResourceWithPostUpdate); ok && rr.ShouldPostUpdate(ctx, req) {
		tflog.SubsystemInfo(ctx, r.ResourceType(), "Start to post update the resource")
		rr.PostUpdate(ctx, req, resp)
		if resp.Diagnostics.HasError() {
			return
		}
		tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Finish to post update the resource")

		r.WritePoll(ctx, WriteOperationPostUpdate, req.Plan, rreq, &rresp)
		*resp = resource.UpdateResponse{
			State:       rresp.State,
			Identity:    rresp.Identity,
			Private:     rresp.Private,
			Diagnostics: rresp.Diagnostics,
		}
	}
}

func (r resourceWrapper) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var timeout timeouts.Value
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("timeouts"), &timeout)...)
	if resp.Diagnostics.HasError() {
		return
	}

	duration, diags := timeout.Delete(ctx, r.Timeout().Delete)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	ctx = tflog.NewSubsystem(ctx, r.Resource.ResourceType())
	ctx = tflog.SubsystemSetField(ctx, r.Resource.ResourceType(), "operation", "Delete")

	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Start to delete the resource")
	r.Resource.Delete(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Finish to delete the resource")

	if rr, ok := r.Resource.(ResourceWithDeletePoll); ok {
		rreq := resource.ReadRequest{
			State:              req.State,
			Private:            req.Private,
			Identity:           req.Identity,
			ProviderMeta:       req.ProviderMeta,
			ClientCapabilities: resource.ReadClientCapabilities{},
		}

		rresp := resource.ReadResponse{
			State:       resp.State,
			Diagnostics: slices.Clone(resp.Diagnostics),
			Identity:    resp.Identity,
			Private:     resp.Private,
			Deferred:    nil,
		}

		retryOption := rr.DeletePollOption(ctx)
		if err := retry.RetryContext(ctx, retryOption, func() *retry.RetryError {
			tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Start to read the resource (post-delete)")

			// Reset rresp before every retry
			rresp = resource.ReadResponse{
				State:       resp.State,
				Diagnostics: slices.Clone(resp.Diagnostics),
				Identity:    resp.Identity,
				Private:     resp.Private,
				Deferred:    nil,
			}

			r.Resource.Read(ctx, rreq, &rresp)
			if !rresp.Diagnostics.HasError() {

			}
			if rr.DeletePollTerminalDiags(rresp.Diagnostics) {
				return nil
			}
			if rr.DeletePollRetryableDiags(rresp.Diagnostics) {
				return retry.RetryableError(errors.New("retry"))
			}
			err := errors.New("no error received but expects one")
			if diags.HasError() {
				err = errorutil.DiagsToError(diags)
			}
			return retry.NonRetryableError(err)
		}); err != nil {
			// Appending the diagnostics to the delete response
			resp.Diagnostics.AddError("Post delete poll", err.Error())
			return
		}
	}
}

func (r resourceWrapper) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.Resource.SetMeta(req.ProviderData.(meta.Meta))
}

// ConfigValidators implements resource.ResourceWithConfigValidators.
func (r resourceWrapper) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	if rr, ok := r.Resource.(resource.ResourceWithConfigValidators); ok {
		return rr.ConfigValidators(ctx)
	}
	return nil
}

// ModifyPlan implements resource.ResourceWithModifyPlan.
func (r resourceWrapper) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if rr, ok := r.Resource.(resource.ResourceWithModifyPlan); ok {
		rr.ModifyPlan(ctx, req, resp)
		return
	}
}

// MoveState implements resource.ResourceWithMoveState.
func (r resourceWrapper) MoveState(ctx context.Context) []resource.StateMover {
	if rr, ok := r.Resource.(resource.ResourceWithMoveState); ok {
		return rr.MoveState(ctx)
	}
	return nil
}

// UpgradeState implements resource.ResourceWithUpgradeState.
func (r resourceWrapper) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	if rr, ok := r.Resource.(resource.ResourceWithUpgradeState); ok {
		return rr.UpgradeState(ctx)
	}
	return nil
}

// ValidateConfig implements resource.ResourceWithValidateConfig.
func (r resourceWrapper) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if rr, ok := r.Resource.(resource.ResourceWithValidateConfig); ok {
		rr.ValidateConfig(ctx, req, resp)
		return
	}
}

// UpgradeIdentity implements resource.ResourceWithUpgradeIdentity.
func (r resourceWrapper) UpgradeIdentity(ctx context.Context) map[int64]resource.IdentityUpgrader {
	if rr, ok := r.Resource.(resource.ResourceWithUpgradeIdentity); ok {
		return rr.UpgradeIdentity(ctx)
	}
	return nil
}

func (r resourceWrapper) setIdentity(ctx context.Context, state tfsdk.State, identity *tfsdk.ResourceIdentity) (diags diag.Diagnostics) {
	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Set the resource identity")
	ident := r.Identity()
	for _, field := range ident.Fields() {
		v := field.Value
		diags.Append(state.GetAttribute(ctx, field.PathState, &v)...)
		if diags.HasError() {
			return
		}
		diags.Append(identity.SetAttribute(ctx, field.PathIdentity, v)...)
		if diags.HasError() {
			return
		}
	}
	return diags
}
