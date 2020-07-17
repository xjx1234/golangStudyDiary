/**
* @Author: XJX
* @Description: 指针操作
* @File: pointer.go
* @Date: 2020/5/20 17:27
 */

package main

import (
	"fmt"
)

func main()  {

	i := 1
	var p *int // 定义指针方式一
	p2 := new(int) //定义指针方式二
	p = &i
	p2 = &i
	fmt.Println(p)
	fmt.Println(p2)

	a := 5
	pointer1 := &a
	fmt.Println(*pointer1) // 5 指针取值
	*pointer1 = 6 // 通过指针修改 a值
	fmt.Println(a) // 6

	b := 1
	pointer2 := &b
	//pointer2++
	fmt.Println(pointer2)
}

