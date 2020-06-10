/**
* @Author: XJX
* @Description: Context演示用例
* @File: context.go
* @Date: 2020/6/10 11:13
 */

package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {

	/*ctx, cancel := context.WithCancel(context.Background())
	randData := cancelFun(ctx)
	for n := range randData {
		if n >= 20 {
			cancel()
			break
		}
	}
	fmt.Println("running ... ")
	time.Sleep(time.Second * 2)*/

	/*ctx, f := context.WithDeadline(context.Background(), time.Now().Add(10))
	ctx, f := context.WithTimeout(context.Background(), 10*time.Second)
	timeOutFun(ctx)
	defer f()*/

	ctx := context.WithValue(context.Background(), "xjx", "hello")
	ctx = context.WithValue(ctx, "id", 1)
	valueFun(ctx)

}

func valueFun(ctx context.Context) {
	id, ok := ctx.Value("id").(int)
	if !ok {
		fmt.Println("id value get fail")
	}
	xjxData := ctx.Value("xjx")
	fmt.Printf("id: %d \n", id)
	fmt.Printf("xjx data: %s \n", xjxData)

}

func timeOutFun(ctx context.Context) {
	n := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("timeOut")
			return
		default:
			incr := rand.Intn(100)
			n += incr
			fmt.Printf("随机数 %d \n", n)
		}
		time.Sleep(time.Second)
	}
}

func cancelFun(ctx context.Context) <-chan int {
	c := make(chan int)
	num := 0
	t := 0
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("耗时 %d 秒， 随机数 %d \n", t, num)
				return
			case c <- num:
				incr := rand.Intn(10)
				num += incr
				if num >= 20 {
					num = 20
				}
				t++
				fmt.Printf("随机数 %d \n", num)
			}
		}
	}()
	return c
}
