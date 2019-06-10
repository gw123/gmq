package base

import (
	"github.com/gw123/GMQ/core/interfaces"
	"encoding/json"
	"os/exec"
	"io"
)

type ExeModule struct {
	BaseModule
	eventNames []string
	execPath   string
	execHandel *exec.Cmd
	outPipe    io.ReadCloser
	inPipe     io.WriteCloser
	errorPipe  io.ReadCloser
}

func NewExeModule() *ExeModule {
	return new(ExeModule)
}

func (this *ExeModule) GetStatus() uint64 {
	return 1
}

func (this *ExeModule) Init(app interfaces.App, config interfaces.ModuleConfig) (err error) {
	this.BaseModule.Init(app, this,config)
	this.execPath = config.GetPath()
	//this.Debug("path :" + path)
	this.execHandel = exec.Command(this.execPath)
	this.errorPipe, err = this.execHandel.StderrPipe()
	if err != nil {
		return err
	}
	this.inPipe, err = this.execHandel.StdinPipe()
	if err != nil {
		return err
	}
	this.outPipe, err = this.execHandel.StdoutPipe()
	if err != nil {
		return err
	}
	this.Start()
	return nil
}


func (this *ExeModule) Handle(event interfaces.Event) error {
	type Event struct {
		MsgId     string
		EventType string
		Payload   string
	}
	ev := Event{
		MsgId:     event.GetMsgId(),
		EventType: event.GetEventName(),
		Payload:   string(event.GetPayload()),
	}
	jsonData, err := json.Marshal(ev)
	this.execHandel.Start()
	this.inPipe.Write([]byte(string(jsonData) + "\n"))
	var runOverFlag = false
	go func() {
		buffer := make([]byte, 1024)
		for ; runOverFlag; {
			_, err := this.outPipe.Read(buffer)
			//len, err := this.outPipe.Read(buffer)
			if err != nil {
				this.Warning("outPipe " + err.Error())
			}
		}
	}()
	this.execHandel.Wait()
	res, err := this.execHandel.Output()
	runOverFlag = true
	if string(res) == "0" {
		//执行成功
		//replay := common.NewResultEvent([]byte("执行成功"))
		//this.App.Pub(replay)
		//this.Info(event.GetMsgId() + " " + event.GetEventName() + " 执行成功")
	} else {
		//replay := common.NewResultEvent([]byte("执行失败"))
		//this.App.Pub(replay)
		this.Error(event.GetMsgId() + " " + event.GetEventName() + " 执行失败" + err.Error())
	}
	return nil
}

func (this *ExeModule) Watch(index int) {

	return
}
