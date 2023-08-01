package go_sorted_set

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortedSet(t *testing.T) {
	assert := assert.New(t)

	zset := New()

	// ZADD: Add elements to the zset
	zset.ZAdd("key", 1, []byte("member1")) // Add "member1" with score 1
	zset.ZAdd("key", 2, []byte("member2")) // Add "member2" with score 2
	zset.ZAdd("key", 3, []byte("member3")) // Add "member3" with score 3
	zset.ZAdd("key", 4, []byte("member4")) // Add "member4" with score 4
	zset.ZAdd("key", 5, []byte("member5")) // Add "member5" with score 5

	// ZSCORE: Get the score of an element
	score, ok := zset.ZScore([]byte("member3"))
	assert.True(ok, "'member3' should exist in the zset")
	assert.Equal(3, score, "Score of 'member3' should be 3")

	// ZRANK: Retrieve the rank of an element
	assert.Equal(3, zset.ZRank([]byte("member3")), "Rank of 'member3' should be 3")

	// ZRANGEBYSCORE: Retrieve elements within a score range
	range1 := zset.ZRangeByScore(2, 4) // Retrieve elements with scores 2 to 4
	assert.Equal(3, len(range1), "There should be 3 elements with scores between 2 and 4")

	// ZREMRANGEBYRANK: Remove elements within a rank range
	removed := zset.ZRemRangeByRank(2, 4) // Remove elements with ranks 2 to 4
	assert.Equal(3, removed, "Should have removed 3 elements")

	// ZCARD: Get the number of elements in the zset
	assert.Equal(2, zset.ZCard(), "Zset should have 2 elements after removal")
}
