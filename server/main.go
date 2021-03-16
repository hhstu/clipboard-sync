package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hhstu/clipboard-sync/pkg"
	"io/ioutil"
	"log"
	"net/http"
)

var data pkg.Object
var token string

func push(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("token") != token {
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err, body)
		return
	}
	log.Println("get msg from ", r.RemoteAddr, "data: ", data.Data)
	fmt.Fprintf(w, "ok")
}

func pull(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("token") != token {
		return
	}
	resp, err := json.Marshal(data)
	if err != nil {
		log.Println(err, resp)
	}

	fmt.Fprintf(w, string(resp))
}
func init() {
	flag.StringVar(&token, "token", "lc666", "token")
	flag.Parse()
}
func main() {
	http.HandleFunc("/push", push)
	http.HandleFunc("/pull", pull)
	http.ListenAndServe(":8080", nil)
}
