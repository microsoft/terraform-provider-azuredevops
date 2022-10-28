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
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/datahelper"
)

// ResourceWorkItem schema and implementation for project workitem ressource
func ResourceWorkItem() *schema.Resource {
	return &schema.Resource{
		Create: ResourceWorkItemCreate,
		Read:   ResourceWorkItemRead,
		Update: ResourceWorkItemUpdate,
		Delete: ResourceWorkItemDelete,
		Schema: map[string]*schema.Schema{
			"title": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Required:     true,
				Optional:     false,
			},
			"project": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Required:     true,
				ForceNew:     true,
			},
			"type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Optional:     false,
				ForceNew:     true,
				Required:     true,
				Description:  "Type of the Work Item",
			},
			"state": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Optional:     true,
				Description:  "state of the Ticket",
			},
			"custom_fields": {
				Type:     schema.TypeMap,
				Optional: true,
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
	project := d.Get("project").(string)
	workItemType := d.Get("type").(string)

	operations := GetPatchOperations(d)

	args := workitemtracking.CreateWorkItemArgs{
		Project:  &project,
		Type:     &workItemType,
		Document: &operations,
	}
	workitem, err := clients.WorkItemTrackingClient.CreateWorkItem(clients.Ctx, args)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", *workitem.Id))
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
	project := d.Get("project").(string)
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
	id, _ := strconv.Atoi(d.Id())
	args := workitemtracking.DeleteWorkItemArgs{
		Id: &id,
	}
	_, err := clients.WorkItemTrackingClient.DeleteWorkItem(clients.Ctx, args)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func GetPatchOperations(d *schema.ResourceData) []webapi.JsonPatchOperation {
	var operations []webapi.JsonPatchOperation
	operations = SetSystemFields(d, operations)
	operations = SetCustomFields(d, operations)
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

func GetFields(d *schema.ResourceData, m *map[string]interface{}) {

	custom_fields := make(map[string]interface{})
	for key, value := range *m {
		v, ok := systemFieldMapping[key]
		if ok {
			d.Set(v, value)
		} else if strings.HasPrefix(key, "Custom.") {
			key_without_custom := strings.ReplaceAll(key, "Custom.", "")
			custom_fields[key_without_custom] = value
		}

	}

	d.Set("custom_fields", custom_fields)
}
