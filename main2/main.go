package main

import (
	"fmt"
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
	"github.com/ignite-laboratories/core/temporal"
	"github.com/ignite-laboratories/core/when"
	"github.com/ignite-laboratories/host/mouse"
)

/**
This example prints the mouse state whenever the cursor
hovers back over any part of the screen that is half the distance
from the furthest point the cursor has traveled since launching.
*/

func main() {
	temporal.Calculation(core.Impulse, when.Frequency(&mouse.SampleRate), false, CalcCoords)
	core.Impulse.Spark()
}

var highestWidth = 0

func CalcCoords(ctx core.Context) *std.MouseState {
	coords := mouse.Sample()
	if coords.Position.X < highestWidth/2 {
		fmt.Println(*coords)
	} else if coords.Position.X > highestWidth {
		highestWidth = coords.Position.X
	}
	return coords
}
