/**
* @Author: XJX
* @Description: 反射演示用例
* @File: reflect.go
* @Date: 2020/6/8 11:42
 */

package main

import (
	"fmt"
	"reflect"
)

func main() {

	var i int = 5
	valueOfI := reflect.ValueOf(i)
	fmt.Println("i value is ", valueOfI)
	typeOfA := reflect.TypeOf(i)
	fmt.Printf("i type:%v; i kind:%v \n", typeOfA.Name(), typeOfA.Kind())

	type XJX int
	type cat struct {
		name string
	}
	var x XJX = 5
	typofX := reflect.TypeOf(x)
	typeOfCat := reflect.TypeOf(cat{"hx"})
	fmt.Printf("x type:%v, kind:%v \n", typofX.Name(), typofX.Kind())
	fmt.Printf("cat type:%v, kind:%v \n", typeOfCat.Name(), typeOfCat.Kind())

	var a int = 1024
	valueofA := reflect.ValueOf(a)
	var getA1 = valueofA.Interface().(int)
	var getA2 = int32(valueofA.Int())
	fmt.Printf("A:%v; A1:%v, A2:%v \n", valueofA, getA1, getA2)

	catP := &cat{
		name: "Kitty",
	}
	typeOfP := reflect.TypeOf(catP)
	fmt.Printf("cat type:%v, kind:%v\n", typeOfP.Name(), typeOfP.Kind())
	typeOfP = typeOfP.Elem()
	fmt.Printf("cat type:%v, kind:%v\n", typeOfP.Name(), typeOfP.Kind())

	type MyCat struct {
		name  string
		color string `json:"color" id:"black"`
	}
	ins := MyCat{
		name:  "A1",
		color: "black",
	}
	typeOfIns := reflect.TypeOf(ins)
	for i := 0; i < typeOfIns.NumField(); i++ {
		fieldType := typeOfIns.Field(i)
		fmt.Printf("name: %v  tag: '%v'\n", fieldType.Name, fieldType.Tag)
	}
	if catType, ok := typeOfIns.FieldByName("color"); ok {
		fmt.Println(catType.Tag.Get("json"), catType.Tag.Get("id"))
	}

	//定义结构体
	type refletc_test struct {
		a, b string
		c    float64
		d    int
		bool
		int32
		next *refletc_test
	}
	rt := refletc_test{
		a:     "xx",
		b:     "zz",
		d:     11,
		bool:  false,
		int32: 11,
		next:  &refletc_test{},
	}
	valueOfRt := reflect.ValueOf(rt)
	fmt.Println("NumField", valueOfRt.NumField()) // 查询结构体字段数
	boolField := valueOfRt.Field(0) // 查询index为0的对象信息
	fmt.Println("Field", boolField.Type())
	fmt.Println("FieldByName(\"bool\").Type", valueOfRt.FieldByName("bool").Type()) //查询 字段名为 bool的成员对象信息
	fmt.Println("FieldByIndex([]int{6, 0}).Type()", valueOfRt.FieldByIndex([]int{6, 0}).Type()) // 多层访问，查询第6个结构体的第一个成员Type

	var y int = 2
	valueOfY := reflect.ValueOf(&y)
	valueOfY = valueOfY.Elem()
	if valueOfY.CanAddr(){
		fmt.Println("y 可以被修改")
		valueOfY.SetInt(3)
		fmt.Printf("y 值为:%v \n", valueOfY.Int())
	}else{
		fmt.Println("y 不能被修改")
	}

	//定义结构体
	type human struct {
		name string
		Age int
	}
	valueOfHuman := reflect.ValueOf(&human{"xjx", 24})
	valueOfHuman = valueOfHuman.Elem()
	valueAge := valueOfHuman.FieldByName("Age")
	valueAge.SetInt(26)
	fmt.Println(valueAge.Int())

	var b int
	// 取变量a的反射类型对象
	typeOfB := reflect.TypeOf(b)
	// 根据反射类型对象创建类型实例
	aIns := reflect.New(typeOfB)
	// 输出Value的类型和种类
	fmt.Println(aIns.Type(), aIns.Kind())


	funcValue := reflect.ValueOf(addNum)
	// 构造函数参数, 传入两个整型值
	paramList := []reflect.Value{reflect.ValueOf(10), reflect.ValueOf(20)}
	// 反射调用函数
	retList := funcValue.Call(paramList)
	// 获取第一个返回值, 取整数值
	fmt.Println(retList[0].Int())

}

func addNum(a, b int) int {
	return a + b
}
