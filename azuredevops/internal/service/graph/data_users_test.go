// +build all core data_sources data_users
// +build !exclude_data_sources !exclude_data_users

package graph

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/ahmetb/go-linq"
	"github.com/terraform-providers/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/converter"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/stretchr/testify/require"
)

var usrList1 = []graph.GraphUser{
	{
		Descriptor:    converter.String("aad.YWU4YzJkMTktOThmYS03ZDhmLWJhNTAtOWI4MWQzYTUxZjcy"),
		DisplayName:   converter.String("Desiree M. Collins"),
		PrincipalName: converter.String("DesireeMCollins@jourrapide.com"),
		Origin:        converter.String("aad"),
		OriginId:      converter.String("bc522d83-7192-4fad-b885-5e1334da4f94"),
	},
	{
		Descriptor:    converter.String("svc.Nzc3NGFjMDMtOGEyOS00NGFjLTg2ZjEtZmE0YmRlZDc4ZGUyOkdpdEh1YiBBcHA6MjE3OGVhMGItMmNlMi00ZDA4LTg4YTMtODdiYWRjYjkzNjE1"),
		DisplayName:   converter.String("GitHub"),
		PrincipalName: converter.String(""),
		Origin:        converter.String("vsts"),
		OriginId:      converter.String("31d16637-81f9-46e1-b305-78f5aa7cc909"),
	},
	{
		Descriptor:    converter.String("aad.NTJmM2YzZmMtZmE4NS00ZGNhLTkwYjUtMDNiY2U0NjBlMmIy"),
		DisplayName:   converter.String("Walter M. Brooks"),
		PrincipalName: converter.String("WalterMBrooks@rhyta.com"),
		Origin:        converter.String("aad"),
		OriginId:      converter.String("8c840d92-f19e-4dfe-8eab-5a1fd67a3a77"),
	},
	{
		Descriptor:    converter.String("msa.MmRkMWZkMWMtNWE0Mi00NDM1LThhM2ItMWQ1NWUwOTA2NzUx"),
		DisplayName:   converter.String("Kerry M. Raymond"),
		PrincipalName: converter.String("KerryMRaymond@teleworm.us"),
		Origin:        converter.String("msa"),
		OriginId:      converter.String("8c840d92-f19e-4dfe-8eab-5a1fd67a3a77"),
	},
}

var usrList2 = []graph.GraphUser{
	{
		Descriptor:    converter.String("svc.Nzc3NGFjMDMtOGEyOS00NGFjLTg2ZjEtZmE0YmRlZDc4ZGUyOkJ1aWxkOmYwNzg2ZGZkLTA4OGYtNDYxOS1hY2NjLTJlYzc1ZTEyNmRjZQ"),
		DisplayName:   converter.String("Project Collection Build Service (unittest)"),
		PrincipalName: converter.String("f0786dfd-088f-4619-accc-2ec75e126dce"),
		Origin:        converter.String("vsts"),
		OriginId:      converter.String("b5f04e6c-87f5-4e19-a5cc-f930bf78dd50"),
	},
	{
		Descriptor:    converter.String("bnd.dXBuOmIwYjg0Yzc1LTI4NjctNGQ0Mi1iNWFhLTg1YjE5Nzc1ZDdlN1xBbGxlbkJNY0tpbm5vbkBkYXlyZXAuY29t"),
		DisplayName:   converter.String("AllenBMcKinnon@dayrep.com"),
		PrincipalName: converter.String("AllenBMcKinnon@dayrep.com"),
		Origin:        converter.String("aad"),
		OriginId:      nil,
	},
	{
		Descriptor:    converter.String("win.Uy0xLTUtMjEtMjA4NTA0NzkxOS0yODAyODgxNjk5LTE5MDg1MTA1MTctNTAw"),
		DisplayName:   converter.String("Local Admin"),
		PrincipalName: converter.String("vm0e8f6bad5c95\\ladmin"),
		Origin:        converter.String("ad"),
		OriginId:      converter.String(""),
	},
}

// verfies that the data source propagates an error from the API correctly
func TestDataSourceUser_Read_TestDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &client.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	expectedArgs := graph.ListUsersArgs{
		SubjectTypes: &[]string{},
	}
	graphClient.
		EXPECT().
		ListUsers(clients.Ctx, expectedArgs).
		Return(nil, errors.New("ListUsers() Failed"))

	resourceData := schema.TestResourceDataRaw(t, DataUsers().Schema, nil)
	err := dataUsersRead(resourceData, clients)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "ListUsers() Failed")
}

func TestDataSourceUser_Read_HandlesContinuationToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &client.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	var calls []*gomock.Call
	calls = append(calls, graphClient.
		EXPECT().
		ListUsers(clients.Ctx, graph.ListUsersArgs{
			SubjectTypes: &[]string{},
		}).
		Return(&graph.PagedGraphUsers{
			GraphUsers:        &usrList1,
			ContinuationToken: &[]string{"2"},
		}, nil).
		Times(1))

	calls = append(calls, graphClient.
		EXPECT().
		ListUsers(clients.Ctx, graph.ListUsersArgs{
			SubjectTypes:      &[]string{},
			ContinuationToken: converter.String("2"),
		}).
		Return(&graph.PagedGraphUsers{
			GraphUsers:        &usrList2,
			ContinuationToken: &[]string{""},
		}, nil).
		Times(1))

	gomock.InOrder(calls...)

	resourceData := schema.TestResourceDataRaw(t, DataUsers().Schema, nil)
	err := dataUsersRead(resourceData, clients)
	require.Nil(t, err)
}

// verifies that a single user can be read successfully
func TestDataSourceUser_Read_TestReadEmptyUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &client.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	/* start writing test here */
	expectedArgs := graph.ListUsersArgs{
		SubjectTypes: &[]string{},
	}
	graphClient.
		EXPECT().
		ListUsers(clients.Ctx, expectedArgs).
		Return(&graph.PagedGraphUsers{
			GraphUsers: &[]graph.GraphUser{},
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataUsers().Schema, nil)
	err := dataUsersRead(resourceData, clients)
	require.Nil(t, err)
	users, ok := resourceData.GetOk("users")
	require.False(t, ok)
	require.NotNil(t, users)
	usersSet, ok := users.(*schema.Set)
	require.True(t, ok)
	require.NotNil(t, usersSet)
	require.Equal(t, 0, usersSet.Len())
}

func TestDataSourceUser_Read_TestFilterByPricipalName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &client.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	/* start writing test here */
	expectedArgs := graph.ListUsersArgs{
		SubjectTypes: &[]string{},
	}
	graphClient.
		EXPECT().
		ListUsers(clients.Ctx, expectedArgs).
		Return(&graph.PagedGraphUsers{
			GraphUsers: &usrList1,
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataUsers().Schema, nil)
	resourceData.Set("principal_name", "DesireeMCollins@jourrapide.com")
	err := dataUsersRead(resourceData, clients)
	require.Nil(t, err)
	users, ok := resourceData.GetOk("users")
	require.True(t, ok)
	require.NotNil(t, users)
	usersSet, ok := users.(*schema.Set)
	require.True(t, ok)
	require.NotNil(t, usersSet)
	require.Equal(t, 1, usersSet.Len())
	u, _ := flattenUser(&usrList1[0])
	require.True(t, usersSet.Contains(u))
}

func TestDataSourceUser_Read_TestFilterByOrigin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &client.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	/* start writing test here */
	expectedArgs := graph.ListUsersArgs{
		SubjectTypes: &[]string{},
	}
	graphClient.
		EXPECT().
		ListUsers(clients.Ctx, expectedArgs).
		Return(&graph.PagedGraphUsers{
			GraphUsers: &usrList1,
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataUsers().Schema, nil)
	resourceData.Set("origin", "aad")
	err := dataUsersRead(resourceData, clients)
	require.Nil(t, err)
	users, ok := resourceData.GetOk("users")
	require.True(t, ok)
	require.NotNil(t, users)
	usersSet, ok := users.(*schema.Set)
	require.True(t, ok)
	require.NotNil(t, usersSet)
	require.Equal(t, 2, usersSet.Len())

	iFound := 0
	for _, elem := range usersSet.List() {
		upn := elem.(map[string]interface{})["principal_name"].(string)
		for _, usr := range usrList1 {
			if strings.EqualFold("aad", *usr.Origin) && strings.EqualFold(upn, *usr.PrincipalName) {
				iFound++
				break
			}
		}
	}
	require.Equal(t, 2, iFound)
}

func TestDataSourceUser_Read_TestFilterByOriginId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &client.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	/* start writing test here */
	expectedArgs := graph.ListUsersArgs{
		SubjectTypes: &[]string{},
	}
	graphClient.
		EXPECT().
		ListUsers(clients.Ctx, expectedArgs).
		Return(&graph.PagedGraphUsers{
			GraphUsers: &usrList1,
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataUsers().Schema, nil)
	resourceData.Set("origin_id", "8c840d92-f19e-4dfe-8eab-5a1fd67a3a77")
	err := dataUsersRead(resourceData, clients)
	require.Nil(t, err)
	users, ok := resourceData.GetOk("users")
	require.True(t, ok)
	require.NotNil(t, users)
	usersSet, ok := users.(*schema.Set)
	require.True(t, ok)
	require.NotNil(t, usersSet)
	require.Equal(t, 2, usersSet.Len())

	iFound := 0
	for _, elem := range usersSet.List() {
		upn := elem.(map[string]interface{})["principal_name"].(string)
		for _, usr := range usrList1 {
			if strings.EqualFold("8c840d92-f19e-4dfe-8eab-5a1fd67a3a77", *usr.OriginId) && strings.EqualFold(upn, *usr.PrincipalName) {
				iFound++
				break
			}
		}
	}
	require.Equal(t, 2, iFound)
}

func TestDataSourceUser_Read_TestFilterByOriginOriginId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &client.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	/* start writing test here */
	expectedArgs := graph.ListUsersArgs{
		SubjectTypes: &[]string{},
	}
	graphClient.
		EXPECT().
		ListUsers(clients.Ctx, expectedArgs).
		Return(&graph.PagedGraphUsers{
			GraphUsers: &usrList1,
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataUsers().Schema, nil)
	resourceData.Set("origin", "aad")
	resourceData.Set("origin_id", "8c840d92-f19e-4dfe-8eab-5a1fd67a3a77")
	err := dataUsersRead(resourceData, clients)
	require.Nil(t, err)
	users, ok := resourceData.GetOk("users")
	require.True(t, ok)
	require.NotNil(t, users)
	usersSet, ok := users.(*schema.Set)
	require.True(t, ok)
	require.NotNil(t, usersSet)
	require.Equal(t, 1, usersSet.Len())

	iFound := 0
	for _, elem := range usersSet.List() {
		upn := elem.(map[string]interface{})["principal_name"].(string)
		for _, usr := range usrList1 {
			if strings.EqualFold("8c840d92-f19e-4dfe-8eab-5a1fd67a3a77", *usr.OriginId) && strings.EqualFold("aad", *usr.Origin) && strings.EqualFold(upn, *usr.PrincipalName) {
				iFound++
				break
			}
		}
	}
	require.Equal(t, 1, iFound)
}

func TestDataSourceUser_Read_TestFilterBySubjectType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &client.AggregatedClient{
		GraphClient: graphClient,
		Ctx:         context.Background(),
	}

	var usrList []graph.GraphUser

	linq.From(usrList2).
		WhereT(func(x interface{}) bool {
			usr := x.(graph.GraphUser)
			return usr.Origin != nil && strings.EqualFold(*usr.Origin, "aad")
		}).
		ToSlice(&usrList)

	/* start writing test here */
	expectedArgs := graph.ListUsersArgs{
		SubjectTypes: &[]string{"aad"},
	}
	graphClient.
		EXPECT().
		ListUsers(clients.Ctx, expectedArgs).
		Return(&graph.PagedGraphUsers{
			GraphUsers: &usrList,
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataUsers().Schema, nil)
	resourceData.Set("subject_types", schema.NewSet(schema.HashString, []interface{}{"aad"}))
	err := dataUsersRead(resourceData, clients)
	require.Nil(t, err)
	users, ok := resourceData.GetOk("users")
	require.True(t, ok)
	require.NotNil(t, users)
	usersSet, ok := users.(*schema.Set)
	require.True(t, ok)
	require.NotNil(t, usersSet)
	require.Equal(t, len(usrList), usersSet.Len())

	iFound := 0
	for _, elem := range usersSet.List() {
		upn := elem.(map[string]interface{})["principal_name"].(string)
		for _, usr := range usrList {
			if strings.EqualFold(upn, *usr.PrincipalName) {
				iFound++
				break
			}
		}
	}
	require.Equal(t, len(usrList), iFound)
}
