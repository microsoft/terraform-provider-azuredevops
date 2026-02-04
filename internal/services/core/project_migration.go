package core

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/featuremanagement"
	"github.com/microsoft/terraform-provider-azuredevops/internal/adocustomtype"
)

var _ resource.ResourceWithUpgradeState = &projectResource{}

type projectResourceModelV0 struct {
	Name              types.String   `tfsdk:"name"`
	Description       types.String   `tfsdk:"description"`
	Visibility        types.String   `tfsdk:"visibility"`
	VersionControl    types.String   `tfsdk:"version_control"`
	WorkItemTemplate  types.String   `tfsdk:"work_item_template"`
	Features          types.Map      `tfsdk:"features"`
	Id                types.String   `tfsdk:"id"`
	ProcessTemplateId types.String   `tfsdk:"process_template_id"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}

func (r *projectResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required: true,
					},
					"description": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"visibility": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"version_control": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"work_item_template": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"features": schema.MapAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
					"id": schema.StringAttribute{
						Computed: true,
					},
					"process_template_id": schema.StringAttribute{
						Computed: true,
					},
					"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
						Create: true,
						Read:   true,
						Update: true,
						Delete: true,
					}),
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var oldState projectResourceModelV0
				resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
				if resp.Diagnostics.HasError() {
					return
				}

				features := projectFeaturesTFModel{}
				for k, v := range oldState.Features.Elements() {
					switch k {
					case "boards":
						features.Boards = types.BoolValue(v.(types.String).ValueString() == string(featuremanagement.ContributedFeatureEnabledValueValues.Enabled))
					case "repositories":
						features.Repos = types.BoolValue(v.(types.String).ValueString() == string(featuremanagement.ContributedFeatureEnabledValueValues.Enabled))
					case "pipelines":
						features.Pipelines = types.BoolValue(v.(types.String).ValueString() == string(featuremanagement.ContributedFeatureEnabledValueValues.Enabled))
					case "testplans":
						features.TestPlans = types.BoolValue(v.(types.String).ValueString() == string(featuremanagement.ContributedFeatureEnabledValueValues.Enabled))
					case "artifacts":
						features.Artifacts = types.BoolValue(v.(types.String).ValueString() == string(featuremanagement.ContributedFeatureEnabledValueValues.Enabled))
					}
				}
				featuresObject, diags := types.ObjectValueFrom(
					ctx,
					map[string]attr.Type{
						"boards":       types.BoolType,
						"repositories": types.BoolType,
						"pipelines":    types.BoolType,
						"testplans":    types.BoolType,
						"artifacts":    types.BoolType,
					},
					features,
				)
				resp.Diagnostics.Append(diags...)
				if resp.Diagnostics.HasError() {
					return
				}

				newState := projectResourceModel{
					Name: adocustomtype.StringCaseInsensitiveValue{
						StringValue: oldState.Name,
					},
					Description:       oldState.Description,
					Visibility:        oldState.Visibility,
					VersionControl:    oldState.VersionControl,
					WorkItemTemplate:  oldState.WorkItemTemplate,
					Features:          featuresObject,
					Id:                oldState.Id,
					ProcessTemplateId: oldState.ProcessTemplateId,
					Timeouts:          oldState.Timeouts,
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
				if resp.Diagnostics.HasError() {
					return
				}
			},
		},
	}
}
