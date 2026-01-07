package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/operations"
	"github.com/microsoft/terraform-provider-azuredevops/internal/adocustomtype"
	"github.com/microsoft/terraform-provider-azuredevops/internal/errorutil"
	"github.com/microsoft/terraform-provider-azuredevops/internal/sdk"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/fwtype"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/pointer"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/retry"
)

var _ sdk.Resource = &projectResource{}

func NewProjectResource() sdk.Resource {
	return &projectResource{}
}

type projectResource struct {
	sdk.ImplSetMeta
	sdk.ImplMetadata
	sdk.ImplLog[*projectResource]
}

type projectModel struct {
	Name              adocustomtype.StringCaseInsensitiveValue `tfsdk:"name"`
	Description       types.String                             `tfsdk:"description"`
	Visibility        types.String                             `tfsdk:"visibility"`
	VersionControl    types.String                             `tfsdk:"version_control"`
	WorkItemTemplate  types.String                             `tfsdk:"work_item_template"`
	Id                types.String                             `tfsdk:"id"`
	ProcessTemplateId types.String                             `tfsdk:"process_template_id"`
	Timeouts          timeouts.Value                           `tfsdk:"timeouts"`
}

func (*projectResource) ResourceType() string {
	return "azuredevops_project"
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
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
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
	var plan projectModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.Log(ctx, "check the presence of the project")

	existing, err := r.Meta.CoreClient.GetProject(ctx, core.GetProjectArgs{ProjectId: plan.Name.ValueStringPointer()})
	if err != nil {
		if !errorutil.WasNotFound(err) {
			resp.Diagnostics.AddError("check the presence of the project", err.Error())
			return
		}
	}
	if !errorutil.WasNotFound(err) {
		if existing.Id == nil {
			resp.Diagnostics.AddError("check the presence of the project", "existing project's id is null")
			return
		}
		resp.Diagnostics = append(resp.Diagnostics, errorutil.ImportAsExistsError(existing.Id.String()))
	}

	r.Log(ctx, "look up the process template id")

	process, err := r.lookupProcess(ctx, func(p core.Process) bool {
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

	r.Log(ctx, "create the project")

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
		resp.Diagnostics.AddError("Queue create project", err.Error())
		return
	}

	r.Log(ctx, "wait for the project creation")

	if err := r.waitOperation(ctx, operationRef); err != nil {
		resp.Diagnostics.AddError("Wait for project delete operation", err.Error())
		return
	}
}

func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state projectModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.Log(ctx, "get the project")

	project, err := r.Meta.CoreClient.GetProject(ctx, core.GetProjectArgs{
		ProjectId:           state.Name.ValueStringPointer(), // Always use "name" to get as "id" is not available at create time
		IncludeCapabilities: pointer.From(true),
	})
	if err != nil && !errorutil.WasNotFound(err) {
		resp.Diagnostics.AddError("Read project", err.Error())
		return
	}
	if errorutil.WasNotFound(err) {
		resp.Diagnostics = append(resp.Diagnostics, sdk.NewDiagResourceNotFound(r.ResourceType(), state.Name.ValueString()))
		return
	}

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

		r.Log(ctx, "look up the template name by id")

		process, err := r.lookupProcess(ctx, func(p core.Process) bool { return p.Id != nil && p.Id.String() == *templateId })
		if err != nil {
			resp.Diagnostics.AddError("Lookup process", err.Error())
			return
		}
		templateName = process.Name
	}

	var id *string
	if project.Id != nil {
		id = pointer.From(project.Id.String())
	}

	state.Id = fwtype.StringValue(id)
	state.Name = adocustomtype.StringCaseInsensitiveValue{
		StringValue: fwtype.StringValue(project.Name),
	}
	state.Description = fwtype.StringValue(project.Description)
	state.Visibility = fwtype.StringValue(project.Visibility)
	state.VersionControl = fwtype.StringValue(versionControl)
	state.WorkItemTemplate = fwtype.StringValue(templateName)
	state.ProcessTemplateId = fwtype.StringValue(templateId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan projectModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state projectModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := uuid.Parse(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Parse resource id as UUID", err.Error())
		return
	}

	r.Log(ctx, "update the project")

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
		resp.Diagnostics.AddError("Queue update project", err.Error())
		return
	}

	r.Log(ctx, "wait for the project update")

	if err := r.waitOperation(ctx, operationRef); err != nil {
		resp.Diagnostics.AddError("Wait for project update operation", err.Error())
		return
	}
}

func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state projectModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.Log(ctx, "delete the prject")

	id, err := uuid.Parse(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Parse resource id as UUID", err.Error())
		return
	}

	operationRef, err := r.Meta.CoreClient.QueueDeleteProject(ctx, core.QueueDeleteProjectArgs{ProjectId: &id})
	if err != nil {
		resp.Diagnostics.AddError("Queue delete project", err.Error())
		return
	}

	r.Log(ctx, "wait for the project deletion")

	if err := r.waitOperation(ctx, operationRef); err != nil {
		resp.Diagnostics.AddError("Wait for project delete operation", err.Error())
		return
	}
}

func (r *projectResource) lookupProcess(ctx context.Context, f func(p core.Process) bool) (*core.Process, error) {
	processes, err := r.Meta.CoreClient.GetProcesses(ctx, core.GetProcessesArgs{})
	if err != nil {
		return nil, err
	}
	if processes == nil {
		return nil, errors.New("unexpected null processes")
	}

	for _, process := range *processes {
		if f(process) {
			return &process, nil
		}
	}

	return nil, errors.New("process not found")
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

		r.Log(ctx, "polling project operation status", map[string]any{
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
