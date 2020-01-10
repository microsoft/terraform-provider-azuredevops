package azuredevops

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/identity"
	"github.com/microsoft/azure-devops-go-api/azuredevops/security"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"
)

var projectSecurityNamespaceID = uuid.MustParse("52d39943-cb85-4d7f-8fa8-c6baac873819")

func resourceProjectPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectPermissionsCreate,
		Read:   resourceProjectPermissionsRead,
		Update: resourceProjectPermissionsUpdate,
		Delete: resourceProjectPermissionsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceProjectPermissionsImporter,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validate.UUID,
				Required:     true,
				ForceNew:     true,
			},
			"principals": {
				Type:     schema.TypeSet,
				Required: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"merge": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"permissions": {
				// Unable to define a validation function, because the
				// keys and values can only be validated with an initialized
				// security client as we must load the security namespace
				// definition and the available permission settings, and a validation
				// function in Terraform only receives the parameter name and the
				// current value as argument
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceProjectPermissionsCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	projectID, ok := d.GetOk("project_id")
	if !ok {
		return fmt.Errorf("Failed to get 'project_id' from schema")
	}
	secns, err := clients.Security.QuerySecurityNamespaces(clients.Ctx, security.QuerySecurityNamespacesArgs{
		SecurityNamespaceId: &projectSecurityNamespaceID,
	})
	if err != nil {
		return err
	}
	if secns == nil || len(*secns) <= 0 || (*secns)[0].Actions == nil || len(*(*secns)[0].Actions) <= 0 {
		return fmt.Errorf("Failed to load security namespace definition with id [%s]", projectSecurityNamespaceID.String())
	}

	actionMap := map[string]security.ActionDefinition{}
	for _, action := range *(*secns)[0].Actions {
		actionMap[*action.Name] = action
	}

	aclToken := fmt.Sprintf("$PROJECT:vstfs:///Classification/TeamProject/%s", projectID.(string))
	aclProject, err := clients.Security.QueryAccessControlLists(clients.Ctx, security.QueryAccessControlListsArgs{
		SecurityNamespaceId: &projectSecurityNamespaceID,
		Token:               converter.String(aclToken),
	})

	if err != nil {
		return err
	}
	if aclProject == nil || len(*aclProject) != 1 {
		return fmt.Errorf("Failed to load current ACL for project [%s]", projectID.(string))
	}

	principalList := d.Get("principals").(*schema.Set).List()
	descriptors := linq.From(principalList).Aggregate(
		func(r interface{}, i interface{}) interface{} {
			if r.(string) == "" {
				return i
			}
			return r.(string) + "," + i.(string)
		},
	)
	idlist, err := clients.Identity.ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
		SubjectDescriptors: converter.String(descriptors.(string)),
	})

	if err != nil {
		return err
	}
	if idlist == nil || len(*idlist) < len(principalList) {
		return fmt.Errorf("Failed to load identity information for defined principals [%s]", descriptors.(string))
	}

	idMap := map[string]identity.Identity{}
	linq.From(*idlist).
		ToMapBy(&idMap,
			func(item interface{}) interface{} { return *item.(identity.Identity).SubjectDescriptor },
			func(item interface{}) interface{} { return item })

	aceSet := schema.NewSet(func(item interface{}) int {
		return hashcode.String(*item.(security.AccessControlEntry).Descriptor)
	}, nil)
	for _, value := range *(*aclProject)[0].AcesDictionary {
		aceSet.Add(value)
	}

	for _, principal := range principalList {
		desc, ok := idMap[principal.(string)]
		if !ok {
			return fmt.Errorf("Unable to resolve id descriptor for principal [%s]", principal)
		}

		log.Printf("[TRACE]Checking ACE list for descriptor [%s]", *desc.Descriptor)
		var aceItem *security.AccessControlEntry
		ace, update := (*(*aclProject)[0].AcesDictionary)[*desc.Descriptor]
		if !update {
			log.Printf("[TRACE]Creating new ACE for subject [%s]", principal)
			// Create new ACE for pricipal
			subjectList, err := clients.GraphClient.LookupSubjects(clients.Ctx, graph.LookupSubjectsArgs{
				SubjectLookup: &graph.GraphSubjectLookup{
					LookupKeys: &[]graph.GraphSubjectLookupKey{
						{
							Descriptor: desc.SubjectDescriptor,
						},
					},
				},
			})
			if err != nil {
				return err
			}
			subject, ok := (*subjectList)[*desc.SubjectDescriptor]
			if !ok {
				return fmt.Errorf("Unable to lookup subject data from subject descriptor [%s]", *desc.SubjectDescriptor)
			}
			if !strings.EqualFold(*subject.SubjectKind, "group") {
				return fmt.Errorf("Dedicated user permissions are currently not implemented. Unable to set project permission for account [%s]", *subject.Descriptor)
			}
			aceItem = new(security.AccessControlEntry)
			aceItem.Allow = new(int)
			aceItem.Deny = new(int)
			aceItem.Descriptor = desc.Descriptor
		} else {
			// update existing ACE for principal
			log.Printf("[TRACE]Updating ACE for descriptor [%s]", *desc.Descriptor)
			aceItem = &ace
		}

		for key, value := range d.Get("permissions").(map[string]interface{}) {
			actionDef, ok := actionMap[key]
			if !ok {
				return fmt.Errorf("Invalid permission [%s]", key)
			}
			if aceItem.Deny == nil {
				aceItem.Deny = new(int)
			}
			if aceItem.Allow == nil {
				aceItem.Allow = new(int)
			}

			if strings.EqualFold("deny", value.(string)) {
				*aceItem.Allow = (*aceItem.Allow) &^ (*actionDef.Bit)
				*aceItem.Deny = (*aceItem.Deny) | (*actionDef.Bit)
			} else if strings.EqualFold("allow", value.(string)) {
				*aceItem.Deny = (*aceItem.Deny) &^ (*actionDef.Bit)
				*aceItem.Allow = (*aceItem.Allow) | (*actionDef.Bit)
			} else if strings.EqualFold("notset", value.(string)) {
				*aceItem.Allow = (*aceItem.Allow) &^ (*actionDef.Bit)
				*aceItem.Deny = (*aceItem.Deny) &^ (*actionDef.Bit)
			} else {
				return fmt.Errorf("Invalid permission action [%s]", value)
			}
		}
		aceSet.Add(*aceItem)
	}

	aceArr := []security.AccessControlEntry{}
	linq.From(aceSet.List()).ToSlice(&aceArr)
	bMerge := d.Get("merge").(bool)
	container := struct {
		Token                *string                        `json:"token,omitempty"`
		Merge                *bool                          `json:"merge,omitempty"`
		AccessControlEntries *[]security.AccessControlEntry `json:"accessControlEntries,omitempty"`
	}{
		Token:                &aclToken,
		Merge:                &bMerge,
		AccessControlEntries: &aceArr,
	}

	log.Printf("[TRACE]SetAccessControlEntries: %s", spew.Sdump(container))
	_, err = clients.Security.SetAccessControlEntries(clients.Ctx, security.SetAccessControlEntriesArgs{
		SecurityNamespaceId: &projectSecurityNamespaceID,
		Container:           container,
	})
	if err != nil {
		return err
	}
	return resourceProjectPermissionsRead(d, m)
}

func resourceProjectPermissionsRead(d *schema.ResourceData, m interface{}) error {
	//clients := m.(*config.AggregatedClient)

	return nil
}

func resourceProjectPermissionsUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceProjectPermissionsRead(d, m)
}

func resourceProjectPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceProjectPermissionsImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	return nil, errors.New("Not implemented")
}
