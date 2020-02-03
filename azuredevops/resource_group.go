package azuredevops

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"scope": {
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
				Optional:     true,
				ForceNew:     true,
			},

			// ***
			// One of
			//     origin_id => GraphGroupOriginIdCreationContext
			//     mail => GraphGroupMailAddressCreationContext
			//     display_name => GraphGroupVstsCreationContext
			// must be specified
			// ***

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
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"origin_id", "mail"},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
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
		},
	}
}

type azDOGraphCreateGroupArgs struct {
	// (required) The subset of the full graph group used to uniquely find the graph subject in an external provider.
	CreationContext interface{}
	// (optional) A descriptor referencing the scope (collection, project) in which the group should be created.
	// If omitted, will be created in the scope of the enclosing account or organization. Valid only for VSTS groups.
	ScopeDescriptor *string
	// (optional) A comma separated list of descriptors referencing groups you want the graph group to join
	GroupDescriptors *[]string
}

func azDOGraphCreateGroup(ctx context.Context, client graph.Client, args azDOGraphCreateGroupArgs) (*graph.GraphGroup, error) {

	if args.CreationContext == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.CreationContext"}
	}
	queryParams := url.Values{}
	if args.ScopeDescriptor != nil {
		queryParams.Add("scopeDescriptor", *args.ScopeDescriptor)
	}
	if args.GroupDescriptors != nil {
		listAsString := strings.Join((*args.GroupDescriptors)[:], ",")
		queryParams.Add("groupDescriptors", listAsString)
	}

	if _, ok := args.CreationContext.(*graph.GraphGroupMailAddressCreationContext); !ok {
		if _, ok := args.CreationContext.(*graph.GraphGroupOriginIdCreationContext); !ok {
			if _, ok := args.CreationContext.(*graph.GraphGroupVstsCreationContext); !ok {
				return nil, fmt.Errorf("Unsupported group creation context")
			}
		}
	}

	body, marshalErr := json.Marshal(args.CreationContext)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationID, _ := uuid.Parse("ebbe6af8-0b91-4c13-8cf1-777c14858188")
	if clientImpl, ok := client.(*graph.ClientImpl); ok {
		resp, err := clientImpl.Client.Send(ctx, http.MethodPost, locationID, "5.1-preview.1", nil, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
		if err != nil {
			return nil, err
		}
		var responseValue graph.GraphGroup
		err = clientImpl.Client.UnmarshalBody(resp, &responseValue)
		return &responseValue, err
	}

	panic("Invalid Azure DevOps Graph client implementation")
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)

	// using: POST https://vssps.dev.azure.com/{organization}/_apis/graph/groups?api-version=5.1-preview.1
	cga := azDOGraphCreateGroupArgs{}
	val, b := d.GetOk("scope")
	if b {
		uuid, _ := uuid.Parse(val.(string))
		desc, err := clients.GraphClient.GetDescriptor(clients.Ctx, graph.GetDescriptorArgs{
			StorageKey: &uuid,
		})
		if err != nil {
			return err
		}
		cga.ScopeDescriptor = desc.Value
	}
	val, b = d.GetOk("origin_id")
	if b {
		if _, b = d.GetOk("mail"); b {
			return fmt.Errorf("Unable to create group with invalid parameters: mail")
		}
		if _, b = d.GetOk("display_name"); b {
			return fmt.Errorf("Unable to create group with invalid parameters: display_name")
		}
		cga.CreationContext = &graph.GraphGroupOriginIdCreationContext{
			OriginId: converter.String(val.(string)),
		}
	} else {
		val, b = d.GetOk("mail")
		if b {
			if _, b = d.GetOk("display_name"); b {
				return fmt.Errorf("Unable to create group with invalid parameters: display_name")
			}
			cga.CreationContext = &graph.GraphGroupMailAddressCreationContext{
				MailAddress: converter.String(val.(string)),
			}
		} else {
			val, b = d.GetOk("display_name")
			if b {
				cga.CreationContext = &graph.GraphGroupVstsCreationContext{
					DisplayName: converter.String(val.(string)),
					Description: converter.String(d.Get("description").(string)),
				}
			} else {
				return fmt.Errorf("INTERNAL ERROR: Unable to determine strategy to create group")
			}
		}
	}
	group, err := azDOGraphCreateGroup(clients.Ctx, clients.GraphClient, cga)
	if err != nil {
		return err
	}
	if group.Descriptor == nil {
		return fmt.Errorf("DevOps REST API returned group object without descriptor")
	}

	stateMembers, exists := d.GetOk("members")
	if exists {
		members := expandGroupMembers(*group.Descriptor, stateMembers.(*schema.Set))
		if err := addMembers(clients, members); err != nil {
			return fmt.Errorf("Error adding group memberships during create: %+v", err)
		}
	}

	d.SetId(*group.Descriptor)

	return resourceGroupRead(d, m)
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)

	// using: GET https://vssps.dev.azure.com/{organization}/_apis/graph/groups/{groupDescriptor}?api-version=5.1-preview.1
	// d.Get("descriptor").(string) => {groupDescriptor}
	getGroupArgs := graph.GetGroupArgs{
		GroupDescriptor: converter.String(d.Id()),
	}
	group, err := clients.GraphClient.GetGroup(clients.Ctx, getGroupArgs)
	if err != nil {
		return err
	}
	if group.Descriptor == nil {
		return fmt.Errorf("DevOps REST API returned group object without descriptor; group %s", d.Id())
	}

	members, err := groupReadMembers(*group.Descriptor, clients)
	if err != nil {
		return err
	}
	return flattenGroup(d, group, members)
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)

	// using: PATCH https://vssps.dev.azure.com/{organization}/_apis/graph/groups/{groupDescriptor}?api-version=5.1-preview.1
	// d.Get("descriptor").(string) => {groupDescriptor}

	d.Partial(true)

	if d.HasChange("description") {
		description := d.Get("description")
		uptGroupArgs := graph.UpdateGroupArgs{
			GroupDescriptor: converter.String(d.Id()),
			PatchDocument: &[]webapi.JsonPatchOperation{
				{
					Op:    &webapi.OperationValues.Replace,
					From:  nil,
					Path:  converter.String("/description"),
					Value: description.(string),
				},
			},
		}

		_, err := clients.GraphClient.UpdateGroup(clients.Ctx, uptGroupArgs)
		if err != nil {
			return err
		}

		d.SetPartial("description")
	}

	if d.HasChange("members") {
		group := d.Id()
		oldData, newData := d.GetChange("members")
		// members that need to be added will be missing from the old data, but present in the new data
		membersToAdd := newData.(*schema.Set).Difference(oldData.(*schema.Set))
		// members that need to be removed will be missing from the new data, but present in the old data
		membersToRemove := oldData.(*schema.Set).Difference(newData.(*schema.Set))
		if err := applyMembershipUpdate(m.(*config.AggregatedClient),
			expandGroupMembers(group, membersToAdd),
			expandGroupMembers(group, membersToRemove)); err != nil {
			return err
		}

		d.SetPartial("members")
	}

	d.Partial(false)
	return resourceGroupRead(d, m)
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)

	// using: DELETE https://vssps.dev.azure.com/{organization}/_apis/graph/groups/{groupDescriptor}?api-version=5.1-preview.1
	// d.Get("descriptor").(string) => {groupDescriptor}
	delGroupArgs := graph.DeleteGroupArgs{
		GroupDescriptor: converter.String(d.Id()),
	}
	err := clients.GraphClient.DeleteGroup(clients.Ctx, delGroupArgs)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

// Convert internal Terraform data structure to an AzDO data structure
func expandGroup(d *schema.ResourceData) (*graph.GraphGroup, *[]graph.GraphMembership, error) {

	group := graph.GraphGroup{
		Descriptor:    converter.String(d.Id()),
		DisplayName:   converter.String(d.Get("display_name").(string)),
		Url:           converter.String(d.Get("url").(string)),
		Origin:        converter.String(d.Get("origin").(string)),
		OriginId:      converter.String(d.Get("origin_id").(string)),
		SubjectKind:   converter.String(d.Get("subject_kind").(string)),
		Domain:        converter.String(d.Get("domain").(string)),
		MailAddress:   converter.String(d.Get("mail").(string)),
		PrincipalName: converter.String(d.Get("principal_name").(string)),
		Description:   converter.String(d.Get("description").(string)),
	}

	stateMembers := d.Get("members")
	dMembers := stateMembers.(*schema.Set).List()
	members := make([]graph.GraphMembership, len(dMembers))
	for i, e := range dMembers {
		members[i] = graph.GraphMembership{
			ContainerDescriptor: group.Descriptor,
			MemberDescriptor:    converter.String(e.(string)),
		}
	}
	return &group, &members, nil
}

func flattenGroup(d *schema.ResourceData, group *graph.GraphGroup, members *[]graph.GraphMembership) error {

	if group.Descriptor != nil {
		d.Set("descriptor", *group.Descriptor)
		d.SetId(*group.Descriptor)
	} else {
		return fmt.Errorf("Group Object does not contain a descriptor")
	}
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
	return nil
}

func groupReadMembers(groupDescriptor string, clients *config.AggregatedClient) (*[]graph.GraphMembership, error) {
	actualMembers, err := clients.GraphClient.ListMemberships(clients.Ctx, graph.ListMembershipsArgs{
		SubjectDescriptor: &groupDescriptor,
		Direction:         &graph.GraphTraversalDirectionValues.Down,
		Depth:             converter.Int(1),
	})
	if err != nil {
		return nil, fmt.Errorf("Error reading group memberships during read: %+v", err)
	}

	members := make([]graph.GraphMembership, len(*actualMembers))
	for i, membership := range *actualMembers {
		members[i] = *buildMembership(groupDescriptor, *membership.MemberDescriptor)
	}

	return &members, nil
}
