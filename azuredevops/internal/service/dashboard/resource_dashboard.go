package dashboard

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/dashboard"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	dashboardextras "github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/dashboardextra"
)

func ResourceDashboard() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDashboardCreate,
		ReadContext:   resourceDashboardRead,
		UpdateContext: resourceDashboardUpdate,
		DeleteContext: resourceDashboardDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), "/")
				if len(idParts) > 3 || len(idParts) < 2 {
					return nil, fmt.Errorf(" Unexpected ID format (%q), Expected: <projetId>/<dasboardId> or <projetId>/<teamId>/<dasboardId>", d.Id())
				}

				d.Set("project_id", idParts[0])
				if len(idParts) == 2 {
					d.SetId(idParts[1])
				}

				if len(idParts) == 3 {
					d.Set("team_id", idParts[1])
					d.SetId(idParts[2])
				}

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},

			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"team_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},

			"refresh_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntInSlice([]int{0, 5}),
			},

			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDashboardCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)

	params := dashboard.CreateDashboardArgs{
		Project: &projectID,
		Dashboard: &dashboard.Dashboard{
			Name:            converter.String(d.Get("name").(string)),
			RefreshInterval: converter.Int(d.Get("refresh_interval").(int)),
			Description:     converter.String(d.Get("description").(string)),
		},
	}

	if v, ok := d.GetOk("team_id"); ok {
		params.Team = converter.String(v.(string))
	}

	resp, err := clients.DashboardClient.CreateDashboard(clients.Ctx, params)

	if err != nil {
		return diag.Errorf(" Creating dashboard in Azure DevOps: %s", err)
	}

	d.SetId(resp.Id.String())
	return resourceDashboardRead(clients.Ctx, d, m)
}

func resourceDashboardRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	id := d.Id()
	dashboardId, err := uuid.Parse(id)
	if err != nil {
		return diag.Errorf(" Parsing dashboard ID: %+v", err)
	}

	params := dashboard.GetDashboardArgs{
		Project:     converter.String(d.Get("project_id").(string)),
		DashboardId: &dashboardId,
	}

	if v, ok := d.GetOk("team_id"); ok {
		params.Team = converter.String(v.(string))
	}

	resp, err := clients.DashboardClient.GetDashboard(clients.Ctx, params)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf(" Getting dashboard with id: %s, %+v", id, err)
	}

	if resp != nil {
		if resp.Name != nil {
			d.Set("name", resp.Name)
		}

		if resp.RefreshInterval != nil {
			d.Set("refresh_interval", resp.RefreshInterval)
		}

		if resp.Description != nil {
			d.Set("description", resp.Description)
		}

		if resp.OwnerId != nil {
			d.Set("owner_id", resp.OwnerId.String())
		}

		if resp.GroupId != nil {
			team, err := clients.CoreClient.GetTeam(clients.Ctx, core.GetTeamArgs{
				ProjectId: converter.String(d.Get("project_id").(string)),
				TeamId:    converter.String(resp.GroupId.String()),
			})

			if err != nil {
				diag.Errorf(" Getting Dashboard Team ID: %s, %+v", resp.GroupId.String(), err)
			}

			if team != nil && team.Name != nil {
				d.Set("team_id", team.Id.String())
			}
		}
	}
	return nil
}

func resourceDashboardUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	id := d.Id()
	dashboardId, err := uuid.Parse(id)
	if err != nil {
		return diag.Errorf(" Parsing Dashboard ID: %+v", err)
	}

	args := dashboard.GetDashboardArgs{
		Project:     converter.String(d.Get("project_id").(string)),
		DashboardId: &dashboardId,
	}
	if v, ok := d.GetOk("team_id"); ok {
		args.Team = converter.String(v.(string))
	}

	existing, err := clients.DashboardClient.GetDashboard(clients.Ctx, args)
	if err != nil {
		return diag.Errorf(" Getting dashboard with ID: %s, %+v", id, err)
	}

	updateArgs := dashboardextras.UpdateDashboardArgs{
		Project:   converter.String(d.Get("project_id").(string)),
		Dashboard: existing,
	}

	if v, ok := d.GetOk("team_id"); ok {
		updateArgs.Team = converter.String(v.(string))
	}

	if d.HasChange("name") {
		updateArgs.Dashboard.Name = converter.String(d.Get("name").(string))
	}

	if d.HasChange("description") {
		updateArgs.Dashboard.Description = converter.String(d.Get("description").(string))
	}

	if d.HasChange("refresh_interval") {
		updateArgs.Dashboard.RefreshInterval = converter.Int(d.Get("refresh_interval").(int))
	}

	_, err = clients.DashboardClientExtra.UpdateDashboard(clients.Ctx, updateArgs)
	if err != nil {
		return diag.Errorf(" Updating dashboard with ID: %s. Error detail: %+v", id, err)
	}
	return resourceDashboardRead(clients.Ctx, d, m)
}

func resourceDashboardDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	id := d.Id()
	dashboardId, err := uuid.Parse(id)
	if err != nil {
		return diag.Errorf(" Parsing dashboard ID: %+v", err)
	}

	params := dashboard.DeleteDashboardArgs{
		Project:     converter.String(d.Get("project_id").(string)),
		DashboardId: &dashboardId,
	}

	if v, ok := d.GetOk("team_id"); ok {
		params.Team = converter.String(v.(string))
	}

	err = clients.DashboardClient.DeleteDashboard(clients.Ctx, params)

	if err != nil {
		var wrapperErr azuredevops.WrappedError
		if errors.As(err, &wrapperErr) {
			if wrapperErr.ErrorCode != nil && *wrapperErr.ErrorCode == 0 && wrapperErr.Message != nil &&
				strings.Contains(*wrapperErr.Message, "VS402433") {
				return diag.Errorf(" You can't delete the last dashboard of a team. You can delete the team "+
					"directly. \nProject ID: %s. \nTeam ID: %s \nDashboard ID: %s \nError detail: %+v", *params.Project, id, *params.Team, err)
			}
		}
		return diag.Errorf(" Deleting dashboard with id %s: %+v", id, err)
	}

	d.SetId("")
	return nil
}
