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
	fmt.Println("[host] - Linux - sparking X window management")

	// Fire up the X server connection
	X, err = xgbutil.NewConn()
	if err != nil {
		panic(err)
	}

	// Spark its "main thread"
	go xevent.Main(X)

	// Set up a thread to clean it up when JanOS "shuts down"
	go func() {
		core.WhileAlive()
		xevent.Quit(X)
	}()
}

// X represents a handle to the underlying x server connection.
var X *xgbutil.XUtil

// Create sparks a new x window and returns a handle to it
func Create() *xwindow.Window {
	atomic.AddInt32(&Count, 1)
	handle, err := xwindow.Generate(X)
	if err != nil {
		panic(err)
	}

	handle.Create(X.RootWin(), 0, 0, 200, 200, xproto.CwEventMask, xproto.EventMaskNoEvent)

	handle.WMGracefulClose(
		func(w *xwindow.Window) {
			xevent.Detach(w.X, w.Id)
			w.Destroy()
			atomic.AddInt32(&Count, -1)
		})

	handle.Map()
	return handle
}
