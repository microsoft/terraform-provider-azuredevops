package graph

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func DataUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIdentitySourceUserRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"descriptor": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"mail_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"origin": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"origin_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"principal_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subject_kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataIdentitySourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	descriptor := d.Get("descriptor").(string)

	user, err := clients.GraphClient.GetUser(clients.Ctx, graph.GetUserArgs{
		UserDescriptor: converter.String(descriptor),
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return diag.Errorf(" User does not exist with descriptor: %s", descriptor)
		}
		return diag.FromErr(err)
	}

	if user == nil {
		return diag.Errorf(" User does not exist with descriptor: %s", descriptor)
	}

	d.SetId(*user.Descriptor)
	d.Set("subject_kind", user.SubjectKind)
	d.Set("principal_name", user.PrincipalName)
	d.Set("mail_address", user.MailAddress)
	d.Set("origin", user.Origin)
	d.Set("origin_id", user.OriginId)
	d.Set("display_name", user.DisplayName)
	d.Set("domain", user.Domain)
	return nil
}
