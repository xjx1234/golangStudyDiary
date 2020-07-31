/**
* @Author: XJX
* @Description:
* @File: RpcClient.go
* @Date: 2020/7/30 15:37
 */

package main

import (
	"fmt"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Params struct {
	Radius float64
}

func main() {

	//基于HTTP RPC获取数据
	httpRpc, err := rpc.DialHTTP("tcp", "127.0.0.1:8081")
	if err != nil {
		log.Fatal(err)
	}
	ret := 0.0
	error := httpRpc.Call("Circular.GetPerimeter", Params{1.1}, &ret)
	if error != nil {
		log.Fatal(error)
	}
	fmt.Printf("Http Perimeter: %v \r\n", ret)

	//基于TCP的RPC获取数据
	tcpRpc, tcpErr := rpc.Dial("tcp", "127.0.0.1:8082")
	if tcpErr != nil {
		log.Fatal(tcpErr)
	}
	tcpRet := 0.0
	tcpErr1 := tcpRpc.Call("Circular.GetArea", Params{2.0}, &tcpRet)
	if tcpErr1 != nil {
		log.Fatal(tcpErr1)
	}
	fmt.Printf("TCP Area: %v \r\n", tcpRet)

	//基于JSONRPC获取数据
	jsonRpc, jsonErr := jsonrpc.Dial("tcp", "127.0.0.1:8083")
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	jsonRet := 0.0
	jsonErr1 := jsonRpc.Call("Circular.GetArea", Params{3.0}, &jsonRet)
	if jsonErr1 != nil {
		log.Fatal(jsonErr1)
	}
	fmt.Printf("JSONRPC Area: %v \r\n", tcpRet)

}
