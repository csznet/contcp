package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
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
		fmt.Fprint(w, string(res))
	})
	http.HandleFunc("/status/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		str := strings.Replace(string(r.URL.Path), "/status/", "", 1)
		res := CheckServer(str)
		fmt.Fprintf(w, "%t", res.Status)
	})
	http.ListenAndServe("0.0.0.0:3344", nil)
}

func CheckServer(uri string) result {
	var res result
	if strings.Count(uri, ":") == 0 {
		// 使用 ICMP ping
		res = pingServer(uri)
	} else {
		// 使用 TCP 连接
		res = tcpCheckServer(uri)
	}
	return res
}

func tcpCheckServer(uri string) result {
	var res result
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
	} else {
		res.Status = false
	}
	return res
}

func pingServer(host string) result {
	var res result
	res.Host = host
	if strings.Count(host, ".") == 0 {
		res.Status = false
		return res
	}
	t1 := time.Now()
	res.Status = pingICMP(host)
	res.TimeStamp = time.Now().Sub(t1).String()
	return res
}

func pingICMP(host string) bool {
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Println("Error listening for ICMP packets:", err)
		return false
	}
	defer c.Close()

	dst, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		fmt.Println("Error resolving IP address:", err)
		return false
	}

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}
	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		fmt.Println("Error marshalling ICMP message:", err)
		return false
	}

	if _, err := c.WriteTo(msgBytes, dst); err != nil {
		fmt.Println("Error writing ICMP message:", err)
		return false
	}

	reply := make([]byte, 1500)
	c.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, _, err := c.ReadFrom(reply)
	if err != nil {
		fmt.Println("Error reading ICMP reply:", err)
		return false
	}

	parsedMsg, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		fmt.Println("Error parsing ICMP message:", err)
		return false
	}

	if parsedMsg.Type == ipv4.ICMPTypeEchoReply {
		return true
	}
	return false
}
