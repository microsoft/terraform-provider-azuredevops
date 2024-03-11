package sdk

import (
	"strings"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
)

// Creates a new Azure DevOps connection instance using a function that returns an authorization header string.
func NewDynamicAuthorizationConnection(organizationUrl string, authProvider func() (string, error)) (*azuredevops.Connection, error) {
	organizationUrl = strings.ToLower(strings.TrimRight(organizationUrl, "/"))
	authorizationString, err := authProvider()
	if err != nil {
		return nil, err
	}
	return &azuredevops.Connection{
		AuthorizationString:     authorizationString,
		BaseUrl:                 organizationUrl,
		SuppressFedAuthRedirect: true,
	}, nil
}
