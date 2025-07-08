package workitemtracking

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

var systemFieldMapping = map[string]string{
	"System.State":         "state",
	"System.Title":         "title",
	"System.WorkItemType":  "type",
	"System.AreaPath":      "area_path",
	"System.IterationPath": "iteration_path",

	"System.Parent": "parent_id",
}

var fieldMapping = map[string]string{
	"title":          "System.Title",
	"type":           "System.WorkItemType",
	"state":          "System.State",
	"area_path":      "System.AreaPath",
	"iteration_path": "System.IterationPath",

	"parent_id": "System.Parent",
}

func ResourceWorkItem() *schema.Resource {
	return &schema.Resource{
		Create:   resourceWorkItemCreate,
		Read:     resourceWorkItemRead,
		Update:   resourceWorkItemUpdate,
		Delete:   resourceWorkItemDelete,
		Importer: tfhelper.ImportProjectQualifiedResource(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"title": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"type": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
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
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
			},
			"area_path": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"iteration_path": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"parent_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"relations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceWorkItemCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	orgName := strings.Split(clients.OrganizationURL, "/")[3]

	var operations []webapi.JsonPatchOperation
	operations = expandSystemFields(d, operations, orgName)
	operations = expandCustomFields(d, operations)
	operations = expandTags(d, operations, webapi.OperationValues.Add)

	args := workitemtracking.CreateWorkItemArgs{
		Project:  converter.String(d.Get("project_id").(string)),
		Type:     converter.String(d.Get("type").(string)),
		Document: &operations,
	}
	workItem, err := clients.WorkItemTrackingClient.CreateWorkItem(clients.Ctx, args)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(*workItem.Id))
	return resourceWorkItemRead(d, m)
}

func resourceWorkItemRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	args := workitemtracking.GetWorkItemArgs{
		Project: converter.String(d.Get("project_id").(string)),
		Id:      &id,
		Expand:  converter.ToPtr(workitemtracking.WorkItemExpandValues.All),
	}
	workItem, err := clients.WorkItemTrackingClient.GetWorkItem(clients.Ctx, args)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	if workItem != nil {
		if workItem.Url != nil {
			d.Set("url", *workItem.Url)
			flattenFields(d, workItem.Fields)
		}

		var relations []map[string]interface{}
		if workItem.Relations != nil {
			for _, v := range *workItem.Relations {
				relations = append(relations, map[string]interface{}{
					"rel": v.Rel,
					"url": v.Url,
				})
			}
		}
		d.Set("relations", relations)
	}
	return nil
}

// resourceWorkItemUpdate update a workitem
func resourceWorkItemUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	project := d.Get("project_id").(string)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	orgName := strings.Split(clients.OrganizationURL, "/")[3]

	var operations []webapi.JsonPatchOperation
	operations = expandSystemFields(d, operations, orgName)
	operations = expandCustomFields(d, operations)
	operations = expandTags(d, operations, webapi.OperationValues.Replace)

	args := workitemtracking.UpdateWorkItemArgs{
		Id:       &id,
		Project:  &project,
		Document: &operations,
	}
	_, err = clients.WorkItemTrackingClient.UpdateWorkItem(clients.Ctx, args)
	if err != nil {
		return fmt.Errorf("Update work item. Project ID: %s, Work Item: %s, Error: %+v", project, d.Id(), err)
	}

	return resourceWorkItemRead(d, m)
}

// resourceWorkItemDelete remove workitem
func resourceWorkItemDelete(d *schema.ResourceData, m interface{}) error {
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
	return nil
}

func expandCustomFields(d *schema.ResourceData, operations []webapi.JsonPatchOperation) []webapi.JsonPatchOperation {
	customFields := d.Get("custom_fields").(map[string]interface{})
	for customFieldName, customFieldValue := range customFields {
		operations = append(operations, webapi.JsonPatchOperation{
			Op:    &webapi.OperationValues.Add,
			From:  nil,
			Path:  converter.String("/fields/Custom." + customFieldName),
			Value: customFieldValue,
		})
	}
	return operations
}

func expandSystemFields(d *schema.ResourceData, operations []webapi.JsonPatchOperation, organizationName string) []webapi.JsonPatchOperation {
	for terraformProperty, apiName := range fieldMapping {
		if terraformProperty == "parent_id" {
			if d.HasChange("parent_id") {
				oldParentId, newParentId := d.GetChange("parent_id")
				if oldParentId.(int) > 0 {
					relations := d.Get("relations").([]interface{})

					// find the parent relationship and delete it
					for idx, relation := range relations {
						if v, ok := relation.(map[string]interface{})["rel"]; ok && strings.EqualFold(v.(string), "System.LinkTypes.Hierarchy-Reverse") {
							operations = append(operations, webapi.JsonPatchOperation{
								Op:   &webapi.OperationValues.Remove,
								From: nil,
								Path: converter.String(fmt.Sprintf("/relations/%d", idx)),
							})
						}
					}
				}

				if newParentId.(int) > 0 {
					operations = append(operations, webapi.JsonPatchOperation{
						Op:   &webapi.OperationValues.Add,
						From: nil,
						Path: converter.String("/relations/-"),
						Value: &map[string]string{
							"rel": "System.LinkTypes.Hierarchy-Reverse",
							"url": fmt.Sprintf("https://dev.azure.com/%s/_apis/wit/workItems/%d", organizationName, newParentId.(int)),
						},
					})
				}
			}
		} else {
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
	}
	return operations
}

func expandTags(d *schema.ResourceData, operations []webapi.JsonPatchOperation, op webapi.Operation) []webapi.JsonPatchOperation {
	tags := d.Get("tags").(*schema.Set).List()
	if len(tags) == 0 {
		operations = append(operations, webapi.JsonPatchOperation{
			Op:    &op,
			From:  nil,
			Path:  converter.String("/fields/System.Tags"),
			Value: "",
		})
	} else {
		operations = append(operations, webapi.JsonPatchOperation{
			Op:    &op,
			From:  nil,
			Path:  converter.String("/fields/System.Tags"),
			Value: strings.Join(tfhelper.ExpandStringList(tags), "; "),
		})
	}

	return operations
}

func flattenFields(d *schema.ResourceData, m *map[string]interface{}) {
	customFields := make(map[string]interface{})
	for key, value := range *m {
		if v, ok := systemFieldMapping[key]; ok {
			d.Set(v, value)
		} else if strings.HasPrefix(key, "Custom.") {
			keyWithoutCustom := strings.ReplaceAll(key, "Custom.", "")
			customFields[keyWithoutCustom] = value
		} else if "System.Tags" == key {
			d.Set("tags", strings.Split(value.(string), "; "))
		}
	}
	d.Set("custom_fields", customFields)
}
