package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/itsmechlark/terraform-provider-cloud66/cloud66"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: cloud66.Provider,
	}

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/itsmechlark/cloud66", opts)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
