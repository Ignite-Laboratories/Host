package mouse

import (
	"fmt"
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
	"github.com/ignite-laboratories/core/temporal"
	"github.com/ignite-laboratories/core/when"
)

// SampleRate provides a standard frequency for sampling the mouse.
var SampleRate = 2048.0

// Coordinates provides an observational dimension that samples the mouse coordinates at the provided frequency.
var Coordinates = temporal.Observer(core.Impulse, when.Frequency(&SampleRate), true, SampleCoordinates)

// Reaction creates a reactionary dimension that samples the mouse and invokes onChange at the provided frequency.
func Reaction(engine *core.Engine, frequency *float64, onChange temporal.Change[std.XY[int]]) *temporal.Dimension[std.XY[int], any] {
	return temporal.Reaction[std.XY[int]](engine, when.Frequency(frequency), false, SampleCoordinates, onChange)
}

// SampleCoordinates gets the current mouse coordinates, or nil if unable to do so.
func SampleCoordinates() *std.XY[int] {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(fmt.Errorf("failed to get mouse position: %v", r))
		}
	}()
	coords := getCoordinates()
	return &coords
}
