//go:build (all || resource_servicehook_subscription) && !exclude_subscriptions

package servicehook

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/servicehooks"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var subscriptionID = uuid.New()

var testResourceSubscription = []servicehooks.Subscription{
	{
		Id:               &subscriptionID,
		PublisherId:      converter.String("tfs"),
		EventType:        converter.String("git.push"),
		ConsumerId:       converter.String("webHooks"),
		ConsumerActionId: converter.String("httpRequest"),
		PublisherInputs: &map[string]string{
			"projectId":  "myprojectid",
			"repository": "myrepo",
		},
		ConsumerInputs: &map[string]string{
			"url": "https://example.com/webhook",
		},
		ResourceVersion: converter.String("1.0"),
		Status:          &servicehooks.SubscriptionStatusValues.Enabled,
	},
	{
		Id:               &subscriptionID,
		PublisherId:      converter.String("pipelines"),
		EventType:        converter.String("ms.vss-pipelines.run-state-changed-event"),
		ConsumerId:       converter.String("azureServiceBus"),
		ConsumerActionId: converter.String("serviceBusQueueMessage"),
		PublisherInputs: &map[string]string{
			"projectId":  "myprojectid",
			"pipelineId": "mypipelineid",
		},
		ConsumerInputs: &map[string]string{
			"connectionString": "Endpoint=sb://test.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=testkey",
			"queueName":        "myqueue",
		},
		ResourceVersion: converter.String("2.0-preview.1"),
		Status:          &servicehooks.SubscriptionStatusValues.Enabled,
	},
}

func TestServicehookSubscription_FlattenExpandRoundTrip(t *testing.T) {
	for _, subscription := range testResourceSubscription {
		resourceData := schema.TestResourceDataRaw(t, ResourceServicehookSubscription().Schema, nil)

		// First flatten the subscription
		err := flattenServicehookSubscription(resourceData, &subscription)
		require.NoError(t, err)

		// Set the consumer inputs (they aren't flattened for security)
		resourceData.Set("consumer_inputs", *subscription.ConsumerInputs)

		// Then expand it back
		subscriptionAfterRoundTrip := expandServicehookSubscription(resourceData)
		subscriptionAfterRoundTrip.Id = subscription.Id

		require.Equal(t, subscription, *subscriptionAfterRoundTrip)
	}
}

func TestServicehookSubscription_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServicehookSubscription()
	for _, subscription := range testResourceSubscription {
		resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
		err := flattenServicehookSubscription(resourceData, &subscription)
		require.NoError(t, err)
		resourceData.Set("consumer_inputs", *subscription.ConsumerInputs)

		mockClient := azdosdkmocks.NewMockServicehooksClient(ctrl)
		clients := &client.AggregatedClient{ServiceHooksClient: mockClient, Ctx: context.Background()}

		subscription.Id = nil
		expectedArgs := servicehooks.CreateSubscriptionArgs{Subscription: &subscription}

		mockClient.
			EXPECT().
			CreateSubscription(clients.Ctx, expectedArgs).
			Return(nil, errors.New("CreateSubscription() Failed")).
			Times(1)

		err = r.Create(resourceData, clients)
		require.Contains(t, err.Error(), "CreateSubscription() Failed")
	}
}

func TestServicehookSubscription_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServicehookSubscription()
	for _, subscription := range testResourceSubscription {
		resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
		resourceData.SetId(subscription.Id.String())
		err := flattenServicehookSubscription(resourceData, &subscription)
		require.NoError(t, err)
		resourceData.Set("consumer_inputs", *subscription.ConsumerInputs)

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

		err = r.Update(resourceData, clients)
		require.Contains(t, err.Error(), "ReplaceSubscription() Failed")
	}
}

func TestServicehookSubscription_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServicehookSubscription()
	for _, subscription := range testResourceSubscription {
		resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
		resourceData.SetId(subscription.Id.String())

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

func TestServicehookSubscription_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServicehookSubscription()
	for _, subscription := range testResourceSubscription {
		resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
		resourceData.SetId(subscription.Id.String())

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

func TestServicehookSubscription_StatusConversion(t *testing.T) {
	tests := []struct {
		input    string
		expected servicehooks.SubscriptionStatus
	}{
		{"enabled", servicehooks.SubscriptionStatusValues.Enabled},
		{"disabled", servicehooks.SubscriptionStatusValues.DisabledByUser},
		{"disabledByUser", servicehooks.SubscriptionStatusValues.DisabledByUser},
		{"disabledBySystem", servicehooks.SubscriptionStatusValues.DisabledBySystem},
		{"onProbation", servicehooks.SubscriptionStatusValues.OnProbation},
		{"invalid", servicehooks.SubscriptionStatusValues.Enabled}, // default
	}

	for _, test := range tests {
		result := convertStatus(test.input)
		require.Equal(t, test.expected, result)
	}
}

func TestServicehookSubscription_StatusConversionFromAPI(t *testing.T) {
	tests := []struct {
		input    servicehooks.SubscriptionStatus
		expected string
	}{
		{servicehooks.SubscriptionStatusValues.Enabled, "enabled"},
		{servicehooks.SubscriptionStatusValues.DisabledByUser, "disabledByUser"},
		{servicehooks.SubscriptionStatusValues.DisabledBySystem, "disabledBySystem"},
		{servicehooks.SubscriptionStatusValues.OnProbation, "onProbation"},
	}

	for _, test := range tests {
		result := convertStatusFromAPI(test.input)
		require.Equal(t, test.expected, result)
	}
}
