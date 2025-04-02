package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/host/graphics"
	"github.com/ignite-laboratories/host/windowing"
	"github.com/remogatto/egl"
	"math/rand"
)

func main() {
	graphics.Setup(windowing.CreateWindow(), red)
	graphics.Setup(windowing.CreateWindow(), blue)
	graphics.Setup(windowing.CreateWindow(), render)
	graphics.Setup(windowing.CreateWindow(), render)
	graphics.Setup(windowing.CreateWindow(), render)

	core.Impulse.StopWhen(windowing.StopPotential)
	core.Impulse.Spark()
}

func red(display egl.Display, surface egl.Surface) {
	for {
		// Clear the window with a background color
		gl.ClearColor(1.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Swap the buffers to display the rendered frame
		egl.SwapBuffers(display, surface)
	}
}

func blue(display egl.Display, surface egl.Surface) {
	for {
		// Clear the window with a background color
		gl.ClearColor(0.0, 0.0, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Swap the buffers to display the rendered frame
		egl.SwapBuffers(display, surface)
	}
}

func render(display egl.Display, surface egl.Surface) {
	r := rand.Float32()
	g := rand.Float32()
	b := rand.Float32()
	for {
		// Clear the window with a background color
		gl.ClearColor(r, g, b, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Swap the buffers to display the rendered frame
		egl.SwapBuffers(display, surface)
	}
}
