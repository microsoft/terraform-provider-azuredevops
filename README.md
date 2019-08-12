- [Background Context](#background-context)
  - [Azure DevOps](#azure-devops)
  - [Go](#go)
  - [Terraform](#terraform)
  - [Docker](#docker)
- [Getting Started](#getting-started)
  - [Dependencies](#dependencies)
  - [Running the samples](#running-the-samples)
    - [Azure DevOps API Sample (`./azdo-api-samples`)](#azure-devops-api-sample-azdo-api-samples)
    - [Terraform Provider Implementation Sample (`./terraform-provider`)](#terraform-provider-implementation-sample-terraform-provider)

# Background Context

In order to be effective on day 1 of this hack you will need to be familiar with Azure DevOps, Go and Terraform. In order to accelerate any gaps in your knowledge, I'm providing some helpful background guides/documentation.

Please spend some time brushing up on any areas that you are not already comfortable in.

## Azure DevOps

- **20 min**: [Create your first pipeline in AzDO](https://docs.microsoft.com/en-us/azure/devops/pipelines/create-first-pipeline?view=azure-devops&tabs=tfs-2018-2)
- **10 min**: [AzDO YAML pipeline schema](https://docs.microsoft.com/en-us/azure/devops/pipelines/yaml-schema?view=azure-devops&tabs=schema) - Note: no need to dive deep here.

## Go

- **30-60 min**: [A Tour of Go](https://tour.golang.org/welcome/1)

## Terraform

- **30-60 min**: [Getting started with Terraform & Azure](https://learn.hashicorp.com/terraform?track=azure)
- **20 min**: [Writing custom providers in Terraform](https://learn.hashicorp.com/terraform/development/writing-custom-terraform-providers)


## Docker

- **10 min**: [Containers Intro](https://www.docker.com/resources/what-container)
- **20-30 min**: [Building Docker Images](https://docs.docker.com/get-started/part2/)

# Getting Started

## Dependencies

Most of the project builds on Docker in order to reduce the chance of workspace incompatability. Therefore, all you need to install are the following:

- WSL and/or `bash` shell. [WSL installation instructions](https://docs.microsoft.com/en-us/windows/wsl/install-win10)
- Docker. [Docker installation instructions](https://runnable.com/docker/getting-started/)


## Running the samples

### Azure DevOps API Sample (`./azdo-api-samples`)

This sample contains an API call sequence to orchestrate the following actions:

- Create an Azure DevOps project, if it does not already exist
- Create an Azure DevOps service connection, if it does not already exist
- Create an Azure DevOps build/release pipeline, if it does not already exist

It is intended to show how the Azure DevOps APIs work at a high level. Further documentation on the Azure DevOps APIs can be [found here](https://docs.microsoft.com/en-us/rest/api/azure/devops/?view=azure-devops-rest-5.1).

The sample can be run by following these steps:

- Create `.env` file

```bash
cd azdo-api-samples/
cp .env.template .env
```
- Fill out `.env` file with relevant environment variables. You will need the following environment variables defined:

| Name | Description | Example Value |
| --- | --- | --- |
 AZDO_PAT | Personal Access Token for hitting API endpoints in Azure DevOps organization | `***` (*sensitive information*) |
 AZDO_ORGANIZATION | Organization in which to create resources in Azure DevOps | `Awsome Organization` |
 AZDO_NEW_PROJECT_NAME | Name of project that will be created | `Awesome Project` |
 AZDO_NEW_PROJECT_DESCRIPTION | Description of project that will be created | `Project for creating awesome things` |
 AZDO_NEW_PIPELINE_NAME | Name of pipeline that will be created | `CI/CD Pipeline` |
 AZDO_PIPELINE_YML_GIT_REPO | GitHub repository hosting pipeline Yaml definition | `nmiodice/terraform-azure-devops-hack` |
 AZDO_PIPELINE_YML_GIT_REPO_BRANCH | Default branch in GitHub for pipeline | `master` |
 AZDO_PIPELINE_YML_FILENAME | Name of Yaml file in GitHub repository | `azdo-api-samples/azure-pipeline.yml` |
 AZDO_GITHUB_SERVICE_CONNECTION_NAME | Name of service connection that will be created | `GitHub Service Connection` |
 AZDO_GITHUB_SERVICE_CONNECTION_PAT | Personal Access Token for authenticating to GitHub | `***` (*sensitive information*) |

- Build/Run the sample
```bash
bash build.sh
docker run -it azdoapis:latest
```

- Navigate to the newly provisioned pipeline in the newly provisioned Azure DevOps project and run the pipeline.

### Terraform Provider Implementation Sample (`./terraform-provider`)

This sample contains a self-contained bare-bones terraform provider implementation and sample usage from Terraform HCL. 

The sample is broken down into the following components:

- `provider-src`: The provider implementation. The provider is named `azuredevops` and exports a single resource `foo`. This sample follows the recommended naming conventions outlined in the documentation for [writing custom providers in Terraform](https://learn.hashicorp.com/terraform/development/writing-custom-terraform-providers).
- `terraform-src`: Terraform `HCL` that uses the sample provider.
- `Dockerfile`: Builds & installs `azuredevops` provider. Installs sample code. Running the resulting image dumps you into a terraform sandbox environment.
- `build.sh`: This builds the docker file
  
The sample can be run like so:

```bash
# These run on your local WSL/bash instance
cd terraform-provider/

# build source, including running unit tests.
bash build.sh
docker run -it azdotf:latest

# These run within the docker container
terraform init                  # initialize workspace
terraform plan                  # show plan 
terraform apply -auto-approve   # apply plan
terraform destroy -auto-approve # destroy resources (these are dummy resources!)
```
