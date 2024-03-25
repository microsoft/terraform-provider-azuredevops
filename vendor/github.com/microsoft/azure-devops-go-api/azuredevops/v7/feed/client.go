// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package feed

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/operations"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

var ResourceAreaId, _ = uuid.Parse("7ab4e64e-c4d8-4f50-ae73-5ef2e21642a5")

type Client interface {
	// [Preview API] Create a feed, a container for various package types.
	CreateFeed(context.Context, CreateFeedArgs) (*Feed, error)
	// [Preview API] Create a new view on the referenced feed.
	CreateFeedView(context.Context, CreateFeedViewArgs) (*FeedView, error)
	// [Preview API] Remove a feed and all its packages. The feed moves to the recycle bin and is reversible.
	DeleteFeed(context.Context, DeleteFeedArgs) error
	// [Preview API] Delete the retention policy for a feed.
	DeleteFeedRetentionPolicies(context.Context, DeleteFeedRetentionPoliciesArgs) error
	// [Preview API] Delete a feed view.
	DeleteFeedView(context.Context, DeleteFeedViewArgs) error
	// [Preview API] Queues a job to remove all package versions from a feed's recycle bin
	EmptyRecycleBin(context.Context, EmptyRecycleBinArgs) (*operations.OperationReference, error)
	// [Preview API] Generate a SVG badge for the latest version of a package.  The generated SVG is typically used as the image in an HTML link which takes users to the feed containing the package to accelerate discovery and consumption.
	GetBadge(context.Context, GetBadgeArgs) (io.ReadCloser, error)
	// [Preview API] Get the settings for a specific feed.
	GetFeed(context.Context, GetFeedArgs) (*Feed, error)
	// [Preview API] Query a feed to determine its current state.
	GetFeedChange(context.Context, GetFeedChangeArgs) (*FeedChange, error)
	// [Preview API] Query to determine which feeds have changed since the last call, tracked through the provided continuationToken. Only changes to a feed itself are returned and impact the continuationToken, not additions or alterations to packages within the feeds.
	GetFeedChanges(context.Context, GetFeedChangesArgs) (*FeedChangesResponse, error)
	// [Preview API] Get the permissions for a feed.
	GetFeedPermissions(context.Context, GetFeedPermissionsArgs) (*[]FeedPermission, error)
	// [Preview API] Get the retention policy for a feed.
	GetFeedRetentionPolicies(context.Context, GetFeedRetentionPoliciesArgs) (*FeedRetentionPolicy, error)
	// [Preview API] Get all feeds in an account where you have the provided role access.
	GetFeeds(context.Context, GetFeedsArgs) (*[]Feed, error)
	// [Preview API] Query for feeds within the recycle bin.
	GetFeedsFromRecycleBin(context.Context, GetFeedsFromRecycleBinArgs) (*[]Feed, error)
	// [Preview API] Get a view by Id.
	GetFeedView(context.Context, GetFeedViewArgs) (*FeedView, error)
	// [Preview API] Get all views for a feed.
	GetFeedViews(context.Context, GetFeedViewsArgs) (*[]FeedView, error)
	// [Preview API] Get all service-wide feed creation and administration permissions.
	GetGlobalPermissions(context.Context, GetGlobalPermissionsArgs) (*[]GlobalPermission, error)
	// [Preview API] Get details about a specific package.
	GetPackage(context.Context, GetPackageArgs) (*Package, error)
	// [Preview API] Get a batch of package changes made to a feed.  The changes returned are 'most recent change' so if an Add is followed by an Update before you begin enumerating, you'll only see one change in the batch.  While consuming batches using the continuation token, you may see changes to the same package version multiple times if they are happening as you enumerate.
	GetPackageChanges(context.Context, GetPackageChangesArgs) (*PackageChangesResponse, error)
	// [Preview API] Get details about all of the packages in the feed. Use the various filters to include or exclude information from the result set.
	GetPackages(context.Context, GetPackagesArgs) (*[]Package, error)
	// [Preview API] Get details about a specific package version.
	GetPackageVersion(context.Context, GetPackageVersionArgs) (*PackageVersion, error)
	// [Preview API] Gets provenance for a package version.
	GetPackageVersionProvenance(context.Context, GetPackageVersionProvenanceArgs) (*PackageVersionProvenance, error)
	// [Preview API] Get a list of package versions, optionally filtering by state.
	GetPackageVersions(context.Context, GetPackageVersionsArgs) (*[]PackageVersion, error)
	// [Preview API] Get information about a package and all its versions within the recycle bin.
	GetRecycleBinPackage(context.Context, GetRecycleBinPackageArgs) (*Package, error)
	// [Preview API] Query for packages within the recycle bin.
	GetRecycleBinPackages(context.Context, GetRecycleBinPackagesArgs) (*[]Package, error)
	// [Preview API] Get information about a package version within the recycle bin.
	GetRecycleBinPackageVersion(context.Context, GetRecycleBinPackageVersionArgs) (*RecycleBinPackageVersion, error)
	// [Preview API] Get a list of package versions within the recycle bin.
	GetRecycleBinPackageVersions(context.Context, GetRecycleBinPackageVersionsArgs) (*[]RecycleBinPackageVersion, error)
	// [Preview API]
	PermanentDeleteFeed(context.Context, PermanentDeleteFeedArgs) error
	// [Preview API]
	QueryPackageMetrics(context.Context, QueryPackageMetricsArgs) (*[]PackageMetrics, error)
	// [Preview API]
	QueryPackageVersionMetrics(context.Context, QueryPackageVersionMetricsArgs) (*[]PackageVersionMetrics, error)
	// [Preview API]
	RestoreDeletedFeed(context.Context, RestoreDeletedFeedArgs) error
	// [Preview API] Update the permissions on a feed.
	SetFeedPermissions(context.Context, SetFeedPermissionsArgs) (*[]FeedPermission, error)
	// [Preview API] Set the retention policy for a feed.
	SetFeedRetentionPolicies(context.Context, SetFeedRetentionPoliciesArgs) (*FeedRetentionPolicy, error)
	// [Preview API] Set service-wide permissions that govern feed creation and administration.
	SetGlobalPermissions(context.Context, SetGlobalPermissionsArgs) (*[]GlobalPermission, error)
	// [Preview API] Change the attributes of a feed.
	UpdateFeed(context.Context, UpdateFeedArgs) (*Feed, error)
	// [Preview API] Update a view.
	UpdateFeedView(context.Context, UpdateFeedViewArgs) (*FeedView, error)
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

// [Preview API] Create a feed, a container for various package types.
func (client *ClientImpl) CreateFeed(ctx context.Context, args CreateFeedArgs) (*Feed, error) {
	if args.Feed == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Feed"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	body, marshalErr := json.Marshal(*args.Feed)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("c65009a7-474a-4ad1-8b42-7d852107ef8c")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Feed
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateFeed function
type CreateFeedArgs struct {
	// (required) A JSON object containing both required and optional attributes for the feed. Name is the only required value.
	Feed *Feed
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Create a new view on the referenced feed.
func (client *ClientImpl) CreateFeedView(ctx context.Context, args CreateFeedViewArgs) (*FeedView, error) {
	if args.View == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.View"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId

	body, marshalErr := json.Marshal(*args.View)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("42a8502a-6785-41bc-8c16-89477d930877")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue FeedView
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateFeedView function
type CreateFeedViewArgs struct {
	// (required) View to be created.
	View *FeedView
	// (required) Name or Id of the feed.
	FeedId *string
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Remove a feed and all its packages. The feed moves to the recycle bin and is reversible.
func (client *ClientImpl) DeleteFeed(ctx context.Context, args DeleteFeedArgs) error {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId

	locationId, _ := uuid.Parse("c65009a7-474a-4ad1-8b42-7d852107ef8c")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteFeed function
type DeleteFeedArgs struct {
	// (required) Name or Id of the feed.
	FeedId *string
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Delete the retention policy for a feed.
func (client *ClientImpl) DeleteFeedRetentionPolicies(ctx context.Context, args DeleteFeedRetentionPoliciesArgs) error {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId

	locationId, _ := uuid.Parse("ed52a011-0112-45b5-9f9e-e14efffb3193")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteFeedRetentionPolicies function
type DeleteFeedRetentionPoliciesArgs struct {
	// (required) Name or ID of the feed.
	FeedId *string
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Delete a feed view.
func (client *ClientImpl) DeleteFeedView(ctx context.Context, args DeleteFeedViewArgs) error {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId
	if args.ViewId == nil || *args.ViewId == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ViewId"}
	}
	routeValues["viewId"] = *args.ViewId

	locationId, _ := uuid.Parse("42a8502a-6785-41bc-8c16-89477d930877")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteFeedView function
type DeleteFeedViewArgs struct {
	// (required) Name or Id of the feed.
	FeedId *string
	// (required) Name or Id of the view.
	ViewId *string
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Queues a job to remove all package versions from a feed's recycle bin
func (client *ClientImpl) EmptyRecycleBin(ctx context.Context, args EmptyRecycleBinArgs) (*operations.OperationReference, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId

	locationId, _ := uuid.Parse("2704e72c-f541-4141-99be-2004b50b05fa")
	resp, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue operations.OperationReference
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the EmptyRecycleBin function
type EmptyRecycleBinArgs struct {
	// (required) Name or Id of the feed
	FeedId *string
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Generate a SVG badge for the latest version of a package.  The generated SVG is typically used as the image in an HTML link which takes users to the feed containing the package to accelerate discovery and consumption.
func (client *ClientImpl) GetBadge(ctx context.Context, args GetBadgeArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId
	if args.PackageId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PackageId"}
	}
	routeValues["packageId"] = (*args.PackageId).String()

	locationId, _ := uuid.Parse("61d885fd-10f3-4a55-82b6-476d866b673f")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "image/svg+xml", nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetBadge function
type GetBadgeArgs struct {
	// (required) Name or Id of the feed.
	FeedId *string
	// (required) Id of the package (GUID Id, not name).
	PackageId *uuid.UUID
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Get the settings for a specific feed.
func (client *ClientImpl) GetFeed(ctx context.Context, args GetFeedArgs) (*Feed, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId

	queryParams := url.Values{}
	if args.IncludeDeletedUpstreams != nil {
		queryParams.Add("includeDeletedUpstreams", strconv.FormatBool(*args.IncludeDeletedUpstreams))
	}
	locationId, _ := uuid.Parse("c65009a7-474a-4ad1-8b42-7d852107ef8c")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Feed
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetFeed function
type GetFeedArgs struct {
	// (required) Name or Id of the feed.
	FeedId *string
	// (optional) Project ID or project name
	Project *string
	// (optional) Include upstreams that have been deleted in the response.
	IncludeDeletedUpstreams *bool
}

// [Preview API] Query a feed to determine its current state.
func (client *ClientImpl) GetFeedChange(ctx context.Context, args GetFeedChangeArgs) (*FeedChange, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId

	locationId, _ := uuid.Parse("29ba2dad-389a-4661-b5d3-de76397ca05b")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue FeedChange
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetFeedChange function
type GetFeedChangeArgs struct {
	// (required) Name or ID of the feed.
	FeedId *string
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Query to determine which feeds have changed since the last call, tracked through the provided continuationToken. Only changes to a feed itself are returned and impact the continuationToken, not additions or alterations to packages within the feeds.
func (client *ClientImpl) GetFeedChanges(ctx context.Context, args GetFeedChangesArgs) (*FeedChangesResponse, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	queryParams := url.Values{}
	if args.IncludeDeleted != nil {
		queryParams.Add("includeDeleted", strconv.FormatBool(*args.IncludeDeleted))
	}
	if args.ContinuationToken != nil {
		queryParams.Add("continuationToken", strconv.FormatUint(*args.ContinuationToken, 10))
	}
	if args.BatchSize != nil {
		queryParams.Add("batchSize", strconv.Itoa(*args.BatchSize))
	}
	locationId, _ := uuid.Parse("29ba2dad-389a-4661-b5d3-de76397ca05b")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue FeedChangesResponse
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetFeedChanges function
type GetFeedChangesArgs struct {
	// (optional) Project ID or project name
	Project *string
	// (optional) If true, get changes for all feeds including deleted feeds. The default value is false.
	IncludeDeleted *bool
	// (optional) A continuation token which acts as a bookmark to a previously retrieved change. This token allows the user to continue retrieving changes in batches, picking up where the previous batch left off. If specified, all the changes that occur strictly after the token will be returned. If not specified or 0, iteration will start with the first change.
	ContinuationToken *uint64
	// (optional) Number of package changes to fetch. The default value is 1000. The maximum value is 2000.
	BatchSize *int
}

// [Preview API] Get the permissions for a feed.
func (client *ClientImpl) GetFeedPermissions(ctx context.Context, args GetFeedPermissionsArgs) (*[]FeedPermission, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId

	queryParams := url.Values{}
	if args.IncludeIds != nil {
		queryParams.Add("includeIds", strconv.FormatBool(*args.IncludeIds))
	}
	if args.ExcludeInheritedPermissions != nil {
		queryParams.Add("excludeInheritedPermissions", strconv.FormatBool(*args.ExcludeInheritedPermissions))
	}
	if args.IdentityDescriptor != nil {
		queryParams.Add("identityDescriptor", *args.IdentityDescriptor)
	}
	if args.IncludeDeletedFeeds != nil {
		queryParams.Add("includeDeletedFeeds", strconv.FormatBool(*args.IncludeDeletedFeeds))
	}
	locationId, _ := uuid.Parse("be8c1476-86a7-44ed-b19d-aec0e9275cd8")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []FeedPermission
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetFeedPermissions function
type GetFeedPermissionsArgs struct {
	// (required) Name or Id of the feed.
	FeedId *string
	// (optional) Project ID or project name
	Project *string
	// (optional) True to include user Ids in the response.  Default is false.
	IncludeIds *bool
	// (optional) True to only return explicitly set permissions on the feed.  Default is false.
	ExcludeInheritedPermissions *bool
	// (optional) Filter permissions to the provided identity.
	IdentityDescriptor *string
	// (optional) If includeDeletedFeeds is true, then feedId must be specified by name and not by Guid.
	IncludeDeletedFeeds *bool
}

// [Preview API] Get the retention policy for a feed.
func (client *ClientImpl) GetFeedRetentionPolicies(ctx context.Context, args GetFeedRetentionPoliciesArgs) (*FeedRetentionPolicy, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId

	locationId, _ := uuid.Parse("ed52a011-0112-45b5-9f9e-e14efffb3193")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue FeedRetentionPolicy
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetFeedRetentionPolicies function
type GetFeedRetentionPoliciesArgs struct {
	// (required) Name or ID of the feed.
	FeedId *string
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Get all feeds in an account where you have the provided role access.
func (client *ClientImpl) GetFeeds(ctx context.Context, args GetFeedsArgs) (*[]Feed, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	queryParams := url.Values{}
	if args.FeedRole != nil {
		queryParams.Add("feedRole", string(*args.FeedRole))
	}
	if args.IncludeDeletedUpstreams != nil {
		queryParams.Add("includeDeletedUpstreams", strconv.FormatBool(*args.IncludeDeletedUpstreams))
	}
	if args.IncludeUrls != nil {
		queryParams.Add("includeUrls", strconv.FormatBool(*args.IncludeUrls))
	}
	locationId, _ := uuid.Parse("c65009a7-474a-4ad1-8b42-7d852107ef8c")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []Feed
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetFeeds function
type GetFeedsArgs struct {
	// (optional) Project ID or project name
	Project *string
	// (optional) Filter by this role, either Administrator(4), Contributor(3), or Reader(2) level permissions.
	FeedRole *FeedRole
	// (optional) Include upstreams that have been deleted in the response.
	IncludeDeletedUpstreams *bool
	// (optional) Resolve names if true
	IncludeUrls *bool
}

// [Preview API] Query for feeds within the recycle bin.
func (client *ClientImpl) GetFeedsFromRecycleBin(ctx context.Context, args GetFeedsFromRecycleBinArgs) (*[]Feed, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	locationId, _ := uuid.Parse("0cee643d-beb9-41f8-9368-3ada763a8344")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []Feed
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetFeedsFromRecycleBin function
type GetFeedsFromRecycleBinArgs struct {
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Get a view by Id.
func (client *ClientImpl) GetFeedView(ctx context.Context, args GetFeedViewArgs) (*FeedView, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId
	if args.ViewId == nil || *args.ViewId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ViewId"}
	}
	routeValues["viewId"] = *args.ViewId

	locationId, _ := uuid.Parse("42a8502a-6785-41bc-8c16-89477d930877")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue FeedView
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetFeedView function
type GetFeedViewArgs struct {
	// (required) Name or Id of the feed.
	FeedId *string
	// (required) Name or Id of the view.
	ViewId *string
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Get all views for a feed.
func (client *ClientImpl) GetFeedViews(ctx context.Context, args GetFeedViewsArgs) (*[]FeedView, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId

	locationId, _ := uuid.Parse("42a8502a-6785-41bc-8c16-89477d930877")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []FeedView
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetFeedViews function
type GetFeedViewsArgs struct {
	// (required) Name or Id of the feed.
	FeedId *string
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Get all service-wide feed creation and administration permissions.
func (client *ClientImpl) GetGlobalPermissions(ctx context.Context, args GetGlobalPermissionsArgs) (*[]GlobalPermission, error) {
	queryParams := url.Values{}
	if args.IncludeIds != nil {
		queryParams.Add("includeIds", strconv.FormatBool(*args.IncludeIds))
	}
	locationId, _ := uuid.Parse("a74419ef-b477-43df-8758-3cd1cd5f56c6")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", nil, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []GlobalPermission
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetGlobalPermissions function
type GetGlobalPermissionsArgs struct {
	// (optional) Set to true to add IdentityIds to the permission objects.
	IncludeIds *bool
}

// [Preview API] Get details about a specific package.
func (client *ClientImpl) GetPackage(ctx context.Context, args GetPackageArgs) (*Package, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId
	if args.PackageId == nil || *args.PackageId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PackageId"}
	}
	routeValues["packageId"] = *args.PackageId

	queryParams := url.Values{}
	if args.IncludeAllVersions != nil {
		queryParams.Add("includeAllVersions", strconv.FormatBool(*args.IncludeAllVersions))
	}
	if args.IncludeUrls != nil {
		queryParams.Add("includeUrls", strconv.FormatBool(*args.IncludeUrls))
	}
	if args.IsListed != nil {
		queryParams.Add("isListed", strconv.FormatBool(*args.IsListed))
	}
	if args.IsRelease != nil {
		queryParams.Add("isRelease", strconv.FormatBool(*args.IsRelease))
	}
	if args.IncludeDeleted != nil {
		queryParams.Add("includeDeleted", strconv.FormatBool(*args.IncludeDeleted))
	}
	if args.IncludeDescription != nil {
		queryParams.Add("includeDescription", strconv.FormatBool(*args.IncludeDescription))
	}
	locationId, _ := uuid.Parse("7a20d846-c929-4acc-9ea2-0d5a7df1b197")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Package
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetPackage function
type GetPackageArgs struct {
	// (required) Name or Id of the feed.
	FeedId *string
	// (required) The package Id (GUID Id, not the package name).
	PackageId *string
	// (optional) Project ID or project name
	Project *string
	// (optional) True to return all versions of the package in the response. Default is false (latest version only).
	IncludeAllVersions *bool
	// (optional) True to return REST Urls with the response. Default is True.
	IncludeUrls *bool
	// (optional) Only applicable for NuGet packages, setting it for other package types will result in a 404. If false, delisted package versions will be returned. Use this to filter the response when includeAllVersions is set to true. Default is unset (do not return delisted packages).
	IsListed *bool
	// (optional) Only applicable for Nuget packages. Use this to filter the response when includeAllVersions is set to true.  Default is True (only return packages without prerelease versioning).
	IsRelease *bool
	// (optional) Return deleted or unpublished versions of packages in the response. Default is False.
	IncludeDeleted *bool
	// (optional) Return the description for every version of each package in the response. Default is False.
	IncludeDescription *bool
}

// [Preview API] Get a batch of package changes made to a feed.  The changes returned are 'most recent change' so if an Add is followed by an Update before you begin enumerating, you'll only see one change in the batch.  While consuming batches using the continuation token, you may see changes to the same package version multiple times if they are happening as you enumerate.
func (client *ClientImpl) GetPackageChanges(ctx context.Context, args GetPackageChangesArgs) (*PackageChangesResponse, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId

	queryParams := url.Values{}
	if args.ContinuationToken != nil {
		queryParams.Add("continuationToken", strconv.FormatUint(*args.ContinuationToken, 10))
	}
	if args.BatchSize != nil {
		queryParams.Add("batchSize", strconv.Itoa(*args.BatchSize))
	}
	locationId, _ := uuid.Parse("323a0631-d083-4005-85ae-035114dfb681")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue PackageChangesResponse
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetPackageChanges function
type GetPackageChangesArgs struct {
	// (required) Name or Id of the feed.
	FeedId *string
	// (optional) Project ID or project name
	Project *string
	// (optional) A continuation token which acts as a bookmark to a previously retrieved change. This token allows the user to continue retrieving changes in batches, picking up where the previous batch left off. If specified, all the changes that occur strictly after the token will be returned. If not specified or 0, iteration will start with the first change.
	ContinuationToken *uint64
	// (optional) Number of package changes to fetch. The default value is 1000. The maximum value is 2000.
	BatchSize *int
}

// [Preview API] Get details about all of the packages in the feed. Use the various filters to include or exclude information from the result set.
func (client *ClientImpl) GetPackages(ctx context.Context, args GetPackagesArgs) (*[]Package, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId

	queryParams := url.Values{}
	if args.ProtocolType != nil {
		queryParams.Add("protocolType", *args.ProtocolType)
	}
	if args.PackageNameQuery != nil {
		queryParams.Add("packageNameQuery", *args.PackageNameQuery)
	}
	if args.NormalizedPackageName != nil {
		queryParams.Add("normalizedPackageName", *args.NormalizedPackageName)
	}
	if args.IncludeUrls != nil {
		queryParams.Add("includeUrls", strconv.FormatBool(*args.IncludeUrls))
	}
	if args.IncludeAllVersions != nil {
		queryParams.Add("includeAllVersions", strconv.FormatBool(*args.IncludeAllVersions))
	}
	if args.IsListed != nil {
		queryParams.Add("isListed", strconv.FormatBool(*args.IsListed))
	}
	if args.GetTopPackageVersions != nil {
		queryParams.Add("getTopPackageVersions", strconv.FormatBool(*args.GetTopPackageVersions))
	}
	if args.IsRelease != nil {
		queryParams.Add("isRelease", strconv.FormatBool(*args.IsRelease))
	}
	if args.IncludeDescription != nil {
		queryParams.Add("includeDescription", strconv.FormatBool(*args.IncludeDescription))
	}
	if args.Top != nil {
		queryParams.Add("$top", strconv.Itoa(*args.Top))
	}
	if args.Skip != nil {
		queryParams.Add("$skip", strconv.Itoa(*args.Skip))
	}
	if args.IncludeDeleted != nil {
		queryParams.Add("includeDeleted", strconv.FormatBool(*args.IncludeDeleted))
	}
	if args.IsCached != nil {
		queryParams.Add("isCached", strconv.FormatBool(*args.IsCached))
	}
	if args.DirectUpstreamId != nil {
		queryParams.Add("directUpstreamId", (*args.DirectUpstreamId).String())
	}
	locationId, _ := uuid.Parse("7a20d846-c929-4acc-9ea2-0d5a7df1b197")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []Package
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetPackages function
type GetPackagesArgs struct {
	// (required) Name or Id of the feed.
	FeedId *string
	// (optional) Project ID or project name
	Project *string
	// (optional) One of the supported artifact package types.
	ProtocolType *string
	// (optional) Filter to packages that contain the provided string. Characters in the string must conform to the package name constraints.
	PackageNameQuery *string
	// (optional) [Obsolete] Used for legacy scenarios and may be removed in future versions.
	NormalizedPackageName *string
	// (optional) True to return REST Urls with the response. Default is True.
	IncludeUrls *bool
	// (optional) True to return all versions of the package in the response. Default is false (latest version only).
	IncludeAllVersions *bool
	// (optional) Only applicable for NuGet packages, setting it for other package types will result in a 404. If false, delisted package versions will be returned. Use this to filter the response when includeAllVersions is set to true. Default is unset (do not return delisted packages).
	IsListed *bool
	// (optional) Changes the behavior of $top and $skip to return all versions of each package up to $top. Must be used in conjunction with includeAllVersions=true
	GetTopPackageVersions *bool
	// (optional) Only applicable for Nuget packages. Use this to filter the response when includeAllVersions is set to true. Default is True (only return packages without prerelease versioning).
	IsRelease *bool
	// (optional) Return the description for every version of each package in the response. Default is False.
	IncludeDescription *bool
	// (optional) Get the top N packages (or package versions where getTopPackageVersions=true)
	Top *int
	// (optional) Skip the first N packages (or package versions where getTopPackageVersions=true)
	Skip *int
	// (optional) Return deleted or unpublished versions of packages in the response. Default is False.
	IncludeDeleted *bool
	// (optional) [Obsolete] Used for legacy scenarios and may be removed in future versions.
	IsCached *bool
	// (optional) Filter results to return packages from a specific upstream.
	DirectUpstreamId *uuid.UUID
}

// [Preview API] Get details about a specific package version.
func (client *ClientImpl) GetPackageVersion(ctx context.Context, args GetPackageVersionArgs) (*PackageVersion, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId
	if args.PackageId == nil || *args.PackageId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PackageId"}
	}
	routeValues["packageId"] = *args.PackageId
	if args.PackageVersionId == nil || *args.PackageVersionId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PackageVersionId"}
	}
	routeValues["packageVersionId"] = *args.PackageVersionId

	queryParams := url.Values{}
	if args.IncludeUrls != nil {
		queryParams.Add("includeUrls", strconv.FormatBool(*args.IncludeUrls))
	}
	if args.IsListed != nil {
		queryParams.Add("isListed", strconv.FormatBool(*args.IsListed))
	}
	if args.IsDeleted != nil {
		queryParams.Add("isDeleted", strconv.FormatBool(*args.IsDeleted))
	}
	locationId, _ := uuid.Parse("3b331909-6a86-44cc-b9ec-c1834c35498f")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue PackageVersion
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetPackageVersion function
type GetPackageVersionArgs struct {
	// (required) Name or Id of the feed.
	FeedId *string
	// (required) Id of the package (GUID Id, not name).
	PackageId *string
	// (required) Id of the package version (GUID Id, not name).
	PackageVersionId *string
	// (optional) Project ID or project name
	Project *string
	// (optional) True to include urls for each version. Default is true.
	IncludeUrls *bool
	// (optional) Only applicable for NuGet packages. If false, delisted package versions will be returned.
	IsListed *bool
	// (optional) This does not have any effect on the requested package version, for other versions returned specifies whether to return only deleted or non-deleted versions of packages in the response. Default is unset (return all versions).
	IsDeleted *bool
}

// [Preview API] Gets provenance for a package version.
func (client *ClientImpl) GetPackageVersionProvenance(ctx context.Context, args GetPackageVersionProvenanceArgs) (*PackageVersionProvenance, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId
	if args.PackageId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PackageId"}
	}
	routeValues["packageId"] = (*args.PackageId).String()
	if args.PackageVersionId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PackageVersionId"}
	}
	routeValues["packageVersionId"] = (*args.PackageVersionId).String()

	locationId, _ := uuid.Parse("0aaeabd4-85cd-4686-8a77-8d31c15690b8")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue PackageVersionProvenance
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetPackageVersionProvenance function
type GetPackageVersionProvenanceArgs struct {
	// (required) Name or Id of the feed.
	FeedId *string
	// (required) Id of the package (GUID Id, not name).
	PackageId *uuid.UUID
	// (required) Id of the package version (GUID Id, not name).
	PackageVersionId *uuid.UUID
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Get a list of package versions, optionally filtering by state.
func (client *ClientImpl) GetPackageVersions(ctx context.Context, args GetPackageVersionsArgs) (*[]PackageVersion, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId
	if args.PackageId == nil || *args.PackageId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PackageId"}
	}
	routeValues["packageId"] = *args.PackageId

	queryParams := url.Values{}
	if args.IncludeUrls != nil {
		queryParams.Add("includeUrls", strconv.FormatBool(*args.IncludeUrls))
	}
	if args.IsListed != nil {
		queryParams.Add("isListed", strconv.FormatBool(*args.IsListed))
	}
	if args.IsDeleted != nil {
		queryParams.Add("isDeleted", strconv.FormatBool(*args.IsDeleted))
	}
	locationId, _ := uuid.Parse("3b331909-6a86-44cc-b9ec-c1834c35498f")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []PackageVersion
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetPackageVersions function
type GetPackageVersionsArgs struct {
	// (required) Name or Id of the feed.
	FeedId *string
	// (required) Id of the package (GUID Id, not name).
	PackageId *string
	// (optional) Project ID or project name
	Project *string
	// (optional) True to include urls for each version. Default is true.
	IncludeUrls *bool
	// (optional) Only applicable for NuGet packages. If false, delisted package versions will be returned.
	IsListed *bool
	// (optional) If set specifies whether to return only deleted or non-deleted versions of packages in the response. Default is unset (return all versions).
	IsDeleted *bool
}

// [Preview API] Get information about a package and all its versions within the recycle bin.
func (client *ClientImpl) GetRecycleBinPackage(ctx context.Context, args GetRecycleBinPackageArgs) (*Package, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId
	if args.PackageId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PackageId"}
	}
	routeValues["packageId"] = (*args.PackageId).String()

	queryParams := url.Values{}
	if args.IncludeUrls != nil {
		queryParams.Add("includeUrls", strconv.FormatBool(*args.IncludeUrls))
	}
	locationId, _ := uuid.Parse("2704e72c-f541-4141-99be-2004b50b05fa")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Package
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetRecycleBinPackage function
type GetRecycleBinPackageArgs struct {
	// (required) Name or Id of the feed.
	FeedId *string
	// (required) The package Id (GUID Id, not the package name).
	PackageId *uuid.UUID
	// (optional) Project ID or project name
	Project *string
	// (optional) True to return REST Urls with the response.  Default is True.
	IncludeUrls *bool
}

// [Preview API] Query for packages within the recycle bin.
func (client *ClientImpl) GetRecycleBinPackages(ctx context.Context, args GetRecycleBinPackagesArgs) (*[]Package, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId

	queryParams := url.Values{}
	if args.ProtocolType != nil {
		queryParams.Add("protocolType", *args.ProtocolType)
	}
	if args.PackageNameQuery != nil {
		queryParams.Add("packageNameQuery", *args.PackageNameQuery)
	}
	if args.IncludeUrls != nil {
		queryParams.Add("includeUrls", strconv.FormatBool(*args.IncludeUrls))
	}
	if args.Top != nil {
		queryParams.Add("$top", strconv.Itoa(*args.Top))
	}
	if args.Skip != nil {
		queryParams.Add("$skip", strconv.Itoa(*args.Skip))
	}
	if args.IncludeAllVersions != nil {
		queryParams.Add("includeAllVersions", strconv.FormatBool(*args.IncludeAllVersions))
	}
	locationId, _ := uuid.Parse("2704e72c-f541-4141-99be-2004b50b05fa")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []Package
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetRecycleBinPackages function
type GetRecycleBinPackagesArgs struct {
	// (required) Name or Id of the feed.
	FeedId *string
	// (optional) Project ID or project name
	Project *string
	// (optional) Type of package (e.g. NuGet, npm, ...).
	ProtocolType *string
	// (optional) Filter to packages matching this name.
	PackageNameQuery *string
	// (optional) True to return REST Urls with the response.  Default is True.
	IncludeUrls *bool
	// (optional) Get the top N packages.
	Top *int
	// (optional) Skip the first N packages.
	Skip *int
	// (optional) True to return all versions of the package in the response.  Default is false (latest version only).
	IncludeAllVersions *bool
}

// [Preview API] Get information about a package version within the recycle bin.
func (client *ClientImpl) GetRecycleBinPackageVersion(ctx context.Context, args GetRecycleBinPackageVersionArgs) (*RecycleBinPackageVersion, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId
	if args.PackageId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PackageId"}
	}
	routeValues["packageId"] = (*args.PackageId).String()
	if args.PackageVersionId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PackageVersionId"}
	}
	routeValues["packageVersionId"] = (*args.PackageVersionId).String()

	queryParams := url.Values{}
	if args.IncludeUrls != nil {
		queryParams.Add("includeUrls", strconv.FormatBool(*args.IncludeUrls))
	}
	locationId, _ := uuid.Parse("aceb4be7-8737-4820-834c-4c549e10fdc7")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue RecycleBinPackageVersion
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetRecycleBinPackageVersion function
type GetRecycleBinPackageVersionArgs struct {
	// (required) Name or Id of the feed.
	FeedId *string
	// (required) The package Id (GUID Id, not the package name).
	PackageId *uuid.UUID
	// (required) The package version Id 9guid Id, not the version string).
	PackageVersionId *uuid.UUID
	// (optional) Project ID or project name
	Project *string
	// (optional) True to return REST Urls with the response.  Default is True.
	IncludeUrls *bool
}

// [Preview API] Get a list of package versions within the recycle bin.
func (client *ClientImpl) GetRecycleBinPackageVersions(ctx context.Context, args GetRecycleBinPackageVersionsArgs) (*[]RecycleBinPackageVersion, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId
	if args.PackageId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PackageId"}
	}
	routeValues["packageId"] = (*args.PackageId).String()

	queryParams := url.Values{}
	if args.IncludeUrls != nil {
		queryParams.Add("includeUrls", strconv.FormatBool(*args.IncludeUrls))
	}
	locationId, _ := uuid.Parse("aceb4be7-8737-4820-834c-4c549e10fdc7")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []RecycleBinPackageVersion
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetRecycleBinPackageVersions function
type GetRecycleBinPackageVersionsArgs struct {
	// (required) Name or Id of the feed.
	FeedId *string
	// (required) The package Id (GUID Id, not the package name).
	PackageId *uuid.UUID
	// (optional) Project ID or project name
	Project *string
	// (optional) True to return REST Urls with the response.  Default is True.
	IncludeUrls *bool
}

// [Preview API]
func (client *ClientImpl) PermanentDeleteFeed(ctx context.Context, args PermanentDeleteFeedArgs) error {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId

	locationId, _ := uuid.Parse("0cee643d-beb9-41f8-9368-3ada763a8344")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the PermanentDeleteFeed function
type PermanentDeleteFeedArgs struct {
	// (required)
	FeedId *string
	// (optional) Project ID or project name
	Project *string
}

// [Preview API]
func (client *ClientImpl) QueryPackageMetrics(ctx context.Context, args QueryPackageMetricsArgs) (*[]PackageMetrics, error) {
	if args.PackageIdQuery == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PackageIdQuery"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId

	body, marshalErr := json.Marshal(*args.PackageIdQuery)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("bddc9b3c-8a59-4a9f-9b40-ee1dcaa2cc0d")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []PackageMetrics
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryPackageMetrics function
type QueryPackageMetricsArgs struct {
	// (required)
	PackageIdQuery *PackageMetricsQuery
	// (required)
	FeedId *string
	// (optional) Project ID or project name
	Project *string
}

// [Preview API]
func (client *ClientImpl) QueryPackageVersionMetrics(ctx context.Context, args QueryPackageVersionMetricsArgs) (*[]PackageVersionMetrics, error) {
	if args.PackageVersionIdQuery == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PackageVersionIdQuery"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId
	if args.PackageId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PackageId"}
	}
	routeValues["packageId"] = (*args.PackageId).String()

	body, marshalErr := json.Marshal(*args.PackageVersionIdQuery)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("e6ae8caa-b6a8-4809-b840-91b2a42c19ad")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []PackageVersionMetrics
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryPackageVersionMetrics function
type QueryPackageVersionMetricsArgs struct {
	// (required)
	PackageVersionIdQuery *PackageVersionMetricsQuery
	// (required)
	FeedId *string
	// (required)
	PackageId *uuid.UUID
	// (optional) Project ID or project name
	Project *string
}

// [Preview API]
func (client *ClientImpl) RestoreDeletedFeed(ctx context.Context, args RestoreDeletedFeedArgs) error {
	if args.PatchJson == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.PatchJson"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId

	body, marshalErr := json.Marshal(*args.PatchJson)
	if marshalErr != nil {
		return marshalErr
	}
	locationId, _ := uuid.Parse("0cee643d-beb9-41f8-9368-3ada763a8344")
	_, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json-patch+json", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the RestoreDeletedFeed function
type RestoreDeletedFeedArgs struct {
	// (required)
	PatchJson *[]webapi.JsonPatchOperation
	// (required)
	FeedId *string
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Update the permissions on a feed.
func (client *ClientImpl) SetFeedPermissions(ctx context.Context, args SetFeedPermissionsArgs) (*[]FeedPermission, error) {
	if args.FeedPermission == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.FeedPermission"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId

	body, marshalErr := json.Marshal(*args.FeedPermission)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("be8c1476-86a7-44ed-b19d-aec0e9275cd8")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []FeedPermission
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the SetFeedPermissions function
type SetFeedPermissionsArgs struct {
	// (required) Permissions to set.
	FeedPermission *[]FeedPermission
	// (required) Name or Id of the feed.
	FeedId *string
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Set the retention policy for a feed.
func (client *ClientImpl) SetFeedRetentionPolicies(ctx context.Context, args SetFeedRetentionPoliciesArgs) (*FeedRetentionPolicy, error) {
	if args.Policy == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Policy"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId

	body, marshalErr := json.Marshal(*args.Policy)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("ed52a011-0112-45b5-9f9e-e14efffb3193")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue FeedRetentionPolicy
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the SetFeedRetentionPolicies function
type SetFeedRetentionPoliciesArgs struct {
	// (required) Feed retention policy.
	Policy *FeedRetentionPolicy
	// (required) Name or ID of the feed.
	FeedId *string
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Set service-wide permissions that govern feed creation and administration.
func (client *ClientImpl) SetGlobalPermissions(ctx context.Context, args SetGlobalPermissionsArgs) (*[]GlobalPermission, error) {
	if args.GlobalPermissions == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.GlobalPermissions"}
	}
	body, marshalErr := json.Marshal(*args.GlobalPermissions)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("a74419ef-b477-43df-8758-3cd1cd5f56c6")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []GlobalPermission
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the SetGlobalPermissions function
type SetGlobalPermissionsArgs struct {
	// (required) New permissions for the organization.
	GlobalPermissions *[]GlobalPermission
}

// [Preview API] Change the attributes of a feed.
func (client *ClientImpl) UpdateFeed(ctx context.Context, args UpdateFeedArgs) (*Feed, error) {
	if args.Feed == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Feed"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId

	body, marshalErr := json.Marshal(*args.Feed)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("c65009a7-474a-4ad1-8b42-7d852107ef8c")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Feed
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateFeed function
type UpdateFeedArgs struct {
	// (required) A JSON object containing the feed settings to be updated.
	Feed *FeedUpdate
	// (required) Name or Id of the feed.
	FeedId *string
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Update a view.
func (client *ClientImpl) UpdateFeedView(ctx context.Context, args UpdateFeedViewArgs) (*FeedView, error) {
	if args.View == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.View"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FeedId == nil || *args.FeedId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FeedId"}
	}
	routeValues["feedId"] = *args.FeedId
	if args.ViewId == nil || *args.ViewId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ViewId"}
	}
	routeValues["viewId"] = *args.ViewId

	body, marshalErr := json.Marshal(*args.View)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("42a8502a-6785-41bc-8c16-89477d930877")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue FeedView
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateFeedView function
type UpdateFeedViewArgs struct {
	// (required) New settings to apply to the specified view.
	View *FeedView
	// (required) Name or Id of the feed.
	FeedId *string
	// (required) Name or Id of the view.
	ViewId *string
	// (optional) Project ID or project name
	Project *string
}