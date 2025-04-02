//go:build linux

package graphics

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/ignite-laboratories/core"
	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgbutil/xwindow"
	"github.com/remogatto/egl"
	"log"
	"runtime"
)

func init() {
	fmt.Println("[host] - Linux - egl graphics bridge")
}

func Setup(window *xwindow.Window, loop func(display egl.Display, surface egl.Surface)) {
	go func() {
		runtime.LockOSThread()

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
		config := chooseEGLConfig(display)

		// Initialize EGL and OpenGL for this window
		context, surface := initializeEGLSurface(display, config, window.Id)
		defer egl.DestroySurface(display, surface)
		defer egl.DestroyContext(display, context)

		if !egl.MakeCurrent(display, surface, surface, context) {
			log.Fatalf("Failed to make EGL context current: %v", egl.GetError())
		}

		// Initialize OpenGL
		if err := gl.Init(); err != nil {
			log.Fatalf("Failed to initialize OpenGL: %v", err)
		}
		gl.Viewport(0, 0, int32(500), int32(500))

		// Start rendering
		for core.Alive {
			loop(display, surface)
		}
	}()
}

// initializeEGLSurface sets up an EGL surface and context tied to the X11 window
func initializeEGLSurface(display egl.Display, config egl.Config, window xproto.Window) (egl.Context, egl.Surface) {
	// Create EGL surface for the X11 window
	surface := egl.CreateWindowSurface(display, config, egl.NativeWindowType(uintptr(window)), nil)
	if surface == egl.NO_SURFACE {
		log.Fatalf("Failed to create EGL surface: %v", egl.GetError())
	}

	// Create an EGL OpenGL context
	context := egl.CreateContext(display, config, egl.NO_CONTEXT, nil)
	if context == egl.NO_CONTEXT {
		log.Fatalf("Failed to create EGL context: %v", egl.GetError())
	}

	return context, surface
}

// chooseEGLConfig selects an appropriate configuration for rendering
func chooseEGLConfig(display egl.Display) egl.Config {
	attribs := []int32{
		egl.RED_SIZE, 8,
		egl.GREEN_SIZE, 8,
		egl.BLUE_SIZE, 8,
		egl.DEPTH_SIZE, 24,
		egl.STENCIL_SIZE, 8,
		egl.SURFACE_TYPE, egl.WINDOW_BIT,
		egl.RENDERABLE_TYPE, egl.OPENGL_BIT,
		egl.NONE,
	}

	var config egl.Config
	var numConfigs int32
	if !egl.ChooseConfig(display, attribs, &config, 1, &numConfigs) || numConfigs == 0 {
		log.Fatal("Failed to choose EGL configuration")
	}

	return config
}
