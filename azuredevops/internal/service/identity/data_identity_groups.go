package identity

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// DataIdentityGroups schema and implementation for group data source
func DataIdentityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIdentityGroupsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"groups": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      getIdentityGroupHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityGroupsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)

	// Get groups in specified project id
	groups, err := getIdentityGroupsWithProjectID(clients, projectID)
	if err != nil {
		errMsg := "Error finding groups"
		if projectID != "" {
			errMsg = fmt.Sprintf("%s for project with ID %s", errMsg, projectID)
		}
		return fmt.Errorf("%s. Error: %v", errMsg, err)
	}

	// With project groups flatten results
	flattenedGroups, err := flattenIdentityGroups(&groups)
	if err != nil {
		return fmt.Errorf("Error flattening groups. Error: %w", err)
	}

	// Set id and group list for groups data resource
	d.SetId("groups-" + uuid.New().String())
	d.Set("groups", flattenedGroups)
	return nil
}

// Get Groups with Scope of Project ID
// Get Groups with Scope of Project ID
func getIdentityGroupsWithProjectID(clients *client.AggregatedClient, projectID string) ([]identity.Identity, error) {
	response, err := clients.IdentityClient.ListGroups(clients.Ctx, identity.ListGroupsArgs{
		ScopeIds: &projectID,
	})
	if err != nil {
		return nil, fmt.Errorf("Error getting groups: %v", err)
	}
	return *response, nil
}

// flatten function
func flattenIdentityGroups(groups *[]identity.Identity) ([]interface{}, error) {
	if groups == nil {
		return nil, fmt.Errorf("Input Groups Parameter is nil")
	}
	results := make([]interface{}, len(*groups))
	for i, group := range *groups {
		groupMap := make(map[string]interface{})

		if group.Descriptor != nil {
			groupMap["id"] = *group.Descriptor
		} else {
			return nil, fmt.Errorf("Group Object does not contain an id")
		}
		if group.ProviderDisplayName != nil {
			groupMap["name"] = *group.ProviderDisplayName
		}
		results[i] = groupMap
	}
	return results, nil
}

func getIdentityGroupHash(v interface{}) int {
	group := v.(map[string]interface{})
	groupID := group["id"].(string)
	return tfhelper.HashString(groupID)
}
