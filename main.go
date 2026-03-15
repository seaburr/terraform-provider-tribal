package main

import (
	"context"
	"flag"
	"log"

	"github.com/seaburr/terraform-provider-tribal/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var version string = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/seaburr/tribal",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
