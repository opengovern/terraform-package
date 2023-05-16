// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package terraform

import (
	"hashicorp/terraform/external/addrs"
	"hashicorp/terraform/external/configs"
)

// GraphNodeAttachProviderMetaConfigs is an interface that must be implemented
// by nodes that want provider meta configurations attached.
type GraphNodeAttachProviderMetaConfigs interface {
	GraphNodeConfigResource

	// Sets the configuration
	AttachProviderMetaConfigs(map[addrs.Provider]*configs.ProviderMeta)
}
