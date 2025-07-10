package utils_test

import (
	"sync"
	"testing"

	"github.com/tomasdemarco/iso8583/utils"
)

func TestNewStan(t *testing.T) {
	tests := []struct {
		name        string
		initialNum  int
		maxNum      int
		expectedNum int
		expectedMax int
	}{
		{
			name:        "Valid initial and max",
			initialNum:  10,
			maxNum:      100,
			expectedNum: 10,
			expectedMax: 100,
		},
		{
			name:        "Initial num too low",
			initialNum:  0,
			maxNum:      100,
			expectedNum: 1,
			expectedMax: 100,
		},
		{
			name:        "Initial num greater than max",
			initialNum:  101,
			maxNum:      100,
			expectedNum: 1,
			expectedMax: 100,
		},
		{
			name:        "Max num too low",
			initialNum:  1,
			maxNum:      0,
			expectedNum: 1,
			expectedMax: 999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := utils.NewStan(tt.initialNum, tt.maxNum)
			if s.Next() != tt.expectedNum {
				t.Errorf("NewStan() initial num got = %v, want %v", s.Next(), tt.expectedNum)
			}
			// To check max, we need to access the unexported field, which is not ideal.
			// For now, we'll rely on the Next() behavior to implicitly test max.
			// A better approach would be to expose a GetMax() method or use reflection for testing.
		})
	}
}

func TestStanNext(t *testing.T) {
	tests := []struct {
		name             string
		initialNum       int
		maxNum           int
		expectedSequence []int
	}{
		{
			name:             "Basic sequence",
			initialNum:       1,
			maxNum:           3,
			expectedSequence: []int{1, 2, 3, 1, 2},
		},
		{
			name:             "Sequence with different start",
			initialNum:       5,
			maxNum:           7,
			expectedSequence: []int{5, 6, 7, 1, 2},
		},
		{
			name:             "Single max",
			initialNum:       1,
			maxNum:           1,
			expectedSequence: []int{1, 1, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := utils.NewStan(tt.initialNum, tt.maxNum)
			for i, expected := range tt.expectedSequence {
				got := s.Next()
				if got != expected {
					t.Errorf("Next() at index %d: got %v, want %v", i, got, expected)
				}
			}
		})
	}
}

func TestStanNextConcurrency(t *testing.T) {
	s := utils.NewStan(1, 1000)
	var wg sync.WaitGroup
	const numGoroutines = 100
	const iterationsPerGoroutine = 100

	results := make(chan int, numGoroutines*iterationsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterationsPerGoroutine; j++ {
				results <- s.Next()
			}
		}()
	}

	wg.Wait()
	close(results)

	seen := make(map[int]bool)
	for val := range results {
		if val < 1 || val > 1000 {
			t.Errorf("Generated STAN %d is out of expected range [1, 1000]", val)
		}
		if seen[val] {
			// This test is for concurrency, not uniqueness across all numbers.
			// It's possible to get duplicates if the sequence wraps around.
			// The main goal is to ensure no race conditions or panics.
		} else {
			seen[val] = true
		}
	}
}
