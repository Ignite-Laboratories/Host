//go:build linux

package graphics

import "C"
import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/ignite-laboratories/host/hydra"
	"log"
	"runtime"
	"strings"
	"unsafe"

	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
)

// SparkRenderableWindow creates a new GL renderable window using EGL.
func SparkRenderableWindow(size std.XY[int], renderer Renderable) *RenderableWindow {
	w := &RenderableWindow{}
	w.Handle = hydra.Create(size)

	go sparkEGLBridge(w.Handle, renderer)
	return w
}

func sparkEGLBridge(handle *hydra.Handle, renderer Renderable) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	//err := hydra.DisableVSync(handle.Display, handle.Window)
	//if err != nil {
	//	log.Printf("Failed to disable VSync: %v", err)
	//}

	screen := hydra.DefaultScreen(handle.Display)

	// Define GLX attributes
	visualAttribs := []int32{
		hydra.GLX_X_RENDERABLE, 1,
		hydra.GLX_RENDER_TYPE, hydra.GLX_RGBA_BIT,
		hydra.GLX_DRAWABLE_TYPE, hydra.GLX_WINDOW_BIT,
		hydra.GLX_X_VISUAL_TYPE, hydra.GLX_TRUE_COLOR,
		hydra.GLX_RED_SIZE, 8,
		hydra.GLX_GREEN_SIZE, 8,
		hydra.GLX_BLUE_SIZE, 8,
		hydra.GLX_DEPTH_SIZE, 24,
		hydra.GLX_DOUBLEBUFFER, 1,
		0, // Null-terminate
	}

	// Choose framebuffer config
	fbConfigs, err := hydra.ChooseFramebufferConfig(handle.Display, screen, visualAttribs)
	if err != nil {
		log.Fatalf("Failed to choose framebuffer config: %v", err)
	}
	fbConfig := fbConfigs[0]

	// Get visual information
	visualInfo, err := hydra.GetVisualFromFBConfig(handle.Display, fbConfig)
	if err != nil {
		log.Fatalf("Failed to get visual info: %v", err)
	}
	defer hydra.FreePointer(unsafe.Pointer(visualInfo))

	// Create GLX context
	context, err := hydra.CreateGLXContext(handle.Display, fbConfig, nil, true, hydra.GLXContextAttributes{
		MajorVersion: 3,
		MinorVersion: 3,
		ProfileMask:  hydra.GLX_CONTEXT_CORE_PROFILE_BIT_ARB,
	})
	if err != nil {
		log.Fatalf("Failed to create GLX context: %v", err)
	}

	// Make the context current
	err = hydra.MakeGLXContextCurrent(handle.Display, handle.Window, context)
	if err != nil {
		log.Fatalf("Failed to make context current: %v", err)
	}

	// Initialize OpenGL using Go-GL
	if err := gl.Init(); err != nil {
		log.Fatalf("Failed to initialize OpenGL: %v", err)
	}

	ver := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println(ver)

	numExtensions := int32(0)
	gl.GetIntegerv(gl.NUM_EXTENSIONS, &numExtensions)

	for i := int32(0); i < numExtensions; i++ {
		extension := gl.GoStr(gl.GetStringi(gl.EXTENSIONS, uint32(i)))
		if strings.Contains(extension, "geometry") {
			fmt.Println(extension)
		}
	}

	// Initialize the renderer
	renderer.Initialize()

	fmt.Println(gl.GoStr(gl.GetString(gl.VERSION)))

	// Main render loop
	for !handle.Destroyed && core.Alive {
		renderer.Render()
		hydra.SwapBuffers(handle.Display, handle.Window)
	}
}
