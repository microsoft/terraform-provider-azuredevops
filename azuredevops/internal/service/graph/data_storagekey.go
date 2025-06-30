package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func DataStorageKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataStorageKeyRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"descriptor": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"storage_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDataStorageKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	storageKey, err := clients.GraphClient.GetStorageKey(clients.Ctx, graph.GetStorageKeyArgs{
		SubjectDescriptor: converter.String(d.Get("descriptor").(string)),
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return diag.Errorf(" The specified descriptor %s does not exist.", storageKey)
		}
		return diag.FromErr(fmt.Errorf("Reading descriptor: %s. Error: %+v", storageKey, err))
	}

	d.SetId(storageKey.Value.String())
	d.Set("storage_key", storageKey.Value.String())
	return nil
}
