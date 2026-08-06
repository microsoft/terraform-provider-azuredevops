package release

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/release"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/validate"
)

const (
	rdVariable          = "variable"
	rdVariableName      = "name"
	rdVariableValue     = "value"
	rdVariableSecretVal = "secret_value"
	rdVariableIsSecret  = "is_secret"
	rdVariableCanOverwr = "allow_override"
)

func ResourceReleaseDefinition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReleaseDefinitionCreate,
		ReadContext:   resourceReleaseDefinitionRead,
		UpdateContext: resourceReleaseDefinitionUpdate,
		DeleteContext: resourceReleaseDefinitionDelete,
		Importer:      tfhelper.ImportProjectQualifiedResource(),
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      `\`,
				ValidateFunc: validate.Path,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"release_name_format": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Release-$(rev:r)",
			},
			"variable_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntAtLeast(1),
				},
			},
			rdVariable: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						rdVariableName: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						rdVariableValue: {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						rdVariableSecretVal: {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
							Default:   "",
						},
						rdVariableIsSecret: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						rdVariableCanOverwr: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"artifact": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alias": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Build",
								"Jenkins",
								"GitHub",
								"Nuget",
								"Team Build (external)",
								"ExternalTFSBuild",
								"Git",
								"TFVC",
								"ExternalTfsXamlBuild",
							}, false),
						},
						"is_primary": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"is_retained": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"definition_reference": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
									},
									"id": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
									},
								},
							},
						},
					},
				},
			},
			"artifact_source_trigger": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"artifact_alias": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"trigger_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_branch": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
									},
									"use_build_definition_branch": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"create_release_on_build_tagging": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"tags": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"schedule_trigger": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"days": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"monday", "tuesday", "wednesday", "thursday",
									"friday", "saturday", "sunday", "all",
								}, true),
								StateFunc: func(v interface{}) string {
									return strings.ToLower(v.(string))
								},
							},
						},
						"start_hours": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      3,
							ValidateFunc: validation.IntBetween(0, 23),
						},
						"start_minutes": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 59),
						},
						"time_zone_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "UTC",
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"only_with_changes": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"stage": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"rank": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"condition": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"condition_type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "event",
										ValidateFunc: validation.StringInSlice([]string{
											"event", "environmentState", "artifact",
										}, false),
									},
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
									},
								},
							},
						},
						"variable_groups": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								ValidateFunc: validation.IntAtLeast(1),
							},
						},
						rdVariable: {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									rdVariableName: {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
									},
									rdVariableValue: {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
									},
									rdVariableSecretVal: {
										Type:      schema.TypeString,
										Optional:  true,
										Sensitive: true,
										Default:   "",
									},
									rdVariableIsSecret: {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									rdVariableCanOverwr: {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
						"retention_policy": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days_to_keep": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      30,
										ValidateFunc: validation.IntAtLeast(0),
									},
									"releases_to_keep": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      3,
										ValidateFunc: validation.IntAtLeast(0),
									},
									"retain_build": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
								},
							},
						},
						"pre_deploy_approval":   approvalSchema(),
						"post_deploy_approval":  approvalSchema(),
						"pre_deployment_gates":  gatesSchema(),
						"post_deployment_gates": gatesSchema(),
						"environment_options": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
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
						},
						"execution_policy": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"concurrency_count": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      0,
										ValidateFunc: validation.IntAtLeast(0),
									},
									"queue_depth_count": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      0,
										ValidateFunc: validation.IntAtLeast(0),
									},
								},
							},
						},
						"environment_trigger": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"trigger_type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  string(release.EnvironmentTriggerTypeValues.RollbackRedeploy),
										ValidateFunc: validation.StringInSlice([]string{
											string(release.EnvironmentTriggerTypeValues.RollbackRedeploy),
											string(release.EnvironmentTriggerTypeValues.DeploymentGroupRedeploy),
										}, false),
									},
									"action": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
									},
									"event_types": {
										Type:     schema.TypeList,
										Required: true,
										MinItems: 1,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"deploy_phase": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
									},
									"phase_type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  string(release.DeployPhaseTypesValues.AgentBasedDeployment),
										ValidateFunc: validation.StringInSlice([]string{
											string(release.DeployPhaseTypesValues.AgentBasedDeployment),
											string(release.DeployPhaseTypesValues.RunOnServer),
											string(release.DeployPhaseTypesValues.MachineGroupBasedDeployment),
										}, false),
									},
									"rank": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"agent_queue_id": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},
									"agent_specification": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotWhiteSpace,
									},
									"skip_artifacts_download": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"timeout_in_minutes": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      0,
										ValidateFunc: validation.IntAtLeast(0),
									},
									"job_cancel_timeout_in_minutes": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      1,
										ValidateFunc: validation.IntAtLeast(0),
									},
									"enable_access_token": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"condition": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "succeeded()",
									},
									"deployment_group_id": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},
									"deployment_health_option": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
									},
									"health_percent": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      0,
										ValidateFunc: validation.IntBetween(0, 100),
									},
									"tags": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"task": taskSchema(),
								},
							},
						},
					},
				},
			},
			"revision": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func taskSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"task_id": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.IsUUID,
					StateFunc: func(v interface{}) string {
						return strings.ToLower(v.(string))
					},
				},
				"version": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "1.*",
				},
				"name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
				"condition": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "succeeded()",
				},
				"inputs": {
					Type:     schema.TypeMap,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"override_inputs": {
					Type:     schema.TypeMap,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"environment": {
					Type:     schema.TypeMap,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"always_run": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"continue_on_error": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"timeout_in_minutes": {
					Type:         schema.TypeInt,
					Optional:     true,
					Default:      0,
					ValidateFunc: validation.IntAtLeast(0),
				},
				"retry_count_on_task_failure": {
					Type:         schema.TypeInt,
					Optional:     true,
					Default:      0,
					ValidateFunc: validation.IntAtLeast(0),
				},
				"ref_name": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
				"definition_type": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func gatesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"is_enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
				"sampling_interval": {
					Type:         schema.TypeInt,
					Optional:     true,
					Default:      15,
					ValidateFunc: validation.IntAtLeast(0),
				},
				"stabilization_time": {
					Type:         schema.TypeInt,
					Optional:     true,
					Default:      0,
					ValidateFunc: validation.IntAtLeast(0),
				},
				"minimum_success_duration": {
					Type:         schema.TypeInt,
					Optional:     true,
					Default:      0,
					ValidateFunc: validation.IntAtLeast(0),
				},
				"timeout": {
					Type:         schema.TypeInt,
					Optional:     true,
					Default:      1440,
					ValidateFunc: validation.IntAtLeast(0),
				},
				"gate": {
					Type:     schema.TypeList,
					Optional: true,
					MinItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"task": taskSchema(),
						},
					},
				},
			},
		},
	}
}

func approvalSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"approval": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"approver_id": {
								Type:         schema.TypeString,
								Optional:     true,
								ValidateFunc: validation.IsUUID,
							},
							"is_automated": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  false,
							},
							"is_notification_on": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  false,
							},
						},
					},
				},
				"approval_options": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"required_approver_count": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntAtLeast(0),
							},
							"release_creator_can_be_approver": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  false,
							},
							"auto_triggered_and_previous_environment_approved_can_be_skipped": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  false,
							},
							"enforce_identity_revalidation": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  false,
							},
							"timeout_in_minutes": {
								Type:         schema.TypeInt,
								Optional:     true,
								Default:      0,
								ValidateFunc: validation.IntBetween(0, 525600),
							},
							"execution_order": {
								Type:     schema.TypeString,
								Optional: true,
								Default:  string(release.ApprovalExecutionOrderValues.BeforeGates),
								ValidateFunc: validation.StringInSlice([]string{
									string(release.ApprovalExecutionOrderValues.BeforeGates),
									string(release.ApprovalExecutionOrderValues.AfterSuccessfulGates),
									string(release.ApprovalExecutionOrderValues.AfterGatesAlways),
								}, false),
							},
						},
					},
				},
			},
		},
	}
}

func resourceReleaseDefinitionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	if err := validateReleaseVariables(d); err != nil {
		return diag.FromErr(err)
	}
	releaseDefinition, projectID := expandReleaseDefinition(d)

	createdReleaseDefinition, err := clients.ReleaseClient.CreateReleaseDefinition(ctx, release.CreateReleaseDefinitionArgs{
		ReleaseDefinition: releaseDefinition,
		Project:           &projectID,
	})
	if err != nil {
		return diag.Errorf("creating release definition: %+v", err)
	}

	d.SetId(strconv.Itoa(*createdReleaseDefinition.Id))
	return resourceReleaseDefinitionRead(ctx, d, m)
}

func resourceReleaseDefinitionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	projectID, releaseDefinitionID, err := tfhelper.ParseProjectIDAndResourceID(d)
	if err != nil {
		return diag.FromErr(err)
	}

	releaseDefinition, err := clients.ReleaseClient.GetReleaseDefinition(ctx, release.GetReleaseDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &releaseDefinitionID,
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return diag.FromErr(flattenReleaseDefinition(d, releaseDefinition, projectID))
}

func resourceReleaseDefinitionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	if err := validateReleaseVariables(d); err != nil {
		return diag.FromErr(err)
	}
	releaseDefinition, projectID := expandReleaseDefinition(d)

	existing, err := clients.ReleaseClient.GetReleaseDefinition(ctx, release.GetReleaseDefinitionArgs{
		Project:      &projectID,
		DefinitionId: releaseDefinition.Id,
	})
	if err != nil {
		return diag.Errorf("reading release definition before update: %+v", err)
	}
	releaseDefinition.Triggers = mergeUnmanagedTriggers(releaseDefinition.Triggers, existing.Triggers)

	_, err = clients.ReleaseClient.UpdateReleaseDefinition(ctx, release.UpdateReleaseDefinitionArgs{
		ReleaseDefinition: releaseDefinition,
		Project:           &projectID,
	})
	if err != nil {
		return diag.Errorf("updating release definition: %+v", err)
	}

	return resourceReleaseDefinitionRead(ctx, d, m)
}

func mergeUnmanagedTriggers(expanded, existing *[]interface{}) *[]interface{} {
	if existing == nil {
		return expanded
	}
	managed := map[string]bool{
		string(release.ReleaseTriggerTypeValues.ArtifactSource): true,
		string(release.ReleaseTriggerTypeValues.Schedule):       true,
	}
	result := make([]interface{}, 0)
	if expanded != nil {
		result = append(result, (*expanded)...)
	}
	for _, item := range *existing {
		trigger, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if !managed[stringFromMap(trigger, "triggerType")] {
			result = append(result, item)
		}
	}
	return &result
}

func resourceReleaseDefinitionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	projectID, releaseDefinitionID, err := tfhelper.ParseProjectIDAndResourceID(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = clients.ReleaseClient.DeleteReleaseDefinition(ctx, release.DeleteReleaseDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &releaseDefinitionID,
	})
	if err != nil {
		return diag.Errorf("deleting release definition: %+v", err)
	}

	d.SetId("")
	return nil
}

func expandReleaseDefinition(d *schema.ResourceData) (*release.ReleaseDefinition, string) {
	projectID := d.Get("project_id").(string)

	releaseDefinition := &release.ReleaseDefinition{
		Name:              converter.String(d.Get("name").(string)),
		Path:              converter.String(d.Get("path").(string)),
		Description:       converter.String(d.Get("description").(string)),
		ReleaseNameFormat: converter.String(d.Get("release_name_format").(string)),
		Variables:         expandReleaseVariables(d.Get(rdVariable).(*schema.Set).List()),
		VariableGroups:    expandReleaseVariableGroups(d.Get("variable_groups").(*schema.Set).List()),
		Artifacts:         expandReleaseArtifacts(d.Get("artifact").(*schema.Set).List()),
		Environments:      expandReleaseStages(d.Get("stage").([]interface{})),
		Triggers:          expandReleaseTriggers(d.Get("artifact_source_trigger").([]interface{}), d.Get("schedule_trigger").([]interface{})),
	}

	if !d.IsNewResource() {
		if id, err := strconv.Atoi(d.Id()); err == nil {
			releaseDefinition.Id = &id
		}
		if revision, ok := d.GetOk("revision"); ok {
			releaseDefinition.Revision = converter.Int(revision.(int))
		}
	}

	return releaseDefinition, projectID
}

func expandReleaseVariables(input []interface{}) *map[string]release.ConfigurationVariableValue {
	variables := map[string]release.ConfigurationVariableValue{}
	for _, item := range input {
		variable := item.(map[string]interface{})
		isSecret := variable[rdVariableIsSecret].(bool)
		value := variable[rdVariableValue].(string)
		if isSecret {
			value = variable[rdVariableSecretVal].(string)
		}
		variables[variable[rdVariableName].(string)] = release.ConfigurationVariableValue{
			Value:         converter.String(value),
			IsSecret:      converter.Bool(isSecret),
			AllowOverride: converter.Bool(variable[rdVariableCanOverwr].(bool)),
		}
	}
	return &variables
}

func validateReleaseVariables(d *schema.ResourceData) error {
	raw := d.GetRawConfig()
	if raw.IsNull() {
		return nil
	}
	config := raw.AsValueMap()

	if variables := config[rdVariable]; !variables.IsNull() {
		for it := variables.ElementIterator(); it.Next(); {
			_, variable := it.Element()
			attrs := variable.AsValueMap()
			if err := validateVariableExclusivity(
				attrs[rdVariableName].AsString(),
				!attrs[rdVariableValue].IsNull(),
				!attrs[rdVariableSecretVal].IsNull(),
				!attrs[rdVariableIsSecret].IsNull(),
			); err != nil {
				return err
			}
		}
	}

	stages := config["stage"]
	if stages.IsNull() {
		return nil
	}
	for sit := stages.ElementIterator(); sit.Next(); {
		_, stage := sit.Element()
		variables := stage.AsValueMap()[rdVariable]
		if variables.IsNull() {
			continue
		}
		for it := variables.ElementIterator(); it.Next(); {
			_, variable := it.Element()
			attrs := variable.AsValueMap()
			if err := validateVariableExclusivity(
				attrs[rdVariableName].AsString(),
				!attrs[rdVariableValue].IsNull(),
				!attrs[rdVariableSecretVal].IsNull(),
				!attrs[rdVariableIsSecret].IsNull(),
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateVariableExclusivity(name string, valueSet, secretValueSet, isSecret bool) error {
	if valueSet && (secretValueSet || isSecret) || secretValueSet != isSecret {
		return fmt.Errorf("`%s` variable can have either only `value` attribute or both `is_secret` and `secret_value` attributes", name)
	}
	return nil
}

func expandReleaseVariableGroups(input []interface{}) *[]int {
	groups := make([]int, 0, len(input))
	for _, item := range input {
		groups = append(groups, item.(int))
	}
	return &groups
}

func expandReleaseArtifacts(input []interface{}) *[]release.Artifact {
	artifacts := make([]release.Artifact, 0, len(input))
	for _, item := range input {
		artifact := item.(map[string]interface{})
		artifacts = append(artifacts, release.Artifact{
			Alias:               converter.String(artifact["alias"].(string)),
			Type:                converter.String(artifact["type"].(string)),
			IsPrimary:           converter.Bool(artifact["is_primary"].(bool)),
			IsRetained:          converter.Bool(artifact["is_retained"].(bool)),
			DefinitionReference: expandArtifactDefinitionReference(artifact["definition_reference"].(*schema.Set).List()),
		})
	}
	return &artifacts
}

func expandArtifactDefinitionReference(input []interface{}) *map[string]release.ArtifactSourceReference {
	refs := map[string]release.ArtifactSourceReference{}
	for _, item := range input {
		ref := item.(map[string]interface{})
		refs[ref["key"].(string)] = release.ArtifactSourceReference{
			Id:   converter.String(ref["id"].(string)),
			Name: converter.String(ref["name"].(string)),
		}
	}
	return &refs
}

func expandReleaseTriggers(artifactSourceInput, scheduleInput []interface{}) *[]interface{} {
	triggers := make([]interface{}, 0, len(artifactSourceInput)+len(scheduleInput))

	for _, item := range artifactSourceInput {
		trigger := item.(map[string]interface{})
		triggers = append(triggers, release.ArtifactSourceTrigger{
			TriggerType:       &release.ReleaseTriggerTypeValues.ArtifactSource,
			ArtifactAlias:     converter.String(trigger["artifact_alias"].(string)),
			TriggerConditions: expandArtifactFilters(trigger["trigger_condition"].([]interface{})),
		})
	}

	for _, item := range scheduleInput {
		trigger := item.(map[string]interface{})
		days := release.ScheduleDays(joinScheduleDays(trigger["days"].(*schema.Set).List()))
		triggers = append(triggers, release.ScheduledReleaseTrigger{
			TriggerType: &release.ReleaseTriggerTypeValues.Schedule,
			Schedule: &release.ReleaseSchedule{
				DaysToRelease:           &days,
				StartHours:              converter.Int(trigger["start_hours"].(int)),
				StartMinutes:            converter.Int(trigger["start_minutes"].(int)),
				TimeZoneId:              converter.String(trigger["time_zone_id"].(string)),
				ScheduleOnlyWithChanges: converter.Bool(trigger["only_with_changes"].(bool)),
			},
		})
	}

	return &triggers
}

func expandArtifactFilters(input []interface{}) *[]release.ArtifactFilter {
	filters := make([]release.ArtifactFilter, 0, len(input))
	for _, item := range input {
		filter := item.(map[string]interface{})
		filters = append(filters, release.ArtifactFilter{
			SourceBranch:                converter.String(filter["source_branch"].(string)),
			UseBuildDefinitionBranch:    converter.Bool(filter["use_build_definition_branch"].(bool)),
			CreateReleaseOnBuildTagging: converter.Bool(filter["create_release_on_build_tagging"].(bool)),
			Tags:                        expandStringSet(filter["tags"].(*schema.Set).List()),
		})
	}
	return &filters
}

func expandStringSet(input []interface{}) *[]string {
	values := make([]string, 0, len(input))
	for _, item := range input {
		values = append(values, item.(string))
	}
	return &values
}

func joinScheduleDays(input []interface{}) string {
	days := make([]string, 0, len(input))
	for _, item := range input {
		days = append(days, strings.ToLower(item.(string)))
	}
	return strings.Join(days, ", ")
}

func flattenReleaseDefinition(d *schema.ResourceData, releaseDefinition *release.ReleaseDefinition, projectID string) error {
	d.Set("project_id", projectID)
	d.Set("name", converter.ToString(releaseDefinition.Name, ""))
	d.Set("path", converter.ToString(releaseDefinition.Path, `\`))
	d.Set("description", converter.ToString(releaseDefinition.Description, ""))
	d.Set("release_name_format", converter.ToString(releaseDefinition.ReleaseNameFormat, ""))
	d.Set("revision", converter.ToInt(releaseDefinition.Revision, 0))

	if err := d.Set(rdVariable, flattenReleaseVariables(releaseDefinition.Variables, d.Get(rdVariable).(*schema.Set).List())); err != nil {
		return fmt.Errorf("setting `variable`: %+v", err)
	}
	if err := d.Set("variable_groups", flattenReleaseVariableGroups(releaseDefinition.VariableGroups)); err != nil {
		return fmt.Errorf("setting `variable_groups`: %+v", err)
	}
	if err := d.Set("artifact", flattenReleaseArtifacts(releaseDefinition.Artifacts, artifactDefinitionKeysFromConfig(d))); err != nil {
		return fmt.Errorf("setting `artifact`: %+v", err)
	}
	artifactSourceTriggers, scheduleTriggers := flattenReleaseTriggers(releaseDefinition.Triggers, d.Get("schedule_trigger").([]interface{}))
	if err := d.Set("artifact_source_trigger", artifactSourceTriggers); err != nil {
		return fmt.Errorf("setting `artifact_source_trigger`: %+v", err)
	}
	if err := d.Set("schedule_trigger", scheduleTriggers); err != nil {
		return fmt.Errorf("setting `schedule_trigger`: %+v", err)
	}
	if err := d.Set("stage", flattenReleaseStages(d, releaseDefinition.Environments)); err != nil {
		return fmt.Errorf("setting `stage`: %+v", err)
	}

	return nil
}

func flattenReleaseVariables(input *map[string]release.ConfigurationVariableValue, stateVars []interface{}) []interface{} {
	if input == nil {
		return nil
	}
	variables := make([]interface{}, 0, len(*input))
	for name, value := range *input {
		isSecret := converter.ToBool(value.IsSecret, false)
		variable := map[string]interface{}{
			rdVariableName:      name,
			rdVariableValue:     converter.ToString(value.Value, ""),
			rdVariableSecretVal: "",
			rdVariableIsSecret:  isSecret,
			rdVariableCanOverwr: converter.ToBool(value.AllowOverride, false),
		}
		if isSecret {
			variable[rdVariableValue] = ""
			if prior := findVariableInState(stateVars, name); prior != nil {
				variable[rdVariableSecretVal] = prior[rdVariableSecretVal]
			}
		}
		variables = append(variables, variable)
	}
	return variables
}

func findVariableInState(stateVars []interface{}, name string) map[string]interface{} {
	for _, item := range stateVars {
		if variable, ok := item.(map[string]interface{}); ok && variable[rdVariableName] == name {
			return variable
		}
	}
	return nil
}

func flattenReleaseVariableGroups(input *[]int) []interface{} {
	if input == nil {
		return nil
	}
	groups := make([]interface{}, 0, len(*input))
	for _, group := range *input {
		groups = append(groups, group)
	}
	return groups
}

func flattenReleaseArtifacts(input *[]release.Artifact, configuredKeys map[string]map[string]bool) []interface{} {
	if input == nil {
		return nil
	}
	artifacts := make([]interface{}, 0, len(*input))
	for _, artifact := range *input {
		alias := converter.ToString(artifact.Alias, "")
		allowedKeys := configuredKeys[alias]
		artifacts = append(artifacts, map[string]interface{}{
			"alias":                alias,
			"type":                 converter.ToString(artifact.Type, ""),
			"is_primary":           converter.ToBool(artifact.IsPrimary, false),
			"is_retained":          converter.ToBool(artifact.IsRetained, false),
			"definition_reference": flattenArtifactDefinitionReference(artifact.DefinitionReference, allowedKeys),
		})
	}
	return artifacts
}

func artifactDefinitionKeysFromConfig(d *schema.ResourceData) map[string]map[string]bool {
	result := map[string]map[string]bool{}
	for _, item := range d.Get("artifact").(*schema.Set).List() {
		artifact := item.(map[string]interface{})
		keys := map[string]bool{}
		for _, ref := range artifact["definition_reference"].(*schema.Set).List() {
			keys[ref.(map[string]interface{})["key"].(string)] = true
		}
		result[artifact["alias"].(string)] = keys
	}
	return result
}

func flattenArtifactDefinitionReference(input *map[string]release.ArtifactSourceReference, allowedKeys map[string]bool) []interface{} {
	if input == nil {
		return nil
	}
	refs := make([]interface{}, 0, len(*input))
	for key, ref := range *input {
		if allowedKeys != nil && !allowedKeys[key] {
			continue
		}
		refs = append(refs, map[string]interface{}{
			"key":  key,
			"id":   converter.ToString(ref.Id, ""),
			"name": converter.ToString(ref.Name, ""),
		})
	}
	return refs
}

func flattenReleaseTriggers(input *[]interface{}, priorSchedules []interface{}) ([]interface{}, []interface{}) {
	artifactSourceTriggers := make([]interface{}, 0)
	scheduleTriggers := make([]interface{}, 0)
	if input == nil {
		return artifactSourceTriggers, scheduleTriggers
	}

	for _, item := range *input {
		trigger, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		switch stringFromMap(trigger, "triggerType") {
		case string(release.ReleaseTriggerTypeValues.ArtifactSource):
			artifactSourceTriggers = append(artifactSourceTriggers, map[string]interface{}{
				"artifact_alias":    stringFromMap(trigger, "artifactAlias"),
				"trigger_condition": flattenArtifactFilters(trigger["triggerConditions"]),
			})
		case string(release.ReleaseTriggerTypeValues.Schedule):
			schedule, _ := trigger["schedule"].(map[string]interface{})
			days := reconcileScheduleDays(
				flattenScheduleDays(schedule["daysToRelease"]),
				priorScheduleDays(priorSchedules, len(scheduleTriggers)),
			)
			scheduleTriggers = append(scheduleTriggers, map[string]interface{}{
				"days":              days,
				"start_hours":       intFromMap(schedule, "startHours"),
				"start_minutes":     intFromMap(schedule, "startMinutes"),
				"time_zone_id":      stringFromMap(schedule, "timeZoneId"),
				"only_with_changes": boolFromMap(schedule, "scheduleOnlyWithChanges"),
			})
		}
	}

	return artifactSourceTriggers, scheduleTriggers
}

func priorScheduleDays(priorSchedules []interface{}, index int) []interface{} {
	if index >= len(priorSchedules) {
		return nil
	}
	schedule, ok := priorSchedules[index].(map[string]interface{})
	if !ok {
		return nil
	}
	if set, ok := schedule["days"].(*schema.Set); ok {
		return set.List()
	}
	return nil
}

func reconcileScheduleDays(flattened, prior []interface{}) []interface{} {
	if len(flattened) == 1 && flattened[0] == "all" && scheduleDaysAreAllSeven(prior) {
		return prior
	}
	return flattened
}

func scheduleDaysAreAllSeven(days []interface{}) bool {
	if len(days) != len(scheduleDayBits) {
		return false
	}
	seen := map[string]bool{}
	for _, d := range days {
		if s, ok := d.(string); ok {
			seen[strings.ToLower(s)] = true
		}
	}
	for _, d := range scheduleDayBits {
		if !seen[d.name] {
			return false
		}
	}
	return true
}

func flattenArtifactFilters(input interface{}) []interface{} {
	rawFilters, ok := input.([]interface{})
	if !ok {
		return nil
	}
	filters := make([]interface{}, 0, len(rawFilters))
	for _, item := range rawFilters {
		filter, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		result := map[string]interface{}{
			"source_branch":                   stringFromMap(filter, "sourceBranch"),
			"use_build_definition_branch":     boolFromMap(filter, "useBuildDefinitionBranch"),
			"create_release_on_build_tagging": boolFromMap(filter, "createReleaseOnBuildTagging"),
		}
		if tags, ok := filter["tags"].([]interface{}); ok {
			result["tags"] = tags
		}
		filters = append(filters, result)
	}
	return filters
}

var scheduleDayBits = []struct {
	bit  int
	name string
}{
	{1, "monday"}, {2, "tuesday"}, {4, "wednesday"}, {8, "thursday"},
	{16, "friday"}, {32, "saturday"}, {64, "sunday"},
}

func flattenScheduleDays(input interface{}) []interface{} {
	switch v := input.(type) {
	case float64:
		return scheduleDaysFromMask(int(v))
	case string:
		if mask, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			return scheduleDaysFromMask(mask)
		}
		return splitScheduleDays(v)
	default:
		return nil
	}
}

func scheduleDaysFromMask(mask int) []interface{} {
	allBits := 0
	for _, d := range scheduleDayBits {
		allBits |= d.bit
	}
	if mask == allBits {
		return []interface{}{"all"}
	}
	days := make([]interface{}, 0, len(scheduleDayBits))
	for _, d := range scheduleDayBits {
		if mask&d.bit != 0 {
			days = append(days, d.name)
		}
	}
	return days
}

func splitScheduleDays(input string) []interface{} {
	if input == "" {
		return nil
	}
	parts := strings.Split(input, ",")
	days := make([]interface{}, 0, len(parts))
	for _, part := range parts {
		if day := strings.ToLower(strings.TrimSpace(part)); day != "" {
			days = append(days, day)
		}
	}
	return days
}

func expandReleaseStages(input []interface{}) *[]release.ReleaseDefinitionEnvironment {
	stages := make([]release.ReleaseDefinitionEnvironment, 0, len(input))
	for i, item := range input {
		stage := item.(map[string]interface{})

		rank := stage["rank"].(int)
		if rank == 0 {
			rank = i + 1
		}

		environment := release.ReleaseDefinitionEnvironment{
			Name:                converter.String(stage["name"].(string)),
			Rank:                converter.Int(rank),
			Conditions:          expandReleaseConditions(stage["condition"].([]interface{}), i == 0),
			Variables:           expandReleaseVariables(stage[rdVariable].(*schema.Set).List()),
			VariableGroups:      expandReleaseVariableGroups(stage["variable_groups"].(*schema.Set).List()),
			RetentionPolicy:     expandReleaseRetentionPolicy(stage["retention_policy"].([]interface{})),
			PreDeployApprovals:  expandReleaseApprovals(stage["pre_deploy_approval"].([]interface{})),
			PostDeployApprovals: expandReleaseApprovals(stage["post_deploy_approval"].([]interface{})),
			PreDeploymentGates:  expandReleaseGates(stage["pre_deployment_gates"].([]interface{})),
			PostDeploymentGates: expandReleaseGates(stage["post_deployment_gates"].([]interface{})),
			EnvironmentOptions:  expandEnvironmentOptions(stage["environment_options"].([]interface{})),
			ExecutionPolicy:     expandExecutionPolicy(stage["execution_policy"].([]interface{})),
			EnvironmentTriggers: expandEnvironmentTriggers(stage["environment_trigger"].([]interface{})),
			DeployPhases:        expandReleaseDeployPhases(stage["deploy_phase"].([]interface{})),
			DeployStep:          &release.ReleaseDefinitionDeployStep{Tasks: &[]release.WorkflowTask{}},
		}

		stages = append(stages, environment)
	}
	return &stages
}

func expandReleaseConditions(input []interface{}, isFirstStage bool) *[]release.Condition {
	conditions := make([]release.Condition, 0, len(input))
	for _, item := range input {
		condition := item.(map[string]interface{})
		conditionType := release.ConditionType(condition["condition_type"].(string))
		conditions = append(conditions, release.Condition{
			ConditionType: &conditionType,
			Name:          converter.String(condition["name"].(string)),
			Value:         converter.String(condition["value"].(string)),
		})
	}

	if len(conditions) == 0 && isFirstStage {
		conditions = append(conditions, release.Condition{
			ConditionType: &release.ConditionTypeValues.Event,
			Name:          converter.String("ReleaseStarted"),
			Value:         converter.String(""),
		})
	}

	return &conditions
}

func expandReleaseRetentionPolicy(input []interface{}) *release.EnvironmentRetentionPolicy {
	policy := release.EnvironmentRetentionPolicy{
		DaysToKeep:     converter.Int(30),
		ReleasesToKeep: converter.Int(3),
		RetainBuild:    converter.Bool(true),
	}
	if len(input) == 1 && input[0] != nil {
		item := input[0].(map[string]interface{})
		policy.DaysToKeep = converter.Int(item["days_to_keep"].(int))
		policy.ReleasesToKeep = converter.Int(item["releases_to_keep"].(int))
		policy.RetainBuild = converter.Bool(item["retain_build"].(bool))
	}
	return &policy
}

func expandEnvironmentOptions(input []interface{}) *release.EnvironmentOptions {
	if len(input) != 1 || input[0] == nil {
		return nil
	}
	item := input[0].(map[string]interface{})
	return &release.EnvironmentOptions{
		AutoLinkWorkItems:            converter.Bool(item["auto_link_work_items"].(bool)),
		BadgeEnabled:                 converter.Bool(item["badge_enabled"].(bool)),
		PublishDeploymentStatus:      converter.Bool(item["publish_deployment_status"].(bool)),
		PullRequestDeploymentEnabled: converter.Bool(item["pull_request_deployment_enabled"].(bool)),
	}
}

func expandExecutionPolicy(input []interface{}) *release.EnvironmentExecutionPolicy {
	if len(input) != 1 || input[0] == nil {
		return nil
	}
	item := input[0].(map[string]interface{})
	return &release.EnvironmentExecutionPolicy{
		ConcurrencyCount: converter.Int(item["concurrency_count"].(int)),
		QueueDepthCount:  converter.Int(item["queue_depth_count"].(int)),
	}
}

func expandEnvironmentTriggers(input []interface{}) *[]release.EnvironmentTrigger {
	triggers := make([]release.EnvironmentTrigger, 0, len(input))
	for _, item := range input {
		trigger := item.(map[string]interface{})
		triggerType := release.EnvironmentTriggerType(trigger["trigger_type"].(string))
		content := release.EnvironmentTriggerContent{
			Action:     converter.String(trigger["action"].(string)),
			EventTypes: expandStringSet(trigger["event_types"].([]interface{})),
		}
		contentBytes, err := json.Marshal(content)
		if err != nil {
			continue
		}
		triggers = append(triggers, release.EnvironmentTrigger{
			TriggerType:    &triggerType,
			TriggerContent: converter.String(string(contentBytes)),
		})
	}
	return &triggers
}

func expandReleaseApprovals(input []interface{}) *release.ReleaseDefinitionApprovals {
	steps := make([]release.ReleaseDefinitionApprovalStep, 0)
	var options *release.ApprovalOptions
	if len(input) == 1 && input[0] != nil {
		block := input[0].(map[string]interface{})
		for i, item := range block["approval"].([]interface{}) {
			approval := item.(map[string]interface{})
			step := release.ReleaseDefinitionApprovalStep{
				Rank:             converter.Int(i + 1),
				IsAutomated:      converter.Bool(approval["is_automated"].(bool)),
				IsNotificationOn: converter.Bool(approval["is_notification_on"].(bool)),
			}
			if id := approval["approver_id"].(string); id != "" {
				step.Approver = &webapi.IdentityRef{Id: converter.String(id)}
			}
			steps = append(steps, step)
		}
		options = expandApprovalOptions(block["approval_options"].([]interface{}))
	}

	if len(steps) == 0 {
		steps = append(steps, release.ReleaseDefinitionApprovalStep{
			Rank:        converter.Int(1),
			IsAutomated: converter.Bool(true),
		})
	}

	return &release.ReleaseDefinitionApprovals{Approvals: &steps, ApprovalOptions: options}
}

func expandApprovalOptions(input []interface{}) *release.ApprovalOptions {
	if len(input) != 1 || input[0] == nil {
		return nil
	}
	item := input[0].(map[string]interface{})
	executionOrder := release.ApprovalExecutionOrder(item["execution_order"].(string))
	return &release.ApprovalOptions{
		RequiredApproverCount:                                   converter.Int(item["required_approver_count"].(int)),
		ReleaseCreatorCanBeApprover:                             converter.Bool(item["release_creator_can_be_approver"].(bool)),
		AutoTriggeredAndPreviousEnvironmentApprovedCanBeSkipped: converter.Bool(item["auto_triggered_and_previous_environment_approved_can_be_skipped"].(bool)),
		EnforceIdentityRevalidation:                             converter.Bool(item["enforce_identity_revalidation"].(bool)),
		TimeoutInMinutes:                                        converter.Int(item["timeout_in_minutes"].(int)),
		ExecutionOrder:                                          &executionOrder,
	}
}

func expandReleaseDeployPhases(input []interface{}) *[]interface{} {
	phases := make([]interface{}, 0, len(input))
	for i, item := range input {
		phase := item.(map[string]interface{})

		rank := phase["rank"].(int)
		if rank == 0 {
			rank = i + 1
		}
		name := converter.String(phase["name"].(string))
		tasks := expandReleaseTasks(phase["task"].([]interface{}))
		phaseType := release.DeployPhaseTypes(phase["phase_type"].(string))

		switch phaseType {
		case release.DeployPhaseTypesValues.RunOnServer:
			phases = append(phases, release.RunOnServerDeployPhase{
				Name:          name,
				Rank:          converter.Int(rank),
				PhaseType:     &release.DeployPhaseTypesValues.RunOnServer,
				WorkflowTasks: tasks,
				DeploymentInput: &release.ServerDeploymentInput{
					Condition:                 converter.String(phase["condition"].(string)),
					TimeoutInMinutes:          converter.Int(phase["timeout_in_minutes"].(int)),
					JobCancelTimeoutInMinutes: converter.Int(phase["job_cancel_timeout_in_minutes"].(int)),
				},
			})
		case release.DeployPhaseTypesValues.MachineGroupBasedDeployment:
			deploymentInput := &release.MachineGroupDeploymentInput{
				Condition:                 converter.String(phase["condition"].(string)),
				TimeoutInMinutes:          converter.Int(phase["timeout_in_minutes"].(int)),
				JobCancelTimeoutInMinutes: converter.Int(phase["job_cancel_timeout_in_minutes"].(int)),
				SkipArtifactsDownload:     converter.Bool(phase["skip_artifacts_download"].(bool)),
				EnableAccessToken:         converter.Bool(phase["enable_access_token"].(bool)),
				Tags:                      expandStringSet(phase["tags"].(*schema.Set).List()),
			}
			if dgID := phase["deployment_group_id"].(int); dgID != 0 {
				deploymentInput.QueueId = converter.Int(dgID)
			}
			if healthOption := phase["deployment_health_option"].(string); healthOption != "" {
				deploymentInput.DeploymentHealthOption = converter.String(healthOption)
				deploymentInput.HealthPercent = converter.Int(phase["health_percent"].(int))
			}
			phases = append(phases, release.MachineGroupBasedDeployPhase{
				Name:            name,
				Rank:            converter.Int(rank),
				PhaseType:       &release.DeployPhaseTypesValues.MachineGroupBasedDeployment,
				WorkflowTasks:   tasks,
				DeploymentInput: deploymentInput,
			})
		default:
			deploymentInput := &release.AgentDeploymentInput{
				SkipArtifactsDownload:     converter.Bool(phase["skip_artifacts_download"].(bool)),
				TimeoutInMinutes:          converter.Int(phase["timeout_in_minutes"].(int)),
				JobCancelTimeoutInMinutes: converter.Int(phase["job_cancel_timeout_in_minutes"].(int)),
				EnableAccessToken:         converter.Bool(phase["enable_access_token"].(bool)),
				Condition:                 converter.String(phase["condition"].(string)),
			}
			if queueID := phase["agent_queue_id"].(int); queueID != 0 {
				deploymentInput.QueueId = converter.Int(queueID)
			}
			if spec := phase["agent_specification"].(string); spec != "" {
				deploymentInput.AgentSpecification = &release.AgentSpecification{Identifier: converter.String(spec)}
			}
			phases = append(phases, release.AgentBasedDeployPhase{
				Name:            name,
				Rank:            converter.Int(rank),
				PhaseType:       &release.DeployPhaseTypesValues.AgentBasedDeployment,
				WorkflowTasks:   tasks,
				DeploymentInput: deploymentInput,
			})
		}
	}
	return &phases
}

func expandReleaseTasks(input []interface{}) *[]release.WorkflowTask {
	tasks := make([]release.WorkflowTask, 0, len(input))
	for _, item := range input {
		task := item.(map[string]interface{})
		workflowTask := release.WorkflowTask{
			TaskId:                  converter.UUID(task["task_id"].(string)),
			Version:                 converter.String(task["version"].(string)),
			Name:                    converter.String(task["name"].(string)),
			Enabled:                 converter.Bool(task["enabled"].(bool)),
			Condition:               converter.String(task["condition"].(string)),
			AlwaysRun:               converter.Bool(task["always_run"].(bool)),
			ContinueOnError:         converter.Bool(task["continue_on_error"].(bool)),
			TimeoutInMinutes:        converter.Int(task["timeout_in_minutes"].(int)),
			RetryCountOnTaskFailure: converter.Int(task["retry_count_on_task_failure"].(int)),
			Inputs:                  expandStringMap(task["inputs"].(map[string]interface{})),
			OverrideInputs:          expandStringMap(task["override_inputs"].(map[string]interface{})),
			Environment:             expandStringMap(task["environment"].(map[string]interface{})),
		}
		if refName := task["ref_name"].(string); refName != "" {
			workflowTask.RefName = converter.String(refName)
		}
		if definitionType := task["definition_type"].(string); definitionType != "" {
			workflowTask.DefinitionType = converter.String(definitionType)
		}
		tasks = append(tasks, workflowTask)
	}
	return &tasks
}

func expandReleaseGates(input []interface{}) *release.ReleaseDefinitionGatesStep {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	block := input[0].(map[string]interface{})

	gates := make([]release.ReleaseDefinitionGate, 0)
	for _, item := range block["gate"].([]interface{}) {
		gate := item.(map[string]interface{})
		gates = append(gates, release.ReleaseDefinitionGate{
			Tasks: expandReleaseTasks(gate["task"].([]interface{})),
		})
	}

	return &release.ReleaseDefinitionGatesStep{
		Gates: &gates,
		GatesOptions: &release.ReleaseDefinitionGatesOptions{
			IsEnabled:              converter.Bool(block["is_enabled"].(bool)),
			SamplingInterval:       converter.Int(block["sampling_interval"].(int)),
			StabilizationTime:      converter.Int(block["stabilization_time"].(int)),
			MinimumSuccessDuration: converter.Int(block["minimum_success_duration"].(int)),
			Timeout:                converter.Int(block["timeout"].(int)),
		},
	}
}

func expandStringMap(input map[string]interface{}) *map[string]string {
	result := make(map[string]string, len(input))
	for k, v := range input {
		result[k] = v.(string)
	}
	return &result
}

func flattenReleaseStages(d *schema.ResourceData, input *[]release.ReleaseDefinitionEnvironment) []interface{} {
	if input == nil {
		return nil
	}
	priorStages := d.Get("stage").([]interface{})
	stages := make([]interface{}, 0, len(*input))
	for _, environment := range *input {
		name := converter.ToString(environment.Name, "")
		stages = append(stages, map[string]interface{}{
			"name":                  name,
			"rank":                  converter.ToInt(environment.Rank, 0),
			"condition":             flattenReleaseConditions(environment.Conditions),
			rdVariable:              flattenReleaseVariables(environment.Variables, priorStageVariables(priorStages, name)),
			"variable_groups":       flattenReleaseVariableGroups(environment.VariableGroups),
			"retention_policy":      flattenReleaseRetentionPolicy(environment.RetentionPolicy),
			"pre_deploy_approval":   flattenReleaseApprovals(environment.PreDeployApprovals),
			"post_deploy_approval":  flattenReleaseApprovals(environment.PostDeployApprovals),
			"pre_deployment_gates":  flattenReleaseGates(environment.PreDeploymentGates),
			"post_deployment_gates": flattenReleaseGates(environment.PostDeploymentGates),
			"environment_options":   flattenEnvironmentOptions(environment.EnvironmentOptions),
			"execution_policy":      flattenExecutionPolicy(environment.ExecutionPolicy),
			"environment_trigger":   flattenEnvironmentTriggers(environment.EnvironmentTriggers),
			"deploy_phase":          flattenReleaseDeployPhases(environment.DeployPhases),
		})
	}
	return stages
}

func priorStageVariables(priorStages []interface{}, name string) []interface{} {
	for _, item := range priorStages {
		stage, ok := item.(map[string]interface{})
		if !ok || stage["name"] != name {
			continue
		}
		if set, ok := stage[rdVariable].(*schema.Set); ok {
			return set.List()
		}
	}
	return nil
}

func flattenReleaseConditions(input *[]release.Condition) []interface{} {
	if input == nil {
		return nil
	}
	conditions := make([]interface{}, 0, len(*input))
	for _, condition := range *input {
		conditions = append(conditions, map[string]interface{}{
			"condition_type": converter.ToString((*string)(condition.ConditionType), ""),
			"name":           converter.ToString(condition.Name, ""),
			"value":          converter.ToString(condition.Value, ""),
		})
	}
	return conditions
}

func flattenReleaseRetentionPolicy(input *release.EnvironmentRetentionPolicy) []interface{} {
	if input == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"days_to_keep":     converter.ToInt(input.DaysToKeep, 0),
		"releases_to_keep": converter.ToInt(input.ReleasesToKeep, 0),
		"retain_build":     converter.ToBool(input.RetainBuild, false),
	}}
}

func flattenEnvironmentOptions(input *release.EnvironmentOptions) []interface{} {
	if input == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"auto_link_work_items":            converter.ToBool(input.AutoLinkWorkItems, false),
		"badge_enabled":                   converter.ToBool(input.BadgeEnabled, false),
		"publish_deployment_status":       converter.ToBool(input.PublishDeploymentStatus, false),
		"pull_request_deployment_enabled": converter.ToBool(input.PullRequestDeploymentEnabled, false),
	}}
}

func flattenExecutionPolicy(input *release.EnvironmentExecutionPolicy) []interface{} {
	if input == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"concurrency_count": converter.ToInt(input.ConcurrencyCount, 0),
		"queue_depth_count": converter.ToInt(input.QueueDepthCount, 0),
	}}
}

func flattenEnvironmentTriggers(input *[]release.EnvironmentTrigger) []interface{} {
	if input == nil || len(*input) == 0 {
		return nil
	}
	triggers := make([]interface{}, 0, len(*input))
	for _, trigger := range *input {
		item := map[string]interface{}{
			"trigger_type": converter.ToString((*string)(trigger.TriggerType), ""),
			"action":       "",
			"event_types":  []interface{}{},
		}
		if trigger.TriggerContent != nil {
			var content release.EnvironmentTriggerContent
			if err := json.Unmarshal([]byte(*trigger.TriggerContent), &content); err == nil {
				item["action"] = converter.ToString(content.Action, "")
				if content.EventTypes != nil {
					events := make([]interface{}, 0, len(*content.EventTypes))
					for _, e := range *content.EventTypes {
						events = append(events, e)
					}
					item["event_types"] = events
				}
			}
		}
		triggers = append(triggers, item)
	}
	return triggers
}

func flattenReleaseGates(input *release.ReleaseDefinitionGatesStep) []interface{} {
	if input == nil || input.Gates == nil || len(*input.Gates) == 0 {
		return nil
	}
	gates := make([]interface{}, 0, len(*input.Gates))
	for _, gate := range *input.Gates {
		gates = append(gates, map[string]interface{}{
			"task": flattenWorkflowTasks(gate.Tasks),
		})
	}
	block := map[string]interface{}{
		"gate": gates,
	}
	if opts := input.GatesOptions; opts != nil {
		block["is_enabled"] = converter.ToBool(opts.IsEnabled, false)
		block["sampling_interval"] = converter.ToInt(opts.SamplingInterval, 0)
		block["stabilization_time"] = converter.ToInt(opts.StabilizationTime, 0)
		block["minimum_success_duration"] = converter.ToInt(opts.MinimumSuccessDuration, 0)
		block["timeout"] = converter.ToInt(opts.Timeout, 0)
	}
	return []interface{}{block}
}

func flattenWorkflowTasks(input *[]release.WorkflowTask) []interface{} {
	if input == nil {
		return nil
	}
	tasks := make([]interface{}, 0, len(*input))
	for _, task := range *input {
		result := map[string]interface{}{
			"version":                     converter.ToString(task.Version, ""),
			"name":                        converter.ToString(task.Name, ""),
			"enabled":                     converter.ToBool(task.Enabled, false),
			"condition":                   converter.ToString(task.Condition, ""),
			"always_run":                  converter.ToBool(task.AlwaysRun, false),
			"continue_on_error":           converter.ToBool(task.ContinueOnError, false),
			"timeout_in_minutes":          converter.ToInt(task.TimeoutInMinutes, 0),
			"retry_count_on_task_failure": converter.ToInt(task.RetryCountOnTaskFailure, 0),
			"ref_name":                    converter.ToString(task.RefName, ""),
			"definition_type":             converter.ToString(task.DefinitionType, ""),
		}
		if task.TaskId != nil {
			result["task_id"] = task.TaskId.String()
		}
		if task.Inputs != nil {
			result["inputs"] = *task.Inputs
		}
		if task.OverrideInputs != nil {
			result["override_inputs"] = *task.OverrideInputs
		}
		if task.Environment != nil {
			result["environment"] = *task.Environment
		}
		tasks = append(tasks, result)
	}
	return tasks
}

func flattenReleaseApprovals(input *release.ReleaseDefinitionApprovals) []interface{} {
	if input == nil || input.Approvals == nil {
		return nil
	}
	approvals := make([]interface{}, 0, len(*input.Approvals))
	for _, step := range *input.Approvals {
		approval := map[string]interface{}{
			"is_automated":       converter.ToBool(step.IsAutomated, false),
			"is_notification_on": converter.ToBool(step.IsNotificationOn, false),
		}
		if step.Approver != nil {
			approval["approver_id"] = converter.ToString(step.Approver.Id, "")
		}
		approvals = append(approvals, approval)
	}
	return []interface{}{map[string]interface{}{
		"approval":         approvals,
		"approval_options": flattenApprovalOptions(input.ApprovalOptions),
	}}
}

func flattenApprovalOptions(input *release.ApprovalOptions) []interface{} {
	if input == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"required_approver_count":                                         converter.ToInt(input.RequiredApproverCount, 0),
		"release_creator_can_be_approver":                                 converter.ToBool(input.ReleaseCreatorCanBeApprover, false),
		"auto_triggered_and_previous_environment_approved_can_be_skipped": converter.ToBool(input.AutoTriggeredAndPreviousEnvironmentApprovedCanBeSkipped, false),
		"enforce_identity_revalidation":                                   converter.ToBool(input.EnforceIdentityRevalidation, false),
		"timeout_in_minutes":                                              converter.ToInt(input.TimeoutInMinutes, 0),
		"execution_order":                                                 converter.ToString((*string)(input.ExecutionOrder), string(release.ApprovalExecutionOrderValues.BeforeGates)),
	}}
}

func flattenReleaseDeployPhases(input *[]interface{}) []interface{} {
	if input == nil {
		return nil
	}
	phases := make([]interface{}, 0, len(*input))
	for _, item := range *input {
		phase, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		phaseType := stringFromMap(phase, "phaseType")
		result := map[string]interface{}{
			"name":       stringFromMap(phase, "name"),
			"rank":       intFromMap(phase, "rank"),
			"phase_type": phaseType,
		}
		if di, ok := phase["deploymentInput"].(map[string]interface{}); ok {
			result["timeout_in_minutes"] = intFromMap(di, "timeoutInMinutes")
			result["job_cancel_timeout_in_minutes"] = intFromMap(di, "jobCancelTimeoutInMinutes")
			result["condition"] = stringFromMap(di, "condition")
			result["skip_artifacts_download"] = boolFromMap(di, "skipArtifactsDownload")
			result["enable_access_token"] = boolFromMap(di, "enableAccessToken")

			switch release.DeployPhaseTypes(phaseType) {
			case release.DeployPhaseTypesValues.MachineGroupBasedDeployment:
				result["deployment_group_id"] = intFromMap(di, "queueId")
				result["deployment_health_option"] = stringFromMap(di, "deploymentHealthOption")
				result["health_percent"] = intFromMap(di, "healthPercent")
				if tags, ok := di["tags"].([]interface{}); ok {
					result["tags"] = tags
				}
			case release.DeployPhaseTypesValues.RunOnServer:
			default:
				result["agent_queue_id"] = intFromMap(di, "queueId")
				if spec, ok := di["agentSpecification"].(map[string]interface{}); ok {
					result["agent_specification"] = stringFromMap(spec, "identifier")
				}
			}
		}
		result["task"] = flattenReleaseTasks(phase["workflowTasks"])
		phases = append(phases, result)
	}
	return phases
}

func flattenReleaseTasks(input interface{}) []interface{} {
	rawTasks, ok := input.([]interface{})
	if !ok {
		return nil
	}
	tasks := make([]interface{}, 0, len(rawTasks))
	for _, item := range rawTasks {
		task, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		tasks = append(tasks, map[string]interface{}{
			"task_id":                     stringFromMap(task, "taskId"),
			"version":                     stringFromMap(task, "version"),
			"name":                        stringFromMap(task, "name"),
			"enabled":                     boolFromMap(task, "enabled"),
			"condition":                   stringFromMap(task, "condition"),
			"always_run":                  boolFromMap(task, "alwaysRun"),
			"continue_on_error":           boolFromMap(task, "continueOnError"),
			"timeout_in_minutes":          intFromMap(task, "timeoutInMinutes"),
			"retry_count_on_task_failure": intFromMap(task, "retryCountOnTaskFailure"),
			"ref_name":                    stringFromMap(task, "refName"),
			"definition_type":             stringFromMap(task, "definitionType"),
			"inputs":                      task["inputs"],
			"override_inputs":             task["overrideInputs"],
			"environment":                 task["environment"],
		})
	}
	return tasks
}

func stringFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func intFromMap(m map[string]interface{}, key string) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	if v, ok := m[key].(int); ok {
		return v
	}
	return 0
}

func boolFromMap(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}
