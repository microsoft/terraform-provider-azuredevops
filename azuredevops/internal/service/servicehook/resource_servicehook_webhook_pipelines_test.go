//go:build all || servicehook
// +build all servicehook

package servicehook

import (
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/servicehooks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

func TestExpandServicehookWebhookPipelines_StageStateChanged(t *testing.T) {
	r := ResourceServicehookWebhookPipelines()
	d := r.TestResourceData()

	projectID := uuid.New().String()
	require.NoError(t, d.Set("project_id", projectID))
	require.NoError(t, d.Set("url", "https://example.azurewebsites.net/api/BuildFailed"))
	require.NoError(t, d.Set("http_headers", map[string]interface{}{
		"x-functions-key": "secret",
	}))
	require.NoError(t, d.Set("stage_state_changed_event", []interface{}{
		map[string]interface{}{
			"stage_state_filter":  "Completed",
			"stage_result_filter": "Canceled",
		},
	}))

	sub := expandServicehookWebhookPipelines(d)

	require.Equal(t, "pipelines", *sub.PublisherId)
	require.Equal(t, "webHooks", *sub.ConsumerId)
	require.Equal(t, "httpRequest", *sub.ConsumerActionId)
	require.Equal(t, "ms.vss-pipelines.stage-state-changed-event", *sub.EventType)

	require.NotNil(t, sub.PublisherInputs)
	require.Equal(t, projectID, (*sub.PublisherInputs)["projectId"])
	require.Equal(t, "Completed", (*sub.PublisherInputs)["stageStateId"])
	require.Equal(t, "Canceled", (*sub.PublisherInputs)["stageResultId"])

	require.NotNil(t, sub.ConsumerInputs)
	require.Equal(t, "https://example.azurewebsites.net/api/BuildFailed", (*sub.ConsumerInputs)["url"])
	require.Equal(t, "x-functions-key:secret", (*sub.ConsumerInputs)["httpHeaders"])
}

func TestFlattenServicehookWebhookPipelines_RoundTrip(t *testing.T) {
	r := ResourceServicehookWebhookPipelines()
	d := r.TestResourceData()

	projectID := uuid.New().String()
	sub := &servicehooks.Subscription{
		ConsumerActionId: converter.String("httpRequest"),
		ConsumerId:       converter.String("webHooks"),
		ConsumerInputs: &map[string]string{
			"url":                    "https://example.com/hook",
			"acceptUntrustedCerts":   "true",
			"resourceDetailsToSend":  "all",
			"messagesToSend":         "all",
			"detailedMessagesToSend": "all",
			"httpHeaders":            "x-functions-key:secret\ncustomId:1",
		},
		EventType:   converter.String("ms.vss-pipelines.stage-state-changed-event"),
		PublisherId: converter.String("pipelines"),
		PublisherInputs: &map[string]string{
			"projectId":     projectID,
			"stageStateId":  "Completed",
			"stageResultId": "Failed",
		},
		ResourceVersion: converter.String("5.1-preview.1"),
	}

	flattenServicehookWebhookPipelines(d, sub)

	require.Equal(t, projectID, d.Get("project_id"))
	require.Equal(t, "https://example.com/hook", d.Get("url"))
	require.Equal(t, true, d.Get("accept_untrusted_certs"))
	require.Equal(t, "5.1-preview.1", d.Get("resource_version"))

	headers := d.Get("http_headers").(map[string]interface{})
	require.Equal(t, "secret", headers["x-functions-key"])
	require.Equal(t, "1", headers["customId"])

	stageList := d.Get("stage_state_changed_event").([]interface{})
	require.Len(t, stageList, 1)
	stage := stageList[0].(map[string]interface{})
	require.Equal(t, "Completed", stage["stage_state_filter"])
	require.Equal(t, "Failed", stage["stage_result_filter"])
}
