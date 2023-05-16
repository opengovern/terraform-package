// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/kaytu-io/terraform-package/external/builtin/providers/terraform"
	"github.com/kaytu-io/terraform-package/external/grpcwrap"
	"github.com/kaytu-io/terraform-package/external/plugin"
	"github.com/kaytu-io/terraform-package/external/tfplugin5"
)

func main() {
	// Provide a binary version of the internal terraform provider for testing
	plugin.Serve(&plugin.ServeOpts{
		GRPCProviderFunc: func() tfplugin5.ProviderServer {
			return grpcwrap.Provider(terraform.NewProvider())
		},
	})
}
