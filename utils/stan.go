// Package utils provides various utility functions used across the ISO 8583 library.
package utils

import (
	"sync"
)

// Stan represents a System Trace Audit Number (STAN) generator.
// It is thread-safe and generates sequential numbers from 1 to 999999, then wraps around.
type Stan struct {
	mu  sync.Mutex
	num int
	max int
}

// NewStan creates and returns a new Stan generator, initialized to initialNum.
// The STAN will wrap around after maxNum.
func NewStan(initialNum, maxNum int) *Stan {
	if initialNum <= 0 || initialNum > maxNum {
		initialNum = 1 // Default to 1 if initialNum is invalid
	}
	if maxNum <= 0 {
		maxNum = 999999 // Default max if maxNum is invalid
	}
	return &Stan{
		num: initialNum,
		max: maxNum,
	}
}

// Next generates and returns the next sequential STAN.
// It is thread-safe and ensures the STAN wraps around after 999999.
func (s *Stan) Next() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer func() {
		s.num++
		if s.num > s.max {
			s.num = 1
		}
	}()

	return s.num
}
