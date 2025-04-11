package hydra

import "github.com/ignite-laboratories/core"

type Manageable interface {
	Initialize()
	Render(core.Context)
	Cleanup()
}
