package timer

import (
	"fmt"
	"time"
)

type Timer struct {
	target       time.Duration
	started      time.Time
	timeEllapsed time.Duration
	isTimeout    bool
}

func New(target time.Duration) Timer {
	return Timer{
		started: time.Now(),
	}
}

func (t Timer) TimeLeft() string {
	timeLeft := t.target - t.timeEllapsed
	minutes := int(timeLeft.Minutes())
	seconds := int(timeLeft.Seconds()) - (minutes * 60)

	if minutes > 0 {
		return fmt.Sprintf("%d minutes %d seconds", minutes, seconds)
	}

	return fmt.Sprintf("%d seconds", seconds)
}

func (t *Timer) Update() time.Duration {
	t.timeEllapsed = time.Since(t.started)
	t.isTimeout = t.timeEllapsed >= t.target
	return t.timeEllapsed
}

func (t *Timer) Reset(target time.Duration) {
	t.target = target
	t.started = time.Now()
	t.timeEllapsed = time.Second * 0
	t.isTimeout = false
}

func (t Timer) IsTimeout() bool {
	return t.isTimeout
}

func (t Timer) PercentageLeft() float64 {
	percentage := t.timeEllapsed.Seconds() / t.target.Seconds()

	if percentage < 0 {
		return float64(0)
	}

	if percentage > 100 {
		return float64(100)
	}

	return percentage
}
