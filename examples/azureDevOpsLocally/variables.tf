// Variaveis relacionados a organização existente
variable "org_url" {
  description = "URL of organization"
  type= string 
}

variable "org_token" {
  description = "TOKEN of organization"
  type= string 
}

// Variaveis relacionado ao projeto ao ser criado
variable "project_name" {
  description = "The name that the project will be called"
  type= string
}

variable "description" {
  description = "Description of project"
  type= string 
}

variable "visibility" {
  description = "Specify privacy of project (private or public)"
  type= string
  default = "private" 
}

variable "version_control" {
  description = "Set system control of project"
  type= string
  default = "Git"  
}

variable "work_item_template" {
  description = "Set work item of project"
  type= string
  default = "Agile"
}
