// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/kaytu-io/terraform-package/external/grpcwrap"
	plugin "github.com/kaytu-io/terraform-package/external/plugin6"
	simple "github.com/kaytu-io/terraform-package/external/provider-simple-v6"
	"github.com/kaytu-io/terraform-package/external/tfplugin6"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		GRPCProviderFunc: func() tfplugin6.ProviderServer {
			return grpcwrap.Provider6(simple.Provider())
		},
	})
}
