# Testing

- [Authoring tests](#authoring-tests)
- [Testing](#testing)
- [Unit Tests](#unit-tests)
- [Acceptance Tests](#acceptance-tests)

Because Terraform plugins are written in Go, unit and integration tests are written using the standard [Go Test](https://golang.org/pkg/testing/) package. The basics of `go test` are not covered in this document but there are many great samples to be found online using your favorite search engine.

Instead, this document focuses on what makes testing for this project unique.

> Note: When naming your unit & acceptance tests, please follow the guidance from Hashicorp found [here](https://www.terraform.io/docs/extend/testing/unit-testing.html).

# Authoring Tests

The Azure DevOps provider applies an approach to separate and group tests by using GO build tags or build constraints. [GO build constraints](https://golang.org/pkg/go/build/#hdr-Build_Constraints).

Thus each `_test.go` files must include a build tag with the following characteristics:

1. The ``// +build`` constraint must include a tag named **all**.
2. The ``// +build`` constraint must include a tag named after the terraform resource or data source which is under test in the specific `_test.go` file.
3. Other build tags can be added as will. The administrators of the Azure DevOps Terraform Provider reserve the right to assign certain tags in the future to organize tests into logical groups.

`_test.go` files which contain test helper routines **must not** include any build tag. Otherwise those routines aren't available during a test run because the GO compiler i.e. `go test` will only honor files that either contain the specified build tag or does not contain any build tag at all.

If HCL code must be created for performing acceptance tests, add a function to `azuredevops\utils\testhelper\hcl.go` and try to reuse existing definitions.

Furthermore use the `test-acc-` prefix for naming Terraform resources or data sources in all acceptance tests. It's preferred to reference the `testhelper.TestAccResourcePrefix` const instead of using strings in acceptance tests.

```go
func TestAccAzureGitRepo_CreateAndUpdate(t *testing.T) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	gitRepoNameFirst := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	gitRepoNameSecond := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	tfRepoNode := "azuredevops_git_repository.gitrepo"

    ...
}
```

# Unit Tests

**Running unit tests**

The unit tests are executed whenever `./scripts/build.sh` is run. This can be run locally, but will also be run on every automated build and will be a gate for any PR against this repository. The tests can also be run in isolation by running the following:

```bash
$ ./scripts/unittest.sh
```

To run only unit tests for a specific resource or data source add the name of the build tag for this Terraform object as parameter to the `unittest.sh` script.

```bash
$ ./scripts/unittest.sh resource_project
```

To run unit tests for multiple resources or data sources or for a logical group of tests you can specify multiple parameters to `unittest.sh`.

**Azure DevOps Client SDK Mocks**

This project has a strong dependency on Microsoft's [Azure DevOps Go SDK](https://github.com/microsoft/azure-devops-go-api). We can mock the behavior of the SDK in our unit tests by using [GoMock](https://github.com/golang/mock), a popular mocking library for Go. This tool allows us to validate business logic against different success/failure modes AzDO services.

In order to use [GoMock](https://github.com/golang/mock) to mock an AzDO SDK, we must first generate a mock for that client. If you are mocking a client already used by the project then it is likely that the mock already exists. Otherwise, you can generate it yourself. The following command will auto-detect all AzDO Go SDKs used by the project and attempt to generate a mock for that client.

```bash
$ ./scripts/generate-mocks.sh
```

**Writing a test using a mock**

There is great documentation on [GoMock's GitHub](https://github.com/golang/mock), but here is test that validates that an error is not swallowed in a certain API failure mode:

> Note: GoMock (and Go in general) is quite verbose!

![Go Mock Example](https://user-images.githubusercontent.com/2497673/67523231-dbc05e00-f673-11e9-91c6-68a6684b3015.png)

Here are some important details:
 - **Lines 92-93**: Defers the call to `Finish()`, which will verify that each expectation (see lines 102-106) set up on the mocks used by the test was met. This will be done *after* the function exits. See [defer behavior in Go](https://tour.golang.org/flowcontrol/12)
 - **Lines 95-96**: Set up test data
 - **Lines 98-99**: Configure mock client(s)
 - **Lines 102-106**: Set an expectation for the mock. In this case, the expectation is that the `CreateDefinition` API will be called. If it is, it will return the specified parameters.
 - **Lines 108-109**: Test response from business logic

# Acceptance Tests

**Running acceptance tests**

> Note: Running acceptance tests provisions and deletes actual resources in AzDO. This can cost money and can be dangerous if you are not running them in isolation!

Integration tests for terraform providers are typically implemented as [Acceptance Tests](https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html). They have a special prefix - `TestAcc` - and will only be run when the `TEST_ACC` environment variable is set. They also rely on some environment variables. The following steps will run configure and run the acceptance tests:

```bash
# AZDO_ORG_SERVICE_URL will be the URL of the AzDO org that you want to provison
# resources inside of.
#   ex: https://dev.azure.com/<your org name>
$ export AZDO_ORG_SERVICE_URL="..."

# AZDO_PERSONAL_ACCESS_TOKEN will be the personal access token that grants access
# to provision and manage resources in Azure DevOps.
#   documentation: https://docs.microsoft.com/en-us/azure/devops/organizations/accounts/use-personal-access-tokens-to-authenticate?view=azure-devops
$ export AZDO_PERSONAL_ACCESS_TOKEN="..."

# Note: AZDO_GITHUB_SERVICE_CONNECTION_PAT is not specifically required
# by the provider, but it is required by the acceptance tests in order to test
# the authentication with GitHub for build definitions hosted in GitHub.
#   documentation: https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line
$ export AZDO_GITHUB_SERVICE_CONNECTION_PAT="..."

$ ./scripts/acctest.sh
```

To run only acceptance tests for a specific resource or data source add the name of the build tag for this Terraform object as parameter to the `acctest.sh` script.

```bash
$ ./scripts/acctest.sh resource_project
```

To run acceptance tests for multiple resources or data sources or for a logical group of tests you can specify multiple parameters to `acctest.sh`.

**Writing an acceptance test**

> Note: The established integration testing pattern for Terraform Providers is to write [Acceptance Tests](https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html). The process is well defined but is complicated. Get started by reading through the excellent [guide](https://www.terraform.io/docs/extend/testing/acceptance-tests/testcase.html) published by Hashicorp.

![Acceptance Test Example](https://user-images.githubusercontent.com/2497673/67523941-49b95500-f675-11e9-8345-21bda99ff1a4.png)

Here are some important details:
 - **Lines 190-192**: Set up resource names. The common prefix, `testAccResourcePrefix`, is used so that it is easy to identify any orphaned test resources in AzDO. This is defined in [provider_test.go](../azuredevops/provider_test.go).
 - **Line 196**: `PreCheck` is a function that verifies that the required environment variables are set. This is configured in [provider_test.go](../azuredevops/provider_test.go).
 - **Line 197**: `Providers` is the actual set of providers being tested. In this case, it is a fully configured `azuredevops` provider. This is configured in [provider_test.go](../azuredevops/provider_test.go).
 - **Line 198**: `CheckDestroy` checks that, after a `terraform destroy` is called, that the resource is actually destroyed from AzDO.
 - **Lines 199-217**: `Steps` is a list of steps that should be run. Each step will execute a `terraform apply` to apply the terraform stanza defined by the `Config` property. It then runs the checks specified by the `Check` property.
