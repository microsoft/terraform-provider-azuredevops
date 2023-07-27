//go:build all || utils || workitemtracking
// +build all utils workitemtracking

package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/testhelper"
	"github.com/stretchr/testify/require"
)

var classificationProjectName = "test-acc-project-0fu72ecbiu"
var classificationProjectID = "9c3a5552-268c-423c-a9cd-7de0b36b7035"

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
		Return(nil, fmt.Errorf(errMsg)).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, CreateClassificationNodeSchema(map[string]*schema.Schema{}), nil)
	resourceData.Set("project_id", classificationProjectID)

	err := ReadClassificationNode(clients, resourceData, workitemtracking.TreeStructureGroupValues.Areas)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), errMsg)
}

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
	require.IsType(t, "", v)
	require.Equal(t, classificationProjectID, v)

	v = resourceData.Get("path")
	require.IsType(t, node.path, v)
	require.Equal(t, node.path, v)

	v = resourceData.Get("has_children")
	require.IsType(t, true, v)
	require.Equal(t, len(node.children) > 0, v)

	v = resourceData.Get("children")
	require.NotNil(t, v)
	require.Len(t, v, len(node.children))
}

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
	id := resourceData.Id()
	require.NotEmpty(t, id)

	var v interface{}

	v = resourceData.Get("project_id")
	require.IsType(t, "", v)
	require.Equal(t, classificationProjectID, v)

	v = resourceData.Get("path")
	require.IsType(t, node.path, v)
	require.Equal(t, node.path, v)

	v = resourceData.Get("has_children")
	require.IsType(t, true, v)
	require.Equal(t, len(node.children) > 0, v)

	v = resourceData.Get("children")
	require.NotNil(t, v)
	require.Len(t, v, len(node.children))

	for i, childNode := range node.children {
		tfNodePrefix := fmt.Sprintf("children.%d.", i)
		v = resourceData.Get(tfNodePrefix + "project_id")
		require.IsType(t, "", v)
		require.Equal(t, classificationProjectID, v)

		v = resourceData.Get(tfNodePrefix + "path")
		require.IsType(t, childNode.path, v)
		require.Equal(t, childNode.path, v)

		v = resourceData.Get(tfNodePrefix + "has_children")
		require.IsType(t, true, v)
		require.Equal(t, len(childNode.children) > 0, v)

		v = resourceData.Get(tfNodePrefix + "children")
		if len(childNode.children) > 0 {
			require.NotNil(t, v)
			require.Len(t, v, len(childNode.children))
		} else {
			require.Nil(t, v)
		}
	}
}
