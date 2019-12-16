// +build all resource_user_entitlement

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
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// if origin_id is provided, it will be used. if principal_name is also supplied, an error will be reported.
func TestAzureDevOpsUserEntitlement_CreateUserEntitlement_DoNotAllowToSetOridinIdAndPrincipalName(t *testing.T) {
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
	require.Regexp(t, "Error creating user entitlement: Error both origin_id and principal_name set. You can not use both", err.Error())
}

// if origin_id is "" and principal_name is supplied, the principal_name will be used.
func TestAzureDevOpsUserEntitlement_CreateUserEntitlement_WithPrincipalName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := azdosdkmocks.NewMockMemberentitlementmanagementClient(ctrl)
	clients := &config.AggregatedClient{
		MemberEntitleManagementClient: client,
		Ctx:                           context.Background(),
	}

	accountLicenseType := licensing.AccountLicenseTypeValues.Express
	origin := "aad" // Default
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
		AddUserEntitlement(gomock.Any(), MatchAddUserEntitlementArgs(memberentitlementmanagement.AddUserEntitlementArgs{
			UserEntitlement: mockUserEntitlement,
		})).
		Return(&memberentitlementmanagement.UserEntitlementsPostResponse{
			IsSuccess:       &expectedIsSuccess,
			UserEntitlement: mockUserEntitlement,
		}, nil).
		Times(1)
	client.EXPECT().GetUserEntitlement(gomock.Any(), memberentitlementmanagement.GetUserEntitlementArgs{
		UserId: mockUserEntitlement.Id,
	}).Return(mockUserEntitlement, nil)

	err := resourceUserEntitlementCreate(resourceData, clients)
	assert.Nil(t, err, "err should not be nil")
}

// if origin_id is "" and principal_name is "", an error will be reported.
func TestAzureDevOpsUserEntitlement_CreateUserEntitlement_Need_OriginID_Or_PrincipalName(t *testing.T) {
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

func TestAzureDevOpsUserEntitlement_CreateUserEntitlement_WithError(t *testing.T) {
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

	// No error but it has a error on the reponse.
	client.
		EXPECT().
		AddUserEntitlement(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("error foo")).
		Times(1)

	err := resourceUserEntitlementCreate(resourceData, clients)
	assert.NotNil(t, err, "err should not be nil")
}

// if the REST-API return the success, but fails on response
func TestAzureDevOpsUserEntitlement_CreateUserEntitlement_WithEarlyAdopter(t *testing.T) {
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

	// No error but it has a error on the reponse.
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

func getMockUserEntitlement(id *uuid.UUID, accountLicenseType licensing.AccountLicenseType, origin string, originID string, principalName string, descriptor string) *memberentitlementmanagement.UserEntitlement {
	subjectKind := "user"
	return &memberentitlementmanagement.UserEntitlement{
		AccessLevel: &licensing.AccessLevel{
			AccountLicenseType: &accountLicenseType,
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

func TestAccAzureDevOpsUserEntitlement_Create(t *testing.T) {
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
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
					testAccCheckUserEntitlementResourceExists(principalName),
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

		if strings.ToLower(*userEntitlement.User.PrincipalName) != strings.ToLower(expectedPrincipalName) {
			return fmt.Errorf("UserEntitlement with ID=%s has PrincipalName=%s, but expected Name=%s", resource.Primary.ID, *userEntitlement.User.PrincipalName, expectedPrincipalName)
		}

		return nil
	}

}

// verifies that all projects referenced in the state are destroyed. This will be invoked
// *after* terraform destroys the resource but *before* the state is wiped clean.
func testAccUserEntitlementCheckDestroy(s *terraform.State) error {
	clients := testAccProvider.Meta().(*config.AggregatedClient)

	// verify that every users referenced in the state does not exist in AzDO
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
			if userEntitlement != nil && userEntitlement.AccessLevel != nil && string(*userEntitlement.AccessLevel.Status) != "none" {
				return fmt.Errorf("Status should be none : %s with readUserEntitlement error %v", string(*userEntitlement.AccessLevel.Status), err)
			}
		}
		if string(*userEntitlement.AccessLevel.Status) != "none" {
			return fmt.Errorf("Status should be none : %s", string(*userEntitlement.AccessLevel.Status))
		}
	}

	return nil
}

type matchAddUserEntitlementArgs struct {
	t memberentitlementmanagement.AddUserEntitlementArgs
}

func MatchAddUserEntitlementArgs(t memberentitlementmanagement.AddUserEntitlementArgs) gomock.Matcher {
	return &matchAddUserEntitlementArgs{t}
}

func (m *matchAddUserEntitlementArgs) Matches(x interface{}) bool {
	args := x.(memberentitlementmanagement.AddUserEntitlementArgs)
	return *args.UserEntitlement.AccessLevel.AccountLicenseType == *m.t.UserEntitlement.AccessLevel.AccountLicenseType &&
		*args.UserEntitlement.User.Origin == *m.t.UserEntitlement.User.Origin &&
		*args.UserEntitlement.User.OriginId == *m.t.UserEntitlement.User.OriginId &&
		*args.UserEntitlement.User.PrincipalName == *m.t.UserEntitlement.User.PrincipalName
}

func (m *matchAddUserEntitlementArgs) String() string {
	return fmt.Sprintf("origin_id: %s, principal_name: %s", *m.t.UserEntitlement.User.OriginId, *m.t.UserEntitlement.User.PrincipalName)
}

func init() {
	InitProvider()
}
