package testclient

import (
	"os"
	"sync"

	"testing"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/terraform-provider-azuredevops/internal/client"
)

var (
	once sync.Once
	c    *client.Client
)

func New(t *testing.T) *client.Client {
	once.Do(func() {
		authProvider := azuredevops.NewAuthProviderPAT(os.Getenv("AZDO_PERSONAL_ACCESS_TOKEN"))
		client, err := client.New(t.Context(), authProvider, os.Getenv("AZDO_ORG_SERVICE_URL"))
		if err != nil {
			t.Fatalf("failed to new client: %v", err)
		}
		c = client
	})
	return c
}
