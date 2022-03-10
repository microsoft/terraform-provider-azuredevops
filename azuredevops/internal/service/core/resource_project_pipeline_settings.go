package core

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/build"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
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
		},
	}
}

func resourceProjectPipelineSettingsCreateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)

	err := configureProjectPipelineGeneralSettings(clients, projectID, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error creating/updating project build general settings: %v", err))
	}
	d.SetId(projectID)
	return nil
}

func resourceProjectPipelineSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectId := d.Id()
	getSettings := build.GetBuildGeneralSettingsArgs{
		Project: converter.String(projectId),
	}

	buildSettings, err := clients.BuildClient.GetBuildGeneralSettings(ctx, getSettings)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error reading project build general settings: %v", err))
	}

	d.Set("project_id", projectId)
	d.Set("enforce_job_scope", buildSettings.EnforceJobAuthScope)
	d.Set("enforce_referenced_repo_scoped_token", buildSettings.EnforceReferencedRepoScopedToken)
	d.Set("enforce_settable_var", buildSettings.EnforceSettableVar)
	d.Set("publish_pipeline_metadata", buildSettings.PublishPipelineMetadata)
	d.Set("status_badges_are_private", buildSettings.StatusBadgesArePrivate)

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

	//lint:ignore SA1019 - must use GetOkExists to determine if the boolean property has been set from the resource
	if val, exists := d.GetOkExists("enforce_job_scope"); exists { //nolint:staticcheck
		settings.NewSettings.EnforceJobAuthScope = converter.Bool(val.(bool))
	}

	//lint:ignore SA1019 - must use GetOkExists to determine if the boolean property has been set from the resource
	if val, exists := d.GetOkExists("enforce_referenced_repo_scoped_token"); exists { //nolint:staticcheck
		settings.NewSettings.EnforceReferencedRepoScopedToken = converter.Bool(val.(bool))
	}

	//lint:ignore SA1019 - must use GetOkExists to determine if the boolean property has been set from the resource
	if val, exists := d.GetOkExists("enforce_settable_var"); exists { //nolint:staticcheck
		settings.NewSettings.EnforceSettableVar = converter.Bool(val.(bool))
	}

	//lint:ignore SA1019 - must use GetOkExists to determine if the boolean property has been set from the resource
	if val, exists := d.GetOkExists("publish_pipeline_metadata"); exists { //nolint:staticcheck
		settings.NewSettings.PublishPipelineMetadata = converter.Bool(val.(bool))
	}

	//lint:ignore SA1019 - must use GetOkExists to determine if the boolean property has been set from the resource
	if val, exists := d.GetOkExists("status_badges_are_private"); exists { //nolint:staticcheck
		settings.NewSettings.StatusBadgesArePrivate = converter.Bool(val.(bool))
	}

	_, err := clients.BuildClient.UpdateBuildGeneralSettings(clients.Ctx, settings)
	if err != nil {
		return err
	}

	return nil
}
