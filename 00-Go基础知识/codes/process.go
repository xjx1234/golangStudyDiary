/**
* @Author: XJX
* @Description: 基本流程控制demo代码
* @File: process.go
* @Date: 2020/5/22 11:11
 */

package main

import "fmt"

func main() {

	i := 10
	if i > 0 {
		fmt.Println("i值大于0\n")
	} else if i > 5 {
		fmt.Println("i值大于5\n")
	} else if i > 10 {
		fmt.Println("i值大于10\n")
	} else {
		fmt.Println("i值小于等于0")
	}

	if num := 10; num%2 == 0 {
		fmt.Println(num, "is even")
	} else {
		fmt.Println(num, "is odd")
	}

	//示例一
	for a := 1; a < 10; a++ {
		fmt.Println(a)
	}

	// 示例二
	var j int
	for ; ; j++ {
		if j > 5 {
			break
		}
	}

	sum := 0
	for {
		sum++
		if sum > 2 {
			break
		}
	}

	x := 5
	for x < 5 {
		fmt.Println("x 值小于5")
	}

	sliceData := []int{1, 2, 3, 4, 5}
	for k, v := range sliceData {
		fmt.Printf("k: %d v: %d \n", k, v)
	}

	for num := 0; num < 20; num++ {
		if num == 10 {
			break
		}
		if num == 5 {
			continue
		}
		fmt.Println(num)
	}

	inputString := "hello"
	switch inputString {
	case "hello":
		fmt.Println("hello")
	case "world":
		fmt.Println("world")
	case "test":
		fmt.Println("test")
	default:
		fmt.Println("hi")
	}

	var a = "mum"
	switch a {
	case "mum", "daddy":
		fmt.Println("family")
	}

	var r int = 11
	switch {
	case r > 10 && r < 20:
		fmt.Println(r)
	}

	var s = "hello"
	switch {
	case s == "hello":
		fmt.Println("hello")
		fallthrough
	case s != "world":
		fmt.Println("world")
	}

	for i := 1; i < 10; i++ {
		for y := 1; y < 10; y++ {
			if y == 5 {
				goto HERE
			}
		}
	}
	HERE:
		fmt.Println("here")

}
