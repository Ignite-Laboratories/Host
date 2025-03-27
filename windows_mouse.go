//go:build windows

package host

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
import "fmt"

func (m mouse) GetCoordinates() (x int, y int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to get mouse position: %v", r)
		}
	}()

	var cX, cY C.int
	C.GetMouseCoordinates(&cX, &cY)
	return int(cX), int(cY), nil
}
