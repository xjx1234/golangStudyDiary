## HTTP服务

对于Golang来说，实现一个简单的`http server`非常容易，只需要短短几行代码。同时有了协程的加持，Go实现的`http server`能够取得非常优秀的性能。这篇文章将会对go标准库`net/http`实现http服务的原理进行较为深入的探究，以此来学习了解网络编程的常见范式以及设计思路。

基于HTTP构建的网络应用包括两个端，即客户端(`Client`)和服务端(`Server`)。两个端的交互行为包括从客户端发出`request`、服务端接受`request`进行处理并返回`response`以及客户端处理`response`。所以http服务器的工作就在于如何接受来自客户端的`request`，并向客户端返回`response`。

### Http Server

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

