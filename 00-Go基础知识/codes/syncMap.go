/**
* @Author: XJX
* @Description:
* @File: syncMap.go
* @Date: 2020/6/9 18:14
 */

package main

import (
	"fmt"
	"sync"
)

var w sync.WaitGroup

/*type MyData struct {
	lock sync.RWMutex
	data map[int]int
}

func writeData(d *MyData) {
	defer d.lock.Unlock()
	defer w.Done()
	d.lock.Lock()
	for i := 0; i < 5000; i++ {
		d.data[i] = i
	}
}

func getData(i int, d *MyData) int {
	defer w.Done()
	for {
		d.lock.RLock()
		if result, ok := d.data[i]; ok {
			d.lock.RUnlock()
			fmt.Println(result)
			return result
		} else {
			d.lock.RUnlock()
			//time.Sleep(time.Microsecond * 10)
			println("waiting....")
		}
	}
}
*/

func writeData(d *sync.Map) {
	defer w.Done()
	for i := 0; i < 5000; i++ {
		d.Store(i, i)
		fmt.Printf("write i:%d \n", i)
	}
}

func getData(i int, d *sync.Map) interface{} {
	defer w.Done()
	for {
		if v, ok := d.Load(i); ok {
			fmt.Printf("v is %v \n", v)
			return v
		} else {
			fmt.Printf("waiting...,statu:%v \n", ok)
		}
	}
}

var myMapData sync.Map

func main() {
	w.Add(1)
	go writeData(&myMapData)
	go getData(1, &myMapData)
	w.Wait()
	/*
		var myMapData = &MyData{data: map[int]int{}}
		w.Add(2)
		go writeData(myMapData)
		go getData(4555, myMapData)
		w.Wait()*/
}
