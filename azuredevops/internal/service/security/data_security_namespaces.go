package security

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/security"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// DataSecurityNamespaces schema and implementation for security namespaces data source
func DataSecurityNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSecurityNamespacesRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"namespaces": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      getSecurityNamespaceHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the security namespace",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the security namespace",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The display name of the security namespace",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the security namespace",
						},
						"actions": {
							Type:        schema.TypeSet,
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
									"namespace_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The namespace ID this action belongs to",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSecurityNamespacesRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	namespaces, err := clients.SecurityClient.QuerySecurityNamespaces(clients.Ctx, security.QuerySecurityNamespacesArgs{})
	if err != nil {
		d.SetId("")
		return fmt.Errorf("querying security namespaces: %v", err)
	}

	if namespaces == nil || len(*namespaces) == 0 {
		d.SetId("")
		return fmt.Errorf("no security namespaces found")
	}

	flattenedNamespaces := flattenSecurityNamespaces(namespaces)
	d.SetId("security-namespaces-" + uuid.New().String())
	err = d.Set("namespaces", flattenedNamespaces)
	if err != nil {
		return err
	}
	return nil
}

func getSecurityNamespaceHash(v interface{}) int {
	return tfhelper.HashString(v.(map[string]interface{})["namespace_id"].(string))
}

func flattenSecurityNamespaces(namespaces *[]security.SecurityNamespaceDescription) []interface{} {
	if namespaces == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, ns := range *namespaces {
		namespace := map[string]interface{}{}

		if ns.NamespaceId != nil {
			namespace["namespace_id"] = ns.NamespaceId.String()
		}
		if ns.Name != nil {
			namespace["name"] = *ns.Name
		}
		if ns.DisplayName != nil {
			namespace["display_name"] = *ns.DisplayName
		}
		// Note: SecurityNamespaceDescription does not have a Description field
		namespace["description"] = ""

		// Flatten actions
		actions := make([]interface{}, 0)
		if ns.Actions != nil {
			for _, action := range *ns.Actions {
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
				if ns.NamespaceId != nil {
					actionMap["namespace_id"] = ns.NamespaceId.String()
				}
				actions = append(actions, actionMap)
			}
		}
		namespace["actions"] = actions

		results = append(results, namespace)
	}
	return results
}
