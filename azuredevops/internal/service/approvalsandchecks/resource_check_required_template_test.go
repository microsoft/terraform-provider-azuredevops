package approvalsandchecks

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/pipelineschecksextras"
	"github.com/stretchr/testify/require"
)

var requiredTemplateCheckID = 123456789
var requiredTemplateEndpointID = uuid.New().String()
var requiredTemplateCheckProjectID = uuid.New().String()

var requiredTemplateTestEndpointType = "endpoint"
var requiredTemplateTestEndpointResource = pipelineschecksextras.Resource{
	Id:   &requiredTemplateEndpointID,
	Type: &requiredTemplateTestEndpointType,
}

var requiredTemplates = []interface{}{
	map[string]interface{}{
		"repositoryType": "git",
		"repositoryName": "proj/repo",
		"repositoryRef":  "refs/heads/master",
		"templatePath":   "templ/other-path.yaml",
	},
}

var requiredTemplateCheckSettings = map[string]interface{}{
	"extendsChecks": requiredTemplates,
}

var requiredTemplateCheckTest = pipelineschecksextras.CheckConfiguration{
	Id:       &requiredTemplateCheckID,
	Type:     approvalAndCheckType.ExtendsCheck,
	Settings: requiredTemplateCheckSettings,
	Resource: &requiredTemplateTestEndpointResource,
	Version:  converter.Int(0),
}

// verifies that the flatten/expand round trip yields the same required template
func TestCheckRequiredTemplate_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceCheckRequiredTemplate().Schema, nil)
	flattenErr := flattenCheckRequiredTemplate(resourceData, &requiredTemplateCheckTest, requiredTemplateCheckProjectID)

	requiredTemplateCheckAfterRoundTrip, projectID, expandErr := expandCheckRequiredTemplate(resourceData)
	requiredTemplateCheckAfterRoundTrip.Id = requiredTemplateCheckTest.Id

	require.Equal(t, requiredTemplateCheckTest, *requiredTemplateCheckAfterRoundTrip)
	require.Equal(t, requiredTemplateCheckProjectID, projectID)
	require.Nil(t, expandErr)
	require.Nil(t, flattenErr)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestCheckRequiredTemplate_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceCheckRequiredTemplate()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.SetId(fmt.Sprintf("%d", *requiredTemplateCheckTest.Id))
	flattenErr := flattenCheckRequiredTemplate(resourceData, &requiredTemplateCheckTest, requiredTemplateCheckProjectID)

	pipelinesChecksClient := azdosdkmocks.NewMockPipelineschecksextrasClient(ctrl)
	clients := &client.AggregatedClient{PipelinesChecksClientExtras: pipelinesChecksClient, Ctx: context.Background()}

	expectedArgs := pipelineschecksextras.AddCheckConfigurationArgs{Configuration: &requiredTemplateCheckTest, Project: &requiredTemplateCheckProjectID}
	//expectedArgs = requiredTemplateCheckTest.Id
	pipelinesChecksClient.
		EXPECT().
		AddCheckConfiguration(clients.Ctx, expectedArgs).
		Return(nil, errors.New("AddCheckConfiguration() Failed")).
		Times(1)

	err := r.Create(resourceData, clients) //nolint:staticcheck
	require.Contains(t, err.Error(), "AddCheckConfiguration() Failed")
	require.Nil(t, flattenErr)
}

// verifies that if an error is produced on a read, it is not swallowed
func TestCheckRequiredTemplate_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceCheckRequiredTemplate()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.SetId(fmt.Sprintf("%d", *requiredTemplateCheckTest.Id))
	flattenErr := flattenCheckRequiredTemplate(resourceData, &requiredTemplateCheckTest, requiredTemplateCheckProjectID)

	pipelinesChecksClient := azdosdkmocks.NewMockPipelineschecksextrasClient(ctrl)
	clients := &client.AggregatedClient{PipelinesChecksClientExtras: pipelinesChecksClient, Ctx: context.Background()}

	expectedArgs := pipelineschecksextras.GetCheckConfigurationArgs{
		Id:      requiredTemplateCheckTest.Id,
		Project: &requiredTemplateCheckProjectID,
		Expand:  converter.ToPtr(pipelineschecksextras.CheckConfigurationExpandParameterValues.Settings),
	}

	pipelinesChecksClient.
		EXPECT().
		GetCheckConfiguration(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetServiceEndpoint() Failed")).
		Times(1)

	err := r.Read(resourceData, clients) //nolint:staticcheck
	require.Contains(t, err.Error(), "GetServiceEndpoint() Failed")
	require.Nil(t, flattenErr)
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestCheckRequiredTemplate_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceCheckRequiredTemplate()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.SetId(fmt.Sprintf("%d", *requiredTemplateCheckTest.Id))
	flattenErr := flattenCheckRequiredTemplate(resourceData, &requiredTemplateCheckTest, requiredTemplateCheckProjectID)

	pipelinesChecksClient := azdosdkmocks.NewMockPipelineschecksextrasClient(ctrl)
	clients := &client.AggregatedClient{PipelinesChecksClientExtras: pipelinesChecksClient, Ctx: context.Background()}

	expectedArgs := pipelineschecksextras.DeleteCheckConfigurationArgs{
		Id:      requiredTemplateCheckTest.Id,
		Project: &requiredTemplateCheckProjectID,
	}

	pipelinesChecksClient.
		EXPECT().
		DeleteCheckConfiguration(clients.Ctx, expectedArgs).
		Return(errors.New("DeleteServiceEndpoint() Failed")).
		Times(1)

	err := r.Delete(resourceData, clients) //nolint:staticcheck
	require.Contains(t, err.Error(), "DeleteServiceEndpoint() Failed")
	require.Nil(t, flattenErr)
}

// verifies that if an error is produced on an update, it is not swallowed
func TestCheckRequiredTemplate_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceCheckRequiredTemplate()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.SetId(fmt.Sprintf("%d", *requiredTemplateCheckTest.Id))
	flattenErr := flattenCheckRequiredTemplate(resourceData, &requiredTemplateCheckTest, requiredTemplateCheckProjectID)

	pipelinesChecksClient := azdosdkmocks.NewMockPipelineschecksextrasClient(ctrl)
	clients := &client.AggregatedClient{PipelinesChecksClientExtras: pipelinesChecksClient, Ctx: context.Background()}

	expectedArgs := pipelineschecksextras.UpdateCheckConfigurationArgs{
		Project:       &requiredTemplateCheckProjectID,
		Configuration: &requiredTemplateCheckTest,
		Id:            &requiredTemplateCheckID,
	}

	pipelinesChecksClient.
		EXPECT().
		UpdateCheckConfiguration(clients.Ctx, expectedArgs).
		Return(nil, errors.New("UpdateServiceEndpoint() Failed")).
		Times(1)

	err := r.Update(resourceData, clients) //nolint:staticcheck
	require.Contains(t, err.Error(), "UpdateServiceEndpoint() Failed")
	require.Nil(t, flattenErr)
}
