//go:build linux

package graphics

import "C"
import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
	"github.com/ignite-laboratories/host/egl"
	"github.com/ignite-laboratories/host/window"
)

func init() {
	fmt.Println("[host] - Linux - sparking EGL bridge")

	// Get the default display
	display = egl.GetDisplay(egl.DEFAULT_DISPLAY)
	if display == egl.Display(egl.NO_DISPLAY) {
		log.Fatalf("Failed to get EGL display: %v", egl.GetEGLError())
	}

	// Initialize EGL
	var major, minor int32
	if !egl.Initialize(display, &major, &minor) {
		log.Fatalf("Failed to initialize EGL: %v", egl.GetEGLError())
	}
	fmt.Printf("EGL initialized. Version %d.%d\n", major, minor)

	// Choose an EGL configuration
	attribs := []int32{
		egl.RED_SIZE, 8,
		egl.GREEN_SIZE, 8,
		egl.BLUE_SIZE, 8,
		egl.ALPHA_SIZE, 8,
		egl.DEPTH_SIZE, 8,
		egl.RENDERABLE_TYPE, egl.OPENGL_BIT,
		egl.NONE,
	}
	configs := make([]egl.Config, 1)
	var numConfigs int32
	if !egl.ChooseConfig(display, attribs, configs, 1, &numConfigs) {
		// The egl.ChooseConfig call failed (e.g., a library issue or invalid attribs)
		log.Fatalf("Failed to choose EGL configuration: %v", egl.GetEGLError())
	}

	// Check if no matching configuration was found
	if numConfigs == 0 {
		log.Fatal("No suitable EGL configuration found")
	}
	config = configs[0]

	var value int32
	egl.GetConfigAttrib(display, config, egl.RENDERABLE_TYPE, &value)
	fmt.Printf("Chosen config RENDERABLE_TYPE: 0x%x\n", value)

	// Cleanup EGL connection when the process quits
	go func() {
		core.WhileAlive()
		egl.Terminate(display)
		fmt.Println("[host] - Linux - closed EGL bridge")
	}()
}

// SparkRenderableWindow creates a new GL renderable window using EGL.
func SparkRenderableWindow(size std.XY[int], renderer Renderable) *RenderableWindow {
	w := &RenderableWindow{}
	w.Handle = window.Create(size)

	go sparkEGLBridge(w.Handle, renderer)
	return w
}

var config egl.Config
var display egl.Display

func sparkEGLBridge(handle *window.Handle, renderer Renderable) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Create an EGL surface for the provided window
	surface := egl.CreateWindowSurface(display, config, egl.NativeWindowType(uintptr(handle.Window.ID)), nil)
	if surface == egl.Surface(egl.NO_SURFACE) {
		log.Fatalf("Failed to create EGL surface: %v", egl.GetEGLError())
	}
	defer egl.DestroySurface(display, surface)

	// Set context attributes for OpenGL 3.3 core
	//contextAttributes := []int32{
	//	egl.CONTEXT_MAJOR_VERSION, 3,
	//	egl.CONTEXT_MINOR_VERSION, 3,
	//	egl.CONTEXT_OPENGL_PROFILE_MASK, egl.CONTEXT_OPENGL_CORE_PROFILE_BIT,
	//	int32(egl.NONE), // Terminator
	//}

	contextAttributes := []int32{
		egl.CONTEXT_MAJOR_VERSION, 3,
		egl.CONTEXT_MINOR_VERSION, 1,
		int32(egl.NONE), // Terminator
	}

	// Create an EGL context
	context := egl.CreateContext(display, config, egl.Context(egl.NO_CONTEXT), contextAttributes)
	if context == egl.Context(egl.NO_CONTEXT) {
		log.Fatalf("Failed to create EGL context: %v", egl.GetEGLError())
	}
	defer egl.DestroyContext(display, context)

	// Make the EGL context current
	if !egl.MakeCurrent(display, surface, surface, context) {
		log.Fatalf("Failed to make EGL context current: %v", egl.GetEGLError())
	}

	// Set swap interval
	if !egl.SwapInterval(display, 0) {
		log.Printf("Failed to set swap interval (non-fatal error): %v", egl.GetEGLError())
	}

	// Initialize the rendering context
	renderer.Initialize()

	// Initialize GL using Go-GL
	if err := gl.Init(); err != nil {
		log.Fatalf("Failed to initialize OpenGL: %v", err)
	}

	// Start the rendering loop
	for !handle.Destroyed && core.Alive {
		renderer.Render()

		// Swap the buffers to display the rendered frame
		if !egl.SwapBuffers(display, surface) {
			log.Printf("Failed to swap buffers: %v", egl.GetEGLError())
		}
	}
}
