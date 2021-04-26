package utils

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/identity"
	"github.com/microsoft/azure-devops-go-api/azuredevops/security"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ActionName type for an permission actions
type ActionName string

// PermissionType type for a single permission
type PermissionType string

type permissionTypeValuesType struct {
	Deny   PermissionType
	Allow  PermissionType
	NotSet PermissionType
}

// PermissionTypeValues valid permission action values
var PermissionTypeValues = permissionTypeValuesType{
	Deny:   "deny",
	Allow:  "allow",
	NotSet: "notset",
}

// SecurityNamespaceID the type of a security namespace id
type SecurityNamespaceID uuid.UUID

type securityNamespaceIDValuesType struct {
	Analytics                      SecurityNamespaceID
	AnalyticsViews                 SecurityNamespaceID
	ReleaseManagement              SecurityNamespaceID
	ReleaseManagement2             SecurityNamespaceID
	AuditLog                       SecurityNamespaceID
	Identity                       SecurityNamespaceID
	WorkItemTrackingAdministration SecurityNamespaceID
	DistributedTask                SecurityNamespaceID
	GitRepositories                SecurityNamespaceID
	VersionControlItems2           SecurityNamespaceID
	EventSubscriber                SecurityNamespaceID
	WorkItemTrackingProvision      SecurityNamespaceID
	ServiceEndpoints               SecurityNamespaceID
	ServiceHooks                   SecurityNamespaceID
	Collection                     SecurityNamespaceID
	Proxy                          SecurityNamespaceID
	Plan                           SecurityNamespaceID
	Process                        SecurityNamespaceID
	AccountAdminSecurity           SecurityNamespaceID
	Library                        SecurityNamespaceID
	Environment                    SecurityNamespaceID
	Project                        SecurityNamespaceID
	EventSubscription              SecurityNamespaceID
	CSS                            SecurityNamespaceID
	TeamLabSecurity                SecurityNamespaceID
	ProjectAnalysisLanguageMetrics SecurityNamespaceID
	Tagging                        SecurityNamespaceID
	MetaTask                       SecurityNamespaceID
	Iteration                      SecurityNamespaceID
	WorkItemQueryFolders           SecurityNamespaceID
	Favorites                      SecurityNamespaceID
	Registry                       SecurityNamespaceID
	Graph                          SecurityNamespaceID
	ViewActivityPaneSecurity       SecurityNamespaceID
	Job                            SecurityNamespaceID
	WorkItemTracking               SecurityNamespaceID
	StrongBox                      SecurityNamespaceID
	Server                         SecurityNamespaceID
	TestManagement                 SecurityNamespaceID
	SettingEntries                 SecurityNamespaceID
	BuildAdministration            SecurityNamespaceID
	Location                       SecurityNamespaceID
	Boards                         SecurityNamespaceID
	UtilizationPermissions         SecurityNamespaceID
	WorkItemsHub                   SecurityNamespaceID
	WebPlatform                    SecurityNamespaceID
	VersionControlPrivileges       SecurityNamespaceID
	Workspaces                     SecurityNamespaceID
	CrossProjectWidgetView         SecurityNamespaceID
	WorkItemTrackingConfiguration  SecurityNamespaceID
	DiscussionThreads              SecurityNamespaceID
	BoardsExternalIntegration      SecurityNamespaceID
	DataProvider                   SecurityNamespaceID
	Social                         SecurityNamespaceID
	Security                       SecurityNamespaceID
	IdentityPicker                 SecurityNamespaceID
	ServicingOrchestration         SecurityNamespaceID
	Build                          SecurityNamespaceID
	DashboardsPrivileges           SecurityNamespaceID
	VersionControlItems            SecurityNamespaceID
}

// SecurityNamespaceIDValues contains all available security namespaces
var SecurityNamespaceIDValues = securityNamespaceIDValuesType{
	Analytics:                      SecurityNamespaceID(uuid.MustParse("58450c49-b02d-465a-ab12-59ae512d6531")),
	AnalyticsViews:                 SecurityNamespaceID(uuid.MustParse("d34d3680-dfe5-4cc6-a949-7d9c68f73cba")),
	ReleaseManagement:              SecurityNamespaceID(uuid.MustParse("7c7d32f7-0e86-4cd6-892e-b35dbba870bd")),
	ReleaseManagement2:             SecurityNamespaceID(uuid.MustParse("c788c23e-1b46-4162-8f5e-d7585343b5de")),
	AuditLog:                       SecurityNamespaceID(uuid.MustParse("a6cc6381-a1ca-4b36-b3c1-4e65211e82b6")),
	Identity:                       SecurityNamespaceID(uuid.MustParse("5a27515b-ccd7-42c9-84f1-54c998f03866")),
	WorkItemTrackingAdministration: SecurityNamespaceID(uuid.MustParse("445d2788-c5fb-4132-bbef-09c4045ad93f")),
	DistributedTask:                SecurityNamespaceID(uuid.MustParse("101eae8c-1709-47f9-b228-0e476c35b3ba")),
	GitRepositories:                SecurityNamespaceID(uuid.MustParse("2e9eb7ed-3c0a-47d4-87c1-0ffdd275fd87")),
	VersionControlItems2:           SecurityNamespaceID(uuid.MustParse("3c15a8b7-af1a-45c2-aa97-2cb97078332e")),
	EventSubscriber:                SecurityNamespaceID(uuid.MustParse("2bf24a2b-70ba-43d3-ad97-3d9e1f75622f")),
	WorkItemTrackingProvision:      SecurityNamespaceID(uuid.MustParse("5a6cd233-6615-414d-9393-48dbb252bd23")),
	ServiceEndpoints:               SecurityNamespaceID(uuid.MustParse("49b48001-ca20-4adc-8111-5b60c903a50c")),
	ServiceHooks:                   SecurityNamespaceID(uuid.MustParse("cb594ebe-87dd-4fc9-ac2c-6a10a4c92046")),
	Collection:                     SecurityNamespaceID(uuid.MustParse("3e65f728-f8bc-4ecd-8764-7e378b19bfa7")),
	Proxy:                          SecurityNamespaceID(uuid.MustParse("cb4d56d2-e84b-457e-8845-81320a133fbb")),
	Plan:                           SecurityNamespaceID(uuid.MustParse("bed337f8-e5f3-4fb9-80da-81e17d06e7a8")),
	Process:                        SecurityNamespaceID(uuid.MustParse("2dab47f9-bd70-49ed-9bd5-8eb051e59c02")),
	AccountAdminSecurity:           SecurityNamespaceID(uuid.MustParse("11238e09-49f2-40c7-94d0-8f0307204ce4")),
	Library:                        SecurityNamespaceID(uuid.MustParse("b7e84409-6553-448a-bbb2-af228e07cbeb")),
	Environment:                    SecurityNamespaceID(uuid.MustParse("83d4c2e6-e57d-4d6e-892b-b87222b7ad20")),
	Project:                        SecurityNamespaceID(uuid.MustParse("52d39943-cb85-4d7f-8fa8-c6baac873819")),
	EventSubscription:              SecurityNamespaceID(uuid.MustParse("58b176e7-3411-457a-89d0-c6d0ccb3c52b")),
	CSS:                            SecurityNamespaceID(uuid.MustParse("83e28ad4-2d72-4ceb-97b0-c7726d5502c3")),
	TeamLabSecurity:                SecurityNamespaceID(uuid.MustParse("9e4894c3-ff9a-4eac-8a85-ce11cafdc6f1")),
	ProjectAnalysisLanguageMetrics: SecurityNamespaceID(uuid.MustParse("fc5b7b85-5d6b-41eb-8534-e128cb10eb67")),
	Tagging:                        SecurityNamespaceID(uuid.MustParse("bb50f182-8e5e-40b8-bc21-e8752a1e7ae2")),
	MetaTask:                       SecurityNamespaceID(uuid.MustParse("f6a4de49-dbe2-4704-86dc-f8ec1a294436")),
	Iteration:                      SecurityNamespaceID(uuid.MustParse("bf7bfa03-b2b7-47db-8113-fa2e002cc5b1")),
	WorkItemQueryFolders:           SecurityNamespaceID(uuid.MustParse("71356614-aad7-4757-8f2c-0fb3bff6f680")),
	Favorites:                      SecurityNamespaceID(uuid.MustParse("fa557b48-b5bf-458a-bb2b-1b680426fe8b")),
	Registry:                       SecurityNamespaceID(uuid.MustParse("4ae0db5d-8437-4ee8-a18b-1f6fb38bd34c")),
	Graph:                          SecurityNamespaceID(uuid.MustParse("c2ee56c9-e8fa-4cdd-9d48-2c44f697a58e")),
	ViewActivityPaneSecurity:       SecurityNamespaceID(uuid.MustParse("dc02bf3d-cd48-46c3-8a41-345094ecc94b")),
	Job:                            SecurityNamespaceID(uuid.MustParse("2a887f97-db68-4b7c-9ae3-5cebd7add999")),
	WorkItemTracking:               SecurityNamespaceID(uuid.MustParse("73e71c45-d483-40d5-bdba-62fd076f7f87")),
	StrongBox:                      SecurityNamespaceID(uuid.MustParse("4a9e8381-289a-4dfd-8460-69028eaa93b3")),
	Server:                         SecurityNamespaceID(uuid.MustParse("1f4179b3-6bac-4d01-b421-71ea09171400")),
	TestManagement:                 SecurityNamespaceID(uuid.MustParse("e06e1c24-e93d-4e4a-908a-7d951187b483")),
	SettingEntries:                 SecurityNamespaceID(uuid.MustParse("6ec4592e-048c-434e-8e6c-8671753a8418")),
	BuildAdministration:            SecurityNamespaceID(uuid.MustParse("302acaca-b667-436d-a946-87133492041c")),
	Location:                       SecurityNamespaceID(uuid.MustParse("2725d2bc-7520-4af4-b0e3-8d876494731f")),
	Boards:                         SecurityNamespaceID(uuid.MustParse("251e12d9-bea3-43a8-bfdb-901b98c0125e")),
	UtilizationPermissions:         SecurityNamespaceID(uuid.MustParse("83abde3a-4593-424e-b45f-9898af99034d")),
	WorkItemsHub:                   SecurityNamespaceID(uuid.MustParse("c0e7a722-1cad-4ae6-b340-a8467501e7ce")),
	WebPlatform:                    SecurityNamespaceID(uuid.MustParse("0582eb05-c896-449a-b933-aa3d99e121d6")),
	VersionControlPrivileges:       SecurityNamespaceID(uuid.MustParse("66312704-deb5-43f9-b51c-ab4ff5e351c3")),
	Workspaces:                     SecurityNamespaceID(uuid.MustParse("93bafc04-9075-403a-9367-b7164eac6b5c")),
	CrossProjectWidgetView:         SecurityNamespaceID(uuid.MustParse("093cbb02-722b-4ad6-9f88-bc452043fa63")),
	WorkItemTrackingConfiguration:  SecurityNamespaceID(uuid.MustParse("35e35e8e-686d-4b01-aff6-c369d6e36ce0")),
	DiscussionThreads:              SecurityNamespaceID(uuid.MustParse("0d140cae-8ac1-4f48-b6d1-c93ce0301a12")),
	BoardsExternalIntegration:      SecurityNamespaceID(uuid.MustParse("5ab15bc8-4ea1-d0f3-8344-cab8fe976877")),
	DataProvider:                   SecurityNamespaceID(uuid.MustParse("7ffa7cf4-317c-4fea-8f1d-cfda50cfa956")),
	Social:                         SecurityNamespaceID(uuid.MustParse("81c27cc8-7a9f-48ee-b63f-df1e1d0412dd")),
	Security:                       SecurityNamespaceID(uuid.MustParse("9a82c708-bfbe-4f31-984c-e860c2196781")),
	IdentityPicker:                 SecurityNamespaceID(uuid.MustParse("a60e0d84-c2f8-48e4-9c0c-f32da48d5fd1")),
	ServicingOrchestration:         SecurityNamespaceID(uuid.MustParse("84cc1aa4-15bc-423d-90d9-f97c450fc729")),
	Build:                          SecurityNamespaceID(uuid.MustParse("33344d9c-fc72-4d6f-aba5-fa317101a7e9")),
	DashboardsPrivileges:           SecurityNamespaceID(uuid.MustParse("8adf73b7-389a-4276-b638-fe1653f7efc7")),
	VersionControlItems:            SecurityNamespaceID(uuid.MustParse("a39371cf-0841-4c16-bbd3-276e341bc052")),
}

// PrincipalPermission describes permissions of a principal
type PrincipalPermission struct {
	SubjectDescriptor string
	Permissions       map[ActionName]PermissionType
}

// SetPrincipalPermission sets permissions for a principal
type SetPrincipalPermission struct {
	Replace             bool
	PrincipalPermission PrincipalPermission
}

// SecurityNamespace an Azure DevOps Security Namespace
type SecurityNamespace struct {
	namespaceID    uuid.UUID
	context        context.Context
	securityClient security.Client
	identityClient identity.Client
	actions        *map[string]security.ActionDefinition
	token          string
}

// TokenCreatorFunc signature for creating namespace tokens
type TokenCreatorFunc func(d *schema.ResourceData, clients *client.AggregatedClient) (string, error)

// NewSecurityNamespace Creates a new instance of a security namespace
func NewSecurityNamespace(d *schema.ResourceData, clients *client.AggregatedClient, namespaceID SecurityNamespaceID, tokenCreator TokenCreatorFunc) (*SecurityNamespace, error) {
	if nil == clients.Ctx {
		return nil, fmt.Errorf("context is nil")
	}
	if nil == clients.SecurityClient {
		return nil, fmt.Errorf("securityClient is nil")
	}
	if nil == clients.IdentityClient {
		return nil, fmt.Errorf("identityClient is nil")
	}
	sn := new(SecurityNamespace)
	sn.context = clients.Ctx
	sn.namespaceID = uuid.UUID(namespaceID)
	sn.securityClient = clients.SecurityClient
	sn.identityClient = clients.IdentityClient
	token, err := tokenCreator(d, clients)
	if err != nil {
		return nil, err
	}
	sn.token = token
	return sn, nil
}

// GetToken return namespace tokens
func (sn *SecurityNamespace) GetToken() string {
	return sn.token
}

func (sn *SecurityNamespace) getActionDefinitions() (*map[string]security.ActionDefinition, error) {
	if sn.actions == nil {
		secns, err := sn.securityClient.QuerySecurityNamespaces(sn.context, security.QuerySecurityNamespacesArgs{
			SecurityNamespaceId: &sn.namespaceID,
		})
		if err != nil {
			return nil, err
		}
		if secns == nil || len(*secns) <= 0 || (*secns)[0].Actions == nil || len(*(*secns)[0].Actions) <= 0 {
			return nil, fmt.Errorf("Failed to load security namespace definition with id [%s]", sn.namespaceID)
		}

		actionMap := map[string]security.ActionDefinition{}
		for _, action := range *(*secns)[0].Actions {
			actionMap[*action.Name] = action
		}
		sn.actions = &actionMap
	}
	return sn.actions, nil
}

func (sn *SecurityNamespace) GetAccessControlList(descriptorList *[]string) (*security.AccessControlList, error) {
	var descriptors *string = nil
	if descriptorList != nil && len(*descriptorList) > 0 {
		val := linq.From(*descriptorList).
			Aggregate(func(r interface{}, i interface{}) interface{} {
				if r.(string) == "" {
					return i
				}
				return r.(string) + "," + i.(string)
			}).(string)
		descriptors = &val
	}

	bTrue := true
	acl, err := sn.securityClient.QueryAccessControlLists(sn.context, security.QueryAccessControlListsArgs{
		SecurityNamespaceId: &sn.namespaceID,
		Token:               &sn.token,
		Descriptors:         descriptors,
		IncludeExtendedInfo: &bTrue,
	})

	if err != nil {
		return nil, err
	}
	if acl == nil || len(*acl) <= 0 {
		return nil, nil
	}
	if len(*acl) != 1 {
		return nil, fmt.Errorf("Failed to load current ACL for token [%s]. Result set contains more than one ACL", sn.token)
	}
	return &(*acl)[0], nil
}

func (sn *SecurityNamespace) getIdentitiesFromSubjects(principal *[]string) (*[]identity.Identity, error) {
	if principal == nil || len(*principal) <= 0 {
		return nil, fmt.Errorf("principal is nil or empty")
	}

	descriptors := linq.From(*principal).
		Aggregate(func(r interface{}, i interface{}) interface{} {
			if r.(string) == "" {
				return i
			}
			return r.(string) + "," + i.(string)
		}).(string)

	idlist, err := sn.identityClient.ReadIdentities(sn.context, identity.ReadIdentitiesArgs{
		SubjectDescriptors: converter.String(descriptors),
	})

	if err != nil {
		return nil, err
	}
	if idlist == nil || len(*idlist) != len(*principal) {
		return nil, fmt.Errorf("Failed to load identity information for defined principals [%s]", descriptors)
	}
	return idlist, nil
}

// SetPrincipalPermissions sets ACLs for specifc token inside a security namespace
func (sn *SecurityNamespace) SetPrincipalPermissions(permissionList *[]SetPrincipalPermission) error {
	if nil == permissionList || len(*permissionList) <= 0 {
		return nil
	}

	permissionMap := map[string]SetPrincipalPermission{}
	linq.From(*permissionList).
		ToMapBy(&permissionMap,
			func(item interface{}) interface{} {
				return item.(SetPrincipalPermission).PrincipalPermission.SubjectDescriptor
			},
			func(item interface{}) interface{} { return item })

	subjectList := make([]string, len(permissionMap))
	linq.From(*permissionList).
		Select(func(item interface{}) interface{} {
			return item.(SetPrincipalPermission).PrincipalPermission.SubjectDescriptor
		}).
		ToSlice(&subjectList)

	idList, err := sn.getIdentitiesFromSubjects(&subjectList)
	if err != nil {
		return err
	}
	idMap := map[string]identity.Identity{}
	linq.From(*idList).
		ToMapBy(&idMap,
			func(item interface{}) interface{} { return *item.(identity.Identity).SubjectDescriptor },
			func(item interface{}) interface{} { return item })

	var descriptorList []string
	linq.From(*idList).
		Select(func(elem interface{}) interface{} {
			return *elem.(identity.Identity).Descriptor
		}).
		ToSlice(&descriptorList)

	acl, err := sn.GetAccessControlList(&descriptorList)
	if err != nil {
		return err
	}
	if acl == nil {
		return fmt.Errorf("Failed to load ACL for token %q", sn.token)
	}
	aceMap := *acl.AcesDictionary

	actionMap, err := sn.getActionDefinitions()
	if err != nil {
		return err
	}

	for subjectDescriptor, principalPermissions := range permissionMap {
		desc, ok := idMap[subjectDescriptor]
		if !ok {
			return fmt.Errorf("Unable to resolve id descriptor for principal [%s]", subjectDescriptor)
		}

		log.Printf("[TRACE] Checking ACE list for descriptor [%s]", subjectDescriptor)
		var aceItem *security.AccessControlEntry
		ace, update := aceMap[*desc.Descriptor]
		if !update {
			log.Printf("[TRACE] Creating new ACE for subject [%s]", subjectDescriptor)
			aceItem = new(security.AccessControlEntry)
			aceItem.Allow = new(int)
			aceItem.Deny = new(int)
			aceItem.Descriptor = desc.Descriptor
		} else {
			// update existing ACE for principal
			log.Printf("[TRACE] Updating ACE for descriptor [%s]", *desc.Descriptor)
			aceItem = &ace
		}

		for key, value := range principalPermissions.PrincipalPermission.Permissions {
			actionDef, ok := (*actionMap)[string(key)]
			if !ok {
				return fmt.Errorf("Invalid permission [%s]", key)
			}
			if aceItem.Deny == nil {
				aceItem.Deny = new(int)
			}
			if aceItem.Allow == nil {
				aceItem.Allow = new(int)
			}

			if strings.EqualFold("deny", string(value)) {
				*aceItem.Allow = (*aceItem.Allow) &^ (*actionDef.Bit)
				*aceItem.Deny = (*aceItem.Deny) | (*actionDef.Bit)
			} else if strings.EqualFold("allow", string(value)) {
				*aceItem.Deny = (*aceItem.Deny) &^ (*actionDef.Bit)
				*aceItem.Allow = (*aceItem.Allow) | (*actionDef.Bit)
			} else if strings.EqualFold("notset", string(value)) {
				*aceItem.Allow = (*aceItem.Allow) &^ (*actionDef.Bit)
				*aceItem.Deny = (*aceItem.Deny) &^ (*actionDef.Bit)
			} else {
				return fmt.Errorf("Invalid permission action [%s]", value)
			}
		}

		bMerge := !principalPermissions.Replace
		container := struct {
			Token                *string                        `json:"token,omitempty"`
			Merge                *bool                          `json:"merge,omitempty"`
			AccessControlEntries *[]security.AccessControlEntry `json:"accessControlEntries,omitempty"`
		}{
			Token:                &sn.token,
			Merge:                &bMerge,
			AccessControlEntries: &[]security.AccessControlEntry{*aceItem},
		}

		_, err = sn.securityClient.SetAccessControlEntries(sn.context, security.SetAccessControlEntriesArgs{
			SecurityNamespaceId: &sn.namespaceID,
			Container:           container,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// GetPrincipalPermissions returns an array of PrincipalPermission for a Security Namespace token an a list of principals
func (sn *SecurityNamespace) GetPrincipalPermissions(principal *[]string) (*[]PrincipalPermission, error) {
	actions, err := sn.getActionDefinitions()
	if err != nil {
		return nil, err
	}

	idList, err := sn.getIdentitiesFromSubjects(principal)
	if err != nil {
		return nil, err
	}

	var descriptorList []string
	linq.From(*idList).
		SelectT(func(elem interface{}) string {
			return *elem.(identity.Identity).Descriptor
		}).
		ToSlice(&descriptorList)
	acl, err := sn.GetAccessControlList(&descriptorList)
	if err != nil {
		return nil, err
	}
	if acl == nil {
		return nil, nil
	}
	idMap := map[string]identity.Identity{}
	linq.From(*idList).
		ToMapBy(&idMap,
			func(item interface{}) interface{} { return *item.(identity.Identity).Descriptor },
			func(item interface{}) interface{} { return item })

	permissions := []PrincipalPermission{}
	for id, ace := range *acl.AcesDictionary {
		subject, ok := idMap[id]
		if !ok {
			return nil, fmt.Errorf("INTERAL ERROR: identity map does not contain an item with key [%s]", id)
		}
		if subject.SubjectDescriptor == nil {
			return nil, fmt.Errorf("Identity %s does not contain a subject descriptor value", id)
		}

		subjectPerm := PrincipalPermission{
			SubjectDescriptor: *(subject.SubjectDescriptor),
			Permissions:       map[ActionName]PermissionType{},
		}
		for actionName, actionDef := range *actions {
			if (*ace.Allow)&(*actionDef.Bit) != 0 {
				subjectPerm.Permissions[ActionName(actionName)] = PermissionTypeValues.Allow
			} else if (*ace.Deny)&(*actionDef.Bit) != 0 {
				subjectPerm.Permissions[ActionName(actionName)] = PermissionTypeValues.Deny
			} else {
				subjectPerm.Permissions[ActionName(actionName)] = PermissionTypeValues.NotSet
			}
		}
		permissions = append(permissions, subjectPerm)
	}
	return &permissions, nil
}

// RemovePrincipalPermissions removes all permissions for given principals and a Security Namespace token
func (sn *SecurityNamespace) RemovePrincipalPermissions(principal *[]string) error {
	idList, err := sn.getIdentitiesFromSubjects(principal)
	if err != nil {
		return err
	}
	acl, err := sn.GetAccessControlList(nil)
	if err != nil {
		return err
	}
	if acl == nil {
		return nil
	}

	val := linq.From(*idList).
		Where(func(i interface{}) bool {
			_, ok := (*acl.AcesDictionary)[*i.(identity.Identity).Descriptor]
			return ok
		}).
		Aggregate(func(r interface{}, i interface{}) interface{} {
			desc := *i.(identity.Identity).Descriptor
			if r.(string) == "" {
				return desc
			}
			return r.(string) + "," + desc
		}).(string)

	log.Printf("[TRACE]RemovePrincipalPermissions: removing the following principals from the ACL %s", val)
	bRet, err := sn.securityClient.RemoveAccessControlEntries(sn.context, security.RemoveAccessControlEntriesArgs{
		SecurityNamespaceId: &sn.namespaceID,
		Token:               &sn.token,
		Descriptors:         &val,
	})
	if err != nil {
		return err
	}
	if !(*bRet) {
		return fmt.Errorf("Failed to remove ACL entries for principals %s", val)
	}
	return nil
}
