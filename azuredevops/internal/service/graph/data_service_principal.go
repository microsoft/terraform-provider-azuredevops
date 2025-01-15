package graph

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
)

// DataServicePrincipal schema and implementation for service principal data source
func DataServicePrincipal() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServicePrincipalRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validation.StringIsNotWhiteSpace,
			},
			"origin_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"descriptor": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Performs a lookup of a service principal. This involves the following actions:
//
//	(1) Identify AzDO graph descriptor for the service principal
//	(2) Get the service principal by descriptor
func dataSourceServicePrincipalRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	searchFilter := converter.String("General")
	displayName := d.Get("display_name").(string)

	// Query ADO for list of identity user with filter
	filteredServicePrincipals, err := getIdentityServicePrincipalsWithFilterValue(clients, searchFilter, displayName)
	if err != nil || filteredServicePrincipals == nil || len(*filteredServicePrincipals) == 0 {
		return fmt.Errorf(" Finding service principal with filter %s. Error: %v", *searchFilter, err)
	}

	flattenedServicePrincipals, err := flattenIdentityServicePrincipals(filteredServicePrincipals)
	if err != nil {
		return fmt.Errorf(" Flatten service principals. Error: %v", err)
	}

	// Filter for the desired user in the FilterUsers results
	targetServicePrincipal := validateIdentityServicePrincipal(flattenedServicePrincipals, displayName)
	if targetServicePrincipal == nil {
		return fmt.Errorf(" Could not find service principal with name: %s", displayName)
	}

	servicePrincipalDescriptor := targetServicePrincipal.SubjectDescriptor

	servicePrincipal, err := getServicePrincipal(clients, servicePrincipalDescriptor)
	if err != nil {
		if servicePrincipalDescriptor != nil {
			return fmt.Errorf(" Finding service principal with Descriptor %s", *servicePrincipalDescriptor)
		}
		return fmt.Errorf(" Finding service principal. Error: %v", err)
	}

	d.SetId(*servicePrincipal.Descriptor)
	d.Set("descriptor", servicePrincipal.Descriptor)
	d.Set("display_name", servicePrincipal.DisplayName)
	d.Set("origin_id", servicePrincipal.OriginId)
	d.Set("origin", servicePrincipal.Origin)

	return nil
}

// Query AZDO for service principals with matching filter and search string
func getIdentityServicePrincipalsWithFilterValue(clients *client.AggregatedClient, searchFilter *string, filterValue string) (*[]identity.Identity, error) {
	// Get list of user with search filter and filter value provided at data source invocation.
	return clients.IdentityClient.ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
		SearchFilter: searchFilter, // Filter to get users
		FilterValue:  &filterValue,
	})
}

// Flatten Query Results
func flattenIdentityServicePrincipals(servicePrincipals *[]identity.Identity) (*[]identity.Identity, error) {
	if servicePrincipals == nil || len(*servicePrincipals) == 0 {
		return nil, fmt.Errorf(" Input Service Principals Parameter is nil")
	}
	results := make([]identity.Identity, 0)
	for _, servicePrincipal := range *servicePrincipals {
		if servicePrincipal.Descriptor == nil {
			return nil, fmt.Errorf(" User Object does not contain an id")
		}
		results = append(results, identity.Identity{
			Id:                  servicePrincipal.Id,
			Descriptor:          servicePrincipal.Descriptor,
			ProviderDisplayName: servicePrincipal.ProviderDisplayName,
			SubjectDescriptor:   servicePrincipal.SubjectDescriptor,
		})
	}
	return &results, nil
}

// Filter results to validate user is correct. Occurs post-flatten due to missing properties based on search-filter.
func validateIdentityServicePrincipal(servicePrincipals *[]identity.Identity, displayName string) *identity.Identity {
	if servicePrincipals == nil || len(*servicePrincipals) == 0 {
		return nil
	}
	for _, servicePrincipal := range *servicePrincipals {
		if strings.Contains(strings.ToLower(*servicePrincipal.ProviderDisplayName), strings.ToLower(displayName)) {
			return &servicePrincipal
		}
	}
	return nil
}

func getServicePrincipal(clients *client.AggregatedClient, servicePrincipalDescriptor *string) (*graph.GraphServicePrincipal, error) {
	return clients.GraphClient.GetServicePrincipal(clients.Ctx, graph.GetServicePrincipalArgs{
		ServicePrincipalDescriptor: servicePrincipalDescriptor,
	})
}
