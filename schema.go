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
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func ParseSchema(jsonSchema []byte) (*Schema, error) {
	schema := &Schema{}
	if err := json.Unmarshal(jsonSchema, schema); err != nil {
		return nil, err
	}
	return schema, nil
}

func (this *Schema) JsonString() string {
	data, err := json.Marshal(this)
	if err != nil {
		panic(err)
	}
	return string(data)
}

/*
{
    "id": "http://json-schema.org/draft-04/schema#",
    "$schema": "http://json-schema.org/draft-04/schema#",
    "description": "Core schema meta-schema",
    "definitions": {
        "schemaArray": {
            "type": "array",
            "minItems": 1,
            "items": { "$ref": "#" }
        },
        "positiveInteger": {
            "type": "integer",
            "minimum": 0
        },
        "positiveIntegerDefault0": {
            "allOf": [ { "$ref": "#/definitions/positiveInteger" }, { "default": 0 } ]
        },
        "simpleTypes": {
            "enum": [ "array", "boolean", "integer", "null", "number", "object", "string" ]
        },
        "stringArray": {
            "type": "array",
            "items": { "type": "string" },
            "minItems": 1,
            "uniqueItems": true
        }
    },
    "type": "object",
    "properties": {
        "id": {
            "type": "string",
            "format": "uri"
        },
        "$schema": {
            "type": "string",
            "format": "uri"
        },
        "title": {
            "type": "string"
        },
        "description": {
            "type": "string"
        },
        "default": {},
        "multipleOf": {
            "type": "number",
            "minimum": 0,
            "exclusiveMinimum": true
        },
        "maximum": {
            "type": "number"
        },
        "exclusiveMaximum": {
            "type": "boolean",
            "default": false
        },
        "minimum": {
            "type": "number"
        },
        "exclusiveMinimum": {
            "type": "boolean",
            "default": false
        },
        "maxLength": { "$ref": "#/definitions/positiveInteger" },
        "minLength": { "$ref": "#/definitions/positiveIntegerDefault0" },
        "pattern": {
            "type": "string",
            "format": "regex"
        },
        "additionalItems": {
            "anyOf": [
                { "type": "boolean" },
                { "$ref": "#" }
            ],
            "default": {}
        },
        "items": {
            "anyOf": [
                { "$ref": "#" },
                { "$ref": "#/definitions/schemaArray" }
            ],
            "default": {}
        },
        "maxItems": { "$ref": "#/definitions/positiveInteger" },
        "minItems": { "$ref": "#/definitions/positiveIntegerDefault0" },
        "uniqueItems": {
            "type": "boolean",
            "default": false
        },
        "maxProperties": { "$ref": "#/definitions/positiveInteger" },
        "minProperties": { "$ref": "#/definitions/positiveIntegerDefault0" },
        "required": { "$ref": "#/definitions/stringArray" },
        "additionalProperties": {
            "anyOf": [
                { "type": "boolean" },
                { "$ref": "#" }
            ],
            "default": {}
        },
        "definitions": {
            "type": "object",
            "additionalProperties": { "$ref": "#" },
            "default": {}
        },
        "properties": {
            "type": "object",
            "additionalProperties": { "$ref": "#" },
            "default": {}
        },
        "patternProperties": {
            "type": "object",
            "additionalProperties": { "$ref": "#" },
            "default": {}
        },
        "dependencies": {
            "type": "object",
            "additionalProperties": {
                "anyOf": [
                    { "$ref": "#" },
                    { "$ref": "#/definitions/stringArray" }
                ]
            }
        },
        "enum": {
            "type": "array",
            "minItems": 1,
            "uniqueItems": true
        },
        "type": {
            "anyOf": [
                { "$ref": "#/definitions/simpleTypes" },
                {
                    "type": "array",
                    "items": { "$ref": "#/definitions/simpleTypes" },
                    "minItems": 1,
                    "uniqueItems": true
                }
            ]
        },
        "allOf": { "$ref": "#/definitions/schemaArray" },
        "anyOf": { "$ref": "#/definitions/schemaArray" },
        "oneOf": { "$ref": "#/definitions/schemaArray" },
        "not": { "$ref": "#" }
    },
    "dependencies": {
        "exclusiveMaximum": [ "maximum" ],
        "exclusiveMinimum": [ "minimum" ]
    },
    "default": {}
}
*/
type Schema struct {
	Id          string      `json:"id,omitempty"`
	Schema      string      `json:"$schema,omitempty"`
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Numeric
	String
	Array
	Object
	Instance

	Ref    string `json:"$ref,omitempty"`
	Format string `json:"format,omitempty"`
}

//http://json-schema.org/latest/json-schema-validation.html#anchor13
type Numeric struct {
	MultipleOf       *float64 `json:"multipleOf,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	ExclusiveMaximum bool     `json:"exclusiveMaximum,omitempty"`
	Minimum          *float64 `json:"minimum,omitempty"`
	ExclusiveMinimum bool     `json:"exclusiveMinimum,omitempty"`
}

func (this Numeric) HasNumericConstraints() bool {
	return this.MultipleOf != nil || this.Maximum != nil || this.Minimum != nil
}

//http://json-schema.org/latest/json-schema-validation.html#anchor25
type String struct {
	MaxLength *uint64 `json:"maxLength,omitempty"`
	MinLength uint64  `json:"minLength,omitempty"`
	Pattern   *string `json:"pattern,omitempty"`
}

func (this String) HasStringConstraints() bool {
	return this.MaxLength != nil || this.MinLength > 0 || this.Pattern != nil
}

//http://json-schema.org/latest/json-schema-validation.html#anchor36
type Array struct {
	AdditionalItems *Additional `json:"additionalItems,omitempty"`
	Items           *Items      `json:"items,omitempty"`
	MaxItems        *uint64     `json:"maxItems,omitempty"`
	MinItems        uint64      `json:"minItems,omitempty"`
	UniqueItems     bool        `json:"uniqueItems,omitempty"`
}

func (this Array) HasArrayConstraints() bool {
	return this.AdditionalItems != nil || this.Items != nil ||
		this.MaxItems != nil || this.MinItems > 0 || this.UniqueItems
}

//http://json-schema.org/latest/json-schema-validation.html#anchor53
type Object struct {
	MaxProperties        *uint64     `json:"maxProperties,omitempty"`
	MinProperties        uint64      `json:"minProperties,omitempty"`
	Required             []string    `json:"required,omitempty"`
	AdditionalProperties *Additional `json:"additionalProperties,omitempty"`
	/*
	   "type": "object",
	   "additionalProperties": { "$ref": "#" },
	   "default": {}
	*/
	//http://json-schema.org/latest/json-schema-validation.html#anchor64
	//  The value of "properties" MUST be an object. Each value of this object MUST be an object, and each object MUST be a valid JSON Schema.
	Properties map[string]*Schema `json:"properties,omitempty"`
	/*
	   "type": "object",
	   "additionalProperties": { "$ref": "#" },
	   "default": {}
	*/
	//http://json-schema.org/latest/json-schema-validation.html#anchor64
	//  The value of "patternProperties" MUST be an object. Each property name of this object SHOULD be a valid regular expression, according to the ECMA 262 regular expression dialect. Each property value of this object MUST be an object, and each object MUST be a valid JSON Schema.
	PatternProperties map[string]*Schema `json:"patternProperties,omitempty"`
	Dependencies      *Dependencies      `json:"dependencies,omitempty"`
}

func (this Object) HasObjectConstraints() bool {
	return this.MaxProperties != nil || this.MinProperties > 0 ||
		this.Required != nil || this.AdditionalProperties != nil ||
		this.Properties != nil || this.PatternProperties != nil ||
		this.Dependencies != nil
}

//http://json-schema.org/latest/json-schema-validation.html#anchor75
type Instance struct {
	/*
	   "type": "object",
	   "additionalProperties": { "$ref": "#" },
	   "default": {}
	*/
	//http://json-schema.org/latest/json-schema-validation.html#anchor94
	//  This keyword's value MUST be an object. Each member value of this object MUST be a valid JSON Schema.
	Definitions map[string]*Schema `json:"definitions,omitempty"`
	/*
	   "type": "array",
	   "minItems": 1,
	   "uniqueItems": true
	*/
	Enum  []interface{} `json:"enum,omitempty"`
	Type  *Type         `json:"type,omitempty"`
	AllOf []*Schema     `json:"allOf,omitempty"`
	AnyOf []*Schema     `json:"anyOf,omitempty"`
	OneOf []*Schema     `json:"oneOf,omitempty"`
	Not   *Schema       `json:"not,omitempty"`
}

func (this Instance) HasInstanceConstraints() bool {
	return this.Definitions != nil || this.Enum != nil ||
		this.Type != nil || this.AllOf != nil || this.AnyOf != nil ||
		this.OneOf != nil || this.Not != nil
}

func (this Instance) GetType() []SimpleType {
	return *this.Type
}

/*
   "anyOf": [
       { "type": "boolean" },
       { "$ref": "#" }
   ],
   "default": {}
*/
//http://json-schema.org/latest/json-schema-validation.html#anchor37
//  The value of "additionalItems" MUST be either a boolean or an object. If it is an object, this object MUST be a valid JSON Schema.
//http://json-schema.org/latest/json-schema-validation.html#anchor49
//  The value of "additionalProperties" MUST be a boolean or an object. If it is an object, it MUST also be a valid JSON Schema.
type Additional struct {
	Bool *bool
	//Typically only the type field of the jsonschema is set.
	Type SimpleType
}

type aSchema struct {
	Type *Type `json:"type"`
}

func (this *Additional) UnmarshalJSON(buf []byte) error {
	var b bool
	dec := json.NewDecoder(bytes.NewBuffer(buf))
	if err := dec.Decode(&b); err == nil {
		*this = Additional{Bool: &b}
		return nil
	}
	s := &aSchema{}
	if err := json.Unmarshal(buf, s); err != nil {
		log.Printf("%v", err)
		return err
	}
	if s.Type == nil {
		return fmt.Errorf("the additional(Items|Properties) field is empty")
	}
	typ := *s.Type
	if len(typ) > 1 {
		return fmt.Errorf("the additional(Items|Properties) field's type field has more than one element")
	}
	if len(typ) == 0 {
		panic(fmt.Errorf("%#v buf = %s", s.Type, string(buf)))
	}
	*this = Additional{Type: typ[0]}
	return nil
}

/*
   "anyOf": [
       { "$ref": "#" },
       { "$ref": "#/definitions/schemaArray" }
   ],
   "default": {}
*/
//http://json-schema.org/latest/json-schema-validation.html#anchor37
//  The value of "items" MUST be either an object or an array. If it is an object, this object MUST be a valid JSON Schema. If it is an array, items of this array MUST be objects, and each of these objects MUST be a valid JSON Schema.
type Items struct {
	Object *Schema
	Array  []*Schema
}

func (this *Items) UnmarshalJSON(buf []byte) error {
	var s *Schema
	if err := json.Unmarshal(buf, &s); err == nil {
		*this = Items{Object: s}
		return nil
	}
	schemas := []*Schema{}
	if err := json.Unmarshal(buf, &schemas); err != nil {
		log.Printf("%v input %s", err, string(buf))
		return err
	}
	*this = Items{Array: schemas}
	return nil
}

/*
   "type": "object",
   "additionalProperties": {
       "anyOf": [
           { "$ref": "#" },
           { "$ref": "#/definitions/stringArray" }
       ]
   }
*/
//http://json-schema.org/latest/json-schema-validation.html#anchor70
//  This keyword's value MUST be an object. Each value of this object MUST be either an object or an array.
//  If the value is an object, it MUST be a valid JSON Schema. This is called a schema dependency.
//  If the value is an array, it MUST have at least one element. Each element MUST be a string, and elements in the array MUST be unique. This is called a property dependency.
//http://spacetelescope.github.io/understanding-json-schema/reference/object.html#dependencies
type Dependencies map[string]*Dependency

type Dependency struct {
	Schema           *Schema
	RequiredProperty []string
}

func (this *Dependency) UnmarshalJSON(buf []byte) error {
	var s *Schema
	if err := json.Unmarshal(buf, &s); err == nil {
		*this = Dependency{Schema: s}
		return nil
	}
	var ss []string
	dec := json.NewDecoder(bytes.NewBuffer(buf))
	if err := dec.Decode(&ss); err != nil {
		log.Printf("%v input %s", err, string(buf))
		return err
	}
	*this = Dependency{RequiredProperty: ss}
	checkUnique := make(map[string]struct{})
	for _, s := range this.RequiredProperty {
		if _, ok := checkUnique[s]; ok {
			err := fmt.Errorf("duplicate found in property dependency list %s", s)
			log.Printf("%v", err)
			return err
		}
		checkUnique[s] = struct{}{}
	}
	return nil
}

/*
"anyOf": [
    { "$ref": "#/definitions/simpleTypes" },
    {
        "type": "array",
        "items": { "$ref": "#/definitions/simpleTypes" },
        "minItems": 1,
        "uniqueItems": true
    }
]
*/
type Type []SimpleType

func (this *Type) HasArray() bool {
	return this.has(TypeArray)
}

func (this *Type) HasNumeric() bool {
	return this.has(TypeInteger) || this.has(TypeNumber)
}

func (this *Type) HasString() bool {
	return this.has(TypeString)
}

func (this *Type) has(s SimpleType) bool {
	if this == nil {
		return false
	}
	for _, t := range *this {
		if t == s {
			return true
		}
	}
	return false
}

func (this *Type) Single() bool {
	if this == nil {
		return false
	}
	return len(*this) == 1
}

func (this *Type) UnmarshalJSON(buf []byte) error {
	t := []SimpleType{}
	decs := json.NewDecoder(bytes.NewBuffer(buf))
	var s string
	if err := decs.Decode(&s); err == nil {
		simpleType, err := newSimpleType(s)
		if err != nil {
			log.Printf("%v", err)
			return err
		}
		t = append(t, simpleType)
		*this = t
		return nil
	} else {
		//log.Printf("type decode err = %v input = %s", err, string(buf))
	}
	var ss []string
	decss := json.NewDecoder(bytes.NewBuffer(buf))
	if err := decss.Decode(&ss); err != nil {
		log.Printf("%v", err)
		return err
	}
	simpleTypes := make(map[string]struct{})
	for _, s := range ss {
		simpleType, err := newSimpleType(s)
		if err != nil {
			log.Printf("%v", err)
			return err
		}
		if _, ok := simpleTypes[s]; ok {
			err := fmt.Errorf("type alternatives are not unique, duplicate %s found", simpleType)
			log.Printf("%v", err)
			return err
		}
		simpleTypes[s] = struct{}{}
		t = append(t, simpleType)
	}
	*this = t
	return nil
}

type SimpleType string

const (
	TypeUnknown = SimpleType("unknown")
	TypeArray   = SimpleType("array")
	TypeBoolean = SimpleType("boolean")
	TypeInteger = SimpleType("integer")
	TypeNull    = SimpleType("null")
	TypeNumber  = SimpleType("number")
	TypeObject  = SimpleType("object")
	TypeString  = SimpleType("string")
)

func newSimpleType(s string) (SimpleType, error) {
	switch s {
	case "array":
		return TypeArray, nil
	case "boolean":
		return TypeBoolean, nil
	case "integer":
		return TypeInteger, nil
	case "null":
		return TypeNull, nil
	case "number":
		return TypeNumber, nil
	case "object":
		return TypeObject, nil
	case "string":
		return TypeString, nil
	}
	err := fmt.Errorf("unknown simpletype %s", s)
	log.Printf("%v", err)
	return TypeUnknown, err
}
