// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package local

import (
	"flag"
	"os"
	"testing"

	_ "hashicorp/terraform/external/logging"
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}
