package build

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/pipelinepermissions"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourcePipelineAuthorization() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelineAuthorizationCreateUpdate,
		Read:   resourcePipelineAuthorizationRead,
		Delete: resourcePipelineAuthorizationDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"endpoint", "queue", "variablegroup", "environment", "repository"}, false),
			},
			"pipeline_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
		},
	}
}

func resourcePipelineAuthorizationCreateUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectId := d.Get("project_id").(string)
	resType := d.Get("type").(string)
	resId := d.Get("resource_id").(string)

	if strings.EqualFold(resType, "repository") {
		resId = projectId + "." + resId
	}

	pipePermissionParams := pipelinepermissions.UpdatePipelinePermisionsForResourceArgs{
		Project:      &projectId,
		ResourceType: &resType,
		ResourceId:   &resId,
	}

	if v, ok := d.GetOk("pipeline_id"); ok {
		pipePermissionParams.ResourceAuthorization = &pipelinepermissions.ResourcePipelinePermissions{
			Pipelines: &[]pipelinepermissions.PipelinePermission{{
				Authorized: converter.ToPtr(true),
				Id:         converter.ToPtr(v.(int)),
			}}}
	} else {
		pipePermissionParams.ResourceAuthorization = &pipelinepermissions.ResourcePipelinePermissions{
			AllPipelines: &pipelinepermissions.Permission{
				Authorized: converter.ToPtr(true),
			}}
	}

	response, err := clients.PipelinePermissionsClient.UpdatePipelinePermisionsForResource(
		clients.Ctx,
		pipePermissionParams,
	)

	if err != nil {
		return fmt.Errorf(" creating authorized resource: %+v", err)
	}
	d.SetId(*response.Resource.Id)

	return resourcePipelineAuthorizationRead(d, m)
}

func resourcePipelineAuthorizationRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectId := d.Get("project_id").(string)
	resType := d.Get("type").(string)
	resId := d.Get("resource_id").(string)

	if strings.EqualFold(resType, "repository") {
		resId = projectId + "." + resId
	}

	resp, err := clients.PipelinePermissionsClient.GetPipelinePermissionsForResource(clients.Ctx,
		pipelinepermissions.GetPipelinePermissionsForResourceArgs{
			Project:      &projectId,
			ResourceType: &resType,
			ResourceId:   &resId,
		},
	)
	if err != nil {
		return fmt.Errorf("%+v", err)
	}

	if resp == nil {
		d.SetId("")
		return nil
	}

	d.Set("type", resp.Resource.Type)
	d.Set("resource_id", resp.Resource.Id)
	if strings.EqualFold(*resp.Resource.Type, "repository") {
		resIds := strings.Split(*resp.Resource.Id, ".")
		d.Set("resource_id", resIds[1])
	}

	if resp.Pipelines != nil && len(*resp.Pipelines) > 0 {
		var exist = false
		for _, pipe := range *resp.Pipelines {
			if *pipe.Id == d.Get("pipeline_id").(int) {
				exist = true
			}
		}
		if !exist {
			d.Set("pipeline_id", nil)
		}
	}

	d.SetId(*resp.Resource.Id)
	return nil
}

func resourcePipelineAuthorizationDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectId := d.Get("project_id").(string)
	resType := d.Get("type").(string)
	resId := d.Get("resource_id").(string)

	if strings.EqualFold(resType, "repository") {
		resId = projectId + "." + resId
	}

	pipePermissionParams := pipelinepermissions.UpdatePipelinePermisionsForResourceArgs{
		Project:      &projectId,
		ResourceType: &resType,
		ResourceId:   &resId,
	}

	if v, ok := d.GetOk("pipeline_id"); ok {
		pipePermissionParams.ResourceAuthorization = &pipelinepermissions.ResourcePipelinePermissions{
			Pipelines: &[]pipelinepermissions.PipelinePermission{{
				Authorized: converter.ToPtr(false),
				Id:         converter.ToPtr(v.(int)),
			}}}
	} else {
		pipePermissionParams.ResourceAuthorization = &pipelinepermissions.ResourcePipelinePermissions{
			AllPipelines: &pipelinepermissions.Permission{
				Authorized: converter.ToPtr(false),
			}}
	}

	_, err := clients.PipelinePermissionsClient.UpdatePipelinePermisionsForResource(
		clients.Ctx,
		pipePermissionParams)

	if err != nil {
		return fmt.Errorf(" deleting authorized resource: %+v", err)
	}

	return nil
}
