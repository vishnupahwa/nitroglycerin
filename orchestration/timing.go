package orchestration

import (
	"log"
	"time"
)

// CalculateStartAt returns the time at the nearest minute after the provided nextMin minutes have passed
// e.g At the epoch, 1970-01-01T00:00:00Z, with nextMin = 2, CalculateStartAt returns 1970-01-01T00:02:00Z in Unix time.
func CalculateStartAt(now time.Time, nextMin int) int64 {
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()+nextMin, 0, 0, now.Location()).Unix()
}

// WaitForStart blocks until the provided time is reached, checking on the minute
func WaitForStart(at int64) {
	startTime := time.Unix(at, 0)
	if time.Now().After(startTime) {
		return
	}

	log.Println("Waiting for execution at " + startTime.String())
	// Tick on the minute
	t := minuteTicker()

	for {
		// Wait for ticker to send
		<-t.C

		// Update the ticker
		t = minuteTicker()
		time.Sleep(1 * time.Nanosecond)
		if time.Now().After(startTime) {
			return
		}

		log.Println("Time remaining: " + time.Until(startTime).String())
	}
}

// minuteTicker returns a new ticker that triggers on the minute
func minuteTicker() *time.Ticker {
	return time.NewTicker(time.Second * time.Duration(60-time.Now().Second()))
}
