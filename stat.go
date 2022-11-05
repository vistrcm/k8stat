package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	badger "github.com/dgraph-io/badger/v3"
	"math"
	"time"
)

type Stat struct {
	averageMem float64
	n          int // number of measurements
	lastUpdate time.Time
}

type Statmap struct {
	db   *badger.DB
	data map[string]*Stat
}

func NewStatMap(db *badger.DB) Statmap {
	return Statmap{
		data: make(map[string]*Stat),
		db:   db,
	}
}

func (s *Statmap) Add(name string, bytes int) {
	if _, ok := (*s).data[name]; !ok {
		(*s).data[name] = NewStat()
	}
	(*s).data[name].addMeasurement(bytes)
	s.save(name, bytes)
}

func (s *Statmap) Print() {
	for pod, v := range (*s).data {
		fmt.Printf("%s: %d (%d samples)\n", pod, v.AvgMem(), v.n)
	}
}

func (s *Statmap) save(name string, bytes int) error {
	err := s.db.Update(func(tnx *badger.Txn) error {
		value, err := intToBytes(bytes)
		if err != nil {
			return err
		}
		e := badger.NewEntry([]byte(name), value)
		err = tnx.SetEntry(e)
		return err
	})
	return err
}

func intToBytes(data int) ([]byte, error) {
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
