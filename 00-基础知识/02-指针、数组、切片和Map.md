## 指针、数组、切片和Map

### 数组

数组是一个由固定长度的特定类型元素组成的序列，一个数组可以由零个或多个元素组成。因为数组的长度是固定的，所以在Go语言中很少直接使用数组。和数组对应的类型是 Slice（切片），Slice 是可以增长和收缩的动态序列，功能也更灵活。所以切片在Go中使用频率更高。

数组的声明语法如下：

> var 数组变量名 [元素数量]Type

如下所示示例，展示了大致数组的定义使用方法：

```go
	var a [3]int  //定义三个整数的数组
	var b [3]int = [3]int{1,2,3} //定义数组并初始化
	c := [...]string{"hello", "word", "ok"} //在数组的定义中，如果在数组长度的位置出现“...”省略号，则表示数组的长度是根据初始化值的个数来计算
	fmt.Println(len(b)) // len函数获取数组长度
	fmt.Println(a[0])
	fmt.Println(c[0])

	// range 可以用来遍历数组
	for _,v := range c{
		fmt.Println(v)
	}

	var x [2][4] int = [2][4]int{{1,2,3,4}, {4,5,6,7}} //多维数组定义并初始化
	fmt.Println(x[0][1])
	fmt.Println(x[0][3])
```

### 切片

切片（slice）是建立在数组之上的更方便，更灵活，更强大的数据结构。切片并不存储任何元素而只是对现有数组的引用。

#### **创建切片**

1. 从数组或者切片生成创建新的切片

   从连续内存区域生成切片是常见的操作，格式如下：

   > slice [开始位置 : 结束位置]

   语法说明如下：

   - slice：表示目标切片对象；
   - 开始位置：对应目标切片对象的索引；
   - 结束位置：对应目标切片的结束索引。

   常见的操作，可以从下面代码中学习:

   ```go
   	a := [5]int{1,2,3,4,5}
   	slice_one := a[0:2] //从数组或者切片中生成切片
   	fmt.Println(slice_one) // [1,2]
   	fmt.Println(a[:]) //输出原切片
   	fmt.Println(a[2:]) //从第三个值开始到最后一个值
   	fmt.Println(a[0:0])//清空切片
   ```

2. 直接声明新切片

   除了可以从原有的数组或者切片中生成切片外，也可以声明一个新的切片，每一种类型都可以拥有其切片类型，表示多个相同类型元素的连续集合，因此切片类型也可以被声明，切片类型声明格式如下：

   > var name []Type

   

   下面代码展示了切片声明的使用过程：

   ```go
   	var strList []string //声明切片
   	var numListEmpty = []int{} //声明并且初始化一个空切片
   	var b  = []int{1,2,3} //简化声明并初始化切片
   	var c []string = []string{"hello", "world"} //声明并初始化切片
   ```

3. make()函数构造切片

   如果需要动态地创建一个切片，可以使用 make() 内建函数，格式如下：

   > make( []Type, size, cap )

其中 Type 是指切片的元素类型，size 指的是为这个类型分配多少个元素，cap 为预分配的元素数量，这个值设定后不影响 size，只是能提前分配空间，降低多次分配空间造成的性能问题。

下面代码展示了切片声明的使用过程：

```go
	makeSlice := make([]int, 5, 5) //使用make创建一个切片
	fmt.Println(makeSlice)
	makeSlice = []int{1,2,3,4,5} //赋值切片
	fmt.Println(makeSlice)
```

#### **修改切片**

```go
	modifySlice := []string{"hello", "china", "!"}
	for k,v := range modifySlice{
		if v == "!"{
			modifySlice[k] = "!!!"
		}
	}
	fmt.Println(modifySlice) //[hello china !!!]
```

上述代码遍历了切片，并且将 ！值修改为 ！！！

#### **追加切片元素**

Go语言的内建函数 append() 可以为切片动态添加元素，代码如下所示：

```go
	var addSlice []int
	addSlice = append(addSlice, 1) //追加1个元素
	fmt.Println(addSlice) // [1]
	addSlice = append(addSlice, 2, 3, 4) //追加3个元素
	fmt.Println(addSlice) // [1 2 3 4]
	addSlice = append(addSlice, []int{5,6,7,8,9}...) // 追加一个切片，切片需要解包
	fmt.Println(addSlice) //[1 2 3 4 5 6 7 8 9]
```

除了在切片的尾部追加，我们还可以在切片的开头添加元素：

```go
	var headAddSlice []int
	headAddSlice = append([]int{0, 1, 2}, headAddSlice...) // 将headAddSlice切片加到其他切片尾部，巧妙的变幻为头部添加功能，注意append第一个参数必须为切片
	fmt.Println(headAddSlice) // [0 1 2]
	headAddSlice = append([]int{3,4}, headAddSlice...)
	fmt.Println(headAddSlice) // [3 4 0 1 2]
```

需要注意的是，在使用 append() 函数为切片动态添加元素时，如果空间不足以容纳足够多的元素，切片就会进行“扩容”，此时新切片的长度会发生改变。在切片开头添加元素一般都会导致内存的重新分配，而且会导致已有元素全部被复制 1 次，因此，从切片的开头添加元素的性能要比从尾部追加元素的性能差很多。切片在扩容时，容量的扩展规律是按容量的 2 倍数进行扩充，例如 1、2、4、8、16……，代码如下：

```go
	var numbers []int
	for i:=0; i<10; i++ {
		numbers = append(numbers, i)
		fmt.Printf("len:%d  cap:%d pointer:%p\n", len(numbers), cap(numbers), numbers) //使用函数 cap()查看切片的容量情况
	}
代码输出如下：
len:1  cap:1 pointer:0xc000056270
len:2  cap:2 pointer:0xc000056290
len:3  cap:4 pointer:0xc000054200
len:4  cap:4 pointer:0xc000054200
len:5  cap:8 pointer:0xc000082100
len:6  cap:8 pointer:0xc000082100
len:7  cap:8 pointer:0xc000082100
len:8  cap:8 pointer:0xc000082100
len:9  cap:16 pointer:0xc00008a000
len:10  cap:16 pointer:0xc00008a000
```

通过查看代码输出，可以发现一个有意思的规律：切片长度 len 并不等于切片的容量 cap。往一个切片中不断添加元素的过程，类似于公司搬家，公司发展初期，资金紧张，人员很少，所以只需要很小的房间即可容纳所有的员工，随着业务的拓展和收入的增加就需要扩充工位，但是办公地的大小是固定的，无法改变，因此公司只能选择搬家，每次搬家就需要将所有的人员转移到新的办公点。

#### **切片复制**

Go语言的内置函数 copy() 可以将一个数组切片复制到另一个数组切片中，如果加入的两个数组切片不一样大，就会按照其中较小的那个数组切片的元素个数进行复制。

copy() 函数的使用格式如下：

> copy( destSlice, srcSlice []T) int

其中 srcSlice 为数据来源切片，destSlice 为复制的目标（也就是将 srcSlice 复制到 destSlice），目标切片必须分配过空间且足够承载复制的元素个数，并且来源和目标的类型必须一致，copy() 函数的返回值表示**实际发生复制的元素个数**。

下面的代码展示了使用 copy() 函数将一个切片复制到另一个切片的过程：

```go
	slice1 := []int{1, 2, 3, 4, 5, 6}
	slice2 := []int{7, 8, 9, 10}
	slice3 := []int{1, 2, 3, 4, 5, 6, 7}
	slice4 := []int{8, 9, 10}
	copy(slice2, slice1)
	fmt.Println(slice2) // [1 2 3]
	copy(slice3, slice4)
	fmt.Println(slice3) // [8 9 10 4 5 6 7]
```

虽然通过循环复制切片元素更直接，不过内置的 copy() 函数使用起来更加方便，copy() 函数的第一个参数是要复制的目标 slice，第二个参数是源 slice，两个 slice 可以共享同一个底层数组，甚至有重叠也没有问题。

**需要注意的是 :  copy复制的只是切片的副本，并不会对源数据有影响**

#### **删除切片元素**

Go语言并没有对删除切片元素提供专用的语法或者接口，需要使用切片本身的特性来删除元素，根据要删除元素的位置有三种情况，分别是从开头位置删除、从中间位置删除和从尾部删除，其中删除切片尾部的元素速度最快。

对于切片删除，下面通过几个例子进行说明：

```go
	// 从开头位置删除
	delSlice := []int{1, 2, 3, 4, 5, 6, 7, 8}
	delSlice = delSlice[1:] // 删除第一个元素
	fmt.Println(delSlice)   // [2 3 4 5 6 7 8]
	delSlice2 := []int{1, 2, 3, 4}
	delSlice2 = append(delSlice2[0:0], delSlice2[1:]...) //删除第一个元素的另外一种实现
	fmt.Println(delSlice2)                               //[2 3 4]
	delSlice3 := []int{1, 2, 3, 4, 5, 6, 7}
	delSlice3 = delSlice3[:copy(delSlice3, delSlice3[1:])] //使用copy方式删除一个元素的实现

	//从尾部删除
	delSlice4 := []int{1, 2, 3, 4, 5, 6, 7}
	delSlice4 = delSlice4[:len(delSlice4)-1] //删除尾部第一个元素

	//从中间删除
	delSlice5 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	delSlice5 = append(delSlice5[:4], delSlice5[5:]...) // 删除第四个元素
	fmt.Println(delSlice5)  //[1 2 3 4 6 7 8 9 10]
```



#### **多维切片**

Go语言中同样允许使用多维切片，声明一个多维数组的语法格式如下：

> var sliceName [][]...[]sliceType

同数组一样，切片也可以有多个维度。具体用发参考以下代码:

```go
	var moreSlice [][]string = [][]string{
		{"C", "C++"},
		{"GO", "RUST"},
		{"PHP"},
	}
	fmt.Println(moreSlice) // [[C C++] [GO RUST] [PHP]]
```

### Map

map 是一种特殊的[数据结构](http://c.biancheng.net/data_structure/)，一种元素对（pair）的无序集合，pair 对应一个 key（索引）和一个 value（值），所以这个结构也称为关联数组或字典，这是一种能够快速寻找值的理想结构，给定 key，就可以迅速找到对应的 value。map 是引用类型，可以使用如下方式声明：

> var mapname map[keytype]valuetype

其中：

- mapname 为 map 的变量名。
- keytype 为键类型。
- valuetype 是键对应的值类型。

#### **Map创建**

我们通过例子来看下Map创建的示例:

```go
	var mapList map[int]int           //声明map
	mapList = map[int]int{1: 1, 2: 2} //赋值map
	fmt.Println(mapList) // map[1:1 2:2]
	var mapList1 map[int]string = map[int]string{1: "hello", 2: "xjx"} //声明并赋值map
	fmt.Println(mapList1) //map[1:hello 2:xjx]
	mapList2 := make(map[int][]int) //使用make创建map
	slice1 := []int{1,2,3,4}
	mapList2 = map[int][]int{1:slice1} //将切片作为值传值给map
	fmt.Println(mapList2) // map[1:[1 2 3 4]]
	mapList3 := map[string]string{"one":"oneData", "two":"twoData", "three":"hello"} //声明并赋值map
	mapList3["one"] = "oneone" // 修改map值
	fmt.Println(mapList3) // map[one:oneone three:hello two:twoData]
```

**特别需要注意的是 map并不能通过new去创建！**

#### **Map遍历元素**

map 的遍历过程使用 for range 循环完成，代码如下：

```go
	mapNew := map[int]string{1: "one", 2: "two"}
	for k, v := range mapNew {
		fmt.Printf("k:%d v:%s\n", k, v)
	}
```

#### **Map元素删除与清空**

使用 delete() 内建函数从 map 中删除一组键值对，delete() 函数的格式如下：

> delete(map, 键)

从 map 中删除一组键值对可以通过下面的代码来学习：

```go
	delMap := map[string]string{"one":"one", "two":"two"}
	delete(delMap, "one") // 删除one键值的元素
	fmt.Println(delMap) //map[two:two]
	delMap = make(map[string]string) // 由于go并没有清空map的函数，所以只能用重新创建一个map来覆盖之前map方案处理
	fmt.Println(delMap) // map[]
```

需要注意的是与[切片]一样，map 是引用类型。当一个 map 赋值给一个新的变量，它们都指向同一个内部数据结构。因此改变其中一个也会反映到另一个：

```go
	animalMap := map[string]string{
		"dog" : "狗",
		"cat" : "猫",
		"pig" : "猪",
		"duck" : "鸭",
	}
	newAnimalMap := animalMap
	delete(newAnimalMap, "pig")
	fmt.Println(animalMap) // map[cat:猫 dog:狗 duck:鸭]
```

### **列表**

列表是一种非连续的存储容器，由多个节点组成，节点通过一些变量记录彼此之间的关系，列表有多种实现方法，如单链表、双链表等。链表是链式的[存储结构]，链表通过指针来连接元素与元素，数组则是把所有元素按次序依次存储。链表的形式有单链表、双链表，循环链表等。具体的链表知识可以查看[百度百科](https://baike.baidu.com/item/%E9%93%BE%E8%A1%A8/9794473?fr=aladdin)，链表相对于切片和数组来说，具有方便数据的删除、插入，长度可变，扩展性好，内存利用率高等优点。

#### **列表创建**

- 通过 container/list 包的 New() 函数初始化 list

  > 变量名 := list.New()

  

- 通过 var 关键字声明初始化 list

  > var 变量名 list.List

#### **列表其他操作**

链表的其他操作，我们通过代码来熟悉熟悉：

```go
	var list1 list.List // list声明方式1
	list2 := list.New() // list声明方式2
	list3 := list.New()

	list1.PushBack("hello") // 添加列表元素到尾部
	list1.PushFront("say")  // 添加列表元素到头部
	element := list1.PushBack("one")
	list1.InsertAfter("fff", element)  //在element点后插入元素
	list1.InsertBefore("ttt", element) //在element点前插入元素

	list2.PushBack("list2")
	list1.PushBackList(list2) //在列表list1后插入列表list2
	list3.PushBack("three")
	list1.PushFrontList(list3) // 在列表list1前插入列表list3

	list1.Remove(element) // 移除element位置

	// 遍历列表
	for i := list1.Front(); i != nil; i = i.Next() {
		fmt.Println(i.Value)
	}
```



### 指针

指针是存储一个变量的内存地址的变量。

Go语言指针与Java C语言指针有很大不同，Go语言只提供了指针允许你控制特定集合的数据结构、分配的数量以及内存访问模式，并不参与指针运算。

指针（pointer）在Go语言中可以被拆分为两个核心概念：

- 类型指针，允许对这个指针类型的数据进行修改，传递数据可以直接使用指针，而无须拷贝数据，类型指针不能进行偏移和运算。
- 切片，由指向起始元素的原始指针、元素数量和容量组成

受益于这样的约束和拆分，Go语言的指针类型变量即拥有指针高效访问的特点，又不会发生指针偏移，从而避免了非法修改关键性数据的问题。同时，垃圾回收也比较容易对不会发生偏移的指针进行检索和回收。切片比原始指针具备更强大的特性，而且更为安全。切片在发生越界时，运行时会报出宕机，并打出堆栈，而原始指针只会崩溃。

#### **指针声明以及赋值**

通过类型作为前缀来定义一个指针’*’

> var name *Type

或者也可以通过new模式去创建一个指针

> name := new(Type)

一般可以这样写，参考以下代码：

```go
	i := 1
	var p *int // 定义指针方式一
	p2 := new(int) //定义指针方式二
	p = &i
	p2 = &i
```

从指针获取值是通过在指针变量前置’*’ 实现的，代码如下：

```go
	a := 5
	pointer1 := &a
	fmt.Println(*pointer1) // 5 指针取值
	*pointer1 = 6 // 通过指针修改 a值
	fmt.Println(a) // 6
```

#### **Go指针不支持运算**

```go
	b := 1
	pointer2 := &b
	pointer2++
```

上面程序会报  invalid operation: pointer2++ (non-numeric type *int) 错误，证明Go指正是不支持参与运算的。



