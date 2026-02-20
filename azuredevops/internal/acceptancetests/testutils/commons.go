package testutils

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// initialized once, so it can be shared by each acceptance test
var provider = azuredevops.Provider()

// GetProvider returns the azuredevops provider
func GetProvider() *schema.Provider {
	return provider
}

func GetProviderFactories() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		//nolint:unparam
		"azuredevops": func() (*schema.Provider, error) {
			return GetProvider(), nil
		},
	}
}

// GetProviders returns a map of all providers needed for the project
func GetProviders() map[string]*schema.Provider {
	return map[string]*schema.Provider{
		"azuredevops": GetProvider(),
	}
}

// ComputeProjectQualifiedResourceImportID returns a function that can be used to construct an import ID of a resource
// that has an import ID in the following form: <project ID>/<resource ID>
func ComputeProjectQualifiedResourceImportID(resourceNode string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceNode]
		if !ok {
			return "", fmt.Errorf("Resource node not found: %s", resourceNode)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["id"]), nil
	}
}

// PreCheck checks that the requisite environment variables are set
func PreCheck(t *testing.T, additionalEnvVars *[]string) {
	requiredEnvVars := []string{
		"AZDO_ORG_SERVICE_URL",
		"AZDO_PERSONAL_ACCESS_TOKEN",
	}
	if additionalEnvVars != nil {
		requiredEnvVars = append(requiredEnvVars, *additionalEnvVars...)
	}
	missing := false
	for _, variable := range requiredEnvVars {
		if _, ok := os.LookupEnv(variable); !ok {
			missing = true
			t.Errorf("`%s` must be set for this acceptance test!", variable)
		}
	}
	if missing {
		t.Fatalf("Some environment variables missing.")
	}
}

// GenerateResourceName generates a random name with a constant prefix, useful for acceptance tests
func GenerateResourceName() string {
	return "test-acc-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
}

// CheckNestedKeyExistsWithValue checks if a property exists with a certain value in an instance state
func CheckNestedKeyExistsWithValue(tfNode string, propertyName string, propertyValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rootModule := s.RootModule()
		resource, ok := rootModule.Resources[tfNode]
		if !ok {
			return fmt.Errorf("Did not find a project in the TF state")
		}

		is := resource.Primary
		if is == nil {
			return fmt.Errorf("No primary instance: %s in %s", tfNode, rootModule.Path)
		}
		if !containsPropertyWithValue(is.Attributes, propertyName, propertyValue) {
			return fmt.Errorf("%s does not contain a pool with %s %s", tfNode, propertyName, propertyValue)
		}
		return nil
	}
}

func containsPropertyWithValue(m map[string]string, property string, value string) bool {
	for k, v := range m {
		if v == value && k[strings.LastIndex(k, ".")+1:] == property {
			return true
		}
	}
	return false
}

func RunTestsInSequence(t *testing.T, tests map[string]map[string]func(t *testing.T)) {
	for group, m := range tests {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

// TestCheckAttr returns a TestCheckFunc that checks the attribute value at key using a predicate function
func TestCheckAttr(addr, key string, fn func(string) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[addr]
		if !ok {
			return fmt.Errorf("resource %s not found", addr)
		}
		v, ok := rs.Primary.Attributes[key]
		if !ok {
			return fmt.Errorf("attribute %s not found", key)
		}
		return fn(v)
	}
}

// TestCheckAttrGreaterThan returns a TestCheckFunc that checks that the attribute value at key is greater than count
func TestCheckAttrGreaterThan(addr, key string, count int) resource.TestCheckFunc {
	return TestCheckAttr(addr, key, func(v string) error {
		actualCount, err := converter.ASCIIToIntPtr(v)
		if err != nil {
			return err
		}
		if *actualCount <= count {
			return fmt.Errorf("expected %s to be greater than %d, got %d", key, count, *actualCount)
		}
		return nil
	})
}
