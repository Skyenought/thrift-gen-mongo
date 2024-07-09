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
	"github.com/hertz-contrib/thrift-gen-mongo/code"
	"github.com/hertz-contrib/thrift-gen-mongo/codegen"
	"github.com/hertz-contrib/thrift-gen-mongo/extract"
)

func getInterfaceMethods() []*extract.InterfaceMethod {
	var methods []*extract.InterfaceMethod

	methods = append(methods, &extract.InterfaceMethod{
		Name:   "MFindOne",
		Params: codegen.GetMFindOneParams(),
		Returns: code.Returns{
			code.IdentType("error"),
		},
	}, &extract.InterfaceMethod{
		Name:   "MFindList",
		Params: codegen.GetMFindListParams(),
		Returns: code.Returns{
			code.IdentType("error"),
		},
	}, &extract.InterfaceMethod{
		Name:   "MFindPageList",
		Params: codegen.GetMFindPageListParams(),
		Returns: code.Returns{
			code.IdentType("error"),
		},
	}, &extract.InterfaceMethod{
		Name:   "MFindSortPageList",
		Params: codegen.GetMFindSortPageListParams(),
		Returns: code.Returns{
			code.IdentType("error"),
		},
	}, &extract.InterfaceMethod{
		Name:   "MInsertOne",
		Params: codegen.GetMInsertOneParams(),
		Returns: code.Returns{
			code.IdentType("(*mongo.InsertOneResult, error)"),
		},
	}, &extract.InterfaceMethod{
		Name:   "MUpdateOne",
		Params: codegen.GetMUpdateOneParams(),
		Returns: code.Returns{
			code.IdentType("(*mongo.UpdateResult, error)"),
		},
	}, &extract.InterfaceMethod{
		Name:   "MUpdateMany",
		Params: codegen.GetMUpdateOneParams(),
		Returns: code.Returns{
			code.IdentType("(*mongo.UpdateResult, error)"),
		},
	}, &extract.InterfaceMethod{
		Name:   "MDeleteOne",
		Params: codegen.GetMDeleteOneParams(),
		Returns: code.Returns{
			code.IdentType("(*mongo.DeleteResult, error)"),
		},
	}, &extract.InterfaceMethod{
		Name:   "MBulkInsert",
		Params: codegen.GetMBulkInsertParams(),
		Returns: code.Returns{
			code.IdentType("(*mongo.BulkWriteResult, error)"),
		},
	}, &extract.InterfaceMethod{
		Name:   "MBulkUpdate",
		Params: codegen.GetMBulkUpdateParams(),
		Returns: code.Returns{
			code.IdentType("(*mongo.BulkWriteResult, error)"),
		},
	}, &extract.InterfaceMethod{
		Name:   "MAggregate",
		Params: codegen.GetMAggregateParams(),
		Returns: code.Returns{
			code.IdentType("error"),
		},
	}, &extract.InterfaceMethod{
		Name:   "MCount",
		Params: codegen.GetMCountParams(),
		Returns: code.Returns{
			code.IdentType("(int64, error)"),
		},
	})

	return methods
}
