package printerModule

import (
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/gworker/jobs"
	"strconv"
	"net"
	"time"
	"github.com/gw123/GMQ/common/common_types"
)

type ScanJob struct {
	jobs.BaseJob
	printerModule interfaces.Module
}

func NewScanJob(ip string, module interfaces.Module) (job *ScanJob) {
	job = new(ScanJob)
	job.CreatedTime = time.Now().Unix()
	job.Flag = jobs.JobFlagNormal
	job.Payload = []byte(ip)
	job.printerModule = module
	return
}

func (this *ScanJob) DoJob() {
	//fmt.Println("执行任务：", this.WorkerName, string(this.Payload))
	ip := string(this.Payload)
	addr := ip + ":" + strconv.Itoa(9100)
	conn, err := net.DialTimeout("tcp", addr, time.Second*2)
	if err != nil {
		return
	}
	defer conn.Close()
	event := common_types.NewEvent("findPrinter", []byte(addr))
	this.printerModule.Pub(event)
	this.printerModule.Info("DoJob " + addr + " 连接成功")
}
