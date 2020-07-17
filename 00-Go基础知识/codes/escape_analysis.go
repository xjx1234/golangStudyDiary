/**
* @Author: XJX
* @Description: 变量逃逸分析
* @File: escape_analysis.go
* @Date: 2020/5/20 10:56
 */

package main

import "fmt"

type Data struct {

}
// 本函数测试入口参数和返回值情况
func dummy(x int) int {
	// 声明一个变量c并赋值
	var c int
	c = x
	return c
}

func ts1() Data{
	var c1 Data
	return c1
}

func ts2()  *Data{
	var c2 Data
	return &c2
}
func main() {
	// 声明a变量并打印
	var a int = 8
	var b int = 7
	fmt.Println(b)
	fmt.Println(a)
	// 打印a变量的值和dummy()函数返回
	fmt.Println(dummy(0))
	ts1()
	ts2()
}
