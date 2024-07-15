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
	"go/format"

	"github.com/hertz-contrib/thrift-gen-mongo/code"
	"github.com/hertz-contrib/thrift-gen-mongo/extract"
	"github.com/hertz-contrib/thrift-gen-mongo/template"
)

func GetUpdateMongoCode(methodRenders []*template.MethodRender, fileContent string) (string, error) {
	tplMongo := &template.Template{}
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

func GetUpdateIfCode(st *extract.IDLExtractStruct, baseRender *template.BaseRender) (string, error) {
	tplIf := &template.Template{}
	tplIf.Renders = append(tplIf.Renders, baseRender)

	methods := code.InterfaceMethods{}
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

func GetNewMongoCode(methodRenders []*template.MethodRender, st *extract.IDLExtractStruct, baseRender *template.BaseRender) (string, error) {
	tplMongo := &template.Template{}
	tplMongo.Renders = append(tplMongo.Renders, baseRender, GetFuncRender(st), GetStructRender(st))
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

func GetNewIfCode(st *extract.IDLExtractStruct, baseRender *template.BaseRender) (string, error) {
	tplIf := &template.Template{}
	tplIf.Renders = append(tplIf.Renders, baseRender)

	methods := code.InterfaceMethods{}
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
