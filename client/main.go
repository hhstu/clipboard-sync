package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/hhstu/clipboard-sync/pkg"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"
)

var data pkg.Object
var server string
var needPush bool
var token string

func init() {
	flag.StringVar(&server, "server", "http://localhost:8080", "server url")
	flag.BoolVar(&needPush, "needPush", true, "push or not")
	flag.StringVar(&token, "token", "lc666", "token")
	flag.Parse()
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
		log.Println(err)
		return
	}
	ret.Header.Add("token", token)
	defer ret.Body.Close()

	body, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var tmp pkg.Object
	err = json.Unmarshal(body, &tmp)
	if err != nil {
		log.Println(err)
		return
	}
	if !reflect.DeepEqual(data, tmp) {
		data = tmp
		clipboard.WriteAll(data.Data)
	} else {
		if needPush {
			tmp, err := clipboard.ReadAll()
			if err != nil {
				log.Println(err)
				return
			}
			if tmp != data.Data {
				data.Data = tmp
				push()
			}
		}
	}

}

func push() {
	push := fmt.Sprintf("%s/push", server)
	dataStr, err := json.Marshal(data)
	req, err := http.NewRequest("POST", push, bytes.NewBuffer(dataStr))
	if err != nil {
		log.Println(err)
		return
	}
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", token)
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Println(err)
	}

}
