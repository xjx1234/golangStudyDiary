/**
* @Author: XJX
* @Description:
* @File: RpcServer.go
* @Date: 2020/7/30 15:13
 */

package main

import (
	"log"
	"math"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
)

type Params struct {
	Radius float64
}

type Circular struct{}

const π = 3.1415

func (c *Circular) GetPerimeter(p Params, perimeter *float64) error {
	*perimeter = π * p.Radius * 2
	return nil
}

func (c *Circular) GetArea(p Params, area *float64) error {
	*area = π * math.Sqrt(p.Radius)
	return nil
}

func main() {
	listener := &Circular{}
	rpc.Register(listener)

	var wg sync.WaitGroup
	wg.Add(3)

	//基于HTTP协议的RPC服务
	go func() {
		rpc.HandleHTTP()
		//http.ListenAndServe(":8081", nil)
		lc, err := net.Listen("tcp", "127.0.0.1:8081")
		if err != nil {
			log.Fatal(err)
			defer wg.Done()
		}
		http.Serve(lc, nil)
	}()
	log.Println("http rpc service start success addr:8081")

	//基于TCP协议的RPC服务
	go func() {
		laddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8082")
		if err != nil {
			log.Fatal(err)
			defer wg.Done()
		} else {
			bg, err1 := net.ListenTCP("tcp", laddr)
			if err1 != nil {
				log.Fatal(err1.Error())
				defer wg.Done()
			}
			//rpc.Accept(bg)
			for {
				conn, err2 := bg.Accept()
				if err2 != nil {
					continue
				}
				go rpc.ServeConn(conn)
			}
		}
	}()
	log.Println("tcp rpc service start success addr:8082")

	//基于JSONRPC协议的RPC服务
	go func() {
		laddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8083")
		if err != nil {
			log.Fatal(err)
			defer wg.Done()
		} else {
			bg, err1 := net.ListenTCP("tcp", laddr)
			if err1 != nil {
				log.Fatal(err1.Error())
				defer wg.Done()
			}
			for {
				conn, err2 := bg.Accept()
				if err2 != nil {
					continue
				}
				go jsonrpc.ServeConn(conn)
			}
		}

	}()
	log.Println("json rpc service start success addr:8083")

	wg.Wait()
}
