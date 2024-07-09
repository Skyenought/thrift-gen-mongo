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

package plugins

import (
	"io"
	"log"
	"os"

	"github.com/hertz-contrib/thrift-gen-mongo/args"
	"github.com/hertz-contrib/thrift-gen-mongo/codegen"
	"github.com/hertz-contrib/thrift-gen-mongo/extract"
	"github.com/hertz-contrib/thrift-gen-mongo/parse"

	"github.com/cloudwego/thriftgo/plugin"
)

func Run() int {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		println("Failed to get input:", err.Error())
		os.Exit(1)
	}

	req, err := plugin.UnmarshalRequest(data)
	if err != nil {
		println("Failed to unmarshal request:", err.Error())
		os.Exit(1)
	}

	if err := HandleRequest(req); err != nil {
		println("Failed to handle request:", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
	return 0
}

func HandleRequest(req *plugin.Request) error {
	a := new(args.Arguments)
	if err := a.Unpack(req.PluginParameters); err != nil {
		log.Printf("[Error]: unpack args failed: %s", err.Error())
		return err
	}

	thriftMeta := &extract.ThriftMeta{
		Req:         req,
		Args:        a,
		ImportPaths: make([]string, 0, 10),
	}

	rawStructs, err := thriftMeta.ParseThriftIdl()
	if err != nil {
		log.Printf("[Error]: parse thrift idl failed: %s", err.Error())
		return err
	}

	operations, err := parse.HandleOperations(rawStructs)
	if err != nil {
		return err
	}

	methodRenders := codegen.HandleCodegen(operations)
	generated, err := buildResponse(a, rawStructs, methodRenders, thriftMeta)
	if err != nil {
		return err
	}

	res := &plugin.Response{
		Contents: generated,
	}

	if err := handleRequest(res); err != nil {
		return err
	}

	if a.GenBase {
		if err = generateBaseMongoFile(a.DaoDir, thriftMeta.ImportPaths, codegen.HandleBaseCodegen(), a.Version); err != nil {
			return err
		}
	}

	return nil
}

func handleRequest(res *plugin.Response) error {
	data, err := plugin.MarshalResponse(res)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func errString(err error) *string {
	es := err.Error()
	return &es
}
