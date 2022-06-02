// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/microsoft/azure-devops-go-api/azuredevops/v6/graph (interfaces: Client)

// Package azdosdkmocks is a generated GoMock package.
package azdosdkmocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	graph "github.com/microsoft/azure-devops-go-api/azuredevops/v6/graph"
	profile "github.com/microsoft/azure-devops-go-api/azuredevops/v6/profile"
)

// MockGraphClient is a mock of Client interface.
type MockGraphClient struct {
	ctrl     *gomock.Controller
	recorder *MockGraphClientMockRecorder
}

// MockGraphClientMockRecorder is the mock recorder for MockGraphClient.
type MockGraphClientMockRecorder struct {
	mock *MockGraphClient
}

// NewMockGraphClient creates a new mock instance.
func NewMockGraphClient(ctrl *gomock.Controller) *MockGraphClient {
	mock := &MockGraphClient{ctrl: ctrl}
	mock.recorder = &MockGraphClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGraphClient) EXPECT() *MockGraphClientMockRecorder {
	return m.recorder
}

// AddMembership mocks base method.
func (m *MockGraphClient) AddMembership(arg0 context.Context, arg1 graph.AddMembershipArgs) (*graph.GraphMembership, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMembership", arg0, arg1)
	ret0, _ := ret[0].(*graph.GraphMembership)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMembership indicates an expected call of AddMembership.
func (mr *MockGraphClientMockRecorder) AddMembership(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMembership", reflect.TypeOf((*MockGraphClient)(nil).AddMembership), arg0, arg1)
}

// CheckMembershipExistence mocks base method.
func (m *MockGraphClient) CheckMembershipExistence(arg0 context.Context, arg1 graph.CheckMembershipExistenceArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckMembershipExistence", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckMembershipExistence indicates an expected call of CheckMembershipExistence.
func (mr *MockGraphClientMockRecorder) CheckMembershipExistence(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckMembershipExistence", reflect.TypeOf((*MockGraphClient)(nil).CheckMembershipExistence), arg0, arg1)
}

// CreateGroup mocks base method.
func (m *MockGraphClient) CreateGroup(arg0 context.Context, arg1 graph.CreateGroupArgs) (*graph.GraphGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGroup", arg0, arg1)
	ret0, _ := ret[0].(*graph.GraphGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGroup indicates an expected call of CreateGroup.
func (mr *MockGraphClientMockRecorder) CreateGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroup", reflect.TypeOf((*MockGraphClient)(nil).CreateGroup), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockGraphClient) CreateUser(arg0 context.Context, arg1 graph.CreateUserArgs) (*graph.GraphUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(*graph.GraphUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockGraphClientMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockGraphClient)(nil).CreateUser), arg0, arg1)
}

// DeleteAvatar mocks base method.
func (m *MockGraphClient) DeleteAvatar(arg0 context.Context, arg1 graph.DeleteAvatarArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAvatar", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAvatar indicates an expected call of DeleteAvatar.
func (mr *MockGraphClientMockRecorder) DeleteAvatar(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAvatar", reflect.TypeOf((*MockGraphClient)(nil).DeleteAvatar), arg0, arg1)
}

// DeleteGroup mocks base method.
func (m *MockGraphClient) DeleteGroup(arg0 context.Context, arg1 graph.DeleteGroupArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGroup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGroup indicates an expected call of DeleteGroup.
func (mr *MockGraphClientMockRecorder) DeleteGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGroup", reflect.TypeOf((*MockGraphClient)(nil).DeleteGroup), arg0, arg1)
}

// DeleteUser mocks base method.
func (m *MockGraphClient) DeleteUser(arg0 context.Context, arg1 graph.DeleteUserArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockGraphClientMockRecorder) DeleteUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockGraphClient)(nil).DeleteUser), arg0, arg1)
}

// GetAvatar mocks base method.
func (m *MockGraphClient) GetAvatar(arg0 context.Context, arg1 graph.GetAvatarArgs) (*profile.Avatar, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvatar", arg0, arg1)
	ret0, _ := ret[0].(*profile.Avatar)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAvatar indicates an expected call of GetAvatar.
func (mr *MockGraphClientMockRecorder) GetAvatar(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvatar", reflect.TypeOf((*MockGraphClient)(nil).GetAvatar), arg0, arg1)
}

// GetDescriptor mocks base method.
func (m *MockGraphClient) GetDescriptor(arg0 context.Context, arg1 graph.GetDescriptorArgs) (*graph.GraphDescriptorResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDescriptor", arg0, arg1)
	ret0, _ := ret[0].(*graph.GraphDescriptorResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDescriptor indicates an expected call of GetDescriptor.
func (mr *MockGraphClientMockRecorder) GetDescriptor(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDescriptor", reflect.TypeOf((*MockGraphClient)(nil).GetDescriptor), arg0, arg1)
}

// GetGroup mocks base method.
func (m *MockGraphClient) GetGroup(arg0 context.Context, arg1 graph.GetGroupArgs) (*graph.GraphGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroup", arg0, arg1)
	ret0, _ := ret[0].(*graph.GraphGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroup indicates an expected call of GetGroup.
func (mr *MockGraphClientMockRecorder) GetGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroup", reflect.TypeOf((*MockGraphClient)(nil).GetGroup), arg0, arg1)
}

// GetMembership mocks base method.
func (m *MockGraphClient) GetMembership(arg0 context.Context, arg1 graph.GetMembershipArgs) (*graph.GraphMembership, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMembership", arg0, arg1)
	ret0, _ := ret[0].(*graph.GraphMembership)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMembership indicates an expected call of GetMembership.
func (mr *MockGraphClientMockRecorder) GetMembership(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMembership", reflect.TypeOf((*MockGraphClient)(nil).GetMembership), arg0, arg1)
}

// GetMembershipState mocks base method.
func (m *MockGraphClient) GetMembershipState(arg0 context.Context, arg1 graph.GetMembershipStateArgs) (*graph.GraphMembershipState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMembershipState", arg0, arg1)
	ret0, _ := ret[0].(*graph.GraphMembershipState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMembershipState indicates an expected call of GetMembershipState.
func (mr *MockGraphClientMockRecorder) GetMembershipState(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMembershipState", reflect.TypeOf((*MockGraphClient)(nil).GetMembershipState), arg0, arg1)
}

// GetProviderInfo mocks base method.
func (m *MockGraphClient) GetProviderInfo(arg0 context.Context, arg1 graph.GetProviderInfoArgs) (*graph.GraphProviderInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProviderInfo", arg0, arg1)
	ret0, _ := ret[0].(*graph.GraphProviderInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProviderInfo indicates an expected call of GetProviderInfo.
func (mr *MockGraphClientMockRecorder) GetProviderInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProviderInfo", reflect.TypeOf((*MockGraphClient)(nil).GetProviderInfo), arg0, arg1)
}

// GetStorageKey mocks base method.
func (m *MockGraphClient) GetStorageKey(arg0 context.Context, arg1 graph.GetStorageKeyArgs) (*graph.GraphStorageKeyResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStorageKey", arg0, arg1)
	ret0, _ := ret[0].(*graph.GraphStorageKeyResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStorageKey indicates an expected call of GetStorageKey.
func (mr *MockGraphClientMockRecorder) GetStorageKey(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStorageKey", reflect.TypeOf((*MockGraphClient)(nil).GetStorageKey), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockGraphClient) GetUser(arg0 context.Context, arg1 graph.GetUserArgs) (*graph.GraphUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(*graph.GraphUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockGraphClientMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockGraphClient)(nil).GetUser), arg0, arg1)
}

// ListGroups mocks base method.
func (m *MockGraphClient) ListGroups(arg0 context.Context, arg1 graph.ListGroupsArgs) (*graph.PagedGraphGroups, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListGroups", arg0, arg1)
	ret0, _ := ret[0].(*graph.PagedGraphGroups)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListGroups indicates an expected call of ListGroups.
func (mr *MockGraphClientMockRecorder) ListGroups(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListGroups", reflect.TypeOf((*MockGraphClient)(nil).ListGroups), arg0, arg1)
}

// ListMemberships mocks base method.
func (m *MockGraphClient) ListMemberships(arg0 context.Context, arg1 graph.ListMembershipsArgs) (*[]graph.GraphMembership, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListMemberships", arg0, arg1)
	ret0, _ := ret[0].(*[]graph.GraphMembership)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListMemberships indicates an expected call of ListMemberships.
func (mr *MockGraphClientMockRecorder) ListMemberships(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListMemberships", reflect.TypeOf((*MockGraphClient)(nil).ListMemberships), arg0, arg1)
}

// ListUsers mocks base method.
func (m *MockGraphClient) ListUsers(arg0 context.Context, arg1 graph.ListUsersArgs) (*graph.PagedGraphUsers, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsers", arg0, arg1)
	ret0, _ := ret[0].(*graph.PagedGraphUsers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsers indicates an expected call of ListUsers.
func (mr *MockGraphClientMockRecorder) ListUsers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockGraphClient)(nil).ListUsers), arg0, arg1)
}

// LookupSubjects mocks base method.
func (m *MockGraphClient) LookupSubjects(arg0 context.Context, arg1 graph.LookupSubjectsArgs) (*map[string]graph.GraphSubject, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupSubjects", arg0, arg1)
	ret0, _ := ret[0].(*map[string]graph.GraphSubject)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LookupSubjects indicates an expected call of LookupSubjects.
func (mr *MockGraphClientMockRecorder) LookupSubjects(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupSubjects", reflect.TypeOf((*MockGraphClient)(nil).LookupSubjects), arg0, arg1)
}

// QuerySubjects mocks base method.
func (m *MockGraphClient) QuerySubjects(arg0 context.Context, arg1 graph.QuerySubjectsArgs) (*[]graph.GraphSubject, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QuerySubjects", arg0, arg1)
	ret0, _ := ret[0].(*[]graph.GraphSubject)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QuerySubjects indicates an expected call of QuerySubjects.
func (mr *MockGraphClientMockRecorder) QuerySubjects(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QuerySubjects", reflect.TypeOf((*MockGraphClient)(nil).QuerySubjects), arg0, arg1)
}

// RemoveMembership mocks base method.
func (m *MockGraphClient) RemoveMembership(arg0 context.Context, arg1 graph.RemoveMembershipArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveMembership", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveMembership indicates an expected call of RemoveMembership.
func (mr *MockGraphClientMockRecorder) RemoveMembership(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveMembership", reflect.TypeOf((*MockGraphClient)(nil).RemoveMembership), arg0, arg1)
}

// RequestAccess mocks base method.
func (m *MockGraphClient) RequestAccess(arg0 context.Context, arg1 graph.RequestAccessArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequestAccess", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RequestAccess indicates an expected call of RequestAccess.
func (mr *MockGraphClientMockRecorder) RequestAccess(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestAccess", reflect.TypeOf((*MockGraphClient)(nil).RequestAccess), arg0, arg1)
}

// SetAvatar mocks base method.
func (m *MockGraphClient) SetAvatar(arg0 context.Context, arg1 graph.SetAvatarArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetAvatar", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetAvatar indicates an expected call of SetAvatar.
func (mr *MockGraphClientMockRecorder) SetAvatar(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAvatar", reflect.TypeOf((*MockGraphClient)(nil).SetAvatar), arg0, arg1)
}

// UpdateGroup mocks base method.
func (m *MockGraphClient) UpdateGroup(arg0 context.Context, arg1 graph.UpdateGroupArgs) (*graph.GraphGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGroup", arg0, arg1)
	ret0, _ := ret[0].(*graph.GraphGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateGroup indicates an expected call of UpdateGroup.
func (mr *MockGraphClientMockRecorder) UpdateGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGroup", reflect.TypeOf((*MockGraphClient)(nil).UpdateGroup), arg0, arg1)
}

// UpdateUser mocks base method.
func (m *MockGraphClient) UpdateUser(arg0 context.Context, arg1 graph.UpdateUserArgs) (*graph.GraphUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1)
	ret0, _ := ret[0].(*graph.GraphUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockGraphClientMockRecorder) UpdateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockGraphClient)(nil).UpdateUser), arg0, arg1)
}
