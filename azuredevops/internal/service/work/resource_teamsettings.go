package work

import (
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/work"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceTeamSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateOrUpdateTeamSettings,
		Read:   resourceReadTeamSettings,
		Update: resourceCreateOrUpdateTeamSettings,
		Delete: resourceDeleteTeamSettings,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
			"bugs_behavior": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "off",
				ValidateFunc: validation.StringInSlice([]string{
					"asRequirements",
					"asTasks",
					"off",
				}, false),
			},
			"working_days": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday",
					}, false),
				},
			},
			"backlog_iteration_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_iteration_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"default_iteration_macro"},
			},
			"default_iteration_macro": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"default_iteration_id"},
			},
		},
	}
}

func resourceCreateOrUpdateTeamSettings(d *schema.ResourceData, m interface{}) error {

	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)
	teamID := d.Get("team_id").(string)

	patch := expandTeamSettingsPatch(d)

	err := updateTeamSettingsInternal(clients, projectID, teamID, patch)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/settings", teamID))

	return resourceReadTeamSettings(d, m)
}

func resourceReadTeamSettings(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID := converter.String(d.Get("project_id").(string))
	teamID := converter.String(d.Get("team_id").(string))

	args := work.GetTeamSettingsArgs{
		Project: projectID,
		Team:    teamID,
	}

	settings, err := clients.WorkClient.GetTeamSettings(clients.Ctx, args)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading team settings: %+v", err)
	}

	if settings.BugsBehavior != nil {
		d.Set("bugs_behavior", string(*settings.BugsBehavior))
	}

	if settings.WorkingDays != nil {
		days := make([]interface{}, len(*settings.WorkingDays))
		for i, day := range *settings.WorkingDays {
			days[i] = string(day)
		}
		d.Set("working_days", schema.NewSet(schema.HashString, days))
	}

	if settings.BacklogIteration != nil && settings.BacklogIteration.Id != nil {
		d.Set("backlog_iteration_id", (*settings.BacklogIteration.Id).String())
	}

	if settings.DefaultIteration != nil && settings.DefaultIteration.Id != nil {
		d.Set("default_iteration_id", (*settings.DefaultIteration.Id).String())
	}
	if settings.DefaultIterationMacro != nil {
		d.Set("default_iteration_macro", *settings.DefaultIterationMacro)
	}

	return nil
}

func resourceDeleteTeamSettings(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID := converter.String(d.Get("project_id").(string))
	teamID := converter.String(d.Get("team_id").(string))

	defaultMacro := "@CurrentIteration"
	defaultBugs := work.BugsBehavior("asRequirements")

	defaultDays := []string{
		"monday", "tuesday", "wednesday", "thursday", "friday",
	}

	resetPatch := work.TeamSettingsPatch{
		BugsBehavior:          &defaultBugs,
		WorkingDays:           &defaultDays,
		DefaultIterationMacro: &defaultMacro,
	}

	err := updateTeamSettingsInternal(clients, *projectID, *teamID, resetPatch)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func updateTeamSettingsInternal(clients *client.AggregatedClient, projectID string, teamID string, patch work.TeamSettingsPatch) error {
	args := work.UpdateTeamSettingsArgs{
		TeamSettingsPatch: &patch,
		Project:           &projectID,
		Team:              &teamID,
	}

	_, err := clients.WorkClient.UpdateTeamSettings(clients.Ctx, args)
	if err != nil {
		return fmt.Errorf("failed to update team settings for team %s: %+v", teamID, err)
	}

	return nil
}

func expandTeamSettingsPatch(d *schema.ResourceData) work.TeamSettingsPatch {
	patch := work.TeamSettingsPatch{}

	if v, ok := d.GetOk("bugs_behavior"); ok {
		behavior := work.BugsBehavior(v.(string))
		patch.BugsBehavior = &behavior
	}

	if v, ok := d.GetOk("working_days"); ok {
		tfDays := v.(*schema.Set).List()
		apiDays := make([]string, 0, len(tfDays))

		for _, day := range tfDays {
			apiDays = append(apiDays, day.(string))
		}

		sort.Slice(apiDays, func(i, j int) bool {
			return dayOrder[apiDays[i]] < dayOrder[apiDays[j]]
		})

		patch.WorkingDays = &apiDays
	}

	if v, ok := d.GetOk("backlog_iteration_id"); ok {
		id, _ := uuid.Parse(v.(string))
		patch.BacklogIteration = &id
	}

	if v, ok := d.GetOk("default_iteration_macro"); ok {
		macro := v.(string)
		patch.DefaultIterationMacro = &macro
	} else if v, ok := d.GetOk("default_iteration_id"); ok {
		id, _ := uuid.Parse(v.(string))
		patch.DefaultIteration = &id
	}

	return patch
}

var dayOrder = map[string]int{
	"monday":    0,
	"tuesday":   1,
	"wednesday": 2,
	"thursday":  3,
	"friday":    4,
	"saturday":  5,
	"sunday":    6,
}
