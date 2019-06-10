package debugModule

import "github.com/gw123/GMQ/core/interfaces"

type DebugModuleProvider struct {
	module interfaces.Module
}

func NewDebugModuleProvider() *DebugModuleProvider {
	this := new(DebugModuleProvider)
	return this
}

func (this *DebugModuleProvider) GetModuleName() string {
	return "Debug"
}

func (this *DebugModuleProvider) Register() {

}

func (this *DebugModuleProvider) GetModule() interfaces.Module {
	if this.module != nil {
		return this.module
	}
	this.module = NewDebugModule()
	return this.module
}

func (this *DebugModuleProvider) GetNewModule() interfaces.Module {
	this.module = NewDebugModule()
	return this.module
}

