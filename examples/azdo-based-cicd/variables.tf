variable "aad_users" {
  description = "A list of AAD user emails that will be granted access to the provisioned project. These are assumed to be part of an AAD group linked to the AzDO org. eg  [\"shanw_cicoria@microsoft.com\"]"
  type        = list(string)
}
