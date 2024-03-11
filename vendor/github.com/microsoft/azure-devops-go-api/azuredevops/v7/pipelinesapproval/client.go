// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package pipelinesapproval

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"net/http"
	"net/url"
	"strings"
)

var ResourceAreaId, _ = uuid.Parse("5b55a9b6-2e0f-40d7-829d-3741d2b8c4e4")

type Client interface {
	// [Preview API] Get an approval.
	GetApproval(context.Context, GetApprovalArgs) (*Approval, error)
	// [Preview API] List Approvals. This can be used to get a set of pending approvals in a pipeline, on an user or for a resource..
	QueryApprovals(context.Context, QueryApprovalsArgs) (*[]Approval, error)
	// [Preview API] Update approvals.
	UpdateApprovals(context.Context, UpdateApprovalsArgs) (*[]Approval, error)
}

type ClientImpl struct {
	Client azuredevops.Client
}

func NewClient(ctx context.Context, connection *azuredevops.Connection) (Client, error) {
	client, err := connection.GetClientByResourceAreaId(ctx, ResourceAreaId)
	if err != nil {
		return nil, err
	}
	return &ClientImpl{
		Client: *client,
	}, nil
}

// [Preview API] Get an approval.
func (client *ClientImpl) GetApproval(ctx context.Context, args GetApprovalArgs) (*Approval, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.ApprovalId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ApprovalId"}
	}
	routeValues["approvalId"] = (*args.ApprovalId).String()

	queryParams := url.Values{}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	locationId, _ := uuid.Parse("37794717-f36f-4d78-b2bf-4dc30d0cfbcd")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Approval
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetApproval function
type GetApprovalArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Id of the approval.
	ApprovalId *uuid.UUID
	// (optional)
	Expand *ApprovalDetailsExpandParameter
}

// [Preview API] List Approvals. This can be used to get a set of pending approvals in a pipeline, on an user or for a resource..
func (client *ClientImpl) QueryApprovals(ctx context.Context, args QueryApprovalsArgs) (*[]Approval, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	queryParams := url.Values{}
	if args.ApprovalIds != nil {
		var stringList []string
		for _, item := range *args.ApprovalIds {
			stringList = append(stringList, item.String())
		}
		listAsString := strings.Join((stringList)[:], ",")
		queryParams.Add("approvalIds", listAsString)
	}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	locationId, _ := uuid.Parse("37794717-f36f-4d78-b2bf-4dc30d0cfbcd")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []Approval
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryApprovals function
type QueryApprovalsArgs struct {
	// (required) Project ID or project name
	Project *string
	// (optional)
	ApprovalIds *[]uuid.UUID
	// (optional)
	Expand *ApprovalDetailsExpandParameter
}

// [Preview API] Update approvals.
func (client *ClientImpl) UpdateApprovals(ctx context.Context, args UpdateApprovalsArgs) (*[]Approval, error) {
	if args.UpdateParameters == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UpdateParameters"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	body, marshalErr := json.Marshal(*args.UpdateParameters)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("37794717-f36f-4d78-b2bf-4dc30d0cfbcd")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []Approval
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateApprovals function
type UpdateApprovalsArgs struct {
	// (required)
	UpdateParameters *[]ApprovalUpdateParameters
	// (required) Project ID or project name
	Project *string
}
