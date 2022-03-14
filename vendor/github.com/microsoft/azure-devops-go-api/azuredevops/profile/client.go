// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package profile

import (
	"context"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"net/http"
	"net/url"
	"strconv"
)

var ResourceAreaId, _ = uuid.Parse("8ccfef3d-2b87-4e99-8ccb-66e343d2daa8")

type Client interface {
	// Gets a user profile.
	GetProfile(context.Context, GetProfileArgs) (*Profile, error)
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

// Gets a user profile.
func (client *ClientImpl) GetProfile(ctx context.Context, args GetProfileArgs) (*Profile, error) {
	routeValues := make(map[string]string)
	if args.Id == nil || *args.Id == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = *args.Id

	queryParams := url.Values{}
	if args.Details != nil {
		queryParams.Add("details", strconv.FormatBool(*args.Details))
	}
	if args.WithAttributes != nil {
		queryParams.Add("withAttributes", strconv.FormatBool(*args.WithAttributes))
	}
	if args.Partition != nil {
		queryParams.Add("partition", *args.Partition)
	}
	if args.CoreAttributes != nil {
		queryParams.Add("coreAttributes", *args.CoreAttributes)
	}
	if args.ForceRefresh != nil {
		queryParams.Add("forceRefresh", strconv.FormatBool(*args.ForceRefresh))
	}
	locationId, _ := uuid.Parse("f83735dc-483f-4238-a291-d45f6080a9af")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Profile
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetProfile function
type GetProfileArgs struct {
	// (required) The ID of the target user profile within the same organization, or 'me' to get the profile of the current authenticated user.
	Id *string
	// (optional) Return public profile information such as display name, email address, country, etc. If false, the withAttributes parameter is ignored.
	Details *bool
	// (optional) If true, gets the attributes (named key-value pairs of arbitrary data) associated with the profile. The partition parameter must also have a value.
	WithAttributes *bool
	// (optional) The partition (named group) of attributes to return.
	Partition *string
	// (optional) A comma-delimited list of core profile attributes to return. Valid values are Email, Avatar, DisplayName, and ContactWithOffers.
	CoreAttributes *string
	// (optional) Not used in this version of the API.
	ForceRefresh *bool
}
