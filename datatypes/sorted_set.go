package datatypes

import (
	"math/rand"
	"strings"
)

const TSortedSet string = "sset"

type node struct {
	value   string
	score   float32
	forward []*node
}

type WithScore struct {
	value string
	score float32
}

func newNode(value string, score float32, level int) *node {
	return &node{
		value,
		score,
		make([]*node, level+1),
	}
}

type skipList struct {
	maxLevel         int
	levelProbability float32
	level            int
	head             *node
}

func NewskipList(maxLevel int, levelProbability float32) skipList {
	return skipList{
		maxLevel,
		levelProbability,
		0,
		nil,
	}
}

func (s skipList) randomLevel() int {
	lvl := 0
	r := rand.Float32()

	for r < s.levelProbability && lvl < s.maxLevel {
		lvl++
		r = rand.Float32()
	}

	return lvl
}

func (s *skipList) insert(value string, score float32) *node {
	if s.head == nil {
		s.head = newNode(value, score, s.maxLevel) // this is temp
		return s.head
	}

	// prepending to list
	if score < s.head.score || (score == s.head.score && strings.Compare(value, s.head.value) == -1) {
		node := newNode(value, score, 0)
		s.head, node = node, s.head
		s.head.forward[0] = node
		return s.head
	}

	current := s.head
	update := make([]*node, s.maxLevel+1)

	for i := s.level; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].score < score {
			current = current.forward[i]
		}

		for current.forward[i] != nil &&
			current.forward[i].score == score &&
			strings.Compare(current.forward[i].value, value) == -1 {
			current = current.forward[i]
		}

		update[i] = current
	}

	current = current.forward[0]

	if current == nil || current.value != value {
		rlevel := s.randomLevel()

		if rlevel > s.level {
			for i := s.level + 1; i < rlevel+1; i++ {
				update[i] = s.head
			}

			s.level = rlevel
		}

		n := newNode(value, score, rlevel)

		for i := 0; i <= rlevel; i++ {
			n.forward[i] = update[i].forward[i]
			update[i].forward[i] = n
		}

		return n
	}

	return nil
}

func (s skipList) RangeWithScores(start, end int) []WithScore {
	var result []WithScore
	current := s.head

	counter := 0

	if start < 0 {
		start = 0
	}

	// end does not need to be reset if it's larger than the list, since iteration will stop anyway at the end of the list
	for current != nil {
		if counter > end {
			return result
		}

		if counter >= start {
			result = append(result, WithScore{current.value, current.score})
		}

		current = current.forward[0]
		counter++
	}

	return result
}

type KvSortedSet struct {
	hMap  map[string]*node
	sList skipList
}

func (ss KvSortedSet) Insert(value string, score float32) {
	_, exists := ss.hMap[value]

	if !exists {
		// TODO delete old node
	}

	n := ss.sList.insert(value, score)
	ss.hMap[value] = n
}

// TODO get this done
func (ss KvSortedSet) Delete(value string) {}
func (ss KvSortedSet) Get(value string) {}
func (ss KvSortedSet) Range(start, end int) {}
func (ss KvSortedSet) RangeWithScores(start, end int) {}

func (ss KvSortedSet) Type() string {
	return TSortedSet
}

// TODO get this working if needed
func (ss KvSortedSet) Response() []byte {
  return []byte("")
}
