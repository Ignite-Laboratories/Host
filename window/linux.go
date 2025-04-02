//go:build linux

package window

import (
	"fmt"
	"github.com/ignite-laboratories/core"
	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgbutil"
	"github.com/jezek/xgbutil/xevent"
	"github.com/jezek/xgbutil/xwindow"
	"sync/atomic"
)

func init() {
	var err error
	fmt.Println("[host] - Linux - x window")

	// Fire up the X server connection
	X, err = xgbutil.NewConn()
	if err != nil {
		panic(err)
	}

	// Setup a thread for it to operate in
	go xevent.Main(X)

	// Setup a thread to terminate it when JanOS "shuts down"
	go func() {
		for core.Alive {
		}
		xevent.Quit(X)
	}()
}

// StopPotential provides a potential that returns true when all of the windows have been globally closed.
func StopPotential(ctx core.Context) bool {
	return Count == 0
}

// X represents the handle to the underlying x server connection.
var X *xgbutil.XUtil

// Count provides the number of open x windows.
var Count int32

// Create creates a new x window and returns a handle to it
func Create() *xwindow.Window {
	atomic.AddInt32(&Count, 1)
	win, err := xwindow.Generate(X)
	if err != nil {
		panic(err)
	}

	win.Create(X.RootWin(), 0, 0, 200, 200, xproto.CwEventMask, xproto.EventMaskNoEvent)

	win.WMGracefulClose(
		func(w *xwindow.Window) {
			xevent.Detach(w.X, w.Id)
			w.Destroy()
			atomic.AddInt32(&Count, -1)
		})

	win.Map()
	return win
}
