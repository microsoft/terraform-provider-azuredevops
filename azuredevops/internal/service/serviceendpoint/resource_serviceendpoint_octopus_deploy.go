package serviceendpoint

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceServiceEndpointOctopusDeploy() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointOctopusDeploy, expandServiceEndpointOctopusDeploy)
	r.Schema["url"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
	}
	r.Schema["api_key"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}
	r.Schema["ignore_ssl_error"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
	return r
}

func expandServiceEndpointOctopusDeploy(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Parameters: &map[string]string{
			"apitoken": d.Get("api_key").(string),
		},
		Scheme: converter.String("Token"),
	}

	serviceEndpoint.Data = &map[string]string{
		"ignoreSslErrors": strconv.FormatBool(d.Get("ignore_ssl_error").(bool)),
	}
	serviceEndpoint.Type = converter.String("OctopusEndpoint")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))
	return serviceEndpoint, projectID, nil
}

func flattenServiceEndpointOctopusDeploy(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("url", *serviceEndpoint.Url)

	ignoreSslErrors, err := strconv.ParseBool((*serviceEndpoint.Data)["ignoreSslErrors"])
	if err != nil {
		panic(fmt.Errorf(" Failed to parse OctopusDeploy.ignore_ssl_error.(Project: %s), (service endpoint:%s) ,Error: %+v", *serviceEndpoint.Name, projectID, err))
	}
	d.Set("ignore_ssl_error", ignoreSslErrors)
}
