package main

import (
	"sort"
)

// Pair is a k/v bucket for HitList.
type Pair struct {
	Key   string
	Value int
}

// HitList is a map for sorted data with struct value types.
type HitList map[int]Pair

// CutTopN returns top n hit sections from an interval data.
func CutTopN(h HitList, n uint) HitList {

	if len(h) <= int(n) {
		return h
	}

	out := HitList{}
	for i := 0; i < int(n); i++ {
		out[i] = h[i]
	}

	return out
}

// RankByHits sorts HitList in reverse order based on Pair Value.
func RankByHits(sectionHits map[string]int) HitList {
	pl := make(HitList, len(sectionHits))
	if len(sectionHits) == 0 {
		return pl
	}
	i := 0
	for k, v := range sectionHits {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func (p HitList) Len() int           { return len(p) }
func (p HitList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p HitList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
