package serviceendpoint

import (
	"fmt"
	"strconv"
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

func ResourceServiceEndpointOctopusDeploy() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointOctopusDeployCreate,
		Read:   resourceServiceEndpointOctopusDeployRead,
		Update: resourceServiceEndpointOctopusDeployUpdate,
		Delete: resourceServiceEndpointOctopusDeployDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResourceUUID(),
		Schema:   baseSchema(),
	}
	r.Schema["url"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
	}
	r.Schema["api_key"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}
	r.Schema["ignore_ssl_error"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
	return r
}

func resourceServiceEndpointOctopusDeployCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, _, err := expandServiceEndpointOctopusDeploy(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointOctopusDeployRead(d, m)
}

func resourceServiceEndpointOctopusDeployRead(d *schema.ResourceData, m interface{}) error {
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

	flattenServiceEndpointOctopusDeploy(d, serviceEndpoint, (*serviceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
	return nil
}

func resourceServiceEndpointOctopusDeployUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectID, err := expandServiceEndpointOctopusDeploy(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	updatedServiceEndpoint, err := updateServiceEndpoint(clients, serviceEndpoint)

	if err != nil {
		return fmt.Errorf("Error updating service endpoint in Azure DevOps: %+v", err)
	}

	flattenServiceEndpointOctopusDeploy(d, updatedServiceEndpoint, projectID)
	return resourceServiceEndpointOctopusDeployRead(d, m)
}

func resourceServiceEndpointOctopusDeployDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectId, err := expandServiceEndpointOctopusDeploy(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, projectId, serviceEndpoint.Id, d.Timeout(schema.TimeoutDelete))
}

func expandServiceEndpointOctopusDeploy(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": d.Get("api_key").(string),
		},
		Scheme: converter.String("Token"),
	}

	serviceEndpoint.Data = &map[string]string{
		"ignoreSslErrors": strconv.FormatBool(d.Get("ignore_ssl_error").(bool)),
	}
	serviceEndpoint.Type = converter.String("OctopusEndpoint")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))
	return serviceEndpoint, projectID, nil
}

func flattenServiceEndpointOctopusDeploy(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("url", *serviceEndpoint.Url)

	if v, ok := (*serviceEndpoint.Data)["ignoreSslErrors"]; ok && v != "" {
		ignoreSslErrors, err := strconv.ParseBool(v)
		if err != nil {
			panic(fmt.Errorf(" Failed to parse OctopusDeploy.ignore_ssl_error.(Project: %s), (service endpoint:%s) ,Error: %+v", *serviceEndpoint.Name, projectID, err))
		}
		d.Set("ignore_ssl_error", ignoreSslErrors)
	}
}
