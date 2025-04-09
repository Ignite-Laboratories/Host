package hydra

/*
#cgo LDFLAGS: -lX11
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <GL/glx.h>
#include <stdlib.h>
*/
import "C"

// --- Event Masks ---
const (
	NoEventMask              = C.NoEventMask
	KeyPressMask             = C.KeyPressMask
	KeyReleaseMask           = C.KeyReleaseMask
	ButtonPressMask          = C.ButtonPressMask
	ButtonReleaseMask        = C.ButtonReleaseMask
	EnterWindowMask          = C.EnterWindowMask
	LeaveWindowMask          = C.LeaveWindowMask
	PointerMotionMask        = C.PointerMotionMask
	ButtonMotionMask         = C.ButtonMotionMask
	KeymapStateMask          = C.KeymapStateMask
	ExposureMask             = C.ExposureMask
	VisibilityChangeMask     = C.VisibilityChangeMask
	StructureNotifyMask      = C.StructureNotifyMask
	ResizeRedirectMask       = C.ResizeRedirectMask
	SubstructureNotifyMask   = C.SubstructureNotifyMask
	SubstructureRedirectMask = C.SubstructureRedirectMask
	FocusChangeMask          = C.FocusChangeMask
	PropertyChangeMask       = C.PropertyChangeMask
	ColormapChangeMask       = C.ColormapChangeMask
	OwnerGrabButtonMask      = C.OwnerGrabButtonMask
)

// --- Event Types ---
const (
	KeyPress         = C.KeyPress
	KeyRelease       = C.KeyRelease
	ButtonPress      = C.ButtonPress
	ButtonRelease    = C.ButtonRelease
	MotionNotify     = C.MotionNotify
	EnterNotify      = C.EnterNotify
	LeaveNotify      = C.LeaveNotify
	FocusIn          = C.FocusIn
	FocusOut         = C.FocusOut
	KeymapNotify     = C.KeymapNotify
	Expose           = C.Expose
	GraphicsExpose   = C.GraphicsExpose
	NoExpose         = C.NoExpose
	VisibilityNotify = C.VisibilityNotify
	CreateNotify     = C.CreateNotify
	DestroyNotify    = C.DestroyNotify
	UnmapNotify      = C.UnmapNotify
	MapNotify        = C.MapNotify
	MapRequest       = C.MapRequest
	ReparentNotify   = C.ReparentNotify
	ConfigureNotify  = C.ConfigureNotify
	ConfigureRequest = C.ConfigureRequest
	GravityNotify    = C.GravityNotify
	ResizeRequest    = C.ResizeRequest
	CirculateNotify  = C.CirculateNotify
	CirculateRequest = C.CirculateRequest
	PropertyNotify   = C.PropertyNotify
	SelectionClear   = C.SelectionClear
	SelectionRequest = C.SelectionRequest
	SelectionNotify  = C.SelectionNotify
	ColormapNotify   = C.ColormapNotify
	ClientMessage    = C.ClientMessage
	MappingNotify    = C.MappingNotify
)

// --- Color Constants ---

// BlackPixel retrieves the black pixel value for a given display and screen.
func BlackPixel(display *C.Display, screen int) C.ulong {
	return C.XBlackPixel(display, C.int(screen))
}

// WhitePixel retrieves the white pixel value for a given display and screen.
func WhitePixel(display *C.Display, screen int) C.ulong {
	return C.XWhitePixel(display, C.int(screen))
}

// --- Common Atoms ---
const (
	XAString       = C.XA_STRING
	WMDeleteWindow = "WM_DELETE_WINDOW" // Used for window close events
	WMProtocols    = "WM_PROTOCOLS"     // For setting window manager protocols
	NetWMName      = "_NET_WM_NAME"     // Used for extended window name properties
	NetWMIconName  = "_NET_WM_ICON_NAME"
)

// --- Mouse State ---
const (
	Button1Mask = 1 << 8  // Left mouse button
	Button2Mask = 1 << 9  // Middle mouse button
	Button3Mask = 1 << 10 // Right mouse button
	Button4Mask = 1 << 11 // Scroll wheel up
	Button5Mask = 1 << 12 // Scroll wheel down
)

// --- GLX ---

const (
	GLX_CONTEXT_MAJOR_VERSION_ARB             = 0x2091
	GLX_CONTEXT_MINOR_VERSION_ARB             = 0x2092
	GLX_CONTEXT_PROFILE_MASK_ARB              = 0x9126
	GLX_CONTEXT_ES2_PROFILE_BIT_EXT           = 0x0004 // This requests OpenGL ES profiles
	GLX_CONTEXT_CORE_PROFILE_BIT_ARB          = 0x1
	GLX_CONTEXT_COMPATIBILITY_PROFILE_BIT_ARB = 0x2
)

const (
	GLX_X_RENDERABLE  = int32(C.GLX_X_RENDERABLE)
	GLX_RENDER_TYPE   = int32(C.GLX_RENDER_TYPE)
	GLX_RGBA_BIT      = int32(C.GLX_RGBA_BIT)
	GLX_DRAWABLE_TYPE = int32(C.GLX_DRAWABLE_TYPE)
	GLX_WINDOW_BIT    = int32(C.GLX_WINDOW_BIT)
	GLX_X_VISUAL_TYPE = int32(C.GLX_X_VISUAL_TYPE)
	GLX_TRUE_COLOR    = int32(C.GLX_TRUE_COLOR)
	GLX_RED_SIZE      = int32(C.GLX_RED_SIZE)
	GLX_GREEN_SIZE    = int32(C.GLX_GREEN_SIZE)
	GLX_BLUE_SIZE     = int32(C.GLX_BLUE_SIZE)
	GLX_DEPTH_SIZE    = int32(C.GLX_DEPTH_SIZE)
	GLX_DOUBLEBUFFER  = int32(C.GLX_DOUBLEBUFFER)
)
