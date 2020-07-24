## websocket实现

### 常见服务器推技术

​		目前传统模式的 Web 系统以客户端发出请求、服务器端响应的方式工作。大部分是基于传统的http协议实现的web服务，但这种模式存在着一些缺陷，比如满足不了实时交互需求，比较耗损服务器性能，传输速度慢等。造成这些缺陷主要是http协议通信只能由客户端发起，做不到服务器主动向客户端推送信息。于是基于这种现象，产生了很多服务器推技术，包括WebSocket，后面我们来解决了目前主流的一些服务器堆技术。

​	常见的服务器推技术有以下几种：

1.  AJAX 轮询

   前端用定时器，每间隔一段时间发送请求来获取数据是否更新，这种方式可兼容ie和支持高级浏览器。通常采取setInterval或者setTimeout实现。通过递归的方法，在获取到数据后每隔一定时间再次发送请求，这样虽然无法保证两次请求间隔为指定时间，但是获取的数据顺序得到保证。

   虽然相对传统web模式，ajax轮询大大改进了用户体验以及数据实时性，但也存在着 页面假死，无畏的网络传输的缺陷。在大数据量以及对实时要求较高的情况下，AJAX还是存在一定弊端。

   ![](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/ajax.jpg)

   主要应用： 常见的对实时要求低，用户量少的 web数据更新

2.  AJAX长轮询 （long-polling）

   客户端像传统轮询一样从服务端请求数据，服务端会阻塞请求不会立刻返回，直到有数据或超时才返回给客户端，然后关闭连接，客户端处理完响应信息后再向服务器发送新的请求。

   ![](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/ajax-poll.jpg)

   AJAX长轮询解决了传统AJAX频繁访问网络请求浪费服务器资源的问题，但存在着长期占有连接，丧失了无状态高并发特点 。

   主要应用： 早期股票系统 早期实时通讯

3. Flash XML socket

   这种方案实现的基础是：

   一、Flash提供了 XMLSocket类。

   二、 JavaScript 和 Flash的紧密结合：在 JavaScript可以直接调用 Flash程序提供的接口。

   但是也存在下面缺陷：

   a) 因为XMLSocket没有HTTP隧道功能，XMLSocket类不能自动穿过防火墙；

   b) 因为是使用套接口，需要设置一个通信端口，防火墙、代理服务器也可能对非HTTP通道端口进行限制；

   主要应用： 网络聊天室，网络互动游戏。

4.  Server-sent Events（sse）

   sse与长轮询机制类似，区别是每个连接不只发送一个消息。客户端发送一个请求，服务端保持这个连接直到有新消息发送回客户端，仍然保持着连接，这样连接就可以消息的再次发送，由服务器单向发送给客户端。

   SSE本质是发送的不是一次性的数据包，而是一个数据流。可以使用 HTTP 301 和 307 重定向与正常的 HTTP 请求一样。服务端连续不断的发送，客户端不会关闭连接，如果连接断开，浏览器会尝试重新连接。如果连接被关闭，客户端可以被告知使用 HTTP 204 无内容响应代码停止重新连接。

   但是 sse只适用于高级浏览器，ie不支持。因为ie上的XMLHttpRequest对象不支持获取部分的响应内容，只有在响应完成之后才能获取其内容。

5. WebSocket

   WebSocket是一种全新的协议，随着HTML5草案的不断完善，越来越多的现代浏览器开始全面支持WebSocket技术了，它将TCP的Socket（套接字）应用在了webpage上，从而使通信双方建立起一个保持在活动状态连接通道。具体技术细节下面章节会细细说明。

### Http &Socket & WebSocket的关联和区别

**Http协议：**简单的对象访问协议，对应于应用层。Http协议是基于TCP链接的。http连接就是所谓的短连接，及客户端向服务器发送一次请求，服务器端相应后连接即会断掉。

**Socket:**  是对TCP/IP协议的封装，Socket本身并不是协议，而是一个调用接口（API），通过Socket，我们才能使用TCP/IP协议。socket连接及时所谓的长连接，理论上客户端和服务端一旦建立连接，则不会主动断掉；但是由于各种环境因素可能会是连接断开，比如说：服务器端或客户端主机down了，网络故障，或者两者之间长时间没有数据传输，网络防火墙可能会断开该链接已释放网络资源。所以当一个socket连接中没有数据的传输，那么为了位置连续的连接需要发送心跳消息，具体心跳消息格式是开发者自己定义的。

**WebSocket协议: **是基于TCP的一种新的网络协议，和http协议一样属于应用层协议，它实现了浏览器与服务器全双工(full-duplex)通信，也就是允许服务器主动发送信息给客户端。我在实现二维码扫描登录时曾使用过，有了它就不需要通过轮询或者建立长连接的方式来使客户端实时获取扫码状态，因为当扫码后，服务器端可以主动发送消息通知客户端。

**webSocket和http的区别 : ** http链接分为短链接和长链接，短链接是每次请求都要重新建立TCP链接，TCP又要三次握手才能建立，然后发送自己的信息。即每一个request对应一个response。长链接是在一定的期限内保持TCP连接不断开。客户端与服务器通信，必须要由客户端发起然后服务器返回结果。客户端是主动的，服务器是被动的。简单的说，WebSocket协议之前，双工通信是通过多个http链接轮询来实现的，这导致了效率低下。WebSocket解决了这个问题，他实现了多路复用，他是全双工通信。在webSocket协议下客服端和浏览器可以同时发送信息。建立了WebSocket之后服务器不必在浏览器发送request请求之后才能发送信息到浏览器。这时的服务器已有主动权想什么时候发就可以什么时候发送信息到客户端，而且信息当中不必再带有head的部分信息了。与http的长链接通信来比，这种方式不仅能降低服务器的压力，而且信息当中也减少了部分多余的信息。

**webSocket和socket的区别:** 就像Java和JavaScript，并没有什么太大的关系，但又不能说完全没关系。可以这么说：命名方面，Socket是一个深入人心的概念，WebSocket借用了这一概念；使用方面，完全两个东西。总之，可以把WebSocket想象成HTTP，HTTP和Socket什么关系，WebSocket和Socket就是什么关系。

### WebSocket服务端实现

​		首先出于兼容性的考虑，WS的握手使用HTTP来实现（此文档中提到未来有可能会使用专用的端口和方法来实现握手），客户端的握手消息就是一个「普通的，带有Upgrade头的，HTTP Request消息」。所以这一个小节到内容大部分都来自于RFC2616，这里只是它的一种应用形式，下面是RFC6455文档中给出的一个客户端握手消息示例：

> HTTP/1.1 101 Switching Protocols 
>
> Upgrade: websocket 
>
> Connection: Upgrade 
>
> Sec-WebSocket-Accept: FavWZqYJf1A1qOejq/JygSoHKDM=

可以看到，前两行跟HTTP的Request的起始行一模一样，而真正在WS的握手过程中起到作用的是下面几个header域。Upgrade：upgrade是HTTP1.1中用于定义转换协议的header域。它表示，如果服务器支持的话，客户端希望使用现有的「网络层」已经建立好的这个「连接（此处是TCP连接）」，切换到另外一个「应用层」（此处是WebSocket）协议。

Connection：HTTP1.1中规定Upgrade只能应用在「直接连接」中，所以带有Upgrade头的HTTP1.1消息必须含有Connection头，因为Connection头的意义就是，任何接收到此消息的人（往往是代理服务器）都要在转发此消息之前处理掉Connection中指定的域（不转发Upgrade域）。

介绍完基本的原理，下面我们讲解代码的部分，此文中实现WebSocket服务我们采用了 `github.com/gorilla/websocket` 库包，代码用了网上的一些例子，首先实现WebSocket我们先升级Http来实现，首先要定义升级部分的参数：

```go
//升级websocket参数配置
var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 10 * time.Second,
	// 控制跨域，可以根据需要写入自己需要的代码
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
```

定义完升级部分参数，我们可以想想，WebSocket连接肯定需要几个必要部分，第一个是 消息， 第二个是连接，我们接着定义这2部分：

```go
var (
	maxConnId int64 //最大ID
)

//连接结构体
type WsConn struct {
	Id         int64  //ID
	WsSocket   *websocket.Conn //连接
	mutex      sync.Mutex //锁
	InChan     chan *WsMessage //收到消息通道
	OutChan    chan *WsMessage //发出信息通道
	ClosedChan chan byte // 关闭连接通道
	isClosed   bool //是否关闭
}

//消息结构体
type WsMessage struct {
	Type int
	Data []byte
}
```

定义完以上部分，就进入到逻辑部分，websocket服务需要处理的部分有 获取信息， 发出信息， 关闭连接，处理初始化连接几块，先从初始化连接部分开始，也就是用户连接上来后的处理部分：

```go
//处理请求
func wsHandler(w http.ResponseWriter, r *http.Request) {
	wsSocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket升级失败" + err.Error())
		return
	}
	maxConnId++
	conn := &WsConn{
		Id:       maxConnId,
		WsSocket: wsSocket,
		InChan:   make(chan *WsMessage, 1000),
		OutChan:  make(chan *WsMessage, 1000),
		ClosedChan: make(chan byte),
		isClosed: false,
	}
	fmt.Println(conn)
	log.Println("total online number:", maxConnId-1)
	go conn.procLoop() //处理心跳 处理通道中的数据读与写
	go conn.wsReadLoop() // websocket收到的数据处理 
	go conn.wsWriteLoop() // 负责websocket发出的数据处理 
}
```

procLoop代码：

```go
func (wsConn *WsConn) procLoop() {

	//心跳机制，为了处理无效连接，这机制大部分场景不需要实现，只是为了演示去实现
	go func() {
		for {
			time.Sleep(5 * time.Second)
			if err := wsConn.wsWrite(websocket.TextMessage, []byte("heartbeat")); err != nil {
				wsConn.wsClose() //关闭连接
				break
			}
		}
	}()

	for {
		msg, err := wsConn.wsRead()
		if err != nil {
			log.Println("获取消息出现错误", err.Error())
			break
		}
		log.Println("接收到消息", string(msg.Data))

		// 以下内容为获取内容后得返回信息，可以根据实际情况修改
		err = wsConn.wsWrite(msg.Type, msg.Data)
		if err != nil {
			log.Println("发送消息给客户端出现错误", err.Error())
			break
		}
	}
}
```

wsReadLoop代码：

```go
func (wsConn *WsConn) wsReadLoop() {
	for {
		msgType, data, err := wsConn.WsSocket.ReadMessage()
		if err != nil {
			goto error
		}
		req := &WsMessage{
			msgType,
			data,
		}
		select {
		case wsConn.InChan <- req:
		case <-wsConn.ClosedChan:
			goto closed
		}
	}
error:
	wsConn.wsClose()
closed:
}
```

wsWriteLoop代码：

```go
func (wsConn *WsConn) wsWriteLoop() {
	for {
		select {
		case msg := <-wsConn.OutChan:
			if err := wsConn.WsSocket.WriteMessage(msg.Type, msg.Data); err != nil {
				goto error
			}
		case <-wsConn.ClosedChan:
			goto closed
		}
	}
error:
	wsConn.wsClose()
closed:
}
```

wsWrite wsRead  wsClose代码：

```go
//关闭连接
func (wsConn *WsConn) wsClose() {
	wsConn.WsSocket.Close()
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.ClosedChan)
	}
}

//读数据
func (wsConn *WsConn) wsRead() (*WsMessage, error) {
	select {
	case msg := <-wsConn.InChan:
		return msg, nil
	case <-wsConn.ClosedChan:
	}
	return nil, errors.New("websocket closed")
}

//写数据
func (wsConn *WsConn) wsWrite(messageType int, data []byte) error {
	select {
	case wsConn.OutChan <- &WsMessage{messageType, data}:
	case <-wsConn.ClosedChan:
		return errors.New("webscoket closed")
	}
	return nil
}
```

最后所以逻辑代码完成，我们按着Http Server方式来启动服务:

```go
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsHandler)
	http.ListenAndServe("127.0.0.1:8888", mux)
}
```



### WebSocket客户端实现

客户端代码分为2种方式，一种是通过JS现实，一种是通过go语言实现，代码我就直接贴下面了，不做过多解析了。

Go客户端：

```go
package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
)

func main() {
	origin := "http://127.0.0.1:8888/"
	url := "ws://127.0.0.1:8888/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ws)
	_, err = ws.Write([]byte("zzzzzzzzzz"))
	fmt.Println(err)

	for {
		var msg = make([]byte, 1000)
		m, err := ws.Read(msg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Receive: %s\n", msg[:m])
	}
}
```

HTML+JS 客户端：

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <script>
        window.addEventListener("load", function(evt) {
            var output = document.getElementById("output");
            var input = document.getElementById("input");
            var ws;
            var print = function(message) {
                var d = document.createElement("div");
                d.innerHTML = message;
                output.appendChild(d);
            };
            document.getElementById("open").onclick = function(evt) {
                if (ws) {
                    return false;
                }
                ws = new WebSocket("ws://127.0.0.1:8888/ws");
                ws.onopen = function(evt) {
                    print("OPEN");
                }
                ws.onclose = function(evt) {
                    print("CLOSE");
                    ws = null;
                }
                ws.onmessage = function(evt) {
                    print("RESPONSE: " + evt.data);
                }
                ws.onerror = function(evt) {
                    print("ERROR: " + evt.data);
                }
                return false;
            };
            document.getElementById("send").onclick = function(evt) {
                if (!ws) {
                    return false;
                }
                print("SEND: " + input.value);
                ws.send(input.value);
                return false;
            };
            document.getElementById("close").onclick = function(evt) {
                if (!ws) {
                    return false;
                }
                ws.close();
                return false;
            };
        });
    </script>
</head>
<body>
<table>
    <tr><td valign="top" width="50%">
        <p>Click "Open" to create a connection to the server,
            "Send" to send a message to the server and "Close" to close the connection.
            You can change the message and send multiple times.
        </p>
            <form>
                <button id="open">Open</button>
                <button id="close">Close</button>
            <input id="input" type="text" value="Hello world!">
            <button id="send">Send</button>
            </form>
    </td><td valign="top" width="50%">
        <div id="output"></div>
    </td></tr></table>
</body>
</html>
```

