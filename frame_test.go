package main

import (
	"reflect"
	"testing"
)

func TestNewFrame(t *testing.T) {

	expected := &Frame{
		PointsQty: 2,
	}

	actual := NewFrame(5, 2)

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Failed GetTopHits test!")
		t.Log("Expected:")
		t.Logf("%+v\n", expected)

		t.Log("Actual:")
		t.Logf("%+v\n", actual)
	}
}

func TestFrame_Rec(t *testing.T) {
	f := NewFrame(6, 2) // frame of 3 items
	f.Rec(6)            // should be removed
	f.Rec(5)
	f.Rec(2)
	f.Rec(1)

	expected := []int{5, 2, 1}

	actual := f.PointHits

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Failed GetTopHits test!")
		t.Log("Expected:")
		t.Logf("%+v\n", expected)

		t.Log("Actual:")
		t.Logf("%+v\n", actual)
		t.Logf("PointQty: %d", f.PointsQty)
	}
}

func TestFrame_recalcAvgTraffic(t *testing.T) {
	f := NewFrame(6, 2) // frame of 3 items
	f.Rec(6)            // should be removed
	f.Rec(5)
	f.Rec(2)
	f.Rec(1)

	expected := 2 // floor( ( 5+2+1 ) / 3 )
	actual := f.AvgTraffic

	if expected != actual {
		t.Errorf("TestFrame_recalcAvgTraffic / (%d, %d): expected %d, actual %d", 6, 2, expected, actual)
	}
}
