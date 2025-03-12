package serviceendpoint

import (
	"fmt"
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
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
	maps.Copy(r.Schema, map[string]*schema.Schema{
		"private_key": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_GCP_SERVICE_CONNECTION_PRIVATE_KEY", nil),
			Description: "Private Key for connecting to the endpoint.",
		},

		"token_uri": {
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_GCP_SERVICE_CONNECTION_TOKEN_URI", nil),
			Description: "The token uri field in the JSON key file for creating the JSON Web Token.",
		},
		"gcp_project_id": {
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_GCP_SERVICE_CONNECTION_GCP_PROJECT_ID", nil),
			Description: "Scope to be provided",
		},
		"client_email": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_GCP_SERVICE_CONNECTION_CLIENT_EMAIL", nil),
			Description: "The client email field in the JSON key file for creating the JSON Web Token.",
		},
		"scope": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_GCP_SERVICE_CONNECTION_SCOPE", nil),
			Description: "Scope to be provided",
		},
	})

	return r
}

func resourceServiceEndpointGcpTerraformCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointGcp(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
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
	if isServiceEndpointDeleted(d, err, serviceEndpoint, getArgs) {
		return nil
	}
	if err != nil {
		return fmt.Errorf(" looking up service endpoint given ID (%s) and project ID (%s): %v", getArgs.EndpointId, *getArgs.Project, err)
	}

	if err = checkServiceConnection(serviceEndpoint); err != nil {
		return err
	}
	flattenServiceEndpointGcp(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointGcpTerraformUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointGcp(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	if _, err = updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf(" Updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointGcpTerraformRead(d, m)
}

func resourceServiceEndpointGcpTerraformDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointGcp(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointGcp(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)
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
	return serviceEndpoint, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointGcp(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)

	d.Set("private_key", d.Get("private_key").(string))
	d.Set("client_email", (*serviceEndpoint.Authorization.Parameters)["Issuer"])
	d.Set("token_uri", (*serviceEndpoint.Authorization.Parameters)["Audience"])
	d.Set("scope", (*serviceEndpoint.Authorization.Parameters)["Scope"])
	d.Set("gcp_project_id", (*serviceEndpoint.Data)["project"])
}
