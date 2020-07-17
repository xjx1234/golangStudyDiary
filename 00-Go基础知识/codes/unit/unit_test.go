/**
* @Author: XJX
* @Description: 单元测试用例
* @File: unit_test.go
* @Date: 2020/6/11 16:20
 */

package main

import (
	"fmt"
	"testing"
)

func TestGetArea(t *testing.T) {
	m := getArea(4, 5)
	if m != 20 {
		t.Error("error")
	}else{
		fmt.Println("ok")
	}
}

func BenchmarkGetArea(t *testing.B) {
	for i := 0; i < t.N; i++ {
		getArea(4, 5)
	}
}
