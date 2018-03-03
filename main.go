package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/sneal/terraform-provider-opsman/opsman"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: opsman.Provider})
}
