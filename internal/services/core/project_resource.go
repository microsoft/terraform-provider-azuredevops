package core

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/operations"
	"github.com/microsoft/terraform-provider-azuredevops/internal/adocustomtype"
	"github.com/microsoft/terraform-provider-azuredevops/internal/framework"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/fwtype"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/pointer"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/retry"
)

var _ framework.Resource = &projectResource{}
var _ framework.ResourceWithPostCreate = &projectResource{}
var _ framework.ResourceWithPostCreatePoll = &projectResource{}
var _ framework.ResourceWithPostUpdate = &projectResource{}
var _ framework.ResourceWithPostUpdatePoll = &projectResource{}

func NewProjectResource() framework.Resource {
	return &projectResource{}
}

type projectResource struct {
	framework.ImplSetMeta
	framework.ImplResourceMetadata
	framework.ImplLog[*projectResource]
}

type projectIdentityModel struct {
	ID types.String `tfsdk:"id"`
}

func (p *projectIdentityModel) Fields() []framework.IdentityField {
	return []framework.IdentityField{
		{
			PathState:    path.Root("id"),
			PathIdentity: path.Root("id"),
			Value:        p.ID,
		},
	}
}

func (p *projectIdentityModel) FromId(id string) error {
	p.ID = types.StringValue(id)
	return nil
}

type projectResourceModel struct {
	Name              adocustomtype.StringCaseInsensitiveValue `tfsdk:"name"`
	Description       types.String                             `tfsdk:"description"`
	Visibility        types.String                             `tfsdk:"visibility"`
	VersionControl    types.String                             `tfsdk:"version_control"`
	WorkItemTemplate  types.String                             `tfsdk:"work_item_template"`
	Features          types.Object                             `tfsdk:"features"`
	Id                types.String                             `tfsdk:"id"`
	ProcessTemplateId types.String                             `tfsdk:"process_template_id"`
	Timeouts          timeouts.Value                           `tfsdk:"timeouts"`
}

func (*projectResource) ResourceType() string {
	return "azuredevops_project"
}

func (r *projectResource) Identity() framework.ResourceIdentity {
	return &projectIdentityModel{}
}

func (r *projectResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
				Description:       "The project UUID",
			},
		},
	}
}

func (r *projectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				CustomType:          adocustomtype.StringCaseInsensitiveType{},
				MarkdownDescription: "The name of the project",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the project",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					// An empty description in request will not be returned in the response.
					stringvalidator.LengthAtLeast(1),
				},
				Default: stringdefault.StaticString(""),
			},
			"visibility": schema.StringAttribute{
				MarkdownDescription: "The visibility of the project. Possible values are `private` and `public`. Defaults to `private`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(string(core.ProjectVisibilityValues.Private)),
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(core.ProjectVisibilityValues.Private),
						string(core.ProjectVisibilityValues.Public),
					),
				},
			},
			"version_control": schema.StringAttribute{
				MarkdownDescription: "The version control system. Possbile values are: `Git` or `Tfvc`. Defaults to `Git`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(string(core.SourceControlTypesValues.Git)),
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(core.SourceControlTypesValues.Git),
						string(core.SourceControlTypesValues.Tfvc),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"work_item_template": schema.StringAttribute{
				MarkdownDescription: "The work item template name. Defaults to the parent organization's default template.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"features": schema.SingleNestedAttribute{
				MarkdownDescription: "Define the status (enabled/disabled) of the project features.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"boards": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
					"repos": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
					"pipelines": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
					"test_plans": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
					"artifacts": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the project",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"process_template_id": schema.StringAttribute{
				MarkdownDescription: "The Process Template ID used by the Project.",
				Computed:            true,
			},
		},
	}
}

func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan projectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.Info(ctx, "look up the process template id")

	process, err := LookupProcess(ctx, r.Meta.CoreClient, func(p core.Process) bool {
		// `work_item_template` is unknown as it is set as O+C.
		if plan.WorkItemTemplate.IsUnknown() {
			return p.IsDefault != nil && *p.IsDefault
		} else {
			return p.Name != nil && *p.Name == plan.WorkItemTemplate.ValueString()
		}
	})
	if err != nil {
		resp.Diagnostics.AddError("Lookup process", err.Error())
		return
	}
	if process.Id == nil {
		resp.Diagnostics.AddError("Lookup process", "unexpected null id")
		return
	}
	processTemplateId := process.Id.String()

	r.Info(ctx, "create the project")

	project := &core.TeamProject{
		Name:        plan.Name.ValueStringPointer(),
		Description: plan.Description.ValueStringPointer(),
		Visibility:  pointer.From(core.ProjectVisibility(plan.Visibility.ValueString())),
		Capabilities: &core.TeamProjectCapabilities{
			Versioncontrol: &core.TeamProjectCapabilitiesVersionControl{
				SourceControlType: (*core.SourceControlTypes)(plan.VersionControl.ValueStringPointer()),
			},
			ProcessTemplate: &core.TeamProjectCapabilitiesProcessTemplate{
				TemplateId: &processTemplateId,
			},
		},
	}

	operationRef, err := r.Meta.CoreClient.QueueCreateProject(ctx, core.QueueCreateProjectArgs{ProjectToCreate: project})
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, framework.NewDiagSdkError("Queue create project", err))
		return
	}

	r.Info(ctx, "wait for the project creation")

	if err := r.waitOperation(ctx, operationRef); err != nil {
		resp.Diagnostics.AddError("Wait for project delete operation", err.Error())
		return
	}
}

func (r *projectResource) PostCreate(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state projectResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var features projectFeaturesTFModel
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("features"), &features)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.Info(ctx, "set the project features")
	if err := setProjectFeature(ctx, r.Meta.FeatureManagementClient, state.Id.ValueString(), features); err != nil {
		resp.Diagnostics.AddError("failed to set project featuers", err.Error())
	}
}

func (r *projectResource) ShouldPostCreate(ctx context.Context, req resource.CreateRequest) bool {
	var features types.Object
	// If not specified (will be unknown as is O+C), no need to set features.
	return !req.Plan.GetAttribute(ctx, path.Root("features"), &features).HasError() && !features.IsUnknown()
}

func (r *projectResource) PostCreatePollCheckers() []framework.PollChecker {
	return r.postWritePollCheckers()
}

func (r *projectResource) PostCreatePollOption(ctx context.Context) retry.RetryOption {
	return r.postWritePollOption(ctx)
}

func (r *projectResource) PostCreatePollRetryableDiags(diags diag.Diagnostics) bool {
	return false
}

func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state projectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.Info(ctx, "get the project")

	projectId := state.Id.ValueString()
	// The id is not available during creation, hence read by name.
	if projectId == "" {
		projectId = state.Name.ValueString()
	}

	project, err := r.Meta.CoreClient.GetProject(ctx, core.GetProjectArgs{
		ProjectId:           &projectId,
		IncludeCapabilities: pointer.From(true),
	})
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, framework.NewDiagSdkError("Read the project", err))
		return
	}
	if project.Id == nil {
		resp.Diagnostics.AddError("Unexpected API response", "project id is null")
		return
	}

	id := project.Id.String()

	var (
		templateId     *string
		templateName   *string
		versionControl *string
	)
	if cap := project.Capabilities; cap != nil {
		if tpl := cap.ProcessTemplate; tpl != nil {
			templateId = cap.ProcessTemplate.TemplateId
		}
		if vc := cap.Versioncontrol; vc != nil {
			versionControl = (*string)(vc.SourceControlType)
		}
	}
	if templateId != nil {
		r.Info(ctx, "look up the template name by id")
		process, err := LookupProcess(ctx, r.Meta.CoreClient, func(p core.Process) bool { return p.Id != nil && p.Id.String() == *templateId })
		if err != nil {
			resp.Diagnostics.AddError("Lookup process", err.Error())
			return
		}
		templateName = process.Name
	}

	r.Info(ctx, "get the project features")
	features, err := getProjectFeatures(ctx, r.Meta.FeatureManagementClient, id)
	if err != nil {
		resp.Diagnostics.AddError("get the project features", err.Error())
		return
	}

	// Set state
	state.Id = types.StringValue(id)
	state.Name = adocustomtype.StringCaseInsensitiveValue{
		StringValue: fwtype.StringValue(project.Name),
	}
	state.Description = fwtype.StringValueOrZero(project.Description)
	state.Visibility = fwtype.StringValue(project.Visibility)
	state.VersionControl = fwtype.StringValue(versionControl)
	state.WorkItemTemplate = fwtype.StringValue(templateName)
	state.Features = *features
	state.ProcessTemplateId = fwtype.StringValue(templateId)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan projectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state projectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Name.Equal(state.Name) || !plan.Description.Equal(state.Description) || !plan.Visibility.Equal(state.Visibility) {
		r.Info(ctx, "update the project")

		id, err := uuid.Parse(state.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Parse resource id as UUID", err.Error())
			return
		}

		project := &core.TeamProject{
			Name:        plan.Name.ValueStringPointer(),
			Description: plan.Description.ValueStringPointer(),
			Visibility:  pointer.From(core.ProjectVisibility(plan.Visibility.ValueString())),
		}

		operationRef, err := r.Meta.CoreClient.UpdateProject(ctx, core.UpdateProjectArgs{
			ProjectUpdate: project,
			ProjectId:     &id,
		})
		if err != nil {
			resp.Diagnostics = append(resp.Diagnostics, framework.NewDiagSdkError("Queue update project", err))
			return
		}

		r.Info(ctx, "wait for the project update")

		if err := r.waitOperation(ctx, operationRef); err != nil {
			resp.Diagnostics.AddError("Wait for project update operation", err.Error())
			return
		}
	}
}

func (r *projectResource) PostUpdate(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state projectResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var features projectFeaturesTFModel
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("features"), &features)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.Info(ctx, "set the project features")
	if err := setProjectFeature(ctx, r.Meta.FeatureManagementClient, state.Id.ValueString(), features); err != nil {
		resp.Diagnostics.AddError("failed to set project featuers", err.Error())
	}
}

func (r *projectResource) ShouldPostUpdate(ctx context.Context, req resource.UpdateRequest) bool {
	var planFeatures types.Object
	if req.Plan.GetAttribute(ctx, path.Root("features"), &planFeatures).HasError() {
		return false
	}
	// If not specified (will be unknown as is O+C), no need to set features.
	if planFeatures.IsUnknown() {
		return false
	}

	var stateFeatures types.Object
	if req.State.GetAttribute(ctx, path.Root("features"), &stateFeatures).HasError() {
		return false
	}

	return !planFeatures.Equal(stateFeatures)
}

func (r *projectResource) PostUpdatePollCheckers() []framework.PollChecker {
	return r.postWritePollCheckers()
}

func (r *projectResource) PostUpdatePollOption(ctx context.Context) retry.RetryOption {
	return r.postWritePollOption(ctx)
}

func (r *projectResource) PostUpdatePollRetryableDiags(diag.Diagnostics) bool {
	return false
}

func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state projectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.Info(ctx, "delete the prject")

	id, err := uuid.Parse(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Parse resource id as UUID", err.Error())
		return
	}

	operationRef, err := r.Meta.CoreClient.QueueDeleteProject(ctx, core.QueueDeleteProjectArgs{ProjectId: &id})
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, framework.NewDiagSdkError("Queue delete project", err))
		return
	}

	r.Info(ctx, "wait for the project deletion")

	if err := r.waitOperation(ctx, operationRef); err != nil {
		resp.Diagnostics.AddError("Wait for project delete operation", err.Error())
		return
	}
}

func (r *projectResource) waitOperation(ctx context.Context, operationRef *operations.OperationReference) error {
	stateConf := &retry.StateChangeConf{
		ContinuousTargetOccurence: 1,
		Delay:                     5 * time.Second,
		MinTimeout:                10 * time.Second,
		Timeout:                   time.Duration(5 * time.Minute),
		Pending: []string{
			string(operations.OperationStatusValues.InProgress),
			string(operations.OperationStatusValues.Queued),
			string(operations.OperationStatusValues.NotSet),
		},
		Target: []string{
			string(operations.OperationStatusValues.Failed),
			string(operations.OperationStatusValues.Succeeded),
			string(operations.OperationStatusValues.Cancelled),
		},
		Refresh: r.pollOperationResult(ctx, operationRef),
	}

	status, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return err
	}
	if status := status.(operations.OperationStatus); status != operations.OperationStatusValues.Succeeded {
		return fmt.Errorf("operation terminated at status: %q", status)
	}
	return nil
}

func (r *projectResource) pollOperationResult(ctx context.Context, operationRef *operations.OperationReference) retry.StateRefreshFunc {
	return func() (any, string, error) {
		ret, err := r.Meta.OperationsClient.GetOperation(ctx, operations.GetOperationArgs{
			OperationId: operationRef.Id,
			PluginId:    operationRef.PluginId,
		})
		if err != nil {
			return nil, string(operations.OperationStatusValues.Failed), err
		}

		r.Info(ctx, "polling project operation status", map[string]any{
			"status": ret.Status,
			"detaul": ret.DetailedMessage,
		})

		status := operations.OperationStatusValues.InProgress
		if ret.Status != nil {
			status = *ret.Status
		}
		return status, string(status), nil
	}
}

func (r *projectResource) postWritePollCheckers() []framework.PollChecker {
	return []framework.PollChecker{
		{
			AttrPath: path.Root("features"),
			Target: types.ObjectNull(map[string]attr.Type{
				"boards":     types.BoolType,
				"repos":      types.BoolType,
				"pipelines":  types.BoolType,
				"test_plans": types.BoolType,
				"artifacts":  types.BoolType,
			}),
		},
	}
}

func (r *projectResource) postWritePollOption(ctx context.Context) retry.RetryOption {
	return retry.NewSimpleRetryOption(ctx, 10, time.Second)
}
