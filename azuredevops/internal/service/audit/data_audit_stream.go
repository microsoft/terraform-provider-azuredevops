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
)

// DataResourceAuditStream returns the audit stream data source
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

	displayName := d.Get("display_name").(string)

	allStreams, err := clients.AuditClient.QueryAllStreams(clients.Ctx, audit.QueryAllStreamsArgs{})
	if err != nil {
		return diag.FromErr(err)
	}

	var foundStream *audit.AuditStream
	if allStreams != nil {
		for _, stream := range *allStreams {
			if *stream.DisplayName == displayName {
				foundStream = &stream
				break
			}
		}
	}

	if foundStream == nil {
		return diag.Errorf("No Audit Stream found with display_name: %s", displayName)
	}

	d.SetId(strconv.Itoa(*foundStream.Id))

	if err := utils.FlattenAuditStream(d, foundStream); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
