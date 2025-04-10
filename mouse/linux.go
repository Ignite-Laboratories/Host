//go:build linux

package mouse

import "C"
import (
	"fmt"
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
	"github.com/ignite-laboratories/host/hydraold"
	"log"
)

func init() {
	fmt.Println("[host] - Linux - sparking X mouse access")
	var err error
	x, err = hydraold.OpenDisplay()
	if err != nil {
		log.Fatalf("Failed to initialize X11: %s\n", err)
	}

	go func() {
		core.WhileAlive()
		hydraold.CloseDisplay(x)
		fmt.Println("[host] - Linux - closed X mouse access")
	}()
}

// X provides a handle to the global connection to the X11 server.
var x *hydraold.Display

// Sample gets the current mouse coordinates, or nil if unable to do so.
func Sample() *std.MouseState {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("failed to get mouse position: %v\n", r)
		}
	}()

	rootWin, _ := hydraold.GetRootWindow(x)
	data, err := hydraold.QueryPointer(x, rootWin)
	if err != nil {
		fmt.Printf("failed to get mouse position: %v", err)
	}

	state := hydraold.PointerQueryToState(data, true)
	return &state
}

// SampleRelative gets the current mouse coordinates relative to a window, or nil if unable to do so.
func SampleRelative(window hydraold.Handle) *std.MouseState {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("failed to get mouse position: %v\n", r)
		}
	}()

	rootWin, _ := hydraold.GetRootWindow(x)
	data, err := hydraold.QueryPointer(x, rootWin)
	if err != nil {
		fmt.Printf("failed to get mouse position: %v", err)
	}

	state := hydraold.PointerQueryToState(data, false)
	return &state
}
