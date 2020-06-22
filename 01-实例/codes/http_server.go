/**
* @Author: XJX
* @Description: http服务端演示代码
* @File: http.go
* @Date: 2020/6/15 12:47
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	name := ParseHttpData(r, "POST", "name")
	w.Write([]byte("Hello " + name.(string) + " !!!"))
}

type indexHandler struct {
	content string
}

type JsonFrom struct {
	Name string `json:"name"`
}

func ParseHttpData(r *http.Request, method string, sname string) interface{} {
	typeContent := r.Header.Get("content-type")
	if strings.Contains(typeContent, "application/json") {
		var jsonData JsonFrom
		bodyByte, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		if err := json.Unmarshal(bodyByte, &jsonData); err != nil {
			fmt.Println(err)
			return ""
		}
		valueJsonData := reflect.ValueOf(&jsonData)
		valueJsonData = valueJsonData.Elem()
		return valueJsonData.FieldByName(sname).String()
	} else if method == "GET" {
		return r.URL.Query().Get(sname)
	} else {
		r.ParseForm()
		r.ParseMultipartForm(128)
		return r.PostForm.Get(sname)
	}
	return ""
}

func (ij *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := ParseHttpData(r, "POST", "Name")
	w.Write([]byte("Hello " + name.(string) + " !!!"))
}

func selfServerIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello selfServerIndex!!!!"))
}

func selfServerIndex2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `<!doctype html>
    <META http-equiv="Content-Type" content="text/html" charset="utf-8">
    <html lang="zh-CN">
            <head>
                    <title>selfServerIndex2</title>
            </head>
            <body>
                <div id="app">selfServerIndex2!</div>
            </body>
    </html>`
	w.Write([]byte(html))
}

func main() {

	//http.HandleFunc("/", myHandler)
	http.Handle("/", &indexHandler{content: ""})
	http.ListenAndServe(":8181", nil)

	/*mux := http.NewServeMux()
	mux.Handle("/test1", http.HandlerFunc(selfServerIndex))
	mux.HandleFunc("/test2", selfServerIndex2)
	http.ListenAndServe(":8888", mux)*/

}
