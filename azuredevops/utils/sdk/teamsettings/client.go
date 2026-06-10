package teamsettings

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
)

var ResourceAreaId, _ = uuid.Parse("1d4f49f9-02b9-4e26-b826-2cdb6195f2a9")
var teamFieldValuesLocationId, _ = uuid.Parse("07ced576-58ed-49e6-9c1e-5cb53ab8bf2a")

type Client interface {
	GetTeamFieldValues(context.Context, GetTeamFieldValuesArgs) (*TeamFieldValues, error)
	UpdateTeamFieldValues(context.Context, UpdateTeamFieldValuesArgs) (*TeamFieldValues, error)
}

type ClientImpl struct {
	Client azuredevops.Client
}

func NewClient(ctx context.Context, connection *azuredevops.Connection) (Client, error) {
	client, err := connection.GetClientByResourceAreaId(ctx, ResourceAreaId)
	if err != nil {
		return nil, err
	}

	return &ClientImpl{Client: *client}, nil
}

func (client *ClientImpl) GetTeamFieldValues(ctx context.Context, args GetTeamFieldValuesArgs) (*TeamFieldValues, error) {
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	if args.Team == nil || *args.Team == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Team"}
	}
	routeValues["team"] = *args.Team

	resp, err := client.Client.Send(ctx, http.MethodGet, teamFieldValuesLocationId, "7.1-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue TeamFieldValues
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

func (client *ClientImpl) UpdateTeamFieldValues(ctx context.Context, args UpdateTeamFieldValuesArgs) (*TeamFieldValues, error) {
	if args.TeamFieldValues == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.TeamFieldValues"}
	}

	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project

	if args.Team == nil || *args.Team == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Team"}
	}
	routeValues["team"] = *args.Team

	body, marshalErr := json.Marshal(args.TeamFieldValues)
	if marshalErr != nil {
		return nil, marshalErr
	}

	resp, err := client.Client.Send(ctx, http.MethodPatch, teamFieldValuesLocationId, "7.1-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue TeamFieldValues
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}
