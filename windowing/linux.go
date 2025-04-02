//go:build linux

package windowing

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
	fmt.Println("[host] - Linux - x windowing")

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

// StopPotential provides a common stopping potential for an impulse engine paired with a windowing context.
func StopPotential(ctx core.Context) bool {
	return Count == 0
}

var X *xgbutil.XUtil
var Count int32

func CreateWindow() *xwindow.Window {
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
