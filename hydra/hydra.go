package hydra

import (
	"fmt"
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"runtime"
	"sync"
)

func init() {
	fmt.Printf("[%v] - sparking SDL2 integration\n", ModuleName)

	var wg sync.WaitGroup
	wg.Add(1)
	// NOTE: The parameters here set the OpenGL specification
	go sparkSDL2(3, 1, false, &wg)
	wg.Wait()
}

var synchro = make(std.Synchro)

var mutex sync.Mutex

// Windows provides the pointer handles to the underlying windows by their unique entity ID.
var Windows = make(map[uint64]*WindowHead)

// HasNoWindows provides a potential that returns true when all the windows have been globally closed.
func HasNoWindows(ctx core.Context) bool {
	return len(Windows) == 0
}

func mainLoop() {
	synchro.Engage() // Listen for external execution

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.WindowEvent:
			// Handle specific window close events
			if e.Event == sdl.WINDOWEVENT_CLOSE {
				mutex.Lock()
				for _, sys := range Windows {
					if sys.WindowID == e.WindowID {
						fmt.Printf("[%v] [%d.%d] closing window\n", ModuleName, sys.WindowID, sys.ID)
						sys.Stop()
						err := sys.Window.Destroy()
						if err != nil {
							fmt.Printf("[%v] Failed to destroy window: %v\n", ModuleName, err)
						}
						delete(Windows, sys.ID)
					}
				}
				mutex.Unlock()
			}
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN {
				switch e.Keysym.Sym {
				case sdl.K_ESCAPE:
					core.ShutdownNow()
				}
			}
		}
	}
}

func sparkSDL2(major int, minor int, coreProfile bool, wg *sync.WaitGroup) {
	runtime.LockOSThread()

	// Initialize SDL
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalf("[%v] Failed to initialize SDL: %v", ModuleName, err)
	}
	defer sdl.Quit()
	driver, _ := sdl.GetCurrentVideoDriver()
	fmt.Printf("Current SDL Video Driver: %s\n", driver)

	// Set OpenGL attributes
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, major)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, minor)
	if coreProfile {
		sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	} else {
		sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_ES)
	}
	sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	sdl.GLSetAttribute(sdl.GL_DEPTH_SIZE, 24)

	wg.Done()

	for core.Alive {
		mainLoop()
	}
	sdl.Quit()
}

func CreateWindow(engine *core.Engine, title string, size std.XY[int], pos std.XY[int], action core.Action, potential core.Potential, muted bool) *WindowHead {
	var window *sdl.Window
	synchro.Send(func() {
		w, err := sdl.CreateWindow(
			title,
			int32(pos.X), int32(pos.Y),
			int32(size.X), int32(size.Y),
			sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE,
		)
		if err != nil {
			log.Fatalf("[%v] Failed to create SDL window: %v", ModuleName, err)
		}
		window = w
	})

	w := &WindowHead{}
	w.WindowID, _ = window.GetID()
	w.Window = window
	w.System = core.CreateSystem(engine, w.impulse, potential, muted)
	Windows[w.ID] = w
	go w.start(action)
	fmt.Printf("[%v] [%d.%d] window created\n", ModuleName, w.WindowID, w.ID)
	return w
}

func CreateFullscreenWindow(engine *core.Engine, title string, action core.Action, potential core.Potential, muted bool) *WindowHead {
	var window *sdl.Window
	synchro.Send(func() {
		w, err := sdl.CreateWindow(
			title,
			sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
			1024, 768,
			sdl.WINDOW_OPENGL|sdl.WINDOW_FULLSCREEN_DESKTOP,
		)
		if err != nil {
			log.Fatalf("[%v] Failed to create SDL window: %v", ModuleName, err)
		}
		window = w
	})

	w := &WindowHead{}
	w.WindowID, _ = window.GetID()
	w.Window = window
	w.System = core.CreateSystem(engine, w.impulse, potential, muted)
	Windows[w.ID] = w
	go w.start(action)
	fmt.Printf("[%v] [%d.%d] window created\n", ModuleName, w.WindowID, w.ID)
	return w
}
