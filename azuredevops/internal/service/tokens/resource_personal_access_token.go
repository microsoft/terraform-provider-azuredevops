package tokens

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	tokens "github.com/microsoft/azure-devops-go-api/azuredevops/v7/tokens"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
)

// ResourcePersonalAccessToken schema and implementation for Personal Access Tokens resource
func ResourcePersonalAccessToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzurePersonalAccessTokenCreate,
		Read:   resourceAzurePersonalAccessTokenRead,
		Update: resourceAzurePersonalAccessTokenUpdate,
		Delete: resourceAzurePersonalAccessTokenRevoke,
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
			},
			"all_orgs": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"authorization_id": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"scope": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: suppress.CaseDifference,
				},
			},
			"token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"target_accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"valid_from": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_to": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},
		},
	}
}

func resourceAzurePersonalAccessTokenCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	token_valid_to := time.Now().UTC().Add(time.Hour * 24 * 30)
	valid_to := d.Get("valid_to").(string)
	if valid_to != "" {
		time, err := time.Parse(time.RFC3339, valid_to)
		if err != nil {
			return fmt.Errorf(" parsing valid to date: %+v", err)
		}
		token_valid_to = time
	}

	scopes := ""
	if d.Get("scope").([]interface{}) != nil {
		for _, scope := range d.Get("scope").([]interface{}) {
			scopes += scope.(string) + " "
		}
		scopes = strings.TrimRight(scopes, " ")
	}

	create_token := tokens.PatTokenCreateRequest{}
	create_token.AllOrgs = converter.Bool(d.Get("all_orgs").(bool))
	create_token.DisplayName = converter.String(d.Get("name").(string))
	create_token.Scope = converter.String(scopes)
	create_token.ValidTo = &azuredevops.Time{Time: token_valid_to}

	args := tokens.CreatePatArgs{Token: &create_token}

	token, err := clients.TokensClient.CreatePat(clients.Ctx, args)
	if err != nil {
		return fmt.Errorf(" creating pat token in Azure DevOps: %+v", err)
	}

	d.SetId((*token.PatToken.AuthorizationId).String())
	return resourceAzurePersonalAccessTokenRead(d, m)
}

func resourceAzurePersonalAccessTokenRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	authorizationID, err := uuid.Parse(d.Id())
	if err != nil || authorizationID == uuid.Nil {
		return fmt.Errorf(" parse token authorization ID: %+v", err)
	}

	token, err := clients.TokensClient.GetPat(clients.Ctx, tokens.GetPatArgs{AuthorizationId: &authorizationID})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(" looking up Personal Access Token with ID %v. Error: %v", authorizationID, err)
	}

	d.SetId((*token.PatToken.AuthorizationId).String())
	d.Set("authorization_id", token.PatToken.AuthorizationId.String())
	d.Set("target_accounts", token.PatToken.TargetAccounts)
	d.Set("scope", strings.Split(*token.PatToken.Scope, " "))
	d.Set("token", *token.PatToken.Token)
	d.Set("valid_to", token.PatToken.ValidTo.String())
	d.Set("valid_from", token.PatToken.ValidFrom.String())
	return nil
}

func resourceAzurePersonalAccessTokenUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	authorizationID, err := uuid.Parse(d.Id())
	if err != nil || authorizationID == uuid.Nil {
		return fmt.Errorf(" parse token authorization ID: %+v", err)
	}

	token_valid_to := time.Now().UTC().Add(time.Hour * 24 * 30)
	valid_to := d.Get("valid_to").(string)
	if valid_to != "" {
		time, err := time.Parse(time.RFC3339, valid_to)
		if err != nil {
			return fmt.Errorf(" parsing valid to date: %+v", err)
		}
		token_valid_to = time
	}

	scopes := ""
	if d.Get("scope").([]interface{}) != nil {
		for _, scope := range d.Get("scope").([]interface{}) {
			scopes += scope.(string) + " "
		}
		scopes = strings.TrimRight(scopes, " ")
	}

	parameter := tokens.UpdatePatArgs{
		Token: &tokens.PatTokenUpdateRequest{
			AllOrgs:         converter.Bool(d.Get("all_orgs").(bool)),
			AuthorizationId: &authorizationID,
			DisplayName:     converter.String(d.Get("name").(string)),
			Scope:           &scopes,
			ValidTo:         &azuredevops.Time{Time: token_valid_to},
		},
	}

	if _, err = clients.TokensClient.UpdatePat(clients.Ctx, parameter); err != nil {
		return fmt.Errorf(" updating Personal Access Token in Azure DevOps: %+v", err)
	}
	return resourceAzurePersonalAccessTokenRead(d, m)
}

func resourceAzurePersonalAccessTokenRevoke(d *schema.ResourceData, m interface{}) error {
	authorizationID, err := uuid.Parse(d.Id())
	if err != nil || authorizationID == uuid.Nil {
		return fmt.Errorf(" parse token authorization ID: %+v", err)
	}

	clients := m.(*client.AggregatedClient)
	if err := clients.TokensClient.Revoke(clients.Ctx, tokens.RevokeArgs{AuthorizationId: &authorizationID}); err != nil {
		return fmt.Errorf(" revoke Personal Access Token in Azure DevOps: %+v", err)
	}
	return nil
}
