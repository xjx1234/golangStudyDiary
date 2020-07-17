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
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"
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

func Ext(fileName string) string {
	return path.Ext(fileName)
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(5)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	extName := Ext(handler.Filename)
	if _, ok := allowExt[extName]; !ok {
		w.Write([]byte(extName + " is not allow ext"))
		return
	}
	if int(handler.Size) > allowMaxSize {
		w.Write([]byte(handler.Filename + " Size bigger +" + string(allowMaxSize)))
		return
	}
	fileexist := checkFileIsExist(filePath)
	if !fileexist {
		err1 := os.Mkdir(filePath, os.ModePerm)
		if err1 != nil {
			w.Write([]byte(err1.Error()))
			return
		}
	}
	rand.Seed(time.Now().UnixNano())
	randInt := strconv.Itoa(rand.Intn(10000000))
	newFile := time.Now().Format("200612150405") + string(randInt) + extName
	f, err := os.OpenFile(filePath+"/"+newFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer f.Close()
	defer file.Close()
	io.Copy(f, file)
	w.Write([]byte("file_upload succ!!!!"))
}

func getContentType(fileName string) (extension, contentType string) {
	arr := strings.Split(fileName, ".")
	if len(arr) >= 2 {
		extension = arr[len(arr)-1]
		switch extension {
		case "jpeg", "jpe", "jpg":
			contentType = "image/jpeg"
		case "png":
			contentType = "image/png"
		case "gif":
			contentType = "image/gif"
		case "mp4":
			contentType = "video/mpeg4"
		case "mp3":
			contentType = "audio/mp3"
		case "wav":
			contentType = "audio/wav"
		case "pdf":
			contentType = "application/pdf"
		case "doc", "":
			contentType = "application/msword"
		}
	}
	contentType = "application/octet-stream"
	return
}

func downFile(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/favicon.ico" {
		return
	}
	fileName := r.URL.Query().Get("file")
	if len(fileName) == 0 {
		w.Write([]byte("file is empty!!!"))
		return
	}
	fileName, err := url.QueryUnescape(fileName)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	f, err := os.Open(filePath + "/" + fileName)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	info, err := f.Stat()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	_, contentType := getContentType(fileName)
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
	f.Seek(0, 0)
	io.Copy(w, f)
}

var allowExt map[string]string = map[string]string{
	".jpg": ".jpg",
}
var allowMaxSize int = 1024 * 1024 * 5
var filePath string = "D:/wamp64/www/fileTest"

func main() {

	//http.HandleFunc("/", myHandler)
	//http.Handle("/", &indexHandler{content: ""})
	//http.ListenAndServe(":8181", nil)

	mux := http.NewServeMux()
	//mux.Handle("/test1", http.HandlerFunc(selfServerIndex))
	//mux.HandleFunc("/test2", selfServerIndex2)
	mux.HandleFunc("/upload", uploadFile)
	mux.HandleFunc("/down", downFile)
	http.ListenAndServe(":8888", mux)

}
