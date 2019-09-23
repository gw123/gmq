package commentModule

import "github.com/gw123/GMQ/core/interfaces"

type CommentModuleProvider struct {
	module interfaces.Module
}

func NewDebugModuleProvider() *CommentModuleProvider {
	this := new(CommentModuleProvider)
	return this
}

func (this *CommentModuleProvider) GetModuleName() string {
	return "Comment"
}

func (this *CommentModuleProvider) Register() {

}

func (this *CommentModuleProvider) GetModule() interfaces.Module {
	if this.module != nil {
		return this.module
	}
	this.module = NewCommentModule()
	return this.module
}

func (this *CommentModuleProvider) GetNewModule() interfaces.Module {
	this.module = NewCommentModule()
	return this.module
}

