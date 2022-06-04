package release

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/release"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/webapi"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/model"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/validate"
)

/*
NOTE : https://www.terraform.io/docs/extend/writing-custom-providers.html
Due to the limitation of tf-11115 it is not possible to nest maps. So the workaround is to let only the innermost data structure be of the type TypeMap: in this case driver_options. The outer data structures are of TypeList which can only have one item.

TODO : Based on the info above any of the Min: 1, Max: 1, items will MOST LIKELY become TypeList instead of TypeSet.
TODO : Otherwise a custom schema.Schema.Set function will need to be created to identify them in the set.
This fixes the behaviour of apply. Otherwise apply will sometimes result in changes.
*/

func ResourceReleaseDefinition() *schema.Resource {

	id := &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}

	rank := &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}

	timeoutInMinutes := &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  0,
	}

	maxExecutionInMinutes := &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      1,
		ValidateFunc: validation.IntAtLeast(1),
	}

	variableGroups := &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type:         schema.TypeInt,
			ValidateFunc: validation.IntAtLeast(1),
		},
		Optional: true,
	}

	configurationVariableValue := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
		},
		"allow_override": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"is_secret": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}

	configurationVariables := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: configurationVariableValue,
		},
	}

	demand := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"value": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	demands := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: demand,
		},
	}

	artifactItems := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	buildArtifactDownloads := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"artifact_alias": {
					Type:     schema.TypeString,
					Required: true,
				},
				"include": artifactItems,
			},
		},
	}

	overrideInputs := &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
	}

	task := map[string]*schema.Schema{
		// NOTE : Due to limitations of ConflictsWith & ExactlyOneOf when using nested structures. https://github.com/hashicorp/terraform-plugin-sdk/issues/71
		// Task will be the only field supplied.
		"task": {
			Type:     schema.TypeString,
			Optional: true,
			ValidateFunc: func(i interface{}, k string) ([]string, []error) {
				v, ok := i.(string)
				if !ok {
					return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
				}
				var validID = regexp.MustCompile(`^.*?@\d+`)
				if match := validID.MatchString(v); !match {
					return nil, []error{fmt.Errorf("invalid task format is name@version")}
				}
				task := strings.Split(v, "@")
				taskName := task[0]
				if _, ok := taskagent.TaskNameToUUID[taskName]; !ok {
					return nil, []error{fmt.Errorf("unkown task %q", taskName)}
				}
				return nil, nil
			},
			// ConflictsWith: []string{"task_id", "version"},
		},
		"always_run": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"condition": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "succeeded()",
		},
		"continue_on_error": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"environment": {
			Type:     schema.TypeMap,
			Optional: true,
		},
		"display_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"override_inputs":    overrideInputs,
		"timeout_in_minutes": timeoutInMinutes,
		"inputs": {
			Type:     schema.TypeMap,
			Optional: true,
		},
		// NOTE : Due to limitations of ConflictsWith & ExactlyOneOf when using nested structures. https://github.com/hashicorp/terraform-plugin-sdk/issues/71
		// It is not possible to have task OR task_id/version
		// If ConflictsWith & ExactlyOneOf support * then this could be done.
		//"task_id": {
		//	Type:     schema.TypeString,
		//	Required: true,
		//	ExactlyOneOf: []string{"task_id", "task"},
		//},
		//"version": {
		//	Type:     schema.TypeString,
		//	Required: true,
		//	ExactlyOneOf: []string{"version", "task"},
		//},
	}

	tasks := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: task,
		},
	}

	releaseDefinitionDeployStep := &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id":   id,
				"task": tasks,
			},
		},
	}

	approvalOptions := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auto_triggered_and_previous_environment_approved_can_be_skipped": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"enforce_identity_revalidation": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"release_creator_can_be_approver": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"required_approver_count": {
					Type:     schema.TypeInt,
					Optional: true,
				},
			},
		},
	}

	releaseDefinitionGatesOptions := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"is_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"minimum_success_duration": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"sampling_interval": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"stabilization_time": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"timeout": {
					Type:     schema.TypeInt,
					Optional: true,
				},
			},
		},
	}

	releaseDefinitionGate := map[string]*schema.Schema{
		"task": tasks,
	}

	releaseDefinitionGates := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: releaseDefinitionGate,
		},
	}

	releaseDefinitionGatesStep := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id":            id,
				"gate":          releaseDefinitionGates,
				"gates_options": releaseDefinitionGatesOptions,
			},
		},
	}

	skipArtifactsDownload := &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}

	environmentOptions := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auto_link_work_items": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"badge_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"publish_deployment_status": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
				"pull_request_deployment_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}

	environmentExecutionPolicy := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"concurrency_count": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  1,
				},
				"queue_depth_count": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  0,
				},
			},
		},
	}

	schedule := map[string]*schema.Schema{
		"days_to_release": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(release.ScheduleDaysValues.All),
				string(release.ScheduleDaysValues.Friday),
				string(release.ScheduleDaysValues.Monday),
				string(release.ScheduleDaysValues.None),
				string(release.ScheduleDaysValues.Saturday),
				string(release.ScheduleDaysValues.Sunday),
				string(release.ScheduleDaysValues.Thursday),
				string(release.ScheduleDaysValues.Tuesday),
				string(release.ScheduleDaysValues.Wednesday),
			}, false),
		},
		"job_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"schedule_only_with_changes": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"start_hours": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"start_minutes": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"time_zone_id": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	schedules := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: schedule,
		},
	}

	releaseDefinitionProperties := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"definition_creation_source": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "ReleaseNew",
				},
				"integrate_jira_work_items": {
					Type:         schema.TypeBool,
					Optional:     true,
					Default:      false,
					AtLeastOneOf: []string{"properties.0.jira_service_endpoint_id"},
				},
				"integrate_boards_work_items": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"jira_service_endpoint_id": {
					Type:         schema.TypeString,
					Optional:     true,
					AtLeastOneOf: []string{"properties.0.integrate_jira_work_items"},
				},
			},
		},
	}

	environmentType := &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ValidateFunc: validation.StringInSlice([]string{
			string(DeploymentTypeValues.Production),
			string(DeploymentTypeValues.Staging),
			string(DeploymentTypeValues.Testing),
			string(DeploymentTypeValues.Development),
			string(DeploymentTypeValues.Unmapped),
		}, false),
	}

	releaseDefinitionEnvironmentProperties := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"boards_environment_type": environmentType,
				"link_boards_work_items": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"jira_environment_type": environmentType,
				"link_jira_work_items": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}

	releaseTrigger := map[string]*schema.Schema{
		"trigger_type": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(release.EnvironmentTriggerTypeValues.Undefined),
				string(release.EnvironmentTriggerTypeValues.DeploymentGroupRedeploy),
				string(release.EnvironmentTriggerTypeValues.RollbackRedeploy),
			}, false),
		},
	}

	releaseTriggers := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: releaseTrigger,
		},
	}

	environmentTrigger := map[string]*schema.Schema{
		"definition_environment_id": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"release_definition_id": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"trigger_content": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"trigger_type": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(release.EnvironmentTriggerTypeValues.Undefined),
				string(release.EnvironmentTriggerTypeValues.DeploymentGroupRedeploy),
				string(release.EnvironmentTriggerTypeValues.RollbackRedeploy),
			}, false),
		},
	}

	environmentTriggers := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: environmentTrigger,
		},
	}

	buildArtifact := map[string]*schema.Schema{
		"project_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"build_pipeline_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"latest": {
			Type:     schema.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"branch": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"tags": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
			ExactlyOneOf: []string{"build_artifact.0.latest", "build_artifact.0.specify"},
		},
		"specify": {
			Type:     schema.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"version": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
			ExactlyOneOf: []string{"build_artifact.0.latest", "build_artifact.0.specify"},
		},
		"alias": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"is_primary": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"is_retained": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
	buildArtifacts := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: buildArtifact,
		},
	}

	approval := &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id":   id,
				"rank": rank,
				"approver_id": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsUUID,
				},
			},
		},
	}

	releaseDefinitionApproval := &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"approval":           approval,
				"timeout_in_minutes": timeoutInMinutes,
			},
		},
		// RequiredWith: []string{"stage.$.pre_deploy_approval.0.approval", "stage.$.pre_deploy_approval.0.timeout_in_minutes", "stage.$.post_deploy_approval.0.approval", "stage.$.post_deploy_approval.0.timeout_in_minutes"}},
		// TODO : How to solve "ignore" if empty and new is isAutomated?
	}

	retentionPolicy := &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"days_to_keep": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  30,
				},
				"releases_to_keep": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  3,
				},
				"retain_build": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
			},
		},
	}

	allowScriptsToAccessOauthToken := &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}

	stages := &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id":   id,
				"rank": rank,
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"owner_id": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.IsUUID,
				},
				"variable":             configurationVariables,
				"variable_groups":      variableGroups,
				"pre_deploy_approval":  releaseDefinitionApproval,
				"post_deploy_approval": releaseDefinitionApproval,
				// NOTE : the ui only allows setting these to the same value for Post and Pre
				"approval_options": approvalOptions,

				"deploy_step": releaseDefinitionDeployStep,

				"job": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"agent": {
								Type:     schema.TypeList,
								Optional: true,
								MinItems: 1,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"name": {
											Type:     schema.TypeString,
											Required: true,
										},
										"override_inputs": overrideInputs,
										"demand":          demands,
										"rank":            rank,
										"agent_pool_hosted_azure_pipelines": {
											Type:     schema.TypeList,
											Optional: true,
											MinItems: 1,
											MaxItems: 1,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"agent_pool_id": {
														Type:     schema.TypeInt,
														Required: true,
													},
													"agent_specification": {
														Type:     schema.TypeString,
														Required: true,
														ValidateFunc: validation.StringInSlice([]string{
															"macOS-10.13",
															"macOS-10.14",
															"macOS-10.15",
															"ubuntu-16.04",
															"ubuntu-18.04",
															"ubuntu-20.04",
															"vs2015-win2012r2",
															"vs2017-win2016",
															"win1803",
															"windows-2019",
														}, false),
													},
												},
											},
										},
										"build_artifact_download": buildArtifactDownloads,
										"agent_pool_private": {
											Type:     schema.TypeList,
											Optional: true,
											MinItems: 1,
											MaxItems: 1,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"agent_pool_id": {
														Type:     schema.TypeString,
														Required: true,
													},
												},
											},
										},
										"timeout_in_minutes":            timeoutInMinutes,
										"max_execution_time_in_minutes": maxExecutionInMinutes,
										"condition": {
											Type:     schema.TypeString,
											Optional: true,
											Default:  "succeeded()",
										},
										"multi_configuration": {
											Type:     schema.TypeList,
											Optional: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"multipliers": {
														Type:        schema.TypeString,
														Required:    true,
														Description: "A list of comma separated configuration variables to use. These are defined on the Variables tab. For example, OperatingSystem, Browser will run the tasks for both variables.",
													},
													"number_of_agents": {
														Type:         schema.TypeInt,
														Required:     true,
														ValidateFunc: validation.IntAtLeast(1),
													},
													"continue_on_error": {
														Type:     schema.TypeBool,
														Optional: true,
														Default:  false,
													},
												},
											},
										},
										"multi_agent": {
											Type:     schema.TypeList,
											Optional: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"max_number_of_agents": {
														Type:         schema.TypeInt,
														Required:     true,
														ValidateFunc: validation.IntAtLeast(1),
													},
													"continue_on_error": {
														Type:     schema.TypeBool,
														Optional: true,
														Default:  false,
													},
												},
											},
										},
										"skip_artifacts_download":             skipArtifactsDownload,
										"allow_scripts_to_access_oauth_token": allowScriptsToAccessOauthToken,
										"task":                                tasks,
									},
								},
							},
							"deployment_group": {
								Type:     schema.TypeList,
								Optional: true,
								MinItems: 1,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"name": {
											Type:     schema.TypeString,
											Required: true,
										},
										"override_inputs": overrideInputs,
										"demand":          demands,
										"rank":            rank,
										"deployment_group_id": {
											Type:     schema.TypeInt,
											Required: true,
										},
										"tags": &model.TagsSchema,
										"multiple": {
											Type:     schema.TypeList,
											Optional: true,
											MinItems: 1,
											MaxItems: 1,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"max_targets_in_parallel": {
														Type:     schema.TypeInt,
														Required: true,
													},
												},
											},
										},
										"allow_scripts_to_access_oauth_token": allowScriptsToAccessOauthToken,
										"timeout_in_minutes":                  timeoutInMinutes,
										"max_execution_time_in_minutes":       maxExecutionInMinutes,
										"condition": {
											Type:     schema.TypeString,
											Required: true,
										},
										"task":                    tasks,
										"skip_artifacts_download": skipArtifactsDownload,
									},
								},
							},
							"agentless": {
								Type:     schema.TypeList,
								Optional: true,
								MinItems: 1,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"name": {
											Type:     schema.TypeString,
											Required: true,
										},
										"override_inputs":               overrideInputs,
										"rank":                          rank,
										"timeout_in_minutes":            timeoutInMinutes,
										"max_execution_time_in_minutes": maxExecutionInMinutes,
										"condition": {
											Type:     schema.TypeString,
											Optional: true,
											Default:  "succeeded()",
										},
										"multi_configuration": {
											Type:     schema.TypeList,
											Optional: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"multipliers": {
														Type:        schema.TypeString,
														Required:    true,
														Description: "A list of comma separated configuration variables to use. These are defined on the Variables tab. For example, OperatingSystem, Browser will run the tasks for both variables.",
													},
													"continue_on_error": {
														Type:     schema.TypeBool,
														Optional: true,
														Default:  false,
													},
												},
											},
										},
										"task": tasks,
									},
								},
							},
						},
					},
				},
				"retention_policy": retentionPolicy,

				// TODO : This is missing from the docs
				// "runOptions": runOptions,
				"environment_options": environmentOptions,
				"demand": {
					Type:       schema.TypeList,
					Optional:   true,
					Deprecated: "Use DeploymentInput.Demands instead",
					Elem: &schema.Resource{
						Schema: demand,
					},
				},
				"after_stage": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"stage_name": {
								Type:     schema.TypeString,
								Required: true,
							},
							"trigger_even_when_stages_partially_succeed": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  false,
							},
						},
					},
					// TODO : deeply nested validation
					// ConflictsWith: []string{"after_release"},
				},

				"after_release": {
					Type:     schema.TypeList,
					MinItems: 1,
					MaxItems: 1,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"event_name": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
					// TODO : deeply nested validation
					// ConflictsWith: []string{"after_stage"},
				},

				"artifact_filter": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"artifact_alias": {
								Type:     schema.TypeString,
								Required: true,
							},
							"include": {
								Type:     schema.TypeList,
								Optional: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"branch_name": {
											Type:     schema.TypeString,
											Required: true,
										},
										"tags": {
											Type:     schema.TypeSet,
											Optional: true,
											Elem: &schema.Schema{
												Type:         schema.TypeString,
												ValidateFunc: validation.NoZeroValues,
											},
										},
									},
								},
							},
							"exclude": {
								Type:     schema.TypeList,
								Optional: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"branch_name": {
											Type:     schema.TypeString,
											Required: true,
										},
									},
								},
							},
						},
					},
				},
				"execution_policy":     environmentExecutionPolicy,
				"schedules":            schedules,
				"properties":           releaseDefinitionEnvironmentProperties,
				"pre_deploy_gate":      releaseDefinitionGatesStep,
				"post_deploy_gate":     releaseDefinitionGatesStep,
				"environment_triggers": environmentTriggers,
				"badge_url": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}

	return &schema.Resource{
		Create: resourceReleaseDefinitionCreate,
		Read:   resourceReleaseDefinitionRead,
		Update: resourceReleaseDefinitionUpdate,
		Delete: resourceReleaseDefinitionDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, buildDefinitionID, err := tfhelper.ParseImportedID(d.Id())
				if err != nil {
					return nil, fmt.Errorf("error parsing the build definition ID from the Terraform resource data: %v", err)
				}
				d.Set("project_id", projectID)
				d.SetId(fmt.Sprintf("%d", buildDefinitionID))

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"revision": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "\\",
				ValidateFunc: validate.Path,
			},
			"variable_groups": variableGroups,
			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"variable": configurationVariables,
			"release_name_format": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Release-$(rev:r)",
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tags":       &model.TagsSchema,
			"properties": releaseDefinitionProperties,
			"comment": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				DiffSuppressFunc: func(_, _, _ string, _ *schema.ResourceData) bool { return true },
			},
			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stage":          stages,
			"build_artifact": buildArtifacts,
			"triggers":       releaseTriggers,
		},
	}
}

func flattenReleaseDefinition(d *schema.ResourceData, releaseDefinition *release.ReleaseDefinition, projectID string) {
	d.SetId(strconv.Itoa(*releaseDefinition.Id))
	d.Set("project_id", projectID)
	d.Set("name", *releaseDefinition.Name)
	d.Set("path", *releaseDefinition.Path)
	d.Set("variable_groups", *releaseDefinition.VariableGroups)
	d.Set("source", *releaseDefinition.Source)
	d.Set("description", converter.ToString(releaseDefinition.Description, ""))
	d.Set("variable", flattenReleaseDefinitionVariables(releaseDefinition.Variables))
	d.Set("release_name_format", *releaseDefinition.ReleaseNameFormat)
	d.Set("url", *releaseDefinition.Url)
	d.Set("is_deleted", *releaseDefinition.IsDeleted)
	d.Set("tags", *releaseDefinition.Tags)
	d.Set("properties", flattenReleaseDefinitionProperties(releaseDefinition.Properties))
	if releaseDefinition.Comment != nil {
		d.Set("comment", *releaseDefinition.Comment)
	}
	d.Set("created_on", releaseDefinition.CreatedOn.Time.Format(time.RFC3339))
	d.Set("modified_on", releaseDefinition.ModifiedOn.Time.Format(time.RFC3339))
	d.Set("stage", flattenReleaseDefinitionEnvironmentList(releaseDefinition.Environments))
	d.Set("build_artifact", flattenReleaseDefinitionArtifactsList(releaseDefinition.Artifacts, release.AgentArtifactTypeValues.Build))
	d.Set("triggers", flattenReleaseDefinitionTriggersList(releaseDefinition.Triggers))

	revision := 0
	if releaseDefinition.Revision != nil {
		revision = *releaseDefinition.Revision
	}
	d.Set("revision", revision)
}

func resourceReleaseDefinitionCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	releaseDefinition, projectID, err := expandReleaseDefinition(d)
	if err != nil {
		return fmt.Errorf("error creating resource Release Definition: %+v", err)
	}

	createdReleaseDefinition, err := createReleaseDefinition(clients, releaseDefinition, projectID)
	if err != nil {
		return fmt.Errorf("error creating resource Release Definition: %+v", err)
	}

	flattenReleaseDefinition(d, createdReleaseDefinition, projectID)
	return nil
}

func resourceReleaseDefinitionRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID, releaseDefinitionID, err := tfhelper.ParseProjectIDAndResourceID(d)

	if err != nil {
		return err
	}

	releaseDefinition, err := clients.ReleaseClient.GetReleaseDefinition(clients.Ctx, release.GetReleaseDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &releaseDefinitionID,
	})

	if err != nil {
		return err
	}

	flattenReleaseDefinition(d, releaseDefinition, projectID)
	return nil
}

func resourceReleaseDefinitionDelete(d *schema.ResourceData, m interface{}) error {
	if d.Id() == "" {
		return nil
	}

	clients := m.(*client.AggregatedClient)
	projectID, releaseDefinitionID, err := tfhelper.ParseProjectIDAndResourceID(d)
	if err != nil {
		return err
	}

	err = clients.ReleaseClient.DeleteReleaseDefinition(m.(*client.AggregatedClient).Ctx, release.DeleteReleaseDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &releaseDefinitionID,
	})

	return err
}

func resourceReleaseDefinitionUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	releaseDefinition, projectID, err := expandReleaseDefinition(d)
	if err != nil {
		return err
	}

	updatedReleaseDefinition, err := clients.ReleaseClient.UpdateReleaseDefinition(m.(*client.AggregatedClient).Ctx, release.UpdateReleaseDefinitionArgs{
		ReleaseDefinition: releaseDefinition,
		Project:           &projectID,
	})

	if err != nil {
		return err
	}

	flattenReleaseDefinition(d, updatedReleaseDefinition, projectID)
	return nil
}

func createReleaseDefinition(clients *client.AggregatedClient, releaseDefinition *release.ReleaseDefinition, project string) (*release.ReleaseDefinition, error) {
	createdBuild, err := clients.ReleaseClient.CreateReleaseDefinition(clients.Ctx, release.CreateReleaseDefinitionArgs{
		ReleaseDefinition: releaseDefinition,
		Project:           &project,
	})

	return createdBuild, err
}

func expandReleaseDefinition(d *schema.ResourceData) (*release.ReleaseDefinition, string, error) {
	projectID := d.Get("project_id").(string)

	// Look for the ID. This may not exist if we are within the context of a "create" operation,
	// so it is OK if it is missing.
	releaseDefinitionID, err := strconv.Atoi(d.Id())
	var releaseDefinitionReference *int = nil
	if err == nil {
		releaseDefinitionReference = &releaseDefinitionID
	}

	createdOn, _ := time.Parse(time.RFC3339, d.Get("created_on").(string))
	modifiedOn, _ := time.Parse(time.RFC3339, d.Get("modified_on").(string))
	source := expandReleaseDefinitionSource(d.Get("source").(string))
	variableGroups := expandIntList(d.Get("variable_groups").([]interface{}))
	environments := expandReleaseDefinitionEnvironmentList(d.Get("stage").([]interface{}))
	variables := expandReleaseConfigurationVariableValueList(d.Get("variable").([]interface{}))
	properties := expandStringMapStringFirstOrNil(d.Get("properties").([]interface{}), expandReleaseDefinitionProperties)
	triggers := expandList(d.Get("triggers").([]interface{}), expandReleaseDefinitionTriggers)
	buildArtifacts := expandReleaseArtifactList(d.Get("build_artifact").([]interface{}), release.AgentArtifactTypeValues.Build)

	artifacts := buildArtifacts

	tags := tfhelper.ExpandStringList(d.Get("tags").([]interface{}))

	releaseDefinition := release.ReleaseDefinition{
		Id:                releaseDefinitionReference,
		Name:              converter.String(d.Get("name").(string)),
		Path:              converter.String(d.Get("path").(string)),
		Revision:          converter.Int(d.Get("revision").(int)),
		Description:       converter.String(d.Get("description").(string)),
		Environments:      &environments,
		Variables:         &variables,
		ReleaseNameFormat: converter.String(d.Get("release_name_format").(string)),
		VariableGroups:    &variableGroups,
		Properties:        properties,
		Artifacts:         &artifacts,
		Url:               converter.String(d.Get("url").(string)),
		Comment:           converter.String(d.Get("comment").(string)),
		Tags:              &tags,
		CreatedOn:         &azuredevops.Time{Time: createdOn},
		ModifiedOn:        &azuredevops.Time{Time: modifiedOn},
		IsDeleted:         converter.Bool(d.Get("is_deleted").(bool)),
		Source:            &source,
		Triggers:          &triggers,
	}

	return &releaseDefinition, projectID, nil
}

// AgentDeploymentInput phase type which contains tasks executed on agent
type AgentDeploymentInput struct {
	// Gets or sets the job condition.
	Condition *string `json:"condition,omitempty"`
	// Gets or sets the job cancel timeout in minutes for deployment which are cancelled by user for this release environment.
	JobCancelTimeoutInMinutes *int `json:"jobCancelTimeoutInMinutes,omitempty"`
	// Gets or sets the override inputs.
	OverrideInputs *map[string]string `json:"overrideInputs,omitempty"`
	// Gets or sets the job execution timeout in minutes for deployment which are queued against this release environment.
	TimeoutInMinutes *int `json:"timeoutInMinutes,omitempty"`
	// Artifacts that downloaded during job execution.
	ArtifactsDownloadInput *release.ArtifactsDownloadInput `json:"artifactsDownloadInput,omitempty"`
	// List demands that needs to meet to execute the job.
	Demands *[]interface{} `json:"demands,omitempty"`
	// Indicates whether to include access token in deployment job or not.
	EnableAccessToken *bool `json:"enableAccessToken,omitempty"`
	// Id of the pool on which job get executed.
	QueueID *int `json:"queueId,omitempty"`
	// Indicates whether artifacts downloaded while job execution or not.
	SkipArtifactsDownload *bool `json:"skipArtifactsDownload,omitempty"`
	// Specification for an agent on which a job gets executed.
	AgentSpecification *release.AgentSpecification `json:"agentSpecification,omitempty"`
	// Gets or sets the image ID.
	ImageID *int `json:"imageId,omitempty"`
	// Gets or sets the parallel execution input.
	ParallelExecution interface{} `json:"parallelExecution,omitempty"`
}

// ServerDeploymentInput phase type which contains tasks executed by server
type ServerDeploymentInput struct {
	// Gets or sets the job condition.
	Condition *string `json:"condition,omitempty"`
	// Gets or sets the job cancel timeout in minutes for deployment which are cancelled by user for this release environment.
	JobCancelTimeoutInMinutes *int `json:"jobCancelTimeoutInMinutes,omitempty"`
	// Gets or sets the override inputs.
	OverrideInputs *map[string]string `json:"overrideInputs,omitempty"`
	// Gets or sets the job execution timeout in minutes for deployment which are queued against this release environment.
	TimeoutInMinutes *int `json:"timeoutInMinutes,omitempty"`
	// Gets or sets the parallel execution input.
	ParallelExecution interface{} `json:"parallelExecution,omitempty"`
}

// ReleaseDeployPhase the deploy phase
type ReleaseDeployPhase struct {
	// Dynamic based on PhaseType
	DeploymentInput interface{} `json:"deploymentInput,omitempty"`
	// WorkflowTasks
	WorkflowTasks *[]release.WorkflowTask `json:"workflowTasks,omitempty"`
	// Gets or sets the reference name of the task.
	RefName *string `json:"refName,omitempty"`
	// Name of the phase.
	Name *string `json:"name,omitempty"`
	// Type of the phase.
	PhaseType *release.DeployPhaseTypes `json:"phaseType,omitempty"`
	// Rank of the phase.
	Rank *int `json:"rank,omitempty"`

	// TODO : Figure out if this is used for Deployment Job Group or all 3 Job Types.
	// Deployment jobs of the phase.
	//DeploymentJobs *[]release.DeploymentJob `json:"deploymentJobs,omitempty"`

	// TODO : Add manual_intervention {} (block) under agentless_job { } (block)
	// List of manual intervention tasks execution information in phase.
	ManualInterventions *[]release.ManualIntervention `json:"manualInterventions,omitempty"`

	// TODO : Remove below properties after a little R&D

	// TODO : Going to remove Id. As it is Deprecated.
	// Deprecated:
	//Id *int `json:"id,omitempty"`

	// TODO : Consider removing ID. As I do not believe you can change the value via the API.
	// TODO : If you can change via API then allow updating. Also explore if this cause a ForceNew/ForceReplace
	// ID of the phase.
	//PhaseId *string `json:"phaseId,omitempty"`

	// TODO : Consider removing RunPlanId. It is stateful data about the current state of the pipeline.
	// TODO : This does not seem like something one would want with terraform.
	// Run Plan ID of the phase.
	//RunPlanId *uuid.UUID `json:"runPlanId,omitempty"`

	// TODO : Consider removing StartedOn. It is stateful data about the current state of the pipeline.
	// TODO : This does not seem like something one would want with terraform.
	// Phase start time.
	//StartedOn *azuredevops.Time `json:"startedOn,omitempty"`

	// TODO : Consider removing Status. It is stateful data about the current state of the pipeline.
	// TODO : This does not seem like something one would want with terraform.
	// Status of the phase.
	//Status *release.DeployPhaseStatus `json:"status,omitempty"`

	// TODO : Consider removing ErrorLog. It is stateful data about the current state of the pipeline.
	// TODO : This does not seem like something one would want with terraform.
	// Phase execution error logs.
	//ErrorLog *string `json:"errorLog,omitempty"`
}

// ArtifactDownloadModeType download type
type ArtifactDownloadModeType string
type artifactDownloadModeTypeValuesType struct {
	Skip      ArtifactDownloadModeType
	Selective ArtifactDownloadModeType
	All       ArtifactDownloadModeType
}

// ArtifactDownloadModeTypeValues enum of download type
var ArtifactDownloadModeTypeValues = artifactDownloadModeTypeValuesType{
	Skip:      "Skip",
	Selective: "Selective",
	All:       "All",
}

// DeploymentHealthOptionType health check type
type DeploymentHealthOptionType string
type deploymentHealthOptionValuesType struct {
	OneTargetAtATime DeploymentHealthOptionType
	Custom           DeploymentHealthOptionType
}

// DeploymentHealthOptionTypeValues enum of health check type
var DeploymentHealthOptionTypeValues = deploymentHealthOptionValuesType{
	OneTargetAtATime: "OneTargetAtATime",
	Custom:           "Custom",
}

// ReleaseDefinitionDemand demand
type ReleaseDefinitionDemand struct {
	// Name of the demand.
	Name *string `json:"name,omitempty"`
	// The value of the demand.
	Value *string `json:"value,omitempty"`
}

// MachineGroupDeploymentMultiple phase type which contains tasks executed on deployment group machines.
type MachineGroupDeploymentMultiple struct {
}

// ReleaseHostedAzurePipelines hosted agent details
type ReleaseHostedAzurePipelines struct {
	AgentSpecification *release.AgentSpecification
	QueueID            *int
}

// DefaultVersionType health check type
type DefaultVersionType string
type defaultVersionValuesType struct {
	Latest                                 DefaultVersionType
	LatestWithBranchAndTags                DefaultVersionType
	LatestWithBuildDefinitionBranchAndTags DefaultVersionType
	SpecificVersion                        DefaultVersionType
	SelectDuringReleaseCreation            DefaultVersionType
}

// DefaultVersionTypeValues enum of health check type
var DefaultVersionTypeValues = defaultVersionValuesType{
	Latest:                                 "latestType",                                 // "Latest"
	LatestWithBranchAndTags:                "latestWithBranchAndTagsType",                // "Latest from a specific branch with tags"
	LatestWithBuildDefinitionBranchAndTags: "latestWithBuildDefinitionBranchAndTagsType", // "Latest from the build pipeline default branch with tags"
	SpecificVersion:                        "specificVersionType",                        // "Specific version"
	SelectDuringReleaseCreation:            "selectDuringReleaseCreationType",            // "Specify at the time of release creation"
}

// DeploymentType deployment type
type DeploymentType string
type deploymentTypeValuesType struct {
	Production  DeploymentType
	Staging     DeploymentType
	Testing     DeploymentType
	Development DeploymentType
	Unmapped    DeploymentType
}

// DeploymentTypeValues enum of download type
var DeploymentTypeValues = deploymentTypeValuesType{
	Production:  "production",
	Staging:     "staging",
	Testing:     "testing",
	Development: "development",
	Unmapped:    "unmapped",
}

type Release interface {
	release.ParallelExecutionInputBase |
		release.ReleaseDefinitionDeployStep |
		release.EnvironmentOptions |
		MachineGroupDeploymentMultiple |
		release.MultiConfigInput |
		release.AgentSpecification |
		release.WorkflowTask |
		release.EnvironmentRetentionPolicy
}

func expandStringMapString(d map[string]interface{}) map[string]string {
	vs := make(map[string]string)
	for k, v := range d {
		if s, ok := v.(string); ok {
			vs[k] = s
		} else if s, ok := v.(int); ok {
			vs[k] = fmt.Sprint(s)
		} else if b, ok := v.(bool); ok {
			vs[k] = strconv.FormatBool(b)
		} else if f, ok := v.(float64); ok {
			vs[k] = strconv.FormatFloat(f, 'E', -1, 64)
		}
	}
	return vs
}

func expandList[T Release | map[string]interface{} | interface{}](d []interface{}, f func(map[string]interface{}) T) []T {
	vs := make([]T, 0, 0)
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, f(val))
		}
	}
	return vs
}

func expandIntList(d []interface{}) []int {
	vs := make([]int, 0, len(d))
	for _, v := range d {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(int))
		}
	}
	return vs
}

func expandFirstOrNil[T Release](d []interface{}, f func(map[string]interface{}) T) *T {
	d2 := expandList(d, f)
	if len(d2) != 1 {
		return nil
	}
	return &d2[0]
}

func expandStringMapStringFirstOrNil[T map[string]interface{}](d []interface{}, f func(map[string]interface{}) T) T {
	d2 := expandList(d, f)
	if len(d2) != 1 {
		return nil
	}
	return d2[0]
}

func expandInterfaceFirstOrNil(d []interface{}, f func(map[string]interface{}) map[string]interface{}) interface{} {
	d2 := expandList(d, f)
	if len(d2) != 1 {
		return nil
	}
	return d2[0]
}

func expandReleaseDefinitionSource(d string) release.ReleaseDefinitionSource {
	switch d {
	case string(release.ReleaseDefinitionSourceValues.RestApi):
		return release.ReleaseDefinitionSourceValues.RestApi
	case string(release.ReleaseDefinitionSourceValues.Ibiza):
		return release.ReleaseDefinitionSourceValues.Ibiza
	case string(release.ReleaseDefinitionSourceValues.PortalExtensionApi):
		return release.ReleaseDefinitionSourceValues.PortalExtensionApi
	case string(release.ReleaseDefinitionSourceValues.UserInterface):
		return release.ReleaseDefinitionSourceValues.UserInterface
	}
	return release.ReleaseDefinitionSourceValues.Undefined
}

func expandReleaseStringProperty(d string) interface{} {
	return map[string]interface{}{
		"$type":  "System.String",
		"$value": d,
	}
}

func expandReleaseEnvironmentProperties(d map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"BoardsEnvironmentType": expandReleaseStringProperty(d["boards_environment_type"].(string)),
		"LinkBoardsWorkItems":   expandReleaseStringProperty(strconv.FormatBool(d["link_boards_work_items"].(bool))),
		"JiraEnvironmentType":   expandReleaseStringProperty(d["jira_environment_type"].(string)),
		"LinkJiraWorkItems":     expandReleaseStringProperty(strconv.FormatBool(d["link_jira_work_items"].(bool))),
	}
}

func expandReleaseDefinitionProperties(d map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"DefinitionCreationSource": expandReleaseStringProperty(d["definition_creation_source"].(string)),
		"IntegrateJiraWorkItems":   expandReleaseStringProperty(strconv.FormatBool(d["integrate_jira_work_items"].(bool))),
		"IntegrateBoardsWorkItems": expandReleaseStringProperty(strconv.FormatBool(d["integrate_boards_work_items"].(bool))),
		"JiraServiceEndpointId":    expandReleaseStringProperty(d["jira_service_endpoint_id"].(string)),
	}
}

func expandReleaseDefinitionTriggers(d map[string]interface{}) interface{} {
	return map[string]interface{}{}
}

func expandReleaseCondition(d map[string]interface{}, t release.ConditionType) release.Condition {
	vs := release.Condition{ConditionType: &t}
	switch t {
	case release.ConditionTypeValues.EnvironmentState:
		vs.Name = converter.String(d["stage_name"].(string))
		vs.Value = converter.String("4")
	case release.ConditionTypeValues.Event:
		vs.Name = converter.String(d["event_name"].(string))
		vs.Value = converter.String("")
	}
	return vs
}
func expandReleaseConditionList(d []interface{}, t release.ConditionType, eventName *string) []release.Condition {
	vs := make([]release.Condition, 0, len(d))
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			if val["event_name"] == "" && eventName != nil {
				val["event_name"] = *eventName
			}
			vs = append(vs, expandReleaseCondition(val, t))
		} else if eventName != nil {
			condition := map[string]interface{}{"event_name": *eventName}
			vs = append(vs, expandReleaseCondition(condition, t))
		}
	}
	return vs
}

func expandArtifactIncludeExclude(d map[string]interface{}, name string) release.Condition {
	tags := make([]interface{}, 0, 0)
	if m, ok := d["tags"].(*schema.Set); ok {
		tags = m.List()
	}
	vs := map[string]interface{}{
		"sourceBranch":                d["branch_name"],
		"tags":                        tags,
		"useBuildDefinitionBranch":    false,
		"createReleaseOnBuildTagging": false,
	}
	value, _ := json.Marshal(vs)
	return release.Condition{
		Name:          converter.String(name),
		ConditionType: &release.ConditionTypeValues.Artifact,
		Value:         converter.String(string(value)),
	}
}
func expandArtifactIncludeExcludeList(d []interface{}, name string) []release.Condition {
	vs := make([]release.Condition, 0, len(d))
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandArtifactIncludeExclude(val, name))
		}
	}
	return vs
}

func expandReleaseConditionArtifactFilter(d map[string]interface{}) []release.Condition {
	includes := expandArtifactIncludeExcludeList(d["include"].([]interface{}), d["artifact_alias"].(string))
	excludes := expandArtifactIncludeExcludeList(d["exclude"].([]interface{}), d["artifact_alias"].(string))
	return append(includes, excludes...)
}
func expandReleaseConditionArtifactFilterList(d []interface{}) []release.Condition {
	vs := make([]release.Condition, 0, 0)
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandReleaseConditionArtifactFilter(val)...)
		}
	}
	return vs
}

func deployPhasesToInterface(deployPhases interface{}) *[]interface{} {
	data, _ := json.Marshal(deployPhases)
	var dp []interface{}
	var _ = json.Unmarshal(data, &dp)
	return &dp
}

func expandJob(d map[string]interface{}, rank int) *ReleaseDeployPhase {
	vs := &ReleaseDeployPhase{}
	if a, ok := d["agent"].([]interface{}); ok && len(a) > 0 {
		vs = expandReleaseDeployPhaseListFirstOrNil(a, release.DeployPhaseTypesValues.AgentBasedDeployment)
	}
	if dg, ok := d["deployment_group"].([]interface{}); ok && len(dg) > 0 {
		vs = expandReleaseDeployPhaseListFirstOrNil(dg, release.DeployPhaseTypesValues.MachineGroupBasedDeployment)
	}
	if al, ok := d["agentless"].([]interface{}); ok && len(al) > 0 {
		vs = expandReleaseDeployPhaseListFirstOrNil(al, release.DeployPhaseTypesValues.RunOnServer)
	}
	vs.Rank = converter.Int(rank)
	return vs
}

func expandJobsList(d []interface{}) []interface{} {
	vs := make([]interface{}, 0, 0)
	for i, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandJob(val, i+1))
		}
	}
	return vs
}

func expandReleaseDefinitionEnvironment(d map[string]interface{}, rank int) release.ReleaseDefinitionEnvironment {
	variableGroups := expandIntList(d["variable_groups"].([]interface{}))
	deployStep := expandFirstOrNil(d["deploy_step"].([]interface{}), expandReleaseDefinitionDeployStep)
	variables := expandReleaseConfigurationVariableValueList(d["variable"].([]interface{}))
	demands := expandList(d["demand"].([]interface{}), expandReleaseDefinitionDemand)
	environmentOptions := expandFirstOrNil(d["environment_options"].([]interface{}), expandReleaseEnvironmentOptions)
	retentionPolicy := expandFirstOrNil(d["retention_policy"].([]interface{}), expandReleaseEnvironmentRetentionPolicy)
	preApprovalOptions := expandReleaseApprovalOptionsListFirstOrNil(d["approval_options"].([]interface{}), release.ApprovalExecutionOrderValues.BeforeGates)
	postApprovalOptions := expandReleaseApprovalOptionsListFirstOrNil(d["approval_options"].([]interface{}), release.ApprovalExecutionOrderValues.AfterSuccessfulGates)
	preDeployApprovals := expandReleaseDefinitionApprovalsListFirstOrNil(d["pre_deploy_approval"].([]interface{}), preApprovalOptions)
	postDeployApprovals := expandReleaseDefinitionApprovalsListFirstOrNil(d["post_deploy_approval"].([]interface{}), postApprovalOptions)
	properties := expandInterfaceFirstOrNil(d["properties"].([]interface{}), expandReleaseEnvironmentProperties)
	deployPhases := expandJobsList(d["job"].([]interface{}))
	preDeploymentGates := expandReleaseDefinitionGatesStepListFirstOrNil(d["pre_deploy_gate"].([]interface{}))
	postDeploymentGates := expandReleaseDefinitionGatesStepListFirstOrNil(d["post_deploy_gate"].([]interface{}))
	afterStageConditions := expandReleaseConditionList(d["after_stage"].([]interface{}), release.ConditionTypeValues.EnvironmentState, nil)
	afterReleaseConditions := expandReleaseConditionList(d["after_release"].([]interface{}), release.ConditionTypeValues.Event, converter.String("ReleaseStarted"))
	artifactFilters := expandReleaseConditionArtifactFilterList(d["artifact_filter"].([]interface{}))
	conditions := append(append(afterStageConditions, afterReleaseConditions...), artifactFilters...)

	return release.ReleaseDefinitionEnvironment{
		Id:                 converter.Int(d["id"].(int)),
		Conditions:         &conditions,
		Demands:            &demands,
		DeployPhases:       deployPhasesToInterface(deployPhases),
		DeployStep:         deployStep,
		EnvironmentOptions: environmentOptions,
		//EnvironmentTriggers: nil,
		//ExecutionPolicy:     nil,
		//Id:                  converter.Int(d["id"].(int)),
		Name:                converter.String(d["name"].(string)),
		PostDeployApprovals: postDeployApprovals,
		PostDeploymentGates: preDeploymentGates,
		PreDeployApprovals:  preDeployApprovals,
		PreDeploymentGates:  postDeploymentGates,
		//ProcessParameters:   nil,
		Properties:      properties,
		QueueId:         nil,
		Rank:            converter.Int(rank),
		RetentionPolicy: retentionPolicy,
		//RunOptions:      nil,
		//Schedules:       nil,
		VariableGroups: &variableGroups,
		Variables:      &variables,
		Owner:          &webapi.IdentityRef{Id: converter.String(d["owner_id"].(string))},
	}
}
func expandReleaseDefinitionEnvironmentList(d []interface{}) []release.ReleaseDefinitionEnvironment {
	vs := make([]release.ReleaseDefinitionEnvironment, 0, len(d))
	for i, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandReleaseDefinitionEnvironment(val, i+1))
		}
	}
	return vs
}

func expandSingleMapList(d []interface{}) map[string]interface{} {
	if len(d) == 1 {
		if ds, ok := d[0].(map[string]interface{}); ok {
			return ds
		}
	}
	return nil
}

func expandReleaseBuildArtifact(d map[string]interface{}) map[string]release.ArtifactSourceReference {
	defaultType := DefaultVersionTypeValues.Latest
	branch := ""
	version := ""
	tags := ""

	latest := expandSingleMapList(d["latest"].([]interface{}))
	if latest != nil {
		branch = latest["branch"].(string)
		tags = latest["tags"].(string)
		if branch != "" && tags != "" {
			defaultType = DefaultVersionTypeValues.LatestWithBranchAndTags
		} else if tags != "" {
			defaultType = DefaultVersionTypeValues.LatestWithBuildDefinitionBranchAndTags
		}
	}

	specify := expandSingleMapList(d["specify"].([]interface{}))
	if specify != nil {
		version = specify["version"].(string)
		if version != "" {
			defaultType = DefaultVersionTypeValues.SpecificVersion
		} else {
			defaultType = DefaultVersionTypeValues.SelectDuringReleaseCreation
		}
	}
	if defaultType == DefaultVersionTypeValues.Latest && len(d["specify"].([]interface{})) == 1 {
		defaultType = DefaultVersionTypeValues.SelectDuringReleaseCreation
	}

	return map[string]release.ArtifactSourceReference{
		//"artifactSourceDefinitionUrl": {Id: converter.String("")},
		"defaultVersionBranch":   {Id: converter.String(branch)},
		"defaultVersionSpecific": {Id: converter.String(version)},
		"defaultVersionTags":     {Id: converter.String(tags)},
		"defaultVersionType":     {Id: converter.String(string(defaultType))},
		"definition":             {Id: converter.String(d["build_pipeline_id"].(string))},
		"definitions":            {Id: converter.String("")},
		"IsMultiDefinitionType":  {Id: converter.String("False")},
		"project":                {Id: converter.String(d["project_id"].(string))},
		//"repository":                  {Id: converter.String("")},
	}
}

func expandReleaseArtifactDefinitionReference(d map[string]interface{}, t release.AgentArtifactType) map[string]release.ArtifactSourceReference {
	switch t {
	case release.AgentArtifactTypeValues.Build:
		return expandReleaseBuildArtifact(d)
	}
	return nil
}

func expandReleaseArtifact(d map[string]interface{}, t release.AgentArtifactType) release.Artifact {
	definitionReference := expandReleaseArtifactDefinitionReference(d, t)
	return release.Artifact{
		Alias:               converter.String(d["alias"].(string)),
		IsPrimary:           converter.Bool(d["is_primary"].(bool)),
		IsRetained:          converter.Bool(d["is_retained"].(bool)),
		Type:                converter.String(strings.Title(string(t))),
		DefinitionReference: &definitionReference,
	}
}
func expandReleaseArtifactList(d []interface{}, t release.AgentArtifactType) []release.Artifact {
	vs := make([]release.Artifact, 0, len(d))
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandReleaseArtifact(val, t))
		}
	}
	return vs
}

func expandReleaseConfigurationVariableValue(d map[string]interface{}) (string, release.ConfigurationVariableValue) {
	return d["name"].(string), release.ConfigurationVariableValue{
		AllowOverride: converter.Bool(d["allow_override"].(bool)),
		Value:         converter.String(d["value"].(string)),
		IsSecret:      converter.Bool(d["is_secret"].(bool)),
	}
}
func expandReleaseConfigurationVariableValueList(d []interface{}) map[string]release.ConfigurationVariableValue {
	vs := make(map[string]release.ConfigurationVariableValue)
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			key, d2 := expandReleaseConfigurationVariableValue(val)
			vs[key] = d2
		}
	}
	return vs
}

func expandReleaseDefinitionDeployStep(d map[string]interface{}) release.ReleaseDefinitionDeployStep {
	tasks := expandList(d["task"].([]interface{}), expandReleaseWorkFlowTask)
	return release.ReleaseDefinitionDeployStep{
		Id:    converter.Int(d["id"].(int)),
		Tasks: &tasks,
	}
}

func expandReleaseDeployPhase(d map[string]interface{}, t release.DeployPhaseTypes) ReleaseDeployPhase {
	workflowTasks := expandList(d["task"].([]interface{}), expandReleaseWorkFlowTask)
	var deploymentInput interface{}
	switch t {
	case release.DeployPhaseTypesValues.AgentBasedDeployment:
		deploymentInput = expandReleaseAgentDeploymentInput(d)
	case release.DeployPhaseTypesValues.MachineGroupBasedDeployment:
		deploymentInput = expandReleaseMachineGroupDeploymentInput(d)
	case release.DeployPhaseTypesValues.RunOnServer:
		deploymentInput = expandReleaseServerDeploymentInput(d)
	}
	return ReleaseDeployPhase{
		DeploymentInput: &deploymentInput,
		Rank:            converter.Int(d["rank"].(int)),
		PhaseType:       &t,
		Name:            converter.String(d["name"].(string)),
		//RefName:         converter.String(d["ref_name"].(string)),
		WorkflowTasks: &workflowTasks,
	}
}
func expandReleaseDeployPhaseList(d []interface{}, t release.DeployPhaseTypes) []ReleaseDeployPhase {
	vs := make([]ReleaseDeployPhase, 0, len(d))
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandReleaseDeployPhase(val, t))
		}
	}
	return vs
}
func expandReleaseDeployPhaseListFirstOrNil(d []interface{}, t release.DeployPhaseTypes) *ReleaseDeployPhase {
	d2 := expandReleaseDeployPhaseList(d, t)
	if len(d2) != 1 {
		return nil
	}
	return &d2[0]
}

func expandReleaseEnvironmentOptions(d map[string]interface{}) release.EnvironmentOptions {
	return release.EnvironmentOptions{
		AutoLinkWorkItems:            converter.Bool(d["auto_link_work_items"].(bool)),
		BadgeEnabled:                 converter.Bool(d["badge_enabled"].(bool)),
		PublishDeploymentStatus:      converter.Bool(d["publish_deployment_status"].(bool)),
		PullRequestDeploymentEnabled: converter.Bool(d["pull_request_deployment_enabled"].(bool)),
	}
}

func expandReleaseMachineGroupDeploymentInputMultiple(d map[string]interface{}) MachineGroupDeploymentMultiple {
	return MachineGroupDeploymentMultiple{}
}

func expandReleaseMultiConfigInput(d map[string]interface{}) release.MultiConfigInput {
	return release.MultiConfigInput{
		Multipliers:           converter.String(d["multipliers"].(string)),
		MaxNumberOfAgents:     converter.Int(d["number_of_agents"].(int)),
		ParallelExecutionType: &release.ParallelExecutionTypesValues.MultiConfiguration,
		ContinueOnError:       converter.Bool(d["continue_on_error"].(bool)),
	}
}

func expandReleaseParallelExecutionInputBase(d map[string]interface{}) release.ParallelExecutionInputBase {
	return release.ParallelExecutionInputBase{
		MaxNumberOfAgents:     converter.Int(d["max_number_of_agents"].(int)),
		ParallelExecutionType: &release.ParallelExecutionTypesValues.MultiMachine,
		ContinueOnError:       converter.Bool(d["continue_on_error"].(bool)),
	}
}

func expandReleaseAgentSpecification(d map[string]interface{}) release.AgentSpecification {
	return release.AgentSpecification{
		Identifier: converter.String(d["agent_specification"].(string)),
	}
}

func expandReleaseHostedAzurePipelines(d map[string]interface{}) ReleaseHostedAzurePipelines {
	agentSpecification := expandReleaseAgentSpecification(d)
	return ReleaseHostedAzurePipelines{
		AgentSpecification: &agentSpecification,
		QueueID:            converter.Int(d["agent_pool_id"].(int)),
	}
}
func expandReleaseHostedAzurePipelinesList(d []interface{}) []ReleaseHostedAzurePipelines {
	vs := make([]ReleaseHostedAzurePipelines, 0, len(d))
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandReleaseHostedAzurePipelines(val))
		}
	}
	return vs
}
func expandReleaseHostedAzurePipelinesListFirstOrNil(d []interface{}) (*release.AgentSpecification, int) {
	d2 := expandReleaseHostedAzurePipelinesList(d)
	if len(d2) != 1 {
		return nil, 0
	}
	return d2[0].AgentSpecification, *d2[0].QueueID
}

func expandReleaseMachineGroupDeploymentInput(d map[string]interface{}) *release.MachineGroupDeploymentInput {
	tags := tfhelper.ExpandStringList(d["tags"].([]interface{}))
	multiple := expandFirstOrNil(d["multiple"].([]interface{}), expandReleaseMachineGroupDeploymentInputMultiple)
	deploymentHealthOption := DeploymentHealthOptionTypeValues.OneTargetAtATime
	if multiple != nil {
		deploymentHealthOption = DeploymentHealthOptionTypeValues.Custom
	}
	return &release.MachineGroupDeploymentInput{
		Condition:                 converter.String(d["condition"].(string)),
		JobCancelTimeoutInMinutes: converter.Int(d["max_execution_time_in_minutes"].(int)),
		OverrideInputs:            nil, // TODO : OverrideInputs
		TimeoutInMinutes:          converter.Int(d["timeout_in_minutes"].(int)),
		ArtifactsDownloadInput:    &release.ArtifactsDownloadInput{},
		EnableAccessToken:         converter.Bool(d["allow_scripts_to_access_oauth_token"].(bool)),
		QueueId:                   converter.Int(d["deployment_group_id"].(int)),
		SkipArtifactsDownload:     converter.Bool(d["skip_artifacts_download"].(bool)),
		DeploymentHealthOption:    converter.String(string(deploymentHealthOption)),
		Tags:                      &tags,
	}
}
func expandReleaseAgentDeploymentInput(d map[string]interface{}) AgentDeploymentInput {
	buildArtifactDownloads := expandReleaseArtifactDownloadInputBaseList(d["build_artifact_download"].([]interface{}), release.AgentArtifactTypeValues.Build)
	downloadInputs := append(buildArtifactDownloads)

	demands := expandList(d["demand"].([]interface{}), expandReleaseDeployPhaseDemand)
	agentPoolPrivate := expandFirstOrNil(d["agent_pool_private"].([]interface{}), expandReleaseAgentSpecification)

	agentPoolHostedAzurePipelines, queueID := expandReleaseHostedAzurePipelinesListFirstOrNil(d["agent_pool_hosted_azure_pipelines"].([]interface{}))
	//if agentPoolPrivate != nil && agentPoolHostedAzurePipelines != nil { // TODO : how to solve
	//	return nil, fmt.Errorf("conflit %s and %s specify only one", "agent_pool_hosted_azure_pipelines", "agent_pool_private")
	//}
	var agentSpecification *release.AgentSpecification
	if agentPoolHostedAzurePipelines != nil {
		agentSpecification = agentPoolHostedAzurePipelines
	} else {
		agentSpecification = agentPoolPrivate
	}

	var parallelExecution interface{} = &release.ExecutionInput{
		ParallelExecutionType: &release.ParallelExecutionTypesValues.None,
	}
	multiConfiguration := expandFirstOrNil(d["multi_configuration"].([]interface{}), expandReleaseMultiConfigInput)
	multiAgent := expandFirstOrNil(d["multi_agent"].([]interface{}), expandReleaseParallelExecutionInputBase)
	//if multiConfiguration != nil && multiAgent != nil { // TODO : how to solve
	//	return nil, fmt.Errorf("conflit %s and %s specify only one", "multi_configuration", "multi_agent")
	//}
	if multiConfiguration != nil {
		parallelExecution = multiConfiguration
	} else if multiAgent != nil {
		parallelExecution = multiAgent
	}

	return AgentDeploymentInput{
		Condition:                 converter.String(d["condition"].(string)),
		JobCancelTimeoutInMinutes: converter.Int(d["max_execution_time_in_minutes"].(int)),
		OverrideInputs:            nil, // TODO : OverrideInputs
		TimeoutInMinutes:          converter.Int(d["timeout_in_minutes"].(int)),
		Demands:                   &demands,
		EnableAccessToken:         converter.Bool(d["allow_scripts_to_access_oauth_token"].(bool)),
		QueueID:                   &queueID,
		SkipArtifactsDownload:     converter.Bool(d["skip_artifacts_download"].(bool)),
		AgentSpecification:        agentSpecification,
		ParallelExecution:         &parallelExecution,
		ArtifactsDownloadInput: &release.ArtifactsDownloadInput{
			DownloadInputs: &downloadInputs,
		},
	}
}
func expandReleaseServerDeploymentInput(d map[string]interface{}) *ServerDeploymentInput {
	var parallelExecution interface{} = &release.ExecutionInput{
		ParallelExecutionType: &release.ParallelExecutionTypesValues.None,
	}
	multiConfiguration := expandFirstOrNil(d["multi_configuration"].([]interface{}), expandReleaseMultiConfigInput)
	if multiConfiguration != nil {
		parallelExecution = multiConfiguration
	}
	return &ServerDeploymentInput{
		Condition:                 converter.String(d["condition"].(string)),
		JobCancelTimeoutInMinutes: converter.Int(d["max_execution_time_in_minutes"].(int)),
		OverrideInputs:            nil, // TODO : OverrideInputs
		TimeoutInMinutes:          converter.Int(d["timeout_in_minutes"].(int)),
		ParallelExecution:         &parallelExecution,
	}
}

func expandReleaseArtifactDownloadInputBase(d map[string]interface{}, t release.AgentArtifactType) release.ArtifactDownloadInputBase {
	mode := ArtifactDownloadModeTypeValues.All
	artifactItems := make([]string, 0, 0)
	artifactItems = tfhelper.ExpandStringList(d["include"].([]interface{}))
	if len(artifactItems) > 0 {
		if artifactItems[0] == "*" {
			artifactItems = []string{}
			mode = ArtifactDownloadModeTypeValues.All
		} else {
			mode = ArtifactDownloadModeTypeValues.Selective
		}
	} else {
		mode = ArtifactDownloadModeTypeValues.Skip
	}
	return release.ArtifactDownloadInputBase{
		Alias:                converter.String(d["artifact_alias"].(string)),
		ArtifactDownloadMode: converter.String(string(mode)),
		ArtifactItems:        &artifactItems,
		ArtifactType:         converter.String(strings.Title(string(t))),
	}
}
func expandReleaseArtifactDownloadInputBaseList(d []interface{}, t release.AgentArtifactType) []release.ArtifactDownloadInputBase {
	vs := make([]release.ArtifactDownloadInputBase, 0, len(d))
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandReleaseArtifactDownloadInputBase(val, t))
		}
	}
	return vs
}

func expandReleaseDeployPhaseDemand(d map[string]interface{}) interface{} {
	name := d["name"].(string)
	configValue := d["value"].(string)
	if len(configValue) > 0 {
		name += " -equals " + configValue
	}
	return name
}

func expandReleaseDefinitionDemand(d map[string]interface{}) interface{} {
	name := d["name"].(string)
	configValue := d["value"].(string)
	if len(configValue) > 0 {
		name += " -equals " + configValue
	}
	return ReleaseDefinitionDemand{
		Name: converter.String(name),
	}
}

func expandReleaseWorkFlowTask(d map[string]interface{}) release.WorkflowTask {
	task := strings.Split(d["task"].(string), "@")
	taskName, version := task[0], task[1]
	taskID := taskagent.TaskNameToUUID[taskName]

	inputs := expandStringMapString(d["inputs"].(map[string]interface{}))
	environment := expandStringMapString(d["environment"].(map[string]interface{}))
	overrideInputs := expandStringMapString(d["override_inputs"].(map[string]interface{}))

	return release.WorkflowTask{
		TaskId:           &taskID,
		Name:             converter.String(d["display_name"].(string)),
		AlwaysRun:        converter.Bool(d["always_run"].(bool)),
		Condition:        converter.String(d["condition"].(string)),
		ContinueOnError:  converter.Bool(d["continue_on_error"].(bool)),
		DefinitionType:   converter.String("task"),
		Enabled:          converter.Bool(d["enabled"].(bool)),
		TimeoutInMinutes: converter.Int(d["timeout_in_minutes"].(int)),
		Environment:      &environment,
		Inputs:           &inputs,
		OverrideInputs:   &overrideInputs,
		Version:          &version,
	}
}

func expandReleaseEnvironmentRetentionPolicy(d map[string]interface{}) release.EnvironmentRetentionPolicy {
	return release.EnvironmentRetentionPolicy{
		DaysToKeep:     converter.Int(d["days_to_keep"].(int)),
		RetainBuild:    converter.Bool(d["retain_build"].(bool)),
		ReleasesToKeep: converter.Int(d["releases_to_keep"].(int)),
	}
}

func automaticReleaseDefinitionApprovalStep() release.ReleaseDefinitionApprovalStep {
	return release.ReleaseDefinitionApprovalStep{
		IsAutomated:      converter.Bool(true),
		IsNotificationOn: converter.Bool(false),
		Rank:             converter.Int(1),
	}
}

func automaticReleaseDefinitionApprovals() release.ReleaseDefinitionApprovals {
	return release.ReleaseDefinitionApprovals{
		Approvals: &[]release.ReleaseDefinitionApprovalStep{automaticReleaseDefinitionApprovalStep()},
	}
}

func expandReleaseDefinitionApprovalStep(d map[string]interface{}, rank int) release.ReleaseDefinitionApprovalStep {
	vs := release.ReleaseDefinitionApprovalStep{
		Id:               converter.Int(d["id"].(int)),
		IsAutomated:      converter.Bool(true),
		IsNotificationOn: converter.Bool(false),
		Rank:             converter.Int(rank),
	}

	if d["approver_id"] != "" {
		vs.IsAutomated = converter.Bool(false)
		vs.Approver = &webapi.IdentityRef{Id: converter.String(d["approver_id"].(string))}
	}
	return vs
}
func expandReleaseDefinitionApprovalStepList(d []interface{}) []release.ReleaseDefinitionApprovalStep {
	vs := make([]release.ReleaseDefinitionApprovalStep, 0, len(d))
	for i, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			if val["id"] != nil {
				vs = append(vs, expandReleaseDefinitionApprovalStep(val, i+1))
			}
		}
	}
	if len(vs) == 0 {
		vs = append(vs, automaticReleaseDefinitionApprovalStep())
	}
	return vs
}

func expandReleaseApprovalOptions(d map[string]interface{}, executionOrder release.ApprovalExecutionOrder) release.ApprovalOptions {
	return release.ApprovalOptions{
		AutoTriggeredAndPreviousEnvironmentApprovedCanBeSkipped: converter.Bool(d["auto_triggered_and_previous_environment_approved_can_be_skipped"].(bool)),
		EnforceIdentityRevalidation:                             converter.Bool(d["enforce_identity_revalidation"].(bool)),
		ExecutionOrder:                                          &executionOrder,
		ReleaseCreatorCanBeApprover:                             converter.Bool(d["release_creator_can_be_approver"].(bool)),
		RequiredApproverCount:                                   converter.Int(d["required_approver_count"].(int)),
	}
}
func expandReleaseApprovalOptionsList(d []interface{}, executionOrder release.ApprovalExecutionOrder) []release.ApprovalOptions {
	vs := make([]release.ApprovalOptions, 0, len(d))
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandReleaseApprovalOptions(val, executionOrder))
		}
	}
	return vs
}
func expandReleaseApprovalOptionsListFirstOrNil(d []interface{}, executionOrder release.ApprovalExecutionOrder) *release.ApprovalOptions {
	d2 := expandReleaseApprovalOptionsList(d, executionOrder)
	if len(d2) != 1 {
		return nil
	}
	return &d2[0]
}

func expandReleaseDefinitionApprovals(d map[string]interface{}, approvalOptions *release.ApprovalOptions) release.ReleaseDefinitionApprovals {
	approvals := expandReleaseDefinitionApprovalStepList(d["approval"].([]interface{}))
	timeoutInMinutes := converter.Int(d["timeout_in_minutes"].(int))
	if timeoutInMinutes != nil && approvalOptions != nil {
		approvalOptions.TimeoutInMinutes = timeoutInMinutes
	}
	return release.ReleaseDefinitionApprovals{
		Approvals:       &approvals,
		ApprovalOptions: approvalOptions,
	}
}
func expandReleaseDefinitionApprovalsList(d []interface{}, approvalOptions *release.ApprovalOptions) []release.ReleaseDefinitionApprovals {
	vs := make([]release.ReleaseDefinitionApprovals, 0, len(d))
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandReleaseDefinitionApprovals(val, approvalOptions))
		}
	}
	if len(vs) == 0 {
		vs = append(vs, automaticReleaseDefinitionApprovals())
	}
	return vs
}
func expandReleaseDefinitionApprovalsListFirstOrNil(d []interface{}, approvalOptions *release.ApprovalOptions) *release.ReleaseDefinitionApprovals {
	d2 := expandReleaseDefinitionApprovalsList(d, approvalOptions)
	if len(d2) != 1 {
		return nil
	}
	return &d2[0]
}

func expandReleaseDefinitionGate(d map[string]interface{}) release.ReleaseDefinitionGate {
	workflowTasks := expandList(d["task"].([]interface{}), expandReleaseWorkFlowTask)
	return release.ReleaseDefinitionGate{
		Tasks: &workflowTasks,
	}
}
func expandReleaseDefinitionGateList(d []interface{}) []release.ReleaseDefinitionGate {
	vs := make([]release.ReleaseDefinitionGate, 0, len(d))
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandReleaseDefinitionGate(val))
		}
	}
	return vs
}

func expandReleaseDefinitionGatesStep(d map[string]interface{}) release.ReleaseDefinitionGatesStep {
	gates := expandReleaseDefinitionGateList(d["gate"].([]interface{}))
	return release.ReleaseDefinitionGatesStep{
		Id:    converter.Int(d["id"].(int)),
		Gates: &gates,
	}
}
func expandReleaseDefinitionGatesStepList(d []interface{}) []release.ReleaseDefinitionGatesStep {
	vs := make([]release.ReleaseDefinitionGatesStep, 0, len(d))
	for _, v := range d {
		if val, ok := v.(map[string]interface{}); ok {
			vs = append(vs, expandReleaseDefinitionGatesStep(val))
		}
	}
	return vs
}
func expandReleaseDefinitionGatesStepListFirstOrNil(d []interface{}) *release.ReleaseDefinitionGatesStep {
	d2 := expandReleaseDefinitionGatesStepList(d)
	if len(d2) != 1 {
		return nil
	}
	return &d2[0]
}

func flattenStringList(list []*string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, *v)
	}
	return vs
}
func flattenStringSet(list []*string) *schema.Set {
	return schema.NewSet(schema.HashString, flattenStringList(list))
}

func flattenIntList(list []*int) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, *v)
	}
	return vs
}
func flattenIntSet(list []*int) *schema.Set {
	return schema.NewSet(schema.HashString, flattenIntList(list))
}

func flattenStringMap(m *map[string]string) map[string]interface{} {
	if m == nil {
		return nil
	}
	vs := map[string]interface{}{}
	for k, v := range *m {
		if v != "" {
			vs[k] = v
		}
	}
	return vs
}

func flattenReleaseDefinitionVariables(m *map[string]release.ConfigurationVariableValue) interface{} {
	if m == nil {
		return nil
	}
	d := make([]map[string]interface{}, len(*m))
	index := 0

	keys := make([]string, 0, len(*m))
	for k := range *m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		m2 := (*m)[k]
		d[index] = map[string]interface{}{
			"name":      k,
			"value":     converter.ToString(m2.Value, ""),
			"is_secret": converter.ToBool(m2.IsSecret, false),
		}
		index = index + 1
	}
	return d
}

func flattenReleaseDefinitionEnvironmentOptions(m *release.EnvironmentOptions) interface{} {
	if m == nil {
		return nil
	}
	return []map[string]interface{}{{
		"auto_link_work_items":            m.AutoLinkWorkItems,
		"badge_enabled":                   m.BadgeEnabled,
		"publish_deployment_status":       m.PublishDeploymentStatus,
		"pull_request_deployment_enabled": m.PullRequestDeploymentEnabled,
	}}
}

func flattenReleasePropertiesMap(d interface{}) interface{} {
	if d2, ok := d.(map[string]interface{}); ok {
		if t, ok := d2["$type"].(string); ok {
			switch t {
			case "System.String":
				if v, ok := d2["$value"].(string); ok {
					if b, err := strconv.ParseBool(v); err == nil {
						return b
					}
					return v
				}
			case "System.Boolean":
				fallthrough
			case "System.Bool":
				return d2["$value"].(bool)
			}
		}
	}
	return nil
}

func flattenReleaseEnvironmentProperties(m interface{}) interface{} {
	if properties, ok := m.(map[string]interface{}); ok {
		d := map[string]interface{}{
			"boards_environment_type": flattenReleasePropertiesMap(properties["BoardsEnvironmentType"]),
			"link_boards_work_items":  flattenReleasePropertiesMap(properties["LinkBoardsWorkItems"]),
			"jira_environment_type":   flattenReleasePropertiesMap(properties["JiraEnvironmentType"]),
			"link_jira_work_items":    flattenReleasePropertiesMap(properties["LinkJiraWorkItems"]),
		}
		return []map[string]interface{}{d}
	}
	return nil
}

func flattenReleaseDefinitionProperties(m interface{}) interface{} {
	if properties, ok := m.(map[string]interface{}); ok {
		d := map[string]interface{}{
			"definition_creation_source":  flattenReleasePropertiesMap(properties["DefinitionCreationSource"]),
			"integrate_jira_work_items":   flattenReleasePropertiesMap(properties["IntegrateBoardsWorkItems"]),
			"integrate_boards_work_items": flattenReleasePropertiesMap(properties["IntegrateJiraWorkItems"]),
			"jira_service_endpoint_id":    flattenReleasePropertiesMap(properties["JiraServiceEndpointId"]),
		}
		return []map[string]interface{}{d}
	}
	return nil
}

func flattenReleaseWorkflowTask(m release.WorkflowTask) map[string]interface{} {
	task := taskagent.TaskUUIDToName[*m.TaskId] + "@" + *m.Version
	return map[string]interface{}{
		"task":               task,
		"display_name":       m.Name,
		"inputs":             flattenStringMap(m.Inputs),
		"condition":          m.Condition,
		"continue_on_error":  m.ContinueOnError,
		"always_run":         m.AlwaysRun,
		"enabled":            m.Enabled,
		"environment":        flattenStringMap(m.Environment),
		"override_inputs":    flattenStringMap(m.OverrideInputs),
		"timeout_in_minutes": m.TimeoutInMinutes,
	}
}
func flattenReleaseWorkflowTaskList(m *[]release.WorkflowTask) []map[string]interface{} {
	ds := make([]map[string]interface{}, 0, len(*m))
	for _, d := range *m {
		ds = append(ds, flattenReleaseWorkflowTask(d))
	}
	return ds
}

func flattenReleaseDeployPhaseDemand(m interface{}) map[string]interface{} {
	if ms, ok := m.(string); ok {
		m2 := strings.Split(ms, "-equals")
		var value string
		if len(m2) == 2 {
			value = strings.TrimSpace(m2[1])
		}
		return map[string]interface{}{
			"name":  strings.TrimSpace(m2[0]),
			"value": value,
		}
	}
	return nil
}
func flattenReleaseDeployPhaseDemandList(m *[]interface{}) []map[string]interface{} {
	ds := make([]map[string]interface{}, 0, len(*m))
	for _, d := range *m {
		ds = append(ds, flattenReleaseDeployPhaseDemand(d))
	}
	return ds
}

func flattenReleaseArtifactDownloadBase(m release.ArtifactDownloadInputBase) map[string]interface{} {
	include := *m.ArtifactItems
	if strings.EqualFold(*m.ArtifactDownloadMode, string(ArtifactDownloadModeTypeValues.All)) {
		include = []string{"*"}
	}
	return map[string]interface{}{
		"artifact_alias": m.Alias,
		"include":        include,
	}
}
func flattenReleaseArtifactDownloadBaseList(m *[]release.ArtifactDownloadInputBase, t release.AgentArtifactType) []map[string]interface{} {
	ds := make([]map[string]interface{}, 0, 0)
	for _, d := range *m {
		if strings.EqualFold(*d.ArtifactType, string(t)) {
			ds = append(ds, flattenReleaseArtifactDownloadBase(d))
		}
	}
	return ds
}

func flattenReleaseDefinitionDeployStep(m *release.ReleaseDefinitionDeployStep) []map[string]interface{} {
	if m == nil {
		return nil
	}
	ms := map[string]interface{}{"id": m.Id}
	if m.Tasks != nil {
		ms["task"] = flattenReleaseWorkflowTaskList(m.Tasks)
	}
	return []map[string]interface{}{ms}
}

func flattenReleaseDefinitionHostedAgentSpecification(m *release.AgentSpecification, rai *release.AgentDeploymentInput) []map[string]interface{} {
	if m == nil {
		return nil
	}
	return []map[string]interface{}{{
		"agent_pool_id":       rai.QueueId,
		"agent_specification": m.Identifier,
	}}
}

func flattenReleaseAgentDeploymentInput(rdp *release.ReleaseDeployPhase, rai *release.AgentDeploymentInput, rwt *[]release.WorkflowTask) map[string]interface{} {
	d := map[string]interface{}{
		"rank":                                rdp.Rank,
		"name":                                rdp.Name,
		"task":                                flattenReleaseWorkflowTaskList(rwt),
		"agent_pool_hosted_azure_pipelines":   flattenReleaseDefinitionHostedAgentSpecification(rai.AgentSpecification, rai),
		"build_artifact_download":             flattenReleaseArtifactDownloadBaseList(rai.ArtifactsDownloadInput.DownloadInputs, release.AgentArtifactTypeValues.Build),
		"timeout_in_minutes":                  rai.TimeoutInMinutes,
		"max_execution_time_in_minutes":       rai.JobCancelTimeoutInMinutes,
		"condition":                           rai.Condition,
		"skip_artifacts_download":             rai.SkipArtifactsDownload,
		"allow_scripts_to_access_oauth_token": rai.EnableAccessToken,
		"demand":                              flattenReleaseDeployPhaseDemandList(rai.Demands),
	}
	return d
}

func flattenReleaseServerInput(rdp *release.ReleaseDeployPhase, rai *release.ServerDeploymentInput, rwt *[]release.WorkflowTask) map[string]interface{} {
	d := map[string]interface{}{
		"rank":                          rdp.Rank,
		"name":                          rdp.Name,
		"task":                          flattenReleaseWorkflowTaskList(rwt),
		"timeout_in_minutes":            rai.TimeoutInMinutes,
		"max_execution_time_in_minutes": rai.JobCancelTimeoutInMinutes,
		"condition":                     rai.Condition,
	}
	return d
}

func unmarshalAgentDeploymentInput(deployPhase []byte, agentDeploymentInput []byte, workflowTasks []byte) (*release.ReleaseDeployPhase, *release.AgentDeploymentInput, *[]release.WorkflowTask) {
	var d release.ReleaseDeployPhase
	_ = json.Unmarshal(deployPhase, &d)

	var d2 release.AgentDeploymentInput
	_ = json.Unmarshal(agentDeploymentInput, &d2)

	var d3 []release.WorkflowTask
	_ = json.Unmarshal(workflowTasks, &d3)
	return &d, &d2, &d3
}

func unmarshalServerDeploymentInput(deployPhase []byte, agentDeploymentInput []byte, workflowTasks []byte) (*release.ReleaseDeployPhase, *release.ServerDeploymentInput, *[]release.WorkflowTask) {
	var d release.ReleaseDeployPhase
	_ = json.Unmarshal(deployPhase, &d)

	var d2 release.ServerDeploymentInput
	_ = json.Unmarshal(agentDeploymentInput, &d2)

	var d3 []release.WorkflowTask
	_ = json.Unmarshal(workflowTasks, &d3)
	return &d, &d2, &d3
}

func flattenReleaseDeploymentInput(m map[string]interface{}) map[string]interface{} {
	deployPhase, err := json.MarshalIndent(m, "", "  ")
	agentDeploymentInput, err2 := json.MarshalIndent(m["deploymentInput"], "", "  ")
	workflowTasks, _ := json.MarshalIndent(m["workflowTasks"], "", "  ")
	if err != nil || err2 != nil {
		return nil
	}

	switch release.DeployPhaseTypes(m["phaseType"].(string)) {
	case release.DeployPhaseTypesValues.AgentBasedDeployment:
		return map[string]interface{}{
			"agent": []map[string]interface{}{
				flattenReleaseAgentDeploymentInput(unmarshalAgentDeploymentInput(deployPhase, agentDeploymentInput, workflowTasks)),
			},
		}
	case release.DeployPhaseTypesValues.RunOnServer:
		return map[string]interface{}{
			"agentless": []map[string]interface{}{
				flattenReleaseServerInput(unmarshalServerDeploymentInput(deployPhase, agentDeploymentInput, workflowTasks)),
			},
		}
	}
	return nil
}

func flattenReleaseDeployPhasesList(m *[]interface{}) []map[string]interface{} {
	ds := make([]map[string]interface{}, 0, 0)
	if m != nil {
		for _, d := range *m {
			ds = append(ds, flattenReleaseDeploymentInput(d.(map[string]interface{})))
		}
	}
	return ds
}

func flattenReleaseDefinitionApprovalStep(m release.ReleaseDefinitionApprovalStep) map[string]interface{} {
	ds := map[string]interface{}{
		"id":   m.Id,
		"rank": m.Rank,
	}
	if m.Approver != nil {
		ds["approver_id"] = m.Approver.Id
	}
	return ds
}
func flattenReleaseDefinitionApprovalStepList(m *[]release.ReleaseDefinitionApprovalStep) []map[string]interface{} {
	ds := make([]map[string]interface{}, 0, 0)
	for _, d := range *m {
		ds = append(ds, flattenReleaseDefinitionApprovalStep(d))
	}
	return ds
}

func flattenApprovalOptions(m1, m2 *release.ApprovalOptions) []map[string]interface{} {
	var m *release.ApprovalOptions
	if m1 != nil && m2 == nil {
		m = m1
	} else if m2 != nil {
		m = m2
	}

	if m == nil {
		return nil
	}

	return []map[string]interface{}{{
		"auto_triggered_and_previous_environment_approved_can_be_skipped": m.AutoTriggeredAndPreviousEnvironmentApprovedCanBeSkipped,
		"enforce_identity_revalidation":                                   m.EnforceIdentityRevalidation,
		"release_creator_can_be_approver":                                 m.ReleaseCreatorCanBeApprover,
		"required_approver_count":                                         m.RequiredApproverCount,
	}}
}

func flattenReleaseDefinitionApprovals(m *release.ReleaseDefinitionApprovals) []map[string]interface{} {
	timeoutInMinutes := 0
	if m.ApprovalOptions != nil {
		timeoutInMinutes = *m.ApprovalOptions.TimeoutInMinutes
	}
	return []map[string]interface{}{{
		"approval":           flattenReleaseDefinitionApprovalStepList(m.Approvals),
		"timeout_in_minutes": timeoutInMinutes,
	}}
}

func flattenReleaseEnvironmentRetentionPolicy(m *release.EnvironmentRetentionPolicy) []map[string]interface{} {
	return []map[string]interface{}{{
		"retain_build":     m.RetainBuild,
		"releases_to_keep": m.ReleasesToKeep,
		"days_to_keep":     m.DaysToKeep,
	}}
}

func flattenReleaseDefinitionGate(m release.ReleaseDefinitionGate) map[string]interface{} {
	return map[string]interface{}{
		"task": flattenReleaseWorkflowTaskList(m.Tasks),
	}
}
func flattenReleaseDefinitionGateList(m *[]release.ReleaseDefinitionGate) []map[string]interface{} {
	ds := make([]map[string]interface{}, 0, 0)
	for _, d := range *m {
		ds = append(ds, flattenReleaseDefinitionGate(d))
	}
	return ds
}

func flattenReleaseReleaseDefinitionGatesOptions(m *release.ReleaseDefinitionGatesOptions) []map[string]interface{} {
	if m == nil {
		return nil
	}
	return []map[string]interface{}{{
		"is_enabled":               m.IsEnabled,
		"minimum_success_duration": m.MinimumSuccessDuration,
		"sampling_interval":        m.SamplingInterval,
		"stabilization_time":       m.StabilizationTime,
		"timeout":                  m.Timeout,
	}}
}
func flattenReleaseDefinitionGatesStep(m *release.ReleaseDefinitionGatesStep) []map[string]interface{} {
	return []map[string]interface{}{{
		"id":            m.Id,
		"gate":          flattenReleaseDefinitionGateList(m.Gates),
		"gates_options": flattenReleaseReleaseDefinitionGatesOptions(m.GatesOptions),
	}}
}

func flattenCondition(m release.Condition, t release.ConditionType) map[string]interface{} {
	switch t {
	case release.ConditionTypeValues.EnvironmentState:
		return map[string]interface{}{
			"stage_name": m.Name,
			"trigger_even_when_stages_partially_succeed": *m.Value != "4",
		}
	case release.ConditionTypeValues.Event:
		return map[string]interface{}{
			"event_name": m.Name,
		}
	}
	return nil
}
func flattenConditionList(m *[]release.Condition, t release.ConditionType) []map[string]interface{} {
	ds := make([]map[string]interface{}, 0, 0)
	for _, d := range *m {
		if *d.ConditionType == t {
			ds = append(ds, flattenCondition(d, t))
		}
	}
	return ds
}

func flattenArtifactFilterInclude(m map[string]interface{}, isInclude bool) map[string]interface{} {
	ds := map[string]interface{}{"branch_name": m["sourceBranch"]}
	if isInclude {
		ds["tags"] = m["tags"]
	}
	return ds
}

func flattenArtifactFilterIncludeExcludeList(m []release.Condition, isInclude bool) []map[string]interface{} {
	ds := make([]map[string]interface{}, 0, 0)
	for _, d := range m {
		var m2 map[string]interface{}
		if err := json.Unmarshal([]byte(*d.Value), &m2); err == nil {
			if branch, ok := m2["sourceBranch"].(string); ok {
				if isInclude == !strings.HasPrefix(branch, "-") {
					ds = append(ds, flattenArtifactFilterInclude(m2, isInclude))
				}
			}
		}
	}
	return ds
}

func filterReleaseCondition(m *[]release.Condition, t release.ConditionType) []release.Condition {
	ds := make([]release.Condition, 0, 0)
	for _, d := range *m {
		if *d.ConditionType == t {
			ds = append(ds, d)
		}
	}
	return ds
}

func flattenArtifactFilter(m []release.Condition, alias string) map[string]interface{} {
	return map[string]interface{}{
		"artifact_alias": alias,
		"include":        flattenArtifactFilterIncludeExcludeList(m, true),
		"exclude":        flattenArtifactFilterIncludeExcludeList(m, false),
	}
}
func flattenArtifactFilterGroup(m map[string][]release.Condition) []map[string]interface{} {
	ds := make([]map[string]interface{}, 0, len(m))
	for k, d := range m {
		ds = append(ds, flattenArtifactFilter(d, k))
	}
	return ds
}
func flattenArtifactFilterList(m *[]release.Condition) []map[string]interface{} {
	artifacts := filterReleaseCondition(m, release.ConditionTypeValues.Artifact)
	groups := map[string][]release.Condition{}
	for _, v := range artifacts {
		groups[*v.Name] = append(groups[*v.Name], v)
	}
	return flattenArtifactFilterGroup(groups)
}

func flattenReleaseDefinitionEnvironment(m release.ReleaseDefinitionEnvironment) interface{} {
	var ownerID *string
	if m.Owner != nil {
		ownerID = m.Owner.Id
	}
	return map[string]interface{}{
		"id":                   m.Id,
		"rank":                 m.Rank,
		"name":                 m.Name,
		"owner_id":             ownerID,
		"variable":             flattenReleaseDefinitionVariables(m.Variables),
		"variable_groups":      m.VariableGroups,
		"approval_options":     flattenApprovalOptions(m.PreDeployApprovals.ApprovalOptions, m.PostDeployApprovals.ApprovalOptions),
		"pre_deploy_approval":  flattenReleaseDefinitionApprovals(m.PreDeployApprovals),
		"deploy_step":          flattenReleaseDefinitionDeployStep(m.DeployStep),
		"post_deploy_approval": flattenReleaseDefinitionApprovals(m.PostDeployApprovals),
		"retention_policy":     flattenReleaseEnvironmentRetentionPolicy(m.RetentionPolicy),
		"pre_deploy_gate":      flattenReleaseDefinitionGatesStep(m.PreDeploymentGates),
		"post_deploy_gate":     flattenReleaseDefinitionGatesStep(m.PostDeploymentGates),
		"after_stage":          flattenConditionList(m.Conditions, release.ConditionTypeValues.EnvironmentState),
		"after_release":        flattenConditionList(m.Conditions, release.ConditionTypeValues.Event),
		"artifact_filter":      flattenArtifactFilterList(m.Conditions),
		"job":                  flattenReleaseDeployPhasesList(m.DeployPhases),
		"properties":           flattenReleaseEnvironmentProperties(m.Properties),
		"environment_options":  flattenReleaseDefinitionEnvironmentOptions(m.EnvironmentOptions),
	}
}
func flattenReleaseDefinitionEnvironmentList(m *[]release.ReleaseDefinitionEnvironment) []interface{} {
	m2 := *m
	sort.Slice(m2, func(i, j int) bool {
		return *m2[i].Rank < *m2[j].Rank
	})
	ds := make([]interface{}, 0, len(*m))
	for _, d := range m2 {
		ds = append(ds, flattenReleaseDefinitionEnvironment(d))
	}
	return ds
}

func flattenReleaseDefinitionTriggers(m interface{}) interface{} {
	return map[string]interface{}{}
}
func flattenReleaseDefinitionTriggersList(m *[]interface{}) []interface{} {
	ds := make([]interface{}, 0, len(*m))
	for _, d := range *m {
		ds = append(ds, flattenReleaseDefinitionTriggers(d))
	}
	return ds
}

func flattenReleaseDefinitionBuildArtifacts(m release.Artifact) map[string]interface{} {
	dr := *m.DefinitionReference
	ds := map[string]interface{}{
		"project_id":        dr["project"].Id,
		"build_pipeline_id": dr["definition"].Id,
		"alias":             m.Alias,
		"is_primary":        m.IsPrimary,
		"is_retained":       m.IsRetained,
	}
	switch *dr["defaultVersionType"].Id {
	case string(DefaultVersionTypeValues.Latest):
		fallthrough
	case string(DefaultVersionTypeValues.LatestWithBranchAndTags):
		fallthrough
	case string(DefaultVersionTypeValues.LatestWithBuildDefinitionBranchAndTags):
		ds["latest"] = []map[string]interface{}{{
			"branch": dr["defaultVersionBranch"].Id,
			"tags":   dr["defaultVersionTags"].Id,
		}}
		break
	case string(DefaultVersionTypeValues.SelectDuringReleaseCreation):
		fallthrough
	case string(DefaultVersionTypeValues.SpecificVersion):
		ds["specify"] = []map[string]interface{}{{
			"version": dr["defaultVersionSpecific"].Id,
		}}
		break
	}
	return ds
}
func flattenReleaseDefinitionBuildArtifactsList(m *[]release.Artifact) []map[string]interface{} {
	ds := make([]map[string]interface{}, 0, len(*m))
	for _, d := range *m {
		ds = append(ds, flattenReleaseDefinitionBuildArtifacts(d))
	}
	return ds
}

func flattenReleaseDefinitionArtifactsList(m *[]release.Artifact, t release.AgentArtifactType) []map[string]interface{} {
	switch t {
	case release.AgentArtifactTypeValues.Build:
		return flattenReleaseDefinitionBuildArtifactsList(m)
	}
	return nil
}
