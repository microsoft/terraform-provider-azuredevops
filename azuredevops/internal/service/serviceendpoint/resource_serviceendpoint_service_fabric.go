package serviceendpoint

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceServiceEndpointServiceFabric schema and implementation for ServiceFabric service endpoint resource
func ResourceServiceEndpointServiceFabric() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointServiceFabric, expandServiceEndpointServiceFabric)

	r.Schema["cluster_endpoint"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "",
	}

	r.Schema["authorization_type"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "",
		ValidateFunc: validation.StringInSlice([]string{
			"Certificate",
			"AzureActiveDirectory",
		}, false),
	}

	secretHashKeyClientCertificate, secretHashSchemaClientCertificate := tfhelper.GenerateSecreteMemoSchema("client_certificate")
	secretHashKeyClientCertificatePassword, secretHashSchemaClientCertificatePassword := tfhelper.GenerateSecreteMemoSchema("client_certificate_password")
	r.Schema["certificate"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"server_certificate_lookup": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "",
					ValidateFunc: validation.StringInSlice([]string{
						"Thumbprint",
						"CommonName",
					}, false),
				},
				"server_certificate_thumbprint": {
					Type:          schema.TypeString,
					Optional:      true,
					Description:   "",
					ConflictsWith: []string{"certificate.0.server_certificate_common_name"},
				},
				"server_certificate_common_name": {
					Type:          schema.TypeString,
					Optional:      true,
					Description:   "",
					ConflictsWith: []string{"certificate.0.server_certificate_thumbprint"},
				},
				"client_certificate": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "",
					Sensitive:        true,
					DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
				},
				"client_certificate_password": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "",
					Sensitive:        true,
					DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
				},
				secretHashKeyClientCertificate:         secretHashSchemaClientCertificate,
				secretHashKeyClientCertificatePassword: secretHashSchemaClientCertificatePassword,
			},
		},
		ConflictsWith: []string{"azure_active_directory"},
	}

	secretHashKeyPassword, secretHashSchemaPassword := tfhelper.GenerateSecreteMemoSchema("password")
	r.Schema["azure_active_directory"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"server_certificate_lookup": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "",
					ValidateFunc: validation.StringInSlice([]string{
						"Thumbprint",
						"CommonName",
					}, false),
				},
				"server_certificate_thumbprint": {
					Type:          schema.TypeString,
					Optional:      true,
					Description:   "",
					ConflictsWith: []string{"certificate.0.server_certificate_common_name"},
				},
				"server_certificate_common_name": {
					Type:          schema.TypeString,
					Optional:      true,
					Description:   "",
					ConflictsWith: []string{"certificate.0.server_certificate_thumbprint"},
				},
				"username": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "",
				},
				"password": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "",
					Sensitive:        true,
					DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
				},
				secretHashKeyPassword: secretHashSchemaPassword,
			},
		},
		ConflictsWith: []string{"certificate"},
	}

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointServiceFabric(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)

	switch d.Get("authorization_type").(string) {
	case "Certificate":
		certificate := d.Get("certificate").([]interface{})[0].(map[string]interface{})
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"certLookup":           certificate["server_certificate_lookup"].(string),
				"servercertthumbprint": certificate["server_certificate_thumbprint"].(string),
				"certificate":          certificate["client_certificate"].(string),
				"certificatepassword":  certificate["client_certificate_password"].(string),
			},
			Scheme: converter.String("Certificate"),
		}
	case "AzureActiveDirectory":
		azureActiveDirectory := d.Get("azure_active_directory").([]interface{})[0].(map[string]interface{})
		serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
			Parameters: &map[string]string{
				"certLookup":           azureActiveDirectory["server_certificate_lookup"].(string),
				"servercertthumbprint": azureActiveDirectory["server_certificate_thumbprint"].(string),
				"username":             azureActiveDirectory["username"].(string),
				"password":             azureActiveDirectory["password"].(string),
			},
			Scheme: converter.String("UsernamePassword"),
		}
	}

	serviceEndpoint.Type = converter.String("servicefabric")
	serviceEndpoint.Url = converter.String(d.Get("cluster_endpoint").(string))
	return serviceEndpoint, projectID, nil
}

func flattenCertificate(serviceEndpoint *serviceendpoint.ServiceEndpoint, hashKeyClientCertificate string, hashValueClientCertificate string, hashKeyClientCertificatePassword string, hashValueClientCertificatePassword string) interface{} {
	return []map[string]interface{}{{
		"server_certificate_lookup":      (*serviceEndpoint.Authorization.Parameters)["certLookup"],
		"server_certificate_thumbprint":  (*serviceEndpoint.Authorization.Parameters)["servercertthumbprint"],
		hashKeyClientCertificate:         hashValueClientCertificate,
		hashKeyClientCertificatePassword: hashValueClientCertificatePassword,
	}}
}

func flattenAzureActiveDirectory(serviceEndpoint *serviceendpoint.ServiceEndpoint, hashKeyPassword string, hashValuePassword string) interface{} {
	return []map[string]interface{}{{
		"server_certificate_lookup":     (*serviceEndpoint.Authorization.Parameters)["certLookup"],
		"server_certificate_thumbprint": (*serviceEndpoint.Authorization.Parameters)["servercertthumbprint"],
		"username":                      (*serviceEndpoint.Authorization.Parameters)["username"],
		hashKeyPassword:                 hashValuePassword,
	}}
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointServiceFabric(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)

	switch *serviceEndpoint.Authorization.Scheme {
	case "Certificate":
		newHashClientCertificate, hashKeyClientCertificate := tfhelper.HelpFlattenSecretNested(d, "certificate", d.Get("certificate.0").(map[string]interface{}), "client_certificate")
		newHashClientCertificatePassword, hashKeyClientCertificatePassword := tfhelper.HelpFlattenSecretNested(d, "certificate", d.Get("certificate.0").(map[string]interface{}), "client_certificate_password")
		certificate := flattenCertificate(serviceEndpoint, hashKeyClientCertificate, newHashClientCertificate, hashKeyClientCertificatePassword, newHashClientCertificatePassword)
		d.Set("certificate", certificate)
		d.Set("authorization_type", "Certificate")
	case "UsernamePassword":
		newHashPassword, hashKeyPassword := tfhelper.HelpFlattenSecretNested(d, "azure_active_directory", d.Get("azure_active_directory.0").(map[string]interface{}), "password")
		azureActiveDirectory := flattenAzureActiveDirectory(serviceEndpoint, hashKeyPassword, newHashPassword)
		d.Set("azure_active_directory", azureActiveDirectory)
		d.Set("authorization_type", "AzureActiveDirectory")
	}

	d.Set("cluster_endpoint", (*serviceEndpoint.Url))
}
