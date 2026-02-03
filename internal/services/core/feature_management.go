package core

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/featuremanagement"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/errorutil"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/pointer"
)

const (
	featureIdBoards    = "ms.vss-work.agile"
	featureIdRepos     = "ms.vss-code.version-control"
	featureIdPipelines = "ms.vss-build.pipelines"
	featureIdPlanTests = "ms.vss-test-web.test"
	featureIdArtifacts = "ms.azure-artifacts.feature"
)

type projectFeaturesTFModel struct {
	Boards    types.Bool `tfsdk:"boards"`
	Repos     types.Bool `tfsdk:"repos"`
	Pipelines types.Bool `tfsdk:"pipelines"`
	TestPlans types.Bool `tfsdk:"test_plans"`
	Artifacts types.Bool `tfsdk:"artifacts"`
}

// Key is the feature id
type projectFeaturesAPIModel map[string]*featuremanagement.ContributedFeatureEnabledValue

func (m projectFeaturesTFModel) ToAPIModel() projectFeaturesAPIModel {
	out := projectFeaturesAPIModel{}

	evaluateVal := func(v types.Bool) *featuremanagement.ContributedFeatureEnabledValue {
		if v.IsNull() {
			return pointer.From(featuremanagement.ContributedFeatureEnabledValueValues.Undefined)
		}
		if v.ValueBool() {
			return pointer.From(featuremanagement.ContributedFeatureEnabledValueValues.Enabled)
		}
		return pointer.From(featuremanagement.ContributedFeatureEnabledValueValues.Disabled)
	}

	out[featureIdBoards] = evaluateVal(m.Boards)
	out[featureIdRepos] = evaluateVal(m.Repos)
	out[featureIdPipelines] = evaluateVal(m.Pipelines)
	out[featureIdPlanTests] = evaluateVal(m.TestPlans)
	out[featureIdArtifacts] = evaluateVal(m.Artifacts)

	return out
}

func (m projectFeaturesAPIModel) ToTFModel() projectFeaturesTFModel {
	evaluateVal := func(v *featuremanagement.ContributedFeatureEnabledValue) types.Bool {
		if v == nil {
			return types.BoolNull()
		}
		switch *v {
		case featuremanagement.ContributedFeatureEnabledValueValues.Undefined:
			return types.BoolNull()
		case featuremanagement.ContributedFeatureEnabledValueValues.Enabled:
			return types.BoolValue(true)
		case featuremanagement.ContributedFeatureEnabledValueValues.Disabled:
			return types.BoolValue(false)
		default:
			panic(fmt.Sprintf("invalid feature value: %v", v))
		}
	}

	return projectFeaturesTFModel{
		Boards:    evaluateVal(m[featureIdBoards]),
		Repos:     evaluateVal(m[featureIdRepos]),
		Pipelines: evaluateVal(m[featureIdPipelines]),
		TestPlans: evaluateVal(m[featureIdPlanTests]),
		Artifacts: evaluateVal(m[featureIdArtifacts]),
	}
}

func setProjectFeature(ctx context.Context, client featuremanagement.Client, projectID string, features projectFeaturesTFModel) error {
	apiFeatures := features.ToAPIModel()
	for featureId, state := range apiFeatures {
		if state == nil || *state == featuremanagement.ContributedFeatureEnabledValueValues.Undefined {
			continue
		}
		_, err := client.SetFeatureStateForScope(ctx, featuremanagement.SetFeatureStateForScopeArgs{
			Feature: &featuremanagement.ContributedFeatureState{
				FeatureId: pointer.From(featureId),
				State:     state,
				Scope: &featuremanagement.ContributedFeatureSettingScope{
					SettingScope: pointer.From("project"),
					UserScoped:   pointer.From(false),
				},
			},
			FeatureId:  pointer.From(featureId),
			UserScope:  pointer.From("host"),
			ScopeName:  pointer.From("project"),
			ScopeValue: &projectID,
		})
		if nil != err {
			return fmt.Errorf("Faild to update project feature %q: %+v", featureId, err)
		}
	}
	return nil
}

func getProjectFeatures(ctx context.Context, client featuremanagement.Client, projectID string) (*types.Object, error) {
	states, err := client.QueryFeatureStatesForNamedScope(ctx, featuremanagement.QueryFeatureStatesForNamedScopeArgs{
		Query: &featuremanagement.ContributedFeatureStateQuery{
			FeatureIds: &[]string{
				featureIdBoards,
				featureIdRepos,
				featureIdPipelines,
				featureIdPlanTests,
				featureIdArtifacts,
			},
		},
		UserScope:  pointer.From("host"),
		ScopeName:  pointer.From("project"),
		ScopeValue: &projectID,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to get project features: %+v", err)
	}

	fstates := states.FeatureStates
	if fstates == nil {
		return nil, fmt.Errorf("unexpected nil feature states returned")
	}

	apiModel := projectFeaturesAPIModel{}
	for featureId, state := range *fstates {
		apiModel[featureId] = state.State
	}

	obj, diags := types.ObjectValueFrom(
		ctx,
		map[string]attr.Type{
			"boards":     types.BoolType,
			"repos":      types.BoolType,
			"pipelines":  types.BoolType,
			"test_plans": types.BoolType,
			"artifacts":  types.BoolType,
		},
		apiModel.ToTFModel(),
	)
	if diags.HasError() {
		return nil, errorutil.DiagsToError(diags)
	}
	return &obj, nil
}
