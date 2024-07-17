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
	"strings"

	"github.com/cloudwego/thriftgo/plugin"
)

type MongoPlugin struct {
	PluginParameters []string
}

func NewMongoPlugin(params string) *MongoPlugin {
	plugin := &MongoPlugin{}
	if params != "" {
		plugin.PluginParameters = strings.Split(params, ",")
	}
	return plugin
}

func (k *MongoPlugin) Invoke(req *plugin.Request) (res *plugin.Response) {
	response, _ := HandleRequest(req)
	return response
}

func (k *MongoPlugin) GetName() string {
	return "Mongo"
}

func (k *MongoPlugin) GetPluginParameters() []string {
	return k.PluginParameters
}
