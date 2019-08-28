package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
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

// Provider - The top level Azure DevOps Provider definition.
func Provider() *schema.Provider {
	p := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			// Why is the key/value named the way they are?
			"azuredevops_foo": resourceFoo(),
		},
		Schema: map[string]*schema.Schema{
			"org_service_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AZDO_ORG_SERVICE_URL", nil),
				Description: "The url of the Azure DevOps instance which should be used.",
			},
			"personal_access_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AZDO_PERSONAL_ACCESS_TOKEN", nil),
				Description: "The personal access token which should be used.",
			},
		},
	}

	p.ConfigureFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		type AzdoConfig struct {
			azdoPAT         string
			organizationURL string
			ctx             context.Context
		}

		azdoConfig := &AzdoConfig{
			azdoPAT:         d.Get("personal_access_token").(string),
			organizationURL: d.Get("org_service_url").(string),
			ctx:             context.Background(),
		}

		if azdoConfig.azdoPAT == "" {
			return nil, fmt.Errorf("the personal access token is required")
		}

		if azdoConfig.organizationURL == "" {
			return nil, fmt.Errorf("the url of the Azure DevOps is required")
		}

		connection := azuredevops.NewPatConnection(azdoConfig.organizationURL, azdoConfig.azdoPAT)

		// client for these APIs (includes CRUD for AzDO projects...):
		//	https://docs.microsoft.com/en-us/rest/api/azure/devops/core/?view=azure-devops-rest-5.1
		coreClient, err := core.NewClient(azdoConfig.ctx, connection)
		if err != nil {
			return nil, err
		}

		// client for these APIs (includes CRUD for AzDO build pipelines...):
		//	https://docs.microsoft.com/en-us/rest/api/azure/devops/build/?view=azure-devops-rest-5.1
		buildClient, err := build.NewClient(azdoConfig.ctx, connection)
		if err != nil {
			return nil, err
		}

		// client for these APIs (monitor async operations...):
		//	https://docs.microsoft.com/en-us/rest/api/azure/devops/operations/operations?view=azure-devops-rest-5.1
		operationsClient := operations.NewClient(azdoConfig.ctx, connection)

		aggregatedClient := &aggregatedClient{
			CoreClient:       coreClient,
			BuildClient:      buildClient,
			OperationsClient: operationsClient,
			ctx:              azdoConfig.ctx,
		}

		log.Printf("Created clients successfully!")
		return aggregatedClient, nil
	}
}
