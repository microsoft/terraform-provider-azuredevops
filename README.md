# Terraform Provider for Azure DevOps (Devops Resource Manager)

The AzureRM Provider supports Terraform 0.12.x and later.

* [Terraform Website](https://www.terraform.io)
* [AzDO Website](https://azure.microsoft.com/en-us/services/devops/)
* [AzDO Provider Documentation](website/docs/index.html.markdown)
* [AzDO Provider Usage Examples](./examples/)

## Usage Example

```hcl
# Make sure to set the following environment variables:
#   AZDO_PERSONAL_ACCESS_TOKEN
#   AZDO_ORG_SERVICE_URL
provider "azuredevops" {
  version = ">= 0.0.1"
}

resource "azuredevops_project" "project" {
  project_name = "My Awesome Project"
  description  = "All of my awesomee things"
}

resource "azuredevops_azure_git_repository" "repository" {
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
    repo_name   = azuredevops_azure_git_repository.repository.name
    branch_name = azuredevops_azure_git_repository.repository.default_branch
    yml_path    = "azure-pipelines.yml"
  }
}
```

## Developer Requirements

* [Terraform](https://www.terraform.io/downloads.html) version 0.12.x +
* [Go](https://golang.org/doc/install) version 1.13.x (to build the provider plugin)

If you're on Windows you'll also need:
* [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)
* [Git Bash for Windows](https://git-scm.com/download/win)

For *GNU32 Make*, make sure its bin path is added to PATH environment variable.*

For *Git Bash for Windows*, at the step of "Adjusting your PATH environment", please choose "Use Git and optional Unix tools from Windows Command Prompt".*

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is **required**). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

First clone the repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-azuredevops`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers/terraform-provider-azuredevops; cd $GOPATH/src/github.com/terraform-providers/terraform-provider-azuredevops
$ git clone git@github.com:terraform-providers/terraform-provider-azuredevops.git
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-azuredevops.git
```

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

In order to run the Unit Tests for the provider, you can run:

```sh
$ make test
```

The majority of tests in the provider are Acceptance Tests - which provisions real resources in Azure. It's possible to run the entire acceptance test suite by running `make testacc` - however it's likely you'll want to run a subset, which you can do using a prefix, by running:

```sh
make testacc SERVICE='resource' TESTARGS='-run=TestAccAzureRMResourceGroup' TESTTIMEOUT='60m'
```

The following Environment Variables must be set in your shell prior to running acceptance tests:

- `AZDO_ORG_SERVICE_URL`
- `AZDO_PERSONAL_ACCESS_TOKEN`
- `AZDO_DOCKERHUB_SERVICE_CONNECTION_EMAIL`
- `AZDO_DOCKERHUB_SERVICE_CONNECTION_PASSWORD`
- `AZDO_DOCKERHUB_SERVICE_CONNECTION_USERNAME`
- `AZDO_GITHUB_SERVICE_CONNECTION_PAT`
- `AZDO_TEST_AAD_USER_EMAIL`

**Note:** Acceptance tests create real resources in Azure DevOps which often cost money to run.
