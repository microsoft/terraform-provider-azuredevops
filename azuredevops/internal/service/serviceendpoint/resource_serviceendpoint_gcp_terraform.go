package serviceendpoint

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointGcp schema and implementation for gcp service endpoint resource
func ResourceServiceEndpointGcp() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointGcpTerraformCreate,
		Read:   resourceServiceEndpointGcpTerraformRead,
		Update: resourceServiceEndpointGcpTerraformUpdate,
		Delete: resourceServiceEndpointGcpTerraformDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResourceUUID(),
		Schema:   baseSchema(),
	}

	r.Schema["private_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Sensitive:   true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_GCP_SERVICE_CONNECTION_PRIVATE_KEY", nil),
		Description: "Private Key for connecting to the endpoint.",
	}
	r.Schema["token_uri"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_GCP_SERVICE_CONNECTION_TOKEN_URI", nil),
		Description: "The token uri field in the JSON key file for creating the JSON Web Token.",
	}
	r.Schema["gcp_project_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_GCP_SERVICE_CONNECTION_GCP_PROJECT_ID", nil),
		Description: "Scope to be provided",
	}
	r.Schema["client_email"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_GCP_SERVICE_CONNECTION_CLIENT_EMAIL", nil),
		Description: "The client email field in the JSON key file for creating the JSON Web Token.",
	}
	r.Schema["scope"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc("AZDO_GCP_SERVICE_CONNECTION_SCOPE", nil),
		Description: "Scope to be provided",
	}
	return r
}

func resourceServiceEndpointGcpTerraformCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, _, err := expandServiceEndpointGcp(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint111(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointGcpTerraformRead(d, m)
}

func resourceServiceEndpointGcpTerraformRead(d *schema.ResourceData, m interface{}) error {
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

	flattenServiceEndpointGcp(d, serviceEndpoint, (*serviceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
	return nil

}

func resourceServiceEndpointGcpTerraformUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectID, err := expandServiceEndpointGcp(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	updatedServiceEndpoint, err := updateServiceEndpoint(clients, serviceEndpoint)

	if err != nil {
		return fmt.Errorf("Error updating service endpoint in Azure DevOps: %+v", err)
	}

	flattenServiceEndpointGcp(d, updatedServiceEndpoint, projectID)
	return resourceServiceEndpointGcpTerraformRead(d, m)

}

func resourceServiceEndpointGcpTerraformDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectId, err := expandServiceEndpointGcp(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, projectId, serviceEndpoint.Id, d.Timeout(schema.TimeoutDelete))

}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointGcp(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"Issuer":     d.Get("client_email").(string),
			"Audience":   d.Get("token_uri").(string),
			"Scope":      d.Get("scope").(string),
			"PrivateKey": d.Get("private_key").(string),
		},
		Scheme: converter.String("JWT"),
	}
	serviceEndpoint.Data = &map[string]string{
		"project": d.Get("gcp_project_id").(string),
	}
	serviceEndpoint.Type = converter.String("GoogleCloudServiceEndpoint")
	serviceEndpoint.Url = converter.String("https://www.googleapis.com/")
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointGcp(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)

	d.Set("private_key", d.Get("private_key").(string))
	d.Set("client_email", (*serviceEndpoint.Authorization.Parameters)["Issuer"])
	d.Set("token_uri", (*serviceEndpoint.Authorization.Parameters)["Audience"])
	d.Set("scope", (*serviceEndpoint.Authorization.Parameters)["Scope"])
	d.Set("gcp_project_id", (*serviceEndpoint.Data)["project"])
}
