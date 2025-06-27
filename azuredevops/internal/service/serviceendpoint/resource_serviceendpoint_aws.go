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
)

// ResourceServiceEndpointAws schema and implementation for aws service endpoint resource
func ResourceServiceEndpointAws() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointAwsCreate,
		Read:   resourceServiceEndpointAwsRead,
		Update: resourceServiceEndpointAwsUpdate,
		Delete: resourceServiceEndpointAwsDelete,
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
		"access_key_id": {
			Type:         schema.TypeString,
			Optional:     true,
			DefaultFunc:  schema.EnvDefaultFunc("AZDO_AWS_SERVICE_CONNECTION_ACCESS_KEY_ID", nil),
			Description:  "The AWS access key ID for signing programmatic requests.",
			RequiredWith: []string{"secret_access_key"},
		},

		"secret_access_key": {
			Type:         schema.TypeString,
			Optional:     true,
			DefaultFunc:  schema.EnvDefaultFunc("AZDO_AWS_SERVICE_CONNECTION_SECRET_ACCESS_KEY", nil),
			Description:  "The AWS secret access key for signing programmatic requests.",
			Sensitive:    true,
			RequiredWith: []string{"access_key_id"},
		},

		"session_token": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_AWS_SERVICE_CONNECTION_SESSION_TOKEN", nil),
			Description: "The AWS session token for signing programmatic requests.",
			Sensitive:   true,
		},
		"role_to_assume": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_AWS_SERVICE_CONNECTION_RTA", nil),
			Description: "The Amazon Resource Name (ARN) of the role to assume.",
		},

		"role_session_name": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_AWS_SERVICE_CONNECTION_RSN", nil),
			Description: "Optional identifier for the assumed role session.",
		},
		"external_id": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_AWS_SERVICE_CONNECTION_EXTERNAL_ID", nil),
			Description: "A unique identifier that is used by third parties when assuming roles in their customers' accounts, aka cross-account role access.",
		},

		"use_oidc": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			DefaultFunc: schema.EnvDefaultFunc("AZDO_AWS_SERVICE_CONNECTION_USE_OIDC", nil),
			Description: "Enable this to attempt getting credentials with OIDC token from Azure Devops.",
		},
	})

	return r
}

func resourceServiceEndpointAwsCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointAws(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointAwsRead(d, m)
}

func resourceServiceEndpointAwsRead(d *schema.ResourceData, m interface{}) error {
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

	return flattenServiceEndpointAws(d, serviceEndpoint)
}

func resourceServiceEndpointAwsUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointAws(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	if _, err = updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf("Updating service endpoint in Azure DevOps: %+v", err)
	}
	return resourceServiceEndpointAwsRead(d, m)
}

func resourceServiceEndpointAwsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointAws(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointAws(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"username":        d.Get("access_key_id").(string),
			"password":        d.Get("secret_access_key").(string),
			"sessionToken":    d.Get("session_token").(string),
			"assumeRoleArn":   d.Get("role_to_assume").(string),
			"roleSessionName": d.Get("role_session_name").(string),
			"externalId":      d.Get("external_id").(string),
			"useOIDC":         strconv.FormatBool(d.Get("use_oidc").(bool)),
		},
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Type = converter.String("aws")
	serviceEndpoint.Url = converter.String("https://aws.amazon.com/")
	return serviceEndpoint, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointAws(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) error {
	doBaseFlattening(d, serviceEndpoint)

	if serviceEndpoint.Authorization != nil && serviceEndpoint.Authorization.Parameters != nil {
		if v, ok := (*serviceEndpoint.Authorization.Parameters)["username"]; ok {
			d.Set("access_key_id", v)
		}

		if v, ok := (*serviceEndpoint.Authorization.Parameters)["assumeRoleArn"]; ok {
			d.Set("role_to_assume", v)
		}

		if v, ok := (*serviceEndpoint.Authorization.Parameters)["roleSessionName"]; ok {
			d.Set("role_session_name", v)
		}

		if v, ok := (*serviceEndpoint.Authorization.Parameters)["externalId"]; ok {
			d.Set("external_id", v)
		}

		if v, ok := (*serviceEndpoint.Authorization.Parameters)["useOIDC"]; ok {
			if v != "" {
				useOIDC, err := strconv.ParseBool(v)
				if err != nil {
					return fmt.Errorf("parse `useOIDC`. Error: %+v", err)
				}
				d.Set("use_oidc", useOIDC)
			}
		}
	}
	return nil
}
