package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
	"github.com/ignite-laboratories/host/graphics"
	"github.com/ignite-laboratories/host/hydraold"
)

type SolidColorWindow struct {
	*graphics.RenderableWindow
	Color std.RGBA
}

func NewSolidColorWindow(color std.RGBA) *SolidColorWindow {
	return &SolidColorWindow{
		Color: color,
	}
}

func (w *SolidColorWindow) Render() {
	// Clear the window with a background color
	gl.ClearColor(w.Color.R, w.Color.G, w.Color.B, w.Color.A)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func main() {
	for i := 0; i < 7; i++ {
		graphics.SparkRenderableWindow(std.XY[int]{X: 640, Y: 480}, NewSolidColorWindow(std.RandomRGB()))
	}
	core.Impulse.StopWhen(hydraold.StopPotential)
	core.Impulse.Spark()
}
