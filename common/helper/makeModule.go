package helper

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"regexp"
	"strings"
)

var moduleContent = `package DebugModule

import (
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/base"
)

type DebugModule struct {
	base.BaseModule
}

func NewDebugModule() *DebugModule {
	this := new(DebugModule)
	return this
}

func (this *DebugModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	this.BaseModule.Init(app, this, config)
	//app.Sub("debug", this)
	return nil
}

func (this *DebugModule) Handle(event interfaces.Msg) error {
	return nil
}

func (this *DebugModule) Watch(index int) {
	return
}
`

var moduleProviderContent = `package DebugModule

import "github.com/gw123/GMQ/core/interfaces"

type DebugModuleProvider struct {
	module interfaces.Module
}

func NewDebugModuleProvider() *DebugModuleProvider {
	this := new(DebugModuleProvider)
	return this
}

func (this *DebugModuleProvider) GetModuleName() string {
	return "DebugModule"
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
`
var reg = regexp.MustCompile("DebugModule")

func MakeModuleContent(moduleName string) string {
	result := reg.ReplaceAllString(moduleContent, moduleName)
	return result
}

func MakeModuleProviderContent(moduleName string) string {
	result := reg.ReplaceAllString(moduleProviderContent, moduleName)
	return result
}

func MakeModule(moduleName string, distDir string) error {
	if distDir == "" {
		distDir = "./modules/"
	}

	if !strings.HasSuffix(distDir, "/") {
		distDir += "/"
	}
	filePath := distDir + moduleName

	fileInfo, err := os.Stat(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil && fileInfo.IsDir() {
		return errors.New("dir is exist")
	}

	fmt.Println("Mkdir :" + filePath)
	err = os.MkdirAll(filePath, 0660)
	if err != nil {
		return err
	}

	f, err := os.Create(filePath + "/" + moduleName + ".go")
	if err != nil {
		return err
	}
	defer f.Close()
	content := MakeModuleContent(moduleName)
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	f2, err := os.Create(filePath + "/" + moduleName + "Provider.go")
	if err != nil {
		return err
	}
	defer f2.Close()

	content2 := MakeModuleProviderContent(moduleName)
	_, err = f2.WriteString(content2)
	if err != nil {
		return err
	}
	return nil
}
