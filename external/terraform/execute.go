// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package terraform

import "github.com/opengovern/terraform-package/external/tfdiags"

// GraphNodeExecutable is the interface that graph nodes must implement to
// enable execution.
type GraphNodeExecutable interface {
	Execute(EvalContext, walkOperation) tfdiags.Diagnostics
}
