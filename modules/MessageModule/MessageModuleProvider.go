package MessageModule

import "github.com/gw123/GMQ/core/interfaces"

type MessageModuleProvider struct {
	module interfaces.Module
}

func NewMessageModuleProvider() *MessageModuleProvider {
	this := new(MessageModuleProvider)
	return this
}

func (this *MessageModuleProvider) GetModuleName() string {
	return "MessageModule"
}

func (this *MessageModuleProvider) Register() {
}

func (this *MessageModuleProvider) GetModule() interfaces.Module {
	if this.module != nil {
		return this.module
	}
	this.module = NewMessageModule()
	return this.module
}

func (this *MessageModuleProvider) GetNewModule() interfaces.Module {
	this.module = NewMessageModule()
	return this.module
}
