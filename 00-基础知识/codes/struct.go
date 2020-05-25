/**
* @Author: XJX
* @Description: 结构体示例
* @File: struct.go
* @Date: 2020/5/25 12:05
 */

package main

import "fmt"

/** 学生结构体 */
type Student struct {
	name string //姓名
	age  int    // 年龄
	sex  string //性别
}

type AnonymousField struct {
	string
	age int
}

func main() {

	/** 方案一 */
	var s Student
	s.age = 11
	s.name = "xjx"
	s.sex = "man"
	fmt.Println(s) //{xjx 11 man}

	/** 方案二 返回的是指针类型对象 */
	var s1 = new(Student)
	s1.name = "xjx"
	s1.age = 31
	s1.sex = "man"
	fmt.Println(s1) // &{xjx 31 man}

	/** 方案三 返回的是指针类型对象 */
	s2 := &Student{}
	s2.sex = "man"
	s2.age = 32
	s2.name = "xjx"
	fmt.Println(s2) // &{xjx 32 man}

	/** 结构体成员变量初始化方案1，此方案可选择性初始化相关值*/
	s3 := Student{
		name: "xjx",
		age:  22,
		sex:  "man",
	}

	s3_1 := Student{
		name: "xxx",
	}
	fmt.Println(s3)   // {xjx 22 man}
	fmt.Println(s3_1) // {xxx 0 }
	/** 结构体成员变量初始化方案2, 此方案必须初始化所有定义键的值*/
	s4 := Student{
		"xjx2",
		11,
		"man",
	}
	fmt.Println(s4) // {xjx2 11 man}

	anonymousStruct := struct {
		firstName, lastName string
		age                 int
	}{
		firstName: "xu",
		lastName:  "jianshe",
		age:       11,
	}
	fmt.Println(anonymousStruct) // {xu jianshe 11}
	printInfo(anonymousStruct)

	anonymousFieldData := AnonymousField{
		string: "hello",
		age:    11,
	}
	fmt.Println(anonymousFieldData.string) // hello

	var onePerson Person
	onePerson.Age = 11
	onePerson.name = "xjx"
	onePerson.Address = Address{
		city:  "HangZhou",
		state: "XiHu",
	}
	fmt.Printf("city: %s, state: %s\n", onePerson.city, onePerson.state)
}

type Address struct {
	city, state string
}

type Person struct {
	name      string
	Age       int
	Address
}

func printInfo(info struct {
	firstName, lastName string
	age                 int
}) {
	fmt.Println(info.age)
}
