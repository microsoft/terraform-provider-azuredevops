package build

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
)

const msgErrorFailedResourceCreate = "error creating authorized resource: %+v"
const msgErrorFailedResourceUpdate = "error updating authorized resource: %+v"
const msgErrorFailedResourceDelete = "error deleting authorized resource: %+v"

// ResourceResourceAuthorization schema and implementation for resource authorization resource
func ResourceResourceAuthorization() *schema.Resource {
	return &schema.Resource{
		Create: resourceResourceAuthorizationCreate,
		Read:   resourceResourceAuthorizationRead,
		Update: resourceResourceAuthorizationUpdate,
		Delete: resourceResourceAuthorizationDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "id of the resource",
				ValidateFunc: validation.NoZeroValues,
			},
			"definition_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "id of the build definition",
				ValidateFunc: validation.NoZeroValues,
			},
			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "endpoint",
				Description:      "type of the resource",
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validation.StringInSlice([]string{"endpoint", "queue", "variablegroup"}, false),
			},
			"authorized": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "indicates whether the resource is authorized for use",
			},
		},
	}
}

func resourceResourceAuthorizationCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	authorizedResource, projectID, definitionID, err := expandAuthorizedResource(d)
	if err != nil {
		return fmt.Errorf(msgErrorFailedResourceCreate, err)
	}

	err = sendAuthorizedResourceToAPI(clients, authorizedResource, projectID, definitionID)
	if err != nil {
		return fmt.Errorf(msgErrorFailedResourceCreate, err)
	}

	return resourceResourceAuthorizationRead(d, m)
}

func resourceResourceAuthorizationRead(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	clients := m.(*client.AggregatedClient)

	authorizedResource, projectID, definitionID, err := expandAuthorizedResource(d)
	if err != nil {
		return err
	}

	if definitionID == 0 {
		if *authorizedResource.Authorized {
			// (attempt) flatten read result from ado
			resourceRefs, err := clients.BuildClient.GetProjectResources(ctx, build.GetProjectResourcesArgs{
				Project: &projectID,
				Type:    authorizedResource.Type,
				Id:      authorizedResource.Id,
			})

			if err != nil {
				return err
			}

			// the authorization does no longer exist
			if len(*resourceRefs) == 0 {
				d.SetId("")
				return nil
			}

			flattenAuthorizedResource(d, &(*resourceRefs)[0], projectID, definitionID)
			return nil
		}
		// flatten structure provided by user-configuration and not read from ado
		flattenAuthorizedResource(d, authorizedResource, projectID, definitionID)
		return nil
	}

	resourceRefs, err := clients.BuildClient.GetDefinitionResources(ctx, build.GetDefinitionResourcesArgs{
		Project:      &projectID,
		DefinitionId: &definitionID,
	})

	if err != nil {
		return err
	}

	for _, resource := range *resourceRefs {
		if resource.Id == authorizedResource.Id {
			flattenAuthorizedResource(d, &resource, projectID, definitionID)
			return nil
		}
	}

	// flatten structure provided by user-configuration and not read from ado
	flattenAuthorizedResource(d, authorizedResource, projectID, definitionID)
	return nil
}

func resourceResourceAuthorizationDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	authorizedResource, projectID, definitionID, err := expandAuthorizedResource(d)
	if err != nil {
		return fmt.Errorf(msgErrorFailedResourceDelete, err)
	}

	// deletion works only by setting authorized to false
	// because the resource to delete might have had this parameter set to true, we overwrite it
	authorizedResource.Authorized = converter.Bool(false)

	err = sendAuthorizedResourceToAPI(clients, authorizedResource, projectID, definitionID)
	if err != nil {
		return fmt.Errorf(msgErrorFailedResourceDelete, err)
	}

	return nil
}

func resourceResourceAuthorizationUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	authorizedResource, projectID, definitionID, err := expandAuthorizedResource(d)
	if err != nil {
		return fmt.Errorf(msgErrorFailedResourceUpdate, err)
	}

	err = sendAuthorizedResourceToAPI(clients, authorizedResource, projectID, definitionID)
	if err != nil {
		return fmt.Errorf(msgErrorFailedResourceUpdate, err)
	}

	return resourceResourceAuthorizationRead(d, m)
}

func flattenAuthorizedResource(d *schema.ResourceData, authorizedResource *build.DefinitionResourceReference, projectID string, definitionID int) {
	d.SetId(*authorizedResource.Id)
	d.Set("resource_id", *authorizedResource.Id)
	d.Set("type", *authorizedResource.Type)
	d.Set("authorized", *authorizedResource.Authorized)
	d.Set("project_id", projectID)
	d.Set("definition_id", definitionID)
}

func expandAuthorizedResource(d *schema.ResourceData) (*build.DefinitionResourceReference, string, int, error) {
	resourceRef := build.DefinitionResourceReference{
		Authorized: converter.Bool(d.Get("authorized").(bool)),
		Id:         converter.String(d.Get("resource_id").(string)),
		Name:       nil,
		Type:       converter.String(d.Get("type").(string)),
	}

	return &resourceRef, d.Get("project_id").(string), d.Get("definition_id").(int), nil
}

func sendAuthorizedResourceToAPI(clients *client.AggregatedClient, resourceRef *build.DefinitionResourceReference, projectID string, definitionID int) error {
	ctx := context.Background()
	if definitionID == 0 {
		_, err := clients.BuildClient.AuthorizeProjectResources(ctx, build.AuthorizeProjectResourcesArgs{
			Resources: &[]build.DefinitionResourceReference{*resourceRef},
			Project:   &projectID,
		})

		return err
	}
	_, err := clients.BuildClient.AuthorizeDefinitionResources(ctx, build.AuthorizeDefinitionResourcesArgs{
		Resources:    &[]build.DefinitionResourceReference{*resourceRef},
		Project:      &projectID,
		DefinitionId: &definitionID,
	})

	return err
}
