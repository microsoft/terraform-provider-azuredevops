//go:build (all || resource_workitemtrackingprocess || resource_workitemtrackingprocess_inherited_page) && !exclude_resource_workitemtrackingprocess
// +build all resource_workitemtrackingprocess resource_workitemtrackingprocess_inherited_page
// +build !exclude_resource_workitemtrackingprocess

package workitemtrackingprocess

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func getInheritedPageResourceData(t *testing.T, input map[string]any) *schema.ResourceData {
	return schema.TestResourceDataRaw(t, ResourceInheritedPage().Schema, input)
}

func createProcessWorkItemTypeWithPage(witRefName string, page workitemtrackingprocess.Page) *workitemtrackingprocess.ProcessWorkItemType {
	return &workitemtrackingprocess.ProcessWorkItemType{
		ReferenceName: &witRefName,
		Layout: &workitemtrackingprocess.FormLayout{
			Pages: &[]workitemtrackingprocess.Page{
				page,
			},
		},
	}
}

func TestInheritedPage_Create_Validation(t *testing.T) {
	processId := uuid.New()
	witRefName := "MyProcess.MyWorkItemType"
	existingPageId := "page-1"
	inherited := true
	notInherited := false

	tests := []struct {
		name               string
		pageId             string
		returnWorkItemType *workitemtrackingprocess.ProcessWorkItemType
		returnError        error
		expectedError      string
	}{
		{
			name:          "API error",
			pageId:        existingPageId,
			returnError:   fmt.Errorf("API error"),
			expectedError: "getting work item type",
		},
		{
			name:          "nil work item type",
			pageId:        existingPageId,
			expectedError: "work item type is nil",
		},
		{
			name:   "nil layout",
			pageId: existingPageId,
			returnWorkItemType: &workitemtrackingprocess.ProcessWorkItemType{
				ReferenceName: &witRefName,
				Layout:        nil,
			},
			expectedError: "work item type layout is nil",
		},
		{
			name:   "page not found",
			pageId: "non-existent-page",
			returnWorkItemType: createProcessWorkItemTypeWithPage(witRefName, workitemtrackingprocess.Page{
				Id:        &existingPageId,
				Inherited: &inherited,
			}),
			expectedError: "page non-existent-page not found in layout",
		},
		{
			name:   "page not inherited",
			pageId: existingPageId,
			returnWorkItemType: createProcessWorkItemTypeWithPage(witRefName, workitemtrackingprocess.Page{
				Id:        &existingPageId,
				Inherited: &notInherited,
			}),
			expectedError: "page page-1 is not inherited",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := azdosdkmocks.NewMockWorkitemtrackingprocessClient(ctrl)
			clients := &client.AggregatedClient{WorkItemTrackingProcessClient: mockClient, Ctx: context.Background()}

			mockClient.EXPECT().GetProcessWorkItemType(clients.Ctx, gomock.Any()).Return(tt.returnWorkItemType, tt.returnError).Times(1)

			d := getInheritedPageResourceData(t, map[string]any{
				"process_id":        processId.String(),
				"work_item_type_id": witRefName,
				"page_id":           tt.pageId,
			})

			diags := createResourceInheritedPage(context.Background(), d, clients)
			assert.NotEmpty(t, diags)
			assert.Contains(t, diags[0].Summary, tt.expectedError)
		})
	}
}
