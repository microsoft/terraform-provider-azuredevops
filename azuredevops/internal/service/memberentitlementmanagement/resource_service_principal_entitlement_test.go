//go:build (all || resource_service_principal_entitlement) && !exclude_resource_service_principal_entitlement
// +build all resource_service_principal_entitlement
// +build !exclude_resource_service_principal_entitlement

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
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/licensing"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/memberentitlementmanagement"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// if the REST-API return the failure, it should fail.
func TestServicePrincipalEntitlement_CreateServicePrincipalEntitlement_WithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	originId := "837934ba-f473-45d8-ab55-83e081171cad"
	origin := "aad"

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.Set("account_license_type", "express")
	resourceData.Set("origin_id", originId)
	resourceData.Set("origin", origin)

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

	originId := "837934ba-f473-45d8-ab55-83e081171cad"
	origin := "aad"

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.Set("account_license_type", "earlyAdopter")
	resourceData.Set("origin_id", originId)
	resourceData.Set("origin", origin)

	var expectedKey interface{} = 5000
	var expectedValue interface{} = "A service principal cannot be assigned an Account-EarlyAdopter license."
	expectedErrors := []azuredevops.KeyValuePair{
		{
			Key:   &expectedKey,
			Value: &expectedValue,
		},
	}
	expectedIsSuccess := false
	operationResult := memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{
		IsSuccess: &expectedIsSuccess,
		Errors:    &expectedErrors,
	}

	// No error but it has an error on the response.
	memberEntitlementClient.
		EXPECT().
		AddServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.ServicePrincipalEntitlementOperationReference{
			Results: &[]memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{operationResult},
		}, nil).
		Times(1)

	err := resourceServicePrincipalEntitlementCreate(resourceData, clients)
	require.Contains(t, err.Error(), "A service principal cannot be assigned an Account-EarlyAdopter license.")
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
	principalName := "[contoso]\\PrincipalName"
	displayName := "displayName"
	descriptor := "baz"
	id := uuid.New()
	mockServicePrincipalEntitlement := getMockServicePrincipalEntitlement(&id, accountLicenseType, origin, originID, principalName, displayName, descriptor)
	expectedIsSuccess := true
	operationResult := memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{
		IsSuccess: &expectedIsSuccess,
		Result:    mockServicePrincipalEntitlement,
	}

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
		Return(&memberentitlementmanagement.ServicePrincipalEntitlementOperationReference{
			Results: &[]memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{operationResult},
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
	resourceData.Set("displayName", displayName)
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
	principalName := "[contoso]\\PrinicipalName"
	displayName := "displayName"
	descriptor := "baz"
	id := uuid.New()
	mockServicePrincipalEntitlement := getMockServicePrincipalEntitlement(&id, accountLicenseType, origin, originID, principalName, displayName, descriptor)
	expectedIsSuccess := true
	operationResult := memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{
		IsSuccess: &expectedIsSuccess,
		Result:    mockServicePrincipalEntitlement,
	}

	memberEntitlementClient.
		EXPECT().
		AddServicePrincipalEntitlement(gomock.Any(), MatchAddServicePrincipalEntitlementArgs(t, memberentitlementmanagement.AddServicePrincipalEntitlementArgs{
			ServicePrincipalEntitlement: mockServicePrincipalEntitlement,
		})).
		Return(&memberentitlementmanagement.ServicePrincipalEntitlementOperationReference{
			Results: &[]memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{operationResult},
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetServicePrincipalEntitlement(gomock.Any(), memberentitlementmanagement.GetServicePrincipalEntitlementArgs{
			ServicePrincipalId: mockServicePrincipalEntitlement.Id,
		}).
		Return(mockServicePrincipalEntitlement, nil)

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.Set("display_name", displayName)
	resourceData.Set("account_license_type", "basic")

	err := resourceServicePrincipalEntitlementCreate(resourceData, clients)
	assert.Nil(t, err, "err should be nil")
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

	id := uuid.New()
	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.SetId(id.String())

	mockServicePrincipalEntitlement := getMockServicePrincipalEntitlement(&id, "", "", "", "", "", "")
	memberEntitlementClient.
		EXPECT().
		GetServicePrincipalEntitlement(gomock.Any(), memberentitlementmanagement.GetServicePrincipalEntitlementArgs{
			ServicePrincipalId: mockServicePrincipalEntitlement.Id,
		}).
		Return(mockServicePrincipalEntitlement, nil)

	d, err := importServicePrincipalEntitlement(resourceData, clients)
	assert.Nil(t, err)
	assert.NotNil(t, d)
	assert.Len(t, d, 1)
	assert.Equal(t, id.String(), d[0].Id())
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
	assert.Contains(t, err.Error(), "Only UUID values can used for import")
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

	operationResult := memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{
		IsSuccess:          &expectedIsSuccess,
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
		Result: nil,
	}

	memberEntitlementClient.
		EXPECT().
		AddServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.ServicePrincipalEntitlementOperationReference{
			Results: &[]memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{operationResult},
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.Set("principal_name", "[contoso]\\Test")

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
	operationResult := memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{
		IsSuccess:          &expectedIsSuccess,
		ServicePrincipalId: &id,
		Errors:             nil,
		Result:             nil,
	}

	memberEntitlementClient.
		EXPECT().
		AddServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.ServicePrincipalEntitlementOperationReference{
			Results: &[]memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{operationResult},
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.Set("principal_name", "[contoso]\\PrincipalName")

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

	operationResult := memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{
		IsSuccess:          &expectedIsSuccess,
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
		Result: nil,
	}

	memberEntitlementClient.
		EXPECT().
		UpdateServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.ServicePrincipalEntitlementOperationReference{
			Results: &[]memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{operationResult},
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.SetId(id.String())
	resourceData.Set("principal_name", "[contoso]\\PrincipalName")

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
	operationResult := memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{
		IsSuccess:          &expectedIsSuccess,
		Errors:             nil,
		Result:             nil,
		ServicePrincipalId: &id,
	}

	memberEntitlementClient.
		EXPECT().
		UpdateServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.ServicePrincipalEntitlementOperationReference{
			Results: &[]memberentitlementmanagement.ServicePrincipalEntitlementOperationResult{operationResult},
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetServicePrincipalEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, ResourceServicePrincipalEntitlement().Schema, nil)
	resourceData.SetId(id.String())
	resourceData.Set("principal_name", "[contoso]\\PrincipalName")

	err := resourceServicePrincipalEntitlementUpdate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	assert.Contains(t, err.Error(), "Unknown API error")
}

func getMockServicePrincipalEntitlement(id *uuid.UUID, accountLicenseType licensing.AccountLicenseType, origin string, originID string, principalName string, displayName string, descriptor string) *memberentitlementmanagement.ServicePrincipalEntitlement {
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
			DisplayName:   &displayName,
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
	m.t.Logf("MatchAddServicePrincipalEntitlementArgs:\nVALUE: account_license_type: [%s], licensing_source: [%s], origin: [%s], origin_id: [%s], display_name: [%s]\n  REF: account_license_type: [%s], licensing_source: [%s], origin: [%s], origin_id: [%s], display_name: [%s], principal_name: [%s]\n",
		*args.ServicePrincipalEntitlement.AccessLevel.AccountLicenseType,
		*args.ServicePrincipalEntitlement.AccessLevel.LicensingSource,
		*args.ServicePrincipalEntitlement.ServicePrincipal.Origin,
		*args.ServicePrincipalEntitlement.ServicePrincipal.OriginId,
		*args.ServicePrincipalEntitlement.ServicePrincipal.DisplayName,
		*m.x.ServicePrincipalEntitlement.AccessLevel.AccountLicenseType,
		*m.x.ServicePrincipalEntitlement.AccessLevel.LicensingSource,
		*m.x.ServicePrincipalEntitlement.ServicePrincipal.Origin,
		*m.x.ServicePrincipalEntitlement.ServicePrincipal.OriginId,
		*m.x.ServicePrincipalEntitlement.ServicePrincipal.DisplayName,
		*m.x.ServicePrincipalEntitlement.ServicePrincipal.PrincipalName)

	return *args.ServicePrincipalEntitlement.AccessLevel.AccountLicenseType == *m.x.ServicePrincipalEntitlement.AccessLevel.AccountLicenseType &&
		*args.ServicePrincipalEntitlement.ServicePrincipal.Origin == *m.x.ServicePrincipalEntitlement.ServicePrincipal.Origin &&
		*args.ServicePrincipalEntitlement.ServicePrincipal.OriginId == *m.x.ServicePrincipalEntitlement.ServicePrincipal.OriginId
}

func (m *matchAddServicePrincipalEntitlementArgs) String() string {
	return fmt.Sprintf("account_license_type: [%s], licensing_source: [%s], origin: [%s], origin_id: [%s], display_name: [%s]",
		*m.x.ServicePrincipalEntitlement.AccessLevel.AccountLicenseType,
		*m.x.ServicePrincipalEntitlement.AccessLevel.LicensingSource,
		*m.x.ServicePrincipalEntitlement.ServicePrincipal.Origin,
		*m.x.ServicePrincipalEntitlement.ServicePrincipal.OriginId,
		*m.x.ServicePrincipalEntitlement.ServicePrincipal.DisplayName)
}
