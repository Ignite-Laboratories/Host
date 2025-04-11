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
	fmt.Printf("[%v] sparking SDL integration\n", ModuleName)

	var wg sync.WaitGroup
	wg.Add(1)
	// NOTE: The parameters here set the OpenGL specification
	go sparkSDL2(3, 1, false, &wg)
	wg.Wait()
}

var synchro = make(std.Synchro)

var mutex sync.Mutex

// Windows provides the pointer handles to the underlying windows by their unique entity ID.
var Windows = make(map[uint64]*Head)

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
							fmt.Printf("[%v] failed to destroy window: %v\n", ModuleName, err)
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
					go core.ShutdownNow()
				}
			}
		}
	}
}

func sparkSDL2(major int, minor int, coreProfile bool, wg *sync.WaitGroup) {
	runtime.LockOSThread()

	// Initialize SDL
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalf("[%v] failed to initialize SDL: %v", ModuleName, err)
	}
	defer sdl.Quit()
	driver, _ := sdl.GetCurrentVideoDriver()
	fmt.Printf("[%v] SDL video driver: %s\n", ModuleName, driver)

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

	fmt.Printf("[%v] SDL integration stopped\n", ModuleName)
}

func CreateWindow(engine *core.Engine, title string, size *std.XY[int], pos *std.XY[int], initialize func(), action core.Action, potential core.Potential, muted bool) *Head {
	var window *sdl.Window
	synchro.Send(func() {
		var posX = sdl.WINDOWPOS_UNDEFINED
		var posY = sdl.WINDOWPOS_UNDEFINED
		if pos != nil {
			posX = pos.X
			posY = pos.Y
		}

		var sizeX = 640
		var sizeY = 480
		if size != nil {
			sizeX = size.X
			sizeY = size.Y
		}

		w, err := sdl.CreateWindow(
			title,
			int32(posX), int32(posY),
			int32(sizeX), int32(sizeY),
			sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE,
		)
		if err != nil {
			log.Fatalf("[%v] failed to create SDL window: %v", ModuleName, err)
		}
		window = w
	})

	w := &Head{}
	w.WindowID, _ = window.GetID()
	w.Window = window
	w.System = core.CreateSystem(engine, w.impulse, potential, muted)
	Windows[w.ID] = w
	go w.start(initialize, action)
	fmt.Printf("[%v] [%d.%d] window created\n", ModuleName, w.WindowID, w.ID)
	return w
}

func CreateFullscreenWindow(engine *core.Engine, title string, initialize func(), action core.Action, potential core.Potential, muted bool) *Head {
	var window *sdl.Window
	synchro.Send(func() {
		w, err := sdl.CreateWindow(
			title,
			sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
			1024, 768,
			sdl.WINDOW_OPENGL|sdl.WINDOW_FULLSCREEN_DESKTOP,
		)
		if err != nil {
			log.Fatalf("[%v] failed to create SDL window: %v", ModuleName, err)
		}
		window = w
	})

	w := &Head{}
	w.WindowID, _ = window.GetID()
	w.Window = window
	w.System = core.CreateSystem(engine, w.impulse, potential, muted)
	Windows[w.ID] = w
	go w.start(initialize, action)
	fmt.Printf("[%v] [%d.%d] window created\n", ModuleName, w.WindowID, w.ID)
	return w
}
