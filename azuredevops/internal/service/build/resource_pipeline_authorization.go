package build

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pipeline_project_id": {
				Type:     schema.TypeString,
				Optional: true,
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
				ValidateFunc: validation.StringInSlice([]string{"endpoint", "queue", "variablegroup", "securefile", "environment", "repository"}, false),
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
	pipelineProjectId := projectId
	if d.Get("pipeline_project_id").(string) != "" {
		pipelineProjectId = d.Get("pipeline_project_id").(string)
	}

	resType := d.Get("type").(string)
	resId := d.Get("resource_id").(string)

	if strings.EqualFold(resType, "repository") {
		resId = projectId + "." + resId
	}

	pipePermissionParams := pipelinepermissions.UpdatePipelinePermisionsForResourceArgs{
		Project:      &pipelineProjectId,
		ResourceType: &resType,
		ResourceId:   &resId,
	}

	if v, ok := d.GetOk("pipeline_id"); ok {
		pipePermissionParams.ResourceAuthorization = &pipelinepermissions.ResourcePipelinePermissions{
			Pipelines: &[]pipelinepermissions.PipelinePermission{{
				Authorized: converter.ToPtr(true),
				Id:         converter.ToPtr(v.(int)),
			}},
		}
	} else {
		pipePermissionParams.ResourceAuthorization = &pipelinepermissions.ResourcePipelinePermissions{
			AllPipelines: &pipelinepermissions.Permission{
				Authorized: converter.ToPtr(true),
			},
		}
	}

	response, err := clients.PipelinePermissionsClient.UpdatePipelinePermisionsForResource(
		clients.Ctx,
		pipePermissionParams,
	)
	if err != nil {
		return fmt.Errorf("creating authorized resource: %+v", err)
	}

	// ensure authorization is complete
	stateConf := &retry.StateChangeConf{
		ContinuousTargetOccurence: 1,
		Delay:                     5 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"waiting"},
		Target:                    []string{"succeed", "failed"},
		Refresh:                   checkPipelineAuthorization(clients, d, pipePermissionParams),
		Timeout:                   d.Timeout(schema.TimeoutCreate),
	}

	if _, err := stateConf.WaitForStateContext(clients.Ctx); err != nil {
		return fmt.Errorf("waiting for pipeline authorization ready. %v ", err)
	}

	d.SetId(*response.Resource.Id)

	return resourcePipelineAuthorizationRead(d, m)
}

func resourcePipelineAuthorizationRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectId := d.Get("project_id").(string)
	pipelineProjectId := projectId
	if d.Get("pipeline_project_id").(string) != "" {
		pipelineProjectId = d.Get("pipeline_project_id").(string)
	}

	resType := d.Get("type").(string)
	resId := d.Get("resource_id").(string)

	if strings.EqualFold(resType, "repository") {
		resId = projectId + "." + resId
	}

	resp, err := clients.PipelinePermissionsClient.GetPipelinePermissionsForResource(clients.Ctx,
		pipelinepermissions.GetPipelinePermissionsForResourceArgs{
			Project:      &pipelineProjectId,
			ResourceType: &resType,
			ResourceId:   &resId,
		},
	)
	if err != nil {
		return fmt.Errorf("%+v", err)
	}

	if resp == nil || (resp.AllPipelines == nil && len(*resp.Pipelines) == 0) {
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
		exist := false
		for _, pipe := range *resp.Pipelines {
			if *pipe.Id == d.Get("pipeline_id").(int) {
				exist = true
			}
		}
		if !exist {
			d.Set("pipeline_id", nil)
		}
	}

	return nil
}

func resourcePipelineAuthorizationDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectId := d.Get("project_id").(string)
	pipelineProjectId := projectId
	if d.Get("pipeline_project_id").(string) != "" {
		pipelineProjectId = d.Get("pipeline_project_id").(string)
	}

	resType := d.Get("type").(string)
	resId := d.Get("resource_id").(string)

	if strings.EqualFold(resType, "repository") {
		resId = projectId + "." + resId
	}

	pipePermissionParams := pipelinepermissions.UpdatePipelinePermisionsForResourceArgs{
		Project:      &pipelineProjectId,
		ResourceType: &resType,
		ResourceId:   &resId,
	}

	if v, ok := d.GetOk("pipeline_id"); ok {
		pipePermissionParams.ResourceAuthorization = &pipelinepermissions.ResourcePipelinePermissions{
			Pipelines: &[]pipelinepermissions.PipelinePermission{{
				Authorized: converter.ToPtr(false),
				Id:         converter.ToPtr(v.(int)),
			}},
		}
	} else {
		pipePermissionParams.ResourceAuthorization = &pipelinepermissions.ResourcePipelinePermissions{
			AllPipelines: &pipelinepermissions.Permission{
				Authorized: converter.ToPtr(false),
			},
		}
	}

	_, err := clients.PipelinePermissionsClient.UpdatePipelinePermisionsForResource(
		clients.Ctx,
		pipePermissionParams)
	if err != nil {
		return fmt.Errorf("deleting authorized resource: %+v", err)
	}

	return nil
}

func checkPipelineAuthorization(clients *client.AggregatedClient, d *schema.ResourceData, params pipelinepermissions.UpdatePipelinePermisionsForResourceArgs) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		projectId := d.Get("project_id").(string)
		resourceType := d.Get("type").(string)
		resourceId := d.Get("resource_id").(string)
		pipelineProjectId := projectId
		if d.Get("pipeline_project_id").(string) != "" {
			pipelineProjectId = d.Get("pipeline_project_id").(string)
		}

		if strings.EqualFold(resourceType, "repository") {
			resourceId = projectId + "." + resourceId
		}

		resp, err := clients.PipelinePermissionsClient.GetPipelinePermissionsForResource(clients.Ctx,
			pipelinepermissions.GetPipelinePermissionsForResourceArgs{
				Project:      &pipelineProjectId,
				ResourceType: &resourceType,
				ResourceId:   &resourceId,
			},
		)
		if err != nil {
			return nil, "failed", err
		}

		// check pipeline authorization if pipeline_id exist
		if pipeId, ok := d.GetOk("pipeline_id"); ok {
			if resp.Pipelines != nil && len(*resp.Pipelines) > 0 {
				for _, pipe := range *resp.Pipelines {
					if *pipe.Id == pipeId.(int) {
						return resp, "succeed", err
					}
				}
				// reapply for authorization
				_, err = clients.PipelinePermissionsClient.UpdatePipelinePermisionsForResource(
					clients.Ctx,
					params,
				)
				return nil, "waiting", err
			}
		} else {
			// check all pipeline authorization
			if resp.AllPipelines != nil && resp.AllPipelines.Authorized != nil && *resp.AllPipelines.Authorized {
				return resp, "succeed", err
			}
			// reapply for authorization
			_, err = clients.PipelinePermissionsClient.UpdatePipelinePermisionsForResource(
				clients.Ctx,
				params,
			)
			return nil, "waiting", err
		}

		return resp, "succeed", nil
	}
}
