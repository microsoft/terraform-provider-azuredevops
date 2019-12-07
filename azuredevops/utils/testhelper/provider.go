package testhelper

import (
	"os"
	"testing"
)

// TestAccPreCheck pre-check to validate that the correct environment variables are set prior to running any acceptance test
func TestAccPreCheck(t *testing.T) {
	requiredEnvVars := []string{
		"AZDO_ORG_SERVICE_URL",
		"AZDO_PERSONAL_ACCESS_TOKEN",
		"AZDO_GITHUB_SERVICE_CONNECTION_PAT",
		"AZDO_TEST_AAD_USER_EMAIL",
		"AZDO_DOCKERHUB_SERVICE_CONNECTION_USERNAME",
		"AZDO_DOCKERHUB_SERVICE_CONNECTION_EMAIL",
		"AZDO_DOCKERHUB_SERVICE_CONNECTION_PASSWORD",
	}

	for _, variable := range requiredEnvVars {
		if os.Getenv(variable) == "" {
			t.Fatalf("`%s` must be set for acceptance tests!", variable)
		}
	}
}

// TestAccResourcePrefix the default prefix for Terrfaorm objects in acceptance tests
const TestAccResourcePrefix = "test-acc-"
