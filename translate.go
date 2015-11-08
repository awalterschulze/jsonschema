//  Copyright 2015 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package jsonschema

import (
	"encoding/json"
	"fmt"
	"github.com/katydid/katydid/funcs"
	"github.com/katydid/katydid/relapse/ast"
	"github.com/katydid/katydid/relapse/combinator"
	"sort"
)

//TODO
func TranslateDraft4(jsonSchema []byte) (*relapse.Grammar, error) {
	schema := &Schema{}
	if err := json.Unmarshal(jsonSchema, schema); err != nil {
		panic(err)
	}
	p, err := translate(schema)
	if err != nil {
		return nil, err
	}
	return relapse.NewGrammar(relapse.RefLookup(map[string]*relapse.Pattern{"main": p})), nil
}

func translate(schema *Schema) (*relapse.Pattern, error) {
	if schema.Id != nil {
		return nil, fmt.Errorf("id not supported")
	}
	if schema.Default != nil {
		return nil, fmt.Errorf("default not supported")
	}
	if schema.HasNumericConstraints() || (schema.Type.Single() &&
		schema.Type.HasNumeric()) {
		p, err := translateNumeric(schema)
		return p, err
	}
	if schema.HasStringConstraints() || (schema.Type.Single() &&
		schema.Type.HasString()) {
		p, err := translateString(schema)
		return p, err
	}
	if schema.HasArrayConstraints() || (schema.Type.Single() &&
		schema.Type.HasArray()) {
		return nil, fmt.Errorf("array not supported")
	}
	if schema.HasObjectConstraints() {
		p, err := translateObject(schema)
		return p, err
	}
	if schema.HasInstanceConstraints() {
		p, err := translateInstance(schema)
		return p, err
	}

	if schema.Ref != nil {
		return nil, fmt.Errorf("ref not supported")
	}
	if schema.Format != nil {
		return nil, fmt.Errorf("format not supported")
	}
	return relapse.NewEmptySet(), nil
}

func translates(schemas []*Schema) ([]*relapse.Pattern, error) {
	ps := make([]*relapse.Pattern, len(schemas))
	for i := range schemas {
		var err error
		ps[i], err = translate(schemas[i])
		if err != nil {
			return nil, err
		}
	}
	return ps, nil
}

func rest(patterns []*relapse.Pattern, index int) []*relapse.Pattern {
	return append(patterns[:index], patterns[index+1:]...)
}

func translateInstance(schema *Schema) (*relapse.Pattern, error) {
	if len(schema.Definitions) > 0 {
		return nil, fmt.Errorf("definitions not supported")
	}
	if len(schema.Enum) > 0 {
		return nil, fmt.Errorf("enum not supported")
	}
	if schema.Type != nil {
		types := *schema.Type
		if len(types) == 1 {
			return translateType(types[0])
		}
		ps := make([]*relapse.Pattern, len(types))
		for i := range types {
			var err error
			ps[i], err = translateType(types[i])
			if err != nil {
				return nil, err
			}
		}
		return relapse.NewOr(ps...), nil
	}
	if len(schema.AllOf) > 0 {
		ps, err := translates(schema.AllOf)
		if err != nil {
			return nil, err
		}
		return relapse.NewAnd(ps...), nil
	}
	if len(schema.AnyOf) > 0 {
		ps, err := translates(schema.AnyOf)
		if err != nil {
			return nil, err
		}
		return relapse.NewOr(ps...), nil
	}
	if len(schema.OneOf) > 0 {
		ps, err := translates(schema.OneOf)
		if err != nil {
			return nil, err
		}
		if len(ps) == 0 {
			return nil, fmt.Errorf("oneof of zero schemas not supported")
		}
		if len(ps) == 1 {
			return ps[0], nil
		}
		orps := make([]*relapse.Pattern, len(ps))
		for i := range ps {
			other := rest(ps, i)
			orps[i] = relapse.NewAnd(
				ps[i],
				relapse.NewNot(
					relapse.NewOr(other...),
				),
			)
		}
		return relapse.NewOr(orps...), nil
	}
	if schema.Not != nil {
		p, err := translate(schema.Not)
		if err != nil {
			return nil, err
		}
		return relapse.NewNot(p), nil
	}
	panic("unreachable object")
}

func translateType(typ SimpleType) (*relapse.Pattern, error) {
	switch typ {
	case TypeArray:
		return nil, fmt.Errorf("type array not supported")
	case TypeBoolean:
		return combinator.Value(funcs.TypeBool(funcs.BoolVar())), nil
	case TypeInteger:
		return combinator.Value(funcs.TypeDouble(Integer())), nil
	case TypeNull:
		return relapse.NewEmpty(), nil
	case TypeNumber:
		return combinator.Value(funcs.TypeDouble(Number())), nil
	case TypeObject:
		return nil, fmt.Errorf("type object not supported")
	case TypeString:
		return combinator.Value(funcs.TypeString(funcs.StringVar())), nil
	}
	panic(fmt.Sprintf("unknown simpletype: %s", typ))
}

func translateObject(schema *Schema) (*relapse.Pattern, error) {
	if schema.MaxProperties != nil {
		return nil, fmt.Errorf("maxProperties not supported")
	}
	if schema.MinProperties > 0 {
		return nil, fmt.Errorf("minProperties not supported")
	}
	required := make(map[string]struct{})
	for _, req := range schema.Required {
		required[req] = struct{}{}
	}
	requiredIf := make(map[string][]string)
	moreProperties := make(map[string]*Schema)
	if schema.Dependencies != nil {
		deps := *schema.Dependencies
		for name, dep := range deps {
			if len(dep.RequiredProperty) > 0 {
				requiredIf[name] = deps[name].RequiredProperty
			} else {
				moreProperties[name] = deps[name].Schema
			}
		}
	}
	additional := relapse.NewZAny()
	if schema.AdditionalProperties != nil {
		if schema.AdditionalProperties.Bool != nil && !(*schema.AdditionalProperties.Bool) {
			additional = relapse.NewEmpty()
		} else if schema.AdditionalProperties.Type != TypeUnknown {
			typ, err := translateType(schema.AdditionalProperties.Type)
			if err != nil {
				return nil, err
			}
			additional = relapse.NewZeroOrMore(
				relapse.NewTreeNode(relapse.NewAnyName(), typ),
			)
		}
	}
	names := []string{}
	for name, _ := range schema.Properties {
		names = append(names, name)
	}
	sort.Strings(names)
	patterns := make(map[string]*relapse.Pattern)
	for _, name := range names {
		child, err := translate(schema.Properties[name])
		if err != nil {
			return nil, err
		}
		patterns[name] = relapse.NewTreeNode(relapse.NewStringName(name), child)
	}
	_ = additional
	for _, name := range names {
		if requires, ok := requiredIf[name]; ok {
			_ = requires
		}
		if s, ok := moreProperties[name]; ok {
			_ = s
		}
		if _, ok := required[name]; !ok {
			//prop = optional(prop)
		}
	}
	if len(schema.PatternProperties) > 0 {
		return nil, fmt.Errorf("patternProperties not supported")
	}
	return nil, fmt.Errorf("object not fully supported")
}

func optional(p *relapse.Pattern) *relapse.Pattern {
	return relapse.NewOr(relapse.NewEmpty(), p)
}

func translateNumeric(schema *Schema) (*relapse.Pattern, error) {
	v := Number()
	if schema.Type != nil {
		if len(*schema.Type) > 1 {
			return nil, fmt.Errorf("list of types not supported with numeric constraints %#v", schema)
		}
		if schema.GetType()[0] == TypeInteger {
			v = Integer()
		} else if schema.GetType()[0] != TypeNumber {
			return nil, fmt.Errorf("%v not supported with numeric constraints", schema.GetType()[0])
		}
	}
	list := []funcs.Bool{}
	if schema.MultipleOf != nil {
		mult := MultipleOf(v, funcs.DoubleConst(*schema.MultipleOf))
		list = append(list, mult)
	}
	if schema.Maximum != nil {
		lt := funcs.DoubleLE(v, funcs.DoubleConst(*schema.Maximum))
		if schema.ExclusiveMaximum {
			lt = funcs.DoubleLt(v, funcs.DoubleConst(*schema.Maximum))
		}
		list = append(list, lt)
	}
	if schema.Minimum != nil {
		lt := funcs.DoubleGE(v, funcs.DoubleConst(*schema.Minimum))
		if schema.ExclusiveMinimum {
			lt = funcs.DoubleGt(v, funcs.DoubleConst(*schema.Minimum))
		}
		list = append(list, lt)
	}
	if len(list) == 0 {
		return combinator.Value(funcs.TypeDouble(v)), nil
	}
	return combinator.Value(and(list)), nil
}

func and(list []funcs.Bool) funcs.Bool {
	if len(list) == 0 {
		panic("unreachable")
	}
	if len(list) == 1 {
		return list[0]
	}
	return funcs.And(list[0], and(list[1:]))
}

func translateString(schema *Schema) (*relapse.Pattern, error) {
	v := funcs.StringVar()
	if schema.Type != nil {
		if len(*schema.Type) > 1 {
			return nil, fmt.Errorf("list of types not supported with string constraints %#v", schema)
		}
		if schema.GetType()[0] != TypeString {
			return nil, fmt.Errorf("%v not supported with string constraints", schema.GetType()[0])
		}
	}
	list := []funcs.Bool{}
	if schema.MaxLength != nil {
		list = append(list, MaxLength(v, *schema.MaxLength))
	}
	if schema.MinLength > 0 {
		list = append(list, MinLength(v, schema.MinLength))
	}
	if schema.Pattern != nil {
		list = append(list, funcs.Regex(funcs.StringConst(*schema.Pattern), v))
	}
	if len(list) == 0 {
		return combinator.Value(funcs.TypeString(v)), nil
	}
	return combinator.Value(and(list)), nil
}

func translateArray(schema *Schema) (*relapse.Pattern, error) {
	if schema.Type != nil {
		if len(*schema.Type) > 1 {
			return nil, fmt.Errorf("list of types not supported with array constraints %#v", schema)
		}
		if schema.GetType()[0] != TypeArray {
			return nil, fmt.Errorf("%v not supported with array constraints", schema.GetType()[0])
		}
	}
	if schema.UniqueItems {
		return nil, fmt.Errorf("uniqueItems are not supported")
	}
	if schema.MaxItems != nil {
		return nil, fmt.Errorf("maxItems are not supported")
	}
	if schema.MinItems > 0 {
		return nil, fmt.Errorf("minItems are not supported")
	}
	additionalItems := true
	if schema.AdditionalItems != nil {
		if schema.Items == nil {
			//any
		}
		if schema.AdditionalItems.Bool != nil {
			additionalItems = *schema.AdditionalItems.Bool
		}
		if !additionalItems && (schema.MaxLength != nil || schema.MinLength > 0) {
			return nil, fmt.Errorf("additionalItems: false and (maxItems|minItems) are not supported together")
		}
		return nil, fmt.Errorf("additionalItems are not supported")
	}
	if schema.Items != nil {
		if schema.Items.Object != nil {
			if schema.Items.Object.Type == nil {
				//any
			} else {
				typ := schema.Items.Object.GetType()[0]
				_ = typ
			}
			//TODO this specifies the type of every item in the list
		} else if schema.Items.Array != nil {
			if !additionalItems {
				//TODO this specifies the length of the list as well as each ordered element's type
				//  if no type is set then any type is accepted
				maxLength := len(schema.Items.Array)
				_ = maxLength
			} else {
				//TODO this specifies the types of the first few ordered items in the list
				//  if no type is set then any type is accepted
			}

		}
		return nil, fmt.Errorf("items are not supported")
	}
	return nil, nil
}
