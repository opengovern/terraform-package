// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	localexec "hashicorp/terraform/external/builtin/provisioners/local-exec"
	"hashicorp/terraform/external/grpcwrap"
	"hashicorp/terraform/external/plugin"
	"hashicorp/terraform/external/tfplugin5"
)

func main() {
	// Provide a binary version of the internal terraform provider for testing
	plugin.Serve(&plugin.ServeOpts{
		GRPCProvisionerFunc: func() tfplugin5.ProvisionerServer {
			return grpcwrap.Provisioner(localexec.New())
		},
	})
}
