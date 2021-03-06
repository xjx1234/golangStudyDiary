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

	familyData := &FamilyMember{
		name: "xjx",
		family: Family{
			fname: "xushijz",
			num:   5,
		},
		address: struct{ city, state string }{city: "Hangzhou", state: "Xihu"},
		child: &FamilyMember{
			name: "xzh",
			family: Family{
				fname: "xushijz",
				num:   5,
			},
			address: struct{ city, state string }{city: "Hangzhou", state: "Xihu"},
			child:   &FamilyMember{},
		},
	}
	fmt.Println(familyData) //&{xjx {xushijz 5} {Hangzhou Xihu} 0xc000086050}

	property := new(Property)
	fmt.Println(property.value) // 0
	property.SetValue(2)
	fmt.Println(property.value) // 2
	fmt.Println(property.GetValue()) //2

	myAddress := Address{
		city: "HangZhou",
		state: "XiHu Area",
	}
	myAddress.ChangeAddress("QuZhou", "QuJiang")
	fmt.Printf("city:" + myAddress.city + " state:" + myAddress.state + "\n")

}








//定义属性结构体
type Property struct {
	value int
}

//设置属性值
func (p *Property) SetValue(v int) {
	p.value = v //修改成员变量值
}

//获取属性值
func (p *Property) GetValue() int {
	return p.value
}

/** 家族信息结构 */
type Family struct {
	fname string //家族名称
	num   int    //家族人数
}

/** 家族成员结构体 */
type FamilyMember struct {
	name    string
	family  Family
	address struct {
		city, state string
	}
	child *FamilyMember
}

type Address struct {
	city, state string
}

func (a Address) ChangeAddress(city, state string){
	a.city = city
	a.state = state
}

type Person struct {
	name string
	Age  int
	Address
}

func printInfo(info struct {
	firstName, lastName string
	age                 int
}) {
	fmt.Println(info.age)
}
