package sfid_test

import (
	"testing"
	"time"

	"github.com/xoctopus/sfid/internal/factory"
	"github.com/xoctopus/sfid/pkg/sfid"
)

var (
	base, _ = time.Parse(factory.Layout, "2025-05-21T00:00:00.000Z")
	g1      = sfid.NewDefaultIDGen(3)
	g2      = factory.NewWorker(2, 1, base, 10, 12)
)

func Benchmark(b *testing.B) {
	b.Run("Generator", func(b *testing.B) {
		for range b.N {
			_ = g1.ID()
		}
	})

	b.Run("Snowflake", func(b *testing.B) {
		for range b.N {
			_ = g2.ID()
		}
	})
}
