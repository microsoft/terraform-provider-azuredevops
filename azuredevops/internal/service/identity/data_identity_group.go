package identity

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// DataIdentityGroup returns the schema and implementation for the group data source
func DataIdentityGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIdentityGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIdentityGroupRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	groupName := d.Get("name").(string)
	projectID := d.Get("project_id").(string)

	// Get groups in specified project ID
	projectGroups, err := getIdentityGroupsWithProjectID(clients, projectID)
	if err != nil {
		errMsg := "Error finding groups"
		if projectID != "" {
			errMsg = fmt.Sprintf("%s for project with ID %s", errMsg, projectID)
		}
		return fmt.Errorf("%s. Error: %v", errMsg, err)
	}

	// Select specific group by name/provider name.
	targetGroup := selectIdentityGroup(&projectGroups, groupName)
	if targetGroup == nil {
		errMsg := fmt.Sprintf("Could not find group with name %s", groupName)
		if projectID != "" {
			errMsg = fmt.Sprintf("%s in project with ID %s", errMsg, projectID)
		}
		return fmt.Errorf(errMsg)
	}

	// Set ID and descriptor for group data resource based on targetGroup output.
	targetGroupID := targetGroup.Id.String()
	d.SetId(targetGroupID)
	d.Set("descriptor", targetGroup.Id)
	return nil
}

// Select Group that match name to Provider Display Name
func selectIdentityGroup(groups *[]identity.Identity, groupName string) *identity.Identity {
	for _, group := range *groups {
		if strings.EqualFold(*group.ProviderDisplayName, groupName) {
			return &group
		}
	}
	return nil
}
