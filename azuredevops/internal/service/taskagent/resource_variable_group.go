package taskagent

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

const (
	vgProjectID         = "project_id"
	vgName              = "name"
	vgDescription       = "description"
	vgAllowAccess       = "allow_access"
	vgVariable          = "variable"
	vgValue             = "value"
	secretVgValue       = "secret_value"
	vgIsSecret          = "is_secret"
	vgKeyVault          = "key_vault"
	vgServiceEndpointID = "service_endpoint_id"
	vgContentType       = "content_type"
	vgEnabled           = "enabled"
	vgExpires           = "expires"
)

const (
	azureKeyVaultType                         = "AzureKeyVault"
	invalidVariableGroupIDErrorMessageFormat  = "Error parsing the variable group ID from the Terraform resource data: %v"
	flatteningVariableGroupErrorMessageFormat = "Error flattening variable group: %v"
	expandingVariableGroupErrorMessageFormat  = "Expanding variable group resource data: %+v"
)

type KeyVaultSecretAttributes struct {
	Enabled       *bool   `json:"enabled,omitempty"`
	Created       *int64  `json:"created,omitempty"`
	Updated       *int64  `json:"updated,omitempty"`
	Expire        *int64  `json:"exp,omitempty"`
	RecoveryLevel *string `json:"recoveryLevel,omitempty"`
}

type KeyVaultSecret struct {
	ContentType              *string `json:"contentType,omitempty"`
	ID                       *string `json:"id,omitempty"`
	KeyVaultSecretAttributes `json:"attributes,omitempty"`
}

type KeyVaultSecretResult struct {
	Value    *[]KeyVaultSecret `json:"value,omitempty"`
	NextLink *string           `json:"nextLink,omitempty"`
}

// ResourceVariableGroup schema and implementation for variable group resource
func ResourceVariableGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceVariableGroupCreate,
		Read:   resourceVariableGroupRead,
		Update: resourceVariableGroupUpdate,
		Delete: resourceVariableGroupDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResource(),
		Schema: map[string]*schema.Schema{
			vgProjectID: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			vgName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			vgDescription: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			vgAllowAccess: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			vgVariable: {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						vgName: {
							Type:     schema.TypeString,
							Required: true,
						},
						vgValue: {
							Type:          schema.TypeString,
							Optional:      true,
							Default:       "",
							ConflictsWith: []string{vgKeyVault},
						},
						secretVgValue: {
							Type:          schema.TypeString,
							Optional:      true,
							Sensitive:     true,
							Default:       "",
							ConflictsWith: []string{vgKeyVault},
						},
						vgIsSecret: {
							Type:          schema.TypeBool,
							Optional:      true,
							Default:       false,
							ConflictsWith: []string{vgKeyVault},
						},
						vgContentType: {
							Type:     schema.TypeString,
							Computed: true,
						},
						vgEnabled: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						vgExpires: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			vgKeyVault: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						vgName: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						vgServiceEndpointID: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"search_depth": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  20,
						},
					},
				},
			},
		},
	}
}

func resourceVariableGroupCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	variableGroupParameters, projectID, err := expandVariableGroupParameters(clients, d)
	if err != nil {
		return fmt.Errorf(expandingVariableGroupErrorMessageFormat, err)
	}

	addedVariableGroup, err := createVariableGroup(clients, variableGroupParameters, projectID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf(" creating variable group in Azure DevOps: %+v", err)
	}

	err = flattenVariableGroup(d, addedVariableGroup, projectID)

	if err != nil {
		return fmt.Errorf(flatteningVariableGroupErrorMessageFormat, err)
	}

	// Update Allow Access with definition Reference
	definitionResourceReferenceArgs := expandAllowAccess(d, addedVariableGroup)
	definitionResourceReference, err := updateDefinitionResourceAuth(clients, definitionResourceReferenceArgs, projectID)
	if err != nil {
		return fmt.Errorf("Error creating definitionResourceReference Azure DevOps object: %+v", err)
	}

	flattenAllowAccess(d, definitionResourceReference)

	return resourceVariableGroupRead(d, m)
}

func resourceVariableGroupRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID, variableGroupID, err := tfhelper.ParseProjectIDAndResourceID(d)
	if err != nil {
		return fmt.Errorf(invalidVariableGroupIDErrorMessageFormat, err)
	}

	variableGroup, err := clients.TaskAgentClient.GetVariableGroup(
		clients.Ctx,
		taskagent.GetVariableGroupArgs{
			GroupId: &variableGroupID,
			Project: &projectID,
		},
	)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error looking up variable group given ID (%v) and project ID (%v): %v", variableGroupID, projectID, err)
	}
	if variableGroup.Id == nil {
		d.SetId("")
		return nil
	}

	err = flattenVariableGroup(d, variableGroup, &projectID)

	if err != nil {
		return fmt.Errorf(flatteningVariableGroupErrorMessageFormat, err)
	}

	//Read the Authorization Resource for get allow access property
	resourceRefType := "variablegroup"
	varGroupID := strconv.Itoa(variableGroupID)

	projectResources, err := clients.BuildClient.GetProjectResources(
		clients.Ctx,
		build.GetProjectResourcesArgs{
			Project: &projectID,
			Type:    &resourceRefType,
			Id:      &varGroupID,
		},
	)

	if err != nil {
		return fmt.Errorf("Error looking up project resources given ID (%v) and project ID (%v): %v", variableGroupID, projectID, err)
	}

	flattenAllowAccess(d, projectResources)
	return nil
}

func resourceVariableGroupUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	variableGroupParams, projectID, err := expandVariableGroupParameters(clients, d)
	if err != nil {
		return fmt.Errorf(expandingVariableGroupErrorMessageFormat, err)
	}

	_, variableGroupID, err := tfhelper.ParseProjectIDAndResourceID(d)
	if err != nil {
		return fmt.Errorf(invalidVariableGroupIDErrorMessageFormat, err)
	}

	updatedVariableGroup, err := updateVariableGroup(clients, variableGroupParams, &variableGroupID, projectID)
	if err != nil {
		return fmt.Errorf("Error updating variable group in Azure DevOps: %+v", err)
	}

	err = flattenVariableGroup(d, updatedVariableGroup, projectID)

	if err != nil {
		return fmt.Errorf(flatteningVariableGroupErrorMessageFormat, err)
	}

	// Update Allow Access
	definitionResourceReferenceArgs := expandAllowAccess(d, updatedVariableGroup)
	definitionResourceReference, err := updateDefinitionResourceAuth(clients, definitionResourceReferenceArgs, projectID)
	if err != nil {
		return fmt.Errorf("Error updating definitionResourceReference Azure DevOps object: %+v", err)
	}

	flattenAllowAccess(d, definitionResourceReference)

	return resourceVariableGroupRead(d, m)
}

func resourceVariableGroupDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID, variableGroupID, err := tfhelper.ParseProjectIDAndResourceID(d)
	if err != nil {
		return fmt.Errorf(invalidVariableGroupIDErrorMessageFormat, err)
	}
	//delete the definition resource (allow access)
	varGroupID := strconv.Itoa(variableGroupID)
	_, err = deleteDefinitionResourceAuth(clients, &varGroupID, &projectID)
	if err != nil {
		return fmt.Errorf("Error deleting the allow access definitionResource for variable group ID (%v) and project ID (%v): %v", variableGroupID, projectID, err)
	}
	//delete the variable group
	return deleteVariableGroup(clients, &projectID, &variableGroupID)
}

// Make the Azure DevOps API call to create the variable group
func createVariableGroup(clients *client.AggregatedClient, variableGroupParams *taskagent.VariableGroupParameters, project *string, timeOut time.Duration) (*taskagent.VariableGroup, error) {
	createdVG, err := clients.TaskAgentClient.AddVariableGroup(
		clients.Ctx,
		taskagent.AddVariableGroupArgs{
			VariableGroupParameters: variableGroupParams,
		})
	if err != nil {
		return nil, err
	}

	stateConf := &retry.StateChangeConf{
		ContinuousTargetOccurence: 2,
		Delay:                     5 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"inProgress"},
		Target:                    []string{"succeeded", "failed"},
		Refresh: func() (result interface{}, state string, err error) {
			createdVG, err = clients.TaskAgentClient.GetVariableGroup(
				clients.Ctx,
				taskagent.GetVariableGroupArgs{
					GroupId: createdVG.Id,
					Project: project,
				},
			)
			if err != nil {
				return createdVG, "failed", fmt.Errorf(" fail to get Variable Group. %v ", err)
			}
			if createdVG != nil && createdVG.Variables != nil && (len(*variableGroupParams.Variables) == len(*createdVG.Variables)) {
				return createdVG, "succeeded", nil
			}
			return createdVG, "inProgress", nil
		},
		Timeout: timeOut,
	}

	if _, err = stateConf.WaitForStateContext(clients.Ctx); err != nil {
		return nil, fmt.Errorf(" waiting for Variable Group ready. %v ", err)
	}
	return createdVG, nil
}

// Make the Azure DevOps API call to update the variable group
func updateVariableGroup(clients *client.AggregatedClient, parameters *taskagent.VariableGroupParameters, variableGroupID *int, project *string) (*taskagent.VariableGroup, error) {
	updatedVariableGroup, err := clients.TaskAgentClient.UpdateVariableGroup(
		clients.Ctx,
		taskagent.UpdateVariableGroupArgs{
			GroupId:                 variableGroupID,
			VariableGroupParameters: parameters,
		})

	return updatedVariableGroup, err
}

// Make the Azure DevOps API call to delete the variable group
func deleteVariableGroup(clients *client.AggregatedClient, projectId *string, variableGroupID *int) error {
	err := clients.TaskAgentClient.DeleteVariableGroup(
		clients.Ctx,
		taskagent.DeleteVariableGroupArgs{
			ProjectIds: &[]string{*projectId},
			GroupId:    variableGroupID,
		})

	return err
}

// Convert internal Terraform data structure to an AzDO data structure
func expandVariableGroupParameters(clients *client.AggregatedClient, d *schema.ResourceData) (*taskagent.VariableGroupParameters, *string, error) {
	projectID := converter.String(d.Get(vgProjectID).(string))
	name := converter.String(d.Get(vgName).(string))
	description := converter.String(d.Get(vgDescription).(string))
	variables := d.Get(vgVariable).(*schema.Set).List()

	// needed to detect if the secret_value attribute is set in the config
	// see https://github.com/hashicorp/terraform-plugin-sdk/issues/741
	for it := d.GetRawConfig().AsValueMap()[vgVariable].ElementIterator(); it.Next(); {
		_, ctyVariable := it.Element()
		ctyVariableAsMap := ctyVariable.AsValueMap()
		name := ctyVariableAsMap[vgName].AsString()
		valueSet := !ctyVariableAsMap[vgValue].IsNull()
		secretValueSet := !ctyVariableAsMap[secretVgValue].IsNull()
		isSecretSet := !ctyVariableAsMap[vgIsSecret].IsNull()

		if valueSet && (secretValueSet || isSecretSet) || secretValueSet != isSecretSet {
			return nil, nil, fmt.Errorf("`%s` variable can have either only `value` attribute or both `is_secret` and `secret_value` attributes", name)
		}
	}

	variableMap := make(map[string]interface{})

	for _, variable := range variables {
		asMap := variable.(map[string]interface{})

		isSecret := converter.Bool(asMap[vgIsSecret].(bool))
		if *isSecret {
			variableMap[asMap[vgName].(string)] = taskagent.VariableValue{
				Value:    converter.String(asMap[secretVgValue].(string)),
				IsSecret: isSecret,
			}
		} else {
			variableMap[asMap[vgName].(string)] = taskagent.VariableValue{
				Value:    converter.String(asMap[vgValue].(string)),
				IsSecret: isSecret,
			}
		}
	}

	projectUUId, err := uuid.Parse(*projectID)
	if err != nil {
		return nil, nil, err
	}

	variableGroup := &taskagent.VariableGroupParameters{
		Name:        name,
		Description: description,
		Variables:   &variableMap,
		VariableGroupProjectReferences: &[]taskagent.VariableGroupProjectReference{
			{
				Description: description,
				Name:        name,
				ProjectReference: &taskagent.ProjectReference{
					Id: &projectUUId,
				},
			},
		},
	}

	keyVault := d.Get(vgKeyVault).([]interface{})

	// Note: this will be of length 1 based on the schema definition above.
	if len(keyVault) == 1 {
		kvConfigures := keyVault[0].(map[string]interface{})
		kvName := kvConfigures[vgName].(string)
		serviceEndpointID := kvConfigures[vgServiceEndpointID].(string)
		depth := kvConfigures["search_depth"].(int)

		serviceEndpointUUID, err := uuid.Parse(serviceEndpointID)
		if err != nil {
			return nil, nil, err
		}

		variableGroup.ProviderData = taskagent.AzureKeyVaultVariableGroupProviderData{
			ServiceEndpointId: &serviceEndpointUUID,
			Vault:             &kvName,
			LastRefreshedOn:   &azuredevops.Time{Time: time.Now()},
		}

		variableGroup.Type = converter.String(azureKeyVaultType)
		kvVariables, invalidVariables, err := searchAzureKVSecrets(clients, *projectID, kvName, serviceEndpointID, variables, depth)
		if err != nil {
			return nil, nil, err
		}

		if len(invalidVariables) > 0 {
			return nil, nil, fmt.Errorf("Invalid Key Vault secret: ( %s ) , can not find in Azure Key Vault or the secret has been disbled: ( %s ) ",
				strings.Join(invalidVariables, ","),
				kvName)
		} else {
			variableGroup.Variables = &kvVariables
		}
	}
	return variableGroup, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenVariableGroup(d *schema.ResourceData, variableGroup *taskagent.VariableGroup, projectID *string) error {
	d.SetId(fmt.Sprintf("%d", *variableGroup.Id))
	d.Set(vgName, *variableGroup.Name)
	d.Set(vgDescription, converter.ToString(variableGroup.Description, ""))
	d.Set(vgProjectID, projectID)

	variables, err := flattenVariables(d, variableGroup)

	if err != nil {
		return err
	}

	if err = d.Set(vgVariable, variables); err != nil {
		return err
	}

	if isKeyVaultVariableGroupType(variableGroup.Type) {
		keyVault, err := flattenKeyVault(d, variableGroup)

		if err != nil {
			return err
		}

		if err = d.Set(vgKeyVault, keyVault); err != nil {
			return err
		}
	}
	return nil
}

func isKeyVaultVariableGroupType(variableGrouptype *string) bool {
	return variableGrouptype != nil && *variableGrouptype == azureKeyVaultType
}

// Convert AzDO Variables data structure to Terraform TypeSet
//
// Note: The AzDO API does not return the value for variables marked as a secret. For this reason
//
//	variables marked as secret will need to be pulled from the state itself
func flattenVariables(d *schema.ResourceData, variableGroup *taskagent.VariableGroup) (interface{}, error) {
	variables := make([]map[string]interface{}, len(*variableGroup.Variables))

	index := 0
	for varName, varVal := range *variableGroup.Variables {
		variableAsJSON, err := json.Marshal(varVal)
		if err != nil {
			return nil, fmt.Errorf("Unable to marshal variable into JSON: %+v", err)
		}

		if isKeyVaultVariableGroupType(variableGroup.Type) {
			variables[index], err = flattenKeyVaultVariable(variableAsJSON, varName)
		} else {
			variables[index], err = flattenVariable(d, variableAsJSON, varName)
		}

		if err != nil {
			return nil, err
		}

		index = index + 1
	}

	return variables, nil
}

func flattenKeyVaultVariable(variableAsJSON []byte, varName string) (map[string]interface{}, error) {
	var variable taskagent.AzureKeyVaultVariableValue
	err := json.Unmarshal(variableAsJSON, &variable)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal variable (%+v): %+v", variable, err)
	}

	variableMap := map[string]interface{}{
		vgName:        varName,
		vgValue:       nil,
		secretVgValue: nil,
		vgIsSecret:    false,
		vgEnabled:     converter.ToBool(variable.Enabled, false),
		vgContentType: converter.ToString(variable.ContentType, ""),
	}
	if variable.Expires != nil {
		variableMap[vgExpires] = variable.Expires.String()
	}
	return variableMap, nil
}

func flattenVariable(d *schema.ResourceData, variableAsJSON []byte, varName string) (map[string]interface{}, error) {
	var variable taskagent.AzureKeyVaultVariableValue
	err := json.Unmarshal(variableAsJSON, &variable)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal variable (%+v): %+v", variable, err)
	}

	isSecret := converter.ToBool(variable.IsSecret, false)
	var val = map[string]interface{}{
		vgName:     varName,
		vgValue:    converter.ToString(variable.Value, ""),
		vgIsSecret: isSecret,
	}

	//read secret variables from state if exist
	if isSecret {
		if stateVal := tfhelper.FindMapInSetWithGivenKeyValue(d, vgVariable, vgName, varName); stateVal != nil {
			val = stateVal
		}
	}
	return val, nil
}

func flattenKeyVault(d *schema.ResourceData, variableGroup *taskagent.VariableGroup) (interface{}, error) {
	providerDataAsJSON, err := json.Marshal(variableGroup.ProviderData)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal provider data into JSON: %+v", err)
	}

	var providerData taskagent.AzureKeyVaultVariableGroupProviderData
	err = json.Unmarshal(providerDataAsJSON, &providerData)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal provider data (%+v): %+v", providerData, err)
	}

	keyVault := []map[string]interface{}{{
		vgName:              providerData.Vault,
		vgServiceEndpointID: providerData.ServiceEndpointId.String(),
	}}

	keyVaultRaw := d.Get("key_vault").([]interface{})
	if len(keyVaultRaw) == 1 {
		kvConfigures := keyVaultRaw[0].(map[string]interface{})
		keyVault[0]["search_depth"] = kvConfigures["search_depth"].(int)
	}

	return keyVault, nil
}

// Convert internal Terraform data structure to an AzDO data structure for Allow Access
func expandAllowAccess(d *schema.ResourceData, createdVariableGroup *taskagent.VariableGroup) []build.DefinitionResourceReference {
	resourceRefType := "variablegroup"
	variableGroupID := strconv.Itoa(*createdVariableGroup.Id)

	var arrayDefinitionResourceReference []build.DefinitionResourceReference

	defResourceRef := build.DefinitionResourceReference{
		Type:       &resourceRefType,
		Authorized: converter.Bool(d.Get(vgAllowAccess).(bool)),
		Name:       createdVariableGroup.Name,
		Id:         &variableGroupID,
	}

	arrayDefinitionResourceReference = append(arrayDefinitionResourceReference, defResourceRef)

	return arrayDefinitionResourceReference
}

// Make the Azure DevOps API call to update the Definition resource = Allow Access
func updateDefinitionResourceAuth(clients *client.AggregatedClient, definitionResource []build.DefinitionResourceReference, project *string) (*[]build.DefinitionResourceReference, error) {
	definitionResourceReference, err := clients.BuildClient.AuthorizeProjectResources(
		clients.Ctx, build.AuthorizeProjectResourcesArgs{
			Resources: &definitionResource,
			Project:   project,
		})

	return definitionResourceReference, err
}

// Make the Azure DevOps API call to delete the resource Auth Authorized=false
func deleteDefinitionResourceAuth(clients *client.AggregatedClient, variableGroupID *string, project *string) (*[]build.DefinitionResourceReference, error) {
	resourceRefType := "variablegroup"
	auth := false
	name := ""

	var ArrayDefinitionResourceReference []build.DefinitionResourceReference

	defResourceRef := build.DefinitionResourceReference{
		Type:       &resourceRefType,
		Authorized: &auth,
		Name:       &name,
		Id:         variableGroupID,
	}

	ArrayDefinitionResourceReference = append(ArrayDefinitionResourceReference, defResourceRef)

	definitionResourceReference, err := clients.BuildClient.AuthorizeProjectResources(
		clients.Ctx, build.AuthorizeProjectResourcesArgs{
			Resources: &ArrayDefinitionResourceReference,
			Project:   project,
		})

	return definitionResourceReference, err
}

// Convert AzDO data structure allow_access to internal Terraform data structure
func flattenAllowAccess(d *schema.ResourceData, definitionResource *[]build.DefinitionResourceReference) {
	variableGroupID := d.Id()
	var allowAccess = false
	if definitionResource != nil {
		for _, authResource := range *definitionResource {
			if variableGroupID == *authResource.Id {
				allowAccess = *authResource.Authorized
			}
		}
	}
	d.Set(vgAllowAccess, allowAccess)
}

func searchAzureKVSecrets(clients *client.AggregatedClient, projectID, kvName, serviceEndpointID string, variables []interface{}, depth int) (kvSecrets map[string]interface{}, invalidSecrets []string, error error) {
	var token, loop, azkvSecretsRaw = "", 0, &KeyVaultSecretResult{}
	kvSecrets = make(map[string]interface{})
	invalidSecrets = make([]string, 0)

	secretNames := make(map[string]string)
	for _, val := range variables {
		name := val.(map[string]interface{})["name"].(string)
		secretNames[name] = name
	}
	for {
		kvSecretsMap := make(map[string]taskagent.AzureKeyVaultVariableValue)
		if azKVSecrets, err := getKVSecretServiceEndpointProxy(clients, kvName, projectID, serviceEndpointID, token); err == nil {
			azkvSecretsRaw, token, err = parseKVSecretResp(azKVSecrets)
			if err != nil {
				return nil, nil, err
			}
			for _, secret := range *azkvSecretsRaw.Value {
				name := getSecretName(*secret.ID)
				kvVariable := taskagent.AzureKeyVaultVariableValue{
					Value:       nil,
					ContentType: secret.ContentType,
					IsSecret:    converter.Bool(true),
					Enabled:     secret.Enabled,
				}

				if secret.ContentType == nil {
					kvVariable.ContentType = converter.String("")
				}
				if secret.Expire != nil {
					kvVariable.Expires = &azuredevops.Time{
						Time: time.Unix(*secret.Expire, 0),
					}
				}
				kvSecretsMap[name] = kvVariable
			}

			// search secret
			for name, secret := range kvSecretsMap {
				if len(secretNames) == 0 {
					break
				}
				if _, ok := secretNames[name]; ok && *secret.Enabled {
					kvSecrets[name] = secret
					delete(secretNames, name)
				}
			}

			// stop search
			if token == "" || loop == depth || len(secretNames) == 0 {
				for k := range secretNames {
					invalidSecrets = append(invalidSecrets, k)
				}
				return kvSecrets, invalidSecrets, err
			}
			loop++
		} else {
			return nil, nil, fmt.Errorf("Failed to get the Azure Key vault secrets. %v ", err)
		}
	}
}

func parseKVSecretResp(azKVSecrets *serviceendpoint.ServiceEndpointRequestResult) (*KeyVaultSecretResult, string, error) {
	if azKVSecrets != nil && *azKVSecrets.StatusCode == "ok" {
		var kvSecrets KeyVaultSecretResult
		secretJson := azKVSecrets.Result.([]interface{})[0].(string)
		if err := json.Unmarshal([]byte(secretJson), &kvSecrets); err != nil {
			return nil, "", fmt.Errorf("Failed to parse the Azure Key value secrets. Service response: %s . %v ", secretJson, err)
		}

		token, err := getSkipToken(kvSecrets.NextLink)
		if err != nil {
			return nil, "", fmt.Errorf(" falied to get skip token, error: %+v", err)
		}
		return &kvSecrets, token, nil
	}
	return nil, "", fmt.Errorf("Failed to get the Azure Key value. Error: ( code: %s, messge: %s )", *azKVSecrets.StatusCode, *azKVSecrets.ErrorMessage)
}

func getKVSecretServiceEndpointProxy(clients *client.AggregatedClient, kvName string, projectID string, serviceEndpointID string, token string) (*serviceendpoint.ServiceEndpointRequestResult, error) {
	reqArgs := serviceendpoint.ExecuteServiceEndpointRequestArgs{
		ServiceEndpointRequest: &serviceendpoint.ServiceEndpointRequest{
			DataSourceDetails: &serviceendpoint.DataSourceDetails{
				DataSourceName: converter.String("AzureKeyVaultSecrets"),
				Parameters: &map[string]string{
					"KeyVaultName": kvName,
				},
			},
			ResultTransformationDetails: &serviceendpoint.ResultTransformationDetails{},
		},
		Project:    &projectID,
		EndpointId: &serviceEndpointID,
	}
	if token != "" {
		(*reqArgs.ServiceEndpointRequest.DataSourceDetails.Parameters)["SkipToken"] = token
		reqArgs.ServiceEndpointRequest.DataSourceDetails.DataSourceName = converter.String("AzureKeyVaultSecretsWithSkipToken")
	}
	return clients.ServiceEndpointClient.ExecuteServiceEndpointRequest(clients.Ctx, reqArgs)
}

func getSecretName(secretID string) (secret string) {
	if len(secretID) == 0 {
		return ""
	}
	secretURL := strings.Split(secretID, "/")
	return secretURL[len(secretURL)-1]
}

func getSkipToken(link *string) (string, error) {
	if link == nil || len(*link) == 0 {
		return "", nil
	}
	linkUrl, err := url.Parse(*link)
	if err != nil {
		return "", err
	}

	params, err := url.ParseQuery(linkUrl.RawQuery)
	if err != nil {
		return "", err
	}

	token := params["$skiptoken"]
	if len(token) > 0 {
		return token[0], nil
	}
	//if skip token not found, just return "" as the skip token
	return "", nil
}
