package graph

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

// DataGroup schema and implementation for group data source
func DataGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
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
			"origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"origin_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Performs a lookup of a project group. This involves the following actions:
//
//	(1) Identify AzDO graph descriptor for the project in which the group exists
//	(2) Query for all AzDO groups that exist within the project. This leverages the AzDO graph descriptor for the project.
//		This involves querying a paginated API, so multiple API calls may be needed for this step.
//	(3) Select group that has the name identified by the schema
func dataSourceGroupRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	groupName, projectID := d.Get("name").(string), d.Get("project_id").(string)

	projectDescriptor, err := getProjectDescriptor(clients, projectID)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return fmt.Errorf("Project with with ID %s was not found. Error: %v", projectID, err)
		}
		return fmt.Errorf("Error finding descriptor for project with ID %s. Error: %v", projectID, err)
	}

	projectGroups, err := getGroupsForDescriptor(clients, projectDescriptor)
	if err != nil {
		errMsg := "Error finding groups"
		if projectID != "" {
			errMsg = fmt.Sprintf("%s for project with ID %s", errMsg, projectID)
		}
		return fmt.Errorf("%s. Error: %v", errMsg, err)
	}

	targetGroup := selectGroup(projectGroups, groupName)
	if targetGroup == nil {
		errMsg := fmt.Sprintf("Could not find group with name %s", groupName)
		if projectID != "" {
			errMsg = fmt.Sprintf("%s in project with ID %s", errMsg, projectID)
		}
		return fmt.Errorf(errMsg)
	}

	d.SetId(*targetGroup.Descriptor)
	d.Set("descriptor", targetGroup.Descriptor)
	d.Set("origin", targetGroup.Origin)
	d.Set("origin_id", targetGroup.OriginId)
	return nil
}

func getProjectDescriptor(clients *client.AggregatedClient, projectID string) (string, error) {
	if projectID == "" {
		return "", nil
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return "", err
	}

	descriptor, err := clients.GraphClient.GetDescriptor(clients.Ctx, graph.GetDescriptorArgs{StorageKey: &projectUUID})
	if err != nil {
		return "", err
	}

	return *descriptor.Value, nil
}

func getGroupsForDescriptor(clients *client.AggregatedClient, projectDescriptor string) (*[]graph.GraphGroup, error) {
	var groups []graph.GraphGroup
	var currentToken string

	for hasMore := true; hasMore; {
		newGroups, latestToken, err := getGroupsWithContinuationToken(clients, projectDescriptor, currentToken)
		currentToken = latestToken
		if err != nil {
			return nil, err
		}

		if newGroups != nil && len(*newGroups) > 0 {
			if projectDescriptor == "" {
				// filter on collection groups
				filteredGroups := []graph.GraphGroup{}
				for _, grp := range *newGroups {
					if grp.Domain == nil {
						continue
					}

					domain := strings.ToLower(*grp.Domain)
					if strings.HasPrefix(domain, "vstfs:///framework/identitydomain") ||
						(strings.HasPrefix(domain, "vstfs:///framework/generic")) {
						filteredGroups = append(filteredGroups, grp)
					}
				}
				groups = append(groups, filteredGroups...)
			} else {
				groups = append(groups, *newGroups...)
			}
		}
		hasMore = currentToken != ""
	}

	return &groups, nil
}

func getGroupsWithContinuationToken(clients *client.AggregatedClient, projectDescriptor string, continuationToken string) (*[]graph.GraphGroup, string, error) {
	args := graph.ListGroupsArgs{}
	if projectDescriptor != "" {
		args.ScopeDescriptor = &projectDescriptor
	}
	if continuationToken != "" {
		args.ContinuationToken = &continuationToken
	}

	response, err := clients.GraphClient.ListGroups(clients.Ctx, args)
	if err != nil {
		return nil, "", err
	}

	if response.ContinuationToken != nil && len(*response.ContinuationToken) > 1 {
		return nil, "", fmt.Errorf("Expected at most 1 continuation token, but found %d", len(*response.ContinuationToken))
	}

	var newToken string
	if response.ContinuationToken != nil && len(*response.ContinuationToken) > 0 {
		newToken = (*response.ContinuationToken)[0]
	}

	return response.GraphGroups, newToken, nil
}

func selectGroup(groups *[]graph.GraphGroup, groupName string) *graph.GraphGroup {
	for _, group := range *groups {
		if strings.EqualFold(*group.DisplayName, groupName) {
			return &group
		}
	}
	return nil
}
