package core

import (
	"github.com/gw123/GMQ/core/interfaces"
	"fmt"
	"time"
	"os"
	"bytes"
	"strings"
	"strconv"
	"sync"
)

/*
 * 模块管理模块
 * 1. 加载模块,卸载模块
 * 2. 管理配置,更新模块配置
*/
type LogManager struct {
	app                   interfaces.App
	fileHandel            *os.File
	isAsync               bool
	interval              int
	filter_log_categories []string
	only_log_categories   []string
	level                 int
	buffer                *bytes.Buffer
	mutex                 sync.Mutex
}

func NewLogManager(app interfaces.App) *LogManager {
	this := new(LogManager)
	this.app = app
	this.buffer = new(bytes.Buffer)
	this.isAsync = true
	this.interval = 1
	var err error

	filter_categories, err := this.app.GetDefaultConfigItem("logFilterCategories")
	if err != nil {
		fmt.Println("LogManager GetDefaultConfigItem", err)
	}

	log_only_categories, err := this.app.GetDefaultConfigItem("logOnlyCategories")
	if err != nil {
		fmt.Println("LogManager GetDefaultConfigItem", err)
	}

	log_interval, err := this.app.GetDefaultConfigItem("logInterval")
	if err != nil {
		fmt.Println("LogManager GetDefaultConfigItem", err)
	}
	if filter_categories != "" {
		this.filter_log_categories = strings.Split(filter_categories, ",")
	}
	if log_only_categories != "" {
		this.only_log_categories = strings.Split(log_only_categories, ",")
	}

	if log_interval == "" {
		this.interval = 1
	} else {
		this.interval, err = strconv.Atoi(log_interval)
	}
	if err != nil {
		fmt.Println("log_interval error ", err)
	}
	if this.interval < 1 {
		this.interval = 1
	}

	date := time.Now().Format("2006-01-02")
	this.fileHandel, err = os.OpenFile(date+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		this.Error("LogManager", err.Error())
	}
	return this
}

func (this *LogManager) SetIsAsync(flag bool) {
	this.isAsync = flag
}

func (this *LogManager) SetIntVal(interval int) {
	this.interval = interval
}

func (this *LogManager) Info(category string, format string, a ...interface{}) {
	this.filter("Info", category, format, a...)
}

func (this *LogManager) Waring(category string, format string, a ...interface{}) {
	this.filter("Waring", category, format, a...)
}

func (this *LogManager) Error(category string, format string, a ...interface{}) {
	this.filter("Error", category, format, a...)
}

func (this *LogManager) Debug(category string, format string, a ...interface{}) {
	this.filter("Debug", category, format, a...)
}

func (this *LogManager) filter(logType, category string, format string, a ...interface{}) {
	//fmt.Println(this.filter_log_categories)
	//fmt.Println(this.only_log_categories)
	var pass = true
	//fmt.Println("category", len(this.only_log_categories))
	if len(this.only_log_categories) != 0 {
		pass = false
		for _, cate := range this.only_log_categories {
			if cate == category {
				pass = true
			}
		}
	} else {

		for _, cate := range this.filter_log_categories {
			if cate == category {
				return
			}
		}
	}

	if !pass {
		return
	}
	tpl := "[%s,%s,%s] %s"
	date := time.Now().Format("2006-01-02 15:04:05")
	data := fmt.Sprintf(tpl, date, logType, category, fmt.Sprintf(format, a...))
	this.Write(data)
}

func (this *LogManager) Write(data string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	//this.fileHandel.Write([]byte(data))
	fmt.Println(data)
	data = data + "\r\n"
	if this.buffer != nil {
		this.buffer.Write([]byte(data))
	} else {
		fmt.Println("LogManager buffer is nil")
		this.fileHandel.Write([]byte(data))
	}
}

func (this *LogManager) Start() {
	go func() {
		var buffer = make([]byte, 4096)
		var len = 1
		var err error
		for ; ; {
			time.Sleep(time.Second * time.Duration(this.interval))
			for ; len >= 0; {
				this.mutex.Lock()
				len, err = this.buffer.Read(buffer)
				this.mutex.Unlock()
				if err != nil {
					if err.Error() == "EOF" {
						break
					}
					fmt.Println("LogManager start", err)
					break
				}
				_, err = this.fileHandel.Write(buffer[0:len])
				if err != nil {
					fmt.Println("LogManager Write", err)
				}
			}
			len = 1
		}
	}()
}

//func (this *ConfigManager) LoadModule1(muduleName string, config []byte) {
//	module := this.Modules[muduleName]
//	module.Init(this.AppInstance, config)
//}
