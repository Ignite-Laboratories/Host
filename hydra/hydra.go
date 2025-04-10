package hydra

import (
	"fmt"
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/core/std"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"runtime"
	"sync"
	"time"
)

func init() {
	var wg sync.WaitGroup
	wg.Add(1)
	// NOTE: The parameters here set the OpenGL specification
	go sparkSDL2(3, 1, false, &wg)

	wg.Wait()
}

type create struct {
	Title  string
	Size   std.XY[int]
	Pos    std.XY[int]
	Window *sdl.Window
}

var bridge = make(chan *std.Synchro[create])

var mutex sync.Mutex

// Windows provides the pointer handles to the underlying windows by their unique entity ID.
var Windows = make(map[uint64]*WindowCtrl)

// When provides a set of convenience potential functions.
var When when

type when int

// HasNoWindows provides a potential that returns true when all the windows have been globally closed.
func (w when) HasNoWindows(ctx core.Context) bool {
	return len(Windows) == 0
}

func sparkSDL2(major int, minor int, coreProfile bool, wg *sync.WaitGroup) {
	runtime.LockOSThread()

	// Initialize SDL
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Fatalf("Failed to initialize SDL: %v", err)
	}
	defer sdl.Quit()

	// Set OpenGL attributes
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, major)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, minor)
	if coreProfile {
		fmt.Println(fmt.Sprintf("[host] - sparking SDL2 integration with OpenGL %d.%d core profile", major, minor))
		sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	} else {
		fmt.Println(fmt.Sprintf("[host] - sparking SDL2 integration with OpenGL ES %d.%d", major, minor))
		sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_ES)
	}
	sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	sdl.GLSetAttribute(sdl.GL_DEPTH_SIZE, 24)

	wg.Done()

	for core.Alive {
		std.SynchroEngage(bridge, func(packet *create) {
			window, err := sdl.CreateWindow(
				packet.Title,
				int32(packet.Pos.X), int32(packet.Pos.Y),
				int32(packet.Size.X), int32(packet.Size.Y),
				sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE,
			)
			if err != nil {
				log.Fatalf("Failed to create SDL window: %v", err)
			}
			packet.Window = window
		})

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.WindowEvent:
				// Handle specific window close events
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					mutex.Lock()
					fmt.Printf("Window %d closing.\n", e.WindowID)
					for _, sys := range Windows {
						if sys.WindowID == e.WindowID {
							sys.Stop()
							err := sys.Window.Destroy()
							if err != nil {
								fmt.Printf("Failed to destroy window: %v\n", err)
							}
							delete(Windows, sys.ID)
						}
					}
					mutex.Unlock()
				}
			}
		}
	}
	time.Sleep(time.Millisecond * 250)
}

func CreateWindow(engine *core.Engine, title string, size std.XY[int], pos std.XY[int], action core.Action, potential core.Potential, muted bool) *WindowCtrl {
	fmt.Println("[host] - sparking new window")

	packet := std.SynchroSend(bridge, &create{
		Title: title,
		Size:  size,
		Pos:   pos,
	})

	window := packet.Window

	w := &WindowCtrl{}
	w.WindowID, _ = window.GetID()
	w.Window = window
	w.System = core.CreateSystem(engine, w.impulse, potential, muted)
	Windows[w.ID] = w
	go w.start(action)
	return w
}
