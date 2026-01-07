package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/internal/client"
)

type TestResource interface {
	Exists(ctx context.Context, client *client.Client, state *terraform.InstanceState) (bool, error)
}
