// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package servicehooks

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/forminput"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/notification"
	"net/http"
	"net/url"
	"strconv"
)

type Client interface {
	// [Preview API] Create a subscription.
	CreateSubscription(context.Context, CreateSubscriptionArgs) (*Subscription, error)
	// [Preview API] Query for service hook subscriptions.
	CreateSubscriptionsQuery(context.Context, CreateSubscriptionsQueryArgs) (*SubscriptionsQuery, error)
	// [Preview API] Sends a test notification. This is useful for verifying the configuration of an updated or new service hooks subscription.
	CreateTestNotification(context.Context, CreateTestNotificationArgs) (*Notification, error)
	// [Preview API] Delete a specific service hooks subscription.
	DeleteSubscription(context.Context, DeleteSubscriptionArgs) error
	// [Preview API] Get a specific consumer service. Optionally filter out consumer actions that do not support any event types for the specified publisher.
	GetConsumer(context.Context, GetConsumerArgs) (*Consumer, error)
	// [Preview API] Get details about a specific consumer action.
	GetConsumerAction(context.Context, GetConsumerActionArgs) (*ConsumerAction, error)
	// [Preview API] Get a specific event type.
	GetEventType(context.Context, GetEventTypeArgs) (*EventTypeDescriptor, error)
	// [Preview API] Get a specific notification for a subscription.
	GetNotification(context.Context, GetNotificationArgs) (*Notification, error)
	// [Preview API] Get a list of notifications for a specific subscription. A notification includes details about the event, the request to and the response from the consumer service.
	GetNotifications(context.Context, GetNotificationsArgs) (*[]Notification, error)
	// [Preview API] Get a specific service hooks publisher.
	GetPublisher(context.Context, GetPublisherArgs) (*Publisher, error)
	// [Preview API] Get a specific service hooks subscription.
	GetSubscription(context.Context, GetSubscriptionArgs) (*Subscription, error)
	// [Preview API]
	GetSubscriptionDiagnostics(context.Context, GetSubscriptionDiagnosticsArgs) (*notification.SubscriptionDiagnostics, error)
	// [Preview API] Get a list of consumer actions for a specific consumer.
	ListConsumerActions(context.Context, ListConsumerActionsArgs) (*[]ConsumerAction, error)
	// [Preview API] Get a list of available service hook consumer services. Optionally filter by consumers that support at least one event type from the specific publisher.
	ListConsumers(context.Context, ListConsumersArgs) (*[]Consumer, error)
	// [Preview API] Get the event types for a specific publisher.
	ListEventTypes(context.Context, ListEventTypesArgs) (*[]EventTypeDescriptor, error)
	// [Preview API] Get a list of publishers.
	ListPublishers(context.Context, ListPublishersArgs) (*[]Publisher, error)
	// [Preview API] Get a list of subscriptions.
	ListSubscriptions(context.Context, ListSubscriptionsArgs) (*[]Subscription, error)
	// [Preview API]
	QueryInputValues(context.Context, QueryInputValuesArgs) (*forminput.InputValuesQuery, error)
	// [Preview API] Query for notifications. A notification includes details about the event, the request to and the response from the consumer service.
	QueryNotifications(context.Context, QueryNotificationsArgs) (*NotificationsQuery, error)
	// [Preview API] Query for service hook publishers.
	QueryPublishers(context.Context, QueryPublishersArgs) (*PublishersQuery, error)
	// [Preview API] Update a subscription. <param name="subscriptionId">ID for a subscription that you wish to update.</param>
	ReplaceSubscription(context.Context, ReplaceSubscriptionArgs) (*Subscription, error)
	// [Preview API]
	UpdateSubscriptionDiagnostics(context.Context, UpdateSubscriptionDiagnosticsArgs) (*notification.SubscriptionDiagnostics, error)
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

// [Preview API] Create a subscription.
func (client *ClientImpl) CreateSubscription(ctx context.Context, args CreateSubscriptionArgs) (*Subscription, error) {
	if args.Subscription == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Subscription"}
	}
	body, marshalErr := json.Marshal(*args.Subscription)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("fc50d02a-849f-41fb-8af1-0a5216103269")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Subscription
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateSubscription function
type CreateSubscriptionArgs struct {
	// (required) Subscription to be created.
	Subscription *Subscription
}

// [Preview API] Query for service hook subscriptions.
func (client *ClientImpl) CreateSubscriptionsQuery(ctx context.Context, args CreateSubscriptionsQueryArgs) (*SubscriptionsQuery, error) {
	if args.Query == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Query"}
	}
	body, marshalErr := json.Marshal(*args.Query)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("c7c3c1cf-9e05-4c0d-a425-a0f922c2c6ed")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue SubscriptionsQuery
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateSubscriptionsQuery function
type CreateSubscriptionsQueryArgs struct {
	// (required)
	Query *SubscriptionsQuery
}

// [Preview API] Sends a test notification. This is useful for verifying the configuration of an updated or new service hooks subscription.
func (client *ClientImpl) CreateTestNotification(ctx context.Context, args CreateTestNotificationArgs) (*Notification, error) {
	if args.TestNotification == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.TestNotification"}
	}
	queryParams := url.Values{}
	if args.UseRealData != nil {
		queryParams.Add("useRealData", strconv.FormatBool(*args.UseRealData))
	}
	body, marshalErr := json.Marshal(*args.TestNotification)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("1139462c-7e27-4524-a997-31b9b73551fe")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", nil, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Notification
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateTestNotification function
type CreateTestNotificationArgs struct {
	// (required)
	TestNotification *Notification
	// (optional) Only allow testing with real data in existing subscriptions.
	UseRealData *bool
}

// [Preview API] Delete a specific service hooks subscription.
func (client *ClientImpl) DeleteSubscription(ctx context.Context, args DeleteSubscriptionArgs) error {
	routeValues := make(map[string]string)
	if args.SubscriptionId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.SubscriptionId"}
	}
	routeValues["subscriptionId"] = (*args.SubscriptionId).String()

	locationId, _ := uuid.Parse("fc50d02a-849f-41fb-8af1-0a5216103269")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteSubscription function
type DeleteSubscriptionArgs struct {
	// (required) ID for a subscription.
	SubscriptionId *uuid.UUID
}

// [Preview API] Get a specific consumer service. Optionally filter out consumer actions that do not support any event types for the specified publisher.
func (client *ClientImpl) GetConsumer(ctx context.Context, args GetConsumerArgs) (*Consumer, error) {
	routeValues := make(map[string]string)
	if args.ConsumerId == nil || *args.ConsumerId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ConsumerId"}
	}
	routeValues["consumerId"] = *args.ConsumerId

	queryParams := url.Values{}
	if args.PublisherId != nil {
		queryParams.Add("publisherId", *args.PublisherId)
	}
	locationId, _ := uuid.Parse("4301c514-5f34-4f5d-a145-f0ea7b5b7d19")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Consumer
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetConsumer function
type GetConsumerArgs struct {
	// (required) ID for a consumer.
	ConsumerId *string
	// (optional)
	PublisherId *string
}

// [Preview API] Get details about a specific consumer action.
func (client *ClientImpl) GetConsumerAction(ctx context.Context, args GetConsumerActionArgs) (*ConsumerAction, error) {
	routeValues := make(map[string]string)
	if args.ConsumerId == nil || *args.ConsumerId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ConsumerId"}
	}
	routeValues["consumerId"] = *args.ConsumerId
	if args.ConsumerActionId == nil || *args.ConsumerActionId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ConsumerActionId"}
	}
	routeValues["consumerActionId"] = *args.ConsumerActionId

	queryParams := url.Values{}
	if args.PublisherId != nil {
		queryParams.Add("publisherId", *args.PublisherId)
	}
	locationId, _ := uuid.Parse("c3428e90-7a69-4194-8ed8-0f153185ee0d")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue ConsumerAction
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetConsumerAction function
type GetConsumerActionArgs struct {
	// (required) ID for a consumer.
	ConsumerId *string
	// (required) ID for a consumerActionId.
	ConsumerActionId *string
	// (optional)
	PublisherId *string
}

// [Preview API] Get a specific event type.
func (client *ClientImpl) GetEventType(ctx context.Context, args GetEventTypeArgs) (*EventTypeDescriptor, error) {
	routeValues := make(map[string]string)
	if args.PublisherId == nil || *args.PublisherId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherId"}
	}
	routeValues["publisherId"] = *args.PublisherId
	if args.EventTypeId == nil || *args.EventTypeId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.EventTypeId"}
	}
	routeValues["eventTypeId"] = *args.EventTypeId

	locationId, _ := uuid.Parse("db4777cd-8e08-4a84-8ba3-c974ea033718")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue EventTypeDescriptor
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetEventType function
type GetEventTypeArgs struct {
	// (required) ID for a publisher.
	PublisherId *string
	// (required)
	EventTypeId *string
}

// [Preview API] Get a specific notification for a subscription.
func (client *ClientImpl) GetNotification(ctx context.Context, args GetNotificationArgs) (*Notification, error) {
	routeValues := make(map[string]string)
	if args.SubscriptionId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.SubscriptionId"}
	}
	routeValues["subscriptionId"] = (*args.SubscriptionId).String()
	if args.NotificationId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.NotificationId"}
	}
	routeValues["notificationId"] = strconv.Itoa(*args.NotificationId)

	locationId, _ := uuid.Parse("0c62d343-21b0-4732-997b-017fde84dc28")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Notification
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetNotification function
type GetNotificationArgs struct {
	// (required) ID for a subscription.
	SubscriptionId *uuid.UUID
	// (required)
	NotificationId *int
}

// [Preview API] Get a list of notifications for a specific subscription. A notification includes details about the event, the request to and the response from the consumer service.
func (client *ClientImpl) GetNotifications(ctx context.Context, args GetNotificationsArgs) (*[]Notification, error) {
	routeValues := make(map[string]string)
	if args.SubscriptionId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.SubscriptionId"}
	}
	routeValues["subscriptionId"] = (*args.SubscriptionId).String()

	queryParams := url.Values{}
	if args.MaxResults != nil {
		queryParams.Add("maxResults", strconv.Itoa(*args.MaxResults))
	}
	if args.Status != nil {
		queryParams.Add("status", string(*args.Status))
	}
	if args.Result != nil {
		queryParams.Add("result", string(*args.Result))
	}
	locationId, _ := uuid.Parse("0c62d343-21b0-4732-997b-017fde84dc28")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []Notification
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetNotifications function
type GetNotificationsArgs struct {
	// (required) ID for a subscription.
	SubscriptionId *uuid.UUID
	// (optional) Maximum number of notifications to return. Default is **100**.
	MaxResults *int
	// (optional) Get only notifications with this status.
	Status *NotificationStatus
	// (optional) Get only notifications with this result type.
	Result *NotificationResult
}

// [Preview API] Get a specific service hooks publisher.
func (client *ClientImpl) GetPublisher(ctx context.Context, args GetPublisherArgs) (*Publisher, error) {
	routeValues := make(map[string]string)
	if args.PublisherId == nil || *args.PublisherId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherId"}
	}
	routeValues["publisherId"] = *args.PublisherId

	locationId, _ := uuid.Parse("1e83a210-5b53-43bc-90f0-d476a4e5d731")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Publisher
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetPublisher function
type GetPublisherArgs struct {
	// (required) ID for a publisher.
	PublisherId *string
}

// [Preview API] Get a specific service hooks subscription.
func (client *ClientImpl) GetSubscription(ctx context.Context, args GetSubscriptionArgs) (*Subscription, error) {
	routeValues := make(map[string]string)
	if args.SubscriptionId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.SubscriptionId"}
	}
	routeValues["subscriptionId"] = (*args.SubscriptionId).String()

	locationId, _ := uuid.Parse("fc50d02a-849f-41fb-8af1-0a5216103269")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Subscription
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetSubscription function
type GetSubscriptionArgs struct {
	// (required) ID for a subscription.
	SubscriptionId *uuid.UUID
}

// [Preview API]
func (client *ClientImpl) GetSubscriptionDiagnostics(ctx context.Context, args GetSubscriptionDiagnosticsArgs) (*notification.SubscriptionDiagnostics, error) {
	routeValues := make(map[string]string)
	if args.SubscriptionId == nil || *args.SubscriptionId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.SubscriptionId"}
	}
	routeValues["subscriptionId"] = *args.SubscriptionId

	locationId, _ := uuid.Parse("3b36bcb5-02ad-43c6-bbfa-6dfc6f8e9d68")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue notification.SubscriptionDiagnostics
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetSubscriptionDiagnostics function
type GetSubscriptionDiagnosticsArgs struct {
	// (required)
	SubscriptionId *string
}

// [Preview API] Get a list of consumer actions for a specific consumer.
func (client *ClientImpl) ListConsumerActions(ctx context.Context, args ListConsumerActionsArgs) (*[]ConsumerAction, error) {
	routeValues := make(map[string]string)
	if args.ConsumerId == nil || *args.ConsumerId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.ConsumerId"}
	}
	routeValues["consumerId"] = *args.ConsumerId

	queryParams := url.Values{}
	if args.PublisherId != nil {
		queryParams.Add("publisherId", *args.PublisherId)
	}
	locationId, _ := uuid.Parse("c3428e90-7a69-4194-8ed8-0f153185ee0d")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []ConsumerAction
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ListConsumerActions function
type ListConsumerActionsArgs struct {
	// (required) ID for a consumer.
	ConsumerId *string
	// (optional)
	PublisherId *string
}

// [Preview API] Get a list of available service hook consumer services. Optionally filter by consumers that support at least one event type from the specific publisher.
func (client *ClientImpl) ListConsumers(ctx context.Context, args ListConsumersArgs) (*[]Consumer, error) {
	queryParams := url.Values{}
	if args.PublisherId != nil {
		queryParams.Add("publisherId", *args.PublisherId)
	}
	locationId, _ := uuid.Parse("4301c514-5f34-4f5d-a145-f0ea7b5b7d19")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", nil, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []Consumer
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ListConsumers function
type ListConsumersArgs struct {
	// (optional)
	PublisherId *string
}

// [Preview API] Get the event types for a specific publisher.
func (client *ClientImpl) ListEventTypes(ctx context.Context, args ListEventTypesArgs) (*[]EventTypeDescriptor, error) {
	routeValues := make(map[string]string)
	if args.PublisherId == nil || *args.PublisherId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherId"}
	}
	routeValues["publisherId"] = *args.PublisherId

	locationId, _ := uuid.Parse("db4777cd-8e08-4a84-8ba3-c974ea033718")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []EventTypeDescriptor
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ListEventTypes function
type ListEventTypesArgs struct {
	// (required) ID for a publisher.
	PublisherId *string
}

// [Preview API] Get a list of publishers.
func (client *ClientImpl) ListPublishers(ctx context.Context, args ListPublishersArgs) (*[]Publisher, error) {
	locationId, _ := uuid.Parse("1e83a210-5b53-43bc-90f0-d476a4e5d731")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", nil, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []Publisher
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ListPublishers function
type ListPublishersArgs struct {
}

// [Preview API] Get a list of subscriptions.
func (client *ClientImpl) ListSubscriptions(ctx context.Context, args ListSubscriptionsArgs) (*[]Subscription, error) {
	queryParams := url.Values{}
	if args.PublisherId != nil {
		queryParams.Add("publisherId", *args.PublisherId)
	}
	if args.EventType != nil {
		queryParams.Add("eventType", *args.EventType)
	}
	if args.ConsumerId != nil {
		queryParams.Add("consumerId", *args.ConsumerId)
	}
	if args.ConsumerActionId != nil {
		queryParams.Add("consumerActionId", *args.ConsumerActionId)
	}
	locationId, _ := uuid.Parse("fc50d02a-849f-41fb-8af1-0a5216103269")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "7.1-preview.1", nil, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []Subscription
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ListSubscriptions function
type ListSubscriptionsArgs struct {
	// (optional) ID for a subscription.
	PublisherId *string
	// (optional) The event type to filter on (if any).
	EventType *string
	// (optional) ID for a consumer.
	ConsumerId *string
	// (optional) ID for a consumerActionId.
	ConsumerActionId *string
}

// [Preview API]
func (client *ClientImpl) QueryInputValues(ctx context.Context, args QueryInputValuesArgs) (*forminput.InputValuesQuery, error) {
	if args.InputValuesQuery == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.InputValuesQuery"}
	}
	routeValues := make(map[string]string)
	if args.PublisherId == nil || *args.PublisherId == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.PublisherId"}
	}
	routeValues["publisherId"] = *args.PublisherId

	body, marshalErr := json.Marshal(*args.InputValuesQuery)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("d815d352-a566-4dc1-a3e3-fd245acf688c")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue forminput.InputValuesQuery
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryInputValues function
type QueryInputValuesArgs struct {
	// (required)
	InputValuesQuery *forminput.InputValuesQuery
	// (required)
	PublisherId *string
}

// [Preview API] Query for notifications. A notification includes details about the event, the request to and the response from the consumer service.
func (client *ClientImpl) QueryNotifications(ctx context.Context, args QueryNotificationsArgs) (*NotificationsQuery, error) {
	if args.Query == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Query"}
	}
	body, marshalErr := json.Marshal(*args.Query)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("1a57562f-160a-4b5c-9185-905e95b39d36")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue NotificationsQuery
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryNotifications function
type QueryNotificationsArgs struct {
	// (required)
	Query *NotificationsQuery
}

// [Preview API] Query for service hook publishers.
func (client *ClientImpl) QueryPublishers(ctx context.Context, args QueryPublishersArgs) (*PublishersQuery, error) {
	if args.Query == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Query"}
	}
	body, marshalErr := json.Marshal(*args.Query)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("99b44a8a-65a8-4670-8f3e-e7f7842cce64")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "7.1-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue PublishersQuery
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryPublishers function
type QueryPublishersArgs struct {
	// (required)
	Query *PublishersQuery
}

// [Preview API] Update a subscription. <param name="subscriptionId">ID for a subscription that you wish to update.</param>
func (client *ClientImpl) ReplaceSubscription(ctx context.Context, args ReplaceSubscriptionArgs) (*Subscription, error) {
	if args.Subscription == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Subscription"}
	}
	routeValues := make(map[string]string)
	if args.SubscriptionId != nil {
		routeValues["subscriptionId"] = (*args.SubscriptionId).String()
	}

	body, marshalErr := json.Marshal(*args.Subscription)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("fc50d02a-849f-41fb-8af1-0a5216103269")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue Subscription
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the ReplaceSubscription function
type ReplaceSubscriptionArgs struct {
	// (required)
	Subscription *Subscription
	// (optional)
	SubscriptionId *uuid.UUID
}

// [Preview API]
func (client *ClientImpl) UpdateSubscriptionDiagnostics(ctx context.Context, args UpdateSubscriptionDiagnosticsArgs) (*notification.SubscriptionDiagnostics, error) {
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
	locationId, _ := uuid.Parse("3b36bcb5-02ad-43c6-bbfa-6dfc6f8e9d68")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue notification.SubscriptionDiagnostics
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateSubscriptionDiagnostics function
type UpdateSubscriptionDiagnosticsArgs struct {
	// (required)
	UpdateParameters *notification.UpdateSubscripitonDiagnosticsParameters
	// (required)
	SubscriptionId *string
}
