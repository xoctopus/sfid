package factory_test

import (
	"sort"
	"sync"
	"sync/atomic"
	"testing"

	. "github.com/xoctopus/x/testx"

	. "github.com/xoctopus/sfid/internal/factory"
)

func BenchmarkSnowflake_ID(b *testing.B) {
	var (
		factories []*Factory
		worker    int
		seq       int
	)

	gen := func() {
		var f1, f2, f4 *Factory
		defer func() {
			if r := recover(); r == nil {
				factories = append(factories, f1, f2, f4)
			} else {
				b.Log(r)
			}
		}()
		f1 = NewFactory(1, base_, worker, seq)
		f2 = NewFactory(2, base_, worker, seq)
		f4 = NewFactory(4, base_, worker, seq)
	}

	// The generated factories can support at least 64 workers and generating up
	// to 128 sfids per gap.
	for bits := 16; bits <= 62; bits++ {
		for worker = 10; worker <= bits; worker++ {
			seq = bits - worker
			if worker < 10 || seq < 10 {
				continue
			}
			gen()
		}
	}
	sort.Slice(factories, func(i, j int) bool {
		fi, fj := factories[i], factories[j]

		if fi.Unit() != fj.Unit() {
			return fi.Unit() < fj.Unit()
		}
		if fi.WorkerBits() != fj.WorkerBits() {
			return fi.WorkerBits() < fj.WorkerBits()
		}
		if fi.SeqBits() != fj.SeqBits() {
			return fi.SeqBits() < fj.SeqBits()
		}
		return fi.GapBits() < fj.GapBits()
	})

	for _, f := range factories {
		b.Run(f.Tag(), func(b *testing.B) {
			s := f.New(1)
			defer func() {
				if e := recover(); e != nil {
					b.Log(e)
				}
			}()
			for range b.N {
				_, _ = s.ID()
			}
		})
	}
}

type SnowflakeTestSuite struct {
	*testing.T
	N    int
	m    sync.Map
	size atomic.Int64
}

func NewSnowflakeTestSuite(t *testing.T, n int) *SnowflakeTestSuite {
	return &SnowflakeTestSuite{T: t, N: n}
}

func (s *SnowflakeTestSuite) ExpectN(n int) {
	Expect(s.T, s.size.Load(), Equal(int64(n)))
}

func (s *SnowflakeTestSuite) Run(sf *Worker) {
	for range s.N {
		if id, err := sf.ID(); err == nil {
			s.m.Store(id, struct{}{})
			s.size.Add(1)
		}
	}
}

func TestSnowflake_ID(t *testing.T) {
	gap, worker, seq := 1, 10, 12
	f := NewFactory(1, base_, 10, 12)
	Expect(t, f.Unit(), Equal(gap))
	Expect(t, f.SeqBits(), Equal(seq))
	Expect(t, f.WorkerBits(), Equal(worker))

	g1 := f.New(1)
	g2 := NewWorker(1, 1, base_, 10, 12)
	Expect(t, g1.WorkerID(), Equal(g2.WorkerID()))
	Expect(t, g1.Tag(), Equal(g2.Tag()))

	t.Run(f.String()+"_1x", func(t *testing.T) {
		g := NewFactory(1, base_, 4, 4).New(1)

		for i := 0; i < 10000; i++ {
			func() {
				defer func() {
					Expect(t, recover(), BeNil[any]())
				}()
				g.MustID()
			}()
		}
	})

	t.Run(f.Tag()+"_1000x", func(t *testing.T) {
		suite := NewSnowflakeTestSuite(t, 1000)
		g := f.New(1)

		con := 1000
		wg := &sync.WaitGroup{}

		for i := 0; i < con; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				suite.Run(g)
			}()
		}

		wg.Wait()
		suite.ExpectN(suite.N * con)
	})

	t.Run(f.Tag()+"InvalidClock", func(t *testing.T) {
		t.Skip("xgo is not support go1.25")
		// g1.ID()
		// now := time.Now()

		// defer func() {
		// 	NewWithT(t).Expect(recover().(string)).To(Equal("invalid system clock, clock moved backwards"))
		// }()

		// mock.Patch(time.Now, func() time.Time { return now.Add(0 - 10*time.Second) })
		// g1.ID()
	})
}
