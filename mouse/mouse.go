package mouse

import (
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
	"github.com/ignite-laboratories/core/temporal"
	"github.com/ignite-laboratories/core/when"
	_ "github.com/ignite-laboratories/host/hydra"
	"github.com/veandco/go-sdl2/sdl"
)

// SampleRate provides a standard frequency for sampling the mouse.
var SampleRate = 2048.0

// State provides an observational dimension that samples the mouse state at the provided frequency.
var State = temporal.Observer(core.Impulse, when.Frequency(&SampleRate), true, Sample)

// Reaction creates a reactionary dimension that samples the mouse and invokes onChange at the provided frequency.
func Reaction(engine *core.Engine, frequency *float64, onChange temporal.Change[std.MouseState]) *temporal.Dimension[std.MouseState, any] {
	return temporal.Reaction[std.MouseState](engine, when.Frequency(frequency), false, Sample, onChange)
}

// Sample gets the current mouse coordinates, or nil if unable to do so.
func Sample() *std.MouseState {
	x, y, st := sdl.GetGlobalMouseState()
	state := sdlStateToStd(st)
	state.Position.X = int(x)
	state.Position.Y = int(y)
	return &state
}

func sdlStateToStd(state uint32) (out std.MouseState) {
	if state&sdl.ButtonLMask() != 0 {
		out.Buttons.Left = true
	}

	if state&sdl.ButtonMMask() != 0 {
		out.Buttons.Middle = true
	}

	if state&sdl.ButtonRMask() != 0 {
		out.Buttons.Right = true
	}

	if state&sdl.ButtonX1Mask() != 0 {
		out.Buttons.Extra1 = true
	}

	if state&sdl.ButtonX2Mask() != 0 {
		out.Buttons.Extra2 = true
	}

	return out
}

// SampleRelative gets the current mouse coordinates relative to a window, or nil if unable to do so.
//func SampleRelative(window hydraold.Handle) *std.MouseState {
//	defer func() {
//		if r := recover(); r != nil {
//			fmt.Printf("failed to get mouse position: %v\n", r)
//		}
//	}()
//
//	rootWin, _ := host.GetRootWindow(x)
//	data, err := host.QueryPointer(x, rootWin)
//	if err != nil {
//		fmt.Printf("failed to get mouse position: %v", err)
//	}
//
//	state := host.PointerQueryToState(data, false)
//	return &state
//}
