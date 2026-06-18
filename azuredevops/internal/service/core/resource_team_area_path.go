package core

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/work"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

var teamAreaPathMutexes = struct {
	sync.Mutex
	m map[string]*sync.Mutex
}{m: make(map[string]*sync.Mutex)}

func getTeamAreaPathMutex(projectID, teamID string) *sync.Mutex {
	key := projectID + "/" + teamID
	teamAreaPathMutexes.Lock()
	defer teamAreaPathMutexes.Unlock()
	if mu, ok := teamAreaPathMutexes.m[key]; ok {
		return mu
	}
	mu := &sync.Mutex{}
	teamAreaPathMutexes.m[key] = mu
	return mu
}

func ResourceTeamAreaPath() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamAreaPathCreate,
		ReadContext:   resourceTeamAreaPathRead,
		UpdateContext: resourceTeamAreaPathUpdate,
		DeleteContext: resourceTeamAreaPathDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: importTeamAreaPathState,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"team_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"area_path": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"include_children": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceTeamAreaPathCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	teamID := d.Get("team_id").(string)
	areaPath := d.Get("area_path").(string)
	includeChildren := d.Get("include_children").(bool)

	mu := getTeamAreaPathMutex(projectID, teamID)
	mu.Lock()
	defer mu.Unlock()

	current, err := clients.WorkClient.GetTeamFieldValues(ctx, work.GetTeamFieldValuesArgs{
		Project: &projectID,
		Team:    &teamID,
	})
	if err != nil {
		return diag.Errorf("reading team field values for team %s: %+v", teamID, err)
	}

	values := normalizeValues(current.Values)

	found := false
	for i := range values {
		if values[i].Value != nil && strings.EqualFold(*values[i].Value, areaPath) {
			values[i].IncludeChildren = &includeChildren
			found = true
			break
		}
	}
	if !found {
		values = append(values, work.TeamFieldValue{
			Value:           &areaPath,
			IncludeChildren: &includeChildren,
		})
	}

	defaultValue := current.DefaultValue
	if defaultValue == nil || *defaultValue == "" {
		defaultValue = &areaPath
	}

	patch := &work.TeamFieldValuesPatch{
		DefaultValue: defaultValue,
		Values:       &values,
	}

	result, err := clients.WorkClient.UpdateTeamFieldValues(ctx, work.UpdateTeamFieldValuesArgs{
		Project: &projectID,
		Team:    &teamID,
		Patch:   patch,
	})
	if err != nil {
		return diag.Errorf("updating team field values for team %s: %+v", teamID, err)
	}

	resultValues := normalizeValues(result.Values)
	pathFound := false
	for _, v := range resultValues {
		if v.Value != nil && strings.EqualFold(*v.Value, areaPath) {
			pathFound = true
			break
		}
	}
	if !pathFound {
		return diag.Errorf("area path %q was not accepted by the API — ensure the area path node exists in project %s before assigning it to a team", areaPath, projectID)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", projectID, teamID, areaPath))

	return resourceTeamAreaPathRead(ctx, d, m)
}

func resourceTeamAreaPathRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	teamID := d.Get("team_id").(string)
	areaPath := d.Get("area_path").(string)

	stateConf := &retry.StateChangeConf{
		Pending: []string{"notfound"},
		Target:  []string{"found"},
		Refresh: func() (interface{}, string, error) {
			current, err := clients.WorkClient.GetTeamFieldValues(ctx, work.GetTeamFieldValuesArgs{
				Project: &projectID,
				Team:    &teamID,
			})
			if err != nil {
				return nil, "", fmt.Errorf("reading team field values for team %s: %+v", teamID, err)
			}

			for _, v := range normalizeValues(current.Values) {
				if v.Value != nil && strings.EqualFold(*v.Value, areaPath) {
					return &v, "found", nil
				}
			}
			return nil, "notfound", nil
		},
		Timeout:    1 * time.Minute,
		MinTimeout: 5 * time.Second,
		Delay:      1 * time.Second,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		d.SetId("")
		return nil
	}

	v := result.(*work.TeamFieldValue)
	includeChildren := false
	if v.IncludeChildren != nil {
		includeChildren = *v.IncludeChildren
	}
	d.Set("include_children", includeChildren)
	return nil
}

func resourceTeamAreaPathUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	teamID := d.Get("team_id").(string)
	areaPath := d.Get("area_path").(string)
	includeChildren := d.Get("include_children").(bool)

	mu := getTeamAreaPathMutex(projectID, teamID)
	mu.Lock()
	defer mu.Unlock()

	current, err := clients.WorkClient.GetTeamFieldValues(ctx, work.GetTeamFieldValuesArgs{
		Project: &projectID,
		Team:    &teamID,
	})
	if err != nil {
		return diag.Errorf("reading team field values for team %s: %+v", teamID, err)
	}

	values := normalizeValues(current.Values)
	found := false
	for i := range values {
		if values[i].Value != nil && strings.EqualFold(*values[i].Value, areaPath) {
			values[i].IncludeChildren = &includeChildren
			found = true
			break
		}
	}
	if !found {
		return diag.Errorf("area path %q not found in team %s field values during update", areaPath, teamID)
	}

	patch := &work.TeamFieldValuesPatch{
		DefaultValue: current.DefaultValue,
		Values:       &values,
	}

	_, err = clients.WorkClient.UpdateTeamFieldValues(ctx, work.UpdateTeamFieldValuesArgs{
		Project: &projectID,
		Team:    &teamID,
		Patch:   patch,
	})
	if err != nil {
		return diag.Errorf("updating team field values for team %s: %+v", teamID, err)
	}

	return resourceTeamAreaPathRead(ctx, d, m)
}

func resourceTeamAreaPathDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	teamID := d.Get("team_id").(string)
	areaPath := d.Get("area_path").(string)

	mu := getTeamAreaPathMutex(projectID, teamID)
	mu.Lock()
	defer mu.Unlock()

	current, err := clients.WorkClient.GetTeamFieldValues(ctx, work.GetTeamFieldValuesArgs{
		Project: &projectID,
		Team:    &teamID,
	})
	if err != nil {
		return diag.Errorf("reading team field values for team %s: %+v", teamID, err)
	}
	if current.Values == nil || len(*current.Values) == 0 {
		return diag.Errorf("team %s has no area path values to delete", teamID)
	}

	var remaining []work.TeamFieldValue
	for _, v := range *current.Values {
		if v.Value != nil && strings.EqualFold(*v.Value, areaPath) {
			continue
		}
		remaining = append(remaining, v)
	}

	if current.DefaultValue != nil && strings.EqualFold(*current.DefaultValue, areaPath) {
		if len(remaining) == 0 {
			return nil
		}
		current.DefaultValue = remaining[0].Value
	}

	patch := &work.TeamFieldValuesPatch{
		DefaultValue: current.DefaultValue,
		Values:       &remaining,
	}

	_, err = clients.WorkClient.UpdateTeamFieldValues(ctx, work.UpdateTeamFieldValuesArgs{
		Project: &projectID,
		Team:    &teamID,
		Patch:   patch,
	})
	if err != nil {
		return diag.Errorf("removing area path %q from team %s: %+v", areaPath, teamID, err)
	}

	return nil
}

func importTeamAreaPathState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("unexpected ID format (%q), expected: <project_id>/<team_id>/<area_path>", d.Id())
	}

	d.Set("project_id", parts[0])
	d.Set("team_id", parts[1])
	d.Set("area_path", parts[2])
	d.SetId(fmt.Sprintf("%s/%s/%s", parts[0], parts[1], parts[2]))

	return []*schema.ResourceData{d}, nil
}

func normalizeValues(v *[]work.TeamFieldValue) []work.TeamFieldValue {
	if v == nil {
		return []work.TeamFieldValue{}
	}
	return *v
}
