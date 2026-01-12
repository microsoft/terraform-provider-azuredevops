package acceptance

import (
	"os"
	"testing"
)

func PreCheck(t *testing.T) {
	variables := []string{
		"AZDO_PERSONAL_ACCESS_TOKEN",
		"AZDO_ORG_SERVICE_URL",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Fatalf("`%s` must be set for acceptance tests!", variable)
		}
	}
}
