package sfid

import "github.com/xoctopus/sfid/internal/factory"

type (
	Factory = factory.Factory
	Worker  = factory.Worker
)

var (
	NewFactory = factory.NewFactory
	NewWorker  = factory.NewWorker
)
