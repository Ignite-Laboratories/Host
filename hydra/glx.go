//go:build linux

package hydra

import "C"
import (
	"errors"
	"strings"
	"unsafe"
)

/*
#cgo LDFLAGS: -lX11 -lGL
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <GL/gl.h>
#include <GL/glx.h>
#include <stdlib.h>

// GLX Context Extensions (for OpenGL versions > 2.1)
typedef struct {
    int contextMajorVersion;
    int contextMinorVersion;
    int contextFlags;
    int profileMask;
} GLXContextAttributes;

PFNGLXCREATECONTEXTATTRIBSARBPROC glXCreateContextAttribsARB = 0;

// Creates a GLX context with specified attributes
GLXContext createGLXContext(Display *display, GLXFBConfig config, GLXContext shareList, Bool direct, GLXContextAttributes attribs) {
    int attribList[] = {
        0x2091, attribs.contextMajorVersion, // GLX_CONTEXT_MAJOR_VERSION_ARB
        0x2092, attribs.contextMinorVersion, // GLX_CONTEXT_MINOR_VERSION_ARB
        0x9126, attribs.profileMask,         // GLX_CONTEXT_PROFILE_MASK_ARB
        0 // Null-terminate the list
    };

    if (!glXCreateContextAttribsARB) {
        glXCreateContextAttribsARB = (PFNGLXCREATECONTEXTATTRIBSARBPROC)glXGetProcAddressARB((const GLubyte *)"glXCreateContextAttribsARB");
    }

    return glXCreateContextAttribsARB(display, config, shareList, direct, attribList);
}


// GLX function pointers for swap interval extensions
void (*glXSwapIntervalEXT)(Display *, GLXDrawable, int) = NULL;
int (*glXSwapIntervalMESA)(int) = NULL;
int (*glXSwapIntervalSGI)(int) = NULL;
*/
import "C"

//
// --- Types and Constants ---
//

type GLXFBConfig = C.GLXFBConfig
type GLXContext = C.GLXContext
type GLXDrawable = C.GLXDrawable
type GLXContextAttributes struct {
	MajorVersion int
	MinorVersion int
	ProfileMask  int
}

//
// --- GLX Wrappers ---
//

// ChooseFramebufferConfig retrieves a list of framebuffer configurations.
func ChooseFramebufferConfig(display *Display, screen int, visualAttribs []int32) ([]GLXFBConfig, error) {
	cAttributes := (*C.int)(unsafe.Pointer(&visualAttribs[0]))
	var fbCount C.int
	fbConfigs := C.glXChooseFBConfig(display.Ptr, C.int(screen), cAttributes, &fbCount)
	if fbConfigs == nil {
		return nil, errors.New("glXChooseFBConfig failed to retrieve framebuffer configurations")
	}

	// Convert C results to Go slice
	goConfigs := (*[1 << 28]C.GLXFBConfig)(unsafe.Pointer(fbConfigs))[:fbCount:fbCount]
	return goConfigs, nil
}

// GetVisualFromFBConfig retrieves the visual information for a framebuffer configuration.
func GetVisualFromFBConfig(display *Display, config GLXFBConfig) (*C.XVisualInfo, error) {
	visualInfo := C.glXGetVisualFromFBConfig(display.Ptr, config)
	if visualInfo == nil {
		return nil, errors.New("glXGetVisualFromFBConfig failed to retrieve visual information")
	}
	return visualInfo, nil
}

// CreateGLXContext creates a new GLX rendering context with specified attributes.
func CreateGLXContext(display *Display, config GLXFBConfig, shareContext GLXContext, direct bool, attribs GLXContextAttributes) (GLXContext, error) {
	cAttribs := C.GLXContextAttributes{
		contextMajorVersion: C.int(attribs.MajorVersion),
		contextMinorVersion: C.int(attribs.MinorVersion),
		profileMask:         C.int(attribs.ProfileMask),
	}

	directC := C.Bool(0)
	if direct {
		directC = C.True
	}

	context := C.createGLXContext(display.Ptr, config, shareContext, directC, cAttribs)
	if context == nil {
		return nil, errors.New("createGLXContext failed to create GLX context")
	}
	return context, nil
}

// MakeGLXContextCurrent makes the specified GLX context current on the drawable.
func MakeGLXContextCurrent(display *Display, window *Window, context GLXContext) error {
	if C.glXMakeCurrent(display.Ptr, window.ID, context) == 0 {
		return errors.New("glXMakeCurrent failed to make context current")
	}
	return nil
}

// SwapBuffers swaps the front and back buffers of the specified drawable.
func SwapBuffers(display *Display, window *Window) {
	C.glXSwapBuffers(display.Ptr, window.ID)
}

// DisableVSync disables VSync using available GLX extensions.
func DisableVSync(display *Display, window *Window) error {
	// Check for available extensions
	extensions := C.GoString(C.glXQueryExtensionsString(display.Ptr, C.XDefaultScreen(display.Ptr)))

	// glXSwapIntervalEXT (GLX_EXT_swap_control)
	if ext := "GLX_EXT_swap_control"; containsExtension(extensions, ext) {
		if C.glXSwapIntervalEXT == nil {
			ptr := C.glXGetProcAddressARB((*C.GLubyte)(unsafe.Pointer(C.CString("glXSwapIntervalEXT"))))
			if ptr == nil {
				return errors.New("failed to load glXSwapIntervalEXT")
			}
			C.glXSwapIntervalEXT = (*[0]byte)(ptr)
		}
		if C.glXSwapIntervalEXT != nil {
			swapIntervalEXT := *(*func(*C.Display, C.GLXDrawable, C.int))(unsafe.Pointer(&C.glXSwapIntervalEXT))
			swapIntervalEXT(display.Ptr, window.ID, 0) // Set swap interval to 0
			return nil
		}
	}

	// glXSwapIntervalMESA (GLX_MESA_swap_control)
	if ext := "GLX_MESA_swap_control"; containsExtension(extensions, ext) {
		if C.glXSwapIntervalMESA == nil {
			ptr := C.glXGetProcAddressARB((*C.GLubyte)(unsafe.Pointer(C.CString("glXSwapIntervalMESA"))))
			if ptr == nil {
				return errors.New("failed to load glXSwapIntervalMESA")
			}
			C.glXSwapIntervalMESA = (*[0]byte)(ptr)
		}
		if C.glXSwapIntervalMESA != nil {
			swapIntervalMESA := *(*func(C.uint) C.int)(unsafe.Pointer(&C.glXSwapIntervalMESA))
			ret := swapIntervalMESA(0) // Set swap interval to 0
			if ret == 0 {
				return nil
			}
			return errors.New("failed to set swap interval using GLX_MESA_swap_control")
		}
	}

	// glXSwapIntervalSGI (GLX_SGI_swap_control)
	if ext := "GLX_SGI_swap_control"; containsExtension(extensions, ext) {
		if C.glXSwapIntervalSGI == nil {
			ptr := C.glXGetProcAddressARB((*C.GLubyte)(unsafe.Pointer(C.CString("glXSwapIntervalSGI"))))
			if ptr == nil {
				return errors.New("failed to load glXSwapIntervalSGI")
			}
			C.glXSwapIntervalSGI = (*[0]byte)(ptr)
		}
		if C.glXSwapIntervalSGI != nil {
			swapIntervalSGI := *(*func(C.int) C.int)(unsafe.Pointer(&C.glXSwapIntervalSGI))
			ret := swapIntervalSGI(0) // Set swap interval to 0
			if ret == 0 {
				return nil
			}
			return errors.New("failed to set swap interval using GLX_SGI_swap_control")
		}
	}

	// If no supported extension is found
	return errors.New("no supported GLX swap control extension found to disable VSync")
}

// Helper function to check for a specific extension
func containsExtension(extensionString, extension string) bool {
	extList := strings.Split(extensionString, " ")
	for _, ext := range extList {
		if ext == extension {
			return true
		}
	}
	return false
}
