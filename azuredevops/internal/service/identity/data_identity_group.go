package identity

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// DataGroup schema and implementation for group data source
func DataIdentityGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataIdentitySourceGroupRead,
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
			"descriptor": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataIdentitySourceGroupRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	groupName, projectID := d.Get("name").(string), d.Get("project_id").(string)

	projectGroups, err := getIdentityGroupsWithProjectDescriptor(clients, projectID)
	if err != nil {
		errMsg := "Error finding groups"
		if projectID != "" {
			errMsg = fmt.Sprintf("%s for project with ID %s", errMsg, projectID)
		}
		return fmt.Errorf("%s. Error: %v", errMsg, err)
	}

	targetGroup := selectIdentityGroup(projectGroups, groupName)
	if targetGroup == nil {
		errMsg := fmt.Sprintf("Could not find group with name %s", groupName)
		if projectID != "" {
			errMsg = fmt.Sprintf("%s in project with ID %s", errMsg, projectID)
		}
		return fmt.Errorf(errMsg)
	}

	d.SetId(*targetGroup.Descriptor)
	d.Set("descriptor", targetGroup.Descriptor)
	return nil
}

func getIdentityGroupsWithProjectDescriptor(clients *client.AggregatedClient, projectDescriptor string) (*[]identity.Identity, error) {
	response, err := clients.IdentityClient.ListGroups(clients.Ctx, identity.ListGroupsArgs{
		ScopeIds: &projectDescriptor,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func selectIdentityGroup(groups *[]identity.Identity, groupName string) *identity.Identity {
	for _, group := range *groups {
		if strings.EqualFold(*group.ProviderDisplayName, groupName) {
			return &group
		}
	}
	return nil
}
