//go:build (all || security || data_sources || data_security_namespace_token) && (!exclude_data_sources || !exclude_security || !exclude_data_security_namespace_token)
// +build all security data_sources data_security_namespace_token
// +build !exclude_data_sources !exclude_security !exclude_data_security_namespace_token

package security

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/security"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/testhelper"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestDataSecurityNamespaceToken_GitRepositories_ProjectOnly tests token generation for Git Repositories namespace with only project_id
func TestDataSecurityNamespaceToken_GitRepositories_ProjectOnly(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	projectID := testhelper.CreateUUID()

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.GitRepositories).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id": projectID.String(),
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("repoV2/%s", projectID.String()), resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_GitRepositories_WithRepository tests token generation with repository_id
func TestDataSecurityNamespaceToken_GitRepositories_WithRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	projectID := testhelper.CreateUUID()
	repoID := testhelper.CreateUUID()

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.GitRepositories).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id":    projectID.String(),
		"repository_id": repoID.String(),
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("repoV2/%s/%s", projectID.String(), repoID.String()), resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_GitRepositories_WithRefName tests token generation with ref_name
func TestDataSecurityNamespaceToken_GitRepositories_WithRefName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	projectID := testhelper.CreateUUID()
	repoID := testhelper.CreateUUID()
	refName := "refs/heads/main"

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.GitRepositories).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id":    projectID.String(),
		"repository_id": repoID.String(),
		"ref_name":      refName,
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)

	// Verify token starts correctly (exact encoding depends on EncodeUtf16HexString)
	token := resourceData.Get("token").(string)
	require.Contains(t, token, fmt.Sprintf("repoV2/%s/%s/refs/heads/", projectID.String(), repoID.String()))
}

// TestDataSecurityNamespaceToken_GitRepositories_RefNameWithoutRepo tests error when ref_name is provided without repository_id
func TestDataSecurityNamespaceToken_GitRepositories_RefNameWithoutRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	projectID := testhelper.CreateUUID()

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.GitRepositories).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id": projectID.String(),
		"ref_name":   "refs/heads/main",
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "ref_name provided without repository_id")
}

// TestDataSecurityNamespaceToken_Project tests token generation for Project namespace
func TestDataSecurityNamespaceToken_Project(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	projectID := testhelper.CreateUUID()

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.Project).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id": projectID.String(),
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("$PROJECT:vstfs:///Classification/TeamProject/%s", projectID.String()), resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_Build_ProjectOnly tests token generation for Build namespace with only project_id
func TestDataSecurityNamespaceToken_Build_ProjectOnly(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	projectID := testhelper.CreateUUID()

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.Build).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id": projectID.String(),
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, projectID.String(), resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_Build_WithPath tests token generation for Build namespace with path
func TestDataSecurityNamespaceToken_Build_WithPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	projectID := testhelper.CreateUUID()

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.Build).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id": projectID.String(),
		"path":       "\\MyFolder\\SubFolder",
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("%s/MyFolder/SubFolder", projectID.String()), resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_Build_WithDefinitionID tests token generation for Build namespace with definition_id
func TestDataSecurityNamespaceToken_Build_WithDefinitionID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	buildClient := azdosdkmocks.NewMockBuildClient(ctrl)
	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		BuildClient:    buildClient,
		Ctx:            context.Background(),
	}

	projectID := testhelper.CreateUUID()
	definitionID := 123

	// Mock the GetDefinition call to return a build definition with a path
	buildClient.
		EXPECT().
		GetDefinition(clients.Ctx, gomock.Any()).
		Return(&build.BuildDefinition{
			Id:   converter.Int(definitionID),
			Path: converter.String("\\MyFolder"),
			Name: converter.String("Test Pipeline"),
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.Build).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id":    projectID.String(),
		"definition_id": "123",
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("%s/MyFolder/123", projectID.String()), resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_Build_WithDefinitionIDRootPath tests token generation when definition is at root
func TestDataSecurityNamespaceToken_Build_WithDefinitionIDRootPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	buildClient := azdosdkmocks.NewMockBuildClient(ctrl)
	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		BuildClient:    buildClient,
		Ctx:            context.Background(),
	}

	projectID := testhelper.CreateUUID()
	definitionID := 456

	// Mock the GetDefinition call to return a build definition at root path
	buildClient.
		EXPECT().
		GetDefinition(clients.Ctx, gomock.Any()).
		Return(&build.BuildDefinition{
			Id:   converter.Int(definitionID),
			Path: converter.String("\\"),
			Name: converter.String("Root Pipeline"),
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.Build).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id":    projectID.String(),
		"definition_id": "456",
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("%s/456", projectID.String()), resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_CSS tests token generation for CSS (Areas) namespace
func TestDataSecurityNamespaceToken_CSS(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	witClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)

	clients := &client.AggregatedClient{
		SecurityClient:         securityClient,
		WorkItemTrackingClient: witClient,
		Ctx:                    context.Background(),
	}

	projectID := testhelper.CreateUUID()
	rootNodeID := testhelper.CreateUUID()

	// Mock the root classification node call
	witClient.
		EXPECT().
		GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
			Project:        converter.String(projectID.String()),
			StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
			Depth:          converter.Int(1),
		}).
		Return(&workitemtracking.WorkItemClassificationNode{
			Identifier: rootNodeID,
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.CSS).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id": projectID.String(),
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("vstfs:///Classification/Node/%s", rootNodeID.String()), resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_CSS_WithPath tests token generation for CSS namespace with path
func TestDataSecurityNamespaceToken_CSS_WithPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	witClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)

	clients := &client.AggregatedClient{
		SecurityClient:         securityClient,
		WorkItemTrackingClient: witClient,
		Ctx:                    context.Background(),
	}

	projectID := testhelper.CreateUUID()
	rootNodeID := testhelper.CreateUUID()
	childNodeID := testhelper.CreateUUID()

	// Mock the root classification node call
	witClient.
		EXPECT().
		GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
			Project:        converter.String(projectID.String()),
			StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
			Depth:          converter.Int(1),
		}).
		Return(&workitemtracking.WorkItemClassificationNode{
			Identifier:  rootNodeID,
			HasChildren: converter.Bool(true),
		}, nil).
		Times(1)

	// Mock the child node call
	witClient.
		EXPECT().
		GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
			Project:        converter.String(projectID.String()),
			Path:           converter.String("ChildArea"),
			StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
			Depth:          converter.Int(1),
		}).
		Return(&workitemtracking.WorkItemClassificationNode{
			Identifier: childNodeID,
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.CSS).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id": projectID.String(),
		"path":       "/ChildArea",
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	expectedToken := fmt.Sprintf("vstfs:///Classification/Node/%s:vstfs:///Classification/Node/%s", rootNodeID.String(), childNodeID.String())
	require.Equal(t, expectedToken, resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_Iteration tests token generation for Iteration namespace
func TestDataSecurityNamespaceToken_Iteration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	witClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)

	clients := &client.AggregatedClient{
		SecurityClient:         securityClient,
		WorkItemTrackingClient: witClient,
		Ctx:                    context.Background(),
	}

	projectID := testhelper.CreateUUID()
	rootNodeID := testhelper.CreateUUID()

	// Mock the root classification node call
	witClient.
		EXPECT().
		GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
			Project:        converter.String(projectID.String()),
			StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
			Depth:          converter.Int(1),
		}).
		Return(&workitemtracking.WorkItemClassificationNode{
			Identifier: rootNodeID,
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.Iteration).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id": projectID.String(),
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("vstfs:///Classification/Node/%s", rootNodeID.String()), resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_Tagging_WithProject tests token generation for Tagging namespace with project
func TestDataSecurityNamespaceToken_Tagging_WithProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	projectID := testhelper.CreateUUID()

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.Tagging).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id": projectID.String(),
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("/%s", projectID.String()), resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_Tagging_NoProject tests token generation for Tagging namespace without project
func TestDataSecurityNamespaceToken_Tagging_NoProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.Tagging).String())

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "", resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_ServiceHooks tests token generation for ServiceHooks namespace
func TestDataSecurityNamespaceToken_ServiceHooks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	projectID := testhelper.CreateUUID()

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.ServiceHooks).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id": projectID.String(),
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("PublisherSecurity/%s", projectID.String()), resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_WorkItemQueryFolders tests token generation for WorkItemQueryFolders namespace
func TestDataSecurityNamespaceToken_WorkItemQueryFolders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	witClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)

	clients := &client.AggregatedClient{
		SecurityClient:         securityClient,
		WorkItemTrackingClient: witClient,
		Ctx:                    context.Background(),
	}

	projectID := testhelper.CreateUUID()
	sharedQueriesID := testhelper.CreateUUID()
	folderID := testhelper.CreateUUID()

	// Mock the Shared Queries folder call
	witClient.
		EXPECT().
		GetQuery(clients.Ctx, workitemtracking.GetQueryArgs{
			Project: converter.String(projectID.String()),
			Query:   converter.String("Shared Queries"),
			Depth:   converter.Int(1),
		}).
		Return(&workitemtracking.QueryHierarchyItem{
			Id:   sharedQueriesID,
			Name: converter.String("Shared Queries"),
			Children: &[]workitemtracking.QueryHierarchyItem{
				{
					Id:   folderID,
					Name: converter.String("MyFolder"),
				},
			},
		}, nil).
		Times(1)

	// Mock the folder call
	witClient.
		EXPECT().
		GetQuery(clients.Ctx, workitemtracking.GetQueryArgs{
			Project: converter.String(projectID.String()),
			Query:   converter.String(folderID.String()),
			Depth:   converter.Int(1),
		}).
		Return(&workitemtracking.QueryHierarchyItem{
			Id:   folderID,
			Name: converter.String("MyFolder"),
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.WorkItemQueryFolders).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id": projectID.String(),
		"path":       "MyFolder",
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	expectedToken := fmt.Sprintf("$/%s/%s/%s", projectID.String(), sharedQueriesID.String(), folderID.String())
	require.Equal(t, expectedToken, resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_Analytics tests token generation for Analytics namespace
func TestDataSecurityNamespaceToken_Analytics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	projectID := testhelper.CreateUUID()

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.Analytics).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id": projectID.String(),
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("$/%s", projectID.String()), resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_AnalyticsViews tests token generation for AnalyticsViews namespace
func TestDataSecurityNamespaceToken_AnalyticsViews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	projectID := testhelper.CreateUUID()

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.AnalyticsViews).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id": projectID.String(),
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("$/Shared/%s", projectID.String()), resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_Collection tests token generation for Collection namespace
func TestDataSecurityNamespaceToken_Collection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.Collection).String())

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "NAMESPACE:", resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_Process tests token generation for Process namespace
func TestDataSecurityNamespaceToken_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.Process).String())

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "$PROCESS:", resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_Process_WithTemplate tests token generation for Process namespace with template
func TestDataSecurityNamespaceToken_Process_WithTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	templateID := testhelper.CreateUUID()

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.Process).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"workitem_template_id": templateID.String(),
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("$PROCESS:%s:", templateID.String()), resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_Process_WithProcessAndTemplate tests token generation for Process namespace with both identifiers
func TestDataSecurityNamespaceToken_Process_WithProcessAndTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	processID := testhelper.CreateUUID()
	templateID := testhelper.CreateUUID()

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.Process).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"process_id":           processID.String(),
		"workitem_template_id": templateID.String(),
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("$PROCESS:%s:%s:", processID.String(), templateID.String()), resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_Process_ProcessWithoutTemplate tests error when process_id without workitem_template_id
func TestDataSecurityNamespaceToken_Process_ProcessWithoutTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	processID := testhelper.CreateUUID()

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.Process).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"process_id": processID.String(),
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "process_id provided without workitem_template_id")
}

// TestDataSecurityNamespaceToken_AuditLog tests token generation for AuditLog namespace
func TestDataSecurityNamespaceToken_AuditLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.AuditLog).String())

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "AllPermissions", resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_BuildAdministration tests token generation for BuildAdministration namespace
func TestDataSecurityNamespaceToken_BuildAdministration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.BuildAdministration).String())

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "BuildPrivileges", resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_Server tests token generation for Server namespace
func TestDataSecurityNamespaceToken_Server(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.Server).String())

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "FrameworkGlobalSecurity", resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_VersionControlPrivileges tests token generation for VersionControlPrivileges namespace
func TestDataSecurityNamespaceToken_VersionControlPrivileges(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.VersionControlPrivileges).String())

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "Global", resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_ServiceEndpoints tests error for ServiceEndpoints namespace
func TestDataSecurityNamespaceToken_ServiceEndpoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.ServiceEndpoints).String())

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "service Endpoints namespace uses role assignments")
}

// TestDataSecurityNamespaceToken_ByName tests token generation using namespace_name instead of namespace_id
func TestDataSecurityNamespaceToken_ByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)

	clients := &client.AggregatedClient{
		SecurityClient: securityClient,
		Ctx:            context.Background(),
	}

	projectID := testhelper.CreateUUID()
	gitRepoNamespaceID := uuid.UUID(utils.SecurityNamespaceIDValues.GitRepositories)

	// Mock the QuerySecurityNamespaces call
	securityClient.
		EXPECT().
		QuerySecurityNamespaces(clients.Ctx, security.QuerySecurityNamespacesArgs{}).
		Return(&[]security.SecurityNamespaceDescription{
			{
				NamespaceId: &gitRepoNamespaceID,
				Name:        converter.String("Git Repositories"),
			},
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_name", "Git Repositories")
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id": projectID.String(),
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, fmt.Sprintf("repoV2/%s", projectID.String()), resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_ByName_NotFound tests error when namespace name is not found
func TestDataSecurityNamespaceToken_ByName_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)

	clients := &client.AggregatedClient{
		SecurityClient: securityClient,
		Ctx:            context.Background(),
	}

	// Mock the QuerySecurityNamespaces call
	securityClient.
		EXPECT().
		QuerySecurityNamespaces(clients.Ctx, security.QuerySecurityNamespacesArgs{}).
		Return(&[]security.SecurityNamespaceDescription{}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_name", "NonExistent Namespace")

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "namespace with name 'NonExistent Namespace' not found")
}

// TestDataSecurityNamespaceToken_MissingRequiredIdentifiers tests error when required identifiers are missing
func TestDataSecurityNamespaceToken_MissingRequiredIdentifiers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.GitRepositories).String())
	// Missing project_id

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "missing required identifiers")
	require.Contains(t, err.Error(), "project_id")
}

// TestDataSecurityNamespaceToken_UnsupportedNamespace tests error for unsupported namespace
func TestDataSecurityNamespaceToken_UnsupportedNamespace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	// Use a namespace that's not in the templates
	unsupportedNamespaceID := uuid.UUID(utils.SecurityNamespaceIDValues.Identity)

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", unsupportedNamespaceID.String())

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "unable to generate token for namespace")
}

// TestDataSecurityNamespaceToken_ReturnIdentifierInfo tests returning identifier information instead of generating a token
func TestDataSecurityNamespaceToken_ReturnIdentifierInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.GitRepositories).String())
	resourceData.Set("return_identifier_info", true)

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.Nil(t, err)

	// Check that required_identifiers and optional_identifiers are populated
	requiredIdentifiers := resourceData.Get("required_identifiers").([]interface{})
	require.Len(t, requiredIdentifiers, 1)
	require.Equal(t, "project_id", requiredIdentifiers[0])

	optionalIdentifiers := resourceData.Get("optional_identifiers").([]interface{})
	require.Len(t, optionalIdentifiers, 2)
	require.Contains(t, optionalIdentifiers, "repository_id")
	require.Contains(t, optionalIdentifiers, "ref_name")

	// Token should not be set
	require.Empty(t, resourceData.Get("token"))
}

// TestDataSecurityNamespaceToken_ReturnIdentifierInfo_UnsupportedNamespace tests error when requesting identifier info for unsupported namespace
func TestDataSecurityNamespaceToken_ReturnIdentifierInfo_UnsupportedNamespace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		Ctx:            context.Background(),
	}

	// Use a namespace that's not in the templates
	unsupportedNamespaceID := uuid.UUID(utils.SecurityNamespaceIDValues.Identity)

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", unsupportedNamespaceID.String())
	resourceData.Set("return_identifier_info", true)

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "no template information available for namespace")
}

// TestDataSecurityNamespaceToken_QueryError tests error handling when API query fails
func TestDataSecurityNamespaceToken_QueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)

	clients := &client.AggregatedClient{
		SecurityClient: securityClient,
		Ctx:            context.Background(),
	}

	// Mock the QuerySecurityNamespaces call to return an error
	securityClient.
		EXPECT().
		QuerySecurityNamespaces(clients.Ctx, security.QuerySecurityNamespacesArgs{}).
		Return(nil, fmt.Errorf("API error")).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_name", "Git Repositories")

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "querying security namespaces")
}

// TestDataSecurityNamespaceToken_ClassificationNode_APIError tests error handling when classification node API fails
func TestDataSecurityNamespaceToken_ClassificationNode_APIError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	witClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)

	clients := &client.AggregatedClient{
		SecurityClient:         securityClient,
		WorkItemTrackingClient: witClient,
		Ctx:                    context.Background(),
	}

	projectID := testhelper.CreateUUID()

	// Mock the root classification node call to fail
	witClient.
		EXPECT().
		GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
			Project:        converter.String(projectID.String()),
			StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
			Depth:          converter.Int(1),
		}).
		Return(nil, fmt.Errorf("API error")).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, nil)
	resourceData.Set("namespace_id", uuid.UUID(utils.SecurityNamespaceIDValues.CSS).String())
	resourceData.Set("identifiers", map[string]interface{}{
		"project_id": projectID.String(),
	})

	err := dataSecurityNamespaceTokenRead(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "error getting root classification node")
}
