// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package differ

import (
	"github.com/zclconf/go-cty/cty"

	"github.com/kaytu-io/terraform-package/external/command/jsonformat/computed"
	"github.com/kaytu-io/terraform-package/external/command/jsonformat/computed/renderers"
	"github.com/kaytu-io/terraform-package/external/command/jsonformat/structured"
)

func computeAttributeDiffAsPrimitive(change structured.Change, ctype cty.Type) computed.Diff {
	return asDiff(change, renderers.Primitive(change.Before, change.After, ctype))
}
