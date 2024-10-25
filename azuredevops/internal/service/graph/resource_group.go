package graph

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceGroup schema and implementation for group resource
func ResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"scope": {
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
				Optional:     true,
				ForceNew:     true,
			},

			"origin_id": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.NoZeroValues,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"mail", "display_name", "scope"},
			},

			"mail": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.NoZeroValues,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"origin_id", "display_name", "scope"},
			},

			"display_name": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.NoZeroValues,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"origin_id", "mail"},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"members": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Computed:   true,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Set:        schema.HashString,
			},

			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"origin": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subject_kind": {
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

			"descriptor": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	var scopeDescriptor *string
	if val, ok := d.GetOk("scope"); ok {
		scopeUid, _ := uuid.Parse(val.(string))
		desc, err := clients.GraphClient.GetDescriptor(clients.Ctx, graph.GetDescriptorArgs{
			StorageKey: &scopeUid,
		})
		if err != nil {
			return err
		}
		scopeDescriptor = desc.Value
	}

	var group *graph.GraphGroup
	var err error
	if v, ok := d.GetOk("display_name"); ok {
		param := graph.CreateGroupVstsArgs{
			CreationContext: &graph.GraphGroupVstsCreationContext{
				DisplayName: converter.String(v.(string)),
				Description: converter.String(d.Get("description").(string)),
			},
			ScopeDescriptor: scopeDescriptor,
		}
		group, err = clients.GraphClient.CreateGroupVsts(clients.Ctx, param)
		if err != nil {
			return err
		}
	} else if v, ok := d.GetOk("origin_id"); ok {
		param := graph.CreateGroupOriginIdArgs{
			CreationContext: &graph.GraphGroupOriginIdCreationContext{
				OriginId: converter.String(v.(string)),
			},
			ScopeDescriptor: scopeDescriptor,
		}
		group, err = clients.GraphClient.CreateGroupOriginId(clients.Ctx, param)
		if err != nil {
			return err
		}
	} else if v, ok := d.GetOk("mail"); ok {
		param := graph.CreateGroupMailAddressArgs{
			CreationContext: &graph.GraphGroupMailAddressCreationContext{
				MailAddress: converter.String(v.(string)),
			},
			ScopeDescriptor: scopeDescriptor,
		}
		group, err = clients.GraphClient.CreateGroupMailAddress(clients.Ctx, param)
		if err != nil {
			return err
		}
	}

	stateMembers, ok := d.GetOk("members")
	if ok {
		members := expandGroupMembers(*group.Descriptor, stateMembers.(*schema.Set))
		if err := addMembers(clients, members); err != nil {
			return fmt.Errorf(" adding group memberships during create: %+v", err)
		}
	}

	d.SetId(*group.Descriptor)
	return resourceGroupRead(d, m)
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	group, err := clients.GraphClient.GetGroup(
		clients.Ctx,
		graph.GetGroupArgs{GroupDescriptor: converter.String(d.Id())})

	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	if group.IsDeleted != nil && *group.IsDeleted {
		d.SetId("")
		return nil
	}

	members, err := groupReadMembers(*group.Descriptor, clients)
	if err != nil {
		return err
	}

	flattenGroup(d, group, members)

	storageKey, err := clients.GraphClient.GetStorageKey(clients.Ctx, graph.GetStorageKeyArgs{
		SubjectDescriptor: group.Descriptor,
	})
	if err != nil {
		return err
	}

	if storageKey.Value != nil {
		d.Set("group_id", storageKey.Value.String())
	}

	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	var operations []webapi.JsonPatchOperation

	if d.HasChange("display_name") {
		displayName := d.Get("display_name")
		patchDisplayNameOperation := webapi.JsonPatchOperation{
			Op:    &webapi.OperationValues.Replace,
			From:  nil,
			Path:  converter.String("/displayName"),
			Value: displayName.(string),
		}
		operations = append(operations, patchDisplayNameOperation)
	}

	if d.HasChange("description") {
		description := d.Get("description")
		patchDescriptionOperation := webapi.JsonPatchOperation{
			Op:    &webapi.OperationValues.Replace,
			From:  nil,
			Path:  converter.String("/description"),
			Value: description.(string),
		}
		operations = append(operations, patchDescriptionOperation)
	}

	if len(operations) > 0 {
		uptGroupArgs := graph.UpdateGroupArgs{
			GroupDescriptor: converter.String(d.Id()),
			PatchDocument:   &operations,
		}

		_, err := clients.GraphClient.UpdateGroup(clients.Ctx, uptGroupArgs)
		if err != nil {
			return err
		}
	}

	if d.HasChange("members") {
		group := d.Id()
		oldData, newData := d.GetChange("members")
		// members that need to be added will be missing from the old data, but present in the new data
		membersToAdd := newData.(*schema.Set).Difference(oldData.(*schema.Set))
		// members that need to be removed will be missing from the new data, but present in the old data
		membersToRemove := oldData.(*schema.Set).Difference(newData.(*schema.Set))
		if err := applyMembershipUpdate(m.(*client.AggregatedClient),
			expandGroupMembers(group, membersToAdd),
			expandGroupMembers(group, membersToRemove)); err != nil {
			return err
		}
	}

	return resourceGroupRead(d, m)
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	stateConf := &retry.StateChangeConf{
		Pending: []string{"Waiting"},
		Target:  []string{"Succeed", "Failed"},
		Refresh: func() (interface{}, string, error) {
			err := clients.GraphClient.DeleteGroup(clients.Ctx, graph.DeleteGroupArgs{
				GroupDescriptor: converter.String(d.Id()),
			})
			if err != nil {
				if utils.ResponseWasNotFound(err) {
					return "", "Succeed", nil
				}
				return nil, "Failed", err
			}

			return nil, "Waiting", nil
		},
		Timeout:    60 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(clients.Ctx); err != nil {
		return fmt.Errorf(" Waiting for group delete. %v ", err)
	}

	return nil
}

func flattenGroup(d *schema.ResourceData, group *graph.GraphGroup, members *[]graph.GraphMembership) {
	d.Set("descriptor", *group.Descriptor)

	if group.DisplayName != nil {
		d.Set("display_name", *group.DisplayName)
	}
	if group.Url != nil {
		d.Set("url", *group.Url)
	}
	if group.Origin != nil {
		d.Set("origin", *group.Origin)
	}
	if group.OriginId != nil {
		d.Set("origin_id", *group.OriginId)
	}
	if group.SubjectKind != nil {
		d.Set("subject_kind", *group.SubjectKind)
	}
	if group.Domain != nil {
		d.Set("domain", *group.Domain)
	}
	if group.MailAddress != nil {
		d.Set("mail", *group.MailAddress)
	}
	if group.PrincipalName != nil {
		d.Set("principal_name", *group.PrincipalName)
	}
	if group.Description != nil {
		d.Set("description", *group.Description)
	}
	if members != nil {
		dMembers := make([]string, len(*members))
		for i, e := range *members {
			dMembers[i] = *e.MemberDescriptor
		}
		d.Set("members", dMembers)
	}

	if projectId := domain2ProjectID(*group.Domain); projectId != "" {
		d.Set("scope", projectId)
	}
}

func groupReadMembers(groupDescriptor string, clients *client.AggregatedClient) (*[]graph.GraphMembership, error) {
	actualMembers, err := clients.GraphClient.ListMemberships(clients.Ctx, graph.ListMembershipsArgs{
		SubjectDescriptor: &groupDescriptor,
		Direction:         &graph.GraphTraversalDirectionValues.Down,
		Depth:             converter.Int(1),
	})
	if err != nil {
		return nil, fmt.Errorf(" Reading group memberships: %+v", err)
	}

	members := make([]graph.GraphMembership, len(*actualMembers))
	for i, membership := range *actualMembers {
		members[i] = graph.GraphMembership{
			ContainerDescriptor: &groupDescriptor,
			MemberDescriptor:    membership.MemberDescriptor,
		}
	}

	return &members, nil
}

func domain2ProjectID(domain string) (projectID string) {
	if domain == "" {
		return ""
	}
	if !strings.HasPrefix(domain, "vstfs:///Classification/TeamProject") {
		return ""
	}
	return domain[36:]
}
