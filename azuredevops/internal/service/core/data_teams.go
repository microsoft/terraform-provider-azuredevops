package core

import (
	"fmt"
	"math/rand"

	"github.com/ahmetb/go-linq"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func DataTeams() *schema.Resource {
	return &schema.Resource{
		Read: dataTeamsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},
			"teams": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"administrators": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
							Set:      schema.HashString,
						},
						"members": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
							Set:      schema.HashString,
						},
					},
				},
			},
		},
	}
}

func dataTeamsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	var projectIDList []string
	data, ok := d.GetOk("project_id")
	if ok {
		projectIDList = []string{
			data.(string),
		}
	} else {
		projectList, err := getProjectsForStateAndName(clients, string(core.ProjectStateValues.All), "")
		if err != nil {
			return err
		}
		linq.From(projectList).
			Select(func(e interface{}) interface{} {
				return e.(core.TeamProjectReference).Id.String()
			}).
			ToSlice(&projectIDList)
	}

	result := make([]interface{}, 0)
	for _, projectID := range projectIDList {
		teamList, err := clients.CoreClient.GetTeams(clients.Ctx, core.GetTeamsArgs{
			ProjectId:      converter.String(projectID),
			Mine:           converter.Bool(false),
			ExpandIdentity: converter.Bool(false),
		})

		if err != nil {
			return err
		}

		if teamList == nil || len(*teamList) <= 0 {
			continue
		}

		teams := make([]interface{}, len(*teamList))
		for i, team := range *teamList {
			members, err := readTeamMembers(clients, &team)
			if err != nil {
				return err
			}
			administrators, err := readTeamAdministrators(d, clients, &team)
			if err != nil {
				return err
			}

			s := make(map[string]interface{})

			if v := team.ProjectId; v != nil {
				s["project_id"] = v.String()
			}
			if v := team.Id; v != nil {
				s["id"] = v.String()
			}
			if v := team.Name; v != nil {
				s["name"] = *v
			}
			if v := team.Description; v != nil {
				s["description"] = *v
			}
			if administrators != nil {
				s["administrators"] = administrators
			}
			if members != nil {
				s["members"] = members
			}

			teams[i] = s
		}

		result = append(result, teams...)
	}
	// The ID for this resource is meaningless so we can just assign a random ID
	d.SetId(fmt.Sprintf("%d", rand.Int()))

	if err := d.Set("teams", result); err != nil {
		return fmt.Errorf("Error setting `teams`: %+v", err)
	}

	return nil
}
