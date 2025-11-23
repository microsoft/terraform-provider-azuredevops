package work

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/work"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/mock/gomock"
)

type MockWorkClient struct {
	mock.Mock
}

func (m *MockWorkClient) GetTeamSettings(ctx context.Context, args work.GetTeamSettingsArgs) (*work.TeamSetting, error) {
	ret := m.Called(ctx, args)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*work.TeamSetting), ret.Error(1)
}

func (m *MockWorkClient) UpdateTeamSettings(ctx context.Context, args work.UpdateTeamSettingsArgs) (*work.TeamSetting, error) {
	ret := m.Called(ctx, args)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*work.TeamSetting), ret.Error(1)
}

func TestExpandTeamSettingsPatch(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceTeamSettings().Schema, map[string]interface{}{
		"bugs_behavior":        "asTasks",
		"working_days":         []interface{}{"monday", "tuesday"},
		"backlog_iteration_id": "00000000-0000-0000-0000-000000000001",
	})

	// test with default_iteration_macro
	resourceData.Set("default_iteration_macro", "@CurrentIteration")
	patch := expandTeamSettingsPatch(resourceData)

	assert.Equal(t, work.BugsBehavior("asTasks"), *patch.BugsBehavior)
	assert.Equal(t, []string{"monday", "tuesday"}, *patch.WorkingDays)
	expectedBacklogId, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")
	assert.Equal(t, expectedBacklogId, *patch.BacklogIteration)
	assert.Equal(t, "@CurrentIteration", *patch.DefaultIterationMacro)
	assert.Nil(t, patch.DefaultIteration)

	resourceDataWithDefaultIterationId := schema.TestResourceDataRaw(t, ResourceTeamSettings().Schema, map[string]interface{}{
		"default_iteration_id": "00000000-0000-0000-0000-000000000002",
	})
	patchWithDefaultIterationId := expandTeamSettingsPatch(resourceDataWithDefaultIterationId)
	expectedDefaultIterationId, _ := uuid.Parse("00000000-0000-0000-0000-000000000002")
	assert.Equal(t, expectedDefaultIterationId, *patchWithDefaultIterationId.DefaultIteration)
	assert.Nil(t, patchWithDefaultIterationId.DefaultIterationMacro)
}

func TestTeamSettings_CreateOrUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workClient := azdosdkmocks.NewMockWorkClient(ctrl)
	clients := &client.AggregatedClient{WorkClient: workClient, Ctx: context.Background()}

	projectID := uuid.New().String()
	teamID := uuid.New().String()
	backlogIterationID := uuid.New()
	defaultIterationID := uuid.New()

	bugsBehavior := work.BugsBehavior("asTasks")

	resourceData := schema.TestResourceDataRaw(t, ResourceTeamSettings().Schema, map[string]interface{}{
		"project_id":           projectID,
		"team_id":              teamID,
		"bugs_behavior":        string(bugsBehavior),
		"working_days":         []interface{}{"monday", "tuesday"},
		"backlog_iteration_id": backlogIterationID.String(),
		"default_iteration_id": defaultIterationID.String(),
	})

	expectedSettings := &work.TeamSetting{
		BugsBehavior: &bugsBehavior,
		WorkingDays:  &[]string{"monday", "tuesday"},
		BacklogIteration: &work.TeamSettingsIteration{
			Id: &backlogIterationID,
		},
		DefaultIteration: &work.TeamSettingsIteration{
			Id: &defaultIterationID,
		},
	}

	workClient.
		EXPECT().
		UpdateTeamSettings(clients.Ctx, gomock.Any()).
		Return(expectedSettings, nil).
		Times(1)

	workClient.
		EXPECT().
		GetTeamSettings(clients.Ctx, work.GetTeamSettingsArgs{
			Project: &projectID,
			Team:    &teamID,
		}).
		Return(expectedSettings, nil).
		Times(1)

	err := resourceCreateOrUpdateTeamSettings(resourceData, clients)

	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%s/settings", teamID), resourceData.Id())

	assert.Equal(t, string(bugsBehavior), resourceData.Get("bugs_behavior"))
	assert.Equal(t, backlogIterationID.String(), resourceData.Get("backlog_iteration_id"))

	days := resourceData.Get("working_days").(*schema.Set).List()
	assert.Contains(t, days, "monday")
	assert.Contains(t, days, "tuesday")
}

func TestTeamSettings_Read(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workClient := azdosdkmocks.NewMockWorkClient(ctrl)
	clients := &client.AggregatedClient{
		WorkClient: workClient,
		Ctx:        context.Background(),
	}

	projectID := uuid.New().String()
	teamID := uuid.New().String()

	resourceData := schema.TestResourceDataRaw(t, ResourceTeamSettings().Schema, map[string]interface{}{
		"project_id": projectID,
		"team_id":    teamID,
	})
	resourceData.SetId(fmt.Sprintf("%s/settings", teamID))

	bugsBehavior := work.BugsBehavior("asTasks")
	backlogIterationID := uuid.New()
	defaultIterationID := uuid.New()

	expectedSettings := &work.TeamSetting{
		BugsBehavior: &bugsBehavior,
		WorkingDays:  &[]string{"monday", "tuesday"},
		BacklogIteration: &work.TeamSettingsIteration{
			Id: &backlogIterationID,
		},
		DefaultIteration: &work.TeamSettingsIteration{
			Id: &defaultIterationID,
		},
		DefaultIterationMacro: converter.String("@CurrentIteration"),
	}

	getArgs := work.GetTeamSettingsArgs{
		Project: &projectID,
		Team:    &teamID,
	}

	workClient.
		EXPECT().
		GetTeamSettings(clients.Ctx, getArgs).
		Return(expectedSettings, nil).
		Times(1)

	err := resourceReadTeamSettings(resourceData, clients)
	assert.NoError(t, err)

	assert.Equal(t, "asTasks", resourceData.Get("bugs_behavior").(string))
	assert.Equal(t, backlogIterationID.String(), resourceData.Get("backlog_iteration_id").(string))
	assert.Equal(t, defaultIterationID.String(), resourceData.Get("default_iteration_id").(string))
	assert.Equal(t, "@CurrentIteration", resourceData.Get("default_iteration_macro").(string))

	days := resourceData.Get("working_days").(*schema.Set)
	assert.True(t, days.Contains("monday"))
	assert.True(t, days.Contains("tuesday"))
}

func TestTeamSettings_Read_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workClient := azdosdkmocks.NewMockWorkClient(ctrl)
	clients := &client.AggregatedClient{
		WorkClient: workClient,
		Ctx:        context.Background(),
	}

	projectID := uuid.New().String()
	teamID := uuid.New().String()

	resourceData := schema.TestResourceDataRaw(t, ResourceTeamSettings().Schema, map[string]interface{}{
		"project_id": projectID,
		"team_id":    teamID,
	})
	resourceData.SetId(fmt.Sprintf("%s/settings", teamID))

	notFoundCode := 404
	notFoundErr := azuredevops.WrappedError{
		StatusCode: &notFoundCode,
	}

	getArgs := work.GetTeamSettingsArgs{
		Project: &projectID,
		Team:    &teamID,
	}

	workClient.
		EXPECT().
		GetTeamSettings(clients.Ctx, getArgs).
		Return(nil, notFoundErr).
		Times(1)

	err := resourceReadTeamSettings(resourceData, clients)
	assert.NoError(t, err)
	assert.Equal(t, "", resourceData.Id())
}

func TestTeamSettings_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workClient := azdosdkmocks.NewMockWorkClient(ctrl)
	clients := &client.AggregatedClient{
		WorkClient: workClient,
		Ctx:        context.Background(),
	}

	projectID := uuid.New().String()
	teamID := uuid.New().String()

	resourceData := schema.TestResourceDataRaw(t, ResourceTeamSettings().Schema, map[string]interface{}{
		"project_id": projectID,
		"team_id":    teamID,
	})
	resourceData.SetId(fmt.Sprintf("%s/settings", teamID))

	workClient.
		EXPECT().
		UpdateTeamSettings(clients.Ctx, gomock.Any()).
		Return(&work.TeamSetting{}, nil).
		Times(1)

	err := resourceDeleteTeamSettings(resourceData, clients)
	assert.NoError(t, err)
	assert.Equal(t, "", resourceData.Id())
}

func TestUpdateTeamSettingsInternal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workClient := azdosdkmocks.NewMockWorkClient(ctrl)
	clients := &client.AggregatedClient{
		WorkClient: workClient,
		Ctx:        context.Background(),
	}

	projectID := "pID"
	teamID := "tID"
	patch := work.TeamSettingsPatch{}

	workClient.
		EXPECT().
		UpdateTeamSettings(clients.Ctx, gomock.Any()).
		Return(&work.TeamSetting{}, nil).
		Times(1)

	err := updateTeamSettingsInternal(clients, projectID, teamID, patch)
	assert.NoError(t, err)

	expectedErr := fmt.Errorf("boom")

	workClient.
		EXPECT().
		UpdateTeamSettings(clients.Ctx, gomock.Any()).
		Return(nil, expectedErr).
		Times(1)

	err = updateTeamSettingsInternal(clients, projectID, teamID, patch)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedErr.Error())
}
