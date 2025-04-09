//go:build linux

package hydra

/*
#cgo LDFLAGS: -lX11
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"github.com/ignite-laboratories/core/std"
	"unsafe"
)

//
// --- Types and State ---
//

type Atom = C.Atom

// Display represents a connection to an X11 server.
type Display struct {
	Ptr *C.Display
}

// Window represents a single X11-managed window.
type Window struct {
	ID C.Window
}

// Event represents an X11 event.
type Event struct {
	Type int
	Raw  C.XEvent // Raw event data (if needed for advanced use cases)
}

//
// --- Initialize and Cleanup ---
//

// OpenDisplay opens a connection to the X11 server.
func OpenDisplay() (*Display, error) {
	ptr := C.XOpenDisplay(nil)
	if ptr == nil {
		return nil, errors.New("[x11] - Failed to open X11 display")
	}
	return &Display{Ptr: ptr}, nil
}

// CloseDisplay closes the connection to the X11 server.
func CloseDisplay(display *Display) {
	if display.Ptr != nil {
		C.XCloseDisplay(display.Ptr)
		display.Ptr = nil
	}
}

//
// --- Mouse Interaction ---
//

type PointerQuery struct {
	Parent Window
	Child  Window
	RootX  int
	RootY  int
	WinX   int
	WinY   int
	Mask   uint
}

// QueryPointer retrieves the pointer's current position and its relationship to the specified window.
func QueryPointer(display *Display, window Window) (PointerQuery, error) {
	// Prepare C variables for receiving data
	var rootReturn, childReturn C.Window
	var rootXReturn, rootYReturn C.int
	var winXReturn, winYReturn C.int
	var maskReturn C.uint

	// Call XQueryPointer
	success := C.XQueryPointer(
		display.Ptr,  // Display pointer
		window.ID,    // Window ID to query
		&rootReturn,  // Pointer to receive parent (root) window
		&childReturn, // Pointer to receive child window
		&rootXReturn, // Pointer to root X coordinate
		&rootYReturn, // Pointer to root Y coordinate
		&winXReturn,  // Pointer to window-relative X coordinate
		&winYReturn,  // Pointer to window-relative Y coordinate
		&maskReturn,  // Pointer to receive button/key mask
	)

	// If XQueryPointer returns 0, it failed
	if success == 0 {
		return PointerQuery{}, errors.New("XQueryPointer failed to query pointer position")
	}

	query := PointerQuery{
		Parent: Window{ID: rootReturn},
		Child:  Window{ID: childReturn},
		RootX:  int(rootXReturn),
		RootY:  int(rootYReturn),
		WinX:   int(winXReturn),
		WinY:   int(winYReturn),
		Mask:   uint(maskReturn),
	}

	// Convert the retrieved values and return them
	return query, nil
}

// PointerQueryToState converts the provided query to a std.MouseState.
//
// root indicates if the state should grab the global coordinates or windows coordinate.
func PointerQueryToState(query PointerQuery, root bool) (state std.MouseState) {
	if root {
		state.Position = std.XY[int]{
			X: query.RootX,
			Y: query.RootY,
		}
	} else {
		state.Position = std.XY[int]{
			X: query.WinX,
			Y: query.WinY,
		}
	}

	if query.Mask&Button1Mask != 0 {
		state.Buttons.Left = true
	}
	if query.Mask&Button2Mask != 0 {
		state.Buttons.Middle = true
	}
	if query.Mask&Button3Mask != 0 {
		state.Buttons.Right = true
	}

	return state
}

//
// --- Basic Display/Screen Handling ---
//

// DefaultScreen retrieves the default screen for the display.
func DefaultScreen(display *Display) int {
	return int(C.XDefaultScreen(display.Ptr))
}

// RootWindow retrieves the root window for the default screen.
func RootWindow(display *Display) Window {
	screen := C.XDefaultScreen(display.Ptr)
	root := C.XRootWindow(display.Ptr, C.int(screen))
	return Window{ID: root}
}

//
// --- Window Management ---
//

// CreateWindow creates a basic X11 window.
func CreateWindow(display *Display, x, y, width, height int) (*Window, error) {
	screen := C.XDefaultScreen(display.Ptr)
	root := C.XRootWindow(display.Ptr, C.int(screen))

	// Create the window
	window := C.XCreateSimpleWindow(
		display.Ptr,
		root,
		C.int(x), C.int(y),
		C.uint(width), C.uint(height),
		C.uint(1), // border width
		C.XBlackPixel(display.Ptr, C.int(screen)),
		C.XWhitePixel(display.Ptr, C.int(screen)),
	)

	if window == 0 {
		return nil, errors.New("[x11] - Failed to create window")
	}

	return &Window{ID: window}, nil
}

// DestroyWindow destroys a created window.
func DestroyWindow(display *Display, win *Window) {
	C.XDestroyWindow(display.Ptr, win.ID)
}

// ShowWindow maps (shows) the window on the screen.
func ShowWindow(display *Display, win *Window) {
	C.XMapWindow(display.Ptr, win.ID)
}

// SetWindowProtocols sets standard window close behavior (WM_DELETE_WINDOW).
func SetWindowProtocols(display *Display, win *Window) error {
	atom := GetAtom(display, "WM_DELETE_WINDOW")
	if atom == 0 {
		return errors.New("[x11] - Failed to retrieve WM_DELETE_WINDOW atom")
	}

	status := C.XSetWMProtocols(display.Ptr, win.ID, &atom, 1)
	if status == 0 {
		return errors.New("[x11] - Failed to set WM_DELETE_WINDOW protocol")
	}
	return nil
}

func StoreName(display *Display, win *Window, name string) {
	C.XStoreName(display.Ptr, win.ID, C.CString(name))
}

//
// --- Event Handling ---
//

// WaitForEvent blocks until an event occurs and returns it.
func WaitForEvent(display *Display) (*Event, error) {
	var raw C.XEvent
	C.XNextEvent(display.Ptr, &raw)

	eventType := getEventType(&raw)
	return &Event{
		Type: eventType,
		Raw:  raw,
	}, nil
}

// getEventType safely extracts the event type from an XEvent.
func getEventType(event *C.XEvent) int {
	anyEvent := (*C.XAnyEvent)(unsafe.Pointer(event))
	return int(anyEvent._type)
}

// GetEventWindow extracts the window ID from an event.
func GetEventWindow(e *Event) *Window {
	anyEvent := (*C.XAnyEvent)(unsafe.Pointer(&e.Raw))
	return &Window{ID: anyEvent.window}
}

//
// --- Atom Management ---
//

// GetAtom retrieves an existing Atom or creates it if it doesn't exist.
func GetAtom(display *Display, name string) C.Atom {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return C.XInternAtom(display.Ptr, cName, 0)
}

//
// --- Utility Functions ---
//

// SelectInput sets the input mask for the window.
func SelectInput(display *Display, win *Window, eventMask int64) {
	C.XSelectInput(display.Ptr, win.ID, C.long(eventMask))
}

// Flush forces the X server to process commands in its buffer.
func Flush(display *Display) {
	C.XFlush(display.Ptr)
}

// GetClientMessageData retrieves the message data for a ClientMessage event.
func GetClientMessageData(e *Event) ([5]int64, error) {
	if e.Type != C.ClientMessage {
		return [5]int64{}, errors.New("event is not a ClientMessage")
	}
	message := (*C.XClientMessageEvent)(unsafe.Pointer(&e.Raw))
	// Convert the 5 long array (data) to a Go-compatible [5]int64
	return *(*[5]int64)(unsafe.Pointer(&message.data)), nil
}

// GetRootWindow retrieves the root window for the default screen of the given display.
func GetRootWindow(display *Display) (Window, error) {
	if display.Ptr == nil {
		return Window{}, errors.New("[x11] - Invalid display connection")
	}

	// Get the default screen
	screen := C.XDefaultScreen(display.Ptr)

	// Retrieve the root window for the default screen
	root := C.XRootWindow(display.Ptr, C.int(screen))

	return Window{ID: root}, nil
}

// FreePointer safely frees X11 allocated memory using XFree.
func FreePointer(ptr unsafe.Pointer) {
	if ptr != nil {
		C.XFree(ptr)
	}
}
