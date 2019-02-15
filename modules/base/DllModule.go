// +build !windows

package base

import "C"
import (
	"plugin"
	"github.com/gw123/GMQ/core/interfaces"
)

/*
char * getModuleVersion()
int stop()
int handle(const char *event)
int start(const char *config, pub_f handel)
*/

type DllModule struct {
	BaseModule
	eventNames []string
	moduleDll  *plugin.Plugin
	handel     plugin.Symbol
	start      plugin.Symbol
	stop       plugin.Symbol
	getVersion plugin.Symbol
}

func NewDllModule() *DllModule {
	return new(DllModule)
}

func (this *DllModule) GetStatus() uint64 {
	return 1
}

func (this *DllModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	this.BaseModule.Init(app, config)
	var err error
	path := config.GetPath()
	//this.Debug("path :" + path)
	this.moduleDll, err = plugin.Open(path)

	if err != nil {
		this.Error("LoadDLL faild " + this.GetModuleName())
		return err
	}

	this.handel, err = this.moduleDll.Lookup("handle")
	if err != nil {
		this.Error("FindProc handel faild " + this.GetModuleName())
		return err
	}

	this.start, err = this.moduleDll.Lookup("start")
	if err != nil {
		this.Error("FindProc start faild " + this.GetModuleName())
		return err
	}

	this.stop, err = this.moduleDll.Lookup("stop")
	if err != nil {
		this.Error("FindProc stop faild " + this.GetModuleName())
		return err
	}
	//C.GoString((*C.char)(unsafe.Pointer(version)))
	this.getVersion, err = this.moduleDll.Lookup("getModuleVersion")
	if err != nil {
		this.Warning("FindProc getVersion faild " + this.GetModuleName())
		this.Version = ""
	} else {
		//version, _:= this.getVersion.(C.getModuleVersion)
		//this.Version = C.GoString((*C.char)(unsafe.Pointer(version)))
	}

	return nil
}

func (this *DllModule) Start() {

}

func (this *DllModule) UnInit() (err error) {
	return err
}
