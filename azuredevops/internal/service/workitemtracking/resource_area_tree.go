package workitemtracking

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// areaTree is the recursive shape of the "paths" attribute: every key is an
// area name and its value is the (possibly empty) tree of that area's
// children, allowing area hierarchies of arbitrary depth to be expressed
// directly as a nested object rather than as a flat list/set of
// slash-separated path strings.
type areaTree map[string]areaTree

// ResourceAreaTree manages an entire area path hierarchy for a project as a
// single resource. Unlike azuredevops_area (which manages one node per
// resource instance), azuredevops_area_tree takes the complete, desired set
// of area paths and owns their full lifecycle - creation, renaming/pruning
// on update, and teardown on delete - in one place. Because there is only
// ever one resource instance per project (no for_each over individual
// nodes), there is no possibility of the "Cycle" error that arises when
// Terraform resources of the same type try to reference sibling instances
// of themselves, and orphaned ancestor nodes are correctly pruned on update
// since this resource has full visibility into both the old and new
// desired state.
func ResourceAreaTree() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAreaTreeCreate,
		ReadContext:   resourceAreaTreeRead,
		UpdateContext: resourceAreaTreeUpdate,
		DeleteContext: resourceAreaTreeDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceAreaTreeImport,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			// The full area path hierarchy to manage, expressed as an
			// infinitely-deep JSON object tree where every key is an area
			// name and its value is the (possibly empty) object describing
			// that area's children, e.g.:
			//   jsonencode({
			//     "Team A" = { "Sub Area" = {} }
			//     "Team B" = {}
			//   })
			// Every node in the tree - including intermediate/ancestor
			// nodes like "Team A" above - is created and removed
			// automatically as part of managing this resource.
			"paths": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validatePaths,
			},
			// Map of every managed area path (including implied ancestors)
			// to its integer node ID.
			"area_ids": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Map of every managed area path (including implied ancestors)
			// to its full canonical Azure DevOps path
			// (e.g. "ProjectName\Area\Team A\Sub Area"), suitable for direct
			// use as the "path" of an azuredevops_team area assignment.
			"area_paths": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// normalizeAreaPath trims surrounding whitespace from a single area path
// segment (an area name as it appears as a key in the paths tree).
func normalizeAreaPath(raw string) string {
	return strings.TrimSpace(raw)
}

// validateAreaPath validates a single area path segment (a key in the paths
// tree), reusing the same naming rules Azure DevOps enforces for
// classification node names.
func validateAreaPath(v interface{}, k string) (warnings []string, errors []error) {
	value, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("%q must be a string", k))
		return warnings, errors
	}
	return validateClassificationNodeName(normalizeAreaPath(value), k)
}

// validatePaths validates the raw "paths" attribute: it must be a JSON
// object that parses into an areaTree, and every key at every depth must be
// a valid area path segment.
func validatePaths(v interface{}, k string) (warnings []string, errors []error) {
	value, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("%q must be a string", k))
		return warnings, errors
	}

	tree, err := expandAreaTree(value)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q is invalid: %w", k, err))
		return warnings, errors
	}

	if _, err := flattenAreaTree(tree); err != nil {
		errors = append(errors, fmt.Errorf("%q is invalid: %w", k, err))
	}

	return warnings, errors
}

// expandAreaTree parses the raw JSON representation of the "paths"
// attribute into an areaTree. An empty/whitespace-only string is treated as
// an empty tree.
func expandAreaTree(raw string) (areaTree, error) {
	tree := areaTree{}

	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return tree, nil
	}

	if err := json.Unmarshal([]byte(trimmed), &tree); err != nil {
		return nil, fmt.Errorf("parsing paths as a JSON object tree: %w", err)
	}

	return tree, nil
}

// flattenAreaTree walks tree depth-first and returns the full, deduplicated
// closure of area paths - every node plus every one of its ancestors -
// sorted ascending by depth (shallowest first) so that processing the list
// in order always creates parents before their children.
func flattenAreaTree(tree areaTree) ([]string, error) {
	seen := map[string]bool{}
	var closure []string

	var walk func(node areaTree, prefix []string) error
	walk = func(node areaTree, prefix []string) error {
		for rawName, children := range node {
			name := normalizeAreaPath(rawName)
			if name == "" {
				return fmt.Errorf("area path segment must not be empty or whitespace (parent: %q)", strings.Join(prefix, "/"))
			}
			if _, errs := validateClassificationNodeName(name, "paths"); len(errs) > 0 {
				return errs[0]
			}

			full := append(append([]string{}, prefix...), name)
			key := strings.Join(full, "/")
			if !seen[key] {
				seen[key] = true
				closure = append(closure, key)
			}

			if err := walk(children, full); err != nil {
				return err
			}
		}
		return nil
	}

	if err := walk(tree, nil); err != nil {
		return nil, err
	}

	sortByDepth(closure, false)
	return closure, nil
}

// sortByDepth sorts area path keys by their nesting depth (number of "/"
// separators). When descending is true, deepest paths sort first, which is
// the safe order for deletion (children before parents).
func sortByDepth(paths []string, descending bool) {
	sort.Slice(paths, func(i, j int) bool {
		di, dj := strings.Count(paths[i], "/"), strings.Count(paths[j], "/")
		if di != dj {
			if descending {
				return di > dj
			}
			return di < dj
		}
		if descending {
			return paths[i] > paths[j]
		}
		return paths[i] < paths[j]
	})
}

func toStringSet(paths []string) map[string]bool {
	set := make(map[string]bool, len(paths))
	for _, p := range paths {
		set[p] = true
	}
	return set
}

func resourceAreaTreeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	tree, err := expandAreaTree(d.Get("paths").(string))
	if err != nil {
		return diag.Errorf("parsing paths for project %q: %+v", projectID, err)
	}
	closure, err := flattenAreaTree(tree)
	if err != nil {
		return diag.Errorf("parsing paths for project %q: %+v", projectID, err)
	}

	if err := ensureAreaNodes(ctx, clients, projectID, closure, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("creating area tree for project %q: %+v", projectID, err)
	}

	d.SetId(projectID)

	return resourceAreaTreeRead(ctx, d, m)
}

// ensureAreaNodes idempotently creates every path in paths, in order.
// paths must already be sorted ascending by depth (see flattenAreaTree) so
// that, for every entry, its parent either already exists in Azure DevOps or
// was created earlier in this same call.
func ensureAreaNodes(ctx context.Context, clients *client.AggregatedClient, projectID string, paths []string, timeout time.Duration) error {
	for _, p := range paths {
		segments := strings.Split(p, "/")
		name := segments[len(segments)-1]

		var parentArg *string
		if len(segments) > 1 {
			parent := strings.Join(segments[:len(segments)-1], "/")
			parentArg = &parent
		}

		err := utils.RetryOn(ctx, timeout, func() error {
			_, err := clients.WorkItemTrackingClient.CreateOrUpdateClassificationNode(clients.Ctx, workitemtracking.CreateOrUpdateClassificationNodeArgs{
				Project:        &projectID,
				StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
				Path:           parentArg,
				PostedNode: &workitemtracking.WorkItemClassificationNode{
					Name: &name,
				},
			})
			return err
		}, utils.RetryOnNotFound, utils.RetryOnUnexpectedException)
		if err != nil {
			return fmt.Errorf("ensuring area path %q: %w", p, err)
		}
	}
	return nil
}

// deleteAreaNodes deletes every path in paths, in order. paths must already
// be sorted descending by depth (deepest first) so that children are always
// removed before their parents. Missing nodes (already deleted, e.g. by a
// previous partially-applied run) are treated as success.
func deleteAreaNodes(ctx context.Context, clients *client.AggregatedClient, projectID string, paths []string, timeout time.Duration) error {
	for _, p := range paths {
		apiPath := p
		err := utils.RetryOn(ctx, timeout, func() error {
			err := clients.WorkItemTrackingClient.DeleteClassificationNode(clients.Ctx, workitemtracking.DeleteClassificationNodeArgs{
				Project:        &projectID,
				StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
				Path:           &apiPath,
			})
			if err != nil && utils.ResponseWasNotFound(err) {
				return nil
			}
			return err
		}, utils.RetryOnUnexpectedException)
		if err != nil {
			return fmt.Errorf("deleting area path %q: %w", p, err)
		}
	}
	return nil
}

func resourceAreaTreeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	tree, err := expandAreaTree(d.Get("paths").(string))
	if err != nil {
		return diag.Errorf("parsing paths for project %q: %+v", projectID, err)
	}
	closure, err := flattenAreaTree(tree)
	if err != nil {
		return diag.Errorf("parsing paths for project %q: %+v", projectID, err)
	}

	areaIDs := map[string]interface{}{}
	areaPaths := map[string]interface{}{}

	for _, p := range closure {
		apiPath := p
		node, err := clients.WorkItemTrackingClient.GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
			Project:        &projectID,
			StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
			Path:           &apiPath,
		})
		if err != nil {
			if utils.ResponseWasNotFound(err) {
				// Node drifted away out-of-band. Omit it from the computed
				// maps; the resulting diff on area_ids/area_paths will
				// surface the drift so the next apply can recreate it.
				continue
			}
			return diag.Errorf("reading area path %q: %+v", p, err)
		}
		if node.Id != nil {
			areaIDs[p] = strconv.Itoa(*node.Id)
		}
		if node.Path != nil {
			areaPaths[p] = *node.Path
		}
	}

	if len(areaIDs) == 0 {
		// Nothing in the desired tree exists anymore; treat as deleted.
		d.SetId("")
		return nil
	}

	d.Set("area_ids", areaIDs)
	d.Set("area_paths", areaPaths)

	return nil
}

// diffAreaTreeClosures compares the old and new desired area tree closures
// and returns the paths to create (ascending by depth) and the paths to
// remove (descending by depth, safe for deletion).
func diffAreaTreeClosures(oldClosure, newClosure []string) (toCreate, toRemove []string) {
	oldSet := toStringSet(oldClosure)
	newSet := toStringSet(newClosure)

	for _, p := range newClosure {
		if !oldSet[p] {
			toCreate = append(toCreate, p)
		}
	}
	// newClosure is already sorted ascending by depth, and filtering
	// preserves that order.

	for _, p := range oldClosure {
		if !newSet[p] {
			toRemove = append(toRemove, p)
		}
	}
	// A path being removed can never have a surviving descendant in
	// newClosure: if it did, that descendant's own ancestor closure would
	// include this path too, keeping it in newSet. So it's always safe to
	// delete toRemove entries deepest-first without orphaning a node that
	// still has live children.
	sortByDepth(toRemove, true)

	return toCreate, toRemove
}

func resourceAreaTreeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	oldRaw, newRaw := d.GetChange("paths")

	oldTree, err := expandAreaTree(oldRaw.(string))
	if err != nil {
		return diag.Errorf("parsing previous paths for project %q: %+v", projectID, err)
	}
	newTree, err := expandAreaTree(newRaw.(string))
	if err != nil {
		return diag.Errorf("parsing paths for project %q: %+v", projectID, err)
	}

	oldClosure, err := flattenAreaTree(oldTree)
	if err != nil {
		return diag.Errorf("parsing previous paths for project %q: %+v", projectID, err)
	}
	newClosure, err := flattenAreaTree(newTree)
	if err != nil {
		return diag.Errorf("parsing paths for project %q: %+v", projectID, err)
	}

	toCreate, toRemove := diffAreaTreeClosures(oldClosure, newClosure)

	if err := deleteAreaNodes(ctx, clients, projectID, toRemove, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return diag.Errorf("pruning area tree for project %q: %+v", projectID, err)
	}
	if err := ensureAreaNodes(ctx, clients, projectID, toCreate, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return diag.Errorf("updating area tree for project %q: %+v", projectID, err)
	}

	return resourceAreaTreeRead(ctx, d, m)
}

func resourceAreaTreeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	tree, err := expandAreaTree(d.Get("paths").(string))
	if err != nil {
		return diag.Errorf("parsing paths for project %q: %+v", projectID, err)
	}
	closure, err := flattenAreaTree(tree)
	if err != nil {
		return diag.Errorf("parsing paths for project %q: %+v", projectID, err)
	}
	sortByDepth(closure, true)

	if err := deleteAreaNodes(ctx, clients, projectID, closure, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("deleting area tree for project %q: %+v", projectID, err)
	}

	return nil
}

func resourceAreaTreeImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	projectID := d.Id()
	clients := m.(*client.AggregatedClient)

	node, err := clients.WorkItemTrackingClient.GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
		Depth:          converter.Int(1000),
	})
	if err != nil {
		return nil, fmt.Errorf("reading area tree for project %q: %w", projectID, err)
	}

	tree := buildAreaTree(node)
	raw, err := json.Marshal(tree)
	if err != nil {
		return nil, fmt.Errorf("encoding area tree for project %q: %w", projectID, err)
	}

	d.SetId(projectID)
	d.Set("project_id", projectID)
	d.Set("paths", string(raw))

	return []*schema.ResourceData{d}, nil
}

// buildAreaTree converts the API's classification node tree (as returned by
// GetClassificationNode with a large Depth) into an areaTree suitable for
// storing in the "paths" attribute.
func buildAreaTree(n *workitemtracking.WorkItemClassificationNode) areaTree {
	tree := areaTree{}
	if n == nil || n.Children == nil {
		return tree
	}
	for _, child := range *n.Children {
		name := ""
		if child.Name != nil {
			name = *child.Name
		}
		c := child
		tree[name] = buildAreaTree(&c)
	}
	return tree
}
