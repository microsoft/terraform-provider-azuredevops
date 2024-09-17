package serviceendpoint

import (
	"fmt"
	"maps"
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

// ResourceServiceEndpointSonarQube schema and implementation for SonarQube service endpoint resource
func ResourceServiceEndpointSonarQube() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointSonarQubeCreate,
		Read:   resourceServiceEndpointSonarQubeRead,
		Update: resourceServiceEndpointSonarQubeUpdate,
		Delete: resourceServiceEndpointSonarQubeDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResourceUUID(),
		Schema:   baseSchema(),
	}

	maps.Copy(r.Schema, map[string]*schema.Schema{
		"url": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			Description:  "Url for the SonarQube Server",
		},

		"token": {
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
			Description:  "Authentication Token generated through SonarQube (go to My Account > Security > Generate Tokens)",
		},
	})

	return r
}

func resourceServiceEndpointSonarQubeCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, _, err := expandServiceEndpointSonarQube(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointSonarQubeRead(d, m)
}

func resourceServiceEndpointSonarQubeRead(d *schema.ResourceData, m interface{}) error {
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

	flattenServiceEndpointSonarQube(d, serviceEndpoint, (*serviceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	return nil
}

func resourceServiceEndpointSonarQubeUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectID, err := expandServiceEndpointSonarQube(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	updatedServiceEndpoint, err := updateServiceEndpoint(clients, serviceEndpoint)

	if err != nil {
		return fmt.Errorf("Error updating service endpoint in Azure DevOps: %+v", err)
	}

	flattenServiceEndpointSonarQube(d, updatedServiceEndpoint, projectID.String())
	return resourceServiceEndpointSonarQubeRead(d, m)
}

func resourceServiceEndpointSonarQubeDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectId, err := expandServiceEndpointSonarQube(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, projectId, serviceEndpoint.Id, d.Timeout(schema.TimeoutDelete))
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointSonarQube(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Scheme: converter.String("UsernamePassword"),
		Parameters: &map[string]string{
			"username": d.Get("token").(string),
		},
	}
	serviceEndpoint.Type = converter.String("sonarqube")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointSonarQube(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID string) {
	doBaseFlattening(d, serviceEndpoint, projectID)

	d.Set("url", *serviceEndpoint.Url)
}
