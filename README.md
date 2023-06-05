# 队列-Queue

提供各种队列

## 循环队列 - Circular Queue

### 特性

+ 支持在队列末尾插入元素，并按照循环方式覆盖最旧的元素。
+ 支持获取最新插入的元素。
+ 支持按照指定偏移量获取队列中的元素。
+ 支持获取最新的一批元素，并按照插入顺序排序。
+ 支持获取所有元素，并按照插入顺序排序。

### 安装

使用 Go 模块引入该包：

```shell
go get github.com/sumengzs/queue
```

导入包：

```go
import "github.com/sumengzs/queue"
```

### 使用示例

```go
package main

import (
	"fmt"
	"github.com/sumengzs/queue"
)

func main() {
	// 创建一个容量为 5 的循环队列
	q := queue.NewCircularQueue(5)

	// 入队操作
	q.Put("A")
	q.Put("B")
	q.Put("C")

	// 获取最新插入的元素
	lastElement := q.Get()
	fmt.Println("Last Element:", lastElement) // Output: Last Element: C

	// 按照指定偏移量获取元素
	element := q.GetPoint(2)
	fmt.Println("Element at offset 2:", element) // Output: Element at offset 2: A

	// 获取最新的两个元素，并按照插入顺序排序
	latestElements := q.Gets(2)
	fmt.Println("Latest Elements:", latestElements) // Output: Latest Elements: [B C]

	// 获取所有元素，并按照插入顺序排序
	allElements := q.GetAll()
	fmt.Println("All Elements:", allElements) // Output: All Elements: [A B C]
}
```

### 贡献

欢迎贡献代码、提出问题或提供改进建议。请在 GitHub 上提交问题或拉取请求。

### 许可证

该项目基于 Apache 许可证进行分发。有关详细信息，请参阅 LICENSE 文件。
