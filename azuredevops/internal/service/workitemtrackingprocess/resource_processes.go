package workitemtrackingprocess

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceProcesses() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceProcess,
		ReadContext:   readResourceProcess,
		UpdateContext: updateResourceProcess,
		DeleteContext: deleteResourceProcess,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"type_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "The ID of the process",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Name of the process",
			},
			"parent_process_type_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "ID of the parent process",
			},
			"reference_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
			},
		},
	}
}

func createResourceProcess(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return readResourceProcess(ctx, d, m)
}

func readResourceProcess(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func updateResourceProcess(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return readResourceProcess(ctx, d, m)
}

func deleteResourceProcess(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
