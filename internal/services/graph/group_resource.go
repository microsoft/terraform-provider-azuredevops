package graph

import (
	"context"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
	"github.com/microsoft/terraform-provider-azuredevops/internal/adocustomtype"
	"github.com/microsoft/terraform-provider-azuredevops/internal/adovalidator"
	"github.com/microsoft/terraform-provider-azuredevops/internal/framework"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/errorutil"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/fwtype"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/pointer"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/retry"
)

var _ framework.ResourceWithPostUpdatePoll = &groupResource{}

func NewGroupResource() framework.Resource {
	return &groupResource{}
}

type groupResource struct {
	framework.ImplSetMeta
	framework.ImplMetadata
	framework.ImplLog[*groupResource]
}

type groupIdentityModel struct {
	Id types.String `tfsdk:"id"`
}

func (p *groupIdentityModel) Fields() []framework.IdentityField {
	return []framework.IdentityField{
		{
			PathState:    path.Root("id"),
			PathIdentity: path.Root("id"),
			Value:        p.Id,
		},
	}
}

func (p *groupIdentityModel) FromId(id string) error {
	p.Id = types.StringValue(id)
	return nil
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
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
				Description:       "The descriptor of the group",
			},
		},
	}
}

type groupModel struct {
	OriginId      types.String                  `tfsdk:"origin_id"`
	Mail          types.String                  `tfsdk:"mail"`
	DisplayName   types.String                  `tfsdk:"display_name"`
	Description   types.String                  `tfsdk:"description"`
	Scope         adocustomtype.StringUUIDValue `tfsdk:"scope"`
	Url           types.String                  `tfsdk:"url"`
	Origin        types.String                  `tfsdk:"origin"`
	SubjectKind   types.String                  `tfsdk:"subject_kind"`
	Domain        types.String                  `tfsdk:"domain"`
	PrincipalName types.String                  `tfsdk:"principal_name"`
	StorageKey    types.String                  `tfsdk:"storage_key"`
	Id            types.String                  `tfsdk:"id"`
	Timeouts      timeouts.Value                `tfsdk:"timeouts"`
}

func (r *groupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"origin_id": schema.StringAttribute{
				MarkdownDescription: "This will create a new graph group that is derived from the object id of an AAD group.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					adovalidator.StringIsUUID(),
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("display_name"),
						path.MatchRoot("mail"),
						path.MatchRoot("origin_id"),
					),
				},
			},
			"mail": schema.StringAttribute{
				MarkdownDescription: "This will create a new graph group that is derived from the email of an AAD group.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("display_name"),
						path.MatchRoot("mail"),
						path.MatchRoot("origin_id"),
					),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the VSTS group.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("display_name"),
						path.MatchRoot("mail"),
						path.MatchRoot("origin_id"),
					),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the group.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("origin_id"),
						path.MatchRoot("mail"),
					),
				},
			},
			"scope": schema.StringAttribute{
				CustomType:          adocustomtype.StringUUIDType{},
				MarkdownDescription: "The id of the scope (e.g. project) in which the group should be created. If omitted, will be created in the scope of the enclosing account or organization. Valid only for VSTS groups.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					adovalidator.StringIsUUID(),
					stringvalidator.ConflictsWith(
						path.MatchRoot("origin_id"),
						path.MatchRoot("mail"),
					),
				},
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

			"storage_key": schema.StringAttribute{
				MarkdownDescription: "The storage key of the group's descriptor.",
				Computed:            true,
			},

			"id": schema.StringAttribute{
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

	var scopeDescriptor *string
	if scope := plan.Scope.ValueString(); scope != "" {
		r.Info(ctx, "get the descriptor for scope", map[string]any{"scope": scope})
		desc, err := r.Meta.GraphClient.GetDescriptor(ctx, graph.GetDescriptorArgs{
			StorageKey: &plan.Scope.UUID,
		})
		if err != nil {
			resp.Diagnostics = append(resp.Diagnostics, framework.NewDiagSdkError("Get descriptor of the scope", err))
			return
		}
		scopeDescriptor = desc.Value
	}

	var (
		group *graph.GraphGroup
		err   error
	)

	r.Info(ctx, "create the group")
	switch {
	case !plan.DisplayName.IsNull() && !plan.DisplayName.IsUnknown():
		{
			// Creating a VSTS group.
			param := graph.CreateGroupVstsArgs{
				ScopeDescriptor: scopeDescriptor,
				CreationContext: &graph.GraphGroupVstsCreationContext{
					DisplayName: plan.DisplayName.ValueStringPointer(),
					Description: plan.Description.ValueStringPointer(),
				},
			}
			group, err = r.Meta.GraphClient.CreateGroupVsts(ctx, param)
		}
	case !plan.OriginId.IsNull() && !plan.OriginId.IsUnknown():
		{
			// Creating a group derived from an AAD group by object id.
			param := graph.CreateGroupOriginIdArgs{
				ScopeDescriptor: scopeDescriptor,
				CreationContext: &graph.GraphGroupOriginIdCreationContext{
					OriginId: plan.OriginId.ValueStringPointer(),
				},
			}
			group, err = r.Meta.GraphClient.CreateGroupOriginId(ctx, param)
		}
	case !plan.Mail.IsNull() && !plan.Mail.IsUnknown():
		{
			// Creating a group derived from an AAD group by mail address.
			param := graph.CreateGroupMailAddressArgs{
				ScopeDescriptor: scopeDescriptor,
				CreationContext: &graph.GraphGroupMailAddressCreationContext{
					MailAddress: plan.Mail.ValueStringPointer(),
				},
			}
			group, err = r.Meta.GraphClient.CreateGroupMailAddress(ctx, param)
		}
	}
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, framework.NewDiagSdkError("Creating the group", err))
		return
	}

	// Set id related attributes to the state to be used by the read.
	plan.Id = fwtype.StringValue(group.Descriptor)
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

	group, err := r.Meta.GraphClient.GetGroup(ctx, graph.GetGroupArgs{GroupDescriptor: state.Id.ValueStringPointer()})
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, framework.NewDiagSdkError("Get the group", err))
		return
	}

	if pointer.To(group.IsDeleted) {
		resp.Diagnostics = append(resp.Diagnostics,
			framework.NewDiagSdkErrorWithCode("Group is explicitly deleted", http.StatusNotFound))
		return
	}

	// Set state
	state.OriginId = fwtype.StringValue(group.OriginId)
	state.Mail = fwtype.StringValue(group.MailAddress)
	state.DisplayName = fwtype.StringValue(group.DisplayName)
	state.Description = fwtype.StringValue(group.Description)
	state.Url = fwtype.StringValue(group.Url)
	state.Origin = fwtype.StringValue(group.Origin)
	state.SubjectKind = fwtype.StringValue(group.SubjectKind)
	state.Domain = fwtype.StringValue(group.Domain)
	state.PrincipalName = fwtype.StringValue(group.PrincipalName)

	state.Scope = adocustomtype.StringUUIDValue{StringValue: types.StringNull()}
	if domain := group.Domain; domain != nil {
		// The domain can be:
		// - Organization scope: vstfs:///Framework/IdentityDomain/<uuid>
		// - Project scope: vstfs:///Classification/TeamProject/<uuid>
		// - Other unknown cases.
		//
		// We simply cut the last segment and regard it as the scope if is an uuid.
		l := strings.Split(*domain, "/")
		scope, diags := adocustomtype.StringUUIDType{}.ValueFromString(ctx, types.StringValue(l[len(l)-1]))
		if !diags.HasError() {
			state.Scope = scope.(adocustomtype.StringUUIDValue)
		}
	}

	storageKey, err := r.GraphClient.GetStorageKey(ctx, graph.GetStorageKeyArgs{
		SubjectDescriptor: group.Descriptor,
	})
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, framework.NewDiagSdkError("Read the storage key", err))
		return
	}
	if id := storageKey.Value; id != nil {
		state.StorageKey = types.StringValue(id.String())
	}

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
			GroupDescriptor: plan.Id.ValueStringPointer(),
			PatchDocument:   &operations,
		}); err != nil {
			resp.Diagnostics = append(resp.Diagnostics, framework.NewDiagSdkError("Update the group", err))
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
		GroupDescriptor: state.Id.ValueStringPointer(),
	}); err != nil {
		resp.Diagnostics = append(resp.Diagnostics, framework.NewDiagSdkError("Delete the group", err))
		return
	}
}

func (r *groupResource) PostUpdatePollRetryOption(ctx context.Context) retry.RetryOption {
	return retry.NewSimpleRetryOption(ctx, 5)
}

func (r *groupResource) PostUpdateCheck(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) bool {
	var stateModel groupModel
	if err := errorutil.DiagsToError(state.Get(ctx, &stateModel)); err != nil {
		return false
	}

	var planModel groupModel
	if err := errorutil.DiagsToError(plan.Get(ctx, &planModel)); err != nil {
		return false
	}

	return planModel.DisplayName.Equal(stateModel.DisplayName) && planModel.Description.Equal(stateModel.Description)
}

func (r *groupResource) PostUpdateRetryableDiag(diag.Diagnostic) bool {
	return false
}
