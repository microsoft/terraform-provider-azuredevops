package sdk

import (
	"strings"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
)

func NewConnection(organizationUrl string, authProvider azuredevops.AuthProvider) (*azuredevops.Connection, error) {
	organizationUrl = strings.ToLower(strings.TrimRight(organizationUrl, "/"))
	return &azuredevops.Connection{
		AuthProvider:            authProvider,
		BaseUrl:                 organizationUrl,
		SuppressFedAuthRedirect: true,
	}, nil
}
