package serviceendpoint

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceServiceEndpointSSH() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointSSH, expandServiceEndpointSSH)
	r.Schema["host"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotEmpty,
		Description:  "The Organization Url.",
	}

	r.Schema["username"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}

	r.Schema["port"] = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntAtLeast(0),
		Default:      22,
	}

	r.Schema["password"] = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		Sensitive:        true,
		DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
		ValidateFunc:     validation.StringIsNotEmpty,
	}

	r.Schema["private_key"] = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		Sensitive:        true,
		DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
		ValidateFunc:     validation.StringIsNotEmpty,
	}

	privateKeyHashKey, privateKeyHashSchema := tfhelper.GenerateSecreteMemoSchema("private_key")
	r.Schema[privateKeyHashKey] = privateKeyHashSchema
	pwdHashKey, pwdHashSchema := tfhelper.GenerateSecreteMemoSchema("password")
	r.Schema[pwdHashKey] = pwdHashSchema
	return r
}

func expandServiceEndpointSSH(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *uuid.UUID, error) {
	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Scheme: converter.String("UsernamePassword"),
	}
	serviceEndpoint.Type = converter.String("ssh")
	parameters := map[string]string{}
	parameters["username"] = d.Get("username").(string)
	if pwd, ok := d.GetOk("password"); ok {
		parameters["password"] = pwd.(string)
	}
	serviceEndpoint.Authorization.Parameters = &parameters

	data := map[string]string{}
	data["Host"] = d.Get("host").(string)
	if port, ok := d.GetOk("port"); ok {
		data["Port"] = strconv.Itoa(port.(int))
	}
	if privateKey, ok := d.GetOk("private_key"); ok {
		data["PrivateKey"] = privateKey.(string)
	}
	serviceEndpoint.Data = &data

	return serviceEndpoint, projectID, nil
}

func flattenServiceEndpointSSH(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *uuid.UUID) {
	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("host", (*serviceEndpoint.Data)["Host"])
	if portStr, ok := (*serviceEndpoint.Data)["Port"]; ok {
		port, _ := strconv.ParseInt(portStr, 10, 64)
		d.Set("port", port)
	}
	d.Set("username", (*serviceEndpoint.Authorization.Parameters)["username"])
	tfhelper.HelpFlattenSecret(d, "private_key")
	tfhelper.HelpFlattenSecret(d, "password")
}
