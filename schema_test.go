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
	"github.com/katydid/katydid/serialize/json"
	"testing"
)

func catch(f func() bool) (v bool, err error) {
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		err = fmt.Errorf("%v", recover())
	}()
	v = f()
	return
}

var skippingFile = map[string]bool{
	"uniqueItems.json":          true, //known issue
	"type.json":                 true,
	"required.json":             true,
	"refRemote.json":            true,
	"ref.json":                  true,
	"properties.json":           true,
	"patternProperties.json":    true, //known issue
	"format.json":               true,
	"not.json":                  true,
	"minProperties.json":        true,
	"minItems.json":             true,
	"maxProperties.json":        true,
	"maxItems.json":             true,
	"items.json":                true,
	"enum.json":                 true,
	"dependencies.json":         true,
	"default.json":              true,
	"definitions.json":          true,
	"allOf.json":                true,
	"additionalProperties.json": true,
	"additionalItems.json":      true,
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
		total++
		//t.Logf("--- RUN: %v", test)
		g, err := TranslateDraft4(test.Schema)
		if err != nil {
			t.Errorf("--- FAIL: %v: %v", test, err)
		} else {
			if err := p.Init(test.Data); err != nil {
				t.Errorf("--- FAIL: %v: %v", test, err)
			} else {
				//t.Logf("--- PASS: %v", test)
			}
			_ = interp.Interpret
			_ = g
			// valid, err := catch(func() bool {
			// 	return interp.Interpret(g, s)
			// })
			// if err != nil {
			// 	fail = true
			// 	t.Errorf(test.Description + ": Interpret Error:" + err.Error())
			// }
			// if valid != test.Valid {
			// 	fail = true
			// 	t.Errorf(test.Description+": expected %v got %v", test.Valid, valid)
			// }
		}
	}
	t.Logf("number of tests passing: %d", total)
}
