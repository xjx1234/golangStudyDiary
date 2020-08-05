/**
* @Author: XJX
* @Description: GRPC服务端
* @File: GRpcServer.go
* @Date: 2020/8/3 11:23
 */

package main

import (
	"calculator"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Cal struct{}

const (
	Port = ":8808"
)

func (c *Cal) Add(ctx context.Context, p *calculator.CalParams) (*calculator.ResultRes, error) {
	last := p.P1 + p.P2
	return &calculator.ResultRes{Res: last}, nil
}

func (c *Cal) Sub(ctx context.Context, p *calculator.CalParams) (*calculator.ResultRes, error) {
	last := p.P1 - p.P2
	return &calculator.ResultRes{Res: last}, nil
}

func (c *Cal) Multiplication(ctx context.Context, p *calculator.CalParams) (*calculator.ResultRes, error) {
	last := p.P1 * p.P2
	return &calculator.ResultRes{Res: last}, nil
}

func (c *Cal) Division(ctx context.Context, p *calculator.CalParams) (*calculator.ResultRes, error) {
	last := p.P1 / p.P2
	return &calculator.ResultRes{Res: last}, nil
}

func main() {
	lis, err := net.Listen("tcp", Port)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	calculator.RegisterCalculateServer(s, &Cal{})
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
		return
	}
}
