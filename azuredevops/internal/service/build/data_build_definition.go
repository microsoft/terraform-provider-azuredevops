package build

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/build"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/validate"
)

// DataBuildDefinition schema and implementation for Git repository data source
func DataBuildDefinition() *schema.Resource {
	filterSchema := map[string]*schema.Schema{
		"include": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
		},
		"exclude": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}

	branchFilter := &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: filterSchema,
		},
	}

	pathFilter := &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: filterSchema,
		},
	}

	return &schema.Resource{
		Read: dataSourceGitRepositoryRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      `\`,
				ValidateFunc: validate.Path,
			},
			"revision": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"variable_groups": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntAtLeast(1),
				},
			},
			bdVariable: {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						bdVariableName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						bdVariableValue: {
							Type:     schema.TypeString,
							Computed: true,
						},
						bdSecretVariableValue: {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						bdVariableIsSecret: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						bdVariableAllowOverride: {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"agent_pool_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repository": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"yml_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repo_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repo_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"branch_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"github_enterprise_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"report_build_status": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"ci_trigger": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use_yaml": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"override": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"batch": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"branch_filter": branchFilter,
									"max_concurrent_builds_per_branch": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"path_filter": pathFilter,
									"polling_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"polling_job_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"pull_request_trigger": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use_yaml": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"initial_branch": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"override": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_cancel": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"branch_filter": branchFilter,
									"path_filter":   pathFilter,
								},
							},
						},
						"forks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"share_secrets": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"comment_required": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"schedules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"branch_filter": branchFilter,
						"days_to_build": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}, false),
							},
						},
						"schedule_only_with_changes": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"start_hours": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"start_minutes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"time_zone": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"schedule_job_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"queue_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGitRepositoryRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	name := d.Get("name").(string)
	path := d.Get("path").(string)
	projectID := d.Get("project_id").(string)

	buildDefinitions, err := getBuildDefinitionsByNameAndProject(clients, name, path, projectID)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return fmt.Errorf("Build Definition with name %s does not exist in project %s in %s path", name, projectID, path)
		}
		return fmt.Errorf("Error finding build definitions. Error: %v", err)
	}
	if buildDefinitions == nil || 0 >= len(*buildDefinitions) {
		return fmt.Errorf("Build Definition with name %s does not exist in project %s in %s path", name, projectID, path)
	}
	if 1 < len(*buildDefinitions) {
		return fmt.Errorf("Multiple build definitions with name %s found in project %s", name, projectID)
	}

	flattenBuildDefinition(d, &(*buildDefinitions)[0], projectID)

	return nil
}

func getBuildDefinitionsByNameAndProject(clients *client.AggregatedClient, name string, path string, projectID string) (*[]build.BuildDefinition, error) {
	getArgs := build.GetDefinitionsArgs{
		Project: &projectID,
		Name:    converter.String(name),
	}

	if path != `\` {
		getArgs.Path = converter.String(path)
	}

	builds, err := clients.BuildClient.GetDefinitions(clients.Ctx, getArgs)
	if err != nil {
		return nil, err
	}
	var buildDefinitions []build.BuildDefinition
	for _, buildDefinition := range builds.Value {
		build, err := clients.BuildClient.GetDefinition(clients.Ctx, build.GetDefinitionArgs{
			Project:      &projectID,
			DefinitionId: buildDefinition.Id,
		})
		if err != nil {
			return nil, err
		}

		buildDefinitions = append(buildDefinitions, *build)
	}

	return &buildDefinitions, nil
}
