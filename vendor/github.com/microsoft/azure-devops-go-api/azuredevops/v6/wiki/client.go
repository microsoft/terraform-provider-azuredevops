// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package wiki

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/git"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

var ResourceAreaId, _ = uuid.Parse("bf7d82a0-8aa5-4613-94ef-6172a5ea01f3")

type Client interface {
	// Creates an attachment in the wiki.
	CreateAttachment(context.Context, CreateAttachmentArgs) (*WikiAttachmentResponse, error)
	// Creates or edits a wiki page.
	CreateOrUpdatePage(context.Context, CreateOrUpdatePageArgs) (*WikiPageResponse, error)
	// Creates a page move operation that updates the path and order of the page as provided in the parameters.
	CreatePageMove(context.Context, CreatePageMoveArgs) (*WikiPageMoveResponse, error)
	// Creates the wiki resource.
	CreateWiki(context.Context, CreateWikiArgs) (*WikiV2, error)
	// Deletes a wiki page.
	DeletePage(context.Context, DeletePageArgs) (*WikiPageResponse, error)
	// [Preview API] Deletes a wiki page.
	DeletePageById(context.Context, DeletePageByIdArgs) (*WikiPageResponse, error)
	// Deletes the wiki corresponding to the wiki ID or wiki name provided.
	DeleteWiki(context.Context, DeleteWikiArgs) (*WikiV2, error)
	// Gets all wikis in a project or collection.
	GetAllWikis(context.Context, GetAllWikisArgs) (*[]WikiV2, error)
	// Gets metadata or content of the wiki page for the provided path. Content negotiation is done based on the `Accept` header sent in the request.
	GetPage(context.Context, GetPageArgs) (*WikiPageResponse, error)
	// [Preview API] Gets metadata or content of the wiki page for the provided page id. Content negotiation is done based on the `Accept` header sent in the request.
	GetPageById(context.Context, GetPageByIdArgs) (*WikiPageResponse, error)
	// [Preview API] Gets metadata or content of the wiki page for the provided page id. Content negotiation is done based on the `Accept` header sent in the request.
	GetPageByIdText(context.Context, GetPageByIdTextArgs) (io.ReadCloser, error)
	// [Preview API] Gets metadata or content of the wiki page for the provided page id. Content negotiation is done based on the `Accept` header sent in the request.
	GetPageByIdZip(context.Context, GetPageByIdZipArgs) (io.ReadCloser, error)
	// [Preview API] Returns page detail corresponding to Page ID.
	GetPageData(context.Context, GetPageDataArgs) (*WikiPageDetail, error)
	// [Preview API] Returns pageable list of Wiki Pages
	GetPagesBatch(context.Context, GetPagesBatchArgs) (*GetPagesBatchResponseValue, error)
	// Gets metadata or content of the wiki page for the provided path. Content negotiation is done based on the `Accept` header sent in the request.
	GetPageText(context.Context, GetPageTextArgs) (io.ReadCloser, error)
	// Gets metadata or content of the wiki page for the provided path. Content negotiation is done based on the `Accept` header sent in the request.
	GetPageZip(context.Context, GetPageZipArgs) (io.ReadCloser, error)
	// Gets the wiki corresponding to the wiki ID or wiki name provided.
	GetWiki(context.Context, GetWikiArgs) (*WikiV2, error)
	// [Preview API] Edits a wiki page.
	UpdatePageById(context.Context, UpdatePageByIdArgs) (*WikiPageResponse, error)
	// Updates the wiki corresponding to the wiki ID or wiki name provided using the update parameters.
	UpdateWiki(context.Context, UpdateWikiArgs) (*WikiV2, error)
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

// Creates an attachment in the wiki.
func (client *ClientImpl) CreateAttachment(ctx context.Context, args CreateAttachmentArgs) (*WikiAttachmentResponse, error) {
	if args.UploadStream == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UploadStream"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WikiIdentifier == nil || *args.WikiIdentifier == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.WikiIdentifier"}
	}
	routeValues["wikiIdentifier"] = *args.WikiIdentifier

	queryParams := url.Values{}
	if args.Name == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "name"}
	}
	queryParams.Add("name", *args.Name)
	if args.VersionDescriptor != nil {
		if args.VersionDescriptor.VersionType != nil {
			queryParams.Add("versionDescriptor.versionType", string(*args.VersionDescriptor.VersionType))
		}
		if args.VersionDescriptor.Version != nil {
			queryParams.Add("versionDescriptor.version", *args.VersionDescriptor.Version)
		}
		if args.VersionDescriptor.VersionOptions != nil {
			queryParams.Add("versionDescriptor.versionOptions", string(*args.VersionDescriptor.VersionOptions))
		}
	}
	locationId, _ := uuid.Parse("c4382d8d-fefc-40e0-92c5-49852e9e17c0")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "6.0", routeValues, queryParams, args.UploadStream, "application/octet-stream", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseBodyValue WikiAttachment
	err = client.Client.UnmarshalBody(resp, &responseBodyValue)

	var responseValue *WikiAttachmentResponse
	if err == nil {
		responseValue = &WikiAttachmentResponse{
			Attachment: &responseBodyValue,
			ETag:       &[]string{resp.Header.Get("ETag")},
		}
	}

	return responseValue, err
}

// Arguments for the CreateAttachment function
type CreateAttachmentArgs struct {
	// (required) Stream to upload
	UploadStream io.Reader
	// (required) Project ID or project name
	Project *string
	// (required) Wiki ID or wiki name.
	WikiIdentifier *string
	// (required) Wiki attachment name.
	Name *string
	// (optional) GitVersionDescriptor for the page. (Optional in case of ProjectWiki).
	VersionDescriptor *git.GitVersionDescriptor
}

// Creates or edits a wiki page.
func (client *ClientImpl) CreateOrUpdatePage(ctx context.Context, args CreateOrUpdatePageArgs) (*WikiPageResponse, error) {
	if args.Parameters == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Parameters"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WikiIdentifier == nil || *args.WikiIdentifier == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.WikiIdentifier"}
	}
	routeValues["wikiIdentifier"] = *args.WikiIdentifier

	queryParams := url.Values{}
	if args.Path == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "path"}
	}
	queryParams.Add("path", *args.Path)
	if args.Comment != nil {
		queryParams.Add("comment", *args.Comment)
	}
	if args.VersionDescriptor != nil {
		if args.VersionDescriptor.VersionType != nil {
			queryParams.Add("versionDescriptor.versionType", string(*args.VersionDescriptor.VersionType))
		}
		if args.VersionDescriptor.Version != nil {
			queryParams.Add("versionDescriptor.version", *args.VersionDescriptor.Version)
		}
		if args.VersionDescriptor.VersionOptions != nil {
			queryParams.Add("versionDescriptor.versionOptions", string(*args.VersionDescriptor.VersionOptions))
		}
	}
	additionalHeaders := make(map[string]string)
	if args.Version != nil {
		additionalHeaders["If-Match"] = *args.Version
	}
	body, marshalErr := json.Marshal(*args.Parameters)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("25d3fbc7-fe3d-46cb-b5a5-0b6f79caf27b")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "6.0", routeValues, queryParams, bytes.NewReader(body), "application/json", "application/json", additionalHeaders)
	if err != nil {
		return nil, err
	}

	var responseBodyValue WikiPage
	err = client.Client.UnmarshalBody(resp, &responseBodyValue)

	var responseValue *WikiPageResponse
	if err == nil {
		responseValue = &WikiPageResponse{
			Page: &responseBodyValue,
			ETag: &[]string{resp.Header.Get("ETag")},
		}
	}

	return responseValue, err
}

// Arguments for the CreateOrUpdatePage function
type CreateOrUpdatePageArgs struct {
	// (required) Wiki create or update operation parameters.
	Parameters *WikiPageCreateOrUpdateParameters
	// (required) Project ID or project name
	Project *string
	// (required) Wiki ID or wiki name.
	WikiIdentifier *string
	// (required) Wiki page path.
	Path *string
	// (required) Version of the page on which the change is to be made. Mandatory for `Edit` scenario. To be populated in the If-Match header of the request.
	Version *string
	// (optional) Comment to be associated with the page operation.
	Comment *string
	// (optional) GitVersionDescriptor for the page. (Optional in case of ProjectWiki).
	VersionDescriptor *git.GitVersionDescriptor
}

// Creates a page move operation that updates the path and order of the page as provided in the parameters.
func (client *ClientImpl) CreatePageMove(ctx context.Context, args CreatePageMoveArgs) (*WikiPageMoveResponse, error) {
	if args.PageMoveParameters == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PageMoveParameters"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WikiIdentifier == nil || *args.WikiIdentifier == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.WikiIdentifier"}
	}
	routeValues["wikiIdentifier"] = *args.WikiIdentifier

	queryParams := url.Values{}
	if args.Comment != nil {
		queryParams.Add("comment", *args.Comment)
	}
	if args.VersionDescriptor != nil {
		if args.VersionDescriptor.VersionType != nil {
			queryParams.Add("versionDescriptor.versionType", string(*args.VersionDescriptor.VersionType))
		}
		if args.VersionDescriptor.Version != nil {
			queryParams.Add("versionDescriptor.version", *args.VersionDescriptor.Version)
		}
		if args.VersionDescriptor.VersionOptions != nil {
			queryParams.Add("versionDescriptor.versionOptions", string(*args.VersionDescriptor.VersionOptions))
		}
	}
	body, marshalErr := json.Marshal(*args.PageMoveParameters)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("e37bbe71-cbae-49e5-9a4e-949143b9d910")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0", routeValues, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseBodyValue WikiPageMove
	err = client.Client.UnmarshalBody(resp, &responseBodyValue)

	var responseValue *WikiPageMoveResponse
	if err == nil {
		responseValue = &WikiPageMoveResponse{
			PageMove: &responseBodyValue,
			ETag:     &[]string{resp.Header.Get("ETag")},
		}
	}

	return responseValue, err
}

// Arguments for the CreatePageMove function
type CreatePageMoveArgs struct {
	// (required) Page more operation parameters.
	PageMoveParameters *WikiPageMoveParameters
	// (required) Project ID or project name
	Project *string
	// (required) Wiki ID or wiki name.
	WikiIdentifier *string
	// (optional) Comment that is to be associated with this page move.
	Comment *string
	// (optional) GitVersionDescriptor for the page. (Optional in case of ProjectWiki).
	VersionDescriptor *git.GitVersionDescriptor
}

// Creates the wiki resource.
func (client *ClientImpl) CreateWiki(ctx context.Context, args CreateWikiArgs) (*WikiV2, error) {
	if args.WikiCreateParams == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.WikiCreateParams"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	body, marshalErr := json.Marshal(*args.WikiCreateParams)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("288d122c-dbd4-451d-aa5f-7dbbba070728")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WikiV2
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateWiki function
type CreateWikiArgs struct {
	// (required) Parameters for the wiki creation.
	WikiCreateParams *WikiCreateParametersV2
	// (optional) Project ID or project name
	Project *string
}

// Deletes a wiki page.
func (client *ClientImpl) DeletePage(ctx context.Context, args DeletePageArgs) (*WikiPageResponse, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WikiIdentifier == nil || *args.WikiIdentifier == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.WikiIdentifier"}
	}
	routeValues["wikiIdentifier"] = *args.WikiIdentifier

	queryParams := url.Values{}
	if args.Path == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "path"}
	}
	queryParams.Add("path", *args.Path)
	if args.Comment != nil {
		queryParams.Add("comment", *args.Comment)
	}
	if args.VersionDescriptor != nil {
		if args.VersionDescriptor.VersionType != nil {
			queryParams.Add("versionDescriptor.versionType", string(*args.VersionDescriptor.VersionType))
		}
		if args.VersionDescriptor.Version != nil {
			queryParams.Add("versionDescriptor.version", *args.VersionDescriptor.Version)
		}
		if args.VersionDescriptor.VersionOptions != nil {
			queryParams.Add("versionDescriptor.versionOptions", string(*args.VersionDescriptor.VersionOptions))
		}
	}
	locationId, _ := uuid.Parse("25d3fbc7-fe3d-46cb-b5a5-0b6f79caf27b")
	resp, err := client.Client.Send(ctx, http.MethodDelete, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseBodyValue WikiPage
	err = client.Client.UnmarshalBody(resp, &responseBodyValue)

	var responseValue *WikiPageResponse
	if err == nil {
		responseValue = &WikiPageResponse{
			Page: &responseBodyValue,
			ETag: &[]string{resp.Header.Get("ETag")},
		}
	}

	return responseValue, err
}

// Arguments for the DeletePage function
type DeletePageArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Wiki ID or wiki name.
	WikiIdentifier *string
	// (required) Wiki page path.
	Path *string
	// (optional) Comment to be associated with this page delete.
	Comment *string
	// (optional) GitVersionDescriptor for the page. (Optional in case of ProjectWiki).
	VersionDescriptor *git.GitVersionDescriptor
}

// [Preview API] Deletes a wiki page.
func (client *ClientImpl) DeletePageById(ctx context.Context, args DeletePageByIdArgs) (*WikiPageResponse, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WikiIdentifier == nil || *args.WikiIdentifier == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.WikiIdentifier"}
	}
	routeValues["wikiIdentifier"] = *args.WikiIdentifier
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)

	queryParams := url.Values{}
	if args.Comment != nil {
		queryParams.Add("comment", *args.Comment)
	}
	locationId, _ := uuid.Parse("ceddcf75-1068-452d-8b13-2d4d76e1f970")
	resp, err := client.Client.Send(ctx, http.MethodDelete, locationId, "6.0-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseBodyValue WikiPage
	err = client.Client.UnmarshalBody(resp, &responseBodyValue)

	var responseValue *WikiPageResponse
	if err == nil {
		responseValue = &WikiPageResponse{
			Page: &responseBodyValue,
			ETag: &[]string{resp.Header.Get("ETag")},
		}
	}

	return responseValue, err
}

// Arguments for the DeletePageById function
type DeletePageByIdArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Wiki ID or wiki name.
	WikiIdentifier *string
	// (required) Wiki page ID.
	Id *int
	// (optional) Comment to be associated with this page delete.
	Comment *string
}

// Deletes the wiki corresponding to the wiki ID or wiki name provided.
func (client *ClientImpl) DeleteWiki(ctx context.Context, args DeleteWikiArgs) (*WikiV2, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.WikiIdentifier == nil || *args.WikiIdentifier == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.WikiIdentifier"}
	}
	routeValues["wikiIdentifier"] = *args.WikiIdentifier

	locationId, _ := uuid.Parse("288d122c-dbd4-451d-aa5f-7dbbba070728")
	resp, err := client.Client.Send(ctx, http.MethodDelete, locationId, "6.0", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WikiV2
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the DeleteWiki function
type DeleteWikiArgs struct {
	// (required) Wiki ID or wiki name.
	WikiIdentifier *string
	// (optional) Project ID or project name
	Project *string
}

// Gets all wikis in a project or collection.
func (client *ClientImpl) GetAllWikis(ctx context.Context, args GetAllWikisArgs) (*[]WikiV2, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	locationId, _ := uuid.Parse("288d122c-dbd4-451d-aa5f-7dbbba070728")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WikiV2
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetAllWikis function
type GetAllWikisArgs struct {
	// (optional) Project ID or project name
	Project *string
}

// Gets metadata or content of the wiki page for the provided path. Content negotiation is done based on the `Accept` header sent in the request.
func (client *ClientImpl) GetPage(ctx context.Context, args GetPageArgs) (*WikiPageResponse, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WikiIdentifier == nil || *args.WikiIdentifier == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.WikiIdentifier"}
	}
	routeValues["wikiIdentifier"] = *args.WikiIdentifier

	queryParams := url.Values{}
	if args.Path != nil {
		queryParams.Add("path", *args.Path)
	}
	if args.RecursionLevel != nil {
		queryParams.Add("recursionLevel", string(*args.RecursionLevel))
	}
	if args.VersionDescriptor != nil {
		if args.VersionDescriptor.VersionType != nil {
			queryParams.Add("versionDescriptor.versionType", string(*args.VersionDescriptor.VersionType))
		}
		if args.VersionDescriptor.Version != nil {
			queryParams.Add("versionDescriptor.version", *args.VersionDescriptor.Version)
		}
		if args.VersionDescriptor.VersionOptions != nil {
			queryParams.Add("versionDescriptor.versionOptions", string(*args.VersionDescriptor.VersionOptions))
		}
	}
	if args.IncludeContent != nil {
		queryParams.Add("includeContent", strconv.FormatBool(*args.IncludeContent))
	}
	locationId, _ := uuid.Parse("25d3fbc7-fe3d-46cb-b5a5-0b6f79caf27b")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseBodyValue WikiPage
	err = client.Client.UnmarshalBody(resp, &responseBodyValue)

	var responseValue *WikiPageResponse
	if err == nil {
		responseValue = &WikiPageResponse{
			Page: &responseBodyValue,
			ETag: &[]string{resp.Header.Get("ETag")},
		}
	}

	return responseValue, err
}

// Arguments for the GetPage function
type GetPageArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Wiki ID or wiki name.
	WikiIdentifier *string
	// (optional) Wiki page path.
	Path *string
	// (optional) Recursion level for subpages retrieval. Defaults to `None` (Optional).
	RecursionLevel *git.VersionControlRecursionType
	// (optional) GitVersionDescriptor for the page. Defaults to the default branch (Optional).
	VersionDescriptor *git.GitVersionDescriptor
	// (optional) True to include the content of the page in the response for Json content type. Defaults to false (Optional)
	IncludeContent *bool
}

// [Preview API] Gets metadata or content of the wiki page for the provided page id. Content negotiation is done based on the `Accept` header sent in the request.
func (client *ClientImpl) GetPageById(ctx context.Context, args GetPageByIdArgs) (*WikiPageResponse, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WikiIdentifier == nil || *args.WikiIdentifier == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.WikiIdentifier"}
	}
	routeValues["wikiIdentifier"] = *args.WikiIdentifier
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)

	queryParams := url.Values{}
	if args.RecursionLevel != nil {
		queryParams.Add("recursionLevel", string(*args.RecursionLevel))
	}
	if args.IncludeContent != nil {
		queryParams.Add("includeContent", strconv.FormatBool(*args.IncludeContent))
	}
	locationId, _ := uuid.Parse("ceddcf75-1068-452d-8b13-2d4d76e1f970")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseBodyValue WikiPage
	err = client.Client.UnmarshalBody(resp, &responseBodyValue)

	var responseValue *WikiPageResponse
	if err == nil {
		responseValue = &WikiPageResponse{
			Page: &responseBodyValue,
			ETag: &[]string{resp.Header.Get("ETag")},
		}
	}

	return responseValue, err
}

// Arguments for the GetPageById function
type GetPageByIdArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Wiki ID or wiki name..
	WikiIdentifier *string
	// (required) Wiki page ID.
	Id *int
	// (optional) Recursion level for subpages retrieval. Defaults to `None` (Optional).
	RecursionLevel *git.VersionControlRecursionType
	// (optional) True to include the content of the page in the response for Json content type. Defaults to false (Optional)
	IncludeContent *bool
}

// [Preview API] Gets metadata or content of the wiki page for the provided page id. Content negotiation is done based on the `Accept` header sent in the request.
func (client *ClientImpl) GetPageByIdText(ctx context.Context, args GetPageByIdTextArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WikiIdentifier == nil || *args.WikiIdentifier == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.WikiIdentifier"}
	}
	routeValues["wikiIdentifier"] = *args.WikiIdentifier
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)

	queryParams := url.Values{}
	if args.RecursionLevel != nil {
		queryParams.Add("recursionLevel", string(*args.RecursionLevel))
	}
	if args.IncludeContent != nil {
		queryParams.Add("includeContent", strconv.FormatBool(*args.IncludeContent))
	}
	locationId, _ := uuid.Parse("ceddcf75-1068-452d-8b13-2d4d76e1f970")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, queryParams, nil, "", "text/plain", nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetPageByIdText function
type GetPageByIdTextArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Wiki ID or wiki name..
	WikiIdentifier *string
	// (required) Wiki page ID.
	Id *int
	// (optional) Recursion level for subpages retrieval. Defaults to `None` (Optional).
	RecursionLevel *git.VersionControlRecursionType
	// (optional) True to include the content of the page in the response for Json content type. Defaults to false (Optional)
	IncludeContent *bool
}

// [Preview API] Gets metadata or content of the wiki page for the provided page id. Content negotiation is done based on the `Accept` header sent in the request.
func (client *ClientImpl) GetPageByIdZip(ctx context.Context, args GetPageByIdZipArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WikiIdentifier == nil || *args.WikiIdentifier == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.WikiIdentifier"}
	}
	routeValues["wikiIdentifier"] = *args.WikiIdentifier
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)

	queryParams := url.Values{}
	if args.RecursionLevel != nil {
		queryParams.Add("recursionLevel", string(*args.RecursionLevel))
	}
	if args.IncludeContent != nil {
		queryParams.Add("includeContent", strconv.FormatBool(*args.IncludeContent))
	}
	locationId, _ := uuid.Parse("ceddcf75-1068-452d-8b13-2d4d76e1f970")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, queryParams, nil, "", "application/zip", nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetPageByIdZip function
type GetPageByIdZipArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Wiki ID or wiki name..
	WikiIdentifier *string
	// (required) Wiki page ID.
	Id *int
	// (optional) Recursion level for subpages retrieval. Defaults to `None` (Optional).
	RecursionLevel *git.VersionControlRecursionType
	// (optional) True to include the content of the page in the response for Json content type. Defaults to false (Optional)
	IncludeContent *bool
}

// [Preview API] Returns page detail corresponding to Page ID.
func (client *ClientImpl) GetPageData(ctx context.Context, args GetPageDataArgs) (*WikiPageDetail, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WikiIdentifier == nil || *args.WikiIdentifier == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.WikiIdentifier"}
	}
	routeValues["wikiIdentifier"] = *args.WikiIdentifier
	if args.PageId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PageId"}
	}
	routeValues["pageId"] = strconv.Itoa(*args.PageId)

	queryParams := url.Values{}
	if args.PageViewsForDays != nil {
		queryParams.Add("pageViewsForDays", strconv.Itoa(*args.PageViewsForDays))
	}
	locationId, _ := uuid.Parse("81c4e0fe-7663-4d62-ad46-6ab78459f274")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WikiPageDetail
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetPageData function
type GetPageDataArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Wiki ID or wiki name.
	WikiIdentifier *string
	// (required) Wiki page ID.
	PageId *int
	// (optional) last N days from the current day for which page views is to be returned. It's inclusive of current day.
	PageViewsForDays *int
}

// [Preview API] Returns pageable list of Wiki Pages
func (client *ClientImpl) GetPagesBatch(ctx context.Context, args GetPagesBatchArgs) (*GetPagesBatchResponseValue, error) {
	if args.PagesBatchRequest == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PagesBatchRequest"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WikiIdentifier == nil || *args.WikiIdentifier == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.WikiIdentifier"}
	}
	routeValues["wikiIdentifier"] = *args.WikiIdentifier

	queryParams := url.Values{}
	if args.VersionDescriptor != nil {
		if args.VersionDescriptor.VersionType != nil {
			queryParams.Add("versionDescriptor.versionType", string(*args.VersionDescriptor.VersionType))
		}
		if args.VersionDescriptor.Version != nil {
			queryParams.Add("versionDescriptor.version", *args.VersionDescriptor.Version)
		}
		if args.VersionDescriptor.VersionOptions != nil {
			queryParams.Add("versionDescriptor.versionOptions", string(*args.VersionDescriptor.VersionOptions))
		}
	}
	body, marshalErr := json.Marshal(*args.PagesBatchRequest)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("71323c46-2592-4398-8771-ced73dd87207")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0-preview.1", routeValues, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue GetPagesBatchResponseValue
	responseValue.ContinuationToken = resp.Header.Get(azuredevops.HeaderKeyContinuationToken)
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue.Value)
	return &responseValue, err
}

// Arguments for the GetPagesBatch function
type GetPagesBatchArgs struct {
	// (required) Wiki batch page request.
	PagesBatchRequest *WikiPagesBatchRequest
	// (required) Project ID or project name
	Project *string
	// (required) Wiki ID or wiki name.
	WikiIdentifier *string
	// (optional) GitVersionDescriptor for the page. (Optional in case of ProjectWiki).
	VersionDescriptor *git.GitVersionDescriptor
}

// Return type for the GetPagesBatch function
type GetPagesBatchResponseValue struct {
	Value             []WikiPageDetail
	ContinuationToken string
}

// Gets metadata or content of the wiki page for the provided path. Content negotiation is done based on the `Accept` header sent in the request.
func (client *ClientImpl) GetPageText(ctx context.Context, args GetPageTextArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WikiIdentifier == nil || *args.WikiIdentifier == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.WikiIdentifier"}
	}
	routeValues["wikiIdentifier"] = *args.WikiIdentifier

	queryParams := url.Values{}
	if args.Path != nil {
		queryParams.Add("path", *args.Path)
	}
	if args.RecursionLevel != nil {
		queryParams.Add("recursionLevel", string(*args.RecursionLevel))
	}
	if args.VersionDescriptor != nil {
		if args.VersionDescriptor.VersionType != nil {
			queryParams.Add("versionDescriptor.versionType", string(*args.VersionDescriptor.VersionType))
		}
		if args.VersionDescriptor.Version != nil {
			queryParams.Add("versionDescriptor.version", *args.VersionDescriptor.Version)
		}
		if args.VersionDescriptor.VersionOptions != nil {
			queryParams.Add("versionDescriptor.versionOptions", string(*args.VersionDescriptor.VersionOptions))
		}
	}
	if args.IncludeContent != nil {
		queryParams.Add("includeContent", strconv.FormatBool(*args.IncludeContent))
	}
	locationId, _ := uuid.Parse("25d3fbc7-fe3d-46cb-b5a5-0b6f79caf27b")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "text/plain", nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetPageText function
type GetPageTextArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Wiki ID or wiki name.
	WikiIdentifier *string
	// (optional) Wiki page path.
	Path *string
	// (optional) Recursion level for subpages retrieval. Defaults to `None` (Optional).
	RecursionLevel *git.VersionControlRecursionType
	// (optional) GitVersionDescriptor for the page. Defaults to the default branch (Optional).
	VersionDescriptor *git.GitVersionDescriptor
	// (optional) True to include the content of the page in the response for Json content type. Defaults to false (Optional)
	IncludeContent *bool
}

// Gets metadata or content of the wiki page for the provided path. Content negotiation is done based on the `Accept` header sent in the request.
func (client *ClientImpl) GetPageZip(ctx context.Context, args GetPageZipArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WikiIdentifier == nil || *args.WikiIdentifier == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.WikiIdentifier"}
	}
	routeValues["wikiIdentifier"] = *args.WikiIdentifier

	queryParams := url.Values{}
	if args.Path != nil {
		queryParams.Add("path", *args.Path)
	}
	if args.RecursionLevel != nil {
		queryParams.Add("recursionLevel", string(*args.RecursionLevel))
	}
	if args.VersionDescriptor != nil {
		if args.VersionDescriptor.VersionType != nil {
			queryParams.Add("versionDescriptor.versionType", string(*args.VersionDescriptor.VersionType))
		}
		if args.VersionDescriptor.Version != nil {
			queryParams.Add("versionDescriptor.version", *args.VersionDescriptor.Version)
		}
		if args.VersionDescriptor.VersionOptions != nil {
			queryParams.Add("versionDescriptor.versionOptions", string(*args.VersionDescriptor.VersionOptions))
		}
	}
	if args.IncludeContent != nil {
		queryParams.Add("includeContent", strconv.FormatBool(*args.IncludeContent))
	}
	locationId, _ := uuid.Parse("25d3fbc7-fe3d-46cb-b5a5-0b6f79caf27b")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/zip", nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetPageZip function
type GetPageZipArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Wiki ID or wiki name.
	WikiIdentifier *string
	// (optional) Wiki page path.
	Path *string
	// (optional) Recursion level for subpages retrieval. Defaults to `None` (Optional).
	RecursionLevel *git.VersionControlRecursionType
	// (optional) GitVersionDescriptor for the page. Defaults to the default branch (Optional).
	VersionDescriptor *git.GitVersionDescriptor
	// (optional) True to include the content of the page in the response for Json content type. Defaults to false (Optional)
	IncludeContent *bool
}

// Gets the wiki corresponding to the wiki ID or wiki name provided.
func (client *ClientImpl) GetWiki(ctx context.Context, args GetWikiArgs) (*WikiV2, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.WikiIdentifier == nil || *args.WikiIdentifier == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.WikiIdentifier"}
	}
	routeValues["wikiIdentifier"] = *args.WikiIdentifier

	locationId, _ := uuid.Parse("288d122c-dbd4-451d-aa5f-7dbbba070728")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WikiV2
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWiki function
type GetWikiArgs struct {
	// (required) Wiki ID or wiki name.
	WikiIdentifier *string
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Edits a wiki page.
func (client *ClientImpl) UpdatePageById(ctx context.Context, args UpdatePageByIdArgs) (*WikiPageResponse, error) {
	if args.Parameters == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Parameters"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WikiIdentifier == nil || *args.WikiIdentifier == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.WikiIdentifier"}
	}
	routeValues["wikiIdentifier"] = *args.WikiIdentifier
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)

	queryParams := url.Values{}
	if args.Comment != nil {
		queryParams.Add("comment", *args.Comment)
	}
	additionalHeaders := make(map[string]string)
	if args.Version != nil {
		additionalHeaders["If-Match"] = *args.Version
	}
	body, marshalErr := json.Marshal(*args.Parameters)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("ceddcf75-1068-452d-8b13-2d4d76e1f970")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "6.0-preview.1", routeValues, queryParams, bytes.NewReader(body), "application/json", "application/json", additionalHeaders)
	if err != nil {
		return nil, err
	}

	var responseBodyValue WikiPage
	err = client.Client.UnmarshalBody(resp, &responseBodyValue)

	var responseValue *WikiPageResponse
	if err == nil {
		responseValue = &WikiPageResponse{
			Page: &responseBodyValue,
			ETag: &[]string{resp.Header.Get("ETag")},
		}
	}

	return responseValue, err
}

// Arguments for the UpdatePageById function
type UpdatePageByIdArgs struct {
	// (required) Wiki update operation parameters.
	Parameters *WikiPageCreateOrUpdateParameters
	// (required) Project ID or project name
	Project *string
	// (required) Wiki ID or wiki name.
	WikiIdentifier *string
	// (required) Wiki page ID.
	Id *int
	// (required) Version of the page on which the change is to be made. Mandatory for `Edit` scenario. To be populated in the If-Match header of the request.
	Version *string
	// (optional) Comment to be associated with the page operation.
	Comment *string
}

// Updates the wiki corresponding to the wiki ID or wiki name provided using the update parameters.
func (client *ClientImpl) UpdateWiki(ctx context.Context, args UpdateWikiArgs) (*WikiV2, error) {
	if args.UpdateParameters == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UpdateParameters"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.WikiIdentifier == nil || *args.WikiIdentifier == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.WikiIdentifier"}
	}
	routeValues["wikiIdentifier"] = *args.WikiIdentifier

	body, marshalErr := json.Marshal(*args.UpdateParameters)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("288d122c-dbd4-451d-aa5f-7dbbba070728")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "6.0", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WikiV2
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateWiki function
type UpdateWikiArgs struct {
	// (required) Update parameters.
	UpdateParameters *WikiUpdateParameters
	// (required) Wiki ID or wiki name.
	WikiIdentifier *string
	// (optional) Project ID or project name
	Project *string
}
