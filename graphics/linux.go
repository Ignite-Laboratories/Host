//go:build linux

package graphics

import "C"
import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/ignite-laboratories/host/hydraold"
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
	w.Handle = hydraold.Create(size)

	go sparkEGLBridge(w.Handle, renderer)
	return w
}

func sparkEGLBridge(handle *hydraold.Handle, renderer Renderable) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	//err := hydraold.DisableVSync(handle.Display, handle.Window)
	//if err != nil {
	//	log.Printf("Failed to disable VSync: %v", err)
	//}

	screen := hydraold.DefaultScreen(handle.Display)

	// Define GLX attributes
	visualAttribs := []int32{
		hydraold.GLX_X_RENDERABLE, 1,
		hydraold.GLX_RENDER_TYPE, hydraold.GLX_RGBA_BIT,
		hydraold.GLX_DRAWABLE_TYPE, hydraold.GLX_WINDOW_BIT,
		hydraold.GLX_X_VISUAL_TYPE, hydraold.GLX_TRUE_COLOR,
		hydraold.GLX_RED_SIZE, 8,
		hydraold.GLX_GREEN_SIZE, 8,
		hydraold.GLX_BLUE_SIZE, 8,
		hydraold.GLX_DEPTH_SIZE, 24,
		hydraold.GLX_DOUBLEBUFFER, 1,
		0, // Null-terminate
	}

	// Choose framebuffer config
	fbConfigs, err := hydraold.ChooseFramebufferConfig(handle.Display, screen, visualAttribs)
	if err != nil {
		log.Fatalf("Failed to choose framebuffer config: %v", err)
	}
	fbConfig := fbConfigs[0]

	// Get visual information
	visualInfo, err := hydraold.GetVisualFromFBConfig(handle.Display, fbConfig)
	if err != nil {
		log.Fatalf("Failed to get visual info: %v", err)
	}
	defer hydraold.FreePointer(unsafe.Pointer(visualInfo))

	// Create GLX context
	context, err := hydraold.CreateGLXContext(handle.Display, fbConfig, nil, true, hydraold.GLXContextAttributes{
		MajorVersion: 3,
		MinorVersion: 1,
		ProfileMask:  hydraold.GLX_CONTEXT_ES2_PROFILE_BIT_EXT,
	})
	if err != nil {
		log.Fatalf("Failed to create GLX context: %v", err)
	}

	// Make the context current
	err = hydraold.MakeGLXContextCurrent(handle.Display, handle.Window, context)
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
		hydraold.SwapBuffers(handle.Display, handle.Window)
	}
}
