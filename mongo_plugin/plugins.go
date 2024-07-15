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
	"io"
	"log"
	"os"

	"github.com/hertz-contrib/thrift-gen-mongo/args"
	"github.com/hertz-contrib/thrift-gen-mongo/codegen"
	"github.com/hertz-contrib/thrift-gen-mongo/extract"
	"github.com/hertz-contrib/thrift-gen-mongo/parse"
	"github.com/hertz-contrib/thrift-gen-mongo/utils"

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

	if _, err := HandleRequest(req); err != nil {
		println("Failed to handle request:", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
	return 0
}

func HandleRequest(req *plugin.Request) (*plugin.Response, error) {
	a := new(args.Arguments)

	if err := a.Unpack(req.PluginParameters); err != nil {
		return returnError("[Error]: unpack args failed", err)
	}

	if a.UseGenDir {
		genDir, err := utils.FindGenDir(a.OutDir)
		if err != nil {
			return returnError("[Error]: find gen dir failed", err)
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
		return returnError("[Error]: parse thrift idl failed", err)
	}

	operations, err := parse.HandleOperations(rawStructs)
	if err != nil {
		return returnError("[Error]: handle operations failed", err)
	}

	methodRenders := codegen.HandleCodegen(operations)

	generated, err := buildResponse(a, rawStructs, methodRenders, thriftMeta)
	if err != nil {
		return returnError("[Error]: build response failed", err)
	}

	if a.GenBase {
		generateds, err := generateBaseMongoFile(a.DaoDir, thriftMeta.ImportPaths, codegen.HandleBaseCodegen(), a.Version)
		if err != nil {
			return returnError("[Error]: generate base mongo file failed", err)
		}
		generated = append(generated, generateds...)
	}

	response := &plugin.Response{
		Contents: generated,
	}

	if err := handleResponse(response); err != nil {
		return returnError("[Error]: handle request failed", err)
	}

	return response, nil
}

func errString(err error) *string {
	if err == nil {
		return nil
	}
	es := err.Error()
	return &es
}

func handleResponse(res *plugin.Response) error {
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

func returnError(msg string, err error) (*plugin.Response, error) {
	log.Printf("%s: %s", msg, err.Error())
	return &plugin.Response{
		Error: errString(err),
	}, err
}
