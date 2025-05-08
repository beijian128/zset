package zset

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewZSet(t *testing.T) {
	zset := NewZSet()
	zset.Add("a", 1.0)
	zset.Add("b", 2.0)
	zset.Add("c", 3.0)

	rk := zset.Rank("a", false)
	assert.Equal(t, int64(0), rk)
	rk = zset.Rank("b", false)
	assert.Equal(t, int64(1), rk)
	rk = zset.Rank("c", false)
	assert.Equal(t, int64(2), rk)

	rk = zset.Rank("a", true)
	assert.Equal(t, int64(2), rk)
	rk = zset.Rank("b", true)
	assert.Equal(t, int64(1), rk)
	rk = zset.Rank("c", true)
	assert.Equal(t, int64(0), rk)

	ele, score, ok := zset.GetByRank(0, false)
	t.Log(ele, score, ok)
	assert.Equal(t, "a", ele)
	assert.Equal(t, 1.0, score)
	assert.Equal(t, true, ok)
	ele, score, ok = zset.GetByRank(1, false)
	t.Log(ele, score, ok)
	assert.Equal(t, "b", ele)
	assert.Equal(t, 2.0, score)
	assert.Equal(t, true, ok)

	rks := zset.RangeByScore(2, 10, 0, 10)
	t.Log(rks)
}
