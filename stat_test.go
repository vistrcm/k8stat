package main

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func ExampleStat_AvgMem() {
	s := NewStat()
	data := []int{1}
	for _, d := range data {
		s.addMeasurement(d)
	}
	fmt.Println(s.AvgMem())
	// Output: 1
}

func TestStat_AvgMem(t *testing.T) {
	tests := []struct {
		name string
		data []int
		want int
	}{
		{"one", []int{1}, 1},
		{"multiOne", []int{1, 1, 1, 1, 1}, 1},
		{"hardcodeOne", []int{1, 2, 3, 4, 5}, 3},
		{"hardcodeOne", []int{1, 2, 3, 4, 5}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStat()
			for _, d := range tt.data {
				s.addMeasurement(d)
			}
			if got := s.AvgMem(); got != tt.want {
				t.Errorf("AvgMem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStat_AvgMem_Generated(t *testing.T) {
	t.Run("generated", func(t *testing.T) {
		s := NewStat()
		data := randomList(100000)
		for _, d := range data {
			s.addMeasurement(d)
		}
		want := semiClassicAVG(data)
		if got := s.AvgMem(); !almostEqual(got, want) {
			t.Errorf("AvgMem() = %v, want %v", got, want)
		}
	})
}

func almostEqual(got int, want int) bool {
	equalThreshold := 0.0000001
	if math.Abs(float64(got)-float64(want)) < equalThreshold {
		return true
	}
	return false
}

func randomList(i int) []int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Perm(i + 1)
}

func semiClassicAVG(data []int) int {
	var sum int
	for _, d := range data {
		sum += d
	}
	return int(math.Ceil(float64(sum) / float64(len(data))))
}
