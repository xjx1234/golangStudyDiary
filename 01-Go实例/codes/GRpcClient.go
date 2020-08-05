/**
* @Author: XJX
* @Description: GRPC客户端
* @File: GRpcClient.go
* @Date: 2020/8/5 16:10
 */

package main

import (
	"calculator"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

const (
	Addres = "localhost:8808"
)

func main() {
	conn, err := grpc.Dial(Addres, grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()
	c := calculator.NewCalculateClient(conn)
	r, err := c.Add(context.Background(), &calculator.CalParams{P1: 1.1, P2: 2.2})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r)
	r1, err1 := c.Multiplication(context.Background(), &calculator.CalParams{P1: 1, P2: 2})
	if err1 != nil {
		log.Fatal(err1)
	}
	fmt.Println(r1)
}
