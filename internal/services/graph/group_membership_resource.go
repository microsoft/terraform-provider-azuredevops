package graph

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/terraform-provider-azuredevops/internal/framework"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/errorutil"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/retry"
)

var _ framework.ResourceWithCreatePoll = &groupMembershipResource{}
var _ framework.ResourceWithDeletePoll = &groupMembershipResource{}

func NewGroupMembershipResource() framework.Resource {
	return &groupMembershipResource{}
}

type groupMembershipResource struct {
	framework.ImplSetMeta
	framework.ImplMetadata
	framework.ImplLog[*groupMembershipResource]
}

type groupMembershipIdentityModel struct {
	GroupId  types.String `tfsdk:"group_id"`
	MemberId types.String `tfsdk:"member_id"`
}

func (p *groupMembershipIdentityModel) Fields() []framework.IdentityField {
	return []framework.IdentityField{
		{
			PathState:    path.Root("group_id"),
			PathIdentity: path.Root("group_id"),
			Value:        p.GroupId,
		},
		{
			PathState:    path.Root("member_id"),
			PathIdentity: path.Root("member_id"),
			Value:        p.MemberId,
		},
	}
}

func (p *groupMembershipIdentityModel) FromId(id string) error {
	groupId, memberId, ok := strings.Cut(id, "/")
	if !ok {
		return fmt.Errorf(`invalid id format, expect="<group_id>/<member_id>"`)
	}
	p.GroupId = types.StringValue(groupId)
	p.MemberId = types.StringValue(memberId)
	return nil
}

func (r *groupMembershipResource) ResourceType() string {
	return "azuredevops_group_membership"
}

func (r *groupMembershipResource) Identity() framework.ResourceIdentity {
	return &groupMembershipIdentityModel{}
}

func (r *groupMembershipResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"group_id": identityschema.StringAttribute{
				RequiredForImport: true,
				Description:       "The id of the container group",
			},
			"member_id": identityschema.StringAttribute{
				RequiredForImport: true,
				Description:       "The id of the member user/group",
			},
		},
	}
}

type groupMembershipModel struct {
	GroupId  types.String   `tfsdk:"group_id"`
	MemberId types.String   `tfsdk:"member_id"`
	Id       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

func (r *groupMembershipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				MarkdownDescription: "The id of the container group.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"member_id": schema.StringAttribute{
				MarkdownDescription: "The id of the member group/user.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the group membership.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *groupMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupMembershipModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.Info(ctx, "check group membership existence")

	if err := r.GraphClient.CheckMembershipExistence(ctx, graph.CheckMembershipExistenceArgs{
		SubjectDescriptor:   plan.MemberId.ValueStringPointer(),
		ContainerDescriptor: plan.GroupId.ValueStringPointer(),
	}); err == nil {
		resp.Diagnostics.Append(
			errorutil.ImportAsExistsError(
				r.ResourceType(),
				fmt.Sprintf("%s/%s", plan.GroupId.ValueString(), plan.MemberId.ValueString()),
			),
		)
		return
	} else if !errorutil.WasNotFound(err) {
		resp.Diagnostics = append(resp.Diagnostics, framework.NewDiagSdkError("Existence check", err))
		return
	}

	r.Info(ctx, "create the group membership")

	if _, err := r.GraphClient.AddMembership(ctx, graph.AddMembershipArgs{
		SubjectDescriptor:   plan.MemberId.ValueStringPointer(),
		ContainerDescriptor: plan.GroupId.ValueStringPointer(),
	}); err != nil {
		resp.Diagnostics = append(resp.Diagnostics, framework.NewDiagSdkError("Create the group membership", err))
		return
	}

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupMembershipModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.Info(ctx, "get the group membership")

	if _, err := r.Meta.GraphClient.GetMembership(ctx, graph.GetMembershipArgs{
		SubjectDescriptor:   state.MemberId.ValueStringPointer(),
		ContainerDescriptor: state.GroupId.ValueStringPointer(),
	}); err != nil {
		resp.Diagnostics = append(resp.Diagnostics, framework.NewDiagSdkError("Get the group membership", err))
		return
	}

	state.Id = types.StringValue(fmt.Sprintf("%s/%s", state.GroupId.ValueString(), state.MemberId.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *groupMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.Append(errorutil.NoUpdateError())
}

func (r *groupMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state groupMembershipModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.Info(ctx, "delete the group membership")

	if err := r.Meta.GraphClient.RemoveMembership(ctx, graph.RemoveMembershipArgs{
		SubjectDescriptor:   state.MemberId.ValueStringPointer(),
		ContainerDescriptor: state.GroupId.ValueStringPointer(),
	}); err != nil {
		resp.Diagnostics = append(resp.Diagnostics, framework.NewDiagSdkError("Delete the group membership", err))
		return
	}
}

func (r *groupMembershipResource) CreatePollOption(ctx context.Context) retry.RetryOption {
	return retry.NewSimpleRetryOption(ctx, 10, time.Second)
}

func (r *groupMembershipResource) CreatePollCheckers() []framework.PollChecker {
	return nil
}

func (r *groupMembershipResource) CreatePollRetryableDiags(diags diag.Diagnostics) bool {
	return slices.ContainsFunc(diags, framework.IsDiagResourceNotFound)
}

func (r *groupMembershipResource) DeletePollOption(ctx context.Context) retry.RetryOption {
	return retry.NewSimpleRetryOption(ctx, 10, time.Second)
}

func (r *groupMembershipResource) DeletePollRetryableDiags(diags diag.Diagnostics) bool {
	return !diags.HasError()
}

func (r *groupMembershipResource) DeletePollTerminalDiags(diags diag.Diagnostics) bool {
	return slices.ContainsFunc(diags, framework.IsDiagResourceNotFound)
}
