## Website Scaffolder

This application scaffolds the documentation for a Data Source/Resource.

**Note:** the documentation generated from this application is intended to be a
starting point, which when finished requires human review - rather than
generating a finished product. 

## Example Usage

```
$ go run main.go -name azuredevops_agent_pool -brand-name "Agent Pool" -type "resource" -resource-id "00000000-0000-0000-0000-000000000000" -website-path ../../../../website/
```

## Arguments

* `-name` - (Required) The Name used for the Resource in Terraform e.g.
  `azuredevops_agent_pool`

* `-brand-name` - (Required) The Brand Name used for this Resource in e.g.
  `Agent Pool` or `Agent Queue`

* `-type` - (Required) The Type of Documentation to generate. Possible values
  are `data` (for a Data Source) or `resource` (for a Resource).

* `-resource-id` - (Required when scaffolding a Resource) An  Resource ID which
  can be used as a placeholder in the import documentation.

* `-website-path` - (Required) The path to the `./website` directory in the root
  of this repository.
