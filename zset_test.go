package zset

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZSet_Add(t *testing.T) {
	t.Run("add new element", func(t *testing.T) {
		z := NewZSet()
		added := z.Add("element1", 1.0)
		assert.True(t, added)
		assert.Equal(t, uint64(1), z.Len())

		score, exists := z.Score("element1")
		assert.True(t, exists)
		assert.Equal(t, 1.0, score)
	})

	t.Run("add existing element with same score", func(t *testing.T) {
		z := NewZSet()
		z.Add("element1", 1.0)
		added := z.Add("element1", 1.0)
		assert.False(t, added)
		assert.Equal(t, uint64(1), z.Len())
	})

	t.Run("update existing element with new score", func(t *testing.T) {
		z := NewZSet()
		z.Add("element1", 1.0)
		added := z.Add("element1", 2.0)
		assert.False(t, added) // returns false because it's an update
		assert.Equal(t, uint64(1), z.Len())

		score, exists := z.Score("element1")
		assert.True(t, exists)
		assert.Equal(t, 2.0, score)
	})

	t.Run("add multiple elements", func(t *testing.T) {
		z := NewZSet()
		z.Add("element1", 1.0)
		z.Add("element2", 2.0)
		z.Add("element3", 3.0)
		assert.Equal(t, uint64(3), z.Len())

		score, exists := z.Score("element2")
		assert.True(t, exists)
		assert.Equal(t, 2.0, score)
	})

	t.Run("verify skiplist and dict consistency", func(t *testing.T) {
		z := NewZSet()
		z.Add("element1", 1.0)
		z.Add("element2", 2.0)

		// Verify dict
		score1, exists1 := z.dict["element1"]
		assert.True(t, exists1)
		assert.Equal(t, 1.0, score1)

		// Verify skiplist
		member, score, exists := z.GetByRank(0, false)
		assert.True(t, exists)
		assert.Equal(t, "element1", member)
		assert.Equal(t, 1.0, score)
	})

	t.Run("add element with zero score", func(t *testing.T) {
		z := NewZSet()
		added := z.Add("element1", 0.0)
		assert.True(t, added)

		score, exists := z.Score("element1")
		assert.True(t, exists)
		assert.Equal(t, 0.0, score)
	})

	t.Run("add element with negative score", func(t *testing.T) {
		z := NewZSet()
		added := z.Add("element1", -1.0)
		assert.True(t, added)

		score, exists := z.Score("element1")
		assert.True(t, exists)
		assert.Equal(t, -1.0, score)
	})
}

func TestZSet_Remove(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *ZSet
		ele      string
		want     bool
		wantDict map[string]float64
		wantLen  uint64
	}{
		{
			name: "remove existing element",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1.0)
				z.Add("b", 2.0)
				return z
			},
			ele:      "a",
			want:     true,
			wantDict: map[string]float64{"b": 2.0},
			wantLen:  1,
		},
		{
			name: "remove non-existing element",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1.0)
				return z
			},
			ele:      "b",
			want:     false,
			wantDict: map[string]float64{"a": 1.0},
			wantLen:  1,
		},
		{
			name: "remove from empty set",
			setup: func() *ZSet {
				return NewZSet()
			},
			ele:      "a",
			want:     false,
			wantDict: map[string]float64{},
			wantLen:  0,
		},
		{
			name: "remove last element",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1.0)
				return z
			},
			ele:      "a",
			want:     true,
			wantDict: map[string]float64{},
			wantLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := tt.setup()
			got := z.Remove(tt.ele)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantDict, z.dict)
			assert.Equal(t, tt.wantLen, z.Len())
		})
	}
}

func TestZSet_Score(t *testing.T) {
	tests := []struct {
		name           string
		setup          func() *ZSet
		ele            string
		expectedScore  float64
		expectedExists bool
	}{
		{
			name: "element exists",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("member1", 10.5)
				return z
			},
			ele:            "member1",
			expectedScore:  10.5,
			expectedExists: true,
		},
		{
			name: "element does not exist",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("member1", 10.5)
				return z
			},
			ele:            "nonexistent",
			expectedScore:  0,
			expectedExists: false,
		},
		{
			name: "empty set",
			setup: func() *ZSet {
				return NewZSet()
			},
			ele:            "any",
			expectedScore:  0,
			expectedExists: false,
		},
		{
			name: "multiple elements",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("member1", 10.5)
				z.Add("member2", 20.0)
				z.Add("member3", 30.75)
				return z
			},
			ele:            "member2",
			expectedScore:  20.0,
			expectedExists: true,
		},
		{
			name: "zero score",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("zero", 0.0)
				return z
			},
			ele:            "zero",
			expectedScore:  0.0,
			expectedExists: true,
		},
		{
			name: "negative score",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("negative", -5.5)
				return z
			},
			ele:            "negative",
			expectedScore:  -5.5,
			expectedExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := tt.setup()
			score, exists := z.Score(tt.ele)
			assert.Equal(t, tt.expectedScore, score)
			assert.Equal(t, tt.expectedExists, exists)
		})
	}
}

func TestZSet_Rank(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *ZSet
		ele      string
		reverse  bool
		expected int64
	}{
		{
			name: "element exists - forward rank",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1.0)
				z.Add("b", 2.0)
				z.Add("c", 3.0)
				return z
			},
			ele:      "b",
			reverse:  false,
			expected: 1,
		},
		{
			name: "element exists - reverse rank",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1.0)
				z.Add("b", 2.0)
				z.Add("c", 3.0)
				return z
			},
			ele:      "b",
			reverse:  true,
			expected: 1,
		},
		{
			name: "element not exists",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1.0)
				z.Add("b", 2.0)
				return z
			},
			ele:      "c",
			reverse:  false,
			expected: -1,
		},
		{
			name: "empty set",
			setup: func() *ZSet {
				return NewZSet()
			},
			ele:      "a",
			reverse:  false,
			expected: -1,
		},
		{
			name: "single element - forward rank",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1.0)
				return z
			},
			ele:      "a",
			reverse:  false,
			expected: 0,
		},
		{
			name: "single element - reverse rank",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1.0)
				return z
			},
			ele:      "a",
			reverse:  true,
			expected: 0,
		},
		{
			name: "duplicate scores - forward rank",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1.0)
				z.Add("b", 1.0)
				z.Add("c", 1.0)
				return z
			},
			ele:      "b",
			reverse:  false,
			expected: 1,
		},
		{
			name: "duplicate scores - reverse rank",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1.0)
				z.Add("b", 1.0)
				z.Add("c", 1.0)
				return z
			},
			ele:      "b",
			reverse:  true,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := tt.setup()
			actual := z.Rank(tt.ele, tt.reverse)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestZSet_GetByRank(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *ZSet
		rank      int64
		reverse   bool
		wantEle   string
		wantScore float64
		wantOk    bool
	}{
		{
			name: "empty set",
			setup: func() *ZSet {
				return NewZSet()
			},
			rank:    0,
			reverse: false,
			wantOk:  false,
		},
		{
			name: "rank out of range (negative)",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1.0)
				return z
			},
			rank:    -1,
			reverse: false,
			wantOk:  false,
		},
		{
			name: "rank out of range (too large)",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1.0)
				return z
			},
			rank:    1,
			reverse: false,
			wantOk:  false,
		},
		{
			name: "valid rank forward",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1.0)
				z.Add("b", 2.0)
				z.Add("c", 3.0)
				return z
			},
			rank:      1,
			reverse:   false,
			wantEle:   "b",
			wantScore: 2.0,
			wantOk:    true,
		},
		{
			name: "valid rank reverse",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1.0)
				z.Add("b", 2.0)
				z.Add("c", 3.0)
				return z
			},
			rank:      1,
			reverse:   true,
			wantEle:   "b",
			wantScore: 2.0,
			wantOk:    true,
		},
		{
			name: "first element forward",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1.0)
				z.Add("b", 2.0)
				return z
			},
			rank:      0,
			reverse:   false,
			wantEle:   "a",
			wantScore: 1.0,
			wantOk:    true,
		},
		{
			name: "first element reverse",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1.0)
				z.Add("b", 2.0)
				return z
			},
			rank:      0,
			reverse:   true,
			wantEle:   "b",
			wantScore: 2.0,
			wantOk:    true,
		},
		{
			name: "last element forward",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1.0)
				z.Add("b", 2.0)
				return z
			},
			rank:      1,
			reverse:   false,
			wantEle:   "b",
			wantScore: 2.0,
			wantOk:    true,
		},
		{
			name: "last element reverse",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1.0)
				z.Add("b", 2.0)
				return z
			},
			rank:      1,
			reverse:   true,
			wantEle:   "a",
			wantScore: 1.0,
			wantOk:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := tt.setup()
			ele, score, ok := z.GetByRank(tt.rank, tt.reverse)
			assert.Equal(t, tt.wantOk, ok)
			if tt.wantOk {
				assert.Equal(t, tt.wantEle, ele)
				assert.Equal(t, tt.wantScore, score)
			}
		})
	}
}

func TestRangeByScore(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *ZSet
		min      float64
		max      float64
		offset   int64
		count    int64
		expected []struct {
			Member string
			Score  float64
		}
	}{
		{
			name: "empty set",
			setup: func() *ZSet {
				return NewZSet()
			},
			min:    0,
			max:    10,
			offset: 0,
			count:  10,
			expected: []struct {
				Member string
				Score  float64
			}(nil),
		},
		{
			name: "single element in range",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 5)
				return z
			},
			min:    0,
			max:    10,
			offset: 0,
			count:  10,
			expected: []struct {
				Member string
				Score  float64
			}{
				{"a", 5},
			},
		},
		{
			name: "single element out of range",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 15)
				return z
			},
			min:    0,
			max:    10,
			offset: 0,
			count:  10,
			expected: []struct {
				Member string
				Score  float64
			}(nil),
		},
		{
			name: "multiple elements with offset",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1)
				z.Add("b", 2)
				z.Add("c", 3)
				z.Add("d", 4)
				z.Add("e", 5)
				return z
			},
			min:    1,
			max:    5,
			offset: 2,
			count:  2,
			expected: []struct {
				Member string
				Score  float64
			}{
				{"c", 3},
				{"d", 4},
			},
		},
		{
			name: "multiple elements with negative count",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1)
				z.Add("b", 2)
				z.Add("c", 3)
				return z
			},
			min:    1,
			max:    3,
			offset: 0,
			count:  -1,
			expected: []struct {
				Member string
				Score  float64
			}{
				{"a", 1},
				{"b", 2},
				{"c", 3},
			},
		},
		{
			name: "multiple elements with partial range",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1)
				z.Add("b", 2)
				z.Add("c", 3)
				z.Add("d", 4)
				z.Add("e", 5)
				return z
			},
			min:    2,
			max:    4,
			offset: 0,
			count:  10,
			expected: []struct {
				Member string
				Score  float64
			}{
				{"b", 2},
				{"c", 3},
				{"d", 4},
			},
		},
		{
			name: "negative offset",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1)
				z.Add("b", 2)
				return z
			},
			min:    1,
			max:    2,
			offset: -1,
			count:  10,
			expected: []struct {
				Member string
				Score  float64
			}{
				{"a", 1},
				{"b", 2},
			},
		},
		{
			name: "count larger than available elements",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1)
				z.Add("b", 2)
				return z
			},
			min:    1,
			max:    2,
			offset: 0,
			count:  5,
			expected: []struct {
				Member string
				Score  float64
			}{
				{"a", 1},
				{"b", 2},
			},
		},
		{
			name: "min equals max",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("a", 1)
				z.Add("b", 2)
				z.Add("c", 2)
				z.Add("d", 3)
				return z
			},
			min:    2,
			max:    2,
			offset: 0,
			count:  10,
			expected: []struct {
				Member string
				Score  float64
			}{
				{"b", 2},
				{"c", 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := tt.setup()
			result := z.RangeByScore(tt.min, tt.max, tt.offset, tt.count)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestZSet_Len(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *ZSet
		expected uint64
	}{
		{
			name: "empty set",
			setup: func() *ZSet {
				return NewZSet()
			},
			expected: 0,
		},
		{
			name: "single element",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("element1", 1.0)
				return z
			},
			expected: 1,
		},
		{
			name: "multiple elements",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("element1", 1.0)
				z.Add("element2", 2.0)
				z.Add("element3", 3.0)
				return z
			},
			expected: 3,
		},
		{
			name: "after removal",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("element1", 1.0)
				z.Add("element2", 2.0)
				z.Remove("element1")
				return z
			},
			expected: 1,
		},
		{
			name: "duplicate elements",
			setup: func() *ZSet {
				z := NewZSet()
				z.Add("element1", 1.0)
				z.Add("element1", 2.0) // Should update score but not increase length
				return z
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := tt.setup()
			assert.Equal(t, tt.expected, z.Len())
		})
	}
}
