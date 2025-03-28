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
	"fmt"
)

func GetCoordinates() (xy core.Coordinate[int], err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to get mouse position: %v", r)
		}
	}()

	var cX, cY C.int
	C.GetMouseCoordinates(&cX, &cY)
	return core.Coordinate[int]{
		X: int(cX),
		Y: int(cY),
	}, nil
}
