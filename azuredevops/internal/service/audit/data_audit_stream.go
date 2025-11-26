package audit

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/audit"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/audit/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func DataResourceAuditStream() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataResourceAuditStreamRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: utils.DataAuditStreamSchema(map[string]*schema.Schema{}),
	}
}

func dataResourceAuditStreamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	var diags diag.Diagnostics

	streamID, err := converter.ASCIIToIntPtr(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	listArgs := audit.QueryStreamByIdArgs{
		StreamId: streamID,
	}
	stream, err := clients.AuditClient.QueryStreamById(clients.Ctx, listArgs)
	if err != nil {
		return diag.FromErr(err)
	}

	if stream == nil {
		d.SetId("")
		return diag.Errorf("Geen Audit Stream gevonden met id: %b", streamID)
	}

	d.SetId(strconv.Itoa(*stream.Id))

	if err := utils.FlattenAuditStream(d, stream); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
