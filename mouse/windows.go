//go:build windows

package mouse

/*
#include <windows.h>

void GetMouseCoordinates(int *x, int *y) {
    POINT p;
    GetCursorPos(&p);
    *x = p.x;
    *y = p.y;
}
*/
import "C"
import (
	"github.com/ignite-laboratories/core/std"
)

func getCoordinates() std.XY[int] {
	var cX, cY C.int
	C.GetMouseCoordinates(&cX, &cY)
	return std.XY[int]{
		X: int(cX),
		Y: int(cY),
	}
}
