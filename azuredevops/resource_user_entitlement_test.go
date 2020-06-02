// +build all resource_user_entitlement
// +build !exclude_resource_user_entitlement

package azuredevops

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/licensing"
	"github.com/microsoft/azure-devops-go-api/azuredevops/memberentitlementmanagement"
	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// if origin_id is provided, it will be used. if principal_name is also supplied, an error will be reported.
func TestUserEntitlement_CreateUserEntitlement_DoNotAllowToSetOridinIdAndPrincipalName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &config.AggregatedClient{
		MemberEntitleManagementClient: nil,
		Ctx:                           context.Background(),
	}

	originID := "e97b0e7f-0a61-41ad-860c-748ec5fcb20b"
	principalName := "foobar@microsoft.com"

	resourceData := schema.TestResourceDataRaw(t, resourceUserEntitlement().Schema, nil)
	resourceData.Set("origin_id", originID)
	resourceData.Set("principal_name", principalName)

	err := resourceUserEntitlementCreate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	require.Regexp(t, "Both origin_id and principal_name set. You can not use both", err.Error())
}

// if origin_id is "" and principal_name is supplied, the principal_name will be used.
func TestUserEntitlement_CreateUserEntitlement_WithPrincipalName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &config.AggregatedClient{
		MemberEntitleManagementClient: client,
		Ctx:                           context.Background(),
	}

	accountLicenseType := licensing.AccountLicenseTypeValues.Express
	origin := ""
	originID := ""
	principalName := "foobar@microsoft.com"
	descriptor := "baz"
	id := uuid.New()
	mockUserEntitlement := getMockUserEntitlement(&id, accountLicenseType, origin, originID, principalName, descriptor)

	resourceData := schema.TestResourceDataRaw(t, resourceUserEntitlement().Schema, nil)
	resourceData.Set("principal_name", principalName)

	expectedIsSuccess := true
	client.
		EXPECT().
		AddUserEntitlement(gomock.Any(), MatchAddUserEntitlementArgs(t, memberentitlementmanagement.AddUserEntitlementArgs{
			UserEntitlement: mockUserEntitlement,
		})).
		Return(&memberentitlementmanagement.UserEntitlementsPostResponse{
			IsSuccess:       &expectedIsSuccess,
			UserEntitlement: mockUserEntitlement,
		}, nil).
		Times(1)

	client.
		EXPECT().
		GetUserEntitlement(gomock.Any(), memberentitlementmanagement.GetUserEntitlementArgs{
			UserId: mockUserEntitlement.Id,
		}).
		Return(mockUserEntitlement, nil)

	err := resourceUserEntitlementCreate(resourceData, clients)
	assert.Nil(t, err, "err should not be nil")
}

// if origin_id is "" and principal_name is "", an error will be reported.
func TestUserEntitlement_CreateUserEntitlement_Need_OriginID_Or_PrincipalName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clients := &config.AggregatedClient{
		MemberEntitleManagementClient: nil,
		Ctx:                           context.Background(),
	}

	resourceData := schema.TestResourceDataRaw(t, resourceUserEntitlement().Schema, nil)
	// originID and principalName is not set.

	err := resourceUserEntitlementCreate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	require.Regexp(t, "Use origin_id or principal_name", err.Error())
}

// if the REST-API return the failure, it should fail.

func TestUserEntitlement_CreateUserEntitlement_WithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &config.AggregatedClient{
		MemberEntitleManagementClient: client,
		Ctx:                           context.Background(),
	}

	principalName := "foobar@microsoft.com"

	resourceData := schema.TestResourceDataRaw(t, resourceUserEntitlement().Schema, nil)
	// resourceData.Set("origin_id", originID)
	resourceData.Set("account_license_type", "express")
	resourceData.Set("principal_name", principalName)

	// No error but it has a error on the response.
	client.
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

	client := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &config.AggregatedClient{
		MemberEntitleManagementClient: client,
		Ctx:                           context.Background(),
	}

	principalName := "foobar@microsoft.com"

	resourceData := schema.TestResourceDataRaw(t, resourceUserEntitlement().Schema, nil)
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
	client.
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

	client := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &config.AggregatedClient{
		MemberEntitleManagementClient: client,
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

	client.
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

	client.
		EXPECT().
		GetUserEntitlement(gomock.Any(), memberentitlementmanagement.GetUserEntitlementArgs{
			UserId: mockUserEntitlement.Id,
		}).
		Return(mockUserEntitlement, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, resourceUserEntitlement().Schema, nil)
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

	client := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &config.AggregatedClient{
		MemberEntitleManagementClient: client,
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

	client.
		EXPECT().
		AddUserEntitlement(gomock.Any(), MatchAddUserEntitlementArgs(t, memberentitlementmanagement.AddUserEntitlementArgs{
			UserEntitlement: mockUserEntitlement,
		})).
		Return(&memberentitlementmanagement.UserEntitlementsPostResponse{
			IsSuccess:       &expectedIsSuccess,
			UserEntitlement: mockUserEntitlement,
		}, nil).
		Times(1)

	client.
		EXPECT().
		GetUserEntitlement(gomock.Any(), memberentitlementmanagement.GetUserEntitlementArgs{
			UserId: mockUserEntitlement.Id,
		}).
		Return(mockUserEntitlement, nil)

	resourceData := schema.TestResourceDataRaw(t, resourceUserEntitlement().Schema, nil)
	resourceData.Set("principal_name", principalName)
	resourceData.Set("account_license_type", "basic")

	err := resourceUserEntitlementCreate(resourceData, clients)
	assert.Nil(t, err, "err should be nil")
}

// TestUserEntitlement_Import_TestUPN tests if import is successful using an UPN
func TestUserEntitlement_Import_TestUPN(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &config.AggregatedClient{
		MemberEntitleManagementClient: client,
		Ctx:                           context.Background(),
	}

	accountLicenseType := licensing.AccountLicenseTypeValues.Express
	origin := ""
	originID := ""
	principalName := "foobar@microsoft.com"
	descriptor := "baz"
	id := uuid.New()
	mockUserEntitlement := getMockUserEntitlement(&id, accountLicenseType, origin, originID, principalName, descriptor)

	client.
		EXPECT().
		GetUserEntitlements(gomock.Any(), gomock.Any()).
		Return(&memberentitlementmanagement.PagedGraphMemberList{
			Members: &[]memberentitlementmanagement.UserEntitlement{
				*mockUserEntitlement,
			},
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, resourceUserEntitlement().Schema, nil)
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

	client := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &config.AggregatedClient{
		MemberEntitleManagementClient: client,
		Ctx:                           context.Background(),
	}

	id := uuid.New().String()
	resourceData := schema.TestResourceDataRaw(t, resourceUserEntitlement().Schema, nil)
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

	client := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &config.AggregatedClient{
		MemberEntitleManagementClient: client,
		Ctx:                           context.Background(),
	}

	id := "InvalidValue-a73c5191-e20d"
	resourceData := schema.TestResourceDataRaw(t, resourceUserEntitlement().Schema, nil)
	resourceData.SetId(id)

	d, err := importUserEntitlement(resourceData, clients)
	assert.Nil(t, d)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Only UUID and UPN values can used for import")
}

func TestUserEntitlement_Create_TestErrorFormatting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &config.AggregatedClient{
		MemberEntitleManagementClient: client,
		Ctx:                           context.Background(),
	}

	id, _ := uuid.NewUUID()
	expectedIsSuccess := false
	k1 := interface{}("9999")
	v1 := interface{}("Error1")
	k2 := interface{}("9998")
	v2 := interface{}("Error2")

	client.
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

	client.
		EXPECT().
		GetUserEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, resourceUserEntitlement().Schema, nil)
	resourceData.Set("principal_name", "foobar@microsoft.com")

	err := resourceUserEntitlementCreate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	assert.Contains(t, err.Error(), "(9999) Error1")
	assert.Contains(t, err.Error(), "(9998) Error2")
}

func TestUserEntitlement_Create_TestEmptyErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &config.AggregatedClient{
		MemberEntitleManagementClient: client,
		Ctx:                           context.Background(),
	}

	id, _ := uuid.NewUUID()
	expectedIsSuccess := false

	client.
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

	client.
		EXPECT().
		GetUserEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, resourceUserEntitlement().Schema, nil)
	resourceData.Set("principal_name", "foobar@microsoft.com")

	err := resourceUserEntitlementCreate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
	assert.Contains(t, err.Error(), "Unknown API error")
}

func TestUserEntitlement_Update_TestErrorFormatting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &config.AggregatedClient{
		MemberEntitleManagementClient: client,
		Ctx:                           context.Background(),
	}

	id, _ := uuid.NewUUID()
	expectedIsSuccess := false
	k1 := interface{}("9999")
	v1 := interface{}("Error1")
	k2 := interface{}("9998")
	v2 := interface{}("Error2")

	client.
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

	client.
		EXPECT().
		GetUserEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, resourceUserEntitlement().Schema, nil)
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

	client := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &config.AggregatedClient{
		MemberEntitleManagementClient: client,
		Ctx:                           context.Background(),
	}

	id, _ := uuid.NewUUID()
	expectedIsSuccess := false

	client.
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

	client.
		EXPECT().
		GetUserEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(0)

	resourceData := schema.TestResourceDataRaw(t, resourceUserEntitlement().Schema, nil)
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

// Acceptance Test Patterns
// Create operation with AzDo account (origin_id)
// Create operation with AzDo account (principal_name)
// Create operation with GitHub account (origin_id)
// Create operation with GitHub account (principal_name)

func TestAccUserEntitlement_Create(t *testing.T) {
	tfNode := "azuredevops_user_entitlement.user"
	principalName := os.Getenv("AZDO_TEST_AAD_USER_EMAIL")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testhelper.TestAccPreCheck(t, &[]string{"AZDO_TEST_AAD_USER_EMAIL"}) },
		Providers:    testAccProviders,
		CheckDestroy: testAccUserEntitlementCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testhelper.TestAccUserEntitlementResource(principalName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserEntitlementResourceExists(principalName),
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
				),
			},
		},
	})
}

// Given the principalName of an AzDO userEntitlement, this will return a function that will check whether
// or not the userEntitlement (1) exists in the state and (2) exist in AzDO and (3) has the correct name
func testAccCheckUserEntitlementResourceExists(expectedPrincipalName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources["azuredevops_user_entitlement.user"]
		if !ok {
			return fmt.Errorf("Did not find a UserEntitlement in the TF state")
		}

		clients := testAccProvider.Meta().(*config.AggregatedClient)
		id, err := uuid.Parse(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing UserEntitlement ID, got %s: %v", resource.Primary.ID, err)
		}

		userEntitlement, err := readUserEntitlement(clients, &id)

		if err != nil {
			return fmt.Errorf("UserEntitlement with ID=%s cannot be found!. Error=%v", id, err)
		}

		if !strings.EqualFold(strings.ToLower(*userEntitlement.User.PrincipalName), strings.ToLower(expectedPrincipalName)) {
			return fmt.Errorf("UserEntitlement with ID=%s has PrincipalName=%s, but expected Name=%s", resource.Primary.ID, *userEntitlement.User.PrincipalName, expectedPrincipalName)
		}

		return nil
	}
}

// verifies that all projects referenced in the state are destroyed. This will be invoked
// *after* terraform destroys the resource but *before* the state is wiped clean.
func testAccUserEntitlementCheckDestroy(s *terraform.State) error {
	clients := testAccProvider.Meta().(*config.AggregatedClient)

	//verify that every users referenced in the state does not exist in AzDO
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_user_entitlement" {
			continue
		}

		id, err := uuid.Parse(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing UserEntitlement ID, got %s: %v", resource.Primary.ID, err)
		}

		userEntitlement, err := readUserEntitlement(clients, &id)
		if err != nil {
			if utils.ResponseWasNotFound(err) {
				return nil
			}
			return fmt.Errorf("Bad: Get UserEntitlment :  %+v", err)
		}

		if userEntitlement != nil && userEntitlement.AccessLevel != nil && string(*userEntitlement.AccessLevel.Status) != "none" {
			return fmt.Errorf("Status should be none : %s with readUserEntitlement error %v", string(*userEntitlement.AccessLevel.Status), err)
		}
	}

	return nil
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

func init() {
	InitProvider()
}
