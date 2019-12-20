package gmq

import (
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

func (em *ExeModule) GetStatus() uint64 {
	return 1
}

func (em *ExeModule) Init(app App, config ModuleConfig) (err error) {
	em.BaseModule.Init(app, em,config)
	em.execPath = config.GetPath()
	//em.Debug("path :" + path)
	em.execHandel = exec.Command(em.execPath)
	em.errorPipe, err = em.execHandel.StderrPipe()
	if err != nil {
		return err
	}
	em.inPipe, err = em.execHandel.StdinPipe()
	if err != nil {
		return err
	}
	em.outPipe, err = em.execHandel.StdoutPipe()
	if err != nil {
		return err
	}
	em.Start()
	return nil
}


func (em *ExeModule) Handle(event Msg) error {
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
	em.execHandel.Start()
	em.inPipe.Write([]byte(string(jsonData) + "\n"))
	var runOverFlag = false
	go func() {
		buffer := make([]byte, 1024)
		for ; runOverFlag; {
			_, err := em.outPipe.Read(buffer)
			//len, err := em.outPipe.Read(buffer)
			if err != nil {
				em.Warning("outPipe " + err.Error())
			}
		}
	}()
	em.execHandel.Wait()
	res, err := em.execHandel.Output()
	runOverFlag = true
	if string(res) == "0" {
		//执行成功
		//replay := common_types.NewResultEvent([]byte("执行成功"))
		//em.App.Pub(replay)
		//em.Info(event.GetMsgId() + " " + event.GetEventName() + " 执行成功")
	} else {
		//replay := common_types.NewResultEvent([]byte("执行失败"))
		//em.App.Pub(replay)
		em.Error(event.GetMsgId() + " " + event.GetEventName() + " 执行失败" + err.Error())
	}
	return nil
}

func (em *ExeModule) Watch(index int) {

	return
}
