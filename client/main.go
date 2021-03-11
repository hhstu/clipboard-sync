package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/hhstu/clipboard-sync/pkg"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

var data pkg.Object
var server string

func init() {
	flag.StringVar(&server, "server", "http://localhost:8080", "server url")
}
func main() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			pull()

		}
	}
}
func pull() {
	pull := fmt.Sprintf("%s/pull", server)
	ret, err := http.Get(pull)

	if err != nil {
		panic(err)
	}
	defer ret.Body.Close()

	body, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		panic(err)
	}
	var tmp pkg.Object
	err = json.Unmarshal(body, &tmp)
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(data, tmp) {
		data = tmp
		clipboard.WriteAll(data.Data)
	} else {
		tmp, err := clipboard.ReadAll()
		if err != nil {
			panic(err)
		}
		if tmp != data.Data {
			data.Data = tmp
			push()
		}
	}

}

func push() {
	push := fmt.Sprintf("%s/push", server)
	dataStr, err := json.Marshal(data)
	req, err := http.NewRequest("POST", push, bytes.NewBuffer(dataStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}
