package workitemtracking

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// getAreaTreeResourceData builds *schema.ResourceData for the area tree
// resource from the given raw attribute map.
func getAreaTreeResourceData(t *testing.T, input map[string]interface{}) *schema.ResourceData {
	r := ResourceAreaTree()
	return schema.TestResourceDataRaw(t, r.Schema, input)
}

func mustJSON(t *testing.T, v interface{}) string {
	raw, err := json.Marshal(v)
	require.NoError(t, err)
	return string(raw)
}

func TestNormalizeAreaPath(t *testing.T) {
	assert.Equal(t, "Team A", normalizeAreaPath("  Team A  "))
	assert.Equal(t, "", normalizeAreaPath("   "))
	assert.Equal(t, "Team A", normalizeAreaPath("Team A"))
}

func TestValidateAreaPath(t *testing.T) {
	_, errs := validateAreaPath("Team A", "paths")
	assert.Empty(t, errs)

	_, errs = validateAreaPath("  Team A  ", "paths")
	assert.Empty(t, errs, "surrounding whitespace should be trimmed before validation")

	_, errs = validateAreaPath("", "paths")
	assert.NotEmpty(t, errs, "empty segment should be invalid")

	_, errs = validateAreaPath("Team/A", "paths")
	assert.NotEmpty(t, errs, "segment must not contain a slash")

	_, errs = validateAreaPath(123, "paths")
	assert.NotEmpty(t, errs, "non-string value should be invalid")
}

func TestExpandAreaTree(t *testing.T) {
	tree, err := expandAreaTree(`{"Team A":{"Sub Area":{}},"Team B":{}}`)
	require.NoError(t, err)
	require.Contains(t, tree, "Team A")
	require.Contains(t, tree, "Team B")
	require.Contains(t, tree["Team A"], "Sub Area")

	// empty/whitespace input is treated as an empty tree, not an error.
	tree, err = expandAreaTree("   ")
	require.NoError(t, err)
	assert.Empty(t, tree)

	_, err = expandAreaTree("not json")
	assert.Error(t, err)

	_, err = expandAreaTree(`["Team A"]`)
	assert.Error(t, err, "a JSON array is not a valid tree")
}

func TestFlattenAreaTree(t *testing.T) {
	tree := areaTree{
		"Team A": areaTree{
			"Sub Area 1": areaTree{},
			"Sub Area 2": areaTree{
				"Grandchild": areaTree{},
			},
		},
		"Team B": areaTree{},
	}

	closure, err := flattenAreaTree(tree)
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{
		"Team A",
		"Team A/Sub Area 1",
		"Team A/Sub Area 2",
		"Team A/Sub Area 2/Grandchild",
		"Team B",
	}, closure)

	// Ancestors must always precede their descendants.
	index := map[string]int{}
	for i, p := range closure {
		index[p] = i
	}
	assert.Less(t, index["Team A"], index["Team A/Sub Area 2"])
	assert.Less(t, index["Team A/Sub Area 2"], index["Team A/Sub Area 2/Grandchild"])
}

func TestFlattenAreaTree_EmptySegment(t *testing.T) {
	tree := areaTree{"   ": areaTree{}}
	_, err := flattenAreaTree(tree)
	assert.Error(t, err)
}

func TestFlattenAreaTree_InvalidSegmentName(t *testing.T) {
	tree := areaTree{"Team/A": areaTree{}}
	_, err := flattenAreaTree(tree)
	assert.Error(t, err)
}

func TestFlattenAreaTree_Empty(t *testing.T) {
	closure, err := flattenAreaTree(areaTree{})
	require.NoError(t, err)
	assert.Empty(t, closure)
}

func TestSortByDepth(t *testing.T) {
	paths := []string{"Team B", "Team A/Sub Area", "Team A"}

	ascending := append([]string{}, paths...)
	sortByDepth(ascending, false)
	assert.Equal(t, []string{"Team A", "Team B", "Team A/Sub Area"}, ascending)

	descending := append([]string{}, paths...)
	sortByDepth(descending, true)
	assert.Equal(t, []string{"Team A/Sub Area", "Team B", "Team A"}, descending)
}

func TestBuildAreaTree(t *testing.T) {
	grandchild := workitemtracking.WorkItemClassificationNode{
		Name: converter.String("Grandchild"),
	}
	child := workitemtracking.WorkItemClassificationNode{
		Name:     converter.String("Sub Area"),
		Children: &[]workitemtracking.WorkItemClassificationNode{grandchild},
	}
	root := workitemtracking.WorkItemClassificationNode{
		Name:     converter.String("Area"),
		Children: &[]workitemtracking.WorkItemClassificationNode{child},
	}

	tree := buildAreaTree(&root)
	require.Contains(t, tree, "Sub Area")
	require.Contains(t, tree["Sub Area"], "Grandchild")

	assert.Empty(t, buildAreaTree(nil))
}

func TestResourceAreaTreeCreate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	projectID := uuid.NewString()

	// "Team A" is created first (no parent), then "Team A/Sub Area" (with
	// parent "Team A"), because flattenAreaTree sorts ascending by depth.
	mockClient.EXPECT().
		CreateOrUpdateClassificationNode(clients.Ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, args workitemtracking.CreateOrUpdateClassificationNodeArgs) (*workitemtracking.WorkItemClassificationNode, error) {
			assert.Nil(t, args.Path)
			assert.Equal(t, "Team A", *args.PostedNode.Name)
			return &workitemtracking.WorkItemClassificationNode{}, nil
		}).Times(1)
	mockClient.EXPECT().
		CreateOrUpdateClassificationNode(clients.Ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, args workitemtracking.CreateOrUpdateClassificationNodeArgs) (*workitemtracking.WorkItemClassificationNode, error) {
			require.NotNil(t, args.Path)
			assert.Equal(t, "Team A", *args.Path)
			assert.Equal(t, "Sub Area", *args.PostedNode.Name)
			return &workitemtracking.WorkItemClassificationNode{}, nil
		}).Times(1)

	// Read after create looks up every node in the closure.
	mockClient.EXPECT().
		GetClassificationNode(clients.Ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, args workitemtracking.GetClassificationNodeArgs) (*workitemtracking.WorkItemClassificationNode, error) {
			return &workitemtracking.WorkItemClassificationNode{
				Id:   converter.Int(1),
				Path: converter.String(*args.Path),
			}, nil
		}).Times(2)

	d := getAreaTreeResourceData(t, map[string]interface{}{
		"project_id": projectID,
		"paths": mustJSON(t, areaTree{
			"Team A": areaTree{"Sub Area": areaTree{}},
		}),
	})

	diags := resourceAreaTreeCreate(context.Background(), d, clients)
	assert.Empty(t, diags)
	assert.Equal(t, projectID, d.Id())
}

func TestResourceAreaTreeCreate_InvalidPaths(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	d := getAreaTreeResourceData(t, map[string]interface{}{
		"project_id": uuid.NewString(),
		"paths":      "not json",
	})

	diags := resourceAreaTreeCreate(context.Background(), d, clients)
	require.NotEmpty(t, diags)
	assert.True(t, diags.HasError())
}

func TestResourceAreaTreeDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	projectID := uuid.NewString()

	var deletedOrder []string
	mockClient.EXPECT().
		DeleteClassificationNode(clients.Ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, args workitemtracking.DeleteClassificationNodeArgs) error {
			deletedOrder = append(deletedOrder, *args.Path)
			return nil
		}).Times(2)

	d := getAreaTreeResourceData(t, map[string]interface{}{
		"project_id": projectID,
		"paths": mustJSON(t, areaTree{
			"Team A": areaTree{"Sub Area": areaTree{}},
		}),
	})

	diags := resourceAreaTreeDelete(context.Background(), d, clients)
	assert.Empty(t, diags)
	// Children must be deleted before their parents.
	require.Equal(t, []string{"Team A/Sub Area", "Team A"}, deletedOrder)
}

func TestDiffAreaTreeClosures(t *testing.T) {
	oldClosure := []string{"Team A", "Team A/Sub Area"}
	newClosure := []string{"Team A", "Team B"}

	toCreate, toRemove := diffAreaTreeClosures(oldClosure, newClosure)

	assert.Equal(t, []string{"Team B"}, toCreate)
	assert.Equal(t, []string{"Team A/Sub Area"}, toRemove)
}

func TestDiffAreaTreeClosures_PrunesDeepestFirst(t *testing.T) {
	oldClosure := []string{"Team A", "Team A/Sub Area", "Team A/Sub Area/Grandchild"}
	newClosure := []string{}

	toCreate, toRemove := diffAreaTreeClosures(oldClosure, newClosure)

	assert.Empty(t, toCreate)
	assert.Equal(t, []string{"Team A/Sub Area/Grandchild", "Team A/Sub Area", "Team A"}, toRemove)
}

func TestResourceAreaTreeUpdate_CreatesAndPrunes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	projectID := uuid.NewString()

	mockClient.EXPECT().
		DeleteClassificationNode(clients.Ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, args workitemtracking.DeleteClassificationNodeArgs) error {
			assert.Equal(t, "Team A", *args.Path)
			return nil
		}).Times(1)

	mockClient.EXPECT().
		CreateOrUpdateClassificationNode(clients.Ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, args workitemtracking.CreateOrUpdateClassificationNodeArgs) (*workitemtracking.WorkItemClassificationNode, error) {
			assert.Equal(t, "Team B", *args.PostedNode.Name)
			return &workitemtracking.WorkItemClassificationNode{}, nil
		}).Times(1)

	toCreate, toRemove := diffAreaTreeClosures([]string{"Team A"}, []string{"Team B"})
	require.Equal(t, []string{"Team B"}, toCreate)
	require.Equal(t, []string{"Team A"}, toRemove)

	require.NoError(t, deleteAreaNodes(context.Background(), clients, projectID, toRemove, time.Minute))
	require.NoError(t, ensureAreaNodes(context.Background(), clients, projectID, toCreate, time.Minute))
}

func TestResourceAreaTreeImport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{WorkItemTrackingClient: mockClient, Ctx: context.Background()}

	projectID := uuid.NewString()

	subArea := workitemtracking.WorkItemClassificationNode{Name: converter.String("Sub Area")}
	teamA := workitemtracking.WorkItemClassificationNode{
		Name:     converter.String("Team A"),
		Children: &[]workitemtracking.WorkItemClassificationNode{subArea},
	}
	root := workitemtracking.WorkItemClassificationNode{
		Children: &[]workitemtracking.WorkItemClassificationNode{teamA},
	}

	mockClient.EXPECT().GetClassificationNode(clients.Ctx, gomock.Any()).Return(&root, nil).Times(1)

	d := getAreaTreeResourceData(t, map[string]interface{}{})
	d.SetId(projectID)

	results, err := resourceAreaTreeImport(context.Background(), d, clients)
	require.NoError(t, err)
	require.Len(t, results, 1)

	imported := results[0]
	assert.Equal(t, projectID, imported.Get("project_id"))

	tree, err := expandAreaTree(imported.Get("paths").(string))
	require.NoError(t, err)
	require.Contains(t, tree, "Team A")
	require.Contains(t, tree["Team A"], "Sub Area")
}
