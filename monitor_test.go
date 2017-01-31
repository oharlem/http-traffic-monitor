package main

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/satyrius/gonx"
)

// Test logic for alert escalation case:
//
// 1. Initial state is OK (default).
// 2. Run for the MTF duration. In this test: 2 seconds.
// 3. Imitate incoming traffic adding several lines to the temporary log file once per second.
// 4. To trigger an alert, we set alert threshold level to 2 hits/s and add 2 entries per second.
//
// 5. Expected result:
// - by the end of MTF avg. traffic will be 2 hits/sec,
// - monitor will register that traffic reached a threshold level,
// - monitor will set an alert state,
// - alert escalation message will be sent to msgChan.
func TestMonitor_AlertEscalation(t *testing.T) {
	var err error
	tempLogFile := getTempLoc(".TestMonitor_AlertEscalation.log")
	alertThreshold := 2
	pollInt := 1

	cfg := &Config{
		AlertThreshold: alertThreshold,
		File:           tempLogFile,
		MTF:            2,
		MaxPolls:       2, // As poll interval 1 sec, thus we limit test to 3 sec length.
		TopN:           3,
		PollInt:        pollInt, // Poll once per second.
		ReportInt:      2,       // Irrelevant, as reports are off for this test.
		SendAlerts:     true,
		SendReports:    false,
		SendTicks:      false,
	}

	s := NewSession(alertThreshold, pollInt, gonx.NewParser(parserFormat))

	// Ensure test log file is in place.
	f, err := os.Create(tempLogFile)
	if err != nil {
		t.Fatalf("\n\nCannot open test file: %s\n", err.Error())
	}
	defer f.Close()

	err = s.SetLog(cfg.File)
	if err != nil {
		// Suppress error as a test log file can be absent
	}
	// Also, not deferring a Close method, as the file will be created later.

	doneChan := make(chan struct{})
	msgChan := make(chan msg)

	go Monitor(cfg, s, doneChan, msgChan)

	// Handler for monitor closing.
	go Ctrl(doneChan)

	expected := 1
	actual := 0

	logUpdateTicker := time.NewTicker(time.Second * 1)

	// str mimics one-time entry of 2 log lines
	// In this case, with 1 poll per second, this equals to a traffic of 2 hits/s, while threshold is 2 hits/s.
	// To trigger an alert, number of lines should be equal or higher than the "AlertThreshold" value.
	str := `
	210.166.12.00 - - [28/Jul/1995:13:17:36 -0400] "GET /htbin/cdt_clock.pl HTTP/1.0" 200 503
	198.155.12.16 - - [28/Jul/1995:13:17:09 -0400] "GET /images/NASA-logosmall.gif HTTP/1.0" 400 786
	`

	// Intercept message channel.
	// Instead of a printer we receive messages here to assert results based on message type.
selectLoop:
	for {
		select {
		case <-logUpdateTicker.C:
			{
				if _, err = f.WriteString(str); err != nil {
					t.Fatalf("\n\nCannot writing to test log: %s\n\n", err.Error())
				}
			}

		case msg := <-msgChan:
			if msg.msgType == msgTypeAlertEsc {
				actual++
			}

		case <-doneChan:
			break selectLoop
		}
	}

	if actual != 1 {
		t.Errorf("Expected %d message, got %d", expected, actual)
	}

	if !s.IsAlert() {
		t.Errorf("Expected %d state, got %d", stateAlert, s.State)
	}

	// Cleanup
	err = os.Remove(tempLogFile)
	if err != nil {
		t.Errorf("Could not delete temp file: %s", err.Error())
	}
}

// Test logic for alert deescalation:
//
// 1. Set initial state to Alert.
// 2. Run for the MTF duration. In this test: 2 seconds.
// 3. Do no additions to the log file imitating no traffic to the system.
//
// 4. Expected result:
// - by the end of MTF avg. traffic will be 0 hits/sec,
// - monitor will register that traffic dropped below a threshold level,
// - monitor will set an OK state,
// - alert deescalation message will be sent to msgChan.
func TestMonitor_AlertDeescalation(t *testing.T) {
	var err error
	tempLogFile := getTempLoc(".TestMonitor_AlertDeescalation.log")
	alertThreshold := 2
	pollInt := 1

	cfg := &Config{
		AlertThreshold: alertThreshold,
		File:           tempLogFile,
		MTF:            2,
		MaxPolls:       2, // As poll interval 1 sec, thus we limit test to 3 sec length.
		TopN:           3,
		PollInt:        pollInt, // Poll once per second.
		ReportInt:      2,       // Irrelevant, as reports are off for this test.
		SendAlerts:     true,
		SendReports:    false,
		SendTicks:      false,
	}

	s := NewSession(alertThreshold, pollInt, gonx.NewParser(parserFormat))
	s.SetAlert()

	// Ensure test log file is in place.
	f, err := os.Create(tempLogFile)
	if err != nil {
		t.Fatalf("\n\nCannot open test file: %s\n", err.Error())
	}
	defer f.Close()

	err = s.SetLog(cfg.File)
	if err != nil {
		// Suppress error as a test log file can be absent
	}
	// Also, not deferring a Close method, as the file will be created later.

	doneChan := make(chan struct{})
	msgChan := make(chan msg)

	go Monitor(cfg, s, doneChan, msgChan)

	// Handler for monitor closing.
	go Ctrl(doneChan)

	expected := 1
	actual := 0

	// Intercept message channel.
selectLoop:
	for {
		select {

		case msg := <-msgChan:
			if msg.msgType == msgTypeAlertDeesc {
				actual++
			}

		case <-doneChan:
			break selectLoop
		}
	}

	if actual != 1 {
		t.Errorf("Expected %d message, got %d", expected, actual)
	}

	if s.IsAlert() {
		t.Errorf("Expected %d state, got %d", stateOK, s.State)
	}

	// Cleanup
	err = os.Remove(tempLogFile)
	if err != nil {
		t.Errorf("Could not delete temp file: %s", err.Error())
	}
}

// getTempLoc prepares full temporary file location.
// One of the purposes - workaround between differences of MacOS temp folder ending with a slash,
// and Debian TMPDIR env var being empty and thus os.TempDir() was creating a temp dir
// without an ending slash.
func getTempLoc(filename string) string {
	return strings.TrimRight(os.TempDir(), "/") + "/" + filename
}
