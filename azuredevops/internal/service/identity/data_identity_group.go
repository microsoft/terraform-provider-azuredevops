package identity

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// DataIdentityGroup returns the schema and implementation for the group data source
func DataIdentityGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIdentityGroupRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"descriptor": {
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
		return fmt.Errorf(" Failed to get groups for project with ID: %s. Error: %v", projectID, err)
	}

	// Select specific group by name/provider name.
	targetGroup := selectIdentityGroup(&projectGroups, groupName)
	if targetGroup == nil {
		return fmt.Errorf(" Can not find group with name %s in project with ID %s", groupName, projectID)
	}

	// Set ID and descriptor for group data resource based on targetGroup output.
	d.SetId(targetGroup.Id.String())
	d.Set("descriptor", targetGroup.Descriptor)
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
