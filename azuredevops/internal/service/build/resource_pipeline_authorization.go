package build

import (
	"fmt"

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
				ValidateFunc: validation.StringInSlice([]string{"endpoint", "queue", "variablegroup", "environment"}, false),
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

	pipePermissionParams := pipelinepermissions.UpdatePipelinePermisionsForResourceArgs{
		Project:      converter.String(d.Get("project_id").(string)),
		ResourceType: converter.String(d.Get("type").(string)),
		ResourceId:   converter.String(d.Get("resource_id").(string)),
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

	resp, err := clients.PipelinePermissionsClient.GetPipelinePermissionsForResource(clients.Ctx,
		pipelinepermissions.GetPipelinePermissionsForResourceArgs{
			Project:      converter.String(d.Get("project_id").(string)),
			ResourceType: converter.String(d.Get("type").(string)),
			ResourceId:   converter.String(d.Get("resource_id").(string)),
		},
	)
	if err != nil {
		return fmt.Errorf("%+v", err)
	}

	if resp == nil {
		d.SetId("")
		return nil
	}

	d.Set("resource_id", resp.Resource.Id)
	d.Set("type", resp.Resource.Type)
	if resp.Pipelines != nil && len(*resp.Pipelines) > 0 {
		pipeAuth := (*resp.Pipelines)[0]
		d.Set("pipeline_id", pipeAuth.Id)
	}

	d.SetId(*resp.Resource.Id)
	return nil
}

func resourcePipelineAuthorizationDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	pipePermissionParams := pipelinepermissions.UpdatePipelinePermisionsForResourceArgs{
		Project:      converter.String(d.Get("project_id").(string)),
		ResourceType: converter.String(d.Get("type").(string)),
		ResourceId:   converter.String(d.Get("resource_id").(string)),
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
