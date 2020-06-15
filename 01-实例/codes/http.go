/**
* @Author: XJX
* @Description: http服务演示代码
* @File: http.go
* @Date: 2020/6/15 12:47
 */

package main

import (
	"fmt"
	"net/http"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world!")
}

type indexHandler struct {
	content string
}

func (ij *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, ij.content)
}

func main() {
	//http.HandleFunc("/", myHandler)
	http.Handle("/", &indexHandler{content: "hello world!!"})
	http.ListenAndServe(":8181", nil)
}
