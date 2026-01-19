package main

import (
	"context"
	"flag"
	"log"

	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/microsoft/terraform-provider-azuredevops/internal/provider"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/microsoft/azuredevops",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), func() fwprovider.Provider { return provider.New() }, opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
