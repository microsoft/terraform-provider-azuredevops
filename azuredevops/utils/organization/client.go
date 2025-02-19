package organization

import (
	"context"
	"fmt"
	"net/http"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
)

// This API is not publicly released, and the client is generated based on the
// API:https://<orgName>.vssps.visualstudio.com/_apis/Organization/Collections/me

const baseUrl = "https://%s.vssps.visualstudio.com/_apis/Organization/Collections/me"

type Client interface {
	GetOrganizationMeta(ctx context.Context, organizationName string) (*OrganizationMeta, error)
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

func (c ClientImpl) GetOrganizationMeta(ctx context.Context, organizationName string) (*OrganizationMeta, error) {
	fullUrl := fmt.Sprintf(baseUrl, organizationName)
	req, err := c.Client.CreateRequestMessage(ctx, http.MethodGet, fullUrl, "", nil, "application/json", "", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.SendRequest(req)
	if err != nil {
		return nil, err
	}

	var responseValue OrganizationMeta
	err = c.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}
