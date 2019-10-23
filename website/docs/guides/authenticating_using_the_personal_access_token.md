# Azure DevOps Provider: Authenticating using the Personal Access Token

Azure DevOps provider support personal access token for authenticating to Azure DevOps.

## Create a personal access token

1. Go to your Azure DevOps. Select an organization.
2. Click the icon next to your icon at the right top corner.
3. Select "Personal access tokens".
4. Click "New Token" then create a new personal access token with the access required by your template. This will be driven primarily based on which resources you need to provision in Azure DevOps. A token with Full access scope will work but may provide more access than you need.

5. Copy the personal access token.

## Configure Environment Variables

Set the two environment variables. For more details, see the [Readme](../../../README.md).
`AZDO_PERSONAL_ACCESS_TOKEN` and `AZDO_ORG_SERVICE_URL`. If you use bash, you can try this.

```bash
$ export AZDO_PERSONAL_ACCESS_TOKEN=<Personal Access Token>
$ export AZDO_ORG_SERVICE_URL=https://dev.azure.com/<Your Org Name>
```

## Configuration

Configuration file requires `azuredevops` provider section. Then use any resources and data sources you want.

```hcl
provider "azuredevops" {
  version = ">= 0.0.1"
}

resource "azuredevops_project" "project" {
  project_name       = "Test Project"
  description        = "Test Project Description"
}
```
