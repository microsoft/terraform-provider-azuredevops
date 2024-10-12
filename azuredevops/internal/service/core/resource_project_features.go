package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/featuremanagement"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ProjectFeatureType Project feature in Azure DevOps
type ProjectFeatureType string

type projectFeatureTypeValuesType struct {
	Boards       ProjectFeatureType
	Repositories ProjectFeatureType
	Pipelines    ProjectFeatureType
	TestPlans    ProjectFeatureType
	Artifacts    ProjectFeatureType
}

// ProjectFeatureTypeValues valid projects features in Azure DevOps
var ProjectFeatureTypeValues = projectFeatureTypeValuesType{
	Boards:       "boards",
	Repositories: "repositories",
	Pipelines:    "pipelines",
	TestPlans:    "testplans",
	Artifacts:    "artifacts",
}

var projectFeatureNameMap = map[string]ProjectFeatureType{
	"ms.vss-work.agile":           ProjectFeatureTypeValues.Boards,
	"ms.vss-code.version-control": ProjectFeatureTypeValues.Repositories,
	"ms.vss-build.pipelines":      ProjectFeatureTypeValues.Pipelines,
	"ms.vss-test-web.test":        ProjectFeatureTypeValues.TestPlans,
	"ms.azure-artifacts.feature":  ProjectFeatureTypeValues.Artifacts,
}

var projectFeatureNameMapReverse = map[ProjectFeatureType]string{
	ProjectFeatureTypeValues.Boards:       "ms.vss-work.agile",
	ProjectFeatureTypeValues.Repositories: "ms.vss-code.version-control",
	ProjectFeatureTypeValues.Pipelines:    "ms.vss-build.pipelines",
	ProjectFeatureTypeValues.TestPlans:    "ms.vss-test-web.test",
	ProjectFeatureTypeValues.Artifacts:    "ms.azure-artifacts.feature",
}

// ResourceProjectFeatures schema and implementation for project features
func ResourceProjectFeatures() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectFeaturesCreateUpdate,
		ReadContext:   resourceProjectFeaturesRead,
		UpdateContext: resourceProjectFeaturesCreateUpdate,
		DeleteContext: resourceProjectFeaturesDelete,
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
			"features": {
				Type:         schema.TypeMap,
				Required:     true,
				ValidateFunc: validateProjectFeatures,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceProjectFeaturesCreateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	featureStates := d.Get("features").(map[string]interface{})

	err := updateProjectFeatureStates(ctx, clients.FeatureManagementClient, projectID, &featureStates)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(projectID)
	return resourceProjectFeaturesRead(ctx, d, m)
}

func resourceProjectFeaturesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	featureStates := d.Get("features").(map[string]interface{})
	currentFeatureStates, err := getConfiguredProjectFeatureStates(ctx, clients.FeatureManagementClient, &featureStates, projectID)
	if err != nil {
		return diag.FromErr(err)
	}
	if currentFeatureStates == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf(" failed to retrieve current feature states for project: %s", projectID))
	}
	d.Set("features", currentFeatureStates)
	return nil
}

func resourceProjectFeaturesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	featureStates := d.Get("features").(map[string]interface{})
	for k := range featureStates {
		featureStates[k] = string(featuremanagement.ContributedFeatureEnabledValueValues.Enabled)
	}
	err := updateProjectFeatureStates(ctx, clients.FeatureManagementClient, projectID, &featureStates)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func getConfiguredProjectFeatureStates(ctx context.Context, fc featuremanagement.Client, featureStates *map[string]interface{}, projectID string) (*map[ProjectFeatureType]featuremanagement.ContributedFeatureEnabledValue, error) {
	if featureStates == nil {
		return nil, nil
	}

	currentFeatureStates, err := getProjectFeatureStates(ctx, fc, projectID)
	if err != nil {
		return nil, err
	}

	for k := range *currentFeatureStates {
		if _, ok := (*featureStates)[string(k)]; !ok {
			delete(*currentFeatureStates, k)
		}
	}
	return currentFeatureStates, nil
}

func updateProjectFeatureStates(ctx context.Context, fc featuremanagement.Client, projectID string, featureStates *map[string]interface{}) error {
	if featureStates == nil {
		return nil
	}
	for k, v := range *featureStates {
		enabledValue := featuremanagement.ContributedFeatureEnabledValue(v.(string))
		f, ok := projectFeatureNameMapReverse[ProjectFeatureType(k)]
		if !ok {
			return fmt.Errorf(" unknown feature: %s, available features are: `boards`, `repositories`,`pipelines`,`testplans`,`artifacts`", k)
		}
		//TODO handle response state
		_, err := fc.SetFeatureStateForScope(ctx, featuremanagement.SetFeatureStateForScopeArgs{
			Feature: &featuremanagement.ContributedFeatureState{
				FeatureId: converter.String(f),
				State:     &enabledValue,
				Scope: &featuremanagement.ContributedFeatureSettingScope{
					SettingScope: converter.String("project"),
					UserScoped:   converter.Bool(false),
				},
			},
			FeatureId:  converter.String(f),
			UserScope:  converter.String("host"),
			ScopeName:  converter.String("project"),
			ScopeValue: &projectID,
		})
		if nil != err {
			return fmt.Errorf(" Faild to update project features. Feature type: %s,  Error: %+v", f, err)
		}
	}
	return nil
}

func getProjectFeatureStates(ctx context.Context, fc featuremanagement.Client, projectID string) (*map[ProjectFeatureType]featuremanagement.ContributedFeatureEnabledValue, error) {
	states, err := fc.QueryFeatureStates(ctx, featuremanagement.QueryFeatureStatesArgs{
		Query: &featuremanagement.ContributedFeatureStateQuery{
			FeatureIds: &[]string{
				"ms.vss-work.agile",
				"ms.vss-code.version-control",
				"ms.vss-build.pipelines",
				"ms.vss-test-web.test",
				"ms.azure-artifacts.feature",
			},
			ScopeValues: &map[string]string{
				"project": projectID,
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf(" Get project features error, project: %s, error: %+v", projectID, err)
	}

	featureStates := make(map[ProjectFeatureType]featuremanagement.ContributedFeatureEnabledValue)
	for k, v := range projectFeatureNameMap {
		if state, ok := (*states.FeatureStates)[k]; ok {
			featureStates[v] = *state.State
		}
	}
	return &featureStates, nil
}

func getDefaultProjectFeatureStates(states *map[string]interface{}) (*map[string]interface{}, error) {
	featureStates := map[string]interface{}{}
	for k := range projectFeatureNameMapReverse {
		if states != nil {
			if _, ok := (*states)[string(k)]; !ok {
				continue
			}
		}
		featureStates[string(k)] = string(featuremanagement.ContributedFeatureEnabledValueValues.Enabled)
	}
	return &featureStates, nil
}

func validateProjectFeatures(i interface{}, k string) ([]string, []error) {
	var errors []error
	var warnings []string

	m := i.(map[string]interface{})

	if len(m) <= 0 {
		errors = append(errors, fmt.Errorf("Feature map must contain at least on entry"))
	}
	for feature, state := range m {
		if _, ok := projectFeatureNameMapReverse[ProjectFeatureType(strings.ToLower(feature))]; !ok {
			errors = append(errors, fmt.Errorf("unknown feature: %s, available features are: `boards`, `repositories`,`pipelines`,`testplans`,`artifacts` ", feature))
		}

		if state != string(featuremanagement.ContributedFeatureEnabledValueValues.Enabled) &&
			state != string(featuremanagement.ContributedFeatureEnabledValueValues.Disabled) {
			errors = append(errors, fmt.Errorf("invalid state: %s", state))
		}
	}
	return warnings, errors
}
