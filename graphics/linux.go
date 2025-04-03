//go:build linux

package graphics

/*
#cgo LDFLAGS: -lEGL -lGL
#include <EGL/egl.h>
#include <GL/gl.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
	"github.com/ignite-laboratories/host/window"
	"log"
	"runtime"
)

func init() {
	fmt.Println("[host] - Linux - sparking native EGL graphics bridge")
}

// SparkRenderableWindow creates a new GL renderable window using EGL.
func SparkRenderableWindow(size std.XY[int], renderer Renderable) *RenderableWindow {
	w := &RenderableWindow{}
	w.Handle = window.Create(size)
	go sparkEGLBridge(w.Handle, renderer)
	return w
}

func sparkEGLBridge(handle window.Handle, renderer Renderable) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Initialize EGL
	display := C.eglGetDisplay(C.EGL_DEFAULT_DISPLAY)
	if display == nil {
		log.Fatalf("Failed to get EGL display: %v", getEGLError())
	}
	if C.eglInitialize(display, nil, nil) == C.EGL_FALSE {
		log.Fatalf("Failed to initialize EGL: %v", getEGLError())
	}
	defer C.eglTerminate(display)

	// Choose an EGL configuration
	attribs := []C.EGLint{
		C.EGL_RED_SIZE, 8,
		C.EGL_GREEN_SIZE, 8,
		C.EGL_BLUE_SIZE, 8,
		C.EGL_DEPTH_SIZE, 24,
		C.EGL_STENCIL_SIZE, 8,
		C.EGL_SURFACE_TYPE, C.EGL_WINDOW_BIT,
		C.EGL_RENDERABLE_TYPE, C.EGL_OPENGL_BIT,
		C.EGL_NONE,
	}
	var config C.EGLConfig
	var numConfigs C.EGLint
	if C.eglChooseConfig(display, &attribs[0], &config, 1, &numConfigs) == C.EGL_FALSE || numConfigs == 0 {
		log.Fatalf("Failed to choose EGL configuration: %v", getEGLError())
	}

	// Create an EGL surface for the X11 window
	surface := C.eglCreateWindowSurface(display, config, C.EGLNativeWindowType(uintptr(handle.Window.ID)), nil)
	if surface == nil {
		log.Fatalf("Failed to create EGL surface: %v", getEGLError())
	}
	defer C.eglDestroySurface(display, surface)

	// Create an EGL context
	context := C.eglCreateContext(display, config, nil, nil)
	if context == nil {
		log.Fatalf("Failed to create EGL context: %v", getEGLError())
	}
	defer C.eglDestroyContext(display, context)

	// Make the current context
	if C.eglMakeCurrent(display, surface, surface, context) == C.EGL_FALSE {
		log.Fatalf("Failed to make EGL context current: %v", getEGLError())
	}

	// Initialize OpenGL
	if err := gl.Init(); err != nil {
		log.Fatalf("Failed to initialize OpenGL: %v", err)
	}

	// Start rendering
	for core.Alive {
		renderer.Render()

		// Swap the buffers to display the rendered frame
		C.eglSwapBuffers(display, surface)
	}
}

// getEGLError retrieves the EGL error string.
func getEGLError() string {
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
