package azuredevops

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/suppress"
)

const msgErrorFailedResourceCreate = "error creating authorized resource: %+v"
const msgErrorFailedResourceUpdate = "error updating authorized resource: %+v"
const msgErrorFailedResourceDelete = "error deleting authorized resource: %+v"

func resourceResourceAuthorization() *schema.Resource {
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
				ValidateFunc: validation.IsUUID,
			},
			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "endpoint",
				Description:      "type of the resource",
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validation.StringInSlice([]string{"endpoint"}, false),
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
	clients := m.(*config.AggregatedClient)
	authorizedResource, projectID, err := expandAuthorizedResource(d)
	if err != nil {
		return fmt.Errorf(msgErrorFailedResourceCreate, err)
	}

	err = sendAuthorizedResourceToAPI(clients, authorizedResource, projectID)
	if err != nil {
		return fmt.Errorf(msgErrorFailedResourceCreate, err)
	}

	return resourceResourceAuthorizationRead(d, m)
}

func resourceResourceAuthorizationRead(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	clients := m.(*config.AggregatedClient)

	authorizedResource, projectID, err := expandAuthorizedResource(d)
	if err != nil {
		return err
	}

	if !*authorizedResource.Authorized {
		// flatten structure provided by user-configuration and not read from ado
		flattenAuthorizedResource(d, authorizedResource, projectID)
	} else {
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

		flattenAuthorizedResource(d, &(*resourceRefs)[0], projectID)
		return nil
	}
	return nil
}

func resourceResourceAuthorizationDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	authorizedResource, projectID, err := expandAuthorizedResource(d)
	if err != nil {
		return fmt.Errorf(msgErrorFailedResourceDelete, err)
	}

	// deletion works only by setting authorized to false
	// because the resource to delete might have had this parameter set to true, we overwrite it
	authorizedResource.Authorized = converter.Bool(false)

	err = sendAuthorizedResourceToAPI(clients, authorizedResource, projectID)
	if err != nil {
		return fmt.Errorf(msgErrorFailedResourceDelete, err)
	}

	return err
}

func resourceResourceAuthorizationUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	authorizedResource, projectID, err := expandAuthorizedResource(d)
	if err != nil {
		return fmt.Errorf(msgErrorFailedResourceUpdate, err)
	}

	err = sendAuthorizedResourceToAPI(clients, authorizedResource, projectID)
	if err != nil {
		return fmt.Errorf(msgErrorFailedResourceUpdate, err)
	}

	return resourceResourceAuthorizationRead(d, m)
}

func flattenAuthorizedResource(d *schema.ResourceData, authorizedResource *build.DefinitionResourceReference, projectID string) {
	d.SetId(*authorizedResource.Id)
	d.Set("resource_id", *authorizedResource.Id)
	d.Set("type", *authorizedResource.Type)
	d.Set("authorized", *authorizedResource.Authorized)
	d.Set("project_id", projectID)
}

func expandAuthorizedResource(d *schema.ResourceData) (*build.DefinitionResourceReference, string, error) {
	resourceRef := build.DefinitionResourceReference{
		Authorized: converter.Bool(d.Get("authorized").(bool)),
		Id:         converter.String(d.Get("resource_id").(string)),
		Name:       nil,
		Type:       converter.String(d.Get("type").(string)),
	}

	return &resourceRef, d.Get("project_id").(string), nil
}

func sendAuthorizedResourceToAPI(clients *config.AggregatedClient, resourceRef *build.DefinitionResourceReference, project string) error {
	ctx := context.Background()

	_, err := clients.BuildClient.AuthorizeProjectResources(ctx, build.AuthorizeProjectResourcesArgs{
		Resources: &[]build.DefinitionResourceReference{*resourceRef},
		Project:   &project,
	})

	return err
}
