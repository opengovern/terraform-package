// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"hashicorp/terraform/external/builtin/providers/terraform"
	"hashicorp/terraform/external/grpcwrap"
	"hashicorp/terraform/external/plugin"
	"hashicorp/terraform/external/tfplugin5"
)

func main() {
	// Provide a binary version of the internal terraform provider for testing
	plugin.Serve(&plugin.ServeOpts{
		GRPCProviderFunc: func() tfplugin5.ProviderServer {
			return grpcwrap.Provider(terraform.NewProvider())
		},
	})
}
