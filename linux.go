//go:build linux

package host

import (
	"fmt"
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/host/x11"
	"log"
)

func init() {
	fmt.Println("[host] - Linux - sparking X11 integration")
	var err error
	X, err = x11.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize X11: %s\n", err)
	}

	go func() {
		core.WhileAlive()
		X.Close()
		fmt.Println("[host] - Linux - closed X11 integration")
	}()
}

// X provides a handle to the global connection to the X11 server.
var X *x11.Wrapper
