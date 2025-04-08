//go:build linux

package egl

/*
#cgo LDFLAGS: -lEGL -lGL
#include <EGL/egl.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// EGL Types
type Boolean uintptr
type Display uintptr
type Context uintptr
type Surface uintptr
type Config uintptr
type NativeWindowType uintptr
type NativeDisplayType uintptr

// **Functions Wrapper Section**

// GetDisplay retrieves the EGL display connection for the specified display ID.
func GetDisplay(displayID NativeDisplayType) Display {
	return Display(C.eglGetDisplay(C.EGLNativeDisplayType(displayID)))
}

func QueryString(
	disp Display, name int32) string {
	return C.GoString(C.eglQueryString(
		C.EGLDisplay(unsafe.Pointer(disp)),
		C.EGLint(name)))
}

// Initialize initializes the EGL display connection and retrieves major/minor version numbers.
func Initialize(display Display, major, minor *int32) bool {
	result := C.eglInitialize(C.EGLDisplay(display), (*C.EGLint)(major), (*C.EGLint)(minor))
	return result == C.EGL_TRUE
}

// ChooseConfig selects an EGL frame buffer configuration that matches the desired attributes.
func ChooseConfig(display Display, attribs []int32, configs []Config, configSize int32, numConfig *int32) bool {
	result := C.eglChooseConfig(
		C.EGLDisplay(display),
		(*C.EGLint)(&attribs[0]),
		(*C.EGLConfig)(&configs[0]),
		C.EGLint(configSize),
		(*C.EGLint)(numConfig),
	)
	return result == C.EGL_TRUE
}

func GetConfigAttrib(
	disp Display, config Config,
	attribute int32, value *int32) bool {
	return C.eglGetConfigAttrib(
		C.EGLDisplay(unsafe.Pointer(disp)),
		C.EGLConfig(config),
		C.EGLint(attribute),
		(*C.EGLint)(unsafe.Pointer(value))) == C.EGL_TRUE
}

// CreatePbufferSurface creates an EGL Pixel Buffer Surface.
func CreatePbufferSurface(display Display, config Config, attribs []int32) Surface {
	var attribsPtr *C.EGLint
	if attribs != nil {
		attribsPtr = (*C.EGLint)(&attribs[0])
	}
	return Surface(C.eglCreatePbufferSurface(C.EGLDisplay(display), C.EGLConfig(config), attribsPtr))
}

// CreateContext creates a new EGL rendering context for the specified display and configuration.
func CreateContext(display Display, config Config, shareContext Context, attribs []int32) Context {
	var attribsPtr *C.EGLint
	if attribs != nil {
		attribsPtr = (*C.EGLint)(&attribs[0])
	}
	return Context(C.eglCreateContext(C.EGLDisplay(display), C.EGLConfig(config), C.EGLContext(shareContext), attribsPtr))
}

// MakeCurrent makes an EGL context and surface the current rendering target.
func MakeCurrent(display Display, draw, read Surface, context Context) bool {
	result := C.eglMakeCurrent(
		C.EGLDisplay(display),
		C.EGLSurface(draw),
		C.EGLSurface(read),
		C.EGLContext(context),
	)
	return result == C.EGL_TRUE
}

// SwapBuffers swaps the front and back buffers for the given surface.
func SwapBuffers(display Display, surface Surface) bool {
	result := C.eglSwapBuffers(C.EGLDisplay(display), C.EGLSurface(surface))
	return result == C.EGL_TRUE
}

// DestroyContext destroys the specified EGL context.
func DestroyContext(display Display, context Context) bool {
	result := C.eglDestroyContext(C.EGLDisplay(display), C.EGLContext(context))
	return result == C.EGL_TRUE
}

// Terminate terminates the connection to the display.
func Terminate(display Display) bool {
	result := C.eglTerminate(C.EGLDisplay(display))
	return result == C.EGL_TRUE
}

// SwapInterval sets the swap interval for the EGL surface.
func SwapInterval(display Display, interval int) bool {
	result := C.eglSwapInterval(C.EGLDisplay(display), C.EGLint(interval))
	return result == C.EGL_TRUE
}

// DestroySurface destroys the specified EGL surface and releases associated resources.
func DestroySurface(display Display, surface Surface) bool {
	result := C.eglDestroySurface(C.EGLDisplay(display), C.EGLSurface(surface))
	return result == C.EGL_TRUE
}

// CreateWindowSurface creates an EGL window surface tied to a native window.
func CreateWindowSurface(display Display, config Config, nativeWindow NativeWindowType, attribs []int32) Surface {
	var attribsPtr *C.EGLint
	if attribs != nil {
		attribsPtr = (*C.EGLint)(&attribs[0])
	}
	return Surface(C.eglCreateWindowSurface(
		C.EGLDisplay(display),
		C.EGLConfig(config),
		C.EGLNativeWindowType(nativeWindow),
		attribsPtr,
	))
}

// GetEGLError retrieves the human-readable string for the last EGL error.
func GetEGLError() string {
	errorCode := C.eglGetError()
	switch errorCode {
	case C.EGL_SUCCESS:
		return "EGL_SUCCESS"
	case C.EGL_NOT_INITIALIZED:
		return "EGL_NOT_INITIALIZED"
	case C.EGL_BAD_ACCESS:
		return "EGL_BAD_ACCESS"
	case C.EGL_BAD_ALLOC:
		return "EGL_BAD_ALLOC"
	case C.EGL_BAD_ATTRIBUTE:
		return "EGL_BAD_ATTRIBUTE"
	case C.EGL_BAD_CONTEXT:
		return "EGL_BAD_CONTEXT"
	case C.EGL_BAD_CONFIG:
		return "EGL_BAD_CONFIG"
	case C.EGL_BAD_CURRENT_SURFACE:
		return "EGL_BAD_CURRENT_SURFACE"
	case C.EGL_BAD_DISPLAY:
		return "EGL_BAD_DISPLAY"
	case C.EGL_BAD_SURFACE:
		return "EGL_BAD_SURFACE"
	case C.EGL_BAD_MATCH:
		return "EGL_BAD_MATCH"
	case C.EGL_BAD_PARAMETER:
		return "EGL_BAD_PARAMETER"
	case C.EGL_BAD_NATIVE_PIXMAP:
		return "EGL_BAD_NATIVE_PIXMAP"
	case C.EGL_BAD_NATIVE_WINDOW:
		return "EGL_BAD_NATIVE_WINDOW"
	case C.EGL_CONTEXT_LOST:
		return "EGL_CONTEXT_LOST"
	default:
		return fmt.Sprintf("Unknown EGL error: %x", errorCode)
	}
}
