//go:build (all || resource_servicehook_storage_queue) && !exclude_subscriptions
// +build all resource_servicehook_storage_queue
// +build !exclude_subscriptions

package servicehook

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/servicehooks"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

var subscriptionStorageQueueID = uuid.New()

var testResourceSubscriptionStorageQueue = []servicehooks.Subscription{
	{
		Id:               &subscriptionStorageQueueID,
		ConsumerActionId: converter.String("enqueue"),
		ConsumerId:       converter.String("azureStorageQueue"),
		ConsumerInputs: &map[string]string{
			"accountKey":  "myaccountkey",
			"accountName": "myaccountname",
			"queueName":   "myqueue",
			"ttl":         "604800",
			"visiTimeout": "0",
		},
		EventType:   converter.String("ms.vss-pipelines.stage-state-changed-event"),
		PublisherId: converter.String("pipelines"),
		PublisherInputs: &map[string]string{
			"pipelineId":    "mypipelineid",
			"projectId":     "myprojectid",
			"stageNameId":   "mystagename",
			"stageStateId":  "mystagestatus",
			"stageResultId": "mystageresult",
		},
		ResourceVersion: converter.String("5.1-preview.1"),
	},
}

func TestResourceServicehookStorageQueue_FlattenExpandRoundTrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceServicehookStorageQueue().Schema, nil)
	for _, subscription := range testResourceSubscriptionStorageQueue {
		flattenServicehookStorageQueue(resourceData, &subscription, (*subscription.ConsumerInputs)["accountKey"])
		subscriptionAfterRoundTrip, _ := expandServicehookStorageQueue(resourceData)

		require.Equal(t, subscription, *subscriptionAfterRoundTrip)
	}
}

func TestResourceServicehookStorageQueue_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServicehookStorageQueue()
	for _, subscription := range testResourceSubscriptionStorageQueue {
		resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
		flattenServicehookStorageQueue(resourceData, &subscription, (*subscription.ConsumerInputs)["accountKey"])

		mockClient := azdosdkmocks.NewMockServicehooksClient(ctrl)
		clients := &client.AggregatedClient{ServiceHooksClient: mockClient, Ctx: context.Background()}

		expectedArgs := servicehooks.CreateSubscriptionArgs{Subscription: &subscription}

		mockClient.
			EXPECT().
			CreateSubscription(clients.Ctx, expectedArgs).
			Return(nil, errors.New("CreateSubscription() Failed")).
			Times(1)

		err := r.Create(resourceData, clients)
		require.Contains(t, err.Error(), "CreateSubscription() Failed")
	}
}

func TestResourceServicehookStorageQueue_Update_DoestNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServicehookStorageQueue()
	for _, subscription := range testResourceSubscriptionStorageQueue {
		resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
		flattenServicehookStorageQueue(resourceData, &subscription, (*subscription.ConsumerInputs)["accountKey"])

		mockClient := azdosdkmocks.NewMockServicehooksClient(ctrl)
		clients := &client.AggregatedClient{ServiceHooksClient: mockClient, Ctx: context.Background()}

		expectedArgs := servicehooks.ReplaceSubscriptionArgs{
			Subscription:   &subscription,
			SubscriptionId: subscription.Id,
		}

		mockClient.
			EXPECT().
			ReplaceSubscription(clients.Ctx, expectedArgs).
			Return(nil, errors.New("ReplaceSubscription() Failed")).
			Times(1)

		err := r.Update(resourceData, clients)
		require.Contains(t, err.Error(), "ReplaceSubscription() Failed")
	}
}

func TestResourceServicehookStorageQueue_Read_DoestNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServicehookStorageQueue()
	for _, subscription := range testResourceSubscriptionStorageQueue {
		resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
		flattenServicehookStorageQueue(resourceData, &subscription, (*subscription.ConsumerInputs)["accountKey"])

		mockClient := azdosdkmocks.NewMockServicehooksClient(ctrl)
		clients := &client.AggregatedClient{ServiceHooksClient: mockClient, Ctx: context.Background()}

		expectedArgs := servicehooks.GetSubscriptionArgs{SubscriptionId: subscription.Id}

		mockClient.
			EXPECT().
			GetSubscription(clients.Ctx, expectedArgs).
			Return(nil, errors.New("GetSubscription() Failed")).
			Times(1)

		err := r.Read(resourceData, clients)
		require.Contains(t, err.Error(), "GetSubscription() Failed")
	}
}

func TestResourceServicehookStorageQueue_Delete_DoestNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServicehookStorageQueue()
	for _, subscription := range testResourceSubscriptionStorageQueue {
		resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
		flattenServicehookStorageQueue(resourceData, &subscription, (*subscription.ConsumerInputs)["accountKey"])

		mockClient := azdosdkmocks.NewMockServicehooksClient(ctrl)
		clients := &client.AggregatedClient{ServiceHooksClient: mockClient, Ctx: context.Background()}

		expectedArgs := servicehooks.DeleteSubscriptionArgs{SubscriptionId: subscription.Id}

		mockClient.
			EXPECT().
			DeleteSubscription(clients.Ctx, expectedArgs).
			Return(errors.New("DeleteSubscription() Failed")).
			Times(1)

		err := r.Delete(resourceData, clients)
		require.Contains(t, err.Error(), "DeleteSubscription() Failed")
	}
}
