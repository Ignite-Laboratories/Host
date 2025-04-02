//go:build linux

package graphics

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/host/window"
	"github.com/jezek/xgbutil/xwindow"
	"github.com/remogatto/egl"
	"log"
	"runtime"
)

func init() {
	fmt.Println("[host] - Linux - EGL graphics bridge")
}

// SparkRenderableWindow creates a new GL renderable window using EGL.
func SparkRenderableWindow(renderer Renderable) *RenderableWindow {
	handle := window.Create()
	go sparkEGLBridge(handle, renderer)

	v := &RenderableWindow{}
	v.handle = handle
	return v
}

func sparkEGLBridge(window *xwindow.Window, renderer Renderable) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Initialize EGL
	display := egl.GetDisplay(egl.DEFAULT_DISPLAY)
	if display == egl.NO_DISPLAY {
		log.Fatal("Failed to get EGL display")
	}
	if !egl.Initialize(display, nil, nil) {
		log.Fatal("Failed to initialize EGL")
	}
	defer egl.Terminate(display)

	// Choose an EGL configuration
	attribs := []int32{
		egl.RED_SIZE, 8, // 8-bit Red channel
		egl.GREEN_SIZE, 8, // 8-bit Green channel
		egl.BLUE_SIZE, 8, // 8-bit Blue channel
		egl.DEPTH_SIZE, 24, // 24-bit Depth buffer
		egl.STENCIL_SIZE, 8, // 8-bit Stencil buffer
		egl.SURFACE_TYPE, egl.WINDOW_BIT, // Rendering surface type: window
		egl.RENDERABLE_TYPE, egl.OPENGL_BIT, // Enable OpenGL rendering
		egl.NONE, // End of attribute list

	}
	var config egl.Config
	var numConfigs int32
	if !egl.ChooseConfig(display, attribs, &config, 1, &numConfigs) || numConfigs == 0 {
		log.Fatal("Failed to choose EGL configuration")
	}

	// Create EGL surface for the X11 window
	surface := egl.CreateWindowSurface(display, config, egl.NativeWindowType(uintptr(window.Id)), nil)
	if surface == egl.NO_SURFACE {
		log.Fatalf("Failed to create EGL surface: %v", egl.GetError())
	}
	defer egl.DestroySurface(display, surface)

	// Create an EGL OpenGL context
	context := egl.CreateContext(display, config, egl.NO_CONTEXT, nil)
	if context == egl.NO_CONTEXT {
		log.Fatalf("Failed to create EGL context: %v", egl.GetError())
	}
	defer egl.DestroyContext(display, context)

	// Initialize EGL and OpenGL for this window
	if !egl.MakeCurrent(display, surface, surface, context) {
		log.Fatalf("Failed to make EGL context current: %v", egl.GetError())
	}

	// Initialize OpenGL
	if err := gl.Init(); err != nil {
		log.Fatalf("Failed to initialize OpenGL: %v", err)
	}

	// Start rendering
	for core.Alive {
		renderer.Render()
		// Swap the buffers to display the rendered frame
		egl.SwapBuffers(display, surface)
	}
}
