//go:build linux

package mouse

import "C"
import (
	"fmt"
	"github.com/ignite-laboratories/core/std"
	"github.com/ignite-laboratories/host"
)

func init() {
	fmt.Println("[host] - Linux - sparking X mouse management")
}

// Sample gets the current mouse coordinates, or nil if unable to do so.
func Sample() *std.MouseState {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("failed to get mouse position: %v\n", r)
		}
	}()
	data, err := host.X.QueryPointer(host.X.RootWindow())
	if err != nil {
		fmt.Printf("failed to get mouse position: %v", err)
	}
	return &data
}
