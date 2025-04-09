//go:build linux

package graphics

/*
#cgo LDFLAGS: -lGL -lX11
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <GL/gl.h>
#include <GL/glx.h>
#include <stdlib.h>

int Test() {
	return 1;
}

// GLX Context Extensions (for OpenGL versions > 2.1)
typedef struct {
    int contextMajorVersion;
    int contextMinorVersion;
    int contextFlags;
    int profileMask;
} GLXContextAttributes;

PFNGLXCREATECONTEXTATTRIBSARBPROC glXCreateContextAttribsARB = 0;

GLXContext createGLXContext(Display *display, GLXFBConfig config, GLXContext shareList, Bool direct, GLXContextAttributes attribs) {
    int attribList[] = {
        0x2091, attribs.contextMajorVersion, // GLX_CONTEXT_MAJOR_VERSION_ARB
        0x2092, attribs.contextMinorVersion, // GLX_CONTEXT_MINOR_VERSION_ARB
        0x9126, attribs.profileMask,         // GLX_CONTEXT_PROFILE_MASK_ARB
        0
    };

    if (!glXCreateContextAttribsARB) {
        glXCreateContextAttribsARB = (PFNGLXCREATECONTEXTATTRIBSARBPROC)glXGetProcAddressARB((const GLubyte *) "glXCreateContextAttribsARB");
    }

    return glXCreateContextAttribsARB(display, config, shareList, direct, attribList);
}

*/
import "C"
import "C"
import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"log"
	"runtime"
	"strings"
	"unsafe"

	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
	"github.com/ignite-laboratories/host/window"
)

// SparkRenderableWindow creates a new GL renderable window using EGL.
func SparkRenderableWindow(size std.XY[int], renderer Renderable) *RenderableWindow {
	w := &RenderableWindow{}
	w.Handle = window.Create(size)

	go sparkEGLBridge(w.Handle, renderer)
	return w
}

func sparkEGLBridge(handle *window.Handle, renderer Renderable) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Get the default screen
	screen := C.XDefaultScreen((*C.Display)(handle.Display.Ptr))

	// Define GLX attributes
	visualAttribs := []C.int{
		C.GLX_X_RENDERABLE, 1, // Ensure renderable
		C.GLX_RENDER_TYPE, C.GLX_RGBA_BIT,
		C.GLX_DRAWABLE_TYPE, C.GLX_WINDOW_BIT,
		C.GLX_X_VISUAL_TYPE, C.GLX_TRUE_COLOR,
		C.GLX_RED_SIZE, 8,
		C.GLX_GREEN_SIZE, 8,
		C.GLX_BLUE_SIZE, 8,
		C.GLX_DEPTH_SIZE, 24,
		C.GLX_DOUBLEBUFFER, 1,
		0, // Null-terminate
	}

	// Retrieve framebuffer configs
	var fbCount C.int
	fbConfigs := C.glXChooseFBConfig((*C.Display)(handle.Display.Ptr), screen, &visualAttribs[0], &fbCount)
	if fbConfigs == nil || fbCount == 0 {
		log.Fatal("Failed to retrieve framebuffer config")
	}

	// Cast the pointer to an array and access the first framebuffer config
	fbConfig := (*[1 << 28]C.GLXFBConfig)(unsafe.Pointer(fbConfigs))[:fbCount:fbCount][0]

	// Get a visual from the framebuffer config
	visualInfo := C.glXGetVisualFromFBConfig((*C.Display)(handle.Display.Ptr), fbConfig)
	if visualInfo == nil {
		log.Fatal("Failed to get visual info")
	}
	defer C.XFree(unsafe.Pointer(visualInfo))

	// Create a GLX context for OpenGL ES 3.2
	contextAttribs := C.GLXContextAttributes{
		contextMajorVersion: 3,   // Request OpenGL ES major version 3
		contextMinorVersion: 1,   // Request OpenGL ES minor version 2
		profileMask:         0x4, // GLX_CONTEXT_ES2_PROFILE_BIT_EXT for OpenGL ES
	}
	glxContext := C.createGLXContext((*C.Display)(handle.Display.Ptr), fbConfig, nil, C.True, contextAttribs)
	if glxContext == nil {
		log.Fatal("Failed to create OpenGL 3.3 Core context")
	}

	// Make the context current
	if ok := C.glXMakeCurrent((*C.Display)(handle.Display.Ptr), C.GLXDrawable(handle.Window.ID), glxContext); ok == 0 {
		log.Fatal("Failed to make OpenGL context current")
	}

	// Initialize GL using Go-GL
	if err := gl.Init(); err != nil {
		log.Fatalf("Failed to initialize OpenGL: %v", err)
	}

	// Initialize the rendering context
	renderer.Initialize()

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

	// Start the rendering loop
	for !handle.Destroyed && core.Alive {
		renderer.Render()

		C.glXSwapBuffers((*C.Display)(handle.Display.Ptr), C.GLXDrawable(handle.Window.ID))
	}
}
