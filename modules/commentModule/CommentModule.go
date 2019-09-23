package commentModule

import (
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/base"
)

type CommentModule struct {
	base.BaseModule
	bindAddr string
	server   *CommentServer
}

func NewCommentModule() *CommentModule {
	this := new(CommentModule)
	return this
}

func (this *CommentModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	this.BaseModule.Init(app, this, config)
	this.bindAddr = config.GetItem("bindAddr")
	this.server = NewCommentServer(this, this.bindAddr)
	return nil
}

func (this *CommentModule) BeforeStart() error {
	return this.server.Start()
}

func (this *CommentModule) UnInit() error {
	this.BaseModule.UnInit()
	this.server.Stop()
	return nil
}

func (this *CommentModule) Handle(event interfaces.Event) error {
	//
	return nil
}

func (this *CommentModule) Watch(index int) {

}
