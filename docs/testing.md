# Testing

* [Unit Tests](#unit-tests)
* [Integration Tests](#integration-tests)

# Unit Tests



## Mocks

The Azure DevOps Provider for Terraform has a strong dependency on the Azure DevOps Go API client. To test our providers, we utilize the GoMock mocking framework in addition to Go's built-in `testing` package. The GoMock framework allows us to mock any clients required from the Azure DevOps Go API (e.g., the `CoreClient` or `BuildClient`) so that we can isolate our unit test to the resource provider code by mocking operations of the Azure DevOps Go API client.

> The `mockgen` tool works with interfaces; however, the Azure DevOps Go API clients are not interfaces. There is ongoing work to generate interfaces for those types. In the meantime, this can be worked around by using the tool [`ifacemaker`](https://github.com/vburenin/ifacemaker) to pull the interfaces from the relevant `client.go` file in the Azure DevOps Go API client source (e.g., the `BuildClient` interface was pulled from https://github.com/microsoft/azure-devops-go-api/blob/dev/azuredevops/build/client.go). Some additional manual modification is required as well to ensure proper namespacing of the generated models. See [`config.go`](https://github.com/microsoft/terraform-provider-azuredevops/blob/master/azuredevops/config.go) for examples that produce proper output from the `mockgen` tool.

### Generating Mocks

The `generate-mocks.sh` script in the `scripts` directory can be run to generate the mocks required for testing this project. It should be updated if additional mocks are required and/or if changes to the existing interfaces are made that would require new mocks be generated. You may run the script either from the project root directory or from the scripts directory directly. The script will install GoMock if it has not been installed and will then automatically run `mockgen` for known source/destination files.

The script includes the following steps that might be manually run to support mocking:

1. Installing GoMock Tools

   Install the `mockgen` tool by running:
   ```sh
   go get github.com/golang/mock/mockgen
   ```

2. Running `mockgen`

   `mockgen` generates mock interfaces from a specified source file. The `config.go` source file contains an aggregation of all of the Azure DevOps Go API clients used by the provider. To generate the mock interfaces for this, run:
   ```sh
   mockgen -source=config.go -destination=mock_config.go
   ```

   > IMPORTANT: The mock clients should be regenerated when making any changes to the `config.go` that affect the interfaces. This is a manual step at this time.

### Writing Unit Tests

When writing a unit test with a mocked client, GoMock will use the mocked clients from the `mockgen`-generated output. Additionally, `github.com/stretchr/testify/require` is used to stop test execution and fail if the assertion is not met.

#### SampleClient interface
```go
type SampleClient interface {
    SampleOperation(arg string) (string, error)
}
```

#### SampleClient usage
```go
func foo(client SampleClient, string bar) (string, error) {
    return client.SampleOperation(bar)
}
```

#### Unit Test example with mocked SampleClient

```go
import {
    "errors"
    "testing"

    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/require"
}

func foo_ReturnsErrorWhenGivenNil(t *testing.T) {
    expectedErrorDescription := "Argument Nil"

    // Setup GoMock Controller to defer until assertions may be invoked
    controller := gomock.NewController(t)
    defer controller.Finish()

    // Setup mock client to handle assertion that Client.Operation with nil
    // arguments will return an error and be called exactly once.
    mockClient := NewMockSampleClient(controller)
    mockClient.
        EXPECT().
        SampleOperation(gomock.Eq(gomock.Nil())).
        Return(nil, errors.New(expectedErrorDescription)).
        Times(1)

    // Execute the provider code, providing 'nil' as the argument
    _, actualError := foo(mockClient, nil)
    require.Equal(t, expectedErrorDescription, actualError.Error())
}
```

# Integration Tests

The established integration testing pattern for Terraform Providers is to write [Acceptance Tests](https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html). The process is well defined but can be a tad tricky to understand at fist. Given this, you may want to get started by reading through the excellent [guide](https://www.terraform.io/docs/extend/testing/acceptance-tests/testcase.html) published by Hashicorp.
