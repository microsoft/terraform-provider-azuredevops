package pipelineschecksextras

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/pipelineschecks"
)

var ResourceAreaId, _ = uuid.Parse("4a933897-0488-45af-bd82-6fd3ad33f46a")

type Client interface {
	// [Preview API] Get Check configuration by Id
	GetCheckConfiguration(context.Context, pipelineschecks.GetCheckConfigurationArgs) (*pipelineschecks.CheckConfiguration, error)
}

type ClientImpl struct {
	Client azuredevops.Client
}

func NewClient(ctx context.Context, connection *azuredevops.Connection) (Client, error) {
	client, err := connection.GetClientByResourceAreaId(ctx, ResourceAreaId)
	if err != nil {
		return nil, err
	}
	return &ClientImpl{
		Client: *client,
	}, nil
}

// [Preview API] Get Check configuration by Id
func (client *ClientImpl) GetCheckConfiguration(ctx context.Context, args pipelineschecks.GetCheckConfigurationArgs) (*pipelineschecks.CheckConfiguration, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Id"}
	}
	routeValues["id"] = strconv.Itoa(*args.Id)

	queryParams := url.Values{}
	queryParams.Add("$expand", "1")

	locationId, _ := uuid.Parse("86c8381e-5aee-4cde-8ae4-25c0c7f5eaea")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "5.1-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue pipelineschecks.CheckConfiguration
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}
