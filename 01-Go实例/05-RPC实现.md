## RPC实现

什么是RPC呢？百度百科给出的解释是这样的：“RPC（Remote Procedure Call Protocol）——远程过程调用协议，它是一种通过网络从远程计算机程序上请求服务，而不需要了解底层网络技术的协议”。该协议允许运行于一台计算机的程序调用另一台计算机的子程序，而程序员无需额外地为这个交互作用编程。 如果涉及的软件采用面向对象编程，那么远程过程调用亦可称作远程调用或远程方法调用。

用通俗易懂的语言描述就是：RPC允许跨机器、跨语言调用计算机程序方法。打个比方，我用go语言写了个获取用户信息的方法getUserInfo，并把go程序部署在阿里云服务器上面，现在我有一个部署在腾讯云上面的php项目，需要调用golang的getUserInfo方法获取用户信息，php跨机器调用go方法的过程就是RPC调用。

在Golang中实现RPC的方式大体有三种，分别来看:

### net/rpc

Golang官方的`net/rpc`包使用`encoding/gob`进行编解码，支持 tcp 或 http 数据传输方式。但是由于gob编码是Golang独有的所以它只支持Golang开发的服务器与客户端之间的交互。

下面我们拿一个例子作为演示，该例子是 计算圆周长 以及 面积的示例：

首先定义一个圆的结构体：

```go
type Circular struct{}
```

其次定义好圆所需的参数：

```go
type Params struct {
	Radius float64  //半径
}
const π = 3.1415  //圆形π值
```

定义圆形周长以及面积计算方法：

```go
//周长计算方法
func (c *Circular) GetPerimeter(p Params, perimeter *float64) error {
	*perimeter = π * p.Radius * 2
	return nil
}

//面积计算方法
func (c *Circular) GetArea(p Params, area *float64) error {
	*area = π * math.Sqrt(p.Radius)
	return nil
}
```

下面我们就来看下完整的服务端代码：

1. 基于 Http协议的Rpc 服务端

   ```go
   func main() {
     	listener := &Circular{}
   	rpc.Register(listener)
       rpc.HandleHTTP()
       lc, err := net.Listen("tcp", "127.0.0.1:8081")
       if err != nil {
           log.Fatal(err)
           defer wg.Done()
       }
       http.Serve(lc, nil)
   }
   ```

   也可以简化为：

   ```go
   func main() {
     	listener := &Circular{}
   	rpc.Register(listener)
       rpc.HandleHTTP()
       http.ListenAndServe(":8081", nil)
   }
   ```

   

2. 基于TCP协议的Rpc服务端

   ```go
   laddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8082")
   if err != nil {
       log.Fatal(err)
   } else {
       bg, err1 := net.ListenTCP("tcp", laddr)
       if err1 != nil {
           log.Fatal(err1.Error())
       }
       for {
           conn, err2 := bg.Accept()
           if err2 != nil {
               continue
           }
           go rpc.ServeConn(conn)
       }
   }
   ```

   也可以简化为：

   ```go
   laddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8082")
   if err != nil {
       log.Fatal(err)
   } else {
       bg, err1 := net.ListenTCP("tcp", laddr)
       if err1 != nil {
           log.Fatal(err1.Error())
       }
       rpc.Accept(bg)
   }
   ```

看完服务端代码，下面展示客户端代码

1. 基于 HTTP 协议的 RPC客户端代码

   ```go
   type Params struct {
   	Radius float64
   }
   func main() {
   	httpRpc, err := rpc.DialHTTP("tcp", "127.0.0.1:8081")
   	if err != nil {
   		log.Fatal(err)
   	}
   	ret := 0.0
   	error := httpRpc.Call("Circular.GetPerimeter", Params{1.1}, &ret)
   	if error != nil {
   		log.Fatal(error)
   	}
   	fmt.Printf("Http Perimeter: %v \r\n", ret)
   }
   ```

   

2. 基于 TCP协议的 RPC客户端代码

   ```go
   type Params struct {
   	Radius float64
   }
   
   func main(){
       tcpRpc, tcpErr := rpc.Dial("tcp", "127.0.0.1:8082")
       if tcpErr != nil {
           log.Fatal(tcpErr)
       }
       tcpRet := 0.0
       tcpErr1 := tcpRpc.Call("Circular.GetArea", Params{2.0}, &tcpRet)
       if tcpErr1 != nil {
           log.Fatal(tcpErr1)
       }
       fmt.Printf("TCP Area: %v \r\n", tcpRet)
   }
   ```

   

### net/jsonrpc

Go语言标准库通过`net/rpc/jsonrpc`这个包支持跨语言的RPC，弥补了 `net/rpc` 包不能跨语言的缺陷。同样我们以前面的例子作为示例。

jsonRpc服务端代码：

```go
func main(){
    laddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8083")
    if err != nil {
        log.Fatal(err)
    } else {
        bg, err1 := net.ListenTCP("tcp", laddr)
        if err1 != nil {
            log.Fatal(err1.Error())
        }
        for {
            conn, err2 := bg.Accept()
            if err2 != nil {
                continue
            }
            go jsonrpc.ServeConn(conn)
        }
    }
}
```

可以看出JSONRPC与基于TCP的RPC协议代码基本一致，唯一区别就是服务器连接一处代码： `go jsonrpc.ServeConn(conn)`

jsonRpc客户端代码：

```go
type Params struct {
	Radius float64
}

func main(){
    jsonRpc, jsonErr := jsonrpc.Dial("tcp", "127.0.0.1:8083")
    if jsonErr != nil {
        log.Fatal(jsonErr)
    }
    jsonRet := 0.0
    jsonErr1 := jsonRpc.Call("Circular.GetArea", Params{3.0}, &jsonRet)
    if jsonErr1 != nil {
        log.Fatal(jsonErr1)
    }
    fmt.Printf("JSONRPC Area: %v \r\n", tcpRet)  
}
```

### gRPC

jsonrpc虽然可以支持跨语言但是不支持HTTP传输，而且性能不是太突出，所以在实际生产环境中都不会用标准库里面的方式，而是选择Thrift、gRPC等方案。

gRPC是Google开源的高性能、通用的开源RPC框架，其主要面向移动应用开发并基于HTTP/2协议标准而设计，基于ProtoBuf序列化协议开发，支持Python、Golang、Java等众多开发语言。

Protobuf是Protocol Buffers的简称，它是Google公司开发的一种数据描述语言，类似于XML、JSON等数据描述语言，它非常轻便高效，很适合做数据存储或 RPC 数据交换格式。由于它一次定义，可生成多种语言的代码，非常适合用于通讯协议、数据存储等领域的语言无关、平台无关、可扩展的序列化结构数据格式。

