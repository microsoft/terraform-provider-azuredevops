package serviceendpoint

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceServiceEndpointNuGet() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointNuGetCreate,
		Read:   resourceServiceEndpointNuGetRead,
		Update: resourceServiceEndpointNuGetUpdate,
		Delete: resourceServiceEndpointNuGetDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResourceUUID(),
		Schema:   baseSchema(),
	}

	r.Schema["feed_url"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
	}

	r.Schema["api_key"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Sensitive:     true,
		ValidateFunc:  validation.StringIsNotEmpty,
		ConflictsWith: []string{"personal_access_token", "username", "password"},
		AtLeastOneOf:  []string{"api_key", "personal_access_token", "username", "password"},
	}

	r.Schema["personal_access_token"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Sensitive:     true,
		ValidateFunc:  validation.StringIsNotEmpty,
		ConflictsWith: []string{"api_key", "username", "password"},
	}

	r.Schema["username"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		ValidateFunc:  validation.StringIsNotEmpty,
		ConflictsWith: []string{"personal_access_token", "api_key"},
		RequiredWith:  []string{"password"},
	}

	r.Schema["password"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Sensitive:     true,
		ValidateFunc:  validation.StringIsNotEmpty,
		ConflictsWith: []string{"personal_access_token", "api_key"},
		RequiredWith:  []string{"username"},
	}
	return r
}

func resourceServiceEndpointNuGetCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, _, err := expandServiceEndpointNuGet(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointNuGetRead(d, m)
}

func resourceServiceEndpointNuGetRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	getArgs, err := serviceEndpointGetArgs(d)
	if err != nil {
		return err
	}

	serviceEndpoint, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, *getArgs)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(" looking up service endpoint given ID (%v) and project ID (%v): %v", getArgs.EndpointId, getArgs.Project, err)
	}

	flattenServiceEndpointNuGet(d, serviceEndpoint, (*serviceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
	return nil
}

func resourceServiceEndpointNuGetUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectID, err := expandServiceEndpointNuGet(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	updatedServiceEndpoint, err := updateServiceEndpoint(clients, serviceEndpoint)

	if err != nil {
		return fmt.Errorf("Error updating service endpoint in Azure DevOps: %+v", err)
	}

	flattenServiceEndpointNuGet(d, updatedServiceEndpoint, projectID)
	return resourceServiceEndpointNuGetRead(d, m)
}

func resourceServiceEndpointNuGetDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectId, err := expandServiceEndpointNuGet(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, projectId, serviceEndpoint.Id, d.Timeout(schema.TimeoutDelete))
}

func expandServiceEndpointNuGet(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("externalnugetfeed")
	serviceEndpoint.Url = converter.String(d.Get("feed_url").(string))
	if apiKey := d.Get("api_key"); apiKey != "" {
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"nugetkey": apiKey.(string),
			},
			Scheme: converter.String("None"),
		}
	}

	if pat := d.Get("personal_access_token"); pat != "" {
		serviceEndpoint.Type = converter.String("externalnugetfeed")
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"apitoken": pat.(string),
			},
			Scheme: converter.String("Token"),
		}
	}

	if uname := d.Get("username"); uname != "" {
		serviceEndpoint.Type = converter.String("externalnugetfeed")
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"username": uname.(string),
				"password": d.Get("password").(string),
			},
			Scheme: converter.String("UsernamePassword"),
		}
	}
	return serviceEndpoint, projectID, nil
}

func flattenServiceEndpointNuGet(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("feed_url", *serviceEndpoint.Url)

	switch *serviceEndpoint.Authorization.Scheme {
	case "None":
		d.Set("api_key", d.Get("api_key"))
	case "Token":
		d.Set("personal_access_token", d.Get("personal_access_token"))
	case "UsernamePassword":
		d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
		d.Set("password", d.Get("password"))
		fmt.Printf("")
	}
}
