package webSocketModule

import (
	"bytes"
	"sync"
)

type Container struct {
	Token  string //容器标识
	Buffer bytes.Buffer
	Mutex  sync.Mutex
}

func (this *Container) Write(buf []byte) (err error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	_, err = this.Buffer.Write(buf)
	if err != nil {
		return err
	}
	return
}

func (this *Container) Read() (line []byte, err error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	line, err = this.Buffer.ReadBytes('\n')
	return
}
