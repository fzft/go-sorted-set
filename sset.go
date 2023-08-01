package go_sorted_set

type ZrangeSpec struct {
	Min, Max    int
	Minx, Maxex bool // are min or max exclusive?
}

type Zset struct {
	dict     map[string]int
	skiplist *SkipList
}

func New() *Zset {
	return &Zset{dict: make(map[string]int), skiplist: newSkipList()}
}

func (s *Zset) ZAdd(key string, score int, member []byte) {
	if _, ok := s.dict[string(member)]; ok {
		s.ZRem(key, member)
	}
	s.dict[string(member)] = score
	s.skiplist.Insert(score, member)
}

// ZRangeByScore returns all nodes with score between Min and Max from the sorted set
func (s *Zset) ZRangeByScore(min, max int) []*Node {
	return s.skiplist.GetRangeByScore(ZrangeSpec{Min: min, Max: max})
}

// ZRangeByRank returns all nodes with rank between start and end from the sorted set
func (s *Zset) ZRangeByRank(start, end int) []*Node {
	return s.skiplist.GetRangeByRank(ZrangeSpec{Min: start, Max: end})
}

// ZRank returns the rank of the element with the given member in the sorted set
func (s *Zset) ZRank(member []byte) int {
	if score, ok := s.dict[string(member)]; ok {
		return s.skiplist.GetRank(score, member)
	}
	return -1
}

// ZRem removes the element with the given score and member from the sorted set
func (s *Zset) ZRem(_ string, member []byte) bool {
	if score, ok := s.dict[string(member)]; ok {
		delete(s.dict, string(member))
		s.skiplist.Delete(score, member)
		return true
	}
	return false
}

// ZRemRangeByScore deletes all nodes with score between Min and Max from the sorted set
func (s *Zset) ZRemRangeByScore(min, max int) int {
	return s.skiplist.DeleteRangeByScore(ZrangeSpec{Min: min, Max: max}, s.dict)
}

// ZRemRangeByRank deletes all nodes with rank between start and end from the sorted set
func (s *Zset) ZRemRangeByRank(start, end int) int {
	return s.skiplist.DeleteRangeByRank(ZrangeSpec{Min: start, Max: end}, s.dict)
}

// ZCard returns the number of elements in the sorted set
func (s *Zset) ZCard() int {
	return s.skiplist.length
}

// ZScore returns the score of the element with the given member in the sorted set
func (s *Zset) ZScore(member []byte) (score int, ok bool) {
	score, ok = s.dict[string(member)]
	return
}
