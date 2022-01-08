//go:build (all || resource_auditstream_azuremonitorlogs) && !exclude_auditstreams
// +build all resource_auditstream_azuremonitorlogs
// +build !exclude_auditstreams

package audit

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/audit"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

var azureMonitorDaysToBackfill = 0
var azureMonitorLogsTestStreamEnabled = true
var azureMonitorLogsTestAuditStream = audit.AuditStream{
	ConsumerInputs: &map[string]string{
		"WorkspaceId": "AZURE_MONITOR_LOGS_TEST_workspace_id",
		"SharedKey":   "AZURE_MONITOR_LOGS_TEST_shared_key",
	},
	ConsumerType: converter.String("AzureMonitorLogs"),
	Id:           converter.Int(1),
}

// verifies that the flatten/expand round trip yields the same service endpoint
func TestAuditStreamAzureMonitorLogs_ExpandFlattenRoundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceAuditStreamAzureMonitorLogs().Schema, nil)
	flattenAuditStreamAzureMonitorLogs(resourceData, &azureMonitorLogsTestAuditStream, &azureMonitorDaysToBackfill, &azureMonitorLogsTestStreamEnabled)

	auditStreamAfterRoundTrip, azureMonitorDaysToBackfillAfterRoundTrip, enabled := expandAuditStreamAzureMonitorLogs(resourceData)

	require.Equal(t, azureMonitorLogsTestAuditStream, *auditStreamAfterRoundTrip)
	require.Equal(t, azureMonitorDaysToBackfill, *azureMonitorDaysToBackfillAfterRoundTrip)
	require.Equal(t, azureMonitorLogsTestStreamEnabled, *enabled)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestAuditStreamAzureMonitorLogs_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceAuditStreamAzureMonitorLogs()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenAuditStreamAzureMonitorLogs(resourceData, &azureMonitorLogsTestAuditStream, &azureMonitorDaysToBackfill, &azureMonitorLogsTestStreamEnabled)

	buildClient := azdosdkmocks.NewMockAuditClient(ctrl)
	clients := &client.AggregatedClient{AuditClient: buildClient, Ctx: context.Background()}

	expectedArgs := audit.CreateStreamArgs{Stream: &azureMonitorLogsTestAuditStream, DaysToBackfill: &azureMonitorDaysToBackfill}
	buildClient.
		EXPECT().
		CreateStream(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateStream() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateStream() Failed")
}

// verifies that if an error is produced on read, the error is not swallowed
func TestAuditStreamAzureMonitorLogs_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceAuditStreamAzureMonitorLogs()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenAuditStreamAzureMonitorLogs(resourceData, &azureMonitorLogsTestAuditStream, &azureMonitorDaysToBackfill, &azureMonitorLogsTestStreamEnabled)

	buildClient := azdosdkmocks.NewMockAuditClient(ctrl)
	clients := &client.AggregatedClient{AuditClient: buildClient, Ctx: context.Background()}

	expectedArgs := audit.QueryStreamByIdArgs{StreamId: azureMonitorLogsTestAuditStream.Id}
	buildClient.
		EXPECT().
		QueryStreamById(clients.Ctx, expectedArgs).
		Return(nil, errors.New("QueryStreamById() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "QueryStreamById() Failed")
}

// // verifies that if an error is produced on update, the error is not swallowed
func TestAuditStreamAzureMonitorLogs_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceAuditStreamAzureMonitorLogs()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenAuditStreamAzureMonitorLogs(resourceData, &azureMonitorLogsTestAuditStream, &azureMonitorDaysToBackfill, &azureMonitorLogsTestStreamEnabled)

	buildClient := azdosdkmocks.NewMockAuditClient(ctrl)
	clients := &client.AggregatedClient{AuditClient: buildClient, Ctx: context.Background()}

	expectedArgs := audit.UpdateStreamArgs{Stream: &azureMonitorLogsTestAuditStream}
	buildClient.
		EXPECT().
		UpdateStream(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateStream() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateStream() Failed")
}

// verifies that if an error is produced on delete, the error is not swallowed
func TestAuditStreamAzureMonitorLogs_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceAuditStreamAzureMonitorLogs()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenAuditStreamAzureMonitorLogs(resourceData, &azureMonitorLogsTestAuditStream, &azureMonitorDaysToBackfill, &azureMonitorLogsTestStreamEnabled)

	buildClient := azdosdkmocks.NewMockAuditClient(ctrl)
	clients := &client.AggregatedClient{AuditClient: buildClient, Ctx: context.Background()}

	expectedArgs := audit.DeleteStreamArgs{StreamId: azureMonitorLogsTestAuditStream.Id}
	buildClient.
		EXPECT().
		DeleteStream(clients.Ctx, expectedArgs).
		Return(errors.New("DeleteStream() Failed")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteStream() Failed")
}
