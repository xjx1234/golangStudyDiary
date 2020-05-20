/**
* @Author: XJX
* @Description: 数组使用示例
* @File: array.go
* @Date: 2020/5/20 18:18
 */

package main

import "fmt"

func main()  {
	var a [3]int  //定义三个整数的数组
	var b [3]int = [3]int{1,2,3} //定义数组并初始化
	c := [...]string{"hello", "word", "ok"} //在数组的定义中，如果在数组长度的位置出现“...”省略号，则表示数组的长度是根据初始化值的个数来计算
	fmt.Println(len(b)) // len函数获取数组长度
	fmt.Println(a[0])
	fmt.Println(c[0])

	// range 可以用来遍历数组
	for _,v := range c{
		fmt.Println(v)
	}

	var x [2][4] int = [2][4]int{{1,2,3,4}, {4,5,6,7}} //多维数组定义并初始化
	fmt.Println(x[0][1])
	fmt.Println(x[0][3])

}