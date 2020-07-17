/**
* @Author: XJX
* @Description: //接口实现示例
* @File: interface.go
* @Date: 2020/5/26 17:45
 */

package main

import (
	"fmt"
	"github.com/shopspring/decimal"
)

// 购物接口
type Shopping interface {
	Buy(string)
	Pay(string, float64, int) float64
}

// 商品结构体
type Goods struct {
}

//定义结构体
func (g *Goods) Buy(name string) {
	fmt.Printf("buy goods:%s\n", name)
}

func (g *Goods) Pay(name string, price float64, num int) float64 {
	decimal.DivisionPrecision = 2
	cost := decimal.NewFromFloat(price).Mul(decimal.NewFromFloat(float64(num)))
	fmt.Printf("pay goods:%s, cost:%T\n", name, cost)
	payCost, _ := cost.Float64()
	return payCost
}

//动物结构体
type Animals struct {
	name     string //名称
	skill    string //技能
	age      int    //年龄
	classify string //分类
}

//动物属性接口
type AnimalAttribute interface {
	GetName() string //动物名称
}

//动物技巧接口
type AnimalSkill interface {
	IsOverAge() bool //是否超龄
}

//获取动物名称函数
func (a Animals) GetName() string {
	return a.name
}

//获取动物是否超龄函数
func (a Animals) IsOverAge() bool {
	if a.age > 5 {
		return true
	} else {
		return false
	}
}

func findAnimalType(i interface{}) {
	switch i.(type) {
	case AnimalSkill:
		fmt.Println("AnimalSkill")
	case AnimalAttribute:
		fmt.Println("AnimalAttribute")
	default:
		fmt.Println("unkonw")
	}
}

// 定义服务接口 含有启动和日志两个函数
type Service interface {
	Start()
	Log(string)
}

// 定义Logger 结构体
type Logger struct {
}

//定义 GameService结构体
type GameService struct {
	Logger
}

//实现Log函数
func (l Logger) Log(logData string) {
	fmt.Println(logData)
}

//实现Start函数
func (g GameService) Start() {
	fmt.Println("GameService Start!!!")
}

type Desc struct {
	descContent string
}

type Describer interface {
	Describe() *Desc
}

var any interface{} //定义空接口

// 函数中参数使用空接口
func anyDes(i interface{}) {
	fmt.Printf("Type = %T, value = %v\n", i, i)
}

// 类型断言 示例 函数
func findType(i interface{}) {
	switch i.(type) {
	case int:
		fmt.Println("i type is int")
	case string:
		fmt.Println("i type is string")
	case bool:
		fmt.Println("i type is bool")
	default:
		fmt.Println("i unknow type")
	}
}

func main() {

	var x interface{}
	x = 10
	value, ok := x.(int)
	fmt.Println(value, ok)

	findType(1)
	findType([]int{1, 2, 4})

	any = 1
	anyDes("hello")
	anyDes(any)

	var d1 Describer
	if d1 == nil {
		fmt.Println("d1 is nil")
	}

	var myService GameService
	myService.Start()
	myService.Log("hello log")

	//声明一个Shopping接口
	var doShopping Shopping
	//实例化一个Good
	myGood := new(Goods)
	//将接口赋值
	doShopping = myGood
	//使用Shopping接口进行购物
	doShopping.Buy("beef")
	cost := doShopping.Pay("beef", 12.1, 2)
	fmt.Println(cost)

	myAnimal := Animals{
		name:     "dog",
		age:      10,
		skill:    "卖萌",
		classify: "哺乳类",
	}
	findAnimalType("1111")
	findAnimalType(myAnimal)

	var myName AnimalAttribute
	var mySkill AnimalSkill
	myName = myAnimal
	name := myName.GetName()
	fmt.Printf("My name is %s\n", name)
	mySkill = myAnimal
	isover := mySkill.IsOverAge()
	fmt.Printf("my Age is %t\n", isover)

	e := Employee {
		firstName: "Naveen",
		lastName: "Ramanathan",
		basicPay: 5000,
		pf: 200,
		totalLeaves: 30,
		leavesTaken: 5,
	}
	var empOp EmployeeOperations = e
	empOp.DisplaySalary()
	fmt.Println("\nLeaves left =", empOp.CalculateLeavesLeft())


}

type SalaryCalculator interface {
	DisplaySalary()
}

type LeaveCalculator interface {
	CalculateLeavesLeft() int
}

type EmployeeOperations interface {
	SalaryCalculator
	LeaveCalculator
}

type Employee struct {
	firstName   string
	lastName    string
	basicPay    int
	pf          int
	totalLeaves int
	leavesTaken int
}

func (e Employee) DisplaySalary() {
	fmt.Printf("%s %s has salary $%d", e.firstName, e.lastName, (e.basicPay + e.pf))
}

func (e Employee) CalculateLeavesLeft() int {
	return e.totalLeaves - e.leavesTaken
}
