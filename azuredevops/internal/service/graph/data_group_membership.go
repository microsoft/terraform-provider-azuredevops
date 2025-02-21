package graph

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func DataGroupMembership() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupMembershipRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"group_descriptor": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceGroupMembershipRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	groupDescriptor := d.Get("group_descriptor").(string)

	memberShips, err := clients.GraphClient.ListMemberships(clients.Ctx, graph.ListMembershipsArgs{
		SubjectDescriptor: &groupDescriptor,
		Direction:         &graph.GraphTraversalDirectionValues.Down,
		Depth:             converter.Int(1),
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return diag.Errorf(" Group with with descriptor: %s not found. Error: %v", groupDescriptor, err)
		}
		return diag.Errorf(" Reading group memberships during read. Group descriptor: %s . Error: %+v", groupDescriptor, err)
	}

	members := make([]string, 0)
	for _, membership := range *memberShips {
		members = append(members, *membership.MemberDescriptor)
	}

	d.SetId(groupDescriptor)
	d.Set("members", members)
	return nil
}
