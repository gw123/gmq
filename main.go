package main

import (
	_ "net/http/pprof"
	"io/ioutil"
	"fmt"
	"net/http"
	"log"
	"github.com/gw123/GMQ/core"
)

func main() {
	data, err := ioutil.ReadFile("main.ini")
	if err != nil {
		fmt.Println(err)
		return
	}
	appInstance := core.NewApp(data)
	appInstance.Start()

	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	select {

	}
	//fmt.Println(app.EventQueues)
}
