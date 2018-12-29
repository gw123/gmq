package main

import (
	"github.com/gw123/GMQ/app"
	_ "net/http/pprof"
	"io/ioutil"
	"fmt"
	"net/http"
	"log"
)

func main() {
	data, err := ioutil.ReadFile("main.ini")
	if err != nil {
		fmt.Println(err)
		return
	}
	appInstance := app.NewApp(data)
	appInstance.Start()

	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	select {

	}
	//fmt.Println(app.EventQueues)
}
