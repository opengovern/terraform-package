// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/opengovern/terraform-package/external/grpcwrap"
	"github.com/opengovern/terraform-package/external/plugin"
	simple "github.com/opengovern/terraform-package/external/provider-simple"
	"github.com/opengovern/terraform-package/external/tfplugin5"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		GRPCProviderFunc: func() tfplugin5.ProviderServer {
			return grpcwrap.Provider(simple.Provider())
		},
	})
}
