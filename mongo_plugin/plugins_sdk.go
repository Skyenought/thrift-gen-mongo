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
