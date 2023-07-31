package go_sorted_set

import (
	"crypto/md5"
	"hash"
)

type SortedSet struct {
	skipList  *SkipList
	hashTable map[string]int
	hash      hash.Hash
}

func NewSortedSet() *SortedSet {
	return &SortedSet{
		skipList:  newSkipList(),
		hashTable: make(map[string]int),
		hash:      md5.New(),
	}
}

func (s *SortedSet) ZAdd(key string, score int, member []byte) {
	s.hash.Write(member)
	defer s.hash.Reset()
	s.skipList.Insert(key, score, member)
	s.hashTable[string(s.hash.Sum(nil))] = score
}

func (s *SortedSet) ZRange(start, end int) []*Node {
	return s.skipList.Search(start, end)
}

func (s *SortedSet) ZRem(_ string, member []byte) {
	s.hash.Write(member)
	defer s.hash.Reset()
	if score, ok := s.hashTable[string(s.hash.Sum(nil))]; ok {
		s.skipList.Delete(score, member)
	}
}

func (s *SortedSet) ZRemRangeByScore(start, end int) []*Node {
	nodes := s.skipList.DeleteByScoreRange(start, end)
	for _, node := range nodes {
		s.hash.Write(node.member)
		delete(s.hashTable, string(s.hash.Sum(nil)))
		s.hash.Reset()
	}
	return nodes
}

func (s *SortedSet) ZRemByScore(score int) []*Node {
	nodes := s.skipList.DeleteByScore(score)
	for _, node := range nodes {
		s.hash.Write(node.member)
		delete(s.hashTable, string(s.hash.Sum(nil)))
		s.hash.Reset()
	}
	return nodes
}

func (s *SortedSet) ZIncrBy(key string, score int, member []byte) {
	s.hash.Write(member)
	defer s.hash.Reset()
	if oldScore, ok := s.hashTable[string(s.hash.Sum(nil))]; ok {
		s.skipList.Delete(oldScore, member)
		s.skipList.Insert(key, oldScore+score, member)
		s.hashTable[string(s.hash.Sum(nil))] = oldScore + score
	} else {
		s.skipList.Insert(key, score, member)
		s.hashTable[string(s.hash.Sum(nil))] = score
	}
}
