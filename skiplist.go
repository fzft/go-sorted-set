package go_sorted_set

import (
	"bytes"
	"math/rand"
)

/**
skip list diagram:
Level 2:  H <-----------------> 5
           |                   ^
           |                   |
Level 1:  H <---------> 3 <-----> 5 <--------> 7
           |           ^         ^            ^
           |           |         |            |
Level 0:  H <--> 1 <--> 3 <--> 4 <--> 5 <--> 6 <--> 7 <--> 8

the arrow pointing to the right(-->) represents the forward pointer of the node
at the bottom level, the forward pointer points to the next node in the same level
At higher levels, forward pointers skip over some nodes to allow for faster traversal of the list

The arrows pointing to the left (<--) represent the backward pointers. These make it possible to traverse the list in the opposite direction

H is the head node. It doesn't have a backward pointer, because there's no node before it
*/

const MAX_LEVEL = 16

type Level struct {
	forward *Node
	span    int
}

type Node struct {
	score    int
	member   []byte
	level    []Level
	backward *Node
}

type SkipList struct {
	level  int
	header *Node
	tail   *Node
	length int
}

func newNode(level int, score int, ele []byte) *Node {
	return &Node{
		score:  score,
		member: ele,
		level:  make([]Level, level),
	}
}

func newSkipList() *SkipList {
	zsl := &SkipList{
		level:  1,
		length: 0,
		header: newNode(MAX_LEVEL, 0, nil),
	}
	for i := 0; i < MAX_LEVEL; i++ {
		zsl.header.level[i].forward = nil
		zsl.header.level[i].span = 0
	}
	zsl.header.backward = nil
	zsl.tail = nil
	return zsl
}

func (s *SkipList) randomLevel() (level int) {
	for level = 1; rand.Int31()&1 == 1; level++ {
	}
	if level < MAX_LEVEL {
		return level
	}
	return MAX_LEVEL
}

// Insert inserts a new node into the skip list
func (s *SkipList) Insert(score int, ele []byte) *Node {
	update := make([]*Node, MAX_LEVEL)
	rank := make([]int, MAX_LEVEL)

	x := s.header
	for i := s.level - 1; i >= 0; i-- {
		if i == s.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		if x.level[i].forward != nil {
			for x.level[i].forward != nil && (x.level[i].forward.score < score ||
				(x.level[i].forward.score == score &&
					bytes.Compare(x.level[i].forward.member, ele) < 0)) {
				rank[i] += x.level[i].span
				x = x.level[i].forward
			}
		}
		update[i] = x
	}

	level := s.randomLevel()
	if level > s.level {
		for i := s.level; i < level; i++ {
			rank[i] = 0
			update[i] = s.header
			update[i].level[i].span = s.length
		}
		s.level = level
	}

	x = newNode(level, score, ele)
	for i := 0; i < level; i++ {
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x

		x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	for i := level; i < s.level; i++ {
		update[i].level[i].span++
	}

	if update[0] == s.header {
		x.backward = nil
	} else {
		x.backward = update[0]
	}

	if x.level[0].forward != nil {
		x.level[0].forward.backward = x
	} else {
		s.tail = x
	}
	s.length++

	return x
}

// Delete deletes the node with the given score and element
func (s *SkipList) Delete(score int, ele []byte) bool {
	update := make([]*Node, MAX_LEVEL)

	x := s.header
	for i := s.level - 1; i >= 0; i-- {
		if x.level[i].forward != nil {
			for x.level[i].forward != nil && (x.level[i].forward.score < score ||
				(x.level[i].forward.score == score &&
					bytes.Compare(x.level[i].forward.member, ele) < 0)) {
				x = x.level[i].forward
			}
		}
		update[i] = x
	}

	x = x.level[0].forward
	if x != nil && score == x.score && bytes.Compare(x.member, ele) == 0 {
		s.deleteNode(x, update)
		return true
	}

	return false
}

func (s *SkipList) deleteNode(x *Node, update []*Node) {
	for i := 0; i < s.level; i++ {
		if update[i].level[i].forward == x {
			update[i].level[i].span += x.level[i].span - 1
			update[i].level[i].forward = x.level[i].forward
		} else {
			update[i].level[i].span -= 1
		}
	}

	if x.level[0].forward != nil {
		x.level[0].forward.backward = x.backward
	} else {
		s.tail = update[0]
	}

	for s.level > 1 && s.header.level[s.level-1].forward == nil {
		s.level--
	}

	s.length--
}

// UpdateScore updates the score of the node with the given score and element
func (s *SkipList) UpdateScore(curScore, newScore int, ele []byte) *Node {
	if s.Delete(curScore, ele) {
		return s.Insert(newScore, ele)
	}
	return nil
}

// GetRangeByScore returns all nodes with score between Min and Max from the skip list
func (s *SkipList) GetRangeByScore(rangex ZrangeSpec) []*Node {
	x := s.header
	result := []*Node{}

	// Traverse to the first node with score >= Min.
	for i := s.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && (x.level[i].forward.score < rangex.Min ||
			(rangex.Maxex && x.level[i].forward.score == rangex.Min)) {
			x = x.level[i].forward
		}
	}

	x = x.level[0].forward

	// Add all nodes with score <= Max to the result.
	for x != nil && x.score <= rangex.Max &&
		!(rangex.Maxex && x.score == rangex.Max) {
		result = append(result, x)
		x = x.level[0].forward
	}

	return result
}

// DeleteRangeByScore deletes all nodes with score between Min and Max from the skip list
func (s *SkipList) DeleteRangeByScore(rangex ZrangeSpec, dict map[string]int) int {
	update := make([]*Node, MAX_LEVEL)

	x := s.header
	for i := s.level - 1; i >= 0; i-- {
		if x.level[i].forward != nil {
			for x.level[i].forward != nil && (x.level[i].forward.score < rangex.Min ||
				(x.level[i].forward.score == rangex.Min && rangex.Minx)) {
				x = x.level[i].forward
			}
		}
		update[i] = x
	}

	x = x.level[0].forward
	removed := 0
	for x != nil && (rangex.Maxex && x.score <= rangex.Max || !rangex.Maxex && x.score < rangex.Max) {
		next := x.level[0].forward
		s.deleteNode(x, update)
		delete(dict, string(x.member))
		removed++
		x = next
	}
	return removed
}

// DeleteRangeByRank deletes all nodes with rank between start and end from the skip list
func (s *SkipList) DeleteRangeByRank(rangex ZrangeSpec, dict map[string]int) int {
	start, end := rangex.Min, rangex.Max

	if start < 0 || start > s.length || end < 0 || (end >= s.length && rangex.Maxex) || start > end {
		return 0
	}

	if rangex.Minx {
		start++
	}
	if !rangex.Maxex {
		end++
	}

	update := make([]*Node, s.level)
	x := s.header

	// Find the node ranks to be updated
	for i := s.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].span <= start {
			start -= x.level[i].span
			end -= x.level[i].span
			x = x.level[i].forward
		}
		update[i] = x
	}

	x = x.level[0].forward
	removed := 0

	// Remove nodes within the rank range
	for x != nil && start <= end {
		next := x.level[0].forward
		s.deleteNode(x, update)
		delete(dict, string(x.member))
		x = next
		start++
		removed++
	}

	return removed
}

// GetRangeByRank returns all nodes with rank between start and end from the skip list
func (s *SkipList) GetRangeByRank(rangex ZrangeSpec) []*Node {
	start, end := rangex.Min, rangex.Max

	if start < 0 || start > s.length || end < 0 || (end >= s.length && rangex.Maxex) || start > end {
		return nil
	}

	if rangex.Minx {
		start++
	}
	if !rangex.Maxex {
		end++
	}

	x := s.header
	traversed := 0
	result := []*Node{}

	// Traverse to the first node of the rank range
	for i := s.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && (traversed+x.level[i].span) < start {
			traversed += x.level[i].span
			x = x.level[i].forward
		}
	}

	// Traverse to the end of the rank range, adding each node to the result
	x = x.level[0].forward
	for x != nil && traversed < end {
		result = append(result, x)
		x = x.level[0].forward
		traversed++
	}

	return result
}

// GetRank returns the rank of the node with the given score and element
func (s *SkipList) GetRank(score int, ele []byte) int {
	x := s.header
	rank := 0

	for i := s.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.score < score ||
				(x.level[i].forward.score == score &&
					bytes.Compare(x.level[i].forward.member, ele) <= 0)) {
			rank += x.level[i].span
			x = x.level[i].forward
		}

		// x might be equal to the searched element.
		if bytes.Equal(x.member, ele) {
			return rank
		}
	}

	return -1
}

// GetElementByRank returns the element with the given rank
func (s *SkipList) GetElementByRank(rank int) *Node {
	if rank < 0 || rank >= s.length {
		return nil
	}

	x := s.header
	traversed := 0

	for i := s.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && (traversed+x.level[i].span) <= rank {
			traversed += x.level[i].span
			x = x.level[i].forward
		}

		if traversed == rank {
			return x
		}
	}

	return nil
}
