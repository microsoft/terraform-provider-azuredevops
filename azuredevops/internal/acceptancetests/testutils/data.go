package testutils

import (
	"fmt"
	"regexp"
)

func RequiresImportError(resourceName string) *regexp.Regexp {
	message := "Error creating service endpoint in Azure DevOps: Service connection with name %[1]s already exists. Only a user having Administrator/User role permissions on service connection %[1]s can see it."
	return regexp.MustCompile(fmt.Sprintf(message, resourceName))
}
