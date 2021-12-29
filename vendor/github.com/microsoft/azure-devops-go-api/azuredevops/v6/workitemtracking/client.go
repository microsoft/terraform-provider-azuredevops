// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package workitemtracking

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/webapi"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var ResourceAreaId, _ = uuid.Parse("5264459e-e5e0-4bd8-b118-0985e68a4ec5")

type Client interface {
	// [Preview API] Add a comment on a work item.
	AddComment(context.Context, AddCommentArgs) (*Comment, error)
	// Uploads an attachment.
	CreateAttachment(context.Context, CreateAttachmentArgs) (*AttachmentReference, error)
	// [Preview API] Adds a new reaction to a comment.
	CreateCommentReaction(context.Context, CreateCommentReactionArgs) (*CommentReaction, error)
	// Create a new field.
	CreateField(context.Context, CreateFieldArgs) (*WorkItemField, error)
	// Create new or update an existing classification node.
	CreateOrUpdateClassificationNode(context.Context, CreateOrUpdateClassificationNodeArgs) (*WorkItemClassificationNode, error)
	// Creates a query, or moves a query.
	CreateQuery(context.Context, CreateQueryArgs) (*QueryHierarchyItem, error)
	// [Preview API] Creates a template
	CreateTemplate(context.Context, CreateTemplateArgs) (*WorkItemTemplate, error)
	// Creates a single work item.
	CreateWorkItem(context.Context, CreateWorkItemArgs) (*WorkItem, error)
	// Delete an existing classification node.
	DeleteClassificationNode(context.Context, DeleteClassificationNodeArgs) error
	// [Preview API] Delete a comment on a work item.
	DeleteComment(context.Context, DeleteCommentArgs) error
	// [Preview API] Deletes an existing reaction on a comment.
	DeleteCommentReaction(context.Context, DeleteCommentReactionArgs) (*CommentReaction, error)
	// Deletes the field. To undelete a filed, see "Update Field" API.
	DeleteField(context.Context, DeleteFieldArgs) error
	// Delete a query or a folder. This deletes any permission change on the deleted query or folder and any of its descendants if it is a folder. It is important to note that the deleted permission changes cannot be recovered upon undeleting the query or folder.
	DeleteQuery(context.Context, DeleteQueryArgs) error
	// [Preview API]
	DeleteTag(context.Context, DeleteTagArgs) error
	// [Preview API] Deletes the template with given id
	DeleteTemplate(context.Context, DeleteTemplateArgs) error
	// Deletes the specified work item and sends it to the Recycle Bin, so that it can be restored back, if required. Optionally, if the destroy parameter has been set to true, it destroys the work item permanently. WARNING: If the destroy parameter is set to true, work items deleted by this command will NOT go to recycle-bin and there is no way to restore/recover them after deletion. It is recommended NOT to use this parameter. If you do, please use this parameter with extreme caution.
	DeleteWorkItem(context.Context, DeleteWorkItemArgs) (*WorkItemDelete, error)
	// Destroys the specified work item permanently from the Recycle Bin. This action can not be undone.
	DestroyWorkItem(context.Context, DestroyWorkItemArgs) error
	// Downloads an attachment.
	GetAttachmentContent(context.Context, GetAttachmentContentArgs) (io.ReadCloser, error)
	// Downloads an attachment.
	GetAttachmentZip(context.Context, GetAttachmentZipArgs) (io.ReadCloser, error)
	// Gets the classification node for a given node path.
	GetClassificationNode(context.Context, GetClassificationNodeArgs) (*WorkItemClassificationNode, error)
	// Gets root classification nodes or list of classification nodes for a given list of nodes ids, for a given project. In case ids parameter is supplied you will  get list of classification nodes for those ids. Otherwise you will get root classification nodes for this project.
	GetClassificationNodes(context.Context, GetClassificationNodesArgs) (*[]WorkItemClassificationNode, error)
	// [Preview API] Returns a work item comment.
	GetComment(context.Context, GetCommentArgs) (*Comment, error)
	// [Preview API] Gets reactions of a comment.
	GetCommentReactions(context.Context, GetCommentReactionsArgs) (*[]CommentReaction, error)
	// [Preview API] Returns a list of work item comments, pageable.
	GetComments(context.Context, GetCommentsArgs) (*CommentList, error)
	// [Preview API] Returns a list of work item comments by ids.
	GetCommentsBatch(context.Context, GetCommentsBatchArgs) (*CommentList, error)
	// [Preview API]
	GetCommentVersion(context.Context, GetCommentVersionArgs) (*CommentVersion, error)
	// [Preview API]
	GetCommentVersions(context.Context, GetCommentVersionsArgs) (*[]CommentVersion, error)
	// Gets a deleted work item from Recycle Bin.
	GetDeletedWorkItem(context.Context, GetDeletedWorkItemArgs) (*WorkItemDelete, error)
	// Gets the work items from the recycle bin, whose IDs have been specified in the parameters
	GetDeletedWorkItems(context.Context, GetDeletedWorkItemsArgs) (*[]WorkItemDeleteReference, error)
	// Gets a list of the IDs and the URLs of the deleted the work items in the Recycle Bin.
	GetDeletedWorkItemShallowReferences(context.Context, GetDeletedWorkItemShallowReferencesArgs) (*[]WorkItemDeleteShallowReference, error)
	// [Preview API] Get users who reacted on the comment.
	GetEngagedUsers(context.Context, GetEngagedUsersArgs) (*[]webapi.IdentityRef, error)
	// Gets information on a specific field.
	GetField(context.Context, GetFieldArgs) (*WorkItemField, error)
	// Returns information for all fields. The project ID/name parameter is optional.
	GetFields(context.Context, GetFieldsArgs) (*[]WorkItemField, error)
	// Gets the root queries and their children
	GetQueries(context.Context, GetQueriesArgs) (*[]QueryHierarchyItem, error)
	// Gets a list of queries by ids (Maximum 1000)
	GetQueriesBatch(context.Context, GetQueriesBatchArgs) (*[]QueryHierarchyItem, error)
	// Retrieves an individual query and its children
	GetQuery(context.Context, GetQueryArgs) (*QueryHierarchyItem, error)
	// Gets the results of the query given the query ID.
	GetQueryResultCount(context.Context, GetQueryResultCountArgs) (*int, error)
	// [Preview API] Gets recent work item activities
	GetRecentActivityData(context.Context, GetRecentActivityDataArgs) (*[]AccountRecentActivityWorkItemModel2, error)
	// Gets the work item relation type definition.
	GetRelationType(context.Context, GetRelationTypeArgs) (*WorkItemRelationType, error)
	// Gets the work item relation types.
	GetRelationTypes(context.Context, GetRelationTypesArgs) (*[]WorkItemRelationType, error)
	// Get a batch of work item links
	GetReportingLinksByLinkType(context.Context, GetReportingLinksByLinkTypeArgs) (*ReportingWorkItemLinksBatch, error)
	// Returns a fully hydrated work item for the requested revision
	GetRevision(context.Context, GetRevisionArgs) (*WorkItem, error)
	// Returns the list of fully hydrated work item revisions, paged.
	GetRevisions(context.Context, GetRevisionsArgs) (*[]WorkItem, error)
	// Gets root classification nodes under the project.
	GetRootNodes(context.Context, GetRootNodesArgs) (*[]WorkItemClassificationNode, error)
	// [Preview API]
	GetTag(context.Context, GetTagArgs) (*WorkItemTagDefinition, error)
	// [Preview API]
	GetTags(context.Context, GetTagsArgs) (*[]WorkItemTagDefinition, error)
	// [Preview API] Gets the template with specified id
	GetTemplate(context.Context, GetTemplateArgs) (*WorkItemTemplate, error)
	// [Preview API] Gets template
	GetTemplates(context.Context, GetTemplatesArgs) (*[]WorkItemTemplateReference, error)
	// Returns a single update for a work item
	GetUpdate(context.Context, GetUpdateArgs) (*WorkItemUpdate, error)
	// Returns a the deltas between work item revisions
	GetUpdates(context.Context, GetUpdatesArgs) (*[]WorkItemUpdate, error)
	// [Preview API] Get the list of work item tracking outbound artifact link types.
	GetWorkArtifactLinkTypes(context.Context, GetWorkArtifactLinkTypesArgs) (*[]WorkArtifactLink, error)
	// Returns a single work item.
	GetWorkItem(context.Context, GetWorkItemArgs) (*WorkItem, error)
	// Get a work item icon given the friendly name and icon color.
	GetWorkItemIconJson(context.Context, GetWorkItemIconJsonArgs) (*WorkItemIcon, error)
	// Get a list of all work item icons.
	GetWorkItemIcons(context.Context, GetWorkItemIconsArgs) (*[]WorkItemIcon, error)
	// Get a work item icon given the friendly name and icon color.
	GetWorkItemIconSvg(context.Context, GetWorkItemIconSvgArgs) (io.ReadCloser, error)
	// Get a work item icon given the friendly name and icon color.
	GetWorkItemIconXaml(context.Context, GetWorkItemIconXamlArgs) (io.ReadCloser, error)
	// [Preview API] Returns the next state on the given work item IDs.
	GetWorkItemNextStatesOnCheckinAction(context.Context, GetWorkItemNextStatesOnCheckinActionArgs) (*[]WorkItemNextStateOnTransition, error)
	// Returns a list of work items (Maximum 200)
	GetWorkItems(context.Context, GetWorkItemsArgs) (*[]WorkItem, error)
	// Gets work items for a list of work item ids (Maximum 200)
	GetWorkItemsBatch(context.Context, GetWorkItemsBatchArgs) (*[]WorkItem, error)
	// Returns a single work item from a template.
	GetWorkItemTemplate(context.Context, GetWorkItemTemplateArgs) (*WorkItem, error)
	// Returns a work item type definition.
	GetWorkItemType(context.Context, GetWorkItemTypeArgs) (*WorkItemType, error)
	// Get all work item type categories.
	GetWorkItemTypeCategories(context.Context, GetWorkItemTypeCategoriesArgs) (*[]WorkItemTypeCategory, error)
	// Get specific work item type category by name.
	GetWorkItemTypeCategory(context.Context, GetWorkItemTypeCategoryArgs) (*WorkItemTypeCategory, error)
	// Get a list of fields for a work item type with detailed references.
	GetWorkItemTypeFieldsWithReferences(context.Context, GetWorkItemTypeFieldsWithReferencesArgs) (*[]WorkItemTypeFieldWithReferences, error)
	// Get a field for a work item type with detailed references.
	GetWorkItemTypeFieldWithReferences(context.Context, GetWorkItemTypeFieldWithReferencesArgs) (*WorkItemTypeFieldWithReferences, error)
	// Returns the list of work item types
	GetWorkItemTypes(context.Context, GetWorkItemTypesArgs) (*[]WorkItemType, error)
	// [Preview API] Returns the state names and colors for a work item type.
	GetWorkItemTypeStates(context.Context, GetWorkItemTypeStatesArgs) (*[]WorkItemStateColor, error)
	// [Preview API] Migrates a project to a different process within the same OOB type. For example, you can only migrate a project from agile/custom-agile to agile/custom-agile.
	MigrateProjectsProcess(context.Context, MigrateProjectsProcessArgs) (*ProcessMigrationResultModel, error)
	// Gets the results of the query given the query ID.
	QueryById(context.Context, QueryByIdArgs) (*WorkItemQueryResult, error)
	// Gets the results of the query given its WIQL.
	QueryByWiql(context.Context, QueryByWiqlArgs) (*WorkItemQueryResult, error)
	// [Preview API] Queries work items linked to a given list of artifact URI.
	QueryWorkItemsForArtifactUris(context.Context, QueryWorkItemsForArtifactUrisArgs) (*ArtifactUriQueryResult, error)
	// [Preview API]
	ReadReportingDiscussions(context.Context, ReadReportingDiscussionsArgs) (*ReportingWorkItemRevisionsBatch, error)
	// Get a batch of work item revisions with the option of including deleted items
	ReadReportingRevisionsGet(context.Context, ReadReportingRevisionsGetArgs) (*ReportingWorkItemRevisionsBatch, error)
	// Get a batch of work item revisions. This request may be used if your list of fields is large enough that it may run the URL over the length limit.
	ReadReportingRevisionsPost(context.Context, ReadReportingRevisionsPostArgs) (*ReportingWorkItemRevisionsBatch, error)
	// [Preview API] Replace template contents
	ReplaceTemplate(context.Context, ReplaceTemplateArgs) (*WorkItemTemplate, error)
	// Restores the deleted work item from Recycle Bin.
	RestoreWorkItem(context.Context, RestoreWorkItemArgs) (*WorkItemDelete, error)
	// Searches all queries the user has access to in the current project
	SearchQueries(context.Context, SearchQueriesArgs) (*QueryHierarchyItemsResult, error)
	// Update an existing classification node.
	UpdateClassificationNode(context.Context, UpdateClassificationNodeArgs) (*WorkItemClassificationNode, error)
	// [Preview API] Update a comment on a work item.
	UpdateComment(context.Context, UpdateCommentArgs) (*Comment, error)
	// Update a field.
	UpdateField(context.Context, UpdateFieldArgs) (*WorkItemField, error)
	// Update a query or a folder. This allows you to update, rename and move queries and folders.
	UpdateQuery(context.Context, UpdateQueryArgs) (*QueryHierarchyItem, error)
	// [Preview API]
	UpdateTag(context.Context, UpdateTagArgs) (*WorkItemTagDefinition, error)
	// Updates a single work item.
	UpdateWorkItem(context.Context, UpdateWorkItemArgs) (*WorkItem, error)
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

// [Preview API] Add a comment on a work item.
func (client *ClientImpl) AddComment(ctx context.Context, args AddCommentArgs) (*Comment, error) {
	if args.Request == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Request"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WorkItemId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.WorkItemId"}
	}
	routeValues["workItemId"] = strconv.Itoa(*args.WorkItemId)

	body, marshalErr := json.Marshal(*args.Request)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("608aac0a-32e1-4493-a863-b9cf4566d257")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0-preview.3", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Comment
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the AddComment function
type AddCommentArgs struct {
	// (required) Comment create request.
	Request *CommentCreate
	// (required) Project ID or project name
	Project *string
	// (required) Id of a work item.
	WorkItemId *int
}

// Uploads an attachment.
func (client *ClientImpl) CreateAttachment(ctx context.Context, args CreateAttachmentArgs) (*AttachmentReference, error) {
	if args.UploadStream == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UploadStream"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	queryParams := url.Values{}
	if args.FileName != nil {
		queryParams.Add("fileName", *args.FileName)
	}
	if args.UploadType != nil {
		queryParams.Add("uploadType", *args.UploadType)
	}
	if args.AreaPath != nil {
		queryParams.Add("areaPath", *args.AreaPath)
	}
	locationId, _ := uuid.Parse("e07b5fa4-1499-494d-a496-64b860fd64ff")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0", routeValues, queryParams, args.UploadStream, "application/octet-stream", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue AttachmentReference
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateAttachment function
type CreateAttachmentArgs struct {
	// (required) Stream to upload
	UploadStream io.Reader
	// (optional) Project ID or project name
	Project *string
	// (optional) The name of the file
	FileName *string
	// (optional) Attachment upload type: Simple or Chunked
	UploadType *string
	// (optional) Target project Area Path
	AreaPath *string
}

// [Preview API] Adds a new reaction to a comment.
func (client *ClientImpl) CreateCommentReaction(ctx context.Context, args CreateCommentReactionArgs) (*CommentReaction, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WorkItemId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.WorkItemId"}
	}
	routeValues["workItemId"] = strconv.Itoa(*args.WorkItemId)
	if args.CommentId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.CommentId"}
	}
	routeValues["commentId"] = strconv.Itoa(*args.CommentId)
	if args.ReactionType == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ReactionType"}
	}
	routeValues["reactionType"] = string(*args.ReactionType)

	locationId, _ := uuid.Parse("f6cb3f27-1028-4851-af96-887e570dc21f")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue CommentReaction
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateCommentReaction function
type CreateCommentReactionArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) WorkItem ID
	WorkItemId *int
	// (required) Comment ID
	CommentId *int
	// (required) Type of the reaction
	ReactionType *CommentReactionType
}

// Create a new field.
func (client *ClientImpl) CreateField(ctx context.Context, args CreateFieldArgs) (*WorkItemField, error) {
	if args.WorkItemField == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.WorkItemField"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	body, marshalErr := json.Marshal(*args.WorkItemField)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("b51fd764-e5c2-4b9b-aaf7-3395cf4bdd94")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemField
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateField function
type CreateFieldArgs struct {
	// (required) New field definition
	WorkItemField *WorkItemField
	// (optional) Project ID or project name
	Project *string
}

// Create new or update an existing classification node.
func (client *ClientImpl) CreateOrUpdateClassificationNode(ctx context.Context, args CreateOrUpdateClassificationNodeArgs) (*WorkItemClassificationNode, error) {
	if args.PostedNode == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PostedNode"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.StructureGroup == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.StructureGroup"}
	}
	routeValues["structureGroup"] = string(*args.StructureGroup)
	if args.Path != nil && *args.Path != "" {
		routeValues["path"] = *args.Path
	}

	body, marshalErr := json.Marshal(*args.PostedNode)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("5a172953-1b41-49d3-840a-33f79c3ce89f")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemClassificationNode
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateOrUpdateClassificationNode function
type CreateOrUpdateClassificationNodeArgs struct {
	// (required) Node to create or update.
	PostedNode *WorkItemClassificationNode
	// (required) Project ID or project name
	Project *string
	// (required) Structure group of the classification node, area or iteration.
	StructureGroup *TreeStructureGroup
	// (optional) Path of the classification node.
	Path *string
}

// Creates a query, or moves a query.
func (client *ClientImpl) CreateQuery(ctx context.Context, args CreateQueryArgs) (*QueryHierarchyItem, error) {
	if args.PostedQuery == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PostedQuery"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Query == nil || *args.Query == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Query"}
	}
	routeValues["query"] = *args.Query

	queryParams := url.Values{}
	if args.ValidateWiqlOnly != nil {
		queryParams.Add("validateWiqlOnly", strconv.FormatBool(*args.ValidateWiqlOnly))
	}
	body, marshalErr := json.Marshal(*args.PostedQuery)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("a67d190c-c41f-424b-814d-0e906f659301")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0", routeValues, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue QueryHierarchyItem
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateQuery function
type CreateQueryArgs struct {
	// (required) The query to create.
	PostedQuery *QueryHierarchyItem
	// (required) Project ID or project name
	Project *string
	// (required) The parent id or path under which the query is to be created.
	Query *string
	// (optional) If you only want to validate your WIQL query without actually creating one, set it to true. Default is false.
	ValidateWiqlOnly *bool
}

// [Preview API] Creates a template
func (client *ClientImpl) CreateTemplate(ctx context.Context, args CreateTemplateArgs) (*WorkItemTemplate, error) {
	if args.Template == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Template"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team == nil || *args.Team == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Team"}
	}
	routeValues["team"] = *args.Team

	body, marshalErr := json.Marshal(*args.Template)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("6a90345f-a676-4969-afce-8e163e1d5642")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemTemplate
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateTemplate function
type CreateTemplateArgs struct {
	// (required) Template contents
	Template *WorkItemTemplate
	// (required) Project ID or project name
	Project *string
	// (required) Team ID or team name
	Team *string
}

// Creates a single work item.
func (client *ClientImpl) CreateWorkItem(ctx context.Context, args CreateWorkItemArgs) (*WorkItem, error) {
	if args.Document == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Document"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Type == nil || *args.Type == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Type"}
	}
	routeValues["type"] = *args.Type

	queryParams := url.Values{}
	if args.ValidateOnly != nil {
		queryParams.Add("validateOnly", strconv.FormatBool(*args.ValidateOnly))
	}
	if args.BypassRules != nil {
		queryParams.Add("bypassRules", strconv.FormatBool(*args.BypassRules))
	}
	if args.SuppressNotifications != nil {
		queryParams.Add("suppressNotifications", strconv.FormatBool(*args.SuppressNotifications))
	}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	body, marshalErr := json.Marshal(*args.Document)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("62d3d110-0047-428c-ad3c-4fe872c91c74")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0", routeValues, queryParams, bytes.NewReader(body), "application/json-patch+json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItem
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateWorkItem function
type CreateWorkItemArgs struct {
	// (required) The JSON Patch document representing the work item
	Document *[]webapi.JsonPatchOperation
	// (required) Project ID or project name
	Project *string
	// (required) The work item type of the work item to create
	Type *string
	// (optional) Indicate if you only want to validate the changes without saving the work item
	ValidateOnly *bool
	// (optional) Do not enforce the work item type rules on this update
	BypassRules *bool
	// (optional) Do not fire any notifications for this change
	SuppressNotifications *bool
	// (optional) The expand parameters for work item attributes. Possible options are { None, Relations, Fields, Links, All }.
	Expand *WorkItemExpand
}

// Delete an existing classification node.
func (client *ClientImpl) DeleteClassificationNode(ctx context.Context, args DeleteClassificationNodeArgs) error {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.StructureGroup == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.StructureGroup"}
	}
	routeValues["structureGroup"] = string(*args.StructureGroup)
	if args.Path != nil && *args.Path != "" {
		routeValues["path"] = *args.Path
	}

	queryParams := url.Values{}
	if args.ReclassifyId != nil {
		queryParams.Add("$reclassifyId", strconv.Itoa(*args.ReclassifyId))
	}
	locationId, _ := uuid.Parse("5a172953-1b41-49d3-840a-33f79c3ce89f")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteClassificationNode function
type DeleteClassificationNodeArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Structure group of the classification node, area or iteration.
	StructureGroup *TreeStructureGroup
	// (optional) Path of the classification node.
	Path *string
	// (optional) Id of the target classification node for reclassification.
	ReclassifyId *int
}

// [Preview API] Delete a comment on a work item.
func (client *ClientImpl) DeleteComment(ctx context.Context, args DeleteCommentArgs) error {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WorkItemId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.WorkItemId"}
	}
	routeValues["workItemId"] = strconv.Itoa(*args.WorkItemId)
	if args.CommentId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.CommentId"}
	}
	routeValues["commentId"] = strconv.Itoa(*args.CommentId)

	locationId, _ := uuid.Parse("608aac0a-32e1-4493-a863-b9cf4566d257")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "6.0-preview.3", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteComment function
type DeleteCommentArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Id of a work item.
	WorkItemId *int
	// (required)
	CommentId *int
}

// [Preview API] Deletes an existing reaction on a comment.
func (client *ClientImpl) DeleteCommentReaction(ctx context.Context, args DeleteCommentReactionArgs) (*CommentReaction, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WorkItemId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.WorkItemId"}
	}
	routeValues["workItemId"] = strconv.Itoa(*args.WorkItemId)
	if args.CommentId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.CommentId"}
	}
	routeValues["commentId"] = strconv.Itoa(*args.CommentId)
	if args.ReactionType == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ReactionType"}
	}
	routeValues["reactionType"] = string(*args.ReactionType)

	locationId, _ := uuid.Parse("f6cb3f27-1028-4851-af96-887e570dc21f")
	resp, err := client.Client.Send(ctx, http.MethodDelete, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue CommentReaction
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the DeleteCommentReaction function
type DeleteCommentReactionArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) WorkItem ID
	WorkItemId *int
	// (required) Comment ID
	CommentId *int
	// (required) Type of the reaction
	ReactionType *CommentReactionType
}

// Deletes the field. To undelete a filed, see "Update Field" API.
func (client *ClientImpl) DeleteField(ctx context.Context, args DeleteFieldArgs) error {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FieldNameOrRefName == nil || *args.FieldNameOrRefName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FieldNameOrRefName"}
	}
	routeValues["fieldNameOrRefName"] = *args.FieldNameOrRefName

	locationId, _ := uuid.Parse("b51fd764-e5c2-4b9b-aaf7-3395cf4bdd94")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "6.0", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteField function
type DeleteFieldArgs struct {
	// (required) Field simple name or reference name
	FieldNameOrRefName *string
	// (optional) Project ID or project name
	Project *string
}

// Delete a query or a folder. This deletes any permission change on the deleted query or folder and any of its descendants if it is a folder. It is important to note that the deleted permission changes cannot be recovered upon undeleting the query or folder.
func (client *ClientImpl) DeleteQuery(ctx context.Context, args DeleteQueryArgs) error {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Query == nil || *args.Query == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Query"}
	}
	routeValues["query"] = *args.Query

	locationId, _ := uuid.Parse("a67d190c-c41f-424b-814d-0e906f659301")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "6.0", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteQuery function
type DeleteQueryArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) ID or path of the query or folder to delete.
	Query *string
}

// [Preview API]
func (client *ClientImpl) DeleteTag(ctx context.Context, args DeleteTagArgs) error {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.TagIdOrName == nil || *args.TagIdOrName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.TagIdOrName"}
	}
	routeValues["tagIdOrName"] = *args.TagIdOrName

	locationId, _ := uuid.Parse("bc15bc60-e7a8-43cb-ab01-2106be3983a1")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteTag function
type DeleteTagArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required)
	TagIdOrName *string
}

// [Preview API] Deletes the template with given id
func (client *ClientImpl) DeleteTemplate(ctx context.Context, args DeleteTemplateArgs) error {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team == nil || *args.Team == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Team"}
	}
	routeValues["team"] = *args.Team
	if args.TemplateId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.TemplateId"}
	}
	routeValues["templateId"] = (*args.TemplateId).String()

	locationId, _ := uuid.Parse("fb10264a-8836-48a0-8033-1b0ccd2748d5")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteTemplate function
type DeleteTemplateArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Team ID or team name
	Team *string
	// (required) Template id
	TemplateId *uuid.UUID
}

// Deletes the specified work item and sends it to the Recycle Bin, so that it can be restored back, if required. Optionally, if the destroy parameter has been set to true, it destroys the work item permanently. WARNING: If the destroy parameter is set to true, work items deleted by this command will NOT go to recycle-bin and there is no way to restore/recover them after deletion. It is recommended NOT to use this parameter. If you do, please use this parameter with extreme caution.
func (client *ClientImpl) DeleteWorkItem(ctx context.Context, args DeleteWorkItemArgs) (*WorkItemDelete, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)

	queryParams := url.Values{}
	if args.Destroy != nil {
		queryParams.Add("destroy", strconv.FormatBool(*args.Destroy))
	}
	locationId, _ := uuid.Parse("72c7ddf8-2cdc-4f60-90cd-ab71c14a399b")
	resp, err := client.Client.Send(ctx, http.MethodDelete, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemDelete
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the DeleteWorkItem function
type DeleteWorkItemArgs struct {
	// (required) ID of the work item to be deleted
	Id *int
	// (optional) Project ID or project name
	Project *string
	// (optional) Optional parameter, if set to true, the work item is deleted permanently. Please note: the destroy action is PERMANENT and cannot be undone.
	Destroy *bool
}

// Destroys the specified work item permanently from the Recycle Bin. This action can not be undone.
func (client *ClientImpl) DestroyWorkItem(ctx context.Context, args DestroyWorkItemArgs) error {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.Id == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)

	locationId, _ := uuid.Parse("b70d8d39-926c-465e-b927-b1bf0e5ca0e0")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "6.0", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DestroyWorkItem function
type DestroyWorkItemArgs struct {
	// (required) ID of the work item to be destroyed permanently
	Id *int
	// (optional) Project ID or project name
	Project *string
}

// Downloads an attachment.
func (client *ClientImpl) GetAttachmentContent(ctx context.Context, args GetAttachmentContentArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = (*args.Id).String()

	queryParams := url.Values{}
	if args.FileName != nil {
		queryParams.Add("fileName", *args.FileName)
	}
	if args.Download != nil {
		queryParams.Add("download", strconv.FormatBool(*args.Download))
	}
	locationId, _ := uuid.Parse("e07b5fa4-1499-494d-a496-64b860fd64ff")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/octet-stream", nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetAttachmentContent function
type GetAttachmentContentArgs struct {
	// (required) Attachment ID
	Id *uuid.UUID
	// (optional) Project ID or project name
	Project *string
	// (optional) Name of the file
	FileName *string
	// (optional) If set to <c>true</c> always download attachment
	Download *bool
}

// Downloads an attachment.
func (client *ClientImpl) GetAttachmentZip(ctx context.Context, args GetAttachmentZipArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = (*args.Id).String()

	queryParams := url.Values{}
	if args.FileName != nil {
		queryParams.Add("fileName", *args.FileName)
	}
	if args.Download != nil {
		queryParams.Add("download", strconv.FormatBool(*args.Download))
	}
	locationId, _ := uuid.Parse("e07b5fa4-1499-494d-a496-64b860fd64ff")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/zip", nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetAttachmentZip function
type GetAttachmentZipArgs struct {
	// (required) Attachment ID
	Id *uuid.UUID
	// (optional) Project ID or project name
	Project *string
	// (optional) Name of the file
	FileName *string
	// (optional) If set to <c>true</c> always download attachment
	Download *bool
}

// Gets the classification node for a given node path.
func (client *ClientImpl) GetClassificationNode(ctx context.Context, args GetClassificationNodeArgs) (*WorkItemClassificationNode, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.StructureGroup == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.StructureGroup"}
	}
	routeValues["structureGroup"] = string(*args.StructureGroup)
	if args.Path != nil && *args.Path != "" {
		routeValues["path"] = *args.Path
	}

	queryParams := url.Values{}
	if args.Depth != nil {
		queryParams.Add("$depth", strconv.Itoa(*args.Depth))
	}
	locationId, _ := uuid.Parse("5a172953-1b41-49d3-840a-33f79c3ce89f")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemClassificationNode
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetClassificationNode function
type GetClassificationNodeArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Structure group of the classification node, area or iteration.
	StructureGroup *TreeStructureGroup
	// (optional) Path of the classification node.
	Path *string
	// (optional) Depth of children to fetch.
	Depth *int
}

// Gets root classification nodes or list of classification nodes for a given list of nodes ids, for a given project. In case ids parameter is supplied you will  get list of classification nodes for those ids. Otherwise you will get root classification nodes for this project.
func (client *ClientImpl) GetClassificationNodes(ctx context.Context, args GetClassificationNodesArgs) (*[]WorkItemClassificationNode, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	queryParams := url.Values{}
	if args.Ids == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "ids"}
	}
	var stringList []string
	for _, item := range *args.Ids {
		stringList = append(stringList, strconv.Itoa(item))
	}
	listAsString := strings.Join((stringList)[:], ",")
	queryParams.Add("ids", listAsString)
	if args.Depth != nil {
		queryParams.Add("$depth", strconv.Itoa(*args.Depth))
	}
	if args.ErrorPolicy != nil {
		queryParams.Add("errorPolicy", string(*args.ErrorPolicy))
	}
	locationId, _ := uuid.Parse("a70579d1-f53a-48ee-a5be-7be8659023b9")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkItemClassificationNode
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetClassificationNodes function
type GetClassificationNodesArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Comma separated integer classification nodes ids. It's not required, if you want root nodes.
	Ids *[]int
	// (optional) Depth of children to fetch.
	Depth *int
	// (optional) Flag to handle errors in getting some nodes. Possible options are Fail and Omit.
	ErrorPolicy *ClassificationNodesErrorPolicy
}

// [Preview API] Returns a work item comment.
func (client *ClientImpl) GetComment(ctx context.Context, args GetCommentArgs) (*Comment, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WorkItemId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.WorkItemId"}
	}
	routeValues["workItemId"] = strconv.Itoa(*args.WorkItemId)
	if args.CommentId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.CommentId"}
	}
	routeValues["commentId"] = strconv.Itoa(*args.CommentId)

	queryParams := url.Values{}
	if args.IncludeDeleted != nil {
		queryParams.Add("includeDeleted", strconv.FormatBool(*args.IncludeDeleted))
	}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	locationId, _ := uuid.Parse("608aac0a-32e1-4493-a863-b9cf4566d257")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.3", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Comment
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetComment function
type GetCommentArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Id of a work item to get the comment.
	WorkItemId *int
	// (required) Id of the comment to return.
	CommentId *int
	// (optional) Specify if the deleted comment should be retrieved.
	IncludeDeleted *bool
	// (optional) Specifies the additional data retrieval options for work item comments.
	Expand *CommentExpandOptions
}

// [Preview API] Gets reactions of a comment.
func (client *ClientImpl) GetCommentReactions(ctx context.Context, args GetCommentReactionsArgs) (*[]CommentReaction, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WorkItemId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.WorkItemId"}
	}
	routeValues["workItemId"] = strconv.Itoa(*args.WorkItemId)
	if args.CommentId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.CommentId"}
	}
	routeValues["commentId"] = strconv.Itoa(*args.CommentId)

	locationId, _ := uuid.Parse("f6cb3f27-1028-4851-af96-887e570dc21f")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []CommentReaction
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetCommentReactions function
type GetCommentReactionsArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) WorkItem ID
	WorkItemId *int
	// (required) Comment ID
	CommentId *int
}

// [Preview API] Returns a list of work item comments, pageable.
func (client *ClientImpl) GetComments(ctx context.Context, args GetCommentsArgs) (*CommentList, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WorkItemId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.WorkItemId"}
	}
	routeValues["workItemId"] = strconv.Itoa(*args.WorkItemId)

	queryParams := url.Values{}
	if args.Top != nil {
		queryParams.Add("$top", strconv.Itoa(*args.Top))
	}
	if args.ContinuationToken != nil {
		queryParams.Add("continuationToken", *args.ContinuationToken)
	}
	if args.IncludeDeleted != nil {
		queryParams.Add("includeDeleted", strconv.FormatBool(*args.IncludeDeleted))
	}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	if args.Order != nil {
		queryParams.Add("order", string(*args.Order))
	}
	locationId, _ := uuid.Parse("608aac0a-32e1-4493-a863-b9cf4566d257")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.3", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue CommentList
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetComments function
type GetCommentsArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Id of a work item to get comments for.
	WorkItemId *int
	// (optional) Max number of comments to return.
	Top *int
	// (optional) Used to query for the next page of comments.
	ContinuationToken *string
	// (optional) Specify if the deleted comments should be retrieved.
	IncludeDeleted *bool
	// (optional) Specifies the additional data retrieval options for work item comments.
	Expand *CommentExpandOptions
	// (optional) Order in which the comments should be returned.
	Order *CommentSortOrder
}

// [Preview API] Returns a list of work item comments by ids.
func (client *ClientImpl) GetCommentsBatch(ctx context.Context, args GetCommentsBatchArgs) (*CommentList, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WorkItemId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.WorkItemId"}
	}
	routeValues["workItemId"] = strconv.Itoa(*args.WorkItemId)

	queryParams := url.Values{}
	if args.Ids == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "ids"}
	}
	var stringList []string
	for _, item := range *args.Ids {
		stringList = append(stringList, strconv.Itoa(item))
	}
	listAsString := strings.Join((stringList)[:], ",")
	queryParams.Add("ids", listAsString)
	if args.IncludeDeleted != nil {
		queryParams.Add("includeDeleted", strconv.FormatBool(*args.IncludeDeleted))
	}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	locationId, _ := uuid.Parse("608aac0a-32e1-4493-a863-b9cf4566d257")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.3", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue CommentList
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetCommentsBatch function
type GetCommentsBatchArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Id of a work item to get comments for.
	WorkItemId *int
	// (required) Comma-separated list of comment ids to return.
	Ids *[]int
	// (optional) Specify if the deleted comments should be retrieved.
	IncludeDeleted *bool
	// (optional) Specifies the additional data retrieval options for work item comments.
	Expand *CommentExpandOptions
}

// [Preview API]
func (client *ClientImpl) GetCommentVersion(ctx context.Context, args GetCommentVersionArgs) (*CommentVersion, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WorkItemId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.WorkItemId"}
	}
	routeValues["workItemId"] = strconv.Itoa(*args.WorkItemId)
	if args.CommentId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.CommentId"}
	}
	routeValues["commentId"] = strconv.Itoa(*args.CommentId)
	if args.Version == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Version"}
	}
	routeValues["version"] = strconv.Itoa(*args.Version)

	locationId, _ := uuid.Parse("49e03b34-3be0-42e3-8a5d-e8dfb88ac954")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue CommentVersion
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetCommentVersion function
type GetCommentVersionArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required)
	WorkItemId *int
	// (required)
	CommentId *int
	// (required)
	Version *int
}

// [Preview API]
func (client *ClientImpl) GetCommentVersions(ctx context.Context, args GetCommentVersionsArgs) (*[]CommentVersion, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WorkItemId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.WorkItemId"}
	}
	routeValues["workItemId"] = strconv.Itoa(*args.WorkItemId)
	if args.CommentId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.CommentId"}
	}
	routeValues["commentId"] = strconv.Itoa(*args.CommentId)

	locationId, _ := uuid.Parse("49e03b34-3be0-42e3-8a5d-e8dfb88ac954")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []CommentVersion
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetCommentVersions function
type GetCommentVersionsArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required)
	WorkItemId *int
	// (required)
	CommentId *int
}

// Gets a deleted work item from Recycle Bin.
func (client *ClientImpl) GetDeletedWorkItem(ctx context.Context, args GetDeletedWorkItemArgs) (*WorkItemDelete, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)

	locationId, _ := uuid.Parse("b70d8d39-926c-465e-b927-b1bf0e5ca0e0")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemDelete
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetDeletedWorkItem function
type GetDeletedWorkItemArgs struct {
	// (required) ID of the work item to be returned
	Id *int
	// (optional) Project ID or project name
	Project *string
}

// Gets the work items from the recycle bin, whose IDs have been specified in the parameters
func (client *ClientImpl) GetDeletedWorkItems(ctx context.Context, args GetDeletedWorkItemsArgs) (*[]WorkItemDeleteReference, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	queryParams := url.Values{}
	if args.Ids == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "ids"}
	}
	var stringList []string
	for _, item := range *args.Ids {
		stringList = append(stringList, strconv.Itoa(item))
	}
	listAsString := strings.Join((stringList)[:], ",")
	queryParams.Add("ids", listAsString)
	locationId, _ := uuid.Parse("b70d8d39-926c-465e-b927-b1bf0e5ca0e0")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkItemDeleteReference
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetDeletedWorkItems function
type GetDeletedWorkItemsArgs struct {
	// (required) Comma separated list of IDs of the deleted work items to be returned
	Ids *[]int
	// (optional) Project ID or project name
	Project *string
}

// Gets a list of the IDs and the URLs of the deleted the work items in the Recycle Bin.
func (client *ClientImpl) GetDeletedWorkItemShallowReferences(ctx context.Context, args GetDeletedWorkItemShallowReferencesArgs) (*[]WorkItemDeleteShallowReference, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	locationId, _ := uuid.Parse("b70d8d39-926c-465e-b927-b1bf0e5ca0e0")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkItemDeleteShallowReference
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetDeletedWorkItemShallowReferences function
type GetDeletedWorkItemShallowReferencesArgs struct {
	// (optional) Project ID or project name
	Project *string
}

// [Preview API] Get users who reacted on the comment.
func (client *ClientImpl) GetEngagedUsers(ctx context.Context, args GetEngagedUsersArgs) (*[]webapi.IdentityRef, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WorkItemId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.WorkItemId"}
	}
	routeValues["workItemId"] = strconv.Itoa(*args.WorkItemId)
	if args.CommentId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.CommentId"}
	}
	routeValues["commentId"] = strconv.Itoa(*args.CommentId)
	if args.ReactionType == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ReactionType"}
	}
	routeValues["reactionType"] = string(*args.ReactionType)

	queryParams := url.Values{}
	if args.Top != nil {
		queryParams.Add("$top", strconv.Itoa(*args.Top))
	}
	if args.Skip != nil {
		queryParams.Add("$skip", strconv.Itoa(*args.Skip))
	}
	locationId, _ := uuid.Parse("e33ca5e0-2349-4285-af3d-d72d86781c35")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []webapi.IdentityRef
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetEngagedUsers function
type GetEngagedUsersArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) WorkItem ID.
	WorkItemId *int
	// (required) Comment ID.
	CommentId *int
	// (required) Type of the reaction.
	ReactionType *CommentReactionType
	// (optional)
	Top *int
	// (optional)
	Skip *int
}

// Gets information on a specific field.
func (client *ClientImpl) GetField(ctx context.Context, args GetFieldArgs) (*WorkItemField, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FieldNameOrRefName == nil || *args.FieldNameOrRefName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FieldNameOrRefName"}
	}
	routeValues["fieldNameOrRefName"] = *args.FieldNameOrRefName

	locationId, _ := uuid.Parse("b51fd764-e5c2-4b9b-aaf7-3395cf4bdd94")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemField
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetField function
type GetFieldArgs struct {
	// (required) Field simple name or reference name
	FieldNameOrRefName *string
	// (optional) Project ID or project name
	Project *string
}

// Returns information for all fields. The project ID/name parameter is optional.
func (client *ClientImpl) GetFields(ctx context.Context, args GetFieldsArgs) (*[]WorkItemField, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	queryParams := url.Values{}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	locationId, _ := uuid.Parse("b51fd764-e5c2-4b9b-aaf7-3395cf4bdd94")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkItemField
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetFields function
type GetFieldsArgs struct {
	// (optional) Project ID or project name
	Project *string
	// (optional) Use ExtensionFields to include extension fields, otherwise exclude them. Unless the feature flag for this parameter is enabled, extension fields are always included.
	Expand *GetFieldsExpand
}

// Gets the root queries and their children
func (client *ClientImpl) GetQueries(ctx context.Context, args GetQueriesArgs) (*[]QueryHierarchyItem, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	queryParams := url.Values{}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	if args.Depth != nil {
		queryParams.Add("$depth", strconv.Itoa(*args.Depth))
	}
	if args.IncludeDeleted != nil {
		queryParams.Add("$includeDeleted", strconv.FormatBool(*args.IncludeDeleted))
	}
	locationId, _ := uuid.Parse("a67d190c-c41f-424b-814d-0e906f659301")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []QueryHierarchyItem
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetQueries function
type GetQueriesArgs struct {
	// (required) Project ID or project name
	Project *string
	// (optional) Include the query string (wiql), clauses, query result columns, and sort options in the results.
	Expand *QueryExpand
	// (optional) In the folder of queries, return child queries and folders to this depth.
	Depth *int
	// (optional) Include deleted queries and folders
	IncludeDeleted *bool
}

// Gets a list of queries by ids (Maximum 1000)
func (client *ClientImpl) GetQueriesBatch(ctx context.Context, args GetQueriesBatchArgs) (*[]QueryHierarchyItem, error) {
	if args.QueryGetRequest == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.QueryGetRequest"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	body, marshalErr := json.Marshal(*args.QueryGetRequest)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("549816f9-09b0-4e75-9e81-01fbfcd07426")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []QueryHierarchyItem
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetQueriesBatch function
type GetQueriesBatchArgs struct {
	// (required)
	QueryGetRequest *QueryBatchGetRequest
	// (required) Project ID or project name
	Project *string
}

// Retrieves an individual query and its children
func (client *ClientImpl) GetQuery(ctx context.Context, args GetQueryArgs) (*QueryHierarchyItem, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Query == nil || *args.Query == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Query"}
	}
	routeValues["query"] = *args.Query

	queryParams := url.Values{}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	if args.Depth != nil {
		queryParams.Add("$depth", strconv.Itoa(*args.Depth))
	}
	if args.IncludeDeleted != nil {
		queryParams.Add("$includeDeleted", strconv.FormatBool(*args.IncludeDeleted))
	}
	locationId, _ := uuid.Parse("a67d190c-c41f-424b-814d-0e906f659301")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue QueryHierarchyItem
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetQuery function
type GetQueryArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) ID or path of the query.
	Query *string
	// (optional) Include the query string (wiql), clauses, query result columns, and sort options in the results.
	Expand *QueryExpand
	// (optional) In the folder of queries, return child queries and folders to this depth.
	Depth *int
	// (optional) Include deleted queries and folders
	IncludeDeleted *bool
}

// Gets the results of the query given the query ID.
func (client *ClientImpl) GetQueryResultCount(ctx context.Context, args GetQueryResultCountArgs) (*int, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.Team != nil && *args.Team != "" {
		routeValues["team"] = *args.Team
	}
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = (*args.Id).String()

	queryParams := url.Values{}
	if args.TimePrecision != nil {
		queryParams.Add("timePrecision", strconv.FormatBool(*args.TimePrecision))
	}
	if args.Top != nil {
		queryParams.Add("$top", strconv.Itoa(*args.Top))
	}
	locationId, _ := uuid.Parse("a02355f5-5f8a-4671-8e32-369d23aac83d")
	resp, err := client.Client.Send(ctx, http.MethodHead, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue int
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetQueryResultCount function
type GetQueryResultCountArgs struct {
	// (required) The query ID.
	Id *uuid.UUID
	// (optional) Project ID or project name
	Project *string
	// (optional) Team ID or team name
	Team *string
	// (optional) Whether or not to use time precision.
	TimePrecision *bool
	// (optional) The max number of results to return.
	Top *int
}

// [Preview API] Gets recent work item activities
func (client *ClientImpl) GetRecentActivityData(ctx context.Context, args GetRecentActivityDataArgs) (*[]AccountRecentActivityWorkItemModel2, error) {
	locationId, _ := uuid.Parse("1bc988f4-c15f-4072-ad35-497c87e3a909")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.2", nil, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []AccountRecentActivityWorkItemModel2
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetRecentActivityData function
type GetRecentActivityDataArgs struct {
}

// Gets the work item relation type definition.
func (client *ClientImpl) GetRelationType(ctx context.Context, args GetRelationTypeArgs) (*WorkItemRelationType, error) {
	routeValues := make(map[string]string)
	if args.Relation == nil || *args.Relation == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Relation"}
	}
	routeValues["relation"] = *args.Relation

	locationId, _ := uuid.Parse("f5d33bc9-5b49-4a3c-a9bd-f3cd46dd2165")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemRelationType
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetRelationType function
type GetRelationTypeArgs struct {
	// (required) The relation name
	Relation *string
}

// Gets the work item relation types.
func (client *ClientImpl) GetRelationTypes(ctx context.Context, args GetRelationTypesArgs) (*[]WorkItemRelationType, error) {
	locationId, _ := uuid.Parse("f5d33bc9-5b49-4a3c-a9bd-f3cd46dd2165")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", nil, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkItemRelationType
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetRelationTypes function
type GetRelationTypesArgs struct {
}

// Get a batch of work item links
func (client *ClientImpl) GetReportingLinksByLinkType(ctx context.Context, args GetReportingLinksByLinkTypeArgs) (*ReportingWorkItemLinksBatch, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	queryParams := url.Values{}
	if args.LinkTypes != nil {
		listAsString := strings.Join((*args.LinkTypes)[:], ",")
		queryParams.Add("linkTypes", listAsString)
	}
	if args.Types != nil {
		listAsString := strings.Join((*args.Types)[:], ",")
		queryParams.Add("types", listAsString)
	}
	if args.ContinuationToken != nil {
		queryParams.Add("continuationToken", *args.ContinuationToken)
	}
	if args.StartDateTime != nil {
		queryParams.Add("startDateTime", (*args.StartDateTime).String())
	}
	locationId, _ := uuid.Parse("b5b5b6d0-0308-40a1-b3f4-b9bb3c66878f")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ReportingWorkItemLinksBatch
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetReportingLinksByLinkType function
type GetReportingLinksByLinkTypeArgs struct {
	// (optional) Project ID or project name
	Project *string
	// (optional) A list of types to filter the results to specific link types. Omit this parameter to get work item links of all link types.
	LinkTypes *[]string
	// (optional) A list of types to filter the results to specific work item types. Omit this parameter to get work item links of all work item types.
	Types *[]string
	// (optional) Specifies the continuationToken to start the batch from. Omit this parameter to get the first batch of links.
	ContinuationToken *string
	// (optional) Date/time to use as a starting point for link changes. Only link changes that occurred after that date/time will be returned. Cannot be used in conjunction with 'watermark' parameter.
	StartDateTime *azuredevops.Time
}

// Returns a fully hydrated work item for the requested revision
func (client *ClientImpl) GetRevision(ctx context.Context, args GetRevisionArgs) (*WorkItem, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)
	if args.RevisionNumber == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.RevisionNumber"}
	}
	routeValues["revisionNumber"] = strconv.Itoa(*args.RevisionNumber)

	queryParams := url.Values{}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	locationId, _ := uuid.Parse("a00c85a5-80fa-4565-99c3-bcd2181434bb")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItem
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetRevision function
type GetRevisionArgs struct {
	// (required)
	Id *int
	// (required)
	RevisionNumber *int
	// (optional) Project ID or project name
	Project *string
	// (optional)
	Expand *WorkItemExpand
}

// Returns the list of fully hydrated work item revisions, paged.
func (client *ClientImpl) GetRevisions(ctx context.Context, args GetRevisionsArgs) (*[]WorkItem, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)

	queryParams := url.Values{}
	if args.Top != nil {
		queryParams.Add("$top", strconv.Itoa(*args.Top))
	}
	if args.Skip != nil {
		queryParams.Add("$skip", strconv.Itoa(*args.Skip))
	}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	locationId, _ := uuid.Parse("a00c85a5-80fa-4565-99c3-bcd2181434bb")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkItem
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetRevisions function
type GetRevisionsArgs struct {
	// (required)
	Id *int
	// (optional) Project ID or project name
	Project *string
	// (optional)
	Top *int
	// (optional)
	Skip *int
	// (optional)
	Expand *WorkItemExpand
}

// Gets root classification nodes under the project.
func (client *ClientImpl) GetRootNodes(ctx context.Context, args GetRootNodesArgs) (*[]WorkItemClassificationNode, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	queryParams := url.Values{}
	if args.Depth != nil {
		queryParams.Add("$depth", strconv.Itoa(*args.Depth))
	}
	locationId, _ := uuid.Parse("a70579d1-f53a-48ee-a5be-7be8659023b9")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkItemClassificationNode
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetRootNodes function
type GetRootNodesArgs struct {
	// (required) Project ID or project name
	Project *string
	// (optional) Depth of children to fetch.
	Depth *int
}

// [Preview API]
func (client *ClientImpl) GetTag(ctx context.Context, args GetTagArgs) (*WorkItemTagDefinition, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.TagIdOrName == nil || *args.TagIdOrName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.TagIdOrName"}
	}
	routeValues["tagIdOrName"] = *args.TagIdOrName

	locationId, _ := uuid.Parse("bc15bc60-e7a8-43cb-ab01-2106be3983a1")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemTagDefinition
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetTag function
type GetTagArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required)
	TagIdOrName *string
}

// [Preview API]
func (client *ClientImpl) GetTags(ctx context.Context, args GetTagsArgs) (*[]WorkItemTagDefinition, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	locationId, _ := uuid.Parse("bc15bc60-e7a8-43cb-ab01-2106be3983a1")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkItemTagDefinition
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetTags function
type GetTagsArgs struct {
	// (required) Project ID or project name
	Project *string
}

// [Preview API] Gets the template with specified id
func (client *ClientImpl) GetTemplate(ctx context.Context, args GetTemplateArgs) (*WorkItemTemplate, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team == nil || *args.Team == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Team"}
	}
	routeValues["team"] = *args.Team
	if args.TemplateId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.TemplateId"}
	}
	routeValues["templateId"] = (*args.TemplateId).String()

	locationId, _ := uuid.Parse("fb10264a-8836-48a0-8033-1b0ccd2748d5")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemTemplate
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetTemplate function
type GetTemplateArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Team ID or team name
	Team *string
	// (required) Template Id
	TemplateId *uuid.UUID
}

// [Preview API] Gets template
func (client *ClientImpl) GetTemplates(ctx context.Context, args GetTemplatesArgs) (*[]WorkItemTemplateReference, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team == nil || *args.Team == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Team"}
	}
	routeValues["team"] = *args.Team

	queryParams := url.Values{}
	if args.Workitemtypename != nil {
		queryParams.Add("workitemtypename", *args.Workitemtypename)
	}
	locationId, _ := uuid.Parse("6a90345f-a676-4969-afce-8e163e1d5642")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkItemTemplateReference
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetTemplates function
type GetTemplatesArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Team ID or team name
	Team *string
	// (optional) Optional, When specified returns templates for given Work item type.
	Workitemtypename *string
}

// Returns a single update for a work item
func (client *ClientImpl) GetUpdate(ctx context.Context, args GetUpdateArgs) (*WorkItemUpdate, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)
	if args.UpdateNumber == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UpdateNumber"}
	}
	routeValues["updateNumber"] = strconv.Itoa(*args.UpdateNumber)

	locationId, _ := uuid.Parse("6570bf97-d02c-4a91-8d93-3abe9895b1a9")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemUpdate
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetUpdate function
type GetUpdateArgs struct {
	// (required)
	Id *int
	// (required)
	UpdateNumber *int
	// (optional) Project ID or project name
	Project *string
}

// Returns a the deltas between work item revisions
func (client *ClientImpl) GetUpdates(ctx context.Context, args GetUpdatesArgs) (*[]WorkItemUpdate, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)

	queryParams := url.Values{}
	if args.Top != nil {
		queryParams.Add("$top", strconv.Itoa(*args.Top))
	}
	if args.Skip != nil {
		queryParams.Add("$skip", strconv.Itoa(*args.Skip))
	}
	locationId, _ := uuid.Parse("6570bf97-d02c-4a91-8d93-3abe9895b1a9")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkItemUpdate
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetUpdates function
type GetUpdatesArgs struct {
	// (required)
	Id *int
	// (optional) Project ID or project name
	Project *string
	// (optional)
	Top *int
	// (optional)
	Skip *int
}

// [Preview API] Get the list of work item tracking outbound artifact link types.
func (client *ClientImpl) GetWorkArtifactLinkTypes(ctx context.Context, args GetWorkArtifactLinkTypesArgs) (*[]WorkArtifactLink, error) {
	locationId, _ := uuid.Parse("1a31de40-e318-41cd-a6c6-881077df52e3")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", nil, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkArtifactLink
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWorkArtifactLinkTypes function
type GetWorkArtifactLinkTypesArgs struct {
}

// Returns a single work item.
func (client *ClientImpl) GetWorkItem(ctx context.Context, args GetWorkItemArgs) (*WorkItem, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)

	queryParams := url.Values{}
	if args.Fields != nil {
		listAsString := strings.Join((*args.Fields)[:], ",")
		queryParams.Add("fields", listAsString)
	}
	if args.AsOf != nil {
		queryParams.Add("asOf", (*args.AsOf).String())
	}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	locationId, _ := uuid.Parse("72c7ddf8-2cdc-4f60-90cd-ab71c14a399b")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItem
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWorkItem function
type GetWorkItemArgs struct {
	// (required) The work item id
	Id *int
	// (optional) Project ID or project name
	Project *string
	// (optional) Comma-separated list of requested fields
	Fields *[]string
	// (optional) AsOf UTC date time string
	AsOf *azuredevops.Time
	// (optional) The expand parameters for work item attributes. Possible options are { None, Relations, Fields, Links, All }.
	Expand *WorkItemExpand
}

// Get a work item icon given the friendly name and icon color.
func (client *ClientImpl) GetWorkItemIconJson(ctx context.Context, args GetWorkItemIconJsonArgs) (*WorkItemIcon, error) {
	routeValues := make(map[string]string)
	if args.Icon == nil || *args.Icon == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Icon"}
	}
	routeValues["icon"] = *args.Icon

	queryParams := url.Values{}
	if args.Color != nil {
		queryParams.Add("color", *args.Color)
	}
	if args.V != nil {
		queryParams.Add("v", strconv.Itoa(*args.V))
	}
	locationId, _ := uuid.Parse("4e1eb4a5-1970-4228-a682-ec48eb2dca30")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemIcon
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWorkItemIconJson function
type GetWorkItemIconJsonArgs struct {
	// (required) The name of the icon
	Icon *string
	// (optional) The 6-digit hex color for the icon
	Color *string
	// (optional) The version of the icon (used only for cache invalidation)
	V *int
}

// Get a list of all work item icons.
func (client *ClientImpl) GetWorkItemIcons(ctx context.Context, args GetWorkItemIconsArgs) (*[]WorkItemIcon, error) {
	locationId, _ := uuid.Parse("4e1eb4a5-1970-4228-a682-ec48eb2dca30")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", nil, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkItemIcon
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWorkItemIcons function
type GetWorkItemIconsArgs struct {
}

// Get a work item icon given the friendly name and icon color.
func (client *ClientImpl) GetWorkItemIconSvg(ctx context.Context, args GetWorkItemIconSvgArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.Icon == nil || *args.Icon == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Icon"}
	}
	routeValues["icon"] = *args.Icon

	queryParams := url.Values{}
	if args.Color != nil {
		queryParams.Add("color", *args.Color)
	}
	if args.V != nil {
		queryParams.Add("v", strconv.Itoa(*args.V))
	}
	locationId, _ := uuid.Parse("4e1eb4a5-1970-4228-a682-ec48eb2dca30")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "image/svg+xml", nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetWorkItemIconSvg function
type GetWorkItemIconSvgArgs struct {
	// (required) The name of the icon
	Icon *string
	// (optional) The 6-digit hex color for the icon
	Color *string
	// (optional) The version of the icon (used only for cache invalidation)
	V *int
}

// Get a work item icon given the friendly name and icon color.
func (client *ClientImpl) GetWorkItemIconXaml(ctx context.Context, args GetWorkItemIconXamlArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.Icon == nil || *args.Icon == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Icon"}
	}
	routeValues["icon"] = *args.Icon

	queryParams := url.Values{}
	if args.Color != nil {
		queryParams.Add("color", *args.Color)
	}
	if args.V != nil {
		queryParams.Add("v", strconv.Itoa(*args.V))
	}
	locationId, _ := uuid.Parse("4e1eb4a5-1970-4228-a682-ec48eb2dca30")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "image/xaml+xml", nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetWorkItemIconXaml function
type GetWorkItemIconXamlArgs struct {
	// (required) The name of the icon
	Icon *string
	// (optional) The 6-digit hex color for the icon
	Color *string
	// (optional) The version of the icon (used only for cache invalidation)
	V *int
}

// [Preview API] Returns the next state on the given work item IDs.
func (client *ClientImpl) GetWorkItemNextStatesOnCheckinAction(ctx context.Context, args GetWorkItemNextStatesOnCheckinActionArgs) (*[]WorkItemNextStateOnTransition, error) {
	queryParams := url.Values{}
	if args.Ids == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "ids"}
	}
	var stringList []string
	for _, item := range *args.Ids {
		stringList = append(stringList, strconv.Itoa(item))
	}
	listAsString := strings.Join((stringList)[:], ",")
	queryParams.Add("ids", listAsString)
	if args.Action != nil {
		queryParams.Add("action", *args.Action)
	}
	locationId, _ := uuid.Parse("afae844b-e2f6-44c2-8053-17b3bb936a40")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", nil, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkItemNextStateOnTransition
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWorkItemNextStatesOnCheckinAction function
type GetWorkItemNextStatesOnCheckinActionArgs struct {
	// (required) list of work item ids
	Ids *[]int
	// (optional) possible actions. Currently only supports checkin
	Action *string
}

// Returns a list of work items (Maximum 200)
func (client *ClientImpl) GetWorkItems(ctx context.Context, args GetWorkItemsArgs) (*[]WorkItem, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	queryParams := url.Values{}
	if args.Ids == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "ids"}
	}
	var stringList []string
	for _, item := range *args.Ids {
		stringList = append(stringList, strconv.Itoa(item))
	}
	listAsString := strings.Join((stringList)[:], ",")
	queryParams.Add("ids", listAsString)
	if args.Fields != nil {
		listAsString := strings.Join((*args.Fields)[:], ",")
		queryParams.Add("fields", listAsString)
	}
	if args.AsOf != nil {
		queryParams.Add("asOf", (*args.AsOf).String())
	}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	if args.ErrorPolicy != nil {
		queryParams.Add("errorPolicy", string(*args.ErrorPolicy))
	}
	locationId, _ := uuid.Parse("72c7ddf8-2cdc-4f60-90cd-ab71c14a399b")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkItem
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWorkItems function
type GetWorkItemsArgs struct {
	// (required) The comma-separated list of requested work item ids. (Maximum 200 ids allowed).
	Ids *[]int
	// (optional) Project ID or project name
	Project *string
	// (optional) Comma-separated list of requested fields
	Fields *[]string
	// (optional) AsOf UTC date time string
	AsOf *azuredevops.Time
	// (optional) The expand parameters for work item attributes. Possible options are { None, Relations, Fields, Links, All }.
	Expand *WorkItemExpand
	// (optional) The flag to control error policy in a bulk get work items request. Possible options are {Fail, Omit}.
	ErrorPolicy *WorkItemErrorPolicy
}

// Gets work items for a list of work item ids (Maximum 200)
func (client *ClientImpl) GetWorkItemsBatch(ctx context.Context, args GetWorkItemsBatchArgs) (*[]WorkItem, error) {
	if args.WorkItemGetRequest == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.WorkItemGetRequest"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	body, marshalErr := json.Marshal(*args.WorkItemGetRequest)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("908509b6-4248-4475-a1cd-829139ba419f")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkItem
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWorkItemsBatch function
type GetWorkItemsBatchArgs struct {
	// (required)
	WorkItemGetRequest *WorkItemBatchGetRequest
	// (optional) Project ID or project name
	Project *string
}

// Returns a single work item from a template.
func (client *ClientImpl) GetWorkItemTemplate(ctx context.Context, args GetWorkItemTemplateArgs) (*WorkItem, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Type == nil || *args.Type == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Type"}
	}
	routeValues["type"] = *args.Type

	queryParams := url.Values{}
	if args.Fields != nil {
		queryParams.Add("fields", *args.Fields)
	}
	if args.AsOf != nil {
		queryParams.Add("asOf", (*args.AsOf).String())
	}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	locationId, _ := uuid.Parse("62d3d110-0047-428c-ad3c-4fe872c91c74")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItem
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWorkItemTemplate function
type GetWorkItemTemplateArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) The work item type name
	Type *string
	// (optional) Comma-separated list of requested fields
	Fields *string
	// (optional) AsOf UTC date time string
	AsOf *azuredevops.Time
	// (optional) The expand parameters for work item attributes. Possible options are { None, Relations, Fields, Links, All }.
	Expand *WorkItemExpand
}

// Returns a work item type definition.
func (client *ClientImpl) GetWorkItemType(ctx context.Context, args GetWorkItemTypeArgs) (*WorkItemType, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Type == nil || *args.Type == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Type"}
	}
	routeValues["type"] = *args.Type

	locationId, _ := uuid.Parse("7c8d7a76-4a09-43e8-b5df-bd792f4ac6aa")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemType
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWorkItemType function
type GetWorkItemTypeArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Work item type name
	Type *string
}

// Get all work item type categories.
func (client *ClientImpl) GetWorkItemTypeCategories(ctx context.Context, args GetWorkItemTypeCategoriesArgs) (*[]WorkItemTypeCategory, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	locationId, _ := uuid.Parse("9b9f5734-36c8-415e-ba67-f83b45c31408")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkItemTypeCategory
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWorkItemTypeCategories function
type GetWorkItemTypeCategoriesArgs struct {
	// (required) Project ID or project name
	Project *string
}

// Get specific work item type category by name.
func (client *ClientImpl) GetWorkItemTypeCategory(ctx context.Context, args GetWorkItemTypeCategoryArgs) (*WorkItemTypeCategory, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Category == nil || *args.Category == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Category"}
	}
	routeValues["category"] = *args.Category

	locationId, _ := uuid.Parse("9b9f5734-36c8-415e-ba67-f83b45c31408")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemTypeCategory
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWorkItemTypeCategory function
type GetWorkItemTypeCategoryArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) The category name
	Category *string
}

// Get a list of fields for a work item type with detailed references.
func (client *ClientImpl) GetWorkItemTypeFieldsWithReferences(ctx context.Context, args GetWorkItemTypeFieldsWithReferencesArgs) (*[]WorkItemTypeFieldWithReferences, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Type == nil || *args.Type == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Type"}
	}
	routeValues["type"] = *args.Type

	queryParams := url.Values{}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	locationId, _ := uuid.Parse("bd293ce5-3d25-4192-8e67-e8092e879efb")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkItemTypeFieldWithReferences
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWorkItemTypeFieldsWithReferences function
type GetWorkItemTypeFieldsWithReferencesArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Work item type.
	Type *string
	// (optional) Expand level for the API response. Properties: to include allowedvalues, default value, isRequired etc. as a part of response; None: to skip these properties.
	Expand *WorkItemTypeFieldsExpandLevel
}

// Get a field for a work item type with detailed references.
func (client *ClientImpl) GetWorkItemTypeFieldWithReferences(ctx context.Context, args GetWorkItemTypeFieldWithReferencesArgs) (*WorkItemTypeFieldWithReferences, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Type == nil || *args.Type == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Type"}
	}
	routeValues["type"] = *args.Type
	if args.Field == nil || *args.Field == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Field"}
	}
	routeValues["field"] = *args.Field

	queryParams := url.Values{}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	locationId, _ := uuid.Parse("bd293ce5-3d25-4192-8e67-e8092e879efb")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemTypeFieldWithReferences
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWorkItemTypeFieldWithReferences function
type GetWorkItemTypeFieldWithReferencesArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) Work item type.
	Type *string
	// (required)
	Field *string
	// (optional) Expand level for the API response. Properties: to include allowedvalues, default value, isRequired etc. as a part of response; None: to skip these properties.
	Expand *WorkItemTypeFieldsExpandLevel
}

// Returns the list of work item types
func (client *ClientImpl) GetWorkItemTypes(ctx context.Context, args GetWorkItemTypesArgs) (*[]WorkItemType, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	locationId, _ := uuid.Parse("7c8d7a76-4a09-43e8-b5df-bd792f4ac6aa")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkItemType
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWorkItemTypes function
type GetWorkItemTypesArgs struct {
	// (required) Project ID or project name
	Project *string
}

// [Preview API] Returns the state names and colors for a work item type.
func (client *ClientImpl) GetWorkItemTypeStates(ctx context.Context, args GetWorkItemTypeStatesArgs) (*[]WorkItemStateColor, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Type == nil || *args.Type == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Type"}
	}
	routeValues["type"] = *args.Type

	locationId, _ := uuid.Parse("7c9d7a76-4a09-43e8-b5df-bd792f4ac6aa")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []WorkItemStateColor
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetWorkItemTypeStates function
type GetWorkItemTypeStatesArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) The state name
	Type *string
}

// [Preview API] Migrates a project to a different process within the same OOB type. For example, you can only migrate a project from agile/custom-agile to agile/custom-agile.
func (client *ClientImpl) MigrateProjectsProcess(ctx context.Context, args MigrateProjectsProcessArgs) (*ProcessMigrationResultModel, error) {
	if args.NewProcess == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.NewProcess"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	body, marshalErr := json.Marshal(*args.NewProcess)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("19801631-d4e5-47e9-8166-0330de0ff1e6")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ProcessMigrationResultModel
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the MigrateProjectsProcess function
type MigrateProjectsProcessArgs struct {
	// (required)
	NewProcess *ProcessIdModel
	// (required) Project ID or project name
	Project *string
}

// Gets the results of the query given the query ID.
func (client *ClientImpl) QueryById(ctx context.Context, args QueryByIdArgs) (*WorkItemQueryResult, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.Team != nil && *args.Team != "" {
		routeValues["team"] = *args.Team
	}
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = (*args.Id).String()

	queryParams := url.Values{}
	if args.TimePrecision != nil {
		queryParams.Add("timePrecision", strconv.FormatBool(*args.TimePrecision))
	}
	if args.Top != nil {
		queryParams.Add("$top", strconv.Itoa(*args.Top))
	}
	locationId, _ := uuid.Parse("a02355f5-5f8a-4671-8e32-369d23aac83d")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemQueryResult
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryById function
type QueryByIdArgs struct {
	// (required) The query ID.
	Id *uuid.UUID
	// (optional) Project ID or project name
	Project *string
	// (optional) Team ID or team name
	Team *string
	// (optional) Whether or not to use time precision.
	TimePrecision *bool
	// (optional) The max number of results to return.
	Top *int
}

// Gets the results of the query given its WIQL.
func (client *ClientImpl) QueryByWiql(ctx context.Context, args QueryByWiqlArgs) (*WorkItemQueryResult, error) {
	if args.Wiql == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Wiql"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.Team != nil && *args.Team != "" {
		routeValues["team"] = *args.Team
	}

	queryParams := url.Values{}
	if args.TimePrecision != nil {
		queryParams.Add("timePrecision", strconv.FormatBool(*args.TimePrecision))
	}
	if args.Top != nil {
		queryParams.Add("$top", strconv.Itoa(*args.Top))
	}
	body, marshalErr := json.Marshal(*args.Wiql)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("1a9c53f7-f243-4447-b110-35ef023636e4")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0", routeValues, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemQueryResult
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryByWiql function
type QueryByWiqlArgs struct {
	// (required) The query containing the WIQL.
	Wiql *Wiql
	// (optional) Project ID or project name
	Project *string
	// (optional) Team ID or team name
	Team *string
	// (optional) Whether or not to use time precision.
	TimePrecision *bool
	// (optional) The max number of results to return.
	Top *int
}

// [Preview API] Queries work items linked to a given list of artifact URI.
func (client *ClientImpl) QueryWorkItemsForArtifactUris(ctx context.Context, args QueryWorkItemsForArtifactUrisArgs) (*ArtifactUriQueryResult, error) {
	if args.ArtifactUriQuery == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ArtifactUriQuery"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	body, marshalErr := json.Marshal(*args.ArtifactUriQuery)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("a9a9aa7a-8c09-44d3-ad1b-46e855c1e3d3")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ArtifactUriQueryResult
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryWorkItemsForArtifactUris function
type QueryWorkItemsForArtifactUrisArgs struct {
	// (required) Defines a list of artifact URI for querying work items.
	ArtifactUriQuery *ArtifactUriQuery
	// (optional) Project ID or project name
	Project *string
}

// [Preview API]
func (client *ClientImpl) ReadReportingDiscussions(ctx context.Context, args ReadReportingDiscussionsArgs) (*ReportingWorkItemRevisionsBatch, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	queryParams := url.Values{}
	if args.ContinuationToken != nil {
		queryParams.Add("continuationToken", *args.ContinuationToken)
	}
	if args.MaxPageSize != nil {
		queryParams.Add("$maxPageSize", strconv.Itoa(*args.MaxPageSize))
	}
	locationId, _ := uuid.Parse("4a644469-90c5-4fcc-9a9f-be0827d369ec")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ReportingWorkItemRevisionsBatch
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ReadReportingDiscussions function
type ReadReportingDiscussionsArgs struct {
	// (optional) Project ID or project name
	Project *string
	// (optional)
	ContinuationToken *string
	// (optional)
	MaxPageSize *int
}

// Get a batch of work item revisions with the option of including deleted items
func (client *ClientImpl) ReadReportingRevisionsGet(ctx context.Context, args ReadReportingRevisionsGetArgs) (*ReportingWorkItemRevisionsBatch, error) {
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	queryParams := url.Values{}
	if args.Fields != nil {
		listAsString := strings.Join((*args.Fields)[:], ",")
		queryParams.Add("fields", listAsString)
	}
	if args.Types != nil {
		listAsString := strings.Join((*args.Types)[:], ",")
		queryParams.Add("types", listAsString)
	}
	if args.ContinuationToken != nil {
		queryParams.Add("continuationToken", *args.ContinuationToken)
	}
	if args.StartDateTime != nil {
		queryParams.Add("startDateTime", (*args.StartDateTime).String())
	}
	if args.IncludeIdentityRef != nil {
		queryParams.Add("includeIdentityRef", strconv.FormatBool(*args.IncludeIdentityRef))
	}
	if args.IncludeDeleted != nil {
		queryParams.Add("includeDeleted", strconv.FormatBool(*args.IncludeDeleted))
	}
	if args.IncludeTagRef != nil {
		queryParams.Add("includeTagRef", strconv.FormatBool(*args.IncludeTagRef))
	}
	if args.IncludeLatestOnly != nil {
		queryParams.Add("includeLatestOnly", strconv.FormatBool(*args.IncludeLatestOnly))
	}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	if args.IncludeDiscussionChangesOnly != nil {
		queryParams.Add("includeDiscussionChangesOnly", strconv.FormatBool(*args.IncludeDiscussionChangesOnly))
	}
	if args.MaxPageSize != nil {
		queryParams.Add("$maxPageSize", strconv.Itoa(*args.MaxPageSize))
	}
	locationId, _ := uuid.Parse("f828fe59-dd87-495d-a17c-7a8d6211ca6c")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ReportingWorkItemRevisionsBatch
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ReadReportingRevisionsGet function
type ReadReportingRevisionsGetArgs struct {
	// (optional) Project ID or project name
	Project *string
	// (optional) A list of fields to return in work item revisions. Omit this parameter to get all reportable fields.
	Fields *[]string
	// (optional) A list of types to filter the results to specific work item types. Omit this parameter to get work item revisions of all work item types.
	Types *[]string
	// (optional) Specifies the watermark to start the batch from. Omit this parameter to get the first batch of revisions.
	ContinuationToken *string
	// (optional) Date/time to use as a starting point for revisions, all revisions will occur after this date/time. Cannot be used in conjunction with 'watermark' parameter.
	StartDateTime *azuredevops.Time
	// (optional) Return an identity reference instead of a string value for identity fields.
	IncludeIdentityRef *bool
	// (optional) Specify if the deleted item should be returned.
	IncludeDeleted *bool
	// (optional) Specify if the tag objects should be returned for System.Tags field.
	IncludeTagRef *bool
	// (optional) Return only the latest revisions of work items, skipping all historical revisions
	IncludeLatestOnly *bool
	// (optional) Return all the fields in work item revisions, including long text fields which are not returned by default
	Expand *ReportingRevisionsExpand
	// (optional) Return only the those revisions of work items, where only history field was changed
	IncludeDiscussionChangesOnly *bool
	// (optional) The maximum number of results to return in this batch
	MaxPageSize *int
}

// Get a batch of work item revisions. This request may be used if your list of fields is large enough that it may run the URL over the length limit.
func (client *ClientImpl) ReadReportingRevisionsPost(ctx context.Context, args ReadReportingRevisionsPostArgs) (*ReportingWorkItemRevisionsBatch, error) {
	if args.Filter == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Filter"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}

	queryParams := url.Values{}
	if args.ContinuationToken != nil {
		queryParams.Add("continuationToken", *args.ContinuationToken)
	}
	if args.StartDateTime != nil {
		queryParams.Add("startDateTime", (*args.StartDateTime).String())
	}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	body, marshalErr := json.Marshal(*args.Filter)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("f828fe59-dd87-495d-a17c-7a8d6211ca6c")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0", routeValues, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ReportingWorkItemRevisionsBatch
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ReadReportingRevisionsPost function
type ReadReportingRevisionsPostArgs struct {
	// (required) An object that contains request settings: field filter, type filter, identity format
	Filter *ReportingWorkItemRevisionsFilter
	// (optional) Project ID or project name
	Project *string
	// (optional) Specifies the watermark to start the batch from. Omit this parameter to get the first batch of revisions.
	ContinuationToken *string
	// (optional) Date/time to use as a starting point for revisions, all revisions will occur after this date/time. Cannot be used in conjunction with 'watermark' parameter.
	StartDateTime *azuredevops.Time
	// (optional)
	Expand *ReportingRevisionsExpand
}

// [Preview API] Replace template contents
func (client *ClientImpl) ReplaceTemplate(ctx context.Context, args ReplaceTemplateArgs) (*WorkItemTemplate, error) {
	if args.TemplateContent == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.TemplateContent"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Team == nil || *args.Team == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Team"}
	}
	routeValues["team"] = *args.Team
	if args.TemplateId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.TemplateId"}
	}
	routeValues["templateId"] = (*args.TemplateId).String()

	body, marshalErr := json.Marshal(*args.TemplateContent)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("fb10264a-8836-48a0-8033-1b0ccd2748d5")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "6.0-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemTemplate
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ReplaceTemplate function
type ReplaceTemplateArgs struct {
	// (required) Template contents to replace with
	TemplateContent *WorkItemTemplate
	// (required) Project ID or project name
	Project *string
	// (required) Team ID or team name
	Team *string
	// (required) Template id
	TemplateId *uuid.UUID
}

// Restores the deleted work item from Recycle Bin.
func (client *ClientImpl) RestoreWorkItem(ctx context.Context, args RestoreWorkItemArgs) (*WorkItemDelete, error) {
	if args.Payload == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Payload"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)

	body, marshalErr := json.Marshal(*args.Payload)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("b70d8d39-926c-465e-b927-b1bf0e5ca0e0")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "6.0", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemDelete
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the RestoreWorkItem function
type RestoreWorkItemArgs struct {
	// (required) Paylod with instructions to update the IsDeleted flag to false
	Payload *WorkItemDeleteUpdate
	// (required) ID of the work item to be restored
	Id *int
	// (optional) Project ID or project name
	Project *string
}

// Searches all queries the user has access to in the current project
func (client *ClientImpl) SearchQueries(ctx context.Context, args SearchQueriesArgs) (*QueryHierarchyItemsResult, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	queryParams := url.Values{}
	if args.Filter == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "filter"}
	}
	queryParams.Add("$filter", *args.Filter)
	if args.Top != nil {
		queryParams.Add("$top", strconv.Itoa(*args.Top))
	}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	if args.IncludeDeleted != nil {
		queryParams.Add("$includeDeleted", strconv.FormatBool(*args.IncludeDeleted))
	}
	locationId, _ := uuid.Parse("a67d190c-c41f-424b-814d-0e906f659301")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue QueryHierarchyItemsResult
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the SearchQueries function
type SearchQueriesArgs struct {
	// (required) Project ID or project name
	Project *string
	// (required) The text to filter the queries with.
	Filter *string
	// (optional) The number of queries to return (Default is 50 and maximum is 200).
	Top *int
	// (optional)
	Expand *QueryExpand
	// (optional) Include deleted queries and folders
	IncludeDeleted *bool
}

// Update an existing classification node.
func (client *ClientImpl) UpdateClassificationNode(ctx context.Context, args UpdateClassificationNodeArgs) (*WorkItemClassificationNode, error) {
	if args.PostedNode == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PostedNode"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.StructureGroup == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.StructureGroup"}
	}
	routeValues["structureGroup"] = string(*args.StructureGroup)
	if args.Path != nil && *args.Path != "" {
		routeValues["path"] = *args.Path
	}

	body, marshalErr := json.Marshal(*args.PostedNode)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("5a172953-1b41-49d3-840a-33f79c3ce89f")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "6.0", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemClassificationNode
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateClassificationNode function
type UpdateClassificationNodeArgs struct {
	// (required) Node to create or update.
	PostedNode *WorkItemClassificationNode
	// (required) Project ID or project name
	Project *string
	// (required) Structure group of the classification node, area or iteration.
	StructureGroup *TreeStructureGroup
	// (optional) Path of the classification node.
	Path *string
}

// [Preview API] Update a comment on a work item.
func (client *ClientImpl) UpdateComment(ctx context.Context, args UpdateCommentArgs) (*Comment, error) {
	if args.Request == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Request"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.WorkItemId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.WorkItemId"}
	}
	routeValues["workItemId"] = strconv.Itoa(*args.WorkItemId)
	if args.CommentId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.CommentId"}
	}
	routeValues["commentId"] = strconv.Itoa(*args.CommentId)

	body, marshalErr := json.Marshal(*args.Request)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("608aac0a-32e1-4493-a863-b9cf4566d257")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "6.0-preview.3", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Comment
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateComment function
type UpdateCommentArgs struct {
	// (required) Comment update request.
	Request *CommentUpdate
	// (required) Project ID or project name
	Project *string
	// (required) Id of a work item.
	WorkItemId *int
	// (required)
	CommentId *int
}

// Update a field.
func (client *ClientImpl) UpdateField(ctx context.Context, args UpdateFieldArgs) (*WorkItemField, error) {
	if args.Payload == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Payload"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.FieldNameOrRefName == nil || *args.FieldNameOrRefName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.FieldNameOrRefName"}
	}
	routeValues["fieldNameOrRefName"] = *args.FieldNameOrRefName

	body, marshalErr := json.Marshal(*args.Payload)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("b51fd764-e5c2-4b9b-aaf7-3395cf4bdd94")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "6.0", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemField
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateField function
type UpdateFieldArgs struct {
	// (required) Payload contains desired value of the field's properties
	Payload *UpdateWorkItemField
	// (required) Name/reference name of the field to be updated
	FieldNameOrRefName *string
	// (optional) Project ID or project name
	Project *string
}

// Update a query or a folder. This allows you to update, rename and move queries and folders.
func (client *ClientImpl) UpdateQuery(ctx context.Context, args UpdateQueryArgs) (*QueryHierarchyItem, error) {
	if args.QueryUpdate == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.QueryUpdate"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Query == nil || *args.Query == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Query"}
	}
	routeValues["query"] = *args.Query

	queryParams := url.Values{}
	if args.UndeleteDescendants != nil {
		queryParams.Add("$undeleteDescendants", strconv.FormatBool(*args.UndeleteDescendants))
	}
	body, marshalErr := json.Marshal(*args.QueryUpdate)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("a67d190c-c41f-424b-814d-0e906f659301")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "6.0", routeValues, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue QueryHierarchyItem
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateQuery function
type UpdateQueryArgs struct {
	// (required) The query to update.
	QueryUpdate *QueryHierarchyItem
	// (required) Project ID or project name
	Project *string
	// (required) The ID or path for the query to update.
	Query *string
	// (optional) Undelete the children of this folder. It is important to note that this will not bring back the permission changes that were previously applied to the descendants.
	UndeleteDescendants *bool
}

// [Preview API]
func (client *ClientImpl) UpdateTag(ctx context.Context, args UpdateTagArgs) (*WorkItemTagDefinition, error) {
	if args.TagData == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.TagData"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.TagIdOrName == nil || *args.TagIdOrName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.TagIdOrName"}
	}
	routeValues["tagIdOrName"] = *args.TagIdOrName

	body, marshalErr := json.Marshal(*args.TagData)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("bc15bc60-e7a8-43cb-ab01-2106be3983a1")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "6.0-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItemTagDefinition
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateTag function
type UpdateTagArgs struct {
	// (required)
	TagData *WorkItemTagDefinition
	// (required) Project ID or project name
	Project *string
	// (required)
	TagIdOrName *string
}

// Updates a single work item.
func (client *ClientImpl) UpdateWorkItem(ctx context.Context, args UpdateWorkItemArgs) (*WorkItem, error) {
	if args.Document == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Document"}
	}
	routeValues := make(map[string]string)
	if args.Project != nil && *args.Project != "" {
		routeValues["project"] = *args.Project
	}
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)

	queryParams := url.Values{}
	if args.ValidateOnly != nil {
		queryParams.Add("validateOnly", strconv.FormatBool(*args.ValidateOnly))
	}
	if args.BypassRules != nil {
		queryParams.Add("bypassRules", strconv.FormatBool(*args.BypassRules))
	}
	if args.SuppressNotifications != nil {
		queryParams.Add("suppressNotifications", strconv.FormatBool(*args.SuppressNotifications))
	}
	if args.Expand != nil {
		queryParams.Add("$expand", string(*args.Expand))
	}
	body, marshalErr := json.Marshal(*args.Document)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("72c7ddf8-2cdc-4f60-90cd-ab71c14a399b")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "6.0", routeValues, queryParams, bytes.NewReader(body), "application/json-patch+json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue WorkItem
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateWorkItem function
type UpdateWorkItemArgs struct {
	// (required) The JSON Patch document representing the update
	Document *[]webapi.JsonPatchOperation
	// (required) The id of the work item to update
	Id *int
	// (optional) Project ID or project name
	Project *string
	// (optional) Indicate if you only want to validate the changes without saving the work item
	ValidateOnly *bool
	// (optional) Do not enforce the work item type rules on this update
	BypassRules *bool
	// (optional) Do not fire any notifications for this change
	SuppressNotifications *bool
	// (optional) The expand parameters for work item attributes. Possible options are { None, Relations, Fields, Links, All }.
	Expand *WorkItemExpand
}
