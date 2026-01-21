package graph

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
	"github.com/microsoft/terraform-provider-azuredevops/internal/errorutil"
	"github.com/microsoft/terraform-provider-azuredevops/internal/framework"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/fwtype"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/pointer"
)

var _ framework.Resource = &groupResource{}

func NewGroupResource() framework.Resource {
	return &groupResource{}
}

type groupResource struct {
	framework.ImplSetMeta
	framework.ImplMetadata
	framework.ImplLog[*groupResource]
}

type groupIdentityModel struct {
	Descriptor types.String `tfsdk:"descriptor"`
}

func (p *groupIdentityModel) Fields() []framework.IdentityField {
	return []framework.IdentityField{
		{
			PathState:    path.Root("descriptor"),
			PathIdentity: path.Root("descriptor"),
			Value:        p.Descriptor,
		},
	}
}

func (p *groupIdentityModel) FromId(id string) {
	p.Descriptor = types.StringValue(id)
}

func (r *groupResource) ResourceType() string {
	return "azuredevops_group"
}

func (r *groupResource) Identity() framework.ResourceIdentity {
	return &groupIdentityModel{}
}

func (r *groupResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"descriptor": identityschema.StringAttribute{
				RequiredForImport: true,
				Description:       "The descriptor of the group",
			},
		},
	}
}

type groupModel struct {
	DisplayName   types.String   `tfsdk:"display_name"`
	Description   types.String   `tfsdk:"description"`
	Url           types.String   `tfsdk:"url"`
	Origin        types.String   `tfsdk:"origin"`
	SubjectKind   types.String   `tfsdk:"subject_kind"`
	Domain        types.String   `tfsdk:"domain"`
	PrincipalName types.String   `tfsdk:"principal_name"`
	Descriptor    types.String   `tfsdk:"descriptor"`
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
}

func (r *groupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the group.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the group.",
				Optional:            true,
			},

			"url": schema.StringAttribute{
				MarkdownDescription: "The full route to the group.",
				Computed:            true,
			},

			"origin": schema.StringAttribute{
				MarkdownDescription: "The type of source provider for the group.",
				Computed:            true,
			},

			"subject_kind": schema.StringAttribute{
				MarkdownDescription: "The type of the graph subject for the group.",
				Computed:            true,
			},

			"domain": schema.StringAttribute{
				MarkdownDescription: "The name of the container fof origin for the group.",
				Computed:            true,
			},

			"principal_name": schema.StringAttribute{
				MarkdownDescription: "The principal name of this group from the source provider.",
				Computed:            true,
			},

			"descriptor": schema.StringAttribute{
				MarkdownDescription: "The descriptor of the group",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.Info(ctx, "create the group")

	param := graph.CreateGroupVstsArgs{
		CreationContext: &graph.GraphGroupVstsCreationContext{
			DisplayName: plan.DisplayName.ValueStringPointer(),
			Description: plan.Description.ValueStringPointer(),
		},
	}
	group, err := r.Meta.GraphClient.CreateGroupVsts(ctx, param)
	if err != nil {
		resp.Diagnostics.AddError("Creating the group", err.Error())
		return
	}

	// Set id related attributes to the state to be used by the read.
	plan.Descriptor = fwtype.StringValue(group.Descriptor)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.Info(ctx, "get the group")

	group, err := r.Meta.GraphClient.GetGroup(ctx, graph.GetGroupArgs{GroupDescriptor: state.Descriptor.ValueStringPointer()})
	if err != nil {
		if errorutil.WasNotFound(err) {
			resp.Diagnostics = append(resp.Diagnostics, framework.NewDiagResourceNotFound(r.ResourceType(), state.Descriptor.ValueString()))
			return
		}
		resp.Diagnostics.AddError("Get the group", err.Error())
		return
	}

	if pointer.To(group.IsDeleted) {
		resp.Diagnostics = append(resp.Diagnostics, framework.NewDiagResourceNotFound(r.ResourceType(), state.DisplayName.ValueString()))
		return
	}

	// Set state
	state.DisplayName = fwtype.StringValue(group.DisplayName)
	state.Description = fwtype.StringValue(group.Description)
	state.Url = fwtype.StringValue(group.Url)
	state.Origin = fwtype.StringValue(group.Origin)
	state.SubjectKind = fwtype.StringValue(group.SubjectKind)
	state.Domain = fwtype.StringValue(group.Domain)
	state.PrincipalName = fwtype.StringValue(group.PrincipalName)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan groupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state groupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var operations []webapi.JsonPatchOperation

	if !plan.DisplayName.Equal(state.DisplayName) {
		operations = append(operations, webapi.JsonPatchOperation{
			Op:    &webapi.OperationValues.Replace,
			Path:  pointer.From("/displayName"),
			Value: plan.DisplayName.ValueString(),
		})
	}
	if !plan.Description.Equal(state.Description) {
		operations = append(operations, webapi.JsonPatchOperation{
			Op:    &webapi.OperationValues.Replace,
			Path:  pointer.From("/description"),
			Value: plan.Description.ValueString(),
		})
	}

	if len(operations) > 0 {
		r.Info(ctx, "update the group")
		if _, err := r.Meta.GraphClient.UpdateGroup(ctx, graph.UpdateGroupArgs{
			GroupDescriptor: plan.Descriptor.ValueStringPointer(),
			PatchDocument:   &operations,
		}); err != nil {
			resp.Diagnostics.AddError("Update the group", err.Error())
			return
		}
	} else {
		r.Info(ctx, "nothing to update onthe group")
	}
}

func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state groupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.Info(ctx, "delete the group")

	if err := r.Meta.GraphClient.DeleteGroup(ctx, graph.DeleteGroupArgs{
		GroupDescriptor: state.Descriptor.ValueStringPointer(),
	}); err != nil {
		resp.Diagnostics.AddError("Delete the group", err.Error())
		return
	}
}
