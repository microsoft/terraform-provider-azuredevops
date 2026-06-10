package core

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/sdk/teamsettings"
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
		Create: resourceTeamAreaPathCreate,
		Read:   resourceTeamAreaPathRead,
		Update: resourceTeamAreaPathUpdate,
		Delete: resourceTeamAreaPathDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: importTeamAreaPathState,
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

func resourceTeamAreaPathCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	teamID := d.Get("team_id").(string)
	areaPath := d.Get("area_path").(string)
	includeChildren := d.Get("include_children").(bool)

	mu := getTeamAreaPathMutex(projectID, teamID)
	mu.Lock()
	defer mu.Unlock()

	current, err := clients.TeamSettingsClient.GetTeamFieldValues(clients.Ctx, teamsettings.GetTeamFieldValuesArgs{
		Project: &projectID,
		Team:    &teamID,
	})
	if err != nil {
		return fmt.Errorf("reading team field values for team %s: %w", teamID, err)
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
		values = append(values, teamsettings.TeamFieldValueReference{
			Value:           &areaPath,
			IncludeChildren: &includeChildren,
		})
	}

	defaultValue := current.DefaultValue
	if defaultValue == nil || *defaultValue == "" {
		defaultValue = &areaPath
	}

	patch := &teamsettings.TeamFieldValues{
		DefaultValue: defaultValue,
		Values:       &values,
	}

	result, err := clients.TeamSettingsClient.UpdateTeamFieldValues(clients.Ctx, teamsettings.UpdateTeamFieldValuesArgs{
		Project:         &projectID,
		Team:            &teamID,
		TeamFieldValues: patch,
	})
	if err != nil {
		return fmt.Errorf("updating team field values for team %s: %w", teamID, err)
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
		return fmt.Errorf("area path %q was not accepted by the API — ensure the area path node exists in project %s before assigning it to a team", areaPath, projectID)
	}

	d.SetId(buildTeamAreaPathResourceID(projectID, teamID, areaPath))
	log.Printf("[DEBUG] azuredevops_team_area_path created: %s", d.Id())

	return resourceTeamAreaPathRead(d, m)
}

func resourceTeamAreaPathRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID, teamID, areaPath, err := parseTeamAreaPathResourceID(d.Id())
	if err != nil {
		return err
	}

	current, err := clients.TeamSettingsClient.GetTeamFieldValues(clients.Ctx, teamsettings.GetTeamFieldValuesArgs{
		Project: &projectID,
		Team:    &teamID,
	})
	if err != nil {
		return fmt.Errorf("reading team field values for team %s: %w", teamID, err)
	}

	values := normalizeValues(current.Values)
	for _, v := range values {
		if v.Value == nil {
			continue
		}
		if strings.EqualFold(*v.Value, areaPath) {
			d.Set("project_id", projectID)
			d.Set("team_id", teamID)
			d.Set("area_path", *v.Value)
			includeChildren := false
			if v.IncludeChildren != nil {
				includeChildren = *v.IncludeChildren
			}
			d.Set("include_children", includeChildren)
			return nil
		}
	}

	log.Printf("[DEBUG] azuredevops_team_area_path %s not found, removing from state", d.Id())
	d.SetId("")
	return nil
}

func resourceTeamAreaPathUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	teamID := d.Get("team_id").(string)
	areaPath := d.Get("area_path").(string)
	includeChildren := d.Get("include_children").(bool)

	mu := getTeamAreaPathMutex(projectID, teamID)
	mu.Lock()
	defer mu.Unlock()

	current, err := clients.TeamSettingsClient.GetTeamFieldValues(clients.Ctx, teamsettings.GetTeamFieldValuesArgs{
		Project: &projectID,
		Team:    &teamID,
	})
	if err != nil {
		return fmt.Errorf("reading team field values for team %s: %w", teamID, err)
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
		return fmt.Errorf("area path %q not found in team %s field values during update", areaPath, teamID)
	}

	patch := &teamsettings.TeamFieldValues{
		DefaultValue: current.DefaultValue,
		Values:       &values,
	}

	_, err = clients.TeamSettingsClient.UpdateTeamFieldValues(clients.Ctx, teamsettings.UpdateTeamFieldValuesArgs{
		Project:         &projectID,
		Team:            &teamID,
		TeamFieldValues: patch,
	})
	if err != nil {
		return fmt.Errorf("updating team field values for team %s: %w", teamID, err)
	}

	return resourceTeamAreaPathRead(d, m)
}

func resourceTeamAreaPathDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	teamID := d.Get("team_id").(string)
	areaPath := d.Get("area_path").(string)

	mu := getTeamAreaPathMutex(projectID, teamID)
	mu.Lock()
	defer mu.Unlock()

	current, err := clients.TeamSettingsClient.GetTeamFieldValues(clients.Ctx, teamsettings.GetTeamFieldValuesArgs{
		Project: &projectID,
		Team:    &teamID,
	})
	if err != nil {
		return fmt.Errorf("reading team field values for team %s: %w", teamID, err)
	}

	values := normalizeValues(current.Values)
	var remaining []teamsettings.TeamFieldValueReference
	for _, v := range values {
		if v.Value != nil && strings.EqualFold(*v.Value, areaPath) {
			continue
		}
		remaining = append(remaining, v)
	}

	if current.DefaultValue != nil && strings.EqualFold(*current.DefaultValue, areaPath) {
		if len(remaining) == 0 {
			log.Printf("[DEBUG] azuredevops_team_area_path %s is the only value and the default; removing from state only", d.Id())
			return nil
		}
		current.DefaultValue = remaining[0].Value
	}

	patch := &teamsettings.TeamFieldValues{
		DefaultValue: current.DefaultValue,
		Values:       &remaining,
	}

	_, err = clients.TeamSettingsClient.UpdateTeamFieldValues(clients.Ctx, teamsettings.UpdateTeamFieldValuesArgs{
		Project:         &projectID,
		Team:            &teamID,
		TeamFieldValues: patch,
	})
	if err != nil {
		return fmt.Errorf("removing area path %q from team %s: %w", areaPath, teamID, err)
	}

	log.Printf("[DEBUG] azuredevops_team_area_path deleted: %s", d.Id())
	return nil
}

func importTeamAreaPathState(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	projectID, teamID, areaPath, err := parseTeamAreaPathResourceID(d.Id())
	if err != nil {
		return nil, err
	}

	_ = d.Set("project_id", projectID)
	_ = d.Set("team_id", teamID)
	_ = d.Set("area_path", areaPath)
	d.SetId(buildTeamAreaPathResourceID(projectID, teamID, areaPath))

	return []*schema.ResourceData{d}, nil
}

func buildTeamAreaPathResourceID(projectID, teamID, areaPath string) string {
	return fmt.Sprintf("%s/%s/%s", projectID, teamID, url.QueryEscape(areaPath))
}

func parseTeamAreaPathResourceID(id string) (string, string, string, error) {
	parts := strings.SplitN(id, "/", 3)
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("unexpected ID format (%q), expected: <project_id>/<team_id>/<url-escaped-area_path>", id)
	}

	areaPath, err := url.QueryUnescape(parts[2])
	if err != nil {
		return "", "", "", fmt.Errorf("failed to decode area_path from ID %q: %w", id, err)
	}

	return parts[0], parts[1], areaPath, nil
}

func normalizeValues(v *[]teamsettings.TeamFieldValueReference) []teamsettings.TeamFieldValueReference {
	if v == nil {
		return []teamsettings.TeamFieldValueReference{}
	}
	return *v
}
