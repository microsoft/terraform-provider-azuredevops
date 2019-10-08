package azuredevops

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/require"
)

var provider = Provider()

// TestAzureDevOpsProvider_foo

func TestAzureDevOpsProvider_HasChildResources(t *testing.T) {
	expectedResources := []string{
		"azuredevops_build_definition",
		"azuredevops_project",
		"azuredevops_serviceendpoint",
	}

	resources := provider.ResourcesMap
	require.Equal(t, len(expectedResources), len(resources), "There are an unexpected number of registered resources")

	for _, resource := range expectedResources {
		require.Contains(t, resources, resource, "An expected resource was not registered")
		require.NotNil(t, resources[resource], "A resource cannot have a nil schema")
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
			expectedValue := "foo-env-var"
			os.Setenv(test.defaultEnvVar, expectedValue)

			actualValue, err := schema[test.name].DefaultFunc()
			require.Nil(t, err, "An error occurred when getting the default value from the environment")
			require.Equal(t, expectedValue, actualValue, "The default value pulled from the environment has the wrong value")
		}
	}
}

// The configuration below can be used for every acceptance test in the project. For that reason,
// it will be defined once and will live in this file.

// pre-check to validate that the correct environment variables are set prior to running any acceptance test
func testAccPreCheck(t *testing.T) {
	requiredEnvVars := []string{
		"AZDO_ORG_SERVICE_URL",
		"AZDO_PERSONAL_ACCESS_TOKEN",
	}

	for _, variable := range requiredEnvVars {
		if os.Getenv(variable) == "" {
			t.Fatalf("`%s` must be set for acceptance tests!", variable)
		}
	}
}

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = provider
	testAccProviders = map[string]terraform.ResourceProvider{
		"azuredevops": testAccProvider,
	}
}
