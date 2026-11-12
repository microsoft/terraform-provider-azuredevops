package organization

import (
	"context"
	"net/http"
	"net/url"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
)

// This API is not publicly released, and the client is generated based on the
// API:https://<orgName>.vssps.visualstudio.com/_apis/Organization/Collections/me

type Client interface {
	GetOrganization(ctx context.Context, organizationName string) (*Organization, error)
}

type ClientImpl struct {
	Client azuredevops.Client
}

func NewClient(ctx context.Context, connection *azuredevops.Connection) Client {
	client := connection.GetClientByUrl(connection.BaseUrl)
	return &ClientImpl{
		Client: *client,
	}
}

func (c ClientImpl) GetOrganization(ctx context.Context, organizationName string) (*Organization, error) {
	fullUrl := (&url.URL{
		Scheme: "https",
		Host:   organizationName + ".vssps.visualstudio.com",
		Path:   "/_apis/Organization/Collections/me",
	}).String()
	req, err := c.Client.CreateRequestMessage(ctx, http.MethodGet, fullUrl, "", nil, "application/json", "", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.SendRequest(req)
	if err != nil {
		return nil, err
	}

	var responseValue Organization
	err = c.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}
