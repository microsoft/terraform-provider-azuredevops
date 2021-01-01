package serviceendpoint

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	}

	secretHashKey1, secretHashSchema1 := tfhelper.GenerateSecreteMemoSchema("client_certificate")
	secretHashKey2, secretHashSchema2 := tfhelper.GenerateSecreteMemoSchema("client_certificate_password")
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
				},
				"server_certificate_thumbprint": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "",
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
				secretHashKey1: secretHashSchema1,
				secretHashKey2: secretHashSchema2,
			},
		},
	}

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointServiceFabric(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)

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

	serviceEndpoint.Type = converter.String("servicefabric")
	serviceEndpoint.Url = converter.String(d.Get("cluster_endpoint").(string))
	return serviceEndpoint, projectID, nil
}

func flattenCertificate(serviceEndpoint *serviceendpoint.ServiceEndpoint, hashKey1 string, hashValue1 string, hashKey2 string, hashValue2 string) interface{} {
	return []map[string]interface{}{{
		"server_certificate_lookup":     (*serviceEndpoint.Authorization.Parameters)["certLookup"],
		"server_certificate_thumbprint": (*serviceEndpoint.Authorization.Parameters)["servercertthumbprint"],
		hashKey1:                        hashValue1,
		hashKey2:                        hashValue2,
	}}
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointServiceFabric(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)

	newHash1, hashKey1 := tfhelper.HelpFlattenSecretNested(d, "certificate", d.Get("certificate.0").(map[string]interface{}), "client_certificate")
	newHash2, hashKey2 := tfhelper.HelpFlattenSecretNested(d, "certificate", d.Get("certificate.0").(map[string]interface{}), "client_certificate_password")
	certificate := flattenCertificate(serviceEndpoint, hashKey1, newHash1, hashKey2, newHash2)
	d.Set("certificate", certificate)

	d.Set("authorization_type", (*serviceEndpoint.Authorization.Scheme))
	d.Set("cluster_endpoint", (*serviceEndpoint.Url))
}
