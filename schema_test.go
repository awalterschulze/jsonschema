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

func TestDraft4(t *testing.T) {
	tests := buildTests(t)
	t.Logf("number of tests %d", len(tests))
	s := json.NewJsonParser()
	for _, test := range tests {
		t.Logf(test.Description)
		if err := s.Init(test.Data); err != nil {
			panic(err)
		}
		g, err := TranslateDraft4(test.Schema)
		if err != nil {
			t.Errorf(test.Description + ": failed to convert:" + err.Error())
		}
		valid, err := catch(func() bool {
			return interp.Interpret(g, s)
		})
		if err != nil {
			t.Errorf(test.Description + ": Interpret Error:" + err.Error())
		}
		if valid != test.Valid {
			t.Errorf(test.Description+": expected %v got %v", test.Valid, valid)
		}
	}
}
