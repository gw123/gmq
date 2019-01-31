package printerModule

import (
	"github.com/gw123/GMQ/modules/base"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/net_tool/net_utils"
	"github.com/gw123/gworker"
	"github.com/gw123/GMQ/common/common_types"
	"github.com/gw123/GMQ/modules/printerModule/PrinterModels"
	"encoding/json"
	"github.com/fpay/openwrt-gateway-go/src/gateway/model"
	"time"
	"strconv"
	"github.com/fpay/openwrt-gateway-go/src/gateway/util"
	"encoding/base64"
)

type PrinterModule struct {
	base.BaseModule
	printerJobs         chan interfaces.Event
	doPrinterJobEndChan chan uint8
	GwDeviceName        string
	printers            []*PrinterModels.Printer `json:"printers"`
}

func NewPrinterModule() *PrinterModule {
	this := new(PrinterModule)
	this.printerJobs = make(chan interfaces.Event, 1024)
	this.doPrinterJobEndChan = make(chan uint8, 1)
	return this
}

func (this *PrinterModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	this.BaseModule.Init(app, config)
	app.Sub("printe", this)
	app.Sub("findPrinter", this)
	app.Sub("scanPrinterOver", this)
	go this.doScan()
	return nil
}

func (this *PrinterModule) doScan() {
	this.Info("start scanPirinter")
	ipList := net_utils.GetIpList()
	group := gworker.NewWorkerGroup(130)
	group.Start()

	for _, ip := range ipList {
		job := NewScanJob(ip, this)
		group.DispatchJob(job)
	}
	group.WaitEmpty()
	event := common_types.NewEvent("scanPrinterOver", []byte(""))
	this.Pub(event)
}

func (this *PrinterModule) doPrinterJob() {
	for {
		printerJob := <-this.printerJobs
		if printerJob.GetEventName() == "printe" {
			this.printe(printerJob.GetPayload())
		}
		if printerJob.GetEventName() == "stop" {
			break
		}
	}
	this.doPrinterJobEndChan <- 0
}

func (this *PrinterModule) printe(message []byte) error {
	this.Info("打印" + string(message))
	printeMsg := common_types.LhMsg{}
	err := json.Unmarshal(message, &printeMsg)

	if err != nil {
		return err
	}
	for _, deviceName := range printeMsg.DeviceNames {
		printer := this.FindPrinterByDeviceName(deviceName)
		if printer == nil {
			printeResponse := common_types.NewResultEvent([]byte("打印失败,设备不存在:" + deviceName))
			printeResponse.MsgId = printeMsg.MsgId
			this.Pub(printeResponse)
			continue
		}
		var err error
		switch printeMsg.EventName {
		case "print_base64":
			data, err := base64.StdEncoding.DecodeString(string(printeMsg.Payload))
			if err == nil {
				err = printer.PrinteRaw(data)
			}
			break
		case "print":
			err = printer.PrinteXml(string(printeMsg.Payload))
			break

		case "print_raw":
			err = printer.PrinteRaw([]byte(printeMsg.Payload))
			break
		}

		//上报打印结果
		printeResponse := common_types.NewResultEvent([]byte("打印失败,设备不存在:" + deviceName))
		printeResponse.MsgId = printeMsg.MsgId
		if err != nil {
			printeResponse.Payload = "打印失败" + err.Error()
			this.Pub(printeResponse)
		} else {
			printeResponse.Payload = "打印成功"
			this.Pub(printeResponse)
		}
	}

	return nil
}

func (this *PrinterModule) Start() {
	for ; ; {
		event := this.BaseModule.Pop()
		if event.GetEventName() == "stop" && string(event.GetPayload()) == this.GetModuleName() {
			this.service(event)
			<-this.doPrinterJobEndChan
			this.Info("all printe job run over,PrinterModule stop")
			break
		}

		err := this.service(event)
		if err != nil {
			this.Error("PrinterModule service " + err.Error())
		}
	}
}

func (this *PrinterModule) service(event interfaces.Event) error {
	this.Info(event.GetEventName() + ", " + event.GetMsgId() + " ," + string(event.GetPayload()))
	switch event.GetEventName() {
	case "printe":
		this.printerJobs <- event
		break
	case "findPrinter":
		break
	case "scanPrinterOver":
		break
	}
	return nil
}

func (this *PrinterModule) FindPrinter(printerType uint, addr string) *PrinterModels.Printer {
	for _, printer := range this.printers {
		if printer.Address == addr {
			return printer
		}
	}
	return nil
}

func (this *PrinterModule) FindPrinterByDeviceName(deviceName string) *PrinterModels.Printer {
	for _, printer := range this.printers {
		if printer.DeviceName == deviceName {
			return printer
		}
	}
	return nil
}

func (this *PrinterModule) CheckIsNewPrinter(printerType uint, addr string) bool {
	for _, printer := range this.printers {
		if printer.Address == addr && printerType == printer.PrinterType {
			return false
		}
	}
	return true
}

func (this *PrinterModule) AppendNewPrinter(printer *PrinterModels.Printer) {
	if printer == nil {
		this.Error("AppendNewPrinter 参数验证失败")
		return
	}

	p := this.FindPrinter(printer.PrinterType, printer.Address)
	if p == nil {
		//新设备
		printer.LastUpateTime = time.Now().Unix()
		id := 1 + len(this.printers)
		printer.DeviceName = this.GwDeviceName + "_sub" + strconv.Itoa(id)
		if id < 10 {
			printer.DeviceName = this.GwDeviceName + "_sub0" + strconv.Itoa(id)
		}

		printer.Id = id
		printer.DestPort = util.GetWanIp() + ":" + strconv.Itoa(printer.Id+9100-1)

		printer.PortName = printer.DeviceName
		printer.Type = "printer"
		if printer.PrinterType == model.USB {
			//fmt.Println("端口映射地址 Adress : " , printer.DestPort)
			go printer.StartServer()
		}
		this.printers = append(this.printers, printer)
		return
	}

	if p != nil && p.Status == model.OFFLINE {
		p.LastUpateTime = time.Now().Unix()
		p.Status = model.ONLINE
		return
	}
}
