# Roadmap

This document describes the project roadmap so that people in the OSS community can get a glimpse into the short/medium term goals. Please keep in mind that this document is aspirational in nature and the scope and timeline may change based on a number of factors. If you have questions about the content here, feel free to open an issue.

## V1 Milestone - ETA late Q4

The V1 milestone has a strong focus on being able to provision a baseline set of features all surrounding CI/CD. There is an aspirational [Terraform Template](../examples/azdo-based-cicd/main.tf) in the samples folder that exemplifies the feature-set being targeted. The features are as follows:
 - Provision & manage AzDO [Projects](https://docs.microsoft.com/en-us/rest/api/azure/devops/core/projects?view=azure-devops-rest-5.0)
 - Provision & manage AzDO [Build Definitions](https://docs.microsoft.com/en-us/rest/api/azure/devops/build/definitions?view=azure-devops-rest-5.0) in projects (pipelines defined in YML)
 - Provision & manage AzDO [Service Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-5.0)
 - Provision & manage AzDO [Git Repos](https://docs.microsoft.com/en-us/rest/api/azure/devops/git/repositories?view=azure-devops-rest-5.0)
 - Provision & manage AzDO [Git Repos Branch Policies](https://docs.microsoft.com/en-us/rest/api/azure/devops/policy/configurations/create?view=azure-devops-rest-5.0)
 - Provision & manage AzDO [Variable Groups](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/variablegroups?view=azure-devops-rest-5.0)
 - Manage AzDO [Group Membership](https://docs.microsoft.com/en-us/rest/api/azure/devops/graph/memberships?view=azure-devops-rest-5.0). This includes adding user entitlements and adding users to groups
