package azuredevops

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/require"
)

var provider = Provider()

// TestAzureDevOpsProvider_foo

func TestAzureDevOpsProvider_HasChildResources(t *testing.T) {
	expectedResources := []string{
		"azuredevops_build_definition",
		"azuredevops_project",
		"azuredevops_serviceendpoint_github",
		"azuredevops_serviceendpoint_dockerhub",
		"azuredevops_serviceendpoint_azurerm",
		"azuredevops_variable_group",
		"azuredevops_git_repository",
		"azuredevops_user_entitlement",
		"azuredevops_group_membership",
		"azuredevops_group",
		"azuredevops_agent_pool",
	}

	resources := provider.ResourcesMap
	require.Equal(t, len(expectedResources), len(resources), "There are an unexpected number of registered resources")

	for _, resource := range expectedResources {
		require.Contains(t, resources, resource, "An expected resource was not registered")
		require.NotNil(t, resources[resource], "A resource cannot have a nil schema")
	}
}

func TestAzureDevOpsProvider_HasChildDataSources(t *testing.T) {
	expectedDataSources := []string{
		"azuredevops_group",
		"azuredevops_projects",
		"azuredevops_git_repositories",
	}

	dataSources := provider.DataSourcesMap
	require.Equal(t, len(expectedDataSources), len(dataSources), "There are an unexpected number of registered data sources")

	for _, resource := range expectedDataSources {
		require.Contains(t, dataSources, resource, "An expected data source was not registered")
		require.NotNil(t, dataSources[resource], "A data source cannot have a nil schema")
	}
}

func TestAzureDevOpsProvider_SchemaIsValid(t *testing.T) {
	type testParams struct {
		name          string
		required      bool
		defaultEnvVar string
		sensitive     bool
	}

	tests := []testParams{
		{"org_service_url", true, "AZDO_ORG_SERVICE_URL", false},
		{"personal_access_token", true, "AZDO_PERSONAL_ACCESS_TOKEN", true},
	}

	schema := provider.Schema
	require.Equal(t, len(tests), len(schema), "There are an unexpected number of properties in the schema")

	for _, test := range tests {
		require.Contains(t, schema, test.name, "An expected property was not found in the schema")
		require.NotNil(t, schema[test.name], "A property in the schema cannot have a nil value")
		require.Equal(t, test.sensitive, schema[test.name].Sensitive, "A property in the schema has an incorrect sensitivity value")

		if test.defaultEnvVar != "" {
			expectedValue := os.Getenv(test.defaultEnvVar)

			actualValue, err := schema[test.name].DefaultFunc()
			if actualValue == nil {
				actualValue = ""
			}

			require.Nil(t, err, "An error occurred when getting the default value from the environment")
			require.Equal(t, expectedValue, actualValue, "The default value pulled from the environment has the wrong value")
		}
	}
}

func init() {
	InitProvider()
}

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func InitProvider() {
	testAccProvider = provider
	testAccProviders = map[string]terraform.ResourceProvider{
		"azuredevops": testAccProvider,
	}
}
