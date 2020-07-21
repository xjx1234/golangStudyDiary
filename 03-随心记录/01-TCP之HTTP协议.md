## TCP/IP协议之HTTP协议

### TCP/IP协议

TCP/IP是一套用于网络通信的协议集合或者系统。TCP/IP协议模型就有OSI模型分为7层。但其实一般我们所谈到的都是四层的TCP/IP协议栈。如图：



![](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/tcp-ip-iso.jpg)



网络接口层：主要是指一些物理层层次的接口，比如电缆等

网络层：提供了独立于硬件的逻辑寻址，实现物理地址和逻辑地址的转换。网络层协议包括IP协议（网际协议），ICMP协议（互联网控制报文协议），IGMP协议(Internet组协议管理)

传输层：为网络提供了流量控制，错误控制和确认服务。传输层有两个互不相同的传输协议：TCP（传输控制协议）、UDP（用户数据报协议）

应用层：为文件传输，网络排错和Internet操作提供具体的程序应用。

TCP/IP具体是怎么通信的呢？

TCP/IP协议中由上到下的将数据封装成包，然后再由下至上的拆包，如下显示数据的打包拆包以及通讯流程：

![](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/tcp-ip-tx.jpg)

### TCP协议

#### TCP报文

TCP协议是面向连接的、可靠传输、有流量控制，拥塞控制，面向字节流传输等很多优点的协议。在前图中，我们看到传输层的HTTP数据包由 TCP首部 + Http数据 组成，我们这边简称为 TCP报文，如下：

![](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/tcp-ip-bw.jpg)

为了完成三次挥手四次握手，这里需要知道序列号seq、确认应答序号ack（小写字母）、控制位：

**序列号seq ：**

​		因为在TCP是面向字节流的，他会将报文都分成一个个字节，给每个字节进行序号编写，比如一个报文有900个字节组成，那么就会编成1-900个序号，然后分几部分来进行传输，比如第一次传，序列号就是1，传了50个字节， 那么第二次传，序列号就为51，所以序列号就是传输的数据的第一个字节相对所有的字节的位置。

**确认应答ack:** 

　　　　如刚说的例子，第一次传了50个字节给对方，对方也会回应你，其中带有确认应答，就是告诉你下一次要传第51个字节来了，所以这个确认应答就是告诉对方下一次要传第多少个字节了。也就是说告诉序列号下一次从哪里开始

**常见控制位:**

　　　　URG:紧急，当URG为1时，表名紧急指针字段有效，标识该报文是一个紧急报文，传送到目标主机后，不用排队，应该让该报文尽量往下排，让其早点让应用程序给接受。

　　　　ACK:确认，当ACK为1时，确认序号才有效。当ACK为0时，确认序号没用

　　　　PSH：推送，当为1时，当遇到此报文时，会减少数据向上交付，本来想应用进程交付数据是要等到一定的缓存大小才发送的，但是遇到它，就不用在等足够多的数据才向上交付，而是让应用进程早点拿到此报文，这个要和紧急分清楚，紧急是插队，但是提交缓存大小的数据不变，这个推送就要排队，但是遇到他的时候，会减少交付的缓存数据，提前交付。

　　　　RST:复位，报文遇到很严重的差错时，比如TCP连接出错等，会将RST置为1，然后释放连接，全部重新来过。

　　　　SYN：同步，在进行连接的时候，也就是三次握手时用得到，下面会具体讲到，配合ACK一起使用

　　　　FIN：终止，在释放连接时，也就是四次挥手时用的。

大致了解完 TCP协议的报文后，我们来看下TCP最著名的三次握手以及四次挥手机制

#### TCP三次握手

在通信之前，TCP协议会先通过三次握手的机制来确认两端口之间的连接是否可用，如下图所示：

<img src="https://myvoice1.oss-cn-beijing.aliyuncs.com/github/tcp-ip-3.jpg" style="zoom: 43%;" />

为了更直观的查看了解三次握手的过程，我们采用Wireshark抓包工具来了解TCP报文的信息，为了更好的对应相关抓包数据跟TCP报文数据，我们进行了下图的整理对应：

![](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/zhuabao-tcp-dy.png)



我们抓取了一组三次握手的流程，图中可以看到wireshark截获到了三次握手的三个数据包。第四个包才是HTTP的， 这说明HTTP的确是使用TCP建立连接的。如图：

![image-20200721122718263](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/zhuabao-group.png)

1. 第一次握手，客户端想要连接，创建传输控制块TCB，状态变为主动打开。发送给服务器不包含数据内容的连接请求报文。该请求报文首部中同步位SYN=1，同时选择一个初始序列号seq=x（此处序号为 0 ）。然后客户端进入 SYN-SENT （同步已发送）状态，告诉服务器我想和你同步连接。TCP规定，SYN报文段（SYN=1的报文段）不能携带数据，但需要消耗掉一个序号。

   seq=0；ack=0；SYN=1

   ![](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/TCP.png)

2. 第二次握手， TCP服务器收到连接请求报文，如果同意连接则发送确认报文。为了保证下次客户端发送报文时seq序列号是正确的，需要发送确认号ack=x+1（x第一次握手时为0，此处x+1=1），同时确认号ack要生效必须发送ACK=1，再加上同步位SYN=1，序列号seq=y（携带Y个字节，此处为0），然后服务器也进 入SYN-RCVD (同步已收到) 状态，完成同步连接。这个报文也是SYN报文，也不能携带数据，但是同样要消耗一个序号。

   SYN=1; ACK=1;ack=1；seq=0

   ![](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/tcp2.png)

3. 第三次握手，客户端收到确认后还要再向服务器发送确认报文。确认报文已经不是请求报文SYN了，不再包含SYN同步位。发送的内容有序列号seq=x+1（和第二次握手的ACK对应），确认号ack=y+1，ACK=1。客户端发送确认报文以后进入ESTABLISHED（已建立）状态，服务器接收到确认报文以后也进入ESTABLISHED状态。此时TCP连接完成建立。

   SYN=0; ACK=1;ack=1；seq=1

   ![](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/tcp3.png)

然后就可以发送TCP接收到Http的数据包后生成的新数据包了！

> 但是貌似看起来两次握手请求就可以完成事，为什么非要三次握手呢？主要是为了防止已经失效的连接请求报文突然又传到了服务器，从而产生错误。
>
> 如果是两次握手，假设一种情景：客户端发送了第一个请求连接报文并未丢失，只是因为网络问题在网络节点中滞留太久了。由于客户端迟迟没有收到确认报文，以为服务器没有收到。于是再发送一条请求连接报文，此时一路畅通完成两次握手建立连接，传输数据，关闭连接。然后那个前一条龟速的请求报文终于走到了服务器，再次和服务器建立连接，这就造成了不必要的资源浪费。如果是三次握手，就算那一条龟速的请求报文最后到达了服务器，然后服务器也发送了确认连接报文，但是此时客户端已经不会再发出确认报文了，服务器也接受不到确认报文，于是无法建立连接。

#### TCP四次挥手

数据传输完毕后，双方都可释放连接。最开始的时候，客户端和服务器都是处于ESTABLISHED状态，然后客户端主动关闭，服务器被动关闭。

<img src="https://myvoice1.oss-cn-beijing.aliyuncs.com/github/tcp-bb.jpg" style="zoom:60%;" />

​    同样我们抓取整个数据流程作为分析使用：

​	![](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/tcp-ball.png)

1. 第一次挥手，客户端从ESTABLISHED状态变为主动关闭状态，客户端发送请求释放连接报文给服务器，FIN=1，seq=u（等于前面已经传送过来的数据的最后一个字节的序号加1），此时客户端进入FIN-WAIT-1（终止等待1）状态。 TCP规定，FIN报文段即使不携带数据，也要消耗一个序号。

   FIN=1; seq=1;ACK=1

   ![](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/tcp-bb1.png)

2. 第二次挥手，服务器接收到客户端发来的请求释放报文以后，发送确认报文告诉客户端我收到了你的请求，内容差不多就是seq=v，ack=u+1，ACK=1，此时服务器进入CLOSE-WAIT（关闭等待）状态。为什么是CLOSE-WAIT状态？可能自己服务器这端还有数据没有发送完，所以这个时候整个TCP的连接就变成了半关闭状态。服务器还能发送数据，客户端也能接收数据，但客户端不能再发送数据了，只能发送确认报文。客户端接收到服务器传来的确认报文以后，进入 FIN-WAIT-1（终止等待2）状态，等待服务器发送连接释放的报文（在这之前，还需要接受服务器没有发送完的最后的数据）。 

   ACK=1;ack=2;seq=2

   ![](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/tcp-bb2.png)

3. 第三次挥手，服务器所有的数据都发送完了，认为可以关闭连接了，于是向客户端发送连接释放报文，内容FIN=1，seq=w，ack=u+1（客户端没发送消息，所以提醒客户端下一次还是从u+1开始发送序列），ACK=1。此时服务器进入了 LAST-ACK（最后确认）状态，等待客户端发送确认报文。

   ACK=1;FIN=1;seq=2;ack=2

   ![](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/tcp-bb3.png)

4. 第三次挥手，客户端接收到了服务器发送的连接释放报文，必须发出确认。确认报文seq=u+1，ack=w+1，ACK=1。此时客户端进入 TIME-WAIT （时间等待）状态，但是没有立马关闭。此时TCP连接还没有释放，必须经过2MSL（最长报文段寿命）的时间后，当客户端撤销相应的TCB后，才进入CLOSED状态。因为这个确认报文可能丢失。服务器收不到确认报文心想这可能是我没传到或者丢失了啊，于是服务器再传一个FIN，然后客户端再重新发送一个确认报文。然后刷新2∗∗MSL时间。直到这个时间内收不到FIN连接释放报文，客户端撤销TCB进入CLOSE状态。而服务器，在接收到确认报文的时候就立马变为CLOSE状态了。所以服务器结束TCP连接的时间略早于客户端。

   ACK=1;seq=2;ack=3

   ![](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/tcp-bb4.png)

> 万一确认连接以后客户端故障怎么办？
>
> TCP设有一个保活计时器。显然客户端故障时服务器不会智障般等下去，白白浪费资源。服务器每次收到一次客户端的请求以后都会刷新这个保活计时器，时间通常设置为2小时。若2个小时依旧没有收到客户端的任何数据，服务器会发送一个探测报文段，每隔75分钟发一个，如果连发十个都没有数据反应，那么服务器就知道客户端故障了，关闭连接。

### HTTP协议

​		HTTP协议是Hyper Text Transfer Protocol（超文本传输协议）的缩写,是用于从万维网（WWW:World Wide Web ）服务器传输超文本到本地浏览器的传送协议。

​		HTTP是一个基于TCP/IP通信协议来传递数据（HTML 文件, 图片文件, 查询结果等）。

#### Http请求报文

http请求报文格式：

![](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/http-bw.jpg)

http请求报文示例：

```
/** 请求行 */
GET /tool/block.php?action=wallet HTTP/1.1
/** 请求头 */
Host: xjx.cn
Connection: keep-alive
Pragma: no-cache
Cache-Control: no-cache
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8
Accept-Encoding: gzip, deflate
Accept-Language: zh-CN,zh;q=0.9,en;q=0.8
Cookie: __clickidc=154744012215372025
/** 请求数据 */
action=wallet
```

**请求行：**

请求行有以下信息组成：

> [请求方法] [URL] [版本号]		
示例：

> GET /tool/block.php?action=wallet HTTP/1.1

Http1.1 请求报文的方法有以下几种：

GET：获取资源
POST：传输实体主体
PUT：传输文件
HEAD：获得报文首部（相当于GET方法获得的资源去掉正文）
DELETE：删除文件
OPTIONS：询问支持的方法（客户端问服务器）
TRACE：追踪路径
OCONNECT：要求用隧道协议连接代理

GET 方法和 POST 方法核心点：

1. 传参的数据量不一样，一个通过 url，一个通过正文，所以 POST 能传更多的数据；
2. GET 方法和 POST 方法传参位置上，可靠性问题。

**请求头部：**

下面列举常见的请求头部信息如下：

|       Accept        |                 指定客户端能够接收的内容类型                 |              Accept: text/plain, text/html              |
| :-----------------: | :----------------------------------------------------------: | :-----------------------------------------------------: |
|   Accept-Charset    |                 浏览器可以接受的字符编码集。                 |               Accept-Charset: iso-8859-5                |
|   Accept-Encoding   |     指定浏览器可以支持的web服务器返回内容压缩编码类型。      |             Accept-Encoding: compress, gzip             |
|   Accept-Language   |                      浏览器可接受的语言                      |                 Accept-Language: en,zh                  |
|    Accept-Ranges    |           可以请求网页实体的一个或者多个子范围字段           |                  Accept-Ranges: bytes                   |
|    Authorization    |                      HTTP授权的授权证书                      |    Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==    |
|    Cache-Control    |                 指定请求和响应遵循的缓存机制                 |                 Cache-Control: no-cache                 |
|     Connection      |      表示是否需要持久连接。（HTTP 1.1默认进行持久连接）      |                    Connection: close                    |
|       Cookie        | HTTP请求发送时，会把保存在该请求域名下的所有cookie值一起发送给web服务器。 |              Cookie: $Version=1; Skin=new;              |
|   Content-Length    |                        请求的内容长度                        |                   Content-Length: 348                   |
|    Content-Type     |                  请求的与实体对应的MIME信息                  |     Content-Type: application/x-www-form-urlencoded     |
|        Date         |                     请求发送的日期和时间                     |           Date: Tue, 15 Nov 2010 08:12:31 GMT           |
|       Expect        |                    请求的特定的服务器行为                    |                  Expect: 100-continue                   |
|        From         |                    发出请求的用户的Email                     |                  From: user@email.com                   |
|        Host         |                指定请求的服务器的域名和端口号                |                   Host: www.zcmhi.com                   |
|      If-Match       |                只有请求内容与实体相匹配才有效                |      If-Match: “737060cd8c284d8af7ad3082f209582d”       |
|  If-Modified-Since  | 如果请求的部分在指定时间之后被修改则请求成功，未被修改则返回304代码 |    If-Modified-Since: Sat, 29 Oct 2010 19:43:31 GMT     |
|    If-None-Match    | 如果内容未改变返回304代码，参数为服务器先前发送的Etag，与服务器回应的Etag比较判断是否改变 |    If-None-Match: “737060cd8c284d8af7ad3082f209582d”    |
|      If-Range       | 如果实体未改变，服务器发送客户端丢失的部分，否则发送整个实体。参数也为Etag |      If-Range: “737060cd8c284d8af7ad3082f209582d”       |
| If-Unmodified-Since |           只在实体在指定时间之后未被修改才请求成功           |   If-Unmodified-Since: Sat, 29 Oct 2010 19:43:31 GMT    |
|    Max-Forwards     |               限制信息通过代理和网关传送的时间               |                    Max-Forwards: 10                     |
|       Pragma        |                    用来包含实现特定的指令                    |                    Pragma: no-cache                     |
| Proxy-Authorization |                     连接到代理的授权证书                     | Proxy-Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ== |
|        Range        |                 只请求实体的一部分，指定范围                 |                  Range: bytes=500-999                   |
|       Referer       |         先前网页的地址，当前请求网页紧随其后,即来路          |     Referer: http://www.zcmhi.com/archives/71.html      |
|         TE          |   客户端愿意接受的传输编码，并通知服务器接受接受尾加头信息   |               TE: trailers,deflate;q=0.5                |
|       Upgrade       |    向服务器指定某种传输协议以便服务器进行转换（如果支持）    |     Upgrade: HTTP/2.0, SHTTP/1.3, IRC/6.9, RTA/x11      |
|     User-Agent      |            User-Agent的内容包含发出请求的用户信息            |          User-Agent: Mozilla/5.0 (Linux; X11)           |
|         Via         |            通知中间网关或代理服务器地址，通信协议            |       Via: 1.0 fred, 1.1 nowhere.com (Apache/1.1)       |
|       Warning       |                    关于消息实体的警告信息                    |             Warn: 199 Miscellaneous warning             |

#### Http响应报文

HTTP的响应报文也由三部分组成（响应行+响应头+响应正文）

![](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/http-xy.png)

**状态行：**

由HTTP协议版本号， 状态码， 状态消息 三部分组成

常见的状态码如下表：

状态码	    类别	                                   原因短语
1XX	Informational (信息性状态码)	接收的请求正在处理
2XX	Success (成功状态码)	请求正常处理完毕
3XX	Redirection (重定向状态码)	需要进行附加操作以完成请求
4XX	Client Error (客户端错误状态码)	服务器无法处理请求
5XX	Server Error (服务器错误状态码)	服务器处理请求出错

需要了解更多的响应码，参考： [HTTP状态码](https://www.runoob.com/http/http-status-codes.html)

**响应头：**

下面列举了常见的相应头信息描述: 

|      应答头      |                             说明                             |
| :--------------: | :----------------------------------------------------------: |
|      Allow       |          服务器支持哪些请求方法（如GET、POST等）。           |
| Content-Encoding | 文档的编码（Encode）方法。只有在解码之后才可以得到Content-Type头指定的内容类型。利用gzip压缩文档能够显著地减少HTML文档的下载时间。Java的GZIPOutputStream可以很方便地进行gzip压缩，但只有Unix上的Netscape和Windows上的IE 4、IE 5才支持它。因此，Servlet应该通过查看Accept-Encoding头（即request.getHeader("Accept-Encoding")）检查浏览器是否支持gzip，为支持gzip的浏览器返回经gzip压缩的HTML页面，为其他浏览器返回普通页面。 |
|  Content-Length  | 表示内容长度。只有当浏览器使用持久HTTP连接时才需要这个数据。如果你想要利用持久连接的优势，可以把输出文档写入 ByteArrayOutputStream，完成后查看其大小，然后把该值放入Content-Length头，最后通过byteArrayStream.writeTo(response.getOutputStream()发送内容。 |
|   Content-Type   | 表示后面的文档属于什么MIME类型。Servlet默认为text/plain，但通常需要显式地指定为text/html。由于经常要设置Content-Type，因此HttpServletResponse提供了一个专用的方法setContentType。 |
|       Date       | 当前的GMT时间。你可以用setDateHeader来设置这个头以避免转换时间格式的麻烦。 |
|     Expires      |       应该在什么时候认为文档已经过期，从而不再缓存它？       |
|  Last-Modified   | 文档的最后改动时间。客户可以通过If-Modified-Since请求头提供一个日期，该请求将被视为一个条件GET，只有改动时间迟于指定时间的文档才会返回，否则返回一个304（Not Modified）状态。Last-Modified也可用setDateHeader方法来设置。 |
|     Location     | 表示客户应当到哪里去提取文档。Location通常不是直接设置的，而是通过HttpServletResponse的sendRedirect方法，该方法同时设置状态代码为302。 |
|     Refresh      | 表示浏览器应该在多少时间之后刷新文档，以秒计。除了刷新当前文档之外，你还可以通过setHeader("Refresh", "5; URL=http://host/path")让浏览器读取指定的页面。  注意这种功能通常是通过设置HTML页面HEAD区的＜META HTTP-EQUIV="Refresh" CONTENT="5;URL=http://host/path"＞实现，这是因为，自动刷新或重定向对于那些不能使用CGI或Servlet的HTML编写者十分重要。但是，对于Servlet来说，直接设置Refresh头更加方便。   注意Refresh的意义是"N秒之后刷新本页面或访问指定页面"，而不是"每隔N秒刷新本页面或访问指定页面"。因此，连续刷新要求每次都发送一个Refresh头，而发送204状态代码则可以阻止浏览器继续刷新，不管是使用Refresh头还是＜META HTTP-EQUIV="Refresh" ...＞。   注意Refresh头不属于HTTP 1.1正式规范的一部分，而是一个扩展，但Netscape和IE都支持它。 |
|      Server      | 服务器名字。Servlet一般不设置这个值，而是由Web服务器自己设置。 |
|    Set-Cookie    | 设置和页面关联的Cookie。Servlet不应使用response.setHeader("Set-Cookie", ...)，而是应使用HttpServletResponse提供的专用方法addCookie。参见下文有关Cookie设置的讨论。 |
| WWW-Authenticate | 客户应该在Authorization头中提供什么类型的授权信息？在包含401（Unauthorized）状态行的应答中这个头是必需的。例如，response.setHeader("WWW-Authenticate", "BASIC realm=＼"executives＼"")。  注意Servlet一般不进行这方面的处理，而是让Web服务器的专门机制来控制受密码保护页面的访问（例如.htaccess）。 |