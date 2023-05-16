// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package collections

import (
	"hashicorp/terraform/external/command/jsonformat/computed"
	"hashicorp/terraform/external/plans"
)

type ProcessKey func(key string) computed.Diff

func TransformMap[Input any](before, after map[string]Input, keys []string, process ProcessKey) (map[string]computed.Diff, plans.Action) {
	current := plans.NoOp
	if before != nil && after == nil {
		current = plans.Delete
	}
	if before == nil && after != nil {
		current = plans.Create
	}

	elements := make(map[string]computed.Diff)
	for _, key := range keys {
		elements[key] = process(key)
		current = CompareActions(current, elements[key].Action)
	}

	return elements, current
}
