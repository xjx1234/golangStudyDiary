/**
* @Author: XJX
* @Description:
* @File: flag.go
* @Date: 2020/6/10 15:04
 */

package main

import (
	"flag"
	"fmt"
	"strings"
)

var name = flag.String("name", "", "姓名")
var age = flag.Int("age", 0, "年龄")
var mrs bool

func Init() {
	flag.BoolVar(&mrs, "mrs", false, "婚否")
}

type sliceValue []string

func newSliceValue(vals []string, p *[]string) *sliceValue{
	*p = vals
	return (*sliceValue)(p)
}

func (s *sliceValue) String() string{
	*s = sliceValue(strings.Split("default", ","))
	return "none of my business"
}

func (s *sliceValue) Set(val string) error{
	*s = sliceValue(strings.Split(val, ","))
	return nil
}



func main() {
	//Init()
	var languages []string
	flag.Var(newSliceValue([]string{}, &languages), "slice", "i like")
	flag.Parse()
	fmt.Println(languages)
	/*fmt.Printf("args=%s, num=%d\n", flag.Args(), flag.NArg())
	for i := 0; i != flag.NArg(); i++ {
		fmt.Printf("arg[%d]=%s\n", i, flag.Arg(i))
	}
	fmt.Println("name=", *name)
	fmt.Println("age=", *age)
	fmt.Println("mrs=", mrs)
	/*flag.Int("age", 11, "年龄")
	flag.String("name", "xjx", "姓名")
	flag.Bool("mrs", true, "婚否")

	var name string
	var age int
	var mrs bool
	flag.StringVar(&name, "name", "xjx", "姓名")
	flag.IntVar(&age, "age", 0, "年龄")
	flag.BoolVar(&mrs, "mrs", false, "婚否")*/

}
