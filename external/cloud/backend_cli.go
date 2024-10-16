// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloud

import (
	"github.com/opengovern/terraform-package/external/backend"
	"github.com/opengovern/terraform-package/external/command/jsonformat"
)

// CLIInit implements backend.CLI
func (b *Cloud) CLIInit(opts *backend.CLIOpts) error {
	if cli, ok := b.local.(backend.CLI); ok {
		if err := cli.CLIInit(opts); err != nil {
			return err
		}
	}

	b.CLI = opts.CLI
	b.CLIColor = opts.CLIColor
	b.ContextOpts = opts.ContextOpts
	b.runningInAutomation = opts.RunningInAutomation
	b.input = opts.Input
	b.renderer = &jsonformat.Renderer{
		Streams:  opts.Streams,
		Colorize: opts.CLIColor,
	}

	return nil
}
