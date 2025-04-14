package hydra

import "github.com/ignite-laboratories/core"

type Manageable interface {
	Initialize()
	Impulse(core.Context)
	Cleanup()
}
