package build

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/build"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
)

const msgErrorFailedResourceCreate = "error creating authorized resource: %+v"
const msgErrorFailedResourceUpdate = "error updating authorized resource: %+v"
const msgErrorFailedResourceDelete = "error deleting authorized resource: %+v"
const msgErrorAuthorizationNoLongerExists = "[WARN] The authorization with ID '%s' no longer exists. Setting Id to empty \n"

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
				Required:     true,
				Description:  "id of the resource",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"definition_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "id of the build definition",
				ValidateFunc: validation.IntAtLeast(1),
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
	authorizedResource, projectID, definitionID := expandAuthorizedResource(d)

	err := sendAuthorizedResourceToAPI(clients, authorizedResource, projectID, definitionID)
	if err != nil {
		return fmt.Errorf(msgErrorFailedResourceCreate, err)
	}

	return resourceResourceAuthorizationRead(d, m)
}

func resourceResourceAuthorizationRead(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	clients := m.(*client.AggregatedClient)

	authorizedResource, projectID, definitionID := expandAuthorizedResource(d)

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

			if len(*resourceRefs) == 0 {
				log.Printf(msgErrorAuthorizationNoLongerExists, *(authorizedResource.Id))
				d.SetId("")
				return nil
			}

			flattenAuthorizedResource(d, &(*resourceRefs)[0], projectID, definitionID)
		} else {
			// flatten structure provided by user-configuration and not read from ado
			flattenAuthorizedResource(d, authorizedResource, projectID, definitionID)
		}
	} else {
		resourceRefs, err := clients.BuildClient.GetDefinitionResources(ctx, build.GetDefinitionResourcesArgs{
			Project:      &projectID,
			DefinitionId: &definitionID,
		})

		if err != nil {
			if utils.ResponseWasNotFound(err) {
				log.Printf(msgErrorAuthorizationNoLongerExists, *(authorizedResource.Id))
				d.SetId("")
				return nil
			}
			return err
		}

		if len(*resourceRefs) == 0 {
			log.Printf(msgErrorAuthorizationNoLongerExists, *(authorizedResource.Id))
			d.SetId("")
			return nil
		}

		for _, resource := range *resourceRefs {
			if resource.Id == authorizedResource.Id {
				flattenAuthorizedResource(d, &resource, projectID, definitionID)
				return nil
			}
		}

		// flatten structure provided by user-configuration and not read from ado
		flattenAuthorizedResource(d, authorizedResource, projectID, definitionID)
	}

	return nil
}

func resourceResourceAuthorizationDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	authorizedResource, projectID, definitionID := expandAuthorizedResource(d)

	// deletion works only by setting authorized to false
	// because the resource to delete might have had this parameter set to true, we overwrite it
	authorizedResource.Authorized = converter.Bool(false)

	err := sendAuthorizedResourceToAPI(clients, authorizedResource, projectID, definitionID)
	if err != nil {
		return fmt.Errorf(msgErrorFailedResourceDelete, err)
	}

	return nil
}

func resourceResourceAuthorizationUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	authorizedResource, projectID, definitionID := expandAuthorizedResource(d)

	err := sendAuthorizedResourceToAPI(clients, authorizedResource, projectID, definitionID)
	if err != nil {
		return fmt.Errorf(msgErrorFailedResourceUpdate, err)
	}

	return resourceResourceAuthorizationRead(d, m)
}

func flattenAuthorizedResource(d *schema.ResourceData, authorizedResource *build.DefinitionResourceReference, projectID string, definitionID int) {
	d.SetId(*authorizedResource.Id)
	d.Set("resource_id", authorizedResource.Id)
	d.Set("type", authorizedResource.Type)
	d.Set("authorized", authorizedResource.Authorized)
	d.Set("project_id", projectID)
	d.Set("definition_id", definitionID)
}

func expandAuthorizedResource(d *schema.ResourceData) (*build.DefinitionResourceReference, string, int) {
	resourceRef := build.DefinitionResourceReference{
		Authorized: converter.Bool(d.Get("authorized").(bool)),
		Id:         converter.String(d.Get("resource_id").(string)),
		Name:       nil,
		Type:       converter.String(d.Get("type").(string)),
	}

	return &resourceRef, d.Get("project_id").(string), d.Get("definition_id").(int)
}

func sendAuthorizedResourceToAPI(clients *client.AggregatedClient, resourceRef *build.DefinitionResourceReference, projectID string, definitionID int) error {
	ctx := context.Background()
	var err error
	if definitionID == 0 {
		_, err = clients.BuildClient.AuthorizeProjectResources(ctx, build.AuthorizeProjectResourcesArgs{
			Resources: &[]build.DefinitionResourceReference{*resourceRef},
			Project:   &projectID,
		})
	} else {
		_, err = clients.BuildClient.AuthorizeDefinitionResources(ctx, build.AuthorizeDefinitionResourcesArgs{
			Resources:    &[]build.DefinitionResourceReference{*resourceRef},
			Project:      &projectID,
			DefinitionId: &definitionID,
		})
	}

	return err
}
