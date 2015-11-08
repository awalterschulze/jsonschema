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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const testPath = "./JSON-Schema-Test-Suite/tests/draft4/"

func getFileNames() []string {
	files := []string{}
	filepath.Walk(testPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files
}

type SchemaTest struct {
	Description string
	Schema      interface{}
	Tests       []*SchemaTesty
}

type SchemaTesty struct {
	Description string
	Data        interface{}
	Valid       bool
}

func buildTests(t *testing.T) []Test {
	tests := []Test{}
	filenames := getFileNames()
	t.Logf("number of test files: %d", len(filenames))
	for _, filename := range filenames {
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		var schemaTests []*SchemaTest
		if err := json.Unmarshal(content, &schemaTests); err != nil {
			panic(filename + ":" + err.Error())
		}
		for i := range schemaTests {
			schemaStr, err := json.Marshal(schemaTests[i].Schema)
			if err != nil {
				panic(err)
			}
			schemaDesc := schemaTests[i].Description
			for j := range schemaTests[i].Tests {
				dataStr, err := json.Marshal(schemaTests[i].Tests[j].Data)
				if err != nil {
					panic(err)
				}
				tests = append(tests, Test{
					Filename:    filepath.Base(filename),
					Description: schemaDesc + ":" + schemaTests[i].Tests[j].Description,
					Schema:      schemaStr,
					Data:        dataStr,
					Valid:       schemaTests[i].Tests[j].Valid,
				})
			}
		}
	}
	return tests
}

type Test struct {
	Filename    string
	Description string
	Schema      []byte
	Data        []byte
	Valid       bool
}

func (this Test) String() string {
	return this.Filename + ":" + this.Description
}
