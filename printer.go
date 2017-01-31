package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	ct "github.com/daviddengcn/go-colortext"
)

var (
	reportTimeFormat = time.RFC3339
)

// Printer is a handler for stdout outputs.
func Printer(cfg *Config, doneChan <-chan struct{}, msgChan <-chan msg) {

printerLoop:
	for {
		select {
		case <-doneChan:
			break printerLoop

		case m := <-msgChan:
			switch m.msgType {
			case msgTypeError:
				printErr(m.body)
			case msgTypeAlertEsc:
				printAlertEsc(m.traffic, m.time)
			case msgTypeAlertDeesc:
				printAlertDeesc(m.traffic, m.time)
			case msgTypePoint:
				printPoint(m.traffic, m.threshold)
			case msgTypeReport:
				printReport(cfg, m.report)
			}
		}
	}

}

// printAlertEsc prints alert escalation message.
func printAlertEsc(tr int, t time.Time) {
	printBigMsg("\u00B7 High traffic generated an alert - hits = %d, triggered at %s", tr, t, ct.Red)

}

// printAlertDeesc prints alert de-escalation message.
func printAlertDeesc(tr int, t time.Time) {
	printBigMsg("\u00B7 High traffic alert recovered. Current hits = %d. At %s", tr, t, ct.Green)
}

// printReport prints out a report block.
func printReport(cfg *Config, r *Report) {
	if len(r.TopSectionHits) == 0 && len(r.StatusCodes) == 0 && r.TotalHits == 0 {
		fmt.Print("REPORT: no changes\n\n")
	} else {
		fmt.Print("\n\n")
		printHR()

		fmt.Printf("REPORT: %s\n\n", r.Time.Format(reportTimeFormat))

		printSections(cfg.TopN, r.TopSectionHits)

		printSummary(cfg, r)

		fmt.Print("\n\n")
	}
}

// printPoint prints out a point data.
func printPoint(traffic, threshold int) {
	fmt.Print(" \u00B7  hits avg: ")
	if traffic >= threshold {
		printRed("%6d", traffic)
	} else {
		fmt.Printf("%6d", traffic)
	}
	fmt.Printf("  / %2d\n", threshold)
}

func printErr(s string, a ...interface{}) {
	fmt.Print("\n")
	labelYellow(s, a...)
	fmt.Print("\n")
}

// printHR prints a horizontal ruler to stdout.
func printHR() {
	ct.ChangeColor(ct.White, false, ct.Black, false)
	fmt.Print(strings.Repeat("-", 80))
	ct.ResetColor()
	fmt.Print("\n")
}

// Labels (fg & bg).

// labelRed - white fg, red bg.
func labelRed(s string, a ...interface{}) {
	ct.ChangeColor(ct.White, true, ct.Red, true)
	fmt.Printf(s, a...)
	ct.ResetColor()
}

// labelYellow - black fg, yellow bg.
func labelYellow(s string, a ...interface{}) {
	ct.ChangeColor(ct.Black, false, ct.Yellow, false)
	fmt.Printf(s, a...)
	ct.ResetColor()
}

// Foreground text coloring only.

// printRed - red fg.
func printRed(s string, a ...interface{}) {
	ct.ChangeColor(ct.Red, true, ct.Black, false)
	fmt.Printf(s, a...)
	ct.ResetColor()
}

// printSections prints out a top sections block of a report.
func printSections(n uint, s map[int]Pair) {
	if len(s) == 0 {
		fmt.Printf("Top %d sections: no entries\n", n)
	} else {
		fmt.Printf("Top %d sections\n", n)
		fmt.Print("| sections                                                       | count")
		fmt.Print("\n")
		printHR()
		for i := 0; i < len(s); i++ {
			fmt.Print("| " + rightPad2Len(s[i].Key, " ", 63))
			fmt.Print("| " + rightPad2Len(strconv.Itoa(s[i].Value), " ", 10))
			fmt.Print("\n")
		}
	}
	fmt.Print("\n")
}

// printSummary prints out a summary part of a report.
func printSummary(cfg *Config, r *Report) {

	fmt.Print("Summary:\n")
	fmt.Print("| hits total | hits/s     | 2xx        | 3xx        | 4xx        | 5xx        ")
	fmt.Print("\n")
	printHR()
	fmt.Print("| " + rightPad2Len(strconv.Itoa(r.TotalHits), " ", 11))

	// total hits / seconds for this interval
	hits := int(math.Ceil(float64(r.TotalHits / cfg.ReportInt)))
	fmt.Print("| " + rightPad2Len(strconv.Itoa(hits), " ", 11))

	fmt.Print("| " + rightPad2Len(strconv.Itoa(r.StatusCodes[2]), " ", 11))
	fmt.Print("| " + rightPad2Len(strconv.Itoa(r.StatusCodes[3]), " ", 11))

	if r.StatusCodes[4] > 0 {
		fmt.Print("| ")
		labelYellow("%s", rightPad2Len(strconv.Itoa(r.StatusCodes[4]), " ", 11))
	} else {
		fmt.Print("| 0          ")
	}

	if r.StatusCodes[5] > 0 {
		fmt.Print("| ")
		labelRed("%s", rightPad2Len(strconv.Itoa(r.StatusCodes[5]), " ", 11))
	} else {
		fmt.Print("| 0          ")
	}

	fmt.Print("\n")
}

func rightPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

func printBigMsg(s string, tr int, t time.Time, bg ct.Color) {
	fmt.Print("\n")

	ct.ChangeColor(ct.White, true, bg, true)
	fmt.Printf("%87s", " ")
	ct.ResetColor()

	fmt.Print("\n")

	str := fmt.Sprintf(s, tr, t.Format(time.RFC3339))
	size := len(str)
	remaining := strconv.Itoa(86 - size)

	ct.ChangeColor(ct.White, true, bg, true)
	fmt.Printf(` %s %`+remaining+`s`, str, " ")
	ct.ResetColor()

	fmt.Print("\n")

	ct.ChangeColor(ct.White, true, bg, true)
	fmt.Printf("%87s", " ")
	ct.ResetColor()

	fmt.Print("\n\n")
}
