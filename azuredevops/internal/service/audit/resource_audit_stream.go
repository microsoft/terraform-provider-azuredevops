package audit

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/audit"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	streamutils "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/audit/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceAuditStream() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuditStreamCreate,
		ReadContext:   resourceAuditStreamRead,
		UpdateContext: resourceAuditStreamUpdate,
		DeleteContext: resourceAuditStreamDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResourceInteger(),

		Schema: streamutils.ResourceAuditStreamSchema(map[string]*schema.Schema{}),
	}
}
func resourceAuditStreamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	stream := streamutils.ExpandAuditStream(d)

	createPayload := audit.CreateStreamArgs{
		Stream: &stream,
	}

	createdStream, err := clients.AuditClient.CreateStream(clients.Ctx, createPayload)

	if err != nil {
		return diag.FromErr(err)
	}

	streamID := strconv.Itoa(*createdStream.Id)
	d.SetId(streamID)

	return resourceAuditStreamRead(ctx, d, m)
}

func resourceAuditStreamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	var diags diag.Diagnostics

	streamID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	stream, err := clients.AuditClient.QueryStreamById(clients.Ctx, audit.QueryStreamByIdArgs{
		StreamId: &streamID,
	})

	if err != nil {
		return diag.FromErr(err)
	}

	if err := streamutils.FlattenAuditStream(d, stream); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceAuditStreamUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	streamID, _ := strconv.Atoi(d.Id())

	stream := streamutils.ExpandAuditStream(d)
	stream.Id = &streamID

	updatePayload := audit.UpdateStreamArgs{
		Stream: &stream,
	}

	_, err := clients.AuditClient.UpdateStream(clients.Ctx, updatePayload)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAuditStreamRead(ctx, d, m)
}

func resourceAuditStreamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	streamID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = clients.AuditClient.DeleteStream(clients.Ctx, audit.DeleteStreamArgs{
		StreamId: &streamID,
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) {

			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
