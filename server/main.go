package main

import (
	"encoding/json"
	"fmt"
	"github.com/hhstu/clipboard-sync/pkg"
	"io/ioutil"
	"log"
	"net/http"
)

var data pkg.Object

func push(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	log.Println("get msg from ", r.RemoteAddr, "data: ", data.Data)
	fmt.Fprintf(w, "ok")
}

func pull(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, string(resp))
}

func main() {
	http.HandleFunc("/push", push)
	http.HandleFunc("/pull", pull)
	http.ListenAndServe(":8080", nil)
}
