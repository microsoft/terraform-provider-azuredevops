package core

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/terraform-provider-azuredevops/internal/adovalidator"
	"github.com/microsoft/terraform-provider-azuredevops/internal/framework"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/fwtype"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/pointer"
)

var _ framework.DataSource = &projectDataSource{}

func NewProjectDataSource() framework.DataSource {
	return &projectDataSource{}
}

type projectDataSource struct {
	framework.ImplSetMeta
	framework.ImplDataSourceMetadata
	framework.ImplLog[*projectDataSource]
}

type projectDataSourceModel struct {
	Name              types.String   `tfsdk:"name"`
	ProjectId         types.String   `tfsdk:"project_id"`
	Description       types.String   `tfsdk:"description"`
	Visibility        types.String   `tfsdk:"visibility"`
	VersionControl    types.String   `tfsdk:"version_control"`
	WorkItemTemplate  types.String   `tfsdk:"work_item_template"`
	ProcessTemplateId types.String   `tfsdk:"process_template_id"`
	Features          types.Map      `tfsdk:"features"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}

func (d *projectDataSource) ResourceType() string {
	return "azuredevops_project"
}

func (d *projectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the project.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("name"),
						path.MatchRoot("project_id"),
					),
				},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The id of the project.",
				Optional:            true,
				Validators: []validator.String{
					adovalidator.StringIsUUID(),
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("name"),
						path.MatchRoot("project_id"),
					),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the project.",
				Computed:            true,
			},
			"visibility": schema.StringAttribute{
				MarkdownDescription: "The visibility of the project.",
				Computed:            true,
			},
			"version_control": schema.StringAttribute{
				MarkdownDescription: "The version control system of the project.",
				Computed:            true,
			},
			"work_item_template": schema.StringAttribute{
				MarkdownDescription: "The work item template name.",
				Computed:            true,
			},
			"process_template_id": schema.StringAttribute{
				MarkdownDescription: "The Process Template ID used by the project.",
				Computed:            true,
			},
			"features": schema.MapAttribute{
				MarkdownDescription: "The features of this project.",
				Computed:            true,
			},
		},
	}
}

func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config projectDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	d.Info(ctx, "get the project")

	var projectId string
	switch {
	case !config.ProjectId.IsNull():
		projectId = config.ProjectId.ValueString()
	case !config.Name.IsNull():
		projectId = config.Name.ValueString()
	}

	project, err := d.Meta.CoreClient.GetProject(ctx, core.GetProjectArgs{
		ProjectId:           &projectId,
		IncludeCapabilities: pointer.From(true),
	})
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, framework.NewDiagSdkError("Read the project", err))
		return
	}

	var id *string
	if project.Id != nil {
		id = pointer.From(project.Id.String())
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
		process, err := LookupProcess(ctx, d.Meta.CoreClient, func(p core.Process) bool { return p.Id != nil && p.Id.String() == *templateId })
		if err != nil {
			resp.Diagnostics.AddError("Lookup process", err.Error())
			return
		}
		templateName = process.Name
	}

	config = projectDataSourceModel{
		Name:              fwtype.StringValue(project.Name),
		ProjectId:         fwtype.StringValue(id),
		Description:       fwtype.StringValue(project.Description),
		Visibility:        fwtype.StringValue(project.Visibility),
		VersionControl:    fwtype.StringValue(versionControl),
		WorkItemTemplate:  fwtype.StringValue(templateName),
		ProcessTemplateId: fwtype.StringValue(templateId),
		Features:          "", //TODO,
		Timeouts:          config.Timeouts,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}
