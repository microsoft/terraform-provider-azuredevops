package securityroles

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/sdk/securityroles"
)

func DataSecurityRoleDefinitions() *schema.Resource {
	return &schema.Resource{
		Read: dataSecurityRoleDefinitionsRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"scope": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
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

func dataSecurityRoleDefinitionsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	scope := d.Get("scope").(string)

	defs, err := clients.SecurityRolesClient.ListSecurityRoleDefinitions(clients.Ctx, &securityroles.ListSecurityRoleDefinitionsArgs{
		Scope: &scope,
	})
	if err != nil {
		d.SetId("")
		return fmt.Errorf("finding security role definitions for scope: %s. Error: %v", scope, err)
	}

	if defs == nil || len(*defs) == 0 {
		d.SetId("")
		return fmt.Errorf("no role definition found at scope: %s", scope)
	}

	fdefs := flattenSRD(defs)
	d.SetId("secroledefs-" + uuid.New().String())
	d.Set("definitions", fdefs)
	return nil
}

func getSRDHash(v interface{}) int {
	return tfhelper.HashString(v.(map[string]interface{})["identifier"].(string))
}

func flattenSRD(srds *[]securityroles.SecurityRoleDefinition) []interface{} {
	if srds == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, srd := range *srds {
		s := map[string]interface{}{}

		if srd.Identifier != nil {
			s["identifier"] = *srd.Identifier
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

		results = append(results, s)
	}
	return results
}
