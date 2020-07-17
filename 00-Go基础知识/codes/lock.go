/**
* @Author: XJX
* @Description: 锁示例
* @File: lock.go
* @Date: 2020/6/3 16:24
 */

package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	count     int
	countLock sync.Mutex
	//countLock sync.RWMutex
)

type value struct {
	value int
	mu    sync.Mutex
}

var wgs sync.WaitGroup

func addValue(v1, v2 *value) {
	defer wgs.Done()
	v1.mu.Lock()
	defer v1.mu.Unlock()
	time.Sleep(2 * time.Second)
	v2.mu.Lock()
	defer v2.mu.Unlock()
	fmt.Printf("SUM:%v \n", v1.value+v2.value)
}

//获取count值
func getCount() int {
	//defer countLock.RUnlock()
	defer countLock.Unlock()
	//countLock.RLock()
	countLock.Lock()
	return count
}

//设置count值
func setCount(num int) {
	defer countLock.Unlock()
	countLock.Lock()
	count = num
}

func main() {

	wgs.Add(2)
	var a, b value
	a.value = 5
	b.value = 4
	go addValue(&a, &b)
	go addValue(&b, &a)
	wgs.Wait()

	//setCount(2)
	//fmt.Println(getCount())
}
