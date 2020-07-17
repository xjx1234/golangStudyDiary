## 包(Package)

### 包的概念

#### GOPATH

了解包前，我们需要了解一个概念：GOPATH，GOPATH 是 Go语言中使用的一个环境变量，它使用绝对路径提供项目的工作目录。工作目录是一个工程开发的相对参考目录，好比当你要在公司编写一套服务器代码，你的工位所包含的桌面、计算机及椅子就是你的工作区。工作区的概念与工作目录的概念也是类似的。如果不使用工作目录的概念，在多人开发时，每个人有一套自己的目录结构，读取配置文件的位置不统一，输出的二进制运行文件也不统一，这样会导致开发的标准不统一，影响开发效率。在命令行中运行`go env`后，命令行将提示以下信息：

```
$ go env
set GOARCH=amd64
set GOBIN=C:\Go\bin
set GOCACHE=C:\Users\fuzamei\AppData\Local\go-build
set GOEXE=.exe
set GOFLAGS=
set GOHOSTARCH=amd64
set GOHOSTOS=windows
set GOOS=windows
set GOPATH=D:\wamp64\www\golang
set GOPROXY=
set GORACE=
set GOROOT=C:\Go
set GOTMPDIR=
set GOTOOLDIR=C:\Go\pkg\tool\windows_amd64
set GCCGO=gccgo
set CC=gcc
set CXX=g++
set CGO_ENABLED=1
set GOMOD=
set CGO_CFLAGS=-g -O2
set CGO_CPPFLAGS=
set CGO_CXXFLAGS=-g -O2
set CGO_FFLAGS=-g -O2
set CGO_LDFLAGS=-g -O2
set PKG_CONFIG=pkg-config

```

命令行说明如下：

- 执行 go env 指令，将输出当前 Go 开发包的环境变量状态。
- GOARCH 表示目标处理器架构。
- GOBIN 表示编译器和链接器的安装位置。
- GOOS 表示目标操作系统。
- GOPATH 表示当前工作目录。
- GOROOT 表示 Go 开发包的安装目录

在 Go 1.8 版本之前，GOPATH 环境变量默认是空的。从 Go 1.8 版本开始，Go 开发包在安装完成后，将 GOPATH 赋予了一个默认的目录，参见下表:

|     平台     |   GOPATH默认值   |        示例        |
| :----------: | :--------------: | :----------------: |
| Windows 平台 | %USERPROFILE%/go | C:\Users\用户名\go |
|  Unix 平台   |     $HOME/go     |  /home/用户名/go   |

在 GOPATH 指定的工作目录下，代码总是会保存在 $GOPATH/src 目录下。在工程经过 go build、go install 或 go get 等指令后，会将产生的二进制可执行文件放在 $GOPATH/bin 目录下，生成的中间缓存文件会被保存在 $GOPATH/pkg 下。

关于GOPATH设置的方法，在这里就不一一列举了，参考文章： https://blog.csdn.net/chenjh213/article/details/51381024

#### 包的概念

Go语言是使用包来组织源代码的，包（package）是多个 Go 源码的集合，是一种高级的代码复用方案。Go语言中为我们提供了很多内置包，如 fmt、os、io 等。任何源代码文件必须属于某个包，同时源码文件的第一行有效代码必须是`package pacakgeName `语句，通过该语句声明自己所在的包。

Go语言的包借助了目录树的组织形式，一般包的名称就是其源文件所在目录的名称，虽然Go语言没有强制要求包名必须和其所在的目录名同名，但还是建议包名和所在目录同名，这样结构更清晰。

包可以定义在很深的目录中，包名的定义是不包括目录路径的，但是包在引用时一般使用全路径引用。比如在`GOPATH/src/a/b/ `下定义一个包 c。在包 c 的源码中只需声明为`package c`，而不是声明为`package a/b/c`，但是在导入 c 包时，需要带上路径，例如`import "a/b/c"`。

包的习惯用法：

- 包名一般是小写的，使用一个简短且有意义的名称。
- 包名一般要和所在的目录同名，也可以不同，包名中不能包含`- `等特殊符号。
- 包一般使用域名作为目录名称，这样能保证包名的唯一性，比如 GitHub 项目的包一般会放到`GOPATH/src/github.com/userName/projectName `目录下。
- 包名为 main 的包为应用程序的入口包，编译不包含 main 包的源码文件时不会得到可执行文件。
- 一个文件夹下的所有源码文件只能属于同一个包，同样属于同一个包的源码文件不能放在多个文件夹下。

### 包的导入

导入有两种基本格式，即单行导入和多行导入，两种导入方法的导入代码效果是一致的:

- 单行导入

  单行导入格式如下：

  > import "包1"

  > import "包2"

  

- 多行导入

  当多行导入时，包名在 import 中的顺序不影响导入效果，格式如下：

  >import(
  >  "包1"
  >  "包2"
  >  …
  >)

单行包导入示例:

```go
//单行示例
import "fmt"
import "time"
func main(){
	fmt.Println(time.Now())
}

//多行导入示例
import (
	"fmt"
	"time"
)
func main(){
	fmt.Println(time.Now())
}
```



### 常见的内置包以及自定义包

#### 内置包

标准的Go语言代码库中包含了大量的包，并且在安装 Go 的时候多数会自动安装到系统中。我们可以在 $GOROOT/src/pkg 目录中查看这些包。下面简单介绍一些我们开发中常用的包。这些包只是其中的一小部分。

#### 1) fmt

fmt 包实现了格式化的标准输入输出，这与C语言中的 printf 和 scanf 类似。其中的 fmt.Printf() 和 fmt.Println() 是开发者使用最为频繁的函数。

格式化短语派生于C语言，一些短语（%- 序列）是这样使用：

- %v：默认格式的值。当打印结构时，加号（%+v）会增加字段名；
- %#v：Go样式的值表达；
- %T：带有类型的 Go 样式的值表达。

#### 2) io

这个包提供了原始的 I/O 操作界面。它主要的任务是对 os 包这样的原始的 I/O 进行封装，增加一些其他相关，使其具有抽象功能用在公共的接口上。

#### 3) bufio

bufio 包通过对 io 包的封装，提供了数据缓冲功能，能够一定程度减少大块数据读写带来的开销。

在 bufio 各个组件内部都维护了一个缓冲区，数据读写操作都直接通过缓存区进行。当发起一次读写操作时，会首先尝试从缓冲区获取数据，只有当缓冲区没有数据时，才会从数据源获取数据更新缓冲。

#### 4) sort

sort 包提供了用于对切片和用户定义的集合进行排序的功能。

#### 5) strconv

strconv 包提供了将字符串转换成基本数据类型，或者从基本数据类型转换为字符串的功能。

#### 6) os

os 包提供了不依赖平台的操作系统函数接口，设计像 Unix 风格，但错误处理是 go 风格，当 os 包使用时，如果失败后返回错误类型而不是错误数量。

#### 7) sync

sync 包实现多线程中锁机制以及其他同步互斥机制。

#### 8) flag

flag 包提供命令行参数的规则定义和传入参数解析的功能。绝大部分的命令行程序都需要用到这个包。

#### 9) encoding/json

JSON 目前广泛用做网络程序中的通信格式。encoding/json 包提供了对 JSON 的基本支持，比如从一个对象序列化为 JSON 字符串，或者从 JSON 字符串反序列化出一个具体的对象等。

#### 10) html/template

主要实现了 web 开发中生成 html 的 template 的一些函数。

#### 11) net/http

net/http 包提供 HTTP 相关服务，主要包括 http 请求、响应和 URL 的解析，以及基本的 http 客户端和扩展的 http 服务。

通过 net/http 包，只需要数行代码，即可实现一个爬虫或者一个 Web 服务器，这在传统语言中是无法想象的。

#### 12) reflect

reflect 包实现了运行时反射，允许程序通过抽象类型操作对象。通常用于处理静态类型 interface{} 的值，并且通过 Typeof 解析出其动态类型信息，通常会返回一个有接口类型 Type 的对象。

#### 13) os/exec

os/exec 包提供了执行自定义 linux 命令的相关实现。

#### 14) strings

strings 包主要是处理字符串的一些函数集合，包括合并、查找、分割、比较、后缀检查、索引、大小写处理等等。

strings 包与 bytes 包的函数接口功能基本一致。

#### 15) bytes

bytes 包提供了对字节切片进行读写操作的一系列函数。字节切片处理的函数比较多，分为基本处理函数、比较函数、后缀检查函数、索引函数、分割函数、大小写处理函数和子切片处理函数等。

#### 16) log

log 包主要用于在程序中输出日志。

log 包中提供了三类日志输出接口，Print、Fatal 和 Panic。

- Print 是普通输出；
- Fatal 是在执行完 Print 后，执行 os.Exit(1)；
- Panic 是在执行完 Print 后调用 panic() 方法。

#### 自定义包

包是Go语言中代码组成和代码编译的主要方式。下面我们要建立一个demo名称的自定义包为例子，我们创建的自定义的包需要将其放在 GOPATH 的 src 目录下（也可以是 src 目录下的某个子目录），而且两个不同的包不能放在同一目录下，这样会引起编译错误。

1.  在GOPATH/src目录下建立 demo目录，然后在此目录下创建demo.go的go文件，当然go文件名称可以不用跟demo目录一致，件的名字也没有任何规定（但后缀必须是 .go），这里我们假设包名就是 .go 的文件名（如果一个包有多个 .go 文件，则其中会有一个 .go 文件的文件名和包名相同）

2. 在demo.go  (GOPATH/src/demo) 文件中创建内容，内容如下:

   ```go
   package demo //标识包名
   
   func MyAdd(num1, num2 int) int{
   	result := num1 + num2
   	return result
   }
   ```

   从上面代码可以看出，自定义包首先要标记该文件归属的包，格式为：

   > package 包名

   包的特性如下：

   - 一个目录下的同级文件归属一个包。
   - 包名可以与其目录不同名。
   - 包名为 main 的包为应用程序的入口包，编译源码没有 main 包时，将无法编译输出可执行的文件。

   3. 在main包 (GOPATH/src/main/main.go) 的main函数中调用自定义包的函数，代码如下:

      ```go
      package main
      
      import (
      	"demo"
      	"fmt"
      )
      
      func main() {
      	data := demo.MyAdd(1,2)
      	fmt.Println(data)
      }
      ```

      上述代码结构输出为:

      > 3

      对引用自定义包需要注意以下几点：

      - 如果项目的目录不在 GOPATH 环境变量中，则需要把项目移到 GOPATH 所在的目录中，或者将项目所在的目录设置到 GOPATH 环境变量中，否则无法完成编译；
      - 使用 import 语句导入包时，使用的是包所属文件夹的名称；
      - 包中的函数名第一个字母要大写，否则无法在外部调用；
      - 自定义包的包名不必与其所在文件夹的名称保持一致，但为了便于维护，建议保持一致；
      - 调用自定义包时使用 `包名 . 函数名` 的方式，如上例：demo.MyAdd()。

### 包的依赖管理

go module 是Go语言从 1.11 版本之后官方推出的版本管理工具，并且从 Go1.13 版本开始，go module 成为了Go语言默认的依赖管理工具。

常用的`go mod`命令如下表所示：

|        命令        |                      作用                      |
| :----------------: | :--------------------------------------------: |
|    go mod init     |       初始化当前文件夹，创建 go.mod 文件       |
|  go mod download   | 下载依赖包到本地（默认为 GOPATH/pkg/mod 目录） |
|    go mod edit     |                编辑 go.mod 文件                |
|    go mod graph    |                 打印模块依赖图                 |
|    go mod tidy     |           增加缺少的包，删除无用的包           |
|   go mod vendor    |           将依赖复制到 vendor 目录下           |
|   go mod verify    |                    校验依赖                    |
|     go mod why     |               解释为什么需要依赖               |
| go clean -modcache |                    清理缓存                    |

下面我们具体来看下在一个项目中如何使用 go module 管理依赖的：

1.  首先设置 Modules 开启

   那么如何设置Modules呢? 不同环境下按下面命令在命令行模式下执行即可。

> windows： set GO111MODULE=on

> linux: export GO111MODULE=on

2. 设置代理服务器，国内的网络有防火墙的存在，这导致有些Go语言的第三方包我们无法直接通过`go get`命令获取。

   GOPROXY 是Go语言官方提供的一种通过中间代理商来为用户提供包下载服务的方式。要使用 GOPROXY 只需要设置环境变量 GOPROXY 即可。

   目前公开的代理服务器的地址有：

   - goproxy.io；

   - goproxy.cn：（推荐）由国内的七牛云提供。

   - https://mirrors.aliyun.com/goproxy/ 阿里云提供

     

   Windows 下设置 GOPROXY 的命令为：

> set GOPROXY=https://goproxy.cn

​	  MacOS 或 Linux 下设置 GOPROXY 的命令为：

> export GOPROXY=https://goproxy.cn

​	  Go语言在 1.13 版本之后 GOPROXY 默认值为 https://proxy.golang.org，在国内可能会存在下载慢或者无法访问的情况，所以十分建议大家将 GOPROXY 设置为国内的 goproxy.cn。

3. 初始化项目，使用`go mod init`初始化生成 go.mod 文件

   go.mod 文件一旦创建后，它的内容将会被 go toolchain 全面掌控，go toolchain 会在各类命令执行时，比如`go get`、`go build`、`go mod`等修改和维护 go.mod 文件。

   go.mod 提供了 module、require、replace 和 exclude 四个命令：

   - module 语句指定包的名字（路径）；
   - require 语句指定的依赖项模块；
   - replace 语句可以替换依赖项模块；
   - exclude 语句可以忽略依赖项模块。

   初始化生成的 go.mod 文件如下所示：

   > module github.com/xjx1234/golangStudyDiary
   >
   > go 1.12

4.  我们将代码添加一个依赖关系，修改代码为：

   ```
   import (
   	"demo"
   	"fmt"
   	"time"
   	"github.com/labstack/echo"
   )
   ```

   上述代码中新增了一个github依赖库，此时我们 `go run package.go` 命令后，发现 go mod 会自动查找依赖自动下载：

   > $ go run package.go
   > go: finding github.com/labstack/gommon/color latest
   > go: finding github.com/labstack/gommon/log latest
   > go: golang.org/x/sys@v0.0.0-20190222072716-a9d3bda3a223: 
   >
   > ...

我们也可以手工更新这些依赖库， 使用命令 `go mod tidy` 即可。

我们也可以使用 go get命令下载指定版本依赖包, 执行`go get `命令，在下载依赖包的同时还可以指定依赖包的版本。

- 运行`go get -u`命令会将项目中的包升级到最新的次要版本或者修订版本；
- 运行`go get -u=patch`命令会将项目中的包升级到最新的修订版本；
- 运行`go get [包名]@[版本号]`命令会下载对应包的指定版本或者将对应包升级到指定的版本。

提示：`go get [包名]@[版本号]`命令中版本号可以是 x.y.z 的形式，例如 go get foo@v1.2.3，也可以是 git 上的分支或 tag，例如 go get foo@master，还可以是 git 提交时的哈希值，例如 go get foo@e3702bed2。

由于某些已知的原因，并不是所有的 package 都能成功下载，比如：golang.org 下的包。modules 可以通过在 go.mod 文件中使用 replace 指令替换成 github 上对应的库，比如：

```
go mod edit -replace=golang.org/x/crypto@v0.0.0=github.com/golang/crypto@latest
```

然后在执行一次 `go mod tidy -v` 即可。

### 包的加载

Go 语言为以上问题提供了一个非常方便的特性：init() 函数。

init() 函数的特性如下：

- 每个源码可以使用 1 个 init() 函数。
- init() 函数会在程序执行前（main() 函数执行前）被自动调用。
- 调用顺序为 main() 中引用的包，以深度优先顺序初始化。

例如，假设有这样的包引用关系：main→A→B→C，那么这些包的 init() 函数调用顺序为：

> C.init→B.init→A.init→main

说明：

- 同一个包中的多个 init() 函数的调用顺序不可预期。
- init() 函数不能被其他函数调用。

Go 程序的启动和加载过程，在执行 main 包的 mian 函数之前， Go 引导程序会先对整个程序的包进行初始化。整个执行的流程如下图所示：

![](https://myvoice1.oss-cn-beijing.aliyuncs.com/github/package_load.gif)



Go语言包的初始化有如下特点：

- 包初始化程序从 main 函数引用的包开始，逐级查找包的引用，直到找到没有引用其他包的包，最终生成一个包引用的有向无环图。
- Go 编译器会将有向无环图转换为一棵树，然后从树的叶子节点开始逐层向上对包进行初始化。
- 单个包的初始化过程如上图所示，先初始化常量，然后是全局变量，最后执行包的 init 函数。

### Context

#### Context初识

在go的1.7之前，context还是非编制的(包golang.org/x/net/context中)，golang团队发现context这个东西还挺好用的，很多地方也都用到了，就把它收编了，**1.7版本正式进入标准库**。随着 Context 包的引入，标准库中很多接口因此加上了 Context 参数，例如 database/sql 包，Context 几乎成为了并发控制和超时控制的标准做法。

context常用的使用姿势：

1. 一个请求对应多个goroutine之间的数据交互
2. 超时控制
3. 上下文控制

Context 包的核心就是 Context 接口，其定义如下：

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
```

这个就是Context的底层数据结构，来分析下：

|   字段   |                             说明                             |
| :------: | :----------------------------------------------------------: |
| Deadline | 返回一个time.Time，表示当前Context应该结束的时间，ok则表示有结束时间 |
|   Done   | 当Context被取消或者超时时候返回的一个close的channel，告诉给context相关的函数要停止当前工作然后返回了。(这个有点像全局广播) |
|   Err    | Err 方法会返回当前 Context 结束的原因，它只会在 Done 返回的 Channel 被关闭时才会返回非空的值 |
|  Value   | Value 方法会从 Context 中返回键对应的值，对于同一个上下文来说，多次调用 Value 并传入相同的 Key 会返回相同的结果，该方法仅用于传递跨 API 和进程间跟请求域的数据。 |

Go语言内置两个函数：Background() 和 TODO()，这两个函数分别返回一个实现了 Context 接口的 background 和 todo。

Background() 主要用于 main 函数、初始化以及测试代码中，作为 Context 这个树结构的最顶层的 Context，也就是根 Context。

TODO()，它目前还不知道具体的使用场景，在不知道该使用什么 Context 的时候，可以使用这个。

background 和 todo 本质上都是 emptyCtx 结构体类型，是一个不可取消，没有设置截止时间，没有携带任何值的 Context。

另外Go语言还带了4个Context实现，分别是:   

- emptyCtx  完全空的Context，实现的函数也都是返回nil，仅仅只是实现了Context的接口
- cancelCtx  继承自Context，同时也实现了canceler接口
- timerCtx 继承自cancelCtx，增加了timeout机制
- valueCtx   存储键值对的数据

为了更方便的创建Context，包里头定义了Background来作为所有Context的根，它是一个emptyCtx的实例。

```go
var (
    background = new(emptyCtx)
    todo       = new(emptyCtx) // 
)
func Background() Context {
    return background
}
```

你可以认为所有的Context是树的结构，Background是树的根，当任一Context被取消的时候，那么继承它的Context 都将被回收。

#### Context实战应用

- **WithCancel**

  WithCancel 的函数源码：

  ```go
  func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
  	c := newCancelCtx(parent)
  	propagateCancel(parent, &c)
  	return &c, func() { c.cancel(true, Canceled) }
  }
  ```

  下面我们用一个示例，来理解此方法的使用，我们设计一个程序，在10以内的数字内产生随机数，当随机数累加超过等于20的时候，停止程序运行，代码如下：

  ```go
  package main
  
  import (
  	"context"
  	"fmt"
  	"math/rand"
  	"time"
  )
  
  func cancelFun(ctx context.Context) <-chan int {
  	c := make(chan int)
  	num := 0
  	t := 0
  	go func() {
  		for {
  			select {
  			case <-ctx.Done():
  				fmt.Printf("耗时 %d 秒， 随机数 %d \n", t, num)
  				return
  			case c <- num:
  				incr := rand.Intn(10)
  				num += incr
  				if num >= 20 {
  					num = 20
  				}
  				t++
  				fmt.Printf("随机数 %d \n", num)
  			}
  		}
  	}()
  	return c
  }
  
  func main(){
      ctx, cancel := context.WithCancel(context.Background())
  	randData := cancelFun(ctx)
  	for n := range randData {
  		if n >= 20 {
  			cancel()
  			break
  		}
  	}
  	fmt.Println("running ... ")
  	time.Sleep(time.Second)
  }
  ```

  输出：

  ```
  $ go run context.go
  随机数 1
  随机数 8
  随机数 15
  随机数 20
  随机数 20
  running ...
  耗时 5 秒， 随机数 20
  ```

- ### WithDeadline & WithTimeout

  WithDeadline与WithTimeout实现的源码:

  ```go
  func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
  	if cur, ok := parent.Deadline(); ok && cur.Before(d) {
  		// The current deadline is already sooner than the new one.
  		return WithCancel(parent)
  	}
  	c := &timerCtx{
  		cancelCtx: newCancelCtx(parent),
  		deadline:  d,
  	}
  	propagateCancel(parent, c)
  	dur := time.Until(d)
  	if dur <= 0 {
  		c.cancel(true, DeadlineExceeded) // deadline has already passed
  		return c, func() { c.cancel(true, Canceled) }
  	}
  	c.mu.Lock()
  	defer c.mu.Unlock()
  	if c.err == nil {
  		c.timer = time.AfterFunc(dur, func() {
  			c.cancel(true, DeadlineExceeded)
  		})
  	}
  	return c, func() { c.cancel(true, Canceled) }
  }
  
  func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
  	return WithDeadline(parent, time.Now().Add(timeout))
  }
  ```

  下面我们用一个示例，来理解此方法的使用，我们设计一个程序，在100以内的数字内产生随机数，当累计时间超过10秒后，停止程序运行，代码如下：

  ```go
  package main
  
  import (
  	"context"
  	"fmt"
  	"math/rand"
  	"time"
  )
  
  func timeOutFun(ctx context.Context) {
  	n := 0
  	for {
  		select {
  		case <-ctx.Done():
  			fmt.Println("timeOut")
  			return
  		default:
  			incr := rand.Intn(100)
  			n += incr
  			fmt.Printf("随机数 %d \n", n)
  		}
  		time.Sleep(time.Second)
  	}
  }
  
  func main(){
  	//ctx, f := context.WithDeadline(context.Background(), time.Now().Add(10))
  	ctx, f := context.WithTimeout(context.Background(), 10*time.Second)
  	timeOutFun(ctx)
  	defer f()
  }
  ```

  程序输出:

  ```go
  $ go run context.go
  随机数 81
  随机数 168
  随机数 215
  随机数 274
  随机数 355
  随机数 373
  随机数 398
  随机数 438
  随机数 494
  随机数 494
  timeOut
  ```

- ### WithValue

  WithValue的实现源码:

  ```go
  func WithValue(parent Context, key, val interface{}) Context {
  	if key == nil {
  		panic("nil key")
  	}
  	if !reflect.TypeOf(key).Comparable() {
  		panic("key is not comparable")
  	}
  	return &valueCtx{parent, key, val}
  }
  ```

  下面我们用一个示例，来理解此方法的使用，我们设计一个程序，将键值为 xjx 和 id 的值带入 context中，并在函数中取出，代码如下：

  ```go
  package main
  
  import (
  	"context"
  	"fmt"
  	"math/rand"
  	"time"
  )
  
  func valueFun(ctx context.Context) {
  	id, ok := ctx.Value("id").(int)
  	if !ok {
  		fmt.Println("id value get fail")
  	}
  	xjxData := ctx.Value("xjx")
  	fmt.Printf("id: %d \n", id)
  	fmt.Printf("xjx: %s \n", xjxData)
  
  }
  
  func main(){
  	ctx := context.WithValue(context.Background(), "xjx", "hello")
  	ctx = context.WithValue(ctx, "id", 1)
  	valueFun(ctx)
  }
  ```

  代码输出:

  ```go
  $ go run context.go
  id: 1
  xjx: hello
  ```

  使用 Context 的注意事项：

  - 不要把 Context 放在结构体中，要以参数的方式显示传递；
  - 以 Context 作为参数的函数方法，应该把 Context 作为第一个参数；
  - 给一个函数方法传递 Context 的时候，不要传递 nil，如果不知道传递什么，就使用 context.TODO；
  - Context 的 Value 相关方法应该传递请求域的必要数据，不应该用于传递可选参数；
  - Context 是线程安全的，可以放心的在多个 Goroutine 中传递。

### 常见包用法简述

#### flag包

在编写命令行程序（工具、server）时，我们有时需要对命令参数进行解析，各种编程语言一般都会提供解析命令行参数的方法或库，以方便程序员使用。在Go语言中的 flag 包中，提供了命令行参数解析的功能。下面我们就来看一下 flag 包可以做什么，它具有什么样的能力。

这里介绍几个概念：

- 命令行参数（或参数）：是指运行程序时提供的参数；
- 已定义命令行参数：是指程序中通过 flag.Type 这种形式定义了的参数；
- 非 flag（non-flag）命令行参数（或保留的命令行参数）：可以简单理解为 flag 包不能解析的参数。

flag 包支持的命令行参数类型有 bool、int、int64、uint、uint64、float、float64、string、duration，如下表所示：

|      参数      |                            有效值                            |
| :------------: | :----------------------------------------------------------: |
|  字符串 flag   |                          合法字符串                          |
|   整数 flag    |           1234、0664、0x1234 等类型，也可以是负数            |
|  浮点数 flag   |                          合法浮点数                          |
| bool 类型 flag |   1、0、t、f、T、F、true、false、TRUE、FALSE、True、False    |
|  时间段 flag   | 任何合法的时间段字符串，如“300ms”、“-1.5h”、“2h45m”，<br/>合法的单位有“ns”、“us”、“µs”、“ms”、“s”、“m”、“h” |

下面来看下 flag包的基本使用:

1.  flag.Type()

   基本格式为:

   > flag.Type(flag 名, 默认值, 帮助信息) *Type

   示例:

   ```go
   	flag.Int("age", 11, "年龄")
   	flag.String("name", "xjx", "姓名")
   	flag.Bool("mrs", true, "婚否")
   ```

   需要注意的是，此时 name、age、mrs 均为对应类型的指针。

2. flag.TypeVar()

   基本格式为：

   > flag.TypeVar(Type 指针, flag 名, 默认值, 帮助信息)

   示例:

   ```go
   	var name string
   	var age int
   	var mrs bool
   	flag.StringVar(&name, "name", "xjx", "姓名")
   	flag.IntVar(&age, "age", 0, "年龄")
   	flag.BoolVar(&mrs, "mrs", false, "婚否")
   ```

   

3. flag.Parse

   通过以上两种方法定义好命令行 flag 参数后，需要通过调用 flag.Parse() 来对命令行参数进行解析。

   支持的命令行参数格式有以下几种：

   - -flag=x；
   - -flag x：只支持非 bool 类型。

   其中，布尔类型的参数必须使用等号的方式指定。

   flag 包的其他函数：

   > flag.Args() //返回命令行参数后的其他参数，以 []string 类型

   >flag.NArg() //返回命令行参数后的其他参数个数

   > flag.NFlag() //返回使用的命令行参 数个数

   

   示例:

   ```go
   package main
   
   import (
   	"flag"
   	"fmt"
   )
   
   var name = flag.String("name", "", "姓名")
   var age = flag.Int("age", 0, "年龄")
   var mrs bool
   
   func Init() {
   	flag.BoolVar(&mrs, "mrs", false, "婚否")
   }
   
   func main() {
   	Init()
   	flag.Parse()
   	fmt.Printf("args=%s, num=%d\n", flag.Args(), flag.NArg())
   	for i := 0; i != flag.NArg(); i++ {
   		fmt.Printf("arg[%d]=%s\n", i, flag.Arg(i))
   	}
   	fmt.Println("name=", *name)
   	fmt.Println("age=", *age)
   	fmt.Println("mrs=", mrs)
   }
   ```

   运行结果如下：

   ```go
   $ go run flag.go  -name="xjx1" -age=33 -mrs=false
   args=[], num=0
   name= xjx1
   age= 33
   mrs= false
   ```

   

4. 自定义value

   自定义是指自己定义一个类型用于flag类型的解析，基本格式为:

   > flag.Var(&flagVal, "name", "help message for flagname")

   

   下面废话不多说直接看示例：

   ```go
   type sliceValue []string
   
   func newSliceValue(vals []string, p *[]string) *sliceValue{
   	*p = vals
   	return (*sliceValue)(p)
   }
   
   /* 
   
   func Var(value Value, name string, usage string) {
   	CommandLine.Var(value, name, usage)
   }
   
   type Value interface {
   	String() string
   	Set(string) error
   }
   */
   
   // 因为flag.Var中第一个Value参数必须实现 String Set 函数
   func (s *sliceValue) String() string{
   	*s = sliceValue(strings.Split("default", ","))
   	return "none of my business"
   }
   
   // 因为flag.Var中第一个Value参数必须实现 String Set 函数
   func (s *sliceValue) Set(val string) error{
   	*s = sliceValue(strings.Split(val, ","))
   	return nil
   }
   
   func main() {
   	var languages []string
   	flag.Var(newSliceValue([]string{}, &languages), "slice", "i like")
   	flag.Parse()
   	fmt.Println(languages)
   }
   
   ```

    运行代码结果为:

   ```go
   $ go run flag.go  -slice go,java,php
   [go java php]
   ```

#### time包

#### os包

#### regexp包

#### big包