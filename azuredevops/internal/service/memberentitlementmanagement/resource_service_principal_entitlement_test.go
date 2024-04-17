//go:build (all || resource_user_entitlement) && !exclude_resource_user_entitlement
// +build all resource_user_entitlement
// +build !exclude_resource_user_entitlement

package memberentitlementmanagement

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
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
)

// if origin_id is provided, it will be used. if principal_name is also supplied, an error will be reported.
func TestServicePrincipalEntitlement_CreateServicePrincipalEntitlement_DoNotAllowToSetOridinIdAndPrincipalName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: nil,
		Ctx:                           context.Background(),
	}

	originID := "e97b0e7f-0a61-41ad-860c-748ec5fcb20b"
	principalName := "foobar@microsoft.com"

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.Set("origin_id", originID)
	resourceData.Set("principal_name", principalName)

	err := resourceServicePrincipalEntitlementCreate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	require.Regexp(t, "Both origin_id and principal_name set. You can not use both", err.Error())
}

// if origin_id is "" and principal_name is supplied, the principal_name will be used.
func TestServicePrincipalEntitlement_CreateServicePrincipalEntitlement_WithPrincipalName(t *testing.T) {
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
	mockServicePrincipalEntitlement := getMockServicePrincipalEntitlement(&id, accountLicenseType, origin, originID, principalName, descriptor)

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.Set("principal_name", principalName)

	expectedIsSuccess := true
	memberEntitlementClient.
		EXPECT().
		AddServicePrincipalEntitlement(gomock.Any(), MatchAddServicePrincipalEntitlementArgs(t, memberentitlementmanagement.AddServicePrincipalEntitlementArgs{
			ServicePrincipalEntitlement: mockServicePrincipalEntitlement,
		})).
		Return(&memberentitlementmanagement.ServicePrincipalEntitlementsPostResponse{
			IsSuccess:                   &expectedIsSuccess,
			ServicePrincipalEntitlement: mockServicePrincipalEntitlement,
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetServicePrincipalEntitlement(gomock.Any(), memberentitlementmanagement.GetServicePrincipalEntitlementArgs{
			ServicePrincipalId: mockServicePrincipalEntitlement.Id,
		}).
		Return(mockServicePrincipalEntitlement, nil)

	err := resourceServicePrincipalEntitlementCreate(resourceData, clients)
	assert.Nil(t, err, "err should not be nil")
}

// if origin_id is "" and principal_name is "", an error will be reported.
func TestServicePrincipalEntitlement_CreateServicePrincipalEntitlement_Need_OriginID_Or_PrincipalName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: nil,
		Ctx:                           context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	// originID and principalName is not set.

	err := resourceServicePrincipalEntitlementCreate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	require.Regexp(t, "Use origin_id or principal_name", err.Error())
}

// if the REST-API return the failure, it should fail.

func TestServicePrincipalEntitlement_CreateServicePrincipalEntitlement_WithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	principalName := "foobar@microsoft.com"

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	// resourceData.Set("origin_id", originID)
	resourceData.Set("account_license_type", "express")
	resourceData.Set("principal_name", principalName)

	// No error but it has a error on the response.
	memberEntitlementClient.
		EXPECT().
		AddServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("error foo")).
		Times(1)

	err := resourceServicePrincipalEntitlementCreate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
}

// if the REST-API return the success, but fails on response
func TestServicePrincipalEntitlement_CreateServicePrincipalEntitlement_WithEarlyAdopter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	principalName := "foobar@microsoft.com"

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
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
		AddServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.ServicePrincipalEntitlementsPostResponse{
			IsSuccess: &expectedIsSuccess,
			OperationResult: &memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{
				IsSuccess: &expectedIsSuccess,
				Errors:    &expectedErrors,
			},
		}, nil).
		Times(1)

	err := resourceServicePrincipalEntitlementCreate(resourceData, clients)
	require.Contains(t, err.Error(), "A user cannot be assigned an Account-EarlyAdopter license.")
}

// TestServicePrincipalEntitlement_Update_TestChangeEntitlement verfies that an entitlement can be changed
func TestServicePrincipalEntitlement_Update_TestChangeEntitlement(t *testing.T) {
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
	mockServicePrincipalEntitlement := getMockServicePrincipalEntitlement(&id, accountLicenseType, origin, originID, principalName, descriptor)
	expectedIsSuccess := true

	memberEntitlementClient.
		EXPECT().
		UpdateServicePrincipalEntitlement(gomock.Any(), memberentitlementmanagement.UpdateServicePrincipalEntitlementArgs{
			ServicePrincipalId: &id,
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
		Return(&memberentitlementmanagement.ServicePrincipalEntitlementsPatchResponse{
			IsSuccess:                   &expectedIsSuccess,
			ServicePrincipalEntitlement: mockServicePrincipalEntitlement,
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetServicePrincipalEntitlement(gomock.Any(), memberentitlementmanagement.GetServicePrincipalEntitlementArgs{
			ServicePrincipalId: mockServicePrincipalEntitlement.Id,
		}).
		Return(mockServicePrincipalEntitlement, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.SetId(id.String())
	resourceData.Set("principal_name", principalName)
	resourceData.Set("account_license_type", string(licensing.AccountLicenseTypeValues.Stakeholder))
	resourceData.Set("licensing_source", string(licensing.LicensingSourceValues.Account))

	err := resourceServicePrincipalEntitlementUpdate(resourceData, clients)
	assert.Nil(t, err)
}

// TestServicePrincipalEntitlement_CreateUpdate_TestBasicEntitlement verifies that the (virtual) Basic entitlement can be set
func TestServicePrincipalEntitlement_CreateUpdate_TestBasicEntitlement(t *testing.T) {
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
	mockServicePrincipalEntitlement := getMockServicePrincipalEntitlement(&id, accountLicenseType, origin, originID, principalName, descriptor)
	expectedIsSuccess := true

	memberEntitlementClient.
		EXPECT().
		AddServicePrincipalEntitlement(gomock.Any(), MatchAddServicePrincipalEntitlementArgs(t, memberentitlementmanagement.AddServicePrincipalEntitlementArgs{
			ServicePrincipalEntitlement: mockServicePrincipalEntitlement,
		})).
		Return(&memberentitlementmanagement.ServicePrincipalEntitlementsPostResponse{
			IsSuccess:                   &expectedIsSuccess,
			ServicePrincipalEntitlement: mockServicePrincipalEntitlement,
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetServicePrincipalEntitlement(gomock.Any(), memberentitlementmanagement.GetServicePrincipalEntitlementArgs{
			ServicePrincipalId: mockServicePrincipalEntitlement.Id,
		}).
		Return(mockServicePrincipalEntitlement, nil)

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.Set("principal_name", principalName)
	resourceData.Set("account_license_type", "basic")

	err := resourceServicePrincipalEntitlementCreate(resourceData, clients)
	assert.Nil(t, err, "err should be nil")
}

// TestServicePrincipalEntitlement_Import_TestUPN tests if import is successful using an UPN
func TestServicePrincipalEntitlement_Import_TestUPN(t *testing.T) {
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

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.SetId(principalName)

	d, err := importServicePrincipalEntitlement(resourceData, clients)
	assert.Nil(t, err)
	assert.NotNil(t, d)
	assert.Len(t, d, 1)
	assert.Equal(t, id.String(), d[0].Id())
}

// TestServicePrincipalEntitlement_Import_TestID tests if import is successful using an UUID
func TestServicePrincipalEntitlement_Import_TestID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	id := uuid.New().String()
	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.SetId(id)

	d, err := importServicePrincipalEntitlement(resourceData, clients)
	assert.Nil(t, err)
	assert.NotNil(t, d)
	assert.Len(t, d, 1)
	assert.Equal(t, id, d[0].Id())
}

// TestServicePrincipalEntitlement_Import_TestInvalidValue tests if only a valid UPN and UUID can be used to import a resource
func TestServicePrincipalEntitlement_Import_TestInvalidValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	id := "InvalidValue-a73c5191-e20d"
	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.SetId(id)

	d, err := importServicePrincipalEntitlement(resourceData, clients)
	assert.Nil(t, d)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Only UUID and UPN values can used for import")
}

func TestServicePrincipalEntitlement_Create_TestErrorFormatting(t *testing.T) {
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
		AddServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.ServicePrincipalEntitlementsPostResponse{
			IsSuccess:                   &expectedIsSuccess,
			ServicePrincipalEntitlement: nil,
			OperationResult: &memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{
				IsSuccess:          &expectedIsSuccess,
				Result:             nil,
				ServicePrincipalId: &id,
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
		GetServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.Set("principal_name", "foobar@microsoft.com")

	err := resourceServicePrincipalEntitlementCreate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	assert.Contains(t, err.Error(), "(9999) Error1")
	assert.Contains(t, err.Error(), "(9998) Error2")
}

func TestServicePrincipalEntitlement_Create_TestEmptyErrors(t *testing.T) {
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
		AddServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.ServicePrincipalEntitlementsPostResponse{
			IsSuccess:                   &expectedIsSuccess,
			ServicePrincipalEntitlement: nil,
			OperationResult: &memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{
				IsSuccess:          &expectedIsSuccess,
				Result:             nil,
				ServicePrincipalId: &id,
				Errors:             nil,
			},
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.Set("principal_name", "foobar@microsoft.com")

	err := resourceServicePrincipalEntitlementCreate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	assert.Contains(t, err.Error(), "Unknown API error")
}

func TestServicePrincipalEntitlement_Update_TestErrorFormatting(t *testing.T) {
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
		UpdateServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.ServicePrincipalEntitlementsPatchResponse{
			IsSuccess:                   &expectedIsSuccess,
			ServicePrincipalEntitlement: nil,
			OperationResults: &[]memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{
				{
					IsSuccess:          &expectedIsSuccess,
					Result:             nil,
					ServicePrincipalId: &id,
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
		GetServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.SetId(id.String())
	resourceData.Set("principal_name", "foobar@microsoft.com")

	err := resourceServicePrincipalEntitlementUpdate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	assert.Contains(t, err.Error(), "(9999) Error1")
	assert.Contains(t, err.Error(), "(9998) Error2")
}

func TestServicePrincipalEntitlement_Update_TestEmptyErrors(t *testing.T) {
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
		UpdateServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.ServicePrincipalEntitlementsPatchResponse{
			IsSuccess:                   &expectedIsSuccess,
			ServicePrincipalEntitlement: nil,
			OperationResults: &[]memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{
				{
					IsSuccess:          &expectedIsSuccess,
					Result:             nil,
					ServicePrincipalId: &id,
					Errors:             nil,
				},
			},
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.SetId(id.String())
	resourceData.Set("principal_name", "foobar@microsoft.com")

	err := resourceServicePrincipalEntitlementUpdate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	assert.Contains(t, err.Error(), "Unknown API error")
}

func getMockServicePrincipalEntitlement(id *uuid.UUID, accountLicenseType licensing.AccountLicenseType, origin string, originID string, principalName string, descriptor string) *memberentitlementmanagement.ServicePrincipalEntitlement {
	subjectKind := "servicePrincipal"
	licensingSource := licensing.LicensingSourceValues.Account

	return &memberentitlementmanagement.ServicePrincipalEntitlement{
		AccessLevel: &licensing.AccessLevel{
			AccountLicenseType: &accountLicenseType,
			LicensingSource:    &licensingSource,
		},
		Id: id,
		ServicePrincipal: &graph.GraphServicePrincipal{
			Origin:        &origin,
			OriginId:      &originID,
			PrincipalName: &principalName,
			SubjectKind:   &subjectKind,
			Descriptor:    &descriptor,
		},
	}
}

type matchAddServicePrincipalEntitlementArgs struct {
	t *testing.T
	x memberentitlementmanagement.AddServicePrincipalEntitlementArgs
}

func MatchAddServicePrincipalEntitlementArgs(t *testing.T, x memberentitlementmanagement.AddServicePrincipalEntitlementArgs) gomock.Matcher {
	return &matchAddServicePrincipalEntitlementArgs{t, x}
}

func (m *matchAddServicePrincipalEntitlementArgs) Matches(x interface{}) bool {
	args := x.(memberentitlementmanagement.AddServicePrincipalEntitlementArgs)
	m.t.Logf("MatchAddServicePrincipalEntitlementArgs:\nVALUE: account_license_type: [%s], licensing_source: [%s], origin: [%s], origin_id: [%s], principal_name: [%s]\n  REF: account_license_type: [%s], licensing_source: [%s], origin: [%s], origin_id: [%s], principal_name: [%s]\n",
		*args.ServicePrincipalEntitlement.AccessLevel.AccountLicenseType,
		*args.ServicePrincipalEntitlement.AccessLevel.LicensingSource,
		*args.ServicePrincipalEntitlement.ServicePrincipal.Origin,
		*args.ServicePrincipalEntitlement.ServicePrincipal.OriginId,
		*args.ServicePrincipalEntitlement.ServicePrincipal.PrincipalName,
		*m.x.ServicePrincipalEntitlement.AccessLevel.AccountLicenseType,
		*m.x.ServicePrincipalEntitlement.AccessLevel.LicensingSource,
		*m.x.ServicePrincipalEntitlement.ServicePrincipal.Origin,
		*m.x.ServicePrincipalEntitlement.ServicePrincipal.OriginId,
		*m.x.ServicePrincipalEntitlement.ServicePrincipal.PrincipalName)

	return *args.ServicePrincipalEntitlement.AccessLevel.AccountLicenseType == *m.x.ServicePrincipalEntitlement.AccessLevel.AccountLicenseType &&
		*args.ServicePrincipalEntitlement.ServicePrincipal.Origin == *m.x.ServicePrincipalEntitlement.ServicePrincipal.Origin &&
		*args.ServicePrincipalEntitlement.ServicePrincipal.OriginId == *m.x.ServicePrincipalEntitlement.ServicePrincipal.OriginId &&
		*args.ServicePrincipalEntitlement.ServicePrincipal.PrincipalName == *m.x.ServicePrincipalEntitlement.ServicePrincipal.PrincipalName
}

func (m *matchAddServicePrincipalEntitlementArgs) String() string {
	return fmt.Sprintf("account_license_type: [%s], licensing_source: [%s], origin: [%s], origin_id: [%s], principal_name: [%s]",
		*m.x.ServicePrincipalEntitlement.AccessLevel.AccountLicenseType,
		*m.x.ServicePrincipalEntitlement.AccessLevel.LicensingSource,
		*m.x.ServicePrincipalEntitlement.ServicePrincipal.Origin,
		*m.x.ServicePrincipalEntitlement.ServicePrincipal.OriginId,
		*m.x.ServicePrincipalEntitlement.ServicePrincipal.PrincipalName)
}
