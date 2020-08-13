// +build all utils securitynamespaces
// +build !exclude_securitynamespaces

package utils

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/identity"
	"github.com/microsoft/azure-devops-go-api/azuredevops/security"
	"github.com/stretchr/testify/assert"
	"github.com/terraform-providers/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/testhelper"
)

type isReadIdentitiesArgs struct{ t identity.ReadIdentitiesArgs }

func IsReadIdentitiesArgs(t identity.ReadIdentitiesArgs) gomock.Matcher {
	return &isReadIdentitiesArgs{t}
}

func (o *isReadIdentitiesArgs) Matches(x interface{}) bool {
	if reflect.TypeOf(x) != reflect.TypeOf(identity.ReadIdentitiesArgs{}) {
		return false
	}
	args := x.(identity.ReadIdentitiesArgs)
	if o.t.Descriptors == nil && args.Descriptors == nil {
		return true
	} else if (o.t.Descriptors == nil && args.Descriptors != nil) || (o.t.Descriptors != nil && args.Descriptors == nil) {
		return false
	}

	argsDescList := strings.Split(*args.Descriptors, ",")
	refDescList := strings.Split(*o.t.Descriptors, ",")
	eq := func(a, b []string) bool {
		if len(a) != len(b) {
			return false
		}
		sort.Strings(a)
		sort.Strings(b)
		for i, v := range a {
			if v != b[i] {
				return false
			}
		}
		return true
	}

	return eq(argsDescList, refDescList)
}

func (o *isReadIdentitiesArgs) String() string {
	return "Equals to an identity.ReadIdentitiesArgs instance"
}

var securityNamespaceDescriptionProjectId = uuid.UUID(SecurityNamespaceIDValues.Project)
var securityNamespaceDescriptionProjectEmpty = []security.SecurityNamespaceDescription{}
var securityNamespaceDescriptionProject = []security.SecurityNamespaceDescription{
	{
		Name:        converter.String("Project"),
		NamespaceId: &securityNamespaceDescriptionProjectId,
		Actions: &[]security.ActionDefinition{
			{
				Name:        converter.String("GENERIC_READ"),
				Bit:         converter.Int(1),
				DisplayName: converter.String("View project-level information"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("GENERIC_WRITE"),
				Bit:         converter.Int(2),
				DisplayName: converter.String("Edit project-level information"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("DELETE"),
				Bit:         converter.Int(4),
				DisplayName: converter.String("Delete team project"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("PUBLISH_TEST_RESULTS"),
				Bit:         converter.Int(8),
				DisplayName: converter.String("Create test runs"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("ADMINISTER_BUILD"),
				Bit:         converter.Int(16),
				DisplayName: converter.String("Administer a build"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("START_BUILD"),
				Bit:         converter.Int(32),
				DisplayName: converter.String("Start a build"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("EDIT_BUILD_STATUS"),
				Bit:         converter.Int(64),
				DisplayName: converter.String("Edit build quality"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("UPDATE_BUILD"),
				Bit:         converter.Int(128),
				DisplayName: converter.String("Write to build operational store"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("DELETE_TEST_RESULTS"),
				Bit:         converter.Int(256),
				DisplayName: converter.String("Delete test runs"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("VIEW_TEST_RESULTS"),
				Bit:         converter.Int(512),
				DisplayName: converter.String("View test runs"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("MANAGE_TEST_ENVIRONMENTS"),
				Bit:         converter.Int(2048),
				DisplayName: converter.String("Manage test environments"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("MANAGE_TEST_CONFIGURATIONS"),
				Bit:         converter.Int(4096),
				DisplayName: converter.String("Manage test configurations"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("WORK_ITEM_DELETE"),
				Bit:         converter.Int(8192),
				DisplayName: converter.String("Delete and restore work items"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("WORK_ITEM_MOVE"),
				Bit:         converter.Int(16384),
				DisplayName: converter.String("Move work items out of this project"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("WORK_ITEM_PERMANENTLY_DELETE"),
				Bit:         converter.Int(32768),
				DisplayName: converter.String("Permanently delete work items"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("RENAME"),
				Bit:         converter.Int(65536),
				DisplayName: converter.String("Rename team project"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("MANAGE_PROPERTIES"),
				Bit:         converter.Int(131072),
				DisplayName: converter.String("Manage project properties"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("MANAGE_SYSTEM_PROPERTIES"),
				Bit:         converter.Int(262144),
				DisplayName: converter.String("Manage system project properties"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("BYPASS_PROPERTY_CACHE"),
				Bit:         converter.Int(524288),
				DisplayName: converter.String("Bypass project property cache"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("BYPASS_RULES"),
				Bit:         converter.Int(1048576),
				DisplayName: converter.String("Bypass rules on work item updates"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("SUPPRESS_NOTIFICATIONS"),
				Bit:         converter.Int(2097152),
				DisplayName: converter.String("Suppress notifications for work item updates"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("UPDATE_VISIBILITY"),
				Bit:         converter.Int(4194304),
				DisplayName: converter.String("Update project visibility"),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("CHANGE_PROCESS"),
				Bit:         converter.Int(8388608),
				DisplayName: converter.String("Change process of team project."),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("AGILETOOLS_BACKLOG"),
				Bit:         converter.Int(16777216),
				DisplayName: converter.String("Agile backlog management."),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
			{
				Name:        converter.String("AGILETOOLS_PLANS"),
				Bit:         converter.Int(33554432),
				DisplayName: converter.String("Agile plans."),
				NamespaceId: &securityNamespaceDescriptionProjectId,
			},
		},
	},
}

var projectID = "9083e944-8e9e-405e-960a-c80180aa71e6"
var projectAccessToken = fmt.Sprintf("$PROJECT:vstfs:///Classification/TeamProject/%s", projectID)
var projectAccessControlListEmpty = []security.AccessControlList{}
var projectAccessControlList = []security.AccessControlList{{
	AcesDictionary: &map[string]security.AccessControlEntry{
		"Microsoft.TeamFoundation.ServiceIdentity;7774ac03-8a29-44ac-86f1-fa4bded78de2:Build:f609b046-3e4a-419a-a5d7-a0840414dc74": {
			Descriptor: converter.String("Microsoft.TeamFoundation.ServiceIdentity;7774ac03-8a29-44ac-86f1-fa4bded78de2:Build:f609b046-3e4a-419a-a5d7-a0840414dc74"),
			Allow:      converter.Int(4745),
			Deny:       converter.Int(0),
		},
		"Microsoft.TeamFoundation.Identity;S-1-9-1551374245-4251810032-2399672646-2899062471-1578266062-0-0-0-0-1": {
			Descriptor: converter.String("Microsoft.TeamFoundation.Identity;S-1-9-1551374245-4251810032-2399672646-2899062471-1578266062-0-0-0-0-1"),
			Allow:      converter.Int(112),
			Deny:       converter.Int(0),
		},
		"Microsoft.TeamFoundation.Identity;S-1-9-1551374245-4251810032-2399672646-2899062471-1578266062-0-0-0-1-1": {
			Descriptor: converter.String("Microsoft.TeamFoundation.Identity;S-1-9-1551374245-4251810032-2399672646-2899062471-1578266062-0-0-0-1-1"),
			Allow:      converter.Int(160),
			Deny:       converter.Int(0),
		},
		"Microsoft.TeamFoundation.Identity;S-1-9-1551374245-4251810032-2399672646-2899062471-1578266062-0-0-0-1-2": {
			Descriptor: converter.String("Microsoft.TeamFoundation.Identity;S-1-9-1551374245-4251810032-2399672646-2899062471-1578266062-0-0-0-1-2"),
			Allow:      converter.Int(112),
			Deny:       converter.Int(0),
		},
		"Microsoft.TeamFoundation.Identity;S-1-9-1551374245-4251810032-2399672646-2899062471-1578266062-0-0-0-4-1": {
			Descriptor: converter.String("Microsoft.TeamFoundation.Identity;S-1-9-1551374245-4251810032-2399672646-2899062471-1578266062-0-0-0-4-1"),
			Allow:      converter.Int(521),
			Deny:       converter.Int(0),
		},
	},
	Token: &projectAccessToken,
},
}

var projectIdentityListEmpty = []identity.Identity{}
var projectIdentityList = []identity.Identity{
	{
		CustomDisplayName:   converter.String("Df609b046-3e4a-419a-a5d7-a0840414dc74 Build Service (ophiosdev)"),
		Descriptor:          converter.String("Microsoft.TeamFoundation.ServiceIdentity;7774ac03-8a29-44ac-86f1-fa4bded78de2:Build:f609b046-3e4a-419a-a5d7-a0840414dc74"),
		Id:                  testhelper.ToUUID("79b8298b-7101-4a53-ad6d-1d3de0b495f1"),
		ProviderDisplayName: converter.String("f609b046-3e4a-419a-a5d7-a0840414dc74"),
		SubjectDescriptor:   converter.String("svc.Nzc3NGFjMDMtOGEyOS00NGFjLTg2ZjEtZmE0YmRlZDc4ZGUyOkJ1aWxkOmY2MDliMDQ2LTNlNGEtNDE5YS1hNWQ3LWEwODQwNDE0ZGM3NA"),
	},
	{
		Descriptor:          converter.String("Microsoft.TeamFoundation.Identity;S-1-9-1551374245-4251810032-2399672646-2899062471-1578266062-0-0-0-0-1"),
		Id:                  testhelper.ToUUID("b555cec9-60cf-4f6e-9626-670f964945c5"),
		ProviderDisplayName: converter.String("[dev]\\Project Collection Administrators"),
		SubjectDescriptor:   converter.String("vssgp.Uy0xLTktMTU1MTM3NDI0NS00MjUxODEwMDMyLTIzOTk2NzI2NDYtMjg5OTA2MjQ3MS0xNTc4MjY2MDYyLTAtMC0wLTAtMQ"),
	},
	{
		Descriptor:          converter.String("Microsoft.TeamFoundation.Identity;S-1-9-1551374245-4251810032-2399672646-2899062471-1578266062-0-0-0-1-1"),
		Id:                  testhelper.ToUUID("3e0ea031-f36c-4c43-ab70-34769dc5ba3a"),
		ProviderDisplayName: converter.String("[dev]\\Project Collection Build Service Accounts"),
		SubjectDescriptor:   converter.String("vssgp.Uy0xLTktMTU1MTM3NDI0NS00MjUxODEwMDMyLTIzOTk2NzI2NDYtMjg5OTA2MjQ3MS0xNTc4MjY2MDYyLTAtMC0wLTEtMQ"),
	},
	{
		Descriptor:          converter.String("Microsoft.TeamFoundation.Identity;S-1-9-1551374245-4251810032-2399672646-2899062471-1578266062-0-0-0-1-2"),
		Id:                  testhelper.ToUUID("e1c911d1-592e-451b-84d2-6c81dbb895c0"),
		ProviderDisplayName: converter.String("[dev]\\Project Collection Build Administrators"),
		SubjectDescriptor:   converter.String("vssgp.Uy0xLTktMTU1MTM3NDI0NS00MjUxODEwMDMyLTIzOTk2NzI2NDYtMjg5OTA2MjQ3MS0xNTc4MjY2MDYyLTAtMC0wLTEtMg"),
	},
	{
		Descriptor:          converter.String("Microsoft.TeamFoundation.Identity;S-1-9-1551374245-4251810032-2399672646-2899062471-1578266062-0-0-0-4-1"),
		Id:                  testhelper.ToUUID("b167c23a-27eb-4c59-aa7c-09794f38a556"),
		ProviderDisplayName: converter.String("[dev]\\Project Collection Test Service Accounts"),
		SubjectDescriptor:   converter.String("vssgp.Uy0xLTktMTU1MTM3NDI0NS00MjUxODEwMDMyLTIzOTk2NzI2NDYtMjg5OTA2MjQ3MS0xNTc4MjY2MDYyLTAtMC0wLTQtMQ"),
	},
}

func TestSecurityNamespace_GetActionDefinitions_HandleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)
	clients := &client.AggregatedClient{
		SecurityClient: securityClient,
		IdentityClient: azdosdkmocks.NewMockIdentityClient(ctrl),
		Ctx:            context.Background(),
	}

	sn, err := NewSecurityNamespace(clients.Ctx, SecurityNamespaceIDValues.Project, clients.SecurityClient, clients.IdentityClient)
	assert.Nil(t, err)
	assert.NotNil(t, sn)

	// QuerySecurityNamespaces
	errMsg := "@@QuerySecurityNamespaces@@failed@@"
	securityClient.
		EXPECT().
		QuerySecurityNamespaces(clients.Ctx, gomock.Any()).
		Return(nil, fmt.Errorf(errMsg)).
		Times(1)

	defs, err := sn.getActionDefinitions()
	assert.Nil(t, defs)
	assert.EqualError(t, err, errMsg)
}

func TestSecurityNamespace_GetActionDefinitions_EnsureExistingValuesUnchanged(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)
	clients := &client.AggregatedClient{
		SecurityClient: securityClient,
		IdentityClient: azdosdkmocks.NewMockIdentityClient(ctrl),
		Ctx:            context.Background(),
	}

	sn, err := NewSecurityNamespace(clients.Ctx, SecurityNamespaceIDValues.Project, clients.SecurityClient, clients.IdentityClient)
	assert.Nil(t, err)
	assert.NotNil(t, sn)

	// QuerySecurityNamespaces
	securityClient.
		EXPECT().
		QuerySecurityNamespaces(clients.Ctx, security.QuerySecurityNamespacesArgs{
			SecurityNamespaceId: &securityNamespaceDescriptionProjectId,
		}).
		Return(&securityNamespaceDescriptionProject, nil).
		Times(1)

	defs1, err := sn.getActionDefinitions()
	assert.Nil(t, err)
	assert.NotNil(t, defs1)

	// ensure second call does not call QuerySecurityNamespaces again
	defs2, err := sn.getActionDefinitions()
	assert.Nil(t, err)
	assert.NotNil(t, defs2)

	// ensure that both calls to getActionDefinitions retun the same values
	assert.Equal(t, defs1, defs2)
}

func TestSecurityNamespace_GetActionDefinitions_EmptyResultError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)
	clients := &client.AggregatedClient{
		SecurityClient: securityClient,
		IdentityClient: azdosdkmocks.NewMockIdentityClient(ctrl),
		Ctx:            context.Background(),
	}

	sn, err := NewSecurityNamespace(clients.Ctx, SecurityNamespaceIDValues.Project, clients.SecurityClient, clients.IdentityClient)
	assert.Nil(t, err)
	assert.NotNil(t, sn)

	// QuerySecurityNamespaces
	securityClient.
		EXPECT().
		QuerySecurityNamespaces(clients.Ctx, security.QuerySecurityNamespacesArgs{
			SecurityNamespaceId: &securityNamespaceDescriptionProjectId,
		}).
		Return(&securityNamespaceDescriptionProjectEmpty, nil).
		Times(1)

	defs, err := sn.getActionDefinitions()
	assert.NotNil(t, err)
	assert.Nil(t, defs)
}

func TestSecurityNamespace_GetActionDefinitions_ValidMapping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)
	clients := &client.AggregatedClient{
		SecurityClient: securityClient,
		IdentityClient: azdosdkmocks.NewMockIdentityClient(ctrl),
		Ctx:            context.Background(),
	}

	sn, err := NewSecurityNamespace(clients.Ctx, SecurityNamespaceIDValues.Project, clients.SecurityClient, clients.IdentityClient)
	assert.Nil(t, err)
	assert.NotNil(t, sn)

	// QuerySecurityNamespaces
	securityClient.
		EXPECT().
		QuerySecurityNamespaces(clients.Ctx, security.QuerySecurityNamespacesArgs{
			SecurityNamespaceId: &securityNamespaceDescriptionProjectId,
		}).
		Return(&securityNamespaceDescriptionProject, nil).
		Times(1)

	defs, err := sn.getActionDefinitions()
	assert.Nil(t, err)
	assert.NotNil(t, defs)
	assert.Equal(t, len(*securityNamespaceDescriptionProject[0].Actions), len(*defs))
	for _, action := range *securityNamespaceDescriptionProject[0].Actions {
		v, ok := (*defs)[*action.Name]
		assert.True(t, ok)
		assert.EqualValues(t, action, v)
	}
}

func TestSecurityNamespace_GetAccessControlList_HandleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)
	clients := &client.AggregatedClient{
		SecurityClient: securityClient,
		IdentityClient: azdosdkmocks.NewMockIdentityClient(ctrl),
		Ctx:            context.Background(),
	}

	sn, err := NewSecurityNamespace(clients.Ctx, SecurityNamespaceIDValues.Project, clients.SecurityClient, clients.IdentityClient)
	assert.Nil(t, err)
	assert.NotNil(t, sn)

	// QueryAccessControlLists
	errMsg := "@@QuerySecurityNamespaces@@failed@@"
	var descriptorList []string
	for _, identity := range projectIdentityList {
		descriptorList = append(descriptorList, *identity.Descriptor)
	}
	securityClient.
		EXPECT().
		QueryAccessControlLists(clients.Ctx, gomock.Any()).
		Return(nil, fmt.Errorf(errMsg)).
		Times(1)

	acl, err := sn.getAccessControlList(&projectAccessToken, &descriptorList)
	assert.NotNil(t, err)
	assert.Nil(t, acl)
	assert.EqualError(t, err, errMsg)
}

func TestSecurityNamespace_GetAccessControlList_EmptyResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)
	clients := &client.AggregatedClient{
		SecurityClient: securityClient,
		IdentityClient: azdosdkmocks.NewMockIdentityClient(ctrl),
		Ctx:            context.Background(),
	}

	sn, err := NewSecurityNamespace(clients.Ctx, SecurityNamespaceIDValues.Project, clients.SecurityClient, clients.IdentityClient)
	assert.Nil(t, err)
	assert.NotNil(t, sn)

	// QueryAccessControlLists
	var descriptors string
	var descriptorList []string
	for _, identity := range projectIdentityList {
		descriptorList = append(descriptorList, *identity.Descriptor)
	}
	descriptors = strings.Join(descriptorList, ",")
	securityClient.
		EXPECT().
		QueryAccessControlLists(clients.Ctx, security.QueryAccessControlListsArgs{
			SecurityNamespaceId: &securityNamespaceDescriptionProjectId,
			Token:               &projectAccessToken,
			Descriptors:         &descriptors,
			IncludeExtendedInfo: converter.Bool(true),
		}).
		Return(&projectAccessControlListEmpty, nil).
		Times(1)

	acl, err := sn.getAccessControlList(&projectAccessToken, &descriptorList)
	assert.Nil(t, err)
	assert.Nil(t, acl)
}

func TestSecurityNamespace_GetAccessControlList_NilResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)
	clients := &client.AggregatedClient{
		SecurityClient: securityClient,
		IdentityClient: azdosdkmocks.NewMockIdentityClient(ctrl),
		Ctx:            context.Background(),
	}

	sn, err := NewSecurityNamespace(clients.Ctx, SecurityNamespaceIDValues.Project, clients.SecurityClient, clients.IdentityClient)
	assert.Nil(t, err)
	assert.NotNil(t, sn)

	// QueryAccessControlLists
	var descriptors string
	var descriptorList []string
	for _, identity := range projectIdentityList {
		descriptorList = append(descriptorList, *identity.Descriptor)
	}
	descriptors = strings.Join(descriptorList, ",")
	securityClient.
		EXPECT().
		QueryAccessControlLists(clients.Ctx, security.QueryAccessControlListsArgs{
			SecurityNamespaceId: &securityNamespaceDescriptionProjectId,
			Token:               &projectAccessToken,
			Descriptors:         &descriptors,
			IncludeExtendedInfo: converter.Bool(true),
		}).
		Return(nil, nil).
		Times(1)

	acl, err := sn.getAccessControlList(&projectAccessToken, &descriptorList)
	assert.Nil(t, err)
	assert.Nil(t, acl)
}

func TestSecurityNamespace_GetAccessControlList_VerifyReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)
	clients := &client.AggregatedClient{
		SecurityClient: securityClient,
		IdentityClient: azdosdkmocks.NewMockIdentityClient(ctrl),
		Ctx:            context.Background(),
	}

	sn, err := NewSecurityNamespace(clients.Ctx, SecurityNamespaceIDValues.Project, clients.SecurityClient, clients.IdentityClient)
	assert.Nil(t, err)
	assert.NotNil(t, sn)

	// QueryAccessControlLists
	var descriptors string
	var descriptorList []string
	for _, identity := range projectIdentityList {
		descriptorList = append(descriptorList, *identity.Descriptor)
	}
	descriptors = strings.Join(descriptorList, ",")
	securityClient.
		EXPECT().
		QueryAccessControlLists(clients.Ctx, security.QueryAccessControlListsArgs{
			SecurityNamespaceId: &securityNamespaceDescriptionProjectId,
			Token:               &projectAccessToken,
			Descriptors:         &descriptors,
			IncludeExtendedInfo: converter.Bool(true),
		}).
		Return(&projectAccessControlList, nil).
		Times(1)

	acl, err := sn.getAccessControlList(&projectAccessToken, &descriptorList)
	assert.Nil(t, err)
	assert.NotNil(t, acl)
	assert.Equal(t, &projectAccessControlList[0], acl)
}

func TestSecurityNamespaces_GetIndentitiesFromSubjects_HandleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		IdentityClient: identityClient,
		Ctx:            context.Background(),
	}

	sn, err := NewSecurityNamespace(clients.Ctx, SecurityNamespaceIDValues.Project, clients.SecurityClient, clients.IdentityClient)
	assert.Nil(t, err)
	assert.NotNil(t, sn)

	// ReadIdentities
	errMsg := "@@ReadIdentities@@failed@@"
	var subjectDescriptorList []string
	for _, identity := range projectIdentityList {
		subjectDescriptorList = append(subjectDescriptorList, *identity.SubjectDescriptor)
	}
	identityClient.
		EXPECT().
		ReadIdentities(clients.Ctx, gomock.Any()).
		Return(nil, fmt.Errorf(errMsg)).
		Times(1)

	idList, err := sn.getIndentitiesFromSubjects(&subjectDescriptorList)
	assert.Nil(t, idList)
	assert.NotNil(t, err)
	assert.EqualError(t, err, errMsg)
}

func TestSecurityNamespaces_GetIndentitiesFromSubjects_HandleEmptyReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		IdentityClient: identityClient,
		Ctx:            context.Background(),
	}

	sn, err := NewSecurityNamespace(clients.Ctx, SecurityNamespaceIDValues.Project, clients.SecurityClient, clients.IdentityClient)
	assert.Nil(t, err)
	assert.NotNil(t, sn)

	// ReadIdentities
	var subjectDescriptors string
	var subjectDescriptorList []string
	for _, identity := range projectIdentityList {
		subjectDescriptorList = append(subjectDescriptorList, *identity.SubjectDescriptor)
	}
	subjectDescriptors = strings.Join(subjectDescriptorList, ",")
	identityClient.
		EXPECT().
		ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
			SubjectDescriptors: &subjectDescriptors,
		}).
		Return(&projectIdentityListEmpty, nil).
		Times(1)

	idList, err := sn.getIndentitiesFromSubjects(&subjectDescriptorList)
	assert.Nil(t, idList)
	assert.NotNil(t, err)
}

func TestSecurityNamespace_GetIndentitiesFromSubjects_VerifyReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{
		SecurityClient: azdosdkmocks.NewMockSecurityClient(ctrl),
		IdentityClient: identityClient,
		Ctx:            context.Background(),
	}

	sn, err := NewSecurityNamespace(clients.Ctx, SecurityNamespaceIDValues.Project, clients.SecurityClient, clients.IdentityClient)
	assert.Nil(t, err)
	assert.NotNil(t, sn)

	// ReadIdentities
	var subjectDescriptors string
	var subjectDescriptorList []string
	for _, identity := range projectIdentityList {
		subjectDescriptorList = append(subjectDescriptorList, *identity.SubjectDescriptor)
	}
	subjectDescriptors = strings.Join(subjectDescriptorList, ",")
	identityClient.
		EXPECT().
		ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
			SubjectDescriptors: &subjectDescriptors,
		}).
		Return(&projectIdentityList, nil).
		Times(1)

	idList, err := sn.getIndentitiesFromSubjects(&subjectDescriptorList)
	assert.NotNil(t, idList)
	assert.Nil(t, err)
	assert.Equal(t, projectIdentityList, *idList)
}

func TestSecurityNamespace_GetPrincipalPermissions_Verify(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)
	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{
		SecurityClient: securityClient,
		IdentityClient: identityClient,
		Ctx:            context.Background(),
	}

	sn, err := NewSecurityNamespace(clients.Ctx, SecurityNamespaceIDValues.Project, clients.SecurityClient, clients.IdentityClient)
	assert.Nil(t, err)
	assert.NotNil(t, sn)

	// getActionDefinitions => QuerySecurityNamespaces
	securityClient.
		EXPECT().
		QuerySecurityNamespaces(clients.Ctx, gomock.Any()).
		Return(&securityNamespaceDescriptionProject, nil).
		Times(1)

	// getIndentitiesFromSubjects => ReadIdentities
	var subjectDescriptorList []string
	subjectDescriptorMap := map[string]string{}
	for _, identity := range projectIdentityList {
		subjectDescriptorList = append(subjectDescriptorList, *identity.SubjectDescriptor)
		subjectDescriptorMap[*identity.SubjectDescriptor] = *identity.SubjectDescriptor
	}
	identityClient.
		EXPECT().
		ReadIdentities(clients.Ctx, gomock.Any()).
		Return(&projectIdentityList, nil).
		Times(1)

	// getAccessControlList => QueryAccessControlLists
	var descriptorList []string
	for _, identity := range projectIdentityList {
		descriptorList = append(descriptorList, *identity.Descriptor)
	}
	securityClient.
		EXPECT().
		QueryAccessControlLists(clients.Ctx, gomock.Any()).
		Return(&projectAccessControlList, nil).
		Times(1)

	token := "GO/UNITTEST/TOKEN"
	perms, err := sn.GetPrincipalPermissions(&token, &subjectDescriptorList)
	assert.NotNil(t, perms)
	assert.Nil(t, err)
	assert.Len(t, *perms, len(subjectDescriptorList))
	for _, v := range *perms {
		_, ok := subjectDescriptorMap[v.SubjectDescriptor]
		assert.True(t, ok)
	}
}
