---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_check_business_hours"
description: |-
  Manages a business hours check.
---

# azuredevops_check_business_hours

Manages a business hours check on a resource within Azure DevOps.

## Example Usage

### Protect a service connection

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_serviceendpoint_generic" "example" {
  project_id            = azuredevops_project.example.id
  server_url            = "https://some-server.example.com"
  username              = "username"
  password              = "password"
  service_endpoint_name = "Example Generic"
  description           = "Managed by Terraform"
}

resource "azuredevops_check_business_hours" "example" {
  project_id           = azuredevops_project.example.id
  display_name         = "Managed by Terraform"
  target_resource_id   = azuredevops_serviceendpoint_generic.example.id
  target_resource_type = "endpoint"
  start_time           = "07:00"
  end_time             = "15:30"
  time_zone            = "UTC"
  monday               = true
  tuesday              = true
}
```

### Protect an environment

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_environment" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Environment"
}

resource "azuredevops_check_business_hours" "example" {
  project_id           = azuredevops_project.example.id
  display_name         = "Managed by Terraform"
  target_resource_id   = azuredevops_environment.example.id
  target_resource_type = "environment"
  start_time           = "07:00"
  end_time             = "15:30"
  time_zone            = "UTC"
  monday               = true
  tuesday              = true
}
```

### Protect an agent queue

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_agent_pool" "example" {
  name = "example-pool"
}

resource "azuredevops_agent_queue" "example" {
  project_id    = azuredevops_project.example.id
  agent_pool_id = azuredevops_agent_pool.example.id
}

resource "azuredevops_check_business_hours" "example" {
  project_id           = azuredevops_project.example.id
  display_name         = "Managed by Terraform"
  target_resource_id   = azuredevops_agent_queue.example.id
  target_resource_type = "queue"
  start_time           = "07:00"
  end_time             = "15:30"
  time_zone            = "UTC"
  monday               = true
  tuesday              = true
}
```

### Protect a repository

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_git_repository" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Empty Git Repository"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_check_business_hours" "example" {
  project_id           = azuredevops_project.example.id
  display_name         = "Managed by Terraform"
  target_resource_id   = "${azuredevops_project.example.id}.${azuredevops_git_repository.example.id}"
  target_resource_type = "repository"
  start_time           = "07:00"
  end_time             = "15:30"
  time_zone            = "UTC"
  monday               = true
  tuesday              = true
}
```

### Protect a variable group

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_variable_group" "example" {
  project_id   = azuredevops_project.example.id
  name         = "Example Variable Group"
  description  = "Example Variable Group Description"
  allow_access = true

  variable {
    name  = "key1"
    value = "val1"
  }

  variable {
    name         = "key2"
    secret_value = "val2"
    is_secret    = true
  }
}

resource "azuredevops_check_business_hours" "example" {
  project_id           = azuredevops_project.example.id
  display_name         = "Managed by Terraform"
  target_resource_id   = azuredevops_variable_group.example.id
  target_resource_type = "variablegroup"
  start_time           = "07:00"
  end_time             = "15:30"
  time_zone            = "UTC"
  monday               = true
  tuesday              = true
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The project ID.
* `target_resource_id` - (Required) The ID of the resource being protected by the check.
* `target_resource_type` - (Required) The type of resource being protected by the check. Valid values: `endpoint`, `environment`, `queue`, `repository`, `securefile`, `variablegroup`.
* `display_name` - (Required) The name of the business hours check displayed in the web UI.
* `start_time` - (Required) The beginning of the time period that this check will be allowed to pass, specified as 24-hour time with leading zeros.
* `end_time` - (Required) The end of the time period that this check will be allowed to pass, specified as 24-hour time with leading zeros.
* `time_zone` - (Required) The time zone this check will be evaluated in. See below for supported values.
* `monday` - (Optional) This check will pass on Mondays. Defaults to `false`.
* `tuesday` - (Optional) This check will pass on Tuesday. Defaults to `false`.
* `wednesday` - (Optional) This check will pass on Wednesdays. Defaults to `false`.
* `thursday` - (Optional) This check will pass on Thursdays. Defaults to `false`.
* `friday` - (Optional) This check will pass on Fridays. Defaults to `false`.
* `saturday` - (Optional) This check will pass on Saturdays. Defaults to `false`.
* `sunday` - (Optional) This check will pass on Sundays. Defaults to `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the check.

## Relevant Links

- [Define approvals and checks](https://learn.microsoft.com/en-us/azure/devops/pipelines/process/approvals?view=azure-devops&tabs=check-pass)

## Import

Importing this resource is not supported.

## Supported Time Zones

- AUS Central Standard Time
- AUS Eastern Standard Time
- Afghanistan Standard Time
- Alaskan Standard Time
- Aleutian Standard Time
- Altai Standard Time
- Arab Standard Time
- Arabian Standard Time
- Arabic Standard Time
- Argentina Standard Time
- Astrakhan Standard Time
- Atlantic Standard Time
- Aus Central W. Standard Time
- Azerbaijan Standard Time
- Azores Standard Time
- Bahia Standard Time
- Bangladesh Standard Time
- Belarus Standard Time
- Bougainville Standard Time
- Canada Central Standard Time
- Cape Verde Standard Time
- Caucasus Standard Time
- Cen. Australia Standard Time
- Central America Standard Time
- Central Asia Standard Time
- Central Brazilian Standard Time
- Central Europe Standard Time
- Central European Standard Time
- Central Pacific Standard Time
- Central Standard Time (Mexico)
- Central Standard Time
- Chatham Islands Standard Time
- China Standard Time
- Cuba Standard Time
- Dateline Standard Time
- E. Africa Standard Time
- E. Australia Standard Time
- E. Europe Standard Time
- E. South America Standard Time
- Easter Island Standard Time
- Eastern Standard Time (Mexico)
- Eastern Standard Time
- Egypt Standard Time
- Ekaterinburg Standard Time
- FLE Standard Time
- Fiji Standard Time
- GMT Standard Time
- GTB Standard Time
- Georgian Standard Time
- Greenland Standard Time
- Greenwich Standard Time
- Haiti Standard Time
- Hawaiian Standard Time
- India Standard Time
- Iran Standard Time
- Israel Standard Time
- Jordan Standard Time
- Kaliningrad Standard Time
- Kamchatka Standard Time
- Korea Standard Time
- Libya Standard Time
- Line Islands Standard Time
- Lord Howe Standard Time
- Magadan Standard Time
- Magallanes Standard Time
- Marquesas Standard Time
- Mauritius Standard Time
- Mid-Atlantic Standard Time
- Middle East Standard Time
- Montevideo Standard Time
- Morocco Standard Time
- Mountain Standard Time (Mexico)
- Mountain Standard Time
- Myanmar Standard Time
- N. Central Asia Standard Time
- Namibia Standard Time
- Nepal Standard Time
- New Zealand Standard Time
- Newfoundland Standard Time
- Norfolk Standard Time
- North Asia East Standard Time
- North Asia Standard Time
- North Korea Standard Time
- Omsk Standard Time
- Pacific SA Standard Time
- Pacific Standard Time (Mexico)
- Pacific Standard Time
- Pakistan Standard Time
- Paraguay Standard Time
- Qyzylorda Standard Time
- Romance Standard Time
- Russia Time Zone 10
- Russia Time Zone 11
- Russia Time Zone 3
- Russian Standard Time
- SA Eastern Standard Time
- SA Pacific Standard Time
- SA Western Standard Time
- SE Asia Standard Time
- Saint Pierre Standard Time
- Sakhalin Standard Time
- Samoa Standard Time
- Sao Tome Standard Time
- Saratov Standard Time
- Singapore Standard Time
- South Africa Standard Time
- South Sudan Standard Time
- Sri Lanka Standard Time
- Sudan Standard Time
- Syria Standard Time
- Taipei Standard Time
- Tasmania Standard Time
- Tocantins Standard Time
- Tokyo Standard Time
- Tomsk Standard Time
- Tonga Standard Time
- Transbaikal Standard Time
- Turkey Standard Time
- Turks And Caicos Standard Time
- US Eastern Standard Time
- US Mountain Standard Time
- UTC
- UTC+12
- UTC+13
- UTC-02
- UTC-08
- UTC-09
- UTC-11
- Ulaanbaatar Standard Time
- Venezuela Standard Time
- Vladivostok Standard Time
- Volgograd Standard Time
- W. Australia Standard Time
- W. Central Africa Standard Time
- W. Europe Standard Time
- W. Mongolia Standard Time
- West Asia Standard Time
- West Bank Standard Time
- West Pacific Standard Time
- Yakutsk Standard Time
- Yukon Standard Time
