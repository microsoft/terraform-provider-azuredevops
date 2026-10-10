// The Azure DevOps Go SDK does not contain an Advanced Security client.
// This file cannot be under "internal", because azdosdkmocks/advancedsecurity_sdk_mock.go depends on it.

package advancedsecurity

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
)

// ResourceAreaId is the "Management" resource area of Advanced Security.
var ResourceAreaId, _ = uuid.Parse("f101720c-9790-45a6-9fb3-494a09fddeeb") //nolint:errcheck

var repoEnablementLocationId, _ = uuid.Parse("d11a1c2b-b904-43dc-b970-bf42486262db") //nolint:errcheck

const apiVersion = "7.2-preview.3"

type Client interface {
	// [Preview API] Get the Advanced Security enablement of a repository.
	GetRepoEnablement(context.Context, GetRepoEnablementArgs) (*RepoEnablementSettings, error)
	// [Preview API] Update the Advanced Security enablement of a repository.
	UpdateRepoEnablement(context.Context, UpdateRepoEnablementArgs) error
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

func repoRouteValues(project, repository *string) (map[string]string, error) {
	if project == nil || *project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	if repository == nil || *repository == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Repository"}
	}
	return map[string]string{
		"project":    *project,
		"repository": *repository,
	}, nil
}

// [Preview API] Get the Advanced Security enablement of a repository.
func (client *ClientImpl) GetRepoEnablement(ctx context.Context, args GetRepoEnablementArgs) (*RepoEnablementSettings, error) {
	routeValues, err := repoRouteValues(args.Project, args.Repository)
	if err != nil {
		return nil, err
	}

	queryParams := url.Values{}
	if args.IncludeAllProperties != nil {
		queryParams.Add("includeAllProperties", strconv.FormatBool(*args.IncludeAllProperties))
	}

	resp, err := client.Client.Send(ctx, http.MethodGet, repoEnablementLocationId, apiVersion, routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue RepoEnablementSettings
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// [Preview API] Update the Advanced Security enablement of a repository.
func (client *ClientImpl) UpdateRepoEnablement(ctx context.Context, args UpdateRepoEnablementArgs) error {
	if args.Settings == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.Settings"}
	}
	routeValues, err := repoRouteValues(args.Project, args.Repository)
	if err != nil {
		return err
	}

	body, err := json.Marshal(args.Settings)
	if err != nil {
		return err
	}

	_, err = client.Client.Send(ctx, http.MethodPatch, repoEnablementLocationId, apiVersion, routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	return err
}
