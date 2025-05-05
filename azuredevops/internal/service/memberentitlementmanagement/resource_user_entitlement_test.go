//go:build (all || resource_user_entitlement) && !exclude_resource_user_entitlement
// +build all resource_user_entitlement
// +build !exclude_resource_user_entitlement

package memberentitlementmanagement

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/licensing"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/memberentitlementmanagement"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// if origin_id is "" and principal_name is supplied, the principal_name will be used.
func TestUserEntitlement_CreateUserEntitlement_WithPrincipalName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	accountLicenseType := licensing.AccountLicenseTypeValues.Express
	origin := ""
	originID := ""
	principalName := "foobar@microsoft.com"
	descriptor := "baz"
	id := uuid.New()
	mockUserEntitlement := getMockUserEntitlement(&id, accountLicenseType, origin, originID, principalName, descriptor)

	resourceData := schema.TestResourceDataRaw(t, ResourceUserEntitlement().Schema, nil)
	resourceData.Set("principal_name", principalName)

	expectedIsSuccess := true
	memberEntitlementClient.
		EXPECT().
		AddUserEntitlement(gomock.Any(), MatchAddUserEntitlementArgs(t, memberentitlementmanagement.AddUserEntitlementArgs{
			UserEntitlement: mockUserEntitlement,
		})).
		Return(&memberentitlementmanagement.UserEntitlementsPostResponse{
			IsSuccess:       &expectedIsSuccess,
			UserEntitlement: mockUserEntitlement,
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetUserEntitlement(gomock.Any(), memberentitlementmanagement.GetUserEntitlementArgs{
			UserId: mockUserEntitlement.Id,
		}).
		Return(mockUserEntitlement, nil)

	err := resourceUserEntitlementCreate(resourceData, clients)
	assert.Nil(t, err, "err should not be nil")
}

//
//// if origin_id is "" and principal_name is "", an error will be reported.
//func TestUserEntitlement_CreateUserEntitlement_Need_OriginID_Or_PrincipalName(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	clients := &client.AggregatedClient{
//		MemberEntitleManagementClient: nil,
//		Ctx:                           context.Background(),
//	}
//
//	resourceData := schema.TestResourceDataRaw(t, ResourceUserEntitlement().Schema, nil)
//	// originID and principalName is not set.
//
//	err := resourceUserEntitlementCreate(resourceData, clients)
//	assert.NotNil(t, err, "err should not be nil")
//	require.Regexp(t, "Use origin_id or principal_name", err.Error())
//}

// if the REST-API return the failure, it should fail.

func TestUserEntitlement_CreateUserEntitlement_WithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	principalName := "foobar@microsoft.com"

	resourceData := schema.TestResourceDataRaw(t, ResourceUserEntitlement().Schema, nil)
	// resourceData.Set("origin_id", originID)
	resourceData.Set("account_license_type", "express")
	resourceData.Set("principal_name", principalName)

	// No error but it has a error on the response.
	memberEntitlementClient.
		EXPECT().
		AddUserEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("error foo")).
		Times(1)

	err := resourceUserEntitlementCreate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
}

// if the REST-API return the success, but fails on response
func TestUserEntitlement_CreateUserEntitlement_WithEarlyAdopter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	principalName := "foobar@microsoft.com"

	resourceData := schema.TestResourceDataRaw(t, ResourceUserEntitlement().Schema, nil)
	// resourceData.Set("origin_id", originID)
	resourceData.Set("account_license_type", "earlyAdopter")
	resourceData.Set("principal_name", principalName)

	var expectedKey interface{} = 5000
	var expectedValue interface{} = "A user cannot be assigned an Account-EarlyAdopter license."
	expectedErrors := []azuredevops.KeyValuePair{
		{
			Key:   &expectedKey,
			Value: &expectedValue,
		},
	}
	expectedIsSuccess := false

	// No error but it has a error on the response.
	memberEntitlementClient.
		EXPECT().
		AddUserEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.UserEntitlementsPostResponse{
			IsSuccess: &expectedIsSuccess,
			OperationResult: &memberentitlementmanagement.UserEntitlementOperationResult{
				IsSuccess: &expectedIsSuccess,
				Errors:    &expectedErrors,
			},
		}, nil).
		Times(1)

	err := resourceUserEntitlementCreate(resourceData, clients)
	require.Contains(t, err.Error(), "A user cannot be assigned an Account-EarlyAdopter license.")
}

// TestUserEntitlement_Update_TestChangeEntitlement verfies that an entitlement can be changed
func TestUserEntitlement_Update_TestChangeEntitlement(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	accountLicenseType := licensing.AccountLicenseTypeValues.Stakeholder
	origin := ""
	originID := ""
	principalName := "foobar@microsoft.com"
	descriptor := "baz"
	id := uuid.New()
	mockUserEntitlement := getMockUserEntitlement(&id, accountLicenseType, origin, originID, principalName, descriptor)
	expectedIsSuccess := true

	memberEntitlementClient.
		EXPECT().
		UpdateUserEntitlement(gomock.Any(), memberentitlementmanagement.UpdateUserEntitlementArgs{
			UserId: &id,
			Document: &[]webapi.JsonPatchOperation{
				{
					Op:   &webapi.OperationValues.Replace,
					From: nil,
					Path: converter.String("/accessLevel"),
					Value: struct {
						AccountLicenseType string `json:"accountLicenseType"`
						LicensingSource    string `json:"licensingSource"`
					}{
						string(licensing.AccountLicenseTypeValues.Stakeholder),
						string(licensing.LicensingSourceValues.Account),
					},
				},
			},
		}).
		Return(&memberentitlementmanagement.UserEntitlementsPatchResponse{
			IsSuccess:       &expectedIsSuccess,
			UserEntitlement: mockUserEntitlement,
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetUserEntitlement(gomock.Any(), memberentitlementmanagement.GetUserEntitlementArgs{
			UserId: mockUserEntitlement.Id,
		}).
		Return(mockUserEntitlement, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceUserEntitlement().Schema, nil)
	resourceData.SetId(id.String())
	resourceData.Set("principal_name", principalName)
	resourceData.Set("account_license_type", string(licensing.AccountLicenseTypeValues.Stakeholder))
	resourceData.Set("licensing_source", string(licensing.LicensingSourceValues.Account))

	err := resourceUserEntitlementUpdate(resourceData, clients)
	assert.Nil(t, err)
}

// TestUserEntitlement_CreateUpdate_TestBasicEntitlement verifies that the (virtual) Basic entitlement can be set
func TestUserEntitlement_CreateUpdate_TestBasicEntitlement(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	accountLicenseType := licensing.AccountLicenseTypeValues.Express
	origin := ""
	originID := ""
	principalName := "foobar@microsoft.com"
	descriptor := "baz"
	id := uuid.New()
	mockUserEntitlement := getMockUserEntitlement(&id, accountLicenseType, origin, originID, principalName, descriptor)
	expectedIsSuccess := true

	memberEntitlementClient.
		EXPECT().
		AddUserEntitlement(gomock.Any(), MatchAddUserEntitlementArgs(t, memberentitlementmanagement.AddUserEntitlementArgs{
			UserEntitlement: mockUserEntitlement,
		})).
		Return(&memberentitlementmanagement.UserEntitlementsPostResponse{
			IsSuccess:       &expectedIsSuccess,
			UserEntitlement: mockUserEntitlement,
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetUserEntitlement(gomock.Any(), memberentitlementmanagement.GetUserEntitlementArgs{
			UserId: mockUserEntitlement.Id,
		}).
		Return(mockUserEntitlement, nil)

	resourceData := schema.TestResourceDataRaw(t, ResourceUserEntitlement().Schema, nil)
	resourceData.Set("principal_name", principalName)
	resourceData.Set("account_license_type", "basic")

	err := resourceUserEntitlementCreate(resourceData, clients)
	assert.Nil(t, err, "err should be nil")
}

// TestUserEntitlement_Import_TestUPN tests if import is successful using an UPN
func TestUserEntitlement_Import_TestUPN(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{
		IdentityClient: identityClient,
		Ctx:            context.Background(),
	}

	principalName := "foobar@microsoft.com"
	id := uuid.New()

	identityClient.
		EXPECT().
		ReadIdentities(gomock.Any(), gomock.Any()).
		Return(&[]identity.Identity{
			{
				Id: &id,
			},
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceUserEntitlement().Schema, nil)
	resourceData.SetId(principalName)

	d, err := importUserEntitlement(resourceData, clients)
	assert.Nil(t, err)
	assert.NotNil(t, d)
	assert.Len(t, d, 1)
	assert.Equal(t, id.String(), d[0].Id())
}

// TestUserEntitlement_Import_TestID tests if import is successful using an UUID
func TestUserEntitlement_Import_TestID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	id := uuid.New().String()
	resourceData := schema.TestResourceDataRaw(t, ResourceUserEntitlement().Schema, nil)
	resourceData.SetId(id)

	d, err := importUserEntitlement(resourceData, clients)
	assert.Nil(t, err)
	assert.NotNil(t, d)
	assert.Len(t, d, 1)
	assert.Equal(t, id, d[0].Id())
}

// TestUserEntilement_Import_TestInvalidValue tests if only a valid UPN and UUID can be used to import a resource
func TestUserEntilement_Import_TestInvalidValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	id := "InvalidValue-a73c5191-e20d"
	resourceData := schema.TestResourceDataRaw(t, ResourceUserEntitlement().Schema, nil)
	resourceData.SetId(id)

	d, err := importUserEntitlement(resourceData, clients)
	assert.Nil(t, d)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Only UUID and UPN values can used for import")
}

func TestUserEntitlement_Create_TestErrorFormatting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	id, _ := uuid.NewUUID()
	expectedIsSuccess := false
	k1 := interface{}("9999")
	v1 := interface{}("Error1")
	k2 := interface{}("9998")
	v2 := interface{}("Error2")

	memberEntitlementClient.
		EXPECT().
		AddUserEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.UserEntitlementsPostResponse{
			IsSuccess:       &expectedIsSuccess,
			UserEntitlement: nil,
			OperationResult: &memberentitlementmanagement.UserEntitlementOperationResult{
				IsSuccess: &expectedIsSuccess,
				Result:    nil,
				UserId:    &id,
				Errors: &[]azuredevops.KeyValuePair{
					{
						Key:   &k1,
						Value: &v1,
					},
					{
						Key:   &k2,
						Value: &v2,
					},
				},
			},
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetUserEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, ResourceUserEntitlement().Schema, nil)
	resourceData.Set("principal_name", "foobar@microsoft.com")

	err := resourceUserEntitlementCreate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	assert.Contains(t, err.Error(), "(9999) Error1")
	assert.Contains(t, err.Error(), "(9998) Error2")
}

func TestUserEntitlement_Create_TestEmptyErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	id, _ := uuid.NewUUID()
	expectedIsSuccess := false

	memberEntitlementClient.
		EXPECT().
		AddUserEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.UserEntitlementsPostResponse{
			IsSuccess:       &expectedIsSuccess,
			UserEntitlement: nil,
			OperationResult: &memberentitlementmanagement.UserEntitlementOperationResult{
				IsSuccess: &expectedIsSuccess,
				Result:    nil,
				UserId:    &id,
				Errors:    nil,
			},
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetUserEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, ResourceUserEntitlement().Schema, nil)
	resourceData.Set("principal_name", "foobar@microsoft.com")

	err := resourceUserEntitlementCreate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	assert.Contains(t, err.Error(), "Unknown API error")
}

func TestUserEntitlement_Update_TestErrorFormatting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	id, _ := uuid.NewUUID()
	expectedIsSuccess := false
	k1 := interface{}("9999")
	v1 := interface{}("Error1")
	k2 := interface{}("9998")
	v2 := interface{}("Error2")

	memberEntitlementClient.
		EXPECT().
		UpdateUserEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.UserEntitlementsPatchResponse{
			IsSuccess:       &expectedIsSuccess,
			UserEntitlement: nil,
			OperationResults: &[]memberentitlementmanagement.UserEntitlementOperationResult{
				{
					IsSuccess: &expectedIsSuccess,
					Result:    nil,
					UserId:    &id,
					Errors: &[]azuredevops.KeyValuePair{
						{
							Key:   &k1,
							Value: &v1,
						},
						{
							Key:   &k2,
							Value: &v2,
						},
					},
				},
			},
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetUserEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, ResourceUserEntitlement().Schema, nil)
	resourceData.SetId(id.String())
	resourceData.Set("principal_name", "foobar@microsoft.com")

	err := resourceUserEntitlementUpdate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	assert.Contains(t, err.Error(), "(9999) Error1")
	assert.Contains(t, err.Error(), "(9998) Error2")
}

func TestUserEntitlement_Update_TestEmptyErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	id, _ := uuid.NewUUID()
	expectedIsSuccess := false

	memberEntitlementClient.
		EXPECT().
		UpdateUserEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.UserEntitlementsPatchResponse{
			IsSuccess:       &expectedIsSuccess,
			UserEntitlement: nil,
			OperationResults: &[]memberentitlementmanagement.UserEntitlementOperationResult{
				{
					IsSuccess: &expectedIsSuccess,
					Result:    nil,
					UserId:    &id,
					Errors:    nil,
				},
			},
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetUserEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, ResourceUserEntitlement().Schema, nil)
	resourceData.SetId(id.String())
	resourceData.Set("principal_name", "foobar@microsoft.com")

	err := resourceUserEntitlementUpdate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	assert.Contains(t, err.Error(), "Unknown API error")
}

func getMockUserEntitlement(id *uuid.UUID, accountLicenseType licensing.AccountLicenseType, origin string, originID string, principalName string, descriptor string) *memberentitlementmanagement.UserEntitlement {
	subjectKind := "user"
	licensingSource := licensing.LicensingSourceValues.Account

	return &memberentitlementmanagement.UserEntitlement{
		AccessLevel: &licensing.AccessLevel{
			AccountLicenseType: &accountLicenseType,
			LicensingSource:    &licensingSource,
		},
		Id: id,
		User: &graph.GraphUser{
			Origin:        &origin,
			OriginId:      &originID,
			PrincipalName: &principalName,
			SubjectKind:   &subjectKind,
			Descriptor:    &descriptor,
		},
	}
}

type matchAddUserEntitlementArgs struct {
	t *testing.T
	x memberentitlementmanagement.AddUserEntitlementArgs
}

func MatchAddUserEntitlementArgs(t *testing.T, x memberentitlementmanagement.AddUserEntitlementArgs) gomock.Matcher {
	return &matchAddUserEntitlementArgs{t, x}
}

func (m *matchAddUserEntitlementArgs) Matches(x interface{}) bool {
	args := x.(memberentitlementmanagement.AddUserEntitlementArgs)
	m.t.Logf("MatchAddUserEntitlementArgs:\nVALUE: account_license_type: [%s], licensing_source: [%s], origin: [%s], origin_id: [%s], principal_name: [%s]\n  REF: account_license_type: [%s], licensing_source: [%s], origin: [%s], origin_id: [%s], principal_name: [%s]\n",
		*args.UserEntitlement.AccessLevel.AccountLicenseType,
		*args.UserEntitlement.AccessLevel.LicensingSource,
		*args.UserEntitlement.User.Origin,
		*args.UserEntitlement.User.OriginId,
		*args.UserEntitlement.User.PrincipalName,
		*m.x.UserEntitlement.AccessLevel.AccountLicenseType,
		*m.x.UserEntitlement.AccessLevel.LicensingSource,
		*m.x.UserEntitlement.User.Origin,
		*m.x.UserEntitlement.User.OriginId,
		*m.x.UserEntitlement.User.PrincipalName)

	return *args.UserEntitlement.AccessLevel.AccountLicenseType == *m.x.UserEntitlement.AccessLevel.AccountLicenseType &&
		*args.UserEntitlement.User.Origin == *m.x.UserEntitlement.User.Origin &&
		*args.UserEntitlement.User.OriginId == *m.x.UserEntitlement.User.OriginId &&
		*args.UserEntitlement.User.PrincipalName == *m.x.UserEntitlement.User.PrincipalName
}

func (m *matchAddUserEntitlementArgs) String() string {
	return fmt.Sprintf("account_license_type: [%s], licensing_source: [%s], origin: [%s], origin_id: [%s], principal_name: [%s]",
		*m.x.UserEntitlement.AccessLevel.AccountLicenseType,
		*m.x.UserEntitlement.AccessLevel.LicensingSource,
		*m.x.UserEntitlement.User.Origin,
		*m.x.UserEntitlement.User.OriginId,
		*m.x.UserEntitlement.User.PrincipalName)
}
