package serviceendpoint

import (
	"fmt"
	"maps"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/validate"
)

// ResourceServiceEndpointJenkins schema and implementation for Jenkins service endpoint resource
func ResourceServiceEndpointJenkins() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointJenkinsCreate,
		Read:   resourceServiceEndpointJenkinsRead,
		Update: resourceServiceEndpointJenkinsUpdate,
		Delete: resourceServiceEndpointJenkinsDelete,
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
			ValidateFunc: validate.Url,
			Description:  "Url for the Jenkins Repository",
		},

		"username": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The Jenkins user name.",
		},

		"password": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "The Jenkins password.",
		},

		"accept_untrusted_certs": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Allows the Jenkins clients to accept self-signed SSL server certificates without installing them into the TFS service role and/or Build Agent computers.",
		},
	})

	return r
}
func resourceServiceEndpointJenkinsCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointJenkins(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointJenkinsRead(d, m)
}

func resourceServiceEndpointJenkinsRead(d *schema.ResourceData, m interface{}) error {
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
	flattenServiceEndpointJenkins(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointJenkinsUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointJenkins(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	if _, err = updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf("Updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointJenkinsRead(d, m)
}

func resourceServiceEndpointJenkinsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointJenkins(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointJenkins(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("Jenkins")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))

	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username": d.Get("username").(string),
			"password": d.Get("password").(string),
		},
		Scheme: converter.String("UsernamePassword"),
	}

	data := map[string]string{}
	data["AcceptUntrustedCerts"] = strconv.FormatBool(d.Get("accept_untrusted_certs").(bool))

	serviceEndpoint.Data = &data

	return serviceEndpoint, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointJenkins(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)
	d.Set("url", *serviceEndpoint.Url)
	d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
	unsecured, err := strconv.ParseBool((*serviceEndpoint.Data)["AcceptUntrustedCerts"])
	if err != nil {
		fmt.Println(err)
		return
	}
	d.Set("accept_untrusted_certs", unsecured)
}
