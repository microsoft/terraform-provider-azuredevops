// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package extensionmanagement

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var ResourceAreaId, _ = uuid.Parse("6c2b0933-3600-42ae-bf8b-93d4f7e83594")

type Client interface {
	// [Preview API] Get an installed extension by its publisher and extension name.
	GetInstalledExtensionByName(context.Context, GetInstalledExtensionByNameArgs) (*InstalledExtension, error)
	// [Preview API] List the installed extensions in the account / project collection.
	GetInstalledExtensions(context.Context, GetInstalledExtensionsArgs) (*[]InstalledExtension, error)
	// [Preview API] Install the specified extension into the account / project collection.
	InstallExtensionByName(context.Context, InstallExtensionByNameArgs) (*InstalledExtension, error)
	// [Preview API] Uninstall the specified extension from the account / project collection.
	UninstallExtensionByName(context.Context, UninstallExtensionByNameArgs) error
	// [Preview API] Update an installed extension. Typically this API is used to enable or disable an extension.
	UpdateInstalledExtension(context.Context, UpdateInstalledExtensionArgs) (*InstalledExtension, error)
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

// [Preview API] Get an installed extension by its publisher and extension name.
func (client *ClientImpl) GetInstalledExtensionByName(ctx context.Context, args GetInstalledExtensionByNameArgs) (*InstalledExtension, error) {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName

	queryParams := url.Values{}
	if args.AssetTypes != nil {
		listAsString := strings.Join((*args.AssetTypes)[:], ":")
		queryParams.Add("assetTypes", listAsString)
	}
	locationId, _ := uuid.Parse("fb0da285-f23e-4b56-8b53-3ef5f9f6de66")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue InstalledExtension
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetInstalledExtensionByName function
type GetInstalledExtensionByNameArgs struct {
	// (required) Name of the publisher. Example: "fabrikam".
	PublisherName *string
	// (required) Name of the extension. Example: "ops-tools".
	ExtensionName *string
	// (optional) Determines which files are returned in the files array.  Provide the wildcard '*' to return all files, or a colon separated list to retrieve files with specific asset types.
	AssetTypes *[]string
}

// [Preview API] List the installed extensions in the account / project collection.
func (client *ClientImpl) GetInstalledExtensions(ctx context.Context, args GetInstalledExtensionsArgs) (*[]InstalledExtension, error) {
	queryParams := url.Values{}
	if args.IncludeDisabledExtensions != nil {
		queryParams.Add("includeDisabledExtensions", strconv.FormatBool(*args.IncludeDisabledExtensions))
	}
	if args.IncludeErrors != nil {
		queryParams.Add("includeErrors", strconv.FormatBool(*args.IncludeErrors))
	}
	if args.AssetTypes != nil {
		listAsString := strings.Join((*args.AssetTypes)[:], ":")
		queryParams.Add("assetTypes", listAsString)
	}
	if args.IncludeInstallationIssues != nil {
		queryParams.Add("includeInstallationIssues", strconv.FormatBool(*args.IncludeInstallationIssues))
	}
	locationId, _ := uuid.Parse("275424d0-c844-4fe2-bda6-04933a1357d8")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", nil, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []InstalledExtension
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetInstalledExtensions function
type GetInstalledExtensionsArgs struct {
	// (optional) If true (the default), include disabled extensions in the results.
	IncludeDisabledExtensions *bool
	// (optional) If true, include installed extensions with errors.
	IncludeErrors *bool
	// (optional) Determines which files are returned in the files array.  Provide the wildcard '*' to return all files, or a colon separated list to retrieve files with specific asset types.
	AssetTypes *[]string
	// (optional)
	IncludeInstallationIssues *bool
}

// [Preview API] Install the specified extension into the account / project collection.
func (client *ClientImpl) InstallExtensionByName(ctx context.Context, args InstallExtensionByNameArgs) (*InstalledExtension, error) {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.Version != nil && *args.Version != "" {
		routeValues["version"] = *args.Version
	}

	locationId, _ := uuid.Parse("fb0da285-f23e-4b56-8b53-3ef5f9f6de66")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue InstalledExtension
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the InstallExtensionByName function
type InstallExtensionByNameArgs struct {
	// (required) Name of the publisher. Example: "fabrikam".
	PublisherName *string
	// (required) Name of the extension. Example: "ops-tools".
	ExtensionName *string
	// (optional)
	Version *string
}

// [Preview API] Uninstall the specified extension from the account / project collection.
func (client *ClientImpl) UninstallExtensionByName(ctx context.Context, args UninstallExtensionByNameArgs) error {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName

	queryParams := url.Values{}
	if args.Reason != nil {
		queryParams.Add("reason", *args.Reason)
	}
	if args.ReasonCode != nil {
		queryParams.Add("reasonCode", *args.ReasonCode)
	}
	locationId, _ := uuid.Parse("fb0da285-f23e-4b56-8b53-3ef5f9f6de66")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the UninstallExtensionByName function
type UninstallExtensionByNameArgs struct {
	// (required) Name of the publisher. Example: "fabrikam".
	PublisherName *string
	// (required) Name of the extension. Example: "ops-tools".
	ExtensionName *string
	// (optional)
	Reason *string
	// (optional)
	ReasonCode *string
}

// [Preview API] Update an installed extension. Typically this API is used to enable or disable an extension.
func (client *ClientImpl) UpdateInstalledExtension(ctx context.Context, args UpdateInstalledExtensionArgs) (*InstalledExtension, error) {
	if args.Extension == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Extension"}
	}
	body, marshalErr := json.Marshal(*args.Extension)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("275424d0-c844-4fe2-bda6-04933a1357d8")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue InstalledExtension
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateInstalledExtension function
type UpdateInstalledExtensionArgs struct {
	// (required)
	Extension *InstalledExtension
}
