package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Point represents data accumulated during last log check.
type Point struct {
	prevSize int64
	size     int64
	diff     int64
	lines    []string
	linesQty int
	reader   *bufio.Reader
	file     *os.File
}

// NewPoint returns a new Point object.
func NewPoint(f *os.File, r *bufio.Reader, prevSize int64) *Point {
	return &Point{
		file:     f,
		prevSize: prevSize,
		reader:   r,
	}
}

// GetChange checks log file for size changes and reads added data into log entry strings.
func (p *Point) GetChange() error {

	stat, err := p.file.Stat()
	if err != nil {
		return err
	}

	p.size = stat.Size()
	p.diff = p.size - p.prevSize

	// If truncated, adjust for a new size and continue from beginning.
	if p.diff > 0 {
		p.lines, err = readIncrement(p.reader)
		if err != nil {
			return fmt.Errorf(" Error reading log chunk: %s ", err.Error())
		}

		p.linesQty = len(p.lines)
	}

	p.prevSize = p.size

	return nil
}

// readIncrement reads a log file from a start position to EOF
// returning result as a slice of (string) log entries.
func readIncrement(r *bufio.Reader) ([]string, error) {

	var data []byte
	var err error
	var out []string

	for {
		data, err = r.ReadBytes('\n')

		if err == nil || err == io.EOF {
			line := strings.TrimSpace(string(data))
			if line != "" {
				out = append(out, line)
			}
		}

		if err != nil {
			if err != io.EOF {
				return out, err
			}

			// EOF
			break
		}
	}

	return out, nil
}
