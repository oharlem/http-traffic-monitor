package main

import (
	"reflect"
	"testing"

	"github.com/satyrius/gonx"
)

func TestNewInterval(t *testing.T) {

	actual := NewSession(2, 2, nil)
	expected := &Session{
		AlertThreshold: 2,
		PollInt:        2,
		Report: &Report{
			StatusCodes: map[uint8]int{
				5: 0,
				4: 0,
				3: 0,
				2: 0,
			},
			TopSectionHits: make(map[int]Pair),
		},
		State: stateOK,
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Failed NewInterval test!")
		t.Log("Expected:")
		t.Logf("%+v\n", expected)

		t.Log("Actual:")
		t.Logf("%+v\n", actual)
	}
}

func TestInterval_AddEntry(t *testing.T) {

	parser := gonx.NewParser(parserFormat)
	s := NewSession(2, 2, parser)

	testString := `182.198.120.1 - - [28/Jul/1995:13:16:47 -0400] "GET /shuttle/technology/sts-newsref/srb.html HTTP/1.0" 200 49553`

	err := s.AddLine(testString)
	if err != nil {
		t.Fatalf("AddEntry should not fail. Error: %+v", err)
	}
	err = s.AddLine(testString)
	if err != nil {
		t.Fatalf("AddEntry should not fail. Error: %+v", err)
	}

	expected := 2
	actual := len(s.Entries)

	if expected != actual {
		t.Errorf("Should have %d entries, got %v", expected, actual)
	}
}

func TestInterval_GetSectionHits(t *testing.T) {

	parser := gonx.NewParser(parserFormat)
	s := NewSession(2, 2, parser)
	var err error

	// Section: /shuttle
	// Hist: 2
	for n := 0; n < 2; n++ {
		err = s.AddLine(`182.198.120.1 - - [28/Jul/1995:13:16:47 -0400] "GET /shuttle/technology/sts-newsref/srb.html HTTP/1.0" 200 49553`)
		if err != nil {
			t.Fatalf("AddEntry should not fail. Error: %+v", err)
		}
	}

	// Section: /images
	// Hist: 3
	for n := 0; n < 3; n++ {
		err = s.AddLine(`198.155.12.13 - - [28/Jul/1995:13:17:09 -0400] "GET /images/NASA-logosmall.gif HTTP/1.0" 200 786`)
		if err != nil {
			t.Fatalf("AddEntry should not fail. Error: %+v", err)
		}
	}

	// Section: /history
	// Hist: 1
	err = s.AddLine(`165.13.14.55 - - [28/Jul/1995:13:17:00 -0400] "GET /history/apollo/apollo-17/apollo-17-info.html HTTP/1.0" 200 1457`)
	if err != nil {
		t.Fatalf("AddEntry should not fail. Error: %+v", err)
	}

	expectedTopN := uint(3)
	s.GetSectionHits(expectedTopN)

	expected := map[int]Pair{
		0: {"/images", 3},
		1: {"/shuttle", 2},
		2: {"/history", 1},
	}

	actual := s.Report.TopSectionHits
	actualTopN := uint(len(s.Report.TopSectionHits))

	if expectedTopN != actualTopN {
		t.Errorf("Expected Top N %d, got %d.", expectedTopN, actualTopN)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Failed GetTopHits test!")
		t.Log("Expected:")
		t.Logf("%+v\n", expected)

		t.Log("Actual:")
		t.Logf("%+v\n", actual)
	}
}

func TestInterval_GetStatusCodes(t *testing.T) {

	parser := gonx.NewParser(parserFormat)
	s := NewSession(2, 2, parser)
	var err error

	// Code: 200
	// Hits: 2
	for n := 0; n < 2; n++ {
		err = s.AddLine(`182.198.120.1 - - [28/Jul/1995:13:16:47 -0400] "GET /shuttle/technology/sts-newsref/srb.html HTTP/1.0" 200 49553`)
		if err != nil {
			t.Fatalf("AddEntry should not fail. Error: %+v", err)
		}
	}

	// Code: 500
	// Hits: 3
	for n := 0; n < 3; n++ {
		err = s.AddLine(`198.155.12.13 - - [28/Jul/1995:13:17:09 -0400] "GET /images/NASA-logosmall.gif HTTP/1.0" 500 786`)
		if err != nil {
			t.Fatalf("AddEntry should not fail. Error: %+v", err)
		}
	}

	// Code: 401
	// Hits: 1
	err = s.AddLine(`165.13.14.55 - - [28/Jul/1995:13:17:00 -0400] "GET /history/apollo/apollo-17/apollo-17-info.html HTTP/1.0" 401 1457`)
	if err != nil {
		t.Fatalf("AddEntry should not fail. Error: %+v", err)
	}

	expected := map[uint8]int{
		5: 3,
		4: 1,
		3: 0,
		2: 2,
	}

	s.GetStatusCodes()
	actual := s.Report.StatusCodes

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Failed GetStatusCodes test!")
		t.Log("Expected:")
		t.Logf("%+v\n", expected)

		t.Log("Actual:")
		t.Logf("%+v\n", actual)
	}
}

func TestSession_ShouldEscalate_True(t *testing.T) {

	threshold := 1
	s := NewSession(threshold, 2, nil)
	s.SetOK()

	expected := true

	traffic := 2

	actual := s.ShouldEscalate(traffic) // traffic = 2, threshold = 1

	if actual != expected {
		t.Errorf("UpdateState(%d, %d): expected %d, actual %d", traffic, threshold, expected, actual)
	}
}

func TestSession_ShouldEscalate_False(t *testing.T) {

	threshold := 3
	s := NewSession(threshold, 2, nil)
	s.SetOK()

	expected := false

	traffic := 2

	actual := s.ShouldEscalate(traffic) // traffic = 2, threshold = 1

	if actual != expected {
		t.Errorf("UpdateState(%d, %d): expected %t, actual %t", traffic, threshold, expected, actual)
	}
}
