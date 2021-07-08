package policy

import (
	"encoding/json"
	"fmt"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/policy/branch"
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/policy"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// Policy type IDs. These are global and can be listed using the following endpoint:
//	https://docs.microsoft.com/en-us/rest/api/azure/devops/policy/types/list?view=azure-devops-rest-5.1
var (
	AuthorEmailPattern = uuid.MustParse("77ed4bd3-b063-4689-934a-175e4d0a78d7")
	FilePathPattern    = uuid.MustParse("51c78909-e838-41a2-9496-c647091e3c61")
	CaseEnforcement    = uuid.MustParse("7ed39669-655c-494e-b4a0-a08b4da0fcce")
	ReservedNames      = uuid.MustParse("db2b9b4c-180d-4529-9701-01541d19f36b")
	PathLength         = uuid.MustParse("001a79cf-fda1-4c4e-9e7c-bac40ee5ead8")
	FileSize           = uuid.MustParse("2e26e725-8201-4edd-8bf5-978563c34a80")
)

// Keys for schema elements
const (
	SchemaProjectID    = "project_id"
	SchemaEnabled      = "enabled"
	SchemaBlocking     = "blocking"
	SchemaSettings     = "settings"
	SchemaScope        = "scope"
	SchemaRepositoryID = "repository_id"
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
			SchemaProjectID: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			SchemaEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			SchemaBlocking: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			SchemaSettings: {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						SchemaScope: {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									branch.SchemaRepositoryID: {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
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
	d.Set(SchemaProjectID, converter.ToString(projectID, ""))
	d.Set(SchemaEnabled, converter.ToBool(policyConfig.IsEnabled, true))
	d.Set(SchemaBlocking, converter.ToBool(policyConfig.IsBlocking, true))
	settings, err := flattenSettings(policyConfig)
	if err != nil {
		return err
	}
	err = d.Set(SchemaSettings, settings)
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
	var scopes []interface{}
	for _, scope := range policySettings.Scopes {
		scopeSetting := map[string]interface{}{}
		if scope.RepositoryID != "" {
			scopeSetting[SchemaRepositoryID] = scope.RepositoryID
			scopes = append(scopes, scopeSetting)
		}
	}

	return []interface{}{
		map[string]interface{}{
			SchemaScope: scopes,
		},
	}, nil
}

// baseExpandFunc expands each of the base elements of the schema
func baseExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	projectID := d.Get(SchemaProjectID).(string)

	policyConfig := policy.PolicyConfiguration{
		IsEnabled:  converter.Bool(d.Get(SchemaEnabled).(bool)),
		IsBlocking: converter.Bool(d.Get(SchemaBlocking).(bool)),
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
	settingsList := d.Get(SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})
	settingsScopes := settings[SchemaScope].([]interface{})

	if len(settingsScopes) == 0 {
		return map[string]interface{}{
			SchemaScope: []map[string]interface{}{
				{
					"repositoryId": "",
				},
			},
		}
	} else {
		scopes := make([]map[string]interface{}, len(settingsScopes))
		for index, scope := range settingsScopes {
			scopeMap := scope.(map[string]interface{})
			scopeSetting := map[string]interface{}{}
			if repoID, ok := scopeMap[SchemaRepositoryID]; ok {
				scopeSetting["repositoryId"] = repoID
			}
			scopes[index] = scopeSetting
		}
		return map[string]interface{}{
			SchemaScope: scopes,
		}
	}

}

func genPolicyCreateFunc(crudArgs *policyCrudArgs) schema.CreateFunc {
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

func genPolicyReadFunc(crudArgs *policyCrudArgs) schema.ReadFunc {
	return func(d *schema.ResourceData, m interface{}) error {
		clients := m.(*client.AggregatedClient)
		projectID := d.Get(SchemaProjectID).(string)
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

func genPolicyUpdateFunc(crudArgs *policyCrudArgs) schema.UpdateFunc {
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

func genPolicyDeleteFunc(crudArgs *policyCrudArgs) schema.DeleteFunc {
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
