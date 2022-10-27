package workitemtracking

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/webapi"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceWorkItem schema and implementation for project workitem ressource
func ResourceWorkItem() *schema.Resource {
	return &schema.Resource{
		Create: ResourceWorkItemCreateOrUpdate,
		Read:   ResourceWorkItemRead,
		Update: ResourceWorkItemCreateOrUpdate,
		Delete: ResourceWorkItemDelete,
		Schema: map[string]*schema.Schema{
			"title": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Required:     true,
				Optional:     false,
				ForceNew:     true,
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
				Optional:     true,
				ForceNew:     true,
				Description:  "Type of the Work Item",
			},
			"state": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Optional:     true,
				Description:  "state of the Ticket",
			},
			"custom_fields": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Optional:     true,
				ForceNew:     true,
			},
		},
	}
}

// ResourceWorkItemCreateOrUpdate create or update workitem
func ResourceWorkItemCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	project := d.Get("project").(string)
	workItemType := d.Get("type").(string)
	state := d.Get("state").(string)

	title := d.Get("title").(string)
	var operations []webapi.JsonPatchOperation
	operations = append(operations, webapi.JsonPatchOperation{
		Op:    &webapi.OperationValues.Add,
		From:  nil,
		Path:  converter.String("/fields/System.Title"),
		Value: title,
	})
	if state != "" {
		operations = append(operations, webapi.JsonPatchOperation{
			Op:    &webapi.OperationValues.Add,
			From:  nil,
			Path:  converter.String("/fields/System.State"),
			Value: state,
		})
	}

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

	mapSystemFields(d, workitem.Fields)

	return nil
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

func mapSystemFields(d *schema.ResourceData, m *map[string]interface{}) {
	biMap := map[string]string{
		"System.State": "state",
		"System.Title": "title",
	}

	for key, value := range *m {
		v, ok := biMap[key]
		if ok {
			d.Set(v, value)
		}
	}
}
