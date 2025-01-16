package serviceendpoint

import (
	"fmt"
	"maps"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceServiceEndpointSSH() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointSSHCreate,
		Read:   resourceServiceEndpointSSHRead,
		Update: resourceServiceEndpointSSHUpdate,
		Delete: resourceServiceEndpointSSHDelete,
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
		"host": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The Organization Url.",
		},

		"username": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"port": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(0),
			Default:      22,
		},

		"password": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"private_key": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	})
	return r
}

func resourceServiceEndpointSSHCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointSSH(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointSSHRead(d, m)
}

func resourceServiceEndpointSSHRead(d *schema.ResourceData, m interface{}) error {
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

	if err = checkServiceConnection(serviceEndpoint); err != nil {
		return err
	}
	return flattenServiceEndpointSSH(d, serviceEndpoint)
}

func resourceServiceEndpointSSHUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointSSH(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	if _, err = updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf(" Updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointSSHRead(d, m)
}

func resourceServiceEndpointSSHDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointSSH(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

func expandServiceEndpointSSH(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Type = converter.String("ssh")
	parameters := map[string]string{}
	parameters["username"] = d.Get("username").(string)
	if pwd, ok := d.GetOk("password"); ok {
		parameters["password"] = pwd.(string)
	}
	serviceEndpoint.Authorization.Parameters = &parameters

	data := map[string]string{}
	data["Host"] = d.Get("host").(string)
	if port, ok := d.GetOk("port"); ok {
		data["Port"] = strconv.Itoa(port.(int))
	}
	if privateKey, ok := d.GetOk("private_key"); ok {
		data["PrivateKey"] = privateKey.(string)
	}
	serviceEndpoint.Data = &data

	return serviceEndpoint, nil
}

func flattenServiceEndpointSSH(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) error {
	doBaseFlattening(d, serviceEndpoint)
	d.Set("host", (*serviceEndpoint.Data)["Host"])
	if portStr, ok := (*serviceEndpoint.Data)["Port"]; ok {
		port, err := strconv.ParseInt(portStr, 10, 64)
		if err != nil {
			return fmt.Errorf(" Parse SSH port: %s. Error: %+v ", portStr, err)
		}
		d.Set("port", port)
	}
	d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
	return nil
}
