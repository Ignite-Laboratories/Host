// Package host provides a toolkit for interfacing with a host architecture.
package host

import (
	"github.com/ignite-laboratories/core"
)

var ModuleName = "host"

func init() {
	core.ModuleReport(ModuleName)
}

func Report() {}
