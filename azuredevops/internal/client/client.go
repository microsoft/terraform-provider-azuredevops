package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/featuremanagement"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/identity"
	"github.com/microsoft/azure-devops-go-api/azuredevops/memberentitlementmanagement"
	"github.com/microsoft/azure-devops-go-api/azuredevops/operations"
	"github.com/microsoft/azure-devops-go-api/azuredevops/policy"
	"github.com/microsoft/azure-devops-go-api/azuredevops/security"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
	"github.com/terraform-providers/terraform-provider-azuredevops/version"
)

// AggregatedClient aggregates all of the underlying clients into a single data
// type. Each client is ready to use and fully configured with the correct
// AzDO PAT/organization
//
// AggregatedClient uses interfaces derived from the underlying client structs to
// allow for mocking to support unit testing of the funcs that invoke the
// Azure DevOps client.
type AggregatedClient struct {
	OrganizationURL               string
	CoreClient                    core.Client
	BuildClient                   build.Client
	GitReposClient                git.Client
	GraphClient                   graph.Client
	OperationsClient              operations.Client
	PolicyClient                  policy.Client
	ServiceEndpointClient         serviceendpoint.Client
	TaskAgentClient               taskagent.Client
	MemberEntitleManagementClient memberentitlementmanagement.Client
	FeatureManagementClient       featuremanagement.Client
	SecurityClient                security.Client
	IdentityClient                identity.Client
	WorkItemTrackingClient        workitemtracking.Client
	Ctx                           context.Context
}

// GetAzdoClient builds and provides a connection to the Azure DevOps API
func GetAzdoClient(azdoPAT string, organizationURL string, tfVersion string) (*AggregatedClient, error) {
	ctx := context.Background()

	if strings.EqualFold(azdoPAT, "") {
		return nil, fmt.Errorf("the personal access token is required")
	}

	if strings.EqualFold(organizationURL, "") {
		return nil, fmt.Errorf("the url of the Azure DevOps is required")
	}

	connection := azuredevops.NewPatConnection(organizationURL, azdoPAT)
	setUserAgent(connection, tfVersion)

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

	// client for these APIs (includes CRUD for AzDO service endpoints a.k.a. service connections...):
	//  https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-5.1
	serviceEndpointClient, err := serviceendpoint.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): serviceendpoint.NewClient failed.")
		return nil, err
	}

	// client for these APIs (includes CRUD for AzDO variable groups):
	taskagentClient, err := taskagent.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): taskagent.NewClient failed.")
		return nil, err
	}

	// client for these APIs:
	//	https://docs.microsoft.com/en-us/rest/api/azure/devops/git/?view=azure-devops-rest-5.1
	gitReposClient, err := git.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): git.NewClient failed.")
		return nil, err
	}

	//  https://docs.microsoft.com/en-us/rest/api/azure/devops/graph/?view=azure-devops-rest-5.1
	graphClient, err := graph.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): graph.NewClient failed.")
		return nil, err
	}

	memberentitlementmanagementClient, err := memberentitlementmanagement.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): memberentitlementmanagement.NewClient failed.")
		return nil, err
	}

	// https://docs.microsoft.com/en-us/rest/api/azure/devops/policy/configurations/create?view=azure-devops-rest-5.1
	policyClient, err := policy.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): policy.NewClient failed.")
		return nil, err
	}

	securityClient := security.NewClient(ctx, connection)
	identityClient, err := identity.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): identity.NewClient failed.")
		return nil, err
	}

	featuremanagementClient := featuremanagement.NewClient(ctx, connection)

	workitemtrackingClient, err := workitemtracking.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): workitemtracking.NewClient failed.")
		return nil, err
	}

	workitemtrackingClient, err := workitemtracking.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): workitemtracking.NewClient failed.")
		return nil, err
	}

	aggregatedClient := &AggregatedClient{
		OrganizationURL:               organizationURL,
		CoreClient:                    coreClient,
		BuildClient:                   buildClient,
		GitReposClient:                gitReposClient,
		GraphClient:                   graphClient,
		OperationsClient:              operationsClient,
		PolicyClient:                  policyClient,
		ServiceEndpointClient:         serviceEndpointClient,
		TaskAgentClient:               taskagentClient,
		MemberEntitleManagementClient: memberentitlementmanagementClient,
		FeatureManagementClient:       featuremanagementClient,
		SecurityClient:                securityClient,
		IdentityClient:                identityClient,
		WorkItemTrackingClient:        workitemtrackingClient,
		Ctx:                           ctx,
	}

	log.Printf("getAzdoClient(): Created core, build, operations, and serviceendpoint clients successfully!")
	return aggregatedClient, nil
}

// setUserAgent set UserAgent for http headers
func setUserAgent(connection *azuredevops.Connection, tfVersion string) {
	tfUserAgent := httpclient.TerraformUserAgent(tfVersion)
	providerUserAgent := fmt.Sprintf("%s terraform-provider-azuredevops/%s", tfUserAgent, version.ProviderVersion)
	connection.UserAgent = strings.TrimSpace(fmt.Sprintf("%s %s", connection.UserAgent, providerUserAgent))

	// append the CloudShell version to the user agent if it exists
	if azureAgent := os.Getenv("AZURE_HTTP_USER_AGENT"); azureAgent != "" {
		connection.UserAgent = fmt.Sprintf("%s %s", connection.UserAgent, azureAgent)
	}

	log.Printf("[DEBUG] AzureRM Client User Agent: %s\n", connection.UserAgent)
}
