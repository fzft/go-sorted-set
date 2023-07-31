package go_sorted_set

import "testing"

func TestSortedSet(t *testing.T) {
	s := NewSortedSet()

	// Test ZAdd
	s.ZAdd("key1", 1, []byte("member1"))
	s.ZAdd("key1", 2, []byte("member2"))
	s.ZAdd("key1", 3, []byte("member3"))

	// Test ZRange
	nodes := s.ZRange(1, 3)
	if len(nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(nodes))
	}

	// Test ZRem
	s.ZRem("key1", []byte("member1"))
	nodes = s.ZRange(1, 3)
	if len(nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(nodes))
	}

	// Test ZRemByScore
	s.ZRemByScore(2)
	nodes = s.ZRange(1, 3)
	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}

	// Test ZIncrBy
	s.ZIncrBy("key1", 1, []byte("member3"))
	nodes = s.ZRange(1, 4)
	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}
	if nodes[0].score != 4 {
		t.Errorf("Expected score 4, got %d", nodes[0].score)
	}

	// Test ZRemRangeByScore
	s.ZRemRangeByScore(1, 4)
	nodes = s.ZRange(1, 4)
	if len(nodes) != 0 {
		t.Errorf("Expected 0 nodes, got %d", len(nodes))
	}
}
