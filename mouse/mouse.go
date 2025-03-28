package mouse

import (
	"fmt"
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
	"github.com/ignite-laboratories/core/temporal"
	"github.com/ignite-laboratories/core/when"
)

func GetCoordinates() std.XY[int] {
	coords, err := getCoordinates()
	if err != nil {
		_ = fmt.Errorf("error getting mouse coordinates: %v", err)
	}
	return coords
}

var SampleRate = 1024.0
var Coordinates = temporal.Calculator(core.Impulse, when.Frequency(&SampleRate), true, coordinateSampler)

func coordinateSampler(ctx core.Context) std.XY[int] {
	coords, err := getCoordinates()
	if err != nil {
		_ = fmt.Errorf("error getting mouse coordinates: %v", err)
	}
	return coords
}
