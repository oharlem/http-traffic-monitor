package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/satyrius/gonx"
)

const (
	// https://en.wikipedia.org/wiki/Common_Log_Format
	// 127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326
	parserFormat = "$remote_addr $user_identifier $remote_user [$time_local] \"$request\" $status $bytes_sent"
)

var (
	cfg *Config
)

func main() {
	cfg = NewConfig()

	s := NewSession(cfg.AlertThreshold, cfg.PollInt, gonx.NewParser(parserFormat))
	err := s.SetLog(cfg.File)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	doneChan := make(chan struct{})
	msgChan := make(chan msg)

	go Monitor(cfg, s, doneChan, msgChan)

	go Ctrl(doneChan)

	go Printer(cfg, doneChan, msgChan)

	<-doneChan

	fmt.Print("\nMonitor stopped.\n")
}

// Ctrl handles monitor shutdown actions.
func Ctrl(doneChan chan<- struct{}) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	for range sigChan {
		doneChan <- struct{}{}
	}
}
