//go:build linux

package x11

/*
#cgo LDFLAGS: -lX11
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"github.com/ignite-laboratories/core/std"
	"log"
	"sync"
	"unsafe"
)

//
// --- Types and State ---
//

type Atom = C.Atom

// Display represents a connection to an X11 server.
type Display struct {
	ptr *C.Display
}

// Window represents a single X11 managed window.
type Window struct {
	ID C.Window
}

// Event represents an X11 event.
type Event struct {
	Type int
	Raw  C.XEvent // Raw event data (if needed for advanced use cases)
}

// Wrapper manages the lifecycle of an X11 connection.
type Wrapper struct {
	display Display
	mutex   sync.Mutex
}

var instance *Wrapper
var once sync.Once

//
// --- Singleton Wrapper ---
//

// Initialize initializes or retrieves the singleton Wrapper instance.
func Initialize() (*Wrapper, error) {
	// Use sync.Once to enforce a single connection to the X11 server.
	once.Do(func() {
		ptr := C.XOpenDisplay(nil)
		if ptr == nil {
			log.Fatal("[x11] - Failed to open X11 display")
		}
		instance = &Wrapper{
			display: Display{ptr: ptr},
		}
	})

	if instance == nil {
		return nil, errors.New("[x11] - Error initializing X11 display")
	}
	return instance, nil
}

// Close cleans up the display connection. Should be called at program termination.
func (w *Wrapper) Close() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.display.ptr != nil {
		C.XCloseDisplay(w.display.ptr)
		fmt.Println("[x11] - Closed X11 display connection")
		w.display.ptr = nil
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
func (w *Wrapper) QueryPointer(window Window) (info std.MouseState, err error) {
	// Prepare C variables for receiving data
	var rootReturn, childReturn C.Window
	var rootXReturn, rootYReturn C.int
	var winXReturn, winYReturn C.int
	var maskReturn C.uint

	// Call XQueryPointer
	success := C.XQueryPointer(
		w.display.ptr, // Display pointer
		window.ID,     // Window ID to query
		&rootReturn,   // Pointer to receive parent (root) window
		&childReturn,  // Pointer to receive child window
		&rootXReturn,  // Pointer to root X coordinate
		&rootYReturn,  // Pointer to root Y coordinate
		&winXReturn,   // Pointer to window-relative X coordinate
		&winYReturn,   // Pointer to window-relative Y coordinate
		&maskReturn,   // Pointer to receive button/key mask
	)

	// If XQueryPointer returns 0, it failed
	if success == 0 {
		return std.MouseState{}, errors.New("XQueryPointer failed to query pointer position")
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

	state := std.MouseState{
		GlobalPosition: std.XY[int]{
			X: query.RootX,
			Y: query.RootY,
		},
		WindowPosition: std.XY[int]{
			X: query.WinX,
			Y: query.WinY,
		},
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

	// Convert the retrieved values and return them
	return state, nil
}

//
// --- Basic Display/Screen Handling ---
//

// DefaultScreen retrieves the default screen for the display.
func (w *Wrapper) DefaultScreen() int {
	return int(C.XDefaultScreen(w.display.ptr))
}

// RootWindow retrieves the root window for the default screen.
func (w *Wrapper) RootWindow() Window {
	screen := C.XDefaultScreen(w.display.ptr)
	root := C.XRootWindow(w.display.ptr, C.int(screen))
	return Window{ID: root}
}

//
// --- Window Management ---
//

// CreateWindow creates a basic X11 window.
func (w *Wrapper) CreateWindow(x, y, width, height int) (*Window, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	screen := w.DefaultScreen()
	root := w.RootWindow()

	// Create the window
	window := C.XCreateSimpleWindow(
		w.display.ptr,
		root.ID,
		C.int(x), C.int(y),
		C.uint(width), C.uint(height),
		C.uint(1), // border width
		C.XBlackPixel(w.display.ptr, C.int(screen)),
		C.XWhitePixel(w.display.ptr, C.int(screen)),
	)

	if window == 0 {
		return nil, errors.New("[x11] - Failed to create window")
	}

	// Return the Window object
	return &Window{ID: window}, nil
}

// DestroyWindow destroys a created window.
func (w *Wrapper) DestroyWindow(win *Window) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	C.XDestroyWindow(w.display.ptr, win.ID)
}

// ShowWindow maps (shows) the window on the screen.
func (w *Wrapper) ShowWindow(win *Window) {
	C.XMapWindow(w.display.ptr, win.ID)
}

// SetWindowProtocols sets standard window close behavior (WM_DELETE_WINDOW).
func (w *Wrapper) SetWindowProtocols(win *Window) error {
	atom := w.Atom(WMDeleteWindow) // Use constant from event package
	if atom == 0 {
		return errors.New("[x11] - Failed to retrieve WM_DELETE_WINDOW atom")
	}

	status := C.XSetWMProtocols(w.display.ptr, win.ID, &atom, 1)
	if status == 0 {
		return errors.New("[x11] - Failed to set WM_DELETE_WINDOW protocol")
	}
	return nil
}

//
// --- Event Handling ---
//

// WaitForEvent blocks until an event occurs and returns it.
func (w *Wrapper) WaitForEvent() (*Event, error) {
	var raw C.XEvent
	C.XNextEvent(w.display.ptr, &raw)

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
func (w *Wrapper) GetEventWindow(e *Event) *Window {
	anyEvent := (*C.XAnyEvent)(unsafe.Pointer(&e.Raw))
	return &Window{ID: anyEvent.window}
}

// GetClientMessageData retrieves the message data for a ClientMessage event.
func (w *Wrapper) GetClientMessageData(e *Event) [5]int64 {
	message := (*C.XClientMessageEvent)(unsafe.Pointer(&e.Raw))
	return *(*[5]int64)(unsafe.Pointer(&message.data))
}

//
// --- Atom Management ---
//

// Atom retrieves an existing Atom or creates it if it doesn't exist.
func (w *Wrapper) Atom(name string) C.Atom {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return C.XInternAtom(w.display.ptr, cName, 0)
}

//
// --- Utility Functions ---
//

// SelectInput sets the input mask for the window.
func (w *Wrapper) SelectInput(win *Window, eventMask int64) {
	C.XSelectInput(w.display.ptr, win.ID, C.long(eventMask))
}

// Flush forces the X server to process commands in its buffer.
func (w *Wrapper) Flush() {
	C.XFlush(w.display.ptr)
}
