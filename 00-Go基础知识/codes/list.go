/**
* @Author: XJX
* @Description:
* @File: list.go
* @Date: 2020/5/21 17:05
 */

package main

import (
	"container/list"
	"fmt"
)

func main() {

	var list1 list.List // list声明方式1
	list2 := list.New() // list声明方式2
	list3 := list.New()

	list1.PushBack("hello") // 添加列表元素到尾部
	list1.PushFront("say")  // 添加列表元素到头部
	element := list1.PushBack("one")
	list1.InsertAfter("fff", element)  //在element点后插入元素
	list1.InsertBefore("ttt", element) //在element点前插入元素

	list2.PushBack("list2")
	list1.PushBackList(list2) //在列表list1后插入列表list2
	list3.PushBack("three")
	list1.PushFrontList(list3) // 在列表list1前插入列表list3

	list1.Remove(element) // 移除element位置

	// 遍历列表
	for i := list1.Front(); i != nil; i = i.Next() {
		fmt.Println(i.Value)
	}

}
