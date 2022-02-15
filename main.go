package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/itsmechlark/terraform-provider-cloud66/cloud66"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cloud66.Provider})
}
