// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package globalref

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/opengovern/terraform-package/external/addrs"
	"github.com/opengovern/terraform-package/external/configs/configload"
	"github.com/opengovern/terraform-package/external/configs/configschema"
	"github.com/opengovern/terraform-package/external/initwd"
	"github.com/opengovern/terraform-package/external/providers"
	"github.com/opengovern/terraform-package/external/registry"
	"github.com/zclconf/go-cty/cty"
)

func testAnalyzer(t *testing.T, fixtureName string) *Analyzer {
	configDir := filepath.Join("testdata", fixtureName)

	loader, cleanup := configload.NewLoaderForTests(t)
	defer cleanup()

	inst := initwd.NewModuleInstaller(loader.ModulesDir(), loader, registry.NewClient(nil, nil))
	_, instDiags := inst.InstallModules(context.Background(), configDir, true, initwd.ModuleInstallHooksImpl{})
	if instDiags.HasErrors() {
		t.Fatalf("unexpected module installation errors: %s", instDiags.Err().Error())
	}
	if err := loader.RefreshModules(); err != nil {
		t.Fatalf("failed to refresh modules after install: %s", err)
	}

	cfg, loadDiags := loader.LoadConfig(configDir)
	if loadDiags.HasErrors() {
		t.Fatalf("unexpected configuration errors: %s", loadDiags.Error())
	}

	resourceTypeSchema := &configschema.Block{
		Attributes: map[string]*configschema.Attribute{
			"string": {Type: cty.String, Optional: true},
			"number": {Type: cty.Number, Optional: true},
			"any":    {Type: cty.DynamicPseudoType, Optional: true},
		},
		BlockTypes: map[string]*configschema.NestedBlock{
			"single": {
				Nesting: configschema.NestingSingle,
				Block: configschema.Block{
					Attributes: map[string]*configschema.Attribute{
						"z": {Type: cty.String, Optional: true},
					},
				},
			},
			"group": {
				Nesting: configschema.NestingGroup,
				Block: configschema.Block{
					Attributes: map[string]*configschema.Attribute{
						"z": {Type: cty.String, Optional: true},
					},
				},
			},
			"list": {
				Nesting: configschema.NestingList,
				Block: configschema.Block{
					Attributes: map[string]*configschema.Attribute{
						"z": {Type: cty.String, Optional: true},
					},
				},
			},
			"map": {
				Nesting: configschema.NestingMap,
				Block: configschema.Block{
					Attributes: map[string]*configschema.Attribute{
						"z": {Type: cty.String, Optional: true},
					},
				},
			},
			"set": {
				Nesting: configschema.NestingSet,
				Block: configschema.Block{
					Attributes: map[string]*configschema.Attribute{
						"z": {Type: cty.String, Optional: true},
					},
				},
			},
		},
	}
	schemas := map[addrs.Provider]*providers.Schemas{
		addrs.MustParseProviderSourceString("hashicorp/test"): {
			ResourceTypes: map[string]*configschema.Block{
				"test_thing": resourceTypeSchema,
			},
			DataSources: map[string]*configschema.Block{
				"test_thing": resourceTypeSchema,
			},
		},
	}

	return NewAnalyzer(cfg, schemas)
}
