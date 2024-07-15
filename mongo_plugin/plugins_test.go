/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mongo_plugin

import (
	"os"
	"testing"

	"github.com/hertz-contrib/thrift-gen-mongo/utils"

	"github.com/hertz-contrib/thrift-gen-mongo/args"
	"github.com/hertz-contrib/thrift-gen-mongo/codegen"
	"github.com/hertz-contrib/thrift-gen-mongo/extract"
	"github.com/hertz-contrib/thrift-gen-mongo/parse"

	"github.com/cloudwego/thriftgo/plugin"
)

func Test_handleRequest(t *testing.T) {
	data, err := os.ReadFile("../plugins/testdata/request_thrift.out")
	if err != nil {
		t.Fatal(err)
	}

	req, err := plugin.UnmarshalRequest(data)
	if err != nil {
		t.Fatal(err)
	}

	a := new(args.Arguments)
	err = a.Unpack(req.PluginParameters)
	if err != nil {
		t.Fatal(err)
	}

	if a.UseGenDir {
		genDir, err := utils.FindGenDir(a.OutDir)
		if err != nil {
			t.Fatal(err)
		}
		a.DaoDir = genDir
		a.ModelDir = genDir
	}

	thriftMeta := &extract.ThriftMeta{
		Req:         req,
		Args:        a,
		ImportPaths: make([]string, 0, 10),
	}

	rawStructs, err := thriftMeta.ParseThriftIDL()
	if err != nil {
		t.Fatal(err)
	}
	operations, err := parse.HandleOperations(rawStructs)
	if err != nil {
		t.Fatal(err)
	}

	methodRenders := codegen.HandleCodegen(operations)
	generated, err := buildResponse(a, rawStructs, methodRenders, thriftMeta)
	if err != nil && len(generated) != 2 {
		t.Fatal(err)
	}
	_, err = generateBaseMongoFile(a.DaoDir, thriftMeta.ImportPaths, codegen.HandleBaseCodegen(), a.Version)
	if err != nil {
		t.Fatal(err)
	}
	if err := handleResponse(&plugin.Response{Contents: generated}); err != nil {
		t.Fatal(err)
	}
}
