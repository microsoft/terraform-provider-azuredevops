package serviceendpoint

import (
	"errors"
	"fmt"
	"maps"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/validate"
)

// ResourceServiceEndpointJFrogPlatformV2 schema and implementation for JFrog Artifactory service endpoint resource
func ResourceServiceEndpointJFrogPlatformV2() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointJFrogPlatformV2Create,
		Read:   resourceServiceEndpointJFrogPlatformV2Read,
		Update: resourceServiceEndpointJFrogPlatformV2Update,
		Delete: resourceServiceEndpointJFrogPlatformV2Delete,
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
			Description:  "Url for the JFrog Artifactory Server",
		},

		"authentication_token": {
			Type:     schema.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"token": {
						Description: "The JFrog Artifactory access token.",
						Type:        schema.TypeString,
						Required:    true,
						Sensitive:   true,
					},
				},
			},
			ExactlyOneOf: []string{"authentication_basic", "authentication_token"},
		},

		"authentication_basic": {
			Type:     schema.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"username": {
						Description: "The JFrog Artifactory user name.",
						Type:        schema.TypeString,
						Required:    true,
						Sensitive:   true,
					},
					"password": {
						Description: "The JFrog Artifactory password.",
						Type:        schema.TypeString,
						Required:    true,
						Sensitive:   true,
					},
				},
			},
		},
	})

	return r
}

func resourceServiceEndpointJFrogPlatformV2Create(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointJFrogPlatformV2(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointJFrogPlatformV2Read(d, m)
}

func resourceServiceEndpointJFrogPlatformV2Read(d *schema.ResourceData, m interface{}) error {
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
	flattenServiceEndpointArtifactoryV2(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointJFrogPlatformV2Update(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointJFrogPlatformV2(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	if _, err = updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf("Updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointJFrogPlatformV2Read(d, m)
}

func resourceServiceEndpointJFrogPlatformV2Delete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointJFrogPlatformV2(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointJFrogPlatformV2(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("jfrogPlatformService")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))
	authScheme := "Token"

	authParams := make(map[string]string)

	if x, ok := d.GetOk("authentication_token"); ok {
		authScheme = "Token"
		msi := x.([]interface{})[0].(map[string]interface{})
		authParams["apitoken"], ok = msi["token"].(string)
		if !ok {
			return nil, errors.New("Unable to read 'token'")
		}
	} else if x, ok := d.GetOk("authentication_basic"); ok {
		authScheme = "UsernamePassword"
		msi := x.([]interface{})[0].(map[string]interface{})
		authParams["username"], ok = msi["username"].(string)
		if !ok {
			return nil, errors.New("Unable to read 'username'")
		}
		authParams["password"], ok = msi["password"].(string)
		if !ok {
			return nil, errors.New("Unable to read 'password'")
		}
	}
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &authParams,
		Scheme:     &authScheme,
	}

	return serviceEndpoint, nil
}
