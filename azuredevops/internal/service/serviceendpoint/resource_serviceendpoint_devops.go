package serviceendpoint

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

const (
	azdoPersonalAccessToken = "personal_access_token"
)

// ResourceServiceEndpointAzureDevOps schema and implementation for Azure DevOps service endpoint resource
func ResourceServiceEndpointAzureDevOps() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointAzureDevOps, expandServiceEndpointAzureDevOps)
	r.Schema["organization"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Default:     "",
		Description: "Azure DevOps organization name",
	}
	authPersonal := &schema.Resource{
		Schema: map[string]*schema.Schema{
			personalAccessToken: {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AZDO_PERSONAL_ACCESS_TOKEN", nil),
				Description:  "The Azure DevOps personal access token which should be used.",
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
	patHashKey, patHashSchema := tfhelper.GenerateSecreteMemoSchema(azdoPersonalAccessToken)
	authPersonal.Schema[patHashKey] = patHashSchema
	r.Schema["auth_personal"] = &schema.Schema{
		Type:          schema.TypeSet,
		Optional:      true,
		MinItems:      1,
		MaxItems:      1,
		Elem:          authPersonal,
		ConflictsWith: []string{},
	}

	return r
}

// Convert internal Terraform data structure to an AzDO data structure:
//
// resource "azuredevops_serviceendpoint_devops" "serviceendpoint" {
// 	project_id             = azuredevops_project.project.id
// 	organization					 = "example"
// 	service_endpoint_name  = "azure-devops-service-connection"
// 	auth_personal {
// 		personal_access_token= "test_token"
// 	}
// 	description = "Managed by Terraform"
// }
func expandServiceEndpointAzureDevOps(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	scheme := "InstallationToken"

	parameters := map[string]string{}

	if config, ok := d.GetOk("auth_personal"); ok {
		scheme = "Token"
		parameters = azdoExpandAuthPersonalSet(config.(*schema.Set))
	}

	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &parameters,
		Scheme:     &scheme,
	}

	serviceEndpoint.Type = converter.String("azdoapi")

	org := d.Get("organization").(string)
	serviceUrl := fmt.Sprint("https://dev.azure.com/", org)
	serviceEndpoint.Url = &serviceUrl

	data := map[string]string{}
	releaseUrl := fmt.Sprint("https://vsrm.dev.azure.com/", org)
	data["releaseUrl"] = releaseUrl
	serviceEndpoint.Data = &data
	return serviceEndpoint, projectID, nil
}

func azdoExpandAuthPersonalSet(d *schema.Set) map[string]string {
	authPerson := make(map[string]string)
	val := d.List()[0].(map[string]interface{}) //auth_personal only have one map configure structure
	authPerson["apitoken"] = val[azdoPersonalAccessToken].(string)
	return authPerson
}

// Convert AzDO data structure to internal Terraform data structure
//
// example for $serviceUrl/_apis/serviceendpoint/endpoints
// {
//   "administratorsGroup": null,
//   "authorization": {
//     "scheme": "Token",
//     "parameters": {
//       "apitoken": "PAT"
//     }
//   },
//   "createdBy": null,
//   "data": {
//     "releaseUrl": "https://vsrm.dev.azure.com/example"
//   },
//   "description": "",
//   "groupScopeId": null,
//   "name": "test",
//   "operationStatus": null,
//   "readersGroup": null,
//   "serviceEndpointProjectReferences": [
//     {
//       "description": "",
//       "name": "test",
//       "projectReference": {
//         "id": "c60cfaaf-63ea-45ac-90e5-2d595275ea42",
//         "name": "test-acc-9k3j7jasp7"
//       }
//     }
//   ],
//   "type": "AZDOAPI",
//   "url": "https://dev.azure.com/example",
//   "isShared": false,
//   "owner": "library"
// }
func flattenServiceEndpointAzureDevOps(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	if strings.EqualFold(*serviceEndpoint.Authorization.Scheme, "Token") {
		authPersonalSet := d.Get("auth_personal").(*schema.Set).List()
		authPersonal := azdoFlattenAuthPerson(d, authPersonalSet)
		if authPersonal != nil {
			d.Set("auth_personal", authPersonal)
		}
	}
}

func azdoFlattenAuthPerson(d *schema.ResourceData, authPersonalSet []interface{}) []interface{} {
	if len(authPersonalSet) == 1 {
		if authPersonal, ok := authPersonalSet[0].(map[string]interface{}); ok {
			newHash, hashKey := tfhelper.HelpFlattenSecretNested(d, "auth_personal", authPersonal, azdoPersonalAccessToken)
			authPersonal[hashKey] = newHash
			return []interface{}{authPersonal}
		}
	}
	return nil
}
