package tfhelper

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	testID      = uuid.New()
	testProject = core.TeamProject{
		Id:          &testID,
		Name:        converter.String("Name"),
		Visibility:  &core.ProjectVisibilityValues.Public,
		Description: converter.String("Description"),
		Capabilities: &map[string]map[string]string{
			"versioncontrol":  {"sourceControlType": "SouceControlType"},
			"processTemplate": {"templateTypeId": testID.String()},
		},
	}
)

type testCase struct {
	Name            string
	projectNameOrID string
	exceptProjectID string
	exceptError     bool
	MockedFunction  func(*azdosdkmocks.MockCoreClientMockRecorder, *client.AggregatedClient, string) *gomock.Call
}

func TestGetRealProjectId(t *testing.T) {
	cases := []testCase{
		{
			Name:            "ProjectName",
			projectNameOrID: "projectName",
			exceptProjectID: testProject.Id.String(),
			exceptError:     false,
			MockedFunction: func(mr *azdosdkmocks.MockCoreClientMockRecorder, clients *client.AggregatedClient, projectNameOrID string) *gomock.Call {
				return mr.GetProject(clients.Ctx, core.GetProjectArgs{
					ProjectId:           &projectNameOrID,
					IncludeCapabilities: converter.Bool(true),
					IncludeHistory:      converter.Bool(false),
				}).Return(&testProject, nil)
			},
		},
		{
			Name:            "ProjectNotExist",
			projectNameOrID: "e0f55995-e2f2",
			exceptProjectID: "",
			exceptError:     true,
			MockedFunction: func(mr *azdosdkmocks.MockCoreClientMockRecorder, clients *client.AggregatedClient, projectNameOrID string) *gomock.Call {
				return mr.GetProject(clients.Ctx, core.GetProjectArgs{
					ProjectId:           &projectNameOrID,
					IncludeCapabilities: converter.Bool(true),
					IncludeHistory:      converter.Bool(false),
				}).Return(nil, fmt.Errorf("Project not found, projectNameOrID: %s ", projectNameOrID))
			},
		},
		{
			Name:            "Legal ProjectID",
			projectNameOrID: "e0f55995-e2f2-4268-9b4f-26295f31a8ad",
			exceptProjectID: "e0f55995-e2f2-4268-9b4f-26295f31a8ad",
			exceptError:     false,
			MockedFunction: func(mr *azdosdkmocks.MockCoreClientMockRecorder, clients *client.AggregatedClient, projectNameOrID string) *gomock.Call {
				return nil
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &client.AggregatedClient{CoreClient: coreClient, Ctx: context.Background()}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing %q..", tc.Name)

		if tc.MockedFunction != nil {
			tc.MockedFunction(coreClient.EXPECT(), clients, tc.projectNameOrID)
		}
		projectID, err := GetRealProjectId(tc.projectNameOrID, clients)
		if tc.exceptError {
			require.NotNil(t, err)
		}
		require.Equal(t, tc.exceptProjectID, projectID)
	}
}
