package main

import (
	"flag"

	"github.com/appclacks/terraform-provider-appclacks/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// Provider documentation generation.
//
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name appclacks

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "enable debug mode")
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{ProviderFunc: provider.Provider})
}
