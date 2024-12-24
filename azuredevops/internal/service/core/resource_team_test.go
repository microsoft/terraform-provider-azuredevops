//go:build (all || core || resource_team) && !exclude_resource_team
// +build all core resource_team
// +build !exclude_resource_team

package core

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/security"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

func TestTeam_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)

	clients := &client.AggregatedClient{
		CoreClient:     coreClient,
		IdentityClient: identityClient,
		SecurityClient: securityClient,
		Ctx:            context.Background(),
	}

	testProjectID := uuid.New()
	testTeamName := "@@TEST TEAM@@"

	coreClient.
		EXPECT().
		CreateTeam(clients.Ctx, core.CreateTeamArgs{
			ProjectId: converter.String(testProjectID.String()),
			Team: &core.WebApiTeam{
				Name: &testTeamName,
			},
		}).
		Return(nil, fmt.Errorf("@@CreateTeam@@failed@@")).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceTeam().Schema, nil)
	resourceData.Set("project_id", testProjectID.String())
	resourceData.Set("name", testTeamName)

	err := resourceTeamCreate(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "@@CreateTeam@@failed@@")
}

func TestTeam_Create_EnsureTeamDeletedOnAddAdministratorsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)

	clients := &client.AggregatedClient{
		CoreClient:     coreClient,
		IdentityClient: identityClient,
		SecurityClient: securityClient,
		Ctx:            context.Background(),
	}

	testProjectID := uuid.New()
	testTeamName := "@@TEST TEAM@@"
	testTeamID := uuid.New()

	coreClient.
		EXPECT().
		CreateTeam(clients.Ctx, core.CreateTeamArgs{
			ProjectId: converter.String(testProjectID.String()),
			Team: &core.WebApiTeam{
				Name: &testTeamName,
			},
		}).
		Return(&core.WebApiTeam{
			Id:        &testTeamID,
			Name:      &testTeamName,
			ProjectId: &testProjectID,
		}, nil).
		Times(1)

	nsID := uuid.UUID(securityhelper.SecurityNamespaceIDValues.Identity)
	securityClient.
		EXPECT().
		QuerySecurityNamespaces(clients.Ctx, security.QuerySecurityNamespacesArgs{
			SecurityNamespaceId: &nsID,
		}).
		Return(&[]security.SecurityNamespaceDescription{
			{
				Actions: &[]security.ActionDefinition{
					{
						Bit:  converter.Int(8),
						Name: converter.String("ManageMembership"),
					},
				},
			},
		}, nil).
		Times(2)

	adminSubjectDescriptor := "aad.ZmY0YWYyMjEtMWFhMi03YWNiLTllNGUtMGIwNzZiYTQ2Y2Yz"
	adminID := uuid.New()

	identityClient.
		EXPECT().
		ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
			SubjectDescriptors: &adminSubjectDescriptor,
		}).
		Return(&[]identity.Identity{
			{
				Descriptor:        converter.String(adminID.String()),
				SubjectDescriptor: &adminSubjectDescriptor,
				IsActive:          converter.Bool(true),
			},
		}, nil).
		Times(1)

	idToken := testProjectID.String() + "\\" + testTeamID.String()
	securityClient.
		EXPECT().
		QueryAccessControlLists(clients.Ctx, security.QueryAccessControlListsArgs{
			SecurityNamespaceId: &nsID,
			Token:               &idToken,
			Descriptors:         converter.String(adminID.String()),
			IncludeExtendedInfo: converter.Bool(true),
		}).
		Return(&[]security.AccessControlList{}, nil).
		Times(1)

	securityClient.
		EXPECT().
		QueryAccessControlLists(clients.Ctx, security.QueryAccessControlListsArgs{
			SecurityNamespaceId: &nsID,
			Token:               &idToken,
			IncludeExtendedInfo: converter.Bool(true),
		}).
		Return(&[]security.AccessControlList{}, nil).
		Times(1)

	coreClient.
		EXPECT().
		DeleteTeam(clients.Ctx, core.DeleteTeamArgs{
			ProjectId: converter.String(testProjectID.String()),
			TeamId:    converter.String(testTeamID.String()),
		}).
		Return(nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceTeam().Schema, nil)
	resourceData.Set("project_id", testProjectID.String())
	resourceData.Set("name", testTeamName)
	resourceData.Set("administrators", schema.NewSet(schema.HashString, []interface{}{
		adminSubjectDescriptor,
	}))

	err := resourceTeamCreate(resourceData, clients)
	require.NotNil(t, err)
}

func TestTeam_Create_EnsureTeamDeletedOnAddMembersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)

	clients := &client.AggregatedClient{
		CoreClient:     coreClient,
		IdentityClient: identityClient,
		SecurityClient: securityClient,
		Ctx:            context.Background(),
	}

	testProjectID := uuid.New()
	testTeamName := "@@TEST TEAM@@"
	testTeamID := uuid.New()

	coreClient.
		EXPECT().
		CreateTeam(clients.Ctx, core.CreateTeamArgs{
			ProjectId: converter.String(testProjectID.String()),
			Team: &core.WebApiTeam{
				Name: &testTeamName,
			},
		}).
		Return(&core.WebApiTeam{
			Id:        &testTeamID,
			Name:      &testTeamName,
			ProjectId: &testProjectID,
		}, nil).
		Times(1)

	memberSubjectDescriptor := "aad.ZmY0YWYyMjEtMWFhMi03YWNiLTllNGUtMGIwNzZiYTQ2Y2Yz"
	memberID := uuid.New()

	identityClient.
		EXPECT().
		ReadMembers(clients.Ctx, identity.ReadMembersArgs{
			ContainerId: converter.String(testTeamID.String()),
		}).
		Return(nil, nil).
		Times(1)

	identityClient.
		EXPECT().
		ReadIdentities(clients.Ctx, identity.ReadIdentitiesArgs{
			SubjectDescriptors: &memberSubjectDescriptor,
		}).
		Return(&[]identity.Identity{
			{
				Id:                &memberID,
				Descriptor:        converter.String(memberID.String()),
				SubjectDescriptor: &memberSubjectDescriptor,
				IsActive:          converter.Bool(true),
			},
		}, nil).
		Times(1)

	identityClient.
		EXPECT().
		AddMember(clients.Ctx, identity.AddMemberArgs{
			ContainerId: converter.String(testTeamID.String()),
			MemberId:    converter.String(memberID.String()),
		}).
		Return(converter.Bool(false), nil).
		Times(1)

	coreClient.
		EXPECT().
		DeleteTeam(clients.Ctx, core.DeleteTeamArgs{
			ProjectId: converter.String(testProjectID.String()),
			TeamId:    converter.String(testTeamID.String()),
		}).
		Return(nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceTeam().Schema, nil)
	resourceData.Set("project_id", testProjectID.String())
	resourceData.Set("name", testTeamName)
	resourceData.Set("members", schema.NewSet(schema.HashString, []interface{}{
		memberSubjectDescriptor,
	}))

	err := resourceTeamCreate(resourceData, clients)
	require.NotNil(t, err)
}

func TestTeam_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)

	clients := &client.AggregatedClient{
		CoreClient:     coreClient,
		IdentityClient: identityClient,
		SecurityClient: securityClient,
		Ctx:            context.Background(),
	}

	testProjectID := uuid.New()
	testTeamName := "@@TEST TEAM@@"
	testTeamID := uuid.New()
	errMsg := "@@GetTeam@@failed@@"

	coreClient.
		EXPECT().
		GetTeam(clients.Ctx, core.GetTeamArgs{
			ProjectId:      converter.String(testProjectID.String()),
			TeamId:         converter.String(testTeamID.String()),
			ExpandIdentity: converter.Bool(false),
		}).
		Return(nil, fmt.Errorf(errMsg)).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceTeam().Schema, nil)
	resourceData.SetId(testTeamID.String())
	resourceData.Set("project_id", testProjectID.String())
	resourceData.Set("name", testTeamName)
	err := resourceTeamRead(resourceData, clients)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), errMsg)
}

func TestTeam_Read_HandlesNotFoundCorrectly(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)

	clients := &client.AggregatedClient{
		CoreClient:     coreClient,
		IdentityClient: identityClient,
		SecurityClient: securityClient,
		Ctx:            context.Background(),
	}

	testProjectID := uuid.New()
	testTeamName := "@@TEST TEAM@@"
	testTeamID := uuid.New()

	coreClient.
		EXPECT().
		GetTeam(clients.Ctx, core.GetTeamArgs{
			ProjectId:      converter.String(testProjectID.String()),
			TeamId:         converter.String(testTeamID.String()),
			ExpandIdentity: converter.Bool(false),
		}).
		Return(nil, azuredevops.WrappedError{
			StatusCode: converter.Int(http.StatusNotFound),
		}).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceTeam().Schema, nil)
	resourceData.SetId(testTeamID.String())
	resourceData.Set("project_id", testProjectID.String())
	resourceData.Set("name", testTeamName)

	err := resourceTeamRead(resourceData, clients)
	require.Nil(t, err)
	require.Zero(t, resourceData.Id())
}

func TestTeam_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)

	clients := &client.AggregatedClient{
		CoreClient:     coreClient,
		IdentityClient: identityClient,
		SecurityClient: securityClient,
		Ctx:            context.Background(),
	}

	testProjectID := uuid.New()
	testTeamName := "@@TEST TEAM@@"
	testTeamID := uuid.New()

	coreClient.
		EXPECT().
		GetTeam(clients.Ctx, core.GetTeamArgs{
			ProjectId:      converter.String(testProjectID.String()),
			TeamId:         converter.String(testTeamID.String()),
			ExpandIdentity: converter.Bool(false),
		}).
		Return(nil, fmt.Errorf("@@GetTeam@@failed@@")).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceTeam().Schema, nil)
	resourceData.SetId(testTeamID.String())
	resourceData.Set("project_id", testProjectID.String())
	resourceData.Set("name", testTeamName)

	err := resourceTeamUpdate(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "@@GetTeam@@failed@@")
	require.NotZero(t, resourceData.Id())
}
