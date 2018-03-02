package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/sneal/terraform-provider-bosh/bosh"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider})
}
