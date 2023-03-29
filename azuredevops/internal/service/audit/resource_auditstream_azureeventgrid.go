package audit

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/audit"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceAuditStreamAzureEventGrid schema and implementation for Azure EventGrid audit resource
func ResourceAuditStreamAzureEventGridTopic() *schema.Resource {
	r := genBaseAuditStreamResource(flattenAuditStreamAzureEventGridTopic, expandAuditStreamAzureEventGridTopic)

	r.Schema["topic_url"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		DefaultFunc:  schema.EnvDefaultFunc("AZDO_AUDIT_EVENTGRID_TOPIC_URL", nil),
		ValidateFunc: validation.IsURLWithHTTPS,
		Description:  "Url for the Azure EventGrid topic that will send events to",
	}

	r.Schema["access_key"] = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		Sensitive:        true,
		DefaultFunc:      schema.EnvDefaultFunc("AZDO_AUDIT_EVENTGRID_TOPIC_KEY", nil),
		DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
		ValidateFunc:     validation.StringIsNotWhiteSpace,
		Description:      "The access key for the Azure EventGrid topic",
	}
	// Add a spot in the schema to store the token secretly
	stSecretHashKey, stSecretHashSchema := tfhelper.GenerateSecreteMemoSchema("access_key")
	r.Schema[stSecretHashKey] = stSecretHashSchema

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandAuditStreamAzureEventGridTopic(d *schema.ResourceData) (*audit.AuditStream, *int, *bool) {
	auditStream, daysToBackfill, enabled := doBaseExpansion(d)
	auditStream.ConsumerType = converter.String("AzureEventGrid")
	auditStream.ConsumerInputs = &map[string]string{
		"EventGridTopicHostname":  d.Get("topic_url").(string),
		"EventGridTopicAccessKey": d.Get("access_key").(string),
	}

	return auditStream, daysToBackfill, enabled
}

// Convert AzDO data structure to internal Terraform data structure
func flattenAuditStreamAzureEventGridTopic(d *schema.ResourceData, auditStream *audit.AuditStream, daysToBackfill *int, enabled *bool) {
	doBaseFlattening(d, auditStream, daysToBackfill, enabled)

	tfhelper.HelpFlattenSecret(d, "access_key")

	d.Set("topic_url", (*auditStream.ConsumerInputs)["EventGridTopicHostname"])
	d.Set("access_key", (*auditStream.ConsumerInputs)["EventGridTopicAccessKey"])
}
