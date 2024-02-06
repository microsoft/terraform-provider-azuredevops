package identity

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// DataUser schema and implementation for user data source
func DataIdentityUser() *schema.Resource {
	return &schema.Resource{
		Read: dataIdentitySourceUserRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"descriptor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"search_filter": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Default:      "General",
			},
		},
	}
}

func dataIdentitySourceUserRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	userName, search_filter := d.Get("name").(string), d.Get("search_filter").(string)

	//https://ado.url.com/_apis/identities?search_filter=General&filterValue=my_user&api-version=7.0
	// Query ADO for list of identity users with filter
	FilterUsers, err := getIdentityUsersWithFilterValue(clients, search_filter, userName)
	if err != nil {
		errMsg := "Error finding user"
		if search_filter != "" {
			errMsg = fmt.Sprintf("%s with filter %s", errMsg, search_filter)
		}
		return fmt.Errorf("%s. Error: %v", errMsg, err)
	}
	// with FilterUsers resultes match and filter for the desired user.
	targetUser := selectIdentityUser(FilterUsers, userName)
	if targetUser == nil {
		errMsg := fmt.Sprintf("Could not find user with name %s", userName)
		if search_filter != "" {
			errMsg = fmt.Sprintf("%s with filter %s", errMsg, search_filter)
		}
		return fmt.Errorf(errMsg)
	}

	// Set id and group list for groups data resource
	d.SetId(*targetUser.Descriptor)
	d.Set("descriptor", targetUser.Descriptor)
	return nil
}

func getIdentityUsersWithFilterValue(clients *client.AggregatedClient, searchfilter string, filtervalue string) (*[]identity.Identity, error) {
	// get list of users with search_filter and filter value provided at data source invokation.
	response, err := clients.IdentityClient.ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
		SearchFilter: &searchfilter, // Filter to get users
		FilterValue:  &filtervalue,  // Search String for user
	})
	results := make([]interface{}, len(*response))

	// for each user that has been provided by response, flatten to only displayName, Descriptor/id and search_filter
	for i, user := range *response {
		s := make(map[string]interface{})

		if user.Id != nil {
			s["descriptor"] = *user.Id
		} else {
			return nil, fmt.Errorf("users Object does not contain a descriptor")
		}
		if user.ProviderDisplayName != nil {
			s["name"] = *user.ProviderDisplayName
		}
		if &filtervalue != nil {
			s["search_filter"] = filtervalue
		}
		results[i] = s
	}
	if err != nil {
		return nil, err
	}

	return response, nil
}

// for list of users filter and return user based on ProviderDisplayName
func selectIdentityUser(users *[]identity.Identity, userName string) *identity.Identity {
	for _, user := range *users {
		if strings.EqualFold(*user.ProviderDisplayName, userName) {
			return &user
		}
	}
	return nil
}
