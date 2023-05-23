package core

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/build"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceProjectPipelineSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectPipelineSettingsCreateUpdate,
		ReadContext:   resourceProjectPipelineSettingsRead,
		UpdateContext: resourceProjectPipelineSettingsCreateUpdate,
		DeleteContext: resourceProjectPipelineSettingsDelete,
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
			"enforce_job_scope": {
				Description: "Limit job authorization scope to current project for non-release pipelines",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"enforce_referenced_repo_scoped_token": {
				Description: "Protect access to repositories in YAML pipelines",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"enforce_settable_var": {
				Description: "Limit variables that can be set at queue time",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"publish_pipeline_metadata": {
				Description: "Publish metadata from pipelines",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"status_badges_are_private": {
				Description: "Disable anonymous access to badges",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"enforce_job_scope_for_release": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceProjectPipelineSettingsCreateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)

	err := configureProjectPipelineGeneralSettings(clients, projectID, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf(" creating/updating project build general settings: %v", err))
	}
	d.SetId(projectID)
	return resourceProjectPipelineSettingsRead(ctx, d, m)
}

func resourceProjectPipelineSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectId := d.Id()
	getSettings := build.GetBuildGeneralSettingsArgs{
		Project: converter.String(projectId),
	}

	buildSettings, err := clients.BuildClient.GetBuildGeneralSettings(ctx, getSettings)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("Error reading project build general settings: %v", err))
	}

	d.Set("project_id", projectId)
	d.Set("enforce_job_scope", buildSettings.EnforceJobAuthScope)
	d.Set("enforce_referenced_repo_scoped_token", buildSettings.EnforceReferencedRepoScopedToken)
	d.Set("enforce_settable_var", buildSettings.EnforceSettableVar)
	d.Set("publish_pipeline_metadata", buildSettings.PublishPipelineMetadata)
	d.Set("status_badges_are_private", buildSettings.StatusBadgesArePrivate)
	d.Set("enforce_job_scope_for_release", buildSettings.EnforceJobAuthScopeForReleases)
	return nil
}

func resourceProjectPipelineSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// nothing to do, as the original settings are unknown.
	return nil
}

func configureProjectPipelineGeneralSettings(clients *client.AggregatedClient, projectId string, d *schema.ResourceData) error {
	settings := build.UpdateBuildGeneralSettingsArgs{
		Project:     converter.String(projectId),
		NewSettings: &build.PipelineGeneralSettings{},
	}

	rawConfig := d.GetRawConfig().AsValueMap()
	enforceJobScope := rawConfig["enforce_job_scope"]
	if !enforceJobScope.IsNull() {
		settings.NewSettings.EnforceJobAuthScope = converter.Bool(enforceJobScope.True())
	}

	enforceReferencedRepoScopedToken := rawConfig["enforce_referenced_repo_scoped_token"]
	if !enforceReferencedRepoScopedToken.IsNull() {
		settings.NewSettings.EnforceReferencedRepoScopedToken = converter.Bool(enforceJobScope.True())
	}

	enforceSettableVar := rawConfig["enforce_settable_var"]
	if !enforceSettableVar.IsNull() {
		settings.NewSettings.EnforceSettableVar = converter.Bool(enforceSettableVar.True())
	}

	publishPipelineMetadata := rawConfig["publish_pipeline_metadata"]
	if !publishPipelineMetadata.IsNull() {
		settings.NewSettings.PublishPipelineMetadata = converter.Bool(publishPipelineMetadata.True())
	}

	statusBadgesArePrivate := rawConfig["status_badges_are_private"]
	if !statusBadgesArePrivate.IsNull() {
		settings.NewSettings.StatusBadgesArePrivate = converter.Bool(statusBadgesArePrivate.True())
	}

	enforceJobAuthScopeForReleases := rawConfig["enforce_job_scope_for_release"]
	if !statusBadgesArePrivate.IsNull() {
		settings.NewSettings.EnforceJobAuthScopeForReleases = converter.Bool(enforceJobAuthScopeForReleases.True())
	}

	_, err := clients.BuildClient.UpdateBuildGeneralSettings(clients.Ctx, settings)
	if err != nil {
		return err
	}

	return nil
}
