/**
* @Author: XJX
* @Description: http客户端演示代码
* @File: http.go
* @Date: 2020/6/15 12:47
 */

package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func getApiKey(baseKey string) string {
	keye := []byte(baseKey)
	first := keye[0:8]
	second := keye[8:]
	mykey := string(second) + string(first)
	t := time.Now().UnixNano() / 1e6
	e := 1*t + 1111111111111
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10)
	r := rand.Intn(10)
	o := rand.Intn(10)
	time := strconv.FormatInt(e, 10) + strconv.Itoa(n) + strconv.Itoa(r) + strconv.Itoa(o)
	apiKey := mykey + "|" + time
	return base64.URLEncoding.EncodeToString([]byte(apiKey))
}

func postFile(filename, apiUrl string) (*http.Response, error) {
	Client := &http.Client{}
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)
	_, err := body_writer.CreateFormFile("uploadfile", filename)
	if err != nil {
		return nil, err
	}
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	boundary := body_writer.Boundary()
	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	request_reader := io.MultiReader(body_buf, fh, close_buf)
	fi, err := fh.Stat()
	if err != nil {
		fmt.Printf("Error Stating file: %s", filename)
		return nil, err
	}
	req, err := http.NewRequest("POST", apiUrl, request_reader)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())
	req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())
	return Client.Do(req)
}

func getFile(apiUrl string) {
	Client := &http.Client{}
	Client.Timeout = time.Second * 60
	request, _ := http.NewRequest("GET", apiUrl, nil)
	resp, err := Client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	out, err := os.Create("D:/wamp64/www/fileTest/1.jpg")
	wt := bufio.NewWriter(out)
	defer out.Close()
	n, err := io.Copy(wt, resp.Body)
	fmt.Println("write", n)
	if err != nil {
		panic(err)
	}
	wt.Flush()
}

//支持 get post格式,支持json格式
func httpDoRequest(apiUrl string, method string, params map[string]interface{}, header map[string]interface{}) (string, error) {
	jsonFlag := false
	httpReq := &http.Request{}
	reader := &bytes.Reader{}
	cType, ok := header["Content-Type"]
	if ok && strings.Contains(fmt.Sprintf("%v", cType), "application/json") {
		jsonFlag = true
	}
	api, _ := url.Parse(apiUrl)
	if len(params) > 0 {
		if jsonFlag {
			bytesData, err := json.Marshal(params)
			if err != nil {
				return "", err
			}
			reader = bytes.NewReader(bytesData)
		} else {
			reqParams := url.Values{}
			for k, v := range params {
				reqParams.Set(k, fmt.Sprint(v))
			}
		}
	}
	if method != "POST"{
		api.RawQuery = reqParams.Encode()
	}
	lastApi := api.String()
	Client := &http.Client{}
	if jsonFlag {
		httpReq, _ = http.NewRequest(method, lastApi, reader)
	} else {
		if method == "POST"{
			body := ioutil.NopCloser(strings.NewReader(reqParams.Encode()))
			httpReq, _ = http.NewRequest(method, apiUrl, body)
		}else{
			httpReq, _ = http.NewRequest(method, lastApi, nil)
		}
	}
	if len(header) > 0 {
		for t, v := range header {
			httpReq.Header.Add(t, fmt.Sprint(v))
		}
	}
	response, _ := Client.Do(httpReq)
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	return string(content), err
}

func main() {
	/*apiUrl := "http://127.0.0.1:8888/upload"
	response, err := postFile("C:/Users/fuzamei/Desktop/2019101319091926888356.jpg", apiUrl)
	fmt.Println(err)
	fmt.Println(response)*/

	/*apiUrl := "http://127.0.0.1:8888/down?file=2020630160250903659.jpg"
	getFile(apiUrl)*/

	/*response, err := http.Get("https://doc.btc.com/v1/poster/production/explorer-banner.json?t=1592798186869")
	fmt.Println(response)
	fmt.Println(err)*/

	/*Client := &http.Client{}
	request, _ := http.NewRequest("GET", "https://doc.btc.com/v1/poster/production/explorer-banner.json?t=1592798186869", nil)
	response, err := Client.Do(request)
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(content))
	fmt.Println(err)*/

	/*urlApi := "https://www.okcoin.cn/api/explorer/v1/eth/transfers"
	var params map[string]interface{}
	var header map[string]interface{}

	params = map[string]interface{}{
		"t":        1592806377542,
		"offset":   0,
		"limit":    20,
		"tranHash": "0x1342bceaee826525e7e8df161cfefab2aac65c982041e0d422948154c505a1e8",
	}

	header = map[string]interface{}{
		"x-apiKey":   getApiKey("a2c903cc-b31e-4547-9299-b6d07b7631ab"),
		"timeout":    10000,
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
	}
	//增加参数 非json格式
	reqParams := url.Values{}
	if len(params) > 0 {
		for k, v := range params {
			reqParams.Set(k, fmt.Sprint(v))
		}
	}
	//如果是json格式
	bytesData, err := json.Marshal(params)
	reader := bytes.NewReader(bytesData)

	api, _ := url.Parse(urlApi)
	api.RawQuery = reqParams.Encode() //json格式可以不需要
	lastApi := api.String()
	Client := &http.Client{}
	//增加 header
	request, _ := http.NewRequest("GET", lastApi, nil)  //非json
	request, _ := http.NewRequest("GET", lastApi, reader) //json
	for t, v := range header {
		request.Header.Add(t, fmt.Sprint(v))
	}
	//增加 cookies
	cookies := &http.Cookie{Name: "OKCOIN"}
	request.AddCookie(cookies)
	response, _ := Client.Do(request)
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(content))
	fmt.Println(err)*/

}
