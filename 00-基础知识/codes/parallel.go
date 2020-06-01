/**
* @Author: XJX
* @Description: 并发章节用例
* @File: parallel.go
* @Date: 2020/6/1 17:49
 */

package main

import (
	"fmt"
	"time"
)

func hello() {
	fmt.Println("hello goroutine")
}

func hello2(ch chan bool) {
	ch <- true
	fmt.Println("hello ch")
}

func main() {
	go hello() // 开启协程调用hello函数
	//开启协程 调用匿名函数
	go func(name string) {
		fmt.Println("hello " + name + " goroutine\n")
	}("anonymous")
	fmt.Println("main function\n")
	time.Sleep(1)

	var ch1 chan int               //声明通道
	ch1 = make(chan int)           //创建通道
	cha2 := make(chan interface{}) //声明并创建通道
	fmt.Printf("%v, %v", ch1, cha2)
	fmt.Println("--------------------")

	ch := make(chan bool)
	go hello2(ch)
	<-ch
	fmt.Println("run end")

	squaresChan := make(chan int)
	cubesChan := make(chan int)
	go calSquares(2, squaresChan)
	go calCubes(3, cubesChan)
	s1, s2 := <-squaresChan, <-cubesChan
	fmt.Printf("the sum is %d\n", s1+s2)

	commChan := make(chan int)
	sendOnlyChan := make(chan<- int)
	sendOnlyChan = commChan
	//recOnlyChan := make(<-chan int)
	//sendOnlyChan <- 1
	//<-recOnlyChan

}

func calSquares(num int, sumData chan int) {
	result := num * num
	sumData <- result
}

func calCubes(num int, sumData chan int) {
	result := num * num * num
	sumData <- result
}
