// Package opengl provides common methods for working with OpenGL
package opengl

import (
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/host"
)

var ModuleName = "opengl"

func init() {
	host.Report()
	core.SubmoduleReport(host.ModuleName, ModuleName)
}

func Report() {}
