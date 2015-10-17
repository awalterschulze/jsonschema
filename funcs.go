// Copyright 2015 Walter Schulze
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
	"github.com/katydid/katydid/funcs"
)

func MultipleOf(n funcs.Double, d funcs.ConstDouble) funcs.Bool {
	return &multipleOf{n, d, 0}
}

//http://json-schema.org/latest/json-schema-validation.html#anchor14
type multipleOf struct {
	N funcs.Double
	D funcs.ConstDouble
	d float64
}

func (this *multipleOf) Init() error {
	f, err := this.D.Eval()
	if err != nil {
		return err
	}
	this.d = f
	return nil
}

func (this *multipleOf) Eval() (bool, error) {
	n, err := this.N.Eval()
	if err != nil {
		return false, err
	}
	v := n / this.d
	return v == float64(int64(v)) || v == float64(uint64(v)), nil
}

func init() {
	funcs.Register("multipleOf", new(multipleOf))
}

func Integer() funcs.Double {
	return &integer{funcs.UintVar(), funcs.IntVar()}
}

type integer struct {
	U funcs.Uint
	I funcs.Int
}

func (this *integer) Eval() (float64, error) {
	u, err := this.U.Eval()
	if err == nil {
		return float64(u), nil
	}
	i, err := this.I.Eval()
	if err == nil {
		return float64(i), nil
	}
	return 0, err
}

func init() {
	funcs.Register("integer", new(integer))
}

func Number() funcs.Double {
	return &number{Integer(), funcs.DoubleVar()}
}

type number struct {
	I funcs.Double
	D funcs.Double
}

func (this *number) Eval() (float64, error) {
	i, err := this.I.Eval()
	if err == nil {
		return i, nil
	}
	return this.D.Eval()
}

func init() {
	funcs.Register("number", new(number))
}

func MaxLength(v funcs.String, n uint64) funcs.Bool {
	return &maxLength{v, funcs.UintConst(n), 0}
}

type maxLength struct {
	S funcs.String
	N funcs.ConstUint
	n uint64
}

func (this *maxLength) Init() error {
	n, err := this.N.Eval()
	if err != nil {
		return err
	}
	this.n = n
	return nil
}

func (this *maxLength) Eval() (bool, error) {
	s, err := this.S.Eval()
	if err != nil {
		return false, err
	}
	return len(s) <= int(this.n), nil
}

func init() {
	funcs.Register("maxLength", new(maxLength))
}

func MinLength(v funcs.String, n uint64) funcs.Bool {
	return &minLength{v, funcs.UintConst(n), 0}
}

type minLength struct {
	S funcs.String
	N funcs.ConstUint
	n uint64
}

func (this *minLength) Init() error {
	n, err := this.N.Eval()
	if err != nil {
		return err
	}
	this.n = n
	return nil
}

func (this *minLength) Eval() (bool, error) {
	s, err := this.S.Eval()
	if err != nil {
		return false, err
	}
	return len(s) >= int(this.n), nil
}

func init() {
	funcs.Register("minLength", new(minLength))
}
