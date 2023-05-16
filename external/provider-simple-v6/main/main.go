// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"hashicorp/terraform/external/grpcwrap"
	plugin "hashicorp/terraform/external/plugin6"
	simple "hashicorp/terraform/external/provider-simple-v6"
	"hashicorp/terraform/external/tfplugin6"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		GRPCProviderFunc: func() tfplugin6.ProviderServer {
			return grpcwrap.Provider6(simple.Provider())
		},
	})
}
