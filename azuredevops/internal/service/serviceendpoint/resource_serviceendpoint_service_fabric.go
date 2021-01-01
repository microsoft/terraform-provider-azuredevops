package serviceendpoint

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

const (
	resourceAttrClusterEndpoint = "cluster_endpoint"
	resourceBlockCertificate    = "certificate"
	resourceBlockAAD            = "azure_active_directory"
)

func makeSchemaCertificate(r *schema.Resource) {
	resourceElemSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"server_certificate_lookup": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "",
				ValidateFunc: validation.StringInSlice([]string{"Thumbprint", "CommonName"}, false),
			},
			"server_certificate_thumbprint": {
				Type:        schema.TypeString, // TODO make this a set
				Optional:    true,
				Description: "",
				// Elem: &schema.Schema{
				// 	Type: schema.TypeString,
				// },
				ConflictsWith: []string{"certificate.0.server_certificate_common_name"},
			},
			"server_certificate_common_name": {
				Type:        schema.TypeString, // TODO make this a set
				Optional:    true,
				Description: "",
				// Elem: &schema.Schema{
				// 	Type: schema.TypeString,
				// },
				ConflictsWith: []string{"certificate.0.server_certificate_thumbprint"},
			},
			"client_certificate": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"client_certificate_password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
		},
	}
	makeProtectedSchema(resourceElemSchema, "client_certificate", "AZDO_SERVICEFABRIC_SERVICE_CONNECTION_CLIENTCERTIFICATE", "")
	makeProtectedSchema(resourceElemSchema, "client_certificate_password", "AZDO_SERVICEFABRIC_SERVICE_CONNECTION_CLIENTCERTIFICATEPASSWORD", "")
	r.Schema[resourceBlockCertificate] = &schema.Schema{
		Type:          schema.TypeSet,
		MaxItems:      1,
		Optional:      true,
		Description:   "'Certificate Based'-type of configuration",
		Elem:          resourceElemSchema,
		ConflictsWith: []string{resourceBlockAAD},
	}
}

func makeSchemaAAD(r *schema.Resource) {
	resourceElemSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"server_certificate_lookup": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "",
				ValidateFunc: validation.StringInSlice([]string{"Thumbprint", "CommonName"}, false),
			},
			"server_certificate_thumbprint": {
				Type:        schema.TypeString, // TODO make this a set
				Optional:    true,
				Description: "",
				// Elem: &schema.Schema{
				// 	Type: schema.TypeString,
				// },
				ConflictsWith: []string{"certificate.0.server_certificate_common_name"},
			},
			"server_certificate_common_name": {
				Type:        schema.TypeString, // TODO make this a set
				Optional:    true,
				Description: "",
				// Elem: &schema.Schema{
				// 	Type: schema.TypeString,
				// },
				ConflictsWith: []string{"certificate.0.server_certificate_thumbprint"},
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
		},
	}
	makeProtectedSchema(resourceElemSchema, "password", "AZDO_SERVICEFABRIC_SERVICE_CONNECTION_AADPASSWORD", "")
	r.Schema[resourceBlockAAD] = &schema.Schema{
		Type:          schema.TypeSet,
		MaxItems:      1,
		Optional:      true,
		Description:   "'Azure Active Directory Based'-type of configuration",
		Elem:          resourceElemSchema,
		ConflictsWith: []string{resourceBlockCertificate},
	}
}

// ResourceServiceEndpointServiceFabric schema and implementation for Service Fabric service endpoint resource
func ResourceServiceEndpointServiceFabric() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointServiceFabric, expandServiceEndpointServiceFabric)
	r.Schema[resourceAttrClusterEndpoint] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Client connection endpoint for the cluster. Prefix the value with 'tcp://'",
	}
	r.Schema[resourceAttrAuthType] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Type of credentials to use",
		ValidateFunc: validation.StringInSlice([]string{"Certificate", "AzureActiveDirectory", "gMSA"}, false),
	}
	makeSchemaCertificate(r)

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointServiceFabric(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Type = converter.String("servicefabric")
	serviceEndpoint.Url = converter.String(d.Get(resourceAttrClusterEndpoint).(string))
	switch d.Get(resourceAttrAuthType).(string) {
	case "Certificate":
		configurationRaw, exists := d.GetOk(resourceBlockCertificate)
		if !exists {
			return nil, nil, fmt.Errorf("Certificate authorization type requires a certificate block")
		}
		configuration := configurationRaw.(*schema.Set).List()[0].(map[string]interface{})
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"certLookup":           configuration["server_certificate_lookup"].(string),
				"servercertthumbprint": configuration["server_certificate_thumbprint"].(string),
				"certificate":          configuration["client_certificate"].(string),
				"certificatepassword":  configuration["client_certificate_password"].(string),
			},
			Scheme: converter.String("Certificate"),
		}
	case "AzureActiveDirectory":
		configurationRaw := d.Get(resourceBlockAAD).(*schema.Set).List()
		configuration := configurationRaw[0].(map[string]interface{})
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"certLookup":           configuration["server_certificate_lookup"].(string),
				"servercertthumbprint": configuration["server_certificate_thumbprint"].(string),
				"username":             configuration["username"].(string),
				"password":             configuration["password"].(string),
			},
			Scheme: converter.String("UsernamePassword"),
		}
	}
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointServiceFabric(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set(resourceAttrClusterEndpoint, (*serviceEndpoint.Url))
	switch *serviceEndpoint.Authorization.Scheme {
	case "Certificate":
		certificateSet := d.Get(resourceBlockCertificate).(*schema.Set).List()
		configuration := certificateSet[0].(map[string]interface{})
		newHashCertificate, hashKeyCertificate := tfhelper.HelpFlattenSecretNested(d, resourceBlockCertificate, configuration, "client_certificate")
		newHashCertificatePassword, hashKeyCertificatePassword := tfhelper.HelpFlattenSecretNested(d, resourceBlockCertificate, configuration, "client_certificate_password")
		certificate := map[string]interface{}{
			"server_certificate_lookup":     (*serviceEndpoint.Authorization.Parameters)["certLookup"],
			"server_certificate_thumbprint": (*serviceEndpoint.Authorization.Parameters)["servercertthumbprint"],
			"client_certificate":            configuration["client_certificate"].(string),
			"client_certificate_password":   configuration["client_certificate_password"].(string),
			hashKeyCertificate:              newHashCertificate,
			hashKeyCertificatePassword:      newHashCertificatePassword,
		}
		certificateList := make([]map[string]interface{}, 1)
		certificateList[0] = certificate
		d.Set(resourceBlockCertificate, certificateList)
		d.Set(resourceAttrAuthType, "Certificate")
	case "UsernamePassword":
		aadSet := d.Get(resourceBlockAAD).(*schema.Set).List()
		configuration := aadSet[0].(map[string]interface{})
		newHashPassword, hashKeyPassword := tfhelper.HelpFlattenSecretNested(d, resourceBlockAAD, configuration, "password")
		aad := map[string]interface{}{
			"server_certificate_lookup":     (*serviceEndpoint.Authorization.Parameters)["certLookup"],
			"server_certificate_thumbprint": (*serviceEndpoint.Authorization.Parameters)["servercertthumbprint"],
			"username":                      configuration["username"].(string),
			hashKeyPassword:                 newHashPassword,
		}
		aadList := make([]map[string]interface{}, 1)
		aadList[0] = aad
		d.Set(resourceBlockAAD, aadList)
		d.Set(resourceAttrAuthType, "AzureActiveDirectory")
	}
}
