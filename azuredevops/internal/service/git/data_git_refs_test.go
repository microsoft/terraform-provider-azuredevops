package git

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGitRefsDataSource_Read_DontSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &client.AggregatedClient{
		GitReposClient: repoClient,
		Ctx:            context.Background(),
	}

	repoClient.
		EXPECT().
		GetRefs(clients.Ctx, gomock.Any()).
		Return(nil, fmt.Errorf("GetRefs error")).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataGitRefs().Schema, nil)
	resourceData.Set("repository_id", "00000000-0000-0000-0000-000000000000")

	err := dataSourceGitRefsRead(clients.Ctx, resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err[0].Summary, "GetRefs error")
}

func TestGitRefsDataSource_Read_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoClient := azdosdkmocks.NewMockGitClient(ctrl)
	clients := &client.AggregatedClient{
		GitReposClient: repoClient,
		Ctx:            context.Background(),
	}

	repoClient.
		EXPECT().
		GetRefs(clients.Ctx, gomock.Any()).
		Return(&git.GetRefsResponseValue{
			Value: []git.GitRef{
				{
					Name:     converter.String("refs/heads/main"),
					ObjectId: converter.String("abcd1234efgh5678"),
				},
				{
					Name:     converter.String("refs/heads/dev"),
					ObjectId: converter.String("1234abcd5678efgh"),
				},
			},
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataGitRefs().Schema, nil)
	resourceData.Set("repository_id", "00000000-0000-0000-0000-000000000000")

	err := dataSourceGitRefsRead(clients.Ctx, resourceData, clients)
	require.Nil(t, err)

	refs := resourceData.Get("refs").([]interface{})
	require.Equal(t, 2, len(refs))

	ref1 := refs[0].(map[string]interface{})
	require.Equal(t, "refs/heads/main", ref1["name"])
	require.Equal(t, "abcd1234efgh5678", ref1["object_id"])
}
