package testhelper

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

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
