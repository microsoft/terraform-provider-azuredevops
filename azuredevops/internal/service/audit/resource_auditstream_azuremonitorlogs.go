package audit

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/audit"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceAuditStreamAzureEventGrid schema and implementation for Azure EventHub audit resource
func ResourceAuditStreamAzureMonitorLogs() *schema.Resource {
	r := genBaseAuditStreamResource(flattenAuditStreamAzureMonitorLogs, expandAuditStreamAzureMonitorLogs)

	r.Schema["workspace_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		DefaultFunc:  schema.EnvDefaultFunc("AZDO_AUDIT_AZURE_MONITORLOGS_WORKSPACE_ID", nil),
		ValidateFunc: validation.StringIsNotWhiteSpace,
		Description:  "Workspace Id for the Azure Monitor Logs instance that will send events to",
	}

	r.Schema["shared_key"] = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		Sensitive:        true,
		DefaultFunc:      schema.EnvDefaultFunc("AZDO_AUDIT_AZURE_MONITORLOGS_SHARED_KEY", nil),
		DiffSuppressFunc: tfhelper.DiffFuncSuppressSecretChanged,
		ValidateFunc:     validation.StringIsNotWhiteSpace,
		Description:      "The shared key for the Azure Monitor Logs instance",
	}
	// Add a spot in the schema to store the token secretly
	stSecretHashKey, stSecretHashSchema := tfhelper.GenerateSecreteMemoSchema("shared_key")
	r.Schema[stSecretHashKey] = stSecretHashSchema

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandAuditStreamAzureMonitorLogs(d *schema.ResourceData) (*audit.AuditStream, *int, error) {
	auditStream, daysToBackfill := doBaseExpansion(d)
	auditStream.ConsumerType = converter.String("AzureMonitorLogs")
	auditStream.ConsumerInputs = &map[string]string{
		"WorkspaceId": d.Get("workspace_id").(string),
		"SharedKey":   d.Get("shared_key").(string),
	}

	return auditStream, daysToBackfill, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenAuditStreamAzureMonitorLogs(d *schema.ResourceData, auditStream *audit.AuditStream, daysToBackfill *int) {
	doBaseFlattening(d, auditStream, daysToBackfill)

	tfhelper.HelpFlattenSecret(d, "shared_key")

	d.Set("workspace_id", (*auditStream.ConsumerInputs)["WorkspaceId"])
	d.Set("shared_key", (*auditStream.ConsumerInputs)["SharedKey"])
}
