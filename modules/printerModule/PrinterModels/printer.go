package PrinterModels

import (
	"github.com/fpay/escpos-go/printer/driver/escpos"
	"github.com/fpay/escpos-go/printer/connection"
	"github.com/pkg/errors"
	"fmt"
	"net"
	"strconv"
	"io/ioutil"
	"time"
)

const USB = 0x01
const Net = 0x02

const ONLINE = 0x01
const OFFLINE = 0x02

type Printer struct {
	Id          int    `json:"id"`
	Type        string `json:"type"`
	PrinterType uint   `json:"portType"`
	Status      uint   `json:"status"`
	Address     string `json:"portName"`

	DestPort      string `json:"destPort"`
	PortName      string
	DeviceName    string `json:"deviceName"`
	driver        *escpos.Escpos
	LastUpateTime int64 `json:"lastUpateTime"`
}

func NewPrinter(printerType uint, addr string) (this *Printer) {
	this = new(Printer)
	if printerType != USB && printerType != Net {
		return nil
	}
	this.PrinterType = printerType
	this.Address = addr
	this.Status = ONLINE
	return this
}

func (this *Printer) GetDriver() (conn *escpos.Escpos, err error) {
	if this == nil {
		return nil, errors.New("empty object ")
	}
	switch this.PrinterType {
	case USB:
		conn, err := connection.NewUsbConnection(this.Address)
		if err != nil {
			return nil, err
		}
		this.driver = escpos.NewEscpos(conn)
		break
	case Net:
		conn, err := connection.NewNetConnection(this.Address + ":9100")
		if err != nil {
			return nil, err
		}
		this.driver = escpos.NewEscpos(conn)
		break
	}
	return this.driver, nil
}

func (this *Printer) PrinteRaw(content []byte) error {
	driver, err := this.GetDriver()
	if err != nil {
		return err
	}
	defer driver.Close()
	err = this.CheckStatusAndWait(driver)
	if err != nil{
		fmt.Println(err)
	}
	_, err = driver.WriteRaw(content)
	if err != nil {
		return err
	}
	driver.Linefeed()
	return nil
}


func (this *Printer) PrinteXml(content string) error {
	root, err := escpos.ParseString(content)
	if err != nil {
		return err
	}
	driver, err := this.GetDriver()
	if err != nil {
		return err
	}
	defer driver.Close()
	err = this.CheckStatusAndWait(driver)
	if err != nil{
		fmt.Println(err)
	}
	err = driver.WriteXml(root)
	if err != nil {
		return err
	}
	driver.Linefeed()
	return nil
}

func (this *Printer) CheckStatusAndWait(escpos *escpos.Escpos) error {
	for i := 0; i < 3; i++ {
		status, _ := escpos.ReadStatus()
		if status {
			return nil
		}
	}
	return errors.New("打印机状态检查失败")
}

func (this *Printer) Print(content string) error {
	driver, err := this.GetDriver()
	if err != nil {
		return err
	}
	defer driver.Close()
	_, err = driver.WriteGbk(content)
	if err != nil {
		return err
	}
	driver.Linefeed()
	return nil
}

func (this *Printer) StartServer() {

	var l net.Listener
	var err error
	//port := 9100 + this.Id
	port := 9100 + this.Id - 1

	for ; ; {
		l, err = net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
		if err != nil {
			fmt.Println("StartServer", err)
			time.Sleep(time.Second * 2)
			continue
		}

		break
	}
	fmt.Println("StartServer", "监听本地:", port)
	defer l.Close()

	for ; ; {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go this.handelRequest(conn)
	}
}

func (this *Printer) handelRequest(conn net.Conn) {
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(time.Second * 30))
	conn.SetWriteDeadline(time.Now().Add(time.Second * 30))
	buffer, err := ioutil.ReadAll(conn)
	if len(buffer) == 0 {
		return
	}
	if err != nil && len(buffer) == 0 {
		return
	}
	this.PrinteRaw(buffer)
	time.Sleep(time.Second * 4)
}

func (this *Printer) Info() string {
	return fmt.Sprintf("PrinterType %d  DeviceName %s  Address %s ", this.PrinterType, this.DeviceName, this.Address)
}
