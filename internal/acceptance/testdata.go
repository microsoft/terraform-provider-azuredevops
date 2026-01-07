package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

type TestData struct {
	// Either the resource type (e.g. azuredevops_project) or the data source type (e.g. data.azuredevops_project)
	ResourceType  string
	ResourceLabel string
	RandomInteger int
	RandomString  string
}

func (d TestData) ResourceAddr() string {
	return d.ResourceType + "." + d.ResourceLabel
}

// BuildTestData generates some test data for the given resource
func BuildTestData(t *testing.T, resourceType string, resourceLabel string) TestData {
	testData := TestData{
		ResourceType:  resourceType,
		ResourceLabel: resourceLabel,
		RandomInteger: acctest.RandInt(),
		RandomString:  acctest.RandString(5),
	}

	return testData
}
