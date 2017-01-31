package main

import (
	"testing"

	"github.com/satyrius/gonx"
)

func TestRequest_ParseEntry(t *testing.T) {

	parser := gonx.NewParser(parserFormat)

	testString := `182.198.120.1 - - [28/Jul/1995:13:16:47 -0400] "GET /shuttle/technology/sts-newsref/srb.html HTTP/1.0" 200 49553`

	r := NewEntry(parser)
	err := r.ParseLine(testString)
	if err != nil {
		t.Fatalf("ParseEntry should not fail. Error: %+v", err)
	}

	expected := "GET"
	actual := r.Method
	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

	expected = "/shuttle/technology/sts-newsref/srb.html"
	actual = r.Path
	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

	expected = "HTTP/1.0"
	actual = r.Protocol
	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

	expected = "/shuttle"
	actual = r.Section
	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

	expected2 := "200"
	actual2 := r.StatusCode
	if expected2 != actual2 {
		t.Errorf("Expected %d, got %d", expected2, actual2)
	}

}
