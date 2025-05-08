package zset

import (
	"math/rand"
	"time"
)

// SKIPLIST_MAXLEVEL 定义跳跃表的最大层数。
const SKIPLIST_MAXLEVEL = 32

// SKIPLIST_P 定义跳跃表节点增加层级的概率。
const SKIPLIST_P = 0.25

// 跳跃表节点
type skiplistNode struct {
	ele      string          // 元素值
	score    float64         // 分数
	backward *skiplistNode   // 后向指针
	level    []skiplistLevel // 层级数组
}

// 跳跃表层级
type skiplistLevel struct {
	forward *skiplistNode // 前向指针
	span    uint64        // 跨度
}

// 跳跃表
type skiplist struct {
	header *skiplistNode // 头节点
	tail   *skiplistNode // 尾节点
	length uint64        // 节点数量
	level  int           // 当前最大层级
}

// ZSet 有序集合，结合哈希表和跳跃表实现。
type ZSet struct {
	dict map[string]float64 // 哈希表，映射元素到分数
	zsl  *skiplist          // 跳跃表，按分数排序元素
}

// 初始化随机数生成器
func init() {
	rand.Seed(time.Now().UnixNano())
}

// createNode 创建一个新的跳跃表节点。
// level: 节点的层级。
// score: 节点的分数。
// ele: 节点的元素值。
// 返回新创建的跳跃表节点指针。
func createNode(level int, score float64, ele string) *skiplistNode {
	node := &skiplistNode{
		ele:   ele,
		score: score,
		level: make([]skiplistLevel, level),
	}
	return node
}

// createSkiplist 创建一个新的跳跃表。
// 返回新创建的跳跃表指针。
func createSkiplist() *skiplist {
	sl := &skiplist{
		level:  1,
		length: 0,
	}
	sl.header = createNode(SKIPLIST_MAXLEVEL, 0, "")
	for j := 0; j < SKIPLIST_MAXLEVEL; j++ {
		sl.header.level[j].forward = nil
		sl.header.level[j].span = 0
	}
	sl.header.backward = nil
	sl.tail = nil
	return sl
}

// NewZSet 创建一个新的有序集合 ZSet。
// 返回新创建的 ZSet 指针。
func NewZSet() *ZSet {
	return &ZSet{
		dict: make(map[string]float64),
		zsl:  createSkiplist(),
	}
}

// randomLevel 随机生成一个跳跃表节点的层级。
// 返回生成的层级。
func randomLevel() int {
	level := 1
	for rand.Float64() < SKIPLIST_P && level < SKIPLIST_MAXLEVEL {
		level++
	}
	return level
}

// insert 向跳跃表中插入一个新节点。
// score: 节点的分数。
// ele: 节点的元素值。
// 返回新插入的节点指针。
func (sl *skiplist) insert(score float64, ele string) *skiplistNode {
	update := make([]*skiplistNode, SKIPLIST_MAXLEVEL)
	rank := make([]uint64, SKIPLIST_MAXLEVEL)

	// 查找插入位置
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		if i == sl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		for x.level[i].forward != nil &&
			(x.level[i].forward.score < score ||
				(x.level[i].forward.score == score && x.level[i].forward.ele < ele)) {
			rank[i] += x.level[i].span
			x = x.level[i].forward
		}
		update[i] = x
	}

	// 随机生成新节点的层级
	level := randomLevel()

	// 如果新节点的层级大于当前跳跃表的层级
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.header
			update[i].level[i].span = sl.length
		}
		sl.level = level
	}

	// 创建新节点
	x = createNode(level, score, ele)

	// 插入节点到跳跃表
	for i := 0; i < level; i++ {
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x

		// 更新跨度
		x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	// 更新高于新节点层级的节点跨度
	for i := level; i < sl.level; i++ {
		update[i].level[i].span++
	}

	// 设置后向指针
	if update[0] == sl.header {
		x.backward = nil
	} else {
		x.backward = update[0]
	}

	if x.level[0].forward != nil {
		x.level[0].forward.backward = x
	} else {
		sl.tail = x
	}

	sl.length++
	return x
}

// Add 向 ZSet 中添加或更新元素。
// ele: 要添加的元素。
// score: 元素的分数。
// 如果元素是新添加的，返回 true；如果元素已存在且分数被更新，返回 false。
func (z *ZSet) Add(ele string, score float64) bool {
	// 检查元素是否已存在
	oldScore, exists := z.dict[ele]

	// 如果元素已存在且分数相同，不做任何操作
	if exists && oldScore == score {
		return false
	}

	// 如果元素已存在，先从跳跃表中删除
	if exists {
		z.zsl.delete(oldScore, ele)
	}

	// 插入新元素到跳跃表
	//node := z.zsl.insert(score, ele)
	z.zsl.insert(score, ele)

	// 更新哈希表
	z.dict[ele] = score

	return !exists
}

// delete 从跳跃表中删除指定分数和元素的节点。
// score: 节点的分数。
// ele: 节点的元素值。
// 如果成功删除，返回 true；否则返回 false。
func (sl *skiplist) delete(score float64, ele string) bool {
	update := make([]*skiplistNode, SKIPLIST_MAXLEVEL)

	// 查找要删除的节点
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.score < score ||
				(x.level[i].forward.score == score && x.level[i].forward.ele < ele)) {
			x = x.level[i].forward
		}
		update[i] = x
	}

	// 获取可能是要删除的节点
	x = x.level[0].forward

	// 检查是否找到了要删除的节点
	if x != nil && x.score == score && x.ele == ele {
		sl.deleteNode(x, update)
		return true
	}

	return false
}

// deleteNode 删除跳跃表中的指定节点。
// x: 要删除的节点。
// update: 记录需要更新的节点。
func (sl *skiplist) deleteNode(x *skiplistNode, update []*skiplistNode) {
	// 更新前向指针和跨度
	for i := 0; i < sl.level; i++ {
		if update[i].level[i].forward == x {
			update[i].level[i].span += x.level[i].span - 1
			update[i].level[i].forward = x.level[i].forward
		} else {
			update[i].level[i].span--
		}
	}

	// 更新后向指针
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x.backward
	} else {
		sl.tail = x.backward
	}

	// 更新跳跃表的最大层级
	for sl.level > 1 && sl.header.level[sl.level-1].forward == nil {
		sl.level--
	}

	sl.length--
}

// Remove 从 ZSet 中删除指定元素。
// ele: 要删除的元素。
// 如果元素存在并成功删除，返回 true；否则返回 false。
func (z *ZSet) Remove(ele string) bool {
	// 检查元素是否存在
	score, exists := z.dict[ele]
	if !exists {
		return false
	}

	// 从跳跃表中删除
	z.zsl.delete(score, ele)

	// 从哈希表中删除
	delete(z.dict, ele)

	return true
}

// Score 获取 ZSet 中指定元素的分数。
// ele: 要获取分数的元素。
// 返回元素的分数和元素是否存在的标志。
func (z *ZSet) Score(ele string) (float64, bool) {
	score, exists := z.dict[ele]
	return score, exists
}

// Rank 获取 ZSet 中指定元素的排名。
// ele: 要获取排名的元素。
// reverse: 是否按降序排名。
// 返回元素的排名（从 0 开始），如果元素不存在返回 -1。
func (z *ZSet) Rank(ele string, reverse bool) int64 {
	score, exists := z.dict[ele]
	if !exists {
		return -1
	}

	rank := z.zsl.getRank(score, ele)
	if rank == 0 {
		return -1
	}

	// 排名从 0 开始
	rank--

	if reverse {
		return int64(z.zsl.length - rank - 1)
	}
	return int64(rank)
}

// GetByRank 获取 ZSet 中指定排名的元素。
// rank: 要获取的排名（从 0 开始）。
// reverse: 是否按降序排名。
// 返回元素、元素的分数和元素是否存在的标志。
func (z *ZSet) GetByRank(rank int64, reverse bool) (string, float64, bool) {
	if rank < 0 || rank >= int64(z.zsl.length) {
		return "", 0, false
	}

	if reverse {
		rank = int64(z.zsl.length) - 1 - rank
	}

	n := z.zsl.getElementByRank(uint64(rank + 1))
	if n == nil {
		return "", 0, false
	}

	return n.ele, n.score, true
}

// getElementByRank 获取跳跃表中指定排名的节点。
// rank: 要获取的排名（从 1 开始）。
// 返回指定排名的节点指针，如果排名无效返回 nil。
func (sl *skiplist) getElementByRank(rank uint64) *skiplistNode {
	if rank == 0 || rank > sl.length {
		return nil
	}

	var traversed uint64 = 0
	x := sl.header

	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && traversed+x.level[i].span <= rank {
			traversed += x.level[i].span
			x = x.level[i].forward
		}
		if traversed == rank {
			return x
		}
	}

	return nil
}

// getRank 获取元素在跳跃表中的排名。
// score: 元素的分数。
// ele: 元素的值。
// 返回元素的排名（从 1 开始），如果元素不存在返回 0。
func (sl *skiplist) getRank(score float64, ele string) uint64 {
	var rank uint64 = 0
	x := sl.header

	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.score < score ||
				(x.level[i].forward.score == score && x.level[i].forward.ele <= ele)) {
			rank += x.level[i].span
			x = x.level[i].forward
		}

		if x != sl.header && x.score == score && x.ele == ele {
			return rank
		}
	}

	return 0
}

// RangeByScore 按分数范围获取 ZSet 中的元素。
// min: 分数范围的最小值。
// max: 分数范围的最大值。
// offset: 跳过的元素数量。
// count: 要获取的元素数量，-1 表示获取所有符合条件的元素。
// 返回符合条件的元素列表。
func (z *ZSet) RangeByScore(min, max float64, offset, count int64) []struct {
	Member string
	Score  float64
} {
	var result []struct {
		Member string
		Score  float64
	}

	// 找到范围的起始节点
	x := z.zsl.header
	if offset < 0 {
		offset = 0
	}

	// 跳到最小分数位置
	for i := z.zsl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.score < min {
			x = x.level[i].forward
		}
	}

	// 向前移动到第一个匹配的元素
	x = x.level[0].forward

	// 跳过 offset 个元素
	var skipped int64 = 0
	for x != nil && skipped < offset {
		if x.score > max {
			break
		}
		skipped++
		x = x.level[0].forward
	}

	// 收集结果
	var returned int64 = 0
	for x != nil && (count < 0 || returned < count) {
		if x.score > max {
			break
		}

		result = append(result, struct {
			Member string
			Score  float64
		}{
			Member: x.ele,
			Score:  x.score,
		})

		returned++
		x = x.level[0].forward
	}

	return result
}

// Len 获取 ZSet 中元素的数量。
// 返回 ZSet 中元素的数量。
func (z *ZSet) Len() uint64 {
	return z.zsl.length
}
