package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/satyrius/gonx"
)

// Session represents a monitoring session and handles all accumulated data.
type Session struct {
	AlertThreshold int
	Entries        []*Entry
	File           *os.File
	Parser         *gonx.Parser
	PollInt        int
	Report         *Report
	State          uint8
}

// Report accumulates data for reports.
type Report struct {
	StatusCodes    map[uint8]int // 4 status code groups: 2xx, 3xx, 4xx, 5xx
	Time           *time.Time
	TopSectionHits map[int]Pair
	TotalHits      int
}

// NewSession returns a new Session object.
func NewSession(threshold, pollInt int, p *gonx.Parser) *Session {
	return &Session{
		AlertThreshold: threshold,
		Parser:         p,
		PollInt:        pollInt,
		Report:         NewReport(nil),
		State:          stateOK,
	}
}

// SetLog adds file to the session.
func (s *Session) SetLog(f string) error {
	var err error
	if f == "" {
		return errors.New("No file provided")
	}
	s.File, err = os.Open(f)
	if err != nil {
		return err
	}

	return nil
}

// Close closes log file opened for this session.
func (s *Session) Close() {
	if s.File != nil {
		s.File.Close()
	}
}

// NewReport represents traffic summary for a user-defined report period.
func NewReport(t *time.Time) *Report {
	return &Report{
		StatusCodes: map[uint8]int{
			5: 0,
			4: 0,
			3: 0,
			2: 0,
		},
		Time:           t,
		TopSectionHits: make(map[int]Pair),
	}
}

// ConsumeLines receives log entries accumulated since last poll,
// converts them into Entry objects and adds to the session buffer.
func (s *Session) ConsumeLines(l []string) error {
	for _, line := range l {

		err := s.AddLine(line)
		if err != nil {
			return fmt.Errorf(" Error adding log entry %s \n Err: %s ", line, err.Error())
		}
	}

	return nil
}

// AddLine adds a log entry as an Entry object to the session buffer.
func (s *Session) AddLine(line string) error {

	r := NewEntry(s.Parser)
	err := r.ParseLine(line)
	if err != nil {
		return err
	}

	s.Entries = append(s.Entries, r)
	s.Report.TotalHits++

	return nil
}

// FlushReport returns interval report and resets accumulated stats.
func (s *Session) FlushReport(cfg *Config, t *time.Time) *Report {

	s.GetSectionHits(cfg.TopN) // sets i.Summary.TopSectionHits

	s.GetStatusCodes() // sets i.Summary.StatusCodes

	out := NewReport(t)

	out.TotalHits = s.Report.TotalHits

	for k, v := range s.Report.StatusCodes {
		out.StatusCodes[k] = v
	}

	for k, v := range s.Report.TopSectionHits {
		out.TopSectionHits[k] = v
	}

	s.reset()

	return out
}

// reset nullifies traffic data accumulated since last report.
func (s *Session) reset() {
	s.Report = NewReport(nil)
	s.Entries = []*Entry{}
}

// GetStatusCodes calculates status code summary.
func (s *Session) GetStatusCodes() {

	for _, e := range s.Entries {
		i, err := strconv.Atoi(e.StatusCode[0:1])
		if err != nil {
			// todo: add error-logging
			continue
		}
		codeGroup := uint8(i)
		_, ok := s.Report.StatusCodes[codeGroup]
		if !ok {
			s.Report.StatusCodes[codeGroup] = 1
		} else {
			s.Report.StatusCodes[codeGroup]++
		}
	}
}

// GetSectionHits calculates top n section hits during the interval poll time.
func (s *Session) GetSectionHits(n uint) {

	sectionHits := make(map[string]int, len(s.Entries))

	// Get hits by section.
	for _, e := range s.Entries {

		_, ok := sectionHits[e.Section]
		if !ok {
			sectionHits[e.Section] = 1
		} else {
			sectionHits[e.Section]++
		}
	}

	// Sort.
	sorted := RankByHits(sectionHits)

	// Limit to top n.
	s.Report.TopSectionHits = CutTopN(sorted, n)
}

// ShouldEscalate returns escalation action code.
func (s *Session) ShouldEscalate(traffic int) bool {

	if !s.IsAlert() && traffic >= s.AlertThreshold {
		s.SetAlert()
		return true
	}

	return false
}

// ShouldDeescalate returns escalation action code.
func (s *Session) ShouldDeescalate(traffic int) bool {

	if s.IsAlert() && traffic < s.AlertThreshold {
		s.SetOK()
		return true
	}

	return false
}

// IsAlert checks for alert state.
func (s *Session) IsAlert() bool {
	return s.State == stateAlert
}

// SetAlert sets current sta te to Alert.
func (s *Session) SetAlert() {
	s.State = stateAlert
}

// SetOK sets current sta te to OK.
func (s *Session) SetOK() {
	s.State = stateOK
}
