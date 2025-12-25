package sfid

import (
	"time"

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

func NewIDGen(worker uint32) IDGen {
	return factory.NewWorker(worker, 1, base, 10, 12)
}
