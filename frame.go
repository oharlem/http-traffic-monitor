package main

import (
	"math"
)

// Frame provides data collected and processed during one poll interval.
type Frame struct {
	PointHits  []int // timeline of hits per each point during MTF
	AvgTraffic int
	PointsQty  int
}

// NewFrame returns a new Frame with a pre-calculated quantity of points per frame.
func NewFrame(mtf, pollInt int) *Frame {

	// We monitor last "mtf" seconds of traffic
	// and checking this time window every pollInt seconds, i.e.
	// we need to store floor(mtf/pollInt) points.
	return &Frame{
		PointsQty: calcPointsPerFrame(mtf, pollInt),
	}
}

// Rec adds quantity of hits for one point.
func (f *Frame) Rec(qty int) {
	f.PointHits = append(f.PointHits, qty)
	if len(f.PointHits) > f.PointsQty {
		// remove 1 element from its head to keep the frame length consistent
		f.PointHits = f.PointHits[1:]
	}

	f.recalcAvgTraffic()
}

// recalcAvgTraffic calculates average traffic volume based on accumulated traffic levels
// for each poll during the user-defined attention span.
func (f *Frame) recalcAvgTraffic() {
	sum := 0
	for _, p := range f.PointHits {
		sum += p
	}

	f.AvgTraffic = calcAvgTraffic(sum, f.PointsQty)
}

// calcAvgTraffic calculates average traffic level per point. Rounded down to nearest int.
func calcAvgTraffic(total, points int) int {
	return int(math.Ceil(float64(total / points)))
}

// calcPointsPerFrame calculates amount of points we need to store for current time frame.
// Based on user configuration. Rounded down to nearest int.
func calcPointsPerFrame(mtf, pollInt int) int {
	return int(math.Floor(float64(mtf / pollInt)))
}
