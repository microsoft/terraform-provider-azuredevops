//go:build all || utils || workitemtracking
// +build all utils workitemtracking

package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/testhelper"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	classificationProjectName = "test-acc-project-0fu72ecbiu"
	classificationProjectID   = "9c3a5552-268c-423c-a9cd-7de0b36b7035"
)

type classificationNodeDefinition struct {
	id         string
	name       string
	pathNative string
	path       string
	children   []*classificationNodeDefinition
}

func newClassificationTestNode(structureType workitemtracking.TreeStructureGroup, parent *classificationNodeDefinition) *classificationNodeDefinition {
	nodeName := "test-acc-node-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	node := classificationNodeDefinition{
		id:         testhelper.CreateUUID().String(),
		name:       nodeName,
		pathNative: "\\" + classificationProjectName + "\\" + string(structureType) + "\\" + nodeName,
		path:       "/" + nodeName,
	}
	if parent != nil {
		node.pathNative = parent.pathNative + "\\" + nodeName
		node.path = parent.path + "/" + nodeName
		parent.children = append(parent.children, &node)
	}
	return &node
}

func newClassificationTestNodes(structureType workitemtracking.TreeStructureGroup, parent *classificationNodeDefinition, size int) *[]*classificationNodeDefinition {
	ret := make([]*classificationNodeDefinition, size)
	for i := 0; i < size; i++ {
		ret[i] = newClassificationTestNode(structureType, parent)
	}
	return &ret
}

func convertClassificationTestNode(testNode *classificationNodeDefinition) *workitemtracking.WorkItemClassificationNode {
	node := workitemtracking.WorkItemClassificationNode{
		Identifier: converter.UUID(testNode.id),
		Name:       converter.String(testNode.name),
		Path:       converter.String(testNode.pathNative),
	}
	if len(testNode.children) > 0 {
		node.HasChildren = converter.Bool(true)
		children := make([]workitemtracking.WorkItemClassificationNode, len(testNode.children))
		for i, v := range testNode.children {
			children[i] = *convertClassificationTestNode(v)
		}
		node.Children = &children
	} else {
		node.HasChildren = converter.Bool(false)
	}
	return &node
}

// --- TESTS ---

// TestClassification_Read_DontSwallowError verifieert dat we echte errors (geen 404) teruggeven
func TestClassification_Read_DontSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	witClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{
		WorkItemTrackingClient: witClient,
		Ctx:                    context.Background(),
	}

	structureType := workitemtracking.TreeStructureGroupValues.Areas
	errMsg := "@@GetClassificationNode@@failed@@"
	witClient.EXPECT().
		GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
			Project:        converter.String(classificationProjectID),
			StructureGroup: &structureType,
			Depth:          converter.Int(1),
		}).
		Return(nil, fmt.Errorf("%s", errMsg)).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, CreateClassificationNodeSchema(map[string]*schema.Schema{}), nil)
	resourceData.Set("project_id", classificationProjectID)

	err := ReadClassificationNode(clients, resourceData, workitemtracking.TreeStructureGroupValues.Areas)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), errMsg)
}

// TestClassification_Read basis test
func TestClassification_Read(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	witClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{
		WorkItemTrackingClient: witClient,
		Ctx:                    context.Background(),
	}

	structureType := workitemtracking.TreeStructureGroupValues.Areas
	node := newClassificationTestNode(structureType, nil)

	witClient.EXPECT().
		GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
			Project:        converter.String(classificationProjectID),
			StructureGroup: &structureType,
			Depth:          converter.Int(1),
		}).
		Return(convertClassificationTestNode(node), nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, CreateClassificationNodeSchema(map[string]*schema.Schema{}), nil)
	resourceData.Set("project_id", classificationProjectID)

	err := ReadClassificationNode(clients, resourceData, structureType)
	require.Nil(t, err)
	id := resourceData.Id()
	require.NotEmpty(t, id)

	var v interface{}
	v = resourceData.Get("project_id")
	require.Equal(t, classificationProjectID, v)

	v = resourceData.Get("path")
	require.Equal(t, node.path, v)
}

// TestClassification_Read_Children test met kinderen
func TestClassification_Read_Children(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	witClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{
		WorkItemTrackingClient: witClient,
		Ctx:                    context.Background(),
	}

	structureType := workitemtracking.TreeStructureGroupValues.Areas
	node := newClassificationTestNode(structureType, nil)
	_ = newClassificationTestNodes(structureType, node, 3)

	witClient.EXPECT().
		GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
			Project:        converter.String(classificationProjectID),
			StructureGroup: &structureType,
			Depth:          converter.Int(1),
		}).
		Return(convertClassificationTestNode(node), nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, CreateClassificationNodeSchema(map[string]*schema.Schema{}), nil)
	resourceData.Set("project_id", classificationProjectID)

	err := ReadClassificationNode(clients, resourceData, structureType)
	require.Nil(t, err)

	v := resourceData.Get("children")
	require.NotNil(t, v)
	require.Len(t, v, len(node.children))
}

// TestClassification_Create verifieert correcte aanroep van CreateOrUpdate
func TestClassification_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	witClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{
		WorkItemTrackingClient: witClient,
		Ctx:                    context.Background(),
	}

	structureType := workitemtracking.TreeStructureGroupValues.Areas
	nodeName := "NewArea"
	projectID := classificationProjectID
	path := "\\Project\\Area"

	// Create arguments matching logic
	expectedArgs := workitemtracking.CreateOrUpdateClassificationNodeArgs{
		Project:        converter.String(projectID),
		StructureGroup: &structureType,
		Path:           converter.String(path),
		PostedNode: &workitemtracking.WorkItemClassificationNode{
			Name: converter.String(nodeName),
		},
	}

	mockReturnNode := workitemtracking.WorkItemClassificationNode{
		Id:          converter.Int(1001),
		Identifier:  converter.UUID("55555555-5555-5555-5555-555555555555"),
		HasChildren: converter.Bool(false),
	}

	witClient.EXPECT().
		CreateOrUpdateClassificationNode(clients.Ctx, expectedArgs).
		Return(&mockReturnNode, nil).
		Times(1)

	// Gebruik CreateClassificationNodeResourceSchema hier!
	resourceData := schema.TestResourceDataRaw(t, CreateClassificationNodeResourceSchema(workitemtracking.TreeStructureGroupValues.Areas), nil)
	resourceData.Set("project_id", projectID)
	resourceData.Set("name", nodeName)
	resourceData.Set("path", path)

	err := CreateOrUpdateClassificationNode(clients, resourceData, structureType)
	require.Nil(t, err)
	require.Equal(t, "55555555-5555-5555-5555-555555555555", resourceData.Id())
	require.Equal(t, 1001, resourceData.Get("node_id"))
}

// TestClassification_Create_WithAttributes verifieert start/einddatum logica
func TestClassification_Create_WithAttributes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	witClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{
		WorkItemTrackingClient: witClient,
		Ctx:                    context.Background(),
	}

	structureType := workitemtracking.TreeStructureGroupValues.Iterations
	nodeName := "Sprint 1"
	startDate := "2023-01-01"
	finishDate := "2023-01-14"

	mockReturnNode := workitemtracking.WorkItemClassificationNode{
		Id:          converter.Int(2002),
		Identifier:  converter.UUID("66666666-6666-6666-6666-666666666666"),
		HasChildren: converter.Bool(false),
	}

	witClient.EXPECT().
		CreateOrUpdateClassificationNode(clients.Ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, args workitemtracking.CreateOrUpdateClassificationNodeArgs) (*workitemtracking.WorkItemClassificationNode, error) {
			require.Equal(t, nodeName, *args.PostedNode.Name)
			require.NotNil(t, args.PostedNode.Attributes)

			attrs := *args.PostedNode.Attributes
			require.Equal(t, startDate, attrs["startDate"])
			require.Equal(t, finishDate, attrs["finishDate"])

			return &mockReturnNode, nil
		}).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, CreateClassificationNodeResourceSchema(workitemtracking.TreeStructureGroupValues.Iterations), nil)
	resourceData.Set("project_id", classificationProjectID)
	resourceData.Set("name", nodeName)
	resourceData.Set("path", "\\Project\\Iteration")

	// Attributes instellen
	resourceData.Set("attributes", []interface{}{
		map[string]interface{}{
			"start_date":  startDate,
			"finish_date": finishDate,
		},
	})

	err := CreateOrUpdateClassificationNode(clients, resourceData, structureType)
	require.Nil(t, err)
}

// // TestClassification_Read_NotFound checkt of resource uit state wordt gehaald bij 404
func TestClassification_Read_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	witClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{
		WorkItemTrackingClient: witClient,
		Ctx:                    context.Background(),
	}

	structureType := workitemtracking.TreeStructureGroupValues.Areas
	errorMsg := "VS402485: The Area/Iteration name is not recognized."

	witClient.EXPECT().
		GetClassificationNode(clients.Ctx, gomock.Any()).
		Return(nil, fmt.Errorf("%s", errorMsg)).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, CreateClassificationNodeSchema(map[string]*schema.Schema{}), nil)
	resourceData.SetId("some-existing-id")
	resourceData.Set("project_id", classificationProjectID)

	err := ReadClassificationNode(clients, resourceData, structureType)

	require.NotNil(t, err)
	require.Equal(t, "", resourceData.Id())
}

// TestClassification_Delete test de delete flow
func TestClassification_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	witClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{
		WorkItemTrackingClient: witClient,
		Ctx:                    context.Background(),
	}

	structureType := workitemtracking.TreeStructureGroupValues.Areas
	pathToDelete := "\\Project\\Area\\ToDelete"
	rootUUID := testhelper.CreateUUID()
	rootID := 999

	// 1. Verwacht GetClassificationNode (Root ophalen)
	witClient.EXPECT().
		GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
			Project:        converter.String(classificationProjectID),
			StructureGroup: &structureType,
			Depth:          converter.Int(1),
			Path:           converter.String(""),
		}).
		Return(&workitemtracking.WorkItemClassificationNode{
			Id:         converter.Int(rootID),
			Identifier: rootUUID,
		}, nil).
		Times(1)

	// 2. Verwacht DeleteClassificationNode
	witClient.EXPECT().
		DeleteClassificationNode(clients.Ctx, workitemtracking.DeleteClassificationNodeArgs{
			Project:        converter.String(classificationProjectID),
			StructureGroup: &structureType,
			Path:           converter.String(pathToDelete),
			ReclassifyId:   converter.Int(rootID),
		}).
		Return(nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, CreateClassificationNodeResourceSchema(workitemtracking.TreeStructureGroupValues.Areas), nil)
	resourceData.Set("project_id", classificationProjectID)
	resourceData.Set("path", pathToDelete)
	resourceData.Set("name", "ToDelete")

	err := DeleteClassificationNode(clients, resourceData, structureType)
	require.Nil(t, err)
}
