package taskagent

import (
	"fmt"
	"strconv"

	"github.com/microsoft/terraform-provider-azuredevops/azdosdk/taskagentkubernetesresource"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

const (
	kubeResClusterName = "cluster_name"
	kubeResName        = "name"
	kubeResNamespace   = "namespace"
	kubeResTags        = "tags"

	kubeResProjectId         = "project_id"
	kubeResEnvironmentId     = "environment_id"
	kubeResServiceEndpointId = "service_endpoint_id"
)

// ResourceKubernetesResource schema and implementation for kubernetes resource
func ResourceKubernetesResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesCreate,
		Read:   resourceKubernetesRead,
		Delete: resourceKubernetesDelete,
		Schema: map[string]*schema.Schema{
			kubeResProjectId: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			kubeResEnvironmentId: {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			kubeResServiceEndpointId: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			kubeResName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			kubeResNamespace: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			kubeResClusterName: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
			kubeResTags: {
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

func resourceKubernetesCreate(d *schema.ResourceData, m interface{}) error {
	project, resource, err := expandKubernetesResource(d)
	if err != nil {
		return fmt.Errorf("Error expanding the Kubernetes resource from state: %+v", err)
	}

	clients := m.(*client.AggregatedClient)
	createdResource, err := clients.TaskAgentKubernetesResourceClient.AddKubernetesResource(clients.Ctx, taskagentkubernetesresource.AddKubernetesResourceArgs{
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
		return fmt.Errorf("Error creating Kubernetes resource in Azure DevOps: %+v", err)
	}

	flattenKubernetesResource(d, project, createdResource)
	return nil
}

func resourceKubernetesRead(d *schema.ResourceData, m interface{}) error {
	project, resource, err := expandKubernetesResource(d)
	if err != nil {
		return fmt.Errorf("Error expanding the Kubernetes resource from state: %+v", err)
	}

	clients := m.(*client.AggregatedClient)
	fetchedResource, err := clients.TaskAgentClient.GetKubernetesResource(clients.Ctx, taskagent.GetKubernetesResourceArgs{
		Project:       converter.String(project.Id.String()),
		EnvironmentId: resource.EnvironmentReference.Id,
		ResourceId:    resource.Id,
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading the Kubernetes resource: %+v", err)
	}

	flattenKubernetesResource(d, project, fetchedResource)
	return nil
}

func resourceKubernetesDelete(d *schema.ResourceData, m interface{}) error {
	project, resource, err := expandKubernetesResource(d)
	if err != nil {
		return fmt.Errorf("Error expanding the Kubernetes resource from state: %+v", err)
	}

	clients := m.(*client.AggregatedClient)
	err = clients.TaskAgentClient.DeleteKubernetesResource(clients.Ctx, taskagent.DeleteKubernetesResourceArgs{
		Project:       converter.String(project.Id.String()),
		EnvironmentId: resource.EnvironmentReference.Id,
		ResourceId:    resource.Id,
	})
	if err != nil {
		return fmt.Errorf("Error deleting Kubernetes environment: %+v", err)
	}

	d.SetId("")
	return nil
}

func expandKubernetesResource(d *schema.ResourceData) (*taskagent.ProjectReference, *taskagent.KubernetesResource, error) {
	projectId, err := uuid.Parse(d.Get(kubeResProjectId).(string))
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to parse project ID to UUID: %s, %+v", d.Get(kubeResProjectId), err)
	}
	project := &taskagent.ProjectReference{Id: &projectId}

	serviceEndpointId, err := uuid.Parse(d.Get(kubeResServiceEndpointId).(string))
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to parse service endpoint ID to UUID: %s, %+v", d.Get(kubeResServiceEndpointId), err)
	}
	tagsSchemaSet := d.Get(kubeResTags).(*schema.Set)
	tags := tfhelper.ExpandStringSet(tagsSchemaSet)

	resource := &taskagent.KubernetesResource{
		EnvironmentReference: &taskagent.EnvironmentReference{
			Id: converter.Int(d.Get(kubeResEnvironmentId).(int)),
		},
		Name:              converter.String(d.Get(kubeResName).(string)),
		Tags:              &tags,
		ClusterName:       converter.String(d.Get(kubeResClusterName).(string)),
		Namespace:         converter.String(d.Get(kubeResNamespace).(string)),
		ServiceEndpointId: &serviceEndpointId,
	}

	// Look for the ID. This may not exist if we are within the context of a "create" operation,
	// so it is OK if it is missing.
	if d.Id() != "" {
		resourceId, err := strconv.Atoi(d.Id())
		if err != nil {
			return nil, nil, fmt.Errorf("Error getting kubernetes resource id: %+v", err)
		}
		resource.Id = &resourceId
	}
	return project, resource, nil
}

func flattenKubernetesResource(d *schema.ResourceData, project *taskagent.ProjectReference, resource *taskagent.KubernetesResource) {
	d.SetId(strconv.Itoa(*resource.Id))
	d.Set(kubeResClusterName, converter.ToString(resource.ClusterName, ""))
	d.Set(kubeResName, *resource.Name)
	d.Set(kubeResNamespace, *resource.Namespace)
	if resource.Tags != nil {
		tags := *resource.Tags
		ifaceTags := make([]interface{}, len(tags))
		for i, tag := range tags {
			ifaceTags[i] = tag
		}
		d.Set(kubeResTags, schema.NewSet(schema.HashString, ifaceTags))
	}
	d.Set(kubeResProjectId, project.Id.String())
	d.Set(kubeResEnvironmentId, *resource.EnvironmentReference.Id)
	d.Set(kubeResServiceEndpointId, resource.ServiceEndpointId.String())
}
