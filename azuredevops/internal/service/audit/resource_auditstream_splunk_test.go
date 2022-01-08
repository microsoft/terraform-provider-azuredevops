//go:build (all || resource_auditstream_splunk) && !exclude_auditstreams
// +build all resource_auditstream_splunk
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

var splunkDaysToBackfill = 0
var splunkTestStreamEnabled = true
var splunkTestAuditStream = audit.AuditStream{
	ConsumerInputs: &map[string]string{
		"SplunkUrl":                 "SPLUNK_TEST_url",
		"SplunkEventCollectorToken": "SPLUNK_TEST_token",
	},
	ConsumerType: converter.String("Splunk"),
	Id:           converter.Int(1),
}

// verifies that the flatten/expand round trip yields the same service endpoint
func TestAuditStreamSplunk_ExpandFlattenRoundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceAuditStreamSplunk().Schema, nil)
	flattenAuditStreamSplunk(resourceData, &splunkTestAuditStream, &splunkDaysToBackfill, &splunkTestStreamEnabled)

	auditStreamAfterRoundTrip, splunkDaysToBackfillAfterRoundTrip, enabled := expandAuditStreamSplunk(resourceData)

	require.Equal(t, splunkTestAuditStream, *auditStreamAfterRoundTrip)
	require.Equal(t, splunkDaysToBackfill, *splunkDaysToBackfillAfterRoundTrip)
	require.Equal(t, splunkTestStreamEnabled, *enabled)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestAuditStreamSplunk_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceAuditStreamSplunk()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenAuditStreamSplunk(resourceData, &splunkTestAuditStream, &splunkDaysToBackfill, &splunkTestStreamEnabled)

	buildClient := azdosdkmocks.NewMockAuditClient(ctrl)
	clients := &client.AggregatedClient{AuditClient: buildClient, Ctx: context.Background()}

	expectedArgs := audit.CreateStreamArgs{Stream: &splunkTestAuditStream, DaysToBackfill: &splunkDaysToBackfill}
	buildClient.
		EXPECT().
		CreateStream(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateStream() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateStream() Failed")
}

// verifies that if an error is produced on read, the error is not swallowed
func TestAuditStreamSplunk_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceAuditStreamSplunk()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenAuditStreamSplunk(resourceData, &splunkTestAuditStream, &splunkDaysToBackfill, &splunkTestStreamEnabled)

	buildClient := azdosdkmocks.NewMockAuditClient(ctrl)
	clients := &client.AggregatedClient{AuditClient: buildClient, Ctx: context.Background()}

	expectedArgs := audit.QueryStreamByIdArgs{StreamId: splunkTestAuditStream.Id}
	buildClient.
		EXPECT().
		QueryStreamById(clients.Ctx, expectedArgs).
		Return(nil, errors.New("QueryStreamById() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "QueryStreamById() Failed")
}

// // verifies that if an error is produced on update, the error is not swallowed
func TestAuditStreamSplunk_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceAuditStreamSplunk()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenAuditStreamSplunk(resourceData, &splunkTestAuditStream, &splunkDaysToBackfill, &splunkTestStreamEnabled)

	buildClient := azdosdkmocks.NewMockAuditClient(ctrl)
	clients := &client.AggregatedClient{AuditClient: buildClient, Ctx: context.Background()}

	expectedArgs := audit.UpdateStreamArgs{Stream: &splunkTestAuditStream}
	buildClient.
		EXPECT().
		UpdateStream(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateStream() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateStream() Failed")
}

// verifies that if an error is produced on delete, the error is not swallowed
func TestAuditStreamSplunk_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceAuditStreamSplunk()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenAuditStreamSplunk(resourceData, &splunkTestAuditStream, &splunkDaysToBackfill, &splunkTestStreamEnabled)

	buildClient := azdosdkmocks.NewMockAuditClient(ctrl)
	clients := &client.AggregatedClient{AuditClient: buildClient, Ctx: context.Background()}

	expectedArgs := audit.DeleteStreamArgs{StreamId: splunkTestAuditStream.Id}
	buildClient.
		EXPECT().
		DeleteStream(clients.Ctx, expectedArgs).
		Return(errors.New("DeleteStream() Failed")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteStream() Failed")
}
