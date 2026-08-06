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
* `artifact` - (Optional) One or more `artifact` blocks as documented below.
* `artifact_source_trigger` - (Optional) One or more `artifact_source_trigger` blocks as documented below. Creates a release automatically when a new version of the linked artifact is available (continuous deployment).
* `schedule_trigger` - (Optional) One or more `schedule_trigger` blocks as documented below. Creates a release on a recurring schedule.

~> **NOTE:** Only `artifact_source_trigger` and `schedule_trigger` are managed by this resource. Other trigger types configured in the Azure DevOps UI (pull request, container image, package, and source repository triggers) are preserved on update but are not managed by Terraform.

* `stage` - (Required) One or more `stage` blocks as documented below.

A `variable` block supports the following:

* `name` - (Required) The name of the variable.
* `value` - (Optional) The value of the variable. Cannot be used together with `is_secret` / `secret_value`.
* `secret_value` - (Optional) The secret value of the variable. Used together with `is_secret = true`. This value is not returned by the API and is stored only in Terraform state.
* `is_secret` - (Optional) Whether the variable is a secret. Defaults to `false`. When `true`, set `secret_value` instead of `value`.
* `allow_override` - (Optional) Whether the variable can be overridden at release time. Defaults to `false`.

An `artifact` block supports the following:

* `alias` - (Required) The alias of the artifact. This is referenced by release variables and triggers (e.g. `$(myapp.BuildNumber)`).
* `type` - (Required) The artifact type. Common values are `Build`, `Git`, `GitHub`, `TFVC`, `Jenkins`, and `Nuget`.
* `is_primary` - (Optional) Whether this is the primary artifact. Defaults to `true`.
* `is_retained` - (Optional) Whether the artifact is retained by the release. Defaults to `false`.
* `definition_reference` - (Required) One or more `definition_reference` blocks describing the artifact source, as documented below.

A `definition_reference` block supports the following:

* `key` - (Required) The reference key. The required keys depend on `type`. For a `Build` artifact these are typically `project`, `definition`, and `defaultVersionType`.
* `id` - (Required) The ID value for the reference key.
* `name` - (Optional) The display name for the reference key.

~> **NOTE:** The Azure DevOps server may enrich `definition_reference` with additional keys it computes (for example `defaultVersionBranch` or `artifactSourceDefinitionUrl`). If a subsequent plan shows a diff, add the reported keys as `definition_reference` blocks to keep the configuration stable.

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

A `stage` block supports the following:

* `name` - (Required) The name of the stage (environment).
* `rank` - (Optional) The order of the stage. Defaults to its position in the list.
* `condition` - (Optional) One or more `condition` blocks controlling when the stage runs, as documented below. When omitted, the **first** stage defaults to starting on release creation (`ReleaseStarted`). At least one stage must start on release creation whenever a trigger is configured.
* `variable` - (Optional) One or more stage-scoped `variable` blocks (same schema as above).
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
