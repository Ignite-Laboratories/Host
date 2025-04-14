package hydra

import (
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

// Head represents a control surface for neural rendering.
type Head struct {
	*core.System

	WindowID uint32
	Window   *sdl.Window
	Driver   *core.Neuron

	synchro    std.Synchro
	manageable Manageable
}

func (w *Head) run() {
	w.manageable.Initialize()

	for core.Alive && w.Alive {
		w.synchro.Engage()

		// Hydra should not be managing anything above 1,000Hz!
		// Why waste the CPU cycles?
		time.Sleep(time.Millisecond)
	}

	w.manageable.Cleanup()
	core.Verbosef(ModuleName, "head [%d.%d] destroyed", w.WindowID, w.ID)
}

func (w *Head) impulse(ctx core.Context) {
	w.synchro.Send(func() {
		w.manageable.Impulse(ctx)

		w.Window.GLSwap()
	})
}
