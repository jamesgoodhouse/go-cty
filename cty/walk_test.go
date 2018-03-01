package cty

import (
	"fmt"
	"testing"
)

func TestWalk(t *testing.T) {
	type Call struct {
		Path string
		Type string
	}

	val := ObjectVal(map[string]Value{
		"string":       StringVal("hello"),
		"number":       NumberIntVal(10),
		"bool":         True,
		"list":         ListVal([]Value{True}),
		"list_empty":   ListValEmpty(Bool),
		"set":          SetVal([]Value{True}),
		"set_empty":    ListValEmpty(Bool),
		"tuple":        TupleVal([]Value{True}),
		"tuple_empty":  EmptyTupleVal,
		"map":          MapVal(map[string]Value{"true": True}),
		"map_empty":    MapValEmpty(Bool),
		"object":       ObjectVal(map[string]Value{"true": True}),
		"object_empty": EmptyObjectVal,
		"null":         NullVal(List(String)),
		"unknown":      NullVal(Map(Bool)),
	})

	gotCalls := map[Call]struct{}{}
	wantCalls := []Call{
		{`cty.Path(nil)`, "object"},
		{`cty.Path{cty.GetAttrStep{Name:"string"}}`, "string"},
		{`cty.Path{cty.GetAttrStep{Name:"number"}}`, "number"},
		{`cty.Path{cty.GetAttrStep{Name:"bool"}}`, "bool"},
		{`cty.Path{cty.GetAttrStep{Name:"list"}}`, "list of bool"},
		{`cty.Path{cty.GetAttrStep{Name:"list"}, cty.IndexStep{Key:cty.NumberIntVal(0)}}`, "bool"},
		{`cty.Path{cty.GetAttrStep{Name:"list_empty"}}`, "list of bool"},
		{`cty.Path{cty.GetAttrStep{Name:"set"}}`, "set of bool"},
		{`cty.Path{cty.GetAttrStep{Name:"set"}, cty.IndexStep{Key:cty.True}}`, "bool"},
		{`cty.Path{cty.GetAttrStep{Name:"set_empty"}}`, "list of bool"},
		{`cty.Path{cty.GetAttrStep{Name:"tuple"}}`, "tuple"},
		{`cty.Path{cty.GetAttrStep{Name:"tuple"}, cty.IndexStep{Key:cty.NumberIntVal(0)}}`, "bool"},
		{`cty.Path{cty.GetAttrStep{Name:"tuple_empty"}}`, "tuple"},
		{`cty.Path{cty.GetAttrStep{Name:"map"}, cty.IndexStep{Key:cty.StringVal("true")}}`, "bool"},
		{`cty.Path{cty.GetAttrStep{Name:"map"}}`, "map of bool"},
		{`cty.Path{cty.GetAttrStep{Name:"map_empty"}}`, "map of bool"},
		{`cty.Path{cty.GetAttrStep{Name:"object"}}`, "object"},
		{`cty.Path{cty.GetAttrStep{Name:"object"}, cty.GetAttrStep{Name:"true"}}`, "bool"},
		{`cty.Path{cty.GetAttrStep{Name:"object_empty"}}`, "object"},
		{`cty.Path{cty.GetAttrStep{Name:"null"}}`, "list of string"},
		{`cty.Path{cty.GetAttrStep{Name:"unknown"}}`, "map of bool"},
	}

	err := Walk(val, func(path Path, val Value) (bool, error) {
		gotCalls[Call{
			Path: fmt.Sprintf("%#v", path),
			Type: val.Type().FriendlyName(),
		}] = struct{}{}
		return true, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(gotCalls) != len(wantCalls) {
		t.Errorf("wrong number of calls %d; want %d", len(gotCalls), len(wantCalls))
	}

	for gotCall := range gotCalls {
		t.Logf("got call {%#q, %q}", gotCall.Path, gotCall.Type)
	}

	for _, wantCall := range wantCalls {
		if _, has := gotCalls[wantCall]; !has {
			t.Errorf("missing call {%#q, %q}", wantCall.Path, wantCall.Type)
		}
	}
}