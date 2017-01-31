package main

import (
	"bufio"
	"time"
)

// Monitor is a monitoring routine that tracks changes to a log file,
// calculates metrics based on accumulated data
// and issues messages based on changes to a log file and/or metrics.
func Monitor(cfg *Config, s *Session, doneChan chan struct{}, msgChan chan<- msg) {

	f := NewFrame(cfg.MTF, cfg.PollInt)
	r := bufio.NewReader(s.File)

	stat, err := s.File.Stat()
	if err != nil {
		panic(err)
	}

	prevSize := stat.Size() // starting file read position

	// Move reader's needle to the position where we stopped reading last time or to the initial position.
	if _, err := s.File.Seek(prevSize, 0); err != nil {
		panic(err)
	}

	// Start tickers
	tickerPolling := time.NewTicker(time.Second * time.Duration(cfg.PollInt))
	tickerReporting := time.NewTicker(time.Second * time.Duration(cfg.ReportInt))

	polls := 0

monitorLoop:
	for {
		select {

		// Main completion handler.
		case <-doneChan:
			tickerPolling.Stop()
			tickerReporting.Stop()
			break monitorLoop

		// Poll ticker.
		case t := <-tickerPolling.C:
			{
				polls++

				// Capture point data.
				p := NewPoint(s.File, r, prevSize)
				err := p.GetChange()
				if err != nil {
					msgChan <- msgErr(err)
				}

				// Register current level of traffic, i.e.
				// quantity of log entries since last poll.
				f.Rec(p.linesQty)

				// Pass entries to the session storage.
				err = s.ConsumeLines(p.lines)
				if err != nil {
					msgChan <- msgErr(err)
				}

				// Print out current point data.
				if cfg.SendTicks {
					msgChan <- msgPoint(f.AvgTraffic, s.AlertThreshold)
				}

				// Monitor alert threshold.
				if cfg.SendAlerts {
					if s.ShouldEscalate(f.AvgTraffic) {
						msgChan <- msgAlertEsc(f.AvgTraffic, t)
					}

					if s.ShouldDeescalate(f.AvgTraffic) {
						msgChan <- msgAlertDeesc(f.AvgTraffic, t)
					}
				}

				if polls == cfg.MaxPolls {
					doneChan <- struct{}{}
				}
			}

		// Reporting ticker.
		case t := <-tickerReporting.C:
			{
				if cfg.SendReports {
					// Get data accumulated during report interval and clean report buffer.
					msgChan <- msgReport(s.FlushReport(cfg, &t))
				}
			}
		}

	}
}
