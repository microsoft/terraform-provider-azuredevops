package serviceendpoint

import (
	"fmt"
	"strconv"
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
	resourceBlockServiceFabricAzureActiveDirectory = "azure_active_directory"
	resourceBlockServiceFabricCertificate          = "certificate"
	resourceBlockServiceFabricNone                 = "none"
)

// ResourceServiceEndpointServiceFabric schema and implementation for ServiceFabric service endpoint resource
func ResourceServiceEndpointServiceFabric() *schema.Resource {
	r := &schema.Resource{
		Create: resourceServiceEndpointServiceFabricCreate,
		Read:   resourceServiceEndpointServiceFabricRead,
		Update: resourceServiceEndpointServiceFabricUpdate,
		Delete: resourceServiceEndpointServiceFabricDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResourceUUID(),
		Schema:   baseSchema(),
	}

	r.Schema["cluster_endpoint"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Client connection endpoint for the cluster. Prefix the value with 'tcp://';. This value overrides the publish profile.",
	}

	r.Schema[resourceBlockServiceFabricCertificate] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"server_certificate_lookup":      servicefabricServerCertificateLookupSchema(),
				"server_certificate_thumbprint":  servicefabricServerCertificateThumbprintSchema(resourceBlockServiceFabricCertificate),
				"server_certificate_common_name": servicefabricServerCertificateCommonNameSchema(resourceBlockServiceFabricCertificate),
				"client_certificate": {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "Base64 encoding of the cluster's client certificate file.",
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"client_certificate_password": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Password for the certificate.",
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
		ConflictsWith: []string{resourceBlockServiceFabricAzureActiveDirectory, resourceBlockServiceFabricNone},
	}

	r.Schema[resourceBlockServiceFabricAzureActiveDirectory] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"server_certificate_lookup":      servicefabricServerCertificateLookupSchema(),
				"server_certificate_thumbprint":  servicefabricServerCertificateThumbprintSchema(resourceBlockServiceFabricAzureActiveDirectory),
				"server_certificate_common_name": servicefabricServerCertificateCommonNameSchema(resourceBlockServiceFabricAzureActiveDirectory),
				"username": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "Specify an Azure Active Directory account.",
				},
				"password": {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "Password for the Azure Active Directory account.",
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
		ConflictsWith: []string{resourceBlockServiceFabricCertificate, resourceBlockServiceFabricNone},
	}

	r.Schema[resourceBlockServiceFabricNone] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"unsecured": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Skip using windows security for authentication.",
				},
				"cluster_spn": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringLenBetween(0, 1024),
					Description:  "Fully qualified domain SPN for gMSA account. This is applicable only if `unsecured` option is disabled.",
				},
			},
		},
		ConflictsWith: []string{resourceBlockServiceFabricCertificate, resourceBlockServiceFabricAzureActiveDirectory},
	}

	return r
}

func resourceServiceEndpointServiceFabricCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, _, err := expandServiceEndpointServiceFabric(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	serviceEndPoint, err := createServiceEndpoint(d, clients, serviceEndpoint)
	if err != nil {
		return err
	}

	d.SetId(serviceEndPoint.Id.String())
	return resourceServiceEndpointServiceFabricRead(d, m)
}

func resourceServiceEndpointServiceFabricRead(d *schema.ResourceData, m interface{}) error {
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

	flattenServiceEndpointServiceFabric(d, serviceEndpoint, (*serviceEndpoint.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
	return nil
}

func resourceServiceEndpointServiceFabricUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectID, err := expandServiceEndpointServiceFabric(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	updatedServiceEndpoint, err := updateServiceEndpoint(clients, serviceEndpoint)

	if err != nil {
		return fmt.Errorf("Error updating service endpoint in Azure DevOps: %+v", err)
	}

	flattenServiceEndpointServiceFabric(d, updatedServiceEndpoint, projectID)
	return resourceServiceEndpointServiceFabricRead(d, m)
}

func resourceServiceEndpointServiceFabricDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	serviceEndpoint, projectId, err := expandServiceEndpointServiceFabric(d)
	if err != nil {
		return fmt.Errorf(errMsgTfConfigRead, err)
	}

	return deleteServiceEndpoint(clients, projectId, serviceEndpoint.Id, d.Timeout(schema.TimeoutDelete))
}

func servicefabricServerCertificateLookupSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ValidateFunc: validation.StringInSlice([]string{
			"Thumbprint",
			"CommonName",
		}, false),
	}
}

func servicefabricServerCertificateThumbprintSchema(blockName string) *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "The thumbprint(s) of the cluster's certificate(s). This is used to verify the identity of the cluster. This value overrides the publish profile. Separate multiple thumbprints with a comma (',')",
		ValidateFunc:  validation.StringIsNotEmpty,
		ConflictsWith: []string{fmt.Sprintf("%s.0.server_certificate_common_name", blockName)},
	}
}

func servicefabricServerCertificateCommonNameSchema(blockName string) *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "The common name(s) of the cluster's certificate(s). This is used to verify the identity of the cluster. This value overrides the publish profile. Separate multiple common names with a comma (',')",
		ValidateFunc:  validation.StringIsNotEmpty,
		ConflictsWith: []string{fmt.Sprintf("%s.0.server_certificate_thumbprint", blockName)},
	}
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointServiceFabric(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("servicefabric")
	serviceEndpoint.Url = converter.String(d.Get("cluster_endpoint").(string))
	certificate, certificateOk := d.GetOk(resourceBlockServiceFabricCertificate)
	if certificateOk {
		configuration := certificate.([]interface{})[0].(map[string]interface{})
		parameters := expandServiceEndpointServiceFabricServerCertificateLookup(configuration)
		parameters["certificate"] = configuration["client_certificate"].(string)
		parameters["certificatepassword"] = configuration["client_certificate_password"].(string)
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &parameters,
			Scheme:     converter.String("Certificate"),
		}
		return serviceEndpoint, projectID, nil
	}

	azureActiveDirectory, azureActiveDirectoryExists := d.GetOk(resourceBlockServiceFabricAzureActiveDirectory)
	if azureActiveDirectoryExists {
		configuration := azureActiveDirectory.([]interface{})[0].(map[string]interface{})
		parameters := expandServiceEndpointServiceFabricServerCertificateLookup(configuration)
		parameters["username"] = configuration["username"].(string)
		parameters["password"] = configuration["password"].(string)
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &parameters,
			Scheme:     converter.String("UsernamePassword"),
		}
		return serviceEndpoint, projectID, nil
	}

	none, noneExists := d.GetOk(resourceBlockServiceFabricNone)
	if noneExists {
		configuration := none.([]interface{})[0].(map[string]interface{})
		parameters := map[string]string{
			"Unsecured":  strconv.FormatBool(configuration["unsecured"].(bool)),
			"ClusterSpn": configuration["cluster_spn"].(string),
		}
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &parameters,
			Scheme:     converter.String("None"),
		}
		return serviceEndpoint, projectID, nil
	}

	return nil, nil, fmt.Errorf("One of %s or %s or %s blocks must be specified", resourceBlockServiceFabricAzureActiveDirectory, resourceBlockServiceFabricCertificate, resourceBlockServiceFabricNone)
}

func expandServiceEndpointServiceFabricServerCertificateLookup(configuration map[string]interface{}) map[string]string {
	certLookup := configuration["server_certificate_lookup"].(string)
	parameters := map[string]string{
		"certLookup": certLookup,
	}
	switch certLookup {
	case "Thumbprint":
		parameters["servercertthumbprint"] = configuration["server_certificate_thumbprint"].(string)
	case "CommonName":
		parameters["servercertcommonname"] = configuration["server_certificate_common_name"].(string)
	}
	return parameters
}

func flattenServiceFabricCertificate(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) interface{} {
	result := flattenServiceEndpointServiceFabricServerCertificateLookup(serviceEndpoint)
	if certificate, ok := d.GetOk(resourceBlockServiceFabricCertificate); ok {
		configuration := certificate.([]interface{})[0].(map[string]interface{})
		if v, ok := configuration["client_certificate"]; ok {
			result[0]["client_certificate"] = v.(string)
		}
		if v, ok := configuration["client_certificate_password"]; ok {
			result[0]["client_certificate_password"] = v.(string)
		}
	}

	return result
}

func flattenServiceFabricAzureActiveDirectory(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint) interface{} {
	result := flattenServiceEndpointServiceFabricServerCertificateLookup(serviceEndpoint)
	result[0]["username"] = (*serviceEndpoint.Authorization.Parameters)["username"]
	if azureActiveDirectory, ok := d.GetOk(resourceBlockServiceFabricAzureActiveDirectory); ok {
		configuration := azureActiveDirectory.([]interface{})[0].(map[string]interface{})
		if v, ok := configuration["password"]; ok {
			result[0]["password"] = v.(string)
		}
	}
	return result
}

func flattenServiceFabricNone(serviceEndpoint *serviceendpoint.ServiceEndpoint) interface{} {
	unsecured, err := strconv.ParseBool((*serviceEndpoint.Authorization.Parameters)["Unsecured"])
	if err != nil {
		return err
	}
	result := []map[string]interface{}{{
		"unsecured":   unsecured,
		"cluster_spn": (*serviceEndpoint.Authorization.Parameters)["ClusterSpn"],
	}}
	return result
}

func flattenServiceEndpointServiceFabricServerCertificateLookup(serviceEndpoint *serviceendpoint.ServiceEndpoint) []map[string]interface{} {
	certLookup := (*serviceEndpoint.Authorization.Parameters)["certLookup"]
	result := []map[string]interface{}{{
		"server_certificate_lookup": certLookup,
	}}
	switch certLookup {
	case "Thumbprint":
		result[0]["server_certificate_thumbprint"] = (*serviceEndpoint.Authorization.Parameters)["servercertthumbprint"]
	case "CommonName":
		result[0]["server_certificate_common_name"] = (*serviceEndpoint.Authorization.Parameters)["servercertcommonname"]
	}
	return result
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointServiceFabric(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)

	switch *serviceEndpoint.Authorization.Scheme {
	case "Certificate":
		certificate := flattenServiceFabricCertificate(d, serviceEndpoint)
		d.Set(resourceBlockServiceFabricCertificate, certificate)
	case "UsernamePassword":
		azureActiveDirectory := flattenServiceFabricAzureActiveDirectory(d, serviceEndpoint)
		d.Set(resourceBlockServiceFabricAzureActiveDirectory, azureActiveDirectory)
	case "None":
		none := flattenServiceFabricNone(serviceEndpoint)
		d.Set(resourceBlockServiceFabricNone, none)
	}

	d.Set("cluster_endpoint", (*serviceEndpoint.Url))
}
