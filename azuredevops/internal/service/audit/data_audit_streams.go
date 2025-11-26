package audit

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/audit"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/audit/utils"
)

func DataResourceAuditStreams() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataResourceAuditStreamsRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: utils.DataAuditStreamsSchema(map[string]*schema.Schema{}),
	}
}

func dataResourceAuditStreamsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	var diags diag.Diagnostics

	allStreams, err := clients.AuditClient.QueryAllStreams(clients.Ctx, audit.QueryAllStreamsArgs{})
	if err != nil {
		return diag.FromErr(err)
	}

	var streamList []any

	if allStreams != nil {
		for _, stream := range *allStreams {
			streamMap := make(map[string]any)
			if err := utils.FlattenSingleAuditStream(streamMap, &stream); err != nil {
				return diag.FromErr(err)
			}
			streamList = append(streamList, streamMap)
		}
	}

	if err := d.Set("streams", streamList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(acctest.RandomWithPrefix("streams-listing"))

	return diags
}
