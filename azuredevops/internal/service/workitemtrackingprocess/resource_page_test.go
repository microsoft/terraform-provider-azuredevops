//go:build (all || resource_workitemtrackingprocess || resource_workitemtrackingprocess_page) && !exclude_resource_workitemtrackingprocess
// +build all resource_workitemtrackingprocess resource_workitemtrackingprocess_page
// +build !exclude_resource_workitemtrackingprocess

package workitemtrackingprocess

import (
	"testing"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
)

func TestPage_FindPageById(t *testing.T) {
	pageId := "target-page"
	pageId1 := "page-1"
	pageId2 := "page-2"

	tests := []struct {
		name     string
		layout   *workitemtrackingprocess.FormLayout
		pageId   string
		expected bool
	}{
		{
			name: "found first page",
			layout: &workitemtrackingprocess.FormLayout{
				Pages: &[]workitemtrackingprocess.Page{
					{Id: &pageId},
				},
			},
			pageId:   pageId,
			expected: true,
		},
		{
			name: "found second page",
			layout: &workitemtrackingprocess.FormLayout{
				Pages: &[]workitemtrackingprocess.Page{
					{Id: &pageId1},
					{Id: &pageId},
				},
			},
			pageId:   pageId,
			expected: true,
		},
		{
			name: "found among multiple pages",
			layout: &workitemtrackingprocess.FormLayout{
				Pages: &[]workitemtrackingprocess.Page{
					{Id: &pageId1},
					{Id: &pageId2},
					{Id: &pageId},
					{Id: converter.String("page-3")},
				},
			},
			pageId:   pageId,
			expected: true,
		},
		{
			name: "not found",
			layout: &workitemtrackingprocess.FormLayout{
				Pages: &[]workitemtrackingprocess.Page{
					{Id: &pageId1},
					{Id: &pageId2},
				},
			},
			pageId:   "nonexistent",
			expected: false,
		},
		{
			name:     "nil layout",
			layout:   nil,
			pageId:   pageId,
			expected: false,
		},
		{
			name: "nil pages",
			layout: &workitemtrackingprocess.FormLayout{
				Pages: nil,
			},
			pageId:   pageId,
			expected: false,
		},
		{
			name: "empty pages",
			layout: &workitemtrackingprocess.FormLayout{
				Pages: &[]workitemtrackingprocess.Page{},
			},
			pageId:   pageId,
			expected: false,
		},
		{
			name: "page with nil id",
			layout: &workitemtrackingprocess.FormLayout{
				Pages: &[]workitemtrackingprocess.Page{
					{Id: nil},
				},
			},
			pageId:   pageId,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findPageById(tt.layout, tt.pageId)
			if tt.expected {
				assert.NotNil(t, result, "expected to find page")
				assert.Equal(t, tt.pageId, *result.Id)
			} else {
				assert.Nil(t, result, "expected not to find page")
			}
		})
	}
}
