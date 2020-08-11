/**
* @Author: XJX
* @Description: UDP server 演示代码
* @File: UDPServer.go
* @Date: 2020/8/10 11:36
 */

package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
)

//定义ip以及端口
const (
	Address string = "127.0.0.1:8808"
)

//定义信息消息体
type UdpMessage struct {
	Addr *net.UDPAddr
	Data []byte
}

//定义连接结构体
type Server struct {
	Conn       *net.UDPConn
	ReadChan   chan *UdpMessage
	WriteChan  chan *UdpMessage
	ClosedChan chan byte
}

//处理UDP信息接收以及发送处理
func (s *Server) procLoop() {
	fmt.Println("procLoop")
	for {
		msg, err := s.read()
		if err != nil {
			log.Println("获取消息出现错误", err.Error())
			break
		}
		log.Println("接收到消息", string(msg.Data))

		// 以下内容为获取内容后得返回信息，可以根据实际情况修改
		err = s.write(msg)
		if err != nil {
			log.Println("发送消息给客户端出现错误", err.Error())
			break
		}
	}
}

//将UDP信息接收后存入管道
func (s *Server) readLoop() {
	fmt.Println("readLoop")
	for {
		data := make([]byte, 2048)
		n, addr, err := s.Conn.ReadFromUDP(data)
		if err != nil {
			goto End
		}
		fmt.Println("have data, data len:" + strconv.Itoa(n))
		req := &UdpMessage{
			addr,
			data[:n],
		}
		select {
		case s.ReadChan <- req:
		case <-s.ClosedChan:
			goto End
		}
	}
End:
	s.udpClosed()
}

//将写管道的信息发送给UDP客户端
func (s *Server) writeLoop() {
	fmt.Println("writeLoop")
	for {
		select {
		case msg := <-s.WriteChan:
			if _, wErr := s.Conn.WriteToUDP(msg.Data, msg.Addr); wErr != nil {
				goto End
			}
		case <-s.ClosedChan:
			goto End
		}
	}
End:
	s.udpClosed()
}

//读处理
func (s *Server) read() (*UdpMessage, error) {
	select {
	case msg := <-s.ReadChan:
		return msg, nil
	case <-s.ClosedChan:
	}
	return nil, errors.New("udp closed")
}

//写处理
func (s *Server) write(d *UdpMessage) error {
	select {
	case s.WriteChan <- d:
	case <-s.ClosedChan:
		return errors.New("udp closed")
	}
	return nil
}

//关闭处理
func (s *Server) udpClosed() {
	s.Conn.Close()
	close(s.ClosedChan)
}

func main() {

	var wg sync.WaitGroup
	wg.Add(1)
	udpAddress, err := net.ResolveUDPAddr("udp", Address)
	if err != nil {
		log.Fatal(err)
	}
	conn, connErr := net.ListenUDP("udp", udpAddress)
	defer conn.Close()
	if connErr != nil {
		log.Fatal(connErr)
	}
	srv := &Server{
		Conn:      conn,
		ReadChan:  make(chan *UdpMessage, 50),
		WriteChan: make(chan *UdpMessage, 50),
	}
	fmt.Println("hello")
	go srv.procLoop()
	go srv.readLoop()
	go srv.writeLoop()
	wg.Wait()
}
