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
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
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
					ValidateFunc: validation.StringInSlice([]string{"allow", "deny", "notset"}, false),
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

	// Resolve the subject descriptor to an identity descriptor
	identityDescriptor, err := resolveIdentityDescriptor(clients, principal)
	if err != nil {
		return fmt.Errorf("resolving identity for principal '%s': %v", principal, err)
	}

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

	// Build action maps and resolve display names to names
	actionMap, err := buildActionMap(namespace, permissions)
	if err != nil {
		return fmt.Errorf("building action map: %v", err)
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
		Descriptor: &identityDescriptor,
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
				Descriptors:         &identityDescriptor,
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

			aceEntry, ok := (*acl.AcesDictionary)[identityDescriptor]
			if !ok {
				return "Waiting", "Waiting", nil
			}

			currentAllow := 0
			currentDeny := 0
			if aceEntry.Allow != nil {
				currentAllow = *aceEntry.Allow
			}
			if aceEntry.Deny != nil {
				currentDeny = *aceEntry.Deny
			}

			// Check if permissions match
			if replace {
				// For replace mode, require exact match
				if currentAllow == allowBits && currentDeny == denyBits {
					return "Synced", "Synced", nil
				}
			} else {
				// For merge mode, check if desired bits are set and notset bits are cleared
				allowMatch := (currentAllow & allowBits) == allowBits
				denyMatch := (currentDeny & denyBits) == denyBits
				notSetCleared := (currentAllow&notSetBits) == 0 && (currentDeny&notSetBits) == 0

				if allowMatch && denyMatch && notSetCleared {
					return "Synced", "Synced", nil
				}
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

	// Resolve the subject descriptor to an identity descriptor
	identityDescriptor, err := resolveIdentityDescriptor(clients, principal)
	if err != nil {
		d.SetId("")
		log.Printf("[INFO] Unable to resolve identity for principal %s. Removing from state: %v", principal, err)
		return nil
	}

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

	// Build action maps and resolve display names to names
	actionMap, err := buildActionMap(namespace, requestedPermissions)
	if err != nil {
		return fmt.Errorf("building action map: %v", err)
	}

	// Query current ACL
	bTrue := true
	acls, err := clients.SecurityClient.QueryAccessControlLists(clients.Ctx, security.QueryAccessControlListsArgs{
		SecurityNamespaceId: &namespaceID,
		Token:               &token,
		Descriptors:         &identityDescriptor,
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

	ace, ok := (*acl.AcesDictionary)[identityDescriptor]
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

		switch {
		case (allowBits & bit) != 0:
			currentPermissions[permName] = "allow"
		case (denyBits & bit) != 0:
			currentPermissions[permName] = "deny"
		default:
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

	// Resolve the subject descriptor to an identity descriptor
	identityDescriptor, err := resolveIdentityDescriptor(clients, principal)
	if err != nil {
		// If we can't resolve the identity, the principal may no longer exist, so nothing to delete
		log.Printf("[INFO] Unable to resolve identity for principal %s during delete, assuming already removed: %v", principal, err)
		return nil
	}

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

	// Build action maps and resolve display names to names
	actionMap, err := buildActionMap(namespace, permissions)
	if err != nil {
		return fmt.Errorf("building action map: %v", err)
	}

	// Read current ACL to get existing permissions
	bTrue := true
	acls, err := clients.SecurityClient.QueryAccessControlLists(clients.Ctx, security.QueryAccessControlListsArgs{
		SecurityNamespaceId: &namespaceID,
		Token:               &token,
		Descriptors:         &identityDescriptor,
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

	ace, ok := (*acl.AcesDictionary)[identityDescriptor]
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
		Descriptor: &identityDescriptor,
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

// buildActionMap creates a mapping from permission names (and DisplayNames) to their bit values.
// It handles the following:
// 1. Maps action Names to bit values
// 2. Maps action DisplayNames to bit values (if DisplayName is unique)
// 3. Validates that no permission is specified twice (once via Name, once via DisplayName)
// 4. Logs ambiguous DisplayNames that map to multiple Names (those are ignored for DisplayName resolution)
func buildActionMap(namespace security.SecurityNamespaceDescription, requestedPermissions map[string]interface{}) (map[string]int, error) {
	actionMap := make(map[string]int)
	nameToDisplayName := make(map[string]string)    // Maps Name -> DisplayName
	displayNameToNames := make(map[string][]string) // Maps DisplayName -> []Name (to detect collisions)

	// First pass: Build all mappings
	if namespace.Actions != nil {
		for _, action := range *namespace.Actions {
			if action.Name != nil && action.Bit != nil {
				name := *action.Name
				bit := *action.Bit

				// Map Name to bit
				actionMap[name] = bit

				// Track DisplayName if it exists
				if action.DisplayName != nil && *action.DisplayName != "" {
					displayName := *action.DisplayName
					nameToDisplayName[name] = displayName
					displayNameToNames[displayName] = append(displayNameToNames[displayName], name)
				}
			}
		}
	}

	// Second pass: Add DisplayName mappings only if they're unambiguous
	ambiguousDisplayNames := make(map[string]bool)
	for displayName, names := range displayNameToNames {
		if len(names) > 1 {
			// Multiple names have the same DisplayName - this is ambiguous
			ambiguousDisplayNames[displayName] = true
			log.Printf("[DEBUG] DisplayName '%s' maps to multiple action names: %v. Will not resolve this DisplayName automatically.",
				displayName, names)
		} else if len(names) == 1 {
			// Unique DisplayName - safe to add to actionMap
			name := names[0]
			actionMap[displayName] = actionMap[name]
		}
	}

	// Third pass: Validate that no permission is specified twice
	// Track which underlying action (by bit value) has been requested
	requestedActionBits := make(map[int][]string) // Maps bit -> []permissionKey (as provided by user)

	for permKey := range requestedPermissions {
		bit, ok := actionMap[permKey]
		if !ok {
			// Permission key not found - will be caught by caller
			continue
		}

		requestedActionBits[bit] = append(requestedActionBits[bit], permKey)
	}

	// Check for duplicates
	for bit, permKeys := range requestedActionBits {
		if len(permKeys) > 1 {
			return nil, fmt.Errorf("permission specified multiple times using different keys: %v (all refer to the same permission bit %d)",
				permKeys, bit)
		}
	}

	// Fourth pass: Validate that ambiguous DisplayNames are not used in requestedPermissions
	for permKey := range requestedPermissions {
		if ambiguousDisplayNames[permKey] {
			// User tried to use an ambiguous DisplayName
			names := displayNameToNames[permKey]
			return nil, fmt.Errorf("permission key '%s' is ambiguous - it matches DisplayName for multiple actions: %v. Please use the action Name instead",
				permKey, names)
		}
	}

	return actionMap, nil
}

// resolveIdentityDescriptor resolves a subject descriptor (e.g., vssgp.Uy0xLTkt...)
// to an identity descriptor (e.g., Microsoft.IdentityModel.Claims.ClaimsIdentity;...)
// which is required by the security API for setting ACEs.
func resolveIdentityDescriptor(clients *client.AggregatedClient, subjectDescriptor string) (string, error) {
	identities, err := clients.IdentityClient.ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
		SubjectDescriptors: &subjectDescriptor,
	})
	if err != nil {
		return "", fmt.Errorf("reading identity for subject descriptor: %v", err)
	}

	if identities == nil || len(*identities) == 0 {
		return "", fmt.Errorf("no identity found for subject descriptor '%s'", subjectDescriptor)
	}

	id := (*identities)[0]
	if id.Descriptor == nil {
		return "", fmt.Errorf("identity descriptor is nil for subject descriptor '%s'", subjectDescriptor)
	}

	// Check if identity is active
	if id.IsActive != nil && !*id.IsActive {
		return "", fmt.Errorf("identity for subject descriptor '%s' is not active", subjectDescriptor)
	}

	return *id.Descriptor, nil
}
