package client

import (
	"context"
	"fmt"
	"os"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/dashboard"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/elastic"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/extensionmanagement"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/featuremanagement"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/feed"
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
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/wiki"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/version"
)

type Client struct {
	OrganizationURL               string
	CoreClient                    core.Client
	BuildClient                   build.Client
	DashboardClient               dashboard.Client
	PipelinesClient               pipelines.Client
	GitReposClient                git.Client
	GraphClient                   graph.Client
	OperationsClient              operations.Client
	PipelinesChecksClient         pipelineschecks.Client
	PipelinePermissionsClient     pipelinepermissions.Client
	PolicyClient                  policy.Client
	ElasticClient                 elastic.Client
	ExtensionManagementClient     extensionmanagement.Client
	ReleaseClient                 release.Client
	ServiceEndpointClient         serviceendpoint.Client
	TaskAgentClient               taskagent.Client
	MemberEntitleManagementClient memberentitlementmanagement.Client
	FeatureManagementClient       featuremanagement.Client
	FeedClient                    feed.Client
	SecurityClient                security.Client
	IdentityClient                identity.Client
	WikiClient                    wiki.Client
	WorkItemTrackingClient        workitemtracking.Client
	WorkItemTrackingProcessClient workitemtrackingprocess.Client
	ServiceHooksClient            servicehooks.Client
}

func New(ctx context.Context, authProvider azuredevops.AuthProvider, organizationURL string) (*Client, error) {
	userAgent := fmt.Sprintf("terraform-provider-azuredevops/%s", version.ProviderVersion)
	// append the CloudShell version to the user agent if it exists
	if azureAgent := os.Getenv("AZURE_HTTP_USER_AGENT"); azureAgent != "" {
		userAgent = fmt.Sprintf("%s %s", userAgent, azureAgent)
	}

	connection := &azuredevops.Connection{
		AuthProvider:            authProvider,
		BaseUrl:                 organizationURL,
		SuppressFedAuthRedirect: true,
		UserAgent:               userAgent,
	}

	coreClient, err := core.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("new core client: %v", err)
	}

	buildClient, err := build.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("new build client: %v", err)
	}

	operationsClient := operations.NewClient(ctx, connection)

	elasticClient := elastic.NewClient(ctx, connection)

	extensionManagementClient, err := extensionmanagement.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("new extension management client: %v", err)
	}

	dashboardClient, err := dashboard.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("new dashboard client: %v", err)
	}

	serviceEndpointClient, err := serviceendpoint.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("new service endpoint client: %v", err)
	}

	taskagentClient, err := taskagent.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("new task agent client: %v", err)
	}

	gitReposClient, err := git.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("new git client: %v", err)
	}

	graphClient, err := graph.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("new graph client: %v", err)
	}

	memberentitlementmanagementClient, err := memberentitlementmanagement.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("new member titlement management client: %v", err)
	}

	policyClient, err := policy.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("new policy client: %v", err)
	}

	releaseClient, err := release.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("new release client: %v", err)
	}

	securityClient := security.NewClient(ctx, connection)

	identityClient, err := identity.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("new identity client: %v", err)
	}

	wikiClient, err := wiki.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("new wiki client: %v", err)
	}

	featuremanagementClient := featuremanagement.NewClient(ctx, connection)

	feedClient, err := feed.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("new feed client: %v", err)
	}

	workitemtrackingClient, err := workitemtracking.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("new work item tracking client: %v", err)
	}

	workitemtrackingprocessClient, err := workitemtrackingprocess.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("new work item tracking process client: %v", err)
	}

	pipelines := pipelines.NewClient(ctx, connection)

	pipelinesChecksClient, err := pipelineschecks.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("new pipelines checks client: %v", err)
	}

	pipelinepermissionsClient, err := pipelinepermissions.NewClient(ctx, connection)
	if err != nil {
		return nil, fmt.Errorf("new pipeline permissions client: %v", err)
	}

	serviceHooksClient := servicehooks.NewClient(ctx, connection)

	client := &Client{
		OrganizationURL:               organizationURL,
		CoreClient:                    coreClient,
		BuildClient:                   buildClient,
		DashboardClient:               dashboardClient,
		ElasticClient:                 elasticClient,
		ExtensionManagementClient:     extensionManagementClient,
		GitReposClient:                gitReposClient,
		GraphClient:                   graphClient,
		OperationsClient:              operationsClient,
		PipelinesClient:               pipelines,
		PipelinesChecksClient:         pipelinesChecksClient,
		PipelinePermissionsClient:     pipelinepermissionsClient,
		PolicyClient:                  policyClient,
		ReleaseClient:                 releaseClient,
		ServiceEndpointClient:         serviceEndpointClient,
		TaskAgentClient:               taskagentClient,
		MemberEntitleManagementClient: memberentitlementmanagementClient,
		FeatureManagementClient:       featuremanagementClient,
		FeedClient:                    feedClient,
		SecurityClient:                securityClient,
		IdentityClient:                identityClient,
		WikiClient:                    wikiClient,
		WorkItemTrackingClient:        workitemtrackingClient,
		WorkItemTrackingProcessClient: workitemtrackingprocessClient,
		ServiceHooksClient:            serviceHooksClient,
	}
	return client, nil
}
