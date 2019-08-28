package main

import (
	"context"
	"fmt"
	"log"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/operations"
)

// Aggregates all of the underlying clients into a single data
// type. Each client is ready to use and fully configured with the correct
// AzDO PAT/organization
type aggregatedClient struct {
	CoreClient       *core.Client
	BuildClient      *build.Client
	OperationsClient *operations.Client
	ctx              context.Context
}

func getAzdoClient(azdoPAT string, organizationURL string) (*aggregatedClient, error) {
	ctx := context.Background()

	if azdoPAT == "" {
		return nil, fmt.Errorf("the personal access token is required")
	}

	if organizationURL == "" {
		return nil, fmt.Errorf("the url of the Azure DevOps is required")
	}

	connection := azuredevops.NewPatConnection(organizationURL, azdoPAT)

	// client for these APIs (includes CRUD for AzDO projects...):
	//	https://docs.microsoft.com/en-us/rest/api/azure/devops/core/?view=azure-devops-rest-5.1
	coreClient, err := core.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): core.NewClient failed.")
		return nil, err
	}

	// client for these APIs (includes CRUD for AzDO build pipelines...):
	//	https://docs.microsoft.com/en-us/rest/api/azure/devops/build/?view=azure-devops-rest-5.1
	buildClient, err := build.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): build.NewClient failed.")
		return nil, err
	}

	// client for these APIs (monitor async operations...):
	//	https://docs.microsoft.com/en-us/rest/api/azure/devops/operations/operations?view=azure-devops-rest-5.1
	operationsClient := operations.NewClient(ctx, connection)

	aggregatedClient := &aggregatedClient{
		CoreClient:       coreClient,
		BuildClient:      buildClient,
		OperationsClient: operationsClient,
		ctx:              ctx,
	}

	log.Printf("getAzdoClient(): Created core, build, and operations clients successfully!")
	return aggregatedClient, nil
}
