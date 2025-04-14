package hydra

import (
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
	"github.com/veandco/go-sdl2/sdl"
	"runtime"
	"sync"
	"time"
)

func init() {
	core.Printf(ModuleName, "sparking SDL integration\n")

	// Fire off SDL2 on its own thread, but wait for it to initialize
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
						core.Printf(ModuleName, "destroying head [%d.%d]\n", sys.WindowID, sys.ID)
						sys.Stop()
						err := sys.Window.Destroy()
						if err != nil {
							core.Printf(ModuleName, "failed to destroy window: %v\n", err)
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
		core.Fatalf(ModuleName, "failed to initialize SDL: %v\n", err)
	}
	defer sdl.Quit()

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
	// Let the rest of things 'initialize' before proceeding
	time.Sleep(time.Millisecond)

	driver, _ := sdl.GetCurrentVideoDriver()
	core.Verbosef(ModuleName, "SDL video driver: %s\n", driver)
	driver = sdl.GetCurrentAudioDriver()
	core.Verbosef(ModuleName, "SDL audio driver: %s\n", driver)

	for core.Alive {
		mainLoop()
	}

	core.Printf(ModuleName, "SDL integration stopped\n")
}

func CreateWindow(engine *core.Engine, title string, size *std.XY[int], pos *std.XY[int], manageable Manageable, potential core.Potential, muted bool) *Head {
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
			core.Fatalf(ModuleName, "failed to create SDL window: %v\n", err)
		}
		window = w
	})

	w := &Head{}
	w.WindowID, _ = window.GetID()
	w.Window = window
	w.System = core.CreateSystem(engine, w.impulse, potential, muted)
	w.synchro = make(std.Synchro)
	w.manageable = manageable
	Windows[w.ID] = w
	go w.run()
	core.Printf(ModuleName, "windowed head [%d.%d] created\n", w.WindowID, w.ID)
	return w
}

func CreateFullscreenWindow(engine *core.Engine, title string, manageable Manageable, potential core.Potential, muted bool) *Head {
	var window *sdl.Window
	synchro.Send(func() {
		w, err := sdl.CreateWindow(
			title,
			sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
			1024, 768,
			sdl.WINDOW_OPENGL|sdl.WINDOW_FULLSCREEN_DESKTOP,
		)
		if err != nil {
			core.Fatalf(ModuleName, "failed to create SDL window: %v\n", err)
		}
		window = w
	})

	w := &Head{}
	w.WindowID, _ = window.GetID()
	w.Window = window
	w.System = core.CreateSystem(engine, w.impulse, potential, muted)
	w.synchro = make(std.Synchro)
	w.manageable = manageable
	Windows[w.ID] = w
	go w.run()
	core.Printf(ModuleName, "fullscreen head [%d.%d] created\n", w.WindowID, w.ID)
	return w
}
