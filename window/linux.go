//go:build linux

package window

import (
	"fmt"
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
	"github.com/ignite-laboratories/host/x11"
	"log"
	"runtime"
)

type Handle struct {
	core.Entity
	Display   *x11.Display
	Window    *x11.Window
	Destroyed bool
}

func SetTitle(handle Handle, title string) {
	x11.StoreName(handle.Display, handle.Window, title)
}

func Create(size std.XY[int]) *Handle {
	// Ensures all OpenGL/Window calls remain on the same OS thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	display, err := x11.OpenDisplay()
	if err != nil {
		log.Fatalf("Failed to initialize X11: %s\n", err)
	}

	// Create the window
	window, err := x11.CreateWindow(display, 0, 0, size.X, size.Y)
	if err != nil {
		log.Fatalf("Failed to create window: %s\n", err)
	}
	if core.Verbose {
		fmt.Printf("[x11] Window created: ID = %d\n", window.ID)
	}

	// Map the new window to an entity ID
	handle := &Handle{Display: display, Window: window}
	handle.ID = core.NextID()
	Handles[handle.ID] = *handle

	// Enable detection of a 'close' event
	x11.SelectInput(display, window, x11.StructureNotifyMask)
	err = x11.SetWindowProtocols(display, window)
	if err != nil {
		log.Fatalf("Failed to set WM_DELETE_WINDOW protocol: %s\n", err)
	}

	// Show the window (map it to the screen)
	x11.ShowWindow(display, window)

	// Launch a separate goroutine to handle incoming events.
	go handleEvents(handle)

	return handle
}

// handleEvents processes incoming X11 events, such as window closing.
func handleEvents(handle *Handle) {
	defer x11.CloseDisplay(handle.Display)

	for !handle.Destroyed && core.Alive {
		// Wait for the next event and retrieve it
		e, err := x11.WaitForEvent(handle.Display)
		if err != nil {
			log.Printf("Failed to wait for event: %s\n", err)
			continue
		}

		switch e.Type {
		case x11.ClientMessage:
			// Retrieve the window and message data
			window := x11.GetEventWindow(e)
			data, _ := x11.GetClientMessageData(e)

			// Check if the event is a `WM_DELETE_WINDOW` request
			wmDeleteAtom := x11.GetAtom(handle.Display, x11.WMDeleteWindow)
			if x11.Atom(data[0]) == wmDeleteAtom {
				if core.Verbose {
					fmt.Printf("[x11] Window closed: ID = %d\n", window.ID)
				}

				// Destroy the window
				x11.DestroyWindow(handle.Display, window)
				x11.Flush(handle.Display)
				delete(Handles, handle.ID)
				handle.Destroyed = true
			}
		}
	}
}
