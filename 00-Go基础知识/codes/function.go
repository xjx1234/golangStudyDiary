/**
* @Author: XJX
* @Description: 函数示例
* @File: function.go
* @Date: 2020/5/22 14:54
 */

package main

import "fmt"

func bill(mileage float64, price float64) float64 {
	return mileage * price
}

func add(x int, y int) int {
	return x + y
} // 常规声明

func sub(x, y int) (z int) {
	z = x - y
	return
} // 返回值命名，只需要return即可

func first(x int, _ int) int {
	return x
} // 空白标识符 _ 可以强调某个参数未被使用

func zero(int, int) int {
	return 0
}

func moreValue(x, y int) (int, int) {
	return (x * y), (x + y)
}

func variableFun() bool {
	return true
}

func variableFun1(x int) int {
	return x
}

func visit(listData []int, f func(int) int) {
	for _, v := range listData {
		f(v)
	}
}

// 定义动物接口 接口定义后续章节会详细说明
type Animal interface {
	Call(interface{})
}

type AnimalCaller func(interface{}) // 将函数定义为类型

//实现Animal的Call
func (a AnimalCaller) Call(p interface{}) {
	a(p)
}

func main() {

	money := bill(3, 4)
	multiplyData, addData := moreValue(2, 3)
	fmt.Println(money)
	fmt.Println(multiplyData, addData)

	/** 函数当做变量使用示例 **/
	var f1 func() bool
	f1 = variableFun
	fmt.Println(f1())

	var f2 func(int) int
	f2 = variableFun1
	fmt.Println(f2(1))

	/** 定义匿名函数为变量 **/
	anonymousFun := func(x int) int {
		return x
	}
	anonymousFun(1) //直接变量调用

	/** 直接运行匿名函数 **/
	func(x int) (y int) {
		y = x
		return
	}(4)

	visit([]int{1, 2, 3, 4, 5, 6}, func(x int) int {
		fmt.Println(x)
		return x
	})

	closureFun := closureShow(1)
	fmt.Println(closureFun())
	fmt.Println(closureFun())

	var myAnimal Animal
	myAnimal = AnimalCaller(func(v interface{}) {
		fmt.Println("animal is ", v)
	})
	myAnimal.Call("dog")

	//findNum(1, 2, 3, 4, 5, 6, 1)
	//findNum(1, []int{1,2,3,4,5,6}...)
	deferShow()
	deferShow2()
	deferShow3()
	deferShow4()
}

// ...Type为接受可变参数
func findNum(num int, nums ...int) {
	isFind := false // 定义是否查询到的标识符变量
	for k, v := range nums {
		if v == num {
			isFind = true
			fmt.Printf("find nums,key:%d num:%d\n", k, v)
			break
		}
	}
	if !isFind {
		fmt.Println("not find this num")
	}
}

func deferShow4() {
	i := 1
	defer func() {
		i++
		fmt.Printf("defer1 i:%d\n", i)
	}()

	defer func() {
		i++
		fmt.Printf("defer2 i:%d\n", i)
	}()
}

func deferShow3() {
	i := 1
	// 下面为闭包操作
	defer func() { //
		i++
		fmt.Printf("defer i: %d\n", i)
	}()
	i = 2
	fmt.Printf("i:%d\n", i)
}

func deferShow2() {
	i := 1
	defer fmt.Printf("defer i:%d\n", i)
	i++
	fmt.Printf("i:%d\n", i)
}

func deferShow() {
	defer fmt.Println("last show")
	fmt.Println("one")
	fmt.Println("two")
}

/** 闭包演示 **/
func closureShow(i int) func() int {
	return func() int {
		i++
		return i
	}
}
