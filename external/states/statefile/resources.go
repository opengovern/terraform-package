package statefile

import (
	"github.com/kaytu-io/terraform-package/external/states"
	"io"
	"sort"

	"github.com/kaytu-io/terraform-package/external/addrs"
	"github.com/kaytu-io/terraform-package/external/configs/hcl2shim"

	ctyjson "github.com/zclconf/go-cty/cty/json"
)

// Returns the list of ARNs for all resources in the state file
func GetResourcesArn(f io.Reader) []string {
	result, err := Read(f)
	if err != nil {
		panic(err)
	}

	state := result.State
	return GetArnsFromStateFile(state)
}

func GetArnsFromStateFile(state *states.State) []string {
	arns := make([]string, 0)
	for _, ms := range state.Modules {
		addrsOrder := make([]addrs.AbsResourceInstance, 0, len(ms.Resources))
		for _, rs := range ms.Resources {
			for ik := range rs.Instances {
				addrsOrder = append(addrsOrder, rs.Addr.Instance(ik))
			}
		}

		sort.Slice(addrsOrder, func(i, j int) bool {
			return addrsOrder[i].Less(addrsOrder[j])
		})

		for _, fakeAbsAddr := range addrsOrder {
			addr := fakeAbsAddr.Resource
			is := ms.ResourceInstance(addr)
			var attributes map[string]string
			if obj := is.Current; obj != nil {
				switch {
				case obj.AttrsFlat != nil:
					attributes = obj.AttrsFlat
				case obj.AttrsJSON != nil:
					ty, err := ctyjson.ImpliedType(obj.AttrsJSON)
					if err == nil {
						val, err := ctyjson.Unmarshal(obj.AttrsJSON, ty)
						if err == nil {
							attributes = hcl2shim.FlatmapValueFromHCL2(val)
						}
					}
				}
			}
			for key, value := range attributes {
				if key == "arn" {
					arns = append(arns, value)
				}
			}
		}
	}
	return arns
}

// Returns the list of resource types
func GetResourcesTypes(f io.Reader) []string {
	result, err := Read(f)
	if err != nil {
		panic(err)
	}

	state := result.State
	return GetResourcesTypesFromState(state)
}

func GetResourcesTypesFromState(state *states.State) []string {
	types := make([]string, 0)

	for _, ms := range state.Modules {
		for _, re := range ms.Resources {
			types = append(types, re.Addr.Resource.Type)
		}
	}
	return types
}
