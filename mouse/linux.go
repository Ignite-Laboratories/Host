//go:build linux

package mouse

import "C"
import (
	"fmt"
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
	"github.com/ignite-laboratories/host/x11"
	"log"
)

func init() {
	fmt.Println("[host] - Linux - sparking X mouse observance")
	var err error
	x, err = x11.OpenDisplay()
	if err != nil {
		log.Fatalf("Failed to initialize X11: %s\n", err)
	}

	go func() {
		core.WhileAlive()
		x11.CloseDisplay(x)
		fmt.Println("[host] - Linux - closed X mouse observance")
	}()
}

// X provides a handle to the global connection to the X11 server.
var x *x11.Display

// Sample gets the current mouse coordinates, or nil if unable to do so.
func Sample() *std.MouseState {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("failed to get mouse position: %v\n", r)
		}
	}()

	rootWin, _ := x11.GetRootWindow(x)
	data, err := x11.QueryPointer(x, rootWin)
	if err != nil {
		fmt.Printf("failed to get mouse position: %v", err)
	}
	return &data
}
