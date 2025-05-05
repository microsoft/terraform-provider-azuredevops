package securityroles

import "github.com/google/uuid"

type SetRoleAssignmentPayload struct {
	UserID   *uuid.UUID `json:"userId"`
	RoleName *string    `json:"roleName"`
}

type SecurityRoleIdentity struct {
	DisplayName *string `json:"displayName"`
	ID          *string `json:"id"`
	UniqueName  *string `json:"uniqueName"`
}

type SecurityRoleDefinition struct {
	DisplayName      *string `json:"displayName"`
	Name             *string `json:"name"`
	AllowPermissions *int    `json:"allowPermissions"`
	DenyPermissions  *int    `json:"denyPermissions"`
	Identifier       *string `json:"identifier"`
	Description      *string `json:"description"`
	Scope            *string `json:"scope"`
}

type SecurityRoleAssignment struct {
	Identity          *SecurityRoleIdentity   `json:"identity"`
	Role              *SecurityRoleDefinition `json:"role"`
	Access            *string                 `json:"access"`
	AccessDisplayName *string                 `json:"accessDisplayName"`
}
