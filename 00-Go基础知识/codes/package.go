/**
* @Author: XJX
* @Description: 包示例演示代码
* @File: package.go
* @Date: 2020/6/9 14:35
 */

package main

import (
	"demo"
	"fmt"
	"time"
	"github.com/labstack/echo"
)

func main() {
	fmt.Println(time.Now())
	data := demo.MyAdd(1,2)
	fmt.Println(data)
	e := echo.New()
	e.Logger.Fatal(e.Start(":1323"))
}
