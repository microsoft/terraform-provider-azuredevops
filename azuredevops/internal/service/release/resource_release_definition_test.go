package release

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
	_tmp "github.com/microsoft/terraform-provider-azuredevops/.tmp"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
	"github.com/stretchr/testify/require"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/release"
)

var testReleaseProjectID = uuid.New().String()

// format and parse to remove monotonic clock
var now, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

var properties interface{}

var testReleaseDefinition = release.ReleaseDefinition{
	Id:             converter.Int(100),
	Revision:       converter.Int(1),
	Name:           converter.String("Name"),
	Path:           converter.String("\\"),
	VariableGroups: &[]int{},
	Source:         &release.ReleaseDefinitionSourceValues.RestApi,
	Description:    converter.String("Description"),
	Variables: &map[string]release.ConfigurationVariableValue{
		"releaseLevel": {
			Value:         converter.String("$(System.DefaultWorkingDirectory)/Directory"),
			IsSecret:      converter.Bool(false),
			AllowOverride: converter.Bool(false),
		},
	},
	Environments: &[]release.ReleaseDefinitionEnvironment{
		{
			Id:   converter.Int(50),
			Name: converter.String("Stage 1"),
			Rank: converter.Int(1),

			Variables: &map[string]release.ConfigurationVariableValue{
				"stageLevel": {
					Value:         converter.String("$(System.DefaultWorkingDirectory)/Directory"),
					IsSecret:      converter.Bool(false),
					AllowOverride: converter.Bool(false),
				},
			},
			VariableGroups: &[]int{},

			PreDeployApprovals: &release.ReleaseDefinitionApprovals{
				Approvals: &[]release.ReleaseDefinitionApprovalStep{{
					Id:               converter.Int(101),
					Rank:             converter.Int(1),
					IsAutomated:      converter.Bool(true),
					IsNotificationOn: converter.Bool(false),
				}},
				ApprovalOptions: &release.ApprovalOptions{
					RequiredApproverCount:                                   converter.Int(1),
					ReleaseCreatorCanBeApprover:                             converter.Bool(false),
					AutoTriggeredAndPreviousEnvironmentApprovedCanBeSkipped: converter.Bool(false),
					EnforceIdentityRevalidation:                             converter.Bool(false),
					TimeoutInMinutes:                                        converter.Int(0),
					ExecutionOrder:                                          &release.ApprovalExecutionOrderValues.BeforeGates,
				},
			},
			PostDeployApprovals: &release.ReleaseDefinitionApprovals{
				Approvals: &[]release.ReleaseDefinitionApprovalStep{{
					Id:               converter.Int(202),
					Rank:             converter.Int(1),
					IsAutomated:      converter.Bool(true),
					IsNotificationOn: converter.Bool(false),
				}},
				ApprovalOptions: &release.ApprovalOptions{
					RequiredApproverCount:                                   converter.Int(1),
					ReleaseCreatorCanBeApprover:                             converter.Bool(false),
					AutoTriggeredAndPreviousEnvironmentApprovedCanBeSkipped: converter.Bool(false),
					EnforceIdentityRevalidation:                             converter.Bool(false),
					TimeoutInMinutes:                                        converter.Int(0),
					ExecutionOrder:                                          &release.ApprovalExecutionOrderValues.AfterSuccessfulGates,
				},
			},
			DeployStep: &release.ReleaseDefinitionDeployStep{
				Id:    converter.Int(303),
				Tasks: &[]release.WorkflowTask{},
			},

			Demands: &[]interface{}{},
			// CurrentRelease: &release.ReleaseShallowReference{},
			Conditions: &[]release.Condition{
				{
					Name:          converter.String("ReleaseStarted"),
					ConditionType: &release.ConditionTypeValues.Event,
					Value:         converter.String(""),
				},
				{
					Name:          converter.String("Build Artifact"),
					ConditionType: &release.ConditionTypeValues.Artifact,
					Value:         converter.String("{\"sourceBranch\":\"master\",\"tags\":[\"major\",\"minor\",\"patch\"],\"useBuildDefinitionBranch\":false,\"createReleaseOnBuildTagging\":false}"),
				},
				{
					Name:          converter.String("Build Artifact"),
					ConditionType: &release.ConditionTypeValues.Artifact,
					Value:         converter.String("{\"sourceBranch\":\"develop\",\"tags\":[],\"useBuildDefinitionBranch\":false,\"createReleaseOnBuildTagging\":false}"),
				},
				{
					Name:          converter.String("Build Artifact"),
					ConditionType: &release.ConditionTypeValues.Artifact,
					Value:         converter.String("{\"sourceBranch\":\"-poc*\",\"tags\":[],\"useBuildDefinitionBranch\":false,\"createReleaseOnBuildTagging\":false}"),
				},
				{
					Name:          converter.String("Another Artifact"),
					ConditionType: &release.ConditionTypeValues.Artifact,
					Value:         converter.String("{\"sourceBranch\":\"master\",\"tags\":[\"tag\"],\"useBuildDefinitionBranch\":false,\"createReleaseOnBuildTagging\":false}"),
				},
			},

			Properties: map[string]interface{}{
				"BoardsEnvironmentType": map[string]interface{}{
					"$type":  "System.String",
					"$value": string(DeploymentTypeValues.Development),
				},
				"LinkBoardsWorkItems": map[string]interface{}{
					"$type":  "System.String",
					"$value": "true",
				},
				"JiraEnvironmentType": map[string]interface{}{
					"$type":  "System.String",
					"$value": string(DeploymentTypeValues.Unmapped),
				},
				"LinkJiraWorkItems": map[string]interface{}{
					"$type":  "System.String",
					"$value": "true",
				},
			},

			//Schedules: &[]release.ReleaseSchedule{},

			RetentionPolicy: &release.EnvironmentRetentionPolicy{
				DaysToKeep:     converter.Int(1),
				ReleasesToKeep: converter.Int(2),
				RetainBuild:    converter.Bool(false),
			},

			PreDeploymentGates: &release.ReleaseDefinitionGatesStep{
				Id:           converter.Int(1),
				GatesOptions: nil,
				Gates:        &[]release.ReleaseDefinitionGate{},
			},

			PostDeploymentGates: &release.ReleaseDefinitionGatesStep{
				Id:           converter.Int(1),
				GatesOptions: nil,
				Gates:        &[]release.ReleaseDefinitionGate{},
			},
			Owner: &webapi.IdentityRef{
				Id: converter.String(testReleaseProjectID),
			},

			DeployPhases: &[]interface{}{
				map[string]interface{}{
					"deploymentInput": map[string]interface{}{
						"parallelExecution": map[string]interface{}{
							"parallelExecutionType": "none",
						},
						"agentSpecification": map[string]interface{}{
							"identifier": "ubuntu-18.04",
						},
						"skipArtifactsDownload": false,
						"artifactsDownloadInput": map[string]interface{}{
							"downloadInputs": []interface{}{
								map[string]interface{}{
									"alias":                "Build Artifact",
									"artifactType":         strings.Title(string(release.AgentArtifactTypeValues.Build)),
									"artifactDownloadMode": string(ArtifactDownloadModeTypeValues.All),
									"artifactItems":        []interface{}{},
								},
							},
						},
						"queueId":                   float64(52),
						"demands":                   []interface{}{"sh", "browser -equals chrome"},
						"enableAccessToken":         false,
						"timeoutInMinutes":          float64(0),
						"jobCancelTimeoutInMinutes": float64(1),
						"condition":                 "succeeded()",
						// "overrideInputs": interface{}{},
					},
					"rank":      float64(1),
					"phaseType": "agentBasedDeployment",
					"name":      "TERRAFORM",
					"workflowTasks": []interface{}{
						map[string]interface{}{
							"environment":      map[string]interface{}{},
							"taskId":           "d9bafed4-0b18-4f58-968d-86655b4d2ce9",
							"version":          "2.*",
							"name":             "Command Line Task",
							"enabled":          true,
							"alwaysRun":        false,
							"continueOnError":  false,
							"timeoutInMinutes": float64(0),
							"definitionType":   "task",
							"overrideInputs":   map[string]interface{}{},
							"condition":        "succeeded()",
							"inputs": map[string]interface{}{
								"script":           "chmod 0777 -R .terraform/plugins",
								"workingDirectory": "$(artifactRoot)",
								"failOnStderr":     "false",
							},
						},
					},
				},
			},
			EnvironmentOptions: &release.EnvironmentOptions{
				AutoLinkWorkItems:            converter.Bool(true),
				BadgeEnabled:                 converter.Bool(true),
				PublishDeploymentStatus:      converter.Bool(true),
				PullRequestDeploymentEnabled: converter.Bool(false),
			},
		},
		{
			Id:   converter.Int(51),
			Name: converter.String("Stage 2"),
			Rank: converter.Int(2),

			Variables:      &map[string]release.ConfigurationVariableValue{},
			VariableGroups: &[]int{},

			PreDeployApprovals: &release.ReleaseDefinitionApprovals{
				Approvals: &[]release.ReleaseDefinitionApprovalStep{{
					Id:               converter.Int(1001),
					Rank:             converter.Int(1),
					IsAutomated:      converter.Bool(true),
					IsNotificationOn: converter.Bool(false),
				}},
			},
			PostDeployApprovals: &release.ReleaseDefinitionApprovals{
				Approvals: &[]release.ReleaseDefinitionApprovalStep{{
					Id:               converter.Int(2002),
					Rank:             converter.Int(1),
					IsAutomated:      converter.Bool(true),
					IsNotificationOn: converter.Bool(false),
				}},
			},
			DeployStep: &release.ReleaseDefinitionDeployStep{
				Id:    converter.Int(3003),
				Tasks: &[]release.WorkflowTask{},
			},

			Demands: &[]interface{}{},
			Conditions: &[]release.Condition{
				{
					Name:          converter.String("Stage 1"),
					ConditionType: &release.ConditionTypeValues.EnvironmentState,
					Value:         converter.String("4"),
				},
			},

			Properties: properties,

			//Schedules: &[]release.ReleaseSchedule{},

			RetentionPolicy: &release.EnvironmentRetentionPolicy{
				DaysToKeep:     converter.Int(1),
				ReleasesToKeep: converter.Int(2),
				RetainBuild:    converter.Bool(false),
			},

			PreDeploymentGates: &release.ReleaseDefinitionGatesStep{
				Id:           converter.Int(1),
				GatesOptions: nil,
				Gates:        &[]release.ReleaseDefinitionGate{},
			},

			PostDeploymentGates: &release.ReleaseDefinitionGatesStep{
				Id:           converter.Int(1),
				GatesOptions: nil,
				Gates:        &[]release.ReleaseDefinitionGate{},
			},
			Owner: &webapi.IdentityRef{
				Id: converter.String(testReleaseProjectID),
			},
			DeployPhases: &[]interface{}{},
		},
	},
	Triggers:          &[]interface{}{},
	Tags:              &[]string{},
	ReleaseNameFormat: converter.String("Release-$(rev:r)"),
	Url:               converter.String(fmt.Sprintf("https://vsrm.dev.azure.com/Demo/%s/_apis/Release/definitions/2", testReleaseProjectID)),
	Properties: map[string]interface{}{
		"DefinitionCreationSource": map[string]interface{}{
			"$type":  "System.String",
			"$value": "ReleaseNew",
		},
		"IntegrateBoardsWorkItems": map[string]interface{}{
			"$type":  "System.String",
			"$value": "true",
		},
		"IntegrateJiraWorkItems": map[string]interface{}{
			"$type":  "System.String",
			"$value": "true",
		},
		"JiraServiceEndpointId": map[string]interface{}{
			"$type":  "System.String",
			"$value": uuid.New().String(),
		},
	},
	IsDeleted:  converter.Bool(false),
	Comment:    converter.String("Changes made by Terraform"),
	CreatedOn:  &azuredevops.Time{Time: now},
	ModifiedOn: &azuredevops.Time{Time: now},
	Artifacts: &[]release.Artifact{
		{
			Alias: converter.String("Build Artifact"),
			DefinitionReference: &map[string]release.ArtifactSourceReference{
				/*
					"artifactSourceDefinitionUrl": {
						Id:   converter.String(""),
						Name: converter.String(""),
					},
				*/
				"defaultVersionBranch": {
					Id: converter.String(""),
					//Name: converter.String(""),
				},
				"defaultVersionSpecific": {
					Id: converter.String(""),
					//Name: converter.String(""),
				},
				"defaultVersionTags": {
					Id: converter.String(""),
					//Name: converter.String(""),
				},
				"defaultVersionType": {
					Id: converter.String("selectDuringReleaseCreationType"),
					//Name: converter.String("Specify at the time of release creation"),
				},
				"definition": {
					Id: converter.String("190"),
					//Name: converter.String("Build Pipeline"),
				},
				"definitions": {
					Id: converter.String(""),
					//Name: converter.String(""),
				},
				"IsMultiDefinitionType": {
					Id: converter.String("False"),
					//Name: converter.String("False"),
				},
				"project": {
					Id: converter.String(testReleaseProjectID),
					//Name: converter.String("Project"),
				},
				/*
					"repository": {
						Id:   converter.String(""),
						Name: converter.String(""),
					},
				*/
			},
			IsPrimary:  converter.Bool(true),
			IsRetained: converter.Bool(false),
			Type:       converter.String(strings.Title(string(release.AgentArtifactTypeValues.Build))),
		},
	},
}

// verifies that the flatten/expand round trip yields the same release definition
func TestAzureDevOpsReleaseDefinition_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceReleaseDefinition().Schema, nil)
	flattenReleaseDefinition(resourceData, &testReleaseDefinition, testReleaseProjectID)

	releaseDefinitionAfterRoundTrip, projectID, err := expandReleaseDefinition(resourceData)

	sortedExpected := sortReleaseDefinition(testReleaseDefinition)
	sortedActual := sortReleaseDefinition(*releaseDefinitionAfterRoundTrip)

	_tmp.Log(sortedExpected)
	_tmp.Log(sortedActual)

	require.Nil(t, err)
	require.Equal(t, sortedExpected, sortedActual)
	require.Equal(t, testReleaseProjectID, projectID)
}

func sortReleaseDefinition(b release.ReleaseDefinition) release.ReleaseDefinition {
	if b.Environments != nil {
		for _, e := range *b.Environments {
			if e.Conditions != nil {
				for _, c := range *e.Conditions {
					if c.Value != nil && strings.Contains(*c.Value, "{") {
						var m map[string]interface{}
						if err := json.Unmarshal([]byte(*c.Value), &m); err == nil {
							if m2, ok := m["tags"].([]interface{}); ok {
								tags := tfhelper.ExpandStringList(m2)
								sort.Strings(tags)
								m["tags"] = tags
							}
							value, _ := json.Marshal(m)
							*c.Value = string(value)
						}
					}
				}
			}
		}
	}
	return b
}
