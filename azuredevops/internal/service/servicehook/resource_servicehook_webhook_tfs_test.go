//go:build (all || resource_servicehook_webhook) && !exclude_subscriptions
// +build all resource_servicehook_webhook
// +build !exclude_subscriptions

package servicehook

import (
	"context"
	"errors"
	"strings"
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

var subscriptionWebhookID = uuid.New()
var createdByID = uuid.New()
var modifiedByID = uuid.New()
var subscriberID = uuid.New()

var enabledStatus = servicehooks.SubscriptionStatus("enabled")
var onProbationStatus = servicehooks.SubscriptionStatus("onProbation")

var testResourceSubscriptionWebhookTfs = []servicehooks.Subscription{
	{
		Id:               &subscriptionWebhookID,
		ConsumerActionId: converter.String("httpRequest"),
		ConsumerId:       converter.String("webHooks"),
		ConsumerInputs: &map[string]string{
			"url":                    "https://example.com/webhook",
			"acceptUntrustedCerts":   "false",
			"resourceDetailsToSend":  "all",
			"messagesToSend":         "all",
			"detailedMessagesToSend": "all",
		},
		EventType:   converter.String("git.push"),
		PublisherId: converter.String("tfs"),
		PublisherInputs: &map[string]string{
			"projectId":  "myprojectid",
			"repository": "myrepositoryid",
			"branch":     "main",
			"pushedBy":   "myuser",
		},
		ResourceVersion: converter.String("latest"),
	},
	{
		Id:               &subscriptionWebhookID,
		ConsumerActionId: converter.String("httpRequest"),
		ConsumerId:       converter.String("webHooks"),
		ConsumerInputs: &map[string]string{
			"url":                    "https://example.com/webhook",
			"acceptUntrustedCerts":   "true",
			"basicAuthPassword":      "user:pass",
			"httpHeaders":            "X-Custom-Header:CustomValue\nAuthorization:Bearer token123",
			"resourceDetailsToSend":  "minimal",
			"messagesToSend":         "all",
			"detailedMessagesToSend": "all",
		},
		EventType:   converter.String("git.pullrequest.created"),
		PublisherId: converter.String("tfs"),
		PublisherInputs: &map[string]string{
			"projectId":                    "myprojectid",
			"repository":                   "myrepositoryid",
			"branch":                       "develop",
			"pullrequestCreatedBy":         "myuser",
			"pullrequestReviewersContains": "reviewergroup",
		},
		ResourceVersion: converter.String("latest"),
	},
	{
		Id:               &subscriptionWebhookID,
		ConsumerActionId: converter.String("httpRequest"),
		ConsumerId:       converter.String("webHooks"),
		ConsumerInputs: &map[string]string{
			"url":                    "https://example.com/webhook",
			"acceptUntrustedCerts":   "false",
			"resourceDetailsToSend":  "none",
			"messagesToSend":         "none",
			"detailedMessagesToSend": "none",
		},
		EventType:   converter.String("workitem.created"),
		PublisherId: converter.String("tfs"),
		PublisherInputs: &map[string]string{
			"projectId":    "myprojectid",
			"workItemType": "Bug",
			"areaPath":     "MyProject\\MyArea",
			"tag":          "urgent",
		},
		ResourceVersion: converter.String("latest"),
	},
	{
		Id:               &subscriptionWebhookID,
		ConsumerActionId: converter.String("httpRequest"),
		ConsumerId:       converter.String("webHooks"),
		ConsumerInputs: &map[string]string{
			"url":                    "https://example.com/webhook",
			"acceptUntrustedCerts":   "false",
			"resourceDetailsToSend":  "all",
			"messagesToSend":         "none",
			"detailedMessagesToSend": "none",
		},
		EventType:   converter.String("build.complete"),
		PublisherId: converter.String("tfs"),
		PublisherInputs: &map[string]string{
			"projectId":      "myprojectid",
			"definitionName": "MyBuildDefinition",
			"buildStatus":    "Succeeded",
		},
		ResourceVersion: converter.String("latest"),
	},
}

func TestServicehookWebhookTfs_FlattenExpandRoundTrip(t *testing.T) {
	for _, subscription := range testResourceSubscriptionWebhookTfs {
		resourceData := schema.TestResourceDataRaw(t, ResourceServicehookWebhookTfs().Schema, nil)
		flattenServicehookWebhookTfs(resourceData, &subscription)

		// For subscriptions with basic auth, we need to manually set the credentials
		// since the flatten function doesn't extract them for security reasons
		if basicAuthUsername, exists := (*subscription.ConsumerInputs)["basicAuthUsername"]; exists {
			resourceData.Set("basic_auth_username", basicAuthUsername)
		}
		if basicAuthPassword, exists := (*subscription.ConsumerInputs)["basicAuthPassword"]; exists {
			resourceData.Set("basic_auth_password", basicAuthPassword)
		}

		subscriptionAfterRoundTrip := expandServicehookWebhookTfs(resourceData)
		subscriptionAfterRoundTrip.Id = subscription.Id

		// Compare everything except ConsumerInputs first
		subscriptionCopy := subscription
		subscriptionAfterRoundTripCopy := *subscriptionAfterRoundTrip
		subscriptionCopy.ConsumerInputs = nil
		subscriptionAfterRoundTripCopy.ConsumerInputs = nil
		require.Equal(t, subscriptionCopy, subscriptionAfterRoundTripCopy)

		// Compare ConsumerInputs separately, handling http headers specially
		if subscription.ConsumerInputs != nil && subscriptionAfterRoundTrip.ConsumerInputs != nil {
			originalInputs := *subscription.ConsumerInputs
			roundTripInputs := *subscriptionAfterRoundTrip.ConsumerInputs

			// Compare all inputs except httpHeaders
			for key, value := range originalInputs {
				if key != "httpHeaders" {
					require.Equal(t, value, roundTripInputs[key], "ConsumerInput %s should match", key)
				}
			}

			// Handle httpHeaders specially - check that all original headers are present
			if originalHeaders, exists := originalInputs["httpHeaders"]; exists {
				roundTripHeaders, rtExists := roundTripInputs["httpHeaders"]
				require.True(t, rtExists, "httpHeaders should exist in round trip result")

				// Split headers by newline and check each one is present
				originalHeaderLines := strings.Split(originalHeaders, "\n")
				roundTripHeaderLines := strings.Split(roundTripHeaders, "\n")

				for _, originalHeader := range originalHeaderLines {
					if strings.TrimSpace(originalHeader) != "" {
						require.Contains(t, roundTripHeaderLines, originalHeader, "Header %s should be present after round trip", originalHeader)
					}
				}
			}
		}
	}
}

func TestServicehookWebhookTfs_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServicehookWebhookTfs()
	for _, subscription := range testResourceSubscriptionWebhookTfs {
		resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
		flattenServicehookWebhookTfs(resourceData, &subscription)

		mockClient := azdosdkmocks.NewMockServicehooksClient(ctrl)
		clients := &client.AggregatedClient{ServiceHooksClient: mockClient, Ctx: context.Background()}
		subscription.Id = nil

		mockClient.
			EXPECT().
			CreateSubscription(clients.Ctx, gomock.AssignableToTypeOf(servicehooks.CreateSubscriptionArgs{})).
			Return(nil, errors.New("CreateSubscription() Failed")).
			Times(1)

		err := r.Create(resourceData, clients)
		require.Contains(t, err.Error(), "CreateSubscription() Failed")
	}
}

func TestServicehookWebhookTfs_Update_DoestNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServicehookWebhookTfs()
	for _, subscription := range testResourceSubscriptionWebhookTfs {
		resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
		resourceData.SetId(subscription.Id.String())
		flattenServicehookWebhookTfs(resourceData, &subscription)

		mockClient := azdosdkmocks.NewMockServicehooksClient(ctrl)
		clients := &client.AggregatedClient{ServiceHooksClient: mockClient, Ctx: context.Background()}

		mockClient.
			EXPECT().
			ReplaceSubscription(clients.Ctx, gomock.AssignableToTypeOf(servicehooks.ReplaceSubscriptionArgs{})).
			Return(nil, errors.New("ReplaceSubscription() Failed")).
			Times(1)

		err := r.Update(resourceData, clients)
		require.Contains(t, err.Error(), "ReplaceSubscription() Failed")
	}
}

func TestServicehookWebhookTfs_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServicehookWebhookTfs()
	for _, subscription := range testResourceSubscriptionWebhookTfs {
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

func TestServicehookWebhookTfs_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServicehookWebhookTfs()
	for _, subscription := range testResourceSubscriptionWebhookTfs {
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
