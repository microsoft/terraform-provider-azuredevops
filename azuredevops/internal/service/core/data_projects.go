package core

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/datahelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// DataProjects schema and implementation for projects data source
func DataProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectsRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				ValidateFunc:     validation.StringIsNotWhiteSpace,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"state": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  core.ProjectStateValues.All,
				ValidateFunc: validation.StringInSlice([]string{
					string(core.ProjectStateValues.Deleting),
					string(core.ProjectStateValues.New),
					string(core.ProjectStateValues.WellFormed),
					string(core.ProjectStateValues.CreatePending),
					string(core.ProjectStateValues.All),
					string(core.ProjectStateValues.Unchanged),
					string(core.ProjectStateValues.Deleted),
				}, true),
			},

			"projects": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      getProjectHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceProjectsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	state := d.Get("state").(string)
	name := d.Get("name").(string)

	projects, err := getProjectsForStateAndName(clients, state, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf(" finding projects with state %s. Error: %v", state, err))
	}

	results := flattenProjectReferences(&projects)

	projectNames, err := datahelper.GetAttributeValues(results, "name")
	if err != nil {
		return diag.FromErr(fmt.Errorf(" failed to get list of project names: %v", err))
	}
	if len(projectNames) <= 0 && name != "" {
		projectNames = append(projectNames, name)
	}
	h := sha1.New()
	if _, err := h.Write([]byte(state + strings.Join(projectNames, "-"))); err != nil {
		return diag.FromErr(fmt.Errorf(" Unable to compute hash for project names: %v", err))
	}
	d.SetId("projects#" + base64.URLEncoding.EncodeToString(h.Sum(nil)))

	err = d.Set("projects", results)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func flattenProjectReferences(input *[]core.TeamProjectReference) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)

	for _, element := range *input {
		output := make(map[string]interface{})
		if element.Name != nil {
			output["name"] = *element.Name
		}

		if element.Id != nil {
			output["project_id"] = element.Id.String()
		}

		if element.Url != nil {
			output["project_url"] = *element.Url
		}

		if element.State != nil {
			output["state"] = string(*element.State)
		}

		results = append(results, output)
	}

	return results
}

func getProjectsForStateAndName(clients *client.AggregatedClient, projectState string, projectName string) ([]core.TeamProjectReference, error) {
	var projects []core.TeamProjectReference
	var currentToken string

	for hasMore := true; hasMore; {
		newProjects, latestToken, err := getProjectsWithContinuationToken(clients, projectState, currentToken)
		currentToken = latestToken
		if err != nil {
			return nil, err
		}

		if projectName != "" {
			for _, project := range newProjects {
				if strings.EqualFold(*project.Name, projectName) {
					return []core.TeamProjectReference{project}, nil
				}
			}
		} else {
			projects = append(projects, newProjects...)
		}
		hasMore = currentToken != ""
	}

	return projects, nil
}

func getProjectsWithContinuationToken(clients *client.AggregatedClient, projectState string, continuationToken string) ([]core.TeamProjectReference, string, error) {
	args := core.GetProjectsArgs{
		StateFilter: converter.ToPtr(core.ProjectState(projectState)),
	}
	if continuationToken != "" {
		token, err := strconv.Atoi(continuationToken)
		if err != nil {
			return nil, "", err
		}
		args.ContinuationToken = &token
	}

	response, err := clients.CoreClient.GetProjects(clients.Ctx, args)
	if err != nil {
		return nil, "", err
	}

	return response.Value, response.ContinuationToken, nil
}

func getProjectHash(v interface{}) int {
	return tfhelper.HashString(v.(map[string]interface{})["project_id"].(string))
}
