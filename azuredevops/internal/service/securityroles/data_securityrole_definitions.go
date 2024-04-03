package securityroles

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/securityroles"
)

func DataSecurityRoleDefinitions() *schema.Resource {
	return &schema.Resource{
		Read: dataSecurityRoleDefinitiosRead,
		Schema: map[string]*schema.Schema{
			"scope": {
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
				Required:     true,
			},
			"definitions": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      getSRDHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"allow_permissions": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"deny_permissions": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"identifier": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"scope": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSecurityRoleDefinitiosRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	scope := d.Get("scope").(string)

	defs, err := clients.SecurityRolesClient.ListSecurityRoleDefinitions(clients.Ctx, &securityroles.ListSecurityRoleDefinitionsArgs{
		Scope: &scope,
	})
	if err != nil {
		errMsg := "Error finding security role definitions"
		if scope != "" {
			errMsg = fmt.Sprintf("%s for scope %s", errMsg, scope)
		}
		return fmt.Errorf("%s. Error: %v", errMsg, err)
	}

	fdefs, err := flattenSRD(defs)
	if err != nil {
		return fmt.Errorf("flattening security role definitions. Error: %w", err)
	}

	d.SetId("secroledefs-" + uuid.New().String())
	d.Set("definitions", fdefs)

	return nil
}

func getSRDHash(v interface{}) int {
	return tfhelper.HashString(v.(map[string]interface{})["identifier"].(string))
}

func flattenSRD(srds *[]securityroles.SecurityRoleDefinition) ([]interface{}, error) {
	if srds == nil {
		return []interface{}{}, nil
	}

	results := make([]interface{}, len(*srds))
	for i, srd := range *srds {
		s := make(map[string]interface{})

		if srd.Identifier != nil {
			s["identifier"] = *srd.Identifier
		} else {
			return nil, fmt.Errorf("Security Role Definition Object does not contain an identifier")
		}
		if srd.DisplayName != nil {
			s["display_name"] = *srd.DisplayName
		}
		if srd.Name != nil {
			s["name"] = *srd.Name
		}
		if srd.Description != nil {
			s["description"] = *srd.Description
		}
		if srd.Scope != nil {
			s["scope"] = *srd.Scope
		}

		if srd.AllowPermissions != nil {
			s["allow_permissions"] = *srd.AllowPermissions
		}

		if srd.DenyPermissions != nil {
			s["deny_permissions"] = *srd.DenyPermissions
		}

		results[i] = s
	}
	return results, nil
}
