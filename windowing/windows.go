//go:build windows

package windowing

import (
	"fmt"
	"github.com/ignite-laboratories/core"
	"github.com/lxn/win"
	"sync/atomic"
	"syscall"
	"unsafe"
)

var count int32

func init() {
	// Initialize the system and launch the message loop in a new goroutine
	fmt.Println("[host] - Windows - WinAPI windowing")

	// Start a goroutine for the main message loop to handle messages for all windows
	go func() {
		var msg win.MSG
		for core.Alive {
			// Retrieve and process messages
			if win.GetMessage(&msg, 0, 0, 0) != 0 {
				win.TranslateMessage(&msg)
				win.DispatchMessage(&msg)
			}
		}
		// Cleanup and shutdown when exiting
		fmt.Println("Shutting down the message loop")
	}()
}

// StopPotential provides a common stopping potential for core impulse engines paired with the windowing system
func StopPotential(ctx core.Context) bool {
	return count == 0
}

// CreateWindow creates a new window using the WinAPI with the lxn/win library
func CreateWindow() win.HWND {
	atomic.AddInt32(&count, 1)

	// Name of the window class
	className := syscall.StringToUTF16Ptr("MyWindowClass")

	// Define a window class
	wndClass := win.WNDCLASSEX{
		CbSize:        uint32(unsafe.Sizeof(win.WNDCLASSEX{})),
		Style:         win.CS_HREDRAW | win.CS_VREDRAW,
		LpfnWndProc:   syscall.NewCallback(windowProc), // Window procedure callback
		HInstance:     win.GetModuleHandle(nil),        // Current instance handle
		HCursor:       win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_ARROW)),
		HbrBackground: win.HBRUSH(win.GetStockObject(win.WHITE_BRUSH)), // White background
		LpszClassName: className,                                       // Class name
	}

	// Register the window class
	if win.RegisterClassEx(&wndClass) == 0 {
		panic("Failed to register window class.")
	}

	// Create the window
	hwnd := win.CreateWindowEx(
		0,                                      // Extended style
		className,                              // Window class name
		syscall.StringToUTF16Ptr("My Window"),  // Window title
		win.WS_OVERLAPPEDWINDOW|win.WS_VISIBLE, // Style: Overlapped and visible
		win.CW_USEDEFAULT, win.CW_USEDEFAULT, 200, 200,
		0, 0, // Parent and menu
		wndClass.HInstance, nil, // Instance and additional parameters
	)

	// Check if window creation failed
	if hwnd == 0 {
		panic("Failed to create window")
	}

	return hwnd
}

// windowProc is the window procedure function for handling system messages
func windowProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_CLOSE:
		// Handle window close
		win.DestroyWindow(hwnd)

	case win.WM_DESTROY:
		// When the window is destroyed, decrement the counter and check shutdown condition
		atomic.AddInt32(&count, -1)
		if count == 0 && !core.Alive {
			win.PostQuitMessage(0) // Signal to quit the application
		}

	default:
		// Default behavior for messages not explicitly handled
		return win.DefWindowProc(hwnd, msg, wParam, lParam)
	}

	return 0 // Message handled
}
