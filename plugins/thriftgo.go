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
	"go/format"
	"os"
	"path/filepath"

	"github.com/hertz-contrib/thrift-gen-mongo/utils"

	"github.com/cloudwego/thriftgo/plugin"
	"github.com/hertz-contrib/thrift-gen-mongo/args"
	"github.com/hertz-contrib/thrift-gen-mongo/code"
	"github.com/hertz-contrib/thrift-gen-mongo/codegen"
	"github.com/hertz-contrib/thrift-gen-mongo/extract"
	"github.com/hertz-contrib/thrift-gen-mongo/template"
)

func buildResponse(args *args.Arguments,
	structs []*extract.IdlExtractStruct,
	methodRenders [][]*template.MethodRender,
	info *extract.ThriftMeta,
) (result []*plugin.Generated, err error) {
	for index, st := range structs {
		// get base render
		baseRender := getBaseRender(st, args.Version)
		// get fileMongoName and fileIfName
		fileMongoName, fileIfName := extract.GetFileName(st.Name, args.DaoDir)

		if st.Update {
			// build update mongo file
			formattedCode, err := getUpdateMongoCode(methodRenders[index], string(st.UpdateCurdFileContent))
			if err != nil {
				return nil, err
			}
			formattedCode, err = codegen.AddMongoImports(formattedCode)
			if err != nil {
				return nil, err
			}
			formattedCode, err = extract.AddMongoModelImports(formattedCode, info.ImportPaths)
			if err != nil {
				return nil, err
			}
			result = append(result, &plugin.Generated{
				Content: formattedCode,
				Name:    &fileMongoName,
			})

			// build update interface file
			formattedCode, err = getUpdateIfCode(st, baseRender)
			if err != nil {
				return nil, err
			}
			formattedCode, err = codegen.AddMongoImports(formattedCode)
			if err != nil {
				return nil, err
			}
			formattedCode, err = extract.AddMongoModelImports(formattedCode, info.ImportPaths)
			if err != nil {
				return nil, err
			}
			result = append(result, &plugin.Generated{
				Content: formattedCode,
				Name:    &fileIfName,
			})
		} else {
			// build new mongo file
			formattedCode, err := getNewMongoCode(methodRenders[index], st, baseRender)
			if err != nil {
				return nil, err
			}
			formattedCode, err = codegen.AddMongoImports(formattedCode)
			if err != nil {
				return nil, err
			}
			formattedCode, err = extract.AddMongoModelImports(formattedCode, info.ImportPaths)
			if err != nil {
				return nil, err
			}
			result = append(result, &plugin.Generated{
				Content: formattedCode,
				Name:    &fileMongoName,
			})

			// build new interface file
			formattedCode, err = getNewIfCode(st, baseRender)
			if err != nil {
				return nil, err
			}
			formattedCode, err = codegen.AddMongoImports(formattedCode)
			if err != nil {
				return nil, err
			}
			formattedCode, err = extract.AddMongoModelImports(formattedCode, info.ImportPaths)
			if err != nil {
				return nil, err
			}
			result = append(result, &plugin.Generated{
				Content: formattedCode,
				Name:    &fileIfName,
			})
		}
	}

	return
}

func getBaseRender(st *extract.IdlExtractStruct, version string) *template.BaseRender {
	pkgName := extract.GetPkgName(st.Name)
	return &template.BaseRender{
		Version:     version,
		PackageName: pkgName,
		Imports:     codegen.BaseMongoImports,
	}
}

func getUpdateMongoCode(methodRenders []*template.MethodRender, fileContent string) (string, error) {
	tplMongo := &template.Template{
		Renders: []template.Render{},
	}
	for _, methodRender := range methodRenders {
		tplMongo.Renders = append(tplMongo.Renders, methodRender)
	}

	buff, err := tplMongo.Build()
	if err != nil {
		return "", err
	}
	data := fileContent + "\n" + buff.String()
	formattedCode, err := format.Source([]byte(data))
	if err != nil {
		return "", err
	}

	return string(formattedCode), nil
}

func getUpdateIfCode(st *extract.IdlExtractStruct, baseRender *template.BaseRender) (string, error) {
	tplIf := &template.Template{
		Renders: []template.Render{},
	}
	tplIf.Renders = append(tplIf.Renders, baseRender)

	methods := make(code.InterfaceMethods, 0, 10)
	for _, preMethod := range st.PreIfMethods {
		methods = append(methods, code.InterfaceMethod{
			Name:    preMethod.Name,
			Params:  preMethod.Params,
			Returns: preMethod.Returns,
		})
	}
	for _, rawMethod := range st.InterfaceInfo.Methods {
		methods = append(methods, code.InterfaceMethod{
			Name:    rawMethod.Name,
			Params:  rawMethod.Params,
			Returns: rawMethod.Returns,
		})
	}

	ifRender := &template.InterfaceRender{
		Name:    st.Name + "Repository",
		Methods: methods,
	}
	tplIf.Renders = append(tplIf.Renders, ifRender)

	buff, err := tplIf.Build()
	if err != nil {
		return "", err
	}
	formattedCode, err := format.Source(buff.Bytes())
	if err != nil {
		return "", err
	}

	return string(formattedCode), nil
}

func getNewMongoCode(methodRenders []*template.MethodRender, st *extract.IdlExtractStruct, baseRender *template.BaseRender) (string, error) {
	tplMongo := &template.Template{
		Renders: []template.Render{},
	}

	tplMongo.Renders = append(tplMongo.Renders, baseRender)
	tplMongo.Renders = append(tplMongo.Renders, codegen.GetFuncRender(st))
	tplMongo.Renders = append(tplMongo.Renders, codegen.GetStructRender(st))
	for _, methodRender := range methodRenders {
		tplMongo.Renders = append(tplMongo.Renders, methodRender)
	}

	buff, err := tplMongo.Build()
	if err != nil {
		return "", err
	}
	formattedCode, err := format.Source(buff.Bytes())
	if err != nil {
		return "", err
	}

	return string(formattedCode), nil
}

func getNewIfCode(st *extract.IdlExtractStruct, baseRender *template.BaseRender) (string, error) {
	tplIf := &template.Template{
		Renders: []template.Render{},
	}
	tplIf.Renders = append(tplIf.Renders, baseRender)

	methods := make(code.InterfaceMethods, 0, 10)
	for _, rawMethod := range st.InterfaceInfo.Methods {
		methods = append(methods, code.InterfaceMethod{
			Name:    rawMethod.Name,
			Params:  rawMethod.Params,
			Returns: rawMethod.Returns,
		})
	}
	ifRender := &template.InterfaceRender{
		Name:    st.Name + "Repository",
		Methods: methods,
	}
	tplIf.Renders = append(tplIf.Renders, ifRender)

	buff, err := tplIf.Build()
	if err != nil {
		return "", err
	}
	formattedCode, err := format.Source(buff.Bytes())
	if err != nil {
		return "", err
	}

	return string(formattedCode), nil
}

func generateBaseMongoFile(daoDir string, importPaths []string, methodRenders []*template.MethodRender, version string) (err error) {
	st := &extract.IdlExtractStruct{
		Name:          "Base",
		StructFields:  []*extract.StructField{},
		InterfaceInfo: &extract.InterfaceInfo{},
		UpdateInfo:    extract.UpdateInfo{},
	}
	st.InterfaceInfo.Methods = getInterfaceMethods()

	baseRender := getBaseRender(st, version)
	fileMongoName, fileIfName := extract.GetFileName(st.Name, daoDir)
	if isExist, _ := utils.PathExist(filepath.Dir(fileMongoName)); !isExist {
		if err := os.MkdirAll(filepath.Dir(fileMongoName), 0o755); err != nil {
			return err
		}
	}
	if isExist, _ := utils.PathExist(filepath.Dir(fileIfName)); !isExist {
		if err := os.MkdirAll(filepath.Dir(fileIfName), 0o755); err != nil {
			return err
		}
	}

	// build new mongo file
	formattedCode, err := getNewMongoCode(methodRenders, st, baseRender)
	if err != nil {
		return err
	}
	formattedCode, err = codegen.AddMongoImports(formattedCode)
	if err != nil {
		return err
	}
	formattedCode, err = extract.AddMongoModelImports(formattedCode, importPaths)
	if err != nil {
		return err
	}
	formattedCode, err = codegen.AddBaseMGoImports(formattedCode)
	if err != nil {
		return err
	}

	if err = utils.CreateFile(fileMongoName, formattedCode); err != nil {
		return err
	}

	// build new interface file
	formattedCode, err = getNewIfCode(st, baseRender)
	if err != nil {
		return err
	}
	formattedCode, err = codegen.AddMongoImports(formattedCode)
	if err != nil {
		return err
	}
	formattedCode, err = extract.AddMongoModelImports(formattedCode, importPaths)
	if err != nil {
		return err
	}

	if err = utils.CreateFile(fileIfName, formattedCode); err != nil {
		return err
	}

	return
}
