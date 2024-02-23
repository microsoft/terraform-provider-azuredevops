package identity

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// DataIdentityUserResource returns the user data source resource
func DataIdentityUser() *schema.Resource {
	return &schema.Resource{
		Read: dataIdentitySourceUserRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"descriptor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"search_filter": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "General",
				ValidateFunc: validation.StringInSlice([]string{"AccountName", "DisplayName", "MailAddress", "General"}, false),
			},
		},
	}
}

func dataIdentitySourceUserRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	userName := d.Get("name").(string)
	searchFilter := d.Get("search_filter").(string)

	// Query ADO for list of identity users with filter
	filterUsers, err := getIdentityUsersWithFilterValue(clients, searchFilter, userName)
	if err != nil {
		return fmt.Errorf(" finding user with filter %s. Error: %v", searchFilter, err)
	}

	flattenUsers, err := FlattenIdentityUsers(filterUsers)
	if err != nil {
		return fmt.Errorf("Error flatten users. Error: %v", err)
	}

	// Filter for the desired user in the FilterUsers results
	targetUser := selectIdentityUser(flattenUsers, userName)
	if targetUser == nil {
		return fmt.Errorf(" Could not find user with name %s with filter %s", userName, searchFilter)
	}

	// Set id and user list for users data resource
	targetUserID := targetUser.Id.String()
	d.SetId(targetUserID)
	d.Set("descriptor", targetUser.Descriptor)
	return nil
}

func getIdentityUsersWithFilterValue(clients *client.AggregatedClient, searchFilter string, filterValue string) (*[]identity.Identity, error) {
	// Get list of users with search filter and filter value provided at data source invocation.
	response, err := clients.IdentityClient.ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
		SearchFilter: &searchFilter, // Filter to get users
		FilterValue:  &filterValue,  // Search String for user
	})

	if err != nil {
		return nil, err
	}
	return response, nil
}

func FlattenIdentityUsers(users *[]identity.Identity) (*[]identity.Identity, error) {
	if users == nil {
		return nil, fmt.Errorf("Input Users Parameter is nil")
	}
	results := make([]identity.Identity, len(*users))
	for i, user := range *users {
		if user.Descriptor == nil {
			return nil, fmt.Errorf("User Object does not contain an id")
		}
		newUser := identity.Identity{
			Descriptor:          user.Descriptor,
			Id:                  user.Id,
			ProviderDisplayName: user.ProviderDisplayName,
			// Add other fields here if needed
		}
		results[i] = newUser
	}
	return &results, nil
}

func selectIdentityUser(users *[]identity.Identity, userName string) *identity.Identity {
	for _, user := range *users {
		if strings.EqualFold(*user.ProviderDisplayName, userName) {
			return &user
		}
	}
	return nil
}
