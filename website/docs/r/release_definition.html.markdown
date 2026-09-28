---
layout: "azuredevops"
page_title: "AzureDevOps: azuredevops_release_definition"
description: |-
  Manages a Release Definition within Azure DevOps organization.
---

# azuredevops_release_definition

Manages a classic Release Definition within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_git_repository" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Repository"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_build_definition" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Build"

  repository {
    repo_type   = "TfsGit"
    repo_id     = azuredevops_git_repository.example.id
    branch_name = azuredevops_git_repository.example.default_branch
    yml_path    = "azure-pipelines.yml"
  }
}

data "azuredevops_agent_queue" "example" {
  project_id = azuredevops_project.example.id
  name       = "Azure Pipelines"
}

resource "azuredevops_release_definition" "example" {
  project_id          = azuredevops_project.example.id
  name                = "Example Release"
  path                = "\\"
  release_name_format = "Release-$(rev:r)"

  variable {
    name  = "environment"
    value = "production"
  }

  secret_variable {
    name  = "api_key"
    value = var.api_key
  }

  artifact {
    alias      = "myapp"
    type       = "Build"
    is_primary = true

    definition_reference {
      key  = "project"
      id   = azuredevops_project.example.id
      name = azuredevops_project.example.name
    }

    definition_reference {
      key  = "definition"
      id   = azuredevops_build_definition.example.id
      name = azuredevops_build_definition.example.name
    }

    definition_reference {
      key  = "defaultVersionType"
      id   = "latestType"
      name = "Latest"
    }
  }

  # Create a release automatically whenever a new build is produced.
  artifact_source_trigger {
    artifact_alias = "myapp"
  }

  # Also create a release every weekday at 03:00 UTC.
  schedule_trigger {
    days              = ["monday", "tuesday", "wednesday", "thursday", "friday"]
    start_hours       = 3
    start_minutes     = 0
    time_zone_id      = "UTC"
    only_with_changes = true
  }

  stage {
    name = "Dev"

    deploy_phase {
      name           = "Agent job"
      agent_queue_id = data.azuredevops_agent_queue.example.id

      task {
        # Command Line task
        task_id = "d9bafed4-0b18-4f58-968d-86655b4d2ce9"
        version = "2.*"
        inputs = {
          script = "echo Hello world"
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID or name of the project. Changing this forces a new resource to be created.
* `name` - (Required) The name of the release definition.
* `path` - (Optional) The folder path of the release definition. Defaults to `\`.
* `description` - (Optional) A description of the release definition.
* `release_name_format` - (Optional) The release name format. Defaults to `Release-$(rev:r)`.
* `variable_groups` - (Optional) A set of variable group IDs to link to the release definition.
* `variable` - (Optional) One or more `variable` blocks as documented below.
* `secret_variable` - (Optional) One or more `secret_variable` blocks as documented below.
* `artifact` - (Optional) One or more `artifact` blocks as documented below.
* `artifact_source_trigger` - (Optional) One or more `artifact_source_trigger` blocks as documented below. Creates a release automatically when a new version of the linked artifact is available (continuous deployment).
* `schedule_trigger` - (Optional) One or more `schedule_trigger` blocks as documented below. Creates a release on a recurring schedule.
* `source_repo_trigger` - (Optional) One or more `source_repo_trigger` blocks as documented below. Creates a release when a commit is pushed to a linked source repository artifact.
* `container_image_trigger` - (Optional) One or more `container_image_trigger` blocks as documented below. Creates a release when a new container image is pushed to a linked container repository artifact.
* `package_trigger` - (Optional) One or more `package_trigger` blocks as documented below. Creates a release when a new version of a linked package artifact is published.
* `pull_request_trigger` - (Optional) One or more `pull_request_trigger` blocks as documented below. Creates a release when a pull request is raised against a linked source repository artifact.

~> **NOTE:** Each trigger references an `artifact` by its `alias` via `artifact_alias`, and the referenced artifact must be of a compatible `type`. Trigger types not yet exposed by this resource are preserved on update rather than removed.

* `stage` - (Required) One or more `stage` blocks as documented below.

A `variable` block supports the following:

* `name` - (Required) The name of the variable.
* `value` - (Optional) The value of the variable.
* `allow_override` - (Optional) Whether the variable can be overridden at release time. Defaults to `false`.

A `secret_variable` block supports the following:

* `name` - (Required) The name of the secret variable.
* `value` - (Optional) The value of the secret variable. This value is marked as sensitive.
* `allow_override` - (Optional) Whether the variable can be overridden at release time. Defaults to `false`.

~> **NOTE:** A name may not be used by both a `variable` and a `secret_variable` block in the same scope.

~> **NOTE:** The service never returns the value of a secret variable, so it is tracked only in Terraform state. On import, `secret_variable` values are empty and the next plan will show a diff until the values are supplied in the configuration.

An `artifact` block supports the following:

* `alias` - (Required) The alias of the artifact. This is referenced by release variables and triggers (e.g. `$(myapp.BuildNumber)`).
* `type` - (Required) The artifact type. Common values are `Build`, `Git`, `GitHub`, `GitHubRelease`, `TFVC`, `Jenkins`, `PackageManagement` (Azure Artifacts feeds), and container registry types such as `DockerHub` and `AzureContainerRepository`. Azure DevOps extensions can register additional artifact types, so this value is not restricted to a fixed list. An unrecognised type is rejected by the service with `VS402864: No artifact type found corresponding to ID <type>`.
* `is_primary` - (Optional) Whether this is the primary artifact. Defaults to `true`.
* `is_retained` - (Optional) Whether the artifact is retained by the release. Defaults to `false`.
* `definition_reference` - (Required) One or more `definition_reference` blocks describing the artifact source, as documented below.

A `definition_reference` block supports the following:

* `key` - (Required) The reference key, for example `project` or `definition`. The keys an artifact requires are determined by its `type`, as described below.
* `id` - (Required) The value for the key. This is the value the service resolves the reference against, usually a GUID or another machine-readable identifier.
* `name` - (Optional) The display name for the value, as shown in the Azure DevOps web UI. This is informational only and is not used to resolve the reference.

Azure DevOps models all artifact types through a single generic structure, so the fields that identify a particular source are supplied as this untyped set of key/value pairs rather than as dedicated attributes. Each `type` requires a different set of keys, which together narrow the artifact down to one specific source and a default version to deploy.

A `Build` artifact, for example, is described by three keys because three separate questions must be answered:

* `project` - which Team Project contains the build pipeline.
* `definition` - which build pipeline within that project produces the artifact.
* `defaultVersionType` - which build to deploy when a release is created manually, for example `latestType` for the most recent build.

Of these, only `project` and `definition` are enforced by the service; `defaultVersionType` is accepted by every artifact type and is defaulted when omitted, but the Azure DevOps web UI always writes it, so most definitions set it explicitly. The keys the service enforces for the most common artifact types are:

* `Build` - `project` and `definition` (the build definition ID).
* `Git` - `project`, `definition` (the repository ID), and `branches` (e.g. `refs/heads/main`). Note that `defaultVersionBranch` is rejected for `Git` artifacts; use `branches` instead.
* `TFVC` - `project` and `definition`.
* `GitHub` - `connection` (the GitHub service connection ID), `definition` (the repository), and `branch`.
* `GitHubRelease` - `connection` (the GitHub service connection ID) and `definition` (the repository).
* `Jenkins` - `connection` (the Jenkins service connection ID) and `definition` (the job name).
* `PackageManagement` - `feed` (the feed ID) and `definition` (the package name).
* `DockerHub` - `connection` (the Docker Hub service connection ID), `definition` (the repository), and `namespaces`.
* `AzureContainerRepository` - `connection` (the Azure Resource Manager service connection ID), `definition` (the repository), `registryurl`, and `resourcegroup`.

Because artifact types can be added by Azure DevOps extensions, this list cannot be exhaustive. To determine the keys for a type not listed above, either:

* Add the artifact to a release definition through the Azure DevOps web UI, then read the `artifacts[].definitionReference` object back from the [Definitions - Get](https://learn.microsoft.com/en-us/rest/api/azure/devops/release/definitions/get) REST API. Its keys map one-to-one onto `definition_reference` blocks, with `id` and `name` taken from each entry.
* Apply the artifact with the keys you expect and let the API report what is missing. It names each problem individually, for example `Input field 'branches' must be present for artifact source: '<alias>'` for a missing key, or `'defaultVersionBranch' is not a valid input field` for one that does not apply to the type.

~> **NOTE:** The Azure DevOps server may enrich `definition_reference` with additional keys it computes (for example `defaultVersionBranch` or `artifactSourceDefinitionUrl`). Keys that are not listed in the configuration are ignored when the resource is read back, so they do not cause a diff. This filtering relies on the configuration, so on import any computed keys are recorded in state as well; if a plan after import shows a diff on `definition_reference`, add the extra keys to the configuration or remove them from state.

An `artifact_source_trigger` block supports the following:

* `artifact_alias` - (Required) The `alias` of the `artifact` this trigger watches. A release is created when a new version of that artifact is available.
* `trigger_condition` - (Optional) One or more `trigger_condition` blocks that filter which artifact versions trigger a release. When omitted, any new version triggers a release.

A `trigger_condition` block supports the following:

* `source_branch` - (Optional) Only trigger for artifacts produced from this source branch (e.g. `refs/heads/main`).
* `use_build_definition_branch` - (Optional) Use the build definition's default branch as the branch filter. Defaults to `false`.
* `create_release_on_build_tagging` - (Optional) Create a release when the build is tagged. Defaults to `false`.
* `tags` - (Optional) A set of build tags to filter on.

A `schedule_trigger` block supports the following:

* `days` - (Required) A set of days on which to create a release. Valid values: `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `sunday`, `all` (case-insensitive; stored in lowercase). Note: specifying all seven days is normalized by the server to `all`.
* `start_hours` - (Optional) The local-time hour (0-23) at which to start. Defaults to `3`.
* `start_minutes` - (Optional) The local-time minute (0-59) at which to start. Defaults to `0`.
* `time_zone_id` - (Optional) The time zone ID for the schedule (e.g. `UTC`). Defaults to `UTC`.
* `only_with_changes` - (Optional) Only create a release if the artifact or definition has changed since the last release. Defaults to `false`.

A `source_repo_trigger` block supports the following:

* `artifact_alias` - (Required) The `alias` of the source repository `artifact` this trigger watches (an artifact of type `Git`, `GitHub`, or `TFVC`).
* `branch_filter` - (Optional) A `branch_filter` block as documented below, restricting which branches trigger a release. When omitted, pushes to any branch trigger a release.

A `branch_filter` block supports the following:

* `include` - (Optional) A set of branches that trigger a release when pushed to (e.g. `refs/heads/main`).
* `exclude` - (Optional) A set of branches that do not trigger a release.

~> **NOTE:** At most one `branch_filter` block may be specified, and it must set at least one of `include` or `exclude`. The service stores branch filters as a single flat list, so multiple blocks cannot be represented and an empty block would produce a persistent diff.

A `container_image_trigger` block supports the following:

* `artifact_alias` - (Required) The `alias` of the container repository `artifact` this trigger watches.
* `tag_filter` - (Optional) A single image tag pattern that triggers a release. When omitted, any new image triggers a release. The service permits at most one tag filter per container image trigger.

~> **NOTE:** Azure DevOps validates the registry credentials when a `container_image_trigger` is attached, so the linked service connection must hold credentials that can authenticate against the registry. Otherwise the API rejects the definition with an error such as `Unable to get access token from Docker Hub repository.`

A `package_trigger` block supports the following:

* `artifact_alias` - (Required) The `alias` of the package `artifact` this trigger watches.

A `pull_request_trigger` block supports the following:

* `artifact_alias` - (Required) The `alias` of the source repository `artifact` this trigger watches.
* `status_policy_name` - (Optional) The name of the policy used to publish the release status back to the pull request.
* `use_artifact_reference` - (Optional) Take the code repository details from the linked artifact rather than specifying them explicitly. Defaults to `true`. Set to `false` when supplying `code_repository_reference`.
* `code_repository_reference` - (Optional) A `code_repository_reference` block as documented below. Only required when `use_artifact_reference` is `false`.
* `trigger_condition` - (Optional) One or more `trigger_condition` blocks that filter which pull requests trigger a release, as documented below. When omitted, pull requests targeting any branch trigger a release.

A `code_repository_reference` block supports the following:

* `system_type` - (Required) The source system hosting the repository. Valid values are `tfsGit` and `gitHub`.
* `repository_reference` - (Optional) One or more `repository_reference` blocks, each supplying a `key`, a `value`, and an optional `display_value`. The required keys depend on `system_type`; for `tfsGit` these are typically `project` and `repository`.

~> **NOTE:** As with `definition_reference`, the server may enrich `repository_reference` with additional computed keys. Only the keys present in your configuration are tracked, so unlisted keys will not produce a diff.

A `trigger_condition` block within `pull_request_trigger` supports the following:

* `target_branch` - (Optional) Only trigger for pull requests targeting this branch (e.g. `refs/heads/main`).
* `tags` - (Optional) A set of tags to filter on.

A `stage` block supports the following:

* `name` - (Required) The name of the stage (environment).
* `rank` - (Optional) The order of the stage. Defaults to its position in the list.
* `condition` - (Optional) One or more `condition` blocks controlling when the stage runs, as documented below. When omitted, the **first** stage defaults to starting on release creation (`ReleaseStarted`). At least one stage must start on release creation whenever a trigger is configured.
* `variable` - (Optional) One or more stage-scoped `variable` blocks (same schema as above).
* `secret_variable` - (Optional) One or more stage-scoped `secret_variable` blocks (same schema as above).
* `variable_groups` - (Optional) A set of variable group IDs scoped to the stage.
* `retention_policy` - (Optional) A `retention_policy` block as documented below.
* `pre_deploy_approval` - (Optional) A `pre_deploy_approval` block as documented below.
* `post_deploy_approval` - (Optional) A `post_deploy_approval` block as documented below.
* `pre_deployment_gates` - (Optional) A `pre_deployment_gates` block as documented below.
* `post_deployment_gates` - (Optional) A `post_deployment_gates` block as documented below.
* `environment_options` - (Optional) An `environment_options` block as documented below.
* `execution_policy` - (Optional) An `execution_policy` block as documented below.
* `environment_trigger` - (Optional) One or more `environment_trigger` blocks as documented below.
* `deploy_phase` - (Required) One or more `deploy_phase` blocks as documented below.

A `condition` block supports the following:

* `condition_type` - (Optional) The type of condition. One of `event`, `environmentState`, or `artifact`. Defaults to `event`.
* `name` - (Required) The condition name. For an `event` condition that starts the stage on release creation, use `ReleaseStarted`. For an `environmentState` condition, use the name of the preceding stage.
* `value` - (Optional) The condition value (for example, the required state for an `environmentState` condition).

~> **Note** Stages run **in parallel** or **sequentially** based solely on their `condition` blocks. Give two or more stages an `event` / `ReleaseStarted` condition to run them in parallel (each starts on release creation). Use an `environmentState` condition that names a preceding stage to run a stage sequentially after it. The `rank` attribute only controls layout, not execution order.

For example, to run `Dev` and `QA` in parallel and then `Prod` after `Dev` succeeds:

```hcl
stage {
  name = "Dev"
  condition {
    condition_type = "event"
    name           = "ReleaseStarted"
  }
  # ...deploy_phase...
}

stage {
  name = "QA"
  condition {
    condition_type = "event"
    name           = "ReleaseStarted"
  }
  # ...deploy_phase...
}

stage {
  name = "Prod"
  condition {
    condition_type = "environmentState"
    name           = "Dev"
    value          = "4" # 4 = succeeded
  }
  # ...deploy_phase...
}
```

A `retention_policy` block supports the following:

* `days_to_keep` - (Optional) Number of days to keep a release. Defaults to `30`.
* `releases_to_keep` - (Optional) Number of releases to keep. Defaults to `3`.
* `retain_build` - (Optional) Whether to retain the associated build. Defaults to `true`.

A `pre_deploy_approval` / `post_deploy_approval` block supports the following:

* `approval` - (Optional) One or more `approval` blocks. When omitted, a single automated approval is configured.
* `approval_options` - (Optional) An `approval_options` block as documented below.

An `approval` block supports the following:

* `approver_id` - (Optional) The ID of the identity that must approve. Omit for an automated approval. This must be the identity's storage-key UUID. When referencing a group from the `azuredevops_group` data source, use its `group_id` attribute (its `id` is a descriptor and `origin_id` is the Azure AD object ID, neither of which is accepted here).
* `is_automated` - (Optional) Whether the approval is automated. Defaults to `false`.
* `is_notification_on` - (Optional) Whether a notification is sent to the approver. Defaults to `false`.

~> **NOTE:** `approver_id` requires the identity's storage-key UUID. To use a group as the approver, reference the `group_id` attribute of the `azuredevops_group` data source:

```hcl
data "azuredevops_group" "approvers" {
  project_id = azuredevops_project.example.id
  name       = "Release Administrators"
}

# ... within a stage block:
pre_deploy_approval {
  approval {
    approver_id = data.azuredevops_group.approvers.group_id
  }
}
```


An `approval_options` block supports the following:

* `required_approver_count` - (Optional) The number of approvals required to move the release forward. `0` means all approvers are required.
* `release_creator_can_be_approver` - (Optional) Whether the user requesting the release or deployment is allowed to approve it. Defaults to `false`.
* `auto_triggered_and_previous_environment_approved_can_be_skipped` - (Optional) Whether the approval can be skipped if the same approver approved the previous stage. Defaults to `false`.
* `enforce_identity_revalidation` - (Optional) Whether to revalidate the identity of the approver before completing the approval. Defaults to `false`.
* `timeout_in_minutes` - (Optional) The approval timeout in minutes. `0` means the default timeout of 30 days. Defaults to `0`.
* `execution_order` - (Optional) The order in which approvals are shown relative to gates. Possible values are `beforeGates`, `afterSuccessfulGates`, and `afterGatesAlways`. Defaults to `beforeGates`.

An `environment_options` block supports the following:

* `auto_link_work_items` - (Optional) Whether to automatically link work items to the deployment. Defaults to `false`.
* `badge_enabled` - (Optional) Whether the deployment badge is enabled for the stage. Defaults to `false`.
* `publish_deployment_status` - (Optional) Whether to publish deployment status to source repositories. Defaults to `true`.
* `pull_request_deployment_enabled` - (Optional) Whether pull request based deployments are enabled for the stage. Defaults to `false`.

An `execution_policy` block supports the following:

* `concurrency_count` - (Optional) The maximum number of deployments to run in parallel for the stage. Set to `0` for unlimited concurrency. Defaults to `0`.
* `queue_depth_count` - (Optional) The number of deployments allowed to be queued for the stage. Defaults to `0`.

An `environment_trigger` block supports the following:

* `trigger_type` - (Optional) The type of environment trigger. One of `rollbackRedeploy` or `deploymentGroupRedeploy`. Defaults to `rollbackRedeploy`.
* `action` - (Required) The action to take when the trigger fires. For an auto-redeploy trigger, use `LatestSuccessfulDeployment` to redeploy the last successful deployment.
* `event_types` - (Required) A list of event types that fire the trigger. For an auto-redeploy trigger, use `["MS.TF.DistributedTask.DeploymentFailed"]`.

A `pre_deployment_gates` / `post_deployment_gates` block supports the following:

* `is_enabled` - (Optional) Whether the gates are enabled. Defaults to `true`.
* `sampling_interval` - (Optional) The time in minutes between re-evaluation of gates. Defaults to `15`.
* `stabilization_time` - (Optional) The delay in minutes before gate evaluation starts. Defaults to `0`.
* `minimum_success_duration` - (Optional) The minimum duration in minutes for successful gate results to remain steady before the gate succeeds. Defaults to `0`.
* `timeout` - (Optional) The timeout in minutes after which gates fail. Defaults to `1440`.
* `gate` - (Optional) One or more `gate` blocks as documented below.

A `gate` block supports the following:

* `task` - (Required) One or more `task` blocks (same schema as the `deploy_phase` `task` block) that implement the gate, for example an "Invoke REST API" or "Query Azure Monitor alerts" task.

A `deploy_phase` block supports the following:

* `name` - (Required) The name of the deploy phase (job).
* `phase_type` - (Optional) The type of deploy phase. One of `agentBasedDeployment` (agent job), `runOnServer` (agentless/server job), or `machineGroupBasedDeployment` (deployment group job). Defaults to `agentBasedDeployment`.
* `rank` - (Optional) The order of the deploy phase. Defaults to its position in the list.
* `timeout_in_minutes` - (Optional) The job execution timeout in minutes. `0` means no timeout. Defaults to `0`.
* `job_cancel_timeout_in_minutes` - (Optional) The job cancel timeout in minutes when a deployment is cancelled. Defaults to `1`.
* `condition` - (Optional) The phase run condition. Defaults to `succeeded()`.
* `task` - (Required) One or more `task` blocks as documented below.

The following arguments apply to `agentBasedDeployment` and `machineGroupBasedDeployment` phases:

* `skip_artifacts_download` - (Optional) Whether to skip downloading artifacts for the phase. Defaults to `false`.
* `enable_access_token` - (Optional) Whether to include the OAuth access token in the deployment job. Defaults to `false`.

The following arguments apply only to `agentBasedDeployment` phases:

* `agent_queue_id` - (Optional) The agent **queue** ID the phase runs on. This is a project-scoped agent queue, not an organization-level agent pool. Source it from the `azuredevops_agent_queue` data source (for example the built-in `Azure Pipelines` queue) rather than hard-coding a numeric ID, since queue IDs differ per project.
* `agent_specification` - (Optional) The agent specification (image) identifier for the phase, for example `ubuntu-22.04`. Required when the phase runs on a Microsoft-hosted agent pool.

The following arguments apply only to `machineGroupBasedDeployment` phases:

* `deployment_group_id` - (Optional) The ID of the deployment group the phase deploys to.
* `tags` - (Optional) A set of deployment target tags used to filter which machines in the deployment group are deployed to.
* `deployment_health_option` - (Optional) The deployment health option, for example `OneTargetAtATime` or `Custom`.
* `health_percent` - (Optional) When `deployment_health_option` is `Custom`, the minimum percentage of targets that must remain healthy. Between `0` and `100`. Defaults to `0`.

A `task` block supports the following:

* `task_id` - (Required) The UUID of the task to run.
* `version` - (Optional) The version of the task. Defaults to `1.*`.
* `name` - (Optional) The display name of the task.
* `enabled` - (Optional) Whether the task is enabled. Defaults to `true`.
* `condition` - (Optional) The task run condition. Defaults to `succeeded()`.
* `inputs` - (Optional) A map of task input values.
* `override_inputs` - (Optional) A map of task input values to override.
* `environment` - (Optional) A map of environment variables to set for the task.
* `always_run` - (Optional) Whether the task runs even when a previous task has failed (unless the deployment was stopped). Defaults to `false`.
* `continue_on_error` - (Optional) Whether the deployment continues when the task fails. Defaults to `false`.
* `timeout_in_minutes` - (Optional) The task timeout in minutes. `0` means no timeout (or the maximum allowed). Defaults to `0`.
* `retry_count_on_task_failure` - (Optional) The number of times to retry the task if it fails. Defaults to `0`.
* `ref_name` - (Optional) The reference name of the task, used to reference its outputs from later tasks.
* `definition_type` - (Optional) The task definition type. For example `task` or `metaTask`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the release definition.
* `revision` - The revision number of the release definition.

## Import

Azure DevOps Release Definitions can be imported using the project name/id and release definition ID:

```sh
terraform import azuredevops_release_definition.example "Example Project/10"
```
