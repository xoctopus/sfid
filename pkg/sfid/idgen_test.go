package sfid_test

import (
	"context"
	"testing"
	"time"

	. "github.com/xoctopus/x/testx"

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
			_, _ = g1.ID()
		}
	})

	b.Run("Snowflake", func(b *testing.B) {
		for range b.N {
			_, _ = g2.ID()
		}
	})
}

func TestInjector(t *testing.T) {
	g := sfid.Must(
		sfid.Carry(g1)(context.Background()),
	)
	Expect(t, g, NotBeNil[sfid.IDGen]())
	g, _ = sfid.From(context.Background())
	Expect(t, g, BeNil[sfid.IDGen]())
}
