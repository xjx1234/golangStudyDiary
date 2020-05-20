/**
* @Author: XJX
* @Description:
* @File: ch1.go
* @Date: 2020/5/19 17:24
 */

package main

import "fmt"

var a string   //声明变量方式一
var b int
var c float32

//批量变量方式
var (
	d int
	e string
	f []float32
	g func() bool
	h struct{
		x int
	}
)

var t string = "hello" //初始化方式一

//批量初始化
var (
	t1 int = 5
	t2 bool =false
	t3 string = "你好"
)

var a3 int = 5 //全局变量

const pi = 3.1415926 //隐式类定义
const b1 string =  "hello" //显式类定义
const (
	b2 = 5
	b3 = "xxxxxxx"
)

// iota使用方法
const (
	Sunday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

const (
	FlagNone = 1 << iota
	FlagRed
	FlagGreen
	FlagBlue
)

type NewType int //讲NewType定义为int
type IntAlias = int //给int类型取别名为 IntAlias

func main() {
	fmt.Println(a)
    fmt.Println(b)
	fmt.Println(c)

	y := 1 //变量最简声明方式
	fmt.Println(y)

	t4 := "hello" //初始化方式三
	fmt.Println(t)
	fmt.Println(t1)
	fmt.Println(t2)
	fmt.Println(t3)
	fmt.Println(t4)

	//多重赋值
	var x1 int = 100
	var x2 int = 200
	x1,x2 = x2,x1
	fmt.Println(x1,x2)

	x3,_ := x1,x2
	fmt.Println(x3)

	//局部变量
	var a1 int = 4
	var a2 int = 3
	fmt.Println(a1+a2)

	fmt.Println(Sunday,Monday,Tuesday,Wednesday,Thursday,Friday,Saturday)
	fmt.Println(FlagNone,FlagRed, FlagGreen, FlagBlue)

	var c1 NewType
	var c2 IntAlias
	fmt.Printf("c1 type %T\n", c1)
	fmt.Printf("c2 type %T\n", c2)


}

//函数中a 和 b参数变量叫做形式参数(形参),形式参数只在函数调用时才会生效，函数调用结束后就会被销毁，在函数未被调用时，函数的形参并不占用实际的存储单元，也没有实际值。
func sum(a int, b int) int{
	num := a+b
	return num
}