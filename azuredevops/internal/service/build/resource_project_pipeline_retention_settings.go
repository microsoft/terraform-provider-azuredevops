package build

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/build"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceProjectPipelineRetentionSettings schema and implementation for project pipeline retention settings resource
func ResourceProjectPipelineRetentionSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectPipelineRetentionSettingsCreateUpdate,
		ReadContext:   resourceProjectPipelineRetentionSettingsRead,
		UpdateContext: resourceProjectPipelineRetentionSettingsCreateUpdate,
		DeleteContext: resourceProjectPipelineRetentionSettingsDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"run_retention": {
				Description: "The number of days to retain pipeline runs",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"artifact_retention": {
				Description: "The number of days to retain artifacts. Artifacts can not live longer than a run, so will be overridden by a shorter run retention setting",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"pull_request_run_retention": {
				Description: "The number of days to retain pull request pipeline runs",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"retain_runs_per_protected_branch": {
				Description: "The number of runs to retain per protected branch",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceProjectPipelineRetentionSettingsCreateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)

	err := configureProjectPipelineRetentionSettings(clients, projectID, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("creating/updating project pipeline retention settings: %v", err))
	}
	d.SetId(projectID)
	return resourceProjectPipelineRetentionSettingsRead(ctx, d, m)
}

func resourceProjectPipelineRetentionSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectId := d.Id()
	getSettings := build.GetRetentionSettingsArgs{
		Project: converter.String(projectId),
	}

	retentionSettings, err := clients.BuildClient.GetRetentionSettings(ctx, getSettings)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("Error reading project pipeline retention settings: %v", err))
	}

	d.Set("project_id", projectId)
	if retentionSettings.PurgeRuns != nil {
		d.Set("run_retention", retentionSettings.PurgeRuns.Value)
	}
	if retentionSettings.PurgeArtifacts != nil {
		d.Set("artifact_retention", retentionSettings.PurgeArtifacts.Value)
	}
	if retentionSettings.PurgePullRequestRuns != nil {
		d.Set("pull_request_run_retention", retentionSettings.PurgePullRequestRuns.Value)
	}
	if retentionSettings.RetainRunsPerProtectedBranch != nil {
		d.Set("retain_runs_per_protected_branch", retentionSettings.RetainRunsPerProtectedBranch.Value)
	}
	return nil
}

func resourceProjectPipelineRetentionSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// nothing to do, as the original settings are unknown.
	return nil
}

func configureProjectPipelineRetentionSettings(clients *client.AggregatedClient, projectId string, d *schema.ResourceData) error {
	settings := build.UpdateRetentionSettingsArgs{
		Project:     converter.String(projectId),
		UpdateModel: &build.UpdateProjectRetentionSettingModel{},
	}

	rawConfig := d.GetRawConfig().AsValueMap()

	runRetention := rawConfig["run_retention"]
	if !runRetention.IsNull() {
		value, _ := runRetention.AsBigFloat().Int64()
		settings.UpdateModel.RunRetention = &build.UpdateRetentionSettingModel{Value: converter.Int(int(value))}
	}

	artifactRetention := rawConfig["artifact_retention"]
	if !artifactRetention.IsNull() {
		value, _ := artifactRetention.AsBigFloat().Int64()
		settings.UpdateModel.ArtifactsRetention = &build.UpdateRetentionSettingModel{Value: converter.Int(int(value))}
	}

	pullRequestRunRetention := rawConfig["pull_request_run_retention"]
	if !pullRequestRunRetention.IsNull() {
		value, _ := pullRequestRunRetention.AsBigFloat().Int64()
		settings.UpdateModel.PullRequestRunRetention = &build.UpdateRetentionSettingModel{Value: converter.Int(int(value))}
	}

	retainRunsPerProtectedBranch := rawConfig["retain_runs_per_protected_branch"]
	if !retainRunsPerProtectedBranch.IsNull() {
		value, _ := retainRunsPerProtectedBranch.AsBigFloat().Int64()
		settings.UpdateModel.RetainRunsPerProtectedBranch = &build.UpdateRetentionSettingModel{Value: converter.Int(int(value))}
	}

	_, err := clients.BuildClient.UpdateRetentionSettings(clients.Ctx, settings)
	if err != nil {
		return err
	}

	return nil
}
