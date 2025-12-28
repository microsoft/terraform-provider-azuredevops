package security

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/security"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// DataSecurityNamespace schema and implementation for security namespace data source
func DataSecurityNamespace() *schema.Resource {
	return &schema.Resource{
		Read: dataSecurityNamespaceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"namespace_id"},
				Description:   "The name of the security namespace",
			},
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Description:   "The ID of the security namespace",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the security namespace",
			},
			"actions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Available actions (permissions) in this namespace",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the action/permission",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The display name of the action/permission",
						},
						"bit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The bit value for this permission",
						},
					},
				},
			},
		},
	}
}

func dataSecurityNamespaceRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	securityNamespacesOnce.Do(func() {
		securityNamespacesCache, securityNamespacesErr = clients.SecurityClient.QuerySecurityNamespaces(clients.Ctx, security.QuerySecurityNamespacesArgs{})
	})

	if securityNamespacesErr != nil {
		return fmt.Errorf("querying security namespaces: %v", securityNamespacesErr)
	}

	if securityNamespacesCache == nil || len(*securityNamespacesCache) == 0 {
		return fmt.Errorf("no security namespaces found")
	}

	name := d.Get("name").(string)
	namespaceID := d.Get("id").(string)

	if name == "" && namespaceID == "" {
		return fmt.Errorf("either 'name' or 'id' must be specified")
	}

	var foundNamespace *security.SecurityNamespaceDescription
	for _, ns := range *securityNamespacesCache {
		if name != "" && ns.Name != nil && strings.EqualFold(*ns.Name, name) {
			foundNamespace = &ns
			break
		}
		if namespaceID != "" && ns.NamespaceId != nil && strings.EqualFold(ns.NamespaceId.String(), namespaceID) {
			foundNamespace = &ns
			break
		}
	}

	if foundNamespace == nil {
		return fmt.Errorf("security namespace not found with name %s or id %s", name, namespaceID)
	}

	if foundNamespace.NamespaceId != nil {
		d.SetId(foundNamespace.NamespaceId.String())
		if err := d.Set("id", foundNamespace.NamespaceId.String()); err != nil {
			return err
		}
	}
	if foundNamespace.Name != nil {
		if err := d.Set("name", *foundNamespace.Name); err != nil {
			return err
		}
	}
	if foundNamespace.DisplayName != nil {
		if err := d.Set("display_name", *foundNamespace.DisplayName); err != nil {
			return err
		}
	}

	actions := make([]interface{}, 0)
	if foundNamespace.Actions != nil {
		for _, action := range *foundNamespace.Actions {
			actionMap := map[string]interface{}{}
			if action.Name != nil {
				actionMap["name"] = *action.Name
			}
			if action.DisplayName != nil {
				actionMap["display_name"] = *action.DisplayName
			}
			if action.Bit != nil {
				actionMap["bit"] = *action.Bit
			}
			actions = append(actions, actionMap)
		}
	}
	err := d.Set("actions", actions)
	if err != nil {
		return err
	}

	return nil
}
