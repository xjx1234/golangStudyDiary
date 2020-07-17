/**
* @Author: XJX
* @Description: 切片演示
* @File: slice.go
* @Date: 2020/5/21 11:01
 */

package main

import "fmt"

func main() {
	a := [5]int{1, 2, 3, 4, 5}
	slice_one := a[0:2]    //从数组或者切片中生成切片
	fmt.Println(slice_one) // [1,2]
	fmt.Println(a[:])      //输出原切片
	fmt.Println(a[2:])     //从第三个值开始到最后一个值
	fmt.Println(a[0:0])    //清空切片

	var strList []string                        //声明切片
	var numListEmpty = []int{}                  //声明并且初始化一个空切片
	var b = []int{1, 2, 3}                      //简化声明并初始化切片
	var c []string = []string{"hello", "world"} //声明并初始化切片
	fmt.Println(b)
	fmt.Println(strList)
	fmt.Println(numListEmpty)
	fmt.Println(c)

	makeSlice := make([]int, 5, 5) //使用make创建一个切片
	fmt.Println(makeSlice)
	makeSlice = []int{1, 2, 3, 4, 5} //赋值切片
	fmt.Println(makeSlice)

	modifySlice := []string{"hello", "china", "!"}
	for k, v := range modifySlice {
		if v == "!" {
			modifySlice[k] = "!!!"
		}
	}
	fmt.Println(modifySlice) //[hello china !!!]

	var addSlice []int
	addSlice = append(addSlice, 1)                       //追加1个元素
	fmt.Println(addSlice)                                // [1]
	addSlice = append(addSlice, 2, 3, 4)                 //追加3个元素
	fmt.Println(addSlice)                                // [1 2 3 4]
	addSlice = append(addSlice, []int{5, 6, 7, 8, 9}...) // 追加一个切片，切片需要解包
	fmt.Println(addSlice)                                //[1 2 3 4 5 6 7 8 9]

	var headAddSlice []int
	headAddSlice = append([]int{0, 1, 2}, headAddSlice...) // 将headAddSlice切片加到其他切片尾部，巧妙的变幻为头部添加功能，注意append第一个参数必须为切片
	fmt.Println(headAddSlice)                              // [0 1 2]
	headAddSlice = append([]int{3, 4}, headAddSlice...)
	fmt.Println(headAddSlice) // [3 4 0 1 2]

	var numbers []int
	for i := 0; i < 10; i++ {
		numbers = append(numbers, i)
		fmt.Printf("len:%d  cap:%d pointer:%p\n", len(numbers), cap(numbers), numbers)
	}

	slice1 := []int{1, 2, 3, 4, 5, 6}
	slice2 := []int{7, 8, 9, 10}
	slice3 := []int{1, 2, 3, 4, 5, 6, 7}
	slice4 := []int{8, 9, 10}
	copy(slice2, slice1)
	fmt.Println(slice2) // [1 2 3]
	copy(slice3, slice4)
	fmt.Println(slice3) // [8 9 10 4 5 6 7]

	// 从开头位置删除
	delSlice := []int{1, 2, 3, 4, 5, 6, 7, 8}
	delSlice = delSlice[1:] // 删除第一个元素
	fmt.Println(delSlice)   // [2 3 4 5 6 7 8]
	delSlice2 := []int{1, 2, 3, 4}
	delSlice2 = append(delSlice2[0:0], delSlice2[1:]...) //删除第一个元素的另外一种实现
	fmt.Println(delSlice2)                               //[2 3 4]
	delSlice3 := []int{1, 2, 3, 4, 5, 6, 7}
	delSlice3 = delSlice3[:copy(delSlice3, delSlice3[1:])] //使用copy方式删除一个元素的实现

	//从尾部删除
	delSlice4 := []int{1, 2, 3, 4, 5, 6, 7}
	delSlice4 = delSlice4[:len(delSlice4)-1] //删除尾部第一个元素

	//从中间删除
	delSlice5 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	delSlice5 = append(delSlice5[:4], delSlice5[5:]...) // 删除第四个元素
	fmt.Println(delSlice5)                              //[1 2 3 4 6 7 8 9 10]

	var moreSlice [][]string = [][]string{
		{"C", "C++"},
		{"GO", "RUST"},
		{"PHP"},
	}
	fmt.Println(moreSlice) // [[C C++] [GO RUST] [PHP]]
}
