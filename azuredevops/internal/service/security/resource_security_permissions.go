package security

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/security"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// ResourceGenericPermissions schema and implementation for generic permission resource
func ResourceGenericPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceGenericPermissionsCreateOrUpdate,
		Read:   resourceGenericPermissionsRead,
		Update: resourceGenericPermissionsCreateOrUpdate,
		Delete: resourceGenericPermissionsDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"namespace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "The ID of the security namespace",
			},
			"token": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "The security token for the resource",
			},
			"principal": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "The descriptor or identity ID of the principal (user or group)",
			},
			"permissions": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"allow", "deny", "notset", "Allow", "Deny", "NotSet", "ALLOW", "DENY", "NOTSET"}, false),
				},
				Description: "Map of permission names to values (allow, deny, or notset)",
			},
			"replace": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Replace existing permissions (true) or merge with existing (false)",
			},
		},
	}
}

func resourceGenericPermissionsCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	namespaceID, err := uuid.Parse(d.Get("namespace_id").(string))
	if err != nil {
		return fmt.Errorf("invalid namespace_id: %v", err)
	}

	token := d.Get("token").(string)
	principal := d.Get("principal").(string)
	permissions := d.Get("permissions").(map[string]interface{})
	replace := d.Get("replace").(bool)

	// Get namespace details to retrieve action definitions
	namespaces, err := clients.SecurityClient.QuerySecurityNamespaces(clients.Ctx, security.QuerySecurityNamespacesArgs{
		SecurityNamespaceId: &namespaceID,
	})
	if err != nil {
		return fmt.Errorf("querying security namespace: %v", err)
	}
	if namespaces == nil || len(*namespaces) == 0 {
		return fmt.Errorf("namespace %s not found", namespaceID.String())
	}

	namespace := (*namespaces)[0]
	actionMap := make(map[string]int)
	if namespace.Actions != nil {
		for _, action := range *namespace.Actions {
			if action.Name != nil && action.Bit != nil {
				actionMap[*action.Name] = *action.Bit
			}
		}
	}

	// Calculate allow and deny bits
	allowBits := 0
	denyBits := 0
	notSetBits := 0

	for permName, permValue := range permissions {
		bit, ok := actionMap[permName]
		if !ok {
			return fmt.Errorf("permission '%s' not found in namespace %s", permName, namespaceID.String())
		}

		permValueStr := strings.ToLower(permValue.(string))
		switch permValueStr {
		case "allow":
			allowBits |= bit
		case "deny":
			denyBits |= bit
		case "notset":
			notSetBits |= bit
		default:
			return fmt.Errorf("invalid permission value '%s' for permission '%s'. Must be allow, deny, or notset", permValue.(string), permName)
		}
	}

	// Build ACE (Access Control Entry)
	ace := security.AccessControlEntry{
		Descriptor: &principal,
		Allow:      &allowBits,
		Deny:       &denyBits,
		ExtendedInfo: &security.AceExtendedInformation{
			EffectiveAllow: &allowBits,
			EffectiveDeny:  &denyBits,
			InheritedAllow: new(int),
			InheritedDeny:  new(int),
		},
	}

	// Create container structure for SetAccessControlEntries
	bMerge := !replace
	container := struct {
		Token                *string                        `json:"token,omitempty"`
		Merge                *bool                          `json:"merge,omitempty"`
		AccessControlEntries *[]security.AccessControlEntry `json:"accessControlEntries,omitempty"`
	}{
		Token:                &token,
		Merge:                &bMerge,
		AccessControlEntries: &[]security.AccessControlEntry{ace},
	}

	// Set ACL
	_, err = clients.SecurityClient.SetAccessControlEntries(clients.Ctx, security.SetAccessControlEntriesArgs{
		SecurityNamespaceId: &namespaceID,
		Container:           container,
	})
	if err != nil {
		return fmt.Errorf("setting permissions: %v", err)
	}

	// Wait for permissions to propagate
	stateConf := &retry.StateChangeConf{
		Pending: []string{"Waiting"},
		Target:  []string{"Synced"},
		Refresh: func() (interface{}, string, error) {
			currentACL, err := clients.SecurityClient.QueryAccessControlLists(clients.Ctx, security.QueryAccessControlListsArgs{
				SecurityNamespaceId: &namespaceID,
				Token:               &token,
				Descriptors:         &principal,
				IncludeExtendedInfo: &[]bool{true}[0],
			})

			if err != nil {
				return nil, "", fmt.Errorf("reading permissions: %v", err)
			}

			if currentACL == nil || len(*currentACL) == 0 {
				return "Waiting", "Waiting", nil
			}

			acl := (*currentACL)[0]
			if acl.AcesDictionary == nil {
				return "Waiting", "Waiting", nil
			}

			aceEntry, ok := (*acl.AcesDictionary)[principal]
			if !ok {
				return "Waiting", "Waiting", nil
			}

			// Check if permissions match
			if aceEntry.Allow != nil && *aceEntry.Allow == allowBits &&
				aceEntry.Deny != nil && *aceEntry.Deny == denyBits {
				return "Synced", "Synced", nil
			}

			return "Waiting", "Waiting", nil
		},
		Timeout:                   5 * time.Minute,
		MinTimeout:                2 * time.Second,
		Delay:                     2 * time.Second,
		ContinuousTargetOccurence: 1,
	}

	if _, err := stateConf.WaitForState(); err != nil { //nolint:staticcheck
		return fmt.Errorf("waiting for permission update: %v", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", namespaceID.String(), token, principal))
	return resourceGenericPermissionsRead(d, m)
}

func resourceGenericPermissionsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	namespaceID, err := uuid.Parse(d.Get("namespace_id").(string))
	if err != nil {
		return fmt.Errorf("invalid namespace_id: %v", err)
	}

	token := d.Get("token").(string)
	principal := d.Get("principal").(string)
	requestedPermissions := d.Get("permissions").(map[string]interface{})

	// Get namespace details
	namespaces, err := clients.SecurityClient.QuerySecurityNamespaces(clients.Ctx, security.QuerySecurityNamespacesArgs{
		SecurityNamespaceId: &namespaceID,
	})
	if err != nil {
		return fmt.Errorf("querying security namespace: %v", err)
	}
	if namespaces == nil || len(*namespaces) == 0 {
		d.SetId("")
		log.Printf("[INFO] Namespace %s not found. Removing from state", namespaceID.String())
		return nil
	}

	namespace := (*namespaces)[0]
	actionMap := make(map[string]int)
	if namespace.Actions != nil {
		for _, action := range *namespace.Actions {
			if action.Name != nil && action.Bit != nil {
				actionMap[*action.Name] = *action.Bit
			}
		}
	}

	// Query current ACL
	bTrue := true
	acls, err := clients.SecurityClient.QueryAccessControlLists(clients.Ctx, security.QueryAccessControlListsArgs{
		SecurityNamespaceId: &namespaceID,
		Token:               &token,
		Descriptors:         &principal,
		IncludeExtendedInfo: &bTrue,
	})

	if err != nil {
		return fmt.Errorf("querying ACL: %v", err)
	}

	if acls == nil || len(*acls) == 0 {
		d.SetId("")
		log.Printf("[INFO] Permissions for token %s not found. Removing from state", token)
		return nil
	}

	acl := (*acls)[0]
	if acl.AcesDictionary == nil {
		d.SetId("")
		log.Printf("[INFO] No ACEs found for principal %s. Removing from state", principal)
		return nil
	}

	ace, ok := (*acl.AcesDictionary)[principal]
	if !ok {
		d.SetId("")
		log.Printf("[INFO] ACE for principal %s not found. Removing from state", principal)
		return nil
	}

	// Convert bits back to permission map
	currentPermissions := make(map[string]interface{})
	allowBits := 0
	denyBits := 0
	if ace.Allow != nil {
		allowBits = *ace.Allow
	}
	if ace.Deny != nil {
		denyBits = *ace.Deny
	}

	// Validate all requested permissions exist in namespace
	for permName := range requestedPermissions {
		bit, ok := actionMap[permName]
		if !ok {
			return fmt.Errorf("permission '%s' not found in namespace %s", permName, namespaceID.String())
		}

		if (allowBits & bit) != 0 {
			currentPermissions[permName] = "allow"
		} else if (denyBits & bit) != 0 {
			currentPermissions[permName] = "deny"
		} else {
			currentPermissions[permName] = "notset"
		}
	}

	if err := d.Set("permissions", currentPermissions); err != nil {
		return fmt.Errorf("setting permissions in state: %v", err)
	}
	return nil
}

func resourceGenericPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	namespaceID, err := uuid.Parse(d.Get("namespace_id").(string))
	if err != nil {
		return fmt.Errorf("invalid namespace_id: %v", err)
	}

	token := d.Get("token").(string)
	principal := d.Get("principal").(string)
	permissions := d.Get("permissions").(map[string]interface{})

	// Get namespace details
	namespaces, err := clients.SecurityClient.QuerySecurityNamespaces(clients.Ctx, security.QuerySecurityNamespacesArgs{
		SecurityNamespaceId: &namespaceID,
	})
	if err != nil {
		return fmt.Errorf("querying security namespace: %v", err)
	}
	if namespaces == nil || len(*namespaces) == 0 {
		return fmt.Errorf("namespace %s not found", namespaceID.String())
	}

	namespace := (*namespaces)[0]
	actionMap := make(map[string]int)
	if namespace.Actions != nil {
		for _, action := range *namespace.Actions {
			if action.Name != nil && action.Bit != nil {
				actionMap[*action.Name] = *action.Bit
			}
		}
	}

	// Validate all permission names before proceeding
	for permName := range permissions {
		if _, ok := actionMap[permName]; !ok {
			return fmt.Errorf("permission '%s' not found in namespace %s", permName, namespaceID.String())
		}
	}

	// Read current ACL to get existing permissions
	bTrue := true
	acls, err := clients.SecurityClient.QueryAccessControlLists(clients.Ctx, security.QueryAccessControlListsArgs{
		SecurityNamespaceId: &namespaceID,
		Token:               &token,
		Descriptors:         &principal,
		IncludeExtendedInfo: &bTrue,
	})

	if err != nil {
		return fmt.Errorf("querying current ACL: %v", err)
	}

	// If no ACL exists, there's nothing to delete
	if acls == nil || len(*acls) == 0 {
		log.Printf("[INFO] No ACL found for token %s, nothing to delete", token)
		return nil
	}

	acl := (*acls)[0]
	if acl.AcesDictionary == nil {
		log.Printf("[INFO] No ACEs found for principal %s, nothing to delete", principal)
		return nil
	}

	ace, ok := (*acl.AcesDictionary)[principal]
	if !ok {
		log.Printf("[INFO] ACE for principal %s not found, nothing to delete", principal)
		return nil
	}

	// Get current allow and deny bits
	currentAllowBits := 0
	currentDenyBits := 0
	if ace.Allow != nil {
		currentAllowBits = *ace.Allow
	}
	if ace.Deny != nil {
		currentDenyBits = *ace.Deny
	}

	// Calculate bits managed by this state file
	managedBits := 0
	for permName := range permissions {
		if bit, ok := actionMap[permName]; ok {
			managedBits |= bit
		}
	}

	// Remove only the managed bits from current permissions
	newAllowBits := currentAllowBits &^ managedBits // Clear managed bits from allow
	newDenyBits := currentDenyBits &^ managedBits   // Clear managed bits from deny

	// Set the updated permissions
	updatedACE := security.AccessControlEntry{
		Descriptor: &principal,
		Allow:      &newAllowBits,
		Deny:       &newDenyBits,
		ExtendedInfo: &security.AceExtendedInformation{
			EffectiveAllow: &newAllowBits,
			EffectiveDeny:  &newDenyBits,
			InheritedAllow: new(int),
			InheritedDeny:  new(int),
		},
	}

	// Create container structure for SetAccessControlEntries
	bMerge := false
	container := struct {
		Token                *string                        `json:"token,omitempty"`
		Merge                *bool                          `json:"merge,omitempty"`
		AccessControlEntries *[]security.AccessControlEntry `json:"accessControlEntries,omitempty"`
	}{
		Token:                &token,
		Merge:                &bMerge,
		AccessControlEntries: &[]security.AccessControlEntry{updatedACE},
	}

	_, err = clients.SecurityClient.SetAccessControlEntries(clients.Ctx, security.SetAccessControlEntriesArgs{
		SecurityNamespaceId: &namespaceID,
		Container:           container,
	})

	if err != nil {
		return fmt.Errorf("removing managed permissions: %v", err)
	}

	log.Printf("[INFO] Successfully removed managed permissions from ACE for principal %s", principal)
	return nil
}
