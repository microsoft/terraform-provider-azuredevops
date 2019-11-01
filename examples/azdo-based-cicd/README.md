# AzureDevops Based CICD with Repo, Branch Policies, Service Connections and Group Membership

## Note: This is aspirational in nature

This template is *aspirational* and is included in order to showcase features that work today, as well as features that we hope will work in the near future.

Please keep in mind the following:
 - The commented out portions of the provider are not yet implemented. As the resources are implemented they will be uncommented in this document.
 - The schema for commented out resources will likely change once they are designed and implemented.

This sample provisions the following infrastructure:
 - AzDO project
 - AzDO Git Repo
 - AzDO Build Definition

This sample has aspirations to provision the following infrastructure:
 - AzureRM service connection, configured with permissions to a subscription in Azure
 - Docker service connection, configured with permissions to an ACR in Azure
