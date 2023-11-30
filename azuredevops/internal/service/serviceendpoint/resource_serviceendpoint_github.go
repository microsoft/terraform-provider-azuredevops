package serviceendpoint

import (
	"fmt"
	"strings"
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

const (
	personalAccessTokenGithub = "personal_access_token"
)

// ResourceServiceEndpointGitHub schema and implementation for github service endpoint resource
func ResourceServiceEndpointGitHub() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointGitHubCreate,
		Read:   resourceServiceEndpointGitHubRead,
		Update: resourceServiceEndpointGitHubUpdate,
		Delete: resourceServiceEndpointGitHubDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResourceUUID(),
		Schema:   baseSchema(),
	}

	r.Schema["auth_personal"] = &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"personal_access_token": {
					Type:         schema.TypeString,
					Required:     true,
					DefaultFunc:  schema.EnvDefaultFunc("AZDO_GITHUB_SERVICE_CONNECTION_PAT", nil),
					Description:  "The GitHub personal access token which should be used.",
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
			},
		},
		ConflictsWith: []string{"auth_oauth"},
	}

	r.Schema["auth_oauth"] = &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MinItems: 1,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"oauth_configuration_id": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
		ConflictsWith: []string{"auth_personal"},
	}

	return r
}

func resourceServiceEndpointGitHubCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, _, err := expandServiceEndpointGitHub(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointGitHubRead(d, m)
}

func resourceServiceEndpointGitHubRead(d *schema.ResourceData, m interface{}) error {
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

	flattenServiceEndpointGitHub(d, serviceEndpoint, (*serviceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
	return nil
}

func resourceServiceEndpointGitHubUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectID, err := expandServiceEndpointGitHub(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	updatedServiceEndpoint, err := updateServiceEndpoint(clients, serviceEndpoint)

	if err != nil {
		return fmt.Errorf("Error updating service endpoint in Azure DevOps: %+v", err)
	}

	flattenServiceEndpointGitHub(d, updatedServiceEndpoint, projectID)
	return resourceServiceEndpointGitHubRead(d, m)
}

func resourceServiceEndpointGitHubDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectId, err := expandServiceEndpointGitHub(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, projectId, serviceEndpoint.Id, d.Timeout(schema.TimeoutDelete))
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointGitHub(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	scheme := "InstallationToken"

	parameters := map[string]string{}

	if config, ok := d.GetOk("auth_personal"); ok {
		scheme = "Token"
		parameters = expandAuthPersonalSetGithub(config.(*schema.Set))
	}

	if config, ok := d.GetOk("auth_oauth"); ok {
		scheme = "OAuth"
		parameters = expandAuthOauthSet(config.(*schema.Set))
	}

	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &parameters,
		Scheme:     &scheme,
	}
	serviceEndpoint.Type = converter.String("github")
	serviceEndpoint.Url = converter.String("https://github.com")

	return serviceEndpoint, projectID, nil
}

func expandAuthPersonalSetGithub(d *schema.Set) map[string]string {
	authPerson := make(map[string]string)
	val := d.List()[0].(map[string]interface{}) //auth_personal only have one map configure structure
	authPerson["AccessToken"] = val[personalAccessTokenGithub].(string)
	return authPerson
}

func expandAuthOauthSet(d *schema.Set) map[string]string {
	authConfig := make(map[string]string)
	val := d.List()[0].(map[string]interface{}) //auth_personal only have one map configure structure
	authConfig["ConfigurationId"] = val["oauth_configuration_id"].(string)
	authConfig["AccessToken"] = ""
	return authConfig
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointGitHub(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	if strings.EqualFold(*serviceEndpoint.Authorization.Scheme, "OAuth") {
		d.Set("auth_oauth", &[]map[string]interface{}{
			{
				"oauth_configuration_id": (*serviceEndpoint.Authorization.Parameters)["ConfigurationId"],
			},
		})
	}
	if strings.EqualFold(*serviceEndpoint.Authorization.Scheme, "Token") {
		authPersonalSet := d.Get("auth_personal").(*schema.Set).List()
		authPersonal := flattenAuthPerson(d, authPersonalSet)
		if authPersonal != nil {
			d.Set("auth_personal", authPersonal)
		}
	}
}

func flattenAuthPerson(d *schema.ResourceData, authPersonalSet []interface{}) []interface{} {
	if len(authPersonalSet) == 1 {
		if authPersonal, ok := authPersonalSet[0].(map[string]interface{}); ok {
			return []interface{}{authPersonal}
		}
	}
	return nil
}
