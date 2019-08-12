
provider "azuredevops" {
  version = ">= 0.0.1"
}

resource "azuredevops_foo" "examplefoo" {
  fookey = "fooValue"
}