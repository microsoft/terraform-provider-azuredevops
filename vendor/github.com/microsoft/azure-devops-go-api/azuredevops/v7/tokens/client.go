// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package tokens

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
)

var ResourceAreaId, _ = uuid.Parse("55967393-20ef-45c6-a96c-b5d5d5986a9a")

type Client interface {
	// [Preview API] Creates a new personal access token (PAT) for the requesting user.
	CreatePat(context.Context, CreatePatArgs) (*PatTokenResult, error)
	// [Preview API] Gets a single personal access token (PAT).
	GetPat(context.Context, GetPatArgs) (*PatTokenResult, error)
	// [Preview API] Revokes a personal access token (PAT) by authorizationId.
	Revoke(context.Context, RevokeArgs) error
	// [Preview API] Updates an existing personal access token (PAT) with the new parameters. To update a token, it must be valid (has not been revoked).
	UpdatePat(context.Context, UpdatePatArgs) (*PatTokenResult, error)
	// [Preview API] Gets a paged list of personal access tokens (PATs) created in this organization. Subsequent calls to the API require the same filtering options to be supplied.
	ListPats(context.Context, ListPatsArgs) (*PagedPatResults, error)
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

// [Preview API] Creates a new personal access token (PAT) for the requesting user.
func (client *ClientImpl) CreatePat(ctx context.Context, args CreatePatArgs) (*PatTokenResult, error) {
	if args.Token == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Token"}
	}

	body, marshalErr := json.Marshal(*args.Token)
	if marshalErr != nil {
		return nil, marshalErr
	}

	locationId, _ := uuid.Parse("55967393-20ef-45c6-a96c-b5d5d5986a9a")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue PatTokenResult
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

type CreatePatArgs struct {
	// (required) The request parameters for creating a personal access token (PAT)
	Token *PatTokenCreateRequest
}

// [Preview API] Gets a single personal access token (PAT).
func (client *ClientImpl) GetPat(ctx context.Context, args GetPatArgs) (*PatTokenResult, error) {
	if args.AuthorizationId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.AuthorizationId"}
	}

	queryParams := url.Values{}
	queryParams.Add("authorizationId", (*args.AuthorizationId).String())
	locationId, _ := uuid.Parse("55967393-20ef-45c6-a96c-b5d5d5986a9a")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", nil, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue PatTokenResult
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

type GetPatArgs struct {
	// (required) The authorizationId identifying a single, unique personal access token (PAT).
	AuthorizationId *uuid.UUID
}

// [Preview API] Revokes a personal access token (PAT) by authorizationId.
func (client *ClientImpl) Revoke(ctx context.Context, args RevokeArgs) error {
	if args.AuthorizationId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.AuthorizationId"}
	}

	queryParams := url.Values{}
	queryParams.Add("authorizationId", (*args.AuthorizationId).String())
	locationId, _ := uuid.Parse("55967393-20ef-45c6-a96c-b5d5d5986a9a")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", nil, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

type RevokeArgs struct {
	// (required) The authorizationId identifying a single, unique personal access token (PAT).
	AuthorizationId *uuid.UUID
}

func (client *ClientImpl) UpdatePat(ctx context.Context, args UpdatePatArgs) (*PatTokenResult, error) {
	if args.Token == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Body"}
	}

	body, marshalErr := json.Marshal(*args.Token)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("55967393-20ef-45c6-a96c-b5d5d5986a9a")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue PatTokenResult
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

type UpdatePatArgs struct {
	// (required) The authorizationId identifying a single, unique personal access token (PAT).
	Token *PatTokenUpdateRequest
}

// [Preview API] Lists of all the session token details of the personal access tokens (PATs) for a particular user.
func (client *ClientImpl) ListPats(ctx context.Context, args ListPatsArgs) (*PagedPatResults, error) {
	queryParams := url.Values{}
	if args.DisplayFilterOption != nil && *args.DisplayFilterOption != "" {
		queryParams.Add("displayFilterOption", string(*args.DisplayFilterOption))
	}
	if args.SortByOption != nil && *args.SortByOption != "" {
		queryParams.Add("sortByOption", string(*args.SortByOption))
	}
	if args.IsSortAscending != nil {
		queryParams.Add("isSortAscending", strconv.FormatBool(*args.IsSortAscending))
	}
	if args.ContinuationToken != nil && *args.ContinuationToken != "" {
		queryParams.Add("continuationToken", *args.ContinuationToken)
	}
	if args.Top != nil {
		queryParams.Add("$top", strconv.Itoa(*args.Top))
	}

	locationId, _ := uuid.Parse("55967393-20ef-45c6-a96c-b5d5d5986a9a")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", nil, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue PagedPatResults
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ListPersonalAccessTokens function
type ListPatsArgs struct {
	// (Optional) Refers to the status of the personal access token (PAT)
	DisplayFilterOption *DisplayFilterOption
	// (Optional) Which field to sort by
	SortByOption *SortByOption
	// (Optional) Ascending or descending
	IsSortAscending *bool
	// (Optional) Where to start returning tokens from
	ContinuationToken *string
	// (Optional) How many tokens to return, limit of 100
	Top *int
}
