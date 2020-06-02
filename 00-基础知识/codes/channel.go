/**
* @Author: XJX
* @Description: 无缓冲以及有缓冲通道示例
* @File: channel.go
* @Date: 2020/6/2 16:18
 */

package main

import "fmt"

func main() {

	/*unbufferedChan := make(chan int) //创建无缓存通道
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

	fmt.Println("main主程序结束")*/

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

}
