package azuredevops

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"
	"strconv"
)

func resourceVariableGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceVariableGroupCreate,
		Read:   resourceVariableGroupRead,
		Update: resourceVariableGroupUpdate,
		Delete: resourceVariableGroupDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, variableGroupID, err := ParseImportedProjectIDAndVariableGroupID(meta.(*config.AggregatedClient), d.Id())
				if err != nil {
					return nil, fmt.Errorf("Error parsing the variable group ID from the Terraform resource data: %v", err)
				}
				d.Set("project_id", projectID)
				d.SetId(fmt.Sprintf("%d", variableGroupID))

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.UUID,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"allow_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"variable": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"is_secret": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
				Required: true,
				MinItems: 1,
				Set: func(i interface{}) int {
					item := i.(map[string]interface{})
					return schema.HashString(item["name"].(string))
				},
			},
		},
	}
}

func resourceVariableGroupCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	variableGroupParameters, projectID := expandVariableGroupParameters(d)

	addedVariableGroup, err := createVariableGroup(clients, variableGroupParameters, projectID)
	if err != nil {
		return fmt.Errorf("Error creating variable group in Azure DevOps: %+v", err)
	}

	flattenVariableGroup(d, addedVariableGroup, projectID)

	// Update Allow Access with definition Reference
	definitionResourceReferenceArgs := expandDefinitionResourceAuth(d, addedVariableGroup)
	definitionResourceReference, err := updateDefinitionResourceAuth(clients, definitionResourceReferenceArgs, projectID)
	if err != nil {
		return fmt.Errorf("Error creating definitionResourceReference Azure DevOps object: %+v", err)
	}

	flattenAllowAccess(d, definitionResourceReference)

	return nil
}

func resourceVariableGroupRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)

	projectID, variableGroupID, err := tfhelper.ParseProjectIDAndResourceID(d)
	if err != nil {
		return fmt.Errorf("Error parsing the variable group ID from the Terraform resource data: %v", err)
	}

	variableGroup, err := clients.TaskAgentClient.GetVariableGroup(
		clients.Ctx,
		taskagent.GetVariableGroupArgs{
			GroupId: &variableGroupID,
			Project: &projectID,
		},
	)
	if err != nil {
		return fmt.Errorf("Error looking up variable group given ID (%v) and project ID (%v): %v", variableGroupID, projectID, err)
	}

	flattenVariableGroup(d, variableGroup, &projectID)

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
	clients := m.(*config.AggregatedClient)
	variableGroupParams, projectID := expandVariableGroupParameters(d)

	_, variableGroupID, err := tfhelper.ParseProjectIDAndResourceID(d)
	if err != nil {
		return fmt.Errorf("Error parsing the variable group ID from the Terraform resource data: %v", err)
	}

	updatedVariableGroup, err := updateVariableGroup(clients, variableGroupParams, &variableGroupID, projectID)
	if err != nil {
		return fmt.Errorf("Error updating variable group in Azure DevOps: %+v", err)
	}

	flattenVariableGroup(d, updatedVariableGroup, projectID)

	// Update Allow Access
	definitionResourceReferenceArgs := expandDefinitionResourceAuth(d, updatedVariableGroup)
	definitionResourceReference, err := updateDefinitionResourceAuth(clients, definitionResourceReferenceArgs, projectID)
	if err != nil {
		return fmt.Errorf("Error updating definitionResourceReference Azure DevOps object: %+v", err)
	}

	flattenAllowAccess(d, definitionResourceReference)

	return nil
}

func resourceVariableGroupDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	projectID, variableGroupID, err := tfhelper.ParseProjectIDAndResourceID(d)
	if err != nil {
		return fmt.Errorf("Error parsing the variable group ID from the Terraform resource data: %v", err)
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
func createVariableGroup(clients *config.AggregatedClient, variableGroupParams *taskagent.VariableGroupParameters, project *string) (*taskagent.VariableGroup, error) {
	createdVariableGroup, err := clients.TaskAgentClient.AddVariableGroup(
		clients.Ctx,
		taskagent.AddVariableGroupArgs{
			Group:   variableGroupParams,
			Project: project,
		})
	return createdVariableGroup, err
}

// Make the Azure DevOps API call to update the variable group
func updateVariableGroup(clients *config.AggregatedClient, parameters *taskagent.VariableGroupParameters, variableGroupID *int, project *string) (*taskagent.VariableGroup, error) {
	updatedVariableGroup, err := clients.TaskAgentClient.UpdateVariableGroup(
		clients.Ctx,
		taskagent.UpdateVariableGroupArgs{
			Project: project,
			GroupId: variableGroupID,
			Group:   parameters,
		})

	return updatedVariableGroup, err
}

// Make the Azure DevOps API call to delete the variable group
func deleteVariableGroup(clients *config.AggregatedClient, project *string, variableGroupID *int) error {
	err := clients.TaskAgentClient.DeleteVariableGroup(
		clients.Ctx,
		taskagent.DeleteVariableGroupArgs{
			Project: project,
			GroupId: variableGroupID,
		})

	return err
}

// Convert internal Terraform data structure to an AzDO data structure
func expandVariableGroupParameters(d *schema.ResourceData) (*taskagent.VariableGroupParameters, *string) {
	projectID := converter.String(d.Get("project_id").(string))
	variables := d.Get("variable").(*schema.Set).List()

	variableMap := make(map[string]taskagent.VariableValue)

	for _, variable := range variables {
		asMap := variable.(map[string]interface{})
		variableMap[asMap["name"].(string)] = taskagent.VariableValue{
			Value:    converter.String(asMap["value"].(string)),
			IsSecret: converter.Bool(asMap["is_secret"].(bool)),
		}
	}

	variableGroup := &taskagent.VariableGroupParameters{
		Name:        converter.String(d.Get("name").(string)),
		Description: converter.String(d.Get("description").(string)),
		Variables:   &variableMap,
	}

	return variableGroup, projectID
}

// Convert AzDO data structure to internal Terraform data structure
func flattenVariableGroup(d *schema.ResourceData, variableGroup *taskagent.VariableGroup, projectID *string) {
	d.SetId(fmt.Sprintf("%d", *variableGroup.Id))
	d.Set("name", *variableGroup.Name)
	d.Set("description", *variableGroup.Description)
	d.Set("variable", flattenVariables(variableGroup))
	d.Set("project_id", projectID)
}

// Convert AzDO Variables data structure to Terraform TypeSet
func flattenVariables(variableGroup *taskagent.VariableGroup) interface{} {
	// Preallocate list of variable prop maps
	variables := make([]map[string]interface{}, len(*variableGroup.Variables))

	index := 0
	for k, v := range *variableGroup.Variables {
		variables[index] = map[string]interface{}{
			"name":      k,
			"value":     converter.ToString(v.Value, ""),
			"is_secret": converter.ToBool(v.IsSecret, false),
		}
		index = index + 1
	}

	return variables
}

// Convert internal Terraform data structure to an AzDO data structure for Allow Access
func expandDefinitionResourceAuth(d *schema.ResourceData, createdVariableGroup *taskagent.VariableGroup) []build.DefinitionResourceReference {

	resourceRefType := "variablegroup"
	variableGroupID := strconv.Itoa(*createdVariableGroup.Id)

	var ArrayDefinitionResourceReference []build.DefinitionResourceReference

	defResourceRef := build.DefinitionResourceReference{
		Type:       &resourceRefType,
		Authorized: converter.Bool(d.Get("allow_access").(bool)),
		Name:       createdVariableGroup.Name,
		Id:         &variableGroupID,
	}

	ArrayDefinitionResourceReference = append(ArrayDefinitionResourceReference, defResourceRef)

	return ArrayDefinitionResourceReference
}

// Make the Azure DevOps API call to update the Definition resource = Allow Access
func updateDefinitionResourceAuth(clients *config.AggregatedClient, definitionResource []build.DefinitionResourceReference, project *string) (*[]build.DefinitionResourceReference, error) {

	definitionResourceReference, err := clients.BuildClient.AuthorizeProjectResources(
		clients.Ctx, build.AuthorizeProjectResourcesArgs{
			Resources: &definitionResource,
			Project:   project,
		})

	return definitionResourceReference, err
}

// Make the Azure DevOps API call to delete the resource Auth Authorized=false
func deleteDefinitionResourceAuth(clients *config.AggregatedClient, variableGroupID *string, project *string) (*[]build.DefinitionResourceReference, error) {

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
	var allowAccess bool
	if len(*definitionResource) > 0 {
		allowAccess = *(*definitionResource)[0].Authorized
	} else {
		allowAccess = false
	}
	d.Set("allow_access", allowAccess)
}

// ParseImportedProjectIDAndVariableGroupID : Parse the Id (projectId/variableGroupId) or (projectName/variableGroupId)
func ParseImportedProjectIDAndVariableGroupID(clients *config.AggregatedClient, id string) (string, int, error) {
	project, resourceID, err := tfhelper.ParseImportedID(id)
	if err != nil {
		return "", 0, err
	}

	// Get the project ID
	currentProject, err := projectRead(clients, project, project)
	if err != nil {
		return "", 0, err
	}

	return currentProject.Id.String(), resourceID, nil
}
