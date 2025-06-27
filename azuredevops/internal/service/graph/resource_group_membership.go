package graph

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
)

// ResourceGroupMembership schema and implementation for group membership resource
func ResourceGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupMembershipCreate,
		Read:   resourceGroupMembershipRead,
		Update: resourceGroupMembershipUpdate,
		Delete: resourceGroupMembershipDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"group": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "add",
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					"add", "overwrite",
				}, true),
			},
			"members": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceGroupMembershipCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	group := d.Get("group").(string)
	mode := d.Get("mode").(string)
	membersToAdd := d.Get("members").(*schema.Set)
	var membersToRemove *schema.Set = nil

	if strings.EqualFold("overwrite", mode) {
		actualMemberships, err := getGroupMemberships(clients, group)
		if err != nil {
			return fmt.Errorf("Reading group memberships during read: %+v", err)
		}
		actualMembershipsSet := getGroupMembershipSet(actualMemberships)
		if err != nil {
			return fmt.Errorf("Converting membership list to set: %+v", err)
		}
		membersToRemove = membersToAdd.Difference(actualMembershipsSet)
	} else {
		membersToRemove = getGroupMembershipSet(nil)
	}

	err := applyMembershipUpdate(m.(*client.AggregatedClient),
		expandGroupMembers(group, membersToAdd),
		expandGroupMembers(group, membersToRemove))
	if err != nil {
		return fmt.Errorf("Adding group memberships during create: %+v", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"Waiting"},
		Target:  []string{"Synched"},
		Refresh: func() (interface{}, string, error) {
			clients := m.(*client.AggregatedClient)
			state := "Waiting"
			actualMemberships, err := getGroupMemberships(clients, group)
			if err != nil {
				return nil, "", fmt.Errorf("Reading group memberships: %+v", err)
			}
			actualMembershipsSet := getGroupMembershipSet(actualMemberships)
			if err != nil {
				return nil, "", fmt.Errorf("Converting membership list to set: %+v", err)
			}
			if (membersToAdd == nil || actualMembershipsSet.Intersection(membersToAdd).Len() <= 0) &&
				(membersToRemove == nil || actualMembershipsSet.Intersection(membersToRemove).Len() <= 0) {
				state = "Synched"
			}

			return state, state, nil
		},
		Timeout:                   60 * time.Minute,
		MinTimeout:                5 * time.Second,
		Delay:                     5 * time.Second,
		ContinuousTargetOccurence: 2,
	}
	if _, err := stateConf.WaitForStateContext(clients.Ctx); err != nil {
		return fmt.Errorf("Error waiting for DevOps synching memberships for group  [%s]: %+v", group, err)
	}

	// The ID for this resource is meaningless so we can just assign a random ID
	d.SetId(fmt.Sprintf("%d", rand.Int()))

	return resourceGroupMembershipRead(d, m)
}

func resourceGroupMembershipRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	group := d.Get("group").(string)

	actualMemberships, err := getGroupMemberships(clients, group)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Eeading group memberships during read: %+v", err)
	}

	mode := d.Get("mode").(string)
	stateMembers := d.Get("members").(*schema.Set)
	members := make([]string, 0)
	for _, membership := range *actualMemberships {
		if strings.EqualFold("overwrite", mode) || stateMembers.Contains(*membership.MemberDescriptor) {
			members = append(members, *membership.MemberDescriptor)
		}
	}

	d.Set("members", members)
	return nil
}

func resourceGroupMembershipUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	if !d.HasChange("members") {
		return nil
	}

	group := d.Get("group").(string)
	oldData, newData := d.GetChange("members")
	// members that need to be added will be missing from the old data, but present in the new data
	membersToAdd := newData.(*schema.Set).Difference(oldData.(*schema.Set))
	// members that need to be removed will be missing from the new data, but present in the old data
	membersToRemove := oldData.(*schema.Set).Difference(newData.(*schema.Set))

	err := applyMembershipUpdate(m.(*client.AggregatedClient),
		expandGroupMembers(group, membersToAdd),
		expandGroupMembers(group, membersToRemove))
	if err != nil {
		return err
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"Waiting"},
		Target:  []string{"Synched"},
		Refresh: func() (interface{}, string, error) {
			state := "Waiting"
			actualMemberships, err := getGroupMemberships(clients, group)
			if err != nil {
				return nil, "", fmt.Errorf("Reading group memberships: %+v", err)
			}
			actualMembershipsSet := getGroupMembershipSet(actualMemberships)
			if actualMembershipsSet.Intersection(membersToAdd).Len() <= 0 &&
				actualMembershipsSet.Intersection(membersToRemove).Len() <= 0 {
				state = "Synched"
			}

			return state, state, nil
		},
		Timeout:                   60 * time.Minute,
		MinTimeout:                5 * time.Second,
		Delay:                     5 * time.Second,
		ContinuousTargetOccurence: 3,
	}
	if _, err := stateConf.WaitForStateContext(clients.Ctx); err != nil {
		return fmt.Errorf("Waiting for DevOps synching memberships for group  [%s]: %+v", group, err)
	}

	return resourceGroupMembershipRead(d, m)
}

func resourceGroupMembershipDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	group := d.Get("group").(string)
	membersToRemove := d.Get("members").(*schema.Set)
	memberships := expandGroupMembers(group, membersToRemove)

	err := removeMembers(clients, memberships)
	stateConf := &retry.StateChangeConf{
		Pending: []string{"Waiting"},
		Target:  []string{"Synched"},
		Refresh: func() (interface{}, string, error) {
			state := "Waiting"
			actualMemberships, err := getGroupMemberships(clients, group)
			if err != nil {
				return nil, "", fmt.Errorf("Reading group memberships: %+v", err)
			}
			actualMembershipsSet := getGroupMembershipSet(actualMemberships)
			if actualMembershipsSet.Intersection(membersToRemove).Len() <= 0 {
				state = "Synched"
			}

			return state, state, nil
		},
		Timeout:                   60 * time.Minute,
		MinTimeout:                5 * time.Second,
		Delay:                     5 * time.Second,
		ContinuousTargetOccurence: 2,
	}
	if _, err := stateConf.WaitForStateContext(clients.Ctx); err != nil {
		return fmt.Errorf("Waiting for DevOps synching memberships for group  [%s]: %+v", group, err)
	}
	if err != nil {
		return fmt.Errorf("Removing group memberships during delete: %+v", err)
	}

	return nil
}

func applyMembershipUpdate(clients *client.AggregatedClient, toAdd *[]graph.GraphMembership, toRemove *[]graph.GraphMembership) error {
	if toRemove != nil && len(*toRemove) > 0 {
		err := removeMembers(clients, toRemove)
		if err != nil {
			return fmt.Errorf("Removing group memberships during update: %+v", err)
		}
	}

	if toAdd != nil && len(*toAdd) > 0 {
		err := addMembers(clients, toAdd)
		if err != nil {
			return fmt.Errorf("Adding group memberships during update: %+v", err)
		}
	}
	return nil
}

// Add members to a group using the AzDO REST API. If any error is encountered, the function immediately returns.
func addMembers(clients *client.AggregatedClient, memberships *[]graph.GraphMembership) error {
	if memberships != nil {
		for _, membership := range *memberships {
			_, err := clients.GraphClient.AddMembership(clients.Ctx, graph.AddMembershipArgs{
				SubjectDescriptor:   membership.MemberDescriptor,
				ContainerDescriptor: membership.ContainerDescriptor,
			})
			if err != nil {
				return fmt.Errorf("Error adding member %s to group %s: %+v",
					converter.ToString(membership.MemberDescriptor, "nil"),
					converter.ToString(membership.ContainerDescriptor, "nil"),
					err)
			}
		}
	}
	return nil
}

// Remove members from a group using the AzDO REST API. If any error is encountered, the function immediately returns.
func removeMembers(clients *client.AggregatedClient, memberships *[]graph.GraphMembership) error {
	if memberships != nil {
		for _, membership := range *memberships {
			err := clients.GraphClient.RemoveMembership(clients.Ctx, graph.RemoveMembershipArgs{
				SubjectDescriptor:   membership.MemberDescriptor,
				ContainerDescriptor: membership.ContainerDescriptor,
			})
			if err != nil {
				return fmt.Errorf("Error removing member from group: %+v", err)
			}
		}
	}
	return nil
}

func expandGroupMembers(group string, memberSet *schema.Set) *[]graph.GraphMembership {
	if memberSet == nil || memberSet.Len() <= 0 {
		return &[]graph.GraphMembership{}
	}
	members := memberSet.List()
	memberships := make([]graph.GraphMembership, len(members))

	for i, member := range members {
		memberships[i] = graph.GraphMembership{
			ContainerDescriptor: &group,
			MemberDescriptor:    converter.String(member.(string)),
		}
	}

	return &memberships
}

func getGroupMemberships(clients *client.AggregatedClient, groupDescriptor string) (*[]graph.GraphMembership, error) {
	return clients.GraphClient.ListMemberships(clients.Ctx, graph.ListMembershipsArgs{
		SubjectDescriptor: &groupDescriptor,
		Direction:         &graph.GraphTraversalDirectionValues.Down,
		Depth:             converter.Int(1),
	})
}

func getGroupMembershipSet(members *[]graph.GraphMembership) *schema.Set {
	set := schema.NewSet(schema.HashString, nil)
	if nil != members {
		for _, member := range *members {
			set.Add(*member.MemberDescriptor)
		}
	}
	return set
}
