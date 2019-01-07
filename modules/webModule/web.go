package main

import (
	"net/http"
	"encoding/json"
	"fmt"
)

func main() {
	handle()
}

func handle()  {
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		type Response struct {
			Status  int       `json:"status"`
			Message string    `json:"message"`
		}
		response := Response{}
		response.Message = "success"
		data, err := json.Marshal(response)
		if err != nil {
			w.Write([]byte("json 压缩失败"))
			return
		}
		//res := strings.Replace(infoPage, "{{$data}}", string(data), 1)
		w.Write(data)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
