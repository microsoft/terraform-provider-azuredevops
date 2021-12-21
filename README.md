# Terraform Provider for Azure DevOps (Devops Resource Manager)

[![Gitter](https://badges.gitter.im/terraform-provider-azuredevops/community.svg)](https://gitter.im/terraform-provider-azuredevops/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/microsoft/terraform-provider-azuredevops)](https://goreportcard.com/report/github.com/microsoft/terraform-provider-azuredevops)

The AzureRM Provider supports Terraform 0.12.x and later.

* [Terraform Website](https://www.terraform.io)
* [Azure DevOps Website](https://azure.microsoft.com/en-us/services/devops/)
* [Provider Documentation](./website/docs/index.html.markdown)
* [Resources Documentation](./website/docs/r/)
* [Data Sources Documentation](./website/docs/d/)
* [Usage Examples](./examples/)
* [Gitter Channel](https://gitter.im/terraform-provider-azuredevops/community)

## Usage Example

```hcl
# Make sure to set the following environment variables:
#   AZDO_PERSONAL_ACCESS_TOKEN
#   AZDO_ORG_SERVICE_URL
terraform {
  required_providers {
    azuredevops = {
      source = "microsoft/azuredevops"
      version = ">=0.1.0"
    }
  }
}

resource "azuredevops_project" "project" {
  name = "My Awesome Project"
  description  = "All of my awesomee things"
}

resource "azuredevops_git_repository" "repository" {
  project_id = azuredevops_project.project.id
  name       = "My Awesome Repo"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_build_definition" "build_definition" {
  project_id = azuredevops_project.project.id
  name       = "My Awesome Build Pipeline"
  path       = "\\"

  repository {
    repo_type   = "TfsGit"
    repo_id     = azuredevops_git_repository.repository.id
    branch_name = azuredevops_git_repository.repository.default_branch
    yml_path    = "azure-pipelines.yml"
  }
}
```

## Developer Requirements

* [Terraform](https://www.terraform.io/downloads.html) version 0.13.x +
* [Go](https://golang.org/doc/install) version 1.16.x (to build the provider plugin)

If you're on Windows you'll also need:

* [Git for Windows](https://git-scm.com/download/win)

If you what to use the `makefile` build strategy on Windows it's required to install

* [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)

For *GNU32 Make*, make sure its bin path is added to PATH environment variable.*

For *Git Bash for Windows*, at the step of "Adjusting your PATH environment", please choose "Use Git and optional Unix tools from Windows Command Prompt".*

As [described below](#build-using-powerShell-scripts) we provide some PowerShell scripts to build the provider on Windows, without the requiremet to install any Unix based tools aside Go.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.16+ is **required**). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

### Using the GOPATH model

First clone the repository to: `$GOPATH/src/github.com/microsoft/terraform-provider-azuredevops`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers && cd "$_"
$ git clone git@github.com:microsoft/terraform-provider-azuredevops.git
$ cd terraform-provider-azuredevops
```

Once you've cloned, run the `./scripts/build.sh` and `./scripts/local-install.sh`, as recommended [here](https://github.com/microsoft/terraform-provider-azuredevops/blob/main/docs/contributing.md#3-build--install-provider).
These commands will sideload the plugin for Terraform.

### Using a directory separate from GOPATH

The infrastructure supports building and testing the provider outside `GOPATH` in an arbitrary directory.
In this scenario all required packages of the provider during build will be managed via the `pkg` in `$GOPATH`. As with the [GOPATH Model](#using-the-gopath-model), you can redefine the `GOPATH` environment variable to prevent existing packages in the current `GOPATH` directory from being changed.

### Build using make

Once inside the provider directory, you can run `make tools` to install the dependent tooling required to compile the provider.

At this point you can compile the provider by running `make build`, which will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-azuredevops
...
```

You can also cross-compile if necessary:

```sh
GOOS=windows GOARCH=amd64 make build
```

#### Unit tests

In order to run the Unit Tests for the provider, you can run:

```sh
$ make test
```

With VSCode Golang extension you can also run and debug the tests using `run test`, `debug test` `run package tests`, `run file tests` buttons.

#### Acceptance tests

The majority of tests in the provider are acceptance tests - which provisions real resources in Azure Devops and Azure. To run any acceptance tests you need to set `AZDO_ORG_SERVICE_URL`, `AZDO_PERSONAL_ACCESS_TOKEN` environment variables, some test have additional environment variables required to run. You can find out the required environment variables by running the test. Most of these variables can be set to dummy values.

The several options to run the tests are:

* Run the entire acceptance test suite

  ```sh
  make testacc
  ```

* Run a subset using a prefix

  ```sh
  make testacc TESTARGS='-run=TestAccBuildDefinitionBitbucket_Create' TESTTAGS='resource_build_definition'
  ```

* With VSCode Golang extension you can also run the tests using `run test`, `run package tests`, `run file tests` buttons above the test

### Scaffolding the Website Documentation

You can scaffold the documentation for a Data Source by running:

```sh
$ make scaffold-website BRAND_NAME="Agent Pool" RESOURCE_NAME="azuredevops_agent_pool" RESOURCE_TYPE="data"
```

You can scaffold the documentation for a Resource by running:

```sh
$ make scaffold-website BRAND_NAME="Agent Pool" RESOURCE_NAME="azuredevops_agent_pool" RESOURCE_TYPE="resource" RESOURCE_ID="00000000-0000-0000-0000-000000000000"
```

>
> `BRAND_NAME` is the human readable name of the object that is handled by a
> Terraform resource or datasource, like `Agent Pool`, `User Entitlement` or `Kubernetes Service Endpoint`
>

### Build using PowerShell scripts

If you like to develop on Windows, we provide a set of PowerShell scripts to build and test the provider.
They don't offer the luxury of a Makefile environment but are quite sufficient to develop on Windows.

#### `scripts\build.ps1`

The `build.ps1`is used to build the provider. Aside this the script runs (if not skipped) the defined unit tests and is able to install the compiled provider locally.

| Parameter   | Description                                                                               |
| ----------- | ----------------------------------------------------------------------------------------- |
| -SkipTests  | Skip running unit tests during build                                                      |
| -Install    | Install the provider locally, after a successful build                                    |
| -DebugBuild | Build the provider with extra debugging information                                       |
| -GoMod      | Control the `-mod` build parameter: Valid values: '' (Empty string), 'vendor', 'readonly' |

#### `scripts\unittest.ps1`

The script is used to execute unit tests. The script is also executed by `build.ps1` if the `-SkipTest` are not specified.

| Parameter   | Description                                                                                                                       |
| ----------- | --------------------------------------------------------------------------------------------------------------------------------- |
| -TestFilter | A GO regular expression which filters the test functions to be executed                                                           |
| -Tag        | Tests in the provider project are organized with GO build tags. The parameter accepts a list of tag names which should be tested. |
| -GoMod      | Control the `-mod` build parameter: Valid values: '' (Empty string), 'vendor', 'readonly'                                         |

#### `scripts\acctest.ps1`

The script is used to execute unit tests.

| Parameter   | Description                                                                                                                       |
| ----------- | --------------------------------------------------------------------------------------------------------------------------------- |
| -TestFilter | A GO regular expression which filters the test functions to be executed                                                           |
| -Tag        | Tests in the provider project are organized with GO build tags. The parameter accepts a list of tag names which should be tested. |
| -GoMod      | Control the `-mod` build parameter: Valid values: '' (Empty string), 'vendor', 'readonly'                                         |

#### `scripts\gofmtcheck.ps1`

To validate if all `.go` files adhere to the required formatting rules, execute `gofmtcheck.ps1`

| Parameter | Description                                                                                                    |
| --------- | -------------------------------------------------------------------------------------------------------------- |
| -Fix      | Fix any formatting rule deviations automatically. If the parameter is not set, the script runs in report mode. |

#### `scripts\lint-check-go.ps1`

Like with `gofmtcheck.ps1` the script validate if all `.go` files adhere to the required formatting rules and if any style mistakes exist. In difference to `gofmtcheck.ps1` the script uses Golint instead of Gofmt.

## Environment variables for acceptance tests

The following Environment Variables must be set in your shell prior to running acceptance tests:

- `AZDO_ORG_SERVICE_URL`
- `AZDO_PERSONAL_ACCESS_TOKEN`
- `AZDO_DOCKERREGISTRY_SERVICE_CONNECTION_EMAIL`
- `AZDO_DOCKERREGISTRY_SERVICE_CONNECTION_PASSWORD`
- `AZDO_DOCKERREGISTRY_SERVICE_CONNECTION_USERNAME`
- `AZDO_GITHUB_SERVICE_CONNECTION_PAT`
- `AZDO_TEST_AAD_USER_EMAIL`

**Note:** Acceptance tests create real resources in Azure DevOps which often cost money to run.
