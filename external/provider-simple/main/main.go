// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"hashicorp/terraform/external/grpcwrap"
	"hashicorp/terraform/external/plugin"
	simple "hashicorp/terraform/external/provider-simple"
	"hashicorp/terraform/external/tfplugin5"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		GRPCProviderFunc: func() tfplugin5.ProviderServer {
			return grpcwrap.Provider(simple.Provider())
		},
	})
}
