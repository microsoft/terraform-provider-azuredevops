package tokens

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/tokens"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
)

// DataAgentQueue schema and implementation for agent queue source
func DataPersonalAccessToken() *schema.Resource {
	return &schema.Resource{
		Read: dataPersonalAccessTokenRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"authorization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"target_accounts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"token": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"valid_from": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"valid_to": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},
		},
	}
}

func dataPersonalAccessTokenRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	authorizationID, err := uuid.Parse(d.Get("authorization_id").(string))
	if err != nil {
		return fmt.Errorf(" parse token authorization ID: %+v", err)
	}

	token, err := clients.TokensClient.GetPat(clients.Ctx, tokens.GetPatArgs{AuthorizationId: &authorizationID})
	if err != nil {
		return fmt.Errorf("Error getting personal access token by authorization ID: %v", err)
	}

	if token == nil {
		d.SetId("")
		return nil
	}

	d.SetId(token.PatToken.AuthorizationId.String())
	d.Set("name", token.PatToken.DisplayName)
	d.Set("authorization_id", token.PatToken.AuthorizationId.String())
	d.Set("scope", token.PatToken.Scope)
	d.Set("target_accounts", token.PatToken.TargetAccounts)
	d.Set("token", token.PatToken.Token)
	d.Set("valid_from", token.PatToken.ValidFrom)
	d.Set("valid_to", token.PatToken.ValidTo)
	return nil
}
