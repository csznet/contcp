package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type result struct {
	Host      string `json:"host"`
	TimeStamp string `json:"time"`
	Status    bool   `json:"status"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		str := strings.Replace(string(r.URL.Path), "/", "", 1)
		res, _ := json.Marshal(CheckServer(str))
		fmt.Fprintf(w, string(res))
	})
	http.ListenAndServe("0.0.0.0:3344", nil)
}

func CheckServer(uri string) result {
	var res result
	if strings.Count(uri, ":") == 0 {
		uri += ":80"
	}
	res.Host = uri
	if strings.Count(uri, ".") == 0 {
		res.Status = false
		return res
	}
	timeout := time.Duration(5 * time.Second)
	t1 := time.Now()
	_, err := net.DialTimeout("tcp", uri, timeout)
	res.TimeStamp = time.Now().Sub(t1).String()
	if err == nil {
		res.Status = true
	}
	return res
}
