package factory

import (
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

func (s *Worker) ID() int64 {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	gap := s.f.Elapsed()
	if s.gap < gap {
		s.gap = gap
		// s.seq = _rand.Uint32N(uint32(s.f.seq.max))
		s.seq = 1
		return s.hi | int64(s.seq)<<s.f.gap.bits | gap
		// return s.f.Build(s.id, s.seq, s.gap)
	}

	if s.gap > gap {
		delta := (s.gap - gap) * int64(s.f.unit)
		must.BeTrueF(
			delta <= int64(5*time.Millisecond),
			"invalid system clock, clock moved backwards %dms", delta/int64(time.Millisecond),
		)
		gap = s.gap
	}

	// in same gap. if reached max sequence, need to wait next gap to make sure
	// generated in ascending order
	s.seq = s.f.Mask(s.seq + 1)
	if s.seq == 0 {
		s.gap++
		s.f.Next(s.gap - gap)
	}
	return s.hi | int64(s.seq)<<s.f.gap.bits | gap
	// return s.f.Build(s.id, s.seq, s.gap)
}

/*
type Worker struct {
	state uint64 // gap | seq
	f     *Factory
	id    uint32
	hi    int64
}

func (s *Worker) ID() int64 {
	seqBits := s.f.seq.bits
	seqMask := uint32(1<<seqBits - 1)

	for {
		state := atomic.LoadUint64(&s.state)
		prevGap := int64(state >> seqBits)
		prevSeq := uint32(state & uint64(seqMask))

		var nextGap int64
		var nextSeq uint32

		if nowGap := s.f.Elapsed(); nowGap > prevGap {
			// gap moved forwards. reset seq
			nextGap = nowGap
			nextSeq = 1
		} else {
			// in same gap or little backwards.
			// add seq
			nextGap = prevGap
			nextSeq = prevSeq + 1

			// seq overflowed. waiting next gap
			if nextSeq > seqMask {
				nextGap = prevGap + 1
				nextSeq = 1
				for s.f.Elapsed() < nextGap {
					runtime.Gosched()
				}
			}
		}

		newState := uint64(nextGap<<seqBits) | uint64(nextSeq)
		// must-one-routine changed state
		if atomic.CompareAndSwapUint64(&s.state, state, newState) {
			return s.hi | (nextGap << seqBits) | int64(nextSeq)
		}
		// failed. try again
	}
}
*/
