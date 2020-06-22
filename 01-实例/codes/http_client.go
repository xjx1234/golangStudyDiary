/**
* @Author: XJX
* @Description: http客户端演示代码
* @File: http.go
* @Date: 2020/6/15 12:47
 */

package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
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

func main() {

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

	urlApi := "https://www.okcoin.cn/api/explorer/v1/eth/transfers"
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
	//增加参数
	reqParams := url.Values{}
	if len(params) > 0 {
		for k, v := range params {
			reqParams.Set(k, fmt.Sprint(v))
		}
	}
	api, _ := url.Parse(urlApi)
	api.RawQuery = reqParams.Encode()
	lastApi := api.String()
	Client := &http.Client{}
	//增加 header
	request, _ := http.NewRequest("GET", lastApi, nil)
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
	fmt.Println(err)
}
