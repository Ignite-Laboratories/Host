//go:build linux

package mouse

/*
#cgo LDFLAGS: -lX11
#include <X11/Xlib.h>

static void GetMouseCoordinates(int *x, int *y) {
    Display *display = XOpenDisplay(NULL);
    if (display == NULL) return;

    Window root, child;
    int root_x, root_y;
    int win_x, win_y;
    unsigned int mask;

    XQueryPointer(display, DefaultRootWindow(display), &root, &child,
                  &root_x, &root_y, &win_x, &win_y, &mask);
    *x = root_x;
    *y = root_y;
    XCloseDisplay(display);
}
*/
import "C"
import (
	"fmt"
	"github.com/ignite-laboratories/core/std"
)

func getCoordinates() (xy std.XY[int], err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to get mouse position: %v", r)
		}
	}()

	var cX, cY C.int
	C.GetMouseCoordinates(&cX, &cY)
	return std.XY[int]{
		X: int(cX),
		Y: int(cY),
	}, nil
}
