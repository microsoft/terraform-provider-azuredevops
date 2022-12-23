package repository

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/policy"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// Policy type IDs. These are global and can be listed using the following endpoint:
//
//	https://docs.microsoft.com/en-us/rest/api/azure/devops/policy/types/list?view=azure-devops-rest-5.1
var (
	AuthorEmailPattern = uuid.MustParse("77ed4bd3-b063-4689-934a-175e4d0a78d7")
	FilePathPattern    = uuid.MustParse("51c78909-e838-41a2-9496-c647091e3c61")
	CaseEnforcement    = uuid.MustParse("7ed39669-655c-494e-b4a0-a08b4da0fcce")
	ReservedNames      = uuid.MustParse("db2b9b4c-180d-4529-9701-01541d19f36b")
	PathLength         = uuid.MustParse("001a79cf-fda1-4c4e-9e7c-bac40ee5ead8")
	FileSize           = uuid.MustParse("2e26e725-8201-4edd-8bf5-978563c34a80")
	CheckCredentials   = uuid.MustParse("e67ae10f-cf9a-40bc-8e66-6b3a8216956e")
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
			"repository_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},
		},
	}
}

type commonPolicySettings struct {
	Scopes []struct {
		RepositoryID string `json:"repositoryId,omitempty"`
	} `json:"scope"`
}

// baseFlattenFunc flattens each of the base elements of the schema
func baseFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	if policyConfig.Id == nil {
		d.SetId("")
		return nil
	}
	d.SetId(strconv.Itoa(*policyConfig.Id))
	d.Set("project_id", converter.ToString(projectID, ""))
	d.Set("enabled", converter.ToBool(policyConfig.IsEnabled, true))
	d.Set("blocking", converter.ToBool(policyConfig.IsBlocking, true))
	repoIds, err := flattenSettings(policyConfig)
	if err != nil {
		return err
	}
	err = d.Set("repository_ids", repoIds)
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

	_ = json.Unmarshal(policyAsJSON, &policySettings)
	var repoIds []interface{}
	for _, scope := range policySettings.Scopes {
		if scope.RepositoryID != "" {
			repoIds = append(repoIds, scope.RepositoryID)
		}
	}
	return repoIds, nil
}

// baseExpandFunc expands each of the base elements of the schema
func baseExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	projectID := d.Get("project_id").(string)

	policyConfig := policy.PolicyConfiguration{
		IsEnabled:  converter.Bool(d.Get("enabled").(bool)),
		IsBlocking: converter.Bool(d.Get("blocking").(bool)),
		Type: &policy.PolicyTypeRef{
			Id: &typeID,
		},
		Settings: expandSettings(d),
	}

	if d.Id() != "" {
		policyID, err := strconv.Atoi(d.Id())
		if err != nil {
			return nil, nil, fmt.Errorf("Error parsing policy configuration ID: (%+v)", err)
		}
		policyConfig.Id = &policyID
	}

	return &policyConfig, &projectID, nil
}

func expandSettings(d *schema.ResourceData) map[string]interface{} {
	repoIds := d.Get("repository_ids").([]interface{})
	if len(repoIds) == 0 {
		return map[string]interface{}{
			"scope": []map[string]interface{}{
				{
					"repositoryId": "",
				},
			},
		}
	}

	scopes := make([]map[string]interface{}, len(repoIds))
	for idx, repoID := range repoIds {
		scopes[idx] = map[string]interface{}{
			"repositoryId": repoID,
		}
	}
	return map[string]interface{}{
		"scope": scopes,
	}
}

//lint:ignore SA1019 schema.CreateFunc is deprecated: Please use the context aware equivalents instead
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
			return fmt.Errorf("Error creating policy in Azure DevOps: %+v", err)
		}

		return crudArgs.FlattenFunc(d, createdPolicy, projectID)
	}
}

//lint:ignore SA1019 SDKv2 migration  - staticcheck's own linter directives are currently being ignored under golanci-lint
func genPolicyReadFunc(crudArgs *policyCrudArgs) schema.ReadFunc { //nolint:staticcheck
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*client.AggregatedClient)
		projectID := d.Get("project_id").(string)
		policyID, err := strconv.Atoi(d.Id())

		if err != nil {
			return fmt.Errorf("Error converting policy ID to an integer: (%+v)", err)
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
			return fmt.Errorf("Error looking up build policy configuration with ID (%v) and project ID (%v): %v", policyID, projectID, err)
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

		updatedPolicy, err := clients.PolicyClient.UpdatePolicyConfiguration(clients.Ctx, policy.UpdatePolicyConfigurationArgs{
			ConfigurationId: policyConfig.Id,
			Configuration:   policyConfig,
			Project:         projectID,
		})

		if err != nil {
			return fmt.Errorf("Error updating policy in Azure DevOps: %+v", err)
		}

		return crudArgs.FlattenFunc(d, updatedPolicy, projectID)
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
			return fmt.Errorf("Error deleting policy in Azure DevOps: %+v", err)
		}

		return nil
	}
}
