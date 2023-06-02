# 队列-Queue

提供各种队列

## 循环队列 - Circular Queue

特性 Feature：

+ 队列大小不可变
+ 数据无法删除
+ 支持单个或者批量查询

```go
package main

import (
	"fmt"
	"github.com/sumengzs/queue"
)

func main() {
	cq := queue.NewCircularQueue(1024)
	for i := 0; i < 1024; i++ {
		cq.Put(i)
	}
	fmt.Println(cq.Get()) // 1023
	fmt.Println(cq.GetPoint(4)) //1020
	fmt.Println(cq.Gets(3)) // [0,1,2]
	fmt.Println(cq.GetAll()) // [0,1,2,...,1023]
}
```

