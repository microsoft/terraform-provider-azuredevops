//go:build (all || resource_check_business_hours) && !exclude_approvalsandchecks
// +build all resource_check_business_hours
// +build !exclude_approvalsandchecks

package approvalsandchecks

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/sdk/pipelineschecksextras"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	CheckBusinessHoursID            = 123456789
	CheckBusinessHoursProjectID     = uuid.New().String()
	CheckBusinessHoursTestProjectID = &CheckBusinessHoursProjectID
)

var CheckBusinessHoursInputs = map[string]interface{}{
	"businessDays": "Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday",
	"timeZone":     "UTC",
	"startTime":    "01:00",
	"endTime":      "02:00",
}

var CheckBusinessHoursSettings = map[string]interface{}{
	"definitionRef": evaluateBusinessHoursDef,
	"displayName":   "Test Business Hours",
	"inputs":        CheckBusinessHoursInputs,
}

var CheckBusinessHoursTest = pipelineschecksextras.CheckConfiguration{
	Id:       &CheckBusinessHoursID,
	Type:     approvalAndCheckType.BusinessHours,
	Settings: CheckBusinessHoursSettings,
	Timeout:  converter.ToPtr(50000),
	Resource: &endpointResource,
	Version:  converter.Int(0),
}

// verifies that the flatten/expand round trip yields the same business hours check
func TestCheckBusinessHours_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceCheckBusinessHours().Schema, nil)
	resourceData.SetId(fmt.Sprintf("%d", CheckBusinessHoursID))
	flattenBusinessHours(resourceData, &CheckBusinessHoursTest, CheckBusinessHoursProjectID)

	CheckBusinessHoursAfterRoundTrip, projectID, err := expandBusinessHours(resourceData)

	require.Equal(t, CheckBusinessHoursTest, *CheckBusinessHoursAfterRoundTrip)
	require.Equal(t, CheckBusinessHoursProjectID, projectID)
	require.Nil(t, err)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestCheckBusinessHours_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceCheckBusinessHours()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.SetId(fmt.Sprintf("%d", CheckBusinessHoursID))
	flattenBusinessHours(resourceData, &CheckBusinessHoursTest, CheckBusinessHoursProjectID)

	pipelinesCheckClient := azdosdkmocks.NewMockPipelineschecksextrasClient(ctrl)
	clients := &client.AggregatedClient{PipelinesChecksClientExtras: pipelinesCheckClient, Ctx: context.Background()}

	expectedArgs := pipelineschecksextras.AddCheckConfigurationArgs{Configuration: &CheckBusinessHoursTest, Project: &CheckBusinessHoursProjectID}
	pipelinesCheckClient.
		EXPECT().
		AddCheckConfiguration(clients.Ctx, expectedArgs).
		Return(nil, errors.New("AddCheckConfiguration() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "AddCheckConfiguration() Failed")
}

// verifies that if an error is produced on a read, it is not swallowed
func TestCheckBusinessHours_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceCheckBusinessHours()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.SetId(fmt.Sprintf("%d", *CheckBusinessHoursTest.Id))
	flattenBusinessHours(resourceData, &CheckBusinessHoursTest, CheckBusinessHoursProjectID)

	pipelinesCheckClient := azdosdkmocks.NewMockPipelineschecksextrasClient(ctrl)
	clients := &client.AggregatedClient{PipelinesChecksClientExtras: pipelinesCheckClient, Ctx: context.Background()}

	expectedArgs := pipelineschecksextras.GetCheckConfigurationArgs{
		Id:      CheckBusinessHoursTest.Id,
		Project: &CheckBusinessHoursProjectID,
		Expand:  converter.ToPtr(pipelineschecksextras.CheckConfigurationExpandParameterValues.Settings),
	}

	pipelinesCheckClient.
		EXPECT().
		GetCheckConfiguration(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetServiceEndpoint() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "GetServiceEndpoint() Failed")
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestCheckBusinessHours_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceCheckBusinessHours()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.SetId(fmt.Sprintf("%d", *CheckBusinessHoursTest.Id))
	flattenBusinessHours(resourceData, &CheckBusinessHoursTest, CheckBusinessHoursProjectID)

	pipelinesCheckClient := azdosdkmocks.NewMockPipelineschecksextrasClient(ctrl)
	clients := &client.AggregatedClient{PipelinesChecksClientExtras: pipelinesCheckClient, Ctx: context.Background()}

	expectedArgs := pipelineschecksextras.DeleteCheckConfigurationArgs{
		Id:      CheckBusinessHoursTest.Id,
		Project: &CheckBusinessHoursProjectID,
	}

	pipelinesCheckClient.
		EXPECT().
		DeleteCheckConfiguration(clients.Ctx, expectedArgs).
		Return(errors.New("DeleteServiceEndpoint() Failed")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteServiceEndpoint() Failed")
}

// verifies that if an error is produced on an update, it is not swallowed
func TestCheckBusinessHours_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceCheckBusinessHours()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.SetId(fmt.Sprintf("%d", *CheckBusinessHoursTest.Id))
	flattenBusinessHours(resourceData, &CheckBusinessHoursTest, CheckBusinessHoursProjectID)

	pipelinesCheckClient := azdosdkmocks.NewMockPipelineschecksextrasClient(ctrl)
	clients := &client.AggregatedClient{PipelinesChecksClientExtras: pipelinesCheckClient, Ctx: context.Background()}

	expectedArgs := pipelineschecksextras.UpdateCheckConfigurationArgs{
		Project:       &CheckBusinessHoursProjectID,
		Configuration: &CheckBusinessHoursTest,
		Id:            &CheckBusinessHoursID,
	}

	pipelinesCheckClient.
		EXPECT().
		UpdateCheckConfiguration(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
}
