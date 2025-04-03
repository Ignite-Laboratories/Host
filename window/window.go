package window

import (
	"github.com/ignite-laboratories/core"
	"sync"
)

// Count provides the number of open x windows.
var Count int32

// Handles provides the pointer handles to the underlying windows by their unique ID.
var Handles = make(map[uint64]uintptr)

// StopPotential provides a potential that returns true when all of the windows have been globally closed.
func StopPotential(ctx core.Context) bool {
	return Count == 0
}

var mutex sync.Mutex
