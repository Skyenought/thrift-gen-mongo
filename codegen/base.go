/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package codegen

import (
	"strings"

	"github.com/hertz-contrib/thrift-gen-mongo/code"
	"github.com/hertz-contrib/thrift-gen-mongo/extract"
	"github.com/hertz-contrib/thrift-gen-mongo/template"
)

func GetInterfaceMethods() []*extract.InterfaceMethod {
	createMethod := func(name string, params []code.Param, returnType string) *extract.InterfaceMethod {
		return &extract.InterfaceMethod{
			Name:    name,
			Params:  params,
			Returns: code.Returns{code.IdentType(returnType)},
		}
	}

	return []*extract.InterfaceMethod{
		createMethod("MFindOne", GetMFindOneParams(), "error"),
		createMethod("MFindList", GetMFindListParams(), "error"),
		createMethod("MFindPageList", GetMFindPageListParams(), "error"),
		createMethod("MFindSortPageList", GetMFindSortPageListParams(), "error"),
		createMethod("MInsertOne", GetMInsertOneParams(), "(*mongo.InsertOneResult, error)"),
		createMethod("MUpdateOne", GetMUpdateOneParams(), "(*mongo.UpdateResult, error)"),
		createMethod("MUpdateMany", GetMUpdateOneParams(), "(*mongo.UpdateResult, error)"),
		createMethod("MDeleteOne", GetMDeleteOneParams(), "(*mongo.DeleteResult, error)"),
		createMethod("MBulkInsert", GetMBulkInsertParams(), "(*mongo.BulkWriteResult, error)"),
		createMethod("MBulkUpdate", GetMBulkUpdateParams(), "(*mongo.BulkWriteResult, error)"),
		createMethod("MAggregate", GetMAggregateParams(), "error"),
		createMethod("MCount", GetMCountParams(), "(int64, error)"),
	}
}

func GetBaseRender(st *extract.IDLExtractStruct, version string) *template.BaseRender {
	packageName := st.PkgName
	if strings.Contains(packageName, "/") || strings.Contains(packageName, "\\") {
		split := strings.FieldsFunc(packageName, func(r rune) bool {
			return r == '/' || r == '\\'
		})
		packageName = split[len(split)-1]
	}

	return &template.BaseRender{
		Version:     version,
		PackageName: packageName,
		Imports:     BaseMongoImports,
	}
}
