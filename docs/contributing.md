# Contributing

- [Contributing](#contributing)
  - [Looking for more?](#looking-for-more)
- [Workspace Setup & Building Project](#workspace-setup--building-project)
  - [1. Install dependencies](#1-install-dependencies)
  - [2. Clone repository](#2-clone-repository)
  - [3. Build & Install Provider](#3-build--install-provider)
  - [4. Run provider locally](#4-run-provider-locally)
- [Development SDLC](#development-sdlc)
  - [1. Pick an issue](#1-pick-an-issue)
  - [2. Repository Structure](#2-repository-structure)
  - [3. Code for Terraform](#3-code-for-terraform)
  - [4. Test changes](#4-test-changes)
  - [5. Debug changes](#5-debug-changes)
  - [6. Document changes](#6-document-changes)
- [Note about CLA](#note-about-cla)

This document is intended to be an introduction to contributing to the `terraform-provider-azuredevops` project. Links to background information about the project and general guidance on Terraform providers are included below:

If you are looking for background information on the project or related technologies (Terraform, Go and Azure DevOps), consider checking out some of these resources first:

* [Introduction to Azure DevOps](https://azure.microsoft.com/en-us/services/devops/)
* [Getting started with Terraform](https://learn.hashicorp.com/terraform#getting-started)
* [Getting started with Go](https://tour.golang.org/welcome/1)
* [README.md for project](../README.md)

If you are familiar with the technologies used for this project but are looking for general guidance on Terraform provider development, consider checking out some of these resources first:

* [Introduction to Provider Development](https://learn.hashicorp.com/terraform/development/writing-custom-terraform-providers)
* [Terraform provider discovery documentation](https://www.terraform.io/docs/extend/how-terraform-works.html#discovery)
* [Terraform Acceptance Testing](https://www.terraform.io/docs/extend/best-practices/testing.html#built-in-patterns)
* [Terraform Schema Behaviors](https://www.terraform.io/docs/extend/schemas/schema-behaviors.html)

If you are still reading, then you are in the right place!

## Looking for more?

If, after reading through the content here, you are seeking more detailed information, you may want to checkout some of the following resources:

* [Getting Started Guide](https://github.com/Azure/terraform/blob/master/provider/CONTRIBUTE.md) written for the `terraform-provider-azurerm` provider. While it targets a different provider there are some great findings that you can read about.
* [Development Environment for Go Lang on Mac](https://medium.com/@tsuyoshiushio/development-environment-for-go-lang-ede316d4512a)

# Workspace Setup & Building Project

This section describes how to get your developer workspace running for the first time so that you're ready to start making contributions. If you have already done this, check out [Development SDLC](#development-sdlc).

> These steps assume you are running with `bash`. If you are using Windows, run all commands using WSL. They are not tested using GitBash.

## 1. Install dependencies

The recommended development environment is Linux or Mac. If you're on Windows you should [install WSL](https://docs.microsoft.com/en-us/windows/wsl/install-win10) so that your environment more closely mirrors a Linux environment.

You will need the following dependencies installed in order to get started:

* [Terraform](https://www.terraform.io/downloads.html) version 0.11.x +
* [Go](https://golang.org/doc/install) version 1.12.x +
* An editor of your choice. We recommend [Visual Studio Code](https://code.visualstudio.com/Download) but any editor will do.


## 2. Clone repository

> Note: This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [GOPATH](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your home directory outside of the standard `GOPATH`.

**Setup your workspace**
```bash
$ DEV_ROOT="$HOME/workspace"
$ mkdir -p "$DEV_ROOT"
$ cd "$DEV_ROOT"
```

**Get the code**

```bash
$ git clone https://github.com/microsoft/terraform-provider-azuredevops.git
$ cd terraform-provider-azuredevops/
```

## 3. Build & Install Provider

**Build & test the azure devops provider**

Running the next command will orchestrate a few things for you:

- Verify that all required Go packages are installed. This may take a few minutes the first time you run the script as the packages will be cached locally. Subsequent runs will be much faster
- Run all unit tests
- Compile the provider codebase

```bash
$ ./scripts/build.sh
...
[INFO] Executing unit tests
...
[INFO] Build finished successfully
```

After this script runs you should see a `./bin/` directory with the compiled terraform provider.

```bash
$ ls -lah ./bin/
...
-rwxrwxrwx 1 ... terraform-provider-azuredevops_v0.0.1
```

**Install the provider**

Terraform provider plugins are not intended to be run directly. You can see this for yourself:

```bash
$ ./bin/terraform-provider-azuredevops_v0.0.1

This binary is a plugin. These are not meant to be executed directly.
Please execute the program that consumes these plugins, which will
load any plugins automatically
```

In order to use the provider locally it must be installed into a location discoverable by Terraform. The `local-install.sh` script does this for you.

```bash
$ ./scripts/local-install.sh
[INFO] Installing provider to /home/$USER/.terraform.d/plugins/
```

To learn more about the plugin discovery process, refer to the official [Terraform provider discovery documentation](https://www.terraform.io/docs/extend/how-terraform-works.html#discovery).

## 4. Run provider locally

> Note: These steps assume you have built the provider locally using the previous steps. The samples in the `./examples/` folder use syntax specific to Terraform 12 + and are not compatible with older versions of Terraform.

**Configure the provider**

You can now use the provider just like you normally would. Try it out by using the project examples:

```bash
$ cd examples/github-based-cicd-simple/

# AZDO_ORG_SERVICE_URL will be the URL of the AzDO org that you want to provison
# resources inside of.
#   ex: https://dev.azure.com/<your org name>
$ export AZDO_ORG_SERVICE_URL="..."

# AZDO_PERSONAL_ACCESS_TOKEN will be the personal access token that grants access
# to provision and manage resources in Azure DevOps.
#   documentation: https://docs.microsoft.com/en-us/azure/devops/organizations/accounts/use-personal-access-tokens-to-authenticate?view=azure-devops
$ export AZDO_PERSONAL_ACCESS_TOKEN="..."

# Note: AZDO_GITHUB_SERVICE_CONNECTION_PAT is not specifically required
# by the provider, but it is required by the example in this folder.
#   documentation: https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line
$ export AZDO_GITHUB_SERVICE_CONNECTION_PAT="..."

$ terraform init
...
Terraform has been successfully initialized!
...
```

The provider has now been initialized and can be used like any other provider in Terraform. You can try for yourself by running any of the `terraform plan|apply|destroy|...` commands.


# Development SDLC

This section outlines typical development scenarios for the repository, relevant areas of the codebase and provides guidance for writing tests and documentation for your changes.

## 1. Pick an issue

Please find an open issue in the backlog that has no asignee and request it to be assigned to you by a core contributor if you don't have the permissions to do it yourself. This will help us track which work items are being worked on.

## 2. Repository Structure

The repository generally follows the [patterns suggested by HashiCorp](https://learn.hashicorp.com/terraform/development/writing-custom-terraform-providers). This pattern is implemented by all major terraform providers.

Here is a simplified version of the repository:
```bash
$ tree -L 2
.
├── azdosdkmocks --------> Generated mocks for AzDO Go SDK
├── azuredevops ---------> Provider implementation
│   ├── config.go -------> AzDO SDK initialization lives here
│   ├── provider.go -----> Exports the AzDO terraform provider
│   ├── data_*.go -------> data_*.go files contain terraform data sources implementations
│   ├── resource_*.go ---> resource_*.go files contain terraform resource implementations
│   └── utils -----------> Utilities used across the codebase
├── docs ----------------> Developer documentation
├── go.mod --------------> Describes project dependencies
├── scripts -------------> All scripts live here
└── website -------------> Client facing documentation
```

## 3. Code for Terraform

After you picked an issue and figured out where you will impelment it, you will quickly realize that HashCorp has take a very opinionated appraoch to building Terraform providers. The following section will outline a few common scenarios for terraform plugin development, and point you towards different pieces of the code that you may need to work with.

If you have not already gone through the guide published by [HashiCorp](https://learn.hashicorp.com/terraform/development/writing-custom-terraform-providers), now is a good time to do so.


**Scenario 1: Change resource or data source schema**

If you need to add, remove or modify the schema of a data source or resource, you will need to first identify the relevant file. The naming scheme used is as follows:
  - `data_foo.go` - an implementation for the Terraform data source for the **foo** Azure DevOps resource
  - `resource_foo.go` - an implementation for the Terraform resource for the **foo** Azure DevOps resource

Open the file and look for the schema. Here is a simple schema found in `data_group.go`. This is fairly simple and only defines three attributes. More complicated ones can be found in the [build definition code](../azuredevops/resource_build_definition.go). The official documentation for the schema can be [found here](https://godoc.org/github.com/bradfeehan/terraform/helper/schema).

![Group Data Source Schema](https://user-images.githubusercontent.com/2497673/67519578-b2500400-f66c-11e9-89f2-725a4341a317.png)

**Scenario 2: Modify an existing resource or data source**

If you need to modify the business logic in an existing resource or data source, you will find the relevant code in one of the `Create`, `Read`, `Update` or `Delete` functions.

![Resource CRUD Functions](https://user-images.githubusercontent.com/2497673/67520080-c5170880-f66d-11e9-81fd-90eccc85eeae.png)

The prototype of these functions are all quite similar. Here is an example of a create function. Keep note of the following details:

 - `d *schema.ResourceData` is passed to the provider by Terraform. It contains the resource configuration specified by the client using the provider, along with any data pulled from the Terraform state.
 - `m interface{}` is, in the case of this provider, a structure containing all of the (intialized) clients needed to make API calls to Azure DevOps.
 - [Flatten/Expand](https://learn.hashicorp.com/terraform/development/writing-custom-terraform-providers#implementing-a-more-complex-read) is a common "idiom" used across terraform providers. It is a standard approach to marshaling and unmarshaling API data structures into the internal terraform state.

![image](https://user-images.githubusercontent.com/2497673/67520284-217a2800-f66e-11e9-87c8-2f87e882eaca.png)


**Scenario 3: Implement a new resource or data source**

If you need to implement a new resource or data source, you should first review a few implementations in order to understand the patterns in the codebase, which generally match the way that other Terraform providers are built.

This scenario is a mix of *Scenario 1* and *Scenario 2.* However, after implementing the schema and CRUD operations, you will need to register your newly created resource or data source. The code for doing this is located in `provider.go`.

![Provider Registration](https://user-images.githubusercontent.com/2497673/67520904-60f54400-f66f-11e9-93ee-43535c72e0da.png)


## 4. Test changes

**Running Unit Tests**

The unit tests are executed whenever `./scripts/build.sh` is run. This can be run locally, but will also be run on every automated build and will be a gate for any PR against this repository. The tests can also be run in isolation by running the following:
```bash
$ ./scripts/unittest.sh
```

**Running Acceptance Tests (Integration Tests)**

> Note: Running acceptance tests provisions and deletes actual resources in AzDO. This can cost money and can be dangerous if you are not running them in isolation!

The acceptance tests for terraform providers are typically implemented as [Acceptance Tests](https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html). Acceptance tests can be invoked by running the following:
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

**Writing your own tests**

There is a lot of context to cover here, so check out our dedicated [testing document](./testing.md) for more information on writing unit and acceptance tests.

## 5. Debug changes

There is a lot of context to cover here, so check out our dedicated [debugging document](debugging.md) for more information on debugging the provider.

## 6. Document changes

Most changes should involve updates to the client-facing reference documentation. Please update the existing documentation to reflect any changes you have made to the codebase.

| Name | Description | Link |
| ---- | ----------- | ---- |
| Index | Table of contents | [index.md](../website/index.md) |
| Resources | Resources reference | [resources](../website/docs/r) |
| Data Sources | Data Sources reference | [data sources](../website/docs/d) |
| Guides | Guide and tutorial docs | [guides](../website/docs/guildes) |

# Note about CLA

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

