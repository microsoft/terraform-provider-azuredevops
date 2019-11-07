package azuredevops

import (
	"fmt"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"math/rand"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
)

func resourceGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupMembershipCreate,
		Read:   resourceGroupMembershipRead,
		Update: resourceGroupMembershipUpdate,
		Delete: resourceGroupMembershipDelete,

		Schema: map[string]*schema.Schema{
			"group": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"members": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				MinItems: 1,
				Required: true,
			},
		},
	}
}

func resourceGroupMembershipCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)
	memberships := expandGroupMemberships(d)

	err := addMemberships(clients, memberships)
	if err != nil {
		return fmt.Errorf("Error adding group memberships during create: %+v", err)
	}

	// The ID for this resource is meaningless so we can just assign a random ID
	d.SetId(fmt.Sprintf("%d", rand.Int()))
	return nil
}

func resourceGroupMembershipUpdate(d *schema.ResourceData, m interface{}) error {
	if !d.HasChange("members") {
		return nil
	}

	group := d.Get("group").(string)
	oldMembers, newMembers := getOldAndNewMemberSetsFromResourceData(d)
	toAdd, toRemove := computeMembershipDiff(group, oldMembers, newMembers)
	return applyMembershipUpdate(m.(*aggregatedClient), toAdd, toRemove)
}

func applyMembershipUpdate(clients *aggregatedClient, toAdd *[]graph.GraphMembership, toRemove *[]graph.GraphMembership) error {
	err := removeMemberships(clients, toRemove)
	if err != nil {
		return fmt.Errorf("Error removing group memberships during update: %+v", err)
	}

	err = addMemberships(clients, toAdd)
	if err != nil {
		return fmt.Errorf("Error adding group memberships during update: %+v", err)
	}

	return nil
}

// Computes the memberships to add and remove. This should only be called during an Update operation
//	The first element returned are memberships to add
// 	The second element returned are memberships to remove
func computeMembershipDiff(group string, oldMembers map[string]bool, newMembers map[string]bool) (*[]graph.GraphMembership, *[]graph.GraphMembership) {
	membersToAdd, membersToRemove := []graph.GraphMembership{}, []graph.GraphMembership{}

	// members that need to be added will be missing from the old data, but present in the new data
	for member := range newMembers {
		if _, exists := oldMembers[member]; !exists {
			membersToAdd = append(membersToAdd, *buildMembership(group, member))
		}
	}

	// members that need to be removed will be missing from the new data, but present in the old data
	for member := range oldMembers {
		if _, exists := newMembers[member]; !exists {
			membersToRemove = append(membersToRemove, *buildMembership(group, member))
		}
	}

	return &membersToAdd, &membersToRemove
}

// Pull the "old" and "new" membership information from the state and convert them into string sets.
//
// If you are curious about the return type, have a read through this article:
//	https://stackoverflow.com/questions/34018908/golang-why-dont-we-have-a-set-datastructure
func getOldAndNewMemberSetsFromResourceData(d *schema.ResourceData) (map[string]bool, map[string]bool) {
	oldData, newData := d.GetChange("members")
	oldMembers := toStringSet(oldData.([]interface{}))
	newMembers := toStringSet(newData.([]interface{}))
	return oldMembers, newMembers
}

// Convert a list of elements into a set of strings
func toStringSet(items ...interface{}) map[string]bool {
	theSet := map[string]bool{}
	for _, item := range items {
		theSet[item.(string)] = true
	}

	return theSet
}

func resourceGroupMembershipDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)
	memberships := expandGroupMemberships(d)

	err := removeMemberships(clients, memberships)
	if err != nil {
		return fmt.Errorf("Error removing group memberships during delete: %+v", err)
	}

	// this marks the resource as deleted
	d.SetId("")
	return nil
}

// Add members to a group using the AzDO REST API. If any error is encountered, the function immediately returns.
func addMemberships(clients *aggregatedClient, memberships *[]graph.GraphMembership) error {
	for _, membership := range *memberships {
		_, err := clients.GraphClient.AddMembership(clients.ctx, graph.AddMembershipArgs{
			SubjectDescriptor:   membership.MemberDescriptor,
			ContainerDescriptor: membership.ContainerDescriptor,
		})

		if err != nil {
			return fmt.Errorf("Error adding member to group: %+v", err)
		}
	}

	return nil
}

// Remove members from a group using the AzDO REST API. If any error is encountered, the function immediately returns.
func removeMemberships(clients *aggregatedClient, memberships *[]graph.GraphMembership) error {
	for _, membership := range *memberships {
		err := clients.GraphClient.RemoveMembership(clients.ctx, graph.RemoveMembershipArgs{
			SubjectDescriptor:   membership.MemberDescriptor,
			ContainerDescriptor: membership.ContainerDescriptor,
		})

		if err != nil {
			return fmt.Errorf("Error removing member from group: %+v", err)
		}
	}

	return nil
}

func expandGroupMemberships(d *schema.ResourceData) *[]graph.GraphMembership {
	group := d.Get("group").(string)
	members := d.Get("members").(*schema.Set).List()
	memberships := make([]graph.GraphMembership, len(members))

	for i, member := range members {
		memberships[i] = *buildMembership(group, member.(string))
	}

	return &memberships
}

func buildMembership(group string, member string) *graph.GraphMembership {
	return &graph.GraphMembership{
		ContainerDescriptor: &group,
		MemberDescriptor:    &member,
	}
}

func resourceGroupMembershipRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)
	group := d.Get("group").(string)

	actualMemberships, err := clients.GraphClient.ListMemberships(clients.ctx, graph.ListMembershipsArgs{
		SubjectDescriptor: &group,
		Direction:         &graph.GraphTraversalDirectionValues.Down,
		Depth:             converter.Int(1),
	})
	if err != nil {
		return fmt.Errorf("Error reading group memberships during read: %+v", err)
	}

	members := make([]string, len(*actualMemberships))
	for i, membership := range *actualMemberships {
		members[i] = *membership.MemberDescriptor
	}

	d.Set("members", members)
	return nil
}
