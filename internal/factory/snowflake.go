package factory

import (
	"fmt"
	"sync"
	"time"

	"github.com/xoctopus/x/misc/must"
)

func NewWorker(worker uint32, unit int, base time.Time, w, s int) *Worker {
	ww := NewFactory(unit, base, w, s).New(worker)
	ww.hi = int64(ww.id) << (ww.f.gap.bits + ww.f.seq.bits)
	return ww
}

type Worker struct {
	hi  int64
	id  uint32
	seq uint32
	gap int64
	f   *Factory
	mtx *sync.Mutex
}

func (s *Worker) WorkerID() uint32 { return s.id }

func (s *Worker) Tag() string { return s.f.Tag() }

func (s *Worker) ID() (int64, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	gap := s.f.Elapsed()
	if s.gap < gap {
		s.gap = gap
		// s.seq = _rand.Uint32N(uint32(s.f.seq.max))
		s.seq = 1
		return s.hi | int64(s.seq)<<s.f.gap.bits | gap, nil
		// return s.f.Build(s.id, s.seq, s.gap)
	}

	if s.gap > gap {
		delta := (s.gap - gap) * int64(s.f.unit)
		if delta <= int64(5*time.Millisecond) {
			return 0, fmt.Errorf("invalid system clock, clock moved backwards %dms", delta/int64(time.Millisecond))
		}
		gap = s.gap
	}

	// in same gap. if reached max sequence, need to wait next gap to make sure
	// generated in ascending order
	s.seq = s.f.Mask(s.seq + 1)
	if s.seq == 0 {
		s.gap++
		s.f.Next(s.gap - gap)
	}
	return s.hi | int64(s.seq)<<s.f.gap.bits | gap, nil
}

func (s *Worker) MustID() int64 {
	return must.NoErrorV(s.ID())
}
