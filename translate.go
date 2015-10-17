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
		return nil, fmt.Errorf("numeric not supported")
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
