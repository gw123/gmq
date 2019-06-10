package webModule

import (
	"github.com/gw123/GMQ/core/interfaces"
)

type WebModuleProvider struct {
	module interfaces.Module
}

func NewWebModuleProvider() *WebModuleProvider {
	this := new(WebModuleProvider)
	return this
}

func (this *WebModuleProvider) GetModuleName() string {
	return "Web"
}

func (this *WebModuleProvider) Register() {

}

func (this *WebModuleProvider) GetModule() interfaces.Module {
	if this.module != nil {
		return this.module
	}
	this.module = NewWebModule()
	return this.module
}

func (this *WebModuleProvider) GetNewModule() interfaces.Module {
	this.module = NewWebModule()
	return this.module
}

