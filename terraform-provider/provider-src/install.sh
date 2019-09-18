echo "Copy build artifact to terraform plugins directory [darwin_amd64]"
mkdir -p ~/.terraform.d/plugins/darwin_amd64
cp -f build/terraform-provider-azuredevops_v0.0.1 ~/.terraform.d/plugins/darwin_amd64

echo "Copy build artifact to terraform plugins directory [linux_amd64]"
mkdir -p ~/.terraform.d/plugins/linux_amd64
cp -f build/terraform-provider-azuredevops_v0.0.1 ~/.terraform.d/plugins/linux_amd64
