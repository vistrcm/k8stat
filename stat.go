package main

import (
	"fmt"
	"math"
	"time"
)

type Stat struct {
	averageMem float64
	n          int // number of measurements
	lastUpdate time.Time
}

type Statmap map[string]*Stat

func NewStatMap() Statmap {
	return make(Statmap)
}

func (s *Statmap) Add(name string, bytes int) {
	if _, ok := (*s)[name]; !ok {
		(*s)[name] = NewStat()
	}
	(*s)[name].addMeasurement(bytes)
}

func (s *Statmap) Print() {
	for pod, v := range *s {
		fmt.Printf("%s: %d (%d samples)\n", pod, v.AvgMem(), v.n)
	}
}

func (s *Stat) AvgMem() int {
	return int(math.Ceil(s.averageMem))
}

func NewStat() *Stat {
	return &Stat{
		n:          0,
		averageMem: 0,
		lastUpdate: time.Now().UTC(),
	}
}

func (s *Stat) addMeasurement(data int) {
	nf := float64(s.n)
	np1f := float64(s.n + 1)

	s.averageMem = s.averageMem*nf/np1f + float64(data)/np1f
	s.n += 1
	s.lastUpdate = time.Now().UTC()
}
