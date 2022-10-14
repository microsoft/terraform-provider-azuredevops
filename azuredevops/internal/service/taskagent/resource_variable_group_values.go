package taskagent

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	v5taskagent "github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceVariableGroupVariables() *schema.Resource {
	return &schema.Resource{
		Create:   resourceVariableGroupVariableCreateUpdate,
		Read:     resourceVariableGroupVariableRead,
		Update:   resourceVariableGroupVariableCreateUpdate,
		Delete:   resourceVariableGroupVariableDelete,
		Importer: tfhelper.ImportProjectQualifiedResource(),
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"variable_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},

			"variable": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:          schema.TypeString,
							Optional:      true,
							Default:       "",
							ConflictsWith: []string{"key_vault"},
						},
						"secret_value": {
							Type:          schema.TypeString,
							Optional:      true,
							Sensitive:     true,
							Default:       "",
							ConflictsWith: []string{"key_vault"},
						},
						"is_secret": {
							Type:          schema.TypeBool,
							Optional:      true,
							Default:       false,
							ConflictsWith: []string{"key_vault"},
						},
						"content_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"expires": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"key_vault": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"service_endpoint_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},
		},
	}
}

func resourceVariableGroupVariableCreateUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectId := d.Get("project_id").(string)
	groupIdRaw := d.Get("variable_group_id").(string)
	groupId, err := strconv.Atoi(groupIdRaw)
	if err != nil {
		return fmt.Errorf(" parsing variable group ID. Error: %+v ", err)
	}

	variables, err := expandVariables(clients, d)
	if err != nil {
		return fmt.Errorf(" expending variable group variables. Error: %+v ", err)
	}

	variableGroups, err := clients.V5TaskAgentClient.GetVariableGroupsById(clients.Ctx,
		v5taskagent.GetVariableGroupsByIdArgs{
			Project:  &projectId,
			GroupIds: &[]int{groupId},
		})
	if err != nil {
		return fmt.Errorf(" Get variable group %d error: %+v", groupId, err)
	}

	if variableGroups == nil && len(*variableGroups) == 0 {
		return fmt.Errorf(" Variable group %d not found. ", groupId)
	}

	variables.Name = (*variableGroups)[0].Name
	variables.Description = (*variableGroups)[0].Description

	variableGroup, err := clients.V5TaskAgentClient.UpdateVariableGroup(
		clients.Ctx,
		v5taskagent.UpdateVariableGroupArgs{
			Group:   variables,
			Project: &projectId,
			GroupId: &groupId,
		})
	if err != nil {
		return fmt.Errorf(" update variable group variables. Error: %+v", err)
	}
	d.SetId(fmt.Sprintf("%d", *variableGroup.Id))

	return resourceVariableGroupVariableRead(d, m)
}

func resourceVariableGroupVariableRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	groupIdRaw := d.Get("variable_group_id").(string)
	groupId, err := strconv.Atoi(groupIdRaw)
	if err != nil {
		return fmt.Errorf(" parsing variable group ID. Error: %+v ", err)
	}

	variableGroup, err := clients.V5TaskAgentClient.GetVariableGroup(
		clients.Ctx,
		v5taskagent.GetVariableGroupArgs{
			GroupId: &groupId,
			Project: &projectID,
		},
	)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error looking up variable group given ID (%v) and project ID (%v): %v ", groupId, projectID, err)
	}
	if variableGroup.Id == nil {
		d.SetId("")
		return nil
	}
	d.Set("project_id", projectID)
	err = flattenVGVariables(d, variableGroup, &projectID)

	if err != nil {
		return fmt.Errorf("Error flattening variables: %+v ", err)
	}

	d.SetId(fmt.Sprintf("%d", *variableGroup.Id))
	return nil
}

func resourceVariableGroupVariableDelete(d *schema.ResourceData, m interface{}) error {

	//TODO remove all variables

	return nil
}

func expandVariables(clients *client.AggregatedClient, d *schema.ResourceData) (*v5taskagent.VariableGroupParameters, error) {
	variables := d.Get("variable").(*schema.Set).List()

	variableMap := make(map[string]interface{})
	for _, variable := range variables {
		asMap := variable.(map[string]interface{})

		isSecret := converter.Bool(asMap["is_secret"].(bool))
		if *isSecret {
			variableMap[asMap["name"].(string)] = v5taskagent.VariableValue{
				Value:    converter.String(asMap["secret_value"].(string)),
				IsSecret: isSecret,
			}
		} else {
			variableMap[asMap["name"].(string)] = v5taskagent.VariableValue{
				Value:    converter.String(asMap["value"].(string)),
				IsSecret: isSecret,
			}
		}
	}

	variableGroup := &v5taskagent.VariableGroupParameters{
		Variables: &variableMap,
	}

	//parse variables from KV
	keyVault := d.Get("key_vault").(*schema.Set).List()
	if len(keyVault) == 1 {
		kvConfigures := keyVault[0].(map[string]interface{})
		kvName := kvConfigures["name"].(string)
		serviceEndpointID := kvConfigures["service_endpoint_id"].(string)

		serviceEndpointUUID, err := uuid.Parse(serviceEndpointID)
		if err != nil {
			return nil, err
		}

		variableGroup.ProviderData = v5taskagent.AzureKeyVaultVariableGroupProviderData{
			ServiceEndpointId: &serviceEndpointUUID,
			Vault:             &kvName,
		}

		projectID := converter.String(d.Get("project_id").(string))
		variableGroup.Type = converter.String("AzureKeyVault")
		kvVariables, invalidVariables, err := searchAzureKVSecrets(clients, *projectID, kvName, serviceEndpointID, variables)
		if err != nil {
			return nil, err
		}

		if len(invalidVariables) > 0 {
			return nil, fmt.Errorf("Invalid Key Vault secret: ( %s ) , can not find in Azure Key Vault: ( %s ) ",
				strings.Join(invalidVariables, ","),
				kvName)
		} else {
			variableGroup.Variables = &kvVariables
		}
	}
	return variableGroup, nil
}

func flattenVGVariables(d *schema.ResourceData, variableGroup *v5taskagent.VariableGroup, projectID *string) error {
	variables := make([]map[string]interface{}, len(*variableGroup.Variables))

	index := 0
	for varName, varVal := range *variableGroup.Variables {
		variableAsJSON, err := json.Marshal(varVal)
		if err != nil {
			return fmt.Errorf("Unable to marshal variable into JSON. Error: %+v ", err)
		}

		if isKeyVaultVariableGroupType(variableGroup.Type) {
			variables[index], err = flattenKeyVaultVariable(variableAsJSON, varName)
		} else {
			variables[index], err = flattenVGVariable(d, variableAsJSON, varName)
		}

		if err != nil {
			return err
		}

		index = index + 1
	}

	if err := d.Set("variable", variables); err != nil {
		return err
	}

	if isKeyVaultVariableGroupType(variableGroup.Type) {
		keyVault, err := flattenKeyVault(d, variableGroup)

		if err != nil {
			return err
		}

		if err = d.Set("key_vault", keyVault); err != nil {
			return err
		}
	}
	return nil
}

func flattenVGVariable(d *schema.ResourceData, variableAsJSON []byte, varName string) (map[string]interface{}, error) {
	var variable v5taskagent.AzureKeyVaultVariableValue
	err := json.Unmarshal(variableAsJSON, &variable)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal variable (%+v): %+v", variable, err)
	}

	isSecret := converter.ToBool(variable.IsSecret, false)
	var val = map[string]interface{}{
		"name":      varName,
		"value":     converter.ToString(variable.Value, ""),
		"is_secret": isSecret,
	}

	//read secret variables from state if exist
	if isSecret {
		if stateVal := tfhelper.FindMapInListWithGivenKeyValue(d, "variable", "name", varName); stateVal != nil {
			val = stateVal
		}
	}
	return val, nil
}
