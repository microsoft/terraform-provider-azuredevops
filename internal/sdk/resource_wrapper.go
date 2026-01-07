package sdk

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-azuredevops/internal/meta"
)

// Implement "Safe" resource interfaces here.
// "Safe" here means calling an interface by doing nothing is effectively the same as if
// this interface is not implemented, during terraform execution.
var _ resource.Resource = resourceWrapper{}
var _ resource.ResourceWithConfigure = resourceWrapper{}
var _ resource.ResourceWithConfigValidators = resourceWrapper{}
var _ resource.ResourceWithImportState = resourceWrapper{}
var _ resource.ResourceWithModifyPlan = resourceWrapper{}
var _ resource.ResourceWithMoveState = resourceWrapper{}
var _ resource.ResourceWithUpgradeState = resourceWrapper{}
var _ resource.ResourceWithValidateConfig = resourceWrapper{}

// The followings are unsafe interfaces. This requires additional wrappers around this resourceWrapper and opt-in.
// var _ resource.ResourceWithIdentity = resourceWrapper{}
// var _ resource.ResourceWithUpgradeIdentity = resourceWrapper{}

type resourceWrapper struct {
	inner Resource
}

func WrapResource(in Resource) func() resource.Resource {
	return func() resource.Resource {
		return resourceWrapper{inner: in}
	}
}

func (r resourceWrapper) Timeout() ResourceTimeout {
	if r, ok := r.inner.(ResourceWithTimeout); ok {
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
	resp.TypeName = r.inner.ResourceType()
}

func (r resourceWrapper) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.inner.Schema(ctx, req, resp)
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

	ctx = tflog.NewSubsystem(ctx, r.inner.ResourceType())
	ctx = tflog.SubsystemSetField(ctx, r.inner.ResourceType(), "operation", "Create")

	r.inner.Create(ctx, req, resp)

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

	r.inner.Read(ctx, rreq, &rresp)

	*resp = resource.CreateResponse{
		State:       rresp.State,
		Identity:    rresp.Identity,
		Private:     rresp.Private,
		Diagnostics: rresp.Diagnostics,
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

	ctx = tflog.NewSubsystem(ctx, r.inner.ResourceType())
	ctx = tflog.SubsystemSetField(ctx, r.inner.ResourceType(), "operation", "Read")

	r.inner.Read(ctx, req, resp)

	// If the resource doesn't exist, remove it from the state and return.
	if slices.ContainsFunc(resp.Diagnostics, IsDiagResourceNotFound) {
		resp.Diagnostics = slices.DeleteFunc(resp.Diagnostics, IsDiagResourceNotFound)
		resp.State.RemoveResource(ctx)
		return
	}
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

	ctx = tflog.NewSubsystem(ctx, r.inner.ResourceType())
	ctx = tflog.SubsystemSetField(ctx, r.inner.ResourceType(), "operation", "Update")

	r.inner.Update(ctx, req, resp)

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

	r.inner.Read(ctx, rreq, &rresp)

	*resp = resource.UpdateResponse{
		State:       rresp.State,
		Identity:    rresp.Identity,
		Private:     rresp.Private,
		Diagnostics: rresp.Diagnostics,
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

	ctx = tflog.NewSubsystem(ctx, r.inner.ResourceType())
	ctx = tflog.SubsystemSetField(ctx, r.inner.ResourceType(), "operation", "Delete")

	r.inner.Delete(ctx, req, resp)
}

// Configure implements resource.ResourceWithConfigure.
func (r resourceWrapper) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.inner.SetMeta(req.ProviderData.(meta.Meta))
}

// ConfigValidators implements resource.ResourceWithConfigValidators.
func (r resourceWrapper) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	if r, ok := r.inner.(resource.ResourceWithConfigValidators); ok {
		return r.ConfigValidators(ctx)
	}
	return nil
}

// ImportState implements resource.ResourceWithImportState.
func (r resourceWrapper) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r, ok := r.inner.(resource.ResourceWithImportState); ok {
		r.ImportState(ctx, req, resp)
		return
	}
}

// ModifyPlan implements resource.ResourceWithModifyPlan.
func (r resourceWrapper) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if r, ok := r.inner.(resource.ResourceWithModifyPlan); ok {
		r.ModifyPlan(ctx, req, resp)
		return
	}
}

// MoveState implements resource.ResourceWithMoveState.
func (r resourceWrapper) MoveState(ctx context.Context) []resource.StateMover {
	if r, ok := r.inner.(resource.ResourceWithMoveState); ok {
		return r.MoveState(ctx)
	}
	return nil
}

// UpgradeState implements resource.ResourceWithUpgradeState.
func (r resourceWrapper) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	if r, ok := r.inner.(resource.ResourceWithUpgradeState); ok {
		return r.UpgradeState(ctx)
	}
	return nil
}

// ValidateConfig implements resource.ResourceWithValidateConfig.
func (r resourceWrapper) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if r, ok := r.inner.(resource.ResourceWithValidateConfig); ok {
		r.ValidateConfig(ctx, req, resp)
		return
	}
}
