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
	"log"
	"os"
	"path/filepath"

	"github.com/cloudwego/thriftgo/plugin"
	"github.com/hertz-contrib/thrift-gen-mongo/args"
	"github.com/hertz-contrib/thrift-gen-mongo/codegen"
	"github.com/hertz-contrib/thrift-gen-mongo/extract"
	"github.com/hertz-contrib/thrift-gen-mongo/template"
	"github.com/hertz-contrib/thrift-gen-mongo/utils"
)

func buildResponse(args *args.Arguments,
	structs []*extract.IDLExtractStruct,
	methodRenders [][]*template.MethodRender,
	info *extract.ThriftMeta,
) (result []*plugin.Generated, err error) {
	addGeneratedFile := func(content, fileName string) error {
		formattedCode, err := codegen.AddMongoImports(content)
		if err != nil {
			return err
		}
		formattedCode, err = extract.AddMongoModelImports(formattedCode, info.ImportPaths)
		if err != nil {
			return err
		}
		result = append(result, &plugin.Generated{
			Content: formattedCode,
			Name:    &fileName,
		})
		return nil
	}

	for index, st := range structs {
		baseRender := codegen.GetBaseRender(st, args.Version)
		fileMongoName, fileIfName := extract.GetFileName(st.PkgName, st.Name, args.DaoDir)

		formattedCode, err := getFileContent(st, methodRenders[index], baseRender)
		if err != nil {
			return nil, err
		}
		if err := addGeneratedFile(formattedCode, fileMongoName); err != nil {
			return nil, err
		}

		formattedCode, err = getInterfaceContent(st, baseRender)
		if err != nil {
			return nil, err
		}
		if err := addGeneratedFile(formattedCode, fileIfName); err != nil {
			return nil, err
		}
	}

	return
}

func getFileContent(st *extract.IDLExtractStruct, methodRenders []*template.MethodRender, baseRender *template.BaseRender) (string, error) {
	var formattedCode string
	var err error
	if st.Update {
		formattedCode, err = codegen.GetUpdateMongoCode(methodRenders, string(st.UpdateCurdFileContent))
	} else {
		formattedCode, err = codegen.GetNewMongoCode(methodRenders, st, baseRender)
	}
	if err != nil {
		return "", err
	}
	return formattedCode, nil
}

func getInterfaceContent(st *extract.IDLExtractStruct, baseRender *template.BaseRender) (string, error) {
	if st.Update {
		return codegen.GetUpdateIfCode(st, baseRender)
	}
	return codegen.GetNewIfCode(st, baseRender)
}

func generateBaseMongoFile(daoDir string, importPaths []string, methodRenders []*template.MethodRender, version string) ([]*plugin.Generated, error) {
	st := newBaseIDLExtractStruct(version)

	fileMongoName, fileIfName := extract.GetFileName(st.PkgName, st.Name, daoDir)
	if err := createDirsIfNotExist(fileMongoName, fileIfName); err != nil {
		return nil, err
	}

	if err := genMongoFile(fileMongoName, methodRenders, st, importPaths, version); err != nil {
		return nil, err
	}

	if err := genInterfaceFile(fileIfName, st, importPaths, version); err != nil {
		return nil, err
	}

	return []*plugin.Generated{
		{
			Content: readContent(fileIfName), // Assuming readContent reads file content into a string.
			Name:    &fileIfName,
		},
	}, nil
}

func newBaseIDLExtractStruct(version string) *extract.IDLExtractStruct {
	return &extract.IDLExtractStruct{
		Name:         "Base",
		PkgName:      "base",
		StructFields: []*extract.StructField{},
		InterfaceInfo: &extract.InterfaceInfo{
			Methods: codegen.GetInterfaceMethods(),
		},
		UpdateInfo: extract.UpdateInfo{},
	}
}

func createDirsIfNotExist(fileMongoName, fileIfName string) error {
	if err := createDirIfNotExist(filepath.Dir(fileMongoName)); err != nil {
		return err
	}
	return createDirIfNotExist(filepath.Dir(fileIfName))
}

func createDirIfNotExist(dir string) error {
	if isExist, _ := utils.PathExist(dir); !isExist {
		return os.MkdirAll(dir, 0o755)
	}
	return nil
}

func genMongoFile(fileMongoName string, methodRenders []*template.MethodRender, st *extract.IDLExtractStruct, importPaths []string, version string) error {
	formattedCode, err := buildMongoFileContent(methodRenders, st, importPaths, version)
	if err != nil {
		return err
	}
	return utils.CreateFile(fileMongoName, formattedCode)
}

func genInterfaceFile(fileIfName string, st *extract.IDLExtractStruct, importPaths []string, version string) error {
	formattedCode, err := buildInterfaceFileContent(st, importPaths, version)
	if err != nil {
		return err
	}
	return utils.CreateFile(fileIfName, formattedCode)
}

func buildMongoFileContent(methodRenders []*template.MethodRender, st *extract.IDLExtractStruct, importPaths []string, version string) (string, error) {
	baseRender := codegen.GetBaseRender(st, version)
	formattedCode, err := codegen.GetNewMongoCode(methodRenders, st, baseRender)
	if err != nil {
		return "", err
	}
	return addCommonImports(formattedCode, importPaths)
}

func buildInterfaceFileContent(st *extract.IDLExtractStruct, importPaths []string, version string) (string, error) {
	baseRender := codegen.GetBaseRender(st, version)
	formattedCode, err := codegen.GetNewIfCode(st, baseRender)
	if err != nil {
		return "", err
	}
	return addCommonImports(formattedCode, importPaths)
}

func addCommonImports(formattedCode string, importPaths []string) (string, error) {
	formattedCode, err := codegen.AddMongoImports(formattedCode)
	if err != nil {
		return "", err
	}
	formattedCode, err = extract.AddMongoModelImports(formattedCode, importPaths)
	if err != nil {
		return "", err
	}
	return codegen.AddBaseMGoImports(formattedCode)
}

func readContent(fileName string) string {
	content, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("failed reading file: %s", err)
	}
	return string(content)
}
