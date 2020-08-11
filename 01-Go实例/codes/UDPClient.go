/**
* @Author: XJX
* @Description: UDP客户端演示代码
* @File: UDPClient.go
* @Date: 2020/8/10 11:36
 */

package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	ServerAddr string = "127.0.0.1:8808"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", ServerAddr)
	if err != nil {
		fmt.Println("x")
		log.Fatal(err)
	}
	conn, lerr := net.DialUDP("udp", nil, udpAddr)
	if lerr != nil {
		log.Fatal(lerr)
	}
	for {
		_, err = conn.Write([]byte("hello xjx")) // 发送数据给UDP服务器
		if err != nil {
			fmt.Println("udp 发送数据失败", err)
		}

		buf := make([]byte, 2048)
		n, readErr := conn.Read(buf)
		if readErr != nil {
			log.Fatal(readErr)
		}
		fmt.Println("client read len:", n)
		fmt.Println("client read data:", string(buf))

		time.Sleep(1000)
	}
}
