// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package gallery

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

var ResourceAreaId, _ = uuid.Parse("69d21c00-f135-441b-b5ce-3626378e0819")

type Client interface {
	// [Preview API]
	AddAssetForEditExtensionDraft(context.Context, AddAssetForEditExtensionDraftArgs) (*ExtensionDraftAsset, error)
	// [Preview API]
	AddAssetForNewExtensionDraft(context.Context, AddAssetForNewExtensionDraftArgs) (*ExtensionDraftAsset, error)
	// [Preview API]
	AssociateAzurePublisher(context.Context, AssociateAzurePublisherArgs) (*AzurePublisher, error)
	// [Preview API]
	CreateCategory(context.Context, CreateCategoryArgs) (*ExtensionCategory, error)
	// [Preview API]
	CreateDraftForEditExtension(context.Context, CreateDraftForEditExtensionArgs) (*ExtensionDraft, error)
	// [Preview API]
	CreateDraftForNewExtension(context.Context, CreateDraftForNewExtensionArgs) (*ExtensionDraft, error)
	// [Preview API]
	CreateExtension(context.Context, CreateExtensionArgs) (*PublishedExtension, error)
	// [Preview API]
	CreateExtensionWithPublisher(context.Context, CreateExtensionWithPublisherArgs) (*PublishedExtension, error)
	// [Preview API]
	CreatePublisher(context.Context, CreatePublisherArgs) (*Publisher, error)
	// [Preview API] Creates a new question for an extension.
	CreateQuestion(context.Context, CreateQuestionArgs) (*Question, error)
	// [Preview API] Creates a new response for a given question for an extension.
	CreateResponse(context.Context, CreateResponseArgs) (*Response, error)
	// [Preview API] Creates a new review for an extension
	CreateReview(context.Context, CreateReviewArgs) (*Review, error)
	// [Preview API]
	CreateSupportRequest(context.Context, CreateSupportRequestArgs) error
	// [Preview API]
	DeleteExtension(context.Context, DeleteExtensionArgs) error
	// [Preview API]
	DeleteExtensionById(context.Context, DeleteExtensionByIdArgs) error
	// [Preview API]
	DeletePublisher(context.Context, DeletePublisherArgs) error
	// [Preview API] Delete publisher asset like logo
	DeletePublisherAsset(context.Context, DeletePublisherAssetArgs) error
	// [Preview API] Deletes an existing question and all its associated responses for an extension. (soft delete)
	DeleteQuestion(context.Context, DeleteQuestionArgs) error
	// [Preview API] Deletes a response for an extension. (soft delete)
	DeleteResponse(context.Context, DeleteResponseArgs) error
	// [Preview API] Deletes a review
	DeleteReview(context.Context, DeleteReviewArgs) error
	// [Preview API]
	ExtensionValidator(context.Context, ExtensionValidatorArgs) error
	// [Preview API]
	FetchDomainToken(context.Context, FetchDomainTokenArgs) (*string, error)
	// [Preview API]
	GenerateKey(context.Context, GenerateKeyArgs) error
	// [Preview API]
	GetAcquisitionOptions(context.Context, GetAcquisitionOptionsArgs) (*AcquisitionOptions, error)
	// [Preview API]
	GetAsset(context.Context, GetAssetArgs) (io.ReadCloser, error)
	// [Preview API]
	GetAssetAuthenticated(context.Context, GetAssetAuthenticatedArgs) (io.ReadCloser, error)
	// [Preview API]
	GetAssetByName(context.Context, GetAssetByNameArgs) (io.ReadCloser, error)
	// [Preview API]
	GetAssetFromEditExtensionDraft(context.Context, GetAssetFromEditExtensionDraftArgs) (io.ReadCloser, error)
	// [Preview API]
	GetAssetFromNewExtensionDraft(context.Context, GetAssetFromNewExtensionDraftArgs) (io.ReadCloser, error)
	// [Preview API]
	GetAssetWithToken(context.Context, GetAssetWithTokenArgs) (io.ReadCloser, error)
	// [Preview API]
	GetCategories(context.Context, GetCategoriesArgs) (*[]string, error)
	// [Preview API]
	GetCategoryDetails(context.Context, GetCategoryDetailsArgs) (*CategoriesResult, error)
	// [Preview API]
	GetCategoryTree(context.Context, GetCategoryTreeArgs) (*ProductCategory, error)
	// [Preview API]
	GetCertificate(context.Context, GetCertificateArgs) (io.ReadCloser, error)
	// [Preview API]
	GetContentVerificationLog(context.Context, GetContentVerificationLogArgs) (io.ReadCloser, error)
	// [Preview API]
	GetExtension(context.Context, GetExtensionArgs) (*PublishedExtension, error)
	// [Preview API]
	GetExtensionById(context.Context, GetExtensionByIdArgs) (*PublishedExtension, error)
	// [Preview API]
	GetExtensionDailyStats(context.Context, GetExtensionDailyStatsArgs) (*ExtensionDailyStats, error)
	// [Preview API] This route/location id only supports HTTP POST anonymously, so that the page view daily stat can be incremented from Marketplace client. Trying to call GET on this route should result in an exception. Without this explicit implementation, calling GET on this public route invokes the above GET implementation GetExtensionDailyStats.
	GetExtensionDailyStatsAnonymous(context.Context, GetExtensionDailyStatsAnonymousArgs) (*ExtensionDailyStats, error)
	// [Preview API] Get install/uninstall events of an extension. If both count and afterDate parameters are specified, count takes precedence.
	GetExtensionEvents(context.Context, GetExtensionEventsArgs) (*ExtensionEvents, error)
	// [Preview API] Returns extension reports
	GetExtensionReports(context.Context, GetExtensionReportsArgs) (interface{}, error)
	// [Preview API] Get all setting entries for the given user/all-users scope
	GetGalleryUserSettings(context.Context, GetGalleryUserSettingsArgs) (*map[string]interface{}, error)
	// [Preview API] This endpoint gets hit when you download a VSTS extension from the Web UI
	GetPackage(context.Context, GetPackageArgs) (io.ReadCloser, error)
	// [Preview API]
	GetPublisher(context.Context, GetPublisherArgs) (*Publisher, error)
	// [Preview API] Get publisher asset like logo as a stream
	GetPublisherAsset(context.Context, GetPublisherAssetArgs) (io.ReadCloser, error)
	// [Preview API]
	GetPublisherWithoutToken(context.Context, GetPublisherWithoutTokenArgs) (*Publisher, error)
	// [Preview API] Returns a list of questions with their responses associated with an extension.
	GetQuestions(context.Context, GetQuestionsArgs) (*QuestionsResult, error)
	// [Preview API] Returns a list of reviews associated with an extension
	GetReviews(context.Context, GetReviewsArgs) (*ReviewsResult, error)
	// [Preview API] Returns a summary of the reviews
	GetReviewsSummary(context.Context, GetReviewsSummaryArgs) (*ReviewSummary, error)
	// [Preview API]
	GetRootCategories(context.Context, GetRootCategoriesArgs) (*ProductCategoriesResult, error)
	// [Preview API]
	GetSigningKey(context.Context, GetSigningKeyArgs) (*string, error)
	// [Preview API]
	GetVerificationLog(context.Context, GetVerificationLogArgs) (io.ReadCloser, error)
	// [Preview API] Increments a daily statistic associated with the extension
	IncrementExtensionDailyStat(context.Context, IncrementExtensionDailyStatArgs) error
	// [Preview API]
	PerformEditExtensionDraftOperation(context.Context, PerformEditExtensionDraftOperationArgs) (*ExtensionDraft, error)
	// [Preview API]
	PerformNewExtensionDraftOperation(context.Context, PerformNewExtensionDraftOperationArgs) (*ExtensionDraft, error)
	// [Preview API] API endpoint to publish extension install/uninstall events. This is meant to be invoked by EMS only for sending us data related to install/uninstall of an extension.
	PublishExtensionEvents(context.Context, PublishExtensionEventsArgs) error
	// [Preview API]
	QueryAssociatedAzurePublisher(context.Context, QueryAssociatedAzurePublisherArgs) (*AzurePublisher, error)
	// [Preview API]
	QueryExtensions(context.Context, QueryExtensionsArgs) (*ExtensionQueryResult, error)
	// [Preview API]
	QueryPublishers(context.Context, QueryPublishersArgs) (*PublisherQueryResult, error)
	// [Preview API] Flags a concern with an existing question for an extension.
	ReportQuestion(context.Context, ReportQuestionArgs) (*Concern, error)
	// [Preview API]
	RequestAcquisition(context.Context, RequestAcquisitionArgs) (*ExtensionAcquisitionRequest, error)
	// [Preview API] Send Notification
	SendNotifications(context.Context, SendNotificationsArgs) error
	// [Preview API] Set all setting entries for the given user/all-users scope
	SetGalleryUserSettings(context.Context, SetGalleryUserSettingsArgs) error
	// [Preview API]
	ShareExtension(context.Context, ShareExtensionArgs) error
	// [Preview API]
	ShareExtensionById(context.Context, ShareExtensionByIdArgs) error
	// [Preview API]
	ShareExtensionWithHost(context.Context, ShareExtensionWithHostArgs) error
	// [Preview API]
	UnshareExtension(context.Context, UnshareExtensionArgs) error
	// [Preview API]
	UnshareExtensionById(context.Context, UnshareExtensionByIdArgs) error
	// [Preview API]
	UnshareExtensionWithHost(context.Context, UnshareExtensionWithHostArgs) error
	// [Preview API] REST endpoint to update an extension.
	UpdateExtension(context.Context, UpdateExtensionArgs) (*PublishedExtension, error)
	// [Preview API]
	UpdateExtensionById(context.Context, UpdateExtensionByIdArgs) (*PublishedExtension, error)
	// [Preview API]
	UpdateExtensionProperties(context.Context, UpdateExtensionPropertiesArgs) (*PublishedExtension, error)
	// [Preview API]
	UpdateExtensionStatistics(context.Context, UpdateExtensionStatisticsArgs) error
	// [Preview API]
	UpdatePayloadInDraftForEditExtension(context.Context, UpdatePayloadInDraftForEditExtensionArgs) (*ExtensionDraft, error)
	// [Preview API]
	UpdatePayloadInDraftForNewExtension(context.Context, UpdatePayloadInDraftForNewExtensionArgs) (*ExtensionDraft, error)
	// [Preview API]
	UpdatePublisher(context.Context, UpdatePublisherArgs) (*Publisher, error)
	// [Preview API] Update publisher asset like logo. It accepts asset file as an octet stream and file name is passed in header values.
	UpdatePublisherAsset(context.Context, UpdatePublisherAssetArgs) (*map[string]string, error)
	// [Preview API] Endpoint to add/modify publisher membership. Currently Supports only addition/modification of 1 user at a time Works only for adding members of same tenant.
	UpdatePublisherMembers(context.Context, UpdatePublisherMembersArgs) (*[]PublisherRoleAssignment, error)
	// [Preview API] Updates an existing question for an extension.
	UpdateQuestion(context.Context, UpdateQuestionArgs) (*Question, error)
	// [Preview API] Updates an existing response for a given question for an extension.
	UpdateResponse(context.Context, UpdateResponseArgs) (*Response, error)
	// [Preview API] Updates or Flags a review
	UpdateReview(context.Context, UpdateReviewArgs) (*ReviewPatch, error)
	// [Preview API]
	UpdateVSCodeWebExtensionStatistics(context.Context, UpdateVSCodeWebExtensionStatisticsArgs) error
	// [Preview API]
	VerifyDomainToken(context.Context, VerifyDomainTokenArgs) error
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

// [Preview API]
func (client *ClientImpl) AddAssetForEditExtensionDraft(ctx context.Context, args AddAssetForEditExtensionDraftArgs) (*ExtensionDraftAsset, error) {
	if args.UploadStream == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UploadStream"}
	}
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.DraftId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DraftId"}
	}
	routeValues["draftId"] = (*args.DraftId).String()
	if args.AssetType == nil || *args.AssetType == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.AssetType"}
	}
	routeValues["assetType"] = *args.AssetType

	locationId, _ := uuid.Parse("f1db9c47-6619-4998-a7e5-d7f9f41a4617")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.1", routeValues, nil, args.UploadStream, "application/octet-stream", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ExtensionDraftAsset
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the AddAssetForEditExtensionDraft function
type AddAssetForEditExtensionDraftArgs struct {
	// (required) Stream to upload
	UploadStream io.Reader
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
	// (required)
	DraftId *uuid.UUID
	// (required)
	AssetType *string
}

// [Preview API]
func (client *ClientImpl) AddAssetForNewExtensionDraft(ctx context.Context, args AddAssetForNewExtensionDraftArgs) (*ExtensionDraftAsset, error) {
	if args.UploadStream == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UploadStream"}
	}
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.DraftId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DraftId"}
	}
	routeValues["draftId"] = (*args.DraftId).String()
	if args.AssetType == nil || *args.AssetType == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.AssetType"}
	}
	routeValues["assetType"] = *args.AssetType

	locationId, _ := uuid.Parse("88c0b1c8-b4f1-498a-9b2a-8446ef9f32e7")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.1", routeValues, nil, args.UploadStream, "application/octet-stream", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ExtensionDraftAsset
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the AddAssetForNewExtensionDraft function
type AddAssetForNewExtensionDraftArgs struct {
	// (required) Stream to upload
	UploadStream io.Reader
	// (required)
	PublisherName *string
	// (required)
	DraftId *uuid.UUID
	// (required)
	AssetType *string
}

// [Preview API]
func (client *ClientImpl) AssociateAzurePublisher(ctx context.Context, args AssociateAzurePublisherArgs) (*AzurePublisher, error) {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName

	queryParams := url.Values{}
	if args.AzurePublisherId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "azurePublisherId"}
	}
	queryParams.Add("azurePublisherId", *args.AzurePublisherId)
	locationId, _ := uuid.Parse("efd202a6-9d87-4ebc-9229-d2b8ae2fdb6d")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue AzurePublisher
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the AssociateAzurePublisher function
type AssociateAzurePublisherArgs struct {
	// (required)
	PublisherName *string
	// (required)
	AzurePublisherId *string
}

// [Preview API]
func (client *ClientImpl) CreateCategory(ctx context.Context, args CreateCategoryArgs) (*ExtensionCategory, error) {
	if args.Category == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Category"}
	}
	body, marshalErr := json.Marshal(*args.Category)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("476531a3-7024-4516-a76a-ed64d3008ad6")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ExtensionCategory
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateCategory function
type CreateCategoryArgs struct {
	// (required)
	Category *ExtensionCategory
}

// [Preview API]
func (client *ClientImpl) CreateDraftForEditExtension(ctx context.Context, args CreateDraftForEditExtensionArgs) (*ExtensionDraft, error) {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName

	locationId, _ := uuid.Parse("02b33873-4e61-496e-83a2-59d1df46b7d8")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ExtensionDraft
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateDraftForEditExtension function
type CreateDraftForEditExtensionArgs struct {
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
}

// [Preview API]
func (client *ClientImpl) CreateDraftForNewExtension(ctx context.Context, args CreateDraftForNewExtensionArgs) (*ExtensionDraft, error) {
	if args.UploadStream == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UploadStream"}
	}
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName

	additionalHeaders := make(map[string]string)
	if args.Product != nil {
		additionalHeaders["X-Market-UploadFileProduct"] = *args.Product
	}
	if args.FileName != nil {
		additionalHeaders["X-Market-UploadFileName"] = *args.FileName
	}
	locationId, _ := uuid.Parse("b3ab127d-ebb9-4d22-b611-4e09593c8d79")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, args.UploadStream, "application/octet-stream", "application/json", additionalHeaders)
	if err != nil {
		return nil, err
	}

	var responseValue ExtensionDraft
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateDraftForNewExtension function
type CreateDraftForNewExtensionArgs struct {
	// (required) Stream to upload
	UploadStream io.Reader
	// (required)
	PublisherName *string
	// (required) Header to pass the product type of the payload file
	Product *string
	// (optional) Header to pass the filename of the uploaded data
	FileName *string
}

// [Preview API]
func (client *ClientImpl) CreateExtension(ctx context.Context, args CreateExtensionArgs) (*PublishedExtension, error) {
	if args.UploadStream == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UploadStream"}
	}
	queryParams := url.Values{}
	if args.ExtensionType != nil {
		queryParams.Add("extensionType", *args.ExtensionType)
	}
	if args.ReCaptchaToken != nil {
		queryParams.Add("reCaptchaToken", *args.ReCaptchaToken)
	}
	locationId, _ := uuid.Parse("a41192c8-9525-4b58-bc86-179fa549d80d")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.2", nil, queryParams, args.UploadStream, "application/octet-stream", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue PublishedExtension
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateExtension function
type CreateExtensionArgs struct {
	// (required) Stream to upload
	UploadStream io.Reader
	// (optional)
	ExtensionType *string
	// (optional)
	ReCaptchaToken *string
}

// [Preview API]
func (client *ClientImpl) CreateExtensionWithPublisher(ctx context.Context, args CreateExtensionWithPublisherArgs) (*PublishedExtension, error) {
	if args.UploadStream == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UploadStream"}
	}
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName

	queryParams := url.Values{}
	if args.ExtensionType != nil {
		queryParams.Add("extensionType", *args.ExtensionType)
	}
	if args.ReCaptchaToken != nil {
		queryParams.Add("reCaptchaToken", *args.ReCaptchaToken)
	}
	locationId, _ := uuid.Parse("e11ea35a-16fe-4b80-ab11-c4cab88a0966")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.2", routeValues, queryParams, args.UploadStream, "application/octet-stream", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue PublishedExtension
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateExtensionWithPublisher function
type CreateExtensionWithPublisherArgs struct {
	// (required) Stream to upload
	UploadStream io.Reader
	// (required)
	PublisherName *string
	// (optional)
	ExtensionType *string
	// (optional)
	ReCaptchaToken *string
}

// [Preview API]
func (client *ClientImpl) CreatePublisher(ctx context.Context, args CreatePublisherArgs) (*Publisher, error) {
	if args.Publisher == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Publisher"}
	}
	body, marshalErr := json.Marshal(*args.Publisher)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("4ddec66a-e4f6-4f5d-999e-9e77710d7ff4")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Publisher
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreatePublisher function
type CreatePublisherArgs struct {
	// (required)
	Publisher *Publisher
}

// [Preview API] Creates a new question for an extension.
func (client *ClientImpl) CreateQuestion(ctx context.Context, args CreateQuestionArgs) (*Question, error) {
	if args.Question == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Question"}
	}
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName

	body, marshalErr := json.Marshal(*args.Question)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("6d1d9741-eca8-4701-a3a5-235afc82dfa4")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Question
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateQuestion function
type CreateQuestionArgs struct {
	// (required) Question to be created for the extension.
	Question *Question
	// (required) Name of the publisher who published the extension.
	PublisherName *string
	// (required) Name of the extension.
	ExtensionName *string
}

// [Preview API] Creates a new response for a given question for an extension.
func (client *ClientImpl) CreateResponse(ctx context.Context, args CreateResponseArgs) (*Response, error) {
	if args.Response == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Response"}
	}
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.QuestionId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.QuestionId"}
	}
	routeValues["questionId"] = strconv.FormatUint(*args.QuestionId, 10)

	body, marshalErr := json.Marshal(*args.Response)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("7f8ae5e0-46b0-438f-b2e8-13e8513517bd")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Response
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateResponse function
type CreateResponseArgs struct {
	// (required) Response to be created for the extension.
	Response *Response
	// (required) Name of the publisher who published the extension.
	PublisherName *string
	// (required) Name of the extension.
	ExtensionName *string
	// (required) Identifier of the question for which response is to be created for the extension.
	QuestionId *uint64
}

// [Preview API] Creates a new review for an extension
func (client *ClientImpl) CreateReview(ctx context.Context, args CreateReviewArgs) (*Review, error) {
	if args.Review == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Review"}
	}
	routeValues := make(map[string]string)
	if args.PubName == nil || *args.PubName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PubName"}
	}
	routeValues["pubName"] = *args.PubName
	if args.ExtName == nil || *args.ExtName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtName"}
	}
	routeValues["extName"] = *args.ExtName

	body, marshalErr := json.Marshal(*args.Review)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("e6e85b9d-aa70-40e6-aa28-d0fbf40b91a3")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Review
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateReview function
type CreateReviewArgs struct {
	// (required) Review to be created for the extension
	Review *Review
	// (required) Name of the publisher who published the extension
	PubName *string
	// (required) Name of the extension
	ExtName *string
}

// [Preview API]
func (client *ClientImpl) CreateSupportRequest(ctx context.Context, args CreateSupportRequestArgs) error {
	if args.CustomerSupportRequest == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.CustomerSupportRequest"}
	}
	body, marshalErr := json.Marshal(*args.CustomerSupportRequest)
	if marshalErr != nil {
		return marshalErr
	}
	locationId, _ := uuid.Parse("8eded385-026a-4c15-b810-b8eb402771f1")
	_, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the CreateSupportRequest function
type CreateSupportRequestArgs struct {
	// (required)
	CustomerSupportRequest *CustomerSupportRequest
}

// [Preview API]
func (client *ClientImpl) DeleteExtension(ctx context.Context, args DeleteExtensionArgs) error {
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
	if args.Version != nil {
		queryParams.Add("version", *args.Version)
	}
	locationId, _ := uuid.Parse("e11ea35a-16fe-4b80-ab11-c4cab88a0966")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.2", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteExtension function
type DeleteExtensionArgs struct {
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
	// (optional)
	Version *string
}

// [Preview API]
func (client *ClientImpl) DeleteExtensionById(ctx context.Context, args DeleteExtensionByIdArgs) error {
	routeValues := make(map[string]string)
	if args.ExtensionId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.ExtensionId"}
	}
	routeValues["extensionId"] = (*args.ExtensionId).String()

	queryParams := url.Values{}
	if args.Version != nil {
		queryParams.Add("version", *args.Version)
	}
	locationId, _ := uuid.Parse("a41192c8-9525-4b58-bc86-179fa549d80d")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.2", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteExtensionById function
type DeleteExtensionByIdArgs struct {
	// (required)
	ExtensionId *uuid.UUID
	// (optional)
	Version *string
}

// [Preview API]
func (client *ClientImpl) DeletePublisher(ctx context.Context, args DeletePublisherArgs) error {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName

	locationId, _ := uuid.Parse("4ddec66a-e4f6-4f5d-999e-9e77710d7ff4")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeletePublisher function
type DeletePublisherArgs struct {
	// (required)
	PublisherName *string
}

// [Preview API] Delete publisher asset like logo
func (client *ClientImpl) DeletePublisherAsset(ctx context.Context, args DeletePublisherAssetArgs) error {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName

	queryParams := url.Values{}
	if args.AssetType != nil {
		queryParams.Add("assetType", *args.AssetType)
	}
	locationId, _ := uuid.Parse("21143299-34f9-4c62-8ca8-53da691192f9")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeletePublisherAsset function
type DeletePublisherAssetArgs struct {
	// (required) Internal name of the publisher
	PublisherName *string
	// (optional) Type of asset. Default value is 'logo'.
	AssetType *string
}

// [Preview API] Deletes an existing question and all its associated responses for an extension. (soft delete)
func (client *ClientImpl) DeleteQuestion(ctx context.Context, args DeleteQuestionArgs) error {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.QuestionId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.QuestionId"}
	}
	routeValues["questionId"] = strconv.FormatUint(*args.QuestionId, 10)

	locationId, _ := uuid.Parse("6d1d9741-eca8-4701-a3a5-235afc82dfa4")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteQuestion function
type DeleteQuestionArgs struct {
	// (required) Name of the publisher who published the extension.
	PublisherName *string
	// (required) Name of the extension.
	ExtensionName *string
	// (required) Identifier of the question to be deleted for the extension.
	QuestionId *uint64
}

// [Preview API] Deletes a response for an extension. (soft delete)
func (client *ClientImpl) DeleteResponse(ctx context.Context, args DeleteResponseArgs) error {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.QuestionId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.QuestionId"}
	}
	routeValues["questionId"] = strconv.FormatUint(*args.QuestionId, 10)
	if args.ResponseId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.ResponseId"}
	}
	routeValues["responseId"] = strconv.FormatUint(*args.ResponseId, 10)

	locationId, _ := uuid.Parse("7f8ae5e0-46b0-438f-b2e8-13e8513517bd")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteResponse function
type DeleteResponseArgs struct {
	// (required) Name of the publisher who published the extension.
	PublisherName *string
	// (required) Name of the extension.
	ExtensionName *string
	// (required) Identifies the question whose response is to be deleted.
	QuestionId *uint64
	// (required) Identifies the response to be deleted.
	ResponseId *uint64
}

// [Preview API] Deletes a review
func (client *ClientImpl) DeleteReview(ctx context.Context, args DeleteReviewArgs) error {
	routeValues := make(map[string]string)
	if args.PubName == nil || *args.PubName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PubName"}
	}
	routeValues["pubName"] = *args.PubName
	if args.ExtName == nil || *args.ExtName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtName"}
	}
	routeValues["extName"] = *args.ExtName
	if args.ReviewId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.ReviewId"}
	}
	routeValues["reviewId"] = strconv.FormatUint(*args.ReviewId, 10)

	locationId, _ := uuid.Parse("e6e85b9d-aa70-40e6-aa28-d0fbf40b91a3")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteReview function
type DeleteReviewArgs struct {
	// (required) Name of the publisher who published the extension
	PubName *string
	// (required) Name of the extension
	ExtName *string
	// (required) Id of the review which needs to be updated
	ReviewId *uint64
}

// [Preview API]
func (client *ClientImpl) ExtensionValidator(ctx context.Context, args ExtensionValidatorArgs) error {
	if args.AzureRestApiRequestModel == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.AzureRestApiRequestModel"}
	}
	body, marshalErr := json.Marshal(*args.AzureRestApiRequestModel)
	if marshalErr != nil {
		return marshalErr
	}
	locationId, _ := uuid.Parse("05e8a5e1-8c59-4c2c-8856-0ff087d1a844")
	_, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the ExtensionValidator function
type ExtensionValidatorArgs struct {
	// (required)
	AzureRestApiRequestModel *AzureRestApiRequestModel
}

// [Preview API]
func (client *ClientImpl) FetchDomainToken(ctx context.Context, args FetchDomainTokenArgs) (*string, error) {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName

	locationId, _ := uuid.Parse("67a609ef-fa74-4b52-8664-78d76f7b3634")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue string
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the FetchDomainToken function
type FetchDomainTokenArgs struct {
	// (required)
	PublisherName *string
}

// [Preview API]
func (client *ClientImpl) GenerateKey(ctx context.Context, args GenerateKeyArgs) error {
	routeValues := make(map[string]string)
	if args.KeyType == nil || *args.KeyType == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.KeyType"}
	}
	routeValues["keyType"] = *args.KeyType

	queryParams := url.Values{}
	if args.ExpireCurrentSeconds != nil {
		queryParams.Add("expireCurrentSeconds", strconv.Itoa(*args.ExpireCurrentSeconds))
	}
	locationId, _ := uuid.Parse("92ed5cf4-c38b-465a-9059-2f2fb7c624b5")
	_, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the GenerateKey function
type GenerateKeyArgs struct {
	// (required)
	KeyType *string
	// (optional)
	ExpireCurrentSeconds *int
}

// [Preview API]
func (client *ClientImpl) GetAcquisitionOptions(ctx context.Context, args GetAcquisitionOptionsArgs) (*AcquisitionOptions, error) {
	routeValues := make(map[string]string)
	if args.ItemId == nil || *args.ItemId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ItemId"}
	}
	routeValues["itemId"] = *args.ItemId

	queryParams := url.Values{}
	if args.InstallationTarget == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "installationTarget"}
	}
	queryParams.Add("installationTarget", *args.InstallationTarget)
	if args.TestCommerce != nil {
		queryParams.Add("testCommerce", strconv.FormatBool(*args.TestCommerce))
	}
	if args.IsFreeOrTrialInstall != nil {
		queryParams.Add("isFreeOrTrialInstall", strconv.FormatBool(*args.IsFreeOrTrialInstall))
	}
	locationId, _ := uuid.Parse("9d0a0105-075e-4760-aa15-8bcf54d1bd7d")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue AcquisitionOptions
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetAcquisitionOptions function
type GetAcquisitionOptionsArgs struct {
	// (required)
	ItemId *string
	// (required)
	InstallationTarget *string
	// (optional)
	TestCommerce *bool
	// (optional)
	IsFreeOrTrialInstall *bool
}

// [Preview API]
func (client *ClientImpl) GetAsset(ctx context.Context, args GetAssetArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.ExtensionId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ExtensionId"}
	}
	routeValues["extensionId"] = (*args.ExtensionId).String()
	if args.Version == nil || *args.Version == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Version"}
	}
	routeValues["version"] = *args.Version
	if args.AssetType == nil || *args.AssetType == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.AssetType"}
	}
	routeValues["assetType"] = *args.AssetType

	queryParams := url.Values{}
	if args.AccountToken != nil {
		queryParams.Add("accountToken", *args.AccountToken)
	}
	if args.AcceptDefault != nil {
		queryParams.Add("acceptDefault", strconv.FormatBool(*args.AcceptDefault))
	}
	additionalHeaders := make(map[string]string)
	if args.AccountTokenHeader != nil {
		additionalHeaders["X-Market-AccountToken"] = *args.AccountTokenHeader
	}
	locationId, _ := uuid.Parse("5d545f3d-ef47-488b-8be3-f5ee1517856c")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/octet-stream", additionalHeaders)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetAsset function
type GetAssetArgs struct {
	// (required)
	ExtensionId *uuid.UUID
	// (required)
	Version *string
	// (required)
	AssetType *string
	// (optional)
	AccountToken *string
	// (optional)
	AcceptDefault *bool
	// (optional) Header to pass the account token
	AccountTokenHeader *string
}

// [Preview API]
func (client *ClientImpl) GetAssetAuthenticated(ctx context.Context, args GetAssetAuthenticatedArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.Version == nil || *args.Version == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Version"}
	}
	routeValues["version"] = *args.Version
	if args.AssetType == nil || *args.AssetType == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.AssetType"}
	}
	routeValues["assetType"] = *args.AssetType

	queryParams := url.Values{}
	if args.AccountToken != nil {
		queryParams.Add("accountToken", *args.AccountToken)
	}
	additionalHeaders := make(map[string]string)
	if args.AccountTokenHeader != nil {
		additionalHeaders["X-Market-AccountToken"] = *args.AccountTokenHeader
	}
	locationId, _ := uuid.Parse("506aff36-2622-4f70-8063-77cce6366d20")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/octet-stream", additionalHeaders)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetAssetAuthenticated function
type GetAssetAuthenticatedArgs struct {
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
	// (required)
	Version *string
	// (required)
	AssetType *string
	// (optional)
	AccountToken *string
	// (optional) Header to pass the account token
	AccountTokenHeader *string
}

// [Preview API]
func (client *ClientImpl) GetAssetByName(ctx context.Context, args GetAssetByNameArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.Version == nil || *args.Version == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Version"}
	}
	routeValues["version"] = *args.Version
	if args.AssetType == nil || *args.AssetType == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.AssetType"}
	}
	routeValues["assetType"] = *args.AssetType

	queryParams := url.Values{}
	if args.AccountToken != nil {
		queryParams.Add("accountToken", *args.AccountToken)
	}
	if args.AcceptDefault != nil {
		queryParams.Add("acceptDefault", strconv.FormatBool(*args.AcceptDefault))
	}
	additionalHeaders := make(map[string]string)
	if args.AccountTokenHeader != nil {
		additionalHeaders["X-Market-AccountToken"] = *args.AccountTokenHeader
	}
	locationId, _ := uuid.Parse("7529171f-a002-4180-93ba-685f358a0482")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/octet-stream", additionalHeaders)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetAssetByName function
type GetAssetByNameArgs struct {
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
	// (required)
	Version *string
	// (required)
	AssetType *string
	// (optional)
	AccountToken *string
	// (optional)
	AcceptDefault *bool
	// (optional) Header to pass the account token
	AccountTokenHeader *string
}

// [Preview API]
func (client *ClientImpl) GetAssetFromEditExtensionDraft(ctx context.Context, args GetAssetFromEditExtensionDraftArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.DraftId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DraftId"}
	}
	routeValues["draftId"] = (*args.DraftId).String()
	if args.AssetType == nil || *args.AssetType == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.AssetType"}
	}
	routeValues["assetType"] = *args.AssetType

	queryParams := url.Values{}
	if args.ExtensionName == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "extensionName"}
	}
	queryParams.Add("extensionName", *args.ExtensionName)
	locationId, _ := uuid.Parse("88c0b1c8-b4f1-498a-9b2a-8446ef9f32e7")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/octet-stream", nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetAssetFromEditExtensionDraft function
type GetAssetFromEditExtensionDraftArgs struct {
	// (required)
	PublisherName *string
	// (required)
	DraftId *uuid.UUID
	// (required)
	AssetType *string
	// (required)
	ExtensionName *string
}

// [Preview API]
func (client *ClientImpl) GetAssetFromNewExtensionDraft(ctx context.Context, args GetAssetFromNewExtensionDraftArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.DraftId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DraftId"}
	}
	routeValues["draftId"] = (*args.DraftId).String()
	if args.AssetType == nil || *args.AssetType == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.AssetType"}
	}
	routeValues["assetType"] = *args.AssetType

	locationId, _ := uuid.Parse("88c0b1c8-b4f1-498a-9b2a-8446ef9f32e7")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/octet-stream", nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetAssetFromNewExtensionDraft function
type GetAssetFromNewExtensionDraftArgs struct {
	// (required)
	PublisherName *string
	// (required)
	DraftId *uuid.UUID
	// (required)
	AssetType *string
}

// [Preview API]
func (client *ClientImpl) GetAssetWithToken(ctx context.Context, args GetAssetWithTokenArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.Version == nil || *args.Version == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Version"}
	}
	routeValues["version"] = *args.Version
	if args.AssetType == nil || *args.AssetType == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.AssetType"}
	}
	routeValues["assetType"] = *args.AssetType
	if args.AssetToken != nil && *args.AssetToken != "" {
		routeValues["assetToken"] = *args.AssetToken
	}

	queryParams := url.Values{}
	if args.AccountToken != nil {
		queryParams.Add("accountToken", *args.AccountToken)
	}
	if args.AcceptDefault != nil {
		queryParams.Add("acceptDefault", strconv.FormatBool(*args.AcceptDefault))
	}
	additionalHeaders := make(map[string]string)
	if args.AccountTokenHeader != nil {
		additionalHeaders["X-Market-AccountToken"] = *args.AccountTokenHeader
	}
	locationId, _ := uuid.Parse("364415a1-0077-4a41-a7a0-06edd4497492")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/octet-stream", additionalHeaders)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetAssetWithToken function
type GetAssetWithTokenArgs struct {
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
	// (required)
	Version *string
	// (required)
	AssetType *string
	// (optional)
	AssetToken *string
	// (optional)
	AccountToken *string
	// (optional)
	AcceptDefault *bool
	// (optional) Header to pass the account token
	AccountTokenHeader *string
}

// [Preview API]
func (client *ClientImpl) GetCategories(ctx context.Context, args GetCategoriesArgs) (*[]string, error) {
	queryParams := url.Values{}
	if args.Languages != nil {
		queryParams.Add("languages", *args.Languages)
	}
	locationId, _ := uuid.Parse("e0a5a71e-3ac3-43a0-ae7d-0bb5c3046a2a")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", nil, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []string
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetCategories function
type GetCategoriesArgs struct {
	// (optional)
	Languages *string
}

// [Preview API]
func (client *ClientImpl) GetCategoryDetails(ctx context.Context, args GetCategoryDetailsArgs) (*CategoriesResult, error) {
	routeValues := make(map[string]string)
	if args.CategoryName == nil || *args.CategoryName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.CategoryName"}
	}
	routeValues["categoryName"] = *args.CategoryName

	queryParams := url.Values{}
	if args.Languages != nil {
		queryParams.Add("languages", *args.Languages)
	}
	if args.Product != nil {
		queryParams.Add("product", *args.Product)
	}
	locationId, _ := uuid.Parse("75d3c04d-84d2-4973-acd2-22627587dabc")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue CategoriesResult
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetCategoryDetails function
type GetCategoryDetailsArgs struct {
	// (required)
	CategoryName *string
	// (optional)
	Languages *string
	// (optional)
	Product *string
}

// [Preview API]
func (client *ClientImpl) GetCategoryTree(ctx context.Context, args GetCategoryTreeArgs) (*ProductCategory, error) {
	routeValues := make(map[string]string)
	if args.Product == nil || *args.Product == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Product"}
	}
	routeValues["product"] = *args.Product
	if args.CategoryId == nil || *args.CategoryId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.CategoryId"}
	}
	routeValues["categoryId"] = *args.CategoryId

	queryParams := url.Values{}
	if args.Lcid != nil {
		queryParams.Add("lcid", strconv.Itoa(*args.Lcid))
	}
	if args.Source != nil {
		queryParams.Add("source", *args.Source)
	}
	if args.ProductVersion != nil {
		queryParams.Add("productVersion", *args.ProductVersion)
	}
	if args.Skus != nil {
		queryParams.Add("skus", *args.Skus)
	}
	if args.SubSkus != nil {
		queryParams.Add("subSkus", *args.SubSkus)
	}
	if args.ProductArchitecture != nil {
		queryParams.Add("productArchitecture", *args.ProductArchitecture)
	}
	locationId, _ := uuid.Parse("1102bb42-82b0-4955-8d8a-435d6b4cedd3")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ProductCategory
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetCategoryTree function
type GetCategoryTreeArgs struct {
	// (required)
	Product *string
	// (required)
	CategoryId *string
	// (optional)
	Lcid *int
	// (optional)
	Source *string
	// (optional)
	ProductVersion *string
	// (optional)
	Skus *string
	// (optional)
	SubSkus *string
	// (optional)
	ProductArchitecture *string
}

// [Preview API]
func (client *ClientImpl) GetCertificate(ctx context.Context, args GetCertificateArgs) (io.ReadCloser, error) {
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

	locationId, _ := uuid.Parse("e905ad6a-3f1f-4d08-9f6d-7d357ff8b7d0")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/octet-stream", nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetCertificate function
type GetCertificateArgs struct {
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
	// (optional)
	Version *string
}

// [Preview API]
func (client *ClientImpl) GetContentVerificationLog(ctx context.Context, args GetContentVerificationLogArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName

	locationId, _ := uuid.Parse("c0f1c7c4-3557-4ffb-b774-1e48c4865e99")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/octet-stream", nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetContentVerificationLog function
type GetContentVerificationLogArgs struct {
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
}

// [Preview API]
func (client *ClientImpl) GetExtension(ctx context.Context, args GetExtensionArgs) (*PublishedExtension, error) {
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
	if args.Version != nil {
		queryParams.Add("version", *args.Version)
	}
	if args.Flags != nil {
		queryParams.Add("flags", string(*args.Flags))
	}
	if args.AccountToken != nil {
		queryParams.Add("accountToken", *args.AccountToken)
	}
	additionalHeaders := make(map[string]string)
	if args.AccountTokenHeader != nil {
		additionalHeaders["X-Market-AccountToken"] = *args.AccountTokenHeader
	}
	locationId, _ := uuid.Parse("e11ea35a-16fe-4b80-ab11-c4cab88a0966")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.2", routeValues, queryParams, nil, "", "application/json", additionalHeaders)
	if err != nil {
		return nil, err
	}

	var responseValue PublishedExtension
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetExtension function
type GetExtensionArgs struct {
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
	// (optional)
	Version *string
	// (optional)
	Flags *ExtensionQueryFlags
	// (optional)
	AccountToken *string
	// (optional) Header to pass the account token
	AccountTokenHeader *string
}

// [Preview API]
func (client *ClientImpl) GetExtensionById(ctx context.Context, args GetExtensionByIdArgs) (*PublishedExtension, error) {
	routeValues := make(map[string]string)
	if args.ExtensionId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ExtensionId"}
	}
	routeValues["extensionId"] = (*args.ExtensionId).String()

	queryParams := url.Values{}
	if args.Version != nil {
		queryParams.Add("version", *args.Version)
	}
	if args.Flags != nil {
		queryParams.Add("flags", string(*args.Flags))
	}
	locationId, _ := uuid.Parse("a41192c8-9525-4b58-bc86-179fa549d80d")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.2", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue PublishedExtension
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetExtensionById function
type GetExtensionByIdArgs struct {
	// (required)
	ExtensionId *uuid.UUID
	// (optional)
	Version *string
	// (optional)
	Flags *ExtensionQueryFlags
}

// [Preview API]
func (client *ClientImpl) GetExtensionDailyStats(ctx context.Context, args GetExtensionDailyStatsArgs) (*ExtensionDailyStats, error) {
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
	if args.Days != nil {
		queryParams.Add("days", strconv.Itoa(*args.Days))
	}
	if args.Aggregate != nil {
		queryParams.Add("aggregate", string(*args.Aggregate))
	}
	if args.AfterDate != nil {
		queryParams.Add("afterDate", (*args.AfterDate).AsQueryParameter())
	}
	locationId, _ := uuid.Parse("ae06047e-51c5-4fb4-ab65-7be488544416")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ExtensionDailyStats
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetExtensionDailyStats function
type GetExtensionDailyStatsArgs struct {
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
	// (optional)
	Days *int
	// (optional)
	Aggregate *ExtensionStatsAggregateType
	// (optional)
	AfterDate *azuredevops.Time
}

// [Preview API] This route/location id only supports HTTP POST anonymously, so that the page view daily stat can be incremented from Marketplace client. Trying to call GET on this route should result in an exception. Without this explicit implementation, calling GET on this public route invokes the above GET implementation GetExtensionDailyStats.
func (client *ClientImpl) GetExtensionDailyStatsAnonymous(ctx context.Context, args GetExtensionDailyStatsAnonymousArgs) (*ExtensionDailyStats, error) {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.Version == nil || *args.Version == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Version"}
	}
	routeValues["version"] = *args.Version

	locationId, _ := uuid.Parse("4fa7adb6-ca65-4075-a232-5f28323288ea")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ExtensionDailyStats
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetExtensionDailyStatsAnonymous function
type GetExtensionDailyStatsAnonymousArgs struct {
	// (required) Name of the publisher
	PublisherName *string
	// (required) Name of the extension
	ExtensionName *string
	// (required) Version of the extension
	Version *string
}

// [Preview API] Get install/uninstall events of an extension. If both count and afterDate parameters are specified, count takes precedence.
func (client *ClientImpl) GetExtensionEvents(ctx context.Context, args GetExtensionEventsArgs) (*ExtensionEvents, error) {
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
	if args.Count != nil {
		queryParams.Add("count", strconv.Itoa(*args.Count))
	}
	if args.AfterDate != nil {
		queryParams.Add("afterDate", (*args.AfterDate).AsQueryParameter())
	}
	if args.Include != nil {
		queryParams.Add("include", *args.Include)
	}
	if args.IncludeProperty != nil {
		queryParams.Add("includeProperty", *args.IncludeProperty)
	}
	locationId, _ := uuid.Parse("3d13c499-2168-4d06-bef4-14aba185dcd5")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ExtensionEvents
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetExtensionEvents function
type GetExtensionEventsArgs struct {
	// (required) Name of the publisher
	PublisherName *string
	// (required) Name of the extension
	ExtensionName *string
	// (optional) Count of events to fetch, applies to each event type.
	Count *int
	// (optional) Fetch events that occurred on or after this date
	AfterDate *azuredevops.Time
	// (optional) Filter options. Supported values: install, uninstall, review, acquisition, sales. Default is to fetch all types of events
	Include *string
	// (optional) Event properties to include. Currently only 'lastContactDetails' is supported for uninstall events
	IncludeProperty *string
}

// [Preview API] Returns extension reports
func (client *ClientImpl) GetExtensionReports(ctx context.Context, args GetExtensionReportsArgs) (interface{}, error) {
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
	if args.Days != nil {
		queryParams.Add("days", strconv.Itoa(*args.Days))
	}
	if args.Count != nil {
		queryParams.Add("count", strconv.Itoa(*args.Count))
	}
	if args.AfterDate != nil {
		queryParams.Add("afterDate", (*args.AfterDate).AsQueryParameter())
	}
	locationId, _ := uuid.Parse("79e0c74f-157f-437e-845f-74fbb4121d4c")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue interface{}
	err = client.Client.UnmarshalBody(resp, responseValue)
	return responseValue, err
}

// Arguments for the GetExtensionReports function
type GetExtensionReportsArgs struct {
	// (required) Name of the publisher who published the extension
	PublisherName *string
	// (required) Name of the extension
	ExtensionName *string
	// (optional) Last n days report. If afterDate and days are specified, days will take priority
	Days *int
	// (optional) Number of events to be returned
	Count *int
	// (optional) Use if you want to fetch events newer than the specified date
	AfterDate *azuredevops.Time
}

// [Preview API] Get all setting entries for the given user/all-users scope
func (client *ClientImpl) GetGalleryUserSettings(ctx context.Context, args GetGalleryUserSettingsArgs) (*map[string]interface{}, error) {
	routeValues := make(map[string]string)
	if args.UserScope == nil || *args.UserScope == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.UserScope"}
	}
	routeValues["userScope"] = *args.UserScope
	if args.Key != nil && *args.Key != "" {
		routeValues["key"] = *args.Key
	}

	locationId, _ := uuid.Parse("9b75ece3-7960-401c-848b-148ac01ca350")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue map[string]interface{}
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetGalleryUserSettings function
type GetGalleryUserSettingsArgs struct {
	// (required) User-Scope at which to get the value. Should be "me" for the current user or "host" for all users.
	UserScope *string
	// (optional) Optional key under which to filter all the entries
	Key *string
}

// [Preview API] This endpoint gets hit when you download a VSTS extension from the Web UI
func (client *ClientImpl) GetPackage(ctx context.Context, args GetPackageArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.Version == nil || *args.Version == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Version"}
	}
	routeValues["version"] = *args.Version

	queryParams := url.Values{}
	if args.AccountToken != nil {
		queryParams.Add("accountToken", *args.AccountToken)
	}
	if args.AcceptDefault != nil {
		queryParams.Add("acceptDefault", strconv.FormatBool(*args.AcceptDefault))
	}
	additionalHeaders := make(map[string]string)
	if args.AccountTokenHeader != nil {
		additionalHeaders["X-Market-AccountToken"] = *args.AccountTokenHeader
	}
	locationId, _ := uuid.Parse("7cb576f8-1cae-4c4b-b7b1-e4af5759e965")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/octet-stream", additionalHeaders)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetPackage function
type GetPackageArgs struct {
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
	// (required)
	Version *string
	// (optional)
	AccountToken *string
	// (optional)
	AcceptDefault *bool
	// (optional) Header to pass the account token
	AccountTokenHeader *string
}

// [Preview API]
func (client *ClientImpl) GetPublisher(ctx context.Context, args GetPublisherArgs) (*Publisher, error) {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName

	queryParams := url.Values{}
	if args.Flags != nil {
		queryParams.Add("flags", strconv.Itoa(*args.Flags))
	}
	locationId, _ := uuid.Parse("4ddec66a-e4f6-4f5d-999e-9e77710d7ff4")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Publisher
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetPublisher function
type GetPublisherArgs struct {
	// (required)
	PublisherName *string
	// (optional)
	Flags *int
}

// [Preview API] Get publisher asset like logo as a stream
func (client *ClientImpl) GetPublisherAsset(ctx context.Context, args GetPublisherAssetArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName

	queryParams := url.Values{}
	if args.AssetType != nil {
		queryParams.Add("assetType", *args.AssetType)
	}
	locationId, _ := uuid.Parse("21143299-34f9-4c62-8ca8-53da691192f9")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/octet-stream", nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetPublisherAsset function
type GetPublisherAssetArgs struct {
	// (required) Internal name of the publisher
	PublisherName *string
	// (optional) Type of asset. Default value is 'logo'.
	AssetType *string
}

// [Preview API]
func (client *ClientImpl) GetPublisherWithoutToken(ctx context.Context, args GetPublisherWithoutTokenArgs) (*Publisher, error) {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName

	locationId, _ := uuid.Parse("215a2ed8-458a-4850-ad5a-45f1dabc3461")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Publisher
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetPublisherWithoutToken function
type GetPublisherWithoutTokenArgs struct {
	// (required)
	PublisherName *string
}

// [Preview API] Returns a list of questions with their responses associated with an extension.
func (client *ClientImpl) GetQuestions(ctx context.Context, args GetQuestionsArgs) (*QuestionsResult, error) {
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
	if args.Count != nil {
		queryParams.Add("count", strconv.Itoa(*args.Count))
	}
	if args.Page != nil {
		queryParams.Add("page", strconv.Itoa(*args.Page))
	}
	if args.AfterDate != nil {
		queryParams.Add("afterDate", (*args.AfterDate).AsQueryParameter())
	}
	locationId, _ := uuid.Parse("c010d03d-812c-4ade-ae07-c1862475eda5")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue QuestionsResult
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetQuestions function
type GetQuestionsArgs struct {
	// (required) Name of the publisher who published the extension.
	PublisherName *string
	// (required) Name of the extension.
	ExtensionName *string
	// (optional) Number of questions to retrieve (defaults to 10).
	Count *int
	// (optional) Page number from which set of questions are to be retrieved.
	Page *int
	// (optional) If provided, results questions are returned which were posted after this date
	AfterDate *azuredevops.Time
}

// [Preview API] Returns a list of reviews associated with an extension
func (client *ClientImpl) GetReviews(ctx context.Context, args GetReviewsArgs) (*ReviewsResult, error) {
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
	if args.Count != nil {
		queryParams.Add("count", strconv.Itoa(*args.Count))
	}
	if args.FilterOptions != nil {
		queryParams.Add("filterOptions", string(*args.FilterOptions))
	}
	if args.BeforeDate != nil {
		queryParams.Add("beforeDate", (*args.BeforeDate).AsQueryParameter())
	}
	if args.AfterDate != nil {
		queryParams.Add("afterDate", (*args.AfterDate).AsQueryParameter())
	}
	locationId, _ := uuid.Parse("5b3f819f-f247-42ad-8c00-dd9ab9ab246d")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ReviewsResult
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetReviews function
type GetReviewsArgs struct {
	// (required) Name of the publisher who published the extension
	PublisherName *string
	// (required) Name of the extension
	ExtensionName *string
	// (optional) Number of reviews to retrieve (defaults to 5)
	Count *int
	// (optional) FilterOptions to filter out empty reviews etcetera, defaults to none
	FilterOptions *ReviewFilterOptions
	// (optional) Use if you want to fetch reviews older than the specified date, defaults to null
	BeforeDate *azuredevops.Time
	// (optional) Use if you want to fetch reviews newer than the specified date, defaults to null
	AfterDate *azuredevops.Time
}

// [Preview API] Returns a summary of the reviews
func (client *ClientImpl) GetReviewsSummary(ctx context.Context, args GetReviewsSummaryArgs) (*ReviewSummary, error) {
	routeValues := make(map[string]string)
	if args.PubName == nil || *args.PubName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PubName"}
	}
	routeValues["pubName"] = *args.PubName
	if args.ExtName == nil || *args.ExtName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtName"}
	}
	routeValues["extName"] = *args.ExtName

	queryParams := url.Values{}
	if args.BeforeDate != nil {
		queryParams.Add("beforeDate", (*args.BeforeDate).AsQueryParameter())
	}
	if args.AfterDate != nil {
		queryParams.Add("afterDate", (*args.AfterDate).AsQueryParameter())
	}
	locationId, _ := uuid.Parse("b7b44e21-209e-48f0-ae78-04727fc37d77")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ReviewSummary
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetReviewsSummary function
type GetReviewsSummaryArgs struct {
	// (required) Name of the publisher who published the extension
	PubName *string
	// (required) Name of the extension
	ExtName *string
	// (optional) Use if you want to fetch summary of reviews older than the specified date, defaults to null
	BeforeDate *azuredevops.Time
	// (optional) Use if you want to fetch summary of reviews newer than the specified date, defaults to null
	AfterDate *azuredevops.Time
}

// [Preview API]
func (client *ClientImpl) GetRootCategories(ctx context.Context, args GetRootCategoriesArgs) (*ProductCategoriesResult, error) {
	routeValues := make(map[string]string)
	if args.Product == nil || *args.Product == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Product"}
	}
	routeValues["product"] = *args.Product

	queryParams := url.Values{}
	if args.Lcid != nil {
		queryParams.Add("lcid", strconv.Itoa(*args.Lcid))
	}
	if args.Source != nil {
		queryParams.Add("source", *args.Source)
	}
	if args.ProductVersion != nil {
		queryParams.Add("productVersion", *args.ProductVersion)
	}
	if args.Skus != nil {
		queryParams.Add("skus", *args.Skus)
	}
	if args.SubSkus != nil {
		queryParams.Add("subSkus", *args.SubSkus)
	}
	locationId, _ := uuid.Parse("31fba831-35b2-46f6-a641-d05de5a877d8")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ProductCategoriesResult
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetRootCategories function
type GetRootCategoriesArgs struct {
	// (required)
	Product *string
	// (optional)
	Lcid *int
	// (optional)
	Source *string
	// (optional)
	ProductVersion *string
	// (optional)
	Skus *string
	// (optional)
	SubSkus *string
}

// [Preview API]
func (client *ClientImpl) GetSigningKey(ctx context.Context, args GetSigningKeyArgs) (*string, error) {
	routeValues := make(map[string]string)
	if args.KeyType == nil || *args.KeyType == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.KeyType"}
	}
	routeValues["keyType"] = *args.KeyType

	locationId, _ := uuid.Parse("92ed5cf4-c38b-465a-9059-2f2fb7c624b5")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue string
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetSigningKey function
type GetSigningKeyArgs struct {
	// (required)
	KeyType *string
}

// [Preview API]
func (client *ClientImpl) GetVerificationLog(ctx context.Context, args GetVerificationLogArgs) (io.ReadCloser, error) {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.Version == nil || *args.Version == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Version"}
	}
	routeValues["version"] = *args.Version

	queryParams := url.Values{}
	if args.TargetPlatform != nil {
		queryParams.Add("targetPlatform", *args.TargetPlatform)
	}
	locationId, _ := uuid.Parse("c5523abe-b843-437f-875b-5833064efe4d")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/octet-stream", nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the GetVerificationLog function
type GetVerificationLogArgs struct {
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
	// (required)
	Version *string
	// (optional)
	TargetPlatform *string
}

// [Preview API] Increments a daily statistic associated with the extension
func (client *ClientImpl) IncrementExtensionDailyStat(ctx context.Context, args IncrementExtensionDailyStatArgs) error {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.Version == nil || *args.Version == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Version"}
	}
	routeValues["version"] = *args.Version

	queryParams := url.Values{}
	if args.StatType == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "statType"}
	}
	queryParams.Add("statType", *args.StatType)
	if args.TargetPlatform != nil {
		queryParams.Add("targetPlatform", *args.TargetPlatform)
	}
	locationId, _ := uuid.Parse("4fa7adb6-ca65-4075-a232-5f28323288ea")
	_, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the IncrementExtensionDailyStat function
type IncrementExtensionDailyStatArgs struct {
	// (required) Name of the publisher
	PublisherName *string
	// (required) Name of the extension
	ExtensionName *string
	// (required) Version of the extension
	Version *string
	// (required) Type of stat to increment
	StatType *string
	// (optional)
	TargetPlatform *string
}

// [Preview API]
func (client *ClientImpl) PerformEditExtensionDraftOperation(ctx context.Context, args PerformEditExtensionDraftOperationArgs) (*ExtensionDraft, error) {
	if args.DraftPatch == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DraftPatch"}
	}
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.DraftId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DraftId"}
	}
	routeValues["draftId"] = (*args.DraftId).String()

	body, marshalErr := json.Marshal(*args.DraftPatch)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("02b33873-4e61-496e-83a2-59d1df46b7d8")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ExtensionDraft
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the PerformEditExtensionDraftOperation function
type PerformEditExtensionDraftOperationArgs struct {
	// (required)
	DraftPatch *ExtensionDraftPatch
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
	// (required)
	DraftId *uuid.UUID
}

// [Preview API]
func (client *ClientImpl) PerformNewExtensionDraftOperation(ctx context.Context, args PerformNewExtensionDraftOperationArgs) (*ExtensionDraft, error) {
	if args.DraftPatch == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DraftPatch"}
	}
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.DraftId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DraftId"}
	}
	routeValues["draftId"] = (*args.DraftId).String()

	body, marshalErr := json.Marshal(*args.DraftPatch)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("b3ab127d-ebb9-4d22-b611-4e09593c8d79")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ExtensionDraft
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the PerformNewExtensionDraftOperation function
type PerformNewExtensionDraftOperationArgs struct {
	// (required)
	DraftPatch *ExtensionDraftPatch
	// (required)
	PublisherName *string
	// (required)
	DraftId *uuid.UUID
}

// [Preview API] API endpoint to publish extension install/uninstall events. This is meant to be invoked by EMS only for sending us data related to install/uninstall of an extension.
func (client *ClientImpl) PublishExtensionEvents(ctx context.Context, args PublishExtensionEventsArgs) error {
	if args.ExtensionEvents == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.ExtensionEvents"}
	}
	body, marshalErr := json.Marshal(*args.ExtensionEvents)
	if marshalErr != nil {
		return marshalErr
	}
	locationId, _ := uuid.Parse("0bf2bd3a-70e0-4d5d-8bf7-bd4a9c2ab6e7")
	_, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the PublishExtensionEvents function
type PublishExtensionEventsArgs struct {
	// (required)
	ExtensionEvents *[]ExtensionEvents
}

// [Preview API]
func (client *ClientImpl) QueryAssociatedAzurePublisher(ctx context.Context, args QueryAssociatedAzurePublisherArgs) (*AzurePublisher, error) {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName

	locationId, _ := uuid.Parse("efd202a6-9d87-4ebc-9229-d2b8ae2fdb6d")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue AzurePublisher
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryAssociatedAzurePublisher function
type QueryAssociatedAzurePublisherArgs struct {
	// (required)
	PublisherName *string
}

// [Preview API]
func (client *ClientImpl) QueryExtensions(ctx context.Context, args QueryExtensionsArgs) (*ExtensionQueryResult, error) {
	if args.ExtensionQuery == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ExtensionQuery"}
	}
	queryParams := url.Values{}
	if args.AccountToken != nil {
		queryParams.Add("accountToken", *args.AccountToken)
	}
	additionalHeaders := make(map[string]string)
	if args.AccountTokenHeader != nil {
		additionalHeaders["X-Market-AccountToken"] = *args.AccountTokenHeader
	}
	body, marshalErr := json.Marshal(*args.ExtensionQuery)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("eb9d5ee1-6d43-456b-b80e-8a96fbc014b6")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", nil, queryParams, bytes.NewReader(body), "application/json", "application/json", additionalHeaders)
	if err != nil {
		return nil, err
	}

	var responseValue ExtensionQueryResult
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryExtensions function
type QueryExtensionsArgs struct {
	// (required)
	ExtensionQuery *ExtensionQuery
	// (optional)
	AccountToken *string
	// (optional) Header to pass the account token
	AccountTokenHeader *string
}

// [Preview API]
func (client *ClientImpl) QueryPublishers(ctx context.Context, args QueryPublishersArgs) (*PublisherQueryResult, error) {
	if args.PublisherQuery == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.PublisherQuery"}
	}
	body, marshalErr := json.Marshal(*args.PublisherQuery)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("2ad6ee0a-b53f-4034-9d1d-d009fda1212e")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue PublisherQueryResult
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryPublishers function
type QueryPublishersArgs struct {
	// (required)
	PublisherQuery *PublisherQuery
}

// [Preview API] Flags a concern with an existing question for an extension.
func (client *ClientImpl) ReportQuestion(ctx context.Context, args ReportQuestionArgs) (*Concern, error) {
	if args.Concern == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Concern"}
	}
	routeValues := make(map[string]string)
	if args.PubName == nil || *args.PubName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PubName"}
	}
	routeValues["pubName"] = *args.PubName
	if args.ExtName == nil || *args.ExtName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtName"}
	}
	routeValues["extName"] = *args.ExtName
	if args.QuestionId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.QuestionId"}
	}
	routeValues["questionId"] = strconv.FormatUint(*args.QuestionId, 10)

	body, marshalErr := json.Marshal(*args.Concern)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("784910cd-254a-494d-898b-0728549b2f10")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Concern
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ReportQuestion function
type ReportQuestionArgs struct {
	// (required) User reported concern with a question for the extension.
	Concern *Concern
	// (required) Name of the publisher who published the extension.
	PubName *string
	// (required) Name of the extension.
	ExtName *string
	// (required) Identifier of the question to be updated for the extension.
	QuestionId *uint64
}

// [Preview API]
func (client *ClientImpl) RequestAcquisition(ctx context.Context, args RequestAcquisitionArgs) (*ExtensionAcquisitionRequest, error) {
	if args.AcquisitionRequest == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.AcquisitionRequest"}
	}
	body, marshalErr := json.Marshal(*args.AcquisitionRequest)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("3adb1f2d-e328-446e-be73-9f6d98071c45")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ExtensionAcquisitionRequest
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the RequestAcquisition function
type RequestAcquisitionArgs struct {
	// (required)
	AcquisitionRequest *ExtensionAcquisitionRequest
}

// [Preview API] Send Notification
func (client *ClientImpl) SendNotifications(ctx context.Context, args SendNotificationsArgs) error {
	if args.NotificationData == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.NotificationData"}
	}
	body, marshalErr := json.Marshal(*args.NotificationData)
	if marshalErr != nil {
		return marshalErr
	}
	locationId, _ := uuid.Parse("eab39817-413c-4602-a49f-07ad00844980")
	_, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the SendNotifications function
type SendNotificationsArgs struct {
	// (required) Denoting the data needed to send notification
	NotificationData *NotificationsData
}

// [Preview API] Set all setting entries for the given user/all-users scope
func (client *ClientImpl) SetGalleryUserSettings(ctx context.Context, args SetGalleryUserSettingsArgs) error {
	if args.Entries == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.Entries"}
	}
	routeValues := make(map[string]string)
	if args.UserScope == nil || *args.UserScope == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.UserScope"}
	}
	routeValues["userScope"] = *args.UserScope

	body, marshalErr := json.Marshal(*args.Entries)
	if marshalErr != nil {
		return marshalErr
	}
	locationId, _ := uuid.Parse("9b75ece3-7960-401c-848b-148ac01ca350")
	_, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the SetGalleryUserSettings function
type SetGalleryUserSettingsArgs struct {
	// (required) A key-value pair of all settings that need to be set
	Entries *map[string]interface{}
	// (required) User-Scope at which to get the value. Should be "me" for the current user or "host" for all users.
	UserScope *string
}

// [Preview API]
func (client *ClientImpl) ShareExtension(ctx context.Context, args ShareExtensionArgs) error {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.AccountName == nil || *args.AccountName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.AccountName"}
	}
	routeValues["accountName"] = *args.AccountName

	locationId, _ := uuid.Parse("a1e66d8f-f5de-4d16-8309-91a4e015ee46")
	_, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the ShareExtension function
type ShareExtensionArgs struct {
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
	// (required)
	AccountName *string
}

// [Preview API]
func (client *ClientImpl) ShareExtensionById(ctx context.Context, args ShareExtensionByIdArgs) error {
	routeValues := make(map[string]string)
	if args.ExtensionId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.ExtensionId"}
	}
	routeValues["extensionId"] = (*args.ExtensionId).String()
	if args.AccountName == nil || *args.AccountName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.AccountName"}
	}
	routeValues["accountName"] = *args.AccountName

	locationId, _ := uuid.Parse("1f19631b-a0b4-4a03-89c2-d79785d24360")
	_, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the ShareExtensionById function
type ShareExtensionByIdArgs struct {
	// (required)
	ExtensionId *uuid.UUID
	// (required)
	AccountName *string
}

// [Preview API]
func (client *ClientImpl) ShareExtensionWithHost(ctx context.Context, args ShareExtensionWithHostArgs) error {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.HostType == nil || *args.HostType == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.HostType"}
	}
	routeValues["hostType"] = *args.HostType
	if args.HostName == nil || *args.HostName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.HostName"}
	}
	routeValues["hostName"] = *args.HostName

	locationId, _ := uuid.Parse("328a3af8-d124-46e9-9483-01690cd415b9")
	_, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the ShareExtensionWithHost function
type ShareExtensionWithHostArgs struct {
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
	// (required)
	HostType *string
	// (required)
	HostName *string
}

// [Preview API]
func (client *ClientImpl) UnshareExtension(ctx context.Context, args UnshareExtensionArgs) error {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.AccountName == nil || *args.AccountName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.AccountName"}
	}
	routeValues["accountName"] = *args.AccountName

	locationId, _ := uuid.Parse("a1e66d8f-f5de-4d16-8309-91a4e015ee46")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the UnshareExtension function
type UnshareExtensionArgs struct {
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
	// (required)
	AccountName *string
}

// [Preview API]
func (client *ClientImpl) UnshareExtensionById(ctx context.Context, args UnshareExtensionByIdArgs) error {
	routeValues := make(map[string]string)
	if args.ExtensionId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.ExtensionId"}
	}
	routeValues["extensionId"] = (*args.ExtensionId).String()
	if args.AccountName == nil || *args.AccountName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.AccountName"}
	}
	routeValues["accountName"] = *args.AccountName

	locationId, _ := uuid.Parse("1f19631b-a0b4-4a03-89c2-d79785d24360")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the UnshareExtensionById function
type UnshareExtensionByIdArgs struct {
	// (required)
	ExtensionId *uuid.UUID
	// (required)
	AccountName *string
}

// [Preview API]
func (client *ClientImpl) UnshareExtensionWithHost(ctx context.Context, args UnshareExtensionWithHostArgs) error {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.HostType == nil || *args.HostType == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.HostType"}
	}
	routeValues["hostType"] = *args.HostType
	if args.HostName == nil || *args.HostName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.HostName"}
	}
	routeValues["hostName"] = *args.HostName

	locationId, _ := uuid.Parse("328a3af8-d124-46e9-9483-01690cd415b9")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the UnshareExtensionWithHost function
type UnshareExtensionWithHostArgs struct {
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
	// (required)
	HostType *string
	// (required)
	HostName *string
}

// [Preview API] REST endpoint to update an extension.
func (client *ClientImpl) UpdateExtension(ctx context.Context, args UpdateExtensionArgs) (*PublishedExtension, error) {
	if args.UploadStream == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UploadStream"}
	}
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
	if args.ExtensionType != nil {
		queryParams.Add("extensionType", *args.ExtensionType)
	}
	if args.ReCaptchaToken != nil {
		queryParams.Add("reCaptchaToken", *args.ReCaptchaToken)
	}
	if args.BypassScopeCheck != nil {
		queryParams.Add("bypassScopeCheck", strconv.FormatBool(*args.BypassScopeCheck))
	}
	locationId, _ := uuid.Parse("e11ea35a-16fe-4b80-ab11-c4cab88a0966")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.2", routeValues, queryParams, args.UploadStream, "application/octet-stream", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue PublishedExtension
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateExtension function
type UpdateExtensionArgs struct {
	// (required) Stream to upload
	UploadStream io.Reader
	// (required) Name of the publisher
	PublisherName *string
	// (required) Name of the extension
	ExtensionName *string
	// (optional)
	ExtensionType *string
	// (optional)
	ReCaptchaToken *string
	// (optional) This parameter decides if the scope change check needs to be invoked or not
	BypassScopeCheck *bool
}

// [Preview API]
func (client *ClientImpl) UpdateExtensionById(ctx context.Context, args UpdateExtensionByIdArgs) (*PublishedExtension, error) {
	routeValues := make(map[string]string)
	if args.ExtensionId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ExtensionId"}
	}
	routeValues["extensionId"] = (*args.ExtensionId).String()

	queryParams := url.Values{}
	if args.ReCaptchaToken != nil {
		queryParams.Add("reCaptchaToken", *args.ReCaptchaToken)
	}
	locationId, _ := uuid.Parse("a41192c8-9525-4b58-bc86-179fa549d80d")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.2", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue PublishedExtension
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateExtensionById function
type UpdateExtensionByIdArgs struct {
	// (required)
	ExtensionId *uuid.UUID
	// (optional)
	ReCaptchaToken *string
}

// [Preview API]
func (client *ClientImpl) UpdateExtensionProperties(ctx context.Context, args UpdateExtensionPropertiesArgs) (*PublishedExtension, error) {
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
	if args.Flags == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "flags"}
	}
	queryParams.Add("flags", string(*args.Flags))
	locationId, _ := uuid.Parse("e11ea35a-16fe-4b80-ab11-c4cab88a0966")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.2", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue PublishedExtension
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateExtensionProperties function
type UpdateExtensionPropertiesArgs struct {
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
	// (required)
	Flags *PublishedExtensionFlags
}

// [Preview API]
func (client *ClientImpl) UpdateExtensionStatistics(ctx context.Context, args UpdateExtensionStatisticsArgs) error {
	if args.ExtensionStatisticsUpdate == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.ExtensionStatisticsUpdate"}
	}
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName

	body, marshalErr := json.Marshal(*args.ExtensionStatisticsUpdate)
	if marshalErr != nil {
		return marshalErr
	}
	locationId, _ := uuid.Parse("a0ea3204-11e9-422d-a9ca-45851cc41400")
	_, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the UpdateExtensionStatistics function
type UpdateExtensionStatisticsArgs struct {
	// (required)
	ExtensionStatisticsUpdate *ExtensionStatisticUpdate
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
}

// [Preview API]
func (client *ClientImpl) UpdatePayloadInDraftForEditExtension(ctx context.Context, args UpdatePayloadInDraftForEditExtensionArgs) (*ExtensionDraft, error) {
	if args.UploadStream == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UploadStream"}
	}
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.DraftId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DraftId"}
	}
	routeValues["draftId"] = (*args.DraftId).String()

	additionalHeaders := make(map[string]string)
	if args.FileName != nil {
		additionalHeaders["X-Market-UploadFileName"] = *args.FileName
	}
	locationId, _ := uuid.Parse("02b33873-4e61-496e-83a2-59d1df46b7d8")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.1", routeValues, nil, args.UploadStream, "application/octet-stream", "application/json", additionalHeaders)
	if err != nil {
		return nil, err
	}

	var responseValue ExtensionDraft
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdatePayloadInDraftForEditExtension function
type UpdatePayloadInDraftForEditExtensionArgs struct {
	// (required) Stream to upload
	UploadStream io.Reader
	// (required)
	PublisherName *string
	// (required)
	ExtensionName *string
	// (required)
	DraftId *uuid.UUID
	// (optional) Header to pass the filename of the uploaded data
	FileName *string
}

// [Preview API]
func (client *ClientImpl) UpdatePayloadInDraftForNewExtension(ctx context.Context, args UpdatePayloadInDraftForNewExtensionArgs) (*ExtensionDraft, error) {
	if args.UploadStream == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UploadStream"}
	}
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.DraftId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.DraftId"}
	}
	routeValues["draftId"] = (*args.DraftId).String()

	additionalHeaders := make(map[string]string)
	if args.FileName != nil {
		additionalHeaders["X-Market-UploadFileName"] = *args.FileName
	}
	locationId, _ := uuid.Parse("b3ab127d-ebb9-4d22-b611-4e09593c8d79")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.1", routeValues, nil, args.UploadStream, "application/octet-stream", "application/json", additionalHeaders)
	if err != nil {
		return nil, err
	}

	var responseValue ExtensionDraft
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdatePayloadInDraftForNewExtension function
type UpdatePayloadInDraftForNewExtensionArgs struct {
	// (required) Stream to upload
	UploadStream io.Reader
	// (required)
	PublisherName *string
	// (required)
	DraftId *uuid.UUID
	// (optional) Header to pass the filename of the uploaded data
	FileName *string
}

// [Preview API]
func (client *ClientImpl) UpdatePublisher(ctx context.Context, args UpdatePublisherArgs) (*Publisher, error) {
	if args.Publisher == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Publisher"}
	}
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName

	body, marshalErr := json.Marshal(*args.Publisher)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("4ddec66a-e4f6-4f5d-999e-9e77710d7ff4")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Publisher
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdatePublisher function
type UpdatePublisherArgs struct {
	// (required)
	Publisher *Publisher
	// (required)
	PublisherName *string
}

// [Preview API] Update publisher asset like logo. It accepts asset file as an octet stream and file name is passed in header values.
func (client *ClientImpl) UpdatePublisherAsset(ctx context.Context, args UpdatePublisherAssetArgs) (*map[string]string, error) {
	if args.UploadStream == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UploadStream"}
	}
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName

	queryParams := url.Values{}
	if args.AssetType != nil {
		queryParams.Add("assetType", *args.AssetType)
	}
	additionalHeaders := make(map[string]string)
	if args.FileName != nil {
		additionalHeaders["X-Market-UploadFileName"] = *args.FileName
	}
	locationId, _ := uuid.Parse("21143299-34f9-4c62-8ca8-53da691192f9")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.1", routeValues, queryParams, args.UploadStream, "application/octet-stream", "application/json", additionalHeaders)
	if err != nil {
		return nil, err
	}

	var responseValue map[string]string
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdatePublisherAsset function
type UpdatePublisherAssetArgs struct {
	// (required) Stream to upload
	UploadStream io.Reader
	// (required) Internal name of the publisher
	PublisherName *string
	// (optional) Type of asset. Default value is 'logo'.
	AssetType *string
	// (optional) Header to pass the filename of the uploaded data
	FileName *string
}

// [Preview API] Endpoint to add/modify publisher membership. Currently Supports only addition/modification of 1 user at a time Works only for adding members of same tenant.
func (client *ClientImpl) UpdatePublisherMembers(ctx context.Context, args UpdatePublisherMembersArgs) (*[]PublisherRoleAssignment, error) {
	if args.RoleAssignments == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.RoleAssignments"}
	}
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName

	queryParams := url.Values{}
	if args.LimitToCallerIdentityDomain != nil {
		queryParams.Add("limitToCallerIdentityDomain", strconv.FormatBool(*args.LimitToCallerIdentityDomain))
	}
	body, marshalErr := json.Marshal(*args.RoleAssignments)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("4ddec66a-e4f6-4f5d-999e-9e77710d7ff4")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []PublisherRoleAssignment
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdatePublisherMembers function
type UpdatePublisherMembersArgs struct {
	// (required) List of user identifiers(email address) and role to be added. Currently only one entry is supported.
	RoleAssignments *[]PublisherUserRoleAssignmentRef
	// (required) The name/id of publisher to which users have to be added
	PublisherName *string
	// (optional) Should cross tenant addtions be allowed or not.
	LimitToCallerIdentityDomain *bool
}

// [Preview API] Updates an existing question for an extension.
func (client *ClientImpl) UpdateQuestion(ctx context.Context, args UpdateQuestionArgs) (*Question, error) {
	if args.Question == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Question"}
	}
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.QuestionId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.QuestionId"}
	}
	routeValues["questionId"] = strconv.FormatUint(*args.QuestionId, 10)

	body, marshalErr := json.Marshal(*args.Question)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("6d1d9741-eca8-4701-a3a5-235afc82dfa4")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Question
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateQuestion function
type UpdateQuestionArgs struct {
	// (required) Updated question to be set for the extension.
	Question *Question
	// (required) Name of the publisher who published the extension.
	PublisherName *string
	// (required) Name of the extension.
	ExtensionName *string
	// (required) Identifier of the question to be updated for the extension.
	QuestionId *uint64
}

// [Preview API] Updates an existing response for a given question for an extension.
func (client *ClientImpl) UpdateResponse(ctx context.Context, args UpdateResponseArgs) (*Response, error) {
	if args.Response == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Response"}
	}
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName
	if args.ExtensionName == nil || *args.ExtensionName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtensionName"}
	}
	routeValues["extensionName"] = *args.ExtensionName
	if args.QuestionId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.QuestionId"}
	}
	routeValues["questionId"] = strconv.FormatUint(*args.QuestionId, 10)
	if args.ResponseId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ResponseId"}
	}
	routeValues["responseId"] = strconv.FormatUint(*args.ResponseId, 10)

	body, marshalErr := json.Marshal(*args.Response)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("7f8ae5e0-46b0-438f-b2e8-13e8513517bd")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Response
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateResponse function
type UpdateResponseArgs struct {
	// (required) Updated response to be set for the extension.
	Response *Response
	// (required) Name of the publisher who published the extension.
	PublisherName *string
	// (required) Name of the extension.
	ExtensionName *string
	// (required) Identifier of the question for which response is to be updated for the extension.
	QuestionId *uint64
	// (required) Identifier of the response which has to be updated.
	ResponseId *uint64
}

// [Preview API] Updates or Flags a review
func (client *ClientImpl) UpdateReview(ctx context.Context, args UpdateReviewArgs) (*ReviewPatch, error) {
	if args.ReviewPatch == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ReviewPatch"}
	}
	routeValues := make(map[string]string)
	if args.PubName == nil || *args.PubName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PubName"}
	}
	routeValues["pubName"] = *args.PubName
	if args.ExtName == nil || *args.ExtName == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ExtName"}
	}
	routeValues["extName"] = *args.ExtName
	if args.ReviewId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.ReviewId"}
	}
	routeValues["reviewId"] = strconv.FormatUint(*args.ReviewId, 10)

	body, marshalErr := json.Marshal(*args.ReviewPatch)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("e6e85b9d-aa70-40e6-aa28-d0fbf40b91a3")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ReviewPatch
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateReview function
type UpdateReviewArgs struct {
	// (required) ReviewPatch object which contains the changes to be applied to the review
	ReviewPatch *ReviewPatch
	// (required) Name of the publisher who published the extension
	PubName *string
	// (required) Name of the extension
	ExtName *string
	// (required) Id of the review which needs to be updated
	ReviewId *uint64
}

// [Preview API]
func (client *ClientImpl) UpdateVSCodeWebExtensionStatistics(ctx context.Context, args UpdateVSCodeWebExtensionStatisticsArgs) error {
	routeValues := make(map[string]string)
	if args.ItemName == nil || *args.ItemName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ItemName"}
	}
	routeValues["itemName"] = *args.ItemName
	if args.Version == nil || *args.Version == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Version"}
	}
	routeValues["version"] = *args.Version
	if args.StatType == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.StatType"}
	}
	routeValues["statType"] = string(*args.StatType)

	locationId, _ := uuid.Parse("205c91a8-7841-4fd3-ae4f-5a745d5a8df5")
	_, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the UpdateVSCodeWebExtensionStatistics function
type UpdateVSCodeWebExtensionStatisticsArgs struct {
	// (required)
	ItemName *string
	// (required)
	Version *string
	// (required)
	StatType *VSCodeWebExtensionStatisicsType
}

// [Preview API]
func (client *ClientImpl) VerifyDomainToken(ctx context.Context, args VerifyDomainTokenArgs) error {
	routeValues := make(map[string]string)
	if args.PublisherName == nil || *args.PublisherName == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherName"}
	}
	routeValues["publisherName"] = *args.PublisherName

	locationId, _ := uuid.Parse("67a609ef-fa74-4b52-8664-78d76f7b3634")
	_, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the VerifyDomainToken function
type VerifyDomainTokenArgs struct {
	// (required)
	PublisherName *string
}
