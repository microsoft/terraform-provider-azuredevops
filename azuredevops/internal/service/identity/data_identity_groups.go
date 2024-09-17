package identity

import (
	"fmt"
	"time"

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
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
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
							Computed: true,
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
		return fmt.Errorf(" failed to get groups for project with ID %s. Error: %v", projectID, err)
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
		return []interface{}{}, nil
	}
	results := make([]interface{}, len(*groups))
	for i, group := range *groups {
		groupMap := make(map[string]interface{})

		if group.Id != nil {
			groupID := *group.Id
			groupMap["id"] = groupID.String()
		} else {
			return nil, fmt.Errorf(" Group Object does not contain an id")
		}
		if group.ProviderDisplayName != nil {
			groupMap["name"] = *group.ProviderDisplayName
		}
		results[i] = groupMap
	}
	return results, nil
}

func getIdentityGroupHash(v interface{}) int {
	return tfhelper.HashString(v.(map[string]interface{})["id"].(string))
}
