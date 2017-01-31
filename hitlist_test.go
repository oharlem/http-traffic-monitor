package main

import (
	"reflect"
	"testing"
)

func TestRankSectionsByHits(t *testing.T) {
	testHitsMap := map[string]int{
		"/foo":    2,
		"/foobar": 5,
		"/barbaz": 4,
		"/baz":    3,
		"/bar":    1,
	}

	expected := HitList{
		0: Pair{"/foobar", 5},
		1: Pair{"/barbaz", 4},
		2: Pair{"/baz", 3},
		3: Pair{"/foo", 2},
		4: Pair{"/bar", 1},
	}

	actual := RankByHits(testHitsMap) // *HitList

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Failed RankSectionsByHits test!")
		t.Log("Expected:")
		t.Logf("%+v\n", expected)

		t.Log("Actual:")
		t.Logf("%+v\n", actual)
	}
}

func TestCutTopN(t *testing.T) {
	testHitList := HitList{
		0: Pair{"/foobar", 5},
		1: Pair{"/barbaz", 4},
		2: Pair{"/baz", 3},
		3: Pair{"/foo", 2},
		4: Pair{"/bar", 1},
	}

	expected := HitList{
		0: Pair{"/foobar", 5},
		1: Pair{"/barbaz", 4},
		2: Pair{"/baz", 3},
	}

	actual := CutTopN(testHitList, 3)

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Failed CutTopN test!")
		t.Log("Expected:")
		t.Logf("%+v\n", expected)

		t.Log("Actual:")
		t.Logf("%+v\n", actual)
	}
}
