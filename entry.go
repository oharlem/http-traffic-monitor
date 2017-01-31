package main

import (
	"fmt"
	"strings"

	"github.com/satyrius/gonx"
)

// Entry represents a request based on a log entry data:
// "GET /shuttle/technology/sts-newsref/srb.html HTTP/1.0" 200
type Entry struct {
	Method     string
	Path       string
	parser     *gonx.Parser
	Section    string
	Protocol   string
	StatusCode string // response code to a given request
}

// NewEntry represents log entry.
func NewEntry(p *gonx.Parser) *Entry {
	return &Entry{
		parser: p,
	}
}

// ParseLine validate and parses a request part of a log entry and creates a Request object on success.
func (r *Entry) ParseLine(line string) error {

	e, err := r.parser.ParseString(line)
	if err != nil {
		return fmt.Errorf(" Error parsing log entry %s \n Err: %s ", line, err.Error())
	}

	rStr, err := e.Field("request")
	if err != nil {
		return err
	}

	parts := strings.Fields(rStr)
	if len(parts) != 3 {
		return fmt.Errorf("Invalid structure of request part: %s", rStr)
	}

	r.Method = parts[0]
	r.Path = parts[1]
	r.Protocol = parts[2]

	partParts := strings.Split(r.Path, "/")

	r.Section = "/"
	if len(partParts) > 0 {
		r.Section += partParts[1]
	}

	r.StatusCode, err = e.Field("status")
	if err != nil {
		return err
	}

	return nil
}
