package hydra

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/ignite-laboratories/core"
	"runtime"
	"sync"
	"time"
)
import "github.com/veandco/go-sdl2/sdl"

// Head represents a control surface for neural rendering to an SDL window.
type Head struct {
	*core.System

	WindowID uint32
	Window   *sdl.Window
	Driver   *core.Neuron

	ready bool
	ctx   core.Context
	mutex sync.Mutex
}

func (w *Head) start(manageable Manageable) {
	runtime.LockOSThread()

	sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)

	// Create OpenGL context
	glContext, err := w.Window.GLCreateContext()
	if err != nil {
		core.Fatalf(ModuleName, "failed to create OpenGL context: %v\n", err)
	}
	defer sdl.GLDeleteContext(glContext)

	if err := sdl.GLSetSwapInterval(-1); err != nil {
		core.Printf(ModuleName, "adaptive v-sync not available, falling back to v-sync\n")
		if err := sdl.GLSetSwapInterval(1); err != nil {
			core.Printf(ModuleName, "standard V-Sync also failed: %v\n", err)
		}
	}

	// Initialize OpenGL
	if err := gl.Init(); err != nil {
		core.Fatalf(ModuleName, "failed to initialize OpenGL: %v", err)
	}

	// Get OpenGL version
	//glVersion := gl.GoStr(gl.GetString(gl.VERSION))
	//fmt.Printf("[%v] [%d.%d] initialized with %s\n", ModuleName, w.WindowID, w.ID, glVersion)
	//
	//fmt.Println("openGL Version:", glVersion)
	//
	//// Get and print extensions
	//numExtensions := int32(0)
	//gl.GetIntegerv(gl.NUM_EXTENSIONS, &numExtensions)
	//
	//for i := int32(0); i < numExtensions; i++ {
	//	extension := gl.GoStr(gl.GetStringi(gl.EXTENSIONS, uint32(i)))
	//	if strings.Contains(extension, "geometry") {
	//		fmt.Println("found geometry-related extension:", extension)
	//	}
	//}

	manageable.Initialize()

	for core.Alive && w.Alive {
		// Busy wait for the next impulse signal
		for core.Alive && w.Alive && !w.ready {
			time.Sleep(time.Millisecond)
		}
		w.mutex.Lock()
		w.ready = false
		ctx := w.ctx
		w.mutex.Unlock()

		manageable.Render(ctx)

		w.Window.GLSwap()
	}

	manageable.Cleanup()
}

func (w *Head) impulse(ctx core.Context) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.ctx = ctx
	w.ready = true
}
