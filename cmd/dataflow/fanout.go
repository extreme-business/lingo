package dataflow

import "sync"

type FanOut[T any] struct {
	In  chan T
	Out []chan T
	wg  sync.WaitGroup
}

func NewFanOut[T any]() *FanOut[T] {
	return &FanOut[T]{
		In:  make(chan T),
		Out: make([]chan T, 0),
		wg:  sync.WaitGroup{},
	}
}

func (s *FanOut[T]) New() chan T {
	out := make(chan T)
	s.Out = append(s.Out, out)
	return out
}

func (s *FanOut[T]) Run() {
	s.wg.Add(1)
	defer s.wg.Done()

	for v := range s.In {
		for _, out := range s.Out {
			out <- v
		}
	}
}

func (s *FanOut[T]) Close() {
	close(s.In)
	s.wg.Wait()

	for _, out := range s.Out {
		close(out)
	}
}
