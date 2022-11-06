package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"time"
)

type Persistence interface {
	Save(names string, stat *Stat) error
}

type Stat struct {
	averageMem float64
	n          int // number of measurements
	lastUpdate time.Time
}

type StatMap struct {
	db   Persistence
	data map[string]*Stat
}

func NewStatMap(db Persistence) StatMap {
	return StatMap{
		data: make(map[string]*Stat),
		db:   db,
	}
}

func (s *StatMap) Add(name string, bytes int) {
	if _, ok := (*s).data[name]; !ok {
		(*s).data[name] = NewStat()
	}
	stat := (*s).data[name]
	stat.addMeasurement(bytes)
	go func() {
		if err := s.db.Save(name, stat); err != nil {
			log.Printf("ERROR saving stat for %q: %v", name, err)
		}
	}()
}

func (s *StatMap) Print() {
	for pod, v := range (*s).data {
		fmt.Printf("%s: %d (%d samples)\n", pod, v.AvgMem(), v.n)
	}
}

func floatToBytes(data float64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, data) // why little? just copy and paste from example
	if err != nil {
		return nil, fmt.Errorf("binary.Write failed: %w", err)
	}
	return buf.Bytes(), nil
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
