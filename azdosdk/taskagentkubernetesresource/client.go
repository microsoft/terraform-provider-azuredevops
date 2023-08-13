package taskagentkubernetesresource

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
)

// Client contains a re-implementation of broken taskagent SDK implementations
type Client interface {
	AddKubernetesResource(context.Context, AddKubernetesResourceArgs) (*taskagent.KubernetesResource, error)
}

type ClientImpl struct {
	Client azuredevops.Client
}

func NewTaskAgentKubernetesResourceClient(ctx context.Context, connection *azuredevops.Connection) (Client, error) {
	client, err := connection.GetClientByResourceAreaId(ctx, taskagent.ResourceAreaId)
	if err != nil {
		return nil, err
	}
	return &ClientImpl{
		Client: *client,
	}, nil
}

// AddKubernetesResource adds a kubernetes resource to an environment.
//
// This re-implementation uses a AddKubernetesResourceArgs rather than a taskagent.AddKubernetesResourceArgs to
// allow for passing the required serviceEndpointId parameter.
func (client *ClientImpl) AddKubernetesResource(ctx context.Context, args AddKubernetesResourceArgs) (*taskagent.KubernetesResource, error) {
	if args.CreateParameters == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.CreateParameters"}
	}
	routeValues := make(map[string]string)
	if args.Project == nil || *args.Project == "" {
		return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
	}
	routeValues["project"] = *args.Project
	if args.EnvironmentId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.EnvironmentId"}
	}
	routeValues["environmentId"] = strconv.Itoa(*args.EnvironmentId)

	body, marshalErr := json.Marshal(*args.CreateParameters)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("73fba52f-15ab-42b3-a538-ce67a9223a04")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0-preview.1", routeValues, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue taskagent.KubernetesResource
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// AddKubernetesResourceArgs is a reimplementation of taskagent.AddKubernetesResourceArgs which takes a
// taskagent.KubernetesResourceCreateParametersExistingEndpoint rather than a
// taskagent.KubernetesResourceCreateParameters
type AddKubernetesResourceArgs struct {
	// (required)
	CreateParameters *taskagent.KubernetesResourceCreateParametersExistingEndpoint
	// (required) Project ID or project name
	Project *string
	// (required)
	EnvironmentId *int
}
