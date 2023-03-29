//go:build (all || resource_auditstream_azureeventgrid) && !exclude_auditstreams
// +build all resource_auditstream_azureeventgrid
// +build !exclude_auditstreams

package audit

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/audit"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

var eventgridDaysToBackfill = 0
var azureEventGridTestStreamEnabled = true
var azureEventGridTestAuditStream = audit.AuditStream{
	ConsumerInputs: &map[string]string{
		"EventGridTopicHostname":  "AZURE_EVENTGRID_TOPIC_TEST_url",
		"EventGridTopicAccessKey": "AZURE_EVENTGRID_TOPIC_TEST_access_key",
	},
	ConsumerType: converter.String("AzureEventGrid"),
	Id:           converter.Int(1),
}

// verifies that the flatten/expand round trip yields the same service endpoint
func TestAuditStreamAzureEventGrid_ExpandFlattenRoundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceAuditStreamAzureEventGridTopic().Schema, nil)
	flattenAuditStreamAzureEventGridTopic(resourceData, &azureEventGridTestAuditStream, &eventgridDaysToBackfill, &azureEventGridTestStreamEnabled)

	auditStreamAfterRoundTrip, eventgridDaysToBackfillAfterRoundTrip, enabled := expandAuditStreamAzureEventGridTopic(resourceData)

	require.Equal(t, azureEventGridTestAuditStream, *auditStreamAfterRoundTrip)
	require.Equal(t, eventgridDaysToBackfill, *eventgridDaysToBackfillAfterRoundTrip)
	require.Equal(t, azureEventGridTestStreamEnabled, *enabled)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestAuditStreamAzureEventGrid_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceAuditStreamAzureEventGridTopic()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenAuditStreamAzureEventGridTopic(resourceData, &azureEventGridTestAuditStream, &eventgridDaysToBackfill, &azureEventGridTestStreamEnabled)

	buildClient := azdosdkmocks.NewMockAuditClient(ctrl)
	clients := &client.AggregatedClient{AuditClient: buildClient, Ctx: context.Background()}

	expectedArgs := audit.CreateStreamArgs{Stream: &azureEventGridTestAuditStream, DaysToBackfill: &eventgridDaysToBackfill}
	buildClient.
		EXPECT().
		CreateStream(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateStream() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateStream() Failed")
}

// verifies that if an error is produced on read, the error is not swallowed
func TestAuditStreamAzureEventGrid_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceAuditStreamAzureEventGridTopic()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenAuditStreamAzureEventGridTopic(resourceData, &azureEventGridTestAuditStream, &eventgridDaysToBackfill, &azureEventGridTestStreamEnabled)

	buildClient := azdosdkmocks.NewMockAuditClient(ctrl)
	clients := &client.AggregatedClient{AuditClient: buildClient, Ctx: context.Background()}

	expectedArgs := audit.QueryStreamByIdArgs{StreamId: azureEventGridTestAuditStream.Id}
	buildClient.
		EXPECT().
		QueryStreamById(clients.Ctx, expectedArgs).
		Return(nil, errors.New("QueryStreamById() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "QueryStreamById() Failed")
}

// // verifies that if an error is produced on update, the error is not swallowed
func TestAuditStreamAzureEventGrid_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceAuditStreamAzureEventGridTopic()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenAuditStreamAzureEventGridTopic(resourceData, &azureEventGridTestAuditStream, &eventgridDaysToBackfill, &azureEventGridTestStreamEnabled)

	buildClient := azdosdkmocks.NewMockAuditClient(ctrl)
	clients := &client.AggregatedClient{AuditClient: buildClient, Ctx: context.Background()}

	expectedArgs := audit.UpdateStreamArgs{Stream: &azureEventGridTestAuditStream}
	buildClient.
		EXPECT().
		UpdateStream(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateStream() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateStream() Failed")
}

// verifies that if an error is produced on delete, the error is not swallowed
func TestAuditStreamAzureEventGrid_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceAuditStreamAzureEventGridTopic()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenAuditStreamAzureEventGridTopic(resourceData, &azureEventGridTestAuditStream, &eventgridDaysToBackfill, &azureEventGridTestStreamEnabled)

	buildClient := azdosdkmocks.NewMockAuditClient(ctrl)
	clients := &client.AggregatedClient{AuditClient: buildClient, Ctx: context.Background()}

	expectedArgs := audit.DeleteStreamArgs{StreamId: azureEventGridTestAuditStream.Id}
	buildClient.
		EXPECT().
		DeleteStream(clients.Ctx, expectedArgs).
		Return(errors.New("DeleteStream() Failed")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteStream() Failed")
}
