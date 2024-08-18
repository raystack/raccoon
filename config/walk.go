package config

import "reflect"

type CfgRef struct {
	Ref  any
	Meta reflect.StructField
}

// Walk returns a pointer to leaf nodes of the configuration
// tree along with it's struct metadata.
func Walk() []CfgRef {
	// todo: use iter.Seq[CfgRef] once we upgrade to go v1.23

	return walk(reflect.ValueOf(cfg))
}

func walk(value reflect.Value) []CfgRef {
	refs := []CfgRef{}
	for idx := range value.NumField() {
		field := value.Field(idx)
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		if field.Kind() == reflect.Struct {
			refs = append(refs, walk(field)...)
			continue
		}
		ref := CfgRef{
			Ref:  field.Addr().Interface(),
			Meta: value.Type().Field(idx),
		}
		refs = append(refs, ref)
	}
	return refs
}
