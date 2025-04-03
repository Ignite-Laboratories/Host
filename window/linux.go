//go:build linux

package window

import (
	"fmt"
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
	"github.com/ignite-laboratories/host"
	"github.com/ignite-laboratories/host/x11"
	"log"
	"runtime"
)

func init() {
	fmt.Println("[host] - Linux - sparking X window management")
}

// handles maps between entity identifiers and created window handles.
var handles = make(map[uintptr]uint64)

func Create(size std.XY[int]) uint64 {
	// Ensures all OpenGL/Window calls remain on the same OS thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Create the window
	window, err := host.X.CreateWindow(0, 0, size.X, size.Y)
	if err != nil {
		log.Fatalf("Failed to create window: %s\n", err)
	}
	if core.Verbose {
		fmt.Printf("[x11] Window created: ID = %d\n", window.ID)
	}

	// Map the new window to an entity ID
	mutex.Lock()
	entityID := core.NextID()
	Handles[entityID] = uintptr(window.ID)
	handles[uintptr(window.ID)] = entityID
	Count++
	mutex.Unlock()

	// Set input events and enable the `WM_DELETE_WINDOW` protocol.
	host.X.SelectInput(window, x11.ExposureMask|x11.KeyPressMask|x11.StructureNotifyMask) // Use event constants
	err = host.X.SetWindowProtocols(window)
	if err != nil {
		log.Fatalf("Failed to set WM_DELETE_WINDOW protocol: %s\n", err)
	}

	// Show the window (map it to the screen)
	host.X.ShowWindow(window)

	// Launch a separate goroutine to handle incoming events.
	go handleEvents()

	return entityID
}

// handleEvents processes incoming X11 events, such as window closing.
func handleEvents() {
	for core.Alive {
		// Wait for the next event and retrieve it
		e, err := host.X.WaitForEvent()
		if err != nil {
			log.Printf("Failed to wait for event: %s\n", err)
			continue
		}

		switch e.Type {
		case x11.ClientMessage:
			// Retrieve the window and message data
			window := host.X.GetEventWindow(e)
			data := host.X.GetClientMessageData(e)

			// Check if the event is a `WM_DELETE_WINDOW` request
			wmDeleteAtom := host.X.Atom(x11.WMDeleteWindow) // Use event.WMDeleteWindow constant
			if x11.Atom(data[0]) == wmDeleteAtom {
				if core.Verbose {
					fmt.Printf("[x11] Window closed: ID = %d\n", window.ID)
				}

				// Destroy the window
				host.X.DestroyWindow(window)
				host.X.Flush()

				// Clean up the internal resources
				mutex.Lock()
				entityID := handles[uintptr(window.ID)]
				delete(handles, uintptr(window.ID))
				delete(Handles, entityID)
				Count--
				mutex.Unlock()
			}
		}
	}
}
