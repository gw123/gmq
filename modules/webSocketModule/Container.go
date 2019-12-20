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

func (c *Container) Write(buf []byte) (err error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	_, err = c.Buffer.Write(buf)
	if err != nil {
		return err
	}
	return
}

func (c *Container) Read() (line []byte, err error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	line, err = c.Buffer.ReadBytes('\n')
	return
}
