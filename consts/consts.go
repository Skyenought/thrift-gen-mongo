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

package consts

import "runtime"

const (
	Slash      = "/"
	BackSlash  = "\\"
	BlackSpace = " "
	Comma      = ";"
	Tilde      = "~"
	LineBreak  = "\n"
)

// File Name
const (
	KitexExtensionYaml = "extensions.yaml"
	LayoutFile         = "layout.yaml"
	PackageLayoutFile  = "package.yaml"
	SuffixGit          = ".git"
	DefaultDbOutFile   = "gen.go"
	Main               = "main.go"
	GoMod              = "go.mod"
	HzFile             = ".hz"
)

const (
	IdlThrift = "thrift"
	IdlProto  = "proto"
)

// SysType is the running program's operating system type
const SysType = runtime.GOOS

const WindowsOS = "windows"

const (
	Go     = "go"
	GOPATH = "GOPATH"
	Env    = "env"
	Mod    = "mod"
	Init   = "init"

	OutDir   = "out_dir"
	Verbose  = "verbose"
	Template = "template"
	Branch   = "branch"
	Name     = "name"

	ModelDir = "model_dir"
	DaoDir   = "dao_dir"

	Service         = "service"
	ServerName      = "server_name"
	ServiceType     = "type"
	Module          = "module"
	IDLPath         = "idl"
	Registry        = "registry"
	Pass            = "pass"
	ProtoSearchPath = "proto_search_path"
	ThriftGo        = "thriftgo"
	Protoc          = "protoc"
	GenBase         = "gen_base"

	ProjectPath   = "project_path"
	HertzRepoUrl  = "hertz_repo_url"
	DSN           = "dsn"
	DBType        = "db_type"
	Tables        = "tables"
	ExcludeTables = "exclude_tables"
	OnlyModel     = "only_model"
	OutFile       = "out_file"
	UnitTest      = "unittest"
	ModelPkgName  = "model_pkg"
	Nullable      = "nullable"
	Signable      = "signable"
	IndexTag      = "index_tag"
	TypeTag       = "type_tag"
	HexTag        = "hex"
	SQLDir        = "sql_dir"
)
