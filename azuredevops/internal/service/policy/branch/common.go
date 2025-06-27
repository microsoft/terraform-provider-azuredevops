package branch

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/policy"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

/**
 * This file contains base functionality that can be leveraged by all policy configuration
 * resources. This is possible because a single API is used for configuring many different
 * policy types.
 */

// Policy type IDs. These are global and can be listed using the following endpoint:
//
//	https://docs.microsoft.com/en-us/rest/api/azure/devops/policy/types/list?view=azure-devops-rest-5.1
var (
	MinReviewerCount  = uuid.MustParse("fa4e907d-c16b-4a4c-9dfa-4906e5d171dd")
	BuildValidation   = uuid.MustParse("0609b952-1397-4640-95ec-e00a01b2c241")
	AutoReviewers     = uuid.MustParse("fd2167ab-b0be-447a-8ec8-39368250530e")
	WorkItemLinking   = uuid.MustParse("40e92b44-2fe1-4dd6-b3d8-74a9c21d0c6e")
	CommentResolution = uuid.MustParse("c6a1889d-b943-4856-b76f-9e46bb6b0df2")
	MergeTypes        = uuid.MustParse("fa4e907d-c16b-4a4c-9dfa-4916e5d171ab")
	StatusCheck       = uuid.MustParse("cbdc66da-9728-4af8-aada-9a5a32e4a226")
)

// policyCrudArgs arguments for genBasePolicyResource
type policyCrudArgs struct {
	FlattenFunc func(d *schema.ResourceData, policy *policy.PolicyConfiguration, projectID *string) error
	ExpandFunc  func(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error)
	PolicyType  uuid.UUID
}

// genBasePolicyResource creates a Resource with the common elements of a build policy
func genBasePolicyResource(crudArgs *policyCrudArgs) *schema.Resource {
	return &schema.Resource{
		Create:   genPolicyCreateFunc(crudArgs),
		Read:     genPolicyReadFunc(crudArgs),
		Update:   genPolicyUpdateFunc(crudArgs),
		Delete:   genPolicyDeleteFunc(crudArgs),
		Importer: tfhelper.ImportProjectQualifiedResourceInteger(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"blocking": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"settings": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scope": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"repository_id": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"repository_ref": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"match_type": {
										Type:             schema.TypeString,
										Optional:         true,
										Default:          "Exact",
										DiffSuppressFunc: suppress.CaseDifference,
										ValidateFunc: validation.StringInSlice([]string{
											"Exact", "Prefix", "DefaultBranch",
										}, true),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

type commonPolicySettings struct {
	Scopes []struct {
		RepositoryID      string `json:"repositoryId,omitempty"`
		RepositoryRefName string `json:"refName,omitempty"`
		MatchType         string `json:"matchKind,omitempty"`
	} `json:"scope"`
}

// baseFlattenFunc flattens each of the base elements of the schema
func baseFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	d.Set("project_id", converter.ToString(projectID, ""))
	d.Set("enabled", converter.ToBool(policyConfig.IsEnabled, true))
	d.Set("blocking", converter.ToBool(policyConfig.IsBlocking, true))
	settings, err := flattenSettings(policyConfig)
	if err != nil {
		return err
	}
	err = d.Set("settings", settings)
	if err != nil {
		return fmt.Errorf("Unable to persist policy settings configuration: %+v", err)
	}
	return nil
}

func flattenSettings(policyConfig *policy.PolicyConfiguration) ([]interface{}, error) {
	policySettings := commonPolicySettings{}
	policyAsJSON, err := json.Marshal(policyConfig.Settings)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal policy settings into JSON: %+v", err)
	}

	err = json.Unmarshal(policyAsJSON, &policySettings)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal policy settings. Error: %+v", err)
	}
	scopes := make([]interface{}, len(policySettings.Scopes))
	for index, scope := range policySettings.Scopes {
		scopeSetting := map[string]interface{}{}
		if scope.RepositoryID != "" {
			scopeSetting["repository_id"] = scope.RepositoryID
		}
		if scope.RepositoryRefName != "" {
			scopeSetting["repository_ref"] = scope.RepositoryRefName
		}
		if scope.MatchType != "" {
			scopeSetting["match_type"] = scope.MatchType
		}
		scopes[index] = scopeSetting
	}
	settings := []interface{}{
		map[string]interface{}{
			"scope": scopes,
		},
	}
	return settings, nil
}

// baseExpandFunc expands each of the base elements of the schema
func baseExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	projectID := d.Get("project_id").(string)
	policySettings, err := expandSettings(d)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing policy configuration settings: (%+v)", err)
	}
	policyConfig := policy.PolicyConfiguration{
		IsEnabled:  converter.Bool(d.Get("enabled").(bool)),
		IsBlocking: converter.Bool(d.Get("blocking").(bool)),
		Type: &policy.PolicyTypeRef{
			Id: &typeID,
		},
		Settings: policySettings,
	}

	if d.Id() != "" {
		policyID, err := strconv.Atoi(d.Id())
		if err != nil {
			return nil, nil, fmt.Errorf("parsing policy configuration ID: (%+v)", err)
		}
		policyConfig.Id = &policyID
	}

	return &policyConfig, &projectID, nil
}

func expandSettings(d *schema.ResourceData) (map[string]interface{}, error) {
	settingsList := d.Get("settings").([]interface{})
	settings := settingsList[0].(map[string]interface{})
	settingsScopes := settings["scope"].([]interface{})

	scopes := make([]map[string]interface{}, len(settingsScopes))
	for index, scope := range settingsScopes {
		scopeMap := scope.(map[string]interface{})

		scopeSetting := map[string]interface{}{}
		if repoID, ok := scopeMap["repository_id"]; ok {
			if repoID == "" {
				scopeSetting["repositoryId"] = nil
			} else {
				scopeSetting["repositoryId"] = repoID
			}
		}
		if repoRef, ok := scopeMap["repository_ref"]; ok {
			if repoRef == "" {
				scopeSetting["refName"] = nil
			} else {
				scopeSetting["refName"] = repoRef
			}
		}
		if matchType, ok := scopeMap["match_type"]; ok {
			if matchType == "" {
				scopeSetting["matchKind"] = nil
			} else {
				scopeSetting["matchKind"] = matchType
			}
		}
		if strings.EqualFold(scopeSetting["matchKind"].(string), "DefaultBranch") && (scopeSetting["repositoryId"] != nil || scopeSetting["refName"] != nil) {
			return nil, fmt.Errorf("neither 'repository_id' nor 'repository_ref' can be set when 'match_type=DefaultBranch'")
		}
		scopes[index] = scopeSetting
	}
	return map[string]interface{}{
		"scope": scopes,
	}, nil
}

//lint:ignore SA1019 SDKv2 migration  - staticcheck's own linter directives are currently being ignored under golanci-lint
func genPolicyCreateFunc(crudArgs *policyCrudArgs) schema.CreateFunc { //nolint:staticcheck
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*client.AggregatedClient)
		policyConfig, projectID, err := crudArgs.ExpandFunc(d, crudArgs.PolicyType)
		if err != nil {
			return err
		}

		createdPolicy, err := clients.PolicyClient.CreatePolicyConfiguration(clients.Ctx, policy.CreatePolicyConfigurationArgs{
			Configuration: policyConfig,
			Project:       projectID,
		})
		if err != nil {
			return fmt.Errorf("creating policy in Azure DevOps: %+v", err)
		}

		d.SetId(strconv.Itoa(*createdPolicy.Id))
		return genPolicyReadFunc(crudArgs)(d, m)
	}
}

//lint:ignore SA1019 SDKv2 migration  - staticcheck's own linter directives are currently being ignored under golanci-lint
func genPolicyReadFunc(crudArgs *policyCrudArgs) schema.ReadFunc { //nolint:staticcheck
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*client.AggregatedClient)
		projectID := d.Get("project_id").(string)
		policyID, err := strconv.Atoi(d.Id())
		if err != nil {
			return fmt.Errorf("converting policy ID to an integer: (%+v)", err)
		}

		policyConfig, err := clients.PolicyClient.GetPolicyConfiguration(clients.Ctx, policy.GetPolicyConfigurationArgs{
			Project:         &projectID,
			ConfigurationId: &policyID,
		})

		if utils.ResponseWasNotFound(err) || (policyConfig != nil && *policyConfig.IsDeleted) {
			d.SetId("")
			return nil
		}

		if err != nil {
			return fmt.Errorf("looking up build policy configuration with ID (%v) and project ID (%v): %v", policyID, projectID, err)
		}

		return crudArgs.FlattenFunc(d, policyConfig, &projectID)
	}
}

//lint:ignore SA1019 SDKv2 migration  - staticcheck's own linter directives are currently being ignored under golanci-lint
func genPolicyUpdateFunc(crudArgs *policyCrudArgs) schema.UpdateFunc { //nolint:staticcheck
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*client.AggregatedClient)
		policyConfig, projectID, err := crudArgs.ExpandFunc(d, crudArgs.PolicyType)
		if err != nil {
			return err
		}

		_, err = clients.PolicyClient.UpdatePolicyConfiguration(clients.Ctx, policy.UpdatePolicyConfigurationArgs{
			ConfigurationId: policyConfig.Id,
			Configuration:   policyConfig,
			Project:         projectID,
		})
		if err != nil {
			return fmt.Errorf("updating policy in Azure DevOps: %+v", err)
		}

		return genPolicyReadFunc(crudArgs)(d, m)
	}
}

//lint:ignore SA1019 SDKv2 migration  - staticcheck's own linter directives are currently being ignored under golanci-lint
func genPolicyDeleteFunc(crudArgs *policyCrudArgs) schema.DeleteFunc { //nolint:staticcheck
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*client.AggregatedClient)
		policyConfig, projectID, err := crudArgs.ExpandFunc(d, crudArgs.PolicyType)
		if err != nil {
			return err
		}

		err = clients.PolicyClient.DeletePolicyConfiguration(clients.Ctx, policy.DeletePolicyConfigurationArgs{
			ConfigurationId: policyConfig.Id,
			Project:         projectID,
		})
		if err != nil {
			return fmt.Errorf("deleting policy in Azure DevOps: %+v", err)
		}

		return nil
	}
}
