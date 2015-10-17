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
		return nil, fmt.Errorf("string not supported")
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
	if len(list) == 0 {
		panic("unreachable")
	}
	if len(list) == 1 {
		return newLeaf(funcs.Sprint(list[0])), nil
	}
	if len(list) == 2 {
		return newLeaf(funcs.Sprint(funcs.And(list[0], list[1]))), nil
	}
	return newLeaf(funcs.Sprint(funcs.And(list[0], funcs.And(list[1], list[2])))), nil
}
