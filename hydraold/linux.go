//go:build linux

package hydraold

import (
	"fmt"
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
	"log"
	"runtime"
)

type Handle struct {
	core.Entity
	Display   *Display
	Window    *Window
	Destroyed bool
}

func SetTitle(handle Handle, title string) {
	StoreName(handle.Display, handle.Window, title)
}

func Create(size std.XY[int]) *Handle {
	// Ensures all OpenGL/Window calls remain on the same OS thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	display, err := OpenDisplay()
	if err != nil {
		log.Fatalf("Failed to initialize X11: %s\n", err)
	}

	// Create the window
	window, err := CreateWindow(display, 0, 0, size.X, size.Y)
	if err != nil {
		log.Fatalf("Failed to create window: %s\n", err)
	}
	if core.Verbose {
		fmt.Printf("[hydra] Window created: ID = %d\n", window.ID)
	}

	// Map the new window to an entity ID
	handle := &Handle{Display: display, Window: window}
	handle.ID = core.NextID()
	Handles[handle.ID] = *handle

	// Enable detection of a 'close' event
	SelectInput(display, window, StructureNotifyMask)
	err = SetWindowProtocols(display, window)
	if err != nil {
		log.Fatalf("Failed to set WM_DELETE_WINDOW protocol: %s\n", err)
	}

	// Show the window (map it to the screen)
	ShowWindow(display, window)

	// Launch a separate goroutine to handle incoming events.
	go handleEvents(handle)

	return handle
}

// handleEvents processes incoming X11 events, such as window closing.
func handleEvents(handle *Handle) {
	defer CloseDisplay(handle.Display)

	for !handle.Destroyed && core.Alive {
		// Wait for the next event and retrieve it
		e, err := WaitForEvent(handle.Display)
		if err != nil {
			log.Printf("Failed to wait for event: %s\n", err)
			continue
		}

		switch e.Type {
		case ClientMessage:
			// Retrieve the window and message data
			window := GetEventWindow(e)
			data, _ := GetClientMessageData(e)

			// Check if the event is a `WM_DELETE_WINDOW` request
			wmDeleteAtom := GetAtom(handle.Display, WMDeleteWindow)
			if Atom(data[0]) == wmDeleteAtom {
				if core.Verbose {
					fmt.Printf("[hydra] Window closed: ID = %d\n", window.ID)
				}

				// Destroy the window
				DestroyWindow(handle.Display, window)
				Flush(handle.Display)
				delete(Handles, handle.ID)
				handle.Destroyed = true
			}
		}
	}
}
