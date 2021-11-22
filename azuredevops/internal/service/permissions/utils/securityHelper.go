package utils

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// SetPrincipalPermissions sets permissions for a specific security namespac
func SetPrincipalPermissions(d *schema.ResourceData, sn *SecurityNamespace, forcePermission *PermissionType, forceReplace bool) error {
	principal, ok := d.GetOk("principal")
	if !ok {
		return fmt.Errorf("Failed to get 'principal' from schema")
	}

	permissions, ok := d.GetOk("permissions")
	if !ok {
		return fmt.Errorf("Failed to get 'permissions' from schema")
	}

	bReplace := d.Get("replace")
	if forceReplace {
		bReplace = forceReplace
	}
	permissionMap := make(map[ActionName]PermissionType, len(permissions.(map[string]interface{})))
	for key, elem := range permissions.(map[string]interface{}) {
		if forcePermission != nil {
			permissionMap[ActionName(key)] = *forcePermission
		} else {
			permissionMap[ActionName(key)] = PermissionType(elem.(string))
		}
	}
	setPermissions := []SetPrincipalPermission{
		{
			Replace: bReplace.(bool),
			PrincipalPermission: PrincipalPermission{
				SubjectDescriptor: principal.(string),
				Permissions:       permissionMap,
			},
		}}

	if err := sn.SetPrincipalPermissions(&setPermissions); err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"Waiting"},
		Target:  []string{"Synched"},
		Refresh: func() (interface{}, string, error) {
			state := "Waiting"
			currentPermissions, err := sn.GetPrincipalPermissions(&[]string{
				principal.(string),
			})
			if err != nil {
				return nil, "", fmt.Errorf("Error reading principal permissions: %+v", err)
			}

			bInsnyc := false
			for key := range ((*currentPermissions)[0]).Permissions {
				bInsnyc = permissionMap[key] == ((*currentPermissions)[0]).Permissions[key]
				if !bInsnyc {
					break
				}
			}
			if bInsnyc {
				state = "Synched"
			}

			return state, state, nil
		},
		Timeout:                   60 * time.Minute,
		MinTimeout:                5 * time.Second,
		Delay:                     5 * time.Second,
		ContinuousTargetOccurence: 1,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(" waiting for permission update. %v ", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", sn.token, principal.(string)))
	return nil
}

// GetPrincipalPermissions gets permissions for a specific security namespac
func GetPrincipalPermissions(d *schema.ResourceData, sn *SecurityNamespace) (*PrincipalPermission, error) {
	principal, ok := d.GetOk("principal")
	if !ok {
		return nil, fmt.Errorf("Failed to get 'principal' from schema")
	}

	permissions, ok := d.GetOk("permissions")
	if !ok {
		return nil, fmt.Errorf("Failed to get 'permissions' from schema")
	}

	principalList := []string{*converter.StringFromInterface(principal)}
	principalPermissions, err := sn.GetPrincipalPermissions(&principalList)
	if err != nil {
		return nil, err
	}
	if principalPermissions == nil || len(*principalPermissions) <= 0 {
		return nil, nil
	}
	if len(*principalPermissions) != 1 {
		return nil, fmt.Errorf("Failed to retrieve current permissions for principal [%s]", principalList[0])
	}
	for key := range ((*principalPermissions)[0]).Permissions {
		if _, ok := permissions.(map[string]interface{})[string(key)]; !ok {
			delete(((*principalPermissions)[0]).Permissions, key)
		}
	}
	return &(*principalPermissions)[0], nil
}
