# Simple GitHub based CICD

This sample provisions the following infrastructure:

- AzDO project
- AzDO build definition that points to `.azdo/azure-pipeline-nightly.yml`
- GitHub service connection that grants AzDO the ability to interact with GitHub via a PAT

Builds will be triggered based on the trigger block in the yaml pipeline specified.
