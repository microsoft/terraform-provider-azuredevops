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

// DataIdentityUser DataIdentityUserResource returns the user data source resource
func DataIdentityUser() *schema.Resource {
	return &schema.Resource{
		Read: dataIdentitySourceUserRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"search_filter": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "General",
				ValidateFunc: validation.StringInSlice([]string{"AccountName", "DisplayName", "MailAddress", "General"}, false),
			},
			"descriptor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject_descriptor": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataIdentitySourceUserRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	userName := d.Get("name").(string)
	searchFilter := d.Get("search_filter").(string)

	// Query ADO for list of identity user with filter
	filterUser, err := getIdentityUsersWithFilterValue(clients, searchFilter, userName)
	if err != nil {
		return fmt.Errorf("Finding user with filter %s. Error: %v", searchFilter, err)
	}

	// Filter for the desired user in the FilterUsers results
	targetUser := validateIdentityUser(filterUser, userName, searchFilter)
	if targetUser == nil {
		return fmt.Errorf("Could not find user with name: %s, with filter: %s", userName, searchFilter)
	}

	// Set id and user list for users data resource
	targetUserID := targetUser.Id.String()
	d.SetId(targetUserID)
	d.Set("descriptor", targetUser.Descriptor)
	d.Set("subject_descriptor", targetUser.SubjectDescriptor)
	return nil
}

// Query AZDO for users with matching filter and search string
func getIdentityUsersWithFilterValue(clients *client.AggregatedClient, searchFilter string, filterValue string) (*[]identity.Identity, error) {
	// Get list of user with search filter and filter value provided at data source invocation.
	response, err := clients.IdentityClient.ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
		SearchFilter: &searchFilter, // Filter to get users
		FilterValue:  &filterValue,  // Search String for user
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Filter results to validate user is correct. Occurs post-flatten due to missing properties based on search-filter.
func validateIdentityUser(users *[]identity.Identity, userName string, searchFilter string) *identity.Identity {
	for _, user := range *users {
		prop := user.Properties.(map[string]interface{})

		switch searchFilter {
		case "General":
			return &user
		case "DisplayName":
			if strings.Contains(strings.ToLower(*user.ProviderDisplayName), strings.ToLower(userName)) {
				return &user
			}
		case "MailAddress":
			if v, ok := prop["Mail"]; ok && v != nil {
				mailProp := v.(map[string]interface{})
				if emailAddress, ok := mailProp["$value"].(string); ok {
					if strings.Contains(strings.ToLower(emailAddress), strings.ToLower(userName)) {
						return &user
					}
				}
			}
		case "AccountName":
			if v, ok := prop["Account"]; ok && v != nil {
				mailProp := v.(map[string]interface{})
				if emailAddress, ok := mailProp["$value"].(string); ok {
					if strings.Contains(strings.ToLower(emailAddress), strings.ToLower(userName)) {
						return &user
					}
				}
			}
		}
	}
	return nil
}
