package sfid

import (
	"context"
	"time"

	"github.com/xoctopus/x/contextx"
	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sfid/internal/factory"
)

var base time.Time

func init() {
	var (
		err   error
		input = "2025-05-21T00:00:00.000Z"
	)
	base, err = time.Parse(factory.Layout, input)
	must.NoErrorF(err, "failed to parse base timestamp: %s", input)
}

type IDGen interface {
	ID() int64
}

func NewDefaultIDGen(worker uint32) IDGen {
	return factory.NewWorker(worker, 1, base, 10, 12)
}

// NewIDGen
// worker: worker id
// unit: snowflake clock unit(milliseconds)
// base: snowflake epoch timestamp
// wbs: worker bits
// sbs: sequence bits
func NewIDGen(worker uint32, unit int, base time.Time, wbs, sbs int) IDGen {
	return factory.NewWorker(worker, unit, base, wbs, sbs)
}

type k struct{}

func With(ctx context.Context, g IDGen) context.Context {
	return context.WithValue(ctx, k{}, g)
}

func From(ctx context.Context) (IDGen, bool) {
	if l, ok := ctx.Value(k{}).(IDGen); ok {
		return l, true
	}
	return nil, false
}

func Must(ctx context.Context) IDGen {
	g, ok := From(ctx)
	must.BeTrueF(ok, "missing sfid.IDGen")
	return g
}

func Carry(g IDGen) contextx.Carrier {
	return func(ctx context.Context) context.Context {
		return With(ctx, g)
	}
}
