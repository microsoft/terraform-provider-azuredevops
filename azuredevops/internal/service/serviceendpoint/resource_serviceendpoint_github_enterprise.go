package serviceendpoint

import (
	"fmt"
	"maps"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointGitHubEnterprise schema and implementation for github-enterprise service endpoint resource
func ResourceServiceEndpointGitHubEnterprise() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointGitHubEnterpriseCreate,
		Read:   resourceServiceEndpointGitHubEnterpriseRead,
		Update: resourceServiceEndpointGitHubEnterpriseUpdate,
		Delete: resourceServiceEndpointGitHubEnterpriseDelete,
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
		"auth_personal": {
			Type:     schema.TypeSet,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"personal_access_token": {
						Type:         schema.TypeString,
						Required:     true,
						Sensitive:    true,
						DefaultFunc:  schema.EnvDefaultFunc("AZDO_GITHUB_ENTERPRISE_SERVICE_CONNECTION_PAT", nil),
						Description:  "The GitHub personal access token which should be used.",
						ValidateFunc: validation.StringIsNotWhiteSpace,
					},
				},
			},
		},

		"url": {
			Type:         schema.TypeString,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			Required:     true,
		},
	})
	return r
}

func resourceServiceEndpointGitHubEnterpriseCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointGitHubEnterprise(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointGitHubEnterpriseRead(d, m)
}

func resourceServiceEndpointGitHubEnterpriseRead(d *schema.ResourceData, m interface{}) error {
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
	flattenServiceEndpointGitHubEnterprise(d, serviceEndpoint)
	return nil
}

func resourceServiceEndpointGitHubEnterpriseUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointGitHubEnterprise(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	if _, err = updateServiceEndpoint(clients, serviceEndpoint); err != nil {
		return fmt.Errorf(" Updating service endpoint in Azure DevOps: %+v", err)
	}

	return resourceServiceEndpointGitHubEnterpriseRead(d, m)
}

func resourceServiceEndpointGitHubEnterpriseDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, err := expandServiceEndpointGitHubEnterprise(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, serviceEndpoint, d.Timeout(schema.TimeoutDelete))
}

func flattenServiceEndpointGitHubEnterprise(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) {
	doBaseFlattening(d, serviceEndpoint)

	if serviceEndpoint != nil {
		if serviceEndpoint.Authorization != nil && serviceEndpoint.Authorization.Scheme != nil {
			if strings.EqualFold(*serviceEndpoint.Authorization.Scheme, "Token") {
				authPersonalSet := d.Get("auth_personal").(*schema.Set).List()
				authPersonal := flattenAuthPersonGithubEnterprise(d, authPersonalSet)
				if authPersonal != nil {
					d.Set("auth_personal", authPersonal)
				}
			}
		}
		if serviceEndpoint.Url != nil {
			d.Set("url", *serviceEndpoint.Url)
		}
	}
}

func flattenAuthPersonGithubEnterprise(d *schema.ResourceData, authPersonalSet []interface{}) []interface{} {
	if len(authPersonalSet) == 1 {
		if authPersonal, ok := authPersonalSet[0].(map[string]interface{}); ok {
			return []interface{}{authPersonal}
		}
	}
	return nil
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointGitHubEnterprise(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpoint := doBaseExpansion(d)

	serviceEndpoint.Type = converter.String("githubenterprise")

	seUrl := d.Get("url").(string)
	serviceEndpoint.Url = converter.String(seUrl)

	scheme := "InstallationToken"
	parameters := map[string]string{}

	if config, ok := d.GetOk("auth_personal"); ok {
		scheme = "Token"
		parameters = expandAuthPersonalSetGithubEnterprise(config.(*schema.Set))
	}

	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &parameters,
		Scheme:     &scheme,
	}

	return serviceEndpoint, nil
}

func expandAuthPersonalSetGithubEnterprise(d *schema.Set) map[string]string {
	authPerson := make(map[string]string)
	val := d.List()[0].(map[string]interface{}) //auth_personal only have one map configure structure

	authPerson["apitoken"] = val["personal_access_token"].(string)
	return authPerson
}
