package permissions

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"testing"

	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"

	"github.com/golang/mock/gomock"
)

func init() {
	/* add code for test setup here */
}

/**
 * Begin unit tests
 */

func TestAzureDevOpsGitPermissions_Create_Test(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &client.AggregatedClient{
		CoreClient: coreClient,
		Ctx:        context.Background(),
	}

	/* start writing test here */
}
