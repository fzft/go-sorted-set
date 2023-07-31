package go_sorted_set

import (
	"bytes"
	"math/rand"
)

const MAX_LEVEL = 16
const MAX_SCORE = 1 << 32

type Node struct {
	score   int
	member  []byte
	key     string
	forward []*Node
}

type SkipList struct {
	level int
	head  *Node
}

func newNode(key string, score int, member []byte, level int) *Node {
	return &Node{
		score:   score,
		member:  member,
		key:     key,
		forward: make([]*Node, level),
	}
}

func newSkipList() *SkipList {
	return &SkipList{
		level: 1,
		head:  newNode("", 0, nil, MAX_LEVEL),
	}
}

func (s *SkipList) randomLevel() (level int) {
	for level = 1; rand.Int31()&1 == 1; level++ {
	}
	if level < MAX_LEVEL {
		return level
	}
	return MAX_LEVEL
}

func (s *SkipList) Insert(key string, score int, member []byte) {
	update := make([]*Node, MAX_LEVEL)
	x := s.head
	for i := s.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && (x.forward[i].score < score || (x.forward[i].score == score && bytes.Compare(x.forward[i].member, member) < 0)) {
			x = x.forward[i]
		}
		update[i] = x
	}
	level := s.randomLevel()
	if level > s.level {
		for i := s.level; i < level; i++ {
			update[i] = s.head
		}
		s.level = level
	}
	x = newNode(key, score, member, level)
	for i := 0; i < level; i++ {
		x.forward[i] = update[i].forward[i]
		update[i].forward[i] = x
	}
}

// Search returns the node with the given score,
// Param:
//
//	start: the start score of the range
//	end: the end score of the range, if end == -1, search until the end
func (s *SkipList) Search(start, end int) []*Node {
	var nodes []*Node

	// Find the first node with score >= start
	x := s.head
	for i := s.level; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].score < start {
			x = x.forward[i]
		}
	}
	x = x.forward[0]

	// Traverse until score > end

	if end == -1 {
		end = MAX_SCORE
	}

	for x != nil && x.score <= end {
		nodes = append(nodes, x)
		x = x.forward[0]
	}

	return nodes
}

func (s *SkipList) Delete(score int, member []byte) {
	update := make([]*Node, MAX_LEVEL)
	x := s.head
	for i := s.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && (x.forward[i].score < score || (x.forward[i].score == score && bytes.Compare(x.forward[i].member, member) < 0)) {
			x = x.forward[i]
		}
		update[i] = x
	}
	x = x.forward[0]
	if x != nil && x.score == score && bytes.Compare(x.member, member) == 0 {
		for i := 0; i < s.level; i++ {
			if update[i].forward[i] != x {
				break
			}
			update[i].forward[i] = x.forward[i]
		}
		for s.level > 1 && s.head.forward[s.level-1] == nil {
			s.level--
		}
	}
}

func (s *SkipList) DeleteByScoreRange(start, end int) []*Node {
	var nodes []*Node

	// Find the first node with score >= start
	x := s.head
	for i := s.level; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].score < start {
			x = x.forward[i]
		}
	}
	x = x.forward[0]

	if end == -1 {
		end = MAX_SCORE
	}

	// Traverse until score > end, removing nodes
	for x != nil && x.score <= end {
		next := x.forward[0] // Remember next node before we change the list
		nodes = append(nodes, x)

		// Remove x from the list
		for i := 0; i <= s.level; i++ {
			if s.head.forward[i] != x {
				break
			}
			s.head.forward[i] = x.forward[i]
		}

		// Reduce level of list if necessary
		for s.level > 0 && s.head.forward[s.level] == nil {
			s.level--
		}

		x = next // Move to next node
	}

	return nodes
}

func (s *SkipList) DeleteByScore(score int) []*Node {
	var nodes []*Node

	// Start from the head node
	x := s.head

	// Traverse each level starting from the highest
	for i := s.level; i >= 0; i-- {
		// Keep moving forward while the next node's score is less than the target score
		for x.forward[i] != nil && x.forward[i].score < score {
			x = x.forward[i]
		}

		// If the next node's score is the target score, remove it
		if x.forward[i] != nil && x.forward[i].score == score {
			target := x.forward[i]

			// Remove the target by updating the forward pointer
			x.forward[i] = target.forward[i]

			// Add the target to the nodes slice if it's not already there
			if i == 0 {
				nodes = append(nodes, target)
			}
		}
	}

	return nodes
}

func (s *SkipList) searchMember(score int, member []byte) *Node {
	// Start from the head node
	x := s.head

	// Traverse each level starting from the highest
	for i := s.level; i >= 0; i-- {
		// Keep moving forward while the next node's score is less than the target score
		for x.forward[i] != nil && (x.forward[i].score < score || (x.forward[i].score == score && bytes.Compare(x.forward[i].member, member) < 0)) {
			x = x.forward[i]
		}
	}

	// Move to the next node, which should have a score >= the target score
	x = x.forward[0]

	// Traverse the nodes with the target score and check the member
	for x != nil && x.score == score {
		if bytes.Equal(x.member, member) {
			return x
		}
		x = x.forward[0]
	}

	// If we reach here, the member was not found
	return nil
}
