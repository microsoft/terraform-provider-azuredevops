package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/elastic"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/featuremanagement"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/memberentitlementmanagement"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/operations"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/pipelinepermissions"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/pipelines"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/pipelineschecks"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/policy"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/release"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/security"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/servicehooks"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/pipelineschecksextras"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/sdk"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/securityroles"
	"github.com/microsoft/terraform-provider-azuredevops/version"
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
	PipelinesClient               pipelines.Client
	GitReposClient                git.Client
	GraphClient                   graph.Client
	OperationsClient              operations.Client
	PipelinesChecksClient         pipelineschecks.Client
	PipelinePermissionsClient     pipelinepermissions.Client
	PipelinesChecksClientExtras   pipelineschecksextras.Client
	PolicyClient                  policy.Client
	ElasticClient                 elastic.Client
	ReleaseClient                 release.Client
	ServiceEndpointClient         serviceendpoint.Client
	TaskAgentClient               taskagent.Client
	MemberEntitleManagementClient memberentitlementmanagement.Client
	FeatureManagementClient       featuremanagement.Client
	SecurityClient                security.Client
	IdentityClient                identity.Client
	WorkItemTrackingClient        workitemtracking.Client
	ServiceHooksClient            servicehooks.Client
	Ctx                           context.Context
	SecurityRolesClient           securityroles.Client
}

// GetAzdoClient builds and provides a connection to the Azure DevOps API
func GetAzdoClient(azdoTokenProvider func() (string, error), organizationURL string, tfVersion string) (*AggregatedClient, error) {
	ctx := context.Background()

	if strings.EqualFold(organizationURL, "") {
		return nil, fmt.Errorf("the url of the Azure DevOps is required")
	}

	connection, err := sdk.NewDynamicAuthorizationConnection(organizationURL, azdoTokenProvider)
	if err != nil {
		return nil, err
	}
	setUserAgent(connection, tfVersion)

	coreClient, err := core.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): core.NewClient failed.")
		return nil, err
	}

	buildClient, err := build.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): build.NewClient failed.")
		return nil, err
	}

	operationsClient := operations.NewClient(ctx, connection)

	elasticClient := elastic.NewClient(ctx, connection)

	serviceEndpointClient, err := serviceendpoint.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): serviceendpoint.NewClient failed.")
		return nil, err
	}

	taskagentClient, err := taskagent.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): taskagent.NewClient failed.")
		return nil, err
	}

	gitReposClient, err := git.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): git.NewClient failed.")
		return nil, err
	}

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

	policyClient, err := policy.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): policy.NewClient failed.")
		return nil, err
	}

	releaseClient, err := release.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): release.NewClient failed.")
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

	pipelines := pipelines.NewClient(ctx, connection)

	pipelinesChecksClient, err := pipelineschecks.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): pipelineschecks.NewClient failed.")
		return nil, err
	}

	pipelinepermissionsClient, err := pipelinepermissions.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): pipelineschecks.NewClient failed.")
		return nil, err
	}

	pipelinesChecksClientExtras, err := pipelineschecksextras.NewClient(ctx, connection)
	if err != nil {
		log.Printf("getAzdoClient(): pipelineschecksextras.NewClient failed.")
		return nil, err
	}

	serviceHooksClient := servicehooks.NewClient(ctx, connection)

	securityRolesClient := securityroles.NewClient(ctx, connection)

	aggregatedClient := &AggregatedClient{
		OrganizationURL:               organizationURL,
		CoreClient:                    coreClient,
		BuildClient:                   buildClient,
		ElasticClient:                 elasticClient,
		GitReposClient:                gitReposClient,
		GraphClient:                   graphClient,
		OperationsClient:              operationsClient,
		PipelinesClient:               pipelines,
		PipelinesChecksClient:         pipelinesChecksClient,
		PipelinePermissionsClient:     pipelinepermissionsClient,
		PipelinesChecksClientExtras:   pipelinesChecksClientExtras,
		PolicyClient:                  policyClient,
		ReleaseClient:                 releaseClient,
		ServiceEndpointClient:         serviceEndpointClient,
		TaskAgentClient:               taskagentClient,
		MemberEntitleManagementClient: memberentitlementmanagementClient,
		FeatureManagementClient:       featuremanagementClient,
		SecurityClient:                securityClient,
		IdentityClient:                identityClient,
		WorkItemTrackingClient:        workitemtrackingClient,
		ServiceHooksClient:            serviceHooksClient,
		SecurityRolesClient:           securityRolesClient,
		Ctx:                           ctx,
	}

	log.Printf("getAzdoClient(): Created core, build, operations, and serviceendpoint clients successfully!")
	return aggregatedClient, nil
}

// setUserAgent set UserAgent for http headers
func setUserAgent(connection *azuredevops.Connection, tfVersion string) {
	providerUserAgent := fmt.Sprintf("terraform-provider-azuredevops/%s", version.ProviderVersion)
	connection.UserAgent = strings.TrimSpace(fmt.Sprintf("%s %s", connection.UserAgent, providerUserAgent))

	// append the CloudShell version to the user agent if it exists
	if azureAgent := os.Getenv("AZURE_HTTP_USER_AGENT"); azureAgent != "" {
		connection.UserAgent = fmt.Sprintf("%s %s", connection.UserAgent, azureAgent)
	}

	log.Printf("[DEBUG] AzureRM Client User Agent: %s\n", connection.UserAgent)
}
