package main

import (
	"flag"
	"math"
)

const (
	// Argument defaults.
	defAlertThreshold = 1000 // hits per interval
	defMTF            = 120  // Monitoring time frame, sec
	defPollInt        = 1    // sec, 1sec - minimum
	defReportInt      = 10   // Default stat summary interval, sec
	defTopN           = 10
	defSendAlerts     = true
	defSendReports    = true
	defSendTicks      = true
)

// Config is a program configuration object.
type Config struct {
	AlertThreshold int
	File           string
	MaxPolls       int
	MTF            int // sec
	PollInt        int // sec
	ReportInt      int // sec
	SendAlerts     bool
	SendReports    bool
	SendTicks      bool
	TopN           uint
}

// NewConfig initializes program configuration and runs basic validation of user-defined arguments.
// By default, all properties are set to def* constants.
func NewConfig() *Config {
	at := flag.Int("alert-threshold", defAlertThreshold, "Alert threshold")
	lf := flag.String("log-file", "", "Log file.")
	mtf := flag.Int("mtf", defMTF, "Monitoring time frame (seconds)")
	pi := flag.Int("poll-interval", defPollInt, "Log polling interval (seconds). 1 sec - mim allowed value.")
	ri := flag.Int("report-interval", defReportInt, "Report interval.")
	sa := flag.Bool("send-alerts", defSendAlerts, "Send alerts")
	sr := flag.Bool("send-reports", defSendReports, "Send reports")
	st := flag.Bool("send-ticks", defSendTicks, "Send tick information")
	tn := flag.Uint("top-n", defTopN, "Number of top section hits displayed during polls")
	flag.Parse()

	if *lf == "" {
		panic("Log file is not provided.")
	}

	if *pi < 1 {
		panic("Invalid polling interval set. Minimal allowed value is 1 second.")
	}

	if *ri < *pi {
		// via pollInt we also check for reportInt to be > 0
		panic("Report interval cannot be smaller than polling interval.")
	}

	if *mtf < *pi {
		panic("Monitoring time frame cannot be smaller than polling interval.")
	}

	return &Config{
		MaxPolls:       math.MaxInt32 - 1,
		AlertThreshold: *at,
		File:           *lf,
		MTF:            *mtf,
		PollInt:        *pi,
		ReportInt:      *ri,
		SendAlerts:     *sa,
		SendReports:    *sr,
		SendTicks:      *st,
		TopN:           *tn,
	}
}
