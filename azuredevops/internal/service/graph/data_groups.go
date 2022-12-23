package graph

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// DataGroups schema and implementation for group data source
func DataGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGroupsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"groups": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      getGroupHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Optional: true,
						},

						"mail_address": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"display_name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"principal_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
func dataSourceGroupsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)

	projectDescriptor, err := getProjectDescriptor(clients, projectID)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return fmt.Errorf("Project with with ID %s was not found. Error: %v", projectID, err)
		}
		return fmt.Errorf("Error finding descriptor for project with ID %s. Error: %v", projectID, err)
	}

	groups, err := getGroupsForDescriptor(clients, projectDescriptor)
	if err != nil {
		errMsg := "Error finding groups"
		if projectID != "" {
			errMsg = fmt.Sprintf("%s for project with ID %s", errMsg, projectID)
		}
		return fmt.Errorf("%s. Error: %v", errMsg, err)
	}

	fgroups, err := flattenGroups(groups)
	if err != nil {
		return fmt.Errorf("Error flatten groups. Error: %w", err)
	}

	d.SetId("groups-" + uuid.New().String())
	d.Set("groups", fgroups)
	return nil
}

func getGroupHash(v interface{}) int {
	return tfhelper.HashString(v.(map[string]interface{})["descriptor"].(string))
}

func flattenGroups(groups *[]graph.GraphGroup) ([]interface{}, error) {
	if groups == nil {
		return []interface{}{}, nil
	}

	results := make([]interface{}, len(*groups))
	for i, group := range *groups {
		s := make(map[string]interface{})

		if group.Descriptor != nil {
			s["descriptor"] = *group.Descriptor
		} else {
			return nil, fmt.Errorf("Group Object does not contain a descriptor")
		}
		if group.DisplayName != nil {
			s["display_name"] = *group.DisplayName
		}
		if group.Url != nil {
			s["url"] = *group.Url
		}
		if group.Origin != nil {
			s["origin"] = *group.Origin
		}
		if group.OriginId != nil {
			s["origin_id"] = *group.OriginId
		}
		if group.Domain != nil {
			s["domain"] = *group.Domain
		}
		if group.MailAddress != nil {
			s["mail_address"] = *group.MailAddress
		}
		if group.PrincipalName != nil {
			s["principal_name"] = *group.PrincipalName
		}
		if group.Description != nil {
			s["description"] = *group.Description
		}

		results[i] = s
	}
	return results, nil
}
