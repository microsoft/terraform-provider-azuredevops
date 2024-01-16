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
	userName, searchFilter := d.Get("name").(string), d.Get("search_filter").(string)

	//https://ado.url.com/_apis/identities?searchFilter=General&filterValue=my_user&api-version=7.0

	FilterUsers, err := getIdentityUsersWithFilterValue(clients, searchFilter, userName)
	if err != nil {
		errMsg := "Error finding user"
		if searchFilter != "" {
			errMsg = fmt.Sprintf("%s with filter %s", errMsg, searchFilter)
		}
		return fmt.Errorf("%s. Error: %v", errMsg, err)
	}

	targetUser := selectIdentityUser(FilterUsers, userName)
	if targetUser == nil {
		errMsg := fmt.Sprintf("Could not find user with name %s", userName)
		if searchFilter != "" {
			errMsg = fmt.Sprintf("%s with filter %s", errMsg, searchFilter)
		}
		return fmt.Errorf(errMsg)
	}

	d.SetId(*targetUser.Descriptor)
	d.Set("descriptor", targetUser.Descriptor)
	return nil
}

func getIdentityUsersWithFilterValue(clients *client.AggregatedClient, searchfilter string, filtervalue string) (*[]identity.Identity, error) {
	response, err := clients.IdentityClient.ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
		SearchFilter: &searchfilter, // Filter to get users
		FilterValue:  &filtervalue,  // Search String for user
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func selectIdentityUser(users *[]identity.Identity, userName string) *identity.Identity {
	for _, user := range *users {
		if strings.EqualFold(*user.ProviderDisplayName, userName) {
			return &user
		}
	}
	return nil
}
