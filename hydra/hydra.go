package hydra

import (
	"github.com/ignite-laboratories/core"
)

// Handles provides the pointer handles to the underlying windows by their unique entity ID.
var Handles = make(map[uint64]Handle)

// StopPotential provides a potential that returns true when all of the windows have been globally closed.
func StopPotential(ctx core.Context) bool {
	return len(Handles) == 0
}
