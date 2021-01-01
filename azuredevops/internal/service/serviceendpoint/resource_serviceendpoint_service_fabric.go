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
	resourceBlockServiceFabricAzureActiveDirectory    = "azure_active_directory"
	resourceBlockServiceFabricCertificate             = "certificate"
	resourceAuthTypeServiceFabricCertificate          = "Certificate"
	resourceAuthTypeServiceFabricAzureActiveDirectory = "AzureActiveDirectory"
)

// ResourceServiceEndpointServiceFabric schema and implementation for ServiceFabric service endpoint resource
func ResourceServiceEndpointServiceFabric() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointServiceFabric, expandServiceEndpointServiceFabric)

	r.Schema["cluster_endpoint"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Client connection endpoint for the cluster. Prefix the value with 'tcp://';. This value overrides the publish profile.",
	}

	r.Schema["authorization_type"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ValidateFunc: validation.StringInSlice([]string{
			resourceAuthTypeServiceFabricCertificate,
			resourceAuthTypeServiceFabricAzureActiveDirectory,
		}, false),
	}

	secretHashKeyClientCertificate, secretHashSchemaClientCertificate := tfhelper.GenerateSecreteMemoSchema("client_certificate")
	secretHashKeyClientCertificatePassword, secretHashSchemaClientCertificatePassword := tfhelper.GenerateSecreteMemoSchema("client_certificate_password")
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
					Type:             schema.TypeString,
					Required:         true,
					Description:      "Base64 encoding of the cluster's client certificate file.",
					Sensitive:        true,
					DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
				},
				"client_certificate_password": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "Password for the certificate.",
					Sensitive:        true,
					DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
				},
				secretHashKeyClientCertificate:         secretHashSchemaClientCertificate,
				secretHashKeyClientCertificatePassword: secretHashSchemaClientCertificatePassword,
			},
		},
		ConflictsWith: []string{resourceBlockServiceFabricAzureActiveDirectory},
	}

	secretHashKeyPassword, secretHashSchemaPassword := tfhelper.GenerateSecreteMemoSchema("password")
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
					Type:        schema.TypeString,
					Required:    true,
					Description: "Specify an Azure Active Directory account.",
				},
				"password": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "Password for the Azure Active Directory account.",
					Sensitive:        true,
					DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
				},
				secretHashKeyPassword: secretHashSchemaPassword,
			},
		},
		ConflictsWith: []string{resourceBlockServiceFabricCertificate},
	}

	return r
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
		Description:   "The thumbprint(s) of the cluster's certificate(s). This is used to verify the identity of the cluster. This value overrides the publish profile. Seperate multiple thumbprints with a comma (',')",
		ConflictsWith: []string{fmt.Sprintf("%s.0.server_certificate_common_name", blockName)},
	}
}

func servicefabricServerCertificateCommonNameSchema(blockName string) *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "The common name(s) of the cluster's certificate(s). This is used to verify the identity of the cluster. This value overrides the publish profile. Seperate multiple common names with a comma (',')",
		ConflictsWith: []string{fmt.Sprintf("%s.0.server_certificate_thumbprint", blockName)},
	}
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointServiceFabric(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)

	switch d.Get("authorization_type").(string) {
	case resourceAuthTypeServiceFabricCertificate:
		data, exists := d.GetOk(resourceBlockServiceFabricCertificate)
		if !exists {
			return nil, nil, fmt.Errorf("%s authorization type requires a %s block", resourceAuthTypeServiceFabricCertificate, resourceBlockServiceFabricCertificate)
		}
		configuration := data.([]interface{})[0].(map[string]interface{})
		certLookup := configuration["server_certificate_lookup"].(string)
		parameters := map[string]string{
			"certLookup":          certLookup,
			"certificate":         configuration["client_certificate"].(string),
			"certificatepassword": configuration["client_certificate_password"].(string),
		}
		switch certLookup {
		case "Thumbprint":
			parameters["servercertthumbprint"] = configuration["server_certificate_thumbprint"].(string)
		case "CommonName":
			parameters["servercertcommonname"] = configuration["server_certificate_common_name"].(string)
		}
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &parameters,
			Scheme:     converter.String("Certificate"),
		}
	case resourceAuthTypeServiceFabricAzureActiveDirectory:
		data, exists := d.GetOk(resourceBlockServiceFabricAzureActiveDirectory)
		if !exists {
			return nil, nil, fmt.Errorf("%s authorization type requires a %s block", resourceAuthTypeServiceFabricAzureActiveDirectory, resourceBlockServiceFabricAzureActiveDirectory)
		}
		configuration := data.([]interface{})[0].(map[string]interface{})
		certLookup := configuration["server_certificate_lookup"].(string)
		parameters := map[string]string{
			"certLookup": certLookup,
			"username":   configuration["username"].(string),
			"password":   configuration["password"].(string),
		}
		switch certLookup {
		case "Thumbprint":
			parameters["servercertthumbprint"] = configuration["server_certificate_thumbprint"].(string)
		case "CommonName":
			parameters["servercertcommonname"] = configuration["server_certificate_common_name"].(string)
		}
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &parameters,
			Scheme:     converter.String("UsernamePassword"),
		}
	}

	serviceEndpoint.Type = converter.String("servicefabric")
	serviceEndpoint.Url = converter.String(d.Get("cluster_endpoint").(string))
	return serviceEndpoint, projectID, nil
}

func flattenCertificate(serviceEndpoint *serviceendpoint.ServiceEndpoint, hashKeyClientCertificate string, hashValueClientCertificate string, hashKeyClientCertificatePassword string, hashValueClientCertificatePassword string) interface{} {
	certLookup := (*serviceEndpoint.Authorization.Parameters)["certLookup"]
	result := []map[string]interface{}{{
		"server_certificate_lookup":      certLookup,
		hashKeyClientCertificate:         hashValueClientCertificate,
		hashKeyClientCertificatePassword: hashValueClientCertificatePassword,
	}}
	switch certLookup {
	case "Thumbprint":
		result[0]["server_certificate_thumbprint"] = (*serviceEndpoint.Authorization.Parameters)["servercertthumbprint"]
	case "CommonName":
		result[0]["server_certificate_common_name"] = (*serviceEndpoint.Authorization.Parameters)["servercertcommonname"]
	}
	return result
}

func flattenAzureActiveDirectory(serviceEndpoint *serviceendpoint.ServiceEndpoint, hashKeyPassword string, hashValuePassword string) interface{} {
	certLookup := (*serviceEndpoint.Authorization.Parameters)["certLookup"]
	result := []map[string]interface{}{{
		"server_certificate_lookup": certLookup,
		"username":                  (*serviceEndpoint.Authorization.Parameters)["username"],
		hashKeyPassword:             hashValuePassword,
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
func flattenServiceEndpointServiceFabric(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)

	switch *serviceEndpoint.Authorization.Scheme {
	case "Certificate":
		newHashClientCertificate, hashKeyClientCertificate := tfhelper.HelpFlattenSecretNested(d, resourceBlockServiceFabricCertificate, d.Get("certificate.0").(map[string]interface{}), "client_certificate")
		newHashClientCertificatePassword, hashKeyClientCertificatePassword := tfhelper.HelpFlattenSecretNested(d, "certificate", d.Get("certificate.0").(map[string]interface{}), "client_certificate_password")
		certificate := flattenCertificate(serviceEndpoint, hashKeyClientCertificate, newHashClientCertificate, hashKeyClientCertificatePassword, newHashClientCertificatePassword)
		d.Set(resourceBlockServiceFabricCertificate, certificate)
		d.Set("authorization_type", resourceAuthTypeServiceFabricCertificate)
	case "UsernamePassword":
		newHashPassword, hashKeyPassword := tfhelper.HelpFlattenSecretNested(d, resourceBlockServiceFabricAzureActiveDirectory, d.Get("azure_active_directory.0").(map[string]interface{}), "password")
		azureActiveDirectory := flattenAzureActiveDirectory(serviceEndpoint, hashKeyPassword, newHashPassword)
		d.Set(resourceBlockServiceFabricAzureActiveDirectory, azureActiveDirectory)
		d.Set("authorization_type", resourceAuthTypeServiceFabricAzureActiveDirectory)
	}

	d.Set("cluster_endpoint", (*serviceEndpoint.Url))
}
