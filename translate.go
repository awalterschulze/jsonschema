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
	"github.com/katydid/katydid/expr/ast"
	exprparser "github.com/katydid/katydid/expr/parser"
	"github.com/katydid/katydid/funcs"
	"github.com/katydid/katydid/relapse/ast"
)

//TODO
func TranslateDraft4(jsonSchema []byte) (*relapse.Grammar, error) {
	schema := &Schema{}
	if err := json.Unmarshal(jsonSchema, schema); err != nil {
		panic(err)
	}
	refs, err := translate(schema)
	if err != nil {
		return nil, err
	}
	return relapse.NewGrammar(refs), nil
}

func translate(schema *Schema) (relapse.RefLookup, error) {
	if schema.Id != nil {
		return nil, fmt.Errorf("id not supported")
	}
	if schema.Default != nil {
		return nil, fmt.Errorf("default not supported")
	}
	if schema.HasNumericConstraints() {
		p, err := translateNumeric(schema)
		return map[string]*relapse.Pattern{"main": p}, err
	}
	if schema.HasStringConstraints() {
		p, err := translateString(schema)
		return map[string]*relapse.Pattern{"main": p}, err
	}
	if schema.HasArrayConstraints() {
		return nil, fmt.Errorf("array not supported")
	}
	if schema.HasObjectConstraints() {
		return nil, fmt.Errorf("object not supported")
	}
	if schema.HasInstanceConstraints() {
		return nil, fmt.Errorf("instance not supported")
	}

	if schema.Ref != nil {
		return nil, fmt.Errorf("ref not supported")
	}
	if schema.Format != nil {
		return nil, fmt.Errorf("format not supported")
	}
	return map[string]*relapse.Pattern{"main": relapse.NewEmptySet()}, nil
}

//TODO rather call a function available in the katydid relapse library
func newLeaf(exprStr string) *relapse.Pattern {
	e, err := exprparser.NewParser().ParseExpr(exprStr)
	if err != nil {
		panic(err)
	}
	return &relapse.Pattern{LeafNode: &relapse.LeafNode{
		RightArrow: &expr.Keyword{Value: "->"},
		Expr:       e,
	}}
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
	return newLeaf(funcs.Sprint(and(list))), nil
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
	return newLeaf(funcs.Sprint(and(list))), nil
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
