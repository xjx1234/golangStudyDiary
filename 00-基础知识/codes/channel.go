/**
* @Author: XJX
* @Description: 无缓冲以及有缓冲通道示例
* @File: channel.go
* @Date: 2020/6/2 16:18
 */

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	unbufferedChan := make(chan int) //创建无缓存通道
	fmt.Printf("leb(c)=%d, cap(c)=%d\n", len(unbufferedChan), cap(unbufferedChan))
	go func() {
		defer fmt.Println("子协程结束")
		for i := 0; i < 3; i++ {
			fmt.Println("here", i)
			unbufferedChan <- i
			fmt.Printf("子进程正在运行[%d]: len(c)=%d,cap(c)=%d \n", i, len(unbufferedChan), cap(unbufferedChan))
		}
	}()

	for i := 0; i < 3; i++ {
		num := <-unbufferedChan
		fmt.Println("num=", num)
	}

	fmt.Println("main主程序结束")

	bufferedChan := make(chan int, 3)
	fmt.Printf("len(c)=%d,cap(c)=%d\n", len(bufferedChan), cap(bufferedChan))
	go func() {
		defer fmt.Println("子协程结束")
		for i := 0; i < 3; i++ {
			bufferedChan <- i
			fmt.Printf("子协程正在运行[%d]: len(c)=%d, cap(c)=%d\n", i, len(bufferedChan), cap(bufferedChan))
		}
	}()
	for i := 0; i < 3; i++ {
		num := <-bufferedChan
		fmt.Println("num=", num)
	}
	fmt.Println("main主程序结束")

	ch1 := make(chan int)
	ch2 := make(chan int)
	go pump1(ch1)
	go pump2(ch2)
	go suck(ch1, ch2)
	time.Sleep(1e9)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go pumpNum(i)
	}
	wg.Wait()

}

var wg sync.WaitGroup

func pumpNum(num int) {
	defer wg.Done()
	fmt.Println(num)
}

func pump1(ch chan int) {
	for i := 0; ; i++ {
		fmt.Println("pump1....")
		ch <- i * 2
	}
}

func pump2(ch chan int) {
	for i := 0; ; i++ {
		fmt.Println("pump2.....")
		ch <- i + 5
	}
}

func suck(ch1, ch2 chan int) {
	for {
		select {
		case v := <-ch1:
			fmt.Printf("Received channel1 : %d \n", v)
		case v := <-ch2:
			fmt.Printf("Received channel2 : %d \n", v)
		}
	}
}
