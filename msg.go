package main

import "time"

type msg struct {
	msgType   string
	body      string
	report    *Report
	time      time.Time
	traffic   int
	threshold int
}

const (
	// Message types.

	msgTypeAlertEsc   = "alertEsc"
	msgTypeAlertDeesc = "alertDeesc"
	msgTypeError      = "err"
	msgTypePoint      = "point"
	msgTypeReport     = "report"
)

// Message objects.

func msgAlertEsc(tr int, t time.Time) msg {
	return msg{
		msgType: msgTypeAlertEsc,
		traffic: tr,
		time:    t,
	}
}

func msgAlertDeesc(tr int, t time.Time) msg {
	return msg{
		msgType: msgTypeAlertDeesc,
		traffic: tr,
		time:    t,
	}
}

func msgErr(err error) msg {
	return msg{
		msgType: msgTypeError,
		body:    err.Error(),
	}
}

func msgPoint(tr, th int) msg {
	return msg{
		msgType:   msgTypePoint,
		threshold: th,
		traffic:   tr,
	}
}

func msgReport(r *Report) msg {
	return msg{
		msgType: msgTypeReport,
		report:  r,
	}
}
