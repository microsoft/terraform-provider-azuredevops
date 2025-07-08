package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

func DataDescriptor() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDescriptorRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"storage_key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"descriptor": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDescriptorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	storageKey := d.Get("storage_key").(string)
	storageKeyUUId, err := uuid.Parse(storageKey)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Invalid storage key: %s. Error: %+v", storageKey, err))
	}

	descriptor, err := clients.GraphClient.GetDescriptor(clients.Ctx, graph.GetDescriptorArgs{StorageKey: &storageKeyUUId})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return diag.Errorf(" The specified storage key %s does not exist.", storageKey)
		}
		return diag.FromErr(fmt.Errorf("Reading storage key: %s. Error: %+v", storageKey, err))
	}

	d.SetId(storageKey)
	d.Set("descriptor", descriptor.Value)
	return nil
}
