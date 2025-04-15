// Package sdl2 provides a way to create impulsable graphical contexts using SDL2
package sdl2

import (
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/host"
)

var ModuleName = "sdl2"

func init() {
	host.Report()
	core.SubmoduleReport(host.ModuleName, ModuleName)
}

func Report() {}
