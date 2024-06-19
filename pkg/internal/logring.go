package internal

import (
	"container/ring"
	"fmt"
)

type LogRing struct {
	ring *ring.Ring
}

func NewLogRing(length int) *LogRing {
	return &LogRing{
		ring.New(length),
	}
}

// Len calculates the non-nil length of the internal ring
func (r *LogRing) Len() int {
	l := 0

	for i := 0; i < r.ring.Len(); i++ {
		if r.ring.Value != nil {
			l++
		}
		r.ring = r.ring.Next()
	}

	return l
}

// Print returns a slice of non-nil ring values
func (r *LogRing) Print() []string {
	var logStrings []string

	for i := 0; i < r.ring.Len(); i++ {
		if r.ring.Value != nil {
			logStrings = append(logStrings, r.ring.Value.(string))
		}
		r.ring = r.ring.Next()
	}

	return logStrings
}

// Truncate removes all the elements of the internal ring
func (r *LogRing) Truncate() *LogRing {
	r.ring.Unlink(r.ring.Len() - 1)
	return r
}

// Append sets a string to the current internal ring value, and then moves the ring pointer forward
func (r *LogRing) Append(line string) {
	r.ring.Value = fmt.Sprintf("%s", line)
	r.ring = r.ring.Next()
}
