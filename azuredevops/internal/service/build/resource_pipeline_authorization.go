package build

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/pipelinepermissions"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

type pipelineAuthorizationId struct {
	projectId         string
	typ               string
	resourceId        string
	pipelineProjectId *string
	pipelineId        *int
}

func (pipelineAuthorizationId) pipelineProjectIdKey() string {
	return "pipeline_project_id"
}

func (pipelineAuthorizationId) pipelineIdKey() string {
	return "pipeline_id"
}

func (id pipelineAuthorizationId) id() string {
	output := fmt.Sprintf("%s/%s/%s", id.projectId, id.typ, id.resourceId)
	if id.pipelineProjectId != nil {
		output += fmt.Sprintf(";%s=%s", id.pipelineProjectIdKey(), *id.pipelineProjectId)
	}
	if id.pipelineId != nil {
		output += fmt.Sprintf(";%s=%s", id.pipelineIdKey(), strconv.Itoa(*id.pipelineId))
	}
	return output
}

func parsePipelineAuthorizationId(id string) (*pipelineAuthorizationId, error) {
	var output pipelineAuthorizationId
	segs := strings.Split(id, ";")
	ssegs := strings.Split(segs[0], "/")
	if len(ssegs) != 3 {
		return nil, fmt.Errorf("invalid id, expect base part to be `<project_id>/<type>/<resource_id>`, got=%s", segs[0])
	}
	output.projectId = ssegs[0]
	output.typ = ssegs[1]
	output.resourceId = ssegs[2]

	pairs := [][2]string{}
	for i, seg := range segs {
		if i == 0 {
			continue
		}
		k, v, ok := strings.Cut(seg, "=")
		if !ok {
			return nil, fmt.Errorf("invalid id, expect optional part to be `<key>=<value>`, got=%s", seg)
		}
		pairs = append(pairs, [2]string{k, v})
	}

	for _, pair := range pairs {
		k, v := pair[0], pair[1]
		switch k {
		case output.pipelineProjectIdKey():
			output.pipelineProjectId = &v
		case output.pipelineIdKey():
			n, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			output.pipelineId = &n
		default:
			return nil, fmt.Errorf("invalid id, unknown optional key %s", k)
		}
	}
	return &output, nil
}

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
				ValidateFunc: validation.StringInSlice([]string{"endpoint", "queue", "variablegroup", "environment", "repository"}, false),
			},
			"pipeline_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				id, err := parsePipelineAuthorizationId(d.Id())
				if err != nil {
					return nil, err
				}

				d.Set("project_id", id.projectId)
				d.Set("type", id.typ)
				d.Set("resource_id", id.resourceId)
				if id.pipelineProjectId != nil {
					d.Set("pipeline_project_id", *id.pipelineProjectId)
				}
				if id.pipelineId != nil {
					d.Set("pipeline_id", *id.pipelineId)
				}
				d.SetId(id.id())
				return []*schema.ResourceData{d}, nil
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

	_, err := clients.PipelinePermissionsClient.UpdatePipelinePermisionsForResource(
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

	if d.IsNewResource() {
		id := pipelineAuthorizationId{
			projectId:  projectId,
			typ:        resType,
			resourceId: d.Get("resource_id").(string),
			pipelineId: nil,
		}
		if v := d.Get("pipeline_project_id").(string); v != "" {
			id.pipelineProjectId = &v
		}
		if v := d.Get("pipeline_id").(int); v != 0 {
			id.pipelineId = &v
		}
		d.SetId(id.id())
	}

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

	if resp == nil || resp.Pipelines == nil || (resp.AllPipelines == nil && len(*resp.Pipelines) == 0) {
		d.SetId("")
		return nil
	}

	var typ, resourceId string
	if resp.Resource != nil {
		if resp.Resource.Type != nil {
			typ = *resp.Resource.Type
		}
		if resp.Resource.Id != nil {
			resourceId = *resp.Resource.Id
		}
	}
	d.Set("type", typ)
	d.Set("resource_id", resourceId)

	if strings.EqualFold(*resp.Resource.Type, "repository") {
		resIds := strings.Split(*resp.Resource.Id, ".")
		if len(resIds) == 2 {
			d.Set("resource_id", resIds[1])
		}
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
