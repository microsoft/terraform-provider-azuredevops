package framework

import (
	"context"
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
		identity.FromId(req.ID)
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
	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Finish to create the resource")

	// Early return, otherwise if we set the state with error diagnostics, the resource will be in tainted state.
	if resp.Diagnostics.HasError() {
		return
	}

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

	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Start to read the resource (post-creation)")
	r.Resource.Read(ctx, rreq, &rresp)
	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Finish to read the resource (post-creation)")

	// Set the identity
	resp.Diagnostics = append(resp.Diagnostics, r.setIdentity(ctx, rresp.State, rresp.Identity)...)

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

	ctx = tflog.NewSubsystem(ctx, r.Resource.ResourceType())
	ctx = tflog.SubsystemSetField(ctx, r.Resource.ResourceType(), "operation", "Read")

	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Start to read the resource")
	r.Resource.Read(ctx, req, resp)
	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Finish to read the resource")

	// If the resource doesn't exist, remove it from the state and return.
	if slices.ContainsFunc(resp.Diagnostics, IsDiagResourceNotFound) {
		tflog.SubsystemWarn(ctx, r.Resource.ResourceType(), "Resource not found, removing it from the state and return")
		resp.Diagnostics = slices.DeleteFunc(resp.Diagnostics, IsDiagResourceNotFound)
		resp.State.RemoveResource(ctx)
		return
	}

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
	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Finish to update the resource")

	// Early return, otherwise if we set the state with error diagnostics, the resource will be in tainted state.
	if resp.Diagnostics.HasError() {
		return
	}

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

	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Start to read the resource (post-update)")
	r.Resource.Read(ctx, rreq, &rresp)
	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Finish to read the resource (post-update)")

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

	ctx = tflog.NewSubsystem(ctx, r.Resource.ResourceType())
	ctx = tflog.SubsystemSetField(ctx, r.Resource.ResourceType(), "operation", "Delete")

	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Start to delete the resource")
	r.Resource.Delete(ctx, req, resp)
	tflog.SubsystemInfo(ctx, r.Resource.ResourceType(), "Finish to delete the resource")
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
