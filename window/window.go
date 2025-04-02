package window

import (
	"github.com/ignite-laboratories/core"
)

// Count provides the number of open x windows.
var Count int32

// StopPotential provides a potential that returns true when all of the windows have been globally closed.
func StopPotential(ctx core.Context) bool {
	return Count == 0
}
