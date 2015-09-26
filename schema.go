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
	Id               *string
	Schema           *string
	Title            *string
	Description      *string
	Default          interface{}
	MultipleOf       *float64
	Maximum          *float64
	ExclusiveMaximum bool
	Minimum          *float64
	ExclusiveMinimum bool
	MaxLength        *uint64
	MinLength        uint64
	Pattern          *string
	/*
	   "anyOf": [
	       { "type": "boolean" },
	       { "$ref": "#" }
	   ],
	   "default": {}
	*/
	AdditionalItems interface{}
	/*
	   "anyOf": [
	       { "$ref": "#" },
	       { "$ref": "#/definitions/schemaArray" }
	   ],
	   "default": {}
	*/
	Items         interface{}
	MaxItems      *uint64
	MinItems      uint64
	UniqueItems   bool
	MaxProperties *uint64
	MinProperties uint64
	Required      []string
	/*
	   "anyOf": [
	       { "type": "boolean" },
	       { "$ref": "#" }
	   ],
	   "default": {}
	*/
	AdditionalProperties interface{}
	/*
	   "type": "object",
	   "additionalProperties": { "$ref": "#" },
	   "default": {}
	*/
	Definitions interface{}
	/*
	   "type": "object",
	   "additionalProperties": { "$ref": "#" },
	   "default": {}
	*/
	Properties interface{}
	/*
	   "type": "object",
	   "additionalProperties": { "$ref": "#" },
	   "default": {}
	*/
	PatternProperties interface{}
	/*
	   "type": "object",
	   "additionalProperties": {
	       "anyOf": [
	           { "$ref": "#" },
	           { "$ref": "#/definitions/stringArray" }
	       ]
	   }
	*/
	Dependencies interface{}
	/*
	   "type": "array",
	   "minItems": 1,
	   "uniqueItems": true
	*/
	Enum interface{}
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
	Type  interface{}
	AllOf []*Schema
	AnyOf []*Schema
	OneOf []*Schema
	Not   *Schema
}
