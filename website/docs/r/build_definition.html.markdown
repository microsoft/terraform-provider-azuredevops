---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_build_definition"
description: |-
  Manages a Build Definition within Azure DevOps organization.
---

# azuredevops_build_definition

Manages a Build Definition within Azure DevOps.

## Example Usage

### Azure DevOps
```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Repository"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_variable_group" "example" {
  project_id   = azuredevops_project.example.id
  name         = "Example Pipeline Variables"
  description  = "Managed by Terraform"
  allow_access = true

  variable {
    name  = "FOO"
    value = "BAR"
  }
}

resource "azuredevops_build_definition" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Build Definition"
  path       = "\\ExampleFolder"

  ci_trigger {
    use_yaml = false
  }

  schedules {
    branch_filter {
      include = ["master"]
      exclude = ["test", "regression"]
    }
    days_to_build              = ["Wed", "Sun"]
    schedule_only_with_changes = true
    start_hours                = 10
    start_minutes              = 59
    time_zone                  = "(UTC) Coordinated Universal Time"
  }

  repository {
    repo_type   = "TfsGit"
    repo_id     = azuredevops_git_repository.example.id
    branch_name = azuredevops_git_repository.example.default_branch
    yml_path    = "azure-pipelines.yml"
  }

  variable_groups = [
    azuredevops_variable_group.example.id
  ]

  variable {
    name  = "PipelineVariable"
    value = "Go Microsoft!"
  }

  variable {
    name         = "PipelineSecret"
    secret_value = "ZGV2cw"
    is_secret    = true
  }
}
```

### GitHub Enterprise
```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_github_enterprise" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example GitHub Enterprise"
  url                   = "https://github.contoso.com"
  description           = "Managed by Terraform"

  auth_personal {
    personal_access_token = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
}

resource "azuredevops_build_definition" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Build Definition"
  path       = "\\ExampleFolder"

  ci_trigger {
    use_yaml = false
  }

  repository {
    repo_type             = "GitHubEnterprise"
    repo_id               = "<GitHub Org>/<Repo Name>"
    github_enterprise_url = "https://github.company.com"
    branch_name           = "master"
    yml_path              = "azure-pipelines.yml"
    service_connection_id = azuredevops_serviceendpoint_github_enterprise.example.id
  }

  schedules {
    branch_filter {
      include = ["main"]
      exclude = ["test", "regression"]
    }
    days_to_build              = ["Wed", "Sun"]
    schedule_only_with_changes = true
    start_hours                = 10
    start_minutes              = 59
    time_zone                  = "(UTC) Coordinated Universal Time"
  }
}
```

### Build Completion Trigger
```hcl
resource "azuredevops_build_definition" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Build Definition"
  path       = "\\ExampleFolder"

  ci_trigger {
    use_yaml = false
  }

  repository {
    repo_type             = "GitHubEnterprise"
    repo_id               = "<GitHub Org>/<Repo Name>"
    github_enterprise_url = "https://github.company.com"
    branch_name           = "main"
    yml_path              = "azure-pipelines.yml"
    service_connection_id = azuredevops_serviceendpoint_github_enterprise.example.id
  }

  build_completion_trigger {
    build_definition_id = 10
    branch_filter {
      include = ["main"]
      exclude = ["test"]
    }
  }
}
```

### Pull Request Trigger
```hcl
data "azuredevops_serviceendpoint_github" "example" {
  project_id          = data.azuredevops_project.example.id
  service_endpoint_id = "00000000-0000-0000-0000-000000000000"
}

resource "azuredevops_build_definition" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Build Definition"
  path       = "\\ExampleFolder"

  ci_trigger {
    use_yaml = false
  }

  repository {
    repo_type             = "GitHub"
    repo_id               = "<GitHub Org>/<Repo Name>"
    branch_name           = "main"
    yml_path              = "azure-pipelines.yml"
    service_connection_id = data.azuredevops_serviceendpoint_github.example.id
  }

  pull_request_trigger {
    override {
      branch_filter {
        include = ["main"]
      }
    }
    forks {
      enabled       = false
      share_secrets = false
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The project ID or project name.
- `name` - (Required) The name of the build definition.
- `repository` - (Required) A `repository` block as documented below.

---
- `path` - (Optional) The folder path of the build definition.
- `agent_pool_name` - (Optional) The agent pool that should execute the build. Defaults to `Azure Pipelines`.
- `ci_trigger` - (Optional) A `ci_trigger` block as documented below.
- `pull_request_trigger` - (Optional) A `pull_request_trigger` block as documented below.
- `build_completion_trigger` - (Optional) A `build_completion_trigger` block as documented below.
- `variable_groups` - (Optional) A list of variable group IDs (integers) to link to the build definition.
- `variable` - (Optional) A list of `variable` blocks, as documented below.
- `features`- (Optional) A `features` blocks as documented below.
- `queue_status`- (Optional) The queue status of the build definition. Valid values: `enabled` or `paused` or `disabled`. Defaults to `enabled`.

---
`features` block supports the following:

  - `skip_first_run` (Optional) Trigger the pipeline to run after the creation. Defaults to `true`.
  
  ~> **Note** The first run(`skip_first_run = false`) will only be triggered on create. If the first run fails, the build definition will still be marked as successfully created. A warning message indicating the inability to run pipeline will be displayed.

---
`variable` block supports the following:

- `name` - (Required) The name of the variable.
- `value` - (Optional) The value of the variable.
- `secret_value` - (Optional) The secret value of the variable. Used when `is_secret` set to `true`.
- `is_secret` - (Optional) True if the variable is a secret. Defaults to `false`.
- `allow_override` - (Optional) True if the variable can be overridden. Defaults to `true`.

---
`repository` block supports the following:

- `branch_name` - (Optional) The branch name for which builds are triggered. Defaults to `master`.
- `repo_id` - (Required) The id of the repository. For `TfsGit` repos, this is simply the ID of the repository. For `Github` repos, this will take the form of `<GitHub Org>/<Repo Name>`. For `Bitbucket` repos, this will take the form of `<Workspace ID>/<Repo Name>`.
- `repo_type` - (Optional) The repository type. Valid values: `GitHub` or `TfsGit` or `Bitbucket` or `GitHub Enterprise`. Defaults to `GitHub`. If `repo_type` is `GitHubEnterprise`, must use existing project and GitHub Enterprise service connection.
- `service_connection_id` - (Optional) The service connection ID. Used if the `repo_type` is `GitHub` or `GitHubEnterprise`.
- `yml_path` - (Required) The path of the Yaml file describing the build definition.
- `github_enterprise_url` - (Optional) The Github Enterprise URL. Used if `repo_type` is `GithubEnterprise`.
- `report_build_status` - (Optional) Report build status. Default is true.

---
`ci_trigger` block supports the following:

- `use_yaml` - (Optional) Use the azure-pipeline file for the build configuration. Defaults to `false`.
- `override` - (Optional) Override the azure-pipeline file and use a this configuration for all builds.

---
`ci_trigger` `override` block supports the following:

- `branch_filter` - (Required) The branches to include and exclude from the trigger. A `branch_filter` block as documented below.
- `batch` - (Optional) If you set batch to true, when a pipeline is running, the system waits until the run is completed, then starts another run with all changes that have not yet been built. Defaults to `true`.
- `path_filter` - (Optional) Specify file paths to include or exclude. Note that the wildcard syntax is different between branches/tags and file paths.
- `max_concurrent_builds_per_branch` - (Optional) The number of max builds per branch. Defaults to `1`.
- `polling_interval` - (Optional) How often the external repository is polled. Defaults to `0`.
- `polling_job_id` - (Computed) This is the ID of the polling job that polls the external repository. Once the build definition is saved/updated, this value is set.

---
`build_completion_trigger` block supports the following:

- `build_definition_id` - (Required) The ID of the build pipeline will be triggered.
- `branch_filter` - (Required) The branches to include and exclude from the trigger. A `branch_filter` block as documented below.

---
`pull_request_trigger` block supports the following:

- `use_yaml` - (Optional) Use the azure-pipeline file for the build configuration. Defaults to `false`.
- `initial_branch` - (Optional) When use_yaml is true set this to the name of the branch that the azure-pipelines.yml exists on. Defaults to `Managed by Terraform`.
- `forks` - (Required) Set permissions for Forked repositories.
- `override` - (Optional) Override the azure-pipeline file and use this configuration for all builds.

---
`forks` block supports the following:

- `enabled` - (Required) Build pull requests from forks of this repository.
- `share_secrets` - (Required) Make secrets available to builds of forks.

---
`pull_request_trigger` `override` block supports the following:

- `branch_filter` - (Required) The branches to include and exclude from the trigger. A `branch_filter` block as documented below.
- `auto_cancel` - (Optional) . Defaults to `true`.
- `path_filter` - (Optional) Specify file paths to include or exclude. Note that the wildcard syntax is different between branches/tags and file paths.

---
`branch_filter` block supports the following:

- `include` - (Optional) List of branch patterns to include.
- `exclude` - (Optional) List of branch patterns to exclude.

---
`path_filter` block supports the following:

- `include` - (Optional) List of path patterns to include.
- `exclude` - (Optional) List of path patterns to exclude.

---
`schedules` block supports the following:

-> **Note:** Schedule pipeline will not use any schedules defined in the YAML file. To use schedules from the YAML file, delete all scheduled triggers.

- `days_to_build`: (Required) When to build. Valid values: `Mon`, `Tue`, `Wed`, `Thu`, `Fri`, `Sat`, `Sun`.
- `schedule_only_with_changes`: (Optional) Schedule builds if the source or pipeline has changed. Defaults to `true`.
- `start_hours`: (Optional) Build start hour. Defaults to `0`. Valid values: `0 ~ 23`.
- `start_minutes`: (Optional) Build start minute. Defaults to `0`. Valid values: `0 ~ 59`.
- `time_zone`: (Optional) Build time zone. Defaults to `(UTC) Coordinated Universal Time`. Valid values: 
  `(UTC-12:00) International Date Line West`,   
  `(UTC-11:00) Coordinated Universal Time-11`,   
  `(UTC-10:00) Aleutian Islands`,   
  `(UTC-10:00) Hawaii`,   
  `(UTC-09:30) Marquesas Islands`,   
  `(UTC-09:00) Alaska`,   
  `(UTC-09:00) Coordinated Universal Time-09`,   
  `(UTC-08:00) Baja California`,   
  `(UTC-08:00) Coordinated Universal Time-08`,   
  `(UTC-08:00) Pacific Time (US &Canada)`,   
  `(UTC-07:00) Arizona`,   
  `(UTC-07:00) Chihuahua, La Paz, Mazatlan`,   
  `(UTC-07:00) Mountain Time (US &Canada)`,   
  `(UTC-07:00) Yukon`,   
  `(UTC-06:00) Central America`,   
  `(UTC-06:00) Central Time (US &Canada)`,   
  `(UTC-06:00) Easter Island`,   
  `(UTC-06:00) Guadalajara, Mexico City, Monterrey`,   
  `(UTC-06:00) Saskatchewan`,   
  `(UTC-05:00) Bogota, Lima, Quito, Rio Branco`,   
  `(UTC-05:00) Chetumal`,   
  `(UTC-05:00) Eastern Time (US &Canada)`,   
  `(UTC-05:00) Haiti`,   
  `(UTC-05:00) Havana`,   
  `(UTC-05:00) Indiana (East)`,   
  `(UTC-05:00) Turks and Caicos`,   
  `(UTC-04:00) Asuncion`,   
  `(UTC-04:00) Atlantic Time (Canada)`,   
  `(UTC-04:00) Caracas`,   
  `(UTC-04:00) Cuiaba`,   
  `(UTC-04:00) Georgetown, La Paz, Manaus, San Juan`,   
  `(UTC-04:00) Santiago`,   
  `(UTC-03:30) Newfoundland`,   
  `(UTC-03:00) Araguaina`,   
  `(UTC-03:00) Brasilia`,   
  `(UTC-03:00) Cayenne, Fortaleza`,   
  `(UTC-03:00) City of Buenos Aires`,   
  `(UTC-03:00) Greenland`,   
  `(UTC-03:00) Montevideo`,   
  `(UTC-03:00) Punta Arenas`,   
  `(UTC-03:00) Saint Pierre and Miquelon`,   
  `(UTC-03:00) Salvador`,   
  `(UTC-02:00) Coordinated Universal Time-02`,   
  `(UTC-02:00) Mid-Atlantic - Old`,   
  `(UTC-01:00) Azores`,   
  `(UTC-01:00) Cabo Verde Is.`,   
  `(UTC) Coordinated Universal Time`,   
  `(UTC+00:00) Dublin, Edinburgh, Lisbon, London`,   
  `(UTC+00:00) Monrovia, Reykjavik`,   
  `(UTC+00:00) Sao Tome`,   
  `(UTC+01:00) Casablanca`,   
  `(UTC+01:00) Amsterdam, Berlin, Bern, Rome, Stockholm, Vienna`,   
  `(UTC+01:00) Belgrade, Bratislava, Budapest, Ljubljana, Prague`,   
  `(UTC+01:00) Brussels, Copenhagen, Madrid, Paris`,   
  `(UTC+01:00) Sarajevo, Skopje, Warsaw, Zagreb`,   
  `(UTC+01:00) West Central Africa`,   
  `(UTC+02:00) Amman`,   
  `(UTC+02:00) Athens, Bucharest`,   
  `(UTC+02:00) Beirut`,   
  `(UTC+02:00) Cairo`,   
  `(UTC+02:00) Chisinau`,   
  `(UTC+02:00) Damascus`,   
  `(UTC+02:00) Gaza, Hebron`,   
  `(UTC+02:00) Harare, Pretoria`,   
  `(UTC+02:00) Helsinki, Kyiv, Riga, Sofia, Tallinn, Vilnius`,   
  `(UTC+02:00) Jerusalem`,   
  `(UTC+02:00) Juba`,   
  `(UTC+02:00) Kaliningrad`,   
  `(UTC+02:00) Khartoum`,   
  `(UTC+02:00) Tripoli`,   
  `(UTC+02:00) Windhoek`,   
  `(UTC+03:00) Baghdad`,   
  `(UTC+03:00) Istanbul`,   
  `(UTC+03:00) Kuwait, Riyadh`,   
  `(UTC+03:00) Minsk`,   
  `(UTC+03:00) Moscow, St. Petersburg`,   
  `(UTC+03:00) Nairobi`,   
  `(UTC+03:00) Volgograd`,   
  `(UTC+03:30) Tehran`,   
  `(UTC+04:00) Abu Dhabi, Muscat`,   
  `(UTC+04:00) Astrakhan, Ulyanovsk`,   
  `(UTC+04:00) Baku`,   
  `(UTC+04:00) Izhevsk, Samara`,   
  `(UTC+04:00) Port Louis`,   
  `(UTC+04:00) Saratov`,   
  `(UTC+04:00) Tbilisi`,   
  `(UTC+04:00) Yerevan`,   
  `(UTC+04:30) Kabul`,   
  `(UTC+05:00) Ashgabat, Tashkent`,   
  `(UTC+05:00) Ekaterinburg`,   
  `(UTC+05:00) Islamabad, Karachi`,   
  `(UTC+05:00) Qyzylorda`,   
  `(UTC+05:30) Chennai, Kolkata, Mumbai, New Delhi`,   
  `(UTC+05:30) Sri Jayawardenepura`,   
  `(UTC+05:45) Kathmandu`,   
  `(UTC+06:00) Astana`,   
  `(UTC+06:00) Dhaka`,   
  `(UTC+06:00) Omsk`,   
  `(UTC+06:30) Yangon (Rangoon)`,   
  `(UTC+07:00) Bangkok, Hanoi, Jakarta`,   
  `(UTC+07:00) Barnaul, Gorno-Altaysk`,   
  `(UTC+07:00) Hovd`,   
  `(UTC+07:00) Krasnoyarsk`,   
  `(UTC+07:00) Novosibirsk`,   
  `(UTC+07:00) Tomsk`,   
  `(UTC+08:00) Beijing, Chongqing, Hong Kong, Urumqi`,   
  `(UTC+08:00) Irkutsk`,   
  `(UTC+08:00) Kuala Lumpur, Singapore`,   
  `(UTC+08:00) Perth`,   
  `(UTC+08:00) Taipei`,   
  `(UTC+08:00) Ulaanbaatar`,   
  `(UTC+08:45) Eucla`,   
  `(UTC+09:00) Chita`,   
  `(UTC+09:00) Osaka, Sapporo, Tokyo`,   
  `(UTC+09:00) Pyongyang`,   
  `(UTC+09:00) Seoul`,   
  `(UTC+09:00) Yakutsk`,   
  `(UTC+09:30) Adelaide`,   
  `(UTC+09:30) Darwin`,   
  `(UTC+10:00) Brisbane`,   
  `(UTC+10:00) Canberra, Melbourne, Sydney`,   
  `(UTC+10:00) Guam, Port Moresby`,   
  `(UTC+10:00) Hobart`,   
  `(UTC+10:00) Vladivostok`,   
  `(UTC+10:30) Lord Howe Island`,   
  `(UTC+11:00) Bougainville Island`,   
  `(UTC+11:00) Chokurdakh`,   
  `(UTC+11:00) Magadan`,   
  `(UTC+11:00) Norfolk Island`,   
  `(UTC+11:00) Sakhalin`,   
  `(UTC+11:00) Solomon Is., New Caledonia`,   
  `(UTC+12:00) Anadyr, Petropavlovsk-Kamchatsky`,   
  `(UTC+12:00) Auckland, Wellington`,   
  `(UTC+12:00) Coordinated Universal Time+12`,   
  `(UTC+12:00) Fiji`,   
  `(UTC+12:00) Petropavlovsk-Kamchatsky - Old`,   
  `(UTC+12:45) Chatham Islands`,   
  `(UTC+13:00) Coordinated Universal Time+13`,   
  `(UTC+13:00) Nuku'alofa`,   
  `(UTC+13:00) Samoa`,   
  `(UTC+14:00) Kiritimati Island`.
- `branch_filter` (Required) block supports the following:
  - `include` - (Optional) List of branch patterns to include.
  - `exclude` - (Optional) List of branch patterns to exclude.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the build definition
- `revision` - The revision of the build definition

---
The `schedules` block exports the following:

- `schedule_job_id` - The ID of the schedule job 

## Remarks

The path attribute can not end in `\` unless the path is the root value of `\`. 

Valid path values (yaml encoded) include:
- `\\`
- `\\ExampleFolder`
- `\\Nested\\Example Folder`

The value of `\\ExampleFolder\\` would be invalid.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Build Definitions](https://docs.microsoft.com/en-us/rest/api/azure/devops/build/definitions?view=azure-devops-rest-7.0)

## Import

Azure DevOps Build Definitions can be imported using the project name/definitions Id or by the project Guid/definitions Id, e.g.

```sh
terraform import azuredevops_build_definition.example "Example Project"/10
```

or

```sh
terraform import azuredevops_build_definition.example 00000000-0000-0000-0000-000000000000/0
```
