package build

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/pipelines"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/model"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/validate"
)

const (
	bdVariable              = "variable"
	bdVariableName          = "name"
	bdVariableValue         = "value"
	bdSecretVariableValue   = "secret_value"
	bdVariableIsSecret      = "is_secret"
	bdVariableAllowOverride = "allow_override"
)

// ResourceBuildDefinition schema and implementation for build definition resource
func ResourceBuildDefinition() *schema.Resource {
	filterSchema := map[string]*schema.Schema{
		"include": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
		},
		"exclude": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}

	branchFilter := &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: filterSchema,
		},
	}

	pathFilter := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: filterSchema,
		},
	}

	return &schema.Resource{
		CreateContext: resourceBuildDefinitionCreate,
		ReadContext:   resourceBuildDefinitionRead,
		UpdateContext: resourceBuildDefinitionUpdate,
		DeleteContext: resourceBuildDefinitionDelete,
		Importer:      tfhelper.ImportProjectQualifiedResource(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"revision": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      `\`,
				ValidateFunc: validate.Path,
			},
			"variable_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntAtLeast(1),
				},
			},
			bdVariable: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						bdVariableName: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						bdVariableValue: {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						bdSecretVariableValue: {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
							Default:   "",
						},
						bdVariableIsSecret: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						bdVariableAllowOverride: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"agent_pool_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Azure Pipelines",
			},
			"repository": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repo_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"repo_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(model.RepoTypeValues.GitHub),
								string(model.RepoTypeValues.TfsGit),
								string(model.RepoTypeValues.Bitbucket),
								string(model.RepoTypeValues.GitHubEnterprise),
								string(model.RepoTypeValues.OtherGit),
							}, false),
						},
						"yml_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"branch_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "master",
						},
						"service_connection_id": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"github_enterprise_url": {
							Type:          schema.TypeString,
							Optional:      true,
							Default:       "",
							ConflictsWith: []string{"repository.0.url"},
						},
						"url": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"repository.0.github_enterprise_url"},
						},
						"report_build_status": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"ci_trigger": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use_yaml": {
							Type:          schema.TypeBool,
							Optional:      true,
							Default:       false,
							ConflictsWith: []string{"ci_trigger.0.override"},
						},
						"override": {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"batch": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"branch_filter": branchFilter,
									"max_concurrent_builds_per_branch": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  1,
									},
									"path_filter": pathFilter,
									"polling_interval": {
										Type:     schema.TypeInt,
										Optional: true,
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
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use_yaml": {
							Type:          schema.TypeBool,
							Optional:      true,
							Default:       false,
							ConflictsWith: []string{"pull_request_trigger.0.override"},
						},
						"initial_branch": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Managed by Terraform",
						},
						"override": {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_cancel": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"branch_filter": branchFilter,
									"path_filter":   pathFilter,
								},
							},
						},
						"forks": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"share_secrets": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"comment_required": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"All", "NonTeamMembers"}, false),
						},
					},
				},
			},
			"build_completion_trigger": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"build_definition_id": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"branch_filter": branchFilter,
					},
				},
			},
			"agent_specification": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"job_authorization_scope": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(build.BuildAuthorizationScopeValues.Project),
					string(build.BuildAuthorizationScopeValues.ProjectCollection),
				}, false),
				Default: string(build.BuildAuthorizationScopeValues.ProjectCollection),
			},
			"jobs": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"ref_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"condition": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"dependencies": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scope": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
						"target": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(model.PipelineJobTypeValues.AgentJob),
											string(model.PipelineJobTypeValues.AgentlessJob),
										}, false),
									},
									"execution_options": {
										Type:     schema.TypeList,
										Required: true,
										MinItems: 1,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(model.JobExecutionOptionsTypeValues.None),
														string(model.JobExecutionOptionsTypeValues.MultiConfiguration),
														string(model.JobExecutionOptionsTypeValues.MultiAgent),
													}, false),
												},
												"max_concurrency": { // needs to be set when executionOptions type is: Multi-Configuration
													Type:         schema.TypeInt,
													Optional:     true,
													Computed:     true,
													ValidateFunc: validation.IntBetween(1, 99),
												},
												"multipliers": { // required when executionOptions type: Multi-Configuration
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"continue_on_error": { // required when executionOptions is: 1, or 2
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
									"demands": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"job_timeout_in_minutes": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 1000000000),
						},
						"job_cancel_timeout_in_minutes": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 60),
						},
						"job_authorization_scope": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(build.BuildAuthorizationScopeValues.Project),
								string(build.BuildAuthorizationScopeValues.ProjectCollection),
							}, false),
							Default: string(build.BuildAuthorizationScopeValues.ProjectCollection),
						},
						"allow_scripts_auth_access_option": { // available when job type is AgentJob(1)
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"schedules": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"branch_filter": branchFilter,
						"days_to_build": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}, false),
							},
						},
						"schedule_only_with_changes": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"start_hours": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 23),
						},
						"start_minutes": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 59),
						},
						"time_zone": {
							Optional:     true,
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice(TimeZones, false),
							Default:      "(UTC) Coordinated Universal Time",
						},
						"schedule_job_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"features": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"skip_first_run": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"queue_status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "enabled",
				ValidateFunc: validation.StringInSlice([]string{
					string(build.DefinitionQueueStatusValues.Enabled),
					string(build.DefinitionQueueStatusValues.Paused),
					string(build.DefinitionQueueStatusValues.Disabled),
				}, false),
			},
		},
	}
}

func resourceBuildDefinitionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	err := validateServiceConnectionIDExistsIfNeeded(d)
	if err != nil {
		return diag.FromErr(err)
	}
	buildDefinition, projectID, err := expandBuildDefinition(d, m)
	if err != nil {
		return diag.Errorf(" Creating Build Definition: %+v", err)
	}

	createdBuildDefinition, err := clients.BuildClient.CreateDefinition(clients.Ctx, build.CreateDefinitionArgs{
		Definition: buildDefinition,
		Project:    &projectID,
	})
	if err != nil {
		return diag.Errorf(" Creating Build Definition: %+v", err)
	}

	var diags diag.Diagnostics = nil
	features := buildDefinitionFeatures(d)
	if len(features) != 0 {
		if v, ok := features["skip_first_run"]; ok {
			if skipFirstRun := v.(bool); !skipFirstRun {
				// trigger the first run
				repo := d.Get("repository").([]interface{})[0].(map[string]interface{})
				branchName := repo["branch_name"].(string)

				branchName = strings.TrimPrefix(branchName, "refs/heads/")

				_, err := clients.PipelinesClient.RunPipeline(clients.Ctx, pipelines.RunPipelineArgs{
					Project:    converter.String(projectID),
					PipelineId: createdBuildDefinition.Id,
					RunParameters: &pipelines.RunPipelineParameters{
						Resources: &pipelines.RunResourcesParameters{
							Repositories: &map[string]pipelines.RepositoryResourceParameters{
								"self": {
									RefName: converter.String("refs/heads/" + branchName),
								},
							},
						},
					},
				})
				if err != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "First run of build definition failed, nothing to trigger",
						Detail:   fmt.Sprintf("Received error: %s\n Try initializing the repository with a valid build definition file", err),
					})
				}
			}
		}
	}

	d.SetId(strconv.Itoa(*createdBuildDefinition.Id))

	readDiag := resourceBuildDefinitionRead(ctx, d, m)
	if readDiag != nil {
		return readDiag
	} else {
		return diags
	}
}

func resourceBuildDefinitionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	projectID, buildDefinitionID, err := tfhelper.ParseProjectIDAndResourceID(d)
	if err != nil {
		return diag.FromErr(err)
	}

	buildDefinition, err := clients.BuildClient.GetDefinition(clients.Ctx, build.GetDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &buildDefinitionID,
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	return diag.FromErr(flattenBuildDefinition(d, buildDefinition, projectID))
}

func resourceBuildDefinitionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	err := validateServiceConnectionIDExistsIfNeeded(d)
	if err != nil {
		return diag.FromErr(err)
	}
	buildDefinition, projectID, err := expandBuildDefinition(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = clients.BuildClient.UpdateDefinition(m.(*client.AggregatedClient).Ctx, build.UpdateDefinitionArgs{
		Definition:   buildDefinition,
		Project:      &projectID,
		DefinitionId: buildDefinition.Id,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceBuildDefinitionRead(ctx, d, m)
}

func resourceBuildDefinitionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	projectID, buildDefinitionID, err := tfhelper.ParseProjectIDAndResourceID(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = clients.BuildClient.DeleteDefinition(m.(*client.AggregatedClient).Ctx, build.DeleteDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &buildDefinitionID,
	})
	return diag.FromErr(err)
}

func flattenBuildDefinition(d *schema.ResourceData, buildDefinition *build.BuildDefinition, projectID string) error {
	d.Set("project_id", projectID)
	d.Set("name", *buildDefinition.Name)
	d.Set("path", *buildDefinition.Path)

	repo, err := flattenRepository(buildDefinition)
	if err != nil {
		return err
	}
	d.Set("repository", repo)

	if buildDefinition.Queue != nil && buildDefinition.Queue.Pool != nil {
		d.Set("agent_pool_name", *buildDefinition.Queue.Pool.Name)
	}

	d.Set("variable_groups", flattenVariableGroups(buildDefinition))
	d.Set(bdVariable, flattenBuildVariables(d, buildDefinition))

	if buildDefinition.Triggers != nil {
		triggers := flattenTriggers(buildDefinition.Triggers)

		if triggers[build.DefinitionTriggerTypeValues.ContinuousIntegration] != nil {
			d.Set("ci_trigger", triggers[build.DefinitionTriggerTypeValues.ContinuousIntegration])
		}

		if triggers[build.DefinitionTriggerTypeValues.PullRequest] != nil {
			d.Set("pull_request_trigger", triggers[build.DefinitionTriggerTypeValues.PullRequest])
		}

		if triggers[build.DefinitionTriggerTypeValues.Schedule] != nil {
			d.Set("schedules", triggers[build.DefinitionTriggerTypeValues.Schedule])
		}

		if triggers[build.DefinitionTriggerTypeValues.BuildCompletion] != nil {
			d.Set("build_completion_trigger", triggers[build.DefinitionTriggerTypeValues.BuildCompletion])
		}
	}

	if buildDefinition.Process != nil {
		pipeJobs, err := flattenBuildDefinitionJobs(buildDefinition.Process)
		if err != nil {
			return fmt.Errorf("Flattening pipeline jobs: %+v", err)
		}

		err = d.Set("jobs", pipeJobs)
		if err != nil {
			return fmt.Errorf("Setting build definition jobs: %+v", err)
		}

		// set agent identifier
		if processMap, ok := buildDefinition.Process.(map[string]interface{}); ok {
			if target, ok := processMap["target"]; ok {
				if spec, ok := target.(map[string]interface{})["agentSpecification"]; ok {
					if agentIdentifier, ok := spec.(map[string]interface{})["identifier"]; ok {
						d.Set("agent_specification", agentIdentifier.(string))
					}
				}
			}
		}
	}

	revision := 0
	if buildDefinition.Revision != nil {
		revision = *buildDefinition.Revision
	}

	d.Set("job_authorization_scope", buildDefinition.JobAuthorizationScope)

	d.Set("revision", revision)
	d.Set("queue_status", *buildDefinition.QueueStatus)
	return nil
}

func flattenBuildDefinitionJobs(input interface{}) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	result := make([]interface{}, 0)
	if v, ok := input.(map[string]interface{}); ok {
		if phases, ok := v["phases"]; ok {
			var jobs []model.PipelineJob

			v, err := json.Marshal(phases)
			if err != nil {
				return nil, fmt.Errorf("Get pipeline jobs byte array: %+v", err)
			}

			err = json.Unmarshal(v, &jobs)
			if err != nil {
				return nil, fmt.Errorf("Convert Pipelins Jobs to PipelineJob: %+v", err)
			}

			for _, job := range jobs {
				var dependencyMap []map[string]interface{}
				if job.Dependencies != nil {
					for _, dependency := range *job.Dependencies {
						dependencyMap = append(dependencyMap, map[string]interface{}{
							"scope": dependency.Scope,
						})
					}
				}

				var targetMap map[string]interface{}
				if job.Target != nil {
					targetMap = map[string]interface{}{
						"demands": job.Target.Demands,
						"type":    model.PipelineJobTypeValueTypeMap[*job.Target.Type],
						"execution_options": []interface{}{
							map[string]interface{}{
								"type": model.JobExecutionOptionsTypeValues.None,
							},
						},
					}

					if job.Target.ExecutionOptions != nil {
						execOptionsMap := map[string]interface{}{
							"continue_on_error": job.Target.ExecutionOptions.ContinueOnError,
							"type":              model.JobExecutionOptionsTypValueTypeMap[*job.Target.ExecutionOptions.Type],
							"max_concurrency":   job.Target.ExecutionOptions.MaxConcurrency,
						}
						if job.Target.ExecutionOptions.Multipliers != nil && len(*job.Target.ExecutionOptions.Multipliers) > 0 {
							execOptionsMap["multipliers"] = (*job.Target.ExecutionOptions.Multipliers)[0]
						}

						targetMap["execution_options"] = []interface{}{execOptionsMap}
					}
				}

				jobConfig := map[string]interface{}{
					"name":                             job.Name,
					"ref_name":                         job.RefName,
					"condition":                        job.Condition,
					"allow_scripts_auth_access_option": job.Target.AllowScriptsAuthAccessOption,
					"job_timeout_in_minutes":           job.JobTimeoutInMinutes,
					"job_cancel_timeout_in_minutes":    job.JobCancelTimeoutInMinutes,
					"job_authorization_scope":          job.JobAuthorizationScope,
					"dependencies":                     dependencyMap,
					"target":                           []interface{}{targetMap},
				}

				result = append(result, jobConfig)
			}
		}
	}

	return result, nil
}

func flattenBuildVariables(d *schema.ResourceData, buildDefinition *build.BuildDefinition) interface{} {
	if buildDefinition.Variables == nil {
		return nil
	}
	variables := make([]map[string]interface{}, len(*buildDefinition.Variables))

	index := 0
	for varName, varVal := range *buildDefinition.Variables {
		var variable map[string]interface{}

		isSecret := converter.ToBool(varVal.IsSecret, false)
		variable = map[string]interface{}{
			bdVariableName:          varName,
			bdVariableValue:         converter.ToString(varVal.Value, ""),
			bdVariableIsSecret:      isSecret,
			bdVariableAllowOverride: converter.ToBool(varVal.AllowOverride, false),
		}

		// read secret variable from state if exist
		if isSecret {
			if stateVal := tfhelper.FindMapInSetWithGivenKeyValue(d, bdVariable, bdVariableName, varName); stateVal != nil {
				variable = stateVal
			}
		}
		variables[index] = variable
		index++
	}

	return variables
}

func flattenVariableGroups(buildDefinition *build.BuildDefinition) []int {
	if buildDefinition.VariableGroups == nil {
		return nil
	}

	variableGroups := make([]int, len(*buildDefinition.VariableGroups))

	for i, variableGroup := range *buildDefinition.VariableGroups {
		variableGroups[i] = *variableGroup.Id
	}

	return variableGroups
}

func flattenRepository(buildDefinition *build.BuildDefinition) (interface{}, error) {
	yamlFilePath := ""
	githubEnterpriseUrl := ""

	// The process member can be of many types -- the only typing information
	// available from the compiler is `interface{}` so we can probe for known
	// implementations
	if processMap, ok := buildDefinition.Process.(map[string]interface{}); ok {
		if v, exist := processMap["yamlFilename"].(string); exist {
			yamlFilePath = v
		}
	}
	if yamlProcess, ok := buildDefinition.Process.(*build.YamlProcess); ok {
		yamlFilePath = *yamlProcess.YamlFilename
	}

	// Set github_enterprise_url value from buildDefinition.Repository URL
	if strings.EqualFold(*buildDefinition.Repository.Type, string(model.RepoTypeValues.GitHubEnterprise)) {
		repoUrl, err := url.Parse(*buildDefinition.Repository.Url)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse repository URL: %+v ", err)
		}
		githubEnterpriseUrl = fmt.Sprintf("%s://%s", repoUrl.Scheme, repoUrl.Host)
	}

	repo := []map[string]interface{}{{
		"yml_path":              yamlFilePath,
		"repo_id":               *buildDefinition.Repository.Id,
		"repo_type":             *buildDefinition.Repository.Type,
		"branch_name":           *buildDefinition.Repository.DefaultBranch,
		"github_enterprise_url": githubEnterpriseUrl,
		"url":                   *buildDefinition.Repository.Url,
	}}

	if buildDefinition.Repository != nil && buildDefinition.Repository.Properties != nil {
		if connectionID, ok := (*buildDefinition.Repository.Properties)["connectedServiceId"]; ok {
			repo[0]["service_connection_id"] = connectionID
		}

		if buildStatus, ok := (*buildDefinition.Repository.Properties)["reportBuildStatus"]; ok {
			reportBuildStatus, err := strconv.ParseBool(buildStatus)
			if err != nil {
				return nil, fmt.Errorf("Unable parse Repository build status. Error: %+v", err)
			}
			repo[0]["report_build_status"] = reportBuildStatus
		}
	}
	return repo, nil
}

func flattenBuildDefinitionBranchOrPathFilter(m []interface{}) []interface{} {
	var include []string
	var exclude []string

	for _, v := range m {
		if v2, ok := v.(string); ok {
			if strings.HasPrefix(v2, "-") {
				exclude = append(exclude, strings.TrimPrefix(v2, "-"))
			} else if strings.HasPrefix(v2, "+") {
				include = append(include, strings.TrimPrefix(v2, "+"))
			}
		}
	}

	return []interface{}{
		map[string]interface{}{
			"include": include,
			"exclude": exclude,
		},
	}
}

func flattenBuildDefinitionContinuousIntegrationTrigger(m interface{}, isYaml bool) interface{} {
	if ms, ok := m.(map[string]interface{}); ok {
		f := map[string]interface{}{
			"use_yaml": isYaml,
		}
		if !isYaml {
			f["override"] = []map[string]interface{}{{
				"batch":                            ms["batchChanges"],
				"branch_filter":                    flattenBuildDefinitionBranchOrPathFilter(ms["branchFilters"].([]interface{})),
				"max_concurrent_builds_per_branch": ms["maxConcurrentBuildsPerBranch"],
				"polling_interval":                 ms["pollingInterval"],
				"polling_job_id":                   ms["pollingJobId"],
				"path_filter":                      flattenBuildDefinitionBranchOrPathFilter(ms["pathFilters"].([]interface{})),
			}}
		}
		return f
	}
	return nil
}

func flattenBuildDefinitionPullRequestTrigger(m interface{}, isYaml bool) interface{} {
	if ms, ok := m.(map[string]interface{}); ok {
		forks := ms["forks"].(map[string]interface{})
		isCommentRequired := ms["isCommentRequiredForPullRequest"].(bool)
		isCommentRequiredNonTeam := ms["requireCommentsForNonTeamMembersOnly"].(bool)

		var commentRequired string
		if isCommentRequired {
			commentRequired = "All"
		}
		if isCommentRequired && isCommentRequiredNonTeam {
			commentRequired = "NonTeamMembers"
		}

		branchFilters := ms["branchFilters"].([]interface{})
		var initialBranch string
		if len(branchFilters) > 0 {
			initialBranch = strings.TrimPrefix(branchFilters[0].(string), "+")
		}

		f := map[string]interface{}{
			"use_yaml":         isYaml,
			"initial_branch":   initialBranch,
			"comment_required": commentRequired,
			"forks": []map[string]interface{}{{
				"enabled":       forks["enabled"],
				"share_secrets": forks["allowSecrets"],
			}},
		}
		if !isYaml {
			f["override"] = []map[string]interface{}{{
				"auto_cancel":   ms["autoCancel"],
				"branch_filter": flattenBuildDefinitionBranchOrPathFilter(branchFilters),
				"path_filter":   flattenBuildDefinitionBranchOrPathFilter(ms["pathFilters"].([]interface{})),
			}}
		}
		return f
	}
	return nil
}

func flattenBuildDefinitionScheduleTrigger(ms map[string]interface{}) []interface{} {
	schedulesResp := ms["schedules"].([]interface{})
	schedules := make([]interface{}, 0)
	for _, schedule := range schedulesResp {
		schedule := schedule.(map[string]interface{})
		branchFilter := flattenBuildDefinitionBranchOrPathFilter(schedule["branchFilters"].([]interface{}))
		scheduleConfig := map[string]interface{}{
			"branch_filter":              branchFilter,
			"schedule_only_with_changes": schedule["scheduleOnlyWithChanges"],
			"start_hours":                schedule["startHours"],
			"start_minutes":              schedule["startMinutes"],
			"time_zone":                  IDToTimeZones[schedule["timeZoneId"].(string)],
			"schedule_job_id":            schedule["scheduleJobId"],
		}

		days := schedule["daysToBuild"]
		switch day := days.(type) {
		case float64:
			scheduleConfig["days_to_build"] = DaysToDate(int(day))
		case string:
			scheduleConfig["days_to_build"] = DaysToDate(DaysToBuild[day])
		}
		schedules = append(schedules, scheduleConfig)
	}
	return schedules
}

func flattenBuildCompletionTrigger(buildCompletionTrigger map[string]interface{}) interface{} {
	buildId := buildCompletionTrigger["definition"].(map[string]interface{})["id"].(float64)
	triggerConfig := map[string]interface{}{
		"branch_filter":       flattenBuildDefinitionBranchOrPathFilter(buildCompletionTrigger["branchFilters"].([]interface{})),
		"build_definition_id": buildId,
	}
	return triggerConfig
}

func flattenTriggers(m *[]interface{}) map[build.DefinitionTriggerType][]interface{} {
	buildTriggers := map[build.DefinitionTriggerType][]interface{}{}
	for _, ds := range *m {
		trigger := ds.(map[string]interface{})
		triggerType := trigger["triggerType"].(string)
		if strings.EqualFold(triggerType, string(build.DefinitionTriggerTypeValues.ContinuousIntegration)) {
			isYaml := false
			if val, ok := trigger["settingsSourceType"]; ok {
				isYaml = int(val.(float64)) == 2
			}
			buildTriggers[build.DefinitionTriggerTypeValues.ContinuousIntegration] = []interface{}{flattenBuildDefinitionContinuousIntegrationTrigger(trigger, isYaml)}
		}
		if strings.EqualFold(triggerType, string(build.DefinitionTriggerTypeValues.PullRequest)) {
			isYaml := false
			if val, ok := trigger["settingsSourceType"]; ok {
				isYaml = int(val.(float64)) == 2
			}
			buildTriggers[build.DefinitionTriggerTypeValues.PullRequest] = []interface{}{flattenBuildDefinitionPullRequestTrigger(trigger, isYaml)}
		}
		if strings.EqualFold(triggerType, string(build.DefinitionTriggerTypeValues.Schedule)) {
			buildTriggers[build.DefinitionTriggerTypeValues.Schedule] = flattenBuildDefinitionScheduleTrigger(trigger)
		}
		if strings.EqualFold(triggerType, string(build.DefinitionTriggerTypeValues.BuildCompletion)) {
			if _, ok := buildTriggers[build.DefinitionTriggerTypeValues.BuildCompletion]; !ok {
				buildTriggers[build.DefinitionTriggerTypeValues.BuildCompletion] = []interface{}{flattenBuildCompletionTrigger(trigger)}
			} else {
				buildTriggers[build.DefinitionTriggerTypeValues.BuildCompletion] = append(
					buildTriggers[build.DefinitionTriggerTypeValues.BuildCompletion],
					flattenBuildCompletionTrigger(trigger))
			}
		}
	}
	return buildTriggers
}

func expandBuildDefinitionBranchOrPathFilter(d map[string]interface{}) []interface{} {
	include := tfhelper.ExpandStringSet(d["include"].(*schema.Set))
	exclude := tfhelper.ExpandStringSet(d["exclude"].(*schema.Set))
	m := make([]interface{}, len(include)+len(exclude))
	i := 0
	for _, v := range include {
		m[i] = "+" + v
		i++
	}
	for _, v := range exclude {
		m[i] = "-" + v
		i++
	}
	return m
}

func expandBuildDefinitionBranchOrPathFilterList(d []interface{}) [][]interface{} {
	vs := make([][]interface{}, 0, len(d))
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandBuildDefinitionBranchOrPathFilter(val))
		}
	}
	return vs
}

func expandBuildDefinitionBranchOrPathFilterSet(configured *schema.Set) []interface{} {
	d2 := expandBuildDefinitionBranchOrPathFilterList(configured.List())
	if len(d2) != 1 {
		return nil
	}
	return d2[0]
}

func expandBuildDefinitionFork(d map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"allowSecrets": d["share_secrets"].(bool),
		"enabled":      d["enabled"].(bool),
	}
}

func expandBuildDefinitionForkList(d []interface{}) []map[string]interface{} {
	vs := make([]map[string]interface{}, 0, len(d))
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandBuildDefinitionFork(val))
		}
	}
	return vs
}

func expandBuildDefinitionForkListFirstOrNil(d []interface{}) map[string]interface{} {
	d2 := expandBuildDefinitionForkList(d)
	if len(d2) != 1 {
		return nil
	}
	return d2[0]
}

func expandBuildDefinitionManualPullRequestTrigger(d map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"branchFilters": expandBuildDefinitionBranchOrPathFilterSet(d["branch_filter"].(*schema.Set)),
		"pathFilters":   expandBuildDefinitionBranchOrPathFilterSet(d["path_filter"].(*schema.Set)),
		"autoCancel":    d["auto_cancel"].(bool),
	}
}

func expandBuildDefinitionManualPullRequestTriggerList(d []interface{}) []map[string]interface{} {
	vs := make([]map[string]interface{}, 0, len(d))
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandBuildDefinitionManualPullRequestTrigger(val))
		}
	}
	return vs
}

func expandBuildDefinitionManualPullRequestTriggerListFirstOrNil(d []interface{}) map[string]interface{} {
	d2 := expandBuildDefinitionManualPullRequestTriggerList(d)
	if len(d2) != 1 {
		return nil
	}
	return d2[0]
}

func expandBuildDefinitionManualContinuousIntegrationTrigger(d map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"batchChanges":                 d["batch"].(bool),
		"branchFilters":                expandBuildDefinitionBranchOrPathFilterSet(d["branch_filter"].(*schema.Set)),
		"maxConcurrentBuildsPerBranch": d["max_concurrent_builds_per_branch"].(int),
		"pathFilters":                  expandBuildDefinitionBranchOrPathFilterSet(d["path_filter"].(*schema.Set)),
		"triggerType":                  string(build.DefinitionTriggerTypeValues.ContinuousIntegration),
		"pollingInterval":              d["polling_interval"].(int),
	}
}

func expandBuildDefinitionManualContinuousIntegrationTriggerList(d []interface{}) []map[string]interface{} {
	vs := make([]map[string]interface{}, 0, len(d))
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandBuildDefinitionManualContinuousIntegrationTrigger(val))
		}
	}
	return vs
}

func expandBuildDefinitionManualContinuousIntegrationTriggerListFirstOrNil(d []interface{}) map[string]interface{} {
	d2 := expandBuildDefinitionManualContinuousIntegrationTriggerList(d)
	if len(d2) != 1 {
		return nil
	}
	return d2[0]
}

func expandBuildDefinitionTrigger(d map[string]interface{}, t build.DefinitionTriggerType, m interface{}, projectID string) (interface{}, error) {
	switch t {
	case build.DefinitionTriggerTypeValues.ContinuousIntegration:
		isYaml := d["use_yaml"].(bool)
		if isYaml {
			return map[string]interface{}{
				"batchChanges":                 false,
				"branchFilters":                []interface{}{},
				"maxConcurrentBuildsPerBranch": 1,
				"pathFilters":                  []interface{}{},
				"triggerType":                  string(t),
				"settingsSourceType":           float64(2),
			}, nil
		}
		return expandBuildDefinitionManualContinuousIntegrationTriggerListFirstOrNil(d["override"].([]interface{})), nil
	case build.DefinitionTriggerTypeValues.PullRequest:
		isYaml := d["use_yaml"].(bool)
		commentRequired := d["comment_required"].(string)
		vs := map[string]interface{}{
			"forks":                                expandBuildDefinitionForkListFirstOrNil(d["forks"].([]interface{})),
			"isCommentRequiredForPullRequest":      len(commentRequired) > 0,
			"requireCommentsForNonTeamMembersOnly": commentRequired == "NonTeamMembers",
			"triggerType":                          string(t),
		}
		if isYaml {
			vs["branchFilters"] = []interface{}{
				"+" + d["initial_branch"].(string),
			}
			vs["pathFilters"] = []interface{}{}
			vs["settingsSourceType"] = float64(2)
		} else {
			override := expandBuildDefinitionManualPullRequestTriggerListFirstOrNil(d["override"].([]interface{}))
			vs["branchFilters"] = override["branchFilters"]
			vs["pathFilters"] = override["pathFilters"]
			vs["autoCancel"] = override["autoCancel"]
		}
		return vs, nil
	case build.DefinitionTriggerTypeValues.Schedule:
		scheduleConfig := map[string]interface{}{
			"branchFilters":           expandBuildDefinitionBranchOrPathFilterSet(d["branch_filter"].(*schema.Set)),
			"scheduleOnlyWithChanges": d["schedule_only_with_changes"],
			"startHours":              d["start_hours"],
			"startMinutes":            d["start_minutes"],
			"timeZoneId":              TimeZoneToID[d["time_zone"].(string)],
			"scheduleJobId":           nil,
		}
		scheduleConfig["daysToBuild"] = DateToDays(d["days_to_build"].([]interface{}))
		return scheduleConfig, nil
	case build.DefinitionTriggerTypeValues.BuildCompletion:
		buildCompleteConfig := map[string]interface{}{
			"branchFilters":           expandBuildDefinitionBranchOrPathFilterSet(d["branch_filter"].(*schema.Set)),
			"requiresSuccessfulBuild": true,
			"triggerType":             string(t),
		}
		clients := m.(*client.AggregatedClient)

		buildDefinition, err := clients.BuildClient.GetDefinition(clients.Ctx, build.GetDefinitionArgs{
			Project:      &projectID,
			DefinitionId: converter.ToPtr(d["build_definition_id"].(int)),
		})
		if err != nil {
			return nil, fmt.Errorf("%+v", err)
		}
		buildCompleteConfig["definition"] = buildDefinition
		return buildCompleteConfig, nil
	}
	return nil, nil
}

func expandBuildDefinitionTriggerList(d []interface{}, t build.DefinitionTriggerType, m interface{}, projectID string) ([]interface{}, error) {
	vs := make([]interface{}, 0, len(d))
	for _, v := range d {
		val, ok := v.(map[string]interface{})
		if ok {
			trigger, err := expandBuildDefinitionTrigger(val, t, m, projectID)
			if err != nil {
				return nil, err
			}
			vs = append(vs, trigger)
		}
	}
	return vs, nil
}

func expandVariableGroups(d *schema.ResourceData) *[]build.VariableGroup {
	variableGroupsInterface := d.Get("variable_groups").(*schema.Set).List()
	variableGroups := make([]build.VariableGroup, len(variableGroupsInterface))

	for i, variableGroup := range variableGroupsInterface {
		variableGroups[i] = *buildVariableGroup(variableGroup.(int))
	}

	return &variableGroups
}

func expandVariables(d *schema.ResourceData) (*map[string]build.BuildDefinitionVariable, error) {
	variables := d.Get(bdVariable)
	if variables == nil {
		return nil, nil
	}

	variablesList := variables.(*schema.Set).List()
	if len(variablesList) == 0 {
		return nil, nil
	}

	expandedVars := map[string]build.BuildDefinitionVariable{}
	for _, variable := range variablesList {
		varAsMap := variable.(map[string]interface{})
		varName := varAsMap[bdVariableName].(string)

		if _, ok := expandedVars[varName]; ok {
			return nil, fmt.Errorf("Unexpectedly found duplicate variable with name %s", varName)
		}

		isSecret := converter.Bool(varAsMap[bdVariableIsSecret].(bool))
		var val *string

		if *isSecret {
			val = converter.String(varAsMap[bdSecretVariableValue].(string))
		} else {
			val = converter.String(varAsMap[bdVariableValue].(string))
		}
		expandedVars[varName] = build.BuildDefinitionVariable{
			AllowOverride: converter.Bool(varAsMap[bdVariableAllowOverride].(bool)),
			IsSecret:      isSecret,
			Value:         val,
		}
	}

	return &expandedVars, nil
}

func expandBuildDefinitionJobs(input []interface{}) (*[]model.PipelineJob, error) {
	if len(input) == 0 {
		return &[]model.PipelineJob{}, nil
	}

	result := make([]model.PipelineJob, 0, len(input))
	for _, jobConfig := range input {
		jobMap := jobConfig.(map[string]interface{})
		job := model.PipelineJob{
			Name:                      converter.String(jobMap["name"].(string)),
			RefName:                   converter.String(jobMap["ref_name"].(string)),
			Condition:                 converter.String(jobMap["condition"].(string)),
			JobTimeoutInMinutes:       converter.Int(jobMap["job_timeout_in_minutes"].(int)),
			JobCancelTimeoutInMinutes: converter.Int(jobMap["job_cancel_timeout_in_minutes"].(int)),
			JobAuthorizationScope:     converter.String(jobMap["job_authorization_scope"].(string)),
		}

		// dependencies
		if depConfig := jobMap["dependencies"].([]interface{}); len(depConfig) > 0 {
			var dependencies []model.JobDependency
			for _, dep := range depConfig {
				depMap := dep.(map[string]interface{})
				dependencies = append(dependencies, model.JobDependency{
					Scope: converter.String(depMap["scope"].(string)),
					Event: converter.String("Completed"),
				})
			}
			job.Dependencies = &dependencies
		}

		// job target
		targetConfig := jobMap["target"].([]interface{})[0]
		targetMap := targetConfig.(map[string]interface{})
		executionOptionsMap := targetMap["execution_options"].([]interface{})[0].(map[string]interface{})

		// set task type (AgentJob or AgentlessJob) and additional options(Allow scripts to access the OAuth token)
		jobType := model.PipelineJobTypeTypeValueMap[targetMap["type"].(string)]

		// execution options
		executionType := model.JobExecutionOptionsType(executionOptionsMap["type"].(string))
		var executeOptions model.JobExecutionOptions

		switch executionType {
		case model.JobExecutionOptionsTypeValues.None:
			executeOptions = model.JobExecutionOptions{
				Type: converter.ToPtr(0),
			}
		case model.JobExecutionOptionsTypeValues.MultiConfiguration:
			executeOptions = model.JobExecutionOptions{
				Type:            converter.ToPtr(1),
				ContinueOnError: converter.ToPtr(executionOptionsMap["continue_on_error"].(bool)),
				Multipliers:     &[]string{executionOptionsMap["multipliers"].(string)},
			}

			if jobType == 1 { // Agent Job
				if v, ok := executionOptionsMap["max_concurrency"]; !ok || v.(int) == 0 {
					return nil, fmt.Errorf("`max_concurrency` must be set when job is `AgentJob`")
				}
				executeOptions.MaxConcurrency = converter.ToPtr(executionOptionsMap["max_concurrency"].(int))
			}

			// TODO
			// if jobType == 2 { // Agentless Job
			//	if v, ok := executionOptionsMap["max_concurrency"]; ok && v.(int) > 0 {
			//		return nil, fmt.Errorf("`max_concurrency` must not be set when job is `AgentlessJob`")
			//	}
			// }

			// AgentlessJob(2)
			if jobType == 2 {
				executeOptions.MaxConcurrency = converter.ToPtr(50) // hard coded
			}
		case model.JobExecutionOptionsTypeValues.MultiAgent: // multi-agent, only available when job type is AgentJob
			if jobType == 2 {
				return nil, fmt.Errorf("`Multi-Agent` is not supported when job Type is `AgentlessJob`")
			}
			executeOptions = model.JobExecutionOptions{
				Type:            converter.ToPtr(2),
				ContinueOnError: converter.ToPtr(executionOptionsMap["continue_on_error"].(bool)),
			}

			if v, ok := executionOptionsMap["multipliers"]; ok {
				if len(v.(string)) > 0 {
					return nil, fmt.Errorf("`multipliers` must not be set when Execution Options Type is `Multi-Agent`")
				}
			}
			if jobType == 1 {
				if v, ok := executionOptionsMap["max_concurrency"]; !ok || v.(int) == 0 {
					return nil, fmt.Errorf("`max_concurrency` must be set when job is `AgentJob`")
				}
				executeOptions.MaxConcurrency = converter.ToPtr(executionOptionsMap["max_concurrency"].(int))
			}
		}

		target := model.JobTarget{
			Type:             converter.ToPtr(jobType),
			ExecutionOptions: &executeOptions,
		}

		// configurations only available when job type is AgentJob(1)
		if *target.Type == 1 {
			target.AllowScriptsAuthAccessOption = converter.ToPtr(jobMap["allow_scripts_auth_access_option"].(bool))
			demands := targetMap["demands"].([]interface{})
			if len(demands) > 0 {
				target.Demands = converter.ToPtr(tfhelper.ExpandStringList(demands))
			}
		} else {
			demands := targetMap["demands"].([]interface{})
			if len(demands) > 0 {
				return nil, fmt.Errorf("`demands` must not be set when Job Type is `AgentlessJob`")
			}
		}
		job.Target = &target

		result = append(result, job)
	}
	return &result, nil
}

func expandBuildDefinition(d *schema.ResourceData, meta interface{}) (*build.BuildDefinition, string, error) {
	projectID := d.Get("project_id").(string)
	repositories := d.Get("repository").([]interface{})

	// Note: If configured, this will be of length 1 based on the schema definition above.
	if len(repositories) != 1 {
		return nil, "", fmt.Errorf("Unexpectedly did not find repository metadata in the resource data")
	}

	repository := repositories[0].(map[string]interface{})

	repoID := repository["repo_id"].(string)
	repoType := repository["repo_type"].(string)
	repoURL := ""
	repoAPIURL := ""

	switch repoType {
	case string(model.RepoTypeValues.GitHub):
		repoURL = fmt.Sprintf("https://github.com/%s.git", repoID)
		repoAPIURL = fmt.Sprintf("https://api.github.com/repos/%s", repoID)
	case string(model.RepoTypeValues.Bitbucket):
		repoURL = fmt.Sprintf("https://bitbucket.org/%s.git", repoID)
		repoAPIURL = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s", repoID)
	case string(model.RepoTypeValues.GitHubEnterprise):
		githubEnterpriseURL := repository["github_enterprise_url"].(string)
		repoURL = fmt.Sprintf("%s/%s.git", githubEnterpriseURL, repoID)
		repoAPIURL = fmt.Sprintf("%s/api/v3/repos/%s", githubEnterpriseURL, repoID)
	case string(model.RepoTypeValues.OtherGit):
		repoURL = repository["url"].(string)
	}

	if strings.EqualFold(repoType, string(model.RepoTypeValues.OtherGit)) {
		if _, ok := repository["service_connection_id"]; !ok {
			return nil, "", fmt.Errorf("`repository.service_connection_id` must be set when `repoType` is `Git`")
		}

		if _, ok := repository["url"]; !ok {
			return nil, "", fmt.Errorf("`repository.service_connection_id` must be set when `repoType` is `Git`")
		}
	}

	var buildTriggers []any

	ciTriggers, err := expandBuildDefinitionTriggerList(
		d.Get("ci_trigger").([]interface{}),
		build.DefinitionTriggerTypeValues.ContinuousIntegration,
		meta,
		projectID,
	)
	if err != nil {
		return nil, "", err
	}
	buildTriggers = append(buildTriggers, ciTriggers...)

	pullRequestTriggers, err := expandBuildDefinitionTriggerList(
		d.Get("pull_request_trigger").([]interface{}),
		build.DefinitionTriggerTypeValues.PullRequest,
		meta,
		projectID,
	)
	if err != nil {
		return nil, "", err
	}
	buildTriggers = append(buildTriggers, pullRequestTriggers...)

	buildCompletionTriggers, err := expandBuildDefinitionTriggerList(
		d.Get("build_completion_trigger").([]interface{}),
		build.DefinitionTriggerTypeValues.BuildCompletion,
		meta,
		projectID,
	)
	if err != nil {
		return nil, "", err
	}
	buildTriggers = append(buildTriggers, buildCompletionTriggers...)

	schedules, err := expandBuildDefinitionTriggerList(
		d.Get("schedules").([]interface{}),
		build.DefinitionTriggerTypeValues.Schedule,
		meta,
		projectID,
	)
	if err != nil {
		return nil, "", err
	}
	if len(schedules) > 0 {
		scheduleTriggers := map[string]interface{}{
			"schedules":   schedules,
			"triggerType": string(build.DefinitionTriggerTypeValues.Schedule),
		}
		buildTriggers = append(buildTriggers, scheduleTriggers)
	}

	// Look for the ID. This may not exist if we are within the context of a "create" operation,
	// so it is OK if it is missing.
	buildDefinitionID, err := strconv.Atoi(d.Id())
	var buildDefinitionReference *int
	if err == nil {
		buildDefinitionReference = &buildDefinitionID
	} else {
		buildDefinitionReference = nil
	}

	variables, err := expandVariables(d)
	if err != nil {
		return nil, "", fmt.Errorf("Expanding varibles: %+v", err)
	}

	queueStatus := build.DefinitionQueueStatus(d.Get("queue_status").(string))

	buildDefinition := build.BuildDefinition{
		Id:       buildDefinitionReference,
		Name:     converter.String(d.Get("name").(string)),
		Path:     converter.String(d.Get("path").(string)),
		Revision: converter.Int(d.Get("revision").(int)),
		Repository: &build.BuildRepository{
			Url:           &repoURL,
			Id:            &repoID,
			Name:          &repoID,
			DefaultBranch: converter.String(repository["branch_name"].(string)),
			Type:          converter.String(repoType),
			Properties: &map[string]string{
				"connectedServiceId": repository["service_connection_id"].(string),
				"apiUrl":             repoAPIURL,
				"reportBuildStatus":  strconv.FormatBool(repository["report_build_status"].(bool)),
			},
		},
		Process: &build.YamlProcess{
			YamlFilename: converter.String(repository["yml_path"].(string)),
		},
		QueueStatus:    &queueStatus,
		Type:           &build.DefinitionTypeValues.Build,
		Quality:        &build.DefinitionQualityValues.Definition,
		VariableGroups: expandVariableGroups(d),
		Variables:      variables,
		Triggers:       &buildTriggers,
	}

	if agentPoolName, ok := d.GetOk("agent_pool_name"); ok {
		buildDefinition.Queue = &build.AgentPoolQueue{
			Name: converter.StringFromInterface(agentPoolName),
			Pool: &build.TaskAgentPoolReference{
				Name: converter.StringFromInterface(agentPoolName),
			},
		}
	}

	// other git need clone the repository information
	if repoType == string(model.RepoTypeValues.OtherGit) {
		(*buildDefinition.Repository.Properties)["fullName"] = "repository"
		(*buildDefinition.Repository.Properties)["cloneUrl"] = repoURL
		buildDefinition.Repository.Clean = converter.String("true")

		jobs, err := expandBuildDefinitionJobs(d.Get("jobs").([]interface{}))
		if err != nil {
			return nil, "", fmt.Errorf("Expanding jobs: %+v", err)
		}

		agentSpecification := d.Get("agent_specification").(string)
		if len(agentSpecification) == 0 {
			return nil, "", fmt.Errorf("Expanding jobs: `agent_specification` must be set when `repo_type` is `Git`")
		}

		buildDefinition.Process = map[string]interface{}{
			"type":   1,
			"phases": jobs,
			"target": map[string]interface{}{
				"agentSpecification": map[string]interface{}{
					"identifier": d.Get("agent_specification"),
				},
			},
		}
	}

	return &buildDefinition, projectID, nil
}

/**
 * certain types of build definitions require a service connection to run. This function
 * returns an error if a service connection was needed but not provided
 */
func validateServiceConnectionIDExistsIfNeeded(d *schema.ResourceData) error {
	repositories := d.Get("repository").([]interface{})
	repository := repositories[0].(map[string]interface{})

	repoType := repository["repo_type"].(string)
	serviceConnectionID := repository["service_connection_id"].(string)

	if strings.EqualFold(repoType, string(model.RepoTypeValues.Bitbucket)) && serviceConnectionID == "" {
		return errors.New("bitbucket repositories need a referenced service connection ID")
	}
	if strings.EqualFold(repoType, string(model.RepoTypeValues.GitHubEnterprise)) && serviceConnectionID == "" {
		return errors.New("GitHub Enterprise repositories need a referenced service connection ID")
	}
	return nil
}

func buildVariableGroup(id int) *build.VariableGroup {
	return &build.VariableGroup{
		Id: &id,
	}
}

func buildDefinitionFeatures(d *schema.ResourceData) map[string]interface{} {
	features := d.Get("features").([]interface{})
	if len(features) != 0 {
		return features[0].(map[string]interface{})
	}
	return nil
}
