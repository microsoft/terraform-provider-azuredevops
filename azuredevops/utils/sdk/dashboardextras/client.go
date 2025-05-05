// This is a copy of github.com/microsoft/azure-devops-go-api/azuredevops/pipelineschecks/client.go
// The existing version does not contain the "Timeout" property on the CheckConfiguration struct

// This file cannot be under "internal", because azdosdkmocks/pipelines_checks_v5_extras_mock.go depends on it.

package dashboardextras

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/dashboard"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils"
)

var ResourceAreaId, _ = uuid.Parse("31c84e0a-3ece-48fd-a29d-100849af99ba") //nolint:errcheck

type Client interface {
	// [Preview API] Update the supplied dashboard.
	UpdateDashboard(context.Context, UpdateDashboardArgs) (*dashboard.Dashboard, error)
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

// [Preview API] Update a dashboard
func (client *ClientImpl) UpdateDashboard(ctx context.Context, args UpdateDashboardArgs) (*dashboard.Dashboard, error) {
	if args.Dashboard == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Dashboard"}
	}

	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	if args.Dashboard.Id == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Dashboard.Id"}
	}
	routeValues["dashboardId"] = (*args.Dashboard.Id).String()

	if args.Team != nil {
		routeValues["team"] = *args.Team
	}

	body, marshalErr := json.Marshal(args.Dashboard)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("454b3e51-2e6e-48d4-ad81-978154089351") //nolint:errcheck
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, utils.ApiVersion, routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue dashboard.Dashboard
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}
