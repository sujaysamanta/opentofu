// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package funcs

import (
	"github.com/opentofu/opentofu/internal/lang/marks"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// SensitiveFunc returns a value identical to its argument except that
// OpenTofu will consider it to be sensitive.
var SensitiveFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "value",
			Type:             cty.DynamicPseudoType,
			AllowUnknown:     true,
			AllowNull:        true,
			AllowMarked:      true,
			AllowDynamicType: true,
		},
	},
	Type: func(args []cty.Value) (cty.Type, error) {
		// This function only affects the value's marks, so the result
		// type is always the same as the argument type.
		return args[0].Type(), nil
	},
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		return args[0].Mark(marks.Sensitive), nil
	},
})

// NonsensitiveFunc takes a sensitive value and returns the same value without
// the sensitive marking, effectively exposing the value.
var NonsensitiveFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "value",
			Type:             cty.DynamicPseudoType,
			AllowUnknown:     true,
			AllowNull:        true,
			AllowMarked:      true,
			AllowDynamicType: true,
		},
	},
	Type: func(args []cty.Value) (cty.Type, error) {
		// This function only affects the value's marks, so the result
		// type is always the same as the argument type.
		return args[0].Type(), nil
	},
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		v, m := args[0].Unmark()
		delete(m, marks.Sensitive) // remove the sensitive marking
		return v.WithMarks(m), nil
	},
})

// IsSensitiveFunc returns whether or not the value is sensitive.
var IsSensitiveFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "value",
			Type:             cty.DynamicPseudoType,
			AllowUnknown:     true,
			AllowNull:        true,
			AllowMarked:      true,
			AllowDynamicType: true,
		},
	},
	Type: func(args []cty.Value) (cty.Type, error) {
		return cty.Bool, nil
	},
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		return cty.BoolVal(args[0].HasMark(marks.Sensitive)), nil
	},
})

// FlipSensitiveFunc flips the sensitivity of a value.
// If the input is sensitive, it returns a non-sensitive value.
// If the input is non-sensitive, it returns a sensitive value.
var FlipSensitiveFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "value",
			Type:             cty.DynamicPseudoType,
			AllowUnknown:     true,
			AllowNull:        true,
			AllowMarked:      true,
			AllowDynamicType: true,
		},
	},
	Type: func(args []cty.Value) (cty.Type, error) {
		// This function only affects the value's marks, so the result
		// type is always the same as the argument type.
		return args[0].Type(), nil
	},
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		v := args[0]
		if v.HasMark(marks.Sensitive) {
			// If the value is sensitive, make it non-sensitive
			rawVal, valueMarks := v.Unmark()
			delete(valueMarks, marks.Sensitive)
			return rawVal.WithMarks(valueMarks), nil
		} else {
			// If the value is not sensitive, make it sensitive
			return v.Mark(marks.Sensitive), nil
		}
	},
})

func Sensitive(v cty.Value) (cty.Value, error) {
	return SensitiveFunc.Call([]cty.Value{v})
}

func Nonsensitive(v cty.Value) (cty.Value, error) {
	return NonsensitiveFunc.Call([]cty.Value{v})
}

func IsSensitive(v cty.Value) (cty.Value, error) {
	return IsSensitiveFunc.Call([]cty.Value{v})
}

func FlipSensitive(v cty.Value) (cty.Value, error) {
	return FlipSensitiveFunc.Call([]cty.Value{v})
}
