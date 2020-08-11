## 基于UDP协议实现的服务端和客户端

### UDP协议概念

用户数据报协议（User Datagram Protocol，缩写为UDP），又称用户数据报文协议，是一个简单的面向数据报(package-oriented)的传输层协议，正式规范为[RFC 768](https://tools.ietf.org/html/rfc768)。UDP只提供数据的不可靠传递，它一旦把应用程序发给网络层的数据发送出去，就不保留数据备份（所以UDP有时候也被认为是不可靠的数据报协议）。UDP在IP数据报的头部仅仅加入了复用和数据校验。

由于缺乏可靠性且属于非连接导向协议，UDP应用一般必须允许一定量的丢包、出错和复制粘贴。但有些应用，比如TFTP，如果需要则必须在应用层增加根本的可靠机制。但是绝大多数UDP应用都不需要可靠机制，甚至可能因为引入可靠机制而降低性能。流媒体（流技术）、即时多媒体游戏和IP电话（VoIP）一定就是典型的UDP应用。如果某个应用需要很高的可靠性，那么可以用传输控制协议（TCP协议）来代替UDP。

由于缺乏拥塞控制（congestion control），需要基于网络的机制来减少因失控和高速UDP流量负荷而导致的拥塞崩溃效应。换句话说，因为UDP发送者不能够检测拥塞，所以像使用包队列和丢弃技术的路由器这样的网络基本设备往往就成为降低UDP过大通信量的有效工具。数据报拥塞控制协议（DCCP）设计成通过在诸如流媒体类型的高速率UDP流中，增加主机拥塞控制，来减小这个潜在的问题。
典型网络上的众多使用UDP协议的关键应用一定程度上是相似的。这些应用包括域名系统（DNS）、简单网络管理协议（SNMP）、动态主机配置协议（DHCP）、路由信息协议（RIP）和某些影音流服务等等。

UDP 协议具备以下特点：

- 没有各种连接：在传输数据前不需要建立连接，也避免了后续的断开连接。
- 不重新排序：对到达顺序混乱的数据包不进行重新排序。
- 没有确认：发送数据包无须等待对方确认。因此，使用 UDP 协议可以随时发送数据，但无法保证数据能否成功被目标主机接收。

UDP在IP报文中的位置如下图所示：

![img](http://c.biancheng.net/uploads/allimg/191111/6-1911111249535K.gif)

UDP 报文中每个字段的含义如下：

- 源端口：这个字段占据 UDP 报文头的前 16 位，通常包含发送数据报的应用程序所使用的 UDP 端口。接收端的应用程序利用这个字段的值作为发送响应的目的地址。这个字段是可选的，所以发送端的应用程序不一定会把自己的端口号写入该字段中。如果不写入端口号，则把这个字段设置为 0。这样，接收端的应用程序就不能发送响应了。
- 目的端口：接收端计算机上 UDP 软件使用的端口，占据 16 位。
- 长度：该字段占据 16 位，表示 UDP 数据报长度，包含 UDP 报文头和 UDP 数据长度。因为 UDP 报文头长度是 8 个字节，所以这个值最小为 8。
- 校验值：该字段占据 16 位，可以检验数据在传输过程中是否被损坏。

为了验证发送的消息使用的是 UDP 协议，可以通过抓包进行查看，如图所示：

![img](http://c.biancheng.net/uploads/allimg/191111/6-19111112543a37.gif)

从数据包可以看到，是 UDP 客户端（源 IP 地址为 192.168.59.132）向 UDP 服务器端（目的 IP 地址为 192.168.59.135）发送的 UDP 数据包，使用的源端口为随机端口 47203，目的端口为 80（UDP 服务器端监听的端口）。

在 User Datagram Protocol 部分中显示了 UDP 数据包的详细信息。可以看到源端口、目的端口，以及包长度为 11 字节、校验值为 0xf878 等信息。

当服务器向客户端发送消息时，使用的也是 UDP 协议。例如，在服务器端回复客户端，输入 hello：

> root@daxueba:~# netwox 90 -P 80
> hi
> hello

通过抓包验证使用的是 UDP 协议，如图所示：

![img](http://c.biancheng.net/uploads/allimg/191111/6-1911111256445a.gif)



### UDP服务端实现

定义连接 ip 端口 以及 消息体结构体

```go
//定义ip以及端口
const (
	Address string = "127.0.0.1:8808"
)

//定义信息消息体
type UdpMessage struct {
	Addr *net.UDPAddr
	Data []byte
}

//定义连接结构体
type Server struct {
	Conn       *net.UDPConn
	ReadChan   chan *UdpMessage
	WriteChan  chan *UdpMessage
	ClosedChan chan byte
}
```

定义 管道读 管道写 管道关闭函数

```go
//管道读入
func (s *Server) read() (*UdpMessage, error) {
	select {
	case msg := <-s.ReadChan:
		return msg, nil
	case <-s.ClosedChan:
	}
	return nil, errors.New("udp closed")
}

//管道写处理
func (s *Server) write(d *UdpMessage) error {
	select {
	case s.WriteChan <- d:
	case <-s.ClosedChan:
		return errors.New("udp closed")
	}
	return nil
}

//管道关闭处理
func (s *Server) udpClosed() {
	s.Conn.Close()
	close(s.ClosedChan)
}
```

定义 读处理 写处理 以及 流程处理

```go
//处理UDP信息接收以及发送处理
func (s *Server) procLoop() {
	fmt.Println("procLoop")
	for {
		msg, err := s.read()
		if err != nil {
			log.Println("获取消息出现错误", err.Error())
			break
		}
		log.Println("接收到消息", string(msg.Data))

		// 以下内容为获取内容后得返回信息，可以根据实际情况修改
		err = s.write(msg)
		if err != nil {
			log.Println("发送消息给客户端出现错误", err.Error())
			break
		}
	}
}

//将UDP信息接收后存入管道
func (s *Server) readLoop() {
	fmt.Println("readLoop")
	for {
		data := make([]byte, 2048)
		n, addr, err := s.Conn.ReadFromUDP(data)
		if err != nil {
			goto End
		}
		fmt.Println("have data, data len:" + strconv.Itoa(n))
		req := &UdpMessage{
			addr,
			data[:n],
		}
		select {
		case s.ReadChan <- req:
		case <-s.ClosedChan:
			goto End
		}
	}
End:
	s.udpClosed()
}

//将写管道的信息发送给UDP客户端
func (s *Server) writeLoop() {
	fmt.Println("writeLoop")
	for {
		select {
		case msg := <-s.WriteChan:
			if _, wErr := s.Conn.WriteToUDP(msg.Data, msg.Addr); wErr != nil {
				goto End
			}
		case <-s.ClosedChan:
			goto End
		}
	}
End:
	s.udpClosed()
}
```

main主体，UDP操作：

```go
func main() {

	var wg sync.WaitGroup
	wg.Add(1)
	udpAddress, err := net.ResolveUDPAddr("udp", Address)
	if err != nil {
		log.Fatal(err)
	}
	conn, connErr := net.ListenUDP("udp", udpAddress)
	defer conn.Close()
	if connErr != nil {
		log.Fatal(connErr)
	}
	srv := &Server{
		Conn:      conn,
		ReadChan:  make(chan *UdpMessage, 50),
		WriteChan: make(chan *UdpMessage, 50),
	}
	fmt.Println("hello")
	go srv.procLoop()
	go srv.readLoop()
	go srv.writeLoop()
	wg.Wait()
}
```



### UDP客户端实现

```go
package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	ServerAddr string = "127.0.0.1:8808"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", ServerAddr)
	if err != nil {
		fmt.Println("x")
		log.Fatal(err)
	}
	conn, lerr := net.DialUDP("udp", nil, udpAddr)
	if lerr != nil {
		log.Fatal(lerr)
	}
	for {
		_, err = conn.Write([]byte("hello xjx")) // 发送数据给UDP服务器
		if err != nil {
			fmt.Println("udp 发送数据失败", err)
		}

		buf := make([]byte, 2048)
		n, readErr := conn.Read(buf)
		if readErr != nil {
			log.Fatal(readErr)
		}
		fmt.Println("client read len:", n)
		fmt.Println("client read data:", string(buf))

		time.Sleep(1000)
	}
}
```

