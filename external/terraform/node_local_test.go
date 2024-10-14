// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package terraform

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/opengovern/terraform-package/external/addrs"
	"github.com/opengovern/terraform-package/external/configs"
	"github.com/opengovern/terraform-package/external/configs/hcl2shim"
	"github.com/opengovern/terraform-package/external/states"
)

func TestNodeLocalExecute(t *testing.T) {
	tests := []struct {
		Value string
		Want  interface{}
		Err   bool
	}{
		{
			"hello!",
			"hello!",
			false,
		},
		{
			"",
			"",
			false,
		},
		{
			"Hello, ${local.foo}",
			nil,
			true, // self-referencing
		},
	}

	for _, test := range tests {
		t.Run(test.Value, func(t *testing.T) {
			expr, diags := hclsyntax.ParseTemplate([]byte(test.Value), "", hcl.Pos{Line: 1, Column: 1})
			if diags.HasErrors() {
				t.Fatal(diags.Error())
			}

			n := &NodeLocal{
				Addr: addrs.LocalValue{Name: "foo"}.Absolute(addrs.RootModuleInstance),
				Config: &configs.Local{
					Expr: expr,
				},
			}
			ctx := &MockEvalContext{
				StateState: states.NewState().SyncWrapper(),

				EvaluateExprResult: hcl2shim.HCL2ValueFromConfigValue(test.Want),
			}

			err := n.Execute(ctx, walkApply)
			if (err != nil) != test.Err {
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				} else {
					t.Errorf("successful Eval; want error")
				}
			}

			ms := ctx.StateState.Module(addrs.RootModuleInstance)
			gotLocals := ms.LocalValues
			wantLocals := map[string]cty.Value{}
			if test.Want != nil {
				wantLocals["foo"] = hcl2shim.HCL2ValueFromConfigValue(test.Want)
			}

			if !reflect.DeepEqual(gotLocals, wantLocals) {
				t.Errorf(
					"wrong locals after Eval\ngot:  %swant: %s",
					spew.Sdump(gotLocals), spew.Sdump(wantLocals),
				)
			}
		})
	}

}
