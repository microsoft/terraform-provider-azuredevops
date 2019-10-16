# Testing

* [Unit Tests](#unit-tests)
* [Integration Tests](#integration-tests)

# Unit Tests

## Mocks

The Azure DevOps Provider for Terraform has a strong dependency on the Azure DevOps Go API client. To test our providers, we utilize the GoMock mocking framework in addition to Go's built-in `testing` package. The GoMock framework allows us to mock any clients required from the Azure DevOps Go API (e.g., the `CoreClient` or `BuildClient`) so that we can isolate our unit test to the resource provider code by mocking operations of the Azure DevOps Go API client.

### Generating Mocks

The `generate-mocks.sh` script in the `scripts` directory can be run to generate the mocks required for testing this project. It should be run whenever a new AzDO SDK is pulled into the project. You may run the script either from the project root directory or from the scripts directory directly. The script will install GoMock if it has not been installed and will then automatically run `mockgen` to generate the mocks for every AzDO SDK that is a dependency of this project.

### Writing Unit Tests

When writing a unit test with a mocked client, GoMock will use the mocked clients from the `mockgen`-generated output. Additionally, `github.com/stretchr/testify/require` is used to stop test execution and fail if the assertion is not met.

When naming your unit tests (and acceptance tests), please follow the guidance from Hashicorp found [here](https://www.terraform.io/docs/extend/testing/unit-testing.html).

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
    "github.com/microsoft/terraform-provider-azuredevops/samplemocks"
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
    mockClient := samplemocks.NewMockSampleClient(controller)
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
