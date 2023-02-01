//go:build (all || resource_group_entitlement) && !exclude_resource_group_entitlement
// +build all resource_group_entitlement
// +build !exclude_resource_group_entitlement

package memberentitlementmanagement

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/licensing"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/memberentitlementmanagement"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/webapi"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// if origin_id is provided, it will be used. if principal_name is also supplied, an error will be reported.
func TestGroupEntitlement_CreateGroupEntitlement_DoNotAllowToSetOriginIdAndPrincipalName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: nil,
		Ctx:                           context.Background(),
	}

	originID := "e97b0e7f-0a61-41ad-860c-748ec5fcb20b"
	principalName := "[contoso]\\PrincipalName"

	resourceData := schema.TestResourceDataRaw(t, ResourceGroupEntitlement().Schema, nil)
	resourceData.Set("origin_id", originID)
	resourceData.Set("principal_name", principalName)

	err := resourceGroupEntitlementCreate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	require.Regexp(t, "Both origin_id and principal_name set. You can not use both", err.Error())
}

// if origin_id is "" and principal_name is supplied, the principal_name will be used.
func TestGroupEntitlement_CreateGroupEntitlement_WithPrincipalName(t *testing.T) {
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
	principalName := "[contoso]\\PrincipalName"
	displayName := ""
	descriptor := "baz"
	id := uuid.New()
	mockGroupEntitlement := getMockGroupEntitlement(&id, accountLicenseType, origin, originID, principalName, displayName, descriptor)

	resourceData := schema.TestResourceDataRaw(t, ResourceGroupEntitlement().Schema, nil)
	resourceData.Set("principal_name", principalName)
	expectedIsSuccess := true
	operationResult := memberentitlementmanagement.GroupOperationResult{
		IsSuccess: &expectedIsSuccess,
		Result:    mockGroupEntitlement,
	}
	memberEntitlementClient.
		EXPECT().
		AddGroupEntitlement(gomock.Any(), MatchAddGroupEntitlementArgs(t, memberentitlementmanagement.AddGroupEntitlementArgs{
			GroupEntitlement: mockGroupEntitlement,
		})).
		Return(&memberentitlementmanagement.GroupEntitlementOperationReference{
			Results: &[]memberentitlementmanagement.GroupOperationResult{operationResult},
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetGroupEntitlement(gomock.Any(), memberentitlementmanagement.GetGroupEntitlementArgs{
			GroupId: mockGroupEntitlement.Id,
		}).
		Return(mockGroupEntitlement, nil)

	err := resourceGroupEntitlementCreate(resourceData, clients)
	assert.Nil(t, err, "err should not be nil")
}

// if origin_id is "" and principal_name is "", an error will be reported.
func TestGroupEntitlement_CreateGroupEntitlement_Need_OriginID_Or_PrincipalName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: nil,
		Ctx:                           context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, ResourceGroupEntitlement().Schema, nil)
	// originID and principalName is not set.

	err := resourceGroupEntitlementCreate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	require.Regexp(t, "Use origin_id or principal_name", err.Error())
}

// if the REST-API return the failure, it should fail.

func TestGroupEntitlement_CreateGroupEntitlement_WithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	principalName := "[contoso]\\PrincipalName"

	resourceData := schema.TestResourceDataRaw(t, ResourceGroupEntitlement().Schema, nil)
	// resourceData.Set("origin_id", originID)
	resourceData.Set("account_license_type", "express")
	resourceData.Set("principal_name", principalName)

	// No error but it has a error on the response.
	memberEntitlementClient.
		EXPECT().
		AddGroupEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("error foo")).
		Times(1)

	err := resourceGroupEntitlementCreate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
}

// if the REST-API return the success, but fails on response
func TestGroupEntitlement_CreateGroupEntitlement_WithEarlyAdopter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	principalName := "[contoso]\\PrincipalName"

	resourceData := schema.TestResourceDataRaw(t, ResourceGroupEntitlement().Schema, nil)
	// resourceData.Set("origin_id", originID)
	resourceData.Set("account_license_type", "earlyAdopter")
	resourceData.Set("principal_name", principalName)

	var expectedKey interface{} = 5000
	var expectedValue interface{} = "A group cannot be assigned an Account-EarlyAdopter license."
	expectedErrors := []azuredevops.KeyValuePair{
		{
			Key:   &expectedKey,
			Value: &expectedValue,
		},
	}
	expectedIsSuccess := false
	operationResult := memberentitlementmanagement.GroupOperationResult{
		IsSuccess: &expectedIsSuccess,
		Errors:    &expectedErrors,
	}

	// No error but it has an error on the response.
	memberEntitlementClient.
		EXPECT().
		AddGroupEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.GroupEntitlementOperationReference{
			Results: &[]memberentitlementmanagement.GroupOperationResult{operationResult},
		}, nil).
		Times(1)

	err := resourceGroupEntitlementCreate(resourceData, clients)
	require.Contains(t, err.Error(), "A group cannot be assigned an Account-EarlyAdopter license.")
}

// TestGroupEntitlement_Update_TestChangeEntitlement verfies that an entitlement can be changed
func TestGroupEntitlement_Update_TestChangeEntitlement(t *testing.T) {
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
	displayName := ""
	descriptor := "baz"
	id := uuid.New()
	mockGroupEntitlement := getMockGroupEntitlement(&id, accountLicenseType, origin, originID, principalName, displayName, descriptor)
	expectedIsSuccess := true
	operationResult := memberentitlementmanagement.GroupOperationResult{
		IsSuccess: &expectedIsSuccess,
		Result:    mockGroupEntitlement,
	}

	memberEntitlementClient.
		EXPECT().
		UpdateGroupEntitlement(gomock.Any(), memberentitlementmanagement.UpdateGroupEntitlementArgs{
			GroupId: &id,
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
		Return(&memberentitlementmanagement.GroupEntitlementOperationReference{
			Results: &[]memberentitlementmanagement.GroupOperationResult{operationResult},
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetGroupEntitlement(gomock.Any(), memberentitlementmanagement.GetGroupEntitlementArgs{
			GroupId: mockGroupEntitlement.Id,
		}).
		Return(mockGroupEntitlement, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceGroupEntitlement().Schema, nil)
	resourceData.SetId(id.String())
	resourceData.Set("principal_name", principalName)
	resourceData.Set("account_license_type", string(licensing.AccountLicenseTypeValues.Stakeholder))
	resourceData.Set("licensing_source", string(licensing.LicensingSourceValues.Account))

	err := resourceGroupEntitlementUpdate(resourceData, clients)
	assert.Nil(t, err)
}

// TestGroupEntitlement_CreateUpdate_TestBasicEntitlement verifies that the (virtual) Basic entitlement can be set
func TestGroupEntitlement_CreateUpdate_TestBasicEntitlement(t *testing.T) {
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
	displayName := ""
	descriptor := "baz"
	id := uuid.New()
	mockGroupEntitlement := getMockGroupEntitlement(&id, accountLicenseType, origin, originID, principalName, displayName, descriptor)
	expectedIsSuccess := true
	operationResult := memberentitlementmanagement.GroupOperationResult{
		IsSuccess: &expectedIsSuccess,
		Result:    mockGroupEntitlement,
	}

	memberEntitlementClient.
		EXPECT().
		AddGroupEntitlement(gomock.Any(), MatchAddGroupEntitlementArgs(t, memberentitlementmanagement.AddGroupEntitlementArgs{
			GroupEntitlement: mockGroupEntitlement,
		})).
		Return(&memberentitlementmanagement.GroupEntitlementOperationReference{
			Results: &[]memberentitlementmanagement.GroupOperationResult{operationResult},
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetGroupEntitlement(gomock.Any(), memberentitlementmanagement.GetGroupEntitlementArgs{
			GroupId: mockGroupEntitlement.Id,
		}).
		Return(mockGroupEntitlement, nil)

	resourceData := schema.TestResourceDataRaw(t, ResourceGroupEntitlement().Schema, nil)
	resourceData.Set("principal_name", principalName)
	resourceData.Set("account_license_type", "basic")

	err := resourceGroupEntitlementCreate(resourceData, clients)
	assert.Nil(t, err, "err should be nil")
}

// TestGroupEntitlement_Import_TestID tests if import is successful using an UUID
func TestGroupEntitlement_Import_TestID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	id := uuid.New().String()
	resourceData := schema.TestResourceDataRaw(t, ResourceGroupEntitlement().Schema, nil)
	resourceData.SetId(id)

	d, err := importGroupEntitlement(resourceData, clients)
	assert.Nil(t, err)
	assert.NotNil(t, d)
	assert.Len(t, d, 1)
	assert.Equal(t, id, d[0].Id())
}

// TestGroupEntitlement_Import_TestInvalidValue tests if only a valid UPN and UUID can be used to import a resource
func TestGroupEntitlement_Import_TestInvalidValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	id := "InvalidValue-a73c5191-e20d"
	resourceData := schema.TestResourceDataRaw(t, ResourceGroupEntitlement().Schema, nil)
	resourceData.SetId(id)

	d, err := importGroupEntitlement(resourceData, clients)
	assert.Nil(t, d)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Only UUID values can used for import")
}

func TestGroupEntitlement_Create_TestErrorFormatting(t *testing.T) {
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

	operationResult := memberentitlementmanagement.GroupOperationResult{
		IsSuccess: &expectedIsSuccess,
		GroupId:   &id,
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
		AddGroupEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.GroupEntitlementOperationReference{
			Results: &[]memberentitlementmanagement.GroupOperationResult{operationResult},
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetGroupEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, ResourceGroupEntitlement().Schema, nil)
	resourceData.Set("principal_name", "[contoso]\\Test")

	err := resourceGroupEntitlementCreate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	assert.Contains(t, err.Error(), "(9999) Error1")
	assert.Contains(t, err.Error(), "(9998) Error2")
}

func TestGroupEntitlement_Create_TestEmptyErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	id, _ := uuid.NewUUID()
	expectedIsSuccess := false
	operationResult := memberentitlementmanagement.GroupOperationResult{
		IsSuccess: &expectedIsSuccess,
		GroupId:   &id,
		Errors:    nil,
		Result:    nil,
	}

	memberEntitlementClient.
		EXPECT().
		AddGroupEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.GroupEntitlementOperationReference{
			Results: &[]memberentitlementmanagement.GroupOperationResult{operationResult},
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetGroupEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, ResourceGroupEntitlement().Schema, nil)
	resourceData.Set("principal_name", "[contoso]\\PrincipalName")

	err := resourceGroupEntitlementCreate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	assert.Contains(t, err.Error(), "Unknown API error")
}

func TestGroupEntitlement_Update_TestErrorFormatting(t *testing.T) {
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

	operationResult := memberentitlementmanagement.GroupOperationResult{
		IsSuccess: &expectedIsSuccess,
		GroupId:   &id,
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
		UpdateGroupEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.GroupEntitlementOperationReference{
			Results: &[]memberentitlementmanagement.GroupOperationResult{operationResult},
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetGroupEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, ResourceGroupEntitlement().Schema, nil)
	resourceData.SetId(id.String())
	resourceData.Set("principal_name", "[contoso]\\PrincipalName")

	err := resourceGroupEntitlementUpdate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	assert.Contains(t, err.Error(), "(9999) Error1")
	assert.Contains(t, err.Error(), "(9998) Error2")
}

func TestGroupEntitlement_Update_TestEmptyErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memberEntitlementClient := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &client.AggregatedClient{
		MemberEntitleManagementClient: memberEntitlementClient,
		Ctx:                           context.Background(),
	}

	id, _ := uuid.NewUUID()
	expectedIsSuccess := false
	operationResult := memberentitlementmanagement.GroupOperationResult{
		IsSuccess: &expectedIsSuccess,
		Errors:    nil,
		Result:    nil,
		GroupId:   &id,
	}

	memberEntitlementClient.
		EXPECT().
		UpdateGroupEntitlement(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.GroupEntitlementOperationReference{
			Results: &[]memberentitlementmanagement.GroupOperationResult{operationResult},
		}, nil).
		Times(1)

	memberEntitlementClient.
		EXPECT().
		GetGroupEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, ResourceGroupEntitlement().Schema, nil)
	resourceData.SetId(id.String())
	resourceData.Set("principal_name", "[contoso]\\PrincipalName")

	err := resourceGroupEntitlementUpdate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	assert.Contains(t, err.Error(), "Unknown API error")
}

func getMockGroupEntitlement(id *uuid.UUID, accountLicenseType licensing.AccountLicenseType, origin string, originID string, principalName string, displayName string, descriptor string) *memberentitlementmanagement.GroupEntitlement {
	subjectKind := "group"
	licensingSource := licensing.LicensingSourceValues.Account

	return &memberentitlementmanagement.GroupEntitlement{
		LicenseRule: &licensing.AccessLevel{
			AccountLicenseType: &accountLicenseType,
			LicensingSource:    &licensingSource,
		},
		Id: id,
		Group: &graph.GraphGroup{
			Origin:        &origin,
			OriginId:      &originID,
			PrincipalName: &principalName,
			DisplayName:   &displayName,
			SubjectKind:   &subjectKind,
			Descriptor:    &descriptor,
		},
	}
}

type matchAddGroupEntitlementArgs struct {
	t *testing.T
	x memberentitlementmanagement.AddGroupEntitlementArgs
}

func MatchAddGroupEntitlementArgs(t *testing.T, x memberentitlementmanagement.AddGroupEntitlementArgs) gomock.Matcher {
	return &matchAddGroupEntitlementArgs{t, x}
}

func (m *matchAddGroupEntitlementArgs) Matches(x interface{}) bool {
	args := x.(memberentitlementmanagement.AddGroupEntitlementArgs)
	m.t.Logf("MatchAddGroupEntitlementArgs:\nVALUE: account_license_type: [%s], licensing_source: [%s], origin: [%s], origin_id: [%s], display_name: [%s], principal_name: [%s]\n  REF: account_license_type: [%s], licensing_source: [%s], origin: [%s], origin_id: [%s], display_name: [%s], principal_name: [%s]\n",
		*args.GroupEntitlement.LicenseRule.AccountLicenseType,
		*args.GroupEntitlement.LicenseRule.LicensingSource,
		*args.GroupEntitlement.Group.Origin,
		*args.GroupEntitlement.Group.OriginId,
		*args.GroupEntitlement.Group.DisplayName,
		*args.GroupEntitlement.Group.PrincipalName,
		*m.x.GroupEntitlement.LicenseRule.AccountLicenseType,
		*m.x.GroupEntitlement.LicenseRule.LicensingSource,
		*m.x.GroupEntitlement.Group.Origin,
		*m.x.GroupEntitlement.Group.OriginId,
		*m.x.GroupEntitlement.Group.DisplayName,
		*m.x.GroupEntitlement.Group.PrincipalName)

	return *args.GroupEntitlement.LicenseRule.AccountLicenseType == *m.x.GroupEntitlement.LicenseRule.AccountLicenseType &&
		*args.GroupEntitlement.Group.Origin == *m.x.GroupEntitlement.Group.Origin &&
		*args.GroupEntitlement.Group.OriginId == *m.x.GroupEntitlement.Group.OriginId &&
		*args.GroupEntitlement.Group.PrincipalName == *m.x.GroupEntitlement.Group.PrincipalName
}

func (m *matchAddGroupEntitlementArgs) String() string {
	return fmt.Sprintf("account_license_type: [%s], licensing_source: [%s], origin: [%s], origin_id: [%s], display_name: [%s], principal_name: [%s]",
		*m.x.GroupEntitlement.LicenseRule.AccountLicenseType,
		*m.x.GroupEntitlement.LicenseRule.LicensingSource,
		*m.x.GroupEntitlement.Group.Origin,
		*m.x.GroupEntitlement.Group.OriginId,
		*m.x.GroupEntitlement.Group.DisplayName,
		*m.x.GroupEntitlement.Group.PrincipalName)
}
