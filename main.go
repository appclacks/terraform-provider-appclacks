package main

import (
	"flag"

	"github.com/appclacks/terraform-provider/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "enable debug mode")
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{ProviderFunc: provider.Provider})
}
