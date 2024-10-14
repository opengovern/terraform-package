// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package terraform

import (
	backendInit "github.com/opengovern/terraform-package/external/backend/init"
)

func init() {
	// Initialize the backends
	backendInit.Init(nil)
}
