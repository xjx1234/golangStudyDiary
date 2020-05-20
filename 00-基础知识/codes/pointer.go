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

type Person struct {
	age  int
	name string
}

func (p Person) sayHi() {
	fmt.Printf("sayHi -- This is %s, my age is %d\n", p.name, p.age)
}

func (p Person) modifyAge(age int) {
	fmt.Printf("modifyAge")
	p.age = age
}

func (p Person) modifyAge2(age int) Person{
	p.age = age
	return p
}

func main() {
	person := Person{20, "XJX"}
	fmt.Printf("person <%s:%d>\n", person.name, person.age)
	person.sayHi()
	person.modifyAge(22)
	person.sayHi()
	newPerson := person.modifyAge2(23)
	fmt.Printf("newPerson <%s:%d>\n", newPerson.name, newPerson.age)
}
