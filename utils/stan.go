package utils

import (
	"sync"
)

type Stan struct {
	mu  sync.Mutex
	num int
}

func NewStan() *Stan {
	return &Stan{
		num: 1,
	}
}

func (s *Stan) Next() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer func() {
		s.num++
		if s.num > 999999 {
			s.num = 1
		}
	}()

	return s.num
}
