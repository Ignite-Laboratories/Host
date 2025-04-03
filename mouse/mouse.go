package mouse

import (
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
	"github.com/ignite-laboratories/core/temporal"
	"github.com/ignite-laboratories/core/when"
)

// SampleRate provides a standard frequency for sampling the mouse.
var SampleRate = 2048.0

// State provides an observational dimension that samples the mouse state at the provided frequency.
var State = temporal.Observer(core.Impulse, when.Frequency(&SampleRate), true, Sample)

// Reaction creates a reactionary dimension that samples the mouse and invokes onChange at the provided frequency.
func Reaction(engine *core.Engine, frequency *float64, onChange temporal.Change[std.MouseState]) *temporal.Dimension[std.MouseState, any] {
	return temporal.Reaction[std.MouseState](engine, when.Frequency(frequency), false, Sample, onChange)
}
