# zset
Go语言实现的有序集合，参考Redis的ZSet数据结构，结合哈希表和跳表实现高效的有序集合操作。

功能特性

• 有序集合操作：按分数维护元素的排序

• 快速查询：使用哈希表实现 O(1) 复杂度的分数查询

• 高效范围查询：使用跳表实现 O(log n) 复杂度的范围操作

• 排名操作：获取元素排名或按排名获取元素

• 分数范围查询：获取指定分数范围内的元素


数据结构

跳表
• 多层链表结构实现高效遍历

• 使用概率 `SKIPLIST_P = 0.25` 随机生成层级

• 最大层级为 32

• 每个节点包含：

• 元素值(string)

• 分数(float64)

• 后向指针

• 带有跨度的前向指针数组


哈希表
• 将元素字符串映射到其分数

• 用于 O(1) 复杂度的分数查询


API 参考

初始化
```go
zset := NewZSet() // 创建一个新的空有序集合
```

核心操作
```go
// 添加或更新元素及其分数
zset.Add(ele string, score float64) bool

// 移除元素
zset.Remove(ele string) bool

// 获取元素的分数
zset.Score(ele string) (float64, bool)

// 获取元素数量
zset.Len() uint64
```

排名操作
```go
// 获取元素排名(从0开始)
// reverse=false: 升序排列(最低分数排名为0)
// reverse=true: 降序排列(最高分数排名为0)
zset.Rank(ele string, reverse bool) int64

// 按排名获取元素
zset.GetByRank(rank int64, reverse bool) (string, float64, bool)
```

范围操作
```go
// 获取分数在[min, max]范围内的元素
// offset: 要跳过的元素数量
// count: 最多返回的元素数量(-1表示无限制)
zset.RangeByScore(min, max float64, offset, count int64) []struct {
    Member string
    Score  float64
}
```

性能特征

| 操作            | 复杂度       |
|----------------|-------------|
| 添加元素        | O(log n)    |
| 移除元素        | O(log n)    |
| 分数查询        | O(1)        |
| 排名查询        | O(log n)    |
| 按排名获取元素  | O(log n)    |
| 分数范围查询    | O(log n + m)| (m = 范围内元素数量)

使用示例

```go
package main

import (
	"fmt"
	"github.com/beijian128/zset" // 替换为实际的导入路径
)

func main() {
	z := zset.NewZSet()
	
	// 添加元素
	z.Add("Alice", 100)
	z.Add("Bob", 75)
	z.Add("Charlie", 125)
	
	// 获取分数
	if score, ok := z.Score("Bob"); ok {
		fmt.Printf("Bob的分数: %.2f\n", score)
	}
	
	// 获取排名
	rank := z.Rank("Charlie", false)
	fmt.Printf("Charlie的排名: %d\n", rank)
	
	// 范围查询
	results := z.RangeByScore(70, 110, 0, -1)
	for _, res := range results {
		fmt.Printf("%s: %.2f\n", res.Member, res.Score)
	}
	
	// 移除元素
	z.Remove("Alice")
	fmt.Printf("总元素数: %d\n", z.Len())
}
```

实现说明

• 跳表采用概率方法生成层级

• 同时维护前向和后向指针以实现高效遍历

• 每个层级跟踪跨度以支持排名操作

• 哈希表确保 O(1) 的分数查询，跳表维护排序

