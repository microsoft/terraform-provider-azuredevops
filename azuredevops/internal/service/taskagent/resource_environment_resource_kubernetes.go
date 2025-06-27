package taskagent

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceEnvironmentKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnvironmentKubernetesCreate,
		Read:   resourceEnvironmentKubernetesRead,
		Delete: resourceEnvironmentKubernetesDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"environment_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"service_endpoint_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Set: schema.HashString,
			},
		},
	}
}

func resourceEnvironmentKubernetesCreate(d *schema.ResourceData, m interface{}) error {
	project, resource, err := expandEnvironmentKubernetesResource(d)
	if err != nil {
		return fmt.Errorf("expanding the Kubernetes resource from state: %+v", err)
	}

	clients := m.(*client.AggregatedClient)
	createdResource, err := clients.TaskAgentClient.AddKubernetesResourcExistingEndpoint(clients.Ctx, taskagent.AddKubernetesResourceArgsExistingEndpoint{
		CreateParameters: &taskagent.KubernetesResourceCreateParametersExistingEndpoint{
			ClusterName:       resource.ClusterName,
			Name:              resource.Name,
			Namespace:         resource.Namespace,
			Tags:              resource.Tags,
			ServiceEndpointId: resource.ServiceEndpointId,
		},
		Project:       converter.String(project.Id.String()),
		EnvironmentId: resource.EnvironmentReference.Id,
	})
	if err != nil {
		return fmt.Errorf("creating Kubernetes resource in Azure DevOps: %+v", err)
	}

	d.SetId(strconv.Itoa(*createdResource.Id))
	return resourceEnvironmentKubernetesRead(d, m)
}

func resourceEnvironmentKubernetesRead(d *schema.ResourceData, m interface{}) error {
	project, resource, err := expandEnvironmentKubernetesResource(d)
	if err != nil {
		return fmt.Errorf("expanding the Kubernetes resource from state: %+v", err)
	}

	resourceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("getting kubernetes resource id: %+v", err)
	}

	clients := m.(*client.AggregatedClient)
	fetchedResource, err := clients.TaskAgentClient.GetKubernetesResource(clients.Ctx, taskagent.GetKubernetesResourceArgs{
		Project:       converter.String(project.Id.String()),
		EnvironmentId: resource.EnvironmentReference.Id,
		ResourceId:    &resourceId,
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading the Kubernetes resource: %+v", err)
	}

	flattenEnvironmentKubernetesResource(d, project, fetchedResource)
	return nil
}

func resourceEnvironmentKubernetesDelete(d *schema.ResourceData, m interface{}) error {
	project, resource, err := expandEnvironmentKubernetesResource(d)
	if err != nil {
		return fmt.Errorf("expanding the Kubernetes resource from state: %+v", err)
	}

	resourceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("getting kubernetes resource id: %+v", err)
	}

	clients := m.(*client.AggregatedClient)
	err = clients.TaskAgentClient.DeleteKubernetesResource(clients.Ctx, taskagent.DeleteKubernetesResourceArgs{
		Project:       converter.String(project.Id.String()),
		EnvironmentId: resource.EnvironmentReference.Id,
		ResourceId:    &resourceId,
	})

	if err != nil {
		return fmt.Errorf("deleting Kubernetes environment: %+v", err)
	}

	return nil
}

func expandEnvironmentKubernetesResource(d *schema.ResourceData) (*taskagent.ProjectReference, *taskagent.KubernetesResource, error) {
	projectId, err := uuid.Parse(d.Get("project_id").(string))
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to parse project ID to UUID: %s, %+v", d.Get("project_id"), err)
	}
	project := &taskagent.ProjectReference{Id: &projectId}

	serviceEndpointId, err := uuid.Parse(d.Get("service_endpoint_id").(string))
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to parse service endpoint ID to UUID: %s, %+v", d.Get("service_endpoint_id"), err)
	}
	tagsSchemaSet := d.Get("tags").(*schema.Set)
	tags := tfhelper.ExpandStringSet(tagsSchemaSet)

	resource := &taskagent.KubernetesResource{
		EnvironmentReference: &taskagent.EnvironmentReference{
			Id: converter.Int(d.Get("environment_id").(int)),
		},
		Name:              converter.String(d.Get("name").(string)),
		Tags:              &tags,
		ClusterName:       converter.String(d.Get("cluster_name").(string)),
		Namespace:         converter.String(d.Get("namespace").(string)),
		ServiceEndpointId: &serviceEndpointId,
	}

	return project, resource, nil
}

func flattenEnvironmentKubernetesResource(d *schema.ResourceData, project *taskagent.ProjectReference, resource *taskagent.KubernetesResource) {
	d.Set("cluster_name", converter.ToString(resource.ClusterName, ""))
	d.Set("name", *resource.Name)
	d.Set("namespace", *resource.Namespace)
	if resource.Tags != nil {
		tags := *resource.Tags
		ifaceTags := make([]interface{}, len(tags))
		for i, tag := range tags {
			ifaceTags[i] = tag
		}
		d.Set("tags", schema.NewSet(schema.HashString, ifaceTags))
	}
	d.Set("project_id", project.Id.String())
	d.Set("environment_id", *resource.EnvironmentReference.Id)
	d.Set("service_endpoint_id", resource.ServiceEndpointId.String())
}
