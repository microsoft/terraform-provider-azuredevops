package serviceendpoint

import (
	"fmt"
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointSonarCloud schema and implementation for SonarCloud service endpoint resource
func ResourceServiceEndpointSonarCloud() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointSonarCloudCreate,
		Read:   resourceServiceEndpointSonarCloudRead,
		Update: resourceServiceEndpointSonarCloudUpdate,
		Delete: resourceServiceEndpointSonarCloudDelete,
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
		"token": {
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
			Description:  "Authentication Token generated through SonarCloud (go to My Account > Security > Generate Tokens)",
		},
	})
	return r
}

func resourceServiceEndpointSonarCloudCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointSonarCloud(d)
	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointSonarCloudRead(d, m)
}

func resourceServiceEndpointSonarCloudRead(d *schema.ResourceData, m interface{}) error {
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
		return fmt.Errorf("looking up service endpoint given ID (%s) and project ID (%s): %v", getArgs.EndpointId, *getArgs.Project, err)
	}

	if err = checkServiceConnection(serviceEndpoint); err != nil {
		return err
	}
	flattenServiceEndpointSonarCloud(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointSonarCloudUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointSonarCloud(d)
	if _, err := updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf("Updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointSonarCloudRead(d, m)
}

func resourceServiceEndpointSonarCloudDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint := expandServiceEndpointSonarCloud(d)
	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointSonarCloud(d *schema.ResourceData) *serviceendpoint.ServiceEndpoint {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Scheme: converter.String("Token"),
		Parameters: &map[string]string{
			"apitoken": d.Get("token").(string),
		},
	}
	serviceEndpoint.Type = converter.String("sonarcloud")
	serviceEndpoint.Url = converter.String("https://sonarcloud.io")
	return serviceEndpoint
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointSonarCloud(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)
}
