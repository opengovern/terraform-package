// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package differ

import (
	"github.com/zclconf/go-cty/cty"

	"hashicorp/terraform/external/command/jsonformat/collections"
	"hashicorp/terraform/external/command/jsonformat/computed"
	"hashicorp/terraform/external/command/jsonformat/computed/renderers"
	"hashicorp/terraform/external/command/jsonformat/structured"
	"hashicorp/terraform/external/command/jsonprovider"
	"hashicorp/terraform/external/plans"
)

func computeAttributeDiffAsMap(change structured.Change, elementType cty.Type) computed.Diff {
	mapValue := change.AsMap()
	elements, current := collections.TransformMap(mapValue.Before, mapValue.After, mapValue.AllKeys(), func(key string) computed.Diff {
		value := mapValue.GetChild(key)
		if !value.RelevantAttributes.MatchesPartial() {
			// Mark non-relevant attributes as unchanged.
			value = value.AsNoOp()
		}
		return ComputeDiffForType(value, elementType)
	})
	return computed.NewDiff(renderers.Map(elements), current, change.ReplacePaths.Matches())
}

func computeAttributeDiffAsNestedMap(change structured.Change, attributes map[string]*jsonprovider.Attribute) computed.Diff {
	mapValue := change.AsMap()
	elements, current := collections.TransformMap(mapValue.Before, mapValue.After, mapValue.ExplicitKeys(), func(key string) computed.Diff {
		value := mapValue.GetChild(key)
		if !value.RelevantAttributes.MatchesPartial() {
			// Mark non-relevant attributes as unchanged.
			value = value.AsNoOp()
		}
		return computeDiffForNestedAttribute(value, &jsonprovider.NestedType{
			Attributes:  attributes,
			NestingMode: "single",
		})
	})
	return computed.NewDiff(renderers.NestedMap(elements), current, change.ReplacePaths.Matches())
}

func computeBlockDiffsAsMap(change structured.Change, block *jsonprovider.Block) (map[string]computed.Diff, plans.Action) {
	mapValue := change.AsMap()
	return collections.TransformMap(mapValue.Before, mapValue.After, mapValue.ExplicitKeys(), func(key string) computed.Diff {
		value := mapValue.GetChild(key)
		if !value.RelevantAttributes.MatchesPartial() {
			// Mark non-relevant attributes as unchanged.
			value = value.AsNoOp()
		}
		return ComputeDiffForBlock(value, block)
	})
}
