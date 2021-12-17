// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"net/http"
	"net/url"
	"strings"
)

type Client interface {
	// Create a new subscription.
	CreateSubscription(context.Context, CreateSubscriptionArgs) (*NotificationSubscription, error)
	// Delete a subscription.
	DeleteSubscription(context.Context, DeleteSubscriptionArgs) error
	// Get a specific event type.
	GetEventType(context.Context, GetEventTypeArgs) (*NotificationEventType, error)
	GetSettings(context.Context, GetSettingsArgs) (*NotificationAdminSettings, error)
	// Get delivery preferences of a notifications subscriber.
	GetSubscriber(context.Context, GetSubscriberArgs) (*NotificationSubscriber, error)
	// Get a notification subscription by its ID.
	GetSubscription(context.Context, GetSubscriptionArgs) (*NotificationSubscription, error)
	// Get the diagnostics settings for a subscription.
	GetSubscriptionDiagnostics(context.Context, GetSubscriptionDiagnosticsArgs) (*SubscriptionDiagnostics, error)
	// Get available subscription templates.
	GetSubscriptionTemplates(context.Context, GetSubscriptionTemplatesArgs) (*[]NotificationSubscriptionTemplate, error)
	// List available event types for this service. Optionally filter by only event types for the specified publisher.
	ListEventTypes(context.Context, ListEventTypesArgs) (*[]NotificationEventType, error)
	// Get a list of diagnostic logs for this service.
	ListLogs(context.Context, ListLogsArgs) (*[]INotificationDiagnosticLog, error)
	// Get a list of notification subscriptions, either by subscription IDs or by all subscriptions for a given user or group.
	ListSubscriptions(context.Context, ListSubscriptionsArgs) (*[]NotificationSubscription, error)
	// Query for subscriptions. A subscription is returned if it matches one or more of the specified conditions.
	QuerySubscriptions(context.Context, QuerySubscriptionsArgs) (*[]NotificationSubscription, error)
	UpdateSettings(context.Context, UpdateSettingsArgs) (*NotificationAdminSettings, error)
	// Update delivery preferences of a notifications subscriber.
	UpdateSubscriber(context.Context, UpdateSubscriberArgs) (*NotificationSubscriber, error)
	// Update an existing subscription. Depending on the type of subscription and permissions, the caller can update the description, filter settings, channel (delivery) settings and more.
	UpdateSubscription(context.Context, UpdateSubscriptionArgs) (*NotificationSubscription, error)
	// Update the diagnostics settings for a subscription.
	UpdateSubscriptionDiagnostics(context.Context, UpdateSubscriptionDiagnosticsArgs) (*SubscriptionDiagnostics, error)
	// Update the specified user's settings for the specified subscription. This API is typically used to opt in or out of a shared subscription. User settings can only be applied to shared subscriptions, like team subscriptions or default subscriptions.
	UpdateSubscriptionUserSettings(context.Context, UpdateSubscriptionUserSettingsArgs) (*SubscriptionUserSettings, error)
}

type ClientImpl struct {
	Client azuredevops.Client
}

func NewClient(ctx context.Context, connection *azuredevops.Connection) Client {
	client := connection.GetClientByUrl(connection.BaseUrl)
	return &ClientImpl{
		Client: *client,
	}
}

// Create a new subscription.
func (client *ClientImpl) CreateSubscription(ctx context.Context, args CreateSubscriptionArgs) (*NotificationSubscription, error) {
	if args.CreateParameters == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.CreateParameters"}
	}
	body, marshalErr := json.Marshal(*args.CreateParameters)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("70f911d6-abac-488c-85b3-a206bf57e165")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "5.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue NotificationSubscription
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateSubscription function
type CreateSubscriptionArgs struct {
	// (required)
	CreateParameters *NotificationSubscriptionCreateParameters
}

// Delete a subscription.
func (client *ClientImpl) DeleteSubscription(ctx context.Context, args DeleteSubscriptionArgs) error {
	routeValues := make(map[string]string)
	if args.SubscriptionId == nil || *args.SubscriptionId == "" {
		return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.SubscriptionId"}
	}
	routeValues["subscriptionId"] = *args.SubscriptionId

	locationId, _ := uuid.Parse("70f911d6-abac-488c-85b3-a206bf57e165")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "5.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteSubscription function
type DeleteSubscriptionArgs struct {
	// (required)
	SubscriptionId *string
}

// Get a specific event type.
func (client *ClientImpl) GetEventType(ctx context.Context, args GetEventTypeArgs) (*NotificationEventType, error) {
	routeValues := make(map[string]string)
	if args.EventType == nil || *args.EventType == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.EventType"}
	}
	routeValues["eventType"] = *args.EventType

	locationId, _ := uuid.Parse("cc84fb5f-6247-4c7a-aeae-e5a3c3fddb21")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue NotificationEventType
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetEventType function
type GetEventTypeArgs struct {
	// (required) The ID of the event type.
	EventType *string
}

func (client *ClientImpl) GetSettings(ctx context.Context, args GetSettingsArgs) (*NotificationAdminSettings, error) {
	locationId, _ := uuid.Parse("cbe076d8-2803-45ff-8d8d-44653686ea2a")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1", nil, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue NotificationAdminSettings
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetSettings function
type GetSettingsArgs struct {
}

// Get delivery preferences of a notifications subscriber.
func (client *ClientImpl) GetSubscriber(ctx context.Context, args GetSubscriberArgs) (*NotificationSubscriber, error) {
	routeValues := make(map[string]string)
	if args.SubscriberId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.SubscriberId"}
	}
	routeValues["subscriberId"] = (*args.SubscriberId).String()

	locationId, _ := uuid.Parse("4d5caff1-25ba-430b-b808-7a1f352cc197")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue NotificationSubscriber
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetSubscriber function
type GetSubscriberArgs struct {
	// (required) ID of the user or group.
	SubscriberId *uuid.UUID
}

// Get a notification subscription by its ID.
func (client *ClientImpl) GetSubscription(ctx context.Context, args GetSubscriptionArgs) (*NotificationSubscription, error) {
	routeValues := make(map[string]string)
	if args.SubscriptionId == nil || *args.SubscriptionId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.SubscriptionId"}
	}
	routeValues["subscriptionId"] = *args.SubscriptionId

	queryParams := url.Values{}
	if args.QueryFlags != nil {
		queryParams.Add("queryFlags", string(*args.QueryFlags))
	}
	locationId, _ := uuid.Parse("70f911d6-abac-488c-85b3-a206bf57e165")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue NotificationSubscription
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetSubscription function
type GetSubscriptionArgs struct {
	// (required)
	SubscriptionId *string
	// (optional)
	QueryFlags *SubscriptionQueryFlags
}

// Get the diagnostics settings for a subscription.
func (client *ClientImpl) GetSubscriptionDiagnostics(ctx context.Context, args GetSubscriptionDiagnosticsArgs) (*SubscriptionDiagnostics, error) {
	routeValues := make(map[string]string)
	if args.SubscriptionId == nil || *args.SubscriptionId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.SubscriptionId"}
	}
	routeValues["subscriptionId"] = *args.SubscriptionId

	locationId, _ := uuid.Parse("20f1929d-4be7-4c2e-a74e-d47640ff3418")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue SubscriptionDiagnostics
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetSubscriptionDiagnostics function
type GetSubscriptionDiagnosticsArgs struct {
	// (required) The id of the notifications subscription.
	SubscriptionId *string
}

// Get available subscription templates.
func (client *ClientImpl) GetSubscriptionTemplates(ctx context.Context, args GetSubscriptionTemplatesArgs) (*[]NotificationSubscriptionTemplate, error) {
	locationId, _ := uuid.Parse("fa5d24ba-7484-4f3d-888d-4ec6b1974082")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1", nil, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []NotificationSubscriptionTemplate
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetSubscriptionTemplates function
type GetSubscriptionTemplatesArgs struct {
}

// List available event types for this service. Optionally filter by only event types for the specified publisher.
func (client *ClientImpl) ListEventTypes(ctx context.Context, args ListEventTypesArgs) (*[]NotificationEventType, error) {
	queryParams := url.Values{}
	if args.PublisherId != nil {
		queryParams.Add("publisherId", *args.PublisherId)
	}
	locationId, _ := uuid.Parse("cc84fb5f-6247-4c7a-aeae-e5a3c3fddb21")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1", nil, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []NotificationEventType
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ListEventTypes function
type ListEventTypesArgs struct {
	// (optional) Limit to event types for this publisher
	PublisherId *string
}

// Get a list of diagnostic logs for this service.
func (client *ClientImpl) ListLogs(ctx context.Context, args ListLogsArgs) (*[]INotificationDiagnosticLog, error) {
	routeValues := make(map[string]string)
	if args.Source == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Source"}
	}
	routeValues["source"] = (*args.Source).String()
	if args.EntryId != nil {
		routeValues["entryId"] = (*args.EntryId).String()
	}

	queryParams := url.Values{}
	if args.StartTime != nil {
		queryParams.Add("startTime", (*args.StartTime).AsQueryParameter())
	}
	if args.EndTime != nil {
		queryParams.Add("endTime", (*args.EndTime).AsQueryParameter())
	}
	locationId, _ := uuid.Parse("991842f3-eb16-4aea-ac81-81353ef2b75c")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []INotificationDiagnosticLog
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ListLogs function
type ListLogsArgs struct {
	// (required) ID specifying which type of logs to check diagnostics for.
	Source *uuid.UUID
	// (optional) The ID of the specific log to query for.
	EntryId *uuid.UUID
	// (optional) Start time for the time range to query in.
	StartTime *azuredevops.Time
	// (optional) End time for the time range to query in.
	EndTime *azuredevops.Time
}

// Get a list of notification subscriptions, either by subscription IDs or by all subscriptions for a given user or group.
func (client *ClientImpl) ListSubscriptions(ctx context.Context, args ListSubscriptionsArgs) (*[]NotificationSubscription, error) {
	queryParams := url.Values{}
	if args.TargetId != nil {
		queryParams.Add("targetId", (*args.TargetId).String())
	}
	if args.Ids != nil {
		listAsString := strings.Join((*args.Ids)[:], ",")
		queryParams.Add("ids", listAsString)
	}
	if args.QueryFlags != nil {
		queryParams.Add("queryFlags", string(*args.QueryFlags))
	}
	locationId, _ := uuid.Parse("70f911d6-abac-488c-85b3-a206bf57e165")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1", nil, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []NotificationSubscription
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ListSubscriptions function
type ListSubscriptionsArgs struct {
	// (optional) User or Group ID
	TargetId *uuid.UUID
	// (optional) List of subscription IDs
	Ids *[]string
	// (optional)
	QueryFlags *SubscriptionQueryFlags
}

// Query for subscriptions. A subscription is returned if it matches one or more of the specified conditions.
func (client *ClientImpl) QuerySubscriptions(ctx context.Context, args QuerySubscriptionsArgs) (*[]NotificationSubscription, error) {
	if args.SubscriptionQuery == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.SubscriptionQuery"}
	}
	body, marshalErr := json.Marshal(*args.SubscriptionQuery)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("6864db85-08c0-4006-8e8e-cc1bebe31675")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "5.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []NotificationSubscription
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QuerySubscriptions function
type QuerySubscriptionsArgs struct {
	// (required)
	SubscriptionQuery *SubscriptionQuery
}

func (client *ClientImpl) UpdateSettings(ctx context.Context, args UpdateSettingsArgs) (*NotificationAdminSettings, error) {
	if args.UpdateParameters == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UpdateParameters"}
	}
	body, marshalErr := json.Marshal(*args.UpdateParameters)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("cbe076d8-2803-45ff-8d8d-44653686ea2a")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "5.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue NotificationAdminSettings
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateSettings function
type UpdateSettingsArgs struct {
	// (required)
	UpdateParameters *NotificationAdminSettingsUpdateParameters
}

// Update delivery preferences of a notifications subscriber.
func (client *ClientImpl) UpdateSubscriber(ctx context.Context, args UpdateSubscriberArgs) (*NotificationSubscriber, error) {
	if args.UpdateParameters == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UpdateParameters"}
	}
	routeValues := make(map[string]string)
	if args.SubscriberId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.SubscriberId"}
	}
	routeValues["subscriberId"] = (*args.SubscriberId).String()

	body, marshalErr := json.Marshal(*args.UpdateParameters)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("4d5caff1-25ba-430b-b808-7a1f352cc197")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "5.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue NotificationSubscriber
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateSubscriber function
type UpdateSubscriberArgs struct {
	// (required)
	UpdateParameters *NotificationSubscriberUpdateParameters
	// (required) ID of the user or group.
	SubscriberId *uuid.UUID
}

// Update an existing subscription. Depending on the type of subscription and permissions, the caller can update the description, filter settings, channel (delivery) settings and more.
func (client *ClientImpl) UpdateSubscription(ctx context.Context, args UpdateSubscriptionArgs) (*NotificationSubscription, error) {
	if args.UpdateParameters == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UpdateParameters"}
	}
	routeValues := make(map[string]string)
	if args.SubscriptionId == nil || *args.SubscriptionId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.SubscriptionId"}
	}
	routeValues["subscriptionId"] = *args.SubscriptionId

	body, marshalErr := json.Marshal(*args.UpdateParameters)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("70f911d6-abac-488c-85b3-a206bf57e165")
	resp, err := client.Client.Send(ctx, http.MethodPatch, locationId, "5.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue NotificationSubscription
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateSubscription function
type UpdateSubscriptionArgs struct {
	// (required)
	UpdateParameters *NotificationSubscriptionUpdateParameters
	// (required)
	SubscriptionId *string
}

// Update the diagnostics settings for a subscription.
func (client *ClientImpl) UpdateSubscriptionDiagnostics(ctx context.Context, args UpdateSubscriptionDiagnosticsArgs) (*SubscriptionDiagnostics, error) {
	if args.UpdateParameters == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UpdateParameters"}
	}
	routeValues := make(map[string]string)
	if args.SubscriptionId == nil || *args.SubscriptionId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.SubscriptionId"}
	}
	routeValues["subscriptionId"] = *args.SubscriptionId

	body, marshalErr := json.Marshal(*args.UpdateParameters)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("20f1929d-4be7-4c2e-a74e-d47640ff3418")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "5.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue SubscriptionDiagnostics
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateSubscriptionDiagnostics function
type UpdateSubscriptionDiagnosticsArgs struct {
	// (required)
	UpdateParameters *UpdateSubscripitonDiagnosticsParameters
	// (required) The id of the notifications subscription.
	SubscriptionId *string
}

// Update the specified user's settings for the specified subscription. This API is typically used to opt in or out of a shared subscription. User settings can only be applied to shared subscriptions, like team subscriptions or default subscriptions.
func (client *ClientImpl) UpdateSubscriptionUserSettings(ctx context.Context, args UpdateSubscriptionUserSettingsArgs) (*SubscriptionUserSettings, error) {
	if args.UserSettings == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UserSettings"}
	}
	routeValues := make(map[string]string)
	if args.SubscriptionId == nil || *args.SubscriptionId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.SubscriptionId"}
	}
	routeValues["subscriptionId"] = *args.SubscriptionId
	if args.UserId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.UserId"}
	}
	routeValues["userId"] = (*args.UserId).String()

	body, marshalErr := json.Marshal(*args.UserSettings)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("ed5a3dff-aeb5-41b1-b4f7-89e66e58b62e")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "5.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue SubscriptionUserSettings
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateSubscriptionUserSettings function
type UpdateSubscriptionUserSettingsArgs struct {
	// (required)
	UserSettings *SubscriptionUserSettings
	// (required)
	SubscriptionId *string
	// (required) ID of the user
	UserId *uuid.UUID
}
