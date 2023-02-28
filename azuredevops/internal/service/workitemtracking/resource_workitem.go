package workitemtracking

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/webapi"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/datahelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceWorkItem schema and implementation for project workitem ressource
func ResourceWorkItem() *schema.Resource {
	return &schema.Resource{
		Create:   ResourceWorkItemCreate,
		Read:     ResourceWorkItemRead,
		Update:   ResourceWorkItemUpdate,
		Delete:   ResourceWorkItemDelete,
		Importer: tfhelper.ImportProjectQualifiedResourceInteger(),
		Schema: map[string]*schema.Schema{
			"title": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Required:     true,
			},
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Required:     true,
				ForceNew:     true,
			},
			"type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				ForceNew:     true,
				Required:     true,
				Description:  "Type of the Work Item",
			},
			"state": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Optional:     true,
				Description:  "State of the Ticket",
			},
			"custom_fields": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
			},
		},
	}
}

var systemFieldMapping = map[string]string{
	"System.State":        "state",
	"System.Title":        "title",
	"System.WorkItemType": "type",
}

// ResourceWorkItemCreateOrUpdate create workitem
func ResourceWorkItemCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	operations := GetPatchOperations(d)

	args := workitemtracking.CreateWorkItemArgs{
		Project:  converter.String(d.Get("project_id").(string)),
		Type:     converter.String(d.Get("type").(string)),
		Document: &operations,
	}
	workitem, err := clients.WorkItemTrackingClient.CreateWorkItem(clients.Ctx, args)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(*workitem.Id))
	return ResourceWorkItemRead(d, m)
}

// ResourceWorkItemRead read workitem
func ResourceWorkItemRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	id, _ := strconv.Atoi(d.Id())
	args := workitemtracking.GetWorkItemArgs{
		Id: &id,
	}
	workitem, err := clients.WorkItemTrackingClient.GetWorkItem(clients.Ctx, args)
	if err != nil {
		return err
	}

	GetFields(d, workitem.Fields)

	return nil
}

// ResourceWorkItemUpdate update a workitem
func ResourceWorkItemUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	project := d.Get("project_id").(string)
	id, _ := strconv.Atoi(d.Id())

	operations := GetPatchOperations(d)

	args := workitemtracking.UpdateWorkItemArgs{
		Id:       &id,
		Project:  &project,
		Document: &operations,
	}
	workitem, err := clients.WorkItemTrackingClient.UpdateWorkItem(clients.Ctx, args)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", *workitem.Id))
	return ResourceWorkItemRead(d, m)
}

// ResourceWorkItemDelete remove workitem
func ResourceWorkItemDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	id, errConvert := strconv.Atoi(d.Id())
	if errConvert != nil {
		return fmt.Errorf("Error getting Workitem Id: %+v", errConvert)
	}
	args := workitemtracking.DeleteWorkItemArgs{
		Id: &id,
	}
	_, err := clients.WorkItemTrackingClient.DeleteWorkItem(clients.Ctx, args)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	d.SetId("")
	return nil
}

func GetPatchOperations(d *schema.ResourceData) []webapi.JsonPatchOperation {
	var operations []webapi.JsonPatchOperation
	operations = SetSystemFields(d, operations)
	operations = SetCustomFields(d, operations)
	operations = SetTags(d, operations)
	return operations
}

func SetCustomFields(d *schema.ResourceData, operations []webapi.JsonPatchOperation) []webapi.JsonPatchOperation {
	custom_fields := d.Get("custom_fields").(map[string]interface{})
	for customFieldName, customFieldValue := range *&custom_fields {
		operations = append(operations, webapi.JsonPatchOperation{
			Op:    &webapi.OperationValues.Add,
			From:  nil,
			Path:  converter.String("/fields/Custom." + customFieldName),
			Value: customFieldValue,
		})
	}
	return operations
}

func SetSystemFields(d *schema.ResourceData, operations []webapi.JsonPatchOperation) []webapi.JsonPatchOperation {
	systemFieldReverseMapping := datahelper.ReverseMap(systemFieldMapping)
	for terraformProperty, apiName := range *&systemFieldReverseMapping {
		value := d.Get(terraformProperty).(string)
		if value != "" {
			operations = append(operations, webapi.JsonPatchOperation{
				Op:    &webapi.OperationValues.Add,
				From:  nil,
				Path:  converter.String("/fields/" + apiName),
				Value: value,
			})
		}
	}
	return operations
}

func SetTags(d *schema.ResourceData, operations []webapi.JsonPatchOperation) []webapi.JsonPatchOperation {
	tags := d.Get("tags").([]interface{})
	if len(tags) == 0 {
		return operations
	}
	operations = append(operations, webapi.JsonPatchOperation{
		Op:    &webapi.OperationValues.Add,
		From:  nil,
		Path:  converter.String("/fields/System.Tags"),
		Value: strings.Join(tfhelper.ExpandStringList(tags), "; "),
	})

	return operations
}

func GetFields(d *schema.ResourceData, m *map[string]interface{}) {
	custom_fields := make(map[string]interface{})
	for key, value := range *m {
		v, ok := systemFieldMapping[key]
		if ok {
			d.Set(v, value)
		} else if strings.HasPrefix(key, "Custom.") {
			key_without_custom := strings.ReplaceAll(key, "Custom.", "")
			custom_fields[key_without_custom] = value
		} else if "System.Tags" == key {
			d.Set("tags", strings.Split(value.(string), "; "))
		}
	}
	d.Set("custom_fields", custom_fields)
}
