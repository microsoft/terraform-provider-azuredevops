# Contributing

This document is intended to be an introduction to contributing to the `terraform-provider-azuredevops` project.

If you are looking for background information on the project or related technologies (Terraform, Go and Azure DevOps), consider checking out some of these resources first:

* [Introduction to Azure DevOps](https://azure.microsoft.com/en-us/services/devops/)
* [Getting started with Terraform](https://learn.hashicorp.com/terraform#getting-started)
* [Getting started with Go](https://tour.golang.org/welcome/1)
* [README.md for project](../README.md)

If you are familiar with the technologies used for this project but are looking for general guidance on Terraform provider development, consider checking out some of these resources first:

* [Introduction to Provider Development](https://learn.hashicorp.com/terraform/development/writing-custom-terraform-providers)
* [Terraform provider discovery documentation](https://www.terraform.io/docs/extend/how-terraform-works.html#discovery)
* [Terraform Acceptance Testing](https://www.terraform.io/docs/extend/best-practices/testing.html#built-in-patterns)

If you are still reading, then you are in the right place!

If, after reading through the content here, you are seeking more detailed information, you may want to checkout [this awesome getting started guide](https://github.com/Azure/terraform/blob/master/provider/CONTRIBUTE.md) that was written for the `terraform-provider-azurerm` project. While it targets a different provider there are some great findings that you can read about

## Install the dependencies

The recommended development environment is Linux or Mac. If you're on Windows you should [install WSL](https://docs.microsoft.com/en-us/windows/wsl/install-win10) so that your environment more closely mirrors a Linux environment.

You will need the following dependencies installed in order to get started:

* [Terraform](https://www.terraform.io/downloads.html) version 0.11.x +
* [Go](https://golang.org/doc/install) version 1.12.x +
* An editor of your choice. We recommend [Visual Studio Code](https://code.visualstudio.com/Download) but any editor will do.


## Building the provider locally

**Note** This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [GOPATH](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your home directory outside of the standard GOPATH.

**Note** These steps assume you are running with `bash`. If you are using Windows, run all commands using WSL. They are not tested using GitBash.

#### Setup your workspace
```bash
$ DEV_ROOT="$HOME/workspace"
$ mkdir -p "$DEV_ROOT"
$ cd "$DEV_ROOT"
```

#### Get the code

```bash
$ git clone https://github.com/microsoft/terraform-provider-azuredevops.git
$ cd terraform-provider-azuredevops/
```

#### Build & test the provider

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

#### Install the provider

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

## Running the provider locally

**Note** These steps assume you have built the provider locally using the previous steps.

#### Configuring the provider

You can then use the provider just like you normally would. You can try it out by using the project examples:

```bash
$ cd examples/

$ export AZDO_ORG_SERVICE_URL="..."
$ export AZDO_PERSONAL_ACCESS_TOKEN="..."

# Note: this one is not specifically required by the provider,
# but it is required by the example in this folder...
$ export AZDO_GITHUB_SERVICE_CONNECTION_PAT="..."

$ terraform init
...
Terraform has been successfully initialized!
...
```

After the provider has been correctly initialized, it can be used like any other provider. You can try for yourself by running any of the `terraform plan|apply|destroy|...` commands.


## Running the provider tests

#### Unit Tests

The unit tests are executed whenever `./scripts/build.sh` is run. This can be run locally, but will also be run on every automated build and will be a gate for any PR against this repository. If you made it this far, you've already ran these!

The tests can also be run in isolation by executing the following command from the [module root](https://blog.golang.org/using-go-modules). This will be the folder containing `go.mod`:
```bash
go test ./...
```

#### Acceptance Tests (Integration Tests)

As of now these are not implemented. Please add a 👍 reaction to our [open issue](https://github.com/microsoft/terraform-provider-azuredevops/issues/59) in order to draw attention to this.

## Note about CLA

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

