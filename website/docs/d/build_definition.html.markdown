---
layout: "azuredevops"
page_title: "AzureDevops: Data Source: azuredevops_build_definition"
description: |-
  Gets information about an existing Build Definition.
---

# Data Source: azuredevops_build_definition

Use this data source to access information about an existing Build Definition.

## Example Usage

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

data "azuredevops_build_definition" "example" {
  project_id = data.azuredevops_project.example.id
  name       = "existing"
}

output "id" {
  value = data.azuredevops_build_definition.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Build Definition.

* `project_id` - (Required) The ID of the project.

---

* `path` - (Optional) The path of the build definition. Default to `\`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Build Definition.

* `agent_pool_name` - The agent pool that should execute the build.

* `ci_trigger` - A `ci_trigger` block as defined below.

* `pull_request_trigger` - A `pull_request_trigger` block as defined below.

* `repository` - A `repository` block as defined below.

* `revision` - The revision of the build definition.

* `schedules` - A `schedules` block as defined below.

* `variable` - A `variable` block as defined below.

* `variable_groups` - A list of variable group IDs.

* `queue_status` - The queue status of the build definition.

* `agent_specification`- The Agent Specification to run the pipelines. Example: `windows-2019`, `windows-latest`, `macos-13` etc.

* `job_authorization_scope`- The job authorization scope for builds queued against this definition.

* `jobs`- A `jobs` blocks as documented below.

---

`jobs` block supports the following:

* `name` - The name of the job.

* `ref_name` - The reference name of the job, can be used to define the job dependencies.

* `condition` - Specifies when this job should run. Can **Custom conditions** to specify more complex conditions. More details: [Pipeline conditions](https://learn.microsoft.com/en-us/azure/devops/pipelines/process/conditions?view=azure-devops)

* `target`- A `target` blocks as documented below.

* `job_timeout_in_minutes` - The job execution timeout (in minutes) for builds queued against this definition.

* `job_cancel_timeout_in_minutes` - The job cancel timeout (in minutes) for builds cancelled by user for this definition.

* `job_authorization_scope`- The job authorization scope for builds queued against this definition.

* `allow_scripts_auth_access_option`- Enables scripts and other processes launched by tasks to access the OAuth token through the `System.AccessToken` variable.

* `dependencies`- A `dependencies` blocks as documented below. Define the job dependencies.

---

`dependencies` block supports the following:

* `scope` The job reference name that depends on. Reference to `jobs.ref_name`

---

`target` block supports the following:

* `type`  The job type.

* `execution_options`- A `execution_options` blocks as documented below.

* `demands` - A list of demands that represents the agent capabilities required by this build. Example: `git`

---

`execution_options` block supports the following:

* `type`- The execution type of the Job.

* `multipliers` -  A list of comma separated configuration variables to use. These are defined on the Variables tab. For example, OperatingSystem, Browser will run the tasks for both variables.

* `max_concurrency` - Limit the number of agents to be used. If job type is `AgentlessJob`, the concurrency is not configurable and is fixed to 50.

* `continue_on_error` - Whether to continue the job when an error occurs.


---

A `branch_filter` block exports the following:

* `exclude` - A `exclude` block as defined below.

* `include` - A `include` block as defined below.

---

A `ci_trigger` block exports the following:

* `override` - A `override` block as defined below.

* `use_yaml` - Use the azure-pipeline file for the build configuration.

---

A `ci_trigger` `override` block supports the following:

* `batch` - If batch is true, when a pipeline is running, the system waits until the run is completed, then starts another run with all changes that have not yet been built.

* `branch_filter` - The branches to include and exclude from the trigger.

* `path_filter` - Specify file paths to include or exclude. Note that the wildcard syntax is different between branches/tags and file paths.

* `max_concurrent_builds_per_branch` - The number of max builds per branch.

* `polling_interval` - How often the external repository is polled.

* `polling_job_id` - This is the ID of the polling job that polls the external repository. Once the build definition is saved/updated, this value is set.

---

A `branch_filter` block supports the following:

* `include` - (Optional) List of branch patterns to include.

* `exclude` - (Optional) List of branch patterns to exclude.

---

A `path_filter` block supports the following:

* `include` - (Optional) List of path patterns to include.
 
* `exclude` - (Optional) List of path patterns to exclude.
 
---

A `pull_request_trigger` block exports the following:

* `comment_required` - Is a comment required on the PR?

* `forks` - A `forks` block as defined above.

* `initial_branch` - When use_yaml is true set this to the name of the branch that the azure-pipelines.yml exists on.

* `override` - A `override` block as defined below.

* `use_yaml` - Use the azure-pipeline file for the build configuration.

---

A `forks` block exports the following:

* `enabled` - Build pull requests from forks of this repository.

* `share_secrets` - Make secrets available to builds of forks.

---

A `pull_request_trigger` `override` block supports the following:

* `auto_cancel` -Should further updates to a PR cancel an in progress validation?

* `branch_filter` - The branches to include and exclude from the trigger. A `branch_filter` block as defined above.

* `path_filter` - The file paths to include or exclude. A `path_filter` block as defined above.

---

A `repository` block exports the following:

* `branch_name` - The branch name for which builds are triggered.

* `github_enterprise_url` - The Github Enterprise URL.

* `repo_id` - The id of the repository.

* `repo_type` - The repository type.

* `report_build_status` - Report build status.

* `service_connection_id` - The service connection ID.

* `yml_path` - The path of the Yaml file describing the build definition.

---

A `schedules` block exports the following:

* `branch_filter` - A `branch_filter` block as defined above.

* `days_to_build` - A list of days to build on.

* `schedule_job_id` - The ID of the schedule job.

* `schedule_only_with_changes` - Schedule builds if the source or pipeline has changed.

* `start_hours` - Build start hour.

* `start_minutes` - Build start minute.

* `time_zone` - Build time zone.

---

A `variable` block exports the following:

* `allow_override` - `true` if the variable can be overridden.

* `is_secret` - `true` if the variable is a secret.

* `name` - The name of the variable.

* `secret_value` - The secret value of the variable.

* `value` - The value of the variable.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minute) Used when retrieving the Build Definition.
