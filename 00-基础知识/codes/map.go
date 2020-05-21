/**
* @Author: XJX
* @Description: map示例
* @File: map.go
* @Date: 2020/5/21 15:56
 */

package main

import "fmt"

func main() {
	var mapList map[int]int                                            //声明map
	mapList = map[int]int{1: 1, 2: 2}                                  //赋值map
	fmt.Println(mapList)                                               // map[1:1 2:2]
	var mapList1 map[int]string = map[int]string{1: "hello", 2: "xjx"} //声明并赋值map
	fmt.Println(mapList1)                                              //map[1:hello 2:xjx]
	mapList2 := make(map[int][]int)                                    //使用make创建map
	slice1 := []int{1, 2, 3, 4}
	mapList2 = map[int][]int{1: slice1}                                                 //将切片作为值传值给map
	fmt.Println(mapList2)                                                               // map[1:[1 2 3 4]]
	mapList3 := map[string]string{"one": "oneData", "two": "twoData", "three": "hello"} //声明并赋值map
	mapList3["one"] = "oneone"                                                          // 修改map值
	fmt.Println(mapList3)                                                               // map[one:oneone three:hello two:twoData]

	mapNew := map[int]string{1: "one", 2: "two"}
	for k, v := range mapNew {
		fmt.Printf("k:%d v:%s\n", k, v)
	}

	delMap := map[string]string{"one":"one", "two":"two"}
	delete(delMap, "one") // 删除one键值的元素
	fmt.Println(delMap) //map[two:two]
	delMap = make(map[string]string) // 由于go并没有清空map的函数，所以只能用重新创建一个map来覆盖之前map方案处理
	fmt.Println(delMap) // map[]

	animalMap := map[string]string{
		"dog" : "狗",
		"cat" : "猫",
		"pig" : "猪",
		"duck" : "鸭",
	}
	newAnimalMap := animalMap
	delete(newAnimalMap, "pig")
	fmt.Println(animalMap) // map[cat:猫 dog:狗 duck:鸭]




}
