package serverNodeModule

import (
	"github.com/gw123/GMQ/core/interfaces"
)

type ServerNodeModuleProvider struct {
	module interfaces.Module
}

func NewServerNodeModuleProvider() *ServerNodeModuleProvider {
	this := new(ServerNodeModuleProvider)
	return this
}

func (this *ServerNodeModuleProvider) GetModuleName() string {
	return "serverNode"
}

func (this *ServerNodeModuleProvider) Register() {

}

func (this *ServerNodeModuleProvider) GetModule() interfaces.Module {
	if this.module != nil {
		return this.module
	}
	this.module = NewServerNodeModule()
	return this.module
}

func (this *ServerNodeModuleProvider) GetNewModule() interfaces.Module {
	this.module = NewServerNodeModule()
	return this.module
}

