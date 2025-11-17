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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Name of the process",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Description of the process",
			},
			"parent_process_type_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "ID of the parent process",
			},
			"reference_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Reference name of process being created. If not specified, server will assign a unique reference name",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Is the process default?",
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Is the process enabled?",
			},
			"customization_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the type of customization on this process. System Process is default process. Inherited Process is modified process that was System process before",
			},
			"expand": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "none",
				ValidateFunc: validation.StringInSlice([]string{"none", "projects"}, false),
				Description:  "Specifies the expand option when getting the process",
			},
			"projects": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the project",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the project",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the project",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Url of the project",
						},
					},
				},
				Computed:    true,
				Description: "Returns associated projects when using the 'projects' expand option",
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
