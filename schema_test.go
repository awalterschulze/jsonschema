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
	"fmt"
	"github.com/katydid/katydid/relapse/interp"
	"github.com/katydid/katydid/serialize/debug"
	"github.com/katydid/katydid/serialize/json"
	"strings"
	"testing"
)

func catch(f func() bool) (v bool, err error) {
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		r := recover()
		if r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	v = f()
	return
}

var skippingFile = map[string]bool{
	"format.json":               true, //optional
	"bignum.json":               true, //optional
	"zeroTerminatedFloats.json": true, //optional
	"uniqueItems.json":          true, //known issue
	"patternProperties.json":    true, //known issue
	"minProperties.json":        true, //known issue?
	"minItems.json":             true, //known issue?
	"maxProperties.json":        true, //known issue?
	"maxItems.json":             true, //known issue?
	"refRemote.json":            true, //known issue?
	"required.json":             true,
	"ref.json":                  true,
	//"properties.json":           true,
	"not.json":                  true, //requires properties and type object
	"items.json":                true,
	"enum.json":                 true, //requires properties and type object
	"dependencies.json":         true,
	"default.json":              true,
	"definitions.json":          true,
	"allOf.json":                true,
	"additionalProperties.json": true,
	"additionalItems.json":      true,
}

var skippingTest = map[string]bool{
	"type.json:object type matches objects:an array is not an object": true, //known issue
	"type.json:array type matches arrays:an object is not an array":   true, //known issue
}

func TestDraft4(t *testing.T) {
	tests := buildTests(t)
	t.Logf("skipping files: %d", len(skippingFile))
	t.Logf("total number of tests: %d", len(tests))
	total := 0

	p := json.NewJsonParser()
	for _, test := range tests {
		if skippingFile[test.Filename] {
			//t.Logf("--- SKIP: %v", test)
			continue
		}
		if skippingTest[test.String()] {
			//t.Logf("--- SKIP: %v", test)
			continue
		}
		//t.Logf("--- RUN: %v", test)
		schema, err := ParseSchema(test.Schema)
		if err != nil {
			t.Errorf("--- FAIL: %v: Parse error %v", test, err)
		} else {
			g, err := TranslateDraft4(schema)
			if err != nil {
				t.Errorf("--- FAIL: %v: Translate error %v", test, err)
			} else {
				if err := p.Init(test.Data); err != nil {
					t.Errorf("--- FAIL: %v: parser Init error %v", test, err)
				}
				_ = interp.Interpret
				_ = g
				valid, err := catch(func() bool {
					return interp.Interpret(g, p)
				})
				if err != nil {
					t.Errorf("--- FAIL: %v: Interpret error %v", test, err)
				} else if valid != test.Valid {
					t.Errorf("--- FAIL: %v: expected %v got %v", test, test.Valid, valid)
				} else {
					//t.Logf("--- PASS: %v", test)
					total++
				}
			}
		}
	}
	t.Logf("number of tests passing: %d", total)
}

func testDebug(t *testing.T, test Test) {
	jsonp := json.NewJsonParser()
	p := debug.NewLogger(jsonp, debug.NewLineLogger())
	t.Logf("Schema = %v", string(test.Schema))
	schema, err := ParseSchema(test.Schema)
	if err != nil {
		t.Fatalf("Parser error %v", err)
	}
	t.Logf("Parsed Schema %v", schema.JsonString())
	g, err := TranslateDraft4(schema)
	if err != nil {
		t.Fatalf("Translate error %v", err)
	}
	t.Logf("Translated = %v", g)
	t.Logf("Input = %v", string(test.Data))
	if err := jsonp.Init(test.Data); err != nil {
		t.Fatalf("parser Init error %v", err)
	}
	valid, err := catch(func() bool {
		return interp.Interpret(g, p)
	})
	if err != nil {
		t.Fatalf("Interpret error %v", err)
	} else if valid != test.Valid {
		t.Fatalf("expected %v got %v", test.Valid, valid)
	}
}

func TestDebug(t *testing.T) {
	tests := buildTests(t)
	for _, test := range tests {
		if !strings.Contains(test.String(), "properties.json:object properties validation:doesn't invalidate other properties") {
			continue
		}
		testDebug(t, test)
	}
}
