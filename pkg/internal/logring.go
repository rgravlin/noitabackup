package internal

import (
	"container/ring"
	"log"
)

type LogRing struct {
	ring *ring.Ring
}

func NewLogRing(length int) *LogRing {
	return &LogRing{
		ring.New(length),
	}
}

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

func (r *LogRing) Append(line string) {
	r.ring.Value = line
	r.ring = r.ring.Next()
}

func (r *LogRing) LogAndAppend(line string) {
	r.Append(line)
	log.Println(line)
}
