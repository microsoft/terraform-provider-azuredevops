echo "Delete build directory"
rm -fr build
echo "Delete provider from terraform plugins directory [darwin_amd64]"
rm -f ~/.terraform.d/plugins/darwin_amd64/terraform-provider-azuredevops_v0.0.1
echo "Delete provider from terraform plugins directory [linux_amd64]"
rm -f ~/.terraform.d/plugins/linux_amd64/terraform-provider-azuredevops_v0.0.1