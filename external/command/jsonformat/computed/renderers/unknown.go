// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package renderers

import (
	"fmt"

	"hashicorp/terraform/external/command/jsonformat/computed"

	"hashicorp/terraform/external/plans"
)

var _ computed.DiffRenderer = (*unknownRenderer)(nil)

func Unknown(before computed.Diff) computed.DiffRenderer {
	return &unknownRenderer{
		before: before,
	}
}

type unknownRenderer struct {
	NoWarningsRenderer

	before computed.Diff
}

func (renderer unknownRenderer) RenderHuman(diff computed.Diff, indent int, opts computed.RenderHumanOpts) string {
	if diff.Action == plans.Create {
		return fmt.Sprintf("(known after apply)%s", forcesReplacement(diff.Replace, opts))
	}

	// Never render null suffix for children of unknown changes.
	opts.OverrideNullSuffix = true
	return fmt.Sprintf("%s -> (known after apply)%s", renderer.before.RenderHuman(indent, opts), forcesReplacement(diff.Replace, opts))
}
