/**
* @Author: XJX
* @Description: websocket客户端
* @File: websocket_client.go
* @Date: 2020/7/14 14:24
 */

package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
)

func main() {
	origin := "http://127.0.0.1:8888/"
	url := "ws://127.0.0.1:8888/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ws)
	_, err = ws.Write([]byte("zzzzzzzzzz"))
	fmt.Println(err)

	for {
		var msg = make([]byte, 1000)
		m, err := ws.Read(msg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Receive: %s\n", msg[:m])
	}
}
