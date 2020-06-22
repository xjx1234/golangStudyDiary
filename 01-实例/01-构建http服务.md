## HTTP服务

### HTTP Server

对于Golang来说，实现一个简单的`http server`非常容易，只需要短短几行代码。同时有了协程的加持，Go实现的`http server`能够取得非常优秀的性能。这篇文章将会对go标准库`net/http`实现http服务的原理进行较为深入的探究，以此来学习了解网络编程的常见范式以及设计思路。

基于HTTP构建的网络应用包括两个端，即客户端(`Client`)和服务端(`Server`)。两个端的交互行为包括从客户端发出`request`、服务端接受`request`进行处理并返回`response`以及客户端处理`response`。所以http服务器的工作就在于如何接受来自客户端的`request`，并向客户端返回`response`。

#### 创建Http Server

典型的http server 的处理流程可以用下图表示：

![](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/bfd4c9f1f3a224e07b27dd7ad7e15b8a.png)

服务器在接收到请求时，首先会进入路由(`router`)，这是一个`Multiplexer`，路由的工作在于为这个`request`找到对应的处理器(`handler`)，处理器对`request`进行处理，并构建`response`。Golang实现的`http server`同样遵循这样的处理流程。

我们先看看Golang如何实现一个简单的`http server`：

```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!!!"))
}

func main() {
	http.HandleFunc("/", myHandler)
	http.ListenAndServe(":8181", nil)
}
```

运行代码后，在本地浏览器打开 `http://127.0.0.1:8181/` 可以看到页面返回 `Hello world !!!` ，这段代码可以看出，首先利用 `http.HandleFunc` 函数在 根路由 `/` 上注册了一个 `myHandler` 的处理函数，然后利用 `http.ListenAndServe` 开启监听，当 / 路由有请求过来，则根据路由执行对应的 handler 函数，即 `myHandler` 函数。下面我们来看下另一种实现 http server 方法:

```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)
type indexHandler struct {
	content string
}

func (ij *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(ij.content))
}

func main() {
	http.Handle("/", &indexHandler{content: "Hello world !!!"})
	http.ListenAndServe(":8181", nil)
}
```

通过两种创建 Http server 模式比较可以看出来， `http.Handle` 与 `http.HandleFunc` 函数都是用于注册路由，其中他们的区别是第二个参数不一致，`http.HandleFunc` 函数的第二个参数是 `handler func(ResponseWriter, *Request)` 类型的处理函数，最终由 `DefaultServeMux.HandleFunc(pattern, handler)` 完成处理; 而 `http.Handle` 的第二个参数是  `handler Handler` ，是一个具有 Handler类型的结构。而Handler是一个接口，需要实现 `ServeHTTP` 函数, 最终由 `DefaultServeMux.Handle(pattern, handler)` 完成处理。

#### 路由注册源码分析

`http.HandleFunc`和`http.Handle`相关源码如下：

```go
// http.HandleFunc 部分源码
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}

func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	mux.Handle(pattern, HandlerFunc(handler))
}
```

```go
// http.Handle 部分源码
func Handle(pattern string, handler Handler) { DefaultServeMux.Handle(pattern, handler) }
```

```go
var DefaultServeMux = &defaultServeMux
var defaultServeMux ServeMux
```

通过源码分析，我们发现最终 这2种实现方式都是通过 调用  `ServeMux` 中的 `Handle` 方法实现路由的注册。这里我们遇到两种类型的对象：`ServeMux`和`Handler`，我们先说`Handler`。我们继续通过源码进行分析：

我们首先来看看 两种方法的 `Handler`:

- http.Handle 部分 `Handler` 是作为 `func Handle` 的第二个参数  `handler Handler`  出现的，跟踪第二个参数的代码可以发现，此处的 `Handler` 类型为接口，需要实现 `ServeHTTP(ResponseWriter, *Request)` 函数，源码如下:

  ```go
  type Handler interface {
  	ServeHTTP(ResponseWriter, *Request)
  }
  ```

- http.HandleFunc 部分的 `Handler` 是通过 `mux.Handle(pattern, HandlerFunc(handler))` 中的 `HandlerFunc` 类型函数实现的，`HandlerFunc`是一个类型，只不过表示的是一个具有`func(ResponseWriter, *Request)`签名的函数类型，并且这种类型实现了`ServeHTTP`方法（在`ServeHTTP`方法中又调用了自身），也就是说这个类型的函数其实就是一个`Handler`类型的对象。利用这种类型转换，我们可以将一个`handler`函数转换为一个`Handler`对象，而不需要定义一个结构体，再让这个结构实现`ServeHTTP`方法。读者可以体会一下这种技巧。， 源码如下：

  ```go
  type HandlerFunc func(ResponseWriter, *Request)
  
  // ServeHTTP calls f(w, r).
  func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
  	f(w, r)
  }
  ```

分析完成 `Handler` ，我们再来看看另外一个 `ServeMux` ,我们通过源码先看看 `ServeMux` 的结构:

```go
type ServeMux struct {
	mu    sync.RWMutex
	m     map[string]muxEntry
	es    []muxEntry // slice of entries sorted from longest to shortest.
	hosts bool       // whether any patterns contain hostnames
}

type muxEntry struct {
	h       Handler
	pattern string
}
```

这里重点关注`ServeMux`中的字段`m`，这是一个`map`，`key`是路由表达式，`value`是一个`muxEntry`结构，`muxEntry`结构体存储了对应的路由表达式和`handler`。值得注意的是，`ServeMux`也实现了`ServeHTTP`方法：

```go
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(StatusBadRequest)
		return
	}
	h, _ := mux.Handler(r)
	h.ServeHTTP(w, r)
}
```

也就是说`ServeMux`结构体也是`Handler`对象，只不过`ServeMux`的`ServeHTTP`方法不是用来处理具体的`request`和构建`response`，而是用来确定路由注册的`handler`。下面我们再看一下`ServeMux`的`Handle`方法具体做了什么：

```go
func (mux *ServeMux) Handle(pattern string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if pattern == "" {
		panic("http: invalid pattern")
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if _, exist := mux.m[pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}
	e := muxEntry{h: handler, pattern: pattern}
	mux.m[pattern] = e
	if pattern[len(pattern)-1] == '/' {
		mux.es = appendSorted(mux.es, e)
	}

	if pattern[0] != '/' {
		mux.hosts = true
	}
}
```

`Handle`方法主要做了两件事情：一个就是向`ServeMux`的`map[string]muxEntry`增加给定的路由匹配规则；然后如果路由表达式以`'/'`结尾，则将对应的`muxEntry`对象加入到`[]muxEntry`中，按照路由表达式长度排序。前者很好理解，但后者可能不太容易看出来有什么作用，这个问题后面再作分析。

#### 自定义ServeMux

综述以上分析，使用 `http.HandleFunc` 与 `http.Handle` 创建的Http Server 服务都是通过默认的 `DefaultServeMux` 来实现的，下面我们可以通过自己创建的 `ServeMux` 来建立 http server服务，代码如下:

```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

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
	mux := http.NewServeMux()
	mux.Handle("/test1", http.HandlerFunc(selfServerIndex))
	mux.HandleFunc("/test2", selfServerIndex2)
	http.ListenAndServe(":8888", mux)
}
```

#### 开启监听服务源码分析

服务的监听与开启是从 `http.ListenAndServe` 函数开始的，下面查看源码：

```go
func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}

func (srv *Server) ListenAndServe() error {
	if srv.shuttingDown() {
		return ErrServerClosed
	}
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}
```

上述代码中，首先 创建一个 `Server`对象，传入地址和handler参数，然后调用对象的 `ListenAndServe` 方法。通过对比， 我们发现监听以及开启服务，都离不开 `Server` 这个结构体，其源码如下：

```go
type Server struct {
	Addr    string  // TCP address to listen on, ":http" if empty
	Handler Handler // handler to invoke, http.DefaultServeMux if nil
	TLSConfig *tls.Config
	ReadTimeout time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout time.Duration
	IdleTimeout time.Duration
	MaxHeaderBytes int
	TLSNextProto map[string]func(*Server, *tls.Conn, Handler)
	ConnState func(net.Conn, ConnState)
	ErrorLog *log.Logger
	disableKeepAlives int32     // accessed atomically.
	inShutdown        int32     // accessed atomically (non-zero means we're in Shutdown)
	nextProtoOnce     sync.Once // guards setupHTTP2_* init
	nextProtoErr      error     // result of http2.ConfigureServer if used
	mu         sync.Mutex
	listeners  map[*net.Listener]struct{}
	activeConn map[*conn]struct{}
	doneChan   chan struct{}
	onShutdown []func()
}
```

在`Server`的`ListenAndServe`方法中，会初始化监听地址`Addr`，同时调用`Listen`方法设置监听。最后将监听的TCP对象传入`Serve`方法：

```go
func (srv *Server) Serve(l net.Listener) error {
    ...
    baseCtx := context.Background() // base is always background, per Issue 16220
    ctx := context.WithValue(baseCtx, ServerContextKey, srv)
    for {
        rw, e := l.Accept() // 等待新的连接建立
        ...
        c := srv.newConn(rw)
        c.setState(c.rwc, StateNew) // before Serve can return
        go c.serve(ctx) // 创建新的协程处理请求
    }
}
```

这里隐去了一些细节，以便了解`Serve`方法的主要逻辑。首先创建一个上下文对象，然后调用`Listener`的`Accept()`等待新的连接建立；一旦有新的连接建立，则调用`Server`的`newConn()`创建新的连接对象，并将连接的状态标志为`StateNew`，然后开启一个新的`goroutine`处理连接请求。

#### 处理连接源码分析

我们继续探索`conn`的`serve()`方法，这个方法同样很长，我们同样只看关键逻辑。坚持一下，马上就要看见大海了。

```go
func (c *conn) serve(ctx context.Context) {
    ...
    for {
        w, err := c.readRequest(ctx)
        if c.r.remain != c.server.initialReadLimitSize() {
            // If we read any bytes off the wire, we're active.
            c.setState(c.rwc, StateActive)
        }
        ...
        serverHandler{c.server}.ServeHTTP(w, w.req)
        w.cancelCtx()
        if c.hijacked() {
            return
        }
        w.finishRequest()
        if !w.shouldReuseConnection() {
            if w.requestBodyLimitHit || w.closedRequestBodyEarly() {
                c.closeWriteAndWait()
            }
            return
        }
        c.setState(c.rwc, StateIdle) // 请求处理结束后，将连接状态置为空闲
        c.curReq.Store((*response)(nil))// 将当前请求置为空
        ...
    }
}
```

当一个连接建立之后，该连接中所有的请求都将在这个协程中进行处理，直到连接被关闭。在`serve()`方法中会循环调用`readRequest()`方法读取下一个请求进行处理，其中最关键的逻辑就是一行代码：

```go
serverHandler{c.server}.ServeHTTP(w, w.req)
```

进一步解释`serverHandler`：

```go
type serverHandler struct {
	srv *Server
}

func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
	handler := sh.srv.Handler
	if handler == nil {
		handler = DefaultServeMux
	}
	if req.RequestURI == "*" && req.Method == "OPTIONS" {
		handler = globalOptionsHandler{}
	}
	handler.ServeHTTP(rw, req)
}
```

在`serverHandler`的`ServeHTTP()`方法里的`sh.srv.Handler`其实就是我们最初在`http.ListenAndServe()`中传入的`Handler`对象，也就是我们自定义的`ServeMux`对象。如果该`Handler`对象为`nil`，则会使用默认的`DefaultServeMux`。最后调用`ServeMux`的`ServeHTTP()`方法匹配当前路由对应的`handler`方法。

后面的逻辑就相对简单清晰了，主要在于调用`ServeMux`的`match`方法匹配到对应的已注册的路由表达式和`handler`。

```go
func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
	handler := sh.srv.Handler
	if handler == nil {
		handler = DefaultServeMux
	}
	if req.RequestURI == "*" && req.Method == "OPTIONS" {
		handler = globalOptionsHandler{}
	}
	handler.ServeHTTP(rw, req)
}

func (mux *ServeMux) handler(host, path string) (h Handler, pattern string) {
    mux.mu.RLock()
    defer mux.mu.RUnlock()

    // Host-specific pattern takes precedence over generic ones
    if mux.hosts {
        h, pattern = mux.match(host + path)
    }
    if h == nil {
        h, pattern = mux.match(path)
    }
    if h == nil {
        h, pattern = NotFoundHandler(), ""
    }
    return
}

// Find a handler on a handler map given a path string.
// Most-specific (longest) pattern wins.
func (mux *ServeMux) match(path string) (h Handler, pattern string) {
    // Check for exact match first.
    v, ok := mux.m[path]
    if ok {
        return v.h, v.pattern
    }

    // Check for longest valid match.  mux.es contains all patterns
    // that end in / sorted from longest to shortest.
    for _, e := range mux.es {
        if strings.HasPrefix(path, e.pattern) {
            return e.h, e.pattern
        }
    }
    return nil, ""
}
```

在`match`方法里我们看到之前提到的mux的`m`字段(类型为`map[string]muxEntry`)和`es`(类型为`[]muxEntry`)。这个方法里首先会利用进行精确匹配，在`map[string]muxEntry`中查找是否有对应的路由规则存在；如果没有匹配的路由规则，则会利用`es`进行近似匹配。

之前提到在注册路由时会把以`'/'`结尾的路由（可称为**节点路由**）加入到`es`字段的`[]muxEntry`中。对于类似`/path1/path2/path3`这样的路由，如果不能找到精确匹配的路由规则，那么则会去匹配和当前路由最接近的已注册的父节点路由，所以如果路由`/path1/path2/`已注册，那么该路由会被匹配，否则继续匹配下一个父节点路由，直到根路由`/`。

由于`[]muxEntry`中的`muxEntry`按照路由表达式从长到短排序，所以进行近似匹配时匹配到的节点路由一定是已注册父节点路由中最相近的。

至此，Go实现的`http server`的大致原理介绍完毕！

#### 获取参数值方案

- GET 参数获取

  ```go
  // r *http.Request
  r.URL.Query().Get(参数名)
  ```

- POST PUT DELETE 参数获取

  ```go
  // r *http.Request
  r.ParseForm() // 支持 x-www-form-urlencoded 类型数据接收
  r.ParseMultipartForm(128) // 支持 from-data 类型数据接收
  r.Form.Get(参数名)
  r.PostForm.Get(参数名)  //post
  ```

- JSON值获取

  ```go
  // r *http.Request
  type JsonFrom struct {
  	Name string `json:"name"`
  }
  
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
      fmt.Println(jsonData.Name)
  ```

  

### HTTP Client

前面我们已经简单的从源码以及示例角度演示了Http Server方面的操作，既然有了Server端操作，肯定也少不了Client端操作，下面我们看下一个最简单的Http Get请求， 代码如下:

```go
	response, err := http.Get("https://doc.btc.com/v1/poster/production/explorer-banner.json?t=1592798186869")
	fmt.Println(response)
	fmt.Println(err)
```

运行代码，输出接口返回值，看似简单的一段代码，却实现了一个http请求，下面我们从源码层面看看Go是如何处理这些方面的。

```go
func Get(url string) (resp *Response, err error) {
	return DefaultClient.Get(url)
}

func (c *Client) Get(url string) (resp *Response, err error) {
	req, err := NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Do(req *Request) (*Response, error) {
	return c.do(req)
}

func (c *Client) do(req *Request) (retres *Response, reterr error) {
    .......
    req = &Request{
        Method:   redirectMethod,
        Response: resp,
        URL:      u,
        Header:   make(Header),
        Host:     host,
        Cancel:   ireq.Cancel,
        ctx:      ireq.ctx,
    }
    ......
    copyHeaders(req)
    ......
    reqs = append(reqs, req)
    var err error
    var didTimeout func() bool
    if resp, didTimeout, err = c.send(req, deadline); err != nil {
        ......
        return nil, uerr(err)
    }
    ...
    req.closeBody()
	}
}
```

源代码中 ，我们可以看到`http.Get` 首先是通过 一个 调用 `DefaultClient` 中的  `Get` 方法去首先请求，而 `DefaultClient.Get` 则是通过 `NewRequest` 方法构造请求数据，然后将数据通过 `Client` 的对象 `do` 方法发到服务器端。在此我们不细述 `do` 方法的实现，我们看下 `DefaultClient` 和 `NewRequest`的代码。

```go
var DefaultClient = &Client{}

func NewRequest(method, url string, body io.Reader) (*Request, error) {
	if method == "" {
		method = "GET"
	}
	if !validMethod(method) {
		return nil, fmt.Errorf("net/http: invalid method %q", method)
	}
	u, err := parseURL(url) 
	if err != nil {
		return nil, err
	}
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	u.Host = removeEmptyPort(u.Host)
	req := &Request{
		Method:     method,
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(Header),
		Body:       rc,
		Host:       u.Host,
	}
	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			req.ContentLength = int64(v.Len())
			buf := v.Bytes()
			req.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return ioutil.NopCloser(r), nil
			}
		case *bytes.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
		case *strings.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
		default:
		}
		if req.GetBody != nil && req.ContentLength == 0 {
			req.Body = NoBody
			req.GetBody = func() (io.ReadCloser, error) { return NoBody, nil }
		}
	}

	return req, nil
}
```

`DefaultClient` 从代码层面看是一个其实就是一个系统 默认的 `&Client{}` 对象， 而 `NewRequest` 则是将用户传入的方法，URL等参数通过各种判断以及组合，组合成一个 `Request` 对象，返回给用户。这时候我们通过绕过一些默认的构造，去实现一个http Get请求试试：

```go
	Client := &http.Client{}
	request, _ := http.NewRequest("GET", "https://doc.btc.com/v1/poster/production/explorer-banner.json?t=1592798186869", nil)
	response, err := Client.Do(request)
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(content))
	fmt.Println(err)
```

通过上述例子，我们还原了 `http.Get` 的代码逻辑，下面我们尝试下更为复杂点的环境，示例代码如下:

```go
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

```

以上例子，没有进行一些封装和简化操作，只是为了更加直观的了解整个http客户端发送的过程，对增加参数以及头部, cookies信息等进行了列举。