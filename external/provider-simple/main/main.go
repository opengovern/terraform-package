// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/kaytu-io/terraform-package/external/grpcwrap"
	"github.com/kaytu-io/terraform-package/external/plugin"
	simple "github.com/kaytu-io/terraform-package/external/provider-simple"
	"github.com/kaytu-io/terraform-package/external/tfplugin5"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		GRPCProviderFunc: func() tfplugin5.ProviderServer {
			return grpcwrap.Provider(simple.Provider())
		},
	})
}
