package utils

import (
	"testing"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/audit"
	"github.com/stretchr/testify/require"
)

func TestExpandAuditStreamStatus(t *testing.T) {
	status := ExpandAuditStreamStatus("enabled")
	require.Equal(t, audit.AuditStreamStatusValues.Enabled, *status)

	status = ExpandAuditStreamStatus("disabledByUser")
	require.Equal(t, audit.AuditStreamStatusValues.DisabledByUser, *status)

	status = ExpandAuditStreamStatus("disabledBySystem")
	require.Equal(t, audit.AuditStreamStatusValues.DisabledBySystem, *status)

	status = ExpandAuditStreamStatus("deleted")
	require.Equal(t, audit.AuditStreamStatusValues.Deleted, *status)

	status = ExpandAuditStreamStatus("backfilling")
	require.Equal(t, audit.AuditStreamStatusValues.Backfilling, *status)

	status = ExpandAuditStreamStatus("unknown")
	require.Equal(t, audit.AuditStreamStatusValues.Enabled, *status)
}

func TestFlattenConsumerInputs(t *testing.T) {
	inputs := map[string]string{
		"url":   "https://example.com",
		"token": "secret",
	}

	flattened := FlattenConsumerInputs(&inputs)
	require.NotNil(t, flattened)
	require.Equal(t, 2, flattened.Len())

	list := flattened.List()
	foundUrl := false
	foundToken := false

	for _, item := range list {
		m := item.(map[string]interface{})
		if m["key"] == "url" {
			require.Equal(t, "https://example.com", m["value"])
			foundUrl = true
		}
		if m["key"] == "token" {
			require.Equal(t, "secret", m["value"])
			foundToken = true
		}
	}

	require.True(t, foundUrl)
	require.True(t, foundToken)
}

func TestFlattenSingleAuditStream(t *testing.T) {
	id := 123
	displayName := "Test Stream"
	consumerType := "splunk"
	status := audit.AuditStreamStatusValues.Enabled

	stream := &audit.AuditStream{
		Id:           &id,
		DisplayName:  &displayName,
		ConsumerType: &consumerType,
		Status:       &status,
		ConsumerInputs: &map[string]string{
			"url": "https://splunk.example.com",
		},
	}

	m := make(map[string]interface{})
	err := FlattenSingleAuditStream(m, stream)
	require.NoError(t, err)

	require.Equal(t, "123", m["id"])
	require.Equal(t, "Test Stream", m["display_name"])
	require.Equal(t, "splunk", m["consumer_type"])
	require.Equal(t, "enabled", m["status"])

	inputs := m["consumer_inputs"].([]interface{})
	require.Equal(t, 1, len(inputs))
	inputMap := inputs[0].(map[string]interface{})
	require.Equal(t, "url", inputMap["key"])
	require.Equal(t, "https://splunk.example.com", inputMap["value"])
}
