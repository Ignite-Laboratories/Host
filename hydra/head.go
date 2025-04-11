package hydra

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/ignite-laboratories/core"
	"log"
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

func (w *Head) start(initialize func(), action core.Action) {
	runtime.LockOSThread()

	// Create OpenGL context
	glContext, err := w.Window.GLCreateContext()
	if err != nil {
		log.Fatalf("[%v] failed to create OpenGL context: %v", ModuleName, err)
	}
	defer sdl.GLDeleteContext(glContext)

	// Enable VSync
	if err := sdl.GLSetSwapInterval(1); err != nil {
		log.Printf("[%v] failed to set VSync: %v", ModuleName, err)
	}

	// Initialize OpenGL
	if err := gl.Init(); err != nil {
		log.Fatalf("[%v] failed to initialize OpenGL: %v", ModuleName, err)
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

	initialize()

	for core.Alive && w.Alive {
		// Busy wait for the next impulse signal
		for core.Alive && w.Alive && !w.ready {
			time.Sleep(time.Millisecond)
		}
		w.mutex.Lock()
		w.ready = false
		ctx := w.ctx
		w.mutex.Unlock()

		action(ctx)

		w.Window.GLSwap()
	}
}

func (w *Head) impulse(ctx core.Context) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.ctx = ctx
	w.ready = true
}
