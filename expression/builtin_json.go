// Copyright 2013 The ql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSES/QL-LICENSE file.

// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package expression

import (
	"regexp"

	"github.com/juju/errors"
	"github.com/pingcap/tidb/context"
	"github.com/pingcap/tidb/util/types"
)

var (
	_ functionClass = &jsonExtractFunctionClass{}
)

var (
	pathExprBodyRe = regexp.MustCompile("(\\.\\w+|\\[[0-9]+\\])")
)

type jsonExtractFunctionClass struct {
	baseFunctionClass
}

func (c *jsonExtractFunctionClass) getFunction(args []Expression, ctx context.Context) (builtinFunc, error) {
	return &builtinJsonExtractSig{newBaseBuiltinFunc(args, ctx)}, errors.Trace(c.verifyArgs(args))
}

type builtinJsonExtractSig struct {
	baseBuiltinFunc
}

func (b *builtinJsonExtractSig) eval(row []types.Datum) (d types.Datum, err error) {
	args, err := b.evalArgs(row)
	if err != nil {
		return d, errors.Trace(err)
	}

	pathExpr, err := args[1].ToString()
	if err == nil {
		err = validatePathExpr(pathExpr)
	}
	if err != nil {
		return d, err
	}

	var j types.Json
	switch args[0].Kind() {
	case types.KindString, types.KindBytes:
		j = types.JsonFromString(args[0].GetString())
	case types.KindMysqlJson:
		j = args[0].GetMysqlJson()
	default:
		return d, errors.New("invalid argument 0")
	}

	retval, err := doJsonExtract(j, pathExpr)
	if err != nil {
		return d, errors.Trace(err)
	}
	d.SetValue(retval)
	return d, err
}

func doJsonExtract(j types.Json, path string) (d types.Datum, err error) {
	// TODO finish this after real json data struct is ok.
	d.SetNull()
	return d, err
}

func validatePathExpr(pathExpr string) error {
	return nil
}
